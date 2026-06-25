package calc

import (
	"github.com/thd-spatial-ai/ignis/internal/models"
)

// CalcLevel8 represents the eighth calculation level with all dependencies
type CalcLevel8 struct {
	Lvl0 *models.TabulaBuildingParameters
	Lvl1 *CalcLevel1
	Lvl2 *CalcLevel2
	Lvl4 *CalcLevel4
	Lvl5 *CalcLevel5
	Lvl6 *CalcLevel6
	Lvl7 *CalcLevel7

	// Calculated attributes
	REnvTotalExactToEstim       float64 `json:"r_EnvTotal_ExactToEstim"`
	FractionEnvelopeRefurbished float64 `json:"Fraction_EnvelopeRefurbished"`
	HTransmissionWall1          float64 `json:"H_Transmission_Wall_1"`
	HTransmissionWall3          float64 `json:"H_Transmission_Wall_3"`
	QSol                        float64 `json:"q_sol"`
}

// NewCalcLevel8 creates a new CalcLevel8 instance and runs calculations
func NewCalcLevel8(lvl0 *models.TabulaBuildingParameters, lvl1 *CalcLevel1, lvl2 *CalcLevel2, lvl4 *CalcLevel4, lvl5 *CalcLevel5, lvl6 *CalcLevel6, lvl7 *CalcLevel7) *CalcLevel8 {
	c := &CalcLevel8{
		Lvl0: lvl0,
		Lvl1: lvl1,
		Lvl2: lvl2,
		Lvl4: lvl4,
		Lvl5: lvl5,
		Lvl6: lvl6,
		Lvl7: lvl7,
	}
	c.Run()
	return c
}

// Run executes all calculations for level 8
func (c *CalcLevel8) Run() {
	c.REnvTotalExactToEstim = c.calcREnvTotalExactToEstim()
	c.FractionEnvelopeRefurbished = c.calcFractionEnvelopeRefurbished()
	c.HTransmissionWall1 = c.calcHTransmissionWall1()
	c.HTransmissionWall3 = c.calcHTransmissionWall3()
	c.QSol = c.calcQSol()
}

// calcREnvTotalExactToEstim calculates ratio of exact to estimated envelope sum
// Excel Formula: IFERROR(A_Exact_Env_Sum/A_Estim_Env_Sum,0)
func (c *CalcLevel8) calcREnvTotalExactToEstim() float64 {
	if c.Lvl7.AEstimEnvSum == 0 {
		return 0
	}
	return c.Lvl1.A_Exact_Env_Sum / c.Lvl7.AEstimEnvSum
}

