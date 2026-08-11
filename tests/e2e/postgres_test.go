//go:build e2e
// +build e2e

package e2e_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/clients/vaultpg"
	"github.com/ruf-dev/artel/internal/config"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	repopg "github.com/ruf-dev/artel/internal/repository/pg"
	svcv1 "github.com/ruf-dev/artel/internal/service/v1"
	"github.com/ruf-dev/artel/internal/transport/auth_api"
	"github.com/ruf-dev/artel/internal/transport/external_connections_api"
	"github.com/ruf-dev/artel/internal/transport/mcp_api"
	"github.com/ruf-dev/artel/internal/transport/mcp_keys_api"
	"github.com/ruf-dev/artel/internal/transport/vaults_api"
	"github.com/ruf-dev/artel/tests/harness"
)

// PostgresSuite exercises the per-vault Postgres feature end to end: the FeaturePostgres
// subscription gate, on-demand database+role provisioning against the shared admin pool,
// BYOK-owned instance resolution (AddPostgresConnection -> syncOwnedPostgresInstance), the
// pg_list_tables/pg_describe_table/pg_query/pg_execute MCP tools, and the non-superuser
// per-vault role's isolation from other vaults' databases.
//
// Like QuotaSuite, this runs with config.EnvironmentConfig.SubscriptionsEnabled = true so
// EnablePostgresDatabase's subscriptions.CheckFeature call actually routes through PaidService
// instead of FreeService's always-allow — the "basic" seed plan (migrations/047) does not set
// "postgres" in its features JSON, so it defaults to false and the gate has something real to
// reject.
type PostgresSuite struct {
	suite.Suite
	db                 *sql.DB
	repos              *repopg.Repos
	svcs               *svcv1.Services
	mcpHdlr            *mcp_api.McpHandler
	postgresInstanceID string

	authClient   pb.AuthAPIClient
	vaultsClient pb.VaultsAPIClient
	ecClient     pb.ExternalConnectionsAPIClient
	mcpKeyClient pb.McpKeysAPIClient
}

func TestPostgres(t *testing.T) {
	suite.Run(t, new(PostgresSuite))
}

func (s *PostgresSuite) SetupSuite() {
	ctx := context.Background()

	s.db = harness.OpenPostgres(s.T())

	cfg := config.EnvironmentConfig{}
	cfg.SubscriptionsEnabled = true

	var credsEncrypted bool
	s.repos, s.svcs, credsEncrypted = harness.BuildServices(s.T(), s.db, cfg)

	pgHost := harness.PostgresHost(s.T())
	pgPort := harness.PostgresPort(s.T())

	s.postgresInstanceID = harness.GetPostgresInstance(s.T(), ctx, s.svcs, pgHost, pgPort)

	s.mcpHdlr = mcp_api.NewMcpHandler(s.svcs.Mcp, s.svcs.Mom)

	s.startGrpcServer(credsEncrypted)
}

// startGrpcServer builds a real *grpc.Server chained with the production auth interceptor,
// registers the AuthAPI, VaultsAPI, ExternalConnectionsAPI and McpKeysAPI implementations onto
// it, and serves it over an in-memory bufconn listener — mirrors ByokStorageSuite/QuotaSuite.
func (s *PostgresSuite) startGrpcServer(credsEncrypted bool) {
	authImpl := auth_api.NewAuthImpl(
		s.svcs.Auth, "", s.svcs.S3Instance, s.svcs.CouchInstance,
		false, credsEncrypted, s.svcs.DockerHost, s.svcs.SetupWizard,
		false,
	)
	vaultsImpl := vaults_api.NewVaultsImpl(s.svcs.Vault, s.svcs.Workbench, false)
	externalConnectionsImpl := external_connections_api.New(s.svcs.ExternalConnectionService(), false)
	mcpKeysImpl := mcp_keys_api.NewMcpKeysImpl(s.svcs.Mcp, s.svcs.Mom, false)

	conn := harness.NewBufconnServer(
		s.T(), s.svcs, authImpl.Register, vaultsImpl.Register, externalConnectionsImpl.Register, mcpKeysImpl.Register,
	)

	s.authClient = pb.NewAuthAPIClient(conn)
	s.vaultsClient = pb.NewVaultsAPIClient(conn)
	s.ecClient = pb.NewExternalConnectionsAPIClient(conn)
	s.mcpKeyClient = pb.NewMcpKeysAPIClient(conn)
}

