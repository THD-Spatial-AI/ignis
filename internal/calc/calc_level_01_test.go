package calc

import "testing"

func TestCalcLevel1_calcACRefEstim(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*CalcLevel1)
		want    float64
	}{
		{
			name: "uses A_C_IntDim when positive",
			setup: func(c *CalcLevel1) {
				c.Lvl0.BasicParameters.Envelope.A_C_IntDim = 100
			},
			want: 100,
		},
		{
			name: "uses 0.85*A_C_ExtDim when IntDim is zero",
			setup: func(c *CalcLevel1) {
				c.Lvl0.BasicParameters.Envelope.A_C_ExtDim = 100
			},
			want: 85,
		},
		{
			name: "uses 1.1*A_C_Living when ExtDim is also zero",
			setup: func(c *CalcLevel1) {
				c.Lvl0.BasicParameters.Envelope.A_C_Living = 100
			},
			want: 110,
		},
		{
			name: "uses 1.4*A_C_Use when Living is also zero",
			setup: func(c *CalcLevel1) {
				c.Lvl0.BasicParameters.Envelope.A_C_Use = 100
			},
			want: 140,
		},
		{
			name: "falls back to 0.85/3*V_C",
			setup: func(c *CalcLevel1) {
				c.Lvl0.BasicParameters.Envelope.V_C = 300
			},
			want: 85,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := newTestParams()
			c := &CalcLevel1{Lvl0: p}
			tt.setup(c)
			got := c.calcACRefEstim()
			if !approxEqual(got, tt.want) {
				t.Errorf("got %.4f, want %.4f", got, tt.want)
			}
		})
	}
}

func TestCalcLevel1_atticCondBranches(t *testing.T) {
	tests := []struct {
		code     string
		pRoof    float64
		qRoof    float64
		pCeil    float64
		qCeil    float64
		fAttic   float64
	}{
		{"C", 1.6, 15, 0, 0, 1.0},
		{"PI", 1.6, 15, 0, 0, 0.5},
		{"NI", 1.6, 15, 0, 0, 0.0},
		{"P", 0.8, 7, 0.6, 3, 0.5},
		{"-", 1.2, 5, 0, 0, 0.0},
		{"", 0, 0, 1.2, 5, 0.0},
	}

	for _, tt := range tests {
		t.Run("code="+tt.code, func(t *testing.T) {
			p := newTestParams()
			p.BasicParameters.BuildingAppearance.Code_AtticCond = tt.code
			c := NewCalcLevel1(p)
			if c.P_Roof != tt.pRoof {
				t.Errorf("P_Roof: got %.1f, want %.1f", c.P_Roof, tt.pRoof)
			}
			if c.Q_Roof != tt.qRoof {
				t.Errorf("Q_Roof: got %.1f, want %.1f", c.Q_Roof, tt.qRoof)
			}
			if c.P_Ceiling != tt.pCeil {
				t.Errorf("P_Ceiling: got %.1f, want %.1f", c.P_Ceiling, tt.pCeil)
			}
			if c.Q_Ceiling != tt.qCeil {
				t.Errorf("Q_Ceiling: got %.1f, want %.1f", c.Q_Ceiling, tt.qCeil)
			}
			if c.F_AtticCond != tt.fAttic {
				t.Errorf("F_AtticCond: got %.1f, want %.1f", c.F_AtticCond, tt.fAttic)
			}
		})
	}
}

func TestCalcLevel1_calcFCellarCond(t *testing.T) {
	tests := []struct {
		code string
		want float64
	}{
		{"C", 1.0},
		{"CI", 1.0},
		{"P", 0.5},
		{"PI", 0.5},
		{"NI", 0.0},
		{"", 0.0},
	}
	for _, tt := range tests {
		t.Run("code="+tt.code, func(t *testing.T) {
			p := newTestParams()
			p.BasicParameters.BuildingAppearance.Code_CellarCond = tt.code
			c := &CalcLevel1{Lvl0: p}
			got := c.calcFCellarCond()
			if got != tt.want {
				t.Errorf("got %.1f, want %.1f", got, tt.want)
			}
		})
	}
}

