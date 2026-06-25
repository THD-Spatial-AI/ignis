package calc

import (
	"github.com/thd-spatial-ai/ignis/internal/models"
	"strings"
)

// CalcLevel2 holds calculated values for level 2
type CalcLevel2 struct {
	Lvl0 *models.TabulaBuildingParameters
	Lvl1 *CalcLevel1

	// Calculated attributes
	A_C_Ref                                  float64 `json:"A_C_Ref"`
	N_Storey_effective                       float64 `json:"n_Storey_effective"`
	N_Storey_effective_envelope              float64 `json:"n_Storey_effective_envelope"`
	Check_ToBeApplied_FloorArea_ExactToEstim int     `json:"Check_ToBeApplied_FloorArea_ExactToEstim"`
	R_Measure_Roof_1                         float64 `json:"R_Measure_Roof_1"`
	R_Measure_Roof_2                         float64 `json:"R_Measure_Roof_2"`
	R_Measure_Wall_1                         float64 `json:"R_Measure_Wall_1"`
	R_Measure_Wall_2                         float64 `json:"R_Measure_Wall_2"`
	R_Measure_Wall_3                         float64 `json:"R_Measure_Wall_3"`
	R_Measure_Floor_1                        float64 `json:"R_Measure_Floor_1"`
	R_Measure_Floor_2                        float64 `json:"R_Measure_Floor_2"`
	U_Measure_Window_1                       float64 `json:"U_Measure_Window_1"`
	U_Measure_Window_2                       float64 `json:"U_Measure_Window_2"`
	U_Measure_Door_1                         float64 `json:"U_Measure_Door_1"`
}

// NewCalcLevel2 creates a new CalcLevel2 instance and runs all calculations
func NewCalcLevel2(lvl0 *models.TabulaBuildingParameters, lvl1 *CalcLevel1) *CalcLevel2 {
	calc := &CalcLevel2{
		Lvl0: lvl0,
		Lvl1: lvl1,
	}
	calc.Run()
	return calc
}

// Run executes all calculation methods in CalcLevel2 and stores output in corresponding attributes
func (c *CalcLevel2) Run() {
	c.A_C_Ref = c.calcACRef()
	c.N_Storey_effective = c.calcNStoreyEffective()
	c.N_Storey_effective_envelope = c.calcNStoreyEffectiveEnvelope()
	c.Check_ToBeApplied_FloorArea_ExactToEstim = c.calcCheckToBeAppliedFloorAreaExactToEstim()
	c.R_Measure_Roof_1 = c.calcRMeasureRoof1()
	c.R_Measure_Roof_2 = c.calcRMeasureRoof2()
	c.R_Measure_Wall_1 = c.calcRMeasureWall1()
	c.R_Measure_Wall_2 = c.calcRMeasureWall2()
	c.R_Measure_Wall_3 = c.calcRMeasureWall3()
	c.R_Measure_Floor_1 = c.calcRMeasureFloor1()
	c.R_Measure_Floor_2 = c.calcRMeasureFloor2()
	c.U_Measure_Window_1 = c.calcUMeasureWindow1()
	c.U_Measure_Window_2 = c.calcUMeasureWindow2()
	c.U_Measure_Door_1 = c.calcUMeasureDoor1()
}

// calcACRef calculates A_C_Ref
// actually measured by applying the TABULA definition, if available; otherwise estimated by applying standard conversion factors
// Formula: IF(A_C_Ref_Input > 0, A_C_Ref_Input, A_C_Ref_Estim)
func (c *CalcLevel2) calcACRef() float64 {
	if c.Lvl0.BasicParameters.Envelope.A_C_Ref_Input > 0 {
		return c.Lvl0.BasicParameters.Envelope.A_C_Ref_Input
	}
	return c.Lvl1.A_C_Ref_Estim
}

// calcNStoreyEffective calculates n_Storey_effective
// Formula: 0.7 * f_AtticCond + n_Storey + f_CellarCond
func (c *CalcLevel2) calcNStoreyEffective() float64 {
	return 0.7*c.Lvl1.F_AtticCond + float64(c.Lvl0.BasicParameters.BuildingAppearance.N_Storey) + c.Lvl1.F_CellarCond
}

// calcNStoreyEffectiveEnvelope calculates n_Storey_effective_envelope
// Formula: IF(RIGHT($Code_AtticCond,1)="I",1,f_AtticCond)*0.7+n_Storey+IF(RIGHT($Code_CellarCond,1)="I",1,f_CellarCond)
func (c *CalcLevel2) calcNStoreyEffectiveEnvelope() float64 {
	atticCond := c.Lvl1.F_AtticCond
	if strings.HasSuffix(c.Lvl0.BasicParameters.BuildingAppearance.Code_AtticCond, "I") {
		atticCond = 1
	}

	cellarCond := c.Lvl1.F_CellarCond
	if strings.HasSuffix(c.Lvl0.BasicParameters.BuildingAppearance.Code_CellarCond, "I") {
		cellarCond = 1
	}

	return atticCond*0.7 + float64(c.Lvl0.BasicParameters.BuildingAppearance.N_Storey) + cellarCond
}

