package calc

import (
	"math"

	"github.com/thd-spatial-ai/ignis/internal/models"
)

// CalcLevel3 holds calculated values for level 3
type CalcLevel3 struct {
	Lvl0 *models.TabulaBuildingParameters
	Lvl2 *CalcLevel2

	// Calculated attributes
	A_C_Storey        float64 `json:"A_C_Storey"`
	V_Estim_C         float64 `json:"V_Estim_C"`
	A_Estim_Door      float64 `json:"A_Estim_Door"`
	R_Before_Roof_1   float64 `json:"R_Before_Roof_1"`
	R_Before_Roof_2   float64 `json:"R_Before_Roof_2"`
	R_Before_Wall_1   float64 `json:"R_Before_Wall_1"`
	R_Before_Wall_2   float64 `json:"R_Before_Wall_2"`
	R_Before_Wall_3   float64 `json:"R_Before_Wall_3"`
	R_Before_Floor_1  float64 `json:"R_Before_Floor_1"`
	R_Before_Floor_2  float64 `json:"R_Before_Floor_2"`
	U_Actual_Window_1 float64 `json:"U_Actual_Window_1"`
	U_Actual_Window_2 float64 `json:"U_Actual_Window_2"`
	U_Actual_Door_1   float64 `json:"U_Actual_Door_1"`
}

// NewCalcLevel3 creates a new CalcLevel3 instance and runs all calculations
func NewCalcLevel3(lvl0 *models.TabulaBuildingParameters, lvl2 *CalcLevel2) *CalcLevel3 {
	calc := &CalcLevel3{
		Lvl0: lvl0,
		Lvl2: lvl2,
	}
	calc.Run()
	return calc
}

// Run executes all calculation methods in CalcLevel3 and stores output in corresponding attributes
func (c *CalcLevel3) Run() {
	c.A_C_Storey = c.calcACStorey()
	c.V_Estim_C = c.calcVEstimC()
	c.A_Estim_Door = c.calcAEstimDoor()
	c.R_Before_Roof_1 = c.calcRBeforeRoof1()
	c.R_Before_Roof_2 = c.calcRBeforeRoof2()
	c.R_Before_Wall_1 = c.calcRBeforeWall1()
	c.R_Before_Wall_2 = c.calcRBeforeWall2()
	c.R_Before_Wall_3 = c.calcRBeforeWall3()
	c.R_Before_Floor_1 = c.calcRBeforeFloor1()
	c.R_Before_Floor_2 = c.calcRBeforeFloor2()
	c.U_Actual_Window_1 = c.calcUActualWindow1()
	c.U_Actual_Window_2 = c.calcUActualWindow2()
	c.U_Actual_Door_1 = c.calcUActualDoor1()
}

// calcACStorey calculates A_C_Storey
// Calculates the effective area of the storey
// Formula: IFERROR(A_C_Ref / n_Storey_effective, 0)
func (c *CalcLevel3) calcACStorey() float64 {
	if c.Lvl2.N_Storey_effective == 0 {
		return 0
	}
	return c.Lvl2.A_C_Ref / c.Lvl2.N_Storey_effective
}

// calcVEstimC calculates V_Estim_C
// Estimates the volume correction based on ceiling height
// Unit: m³
// Formula: ROUND(3/0.85, 1) * f_Corr_CeilingHeight * A_C_Ref
func (c *CalcLevel3) calcVEstimC() float64 {
	return math.Round((3.0/0.85)*10) / 10 * c.Lvl0.AdvancedParameters.PredefinedCodes.F_Corr_CeilingHeight * c.Lvl2.A_C_Ref
}

// calcAEstimDoor calculates A_Estim_Door
// Estimates the area of the door
// Unit: m²
// Formula: 0.01 * A_C_Ref + 1.5
func (c *CalcLevel3) calcAEstimDoor() float64 {
	return 0.01*c.Lvl2.A_C_Ref + 1.5
}

// calcRBeforeRoof1 calculates R_Before_Roof_1
// Element type roof 1
// Unit: m²*K/W
// Formula: IFERROR(IF(ISNUMBER(R_Add_UnheatedSpace_Roof_1), R_Add_UnheatedSpace_Roof_1, 0) + IF(ISNUMBER(U_Roof_1), 1/U_Roof_1 - IF(AND(Code_MeasureType_Roof_1="ReplaceInsulation", NOT(ISERROR(R_Measure_Roof_1))), d_Insulation_Roof_1/0.04, 0), 0), 0)
func (c *CalcLevel3) calcRBeforeRoof1() float64 {
	rAdd := c.Lvl0.AdvancedParameters.AdditionalResistances.R_Add_UnheatedSpace_Roof_1

	var uRoof float64
	if c.Lvl0.AdvancedParameters.Uvalues.U_Roof_1 != 0 {
		uRoof = 1 / c.Lvl0.AdvancedParameters.Uvalues.U_Roof_1
	}

	var rMeasure float64
	if c.Lvl0.AdvancedParameters.MeasureTypes.Code_MeasureType_Roof_1 == "ReplaceInsulation" {
		rMeasure = c.Lvl0.AdvancedParameters.Insulation.D_Insulation_Roof_1 / 0.04
	}

	return rAdd + uRoof - rMeasure
}

