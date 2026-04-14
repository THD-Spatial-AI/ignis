package calc

import (
	"github.com/THD-Spatial-AI/hdcp-go/internal/models"
)

// CalcLevel4 holds calculated values for level 4
type CalcLevel4 struct {
	Lvl0 *models.TabulaBuildingParameters
	Lvl1 *CalcLevel1
	Lvl2 *CalcLevel2
	Lvl3 *CalcLevel3

	// Calculated attributes
	A_Estim_GrossWall_Storey float64 `json:"A_Estim_GrossWall_Storey"`
	A_Estim_Roof             float64 `json:"A_Estim_Roof"`
	A_Estim_UpperCeiling     float64 `json:"A_Estim_UpperCeiling"`
	A_Estim_Floor            float64 `json:"A_Estim_Floor"`
	A_Estim_Window           float64 `json:"A_Estim_Window"`
	A_Calc_Door_1            float64 `json:"A_Calc_Door_1"`
	U_Measure_Roof_1         float64 `json:"U_Measure_Roof_1"`
	U_Measure_Roof_2         float64 `json:"U_Measure_Roof_2"`
	U_Measure_Wall_1         float64 `json:"U_Measure_Wall_1"`
	U_Measure_Wall_2         float64 `json:"U_Measure_Wall_2"`
	U_Measure_Wall_3         float64 `json:"U_Measure_Wall_3"`
	U_Measure_Floor_1        float64 `json:"U_Measure_Floor_1"`
	U_Measure_Floor_2        float64 `json:"U_Measure_Floor_2"`
}

// NewCalcLevel4 creates a new CalcLevel4 instance and runs all calculations
func NewCalcLevel4(lvl0 *models.TabulaBuildingParameters, lvl1 *CalcLevel1, lvl2 *CalcLevel2, lvl3 *CalcLevel3) *CalcLevel4 {
	calc := &CalcLevel4{
		Lvl0: lvl0,
		Lvl1: lvl1,
		Lvl2: lvl2,
		Lvl3: lvl3,
	}
	calc.Run()
	return calc
}

// Run executes all calculation methods in CalcLevel4 and stores output in corresponding attributes
func (c *CalcLevel4) Run() {
	c.A_Estim_GrossWall_Storey = c.calcAEstimGrossWallStorey()
	c.A_Estim_Roof = c.calcAEstimRoof()
	c.A_Estim_UpperCeiling = c.calcAEstimUpperCeiling()
	c.A_Estim_Floor = c.calcAEstimFloor()
	c.A_Estim_Window = c.calcAEstimWindow()
	c.A_Calc_Door_1 = c.calcACalcDoor1()
	c.U_Measure_Roof_1 = c.calcUMeasureRoof1()
	c.U_Measure_Roof_2 = c.calcUMeasureRoof2()
	c.U_Measure_Wall_1 = c.calcUMeasureWall1()
	c.U_Measure_Wall_2 = c.calcUMeasureWall2()
	c.U_Measure_Wall_3 = c.calcUMeasureWall3()
	c.U_Measure_Floor_1 = c.calcUMeasureFloor1()
	c.U_Measure_Floor_2 = c.calcUMeasureFloor2()
}

// calcAEstimGrossWallStorey calculates A_Estim_GrossWall_Storey
// interim quantity
// Unit: m²
// Formula: f_Corr_CeilingHeight * f_ComplexFootprint * (0.7 * A_C_Storey + IF(Code_AttachedNeighbours = "B_N2", 5, IF(Code_AttachedNeighbours = "B_N1", 25, 50)))
func (c *CalcLevel4) calcAEstimGrossWallStorey() float64 {
	var attachedNeigh float64
	switch c.Lvl0.BasicParameters.BuildingAppearance.Code_AttachedNeighbours {
	case "B_N2":
		attachedNeigh = 5
	case "B_N1":
		attachedNeigh = 25
	default:
		attachedNeigh = 50
	}
	return c.Lvl0.AdvancedParameters.PredefinedCodes.F_Corr_CeilingHeight * c.Lvl1.F_ComplexFootprint * (0.7*c.Lvl3.A_C_Storey + attachedNeigh)
}

// calcAEstimRoof calculates A_Estim_Roof
// roof
// Unit: m²
// Formula: f_ComplexRoof * (p_Roof * A_C_Storey + q_Roof)
func (c *CalcLevel4) calcAEstimRoof() float64 {
	return c.Lvl1.F_ComplexRoof * (c.Lvl1.P_Roof*c.Lvl3.A_C_Storey + c.Lvl1.Q_Roof)
}

