//go:build e2e
// +build e2e

// Package harness holds the shared bufconn-server / Postgres / service-wiring boilerplate used by
// every e2e-tagged suite under tests/ (e2e, byok, quota, tract_e2e, gitlab_trigger_e2e). All suites
// share the single Postgres/CouchDB/S3 stack from tests/docker-compose.yaml — this package
// deliberately does not create a database per call or otherwise isolate one suite from another
// beyond the deterministic Slug-derived naming callers apply themselves.
//
// Provisioning (migrations, the CouchDB/S3 admin-pool instances, system_settings setup) happens
// exactly once per full test run, in tests/bootstrap, not per-suite — see ApplyMigrations,
// ResetPostgres/ResetCouchDB/ResetS3, and ProvisionCouchInstance/ProvisionS3Instance. Suites
// themselves only look up what bootstrap already provisioned, via GetCouchInstance/GetS3Instance.
package harness

import (
	"context"
	"database/sql"
	"net"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/clients/couchdb"
	s3client "github.com/ruf-dev/artel/internal/clients/s3"
	"github.com/ruf-dev/artel/internal/config"
	"github.com/ruf-dev/artel/internal/cryptoutil"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware"
	repopg "github.com/ruf-dev/artel/internal/repository/pg"
	svcv1 "github.com/ruf-dev/artel/internal/service/v1"
	"github.com/ruf-dev/artel/migrations"
)

// bufSize is the in-memory bufconn listener buffer size shared by every e2e suite's grpc server.
const bufSize = 1024 * 1024

// couchReadyPollMaxAttempts and couchReadyPollDelay bound waitForCouchDBReady's post-Setup poll —
// see that function's doc comment for why it's needed at all.
const (
	couchReadyPollMaxAttempts = 20
	couchReadyPollDelay       = 200 * time.Millisecond
)

var notAlnumRe = regexp.MustCompile(`[^a-z0-9]+`)

// envOrDefault returns the value of the environment variable key, or def when it is unset/empty.
func envOrDefault(key, def string) string {
	v := os.Getenv(key)
	if v != "" {
		return v
	}

	return def
}

// CouchURL returns the shared e2e CouchDB base URL (COUCH_URL, defaulting to the
// tests/docker-compose.yaml instance on port 15985).
func CouchURL(t *testing.T) string {
	t.Helper()

	return envOrDefault("COUCH_URL", "http://localhost:15985")
}

// CouchCreds returns the shared e2e CouchDB admin credentials (COUCH_USER/COUCH_PASS, defaulting
// to admin/admin — the tests/docker-compose.yaml default).
func CouchCreds(t *testing.T) (user, pass string) {
	t.Helper()

	return envOrDefault("COUCH_USER", "admin"), envOrDefault("COUCH_PASS", "admin")
}

// S3Endpoint returns the shared e2e S3 endpoint (S3_ENDPOINT, defaulting to the
// tests/docker-compose.yaml MinIO instance on port 19000).
func S3Endpoint(t *testing.T) string {
	t.Helper()

	return envOrDefault("S3_ENDPOINT", "localhost:19000")
}

// S3Creds returns the shared e2e S3 access/secret key (S3_ACCESS_KEY/S3_SECRET_KEY, defaulting to
// minioadmin/minioadmin — the literal every suite hardcoded before this helper existed).
func S3Creds(t *testing.T) (accessKey, secretKey string) {
	t.Helper()

	return envOrDefault("S3_ACCESS_KEY", "minioadmin"), envOrDefault("S3_SECRET_KEY", "minioadmin")
}

// OpenPostgres opens the shared e2e Postgres instance (DSN from PG_DSN, defaulting to the
// tests/docker-compose.yaml instance on port 15434) and pings it. Registers a t.Cleanup that
// closes the DB when the test finishes, so callers don't have to. Does NOT run migrations — call
// ApplyMigrations explicitly; tests/bootstrap does this once per full run, suites no longer do.
func OpenPostgres(t *testing.T) *sql.DB {
	t.Helper()

	pgDSN := envOrDefault("PG_DSN", "postgres://artel:artel_db@localhost:15434/artel_db?sslmode=disable")

	db, err := sql.Open("postgres", pgDSN)
	require.NoError(t, err, "open postgres")

	t.Cleanup(func() { _ = db.Close() })

	err = db.Ping()
	require.NoError(t, err, "ping postgres — is the container running?")

	return db
}

