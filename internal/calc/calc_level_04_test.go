package calc

import "testing"

func TestCalcLevel4_calcAEstimGrossWallStorey_attachedNeighbours(t *testing.T) {
	tests := []struct {
		code string
		add  float64
	}{
		{"B_N2", 5},
		{"B_N1", 25},
		{"B_Alone", 50},
		{"", 50},
	}
	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			p := newTestParams()
			p.BasicParameters.BuildingAppearance.Code_AttachedNeighbours = tt.code
			p.AdvancedParameters.PredefinedCodes.F_Corr_CeilingHeight = 1.0
			lvl1 := &CalcLevel1{F_ComplexFootprint: 1.0}
			lvl3 := &CalcLevel3{A_C_Storey: 100}
			c := &CalcLevel4{Lvl0: p, Lvl1: lvl1, Lvl3: lvl3}
			// 1.0 * 1.0 * (0.7*100 + tt.add)
			want := 1.0 * 1.0 * (0.7*100 + tt.add)
			if got := c.calcAEstimGrossWallStorey(); !approxEqual(got, want) {
				t.Errorf("code=%s: got %.2f, want %.2f", tt.code, got, want)
			}
		})
	}
}

func TestCalcLevel4_calcAEstimRoof(t *testing.T) {
	lvl1 := &CalcLevel1{F_ComplexRoof: 1.3, P_Roof: 1.6, Q_Roof: 15}
	lvl3 := &CalcLevel3{A_C_Storey: 50}
	c := &CalcLevel4{Lvl1: lvl1, Lvl3: lvl3}
	// 1.3 * (1.6*50 + 15) = 1.3 * 95 = 123.5
	if got := c.calcAEstimRoof(); !approxEqual(got, 123.5) {
		t.Errorf("got %.4f, want 123.5", got)
	}
}

func TestCalcLevel4_calcAEstimUpperCeiling(t *testing.T) {
	lvl1 := &CalcLevel1{P_Ceiling: 0.6, Q_Ceiling: 3}
	lvl3 := &CalcLevel3{A_C_Storey: 50}
	c := &CalcLevel4{Lvl1: lvl1, Lvl3: lvl3}
	// 0.6*50 + 3 = 33
	if got := c.calcAEstimUpperCeiling(); !approxEqual(got, 33) {
		t.Errorf("got %.4f, want 33", got)
	}
}

func TestCalcLevel4_calcAEstimFloor(t *testing.T) {
	lvl3 := &CalcLevel3{A_C_Storey: 50}
	c := &CalcLevel4{Lvl3: lvl3}
	// 1.2*50 + 5 = 65
	if got := c.calcAEstimFloor(); !approxEqual(got, 65) {
		t.Errorf("got %.4f, want 65", got)
	}
}

func TestCalcLevel4_calcAEstimWindow(t *testing.T) {
	lvl2 := &CalcLevel2{A_C_Ref: 100}
	lvl3 := &CalcLevel3{A_Estim_Door: 2.5}
	c := &CalcLevel4{Lvl2: lvl2, Lvl3: lvl3}
	// 0.18*100 - 2.5 = 15.5
	if got := c.calcAEstimWindow(); !approxEqual(got, 15.5) {
		t.Errorf("got %.4f, want 15.5", got)
	}
}

func TestCalcLevel4_calcACalcDoor1_estimation(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.PredefinedCodes.Code_TypeIntake_EnvelopeArea = "Estimation"
	lvl3 := &CalcLevel3{A_Estim_Door: 3.0}
	c := &CalcLevel4{Lvl0: p, Lvl3: lvl3}
	if got := c.calcACalcDoor1(); got != 3.0 {
		t.Errorf("got %.2f, want 3.0", got)
	}
}

func TestCalcLevel4_calcACalcDoor1_notEstimation(t *testing.T) {
	p := newTestParams()
	p.BasicParameters.Envelope.A_Door_1 = 2.1
	lvl3 := &CalcLevel3{A_Estim_Door: 3.0}
	c := &CalcLevel4{Lvl0: p, Lvl3: lvl3}
	if got := c.calcACalcDoor1(); got != 2.1 {
		t.Errorf("got %.2f, want 2.1", got)
	}
}

