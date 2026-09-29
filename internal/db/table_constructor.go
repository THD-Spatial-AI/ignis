package importer

import (
	"context"
	"fmt"
	"github.com/thd-spatial-ai/ignis/internal/config"
	"github.com/thd-spatial-ai/ignis/internal/utils"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xuri/excelize/v2"
)

// lfsPointerPrefix is how a Git LFS pointer file starts, in place of the
// real binary content, when `git lfs pull` was never run after clone.
const lfsPointerPrefix = "version https://git-lfs.github.com"

// checkXLSXHeader reads the first few bytes of path and fails fast with an
// actionable message if it isn't a real .xlsx (a zip file, signature "PK").
// Without this, excelize.OpenFile's own error ("zip: not a valid zip file")
// gives no hint that the real cause is an unresolved Git LFS pointer.
func checkXLSXHeader(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("opening workbook %s: %w", path, err)
	}
	defer f.Close()

	head := make([]byte, len(lfsPointerPrefix))
	n, err := io.ReadFull(f, head)
	if err != nil && err != io.ErrUnexpectedEOF {
		return fmt.Errorf("reading workbook %s: %w", path, err)
	}
	head = head[:n]

	if len(head) >= 2 && head[0] == 'P' && head[1] == 'K' {
		return nil
	}
	if strings.HasPrefix(string(head), lfsPointerPrefix) {
		return fmt.Errorf("%s is a Git LFS pointer, not the workbook itself — run `git lfs install && git lfs pull`", path)
	}
	return fmt.Errorf("%s is not a valid .xlsx file (missing zip signature)", path)
}

type HeaderInfo struct {
	CellIndex           int
	CellDataType        string
	CellDataValidations []string
}

type TableConstructor struct {
	conn          *pgxpool.Pool
	tx            pgx.Tx // open only for the duration of Run
	cfg           *config.Config
	xlsxFile      *excelize.File
	headers       map[string]*HeaderInfo
	countryCodes  []string // Fixed field name
	countryHelper *utils.TabulaCountryHelper
}

func NewTableConstructor(conn *pgxpool.Pool, cfg *config.Config) *TableConstructor {
	return &TableConstructor{
		conn:          conn,
		cfg:           cfg,
		headers:       make(map[string]*HeaderInfo),
		countryHelper: utils.NewTabulaCountryHelper(),
	}
}

// IsSeeded reports whether schema holds any table. Run writes in a single
// transaction, so a schema with tables is one a Run completed.
func IsSeeded(ctx context.Context, pool *pgxpool.Pool, schema string) (bool, error) {
	var seeded bool
	err := pool.QueryRow(ctx,
		`SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = $1)`,
		schema,
	).Scan(&seeded)
	if err != nil {
		return false, fmt.Errorf("checking for tables in schema %s: %w", schema, err)
	}
	return seeded, nil
}

// Run drops and recreates every country table from the workbook in one
// transaction: on any error nothing is committed and the previous tables stay.
func (tc *TableConstructor) Run() error {
	defer tc.close()

	if err := tc.loadWorkbook(); err != nil {
		return err
	}

	rows, err := tc.xlsxFile.GetRows("Calc.Set.Building")
	if err != nil {
		return err
	}

	tc.extractHeaders(rows)
	tc.extractCountryCodes(rows)

	ctx := context.Background()
	tc.tx, err = tc.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("starting seed transaction: %w", err)
	}
	defer tc.tx.Rollback(ctx) // no-op once committed

	if err := tc.createTables(); err != nil {
		return err
	}
	if err := tc.insertData(rows); err != nil {
		return err
	}
	if err := tc.updateDropdownColumns(); err != nil {
		return err
	}
	if err := tc.tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing seed transaction: %w", err)
	}

	utils.Info.Println("Table construction completed successfully")
	return nil
}

func (tc *TableConstructor) loadWorkbook() error {
	if err := checkXLSXHeader(tc.cfg.Data.ExcelFile); err != nil {
		return err
	}
	file, err := excelize.OpenFile(tc.cfg.Data.ExcelFile)
	tc.xlsxFile = file
	return err
}

