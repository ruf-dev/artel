//go:build e2e
// +build e2e

package e2e_test

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand/v2"
	"net"
	"strings"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	s3client "github.com/ruf-dev/artel/internal/clients/s3"
	"github.com/ruf-dev/artel/internal/config"
	"github.com/ruf-dev/artel/internal/cryptoutil"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware"
	repopg "github.com/ruf-dev/artel/internal/repository/pg"
	svcv1 "github.com/ruf-dev/artel/internal/service/v1"
	"github.com/ruf-dev/artel/internal/transport/external_connections_api"
	"github.com/ruf-dev/artel/internal/transport/vaults_api"
	"github.com/ruf-dev/artel/migrations"
)

const byokBufSize = 1024 * 1024

// ByokStorageSuite is the first suite in this codebase to drive requests through a real
// *grpc.Server (over an in-memory bufconn listener) instead of calling service methods directly
// in Go or hitting an http.Handler via httptest — it exercises the production
// middleware.GrpcAuthInterceptor end to end, so BYOK CouchDB/S3 connections registered via
// AddCouchDBConnection/AddS3Connection are verified reachable from the real API, not just from
// Go code calling the service layer directly.
type ByokStorageSuite struct {
	suite.Suite
	db              *sql.DB
	repos           *repopg.Repos
	svcs            *svcv1.Services
	couchInstanceID string

	grpcServer *grpc.Server
	listener   *bufconn.Listener
	conn       *grpc.ClientConn

	ecClient     pb.ExternalConnectionsAPIClient
	vaultsClient pb.VaultsAPIClient
}

func TestByokStorage(t *testing.T) {
	suite.Run(t, new(ByokStorageSuite))
}

func (s *ByokStorageSuite) SetupSuite() {
	ctx := context.Background()

	pgDSN := envOrDefault("PG_DSN", "postgres://artel:artel_db@localhost:15434/artel_db?sslmode=disable")
	db, err := sql.Open("postgres", pgDSN)
	s.Require().NoError(err, "open postgres")

	err = db.Ping()
	s.Require().NoError(err, "ping postgres — is the container running?")
	s.db = db

	err = migrations.ApplyMigration(db)
	s.Require().NoError(err, "run migrations")

	encKey := make([]byte, 32)
	encryptor, err := cryptoutil.NewAESEncryptor(encKey)
	s.Require().NoError(err, "create AES encryptor")
	s.repos = repopg.New(db, encryptor)

	cfg := config.EnvironmentConfig{}
	s.svcs, err = svcv1.New(s.repos, cfg)
	s.Require().NoError(err, "init services")

	err = s.repos.SystemSettings().CompleteSetup(ctx)
	s.Require().NoError(err, "complete setup so Register/Login aren't gated in tests")

	err = s.repos.SystemSettings().UpdateRegistrationMode(ctx, domain.RegistrationModeSelfRegister)
	s.Require().NoError(err, "enable self-registration so tests can Register directly")

	couchURL := envOrDefault("COUCH_URL", "http://localhost:15985")
	couchUser := envOrDefault("COUCH_USER", "admin")
	couchPass := envOrDefault("COUCH_PASS", "admin")

	s.couchInstanceID, err = s.svcs.CouchInstance.RegisterCouchInstance(ctx, couchURL, couchUser, couchPass)
	s.Require().NoError(err, "register couch instance")

	s.startGrpcServer()
}