// calcAEstimUpperCeiling calculates A_Estim_UpperCeiling
// upper ceiling
// Unit: m²
// Formula: p_Ceiling * A_C_Storey + q_Ceiling
func (c *CalcLevel4) calcAEstimUpperCeiling() float64 {
	return c.Lvl1.P_Ceiling*c.Lvl3.A_C_Storey + c.Lvl1.Q_Ceiling
}

// calcAEstimFloor calculates A_Estim_Floor
// floor above cellar or soil
// Unit: m²
// Formula: 1.2 * A_C_Storey + 5
func (c *CalcLevel4) calcAEstimFloor() float64 {
	return 1.2*c.Lvl3.A_C_Storey + 5
}

// calcAEstimWindow calculates A_Estim_Window
// window
// Unit: m²
// Formula: 0.18 * A_C_Ref - A_Estim_Door
func (c *CalcLevel4) calcAEstimWindow() float64 {
	return 0.18*c.Lvl2.A_C_Ref - c.Lvl3.A_Estim_Door
}

// calcACalcDoor1 calculates A_Calc_Door_1
// door calculation
// Unit: m²
// Formula: IF(Code_TypeIntake_EnvelopeArea = "Estimation", A_Estim_Door, A_Door_1)
func (c *CalcLevel4) calcACalcDoor1() float64 {
	if c.Lvl0.AdvancedParameters.PredefinedCodes.Code_TypeIntake_EnvelopeArea == "Estimation" {
		return c.Lvl3.A_Estim_Door
	}
	return c.Lvl0.BasicParameters.Envelope.A_Door_1
}

// calcUMeasureRoof1 calculates U_Measure_Roof_1
// element type roof 1
// Unit: W/(m²*K)
// Formula: IFERROR(1/(IF(Code_MeasureType_Roof_1 = "Replace", IF(ISNUMBER(R_Add_UnheatedSpace_Roof_1), R_Add_UnheatedSpace_Roof_1, 0), R_Before_Roof_1) + IF(ISNUMBER(R_Measure_Roof_1), R_Measure_Roof_1, 0)), 0)
func (c *CalcLevel4) calcUMeasureRoof1() float64 {
	var rBefore float64
	if c.Lvl0.AdvancedParameters.MeasureTypes.Code_MeasureType_Roof_1 == "Replace" {
		rBefore = c.Lvl0.AdvancedParameters.AdditionalResistances.R_Add_UnheatedSpace_Roof_1
	} else {
		rBefore = c.Lvl3.R_Before_Roof_1
	}

	denominator := rBefore + c.Lvl2.R_Measure_Roof_1
	if denominator != 0 {
		return 1 / denominator
	}
	return 0
}

// calcUMeasureRoof2 calculates U_Measure_Roof_2
// element type roof 2
// Unit: W/(m²*K)
// Formula: IFERROR(1/(IF(Code_MeasureType_Roof_2 = "Replace", IF(ISNUMBER(R_Add_UnheatedSpace_Roof_2), R_Add_UnheatedSpace_Roof_2, 0), R_Before_Roof_2) + IF(ISNUMBER(R_Measure_Roof_2), R_Measure_Roof_2, 0)), 0)
func (c *CalcLevel4) calcUMeasureRoof2() float64 {
	var rBefore float64
	if c.Lvl0.AdvancedParameters.MeasureTypes.Code_MeasureType_Roof_2 == "Replace" {
		rBefore = c.Lvl0.AdvancedParameters.AdditionalResistances.R_Add_UnheatedSpace_Roof_2
	} else {
		rBefore = c.Lvl3.R_Before_Roof_2
	}

	denominator := rBefore + c.Lvl2.R_Measure_Roof_2
	if denominator != 0 {
		return 1 / denominator
	}
	return 0
}

// calcUMeasureWall1 calculates U_Measure_Wall_1
// element type wall 1
// Unit: W/(m²*K)
// Formula: IFERROR(1/(IF(Code_MeasureType_Wall_1 = "Replace", IF(ISNUMBER(R_Add_UnheatedSpace_Wall_1), R_Add_UnheatedSpace_Wall_1, 0), R_Before_Wall_1) + IF(ISNUMBER(R_Measure_Wall_1), R_Measure_Wall_1, 0)), 0)
func (c *CalcLevel4) calcUMeasureWall1() float64 {
	var rBefore float64
	if c.Lvl0.AdvancedParameters.MeasureTypes.Code_MeasureType_Wall_1 == "Replace" {
		rBefore = c.Lvl0.AdvancedParameters.AdditionalResistances.R_Add_UnheatedSpace_Wall_1
	} else {
		rBefore = c.Lvl3.R_Before_Wall_1
	}

	denominator := rBefore + c.Lvl2.R_Measure_Wall_1
	if denominator != 0 {
		return 1 / denominator
	}
	return 0
}

