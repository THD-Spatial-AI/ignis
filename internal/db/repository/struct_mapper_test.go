package repository

import (
	"reflect"
	"testing"
)

func TestInitializeTabulaData(t *testing.T) {
	data := initializeTabulaData()

	if data.BasicParameters == nil || data.BasicParameters.BuildingAppearance == nil || data.BasicParameters.Envelope == nil {
		t.Fatal("expected BasicParameters and its nested pointers to be non-nil")
	}
	if data.AdvancedParameters == nil {
		t.Fatal("expected AdvancedParameters to be non-nil")
	}
	if data.AdvancedParameters.PredefinedCodes.F_Corr_CeilingHeight != 1.0 {
		t.Errorf("F_Corr_CeilingHeight = %v, want 1.0 default", data.AdvancedParameters.PredefinedCodes.F_Corr_CeilingHeight)
	}
}

func TestPopulateStructFromMap_setsFieldsByJSONTag(t *testing.T) {
	data := initializeTabulaData()
	dataMap := map[string]interface{}{
		"Code_BuildingVariant": "DE.N.SFH.01.Gen",
		"n_Storey":             int32(3),
		"A_Roof_1":             float32(45.5),
		"HeatingDays":          int64(180),
		"Theta_e":              -12,
		"n_air_infiltration":   float64(0.5),
	}

	populateStructFromMap(data, dataMap)

	if data.BasicParameters.BuildingAppearance.Code_BuildingVariant != "DE.N.SFH.01.Gen" {
		t.Errorf("Code_BuildingVariant = %q, want %q", data.BasicParameters.BuildingAppearance.Code_BuildingVariant, "DE.N.SFH.01.Gen")
	}
	if data.BasicParameters.BuildingAppearance.N_Storey != 3 {
		t.Errorf("N_Storey = %d, want 3 (from int32 source)", data.BasicParameters.BuildingAppearance.N_Storey)
	}
	if data.BasicParameters.Envelope.A_Roof_1 != float64(float32(45.5)) {
		t.Errorf("A_Roof_1 = %v, want %v (from float32 source)", data.BasicParameters.Envelope.A_Roof_1, float64(float32(45.5)))
	}
	if data.AdvancedParameters.ClimateConditions.HeatingDays != 180 {
		t.Errorf("HeatingDays = %d, want 180 (from int64 source)", data.AdvancedParameters.ClimateConditions.HeatingDays)
	}
	if data.AdvancedParameters.ClimateConditions.Theta_e != -12 {
		t.Errorf("Theta_e = %v, want -12 (from int source, stringified fallback not expected)", data.AdvancedParameters.ClimateConditions.Theta_e)
	}
	if data.AdvancedParameters.AirInfiltration.N_air_infiltration != 0.5 {
		t.Errorf("N_air_infiltration = %v, want 0.5", data.AdvancedParameters.AirInfiltration.N_air_infiltration)
	}
}

func TestPopulateStructFromMap_missingAndNilKeysAreSkipped(t *testing.T) {
	data := initializeTabulaData()
	dataMap := map[string]interface{}{
		"Code_BuildingVariant": nil, // present but nil - must be skipped, not zeroed weirdly
		// "n_Storey" intentionally absent
	}

	populateStructFromMap(data, dataMap)

	if data.BasicParameters.BuildingAppearance.Code_BuildingVariant != "" {
		t.Errorf("Code_BuildingVariant = %q, want zero value when map value is nil", data.BasicParameters.BuildingAppearance.Code_BuildingVariant)
	}
	if data.BasicParameters.BuildingAppearance.N_Storey != 0 {
		t.Errorf("N_Storey = %d, want zero value when key is absent", data.BasicParameters.BuildingAppearance.N_Storey)
	}
}

func TestPopulateStructFromMap_stringFieldNonStringFallback(t *testing.T) {
	data := initializeTabulaData()
	// Code_BuildingVariant is a string field; feeding it a non-string value
	// exercises setFieldValue's fmt.Sprintf fallback branch.
	dataMap := map[string]interface{}{"Code_BuildingVariant": 12345}

	populateStructFromMap(data, dataMap)

	if data.BasicParameters.BuildingAppearance.Code_BuildingVariant != "12345" {
		t.Errorf("Code_BuildingVariant = %q, want %q (stringified fallback)", data.BasicParameters.BuildingAppearance.Code_BuildingVariant, "12345")
	}
}