func TestCalcLevel1_calcFComplexRoof(t *testing.T) {
	tests := []struct{ code string; want float64 }{
		{"Simple", 0.9},
		{"Complex", 1.3},
		{"Standard", 1.0},
	}
	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			p := newTestParams()
			p.BasicParameters.BuildingAppearance.Code_ComplexRoof = tt.code
			c := &CalcLevel1{Lvl0: p}
			if got := c.calcFComplexRoof(); got != tt.want {
				t.Errorf("got %.1f, want %.1f", got, tt.want)
			}
		})
	}
}

func TestCalcLevel1_calcFComplexFootprint(t *testing.T) {
	tests := []struct{ code string; want float64 }{
		{"Simple", 0.9},
		{"Complex", 1.2},
		{"Standard", 1.0},
	}
	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			p := newTestParams()
			p.BasicParameters.BuildingAppearance.Code_ComplexFootprint = tt.code
			c := &CalcLevel1{Lvl0: p}
			if got := c.calcFComplexFootprint(); got != tt.want {
				t.Errorf("got %.1f, want %.1f", got, tt.want)
			}
		})
	}
}

func TestCalcLevel1_calcCodeEstimConstructionBorderWall(t *testing.T) {
	p := newTestParams()
	p.BasicParameters.BuildingAppearance.Code_CellarCond = "P"
	c := &CalcLevel1{Lvl0: p}
	if got := c.calcCodeEstimConstructionBorderWallToCellarOrSoil(); got != "Unh" {
		t.Errorf("expected Unh, got %s", got)
	}

	p.BasicParameters.BuildingAppearance.Code_CellarCond = "C"
	if got := c.calcCodeEstimConstructionBorderWallToCellarOrSoil(); got != "Soil" {
		t.Errorf("expected Soil, got %s", got)
	}
}

func TestCalcLevel1_calcCodeEstimConstructionBorderFloor(t *testing.T) {
	p := newTestParams()
	p.BasicParameters.BuildingAppearance.Code_CellarCond = "-"
	c := &CalcLevel1{Lvl0: p}
	if got := c.calcCodeEstimConstructionBorderFloor(); got != "Soil" {
		t.Errorf("expected Soil, got %s", got)
	}

	p.BasicParameters.BuildingAppearance.Code_CellarCond = "C"
	if got := c.calcCodeEstimConstructionBorderFloor(); got != "Cellar" {
		t.Errorf("expected Cellar, got %s", got)
	}
}

func TestCalcLevel1_calcAExactEnvSum(t *testing.T) {
	p := newTestParams()
	env := p.BasicParameters.Envelope
	env.A_Roof_1 = 10; env.A_Roof_2 = 10
	env.A_Wall_1 = 10; env.A_Wall_2 = 10; env.A_Wall_3 = 10
	env.A_Floor_1 = 10; env.A_Floor_2 = 10
	env.A_Window_1 = 10; env.A_Window_2 = 10
	env.A_Door_1 = 10
	c := &CalcLevel1{Lvl0: p}
	if got := c.calcAExactEnvSum(); got != 100 {
		t.Errorf("got %.1f, want 100", got)
	}
}

func TestCalcLevel1_calcACalcWindow_estimation(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.PredefinedCodes.Code_TypeIntake_EnvelopeArea = "Estimation"
	p.BasicParameters.Envelope.A_Window_2 = 20
	p.BasicParameters.Envelope.A_Window_Horizontal = 5
	p.BasicParameters.Envelope.A_Window_South = 8
	p.BasicParameters.Envelope.A_Window_North = 3
	c := &CalcLevel1{Lvl0: p}
	if got := c.calcACalcWindow2(); got != 0 {
		t.Errorf("calcACalcWindow2: got %.1f, want 0", got)
	}
	if got := c.calcACalcWindowHorizontal(); got != 0 {
		t.Errorf("calcACalcWindowHorizontal: got %.1f, want 0", got)
	}
	if got := c.calcACalcWindowSouth(); got != 0 {
		t.Errorf("calcACalcWindowSouth: got %.1f, want 0", got)
	}
	if got := c.calcACalcWindowNorth(); got != 0 {
		t.Errorf("calcACalcWindowNorth: got %.1f, want 0", got)
	}
}