// calcCheckToBeAppliedFloorAreaExactToEstim calculates Check_ToBeApplied_FloorArea_ExactToEstim
// Formula: IF(f_AtticCond+f_CellarCond=0,1,0)
func (c *CalcLevel2) calcCheckToBeAppliedFloorAreaExactToEstim() int {
	if c.Lvl1.F_AtticCond+c.Lvl1.F_CellarCond == 0 {
		return 1
	}
	return 0
}

// calcRMeasureRoof1 calculates R_Measure_Roof_1
// element type roof 1
// Unit: m²K/W
// Formula: IFERROR(IF(d_Insulation_PredefinedMeasure_Roof_1 <> 0, d_Insulation_Measure_Roof_1 / d_Insulation_PredefinedMeasure_Roof_1, 1) * R_PredefinedMeasure_Roof_1, 0)
func (c *CalcLevel2) calcRMeasureRoof1() float64 {
	if c.Lvl0.AdvancedParameters.InsulationMeasures.D_Insulation_PredefinedMeasure_Roof_1 != 0 {
		return (c.Lvl1.D_Insulation_Measure_Roof_1 / c.Lvl0.AdvancedParameters.InsulationMeasures.D_Insulation_PredefinedMeasure_Roof_1) *
			c.Lvl0.AdvancedParameters.ThermalResistances.R_PredefinedMeasure_Roof_1
	}
	return c.Lvl0.AdvancedParameters.ThermalResistances.R_PredefinedMeasure_Roof_1
}

// calcRMeasureRoof2 calculates R_Measure_Roof_2
// element type roof 2
// Unit: m²K/W
// Formula: IFERROR(IF(d_Insulation_PredefinedMeasure_Roof_2 <> 0, d_Insulation_Measure_Roof_2 / d_Insulation_PredefinedMeasure_Roof_2, 1) * R_PredefinedMeasure_Roof_2, 0)
func (c *CalcLevel2) calcRMeasureRoof2() float64 {
	if c.Lvl0.AdvancedParameters.InsulationMeasures.D_Insulation_PredefinedMeasure_Roof_2 != 0 {
		return (c.Lvl1.D_Insulation_Measure_Roof_2 / c.Lvl0.AdvancedParameters.InsulationMeasures.D_Insulation_PredefinedMeasure_Roof_2) *
			c.Lvl0.AdvancedParameters.ThermalResistances.R_PredefinedMeasure_Roof_2
	}
	return c.Lvl0.AdvancedParameters.ThermalResistances.R_PredefinedMeasure_Roof_2
}

// calcRMeasureWall1 calculates R_Measure_Wall_1
// element type wall 1
// Unit: m²K/W
// Formula: IFERROR(IF(d_Insulation_PredefinedMeasure_Wall_1 <> 0, d_Insulation_Measure_Wall_1 / d_Insulation_PredefinedMeasure_Wall_1, 1) * R_PredefinedMeasure_Wall_1, 0)
func (c *CalcLevel2) calcRMeasureWall1() float64 {
	if c.Lvl0.AdvancedParameters.InsulationMeasures.D_Insulation_PredefinedMeasure_Wall_1 != 0 {
		return (c.Lvl1.D_Insulation_Measure_Wall_1 / c.Lvl0.AdvancedParameters.InsulationMeasures.D_Insulation_PredefinedMeasure_Wall_1) *
			c.Lvl0.AdvancedParameters.ThermalResistances.R_PredefinedMeasure_Wall_1
	}
	return c.Lvl0.AdvancedParameters.ThermalResistances.R_PredefinedMeasure_Wall_1
}

// calcRMeasureWall2 calculates R_Measure_Wall_2
// element type wall 2
// Unit: m²K/W
// Formula: IFERROR(IF(d_Insulation_PredefinedMeasure_Wall_2 <> 0, d_Insulation_Measure_Wall_2 / d_Insulation_PredefinedMeasure_Wall_2, 1) * R_PredefinedMeasure_Wall_2, 0)
func (c *CalcLevel2) calcRMeasureWall2() float64 {
	if c.Lvl0.AdvancedParameters.InsulationMeasures.D_Insulation_PredefinedMeasure_Wall_2 != 0 {
		return (c.Lvl1.D_Insulation_Measure_Wall_2 / c.Lvl0.AdvancedParameters.InsulationMeasures.D_Insulation_PredefinedMeasure_Wall_2) *
			c.Lvl0.AdvancedParameters.ThermalResistances.R_PredefinedMeasure_Wall_2
	}
	return c.Lvl0.AdvancedParameters.ThermalResistances.R_PredefinedMeasure_Wall_2
}

