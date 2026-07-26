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