func (s *PostgresSuite) TearDownSuite() {
	if s.db != nil {
		err := s.db.Close()
		s.NoError(err, "close db")
	}
}

// pgTestUser bundles what every test in this suite needs after registering+authenticating a
// fresh user and creating them a vault.
type pgTestUser struct {
	userUuid  uuid.UUID
	authedCtx context.Context
	vaultUuid uuid.UUID
}

// registerLoginVault registers+logs in a fresh user via the real AuthAPI RPC, activates their
// subscription (required before any authenticated RPC since this suite runs with
// SubscriptionsEnabled=true — see QuotaSuite.setupUser's identical note), and creates them a
// vault via the real VaultsAPI RPC. suffix disambiguates the email/vault name when a single test
// needs more than one user (harness.Slug is derived from t.Name(), which is identical for every
// call within one test method).
func (s *PostgresSuite) registerLoginVault(suffix string) pgTestUser {
	ctx := context.Background()

	slug := harness.Slug(s.T()) + "_" + suffix
	email := slug + "@test.local"
	// A fixed password, not slug-derived: bcrypt rejects passwords over 72 bytes, and slug embeds
	// the full (possibly long, subtest-qualified) t.Name() — uniqueness only matters for email.
	password := "test-password-postgres"

	registerReq := &pb.Register_Request{Email: email, Password: password}

	registerResp, err := s.authClient.Register(ctx, registerReq)
	s.Require().NoError(err)

	userUuid, err := uuid.Parse(registerResp.Id)
	s.Require().NoError(err)
	s.T().Cleanup(func() {
		_ = s.repos.Users().Delete(context.Background(), userUuid)
	})

	// Email/password registration never populates domain.User.Username (defaults to "" — see
	// migrations/007_telegram_auth.sql). Vault creation derives the CouchDB database name from
	// it, and an empty UserName produces a name starting with "-", which CouchDB rejects.
	_, err = s.db.Exec(`UPDATE users SET username = $1 WHERE id = $2`, slug, userUuid)
	s.Require().NoError(err, "set username stand-in for vault creation")

	// With SubscriptionsEnabled=true, middleware.authWithSession calls subscriptionService.
	// CheckActive on every authenticated RPC, so the subscription has to be activated before
	// CreateVault below.
	_, err = s.repos.Subscriptions().Upsert(ctx, userUuid, true)
	s.Require().NoError(err, "activate subscription so authenticated RPCs pass CheckActive")

	passwordCreds := &pb.PasswordCredentials{Email: email, Password: password}
	loginReq := &pb.Login_Request{Method: &pb.Login_Request_Password{Password: passwordCreds}}

	loginResp, err := s.authClient.Login(ctx, loginReq)
	s.Require().NoError(err)
	s.Require().NotEmpty(loginResp.Token)

	authedCtx := harness.AuthedContext(ctx, loginResp.Token)

	createVaultResp, err := s.vaultsClient.CreateVault(authedCtx, &pb.CreateVault_Request{Name: slug + "_vault"})
	s.Require().NoError(err)

	vaultUuid, err := uuid.Parse(createVaultResp.Id)
	s.Require().NoError(err)
	s.T().Cleanup(func() {
		_ = s.svcs.Vault.DeleteVault(context.Background(), vaultUuid)
	})

	return pgTestUser{userUuid: userUuid, authedCtx: authedCtx, vaultUuid: vaultUuid}
}

// setPostgresFeature writes a sparse FeatureOverrides entry for domain.FeaturePostgres on
// userUuid's subscription, the same admin-override mechanism QuotaSuite.overrideQuota uses for
// storage quotas — the seeded "basic" plan (migrations/047_subscription_plans.sql) has no
// "postgres" key in its features JSON, so it defaults to false until overridden here.
func (s *PostgresSuite) setPostgresFeature(userUuid uuid.UUID, enabled bool) {
	ctx := context.Background()
	subsRepo := s.repos.Subscriptions()

	sub, err := subsRepo.GetByUser(ctx, userUuid)
	s.Require().NoError(err, "get subscription")

	if sub.FeatureOverrides == nil {
		sub.FeatureOverrides = map[domain.SubscriptionFeature]bool{}
	}

	sub.FeatureOverrides[domain.FeaturePostgres] = enabled

	_, err = subsRepo.UpsertOverrides(ctx, sub)
	s.Require().NoError(err, "override postgres feature")
}

