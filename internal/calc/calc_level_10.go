package calc

// CalcLevel10 represents the tenth calculation level with all dependencies
type CalcLevel10 struct {
	Lvl1 *CalcLevel1
	Lvl2 *CalcLevel2
	Lvl4 *CalcLevel4
	Lvl5 *CalcLevel5
	Lvl6 *CalcLevel6
	Lvl7 *CalcLevel7
	Lvl9 *CalcLevel9

	// Calculated attributes
	CheckEnvAreaExactToEstim     int     `json:"Check_EnvArea_ExactToEstim"`
	HTransmissionThermalBridging float64 `json:"H_Transmission_ThermalBridging"`
}

// NewCalcLevel10 creates a new CalcLevel10 instance and runs calculations
func NewCalcLevel10(lvl1 *CalcLevel1, lvl2 *CalcLevel2, lvl4 *CalcLevel4, lvl5 *CalcLevel5, lvl6 *CalcLevel6, lvl7 *CalcLevel7, lvl9 *CalcLevel9) *CalcLevel10 {
	c := &CalcLevel10{
		Lvl1: lvl1,
		Lvl2: lvl2,
		Lvl4: lvl4,
		Lvl5: lvl5,
		Lvl6: lvl6,
		Lvl7: lvl7,
		Lvl9: lvl9,
	}
	c.Run()
	return c
}

// Run executes all calculations for level 10
func (c *CalcLevel10) Run() {
	c.CheckEnvAreaExactToEstim = c.calcCheckEnvAreaExactToEstim()
	c.HTransmissionThermalBridging = c.calcHTransmissionThermalBridging()
}

// calcCheckEnvAreaExactToEstim checks if the environmental area is exact to estimated based on various conditions
// Excel Formula: Check_EnvSum_ExactToEstim * IF(Check_ToBeApplied_FloorArea_ExactToEstim=1, Check_FloorArea_ExactToEstim, 1) * Check_WindowArea_ExactToEstim
func (c *CalcLevel10) calcCheckEnvAreaExactToEstim() int {
	result := c.Lvl9.CheckEnvSumExactToEstim

	if c.Lvl2.Check_ToBeApplied_FloorArea_ExactToEstim == 1 {
		result *= c.Lvl6.CheckFloorAreaExactToEstim
	}

	result *= c.Lvl6.CheckWindowAreaExactToEstim

	if result > 1 {
		return 1
	} else if result < 0 {
		return -1
	}
	return result
}

// calcHTransmissionThermalBridging calculates the supplemental heat loss due to thermal bridging
// Excel Formula: SUM(A_Calc_Roof_1:A_Calc_Door_1) * delta_U_ThermalBridging
func (c *CalcLevel10) calcHTransmissionThermalBridging() float64 {
	totalArea := c.Lvl5.A_Calc_Roof_1 +
		c.Lvl5.A_Calc_Roof_2 +
		c.Lvl5.A_Calc_Floor_1 +
		c.Lvl6.ACalcFloor2 +
		c.Lvl7.ACalcWall1 +
		c.Lvl6.ACalcWall2 +
		c.Lvl7.ACalcWall3 +
		c.Lvl5.A_Calc_Window_1 +
		c.Lvl1.A_Calc_Window_2 +
		c.Lvl4.A_Calc_Door_1

	return totalArea * c.Lvl9.DeltaUThermalBridging
}
