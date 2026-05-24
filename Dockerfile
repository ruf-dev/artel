# Stage 1: build React/Bun UI
FROM --platform=$BUILDPLATFORM oven/bun:1 AS ui-builder

WORKDIR /ui

COPY pkg/client/ArtelUI/package.json pkg/client/ArtelUI/bun.lock ./
RUN bun install --frozen-lockfile

COPY pkg/client/ArtelUI/ ./
RUN bun run build

# Stage 2: build Go binary (embeds UI dist)
FROM --platform=$BUILDPLATFORM golang:1.26.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg go mod download

COPY . .
COPY --from=ui-builder /ui/dist ./internal/transport/ui/dist

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 \
    go build -o /deploy/server/service ./cmd/service/main.go && \
    cp -r config /deploy/server/config && \
    if [ -d "./migrations" ];  then \
      cp -r ./migrations /deploy/server/migrations;\
    fi

# Stage 3: minimal scratch image
FROM scratch

LABEL MATRESHKA_CONFIG_ENABLED=true

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /deploy/server/ /app/

WORKDIR /app

EXPOSE 1551

ENTRYPOINT ["/app/service"]
