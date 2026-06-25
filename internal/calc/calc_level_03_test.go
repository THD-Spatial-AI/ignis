package calc

import (
	"math"
	"testing"
)

func TestCalcLevel3_calcACStorey_nonZeroDenom(t *testing.T) {
	lvl2 := &CalcLevel2{A_C_Ref: 100, N_Storey_effective: 2.0}
	c := &CalcLevel3{Lvl2: lvl2}
	if got := c.calcACStorey(); !approxEqual(got, 50) {
		t.Errorf("got %.4f, want 50", got)
	}
}

func TestCalcLevel3_calcACStorey_zeroDenom(t *testing.T) {
	lvl2 := &CalcLevel2{A_C_Ref: 100, N_Storey_effective: 0}
	c := &CalcLevel3{Lvl2: lvl2}
	if got := c.calcACStorey(); got != 0 {
		t.Errorf("got %.4f, want 0", got)
	}
}

func TestCalcLevel3_calcVEstimC(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.PredefinedCodes.F_Corr_CeilingHeight = 1.0
	lvl2 := &CalcLevel2{A_C_Ref: 100}
	c := &CalcLevel3{Lvl0: p, Lvl2: lvl2}
	// math.Round((3/0.85)*10)/10 = 3.5
	want := math.Round((3.0/0.85)*10) / 10 * 1.0 * 100
	if got := c.calcVEstimC(); !approxEqual(got, want) {
		t.Errorf("got %.4f, want %.4f", got, want)
	}
}

func TestCalcLevel3_calcAEstimDoor(t *testing.T) {
	lvl2 := &CalcLevel2{A_C_Ref: 100}
	c := &CalcLevel3{Lvl2: lvl2}
	// 0.01*100 + 1.5 = 2.5
	if got := c.calcAEstimDoor(); !approxEqual(got, 2.5) {
		t.Errorf("got %.4f, want 2.5", got)
	}
}

func TestCalcLevel3_calcRBeforeRoof1_withU(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.Uvalues.U_Roof_1 = 2.0
	p.AdvancedParameters.AdditionalResistances.R_Add_UnheatedSpace_Roof_1 = 0.1
	lvl2 := &CalcLevel2{}
	c := &CalcLevel3{Lvl0: p, Lvl2: lvl2}
	// 0.1 + 1/2.0 - 0 = 0.6
	if got := c.calcRBeforeRoof1(); !approxEqual(got, 0.6) {
		t.Errorf("got %.4f, want 0.6", got)
	}
}

func TestCalcLevel3_calcRBeforeRoof1_replaceInsulation(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.Uvalues.U_Roof_1 = 2.0
	p.AdvancedParameters.MeasureTypes.Code_MeasureType_Roof_1 = "ReplaceInsulation"
	p.AdvancedParameters.Insulation.D_Insulation_Roof_1 = 0.08 // 0.08/0.04 = 2.0
	lvl2 := &CalcLevel2{}
	c := &CalcLevel3{Lvl0: p, Lvl2: lvl2}
	// 0 + 1/2.0 - 2.0 = -1.5
	if got := c.calcRBeforeRoof1(); !approxEqual(got, -1.5) {
		t.Errorf("got %.4f, want -1.5", got)
	}
}

func TestCalcLevel3_calcRBeforeRoof1_zeroU(t *testing.T) {
	p := newTestParams()
	// U_Roof_1 = 0 (zero value)
	lvl2 := &CalcLevel2{}
	c := &CalcLevel3{Lvl0: p, Lvl2: lvl2}
	if got := c.calcRBeforeRoof1(); got != 0 {
		t.Errorf("got %.4f, want 0", got)
	}
}

func TestCalcLevel3_calcRBeforeRoof2(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.Uvalues.U_Roof_2 = 4.0
	p.AdvancedParameters.MeasureTypes.Code_MeasureType_Roof_2 = "ReplaceInsulation"
	p.AdvancedParameters.Insulation.D_Insulation_Roof_2 = 0.04 // /0.04 = 1.0
	lvl2 := &CalcLevel2{}
	c := &CalcLevel3{Lvl0: p, Lvl2: lvl2}
	// 0 + 1/4 - 1.0 = -0.75
	if got := c.calcRBeforeRoof2(); !approxEqual(got, -0.75) {
		t.Errorf("got %.4f, want -0.75", got)
	}
}

