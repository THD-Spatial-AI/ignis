package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/thd-spatial-ai/ignis/internal/config"
	"github.com/thd-spatial-ai/ignis/internal/models"
	"github.com/thd-spatial-ai/ignis/internal/pipeline"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TestResult holds the results of a single pipeline test
type TestResult struct {
	RowID          int
	BuildingID     string
	CalculatedQHND float64
	ExpectedQHND   float64
	Difference     float64
	PercentError   float64
	Passed         bool
	ErrorMessage   string
}

const tolerancePercent = 2.5 // percent, allowed deviation from the TABULA reference q_h_nd

// cfg is loaded in main, not at package init, so the package can be tested
// without a database configuration.
var cfg config.Config

func main() {
	strict := flag.Bool("strict", false, "exit 1 when any building falls outside the tolerance")
	flag.Parse()

	fmt.Println("=== ignis Validation Tool ===")
	startTime := time.Now()
	cfg = config.LoadConfig()

	fmt.Printf("Database: %s@%s:%s/%s\n\n", cfg.DB.User, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)

	// Connect to database
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name, cfg.DB.SSLMode)

	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Test connection
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	fmt.Println("✓ Database connection successful")

	// Get all table names from database
	tableNames := []string{}
	rows, err := pool.Query(context.Background(), fmt.Sprintf(`SELECT table_name FROM information_schema.tables WHERE table_schema='%s' AND table_type='BASE TABLE' ORDER BY table_name`, cfg.DB.Schemas.Tabula))
	if err != nil {
		log.Fatalf("Failed to query table names: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			log.Printf("Failed to scan table name: %v", err)
			continue
		}
		tableNames = append(tableNames, tableName)
	}

	if len(tableNames) == 0 {
		log.Fatalf("No tables found in the database")
	}

	fmt.Printf("Found %d tables in the database:\n", len(tableNames))
	for _, tn := range tableNames {
		fmt.Printf(" - %s\n", tn)
	}

	results := make(map[string][]TestResult, len(tableNames))
	for _, tableName := range tableNames {
		results[tableName] = testAllBuildings(pool, *cfg.DB, tableName)
	}

	failed := report(os.Stdout, tableNames, results)
	fmt.Printf("\n=== Validation completed in %s ===\n", time.Since(startTime))

	pool.Close()
	os.Exit(exitCode(failed, *strict))
}

// report writes each table's pass count, then every failing building, then
// the overall count, and returns the number of buildings that did not pass.
// A building with an ErrorMessage counts as failed.
func report(w io.Writer, tableNames []string, results map[string][]TestResult) int {
	fmt.Fprintf(w, "\nPass rate per table (tolerance %.1f%%):\n", tolerancePercent)
	var failures []string
	total, passed := 0, 0
	for _, tableName := range tableNames {
		tablePassed := 0
		for _, r := range results[tableName] {
			switch {
			case r.ErrorMessage != "":
				failures = append(failures, fmt.Sprintf("  %-15s %-35s %s", tableName, r.BuildingID, r.ErrorMessage))
			case r.Passed:
				tablePassed++
			default:
				failures = append(failures, fmt.Sprintf("  %-15s %-35s calculated %10.2f  expected %10.2f  error %.2f%%",
					tableName, r.BuildingID, r.CalculatedQHND, r.ExpectedQHND, r.PercentError))
			}
		}
		n := len(results[tableName])
		total += n
		passed += tablePassed
		fmt.Fprintf(w, "  %-15s %4d/%-4d %6.1f%%\n", tableName, tablePassed, n, percent(tablePassed, n))
	}

	if len(failures) > 0 {
		fmt.Fprintf(w, "\nFailures (%d):\n", len(failures))
		for _, line := range failures {
			fmt.Fprintln(w, line)
		}
	}
	fmt.Fprintf(w, "\nTotal: %d/%d passed (%.1f%%)\n", passed, total, percent(passed, total))
	return total - passed
}

func percent(part, whole int) float64 {
	if whole == 0 {
		return 0
	}
	return float64(part) / float64(whole) * 100
}

// exitCode is 1 only under -strict with at least one failure, so a plain run
// reports failures without failing the caller.
func exitCode(failed int, strict bool) int {
	if strict && failed > 0 {
		return 1
	}
	return 0
}

