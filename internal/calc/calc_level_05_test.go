package calc

import "testing"

func TestCalcLevel5_calcAEstimWallToCellarOrSoil_suffixI(t *testing.T) {
	p := newTestParams()
	p.BasicParameters.BuildingAppearance.Code_CellarCond = "CI"
	lvl1 := &CalcLevel1{F_CellarCond: 1.0}
	lvl4 := &CalcLevel4{A_Estim_GrossWall_Storey: 100}
	c := &CalcLevel5{Lvl0: p, Lvl1: lvl1, Lvl4: lvl4}
	// suffix "I" → cellarCond = 1, result = 0.5*1*100 = 50
	if got := c.calcAEstimWallToCellarOrSoil(); !approxEqual(got, 50) {
		t.Errorf("got %.4f, want 50", got)
	}
}

func TestCalcLevel5_calcAEstimWallToCellarOrSoil_noSuffixI(t *testing.T) {
	p := newTestParams()
	p.BasicParameters.BuildingAppearance.Code_CellarCond = "P"
	lvl1 := &CalcLevel1{F_CellarCond: 0.5}
	lvl4 := &CalcLevel4{A_Estim_GrossWall_Storey: 100}
	c := &CalcLevel5{Lvl0: p, Lvl1: lvl1, Lvl4: lvl4}
	// no suffix I → uses F_CellarCond 0.5, result = 0.5*0.5*100 = 25
	if got := c.calcAEstimWallToCellarOrSoil(); !approxEqual(got, 25) {
		t.Errorf("got %.4f, want 25", got)
	}
}

func TestCalcLevel5_calcREnvFloorExactToEstim_nonZeroDenom(t *testing.T) {
	p := newTestParams()
	p.BasicParameters.Envelope.A_Floor_1 = 30
	p.BasicParameters.Envelope.A_Floor_2 = 20
	lvl4 := &CalcLevel4{A_Estim_Floor: 100}
	c := &CalcLevel5{Lvl0: p, Lvl4: lvl4}
	if got := c.calcREnvFloorExactToEstim(); !approxEqual(got, 0.5) {
		t.Errorf("got %.4f, want 0.5", got)
	}
}

func TestCalcLevel5_calcREnvFloorExactToEstim_zeroDenom(t *testing.T) {
	p := newTestParams()
	lvl4 := &CalcLevel4{A_Estim_Floor: 0}
	c := &CalcLevel5{Lvl0: p, Lvl4: lvl4}
	if got := c.calcREnvFloorExactToEstim(); got != 0 {
		t.Errorf("got %.4f, want 0", got)
	}
}

func TestCalcLevel5_calcREnvWindowExactToEstim_nonZeroDenom(t *testing.T) {
	p := newTestParams()
	p.BasicParameters.Envelope.A_Window_1 = 10
	p.BasicParameters.Envelope.A_Window_2 = 10
	p.BasicParameters.Envelope.A_Door_1 = 5
	lvl4 := &CalcLevel4{A_Estim_Window: 20}
	lvl3 := &CalcLevel3{A_Estim_Door: 5}
	c := &CalcLevel5{Lvl0: p, Lvl4: lvl4, Lvl3: lvl3}
	// (10+10+5)/(20+5) = 25/25 = 1.0
	if got := c.calcREnvWindowExactToEstim(); !approxEqual(got, 1.0) {
		t.Errorf("got %.4f, want 1.0", got)
	}
}

func TestCalcLevel5_calcREnvWindowExactToEstim_zeroDenom(t *testing.T) {
	p := newTestParams()
	lvl4 := &CalcLevel4{A_Estim_Window: 0}
	lvl3 := &CalcLevel3{A_Estim_Door: 0}
	c := &CalcLevel5{Lvl0: p, Lvl4: lvl4, Lvl3: lvl3}
	if got := c.calcREnvWindowExactToEstim(); got != 0 {
		t.Errorf("got %.4f, want 0", got)
	}
}

func TestCalcLevel5_calcACalcRoof_estimation(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.PredefinedCodes.Code_TypeIntake_EnvelopeArea = "Estimation"
	lvl4 := &CalcLevel4{A_Estim_Roof: 80, A_Estim_UpperCeiling: 60}
	lvl1 := &CalcLevel1{}
	c := &CalcLevel5{Lvl0: p, Lvl1: lvl1, Lvl4: lvl4}
	if got := c.calcACalcRoof1(); !approxEqual(got, 80) {
		t.Errorf("ACalcRoof1: got %.2f, want 80", got)
	}
	if got := c.calcACalcRoof2(); !approxEqual(got, 60) {
		t.Errorf("ACalcRoof2: got %.2f, want 60", got)
	}
}

func TestCalcLevel5_calcACalcRoof_notEstimation(t *testing.T) {
	p := newTestParams()
	p.BasicParameters.Envelope.A_Roof_1 = 50
	p.BasicParameters.Envelope.A_Roof_2 = 40
	lvl4 := &CalcLevel4{A_Estim_Roof: 80, A_Estim_UpperCeiling: 60}
	lvl1 := &CalcLevel1{}
	c := &CalcLevel5{Lvl0: p, Lvl1: lvl1, Lvl4: lvl4}
	if got := c.calcACalcRoof1(); !approxEqual(got, 50) {
		t.Errorf("ACalcRoof1: got %.2f, want 50", got)
	}
	if got := c.calcACalcRoof2(); !approxEqual(got, 40) {
		t.Errorf("ACalcRoof2: got %.2f, want 40", got)
	}
}

