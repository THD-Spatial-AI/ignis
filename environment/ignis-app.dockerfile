# syntax=docker/dockerfile:1

# ---------------------------------------------------------------------------
# Stage 1: builder — has the full Go toolchain + source, never shipped
# ---------------------------------------------------------------------------
ARG GO_VERSION=1.26.5
FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /app

# Cache module downloads separately from source changes
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Passed by the publish workflow on a tag build; default to the same values
# internal/version hard-codes, so a plain `docker build` is unchanged.
ARG VERSION=dev
ARG COMMIT=none
ARG DATE=unknown
# -buildvcs=false: the build context may contain .git, and the toolchain's
# VCS stamping shells out to git, which fails on an ownership mismatch in CI.
# The image is identified by its release tag, so the stamp buys nothing.
# -X: same three variables release.yml stamps into the standalone binaries.
RUN CGO_ENABLED=0 go build -buildvcs=false \
    -ldflags="-s -w \
      -X github.com/thd-spatial-ai/ignis/internal/version.Version=${VERSION} \
      -X github.com/thd-spatial-ai/ignis/internal/version.Commit=${COMMIT} \
      -X github.com/thd-spatial-ai/ignis/internal/version.Date=${DATE}" \
    -o /out/ignis ./cmd/ignis

# ---------------------------------------------------------------------------
# Stage 2: final — only the compiled server binary, no compiler/source/git
# ---------------------------------------------------------------------------
FROM debian:bookworm-slim

# ca-certificates: outbound TLS: wget: used by the Compose healthcheck below.
# No postgresql-client / build-essential / git here — ignis talks to Postgres
# purely over the wire (pgx, pure Go driver), it never shells out to psql.
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates wget \
    && rm -rf /var/lib/apt/lists/*

RUN useradd -m -u 10001 appuser
WORKDIR /app
COPY --from=builder /out/ignis /app/bin/ignis
RUN chown -R appuser:appuser /app

USER appuser

CMD ["./bin/ignis"]
