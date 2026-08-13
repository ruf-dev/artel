all: codegen build-ui

# Golang service scripts

codegen:
	@echo --- Generating proto contracts ---
	moti g
	@echo --- Generating Sqlc queries ---
	sqlc generate

# Build Web User Interface from pkg/client/ArtelUI
build-ui:
	@echo --- Building WebUI ---
	cd pkg/client/ArtelUI &&  bun i &&  bun run build
	@echo --- Moving WebUI to go app to serve ---
	rm -rf internal/transport/ui/dist
	cp -r pkg/client/ArtelUI/dist internal/transport/ui/dist

build-local-container:
	docker buildx build \
			--load \
			--platform linux/arm64 \
			-t artel:local .

lint: .lint-go .lint-react

.lint-go:
	go fmt ./...
	golangci-lint run
.lint-react:
	cd pkg/client/ArtelUI && bun run lint
### Web Client Setup
client-setup:
	cd pkg/client/ArtelUI && bun i

client:
	cd pkg/client/ArtelUI && vite

# Runs the Go backend (-dev flag loads config/dev.yaml) and the Vite dev server
# together, so `make dev-serve` gives a full local dev stack in one command.
# trap kill 0 tears down both processes together on Ctrl+C or exit.
dev-serve:
	@trap 'kill 0' EXIT; \
	(cd pkg/client/ArtelUI && vite) & \
	go run ./cmd/service -dev; \
	wait

# Applies Postgres migrations as a standalone step, independent of app startup — for operators
# who set skip_migrations=true and want to run migrations explicitly instead.
migrate:
	go run ./cmd/migrate -dev

### local dev environment
setup-dev-env:
	docker compose up -d
	./scripts/setup_dev_garage.sh
	docker compose -f tests/docker-compose.yaml up -d test-dockerd
	./scripts/setup_dev_docker_host.sh

### Garage Setup
garage-status:
	docker exec -it artel-garage-s3 /garage status

setup-dev-garage:
	./scripts/setup_dev_garage.sh

### Docker Host (local dind) Setup
# test-dockerd (tests/docker-compose.yaml) is the local docker:dind daemon registered as a
# docker_hosts row for workbench.Service to schedule containers against — see
# docs/workbench/02_docker_topology.md. It lives in tests/docker-compose.yaml, not the root
# docker-compose.yaml, since it's also what `go test -tags e2e ./tests/...` uses; this target
# just also brings it up for non-e2e local dev.
setup-dev-docker-host:
	docker compose -f tests/docker-compose.yaml up -d test-dockerd
	./scripts/setup_dev_docker_host.sh

### E2E testing
# Brings up tests/docker-compose.yaml (postgres, couchdb, minio, test-dockerd) and blocks
# until every service reports healthy (via each service's healthcheck), instead of just
# `up -d` returning as soon as the containers start — this is what previously let e2e tests
# race a backend (e.g. CouchDB) that hadn't finished starting yet.
setup-e2e-env:
	docker compose -f tests/docker-compose.yaml up -d --wait

# Runs the e2e-tagged test suites against a ready environment. The compose stack is left
# running afterward (see setup-e2e-env) so repeated `make test-e2e` runs don't pay the
# container startup cost each time.
#
# Provisioning (migrations, CouchDB/S3 admin-pool instances, system_settings setup) now happens
# exactly once per run via tests/bootstrap's TestEnvSetup, instead of once per suite — so suites no
# longer race each other's goose migrations and can run with Go's default cross-package
# parallelism (no more -p 1). tests/bootstrap is gated behind an extra e2e_bootstrap build tag so
# it's invoked explicitly here rather than picked up by the `./tests/...` wildcard sweep below.
# TestEnvCleanup always runs after, even on suite failure, via capturing $$status and exiting with
# it at the end — leaving the shared stack empty for whoever inspects it next.
test-e2e: setup-e2e-env
	go test -tags "e2e e2e_bootstrap" -count=1 ./tests/bootstrap/... -run TestEnvSetup
	go test -tags e2e ./tests/...; status=$$?; \
	go test -tags "e2e e2e_bootstrap" -count=1 ./tests/bootstrap/... -run TestEnvCleanup; \
	exit $$status