func testAllBuildings(pool *pgxpool.Pool, cfg config.DBConfig, tableName string) []TestResult {
	// Get all row IDs
	query := fmt.Sprintf(`SELECT id FROM %s.%s ORDER BY id`, cfg.Schemas.Tabula, tableName)
	rows, err := pool.Query(context.Background(), query)
	if err != nil {
		log.Fatalf("Failed to query row IDs: %v", err)
	}
	defer rows.Close()

	var rowIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Printf("Failed to scan row ID: %v", err)
			continue
		}
		rowIDs = append(rowIDs, id)
	}

	// Run tests in parallel
	var results []TestResult
	resultsChan := make(chan TestResult, len(rowIDs))

	// Launch goroutines for parallel execution
	for _, rowID := range rowIDs {
		go func(id int) {
			result := runPipelineTest(pool, tableName, id)
			resultsChan <- result
		}(rowID)
	}

	// Collect results
	for range rowIDs {
		results = append(results, <-resultsChan)
	}
	close(resultsChan)

	sort.Slice(results, func(i, j int) bool { return results[i].RowID < results[j].RowID })
	return results
}

func runPipelineTest(pool *pgxpool.Pool, tableName string, rowID int) TestResult {
	result := TestResult{
		RowID: rowID,
	}

	// Load Tabula data and expected Q_h_nd from database
	tabulaData, buildingID, expectedQHND, err := loadTabulaDataFromDB(pool, tableName, rowID)
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("Failed to load data: %v", err)
		return result
	}

	result.BuildingID = buildingID
	result.ExpectedQHND = expectedQHND

	// // Save Tabula model as JSON before passing to pipeline
	// if err := saveTabulaModelAsJSON(tabulaData, buildingID, rowID); err != nil {
	// 	log.Printf("Warning: Failed to save Tabula model as JSON: %v", err)
	// 	// Continue execution even if JSON save fails
	// }

	// Create ignis instance and run pipeline
	logger := pipeline.NewLogger(log.New(os.Stdout, "", 0))
	p := pipeline.NewPipeline(tabulaData, logger)

	// Run the calculation pipeline
	calculatedQHND, err := p.Run()
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("pipeline error: %v", err)
		return result
	}
	result.CalculatedQHND = calculatedQHND

	// Calculate difference and percent error
	result.Difference = calculatedQHND - expectedQHND
	if expectedQHND != 0 {
		result.PercentError = math.Abs(result.Difference / expectedQHND * 100)
	} else if calculatedQHND != 0 {
		// If expected is 0 but calculated is not, still show error
		result.PercentError = 100.0
	}

	// Determine if test passed
	result.Passed = result.PercentError <= tolerancePercent

	return result
}

// loadTabulaDataFromDB loads Tabula building parameters from the database row
func loadTabulaDataFromDB(pool *pgxpool.Pool, tableName string, rowID int) (*models.TabulaBuildingParameters, string, float64, error) {
	// Query to get all column data for the row
	query := fmt.Sprintf(`SELECT * FROM %s.%s WHERE id = $1`, cfg.DB.Schemas.Tabula, tableName)

	rows, err := pool.Query(context.Background(), query, rowID)
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to query building data: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, "", 0, fmt.Errorf("no data found for row ID %d", rowID)
	}

	// Get column names and values
	fieldDescriptions := rows.FieldDescriptions()
	values, err := rows.Values()
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to read row values: %w", err)
	}

	// Build a map of column_name -> value
	// The column names from the database match the JSON tags in the tabula.go model
	dataMap := make(map[string]interface{})
	for i, fd := range fieldDescriptions {
		colName := string(fd.Name)
		if i < len(values) {
			dataMap[colName] = values[i]
		}
	}

	// Extract expected Q_h_nd (the output value to compare against)
	var expectedQHND float64
	if val, ok := dataMap["q_h_nd"]; ok && val != nil {
		switch v := val.(type) {
		case float64:
			expectedQHND = v
		case float32:
			expectedQHND = float64(v)
		}
	}

	// Extract building ID for reference
	buildingID := ""
	if val, ok := dataMap["Code_BuildingVariant"]; ok && val != nil {
		buildingID = fmt.Sprintf("%v", val)
	}

	// Initialize TabulaBuildingParameters with all nested structs
	tabulaData := initializeTabulaData()

	// Use reflection to populate the structs from the database map using JSON tags
	populateStructFromMapUsingReflection(tabulaData, dataMap)

	return tabulaData, buildingID, expectedQHND, nil
}

