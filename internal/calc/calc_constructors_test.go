package calc

// This file exercises the New/Run entry points for levels whose constructors
// were not called by the level-specific test files (which test private methods
// directly). One call per level is enough to bring NewCalcLevelN and Run to 100%.

import "testing"

func TestCalcLevel3_constructor(t *testing.T) {
	p := newTestParams()
	lvl2 := &CalcLevel2{A_C_Ref: 100, N_Storey_effective: 2}
	if c := NewCalcLevel3(p, lvl2); c == nil {
		t.Error("expected non-nil CalcLevel3")
	}
}

func TestCalcLevel6_constructor(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.PredefinedCodes.Code_TypeIntake_EnvelopeArea = "Estimation"
	lvl1 := &CalcLevel1{}
	lvl2 := &CalcLevel2{}
	lvl3 := &CalcLevel3{}
	lvl4 := &CalcLevel4{}
	lvl5 := &CalcLevel5{}
	if c := NewCalcLevel6(p, lvl1, lvl2, lvl3, lvl4, lvl5); c == nil {
		t.Error("expected non-nil CalcLevel6")
	}
}

func TestCalcLevel7_constructor(t *testing.T) {
	p := newTestParams()
	if c := NewCalcLevel7(p, &CalcLevel1{}, &CalcLevel3{}, &CalcLevel4{}, &CalcLevel5{}, &CalcLevel6{}); c == nil {
		t.Error("expected non-nil CalcLevel7")
	}
}

func TestCalcLevel8_constructor(t *testing.T) {
	p := newTestParams()
	if c := NewCalcLevel8(p, &CalcLevel1{}, &CalcLevel2{}, &CalcLevel4{}, &CalcLevel5{}, &CalcLevel6{}, &CalcLevel7{}); c == nil {
		t.Error("expected non-nil CalcLevel8")
	}
}

func TestCalcLevel9_constructor(t *testing.T) {
	p := newTestParams()
	if c := NewCalcLevel9(p, &CalcLevel8{}); c == nil {
		t.Error("expected non-nil CalcLevel9")
	}
}

func TestCalcLevel10_constructor(t *testing.T) {
	if c := NewCalcLevel10(&CalcLevel1{}, &CalcLevel2{}, &CalcLevel4{}, &CalcLevel5{}, &CalcLevel6{}, &CalcLevel7{}, &CalcLevel9{}); c == nil {
		t.Error("expected non-nil CalcLevel10")
	}
}

func TestCalcLevel11_constructor(t *testing.T) {
	if c := NewCalcLevel11(&CalcLevel2{}, &CalcLevel5{}, &CalcLevel6{}, &CalcLevel7{}, &CalcLevel8{}, &CalcLevel10{}); c == nil {
		t.Error("expected non-nil CalcLevel11")
	}
}

func TestCalcLevel12_constructor(t *testing.T) {
	p := newTestParams()
	if c := NewCalcLevel12(p, &CalcLevel1{}, &CalcLevel11{}); c == nil {
		t.Error("expected non-nil CalcLevel12")
	}
}

func TestCalcLevel13_constructor(t *testing.T) {
	if c := NewCalcLevel13(&CalcLevel1{}, &CalcLevel11{}, &CalcLevel12{}); c == nil {
		t.Error("expected non-nil CalcLevel13")
	}
}

func TestCalcLevel14_constructor(t *testing.T) {
	if c := NewCalcLevel14(&CalcLevel13{}); c == nil {
		t.Error("expected non-nil CalcLevel14")
	}
}

func TestCalcLevel15_constructor(t *testing.T) {
	if c := NewCalcLevel15(&CalcLevel1{}, &CalcLevel8{}, &CalcLevel14{}); c == nil {
		t.Error("expected non-nil CalcLevel15")
	}
}

func TestCalcLevel16_constructor(t *testing.T) {
	if c := NewCalcLevel16(&CalcLevel13{}, &CalcLevel15{}); c == nil {
		t.Error("expected non-nil CalcLevel16")
	}
}

// --- Missing branch coverage ---

// Level 2: calcRMeasureRoof2 through calcRMeasureFloor2 each need the
// "predefined == 0 → return R directly" fallback branch tested.
func TestCalcLevel2_calcRMeasure_predefinedZeroFallback(t *testing.T) {
	p := newTestParams()
	// All predefined insulation thicknesses left at zero
	tr := p.AdvancedParameters.ThermalResistances
	tr.R_PredefinedMeasure_Roof_2 = 4.0
	tr.R_PredefinedMeasure_Wall_1 = 3.0
	tr.R_PredefinedMeasure_Wall_2 = 3.0
	tr.R_PredefinedMeasure_Wall_3 = 3.0
	tr.R_PredefinedMeasure_Floor_1 = 2.0
	tr.R_PredefinedMeasure_Floor_2 = 2.0
	lvl1 := &CalcLevel1{}
	c := &CalcLevel2{Lvl0: p, Lvl1: lvl1}
	for name, got := range map[string]float64{
		"Roof2":  c.calcRMeasureRoof2(),
		"Wall1":  c.calcRMeasureWall1(),
		"Wall2":  c.calcRMeasureWall2(),
		"Wall3":  c.calcRMeasureWall3(),
		"Floor1": c.calcRMeasureFloor1(),
		"Floor2": c.calcRMeasureFloor2(),
	} {
		if got == 0 {
			t.Errorf("%s: expected non-zero R when predefined==0 falls back to R_Predefined", name)
		}
	}
}

