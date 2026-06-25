package calc

import "testing"

func TestCalcLevel2_calcACRef_inputTakesPriority(t *testing.T) {
	p := newTestParams()
	p.BasicParameters.Envelope.A_C_Ref_Input = 150
	lvl1 := &CalcLevel1{A_C_Ref_Estim: 100}
	c := &CalcLevel2{Lvl0: p, Lvl1: lvl1}
	if got := c.calcACRef(); got != 150 {
		t.Errorf("got %.1f, want 150", got)
	}
}

func TestCalcLevel2_calcACRef_fallsToEstim(t *testing.T) {
	p := newTestParams()
	lvl1 := &CalcLevel1{A_C_Ref_Estim: 100}
	c := &CalcLevel2{Lvl0: p, Lvl1: lvl1}
	if got := c.calcACRef(); got != 100 {
		t.Errorf("got %.1f, want 100", got)
	}
}

func TestCalcLevel2_calcNStoreyEffective(t *testing.T) {
	p := newTestParams()
	p.BasicParameters.BuildingAppearance.N_Storey = 3
	lvl1 := &CalcLevel1{F_AtticCond: 1.0, F_CellarCond: 0.5}
	c := &CalcLevel2{Lvl0: p, Lvl1: lvl1}
	// 0.7*1.0 + 3 + 0.5 = 4.2
	if got := c.calcNStoreyEffective(); !approxEqual(got, 4.2) {
		t.Errorf("got %.4f, want 4.2", got)
	}
}

func TestCalcLevel2_calcNStoreyEffectiveEnvelope(t *testing.T) {
	tests := []struct {
		name       string
		atticCode  string
		cellarCode string
		fAttic     float64
		fCellar    float64
		nStorey    int
		want       float64
	}{
		{
			name:       "attic suffix I overrides fAttic to 1",
			atticCode:  "PI", cellarCode: "N",
			fAttic: 0.5, fCellar: 0, nStorey: 2,
			want: 1*0.7 + 2 + 0, // attic "PI" ends in "I" → override to 1
		},
		{
			name:       "cellar suffix I overrides fCellar to 1",
			atticCode:  "N", cellarCode: "CI",
			fAttic: 0, fCellar: 1.0, nStorey: 2,
			want: 0*0.7 + 2 + 1, // cellar "CI" ends in "I" → cellar stays 1 (already 1 from func)
		},
		{
			name:       "no suffix I uses f values directly",
			atticCode:  "P", cellarCode: "P",
			fAttic: 0.5, fCellar: 0.5, nStorey: 1,
			want: 0.5*0.7 + 1 + 0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := newTestParams()
			p.BasicParameters.BuildingAppearance.Code_AtticCond = tt.atticCode
			p.BasicParameters.BuildingAppearance.Code_CellarCond = tt.cellarCode
			p.BasicParameters.BuildingAppearance.N_Storey = tt.nStorey
			lvl1 := &CalcLevel1{F_AtticCond: tt.fAttic, F_CellarCond: tt.fCellar}
			c := &CalcLevel2{Lvl0: p, Lvl1: lvl1}
			got := c.calcNStoreyEffectiveEnvelope()
			if !approxEqual(got, tt.want) {
				t.Errorf("got %.4f, want %.4f", got, tt.want)
			}
		})
	}
}

func TestCalcLevel2_calcCheckToBeAppliedFloorAreaExactToEstim(t *testing.T) {
	// sum == 0 → 1
	c := &CalcLevel2{Lvl1: &CalcLevel1{F_AtticCond: 0, F_CellarCond: 0}}
	if got := c.calcCheckToBeAppliedFloorAreaExactToEstim(); got != 1 {
		t.Errorf("expected 1, got %d", got)
	}
	// sum != 0 → 0
	c.Lvl1 = &CalcLevel1{F_AtticCond: 0.5, F_CellarCond: 0}
	if got := c.calcCheckToBeAppliedFloorAreaExactToEstim(); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}

func TestCalcLevel2_calcRMeasureRoof1_withPredefined(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.InsulationMeasures.D_Insulation_PredefinedMeasure_Roof_1 = 0.2
	p.AdvancedParameters.ThermalResistances.R_PredefinedMeasure_Roof_1 = 5.0
	lvl1 := &CalcLevel1{D_Insulation_Measure_Roof_1: 0.4}
	c := &CalcLevel2{Lvl0: p, Lvl1: lvl1}
	// (0.4/0.2) * 5.0 = 10.0
	if got := c.calcRMeasureRoof1(); !approxEqual(got, 10.0) {
		t.Errorf("got %.4f, want 10.0", got)
	}
}

func TestCalcLevel2_calcRMeasureRoof1_noPredefined(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.InsulationMeasures.D_Insulation_PredefinedMeasure_Roof_1 = 0
	p.AdvancedParameters.ThermalResistances.R_PredefinedMeasure_Roof_1 = 5.0
	lvl1 := &CalcLevel1{}
	c := &CalcLevel2{Lvl0: p, Lvl1: lvl1}
	// predefined == 0 → return R directly
	if got := c.calcRMeasureRoof1(); !approxEqual(got, 5.0) {
		t.Errorf("got %.4f, want 5.0", got)
	}
}

