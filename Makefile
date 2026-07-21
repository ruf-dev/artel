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

### Garage Setup
garage-status:
	docker exec -it artel-garage-s3 /garage status

garage-setup:
	@echo INSTEAD OF '9961af1a4239d968' - use your id from `make garage-status`
	docker exec -it artel-garage-s3 /garage layout assign 9961af1a4239d968 -z z1 -c 1G
	docker exec -it artel-garage-s3 /garage layout apply --version 1

garage-gen-api-key:
	docker exec -it artel-garage-s3 /garage key create api-key
