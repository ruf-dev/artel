// Package vaultpg is the client for per-vault Postgres databases: building connection strings for
// a domain.PostgresInstance/vault database+role pair, running the admin-side DDL that provisions a
// vault's role+database (AdminClient), and caching the tenant-side *sql.DB connections the MCP
// executor queries through (ConnPool).
package vaultpg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/utils"
	"go.redsock.ru/rerrors"
)

// Config identifies one Postgres database+role to connect to — either a postgres_instances admin
// database (AdminClient) or a per-vault tenant database (ConnPool).
type Config struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
	SSLMode  string
}

// DSN builds a `postgres://` URL-form DSN — lib/pq accepts both keyword and URL DSN forms, this
// codebase's other Postgres DSN builder (matreshka's resources.Postgres.ConnectionString) also
// uses the URL form. net/url takes care of escaping special characters in user/password/database.
func (c Config) DSN() string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.User, c.Password),
		Host:   fmt.Sprintf("%s:%d", c.Host, c.Port),
		Path:   "/" + c.Database,
	}

	query := u.Query()
	if c.SSLMode != "" {
		query.Set("sslmode", c.SSLMode)
	}
	u.RawQuery = query.Encode()

	return u.String()
}

// AdminClient runs the admin-side DDL (CREATE ROLE/DATABASE, etc.) against a postgres_instances
// row's maintenance database — mirrors couchdb.Client's role in vault.go's ensureCouchUserExists/
// ensureVaultExists, but for Postgres there's no higher-level client library, so this talks
// database/sql directly.
type AdminClient struct {
	db *sql.DB
}

// NewAdminClient opens a fresh connection to instance's admin_database — mirrors newCouchClient in
// internal/service/v1/vault/vault.go (fresh per call, not pooled).
func NewAdminClient(instance domain.PostgresInstance) (*AdminClient, error) {
	cfg := Config{
		Host:     instance.Host,
		Port:     instance.Port,
		Database: instance.AdminDatabase,
		User:     instance.Username,
		Password: instance.Password,
		SSLMode:  instance.SSLMode,
	}

	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, rerrors.Wrap(err, "error opening admin postgres connection")
	}

	c := &AdminClient{
		db: db,
	}

	return c, nil
}

// Close closes the underlying admin connection — callers own the AdminClient's lifetime since
// NewAdminClient opens a fresh connection per call.
func (c *AdminClient) Close() error {
	return c.db.Close()
}

// EnsureRole creates roleName as a LOGIN role with password, or — if it already exists —
// resets its password. Identifiers are quoted with pq.QuoteIdentifier and the password is quoted
// with pq.QuoteLiteral since CREATE ROLE/ALTER ROLE do not accept bind parameters for these
// clauses; callers only ever pass UUID-derived role names and generatePassword()-style random hex,
// never raw user input, but this stays defensive regardless.
func (c *AdminClient) EnsureRole(ctx context.Context, roleName, password string) error {
	quotedRole := pq.QuoteIdentifier(roleName)
	quotedPassword := pq.QuoteLiteral(password)

	createStmt := fmt.Sprintf("CREATE ROLE %s LOGIN PASSWORD %s", quotedRole, quotedPassword)

	_, err := c.db.ExecContext(ctx, createStmt)
	if err == nil {
		return nil
	}

	if !isDuplicateObject(err) {
		return rerrors.Wrap(err, "error creating postgres role")
	}

	alterStmt := fmt.Sprintf("ALTER ROLE %s PASSWORD %s", quotedRole, quotedPassword)

	_, err = c.db.ExecContext(ctx, alterStmt)
	if err != nil {
		return rerrors.Wrap(err, "error updating existing postgres role password")
	}

	return nil
}