// mcpPgTool drives a real tools/call request against the production MCP JSON-RPC handler
// (mirrors QuotaSuite.callWriteFile/e2e_test.go's inline pattern), asserting no rpc-level error
// and returning the single text content block's payload.
func (s *PostgresSuite) mcpPgTool(rawToken, name string, arguments map[string]any) string {
	callParams := map[string]any{
		"name":      name,
		"arguments": arguments,
	}
	body := mcpCall("tools/call", callParams)

	req := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+rawToken)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.mcpHdlr.ServeHTTP(w, req)
	s.Require().Equal(http.StatusOK, w.Code)

	var rpcResp struct {
		Error  *mcpRpcError `json:"error"`
		Result *struct {
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
		} `json:"result"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &rpcResp)
	s.Require().NoError(err)
	s.Require().Nil(rpcResp.Error, "mcp tool %s returned rpc error: %+v", name, rpcResp.Error)
	s.Require().NotNil(rpcResp.Result)
	s.Require().Len(rpcResp.Result.Content, 1)

	return rpcResp.Result.Content[0].Text
}

// TestEnablePostgresDatabase_FeatureGate_RejectsThenSucceeds confirms EnablePostgresDatabase is
// rejected while the caller's plan has FeaturePostgres off (the "basic" plan's default), and
// succeeds once an admin override turns it on — exercises subscriptions.CheckFeature routing
// through the real PaidService rather than the SubscriptionsEnabled=false default that always
// routes through FreeService's always-allow.
func (s *PostgresSuite) TestEnablePostgresDatabase_FeatureGate_RejectsThenSucceeds() {
	user := s.registerLoginVault("gate")

	_, err := s.vaultsClient.EnablePostgresDatabase(user.authedCtx, &pb.EnablePostgresDatabase_Request{
		VaultId: user.vaultUuid.String(),
	})
	s.Require().Error(err, "basic plan has the postgres feature disabled by default")

	st, ok := status.FromError(err)
	s.Require().True(ok, "error should be a well-formed grpc status")
	s.Equal(codes.PermissionDenied, st.Code())
	s.Contains(st.Message(), "not enabled", "error should explain the feature is not enabled")

	s.setPostgresFeature(user.userUuid, true)

	resp, err := s.vaultsClient.EnablePostgresDatabase(user.authedCtx, &pb.EnablePostgresDatabase_Request{
		VaultId: user.vaultUuid.String(),
	})
	s.Require().NoError(err, "enabling postgres should succeed once the feature override is granted")
	s.Equal("ready", resp.Status)
	s.Empty(resp.ErrorMessage)
}

// TestEnablePostgresDatabase_StatusTransitions_ResolvesToAdminPool confirms
// GetPostgresDatabase reports "not_enabled" before EnablePostgresDatabase runs, "ready"
// immediately after (EnablePostgresDatabase provisions synchronously — see
// vault.Service.EnablePostgresDatabase, no polling needed), that the vault resolves to the
// shared admin pool instance (no BYOK connection registered for this user), and that the
// database genuinely exists on the physical Postgres server, not just as a postgres row.
func (s *PostgresSuite) TestEnablePostgresDatabase_StatusTransitions_ResolvesToAdminPool() {
	user := s.registerLoginVault("lifecycle")
	s.setPostgresFeature(user.userUuid, true)

	getResp, err := s.vaultsClient.GetPostgresDatabase(user.authedCtx, &pb.GetPostgresDatabase_Request{
		VaultId: user.vaultUuid.String(),
	})
	s.Require().NoError(err)
	s.Equal("not_enabled", getResp.Status, "no postgres database should exist yet")

	enableResp, err := s.vaultsClient.EnablePostgresDatabase(user.authedCtx, &pb.EnablePostgresDatabase_Request{
		VaultId: user.vaultUuid.String(),
	})
	s.Require().NoError(err)
	s.Equal("ready", enableResp.Status, "EnablePostgresDatabase provisions synchronously, no provisioning-state poll needed")
	s.Empty(enableResp.ErrorMessage)

	getResp, err = s.vaultsClient.GetPostgresDatabase(user.authedCtx, &pb.GetPostgresDatabase_Request{
		VaultId: user.vaultUuid.String(),
	})
	s.Require().NoError(err)
	s.Equal("ready", getResp.Status)

	var instanceID uuid.UUID

	err = s.db.QueryRow(
		`SELECT postgres_instance_id FROM vault_postgres_databases WHERE vault_id = $1`, user.vaultUuid,
	).Scan(&instanceID)
	s.Require().NoError(err)
	s.Equal(s.postgresInstanceID, instanceID.String(), "vault with no BYOK connection should resolve to the shared admin pool instance")

	var dbName string

	err = s.db.QueryRow(
		`SELECT database_name FROM vault_postgres_databases WHERE vault_id = $1`, user.vaultUuid,
	).Scan(&dbName)
	s.Require().NoError(err)

	pgHost := harness.PostgresHost(s.T())
	pgPort := harness.PostgresPort(s.T())
	adminUser, adminPass := harness.PostgresAdminCreds(s.T())

	cfg := vaultpg.Config{Host: pgHost, Port: pgPort, Database: dbName, User: adminUser, Password: adminPass, SSLMode: "disable"}

	adminDB, err := sql.Open("postgres", cfg.DSN())
	s.Require().NoError(err)
	defer func() { _ = adminDB.Close() }()

	err = adminDB.PingContext(context.Background())
	s.Require().NoError(err, "provisioned database should genuinely exist and be reachable, not just recorded in postgres")
}

// TestAddPostgresConnection_SyncsOwnedInstance_VaultResolvesToOwned confirms
// AddPostgresConnection persists a BYOK external_connections row and syncs an owned
// postgres_instances pool row (syncOwnedPostgresInstance), and that a subsequent
// EnablePostgresDatabase call by that same user resolves to their owned instance rather than the
// shared admin pool instance registered by tests/bootstrap.
func (s *PostgresSuite) TestAddPostgresConnection_SyncsOwnedInstance_VaultResolvesToOwned() {
	user := s.registerLoginVault("byok")
	s.setPostgresFeature(user.userUuid, true)

	pgHost := harness.PostgresHost(s.T())
	pgPort := harness.PostgresPort(s.T())
	adminDatabase := harness.PostgresAdminDatabase(s.T())
	adminUser, adminPass := harness.PostgresAdminCreds(s.T())

	// postgres_instances carries no UNIQUE constraint on host/port (unlike couch_instances.url),
	// so this can reuse the exact same physical test-postgres server the admin pool instance
	// already targets — the row stays distinct because it's owner-scoped, not because the
	// connection coordinates differ.
	addReq := &pb.AddPostgresConnection_Request{
		Host:     pgHost,
		Port:     int32(pgPort),
		Database: adminDatabase,
		Username: adminUser,
		Password: adminPass,
		SslMode:  "disable",
	}

	_, err := s.ecClient.AddPostgresConnection(user.authedCtx, addReq)
	s.Require().NoError(err, "AddPostgresConnection rpc should succeed against the shared test postgres server")

	var provider string

	err = s.db.QueryRow(`SELECT provider FROM external_connections WHERE user_id = $1`, user.userUuid).Scan(&provider)
	s.Require().NoError(err, "external_connections row should have been persisted")
	s.Equal(domain.ProviderPostgres, provider)

	var ownedInstanceID uuid.UUID

	err = s.db.QueryRow(`SELECT id FROM postgres_instances WHERE owner_user_id = $1`, user.userUuid).Scan(&ownedInstanceID)
	s.Require().NoError(err, "owned postgres_instances row should have been created")
	s.NotEqual(s.postgresInstanceID, ownedInstanceID.String(), "owned instance must be a distinct row from the shared admin pool instance")

	enableResp, err := s.vaultsClient.EnablePostgresDatabase(user.authedCtx, &pb.EnablePostgresDatabase_Request{
		VaultId: user.vaultUuid.String(),
	})
	s.Require().NoError(err, "EnablePostgresDatabase rpc should succeed")
	s.Equal("ready", enableResp.Status)

	var vaultInstanceID uuid.UUID

	err = s.db.QueryRow(
		`SELECT postgres_instance_id FROM vault_postgres_databases WHERE vault_id = $1`, user.vaultUuid,
	).Scan(&vaultInstanceID)
	s.Require().NoError(err)
	s.Equal(ownedInstanceID, vaultInstanceID, "vault should resolve to the caller's owned BYOK instance")
	s.NotEqual(s.postgresInstanceID, vaultInstanceID.String(), "vault must not fall back to the shared admin pool instance")
}

// TestMcpPostgresRoundTrip_QueryExecuteListTables drives pg_execute/pg_query/pg_list_tables
// through the real MCP JSON-RPC-over-HTTP handler (mirrors e2e_test.go's write_file round trip)
// against a vault with Postgres enabled, confirming the tools actually reach the vault's
// dedicated tenant database and round-trip real rows.
func (s *PostgresSuite) TestMcpPostgresRoundTrip_QueryExecuteListTables() {
	user := s.registerLoginVault("mcp")
	s.setPostgresFeature(user.userUuid, true)

	_, err := s.vaultsClient.EnablePostgresDatabase(user.authedCtx, &pb.EnablePostgresDatabase_Request{
		VaultId: user.vaultUuid.String(),
	})
	s.Require().NoError(err)

	createKeyResp, err := s.mcpKeyClient.CreateMcpKey(user.authedCtx, &pb.CreateMcpKey_Request{
		VaultId: user.vaultUuid.String(),
		Name:    "postgres-mcp-key",
	})
	s.Require().NoError(err)
	s.T().Cleanup(func() {
		keyUuid, parseErr := uuid.Parse(createKeyResp.Key.Id)
		if parseErr == nil {
			_ = s.svcs.Mcp.RevokeKey(context.Background(), keyUuid)
		}
	})

	rawToken := createKeyResp.RawToken

	s.mcpPgTool(rawToken, "pg_execute", map[string]any{
		"sql": "CREATE TABLE items (id serial primary key, name text)",
	})
	s.mcpPgTool(rawToken, "pg_execute", map[string]any{
		"sql": "INSERT INTO items (name) VALUES ('test')",
	})

	queryResultText := s.mcpPgTool(rawToken, "pg_query", map[string]any{"sql": "SELECT * FROM items"})

	var queryResult struct {
		Rows      []map[string]interface{} `json:"rows"`
		Truncated bool                     `json:"truncated"`
	}

	err = json.Unmarshal([]byte(queryResultText), &queryResult)
	s.Require().NoError(err)
	s.Require().Len(queryResult.Rows, 1)
	s.Equal("test", queryResult.Rows[0]["name"])
	s.False(queryResult.Truncated)

	listResultText := s.mcpPgTool(rawToken, "pg_list_tables", map[string]any{})

	var tables []string

	err = json.Unmarshal([]byte(listResultText), &tables)
	s.Require().NoError(err)
	s.Contains(tables, "items")
}

// TestMcpPostgresTool_NoDatabaseLinked_ReturnsUserError confirms a pg_* tool call against a vault
// with no Postgres database enabled fails cleanly with user_errors.NoPostgresDatabaseLinked (424
// Failed Dependency), not a crash or an opaque generic failure.
func (s *PostgresSuite) TestMcpPostgresTool_NoDatabaseLinked_ReturnsUserError() {
	user := s.registerLoginVault("nopg")
	// Deliberately do not call EnablePostgresDatabase for this vault.

	createKeyResp, err := s.mcpKeyClient.CreateMcpKey(user.authedCtx, &pb.CreateMcpKey_Request{
		VaultId: user.vaultUuid.String(),
		Name:    "nopg-mcp-key",
	})
	s.Require().NoError(err)
	s.T().Cleanup(func() {
		keyUuid, parseErr := uuid.Parse(createKeyResp.Key.Id)
		if parseErr == nil {
			_ = s.svcs.Mcp.RevokeKey(context.Background(), keyUuid)
		}
	})

	callParams := map[string]any{"name": "pg_list_tables", "arguments": map[string]any{}}
	body := mcpCall("tools/call", callParams)

	req := httptest.NewRequest(http.MethodPost, "/mcp", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+createKeyResp.RawToken)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.mcpHdlr.ServeHTTP(w, req)
	s.Require().Equal(http.StatusOK, w.Code, "mcp jsonrpc errors are still 200 at the http transport level")

	var rpcResp struct {
		Error *mcpRpcError `json:"error"`
	}

	err = json.Unmarshal(w.Body.Bytes(), &rpcResp)
	s.Require().NoError(err)
	s.Require().NotNil(rpcResp.Error, "pg_list_tables against a vault with no postgres database should return an rpc error")
	s.Contains(
		rpcResp.Error.Details, "vault has no enabled postgres database",
		"error should be user_errors.NoPostgresDatabaseLinked, not a generic/crash failure",
	)
}

// TestPostgresRoleIsolation_CannotCreateInOtherVaultDatabase provisions Postgres for two
// separate vaults (owned by two separate users) on the shared admin pool, then opens a raw
// connection to vault B's database using vault A's granted, non-superuser role credentials and
// attempts to CREATE TABLE in it. Postgres 15+ revokes the CREATE privilege on a fresh database's
// public schema from PUBLIC, leaving only the database's owner (vault B's own role) able to
// create objects there — so this is a deterministic proof that the per-vault role genuinely
// cannot cross vault boundaries, independent of any SELECT-grant subtleties a read-only probe
// would be exposed to.
func (s *PostgresSuite) TestPostgresRoleIsolation_CannotCreateInOtherVaultDatabase() {
	userA := s.registerLoginVault("isoa")
	userB := s.registerLoginVault("isob")
	s.setPostgresFeature(userA.userUuid, true)
	s.setPostgresFeature(userB.userUuid, true)

	_, err := s.vaultsClient.EnablePostgresDatabase(userA.authedCtx, &pb.EnablePostgresDatabase_Request{
		VaultId: userA.vaultUuid.String(),
	})
	s.Require().NoError(err)

	_, err = s.vaultsClient.EnablePostgresDatabase(userB.authedCtx, &pb.EnablePostgresDatabase_Request{
		VaultId: userB.vaultUuid.String(),
	})
	s.Require().NoError(err)

	// GetPostgresDatabase (service layer, not the gRPC client) returns the decrypted role
	// credentials directly — user_context.WithUserContext stands in for the gRPC auth
	// interceptor the same way QuotaSuite.setupUser's userCtx does for LinkS3Bucket.
	ucA := user_context.UserContext{UserUuid: userA.userUuid}
	userCtxA := user_context.WithUserContext(context.Background(), ucA)

	roleA, err := s.svcs.Vault.GetPostgresDatabase(userCtxA, userA.vaultUuid)
	s.Require().NoError(err)
	s.Require().True(roleA.Valid)

	var dbNameB string

	err = s.db.QueryRow(
		`SELECT database_name FROM vault_postgres_databases WHERE vault_id = $1`, userB.vaultUuid,
	).Scan(&dbNameB)
	s.Require().NoError(err)

	pgHost := harness.PostgresHost(s.T())
	pgPort := harness.PostgresPort(s.T())

	cfg := vaultpg.Config{
		Host:     pgHost,
		Port:     pgPort,
		Database: dbNameB,
		User:     roleA.V.RoleUsername,
		Password: roleA.V.RolePassword,
		SSLMode:  "disable",
	}

	crossDB, err := sql.Open("postgres", cfg.DSN())
	s.Require().NoError(err)
	defer func() { _ = crossDB.Close() }()

	_, err = crossDB.ExecContext(context.Background(), "CREATE TABLE isolation_probe (id int)")
	s.Require().Error(err, "vault A's role must not be able to create objects in vault B's database")
	s.Contains(err.Error(), "permission denied", "failure should be a postgres permission error, proving role isolation")
}
