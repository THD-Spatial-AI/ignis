package importer

import (
	"fmt"
	"sort"
	"testing"

	"github.com/thd-spatial-ai/ignis/internal/config"
)

// buildFixtureRows constructs a minimal TABULA-shaped [][]string: a header
// row, a data-type row at index 5, and 20 data rows (indices 12-31) - the
// same layout extractHeaders/getDropdownValues/extractCountryCodes expect
// from the real "Calc.Set.Building" sheet.
func buildFixtureRows() [][]string {
	rows := make([][]string, 32)
	for i := range rows {
		rows[i] = []string{}
	}
	rows[0] = []string{"Code_BuildingVariant", "Code_Country", "Code_ComplexRoof", "Code_Custom", "A_Roof_1", "HeatingDays"}
	rows[5] = []string{"VarChar", "VarChar", "VarChar", "VarChar", "Real", "Integer"}

	countries := []string{"DE", "AT", "FR"}
	customVals := []string{"X", "Y"}
	for i := 12; i < 32; i++ {
		rows[i] = []string{
			fmt.Sprintf("DE.N.SFH.%02d", i),
			countries[i%len(countries)],
			"yes",
			customVals[i%len(customVals)],
			"45.5",
			"180",
		}
	}
	return rows
}

func newTestTableConstructor() *TableConstructor {
	return NewTableConstructor(nil, &config.Config{
		DB: &config.DBConfig{Schemas: &config.Schemas{Tabula: "tabula"}},
	})
}

func TestExtractHeaders_populatesHeaderInfo(t *testing.T) {
	tc := newTestTableConstructor()
	tc.extractHeaders(buildFixtureRows())

	info, ok := tc.headers["A_Roof_1"]
	if !ok {
		t.Fatal(`expected "A_Roof_1" in headers`)
	}
	if info.CellIndex != 4 || info.CellDataType != "Real" {
		t.Errorf("A_Roof_1 header = %+v, want CellIndex=4 CellDataType=Real", info)
	}

	info, ok = tc.headers["HeatingDays"]
	if !ok || info.CellDataType != "Integer" {
		t.Errorf("HeatingDays header = %+v, want CellDataType=Integer", info)
	}
}

func TestExtractHeaders_emptyHeaderNameSkipped(t *testing.T) {
	tc := newTestTableConstructor()
	rows := buildFixtureRows()
	rows[0] = append(rows[0], "")
	tc.extractHeaders(rows)

	if _, ok := tc.headers[""]; ok {
		t.Error("expected empty-string header to be skipped")
	}
}

func TestExtractHeaders_missingOrEmptyDataTypeDefaultsToVarChar(t *testing.T) {
	tc := newTestTableConstructor()
	rows := buildFixtureRows()
	rows[5][0] = "" // Code_BuildingVariant's data type cell is blank
	tc.extractHeaders(rows)

	if tc.headers["Code_BuildingVariant"].CellDataType != "VarChar" {
		t.Errorf("CellDataType = %q, want VarChar fallback", tc.headers["Code_BuildingVariant"].CellDataType)
	}
}

func TestExtractHeaders_shortRowsIsNoop(t *testing.T) {
	tc := newTestTableConstructor()
	tc.extractHeaders([][]string{{"only one row"}})
	if len(tc.headers) != 0 {
		t.Errorf("expected no headers extracted from <6 rows, got %d", len(tc.headers))
	}
}

func TestExtractHeaders_addsDropdownColumn(t *testing.T) {
	tc := newTestTableConstructor()
	tc.extractHeaders(buildFixtureRows())

	info, ok := tc.headers["Code_Custom_val_data"]
	if !ok {
		t.Fatal(`expected "Code_Custom_val_data" dropdown column to be added`)
	}
	if info.CellDataType != "List" {
		t.Errorf("CellDataType = %q, want List", info.CellDataType)
	}
	want := []string{"X", "Y"}
	sort.Strings(info.CellDataValidations)
	if len(info.CellDataValidations) != len(want) {
		t.Fatalf("CellDataValidations = %v, want %v", info.CellDataValidations, want)
	}
	for i, v := range want {
		if info.CellDataValidations[i] != v {
			t.Errorf("CellDataValidations = %v, want %v", info.CellDataValidations, want)
		}
	}
}

func TestGetDropdownValues_knownColumnsAreNotResampled(t *testing.T) {
	tc := newTestTableConstructor()
	rows := buildFixtureRows()
	dropdowns := tc.getDropdownValues(rows, rows[0])

	want := []string{"yes", "no"}
	got := dropdowns["Code_ComplexRoof"]
	if len(got) != len(want) {
		t.Fatalf("Code_ComplexRoof = %v, want the hardcoded %v (not resampled)", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Code_ComplexRoof = %v, want %v", got, want)
		}
	}
}

func TestGetDropdownValues_tooManyDistinctValuesIsExcluded(t *testing.T) {
	tc := newTestTableConstructor()
	rows := make([][]string, 32)
	for i := range rows {
		rows[i] = []string{}
	}
	rows[0] = []string{"Code_HighCardinality"}
	rows[5] = []string{"VarChar"}
	for i := 12; i < 32; i++ {
		// 20 rows, each with a unique value -> exceeds the len<=15 cutoff.
		rows[i] = []string{fmt.Sprintf("val%d", i)}
	}

	dropdowns := tc.getDropdownValues(rows, rows[0])
	if _, exists := dropdowns["Code_HighCardinality"]; exists {
		t.Error("expected column with >15 distinct sampled values to be excluded")
	}
}

func TestExtractCountryCodes_collectsValidCodesSorted(t *testing.T) {
	tc := newTestTableConstructor()
	tc.extractCountryCodes(buildFixtureRows())

	want := []string{"AT", "DE", "FR"}
	if len(tc.countryCodes) != len(want) {
		t.Fatalf("countryCodes = %v, want %v", tc.countryCodes, want)
	}
	for i := range want {
		if tc.countryCodes[i] != want[i] {
			t.Errorf("countryCodes = %v, want %v", tc.countryCodes, want)
		}
	}
}

func TestExtractCountryCodes_invalidCodesIgnored(t *testing.T) {
	tc := newTestTableConstructor()
	rows := buildFixtureRows()
	for i := 12; i < 32; i++ {
		rows[i][1] = "ZZ" // not a real ISO2 code
	}
	tc.extractCountryCodes(rows)

	if len(tc.countryCodes) != 0 {
		t.Errorf("countryCodes = %v, want none for an all-invalid country column", tc.countryCodes)
	}
}

func TestClose_nilXlsxFileIsNoop(t *testing.T) {
	tc := newTestTableConstructor()
	if err := tc.close(); err != nil {
		t.Errorf("close() with nil xlsxFile: unexpected error %v", err)
	}
}

func TestLoadWorkbook_missingFile_returnsError(t *testing.T) {
	tc := NewTableConstructor(nil, &config.Config{
		Data: &config.DataPaths{ExcelFile: "/nonexistent/path/workbook.xlsx"},
	})
	if err := tc.loadWorkbook(); err == nil {
		t.Error("expected error opening a nonexistent workbook")
	}
}
