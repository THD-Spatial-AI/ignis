package utils

import (
	"context"
	"testing"
	"time"

	"github.com/thd-spatial-ai/ignis/internal/config"
)

func TestBuildConnectionString(t *testing.T) {
	cfg := config.Config{
		DB: &config.DBConfig{
			User:     "tester",
			Password: "secret",
			Host:     "db.internal",
			Port:     "5433",
			Name:     "ignis_test",
			SSLMode:  "disable",
		},
	}

	got := BuildConnectionString(cfg)
	want := "postgres://tester:secret@db.internal:5433/ignis_test?sslmode=disable"
	if got != want {
		t.Errorf("BuildConnectionString = %q, want %q", got, want)
	}
}

func TestConnectPool_invalidConnString(t *testing.T) {
	ctx := context.Background()
	pool, err := ConnectPool(ctx, "not-a-valid-connection-string")
	if err == nil {
		if pool != nil {
			pool.Close()
		}
		t.Fatal("expected error for malformed connection string, got nil")
	}
	if pool != nil {
		t.Error("expected nil pool on error")
	}
}

func TestConnectPool_pingFailure(t *testing.T) {
	// A well-formed connection string pointing at a port nothing listens on:
	// pgxpool.New succeeds (it doesn't dial), so this exercises the Ping failure branch.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	pool, err := ConnectPool(ctx, "postgres://user:pass@127.0.0.1:1/nonexistent?sslmode=disable&connect_timeout=1")
	if err == nil {
		if pool != nil {
			pool.Close()
		}
		t.Fatal("expected ping error for unreachable database, got nil")
	}
	if pool != nil {
		t.Error("expected nil pool on ping failure")
	}
}