func TestPopulateStructFromMap_nilPointerTarget(t *testing.T) {
	// populateStruct must not panic when reflect.ValueOf(target) is a nil pointer.
	var data *struct{ Foo string }
	populateStructFromMap(data, map[string]interface{}{"Foo": "bar"})
}

func TestPopulateStructFromMap_nonStructTarget(t *testing.T) {
	// populateStruct must be a no-op (not panic) for a non-struct, non-pointer target.
	x := 5
	populateStructFromMap(&x, map[string]interface{}{"x": 10})
	if x != 5 {
		t.Errorf("x = %d, want unchanged 5", x)
	}
}

func TestPopulateStructFromMap_unexportedFieldSkipped(t *testing.T) {
	type withUnexported struct {
		unexported string //nolint:unused // deliberately unexported to hit the CanSet()==false skip branch
		Exported   string `json:"exported"`
	}
	target := &withUnexported{}
	populateStructFromMap(target, map[string]interface{}{
		"exported": "should be set",
	})
	if target.Exported != "should be set" {
		t.Errorf("Exported = %q, want %q", target.Exported, "should be set")
	}
	if target.unexported != "" {
		t.Errorf("unexported = %q, want zero value (field.CanSet() must be false)", target.unexported)
	}
}

func TestPopulateStructFromMap_valueStructWithTag(t *testing.T) {
	// Covers the "has a JSON tag, but is itself a nested value struct" branch
	// in populateStruct - the tag is present but ignored in favour of recursing.
	type Inner struct {
		Name string `json:"name"`
	}
	type Outer struct {
		Inner Inner `json:"inner"`
	}
	target := &Outer{}
	populateStructFromMap(target, map[string]interface{}{"name": "nested"})
	if target.Inner.Name != "nested" {
		t.Errorf("Inner.Name = %q, want %q", target.Inner.Name, "nested")
	}
}

func TestSetFieldValue_intFieldFromPlainInt(t *testing.T) {
	type s struct {
		N int `json:"n"`
	}
	target := &s{}
	populateStructFromMap(target, map[string]interface{}{"n": 7})
	if target.N != 7 {
		t.Errorf("N = %d, want 7", target.N)
	}
}

func TestSetFieldValue_floatFieldFromInt64(t *testing.T) {
	type s struct {
		F float64 `json:"f"`
	}
	target := &s{}
	populateStructFromMap(target, map[string]interface{}{"f": int64(9)})
	if target.F != 9 {
		t.Errorf("F = %v, want 9", target.F)
	}
}

func TestPopulateStructFromMap_valueStructWithoutTag(t *testing.T) {
	// Covers the "no JSON tag, non-pointer nested struct" branch in populateStruct.
	type Inner struct {
		Name string `json:"name"`
	}
	type Outer struct {
		Inner Inner // no json tag, embedded-by-value struct
	}
	target := &Outer{}
	populateStructFromMap(target, map[string]interface{}{"name": "nested"})
	if target.Inner.Name != "nested" {
		t.Errorf("Inner.Name = %q, want %q", target.Inner.Name, "nested")
	}
}

func TestNormalizeValue(t *testing.T) {
	if got := normalizeValue([]byte("hello")); got != "hello" {
		t.Errorf("normalizeValue([]byte) = %v, want %q", got, "hello")
	}
	if got := normalizeValue(42); got != 42 {
		t.Errorf("normalizeValue(42) = %v, want 42", got)
	}
	if got := normalizeValue(nil); got != nil {
		t.Errorf("normalizeValue(nil) = %v, want nil", got)
	}
}

func TestToFloat64(t *testing.T) {
	cases := []struct {
		in   interface{}
		want float64
	}{
		{float64(1.5), 1.5},
		{float32(2.5), 2.5},
		{int(3), 3},
		{int32(4), 4},
		{int64(5), 5},
		{"6.5", 6.5},
		{"not-a-number", 0},
		{nil, 0},
	}
	for _, tc := range cases {
		if got := toFloat64(tc.in); got != tc.want {
			t.Errorf("toFloat64(%v) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestSetFieldValue_unsettableFieldIsNoop(t *testing.T) {
	type s struct{ unexported string }
	v := reflect.ValueOf(&s{}).Elem().Field(0)
	setFieldValue(v, "value") // must not panic despite CanSet() == false
}