// ApplyMigrations runs the full migration set against db. Called exactly once per full test run,
// from tests/bootstrap's TestEnvSetup — suites themselves no longer call this.
func ApplyMigrations(t *testing.T, db *sql.DB) {
	t.Helper()

	err := migrations.ApplyMigration(db)
	require.NoError(t, err, "run migrations")
}

// ResetPostgres wipes all test-created state from db. This Postgres instance is dedicated to the
// e2e stack (tests/docker-compose.yaml) — nothing else uses it — so this is a full wipe, not
// prefix-matching of test-created rows.
//
// DELETE FROM users cascades through every table carrying an ON DELETE CASCADE FK back to users,
// directly or transitively via vaults/tracts/triggers/external_connections (vaults, mcp_keys,
// tracts, triggers, external_connections, workbenches, sessions, vault_members, subscriptions,
// couch_accounts, vault_invites, task_trackers, email_accounts, telegram_auth, tract_templates,
// tract_runs/tract_run_steps, trigger_links, mcp_spreadsheets, mcp_connectors, user_permissions —
// see migrations/*.sql for the individual FKs). Two things don't cascade off users and need an
// explicit sweep:
//   - couch_instances / s3_instances pool rows (owner_user_id IS NULL — never owned by a user)
//   - mcps rows seeded by e2e tests: mcps.owner_user_id is ON DELETE SET NULL, not CASCADE (see
//     migrations/060_mcp_ownership.sql), so community MoMs created by
//     tests/e2e/mcp_ownership_test.go (named "e2e_ownership_mom_<hex>" via its randomMcpName
//     helper) need an explicit sweep. The "e2e_" prefix deliberately does not match the seeded
//     system MoMs (email, gitlab, ...) from migrations.
func ResetPostgres(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()

	_, err := db.ExecContext(ctx, `DELETE FROM users`)
	require.NoError(t, err, "delete users")

	_, err = db.ExecContext(ctx, `DELETE FROM couch_instances WHERE owner_user_id IS NULL`)
	require.NoError(t, err, "delete pool couch instances")

	_, err = db.ExecContext(ctx, `DELETE FROM s3_instances WHERE owner_user_id IS NULL`)
	require.NoError(t, err, "delete pool s3 instances")

	_, err = db.ExecContext(ctx, `DELETE FROM mcps WHERE name LIKE 'e2e_%'`)
	require.NoError(t, err, "delete e2e-seeded mcps")
}

// ResetCouchDB deletes every database on the CouchDB instance at url except the three CouchDB
// system databases (_users, _replicator, _global_changes), leaving the instance itself — and its
// admin credentials / single-node setup — intact for reuse.
func ResetCouchDB(t *testing.T, ctx context.Context, url, username, password string) {
	t.Helper()

	cfg := couchdb.Config{BaseURL: url, User: username, Password: password}

	client, err := couchdb.New(cfg)
	require.NoError(t, err, "creating couch client to reset databases")

	dbs, err := client.ListDatabases(ctx)
	require.NoError(t, err, "listing couch databases")

	for _, dbName := range dbs {
		if isCouchSystemDB(dbName) {
			continue
		}

		err = client.DeleteDatabase(ctx, dbName)
		require.NoError(t, err, "deleting couch database "+dbName)
	}
}

func isCouchSystemDB(name string) bool {
	switch name {
	case "_users", "_replicator", "_global_changes":
		return true
	default:
		return false
	}
}

// ResetS3 deletes every bucket on the S3-compatible endpoint, emptying each bucket first — S3
// requires a bucket to be empty before it can be removed.
func ResetS3(t *testing.T, ctx context.Context, endpoint, accessKey, secretKey string) {
	t.Helper()

	cfg := s3client.Config{
		Endpoint:  endpoint,
		Region:    "us-east-1",
		AccessKey: accessKey,
		SecretKey: secretKey,
		UseSSL:    false,
		PathStyle: true,
	}

	buckets, err := s3client.ListBuckets(ctx, cfg)
	require.NoError(t, err, "listing s3 buckets")

	for _, bucket := range buckets {
		err = s3client.DeleteBucket(ctx, cfg, bucket)
		require.NoError(t, err, "deleting s3 bucket "+bucket)
	}
}