// Level 2: calcUMeasureDoor1 needs the "not Replace, non-zero denom" path.
func TestCalcLevel2_calcUMeasureDoor1_notReplace_nonZeroDenom(t *testing.T) {
	p := newTestParams()
	// Code stays empty (not "Replace")
	lvl1 := &CalcLevel1{R_Before_Door_1: 1.0, R_Measure_Door_1: 1.0}
	c := &CalcLevel2{Lvl0: p, Lvl1: lvl1}
	// denom = 2.0 → 0.5
	if got := c.calcUMeasureDoor1(); !approxEqual(got, 0.5) {
		t.Errorf("got %.4f, want 0.5", got)
	}
}

// Level 3: calcRBeforeWall2 and calcRBeforeWall3 need the ReplaceInsulation branch.
func TestCalcLevel3_calcRBeforeWall2_replaceInsulation(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.Uvalues.U_Wall_2 = 0
	p.AdvancedParameters.MeasureTypes.Code_MeasureType_Wall_2 = "ReplaceInsulation"
	p.AdvancedParameters.Insulation.D_Insulation_Wall_2 = 0.04 // /0.04 = 1.0
	c := &CalcLevel3{Lvl0: p, Lvl2: &CalcLevel2{}}
	if got := c.calcRBeforeWall2(); !approxEqual(got, -1.0) {
		t.Errorf("got %.4f, want -1.0", got)
	}
}

func TestCalcLevel3_calcRBeforeWall3_replaceInsulation(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.MeasureTypes.Code_MeasureType_Wall_3 = "ReplaceInsulation"
	p.AdvancedParameters.Insulation.D_Insulation_Wall_3 = 0.08 // 2.0
	c := &CalcLevel3{Lvl0: p, Lvl2: &CalcLevel2{}}
	if got := c.calcRBeforeWall3(); !approxEqual(got, -2.0) {
		t.Errorf("got %.4f, want -2.0", got)
	}
}

// Level 4: calcUMeasureRoof2–Floor1 need the zero-denominator "return 0" path,
// and calcUMeasureFloor2 needs the "Replace + non-zero denom" path.
func TestCalcLevel4_calcUMeasure_zeroDenom(t *testing.T) {
	p := newTestParams()
	// No Code set (else branch), all R values zero → denom = 0
	lvl2 := &CalcLevel2{}
	lvl3 := &CalcLevel3{}
	lvl1 := &CalcLevel1{}
	c := NewCalcLevel4(p, lvl1, lvl2, lvl3)
	for name, got := range map[string]float64{
		"Roof2":  c.U_Measure_Roof_2,
		"Wall1":  c.U_Measure_Wall_1,
		"Wall2":  c.U_Measure_Wall_2,
		"Wall3":  c.U_Measure_Wall_3,
		"Floor1": c.U_Measure_Floor_1,
	} {
		if got != 0 {
			t.Errorf("%s: expected 0 when denom==0, got %.4f", name, got)
		}
	}
}

func TestCalcLevel4_calcUMeasureFloor2_replace_nonZeroDenom(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.MeasureTypes.Code_MeasureType_Floor_2 = "Replace"
	p.AdvancedParameters.AdditionalResistances.R_Add_UnheatedSpace_Floor_2 = 1.0
	lvl2 := &CalcLevel2{R_Measure_Floor_2: 1.0}
	lvl3 := &CalcLevel3{}
	lvl1 := &CalcLevel1{}
	c := &CalcLevel4{Lvl0: p, Lvl1: lvl1, Lvl2: lvl2, Lvl3: lvl3}
	// Replace: rBefore = R_Add = 1.0, denom = 1+1 = 2 → 0.5
	if got := c.calcUMeasureFloor2(); !approxEqual(got, 0.5) {
		t.Errorf("got %.4f, want 0.5", got)
	}
}

// Level 4: calcUMeasureRoof2 through calcUMeasureFloor1 need the
// "not Replace, non-zero denom" path.
func TestCalcLevel4_calcUMeasure_notReplace_nonZeroDenom(t *testing.T) {
	p := newTestParams()
	// No Code_MeasureType set → else branch uses R_Before from lvl3
	lvl2 := &CalcLevel2{
		R_Measure_Roof_2: 2.0, R_Measure_Wall_1: 2.0, R_Measure_Wall_2: 2.0,
		R_Measure_Wall_3: 2.0, R_Measure_Floor_1: 2.0,
	}
	lvl3 := &CalcLevel3{
		R_Before_Roof_2: 2.0, R_Before_Wall_1: 2.0, R_Before_Wall_2: 2.0,
		R_Before_Wall_3: 2.0, R_Before_Floor_1: 2.0,
	}
	lvl1 := &CalcLevel1{}
	c := NewCalcLevel4(p, lvl1, lvl2, lvl3)
	// Each: denom = 2+2 = 4 → 0.25
	for name, got := range map[string]float64{
		"Roof2":  c.U_Measure_Roof_2,
		"Wall1":  c.U_Measure_Wall_1,
		"Wall2":  c.U_Measure_Wall_2,
		"Wall3":  c.U_Measure_Wall_3,
		"Floor1": c.U_Measure_Floor_1,
	} {
		if !approxEqual(got, 0.25) {
			t.Errorf("%s: got %.4f, want 0.25", name, got)
		}
	}
}