func TestCalcLevel1_calcACalcWindow_notEstimation(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.PredefinedCodes.Code_TypeIntake_EnvelopeArea = "Exact"
	p.BasicParameters.Envelope.A_Window_2 = 20
	p.BasicParameters.Envelope.A_Window_Horizontal = 5
	p.BasicParameters.Envelope.A_Window_South = 8
	p.BasicParameters.Envelope.A_Window_North = 3
	c := &CalcLevel1{Lvl0: p}
	if got := c.calcACalcWindow2(); got != 20 {
		t.Errorf("calcACalcWindow2: got %.1f, want 20", got)
	}
	if got := c.calcACalcWindowHorizontal(); got != 5 {
		t.Errorf("calcACalcWindowHorizontal: got %.1f, want 5", got)
	}
	if got := c.calcACalcWindowSouth(); got != 8 {
		t.Errorf("calcACalcWindowSouth: got %.1f, want 8", got)
	}
	if got := c.calcACalcWindowNorth(); got != 3 {
		t.Errorf("calcACalcWindowNorth: got %.1f, want 3", got)
	}
}

func TestCalcLevel1_calcDInsulationMeasure_inputOverridesPredefined(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.ActualInsulation.D_Insulation_Input_Measure_Roof_1 = 0.2
	p.AdvancedParameters.InsulationMeasures.D_Insulation_PredefinedMeasure_Roof_1 = 0.1
	c := &CalcLevel1{Lvl0: p}
	if got := c.calcDInsulationMeasureRoof1(); got != 0.2 {
		t.Errorf("got %.2f, want 0.2", got)
	}
}

func TestCalcLevel1_calcDInsulationMeasure_fallsToPredefined(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.InsulationMeasures.D_Insulation_PredefinedMeasure_Roof_1 = 0.1
	// Input = 0 (zero value)
	c := &CalcLevel1{Lvl0: p}
	if got := c.calcDInsulationMeasureRoof1(); got != 0.1 {
		t.Errorf("got %.2f, want 0.1", got)
	}
}

func TestCalcLevel1_calcDInsulationMeasureAllElements(t *testing.T) {
	// Ensure all 7 element types exercise the fallback branch via NewCalcLevel1.
	p := newTestParams()
	p.AdvancedParameters.InsulationMeasures.D_Insulation_PredefinedMeasure_Roof_2 = 0.14
	p.AdvancedParameters.InsulationMeasures.D_Insulation_PredefinedMeasure_Wall_1 = 0.16
	p.AdvancedParameters.InsulationMeasures.D_Insulation_PredefinedMeasure_Wall_2 = 0.12
	p.AdvancedParameters.InsulationMeasures.D_Insulation_PredefinedMeasure_Wall_3 = 0.10
	p.AdvancedParameters.InsulationMeasures.D_Insulation_PredefinedMeasure_Floor_1 = 0.08
	p.AdvancedParameters.InsulationMeasures.D_Insulation_PredefinedMeasure_Floor_2 = 0.06

	p.AdvancedParameters.ActualInsulation.D_Insulation_Input_Measure_Roof_2 = 0.20
	p.AdvancedParameters.ActualInsulation.D_Insulation_Input_Measure_Wall_1 = 0.22
	p.AdvancedParameters.ActualInsulation.D_Insulation_Input_Measure_Wall_2 = 0.24
	p.AdvancedParameters.ActualInsulation.D_Insulation_Input_Measure_Wall_3 = 0.26
	p.AdvancedParameters.ActualInsulation.D_Insulation_Input_Measure_Floor_1 = 0.28
	p.AdvancedParameters.ActualInsulation.D_Insulation_Input_Measure_Floor_2 = 0.30

	c := NewCalcLevel1(p)
	// Input non-zero → override branches
	if c.D_Insulation_Measure_Roof_2 != 0.20 {
		t.Errorf("Roof2: got %.2f, want 0.20", c.D_Insulation_Measure_Roof_2)
	}
	if c.D_Insulation_Measure_Wall_1 != 0.22 {
		t.Errorf("Wall1: got %.2f", c.D_Insulation_Measure_Wall_1)
	}
}

