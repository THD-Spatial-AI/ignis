# syntax=docker/dockerfile:1

# ---------------------------------------------------------------------------
# Stage 1: builder — has the full Go toolchain + source, never shipped
# ---------------------------------------------------------------------------
FROM golang:1.26.5-bookworm AS builder

WORKDIR /app

# Cache module downloads separately from source changes
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /out/ignis ./cmd/ignis

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