// calcRMeasureWall3 calculates R_Measure_Wall_3
// element type wall 3
// Unit: m²K/W
// Formula: IFERROR(IF(d_Insulation_PredefinedMeasure_Wall_3 <> 0, d_Insulation_Measure_Wall_3 / d_Insulation_PredefinedMeasure_Wall_3, 1) * R_PredefinedMeasure_Wall_3, 0)
func (c *CalcLevel2) calcRMeasureWall3() float64 {
	if c.Lvl0.AdvancedParameters.InsulationMeasures.D_Insulation_PredefinedMeasure_Wall_3 != 0 {
		return (c.Lvl1.D_Insulation_Measure_Wall_3 / c.Lvl0.AdvancedParameters.InsulationMeasures.D_Insulation_PredefinedMeasure_Wall_3) *
			c.Lvl0.AdvancedParameters.ThermalResistances.R_PredefinedMeasure_Wall_3
	}
	return c.Lvl0.AdvancedParameters.ThermalResistances.R_PredefinedMeasure_Wall_3
}

// calcRMeasureFloor1 calculates R_Measure_Floor_1
// element type floor 1
// Unit: m²K/W
// Formula: IFERROR(IF(d_Insulation_PredefinedMeasure_Floor_1 <> 0, d_Insulation_Measure_Floor_1 / d_Insulation_PredefinedMeasure_Floor_1, 1) * R_PredefinedMeasure_Floor_1, 0)
func (c *CalcLevel2) calcRMeasureFloor1() float64 {
	if c.Lvl0.AdvancedParameters.InsulationMeasures.D_Insulation_PredefinedMeasure_Floor_1 != 0 {
		return (c.Lvl1.D_Insulation_Measure_Floor_1 / c.Lvl0.AdvancedParameters.InsulationMeasures.D_Insulation_PredefinedMeasure_Floor_1) *
			c.Lvl0.AdvancedParameters.ThermalResistances.R_PredefinedMeasure_Floor_1
	}
	return c.Lvl0.AdvancedParameters.ThermalResistances.R_PredefinedMeasure_Floor_1
}

// calcRMeasureFloor2 calculates R_Measure_Floor_2
// element type floor 2
// Unit: m²K/W
// Formula: IFERROR(IF(d_Insulation_PredefinedMeasure_Floor_2 <> 0, d_Insulation_Measure_Floor_2 / d_Insulation_PredefinedMeasure_Floor_2, 1) * R_PredefinedMeasure_Floor_2, 0)
func (c *CalcLevel2) calcRMeasureFloor2() float64 {
	if c.Lvl0.AdvancedParameters.InsulationMeasures.D_Insulation_PredefinedMeasure_Floor_2 != 0 {
		return (c.Lvl1.D_Insulation_Measure_Floor_2 / c.Lvl0.AdvancedParameters.InsulationMeasures.D_Insulation_PredefinedMeasure_Floor_2) *
			c.Lvl0.AdvancedParameters.ThermalResistances.R_PredefinedMeasure_Floor_2
	}
	return c.Lvl0.AdvancedParameters.ThermalResistances.R_PredefinedMeasure_Floor_2
}

// calcUMeasureWindow1 calculates U_Measure_Window_1
// element type window 1
// Unit: W/(m²K)
// Formula: IFERROR(1 / (IF(Code_MeasureType_Window_1 = "Replace", 0, R_Before_Window_1) + IF(ISNUMBER(R_Measure_Window_1), R_Measure_Window_1, 0)), 0)
func (c *CalcLevel2) calcUMeasureWindow1() float64 {
	rBefore := c.Lvl1.R_Before_Window_1
	if c.Lvl0.AdvancedParameters.MeasureTypes.Code_MeasureType_Window_1 == "Replace" {
		rBefore = 0
	}

	denominator := rBefore + c.Lvl1.R_Measure_Window_1
	if denominator != 0 {
		return 1 / denominator
	}
	return 0
}

// calcUMeasureWindow2 calculates U_Measure_Window_2
// element type window 2
// Unit: W/(m²K)
// Formula: IFERROR(1 / (IF(Code_MeasureType_Window_2 = "Replace", 0, R_Before_Window_2) + IF(ISNUMBER(R_Measure_Window_2), R_Measure_Window_2, 0)), 0)
func (c *CalcLevel2) calcUMeasureWindow2() float64 {
	rBefore := c.Lvl1.R_Before_Window_2
	if c.Lvl0.AdvancedParameters.MeasureTypes.Code_MeasureType_Window_2 == "Replace" {
		rBefore = 0
	}

	denominator := rBefore + c.Lvl1.R_Measure_Window_2
	if denominator != 0 {
		return 1 / denominator
	}
	return 0
}

// calcUMeasureDoor1 calculates U_Measure_Door_1
// element type door 1
// Unit: W/(m²K)
// Formula: IFERROR(1 / (IF(Code_MeasureType_Door_1 = "Replace", 0, R_Before_Door_1) + IF(ISNUMBER(R_Measure_Door_1), R_Measure_Door_1, 0)), 0)
func (c *CalcLevel2) calcUMeasureDoor1() float64 {
	rBefore := c.Lvl1.R_Before_Door_1
	if c.Lvl0.AdvancedParameters.MeasureTypes.Code_MeasureType_Door_1 == "Replace" {
		rBefore = 0
	}

	denominator := rBefore + c.Lvl1.R_Measure_Door_1
	if denominator != 0 {
		return 1 / denominator
	}
	return 0
}