// EnsureDatabase creates dbName owned by ownerRole, tolerating "database already exists" the same
// way ensureVaultExists tolerates user_errors.CouchDbDatabaseAlreadyExists.
func (c *AdminClient) EnsureDatabase(ctx context.Context, dbName, ownerRole string) error {
	quotedDB := pq.QuoteIdentifier(dbName)
	quotedOwner := pq.QuoteIdentifier(ownerRole)

	stmt := fmt.Sprintf("CREATE DATABASE %s OWNER %s", quotedDB, quotedOwner)

	_, err := c.db.ExecContext(ctx, stmt)
	if err != nil {
		if isDuplicateDatabase(err) {
			return nil
		}

		return rerrors.Wrap(err, "error creating postgres database")
	}

	return nil
}

// DropDatabase best-effort drops dbName, tolerating "does not exist" via IF EXISTS.
func (c *AdminClient) DropDatabase(ctx context.Context, dbName string) error {
	quotedDB := pq.QuoteIdentifier(dbName)

	stmt := fmt.Sprintf("DROP DATABASE IF EXISTS %s", quotedDB)

	_, err := c.db.ExecContext(ctx, stmt)
	if err != nil {
		return rerrors.Wrap(err, "error dropping postgres database")
	}

	return nil
}

// DropRole best-effort drops roleName, tolerating "does not exist" via IF EXISTS.
func (c *AdminClient) DropRole(ctx context.Context, roleName string) error {
	quotedRole := pq.QuoteIdentifier(roleName)

	stmt := fmt.Sprintf("DROP ROLE IF EXISTS %s", quotedRole)

	_, err := c.db.ExecContext(ctx, stmt)
	if err != nil {
		return rerrors.Wrap(err, "error dropping postgres role")
	}

	return nil
}

// isDuplicateObject reports whether err is Postgres error code 42710 (duplicate_object), raised by
// CREATE ROLE when the role already exists.
func isDuplicateObject(err error) bool {
	var pqErr *pq.Error

	if errors.As(err, &pqErr) {
		return pqErr.Code == pqerror.DuplicateObject
	}

	return strings.Contains(err.Error(), "already exists")
}

// isDuplicateDatabase reports whether err is Postgres error code 42P04 (duplicate_database),
// raised by CREATE DATABASE when the database already exists.
func isDuplicateDatabase(err error) bool {
	var pqErr *pq.Error

	if errors.As(err, &pqErr) {
		return pqErr.Code == pqerror.DuplicateDatabase
	}

	return strings.Contains(err.Error(), "already exists")
}

// ConnPool caches one tenant *sql.DB per vault — opened and validated on first use, reused
// thereafter, evicted on Close. Used by the MCP executor (a later task) to avoid reopening a
// connection on every tool call.
type ConnPool struct {
	mu    sync.Mutex
	conns map[uuid.UUID]*sql.DB
}

// NewConnPool returns an empty ConnPool ready for use.
func NewConnPool() *ConnPool {
	p := &ConnPool{
		conns: make(map[uuid.UUID]*sql.DB),
	}

	return p
}

// Get returns the cached connection for vaultID, opening and caching one against cfg if none
// exists yet, or if the cached one no longer pings successfully.
func (p *ConnPool) Get(ctx context.Context, vaultID uuid.UUID, cfg Config) (*sql.DB, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	db, ok := p.conns[vaultID]
	if ok {
		err := db.PingContext(ctx)
		if err == nil {
			return db, nil
		}

		utils.CloseWithLog(db, "stale vault postgres connection")
		delete(p.conns, vaultID)
	}

	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, rerrors.Wrap(err, "error opening vault postgres connection")
	}

	err = db.PingContext(ctx)
	if err != nil {
		utils.CloseWithLog(db, "failed vault postgres connection")

		return nil, rerrors.Wrap(err, "error pinging vault postgres connection")
	}

	p.conns[vaultID] = db

	return db, nil
}

// Close closes and evicts vaultID's cached connection, if any — called when Postgres is disabled
// for a vault. No-op if there's nothing cached.
func (p *ConnPool) Close(vaultID uuid.UUID) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	db, ok := p.conns[vaultID]
	if !ok {
		return nil
	}

	delete(p.conns, vaultID)

	err := db.Close()
	if err != nil {
		return rerrors.Wrap(err, "error closing vault postgres connection")
	}

	return nil
}
