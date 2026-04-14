package calc

import (
	"github.com/THD-Spatial-AI/hdcp-go/internal/models"
	"strings"
)

// CalcLevel5 holds calculated values for level 5
type CalcLevel5 struct {
	Lvl0 *models.TabulaBuildingParameters
	Lvl1 *CalcLevel1
	Lvl3 *CalcLevel3
	Lvl4 *CalcLevel4

	// Calculated attributes
	A_Estim_Wall_ToCellarOrSoil float64 `json:"A_Estim_Wall_ToCellarOrSoil"`
	R_EnvFloor_ExactToEstim     float64 `json:"r_EnvFloor_ExactToEstim"`
	R_EnvWindow_ExactToEstim    float64 `json:"r_EnvWindow_ExactToEstim"`
	A_Calc_Roof_1               float64 `json:"A_Calc_Roof_1"`
	A_Calc_Roof_2               float64 `json:"A_Calc_Roof_2"`
	A_Calc_Floor_1              float64 `json:"A_Calc_Floor_1"`
	A_Calc_Window_1             float64 `json:"A_Calc_Window_1"`
	U_Actual_Roof_1             float64 `json:"U_Actual_Roof_1"`
	U_Actual_Roof_2             float64 `json:"U_Actual_Roof_2"`
	U_Actual_Wall_1             float64 `json:"U_Actual_Wall_1"`
	U_Actual_Wall_2             float64 `json:"U_Actual_Wall_2"`
	U_Actual_Wall_3             float64 `json:"U_Actual_Wall_3"`
	U_Actual_Floor_1            float64 `json:"U_Actual_Floor_1"`
	U_Actual_Floor_2            float64 `json:"U_Actual_Floor_2"`
	H_Transmission_Door_1       float64 `json:"H_Transmission_Door_1"`
}

// NewCalcLevel5 creates a new CalcLevel5 instance and runs all calculations
func NewCalcLevel5(lvl0 *models.TabulaBuildingParameters, lvl1 *CalcLevel1, lvl3 *CalcLevel3, lvl4 *CalcLevel4) *CalcLevel5 {
	calc := &CalcLevel5{
		Lvl0: lvl0,
		Lvl1: lvl1,
		Lvl3: lvl3,
		Lvl4: lvl4,
	}
	calc.Run()
	return calc
}

// Run executes all calculation methods in CalcLevel5 and stores output in corresponding attributes
func (c *CalcLevel5) Run() {
	c.A_Estim_Wall_ToCellarOrSoil = c.calcAEstimWallToCellarOrSoil()
	c.R_EnvFloor_ExactToEstim = c.calcREnvFloorExactToEstim()
	c.R_EnvWindow_ExactToEstim = c.calcREnvWindowExactToEstim()
	c.A_Calc_Roof_1 = c.calcACalcRoof1()
	c.A_Calc_Roof_2 = c.calcACalcRoof2()
	c.A_Calc_Floor_1 = c.calcACalcFloor1()
	c.A_Calc_Window_1 = c.calcACalcWindow1()
	c.U_Actual_Roof_1 = c.calcUActualRoof1()
	c.U_Actual_Roof_2 = c.calcUActualRoof2()
	c.U_Actual_Wall_1 = c.calcUActualWall1()
	c.U_Actual_Wall_2 = c.calcUActualWall2()
	c.U_Actual_Wall_3 = c.calcUActualWall3()
	c.U_Actual_Floor_1 = c.calcUActualFloor1()
	c.U_Actual_Floor_2 = c.calcUActualFloor2()
	c.H_Transmission_Door_1 = c.calcHTransmissionDoor1()
}

