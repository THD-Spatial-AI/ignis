package calc

import (
	"github.com/THD-Spatial-AI/hdcp-go/internal/models"
)

// CalcLevel6 represents the sixth calculation level with all dependencies
type CalcLevel6 struct {
	Lvl0 *models.TabulaBuildingParameters
	Lvl1 *CalcLevel1
	Lvl2 *CalcLevel2
	Lvl3 *CalcLevel3
	Lvl4 *CalcLevel4
	Lvl5 *CalcLevel5

	// Calculated attributes
	AEstimWallExtAir            float64 `json:"A_Estim_Wall_ExtAir"`
	CheckFloorAreaExactToEstim  int     `json:"Check_FloorArea_ExactToEstim"`
	CheckWindowAreaExactToEstim int     `json:"Check_WindowArea_ExactToEstim"`
	ACalcWall2                  float64 `json:"A_Calc_Wall_2"`
	ACalcFloor2                 float64 `json:"A_Calc_Floor_2"`
	ACalcWindowEast             float64 `json:"A_Calc_Window_East"`
	ACalcWindowWest             float64 `json:"A_Calc_Window_West"`
	HTransmissionRoof1          float64 `json:"H_Transmission_Roof_1"`
	HTransmissionRoof2          float64 `json:"H_Transmission_Roof_2"`
	HTransmissionFloor1         float64 `json:"H_Transmission_Floor_1"`
	HTransmissionWindow1        float64 `json:"H_Transmission_Window_1"`
	HTransmissionWindow2        float64 `json:"H_Transmission_Window_2"`
	GGlN                        float64 `json:"g_gl_n"`
}

// NewCalcLevel6 creates a new CalcLevel6 instance and runs calculations
func NewCalcLevel6(lvl0 *models.TabulaBuildingParameters, lvl1 *CalcLevel1, lvl2 *CalcLevel2, lvl3 *CalcLevel3, lvl4 *CalcLevel4, lvl5 *CalcLevel5) *CalcLevel6 {
	c := &CalcLevel6{
		Lvl0: lvl0,
		Lvl1: lvl1,
		Lvl2: lvl2,
		Lvl3: lvl3,
		Lvl4: lvl4,
		Lvl5: lvl5,
	}
	c.Run()
	return c
}

// Run executes all calculations for level 6
func (c *CalcLevel6) Run() {
	c.AEstimWallExtAir = c.calcAEstimWallExtAir()
	c.CheckFloorAreaExactToEstim = c.calcCheckFloorAreaExactToEstim()
	c.CheckWindowAreaExactToEstim = c.calcCheckWindowAreaExactToEstim()
	c.ACalcWall2 = c.calcACalcWall2()
	c.ACalcFloor2 = c.calcACalcFloor2()
	c.ACalcWindowEast = c.calcACalcWindowEast()
	c.ACalcWindowWest = c.calcACalcWindowWest()
	c.HTransmissionRoof1 = c.calcHTransmissionRoof1()
	c.HTransmissionRoof2 = c.calcHTransmissionRoof2()
	c.HTransmissionFloor1 = c.calcHTransmissionFloor1()
	c.HTransmissionWindow1 = c.calcHTransmissionWindow1()
	c.HTransmissionWindow2 = c.calcHTransmissionWindow2()
	c.GGlN = c.calcGGlN()
}

// calcAEstimWallExtAir calculates wall area (external air)
// Excel Formula: n_Storey_effective_envelope*A_Estim_GrossWall_Storey - A_Estim_Wall_ToCellarOrSoil - A_Estim_Window - A_Estim_Door
func (c *CalcLevel6) calcAEstimWallExtAir() float64 {
	return (c.Lvl2.N_Storey_effective_envelope * c.Lvl4.A_Estim_GrossWall_Storey) -
		c.Lvl5.A_Estim_Wall_ToCellarOrSoil - c.Lvl4.A_Estim_Window - c.Lvl3.A_Estim_Door
}

// calcCheckFloorAreaExactToEstim checks if the floor area estimation is within plausible limits
// Excel Formula: IF(AND(r_EnvFloor_ExactToEstim >= f_PlausiCrit_FloorArea_LowerLimit, r_EnvFloor_ExactToEstim <= f_PlausiCrit_FloorArea_UpperLimit), 1, 0)
func (c *CalcLevel6) calcCheckFloorAreaExactToEstim() int {
	if c.Lvl5.R_EnvFloor_ExactToEstim >= c.Lvl0.AdvancedParameters.PredefinedCodes.F_PlausiCrit_FloorArea_LowerLimit &&
		c.Lvl5.R_EnvFloor_ExactToEstim <= c.Lvl0.AdvancedParameters.PredefinedCodes.F_PlausiCrit_FloorArea_UpperLimit {
		return 1
	}
	return 0
}