func TestCalcLevel1_calcRBeforeWindow_nonZeroU(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.Uvalues.U_Window_1 = 2.0
	p.AdvancedParameters.Uvalues.U_Window_2 = 4.0
	p.AdvancedParameters.Uvalues.U_Door_1 = 1.0
	c := &CalcLevel1{Lvl0: p}
	if got := c.calcRBeforeWindow1(); !approxEqual(got, 0.5) {
		t.Errorf("RBeforeWindow1: got %.4f, want 0.5", got)
	}
	if got := c.calcRBeforeWindow2(); !approxEqual(got, 0.25) {
		t.Errorf("RBeforeWindow2: got %.4f, want 0.25", got)
	}
	if got := c.calcRBeforeDoor1(); !approxEqual(got, 1.0) {
		t.Errorf("RBeforeDoor1: got %.4f, want 1.0", got)
	}
}

func TestCalcLevel1_calcRBeforeWindow_zeroU(t *testing.T) {
	p := newTestParams()
	// U = 0 (zero values)
	c := &CalcLevel1{Lvl0: p}
	if got := c.calcRBeforeWindow1(); got != 0 {
		t.Errorf("RBeforeWindow1: got %.4f, want 0", got)
	}
	if got := c.calcRBeforeWindow2(); got != 0 {
		t.Errorf("RBeforeWindow2: got %.4f, want 0", got)
	}
	if got := c.calcRBeforeDoor1(); got != 0 {
		t.Errorf("RBeforeDoor1: got %.4f, want 0", got)
	}
}

func TestCalcLevel1_calcHVentilation(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.AirInfiltration.N_air_use = 0.5
	p.AdvancedParameters.AirInfiltration.N_air_infiltration = 0.1
	p.BasicParameters.BuildingAppearance.H_room = 2.5
	c := &CalcLevel1{Lvl0: p}
	// 0.34 * (0.5 + 0.1) * 2.5 = 0.34 * 0.6 * 2.5 = 0.51
	if got := c.calcHVentilation(); !approxEqual(got, 0.51) {
		t.Errorf("got %.6f, want 0.51", got)
	}
}

func TestCalcLevel1_calcSumDeltaTForHeatingDays(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.ClimateConditions.Theta_i = 20
	p.AdvancedParameters.ClimateConditions.Theta_e = -5
	p.AdvancedParameters.ClimateConditions.HeatingDays = 200
	c := &CalcLevel1{Lvl0: p}
	// (20 - (-5)) * 200 = 5000
	if got := c.calcSumDeltaTForHeatingDays(); got != 5000 {
		t.Errorf("got %.1f, want 5000", got)
	}
}

func TestCalcLevel1_calcQInt(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.HeatTransfer.Phi_int = 5.0
	p.AdvancedParameters.ClimateConditions.HeatingDays = 200
	c := &CalcLevel1{Lvl0: p}
	// 5.0 * 0.024 * 200 = 24
	if got := c.calcQInt(); !approxEqual(got, 24) {
		t.Errorf("got %.4f, want 24", got)
	}
}

func TestCalcLevel1_passthroughMethods(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.ThermalResistances.R_PredefinedMeasure_Window_1 = 1.5
	p.AdvancedParameters.ThermalResistances.R_PredefinedMeasure_Window_2 = 2.0
	p.AdvancedParameters.ThermalResistances.R_PredefinedMeasure_Door_1 = 0.8
	p.AdvancedParameters.SolarTransmittance.G_gl_n_PredefinedMeasure_Window_1 = 0.6
	p.AdvancedParameters.SolarTransmittance.G_gl_n_PredefinedMeasure_Window_2 = 0.5
	c := NewCalcLevel1(p)
	if c.R_Measure_Window_1 != 1.5 {
		t.Errorf("R_Measure_Window_1: got %.2f, want 1.5", c.R_Measure_Window_1)
	}
	if c.R_Measure_Window_2 != 2.0 {
		t.Errorf("R_Measure_Window_2: got %.2f, want 2.0", c.R_Measure_Window_2)
	}
	if c.R_Measure_Door_1 != 0.8 {
		t.Errorf("R_Measure_Door_1: got %.2f, want 0.8", c.R_Measure_Door_1)
	}
	if c.G_gl_n_Measure_Window_1 != 0.6 {
		t.Errorf("G_gl_n_Measure_Window_1: got %.2f", c.G_gl_n_Measure_Window_1)
	}
	if c.G_gl_n_Measure_Window_2 != 0.5 {
		t.Errorf("G_gl_n_Measure_Window_2: got %.2f", c.G_gl_n_Measure_Window_2)
	}
}