func (tc *TableConstructor) extractHeaders(rows [][]string) {
	if len(rows) < 6 {
		return
	}

	headerRow, dataTypeRow := rows[0], rows[5]
	dropdowns := tc.getDropdownValues(rows, headerRow)

	for i, header := range headerRow {
		if header == "" {
			continue
		}

		// Include f_Measure_* columns (area fractions for refurbishment measures)
		// These are important for calculating actual U-values after refurbishment

		dataType := "VarChar"
		if i < len(dataTypeRow) && dataTypeRow[i] != "" {
			dataType = dataTypeRow[i]
		}

		// Add regular column
		tc.headers[header] = &HeaderInfo{
			CellIndex:    i,
			CellDataType: dataType,
		}

		// Add dropdown column if exists
		if vals, exists := dropdowns[header]; exists && len(vals) > 0 {
			tc.headers[header+"_val_data"] = &HeaderInfo{
				CellIndex:           i,
				CellDataType:        "List",
				CellDataValidations: vals,
			}
		}
	}
}

func (tc *TableConstructor) getDropdownValues(rows [][]string, headerRow []string) map[string][]string {
	known := map[string][]string{
		"Code_ComplexRoof":  {"yes", "no"},
		"Code_Country":      {"AT", "BE", "BG", "CY", "CZ", "DE", "DK", "ES", "FR", "GB", "GR", "HU", "IE", "IT", "NL", "NO", "PL", "RS", "SE", "SI"},
		"Code_BuildingType": {"SFH", "TH", "MFH", "AB"},
	}

	// Sample 20 rows for Code_ columns
	columnVals := make(map[int]map[string]bool)
	for i, header := range headerRow {
		if strings.HasPrefix(header, "Code_") && known[header] == nil {
			columnVals[i] = make(map[string]bool)
		}
	}

	for i := 12; i < len(rows) && i < 32; i++ { // Sample rows 13-32
		for colIdx, vals := range columnVals {
			if colIdx < len(rows[i]) {
				val := strings.TrimSpace(rows[i][colIdx])
				if val != "" && len(val) <= 20 {
					vals[val] = true
				}
			}
		}
	}

	// Convert to arrays
	for colIdx, vals := range columnVals {
		if len(vals) > 1 && len(vals) <= 15 {
			var arr []string
			for val := range vals {
				arr = append(arr, val)
			}
			sort.Strings(arr)
			known[headerRow[colIdx]] = arr
		}
	}

	return known
}

func (tc *TableConstructor) extractCountryCodes(rows [][]string) {
	codeCol := -1
	for i, header := range rows[0] {
		if strings.ToLower(header) == "code_country" {
			codeCol = i
			break
		}
	}

	codes := make(map[string]bool)
	for i := 12; i < len(rows) && i < 2159; i++ {
		if codeCol < len(rows[i]) {
			code := strings.TrimSpace(rows[i][codeCol])
			if len(code) == 2 && tc.countryHelper.IsValidCountryCode(code) {
				codes[code] = true
			}
		}
	}

	for code := range codes {
		tc.countryCodes = append(tc.countryCodes, code)
	}
	sort.Strings(tc.countryCodes)
}

func (tc *TableConstructor) createTables() error {
	schema := tc.cfg.DB.Schemas.Tabula
	if _, err := tc.tx.Exec(context.Background(), fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema)); err != nil {
		return fmt.Errorf("creating schema %s: %w", schema, err)
	}

	var cols []string
	for header, info := range tc.headers {
		pgType := map[string]string{"VarChar": "VARCHAR", "Text": "TEXT", "Date": "DATE", "Real": "REAL", "Integer": "INTEGER", "List": "VARCHAR[]"}[info.CellDataType]
		if pgType == "" {
			pgType = "VARCHAR"
		}
		cols = append(cols, fmt.Sprintf(`"%s" %s`, header, pgType))
	}

	for _, code := range tc.countryCodes {
		table := fmt.Sprintf("%s.%s", tc.cfg.DB.Schemas.Tabula, tc.countryHelper.CodeToCountry(code))
		if _, err := tc.tx.Exec(context.Background(), fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)); err != nil {
			return fmt.Errorf("dropping table %s: %w", table, err)
		}
		if _, err := tc.tx.Exec(context.Background(), fmt.Sprintf("CREATE TABLE %s (id SERIAL PRIMARY KEY, %s)", table, strings.Join(cols, ", "))); err != nil {
			return fmt.Errorf("creating table %s: %w", table, err)
		}
	}
	return nil
}

