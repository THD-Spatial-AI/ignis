package calc

import "testing"

func TestCalcLevel8_calcREnvTotalExactToEstim_nonZeroDenom(t *testing.T) {
	lvl1 := &CalcLevel1{A_Exact_Env_Sum: 90}
	lvl7 := &CalcLevel7{AEstimEnvSum: 100}
	c := &CalcLevel8{Lvl1: lvl1, Lvl7: lvl7}
	if got := c.calcREnvTotalExactToEstim(); !approxEqual(got, 0.9) {
		t.Errorf("got %.4f, want 0.9", got)
	}
}

func TestCalcLevel8_calcREnvTotalExactToEstim_zeroDenom(t *testing.T) {
	lvl1 := &CalcLevel1{A_Exact_Env_Sum: 90}
	lvl7 := &CalcLevel7{AEstimEnvSum: 0}
	c := &CalcLevel8{Lvl1: lvl1, Lvl7: lvl7}
	if got := c.calcREnvTotalExactToEstim(); got != 0 {
		t.Errorf("got %.4f, want 0", got)
	}
}

func TestCalcLevel8_calcFractionEnvelopeRefurbished_allMeasures(t *testing.T) {
	p := newTestParams()
	// Set all measure fractions to 0.5
	fracs := p.AdvancedParameters.MeasureFractions
	fracs.F_Measure_Roof_1 = 0.5
	fracs.F_Measure_Roof_2 = 0.5
	fracs.F_Measure_Wall_1 = 0.5
	fracs.F_Measure_Wall_2 = 0.5
	fracs.F_Measure_Wall_3 = 0.5
	fracs.F_Measure_Floor_1 = 0.5
	fracs.F_Measure_Floor_2 = 0.5
	fracs.F_Measure_Window_1 = 0.5
	fracs.F_Measure_Window_2 = 0.5
	fracs.F_Measure_Door_1 = 0.5

	// All R_Measure > 0 to trigger all numerator additions
	lvl2 := &CalcLevel2{
		R_Measure_Roof_1: 1, R_Measure_Roof_2: 1,
		R_Measure_Wall_1: 1, R_Measure_Wall_2: 1, R_Measure_Wall_3: 1,
		R_Measure_Floor_1: 1, R_Measure_Floor_2: 1,
	}
	lvl1 := &CalcLevel1{
		R_Measure_Window_1: 1, R_Measure_Window_2: 1, R_Measure_Door_1: 1,
		A_Calc_Window_2: 10,
	}
	// Each area = 10 so denominator = 10*10 = 100, numerator = 10*0.5 = 5 per element × 10 = 50
	const area = 10.0
	lvl5 := &CalcLevel5{
		A_Calc_Roof_1: area, A_Calc_Roof_2: area, A_Calc_Floor_1: area, A_Calc_Window_1: area,
	}
	lvl6 := &CalcLevel6{ACalcWall2: area, ACalcFloor2: area}
	lvl7 := &CalcLevel7{ACalcWall1: area, ACalcWall3: area}
	lvl4 := &CalcLevel4{A_Calc_Door_1: area}

	c := &CalcLevel8{
		Lvl0: p, Lvl1: lvl1, Lvl2: lvl2,
		Lvl4: lvl4, Lvl5: lvl5, Lvl6: lvl6, Lvl7: lvl7,
	}
	got := c.calcFractionEnvelopeRefurbished()
	// numerator = 0.5*10 per element × 10 elements = 50, denominator = 10*10 = 100 → 0.5
	if !approxEqual(got, 0.5) {
		t.Errorf("got %.4f, want 0.5", got)
	}
}

func TestCalcLevel8_calcFractionEnvelopeRefurbished_zeroDenom(t *testing.T) {
	p := newTestParams()
	// All R = 0 so numerator = 0; all areas = 0 so denominator = 0
	c := &CalcLevel8{
		Lvl0: p,
		Lvl1: &CalcLevel1{},
		Lvl2: &CalcLevel2{},
		Lvl4: &CalcLevel4{},
		Lvl5: &CalcLevel5{},
		Lvl6: &CalcLevel6{},
		Lvl7: &CalcLevel7{},
	}
	if got := c.calcFractionEnvelopeRefurbished(); got != 0 {
		t.Errorf("got %.4f, want 0", got)
	}
}

func TestCalcLevel8_calcHTransmissionWall1And3(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.HeatLosses.B_Transmission_Wall_1 = 1.0
	p.AdvancedParameters.HeatLosses.B_Transmission_Wall_3 = 1.0
	lvl5 := &CalcLevel5{U_Actual_Wall_1: 2.0, U_Actual_Wall_3: 3.0}
	lvl7 := &CalcLevel7{ACalcWall1: 10, ACalcWall3: 5}
	c := &CalcLevel8{Lvl0: p, Lvl5: lvl5, Lvl7: lvl7}
	if got := c.calcHTransmissionWall1(); !approxEqual(got, 20) {
		t.Errorf("Wall1: got %.2f, want 20", got)
	}
	if got := c.calcHTransmissionWall3(); !approxEqual(got, 15) {
		t.Errorf("Wall3: got %.2f, want 15", got)
	}
}

func TestCalcLevel8_calcQSol_nonZero(t *testing.T) {
	lvl2 := &CalcLevel2{A_C_Ref: 100}
	lvl7 := &CalcLevel7{QSolHor: 50, QSolEast: 30, QSolSouth: 80, QSolWest: 30, QSolNorth: 10}
	c := &CalcLevel8{Lvl2: lvl2, Lvl7: lvl7}
	// (50+30+80+30+10)/100 = 2.0
	if got := c.calcQSol(); !approxEqual(got, 2.0) {
		t.Errorf("got %.4f, want 2.0", got)
	}
}

func TestCalcLevel8_calcQSol_zeroRef(t *testing.T) {
	lvl2 := &CalcLevel2{A_C_Ref: 0}
	lvl7 := &CalcLevel7{QSolHor: 100}
	c := &CalcLevel8{Lvl2: lvl2, Lvl7: lvl7}
	if got := c.calcQSol(); got != 0 {
		t.Errorf("got %.4f, want 0", got)
	}
}
