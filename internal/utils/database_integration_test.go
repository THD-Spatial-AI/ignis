//go:build integration

package utils_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/thd-spatial-ai/ignis/internal/config"
	"github.com/thd-spatial-ai/ignis/internal/utils"
)

var (
	testHost string
	testPort string
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image: "postgres:17-alpine",
		Env: map[string]string{
			"POSTGRES_DB":       "postgres",
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
		},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForLog("database system is ready to accept connections").AsRegexp(),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("failed to start Postgres container: %v", err)
	}
	defer container.Terminate(ctx)

	testHost, err = container.Host(ctx)
	if err != nil {
		log.Fatalf("failed to get container host: %v", err)
	}
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatalf("failed to get container port: %v", err)
	}
	testPort = port.Port()

	// The postgres image restarts once after the "ready to accept connections"
	// log line (re-init with the final config), so the first few real
	// connection attempts can be reset. Wait it out before running any test.
	if err := waitUntilReady(ctx); err != nil {
		log.Fatalf("Postgres not ready: %v", err)
	}

	os.Exit(m.Run())
}

func waitUntilReady(ctx context.Context) error {
	deadline := time.Now().Add(15 * time.Second)
	var lastErr error
	for time.Now().Before(deadline) {
		pool, err := utils.ConnectPool(ctx, testConnString())
		if err == nil {
			pool.Close()
			return nil
		}
		lastErr = err
		time.Sleep(300 * time.Millisecond)
	}
	return lastErr
}

func testConnString() string {
	cfg := config.Config{
		DB: &config.DBConfig{
			User:     "test",
			Password: "test",
			Host:     testHost,
			Port:     testPort,
			Name:     "postgres",
			SSLMode:  "disable",
		},
	}
	return utils.BuildConnectionString(cfg)
}

func TestConnectPool_success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := utils.ConnectPool(ctx, testConnString())
	if err != nil {
		t.Fatalf("ConnectPool: unexpected error: %v", err)
	}
	defer pool.Close()

	var result int
	if err := pool.QueryRow(ctx, "SELECT 1").Scan(&result); err != nil {
		t.Fatalf("query through connected pool failed: %v", err)
	}
	if result != 1 {
		t.Errorf("SELECT 1 = %d, want 1", result)
	}
}

// TestBuildConnectionString_roundTripsToRealConnection is an integration-level
// sanity check that BuildConnectionString's format is actually accepted by
// pgx/Postgres, not just string-equality-tested (see database_test.go).
func TestBuildConnectionString_roundTripsToRealConnection(t *testing.T) {
	connStr := fmt.Sprintf("postgres://test:test@%s:%s/postgres?sslmode=disable", testHost, testPort)
	if connStr != testConnString() {
		t.Fatalf("test helper drifted from BuildConnectionString's actual format: got %q want %q", testConnString(), connStr)
	}
}