func (tc *TableConstructor) insertData(rows [][]string) error {
	headerMap := make(map[string]int)
	for i, header := range rows[0] {
		headerMap[header] = i
	}

	codeCol := headerMap["Code_Country"]

	// Limit to row 2159 - rows beyond contain incomplete datasets
	for i := 12; i < len(rows) && i < 2159; i++ {
		if codeCol >= len(rows[i]) {
			continue
		}

		code := strings.TrimSpace(rows[i][codeCol])
		if !tc.countryHelper.IsValidCountryCode(code) {
			continue
		}

		if err := tc.insertRow(code, headerMap, rows[i]); err != nil {
			return fmt.Errorf("workbook row %d: %w", i+1, err)
		}
	}
	return nil
}

func (tc *TableConstructor) insertRow(countryCode string, headerMap map[string]int, dataRow []string) error {
	table := fmt.Sprintf("%s.%s", tc.cfg.DB.Schemas.Tabula, tc.countryHelper.CodeToCountry(countryCode))

	var cols, placeholders []string
	var vals []interface{}
	paramIdx := 1

	for header, info := range tc.headers {
		if strings.HasSuffix(header, "_val_data") {
			continue // Skip dropdown columns in initial insert
		}

		cols = append(cols, fmt.Sprintf(`"%s"`, header))
		placeholders = append(placeholders, fmt.Sprintf("$%d", paramIdx))

		var val interface{}
		if colIdx, exists := headerMap[header]; exists && colIdx < len(dataRow) {
			cellVal := strings.TrimSpace(dataRow[colIdx])
			if cellVal != "" {
				switch info.CellDataType {
				case "Real":
					if v, err := strconv.ParseFloat(cellVal, 64); err == nil {
						val = v
					}
				case "Integer":
					if v, err := strconv.Atoi(cellVal); err == nil {
						val = v
					}
				case "Date":
					// Skip invalid dates like "1900-01-00"
					if !strings.Contains(cellVal, "-00") {
						val = cellVal
					}
				default:
					val = cellVal
				}
			}
		}
		vals = append(vals, val)
		paramIdx++
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(cols, ", "), strings.Join(placeholders, ", "))
	if _, err := tc.tx.Exec(context.Background(), query, vals...); err != nil {
		return fmt.Errorf("inserting into %s: %w", table, err)
	}
	return nil
}

func (tc *TableConstructor) updateDropdownColumns() error {
	for _, code := range tc.countryCodes {
		table := fmt.Sprintf("%s.%s", tc.cfg.DB.Schemas.Tabula, tc.countryHelper.CodeToCountry(code))

		var updates []string
		for header, info := range tc.headers {
			if strings.HasSuffix(header, "_val_data") && len(info.CellDataValidations) > 0 {
				var quoted []string
				for _, val := range info.CellDataValidations {
					quoted = append(quoted, fmt.Sprintf("'%s'", strings.ReplaceAll(val, "'", "''")))
				}
				updates = append(updates, fmt.Sprintf(`"%s" = ARRAY[%s]`, header, strings.Join(quoted, ",")))
			}
		}

		if len(updates) > 0 {
			if _, err := tc.tx.Exec(context.Background(), fmt.Sprintf("UPDATE %s SET %s", table, strings.Join(updates, ", "))); err != nil {
				return fmt.Errorf("setting dropdown columns on %s: %w", table, err)
			}
		}
	}
	return nil
}

func (tc *TableConstructor) close() error {
	if tc.xlsxFile != nil {
		return tc.xlsxFile.Close()
	}
	return nil
}