// ProvisionCouchInstance registers url as the e2e admin-pool CouchDB instance and runs its
// Setup() (system dbs + single-node enable) — see
// internal/service/v1/couchinstances/couchinstances.go's RegisterCouchInstance. This is the only
// place that registers the pool instance: it runs exactly once per test invocation, from
// tests/bootstrap, after ResetPostgres has already cleared any leftover pool row and before any
// suite runs — so unlike the old EnsureCouchInstance, there's no concurrent registration to
// tolerate and no AlreadyExists fallback needed.
func ProvisionCouchInstance(
	t *testing.T, ctx context.Context, svcs *svcv1.Services, url, username, password string,
) string {
	t.Helper()

	id, err := svcs.CouchInstance.RegisterCouchInstance(ctx, url, username, password)
	require.NoError(t, err, "register couch instance")

	waitForCouchDBReady(t, ctx, url, username, password)

	return id
}

// waitForCouchDBReady confirms the admin credentials are actually being honored before any suite
// attempts to create a real end-user CouchDB account. RegisterCouchInstance's internal Setup()
// call (enable_single_node) can report success while CouchDB is still propagating the admin
// credentials it just wrote — a write attempted against _users in that window fails with a
// transient "Unauthorized: Name or password is incorrect", even though enableSingleNode itself
// (internal/clients/couchdb/client.go) already retries its own call and returned success. This
// polls a lightweight authenticated request until it stops failing, the same way the old
// EnsureCouchInstance's waitForCouchSetup did for concurrent reusers — the reason has changed
// (settling after a fresh Setup, not waiting on another process) but the fix is the same shape.
func waitForCouchDBReady(t *testing.T, ctx context.Context, url, username, password string) {
	t.Helper()

	cfg := couchdb.Config{BaseURL: url, User: username, Password: password}

	client, err := couchdb.New(cfg)
	require.NoError(t, err, "creating couch client to confirm readiness")

	var lastErr error

	for attempt := 0; attempt < couchReadyPollMaxAttempts; attempt++ {
		_, lastErr = client.ListUsers(ctx)
		if lastErr == nil {
			return
		}

		if attempt < couchReadyPollMaxAttempts-1 {
			time.Sleep(couchReadyPollDelay)
		}
	}

	t.Fatalf("couch instance at %q did not accept authenticated requests in time: %v", url, lastErr)
}

// ProvisionS3Instance registers endpoint as the e2e admin-pool S3 instance. Same one-shot
// guarantee as ProvisionCouchInstance — called once from tests/bootstrap.
func ProvisionS3Instance(
	t *testing.T, ctx context.Context, svcs *svcv1.Services, endpoint, accessKey, secretKey string,
) string {
	t.Helper()

	id, err := svcs.S3Instance.RegisterS3Instance(ctx, endpoint, "us-east-1", accessKey, secretKey, false, true)
	require.NoError(t, err, "register s3 instance")

	return id
}

// GetCouchInstance looks up the id of the already-provisioned pool CouchDB instance at url.
// Suites call this instead of registering their own — provisioning now happens exactly once, in
// tests/bootstrap, before the suite run starts.
func GetCouchInstance(t *testing.T, ctx context.Context, svcs *svcv1.Services, url string) string {
	t.Helper()

	instances, err := svcs.CouchInstance.ListCouchInstances(ctx)
	require.NoError(t, err, "list couch instances")

	for _, instance := range instances {
		if instance.Url == url {
			return instance.Uuid.String()
		}
	}

	t.Fatalf(
		"no couch instance registered for url %q — run the bootstrap setup step first "+
			`(make test-e2e, or go test -tags "e2e e2e_bootstrap" ./tests/bootstrap/... -run TestEnvSetup)`,
		url,
	)

	return ""
}