// startGrpcServer builds a real *grpc.Server chained with the production auth interceptor
// (middleware.GrpcAuthInterceptor — not WithNoAuth), registers the ExternalConnectionsAPI and
// VaultsAPI implementations onto it, and serves it over an in-memory bufconn listener so the
// suite's RPCs travel through the real transport + auth stack without binding a TCP port.
func (s *ByokStorageSuite) startGrpcServer() {
	vaultsImpl := vaults_api.NewVaultsImpl(s.svcs.Vault, s.svcs.Workbench, false)
	externalConnectionsImpl := external_connections_api.New(s.svcs.ExternalConnectionService(), false)

	authOption := middleware.GrpcAuthInterceptor(s.svcs)

	s.grpcServer = grpc.NewServer(authOption)
	vaultsImpl.Register(s.grpcServer)
	externalConnectionsImpl.Register(s.grpcServer)

	s.listener = bufconn.Listen(byokBufSize)

	go func() {
		_ = s.grpcServer.Serve(s.listener)
	}()

	dialer := func(ctx context.Context, _ string) (net.Conn, error) {
		return s.listener.DialContext(ctx)
	}

	conn, err := grpc.NewClient(
		"passthrough:///bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	s.Require().NoError(err, "dial bufconn grpc client")
	s.conn = conn

	s.ecClient = pb.NewExternalConnectionsAPIClient(conn)
	s.vaultsClient = pb.NewVaultsAPIClient(conn)
}

func (s *ByokStorageSuite) TearDownSuite() {
	ctx := context.Background()

	if s.conn != nil {
		_ = s.conn.Close()
	}

	if s.grpcServer != nil {
		s.grpcServer.Stop()
	}

	if s.couchInstanceID != "" {
		err := s.svcs.CouchInstance.DeleteCouchInstance(ctx, s.couchInstanceID)
		s.NoError(err, "delete couch instance")
	}

	if s.db != nil {
		err := s.db.Close()
		s.NoError(err, "close db")
	}
}

// registerAndAuthenticate registers+logs in a fresh user via the service layer directly (mirrors
// e2e_test.go/quota_test.go), then returns a context carrying the session token as outgoing gRPC
// metadata the same way the production auth interceptor expects to read it back (see
// middleware/cookie_annotator.go's authHeader const and auth_interceptor.go's authWithSession).
func (s *ByokStorageSuite) registerAndAuthenticate(namePrefix string) (context.Context, uuid.UUID) {
	ctx := context.Background()

	email := randomEmail()
	password := "test-password-" + namePrefix

	user, err := s.svcs.Auth.Register(ctx, email, password)
	s.Require().NoError(err)
	s.T().Cleanup(func() {
		_ = s.repos.Users().Delete(context.Background(), user.Uuid)
	})

	// Email/password registration never populates domain.User.Username (defaults to "" — see
	// migrations/007_telegram_auth.sql). Vault creation derives the CouchDB database name from
	// it (sanitizeCouchDBName(uc.UserName + "-" + vaultName + "-vault")), and an empty UserName
	// produces a name starting with "-", which CouchDB rejects. Because this suite authenticates
	// through the real gRPC auth interceptor (which reads UserName off the persisted user row,
	// unlike e2e_test.go/quota_test.go which build a UserContext by hand), the stand-in has to be
	// written to the row itself rather than only into a locally built context.
	localPart := strings.SplitN(email, "@", 2)[0]

	_, err = s.db.Exec(`UPDATE users SET username = $1 WHERE id = $2`, localPart, user.Uuid)
	s.Require().NoError(err, "set username stand-in for vault creation")

	session, err := s.svcs.Auth.Login(ctx, email, password)
	s.Require().NoError(err)
	s.Require().NotEmpty(session.Token)

	authedCtx := metadata.AppendToOutgoingContext(ctx, "authorization", session.Token)

	return authedCtx, user.Uuid
}

func (s *ByokStorageSuite) TestAddCouchDBConnection_PersistsAndVaultUsesOwnedInstance() {
	ctx, userUuid := s.registerAndAuthenticate("couch_byok")

	// couch_instances.url carries a UNIQUE constraint (migrations/002_couch_instances.sql), so
	// this can't reuse the exact URL string SetupSuite already registered for the admin pool
	// instance — it would collide on that constraint. "127.0.0.1" resolves to the same published
	// docker-compose port as "localhost" (see tests/docker-compose.yaml), so this still targets
	// the same physical CouchDB container while keeping the row distinct.
	byokCouchURL := strings.Replace(envOrDefault("COUCH_URL", "http://localhost:15985"), "localhost", "127.0.0.1", 1)

	addReq := &pb.AddCouchDBConnection_Request{
		Url:      byokCouchURL,
		Username: envOrDefault("COUCH_USER", "admin"),
		Password: envOrDefault("COUCH_PASS", "admin"),
	}

	_, err := s.ecClient.AddCouchDBConnection(ctx, addReq)
	s.Require().NoError(err, "AddCouchDBConnection rpc should succeed")

	var provider string
	err = s.db.QueryRow(`SELECT provider FROM external_connections WHERE user_id = $1`, userUuid).Scan(&provider)
	s.Require().NoError(err, "external_connections row should have been persisted")
	s.Equal(domain.ProviderCouchDB, provider)

	var ownedInstanceIDs []uuid.UUID
	rows, err := s.db.Query(`SELECT id FROM couch_instances WHERE owner_user_id = $1`, userUuid)
	s.Require().NoError(err)

	for rows.Next() {
		var id uuid.UUID

		err = rows.Scan(&id)
		s.Require().NoError(err)

		ownedInstanceIDs = append(ownedInstanceIDs, id)
	}
	s.Require().NoError(rows.Err())
	s.Require().Len(ownedInstanceIDs, 1, "exactly one owned couch_instances row should exist for this user")

	ownedInstanceID := ownedInstanceIDs[0]

	createResp, err := s.vaultsClient.CreateVault(ctx, &pb.CreateVault_Request{Name: "byok_couch_vault"})
	s.Require().NoError(err, "CreateVault rpc should succeed")

	vaultUuid, err := uuid.Parse(createResp.Id)
	s.Require().NoError(err)
	s.T().Cleanup(func() {
		_ = s.svcs.Vault.DeleteVault(context.Background(), vaultUuid)
	})

	var vaultCouchInstanceID uuid.UUID

	err = s.db.QueryRow(`SELECT couch_instance_id FROM vaults WHERE id = $1`, vaultUuid).Scan(&vaultCouchInstanceID)
	s.Require().NoError(err)

	s.Equal(ownedInstanceID, vaultCouchInstanceID, "vault should resolve to the caller's owned BYOK instance")
	s.NotEqual(s.couchInstanceID, vaultCouchInstanceID.String(), "vault must not fall back to the shared admin pool instance")
}

func (s *ByokStorageSuite) TestAddS3Connection_PersistsAndLinkS3BucketAutoResolves() {
	ctx, userUuid := s.registerAndAuthenticate("s3_byok")

	addReq := &pb.AddS3Connection_Request{
		Endpoint:  envOrDefault("S3_ENDPOINT", "localhost:19000"),
		Region:    "us-east-1",
		AccessKey: "minioadmin",
		SecretKey: "minioadmin",
		UseSsl:    false,
		PathStyle: true,
	}

	_, err := s.ecClient.AddS3Connection(ctx, addReq)
	s.Require().NoError(err, "AddS3Connection rpc should succeed")

	var provider string
	err = s.db.QueryRow(`SELECT provider FROM external_connections WHERE user_id = $1`, userUuid).Scan(&provider)
	s.Require().NoError(err, "external_connections row should have been persisted")
	s.Equal(domain.ProviderS3, provider)

	var ownedS3InstanceID uuid.UUID

	err = s.db.QueryRow(`SELECT id FROM s3_instances WHERE owner_user_id = $1`, userUuid).Scan(&ownedS3InstanceID)
	s.Require().NoError(err, "owned s3_instances row should have been created")

	createResp, err := s.vaultsClient.CreateVault(ctx, &pb.CreateVault_Request{Name: "byok_s3_vault"})
	s.Require().NoError(err, "CreateVault rpc should succeed")

	vaultUuid, err := uuid.Parse(createResp.Id)
	s.Require().NoError(err)
	s.T().Cleanup(func() {
		_ = s.svcs.Vault.DeleteVault(context.Background(), vaultUuid)
	})

	bucketName := fmt.Sprintf("byok-e2e-%08x", rand.Uint32())

	// Empty S3InstanceId exercises the Task 1 transport fix: the handler must treat "" as
	// "not provided" and pass nil through to the service layer's auto-resolve (PickForUser)
	// instead of failing uuid.Parse("").
	linkReq := &pb.LinkS3Bucket_Request{
		VaultId:      vaultUuid.String(),
		S3InstanceId: "",
		BucketName:   bucketName,
	}

	_, err = s.vaultsClient.LinkS3Bucket(ctx, linkReq)
	s.Require().NoError(err, "LinkS3Bucket rpc with empty s3_instance_id should auto-resolve and succeed")

	var vaultS3InstanceID uuid.UUID

	err = s.db.QueryRow(`SELECT s3_instance_id FROM vaults WHERE id = $1`, vaultUuid).Scan(&vaultS3InstanceID)
	s.Require().NoError(err)
	s.Equal(ownedS3InstanceID, vaultS3InstanceID, "vault should resolve to the caller's owned BYOK s3 instance")

	bucketCfg := s3client.Config{
		Endpoint:  envOrDefault("S3_ENDPOINT", "localhost:19000"),
		Region:    "us-east-1",
		AccessKey: "minioadmin",
		SecretKey: "minioadmin",
		UseSSL:    false,
		PathStyle: true,
	}

	bucketClient, err := s3client.New(bucketCfg, bucketName)
	s.Require().NoError(err)

	_, err = bucketClient.List(context.Background())
	s.Require().NoError(err, "bucket should actually exist in MinIO, not just be recorded in postgres")
}

func (s *ByokStorageSuite) TestCreateVault_NoOwnedConnection_FallsBackToPool() {
	ctx, _ := s.registerAndAuthenticate("no_byok")

	createResp, err := s.vaultsClient.CreateVault(ctx, &pb.CreateVault_Request{Name: "pool_fallback_vault"})
	s.Require().NoError(err, "CreateVault rpc should succeed for a user with no BYOK connections")

	vaultUuid, err := uuid.Parse(createResp.Id)
	s.Require().NoError(err)
	s.T().Cleanup(func() {
		_ = s.svcs.Vault.DeleteVault(context.Background(), vaultUuid)
	})

	var vaultCouchInstanceID uuid.UUID

	err = s.db.QueryRow(`SELECT couch_instance_id FROM vaults WHERE id = $1`, vaultUuid).Scan(&vaultCouchInstanceID)
	s.Require().NoError(err)
	s.Equal(s.couchInstanceID, vaultCouchInstanceID.String(), "vault should fall back to the shared admin pool instance")

	var ownerUserID sql.NullString

	err = s.db.QueryRow(`SELECT owner_user_id FROM couch_instances WHERE id = $1`, vaultCouchInstanceID).Scan(&ownerUserID)
	s.Require().NoError(err)
	s.False(ownerUserID.Valid, "the shared admin pool instance must have no owner")
}

func (s *ByokStorageSuite) TestAddCouchDBConnection_BadCredentials_RejectsAndPersistsNothing() {
	ctx, userUuid := s.registerAndAuthenticate("bad_couch")

	addReq := &pb.AddCouchDBConnection_Request{
		Url:      "http://localhost:1",
		Username: "admin",
		Password: "admin",
	}

	_, err := s.ecClient.AddCouchDBConnection(ctx, addReq)
	s.Require().Error(err, "AddCouchDBConnection should reject an unreachable couchdb url")

	st, ok := status.FromError(err)
	s.Require().True(ok, "error should be a well-formed grpc status")
	s.NotEqual(codes.OK, st.Code(), "rejected credentials should surface a non-OK grpc status code")

	var count int

	err = s.db.QueryRow(`SELECT count(*) FROM external_connections WHERE user_id = $1`, userUuid).Scan(&count)
	s.Require().NoError(err)
	s.Equal(0, count, "no external_connections row should be persisted for rejected credentials")

	err = s.db.QueryRow(`SELECT count(*) FROM couch_instances WHERE owner_user_id = $1`, userUuid).Scan(&count)
	s.Require().NoError(err)
	s.Equal(0, count, "no owned couch_instances row should be persisted for rejected credentials")
}
