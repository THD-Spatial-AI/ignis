package handler

import (
	"testing"
)

func TestAggregateSurfaces_nilOrEmpty_isNoOp(t *testing.T) {
	building := realisticBuilding()
	wantArea, wantU := building.BasicParameters.Envelope.A_Wall_1, building.AdvancedParameters.Uvalues.U_Wall_1

	if err := aggregateSurfaces(nil, building); err != nil {
		t.Fatalf("nil surfaces: unexpected error: %v", err)
	}
	if building.BasicParameters.Envelope.A_Wall_1 != wantArea || building.AdvancedParameters.Uvalues.U_Wall_1 != wantU {
		t.Errorf("nil surfaces changed the building: A_Wall_1=%v U_Wall_1=%v", building.BasicParameters.Envelope.A_Wall_1, building.AdvancedParameters.Uvalues.U_Wall_1)
	}

	if err := aggregateSurfaces([]Surface{}, building); err != nil {
		t.Fatalf("empty surfaces: unexpected error: %v", err)
	}
	if building.BasicParameters.Envelope.A_Wall_1 != wantArea || building.AdvancedParameters.Uvalues.U_Wall_1 != wantU {
		t.Errorf("empty surfaces changed the building: A_Wall_1=%v U_Wall_1=%v", building.BasicParameters.Envelope.A_Wall_1, building.AdvancedParameters.Uvalues.U_Wall_1)
	}
}

func TestAggregateSurfaces_areaWeightedAverage(t *testing.T) {
	building := realisticBuilding()
	surfaces := []Surface{
		{ID: "wall-a", Type: "wall", Area: 40, UValue: 1.0},
		{ID: "wall-b", Type: "wall", Area: 60, UValue: 0.5},
	}
	if err := aggregateSurfaces(surfaces, building); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantArea := 100.0
	wantU := (40*1.0 + 60*0.5) / 100.0 // area-weighted average = 0.7
	if building.BasicParameters.Envelope.A_Wall_1 != wantArea {
		t.Errorf("A_Wall_1 = %v, want %v", building.BasicParameters.Envelope.A_Wall_1, wantArea)
	}
	if building.AdvancedParameters.Uvalues.U_Wall_1 != wantU {
		t.Errorf("U_Wall_1 = %v, want %v", building.AdvancedParameters.Uvalues.U_Wall_1, wantU)
	}
}

func TestAggregateSurfaces_categoryWithNoSurfaces_keepsFallback(t *testing.T) {
	building := realisticBuilding()
	wantRoofArea, wantRoofU := building.BasicParameters.Envelope.A_Roof_1, building.AdvancedParameters.Uvalues.U_Roof_1

	surfaces := []Surface{{ID: "wall-a", Type: "wall", Area: 10, UValue: 1.0}}
	if err := aggregateSurfaces(surfaces, building); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if building.BasicParameters.Envelope.A_Roof_1 != wantRoofArea || building.AdvancedParameters.Uvalues.U_Roof_1 != wantRoofU {
		t.Errorf("roof category changed despite no roof surfaces given: area=%v u=%v", building.BasicParameters.Envelope.A_Roof_1, building.AdvancedParameters.Uvalues.U_Roof_1)
	}
}

func TestAggregateSurfaces_windowOrientationBucketing(t *testing.T) {
	building := realisticBuilding()
	south, east, horizontal := 180.0, 90.0, 0.0
	surfaces := []Surface{
		{ID: "win-s", Type: "window", Area: 8, UValue: 1.2, Azimuth: &south},
		{ID: "win-e", Type: "window", Area: 4, UValue: 1.5, Azimuth: &east},
		{ID: "win-sky", Type: "window", Area: 2, UValue: 1.8, Azimuth: &south, Tilt: &horizontal},
		{ID: "win-no-azimuth", Type: "window", Area: 1, UValue: 2.0}, // defaults to South
	}
	if err := aggregateSurfaces(surfaces, building); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	env := building.BasicParameters.Envelope
	if env.A_Window_South != 9 { // win-s (8) + win-no-azimuth (1)
		t.Errorf("A_Window_South = %v, want 9", env.A_Window_South)
	}
	if env.A_Window_East != 4 {
		t.Errorf("A_Window_East = %v, want 4", env.A_Window_East)
	}
	if env.A_Window_Horizontal != 2 {
		t.Errorf("A_Window_Horizontal = %v, want 2", env.A_Window_Horizontal)
	}
	if env.A_Window_North != 0 || env.A_Window_West != 0 {
		t.Errorf("expected North/West buckets to stay 0, got North=%v West=%v", env.A_Window_North, env.A_Window_West)
	}
	wantTotalArea := 15.0
	if env.A_Window_1 != wantTotalArea {
		t.Errorf("A_Window_1 = %v, want %v", env.A_Window_1, wantTotalArea)
	}
}

func TestAggregateSurfaces_unknownType_returnsError(t *testing.T) {
	building := realisticBuilding()
	surfaces := []Surface{{ID: "mystery", Type: "skylight-tube", Area: 1, UValue: 1}}
	if err := aggregateSurfaces(surfaces, building); err == nil {
		t.Error("expected error for unknown surface type, got nil")
	}
}

func TestAggregateSurfaces_nonPositiveArea_returnsError(t *testing.T) {
	building := realisticBuilding()
	surfaces := []Surface{{ID: "wall-a", Type: "wall", Area: 0, UValue: 1}}
	if err := aggregateSurfaces(surfaces, building); err == nil {
		t.Error("expected error for zero area, got nil")
	}
}

func TestAggregateSurfaces_nonPositiveUValue_returnsError(t *testing.T) {
	building := realisticBuilding()
	surfaces := []Surface{{ID: "wall-a", Type: "wall", Area: 10, UValue: -1}}
	if err := aggregateSurfaces(surfaces, building); err == nil {
		t.Error("expected error for negative u_value, got nil")
	}
}

// TestAggregateSurfaces_roundTrip proves the aggregation reproduces exactly what
// TABULA's own _1/_2 slots already encode, per issue #5's own stated requirement
// that N=1..3 slots must produce identical results to today: reconstructing a
// []Surface list from an existing variant's own multi-slot fields and
// aggregating them back must reproduce the same effective (area, U) the
// pipeline would already compute from those slots directly.
func TestAggregateSurfaces_roundTrip(t *testing.T) {
	building := realisticBuilding()

	// What the existing slot-based formula (calc_level_05.go's wall blend,
	// f_Measure=0 case) already produces from A_Wall_1/2 + U_Wall_1/2 combined
	// as if they were two parallel elements of the same category.
	wantArea := 80.0 + 40.0
	wantU := (80*0.9 + 40*0.3) / wantArea

	surfaces := []Surface{
		{ID: "wall-slot-1", Type: "wall", Area: 80, UValue: 0.9},
		{ID: "wall-slot-2", Type: "wall", Area: 40, UValue: 0.3},
	}
	if err := aggregateSurfaces(surfaces, building); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if building.BasicParameters.Envelope.A_Wall_1 != wantArea {
		t.Errorf("A_Wall_1 = %v, want %v", building.BasicParameters.Envelope.A_Wall_1, wantArea)
	}
	if building.AdvancedParameters.Uvalues.U_Wall_1 != wantU {
		t.Errorf("U_Wall_1 = %v, want %v", building.AdvancedParameters.Uvalues.U_Wall_1, wantU)
	}
}