// calcRBeforeRoof2 calculates R_Before_Roof_2
// Element type roof 2
// Unit: m²*K/W
// Formula: IFERROR(IF(ISNUMBER(R_Add_UnheatedSpace_Roof_2), R_Add_UnheatedSpace_Roof_2, 0) + IF(ISNUMBER(U_Roof_2), 1/U_Roof_2 - IF(AND(Code_MeasureType_Roof_2="ReplaceInsulation", NOT(ISERROR(R_Measure_Roof_2))), d_Insulation_Roof_2/0.04, 0), 0), 0)
func (c *CalcLevel3) calcRBeforeRoof2() float64 {
	rAdd := c.Lvl0.AdvancedParameters.AdditionalResistances.R_Add_UnheatedSpace_Roof_2

	var uRoof float64
	if c.Lvl0.AdvancedParameters.Uvalues.U_Roof_2 != 0 {
		uRoof = 1 / c.Lvl0.AdvancedParameters.Uvalues.U_Roof_2
	}

	var rMeasure float64
	if c.Lvl0.AdvancedParameters.MeasureTypes.Code_MeasureType_Roof_2 == "ReplaceInsulation" {
		rMeasure = c.Lvl0.AdvancedParameters.Insulation.D_Insulation_Roof_2 / 0.04
	}

	return rAdd + uRoof - rMeasure
}

// calcRBeforeWall1 calculates R_Before_Wall_1
// Element type wall 1
// Unit: m²*K/W
// Formula: IFERROR(IF(ISNUMBER(R_Add_UnheatedSpace_Wall_1), R_Add_UnheatedSpace_Wall_1, 0) + IF(ISNUMBER(U_Wall_1), 1/U_Wall_1 - IF(AND(Code_MeasureType_Wall_1="ReplaceInsulation", NOT(ISERROR(R_Measure_Wall_1))), d_Insulation_Wall_1/0.04, 0), 0), 0)
func (c *CalcLevel3) calcRBeforeWall1() float64 {
	rAdd := c.Lvl0.AdvancedParameters.AdditionalResistances.R_Add_UnheatedSpace_Wall_1

	var uWall float64
	if c.Lvl0.AdvancedParameters.Uvalues.U_Wall_1 != 0 {
		uWall = 1 / c.Lvl0.AdvancedParameters.Uvalues.U_Wall_1
	}

	var rMeasure float64
	if c.Lvl0.AdvancedParameters.MeasureTypes.Code_MeasureType_Wall_1 == "ReplaceInsulation" {
		rMeasure = c.Lvl0.AdvancedParameters.Insulation.D_Insulation_Wall_1 / 0.04
	}

	return rAdd + uWall - rMeasure
}

// calcRBeforeWall2 calculates R_Before_Wall_2
// Element type wall 2
// Unit: m²*K/W
// Formula: IFERROR(IF(ISNUMBER(R_Add_UnheatedSpace_Wall_2), R_Add_UnheatedSpace_Wall_2, 0) + IF(ISNUMBER(U_Wall_2), 1/U_Wall_2 - IF(AND(Code_MeasureType_Wall_2="ReplaceInsulation", NOT(ISERROR(R_Measure_Wall_2))), d_Insulation_Wall_2/0.04, 0), 0), 0)
func (c *CalcLevel3) calcRBeforeWall2() float64 {
	rAdd := c.Lvl0.AdvancedParameters.AdditionalResistances.R_Add_UnheatedSpace_Wall_2

	var uWall float64
	if c.Lvl0.AdvancedParameters.Uvalues.U_Wall_2 != 0 {
		uWall = 1 / c.Lvl0.AdvancedParameters.Uvalues.U_Wall_2
	}

	var rMeasure float64
	if c.Lvl0.AdvancedParameters.MeasureTypes.Code_MeasureType_Wall_2 == "ReplaceInsulation" {
		rMeasure = c.Lvl0.AdvancedParameters.Insulation.D_Insulation_Wall_2 / 0.04
	}

	return rAdd + uWall - rMeasure
}

// calcRBeforeWall3 calculates R_Before_Wall_3
// Element type wall 3
// Unit: m²*K/W
// Formula: IFERROR(IF(ISNUMBER(R_Add_UnheatedSpace_Wall_3), R_Add_UnheatedSpace_Wall_3, 0) + IF(ISNUMBER(U_Wall_3), 1/U_Wall_3 - IF(AND(Code_MeasureType_Wall_3="ReplaceInsulation", NOT(ISERROR(R_Measure_Wall_3))), d_Insulation_Wall_3/0.04, 0), 0), 0)
func (c *CalcLevel3) calcRBeforeWall3() float64 {
	rAdd := c.Lvl0.AdvancedParameters.AdditionalResistances.R_Add_UnheatedSpace_Wall_3

	var uWall float64
	if c.Lvl0.AdvancedParameters.Uvalues.U_Wall_3 != 0 {
		uWall = 1 / c.Lvl0.AdvancedParameters.Uvalues.U_Wall_3
	}

	var rMeasure float64
	if c.Lvl0.AdvancedParameters.MeasureTypes.Code_MeasureType_Wall_3 == "ReplaceInsulation" {
		rMeasure = c.Lvl0.AdvancedParameters.Insulation.D_Insulation_Wall_3 / 0.04
	}

	return rAdd + uWall - rMeasure
}

