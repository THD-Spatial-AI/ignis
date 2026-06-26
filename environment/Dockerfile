FROM debian:bookworm-slim

ARG DEBIAN_FRONTEND=noninteractive
ENV GO_VERSION=1.26.0

# -----------------------------
# Install system dependencies
# -----------------------------
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential git wget ca-certificates \
    postgresql-client \
    && rm -rf /var/lib/apt/lists/*

# -----------------------------
# Install Go
# -----------------------------
RUN wget -q https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz -O /tmp/go.tar.gz \
    && rm -rf /usr/local/go \
    && tar -C /usr/local -xzf /tmp/go.tar.gz \
    && rm /tmp/go.tar.gz
ENV GOROOT="/usr/local/go"
ENV PATH="${GOROOT}/bin:${PATH}"

# -----------------------------
# Copy and build project (as root)
# -----------------------------
WORKDIR /app
COPY . .
RUN go mod tidy \
    && go build -o bin/app      cmd/app/main.go \
    && go build -o bin/build_db cmd/build_db/main.go \
    && go build -o bin/validate cmd/validate/main.go

# -----------------------------
# Create non-root user and fix permissions
# -----------------------------
RUN useradd -m -u 10001 appuser \
    && mkdir -p /home/appuser/go/bin \
    && chown -R appuser:appuser /home/appuser /app
ENV GOPATH="/home/appuser/go"
ENV PATH="${GOPATH}/bin:${PATH}"

USER appuser

# -----------------------------
# Set up shell environment
# -----------------------------
RUN echo 'export GOROOT=/usr/local/go' >> /home/appuser/.bashrc \
    && echo 'export GOPATH=/home/appuser/go' >> /home/appuser/.bashrc \
    && echo 'export PATH=$GOROOT/bin:$GOPATH/bin:$PATH' >> /home/appuser/.bashrc

# -----------------------------
# Default command
# -----------------------------
CMD ["bash", "-lc", "echo 'Go:' $(go version); exec bash"]
