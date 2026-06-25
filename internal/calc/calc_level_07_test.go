package calc

import "testing"

func TestCalcLevel7_calcAEstimEnvSum(t *testing.T) {
	lvl4 := &CalcLevel4{A_Estim_Roof: 10, A_Estim_UpperCeiling: 5, A_Estim_Window: 8, A_Estim_Floor: 20}
	lvl5 := &CalcLevel5{A_Estim_Wall_ToCellarOrSoil: 15}
	lvl6 := &CalcLevel6{AEstimWallExtAir: 30}
	lvl3 := &CalcLevel3{A_Estim_Door: 2}
	c := &CalcLevel7{Lvl4: lvl4, Lvl5: lvl5, Lvl6: lvl6, Lvl3: lvl3}
	// 10+5+30+15+20+8+2 = 90
	if got := c.calcAEstimEnvSum(); !approxEqual(got, 90) {
		t.Errorf("got %.2f, want 90", got)
	}
}

func TestCalcLevel7_calcACalcWall1_estimation(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.PredefinedCodes.Code_TypeIntake_EnvelopeArea = "Estimation"
	p.BasicParameters.Envelope.A_Wall_1 = 40
	lvl6 := &CalcLevel6{AEstimWallExtAir: 60}
	c := &CalcLevel7{Lvl0: p, Lvl6: lvl6}
	if got := c.calcACalcWall1(); !approxEqual(got, 60) {
		t.Errorf("estimation: got %.2f, want 60", got)
	}
	p.AdvancedParameters.PredefinedCodes.Code_TypeIntake_EnvelopeArea = ""
	if got := c.calcACalcWall1(); !approxEqual(got, 40) {
		t.Errorf("exact: got %.2f, want 40", got)
	}
}

func TestCalcLevel7_calcACalcWall3(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.PredefinedCodes.Code_TypeIntake_EnvelopeArea = "Estimation"
	p.BasicParameters.Envelope.A_Wall_3 = 12
	lvl5 := &CalcLevel5{A_Estim_Wall_ToCellarOrSoil: 20}
	lvl6 := &CalcLevel6{ACalcWall2: 8}
	c := &CalcLevel7{Lvl0: p, Lvl5: lvl5, Lvl6: lvl6}
	// estimation: 20 - 8 = 12
	if got := c.calcACalcWall3(); !approxEqual(got, 12) {
		t.Errorf("estimation: got %.2f, want 12", got)
	}
	p.AdvancedParameters.PredefinedCodes.Code_TypeIntake_EnvelopeArea = ""
	if got := c.calcACalcWall3(); !approxEqual(got, 12) {
		t.Errorf("exact: got %.2f, want 12", got)
	}
}

func TestCalcLevel7_calcHTransmissionWall2(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.HeatLosses.B_Transmission_Wall_2 = 1.0
	lvl5 := &CalcLevel5{U_Actual_Wall_2: 2.0}
	lvl6 := &CalcLevel6{ACalcWall2: 10}
	c := &CalcLevel7{Lvl0: p, Lvl5: lvl5, Lvl6: lvl6}
	if got := c.calcHTransmissionWall2(); !approxEqual(got, 20) {
		t.Errorf("got %.2f, want 20", got)
	}
}

func TestCalcLevel7_calcHTransmissionFloor2(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.HeatLosses.B_Transmission_Floor_2 = 1.0
	lvl5 := &CalcLevel5{U_Actual_Floor_2: 1.5}
	lvl6 := &CalcLevel6{ACalcFloor2: 10}
	c := &CalcLevel7{Lvl0: p, Lvl5: lvl5, Lvl6: lvl6}
	if got := c.calcHTransmissionFloor2(); !approxEqual(got, 15) {
		t.Errorf("got %.2f, want 15", got)
	}
}

func TestCalcLevel7_calcQSol(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.SolarGains.I_Sol_Horizontal = 100
	p.AdvancedParameters.SolarGains.I_Sol_East = 80
	p.AdvancedParameters.SolarGains.I_Sol_South = 120
	p.AdvancedParameters.SolarGains.I_Sol_West = 80
	p.AdvancedParameters.SolarGains.I_Sol_North = 40
	p.AdvancedParameters.HeatTransfer.F_sh_hor = 1.0
	p.AdvancedParameters.HeatTransfer.F_sh_vert = 1.0
	p.AdvancedParameters.HeatTransfer.F_f = 0.0
	p.AdvancedParameters.HeatTransfer.F_w = 1.0

	lvl1 := &CalcLevel1{A_Calc_Window_Horizontal: 2, A_Calc_Window_South: 5, A_Calc_Window_North: 3}
	lvl6 := &CalcLevel6{ACalcWindowEast: 4, ACalcWindowWest: 4, GGlN: 0.6}
	c := &CalcLevel7{Lvl0: p, Lvl1: lvl1, Lvl6: lvl6}

	// QSolHor = 2*100*1*1*1*0.6 = 120
	if got := c.calcQSolHor(); !approxEqual(got, 120) {
		t.Errorf("QSolHor: got %.2f, want 120", got)
	}
	// QSolEast = 4*80*1*1*1*0.6 = 192
	if got := c.calcQSolEast(); !approxEqual(got, 192) {
		t.Errorf("QSolEast: got %.2f, want 192", got)
	}
	// QSolSouth = 5*120*1*1*1*0.6 = 360
	if got := c.calcQSolSouth(); !approxEqual(got, 360) {
		t.Errorf("QSolSouth: got %.2f, want 360", got)
	}
	// QSolWest = 4*80*1*1*1*0.6 = 192
	if got := c.calcQSolWest(); !approxEqual(got, 192) {
		t.Errorf("QSolWest: got %.2f, want 192", got)
	}
	// QSolNorth = 3*40*1*1*1*0.6 = 72
	if got := c.calcQSolNorth(); !approxEqual(got, 72) {
		t.Errorf("QSolNorth: got %.2f, want 72", got)
	}
}
