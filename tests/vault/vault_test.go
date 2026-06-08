package vault_test

import (
	"context"
	"fmt"
	"math/rand/v2"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/ruf-dev/artel/internal/clients/couchdb"
)

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

type VaultDeletionSuite struct {
	suite.Suite
	admin *couchdb.Client
}

func TestVaultDeletion(t *testing.T) {
	suite.Run(t, new(VaultDeletionSuite))
}

func (s *VaultDeletionSuite) SetupSuite() {
	s.admin = couchdb.New(couchdb.Config{
		BaseURL:  envOrDefault("COUCH_URL", "http://localhost:15985"),
		User:     envOrDefault("COUCH_USER", "admin"),
		Password: envOrDefault("COUCH_PASS", "admin"),
	})

	err := s.admin.Setup(context.Background())
	s.Require().NoError(err, "CouchDB system setup failed — is the container running?")
}

func (s *VaultDeletionSuite) TestDeleteDatabase_RemovesDatabaseFromCouch() {
	ctx := context.Background()
	dbName := fmt.Sprintf("test_vault_delete_%08x", rand.Uint32())

	err := s.admin.CreateDatabase(ctx, dbName)
	s.Require().NoError(err)

	dbs, err := s.admin.ListDatabases(ctx)
	s.Require().NoError(err)
	require.Contains(s.T(), dbs, dbName, "database should exist before deletion")

	err = s.admin.DeleteDatabase(ctx, dbName)
	s.Require().NoError(err)

	dbs, err = s.admin.ListDatabases(ctx)
	s.Require().NoError(err)
	require.NotContains(s.T(), dbs, dbName, "database should be gone after deletion")
}