// GetS3Instance looks up the id of the already-provisioned pool S3 instance at endpoint.
func GetS3Instance(t *testing.T, ctx context.Context, svcs *svcv1.Services, endpoint string) string {
	t.Helper()

	instances, err := svcs.S3Instance.ListS3Instances(ctx)
	require.NoError(t, err, "list s3 instances")

	for _, instance := range instances {
		if instance.Endpoint == endpoint {
			return instance.Uuid.String()
		}
	}

	t.Fatalf(
		"no s3 instance registered for endpoint %q — run the bootstrap setup step first "+
			`(make test-e2e, or go test -tags "e2e e2e_bootstrap" ./tests/bootstrap/... -run TestEnvSetup)`,
		endpoint,
	)

	return ""
}

// BuildServices wires a repopg.Repos + svcv1.Services pair on top of db, using the same
// zero-key AES encryptor every e2e suite has used so far. It returns whether that encryptor is
// plaintext — the credsEncrypted flag transport constructors like auth_api.NewAuthImpl take.
func BuildServices(t *testing.T, db *sql.DB, cfg config.EnvironmentConfig) (*repopg.Repos, *svcv1.Services, bool) {
	t.Helper()

	encKey := make([]byte, 32)

	encryptor, err := cryptoutil.NewAESEncryptor(encKey)
	require.NoError(t, err, "create AES encryptor")

	repos := repopg.New(db, encryptor)

	svcs, err := svcv1.New(repos, cfg)
	require.NoError(t, err, "init services")

	return repos, svcs, encryptor.IsPlainText()
}

// CompleteSetup marks first-run setup complete and enables self-registration, the idempotent pair
// every suite needs before it can Register/Login through the real AuthAPI without being gated by
// the setup wizard.
func CompleteSetup(t *testing.T, ctx context.Context, repos *repopg.Repos) {
	t.Helper()

	err := repos.SystemSettings().CompleteSetup(ctx)
	require.NoError(t, err, "complete setup so Register/Login aren't gated in tests")

	err = repos.SystemSettings().UpdateRegistrationMode(ctx, domain.RegistrationModeSelfRegister)
	require.NoError(t, err, "enable self-registration so tests can Register directly")
}

// Slug sanitizes t.Name() into a lowercase alnum+underscore string, for deterministic per-test
// identifiers (usernames, vault names, ...) that stay stable across repeated runs against the one
// shared DB.
func Slug(t *testing.T) string {
	t.Helper()

	return notAlnumRe.ReplaceAllString(strings.ToLower(t.Name()), "_")
}

// AuthedContext carries token as outgoing gRPC metadata, the same way a real client authenticates
// against middleware.GrpcAuthInterceptor.
func AuthedContext(ctx context.Context, token string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", token)
}

// NewBufconnServer builds a real *grpc.Server chained with the production auth interceptor
// (middleware.GrpcAuthInterceptor, with auth ignored on Register/Login), calls each of register
// against it, and serves it over an in-memory bufconn listener so callers' RPCs travel through the
// real transport + auth stack without binding a TCP port. Registers a t.Cleanup that closes the
// dialed connection and stops the server, so callers don't manage that lifecycle by hand. Callers
// build their own typed clients off the returned conn (pb.NewAuthAPIClient(conn), ...) since which
// APIs a suite needs varies per suite.
//
// register is typed grpc.ServiceRegistrar rather than *grpc.Server because that's the exact
// signature every transport impl's Register method already has (e.g. auth_api.AuthImpl.Register),
// so callers can pass those method values directly — *grpc.Server satisfies the interface.
func NewBufconnServer(t *testing.T, svcs *svcv1.Services, register ...func(grpc.ServiceRegistrar)) *grpc.ClientConn {
	t.Helper()

	ignoredAuthPaths := middleware.WithIgnoredPathAuthOption(
		pb.AuthAPI_Register_FullMethodName,
		pb.AuthAPI_Login_FullMethodName,
	)
	authOption := middleware.GrpcAuthInterceptor(svcs, ignoredAuthPaths)

	grpcServer := grpc.NewServer(authOption)
	for _, r := range register {
		r(grpcServer)
	}

	listener := bufconn.Listen(bufSize)

	go func() {
		_ = grpcServer.Serve(listener)
	}()

	dialer := func(ctx context.Context, _ string) (net.Conn, error) {
		return listener.DialContext(ctx)
	}

	conn, err := grpc.NewClient(
		"passthrough:///bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err, "dial bufconn grpc client")

	t.Cleanup(func() {
		_ = conn.Close()
		grpcServer.Stop()
	})

	return conn
}
