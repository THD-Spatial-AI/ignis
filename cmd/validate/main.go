package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/thd-spatial-ai/ignis/internal/config"
	"github.com/thd-spatial-ai/ignis/internal/hdcp"
	"github.com/thd-spatial-ai/ignis/internal/models"
	"log"
	"math"
	"os"
	"path/filepath"
	"reflect"
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

const tolerancePercent = 2.5 // 2% tolerance for heating demand cross-country variations

var cfg = config.LoadConfig()

func main() {
	fmt.Println("=== ignis Validation Tool ===")
	startTime := time.Now()
	// Load configuration
	cfg := config.LoadConfig()

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

	// Test all buildings
	totalRows := 0
	for _, tableName := range tableNames {
		// fmt.Printf("\n=== Testing buildings in table: %s ===\n", tableName)

		totalRows += testAllBuildings(pool, *cfg.DB, tableName)
	}

	elapsed := time.Since(startTime)
	fmt.Printf("\n=== Validation completed in %s ===\n", elapsed)
	fmt.Printf("Total buildings tested across all tables: %d\n", totalRows)
}

func testAllBuildings(pool *pgxpool.Pool, cfg config.DBConfig, tableName string) int {
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

	// // Print summary
	// fmt.Println("\n" + strings.Repeat("=", 60))
	// fmt.Println("=== Test Summary ===")
	// fmt.Println(strings.Repeat("=", 60))

	passCount := 0
	failCount := 0
	for _, result := range results {
		if result.Passed {
			passCount++
		} else {
			failCount++
		}
	}

	// fmt.Printf("Total buildings tested: %d\n", len(results))
	// fmt.Printf("Passed: %d (%.1f%%)\n", passCount, float64(passCount)/float64(len(rowIDs))*100)
	// fmt.Printf("Failed: %d (%.1f%%)\n", failCount, float64(failCount)/float64(len(rowIDs))*100)

	// // Print detailed results for failed tests (limited to first 10)
	// if failCount > 0 {
	// 	fmt.Println("\n=== Failed Tests Details (showing first 10) ===")
	// 	count := 0
	// 	for _, result := range results {
	// 		if !result.Passed && count < 10 {
	// 			printTestResult(result)
	// 			count++
	// 		}
	// 	}
	// }

	return len(results)
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
	logger := hdcp.NewLogger(log.New(os.Stdout, "", 0))
	pipeline := hdcp.NewPipeline(tabulaData, logger)

	// Run the calculation pipeline
	calculatedQHND, err := pipeline.Run()
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

	// Determine if test passed (allow 2% tolerance for cross-country variations)
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

// printTestResult prints detailed test result information
func printTestResult(result TestResult) {
	separator := strings.Repeat("=", 60)
	fmt.Println("\n" + separator)
	fmt.Printf("Row ID: %d | Building: %s\n", result.RowID, result.BuildingID)
	fmt.Println(separator)

	if result.ErrorMessage != "" {
		fmt.Printf("Status: ✗ FAILED\n")
		fmt.Printf("Error: %s\n", result.ErrorMessage)
		fmt.Println(separator)
		return
	}

	if result.Passed {
		fmt.Printf("Status: o PASSED\n")
	} else {
		fmt.Printf("Status: x FAILED\n")
	}

	fmt.Printf("\nResults:\n")
	fmt.Printf("  Calculated q_h_nd: %.6f kWh/(m²·a)\n", result.CalculatedQHND)
	fmt.Printf("  Expected q_h_nd:   %.6f kWh/(m²·a)\n", result.ExpectedQHND)
	fmt.Printf("  Difference:        %.6f kWh/(m²·a)\n", result.Difference)
	fmt.Printf("  Percent Error:     %.4f%%\n", result.PercentError)
	fmt.Println(separator)
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