// calcUMeasureWall2 calculates U_Measure_Wall_2
// element type wall 2
// Unit: W/(m²*K)
// Formula: IFERROR(1/(IF(Code_MeasureType_Wall_2 = "Replace", IF(ISNUMBER(R_Add_UnheatedSpace_Wall_2), R_Add_UnheatedSpace_Wall_2, 0), R_Before_Wall_2) + IF(ISNUMBER(R_Measure_Wall_2), R_Measure_Wall_2, 0)), 0)
func (c *CalcLevel4) calcUMeasureWall2() float64 {
	var rBefore float64
	if c.Lvl0.AdvancedParameters.MeasureTypes.Code_MeasureType_Wall_2 == "Replace" {
		rBefore = c.Lvl0.AdvancedParameters.AdditionalResistances.R_Add_UnheatedSpace_Wall_2
	} else {
		rBefore = c.Lvl3.R_Before_Wall_2
	}

	denominator := rBefore + c.Lvl2.R_Measure_Wall_2
	if denominator != 0 {
		return 1 / denominator
	}
	return 0
}

// calcUMeasureWall3 calculates U_Measure_Wall_3
// element type wall 3
// Unit: W/(m²*K)
// Formula: IFERROR(1/(IF(Code_MeasureType_Wall_3 = "Replace", IF(ISNUMBER(R_Add_UnheatedSpace_Wall_3), R_Add_UnheatedSpace_Wall_3, 0), R_Before_Wall_3) + IF(ISNUMBER(R_Measure_Wall_3), R_Measure_Wall_3, 0)), 0)
func (c *CalcLevel4) calcUMeasureWall3() float64 {
	var rBefore float64
	if c.Lvl0.AdvancedParameters.MeasureTypes.Code_MeasureType_Wall_3 == "Replace" {
		rBefore = c.Lvl0.AdvancedParameters.AdditionalResistances.R_Add_UnheatedSpace_Wall_3
	} else {
		rBefore = c.Lvl3.R_Before_Wall_3
	}

	denominator := rBefore + c.Lvl2.R_Measure_Wall_3
	if denominator != 0 {
		return 1 / denominator
	}
	return 0
}

// calcUMeasureFloor1 calculates U_Measure_Floor_1
// element type floor 1
// Unit: W/(m²K)
// Formula: IFERROR(1/(IF(Code_MeasureType_Floor_1="Replace",IF(ISNUMBER(R_Add_UnheatedSpace_Floor_1),R_Add_UnheatedSpace_Floor_1,0),R_Before_Floor_1)+IF(ISNUMBER(R_Measure_Floor_1),R_Measure_Floor_1,0)),0)
func (c *CalcLevel4) calcUMeasureFloor1() float64 {
	var rBefore float64
	if c.Lvl0.AdvancedParameters.MeasureTypes.Code_MeasureType_Floor_1 == "Replace" {
		rBefore = c.Lvl0.AdvancedParameters.AdditionalResistances.R_Add_UnheatedSpace_Floor_1
	} else {
		rBefore = c.Lvl3.R_Before_Floor_1
	}

	denominator := rBefore + c.Lvl2.R_Measure_Floor_1
	if denominator != 0 {
		return 1 / denominator
	}
	return 0
}

// calcUMeasureFloor2 calculates U_Measure_Floor_2
// element type floor 2
// Unit: W/(m²K)
// Formula: IFERROR(1/(IF(Code_MeasureType_Floor_2="Replace",IF(ISNUMBER(R_Add_UnheatedSpace_Floor_2),R_Add_UnheatedSpace_Floor_2,0),R_Before_Floor_2)+IF(ISNUMBER(R_Measure_Floor_2),R_Measure_Floor_2,0)),0)
func (c *CalcLevel4) calcUMeasureFloor2() float64 {
	var rBefore float64
	if c.Lvl0.AdvancedParameters.MeasureTypes.Code_MeasureType_Floor_2 == "Replace" {
		rBefore = c.Lvl0.AdvancedParameters.AdditionalResistances.R_Add_UnheatedSpace_Floor_2
	} else {
		rBefore = c.Lvl3.R_Before_Floor_2
	}

	denominator := rBefore + c.Lvl2.R_Measure_Floor_2
	if denominator != 0 {
		return 1 / denominator
	}
	return 0
}
