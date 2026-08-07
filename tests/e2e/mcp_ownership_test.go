//go:build e2e
// +build e2e

package e2e_test

import (
	"context"
	"fmt"
	"math/rand/v2"
	"testing"

	_ "github.com/lib/pq"
	"github.com/ruf-dev/artel/internal/config"
	"github.com/ruf-dev/artel/internal/domain"
	repopg "github.com/ruf-dev/artel/internal/repository/pg"
	"github.com/ruf-dev/artel/tests/harness"
	"github.com/stretchr/testify/suite"
)

// McpOwnershipSuite exercises the mcps.owner_user_id / mcps.is_community columns added in
// migration 060: Upsert must persist them, Get and List must round-trip them back into
// domain.McpDefinition, both for a "community" MoM owned by a user and a "system" MoM with no
// owner (nil OwnerUserUuid, IsCommunity false — the shape existing seeded MoMs like email/gitlab
// keep untouched).
type McpOwnershipSuite struct {
	suite.Suite
	repos *repopg.Repos
}

func TestMcpOwnership(t *testing.T) {
	suite.Run(t, new(McpOwnershipSuite))
}

func (s *McpOwnershipSuite) SetupSuite() {
	db := harness.OpenPostgres(s.T())

	cfg := config.EnvironmentConfig{}

	s.repos, _, _ = harness.BuildServices(s.T(), db, cfg)
}

func randomMcpName() string {
	return fmt.Sprintf("e2e_ownership_mom_%08x", rand.Uint32())
}

func (s *McpOwnershipSuite) TestCommunityMomRoundTrip() {
	ctx := context.Background()

	email := harness.Slug(s.T()) + "@test.local"

	user, err := s.repos.Users().Create(ctx, email, "unused-hash")
	s.Require().NoError(err)
	s.T().Cleanup(func() {
		_ = s.repos.Users().Delete(context.Background(), user.Uuid)
	})

	name := randomMcpName()
	def := domain.McpDefinition{
		Name:          name,
		Author:        "e2e-test",
		Description:   "community MoM ownership round trip",
		OwnerUserUuid: &user.Uuid,
		IsCommunity:   true,
	}
	s.T().Cleanup(func() {
		_ = s.repos.McpDefinitions().Delete(context.Background(), name)
	})

	upserted, err := s.repos.McpDefinitions().Upsert(ctx, def)
	s.Require().NoError(err)
	s.Require().NotNil(upserted.OwnerUserUuid)
	s.Equal(user.Uuid, *upserted.OwnerUserUuid)
	s.True(upserted.IsCommunity)

	got, err := s.repos.McpDefinitions().Get(ctx, name)
	s.Require().NoError(err)
	s.Require().True(got.Valid)
	s.Require().NotNil(got.V.OwnerUserUuid)
	s.Equal(user.Uuid, *got.V.OwnerUserUuid)
	s.True(got.V.IsCommunity)

	all, err := s.repos.McpDefinitions().List(ctx)
	s.Require().NoError(err)

	found := false
	for _, d := range all {
		if d.Name != name {
			continue
		}
		found = true
		s.Require().NotNil(d.OwnerUserUuid)
		s.Equal(user.Uuid, *d.OwnerUserUuid)
		s.True(d.IsCommunity)
	}
	s.True(found, "upserted MoM must appear in List")
}

func (s *McpOwnershipSuite) TestSystemMomHasNoOwner() {
	ctx := context.Background()

	name := randomMcpName()
	def := domain.McpDefinition{
		Name:        name,
		Author:      "e2e-test",
		Description: "system MoM has no owner",
	}
	s.T().Cleanup(func() {
		_ = s.repos.McpDefinitions().Delete(context.Background(), name)
	})

	upserted, err := s.repos.McpDefinitions().Upsert(ctx, def)
	s.Require().NoError(err)
	s.Nil(upserted.OwnerUserUuid)
	s.False(upserted.IsCommunity)

	got, err := s.repos.McpDefinitions().Get(ctx, name)
	s.Require().NoError(err)
	s.Require().True(got.Valid)
	s.Nil(got.V.OwnerUserUuid)
	s.False(got.V.IsCommunity)
}