// calcAEstimWallToCellarOrSoil calculates A_Estim_Wall_ToCellarOrSoil
// wall bordering at soil or unheated cellar
// Unit: m²
// Formula: 0.5*IF(RIGHT(Code_CellarCond,1)="I",1,f_CellarCond)*A_Estim_GrossWall_Storey
func (c *CalcLevel5) calcAEstimWallToCellarOrSoil() float64 {
	var cellarCond float64
	if strings.HasSuffix(c.Lvl0.BasicParameters.BuildingAppearance.Code_CellarCond, "I") {
		cellarCond = 1
	} else {
		cellarCond = c.Lvl1.F_CellarCond
	}
	return 0.5 * cellarCond * c.Lvl4.A_Estim_GrossWall_Storey
}

// calcREnvFloorExactToEstim calculates r_EnvFloor_ExactToEstim
// Formula: IFERROR((A_Floor_1+A_Floor_2)/(A_Estim_Floor),0)
func (c *CalcLevel5) calcREnvFloorExactToEstim() float64 {
	if c.Lvl4.A_Estim_Floor == 0 {
		return 0
	}
	return (c.Lvl0.BasicParameters.Envelope.A_Floor_1 + c.Lvl0.BasicParameters.Envelope.A_Floor_2) / c.Lvl4.A_Estim_Floor
}

// calcREnvWindowExactToEstim calculates r_EnvWindow_ExactToEstim
// Formula: IFERROR((A_Window_1+A_Window_2+A_Door_1)/(A_Estim_Window+A_Estim_Door),0)
func (c *CalcLevel5) calcREnvWindowExactToEstim() float64 {
	denominator := c.Lvl4.A_Estim_Window + c.Lvl3.A_Estim_Door
	if denominator == 0 {
		return 0
	}
	return (c.Lvl0.BasicParameters.Envelope.A_Window_1 + c.Lvl0.BasicParameters.Envelope.A_Window_2 + c.Lvl0.BasicParameters.Envelope.A_Door_1) / denominator
}

// calcACalcRoof1 calculates A_Calc_Roof_1
// element type roof 1
// Unit: m²
// Formula: IF($Code_TypeIntake_EnvelopeArea="Estimation",A_Estim_Roof,A_Roof_1)
func (c *CalcLevel5) calcACalcRoof1() float64 {
	if c.Lvl0.AdvancedParameters.PredefinedCodes.Code_TypeIntake_EnvelopeArea == "Estimation" {
		return c.Lvl4.A_Estim_Roof
	}
	return c.Lvl0.BasicParameters.Envelope.A_Roof_1
}

// calcACalcRoof2 calculates A_Calc_Roof_2
// element type roof 2
// Unit: m²
// Formula: IF($Code_TypeIntake_EnvelopeArea="Estimation",A_Estim_UpperCeiling,A_Roof_2)
func (c *CalcLevel5) calcACalcRoof2() float64 {
	if c.Lvl0.AdvancedParameters.PredefinedCodes.Code_TypeIntake_EnvelopeArea == "Estimation" {
		return c.Lvl4.A_Estim_UpperCeiling
	}
	return c.Lvl0.BasicParameters.Envelope.A_Roof_2
}

// calcACalcFloor1 calculates A_Calc_Floor_1
// element type floor 1
// Unit: m²
// Formula: IF($Code_TypeIntake_EnvelopeArea="Estimation",IF(Code_Estim_ConstructionBorder_Floor="Soil",0,A_Estim_Floor),A_Floor_1)
func (c *CalcLevel5) calcACalcFloor1() float64 {
	if c.Lvl0.AdvancedParameters.PredefinedCodes.Code_TypeIntake_EnvelopeArea == "Estimation" {
		if c.Lvl1.Code_Estim_ConstructionBorder_Floor == "Soil" {
			return 0
		}
		return c.Lvl4.A_Estim_Floor
	}
	return c.Lvl0.BasicParameters.Envelope.A_Floor_1
}

// calcACalcWindow1 calculates A_Calc_Window_1
// element type window 1
// Unit: m²
// Formula: IF($Code_TypeIntake_EnvelopeArea="Estimation",A_Estim_Window,A_Window_1)
func (c *CalcLevel5) calcACalcWindow1() float64 {
	if c.Lvl0.AdvancedParameters.PredefinedCodes.Code_TypeIntake_EnvelopeArea == "Estimation" {
		return c.Lvl4.A_Estim_Window
	}
	return c.Lvl0.BasicParameters.Envelope.A_Window_1
}

