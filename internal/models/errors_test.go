package models

import "testing"

func TestValidationError(t *testing.T) {
	err := NewValidationError("A_Roof_1", "must be positive")
	want := "validation error on field 'A_Roof_1': must be positive"
	if err.Error() != want {
		t.Errorf("Error() = %q, want %q", err.Error(), want)
	}
	if err.Field != "A_Roof_1" || err.Message != "must be positive" {
		t.Errorf("unexpected fields: %+v", err)
	}
}

func TestCalculationError(t *testing.T) {
	err := NewCalculationError("Level 7", "division by zero")
	want := "calculation error at level 'Level 7': division by zero"
	if err.Error() != want {
		t.Errorf("Error() = %q, want %q", err.Error(), want)
	}
	if err.Level != "Level 7" || err.Message != "division by zero" {
		t.Errorf("unexpected fields: %+v", err)
	}
}

func TestDatabaseError(t *testing.T) {
	err := NewDatabaseError("query", "connection refused")
	want := "database error during 'query': connection refused"
	if err.Error() != want {
		t.Errorf("Error() = %q, want %q", err.Error(), want)
	}
	if err.Operation != "query" || err.Message != "connection refused" {
		t.Errorf("unexpected fields: %+v", err)
	}
}

func TestNotFoundError(t *testing.T) {
	err := NewNotFoundError("variant", "DE.N.SFH.01")
	want := "variant with ID 'DE.N.SFH.01' not found"
	if err.Error() != want {
		t.Errorf("Error() = %q, want %q", err.Error(), want)
	}
	if err.Resource != "variant" || err.ID != "DE.N.SFH.01" {
		t.Errorf("unexpected fields: %+v", err)
	}
}