func TestCalcLevel4_calcUMeasureRoof1_replace(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.MeasureTypes.Code_MeasureType_Roof_1 = "Replace"
	p.AdvancedParameters.AdditionalResistances.R_Add_UnheatedSpace_Roof_1 = 0.2
	lvl2 := &CalcLevel2{R_Measure_Roof_1: 3.0}
	lvl3 := &CalcLevel3{}
	c := &CalcLevel4{Lvl0: p, Lvl2: lvl2, Lvl3: lvl3}
	// Replace: rBefore = R_Add (0.2), denom = 0.2 + 3.0 = 3.2
	if got := c.calcUMeasureRoof1(); !approxEqual(got, 1.0/3.2) {
		t.Errorf("got %.4f, want %.4f", got, 1.0/3.2)
	}
}

func TestCalcLevel4_calcUMeasureRoof1_notReplace(t *testing.T) {
	p := newTestParams()
	lvl2 := &CalcLevel2{R_Measure_Roof_1: 4.0}
	lvl3 := &CalcLevel3{R_Before_Roof_1: 1.0}
	c := &CalcLevel4{Lvl0: p, Lvl2: lvl2, Lvl3: lvl3}
	// denom = 1.0 + 4.0 = 5.0
	if got := c.calcUMeasureRoof1(); !approxEqual(got, 0.2) {
		t.Errorf("got %.4f, want 0.2", got)
	}
}

func TestCalcLevel4_calcUMeasureRoof1_zeroDenom(t *testing.T) {
	p := newTestParams()
	lvl2 := &CalcLevel2{R_Measure_Roof_1: 0}
	lvl3 := &CalcLevel3{R_Before_Roof_1: 0}
	c := &CalcLevel4{Lvl0: p, Lvl2: lvl2, Lvl3: lvl3}
	if got := c.calcUMeasureRoof1(); got != 0 {
		t.Errorf("got %.4f, want 0", got)
	}
}

func TestCalcLevel4_calcUMeasureAllElements(t *testing.T) {
	// Exercise Replace branch for all remaining elements and zero-denom for Floor2
	p := newTestParams()
	p.AdvancedParameters.MeasureTypes.Code_MeasureType_Roof_2 = "Replace"
	p.AdvancedParameters.MeasureTypes.Code_MeasureType_Wall_1 = "Replace"
	p.AdvancedParameters.MeasureTypes.Code_MeasureType_Wall_2 = "Replace"
	p.AdvancedParameters.MeasureTypes.Code_MeasureType_Wall_3 = "Replace"
	p.AdvancedParameters.MeasureTypes.Code_MeasureType_Floor_1 = "Replace"
	p.AdvancedParameters.MeasureTypes.Code_MeasureType_Floor_2 = "" // not replace

	lvl2 := &CalcLevel2{
		R_Measure_Roof_2: 5, R_Measure_Wall_1: 5, R_Measure_Wall_2: 5,
		R_Measure_Wall_3: 5, R_Measure_Floor_1: 5, R_Measure_Floor_2: 0,
	}
	lvl3 := &CalcLevel3{R_Before_Floor_2: 0} // denom = 0 for Floor2
	lvl1 := &CalcLevel1{}
	c := NewCalcLevel4(p, lvl1, lvl2, lvl3)
	for name, got := range map[string]float64{
		"Roof2":  c.U_Measure_Roof_2,
		"Wall1":  c.U_Measure_Wall_1,
		"Wall2":  c.U_Measure_Wall_2,
		"Wall3":  c.U_Measure_Wall_3,
		"Floor1": c.U_Measure_Floor_1,
	} {
		if got == 0 {
			t.Errorf("%s: got 0, expected non-zero", name)
		}
	}
	if c.U_Measure_Floor_2 != 0 {
		t.Errorf("U_Measure_Floor_2: got %.4f, want 0 (zero denom)", c.U_Measure_Floor_2)
	}
}
