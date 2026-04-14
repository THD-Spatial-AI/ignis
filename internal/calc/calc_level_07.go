package calc

import (
	"github.com/THD-Spatial-AI/hdcp-go/internal/models"
)

// CalcLevel7 represents the seventh calculation level with all dependencies
type CalcLevel7 struct {
	Lvl0 *models.TabulaBuildingParameters
	Lvl1 *CalcLevel1
	Lvl3 *CalcLevel3
	Lvl4 *CalcLevel4
	Lvl5 *CalcLevel5
	Lvl6 *CalcLevel6

	// Calculated attributes
	AEstimEnvSum        float64 `json:"A_Estim_Env_Sum"`
	ACalcWall1          float64 `json:"A_Calc_Wall_1"`
	ACalcWall3          float64 `json:"A_Calc_Wall_3"`
	HTransmissionWall2  float64 `json:"H_Transmission_Wall_2"`
	HTransmissionFloor2 float64 `json:"H_Transmission_Floor_2"`
	QSolHor             float64 `json:"Q_Sol_Hor"`
	QSolEast            float64 `json:"Q_Sol_East"`
	QSolSouth           float64 `json:"Q_Sol_South"`
	QSolWest            float64 `json:"Q_Sol_West"`
	QSolNorth           float64 `json:"Q_Sol_North"`
}

// NewCalcLevel7 creates a new CalcLevel7 instance and runs calculations
func NewCalcLevel7(lvl0 *models.TabulaBuildingParameters, lvl1 *CalcLevel1, lvl3 *CalcLevel3, lvl4 *CalcLevel4, lvl5 *CalcLevel5, lvl6 *CalcLevel6) *CalcLevel7 {
	c := &CalcLevel7{
		Lvl0: lvl0,
		Lvl1: lvl1,
		Lvl3: lvl3,
		Lvl4: lvl4,
		Lvl5: lvl5,
		Lvl6: lvl6,
	}
	c.Run()
	return c
}

// Run executes all calculations for level 7
func (c *CalcLevel7) Run() {
	c.AEstimEnvSum = c.calcAEstimEnvSum()
	c.ACalcWall1 = c.calcACalcWall1()
	c.ACalcWall3 = c.calcACalcWall3()
	c.HTransmissionWall2 = c.calcHTransmissionWall2()
	c.HTransmissionFloor2 = c.calcHTransmissionFloor2()
	c.QSolHor = c.calcQSolHor()
	c.QSolEast = c.calcQSolEast()
	c.QSolSouth = c.calcQSolSouth()
	c.QSolWest = c.calcQSolWest()
	c.QSolNorth = c.calcQSolNorth()
}

// calcAEstimEnvSum calculates sum of estimated envelope areas
// Excel Formula: SUM(A_Estim_Roof:A_Estim_Door)
func (c *CalcLevel7) calcAEstimEnvSum() float64 {
	return c.Lvl4.A_Estim_Roof +
		c.Lvl4.A_Estim_UpperCeiling +
		c.Lvl6.AEstimWallExtAir +
		c.Lvl5.A_Estim_Wall_ToCellarOrSoil +
		c.Lvl4.A_Estim_Floor +
		c.Lvl4.A_Estim_Window +
		c.Lvl3.A_Estim_Door
}

// calcACalcWall1 calculates wall area based on envelope type
// Excel Formula: IF($Code_TypeIntake_EnvelopeArea="Estimation", A_Estim_Wall_ExtAir, A_Wall_1)
func (c *CalcLevel7) calcACalcWall1() float64 {
	if c.Lvl0.AdvancedParameters.PredefinedCodes.Code_TypeIntake_EnvelopeArea == "Estimation" {
		return c.Lvl6.AEstimWallExtAir
	}
	return c.Lvl0.BasicParameters.Envelope.A_Wall_1
}

// calcACalcWall3 calculates wall area 3 based on envelope type
// Excel Formula: IF($Code_TypeIntake_EnvelopeArea="Estimation", A_Estim_Wall_ToCellarOrSoil - A_Calc_Wall_2, A_Wall_3)
func (c *CalcLevel7) calcACalcWall3() float64 {
	if c.Lvl0.AdvancedParameters.PredefinedCodes.Code_TypeIntake_EnvelopeArea == "Estimation" {
		return c.Lvl5.A_Estim_Wall_ToCellarOrSoil - c.Lvl6.ACalcWall2
	}
	return c.Lvl0.BasicParameters.Envelope.A_Wall_3
}