// initializeTabulaData creates a fully initialized TabulaBuildingParameters with all nested structs
func initializeTabulaData() *models.TabulaBuildingParameters {
	data := &models.TabulaBuildingParameters{
		BasicParameters: &models.BasicParameters{
			BuildingAppearance: &models.BuildingThematic{},
			Envelope:           &models.Envelope{},
		},
		AdvancedParameters: &models.AdvancedParameters{
			AirInfiltration:       &models.AirInfiltration{},
			ClimateConditions:     &models.ClimateConditions{},
			Uvalues:               &models.Uvalues{},
			Insulation:            &models.InsulationThicknesses{},
			SolarGains:            &models.SolarGains{},
			ThermalBridges:        &models.ThermalBridgeParameters{},
			HeatLosses:            &models.TransmissionHeatLoss{},
			ThermalResistances:    &models.ThermalResistances{},
			InsulationMeasures:    &models.InsulationPredefinedMeasures{},
			ActualInsulation:      &models.ActualInsulationThicknesses{},
			HeatTransfer:          &models.HeatTransferCoefficients{},
			PredefinedCodes:       &models.PredefinedCodes{},
			MeasureTypes:          &models.MeasureTypeCodes{},
			SolarTransmittance:    &models.SolarEnergyTransmittance{},
			MeasureFractions:      &models.MeasureAreaFractions{},
			AdditionalResistances: &models.AdditionalThermalResistance{},
		},
	}

	// Set constant default values that aren't in the database
	data.AdvancedParameters.PredefinedCodes.F_Corr_CeilingHeight = 1.0

	// NOTE: f_Measure values for all building components are now loaded from the database directly
	// They should NOT be overridden here

	return data
}

// populateStructFromMapUsingReflection walks through all nested structs and populates fields using JSON tags
func populateStructFromMapUsingReflection(target interface{}, dataMap map[string]interface{}) {
	populateStruct(reflect.ValueOf(target), dataMap)
}

func populateStruct(val reflect.Value, dataMap map[string]interface{}) {
	// Dereference pointers
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Skip unexported fields
		if !field.CanSet() {
			continue
		}

		// Get JSON tag
		jsonTag := fieldType.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			// No JSON tag, check if it's a nested struct
			if field.Kind() == reflect.Ptr && field.Type().Elem().Kind() == reflect.Struct {
				populateStruct(field, dataMap)
			} else if field.Kind() == reflect.Struct {
				populateStruct(field, dataMap)
			}
			continue
		}

		// Handle nested structs (pointers to structs)
		if field.Kind() == reflect.Ptr {
			if field.Type().Elem().Kind() == reflect.Struct {
				populateStruct(field, dataMap)
				continue
			}
		} else if field.Kind() == reflect.Struct {
			populateStruct(field, dataMap)
			continue
		}

		// Get value from map using JSON tag as key
		dbValue, ok := dataMap[jsonTag]
		if !ok || dbValue == nil {
			continue
		}

		// Set the field value
		setFieldValue(field, dbValue)
	}
}

func setFieldValue(field reflect.Value, value interface{}) {
	if !field.CanSet() {
		return
	}

	switch field.Kind() {
	case reflect.String:
		if v, ok := value.(string); ok {
			field.SetString(v)
		} else {
			field.SetString(fmt.Sprintf("%v", value))
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch v := value.(type) {
		case int:
			field.SetInt(int64(v))
		case int64:
			field.SetInt(v)
		case int32:
			field.SetInt(int64(v))
		case float64:
			field.SetInt(int64(v))
		case float32:
			field.SetInt(int64(v))
		}
	case reflect.Float32, reflect.Float64:
		switch v := value.(type) {
		case float64:
			field.SetFloat(v)
		case float32:
			field.SetFloat(float64(v))
		case int:
			field.SetFloat(float64(v))
		case int64:
			field.SetFloat(float64(v))
		}
	}
}

// saveTabulaModelAsJSON saves the TabulaBuildingParameters as a JSON file
func saveTabulaModelAsJSON(tabulaData *models.TabulaBuildingParameters, buildingID string, rowID int) error {
	// Create output directory if it doesn't exist
	outputDir := "data/tabula_models"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create filename based on building ID and row ID
	sanitizedBuildingID := strings.ReplaceAll(buildingID, "/", "_")
	sanitizedBuildingID = strings.ReplaceAll(sanitizedBuildingID, ".", "_")
	filename := fmt.Sprintf("%s_row_%d.json", sanitizedBuildingID, rowID)
	filePath := filepath.Join(outputDir, filename)

	// Marshal to JSON with indentation for readability
	jsonData, err := json.MarshalIndent(tabulaData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tabula data to JSON: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	return nil
}
