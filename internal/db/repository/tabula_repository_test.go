package repository

import "testing"

func TestTabulaRepository_qualifyTable(t *testing.T) {
	cases := []struct {
		schema, table, want string
	}{
		{"tabula", "germany", `"tabula"."germany"`},
		{"", "germany", `"germany"`},
	}
	for _, tc := range cases {
		r := NewTabulaRepository(nil, tc.schema)
		if got := r.qualifyTable(tc.table); got != tc.want {
			t.Errorf("qualifyTable(schema=%q, table=%q) = %q, want %q", tc.schema, tc.table, got, tc.want)
		}
	}
}

func TestNewTabulaRepository(t *testing.T) {
	r := NewTabulaRepository(nil, "tabula")
	if r.schema != "tabula" {
		t.Errorf("schema = %q, want %q", r.schema, "tabula")
	}
}

func TestErrVariantNotFound(t *testing.T) {
	if ErrVariantNotFound.Error() != "tabula variant not found" {
		t.Errorf("ErrVariantNotFound.Error() = %q", ErrVariantNotFound.Error())
	}
}