// calcUActualRoof1 calculates U_Actual_Roof_1
// element type roof 1
// Unit: W/(m²*K)
// Formula: IF(U_Roof_1>0,(1-f_Measure_Roof_1)*1/(1/U_Roof_1+R_Add_UnheatedSpace_Roof_1),0)+f_Measure_Roof_1*U_Measure_Roof_1
func (c *CalcLevel5) calcUActualRoof1() float64 {
	var result float64
	if c.Lvl0.AdvancedParameters.Uvalues.U_Roof_1 > 0 {
		result = (1 - c.Lvl0.AdvancedParameters.MeasureFractions.F_Measure_Roof_1) * 1 / (1/c.Lvl0.AdvancedParameters.Uvalues.U_Roof_1 + c.Lvl0.AdvancedParameters.AdditionalResistances.R_Add_UnheatedSpace_Roof_1)
	}
	return result + c.Lvl0.AdvancedParameters.MeasureFractions.F_Measure_Roof_1*c.Lvl4.U_Measure_Roof_1
}

// calcUActualRoof2 calculates U_Actual_Roof_2
// element type roof 2
// Unit: W/(m²*K)
// Formula: IF(U_Roof_2>0,(1-f_Measure_Roof_2)*1/(1/U_Roof_2+R_Add_UnheatedSpace_Roof_2),0)+f_Measure_Roof_2*U_Measure_Roof_2
func (c *CalcLevel5) calcUActualRoof2() float64 {
	var result float64
	if c.Lvl0.AdvancedParameters.Uvalues.U_Roof_2 > 0 {
		result = (1 - c.Lvl0.AdvancedParameters.MeasureFractions.F_Measure_Roof_2) * 1 / (1/c.Lvl0.AdvancedParameters.Uvalues.U_Roof_2 + c.Lvl0.AdvancedParameters.AdditionalResistances.R_Add_UnheatedSpace_Roof_2)
	}
	return result + c.Lvl0.AdvancedParameters.MeasureFractions.F_Measure_Roof_2*c.Lvl4.U_Measure_Roof_2
}

// calcUActualWall1 calculates U_Actual_Wall_1
// element type wall 1
// Unit: W/(m²*K)
// Formula: IF(U_Wall_1>0,(1-f_Measure_Wall_1)*1/(1/U_Wall_1+R_Add_UnheatedSpace_Wall_1),0)+f_Measure_Wall_1*U_Measure_Wall_1
func (c *CalcLevel5) calcUActualWall1() float64 {
	var result float64
	if c.Lvl0.AdvancedParameters.Uvalues.U_Wall_1 > 0 {
		result = (1 - c.Lvl0.AdvancedParameters.MeasureFractions.F_Measure_Wall_1) * 1 / (1/c.Lvl0.AdvancedParameters.Uvalues.U_Wall_1 + c.Lvl0.AdvancedParameters.AdditionalResistances.R_Add_UnheatedSpace_Wall_1)
	}
	return result + c.Lvl0.AdvancedParameters.MeasureFractions.F_Measure_Wall_1*c.Lvl4.U_Measure_Wall_1
}

// calcUActualWall2 calculates U_Actual_Wall_2
// element type wall 2
// Unit: W/(m²*K)
// Formula: IF(U_Wall_2>0,(1-f_Measure_Wall_2)*1/(1/U_Wall_2+R_Add_UnheatedSpace_Wall_2),0)+f_Measure_Wall_2*U_Measure_Wall_2
func (c *CalcLevel5) calcUActualWall2() float64 {
	var result float64
	if c.Lvl0.AdvancedParameters.Uvalues.U_Wall_2 > 0 {
		result = (1 - c.Lvl0.AdvancedParameters.MeasureFractions.F_Measure_Wall_2) * 1 / (1/c.Lvl0.AdvancedParameters.Uvalues.U_Wall_2 + c.Lvl0.AdvancedParameters.AdditionalResistances.R_Add_UnheatedSpace_Wall_2)
	}
	return result + c.Lvl0.AdvancedParameters.MeasureFractions.F_Measure_Wall_2*c.Lvl4.U_Measure_Wall_2
}