func TestCalcLevel5_calcACalcFloor1_estimation_soil(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.PredefinedCodes.Code_TypeIntake_EnvelopeArea = "Estimation"
	lvl1 := &CalcLevel1{Code_Estim_ConstructionBorder_Floor: "Soil"}
	lvl4 := &CalcLevel4{A_Estim_Floor: 50}
	c := &CalcLevel5{Lvl0: p, Lvl1: lvl1, Lvl4: lvl4}
	if got := c.calcACalcFloor1(); got != 0 {
		t.Errorf("got %.2f, want 0", got)
	}
}

func TestCalcLevel5_calcACalcFloor1_estimation_cellar(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.PredefinedCodes.Code_TypeIntake_EnvelopeArea = "Estimation"
	lvl1 := &CalcLevel1{Code_Estim_ConstructionBorder_Floor: "Cellar"}
	lvl4 := &CalcLevel4{A_Estim_Floor: 50}
	c := &CalcLevel5{Lvl0: p, Lvl1: lvl1, Lvl4: lvl4}
	if got := c.calcACalcFloor1(); !approxEqual(got, 50) {
		t.Errorf("got %.2f, want 50", got)
	}
}

func TestCalcLevel5_calcACalcFloor1_notEstimation(t *testing.T) {
	p := newTestParams()
	p.BasicParameters.Envelope.A_Floor_1 = 30
	lvl1 := &CalcLevel1{}
	lvl4 := &CalcLevel4{A_Estim_Floor: 50}
	c := &CalcLevel5{Lvl0: p, Lvl1: lvl1, Lvl4: lvl4}
	if got := c.calcACalcFloor1(); !approxEqual(got, 30) {
		t.Errorf("got %.2f, want 30", got)
	}
}

func TestCalcLevel5_calcACalcWindow1(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.PredefinedCodes.Code_TypeIntake_EnvelopeArea = "Estimation"
	p.BasicParameters.Envelope.A_Window_1 = 10
	lvl1 := &CalcLevel1{}
	lvl4 := &CalcLevel4{A_Estim_Window: 20}
	c := &CalcLevel5{Lvl0: p, Lvl1: lvl1, Lvl4: lvl4}
	if got := c.calcACalcWindow1(); !approxEqual(got, 20) {
		t.Errorf("estimation: got %.2f, want 20", got)
	}
	p.AdvancedParameters.PredefinedCodes.Code_TypeIntake_EnvelopeArea = ""
	if got := c.calcACalcWindow1(); !approxEqual(got, 10) {
		t.Errorf("exact: got %.2f, want 10", got)
	}
}

func TestCalcLevel5_calcUActualRoof_positiveU(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.Uvalues.U_Roof_1 = 2.0
	p.AdvancedParameters.AdditionalResistances.R_Add_UnheatedSpace_Roof_1 = 0
	p.AdvancedParameters.MeasureFractions.F_Measure_Roof_1 = 0.5
	lvl4 := &CalcLevel4{U_Measure_Roof_1: 1.0}
	c := &CalcLevel5{Lvl0: p, Lvl4: lvl4}
	// (1-0.5) * 1/(1/2+0) + 0.5*1.0 = 0.5*2.0 + 0.5 = 1.5
	if got := c.calcUActualRoof1(); !approxEqual(got, 1.5) {
		t.Errorf("got %.4f, want 1.5", got)
	}
}

func TestCalcLevel5_calcUActualRoof_zeroU(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.Uvalues.U_Roof_1 = 0
	p.AdvancedParameters.MeasureFractions.F_Measure_Roof_1 = 0.5
	lvl4 := &CalcLevel4{U_Measure_Roof_1: 1.0}
	c := &CalcLevel5{Lvl0: p, Lvl4: lvl4}
	// U=0 → result = 0, + 0.5*1.0 = 0.5
	if got := c.calcUActualRoof1(); !approxEqual(got, 0.5) {
		t.Errorf("got %.4f, want 0.5", got)
	}
}

func TestCalcLevel5_calcUActualAllElements(t *testing.T) {
	// Exercise all remaining U_Actual branches via NewCalcLevel5
	p := newTestParams()
	p.AdvancedParameters.Uvalues.U_Roof_2 = 1.0
	p.AdvancedParameters.Uvalues.U_Wall_1 = 1.0
	p.AdvancedParameters.Uvalues.U_Wall_2 = 1.0
	p.AdvancedParameters.Uvalues.U_Wall_3 = 1.0
	p.AdvancedParameters.Uvalues.U_Floor_1 = 1.0
	p.AdvancedParameters.Uvalues.U_Floor_2 = 1.0
	lvl1 := &CalcLevel1{}
	lvl3 := &CalcLevel3{A_Estim_Door: 2}
	lvl4 := &CalcLevel4{A_Estim_Floor: 50}
	c := NewCalcLevel5(p, lvl1, lvl3, lvl4)
	// Just verify they were computed (non-zero from U=1 and fraction=0)
	for name, got := range map[string]float64{
		"Roof2":  c.U_Actual_Roof_2,
		"Wall1":  c.U_Actual_Wall_1,
		"Wall2":  c.U_Actual_Wall_2,
		"Wall3":  c.U_Actual_Wall_3,
		"Floor1": c.U_Actual_Floor_1,
		"Floor2": c.U_Actual_Floor_2,
	} {
		if got == 0 {
			t.Errorf("%s: expected non-zero when U=1.0", name)
		}
	}
}

func TestCalcLevel5_calcHTransmissionDoor1(t *testing.T) {
	lvl3 := &CalcLevel3{U_Actual_Door_1: 2.0}
	lvl4 := &CalcLevel4{A_Calc_Door_1: 3.0}
	c := &CalcLevel5{Lvl3: lvl3, Lvl4: lvl4}
	if got := c.calcHTransmissionDoor1(); !approxEqual(got, 6.0) {
		t.Errorf("got %.4f, want 6.0", got)
	}
}