// calcRBeforeFloor1 calculates R_Before_Floor_1
// Element type floor 1
// Unit: m²*K/W
// Formula: IFERROR(IF(ISNUMBER(R_Add_UnheatedSpace_Floor_1), R_Add_UnheatedSpace_Floor_1, 0) + IF(ISNUMBER(U_Floor_1), 1/U_Floor_1 - IF(AND(Code_MeasureType_Floor_1="ReplaceInsulation", NOT(ISERROR(R_Measure_Floor_1))), d_Insulation_Floor_1/0.04, 0), 0), 0)
func (c *CalcLevel3) calcRBeforeFloor1() float64 {
	rAdd := c.Lvl0.AdvancedParameters.AdditionalResistances.R_Add_UnheatedSpace_Floor_1

	var uFloor float64
	if c.Lvl0.AdvancedParameters.Uvalues.U_Floor_1 != 0 {
		uFloor = 1 / c.Lvl0.AdvancedParameters.Uvalues.U_Floor_1
	}

	var rMeasure float64
	if c.Lvl0.AdvancedParameters.MeasureTypes.Code_MeasureType_Floor_1 == "ReplaceInsulation" {
		rMeasure = c.Lvl0.AdvancedParameters.Insulation.D_Insulation_Floor_1 / 0.04
	}

	return rAdd + uFloor - rMeasure
}

// calcRBeforeFloor2 calculates R_Before_Floor_2
// Element type floor 2
// Unit: m²*K/W
// Formula: IFERROR(IF(ISNUMBER(R_Add_UnheatedSpace_Floor_2), R_Add_UnheatedSpace_Floor_2, 0) + IF(ISNUMBER(U_Floor_2), 1/U_Floor_2 - IF(AND(Code_MeasureType_Floor_2="ReplaceInsulation", NOT(ISERROR(R_Measure_Floor_2))), d_Insulation_Floor_2/0.04, 0), 0), 0)
func (c *CalcLevel3) calcRBeforeFloor2() float64 {
	rAdd := c.Lvl0.AdvancedParameters.AdditionalResistances.R_Add_UnheatedSpace_Floor_2

	var uFloor float64
	if c.Lvl0.AdvancedParameters.Uvalues.U_Floor_2 != 0 {
		uFloor = 1 / c.Lvl0.AdvancedParameters.Uvalues.U_Floor_2
	}

	var rMeasure float64
	if c.Lvl0.AdvancedParameters.MeasureTypes.Code_MeasureType_Floor_2 == "ReplaceInsulation" {
		rMeasure = c.Lvl0.AdvancedParameters.Insulation.D_Insulation_Floor_2 / 0.04
	}

	return rAdd + uFloor - rMeasure
}

// calcUActualWindow1 calculates U_Actual_Window_1
// Actual U-value for window 1
// Unit: W/(m²*K)
// Formula: (1-f_Measure_Window_1)*U_Window_1 + f_Measure_Window_1*U_Measure_Window_1
func (c *CalcLevel3) calcUActualWindow1() float64 {
	return (1-c.Lvl0.AdvancedParameters.MeasureFractions.F_Measure_Window_1)*c.Lvl0.AdvancedParameters.Uvalues.U_Window_1 +
		c.Lvl0.AdvancedParameters.MeasureFractions.F_Measure_Window_1*c.Lvl2.U_Measure_Window_1
}

// calcUActualWindow2 calculates U_Actual_Window_2
// Actual U-value for window 2
// Unit: W/(m²*K)
// Formula: (1-f_Measure_Window_2)*U_Window_2 + f_Measure_Window_2*U_Measure_Window_2
func (c *CalcLevel3) calcUActualWindow2() float64 {
	return (1-c.Lvl0.AdvancedParameters.MeasureFractions.F_Measure_Window_2)*c.Lvl0.AdvancedParameters.Uvalues.U_Window_2 +
		c.Lvl0.AdvancedParameters.MeasureFractions.F_Measure_Window_2*c.Lvl2.U_Measure_Window_2
}

// calcUActualDoor1 calculates U_Actual_Door_1
// Actual U-value for door 1
// Unit: W/(m²*K)
// Formula: (1-f_Measure_Door_1)*U_Door_1 + f_Measure_Door_1*U_Measure_Door_1
func (c *CalcLevel3) calcUActualDoor1() float64 {
	return (1-c.Lvl0.AdvancedParameters.MeasureFractions.F_Measure_Door_1)*c.Lvl0.AdvancedParameters.Uvalues.U_Door_1 +
		c.Lvl0.AdvancedParameters.MeasureFractions.F_Measure_Door_1*c.Lvl2.U_Measure_Door_1
}