// calcUActualWall3 calculates U_Actual_Wall_3
// element type wall 3
// Unit: W/(m²*K)
// Formula: IF(U_Wall_3>0,(1-f_Measure_Wall_3)*1/(1/U_Wall_3+R_Add_UnheatedSpace_Wall_3),0)+f_Measure_Wall_3*U_Measure_Wall_3
func (c *CalcLevel5) calcUActualWall3() float64 {
	var result float64
	if c.Lvl0.AdvancedParameters.Uvalues.U_Wall_3 > 0 {
		result = (1 - c.Lvl0.AdvancedParameters.MeasureFractions.F_Measure_Wall_3) * 1 / (1/c.Lvl0.AdvancedParameters.Uvalues.U_Wall_3 + c.Lvl0.AdvancedParameters.AdditionalResistances.R_Add_UnheatedSpace_Wall_3)
	}
	return result + c.Lvl0.AdvancedParameters.MeasureFractions.F_Measure_Wall_3*c.Lvl4.U_Measure_Wall_3
}

// calcUActualFloor1 calculates U_Actual_Floor_1
// element type floor 1
// Unit: W/(m²*K)
// Formula: IF(U_Floor_1>0,(1-f_Measure_Floor_1)*1/(1/U_Floor_1+R_Add_UnheatedSpace_Floor_1),0)+f_Measure_Floor_1*U_Measure_Floor_1
func (c *CalcLevel5) calcUActualFloor1() float64 {
	var result float64
	if c.Lvl0.AdvancedParameters.Uvalues.U_Floor_1 > 0 {
		result = (1 - c.Lvl0.AdvancedParameters.MeasureFractions.F_Measure_Floor_1) * 1 / (1/c.Lvl0.AdvancedParameters.Uvalues.U_Floor_1 + c.Lvl0.AdvancedParameters.AdditionalResistances.R_Add_UnheatedSpace_Floor_1)
	}
	return result + c.Lvl0.AdvancedParameters.MeasureFractions.F_Measure_Floor_1*c.Lvl4.U_Measure_Floor_1
}

// calcUActualFloor2 calculates U_Actual_Floor_2
// element type floor 2
// Unit: W/(m²*K)
// Formula: IF(U_Floor_2>0,(1-f_Measure_Floor_2)*1/(1/U_Floor_2+R_Add_UnheatedSpace_Floor_2),0)+f_Measure_Floor_2*U_Measure_Floor_2
func (c *CalcLevel5) calcUActualFloor2() float64 {
	var result float64
	if c.Lvl0.AdvancedParameters.Uvalues.U_Floor_2 > 0 {
		result = (1 - c.Lvl0.AdvancedParameters.MeasureFractions.F_Measure_Floor_2) * 1 / (1/c.Lvl0.AdvancedParameters.Uvalues.U_Floor_2 + c.Lvl0.AdvancedParameters.AdditionalResistances.R_Add_UnheatedSpace_Floor_2)
	}
	return result + c.Lvl0.AdvancedParameters.MeasureFractions.F_Measure_Floor_2*c.Lvl4.U_Measure_Floor_2
}

// calcHTransmissionDoor1 calculates H_Transmission_Door_1
// element type door 1
// Unit: W/K
// Formula: IF(ISERROR(U_Actual_Door_1*A_Calc_Door_1*1),0,U_Actual_Door_1*A_Calc_Door_1*1)
func (c *CalcLevel5) calcHTransmissionDoor1() float64 {
	return c.Lvl3.U_Actual_Door_1 * c.Lvl4.A_Calc_Door_1 * 1
}
