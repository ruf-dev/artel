all: codegen build-ui

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

up:
	docker compose up -d

down:
	docker compose down

build-local-container:
	docker buildx build \
			--load \
			--platform linux/arm64 \
			-t artel:local .

lint:
	golangci-lint run ./...