// calcHTransmissionWall2 calculates heat transmission through wall 2
// Excel Formula: IF(ISERROR(U_Actual_Wall_2 * A_Calc_Wall_2 * b_Transmission_Wall_2), 0, U_Actual_Wall_2 * A_Calc_Wall_2 * b_Transmission_Wall_2)
func (c *CalcLevel7) calcHTransmissionWall2() float64 {
	result := c.Lvl5.U_Actual_Wall_2 * c.Lvl6.ACalcWall2 * c.Lvl0.AdvancedParameters.HeatLosses.B_Transmission_Wall_2
	return result
}

// calcHTransmissionFloor2 calculates heat transmission through floor 2
// Excel Formula: IF(ISERROR(U_Actual_Floor_2 * A_Calc_Floor_2 * b_Transmission_Floor_2), 0, U_Actual_Floor_2 * A_Calc_Floor_2 * b_Transmission_Floor_2)
func (c *CalcLevel7) calcHTransmissionFloor2() float64 {
	result := c.Lvl5.U_Actual_Floor_2 * c.Lvl6.ACalcFloor2 * c.Lvl0.AdvancedParameters.HeatLosses.B_Transmission_Floor_2
	return result
}

// calcQSolHor calculates solar energy gain through horizontal windows
// Excel Formula: A_Calc_Window_Horizontal * I_Sol_Hor * F_sh_hor * (1 - F_f) * F_w * g_gl_n
func (c *CalcLevel7) calcQSolHor() float64 {
	return c.Lvl1.A_Calc_Window_Horizontal *
		c.Lvl0.AdvancedParameters.SolarGains.I_Sol_Horizontal *
		c.Lvl0.AdvancedParameters.HeatTransfer.F_sh_hor *
		(1 - c.Lvl0.AdvancedParameters.HeatTransfer.F_f) *
		c.Lvl0.AdvancedParameters.HeatTransfer.F_w *
		c.Lvl6.GGlN
}

// calcQSolEast calculates solar energy gain through east windows
// Excel Formula: A_Calc_Window_East * I_Sol_East * F_sh_vert * (1 - F_f) * F_w * g_gl_n
func (c *CalcLevel7) calcQSolEast() float64 {
	return c.Lvl6.ACalcWindowEast *
		c.Lvl0.AdvancedParameters.SolarGains.I_Sol_East *
		c.Lvl0.AdvancedParameters.HeatTransfer.F_sh_vert *
		(1 - c.Lvl0.AdvancedParameters.HeatTransfer.F_f) *
		c.Lvl0.AdvancedParameters.HeatTransfer.F_w *
		c.Lvl6.GGlN
}

// calcQSolSouth calculates solar energy gain through south windows
// Excel Formula: A_Calc_Window_South * I_Sol_South * F_sh_vert * (1 - F_f) * F_w * g_gl_n
func (c *CalcLevel7) calcQSolSouth() float64 {
	return c.Lvl1.A_Calc_Window_South *
		c.Lvl0.AdvancedParameters.SolarGains.I_Sol_South *
		c.Lvl0.AdvancedParameters.HeatTransfer.F_sh_vert *
		(1 - c.Lvl0.AdvancedParameters.HeatTransfer.F_f) *
		c.Lvl0.AdvancedParameters.HeatTransfer.F_w *
		c.Lvl6.GGlN
}

// calcQSolWest calculates solar energy gain through west windows
// Excel Formula: A_Calc_Window_West * I_Sol_West * F_sh_vert * (1 - F_f) * F_w * g_gl_n
func (c *CalcLevel7) calcQSolWest() float64 {
	return c.Lvl6.ACalcWindowWest *
		c.Lvl0.AdvancedParameters.SolarGains.I_Sol_West *
		c.Lvl0.AdvancedParameters.HeatTransfer.F_sh_vert *
		(1 - c.Lvl0.AdvancedParameters.HeatTransfer.F_f) *
		c.Lvl0.AdvancedParameters.HeatTransfer.F_w *
		c.Lvl6.GGlN
}

// calcQSolNorth calculates solar energy gain through north windows
// Excel Formula: A_Calc_Window_North * I_Sol_North * F_sh_vert * (1 - F_f) * F_w * g_gl_n
func (c *CalcLevel7) calcQSolNorth() float64 {
	return c.Lvl1.A_Calc_Window_North *
		c.Lvl0.AdvancedParameters.SolarGains.I_Sol_North *
		c.Lvl0.AdvancedParameters.HeatTransfer.F_sh_vert *
		(1 - c.Lvl0.AdvancedParameters.HeatTransfer.F_f) *
		c.Lvl0.AdvancedParameters.HeatTransfer.F_w *
		c.Lvl6.GGlN
}