func TestCalcLevel2_calcRMeasureAllElements(t *testing.T) {
	// Exercise all remaining R_Measure elements via NewCalcLevel2 with non-zero predefined
	p := newTestParams()
	const rRef = 4.0
	const dPre = 0.1
	const dIn = 0.2 // ratio = 2 → R = 8.0
	im := p.AdvancedParameters.InsulationMeasures
	im.D_Insulation_PredefinedMeasure_Roof_2 = dPre
	im.D_Insulation_PredefinedMeasure_Wall_1 = dPre
	im.D_Insulation_PredefinedMeasure_Wall_2 = dPre
	im.D_Insulation_PredefinedMeasure_Wall_3 = dPre
	im.D_Insulation_PredefinedMeasure_Floor_1 = dPre
	im.D_Insulation_PredefinedMeasure_Floor_2 = dPre
	tr := p.AdvancedParameters.ThermalResistances
	tr.R_PredefinedMeasure_Roof_2 = rRef
	tr.R_PredefinedMeasure_Wall_1 = rRef
	tr.R_PredefinedMeasure_Wall_2 = rRef
	tr.R_PredefinedMeasure_Wall_3 = rRef
	tr.R_PredefinedMeasure_Floor_1 = rRef
	tr.R_PredefinedMeasure_Floor_2 = rRef

	lvl1 := &CalcLevel1{
		D_Insulation_Measure_Roof_2:  dIn,
		D_Insulation_Measure_Wall_1:  dIn,
		D_Insulation_Measure_Wall_2:  dIn,
		D_Insulation_Measure_Wall_3:  dIn,
		D_Insulation_Measure_Floor_1: dIn,
		D_Insulation_Measure_Floor_2: dIn,
	}
	c := NewCalcLevel2(p, lvl1)
	want := (dIn / dPre) * rRef // 8.0
	for _, got := range []float64{
		c.R_Measure_Roof_2, c.R_Measure_Wall_1, c.R_Measure_Wall_2,
		c.R_Measure_Wall_3, c.R_Measure_Floor_1, c.R_Measure_Floor_2,
	} {
		if !approxEqual(got, want) {
			t.Errorf("R_Measure element: got %.4f, want %.4f", got, want)
		}
	}
}

func TestCalcLevel2_calcUMeasureWindow1_replace(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.MeasureTypes.Code_MeasureType_Window_1 = "Replace"
	lvl1 := &CalcLevel1{R_Before_Window_1: 0.5, R_Measure_Window_1: 2.0}
	c := &CalcLevel2{Lvl0: p, Lvl1: lvl1}
	// Replace → rBefore = 0, denom = 2.0, result = 0.5
	if got := c.calcUMeasureWindow1(); !approxEqual(got, 0.5) {
		t.Errorf("got %.4f, want 0.5", got)
	}
}

func TestCalcLevel2_calcUMeasureWindow1_notReplace(t *testing.T) {
	p := newTestParams()
	lvl1 := &CalcLevel1{R_Before_Window_1: 0.5, R_Measure_Window_1: 0.5}
	c := &CalcLevel2{Lvl0: p, Lvl1: lvl1}
	// denom = 1.0, result = 1.0
	if got := c.calcUMeasureWindow1(); !approxEqual(got, 1.0) {
		t.Errorf("got %.4f, want 1.0", got)
	}
}

func TestCalcLevel2_calcUMeasureWindow1_zeroDenom(t *testing.T) {
	p := newTestParams()
	lvl1 := &CalcLevel1{R_Before_Window_1: 0, R_Measure_Window_1: 0}
	c := &CalcLevel2{Lvl0: p, Lvl1: lvl1}
	if got := c.calcUMeasureWindow1(); got != 0 {
		t.Errorf("got %.4f, want 0", got)
	}
}

func TestCalcLevel2_calcUMeasureWindow2AndDoor1(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.MeasureTypes.Code_MeasureType_Window_2 = "Replace"
	p.AdvancedParameters.MeasureTypes.Code_MeasureType_Door_1 = "Replace"
	lvl1 := &CalcLevel1{
		R_Before_Window_2: 1.0, R_Measure_Window_2: 3.0,
		R_Before_Door_1: 1.0, R_Measure_Door_1: 0,
	}
	c := &CalcLevel2{Lvl0: p, Lvl1: lvl1}
	// Window2 Replace: denom = 0+3 = 3, result = 1/3
	if got := c.calcUMeasureWindow2(); !approxEqual(got, 1.0/3) {
		t.Errorf("UMeasureWindow2: got %.4f, want %.4f", got, 1.0/3)
	}
	// Door1 Replace: denom = 0+0 = 0 → 0
	if got := c.calcUMeasureDoor1(); got != 0 {
		t.Errorf("UMeasureDoor1: got %.4f, want 0", got)
	}
}