// calcCheckWindowAreaExactToEstim checks if the window area estimation is within plausible limits
// Excel Formula: IF(AND(r_EnvWindow_ExactToEstim >= f_PlausiCrit_WindowArea_LowerLimit, r_EnvWindow_ExactToEstim <= f_PlausiCrit_WindowArea_UpperLimit), 1, 0)
func (c *CalcLevel6) calcCheckWindowAreaExactToEstim() int {
	if c.Lvl5.R_EnvWindow_ExactToEstim >= c.Lvl0.AdvancedParameters.PredefinedCodes.F_PlausiCrit_WindowArea_LowerLimit &&
		c.Lvl5.R_EnvWindow_ExactToEstim <= c.Lvl0.AdvancedParameters.PredefinedCodes.F_PlausiCrit_WindowArea_UpperLimit {
		return 1
	}
	return 0
}

// calcACalcWall2 calculates wall area based on envelope type and construction border
// Excel Formula: IF($Code_TypeIntake_EnvelopeArea="Estimation", IF(Code_Estim_ConstructionBorder_Wall_ToCellarOrSoil="Soil", 0, A_Estim_Wall_ToCellarOrSoil), A_Wall_2)
func (c *CalcLevel6) calcACalcWall2() float64 {
	if c.Lvl0.AdvancedParameters.PredefinedCodes.Code_TypeIntake_EnvelopeArea == "Estimation" {
		if c.Lvl1.Code_Estim_ConstructionBorder_Wall_ToCellarOrSoil == "Soil" {
			return 0
		}
		return c.Lvl5.A_Estim_Wall_ToCellarOrSoil
	}
	return c.Lvl0.BasicParameters.Envelope.A_Wall_2
}

// calcACalcFloor2 calculates floor area based on envelope type
// Excel Formula: IF($Code_TypeIntake_EnvelopeArea="Estimation", A_Estim_Floor - A_Calc_Floor_1, A_Floor_2)
func (c *CalcLevel6) calcACalcFloor2() float64 {
	if c.Lvl0.AdvancedParameters.PredefinedCodes.Code_TypeIntake_EnvelopeArea == "Estimation" {
		return c.Lvl4.A_Estim_Floor - c.Lvl5.A_Calc_Floor_1
	}
	return c.Lvl0.BasicParameters.Envelope.A_Floor_2
}

// calcACalcWindowEast calculates east window area based on envelope type
// Excel Formula: IF($Code_TypeIntake_EnvelopeArea="Estimation", 0.5 * SUM($A_Calc_Window_1:$A_Calc_Window_2), A_Window_East)
func (c *CalcLevel6) calcACalcWindowEast() float64 {
	if c.Lvl0.AdvancedParameters.PredefinedCodes.Code_TypeIntake_EnvelopeArea == "Estimation" {
		return 0.5 * (c.Lvl5.A_Calc_Window_1 + c.Lvl1.A_Calc_Window_2)
	}
	return c.Lvl0.BasicParameters.Envelope.A_Window_East
}

// calcACalcWindowWest calculates west window area based on envelope type
// Excel Formula: IF($Code_TypeIntake_EnvelopeArea="Estimation", 0.5 * SUM($A_Calc_Window_1:$A_Calc_Window_2), A_Window_West)
func (c *CalcLevel6) calcACalcWindowWest() float64 {
	if c.Lvl0.AdvancedParameters.PredefinedCodes.Code_TypeIntake_EnvelopeArea == "Estimation" {
		return 0.5 * (c.Lvl5.A_Calc_Window_1 + c.Lvl1.A_Calc_Window_2)
	}
	return c.Lvl0.BasicParameters.Envelope.A_Window_West
}

// calcHTransmissionRoof1 calculates heat transmission through roof 1
// Excel Formula: IF(ISERROR(U_Actual_Roof_1 * A_Calc_Roof_1 * b_Transmission_Roof_1), 0, U_Actual_Roof_1 * A_Calc_Roof_1 * b_Transmission_Roof_1)
func (c *CalcLevel6) calcHTransmissionRoof1() float64 {
	result := c.Lvl5.U_Actual_Roof_1 * c.Lvl5.A_Calc_Roof_1 * c.Lvl0.AdvancedParameters.HeatLosses.B_Transmission_Roof_1
	return result
}

