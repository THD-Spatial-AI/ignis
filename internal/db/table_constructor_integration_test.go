//go:build integration

package importer_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/thd-spatial-ai/ignis/internal/config"
	importer "github.com/thd-spatial-ai/ignis/internal/db"
	"github.com/thd-spatial-ai/ignis/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xuri/excelize/v2"
)

var testPool *pgxpool.Pool

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

	os.Exit(m.Run())
}

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

// writeFixtureWorkbook creates a minimal xlsx with a "Calc.Set.Building"
// sheet in the same layout as the real TABULA workbook: header row 1, data
// type row 6, data starting row 13 - matching what extractHeaders/
// extractCountryCodes (see table_constructor.go) expect.
func writeFixtureWorkbook(t *testing.T) string {
	t.Helper()
	f := excelize.NewFile()
	sheet := "Calc.Set.Building"
	if _, err := f.NewSheet(sheet); err != nil {
		t.Fatal(err)
	}

	headers := []string{"Code_BuildingVariant", "Code_Country", "Code_ComplexRoof", "A_Roof_1", "HeatingDays"}
	dataTypes := []string{"VarChar", "VarChar", "VarChar", "Real", "Integer"}
	for col, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheet, cell, h)
	}
	for col, dt := range dataTypes {
		cell, _ := excelize.CoordinatesToCellName(col+1, 6)
		f.SetCellValue(sheet, cell, dt)
	}

	countries := []string{"DE", "AT"}
	for i := 0; i < 20; i++ {
		row := 13 + i // rows[12] in 0-indexed = row 13 in 1-indexed sheet
		values := []interface{}{
			fmt.Sprintf("%s.N.SFH.%02d", countries[i%len(countries)], i),
			countries[i%len(countries)],
			"yes",
			45.5,
			180,
		}
		for col, v := range values {
			cell, _ := excelize.CoordinatesToCellName(col+1, row)
			f.SetCellValue(sheet, cell, v)
		}
	}

	path := filepath.Join(t.TempDir(), "fixture.xlsx")
	if err := f.SaveAs(path); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestTableConstructor_Run_createsAndPopulatesCountryTables(t *testing.T) {
	xlsxPath := writeFixtureWorkbook(t)
	cfg := &config.Config{
		Data: &config.DataPaths{ExcelFile: xlsxPath},
		DB:   &config.DBConfig{Schemas: &config.Schemas{Tabula: "tc_test"}},
	}

	tc := importer.NewTableConstructor(testPool, cfg)
	if err := tc.Run(); err != nil {
		t.Fatalf("Run() unexpected error: %v", err)
	}

	ctx := context.Background()

	var tableCount int
	err := testPool.QueryRow(ctx,
		`SELECT count(*) FROM information_schema.tables WHERE table_schema = $1 AND table_name = $2`,
		"tc_test", "germany",
	).Scan(&tableCount)
	if err != nil {
		t.Fatalf("checking table existence: %v", err)
	}
	if tableCount != 1 {
		t.Fatalf("expected tc_test.germany to exist, tableCount=%d", tableCount)
	}

	var rowCount int
	if err := testPool.QueryRow(ctx, `SELECT count(*) FROM tc_test.germany`).Scan(&rowCount); err != nil {
		t.Fatalf("counting rows: %v", err)
	}
	if rowCount != 10 { // half of the 20 fixture rows are DE
		t.Errorf("row count = %d, want 10", rowCount)
	}

	var aRoof1 float64
	if err := testPool.QueryRow(ctx, `SELECT "A_Roof_1" FROM tc_test.germany LIMIT 1`).Scan(&aRoof1); err != nil {
		t.Fatalf("reading A_Roof_1: %v", err)
	}
	if aRoof1 != 45.5 {
		t.Errorf("A_Roof_1 = %v, want 45.5", aRoof1)
	}
}

func TestTableConstructor_Run_missingWorkbook_returnsError(t *testing.T) {
	cfg := &config.Config{
		Data: &config.DataPaths{ExcelFile: "/nonexistent/workbook.xlsx"},
		DB:   &config.DBConfig{Schemas: &config.Schemas{Tabula: "tc_test"}},
	}
	tc := importer.NewTableConstructor(testPool, cfg)
	if err := tc.Run(); err == nil {
		t.Fatal("expected error for missing workbook")
	}
}
