//go:build e2e && e2e_bootstrap
// +build e2e,e2e_bootstrap

// Package bootstrap holds the one-time setup/cleanup steps that used to be duplicated across every
// e2e suite's SetupSuite/TearDownSuite: applying migrations, wiping any state left over from a
// crashed prior run, and registering the shared CouchDB/S3 admin-pool instances. These now happen
// exactly once per full test run rather than once per suite/package.
//
// TestEnvSetup and TestEnvCleanup are invoked explicitly by `make test-e2e` (see the Makefile),
// before and after the real suite run (`go test -tags e2e ./tests/...`) — they are gated behind
// the additional e2e_bootstrap build tag (on top of e2e) specifically so `go test -tags e2e
// ./tests/...`'s normal package sweep never picks them up as a side effect; running this package
// requires the explicit `-tags "e2e e2e_bootstrap"`.
package bootstrap

import (
	"context"
	"testing"

	"github.com/ruf-dev/artel/internal/config"
	"github.com/ruf-dev/artel/tests/harness"
)

// TestEnvSetup prepares the shared e2e stack for a fresh suite run: applies migrations, wipes any
// state a previous (possibly crashed) run left behind, then provisions the CouchDB/S3 admin-pool
// instances every suite looks up via harness.GetCouchInstance/GetS3Instance.
func TestEnvSetup(t *testing.T) {
	ctx := context.Background()

	db := harness.OpenPostgres(t)
	harness.ApplyMigrations(t, db)

	couchURL := harness.CouchURL(t)
	couchUser, couchPass := harness.CouchCreds(t)
	s3Endpoint := harness.S3Endpoint(t)
	s3AccessKey, s3SecretKey := harness.S3Creds(t)

	// Reset first, in case a previous run crashed and left garbage behind.
	harness.ResetPostgres(t, ctx, db)
	harness.ResetCouchDB(t, ctx, couchURL, couchUser, couchPass)
	harness.ResetS3(t, ctx, s3Endpoint, s3AccessKey, s3SecretKey)

	repos, svcs, _ := harness.BuildServices(t, db, config.EnvironmentConfig{})

	harness.CompleteSetup(t, ctx, repos)

	harness.ProvisionCouchInstance(t, ctx, svcs, couchURL, couchUser, couchPass)
	harness.ProvisionS3Instance(t, ctx, svcs, s3Endpoint, s3AccessKey, s3SecretKey)
}

// TestEnvCleanup wipes all test-created state after the suite run finishes, leaving the shared
// stack empty for whoever inspects it next.
func TestEnvCleanup(t *testing.T) {
	ctx := context.Background()

	db := harness.OpenPostgres(t)

	couchURL := harness.CouchURL(t)
	couchUser, couchPass := harness.CouchCreds(t)
	s3Endpoint := harness.S3Endpoint(t)
	s3AccessKey, s3SecretKey := harness.S3Creds(t)

	harness.ResetPostgres(t, ctx, db)
	harness.ResetCouchDB(t, ctx, couchURL, couchUser, couchPass)
	harness.ResetS3(t, ctx, s3Endpoint, s3AccessKey, s3SecretKey)
}