// calcHTransmissionRoof2 calculates heat transmission through roof 2
// Excel Formula: IF(ISERROR(U_Actual_Roof_2 * A_Calc_Roof_2 * b_Transmission_Roof_2), 0, U_Actual_Roof_2 * A_Calc_Roof_2 * b_Transmission_Roof_2)
func (c *CalcLevel6) calcHTransmissionRoof2() float64 {
	result := c.Lvl5.U_Actual_Roof_2 * c.Lvl5.A_Calc_Roof_2 * c.Lvl0.AdvancedParameters.HeatLosses.B_Transmission_Roof_2
	return result
}

// calcHTransmissionFloor1 calculates heat transmission through floor 1
// Excel Formula: IF(ISERROR(U_Actual_Floor_1 * A_Calc_Floor_1 * b_Transmission_Floor_1), 0, U_Actual_Floor_1 * A_Calc_Floor_1 * b_Transmission_Floor_1)
func (c *CalcLevel6) calcHTransmissionFloor1() float64 {
	result := c.Lvl5.U_Actual_Floor_1 * c.Lvl5.A_Calc_Floor_1 * c.Lvl0.AdvancedParameters.HeatLosses.B_Transmission_Floor_1
	return result
}

// calcHTransmissionWindow1 calculates heat transmission through window 1
// Excel Formula: IF(ISERROR(U_Actual_Window_1 * A_Calc_Window_1 * 1), 0, U_Actual_Window_1 * A_Calc_Window_1 * 1)
func (c *CalcLevel6) calcHTransmissionWindow1() float64 {
	result := c.Lvl3.U_Actual_Window_1 * c.Lvl5.A_Calc_Window_1
	return result
}

// calcHTransmissionWindow2 calculates heat transmission through window 2
// Excel Formula: IF(ISERROR(U_Actual_Window_2 * A_Calc_Window_2 * 1), 0, U_Actual_Window_2 * A_Calc_Window_2 * 1)
func (c *CalcLevel6) calcHTransmissionWindow2() float64 {
	result := c.Lvl3.U_Actual_Window_2 * c.Lvl1.A_Calc_Window_2
	return result
}

// calcGGlN calculates average for both window types, considering refurbished state
// Excel Formula: IFERROR((IF(LEN(Code_Measure_Window_1)>1, IF(ISERROR(g_gl_n_Measure_Window_1), 0, g_gl_n_Measure_Window_1), IF(ISERROR(g_gl_n_Window_1), 0, g_gl_n_Window_1)) * A_Calc_Window_1 + IF(LEN(Code_Measure_Window_2)>1, IF(ISERROR(g_gl_n_Measure_Window_2), 0, g_gl_n_Measure_Window_2), IF(ISERROR(g_gl_n_Window_2), 0, g_gl_n_Window_2)) * A_Calc_Window_2) / (A_Calc_Window_1 + A_Calc_Window_2), 0)
func (c *CalcLevel6) calcGGlN() float64 {
	// Determine value for window 1
	// Note: Code_Measure can be empty string or "0" to indicate no measure
	var valueWindow1 float64
	code1 := c.Lvl0.AdvancedParameters.MeasureTypes.Code_Measure_Window_1
	if code1 != "" && code1 != "0" {
		valueWindow1 = c.Lvl1.G_gl_n_Measure_Window_1
	} else {
		valueWindow1 = c.Lvl0.AdvancedParameters.SolarTransmittance.G_gl_n_Window_1
	}
	weightedValueWindow1 := valueWindow1 * c.Lvl5.A_Calc_Window_1

	// Determine value for window 2
	var valueWindow2 float64
	code2 := c.Lvl0.AdvancedParameters.MeasureTypes.Code_Measure_Window_2
	if code2 != "" && code2 != "0" {
		valueWindow2 = c.Lvl1.G_gl_n_Measure_Window_2
	} else {
		valueWindow2 = c.Lvl0.AdvancedParameters.SolarTransmittance.G_gl_n_Window_2
	}
	weightedValueWindow2 := valueWindow2 * c.Lvl1.A_Calc_Window_2

	// Calculate final average
	denominator := c.Lvl5.A_Calc_Window_1 + c.Lvl1.A_Calc_Window_2
	if denominator == 0 {
		return 0
	}
	return (weightedValueWindow1 + weightedValueWindow2) / denominator
}