// calcFractionEnvelopeRefurbished calculates fraction of the envelope that has been refurbished
// Excel Formula: IFERROR((IF(R_Measure_Roof_1>0,f_Measure_Roof_1*A_Calc_Roof_1,0)+IF(R_Measure_Roof_2>0,f_Measure_Roof_2*A_Calc_Roof_2,0)+IF(R_Measure_Wall_1>0,f_Measure_Wall_1*A_Calc_Wall_1,0)+IF(R_Measure_Wall_2>0,f_Measure_Wall_2*A_Calc_Wall_2,0)+IF(R_Measure_Wall_3>0,f_Measure_Wall_3*A_Calc_Wall_3,0)+IF(R_Measure_Floor_1>0,f_Measure_Floor_1*A_Calc_Floor_1,0)+IF(R_Measure_Floor_2>0,f_Measure_Floor_2*A_Calc_Floor_2,0)+IF(R_Measure_Window_1>0,f_Measure_Window_1*A_Calc_Window_1,0)+IF(R_Measure_Window_2>0,f_Measure_Window_2*A_Calc_Window_2,0)+IF(R_Measure_Door_1>0,f_Measure_Door_1*A_Calc_Door_1,0))/SUM(A_Calc_Roof_1:A_Calc_Door_1),0)
func (c *CalcLevel8) calcFractionEnvelopeRefurbished() float64 {
	numerator := 0.0

	if c.Lvl2.R_Measure_Roof_1 > 0 {
		numerator += c.Lvl0.AdvancedParameters.MeasureFractions.F_Measure_Roof_1 * c.Lvl5.A_Calc_Roof_1
	}
	if c.Lvl2.R_Measure_Roof_2 > 0 {
		numerator += c.Lvl0.AdvancedParameters.MeasureFractions.F_Measure_Roof_2 * c.Lvl5.A_Calc_Roof_2
	}
	if c.Lvl2.R_Measure_Wall_1 > 0 {
		numerator += c.Lvl0.AdvancedParameters.MeasureFractions.F_Measure_Wall_1 * c.Lvl7.ACalcWall1
	}
	if c.Lvl2.R_Measure_Wall_2 > 0 {
		numerator += c.Lvl0.AdvancedParameters.MeasureFractions.F_Measure_Wall_2 * c.Lvl6.ACalcWall2
	}
	if c.Lvl2.R_Measure_Wall_3 > 0 {
		numerator += c.Lvl0.AdvancedParameters.MeasureFractions.F_Measure_Wall_3 * c.Lvl7.ACalcWall3
	}
	if c.Lvl2.R_Measure_Floor_1 > 0 {
		numerator += c.Lvl0.AdvancedParameters.MeasureFractions.F_Measure_Floor_1 * c.Lvl5.A_Calc_Floor_1
	}
	if c.Lvl2.R_Measure_Floor_2 > 0 {
		numerator += c.Lvl0.AdvancedParameters.MeasureFractions.F_Measure_Floor_2 * c.Lvl6.ACalcFloor2
	}
	if c.Lvl1.R_Measure_Window_1 > 0 {
		numerator += c.Lvl0.AdvancedParameters.MeasureFractions.F_Measure_Window_1 * c.Lvl5.A_Calc_Window_1
	}
	if c.Lvl1.R_Measure_Window_2 > 0 {
		numerator += c.Lvl0.AdvancedParameters.MeasureFractions.F_Measure_Window_2 * c.Lvl1.A_Calc_Window_2
	}
	if c.Lvl1.R_Measure_Door_1 > 0 {
		numerator += c.Lvl0.AdvancedParameters.MeasureFractions.F_Measure_Door_1 * c.Lvl4.A_Calc_Door_1
	}

	denominator := c.Lvl5.A_Calc_Roof_1 +
		c.Lvl5.A_Calc_Roof_2 +
		c.Lvl7.ACalcWall1 +
		c.Lvl6.ACalcWall2 +
		c.Lvl7.ACalcWall3 +
		c.Lvl5.A_Calc_Floor_1 +
		c.Lvl6.ACalcFloor2 +
		c.Lvl5.A_Calc_Window_1 +
		c.Lvl1.A_Calc_Window_2 +
		c.Lvl4.A_Calc_Door_1

	if denominator == 0 {
		return 0
	}
	return numerator / denominator
}

// calcHTransmissionWall1 calculates heat transmission through wall 1
// Excel Formula: IF(ISERROR(U_Actual_Wall_1*A_Calc_Wall_1*b_Transmission_Wall_1),0,U_Actual_Wall_1*A_Calc_Wall_1*b_Transmission_Wall_1)
func (c *CalcLevel8) calcHTransmissionWall1() float64 {
	result := c.Lvl5.U_Actual_Wall_1 * c.Lvl7.ACalcWall1 * c.Lvl0.AdvancedParameters.HeatLosses.B_Transmission_Wall_1
	return result
}

// calcHTransmissionWall3 calculates heat transmission through wall 3
// Excel Formula: IF(ISERROR(U_Actual_Wall_3*A_Calc_Wall_3*b_Transmission_Wall_3),0,U_Actual_Wall_3*A_Calc_Wall_3*b_Transmission_Wall_3)
func (c *CalcLevel8) calcHTransmissionWall3() float64 {
	result := c.Lvl5.U_Actual_Wall_3 * c.Lvl7.ACalcWall3 * c.Lvl0.AdvancedParameters.HeatLosses.B_Transmission_Wall_3
	return result
}

// calcQSol calculates total solar energy gain per unit area
// Excel Formula: IFERROR(SUM(Q_Sol_Hor:Q_Sol_North)/A_C_Ref,0)
func (c *CalcLevel8) calcQSol() float64 {
	totalSolarEnergy := c.Lvl7.QSolHor + c.Lvl7.QSolEast + c.Lvl7.QSolSouth + c.Lvl7.QSolWest + c.Lvl7.QSolNorth
	if c.Lvl2.A_C_Ref == 0 {
		return 0
	}
	return totalSolarEnergy / c.Lvl2.A_C_Ref
}
