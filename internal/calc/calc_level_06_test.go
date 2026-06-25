package calc

import "testing"

func TestCalcLevel6_calcAEstimWallExtAir(t *testing.T) {
	lvl2 := &CalcLevel2{N_Storey_effective_envelope: 2.0}
	lvl4 := &CalcLevel4{A_Estim_GrossWall_Storey: 50, A_Estim_Window: 10}
	lvl5 := &CalcLevel5{A_Estim_Wall_ToCellarOrSoil: 5}
	lvl3 := &CalcLevel3{A_Estim_Door: 2}
	c := &CalcLevel6{Lvl2: lvl2, Lvl4: lvl4, Lvl5: lvl5, Lvl3: lvl3}
	// 2.0*50 - 5 - 10 - 2 = 83
	if got := c.calcAEstimWallExtAir(); !approxEqual(got, 83) {
		t.Errorf("got %.4f, want 83", got)
	}
}

func TestCalcLevel6_calcCheckFloorAreaExactToEstim_within(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.PredefinedCodes.F_PlausiCrit_FloorArea_LowerLimit = 0.5
	p.AdvancedParameters.PredefinedCodes.F_PlausiCrit_FloorArea_UpperLimit = 1.5
	lvl5 := &CalcLevel5{R_EnvFloor_ExactToEstim: 1.0}
	c := &CalcLevel6{Lvl0: p, Lvl5: lvl5}
	if got := c.calcCheckFloorAreaExactToEstim(); got != 1 {
		t.Errorf("got %d, want 1", got)
	}
}

func TestCalcLevel6_calcCheckFloorAreaExactToEstim_outside(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.PredefinedCodes.F_PlausiCrit_FloorArea_LowerLimit = 0.5
	p.AdvancedParameters.PredefinedCodes.F_PlausiCrit_FloorArea_UpperLimit = 1.5
	lvl5 := &CalcLevel5{R_EnvFloor_ExactToEstim: 2.0}
	c := &CalcLevel6{Lvl0: p, Lvl5: lvl5}
	if got := c.calcCheckFloorAreaExactToEstim(); got != 0 {
		t.Errorf("got %d, want 0", got)
	}
}

func TestCalcLevel6_calcCheckWindowAreaExactToEstim(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.PredefinedCodes.F_PlausiCrit_WindowArea_LowerLimit = 0.5
	p.AdvancedParameters.PredefinedCodes.F_PlausiCrit_WindowArea_UpperLimit = 1.5
	lvl5ok := &CalcLevel5{R_EnvWindow_ExactToEstim: 1.0}
	c := &CalcLevel6{Lvl0: p, Lvl5: lvl5ok}
	if got := c.calcCheckWindowAreaExactToEstim(); got != 1 {
		t.Errorf("within: got %d, want 1", got)
	}
	lvl5out := &CalcLevel5{R_EnvWindow_ExactToEstim: 3.0}
	c.Lvl5 = lvl5out
	if got := c.calcCheckWindowAreaExactToEstim(); got != 0 {
		t.Errorf("outside: got %d, want 0", got)
	}
}

func TestCalcLevel6_calcACalcWall2_estimation_soil(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.PredefinedCodes.Code_TypeIntake_EnvelopeArea = "Estimation"
	lvl1 := &CalcLevel1{Code_Estim_ConstructionBorder_Wall_ToCellarOrSoil: "Soil"}
	lvl5 := &CalcLevel5{A_Estim_Wall_ToCellarOrSoil: 30}
	c := &CalcLevel6{Lvl0: p, Lvl1: lvl1, Lvl5: lvl5}
	if got := c.calcACalcWall2(); got != 0 {
		t.Errorf("Soil: got %.2f, want 0", got)
	}
}

func TestCalcLevel6_calcACalcWall2_estimation_notSoil(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.PredefinedCodes.Code_TypeIntake_EnvelopeArea = "Estimation"
	lvl1 := &CalcLevel1{Code_Estim_ConstructionBorder_Wall_ToCellarOrSoil: "Unh"}
	lvl5 := &CalcLevel5{A_Estim_Wall_ToCellarOrSoil: 30}
	c := &CalcLevel6{Lvl0: p, Lvl1: lvl1, Lvl5: lvl5}
	if got := c.calcACalcWall2(); !approxEqual(got, 30) {
		t.Errorf("Unh: got %.2f, want 30", got)
	}
}

func TestCalcLevel6_calcACalcWall2_notEstimation(t *testing.T) {
	p := newTestParams()
	p.BasicParameters.Envelope.A_Wall_2 = 25
	lvl1 := &CalcLevel1{}
	lvl5 := &CalcLevel5{}
	c := &CalcLevel6{Lvl0: p, Lvl1: lvl1, Lvl5: lvl5}
	if got := c.calcACalcWall2(); !approxEqual(got, 25) {
		t.Errorf("got %.2f, want 25", got)
	}
}

func TestCalcLevel6_calcACalcFloor2(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.PredefinedCodes.Code_TypeIntake_EnvelopeArea = "Estimation"
	p.BasicParameters.Envelope.A_Floor_2 = 20
	lvl4 := &CalcLevel4{A_Estim_Floor: 80}
	lvl5 := &CalcLevel5{A_Calc_Floor_1: 30}
	c := &CalcLevel6{Lvl0: p, Lvl4: lvl4, Lvl5: lvl5}
	if got := c.calcACalcFloor2(); !approxEqual(got, 50) {
		t.Errorf("estimation: got %.2f, want 50", got)
	}
	p.AdvancedParameters.PredefinedCodes.Code_TypeIntake_EnvelopeArea = ""
	if got := c.calcACalcFloor2(); !approxEqual(got, 20) {
		t.Errorf("exact: got %.2f, want 20", got)
	}
}

