//go:build integration

package repository_test

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/thd-spatial-ai/ignis/internal/config"
	"github.com/thd-spatial-ai/ignis/internal/db/repository"
	"github.com/thd-spatial-ai/ignis/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testPool *pgxpool.Pool

// TestMain starts a single Postgres container for the whole package and seeds
// tabula.germany with quoted, mixed-case columns ("Code_BuildingVariant", ...)
// - the shape TableConstructor actually produces in production, and what
// TabulaRepository's quoted queries expect.
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

	host, err := container.Host(ctx)
	if err != nil {
		log.Fatalf("failed to get container host: %v", err)
	}
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatalf("failed to get container port: %v", err)
	}

	connStr := utils.BuildConnectionString(config.Config{
		DB: &config.DBConfig{User: "test", Password: "test", Host: host, Port: port.Port(), Name: "postgres", SSLMode: "disable"},
	})

	pool, err := waitForPool(ctx, connStr)
	if err != nil {
		log.Fatalf("Postgres not ready: %v", err)
	}
	testPool = pool
	defer testPool.Close()

	if err := seedFixtures(ctx, testPool); err != nil {
		log.Fatalf("failed to seed fixtures: %v", err)
	}

	os.Exit(m.Run())
}

// waitForPool retries ConnectPool for up to 15s: the postgres image restarts
// once after its "ready" log line, so the first real connection attempts can
// be reset.
func waitForPool(ctx context.Context, connStr string) (*pgxpool.Pool, error) {
	deadline := time.Now().Add(15 * time.Second)
	var lastErr error
	for time.Now().Before(deadline) {
		pool, err := utils.ConnectPool(ctx, connStr)
		if err == nil {
			return pool, nil
		}
		lastErr = err
		time.Sleep(300 * time.Millisecond)
	}
	return nil, lastErr
}

func seedFixtures(ctx context.Context, pool *pgxpool.Pool) error {
	statements := []string{
		`CREATE SCHEMA IF NOT EXISTS tabula`,
		`CREATE TABLE tabula.germany (
			id SERIAL PRIMARY KEY,
			"Code_BuildingVariant" VARCHAR,
			"A_Roof_1" REAL,
			"HeatingDays" INTEGER,
			"q_h_nd" REAL
		)`,
		`INSERT INTO tabula.germany ("Code_BuildingVariant", "A_Roof_1", "HeatingDays", "q_h_nd") VALUES
			('DE.N.SFH.01.Gen', 75.5, 185, 123.45),
			('DE.N.SFH.01.ReEx', 75.5, 185, 98.70),
			('DE.N.MFH.01.Gen', 200.0, 185, 88.10)`,
	}
	for _, stmt := range statements {
		if _, err := pool.Exec(ctx, stmt); err != nil {
			return err
		}
	}
	return nil
}

// --- TabulaRepository ---

func TestTabulaRepository_ListVariants_success(t *testing.T) {
	r := repository.NewTabulaRepository(testPool, "tabula")
	variants, err := r.ListVariants(context.Background(), "germany")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"DE.N.MFH.01.Gen", "DE.N.SFH.01.Gen", "DE.N.SFH.01.ReEx"}
	if len(variants) != len(want) {
		t.Fatalf("variants = %v, want %v", variants, want)
	}
	for i := range want {
		if variants[i] != want[i] {
			t.Errorf("variants[%d] = %q, want %q", i, variants[i], want[i])
		}
	}
}

func TestTabulaRepository_ListVariants_queryError(t *testing.T) {
	r := repository.NewTabulaRepository(testPool, "tabula")
	_, err := r.ListVariants(context.Background(), "nonexistent_table")
	if err == nil {
		t.Fatal("expected error for nonexistent table")
	}
}

func TestTabulaRepository_MatchVariants_success(t *testing.T) {
	r := repository.NewTabulaRepository(testPool, "tabula")
	codes, err := r.MatchVariants(context.Background(), "germany", "DE.N.SFH.01")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(codes) != 2 {
		t.Fatalf("codes = %v, want 2 SFH.01 variants", codes)
	}
}

func TestTabulaRepository_MatchVariants_noMatches(t *testing.T) {
	r := repository.NewTabulaRepository(testPool, "tabula")
	codes, err := r.MatchVariants(context.Background(), "germany", "DE.N.SFH.99")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(codes) != 0 {
		t.Errorf("codes = %v, want none", codes)
	}
}

func TestTabulaRepository_GetVariant_success(t *testing.T) {
	r := repository.NewTabulaRepository(testPool, "tabula")
	params, buildingID, expectedQHND, err := r.GetVariant(context.Background(), "germany", "DE.N.SFH.01.Gen")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buildingID != "DE.N.SFH.01.Gen" {
		t.Errorf("buildingID = %q, want %q", buildingID, "DE.N.SFH.01.Gen")
	}
	if expectedQHND != float64(float32(123.45)) {
		t.Errorf("expectedQHND = %v, want %v", expectedQHND, float64(float32(123.45)))
	}
	if params.BasicParameters.Envelope.A_Roof_1 != float64(float32(75.5)) {
		t.Errorf("A_Roof_1 = %v, want %v", params.BasicParameters.Envelope.A_Roof_1, float64(float32(75.5)))
	}
	if params.AdvancedParameters.ClimateConditions.HeatingDays != 185 {
		t.Errorf("HeatingDays = %d, want 185", params.AdvancedParameters.ClimateConditions.HeatingDays)
	}
}

func TestTabulaRepository_GetVariant_notFound(t *testing.T) {
	r := repository.NewTabulaRepository(testPool, "tabula")
	_, _, _, err := r.GetVariant(context.Background(), "germany", "DE.N.SFH.99.Gen")
	if !errors.Is(err, repository.ErrVariantNotFound) {
		t.Errorf("err = %v, want ErrVariantNotFound", err)
	}
}

func TestTabulaRepository_GetVariant_queryError(t *testing.T) {
	r := repository.NewTabulaRepository(testPool, "tabula")
	_, _, _, err := r.GetVariant(context.Background(), "nonexistent_table", "DE.N.SFH.01.Gen")
	if err == nil {
		t.Fatal("expected error for nonexistent table")
	}
	if errors.Is(err, repository.ErrVariantNotFound) {
		t.Error("a query error against a missing table should not be reported as ErrVariantNotFound")
	}
}