func TestCalcLevel3_calcRBeforeWallsAndFloors(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.Uvalues.U_Wall_1 = 1.0
	p.AdvancedParameters.Uvalues.U_Wall_2 = 1.0
	p.AdvancedParameters.Uvalues.U_Wall_3 = 1.0
	p.AdvancedParameters.Uvalues.U_Floor_1 = 1.0
	p.AdvancedParameters.Uvalues.U_Floor_2 = 1.0
	// No ReplaceInsulation → rMeasure = 0
	lvl2 := &CalcLevel2{}
	c := &CalcLevel3{Lvl0: p, Lvl2: lvl2}
	for name, got := range map[string]float64{
		"Wall1":  c.calcRBeforeWall1(),
		"Wall2":  c.calcRBeforeWall2(),
		"Wall3":  c.calcRBeforeWall3(),
		"Floor1": c.calcRBeforeFloor1(),
		"Floor2": c.calcRBeforeFloor2(),
	} {
		if !approxEqual(got, 1.0) {
			t.Errorf("%s: got %.4f, want 1.0", name, got)
		}
	}
}

func TestCalcLevel3_calcRBeforeWall_replaceInsulation(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.Uvalues.U_Wall_1 = 0 // zero U → uWall=0
	p.AdvancedParameters.MeasureTypes.Code_MeasureType_Wall_1 = "ReplaceInsulation"
	p.AdvancedParameters.Insulation.D_Insulation_Wall_1 = 0.04 // /0.04 = 1.0
	lvl2 := &CalcLevel2{}
	c := &CalcLevel3{Lvl0: p, Lvl2: lvl2}
	// 0 + 0 - 1.0 = -1.0
	if got := c.calcRBeforeWall1(); !approxEqual(got, -1.0) {
		t.Errorf("got %.4f, want -1.0", got)
	}
}

func TestCalcLevel3_calcRBeforeFloor_replaceInsulation(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.Uvalues.U_Floor_1 = 0
	p.AdvancedParameters.MeasureTypes.Code_MeasureType_Floor_1 = "ReplaceInsulation"
	p.AdvancedParameters.Insulation.D_Insulation_Floor_1 = 0.08 // /0.04 = 2.0
	lvl2 := &CalcLevel2{}
	c := &CalcLevel3{Lvl0: p, Lvl2: lvl2}
	if got := c.calcRBeforeFloor1(); !approxEqual(got, -2.0) {
		t.Errorf("got %.4f, want -2.0", got)
	}
}

func TestCalcLevel3_calcRBeforeFloor2_replaceInsulation(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.Uvalues.U_Floor_2 = 0
	p.AdvancedParameters.MeasureTypes.Code_MeasureType_Floor_2 = "ReplaceInsulation"
	p.AdvancedParameters.Insulation.D_Insulation_Floor_2 = 0.04
	lvl2 := &CalcLevel2{}
	c := &CalcLevel3{Lvl0: p, Lvl2: lvl2}
	if got := c.calcRBeforeFloor2(); !approxEqual(got, -1.0) {
		t.Errorf("got %.4f, want -1.0", got)
	}
}

func TestCalcLevel3_calcUActualWindows(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.Uvalues.U_Window_1 = 2.0
	p.AdvancedParameters.Uvalues.U_Window_2 = 3.0
	p.AdvancedParameters.Uvalues.U_Door_1 = 1.0
	p.AdvancedParameters.MeasureFractions.F_Measure_Window_1 = 0.5
	p.AdvancedParameters.MeasureFractions.F_Measure_Window_2 = 0.25
	p.AdvancedParameters.MeasureFractions.F_Measure_Door_1 = 0.0
	lvl2 := &CalcLevel2{U_Measure_Window_1: 1.0, U_Measure_Window_2: 2.0, U_Measure_Door_1: 0.5}
	c := &CalcLevel3{Lvl0: p, Lvl2: lvl2}
	// Window1: (1-0.5)*2.0 + 0.5*1.0 = 1.5
	if got := c.calcUActualWindow1(); !approxEqual(got, 1.5) {
		t.Errorf("UActualWindow1: got %.4f, want 1.5", got)
	}
	// Window2: (1-0.25)*3.0 + 0.25*2.0 = 2.75
	if got := c.calcUActualWindow2(); !approxEqual(got, 2.75) {
		t.Errorf("UActualWindow2: got %.4f, want 2.75", got)
	}
	// Door1: (1-0)*1.0 + 0*0.5 = 1.0
	if got := c.calcUActualDoor1(); !approxEqual(got, 1.0) {
		t.Errorf("UActualDoor1: got %.4f, want 1.0", got)
	}
}