func TestCalcLevel6_calcACalcWindowEastWest(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.PredefinedCodes.Code_TypeIntake_EnvelopeArea = "Estimation"
	p.BasicParameters.Envelope.A_Window_East = 5
	p.BasicParameters.Envelope.A_Window_West = 6
	lvl5 := &CalcLevel5{A_Calc_Window_1: 10}
	lvl1 := &CalcLevel1{A_Calc_Window_2: 6}
	c := &CalcLevel6{Lvl0: p, Lvl1: lvl1, Lvl5: lvl5}
	// Estimation: 0.5*(10+6) = 8
	if got := c.calcACalcWindowEast(); !approxEqual(got, 8) {
		t.Errorf("East estimation: got %.2f, want 8", got)
	}
	if got := c.calcACalcWindowWest(); !approxEqual(got, 8) {
		t.Errorf("West estimation: got %.2f, want 8", got)
	}
	p.AdvancedParameters.PredefinedCodes.Code_TypeIntake_EnvelopeArea = ""
	if got := c.calcACalcWindowEast(); !approxEqual(got, 5) {
		t.Errorf("East exact: got %.2f, want 5", got)
	}
	if got := c.calcACalcWindowWest(); !approxEqual(got, 6) {
		t.Errorf("West exact: got %.2f, want 6", got)
	}
}

func TestCalcLevel6_calcHTransmissions(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.HeatLosses.B_Transmission_Roof_1 = 1.0
	p.AdvancedParameters.HeatLosses.B_Transmission_Roof_2 = 1.0
	p.AdvancedParameters.HeatLosses.B_Transmission_Floor_1 = 1.0
	lvl5 := &CalcLevel5{
		U_Actual_Roof_1:  2.0, A_Calc_Roof_1: 10,
		U_Actual_Roof_2:  3.0, A_Calc_Roof_2: 10,
		U_Actual_Floor_1: 1.5, A_Calc_Floor_1: 10,
		A_Calc_Window_1: 10,
	}
	lvl3 := &CalcLevel3{U_Actual_Window_1: 2.0, U_Actual_Window_2: 1.0}
	lvl1 := &CalcLevel1{A_Calc_Window_2: 5}
	c := &CalcLevel6{Lvl0: p, Lvl5: lvl5, Lvl3: lvl3, Lvl1: lvl1}
	if got := c.calcHTransmissionRoof1(); !approxEqual(got, 20) {
		t.Errorf("Roof1: got %.2f, want 20", got)
	}
	if got := c.calcHTransmissionRoof2(); !approxEqual(got, 30) {
		t.Errorf("Roof2: got %.2f, want 30", got)
	}
	if got := c.calcHTransmissionFloor1(); !approxEqual(got, 15) {
		t.Errorf("Floor1: got %.2f, want 15", got)
	}
	if got := c.calcHTransmissionWindow1(); !approxEqual(got, 2.0*10) {
		t.Errorf("Window1: got %.2f", got)
	}
	if got := c.calcHTransmissionWindow2(); !approxEqual(got, 1.0*5) {
		t.Errorf("Window2: got %.2f, want 5", got)
	}
}

func TestCalcLevel6_calcGGlN_withMeasureCode(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.MeasureTypes.Code_Measure_Window_1 = "M1"
	p.AdvancedParameters.MeasureTypes.Code_Measure_Window_2 = "M2"
	p.AdvancedParameters.SolarTransmittance.G_gl_n_Window_1 = 0.5
	p.AdvancedParameters.SolarTransmittance.G_gl_n_Window_2 = 0.4
	lvl1 := &CalcLevel1{
		G_gl_n_Measure_Window_1: 0.6,
		G_gl_n_Measure_Window_2: 0.7,
		A_Calc_Window_2:         10,
	}
	lvl5 := &CalcLevel5{A_Calc_Window_1: 10}
	c := &CalcLevel6{Lvl0: p, Lvl1: lvl1, Lvl5: lvl5}
	// (0.6*10 + 0.7*10) / (10+10) = 13/20 = 0.65
	if got := c.calcGGlN(); !approxEqual(got, 0.65) {
		t.Errorf("got %.4f, want 0.65", got)
	}
}

func TestCalcLevel6_calcGGlN_noMeasureCode(t *testing.T) {
	p := newTestParams()
	// Code = "" and "0" → use direct g_gl_n values
	p.AdvancedParameters.MeasureTypes.Code_Measure_Window_1 = ""
	p.AdvancedParameters.MeasureTypes.Code_Measure_Window_2 = "0"
	p.AdvancedParameters.SolarTransmittance.G_gl_n_Window_1 = 0.5
	p.AdvancedParameters.SolarTransmittance.G_gl_n_Window_2 = 0.4
	lvl1 := &CalcLevel1{A_Calc_Window_2: 10}
	lvl5 := &CalcLevel5{A_Calc_Window_1: 10}
	c := &CalcLevel6{Lvl0: p, Lvl1: lvl1, Lvl5: lvl5}
	// (0.5*10 + 0.4*10) / 20 = 0.45
	if got := c.calcGGlN(); !approxEqual(got, 0.45) {
		t.Errorf("got %.4f, want 0.45", got)
	}
}

func TestCalcLevel6_calcGGlN_zeroDenom(t *testing.T) {
	p := newTestParams()
	lvl1 := &CalcLevel1{A_Calc_Window_2: 0}
	lvl5 := &CalcLevel5{A_Calc_Window_1: 0}
	c := &CalcLevel6{Lvl0: p, Lvl1: lvl1, Lvl5: lvl5}
	if got := c.calcGGlN(); got != 0 {
		t.Errorf("got %.4f, want 0", got)
	}
}
