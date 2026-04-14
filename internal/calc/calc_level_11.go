package calc

// CalcLevel11 represents the eleventh calculation level with all dependencies
type CalcLevel11 struct {
	Lvl2  *CalcLevel2
	Lvl5  *CalcLevel5
	Lvl6  *CalcLevel6
	Lvl7  *CalcLevel7
	Lvl8  *CalcLevel8
	Lvl10 *CalcLevel10

	// Calculated attributes
	HTransmission float64 `json:"h_Transmission"`
}

// NewCalcLevel11 creates a new CalcLevel11 instance and runs calculations
func NewCalcLevel11(lvl2 *CalcLevel2, lvl5 *CalcLevel5, lvl6 *CalcLevel6, lvl7 *CalcLevel7, lvl8 *CalcLevel8, lvl10 *CalcLevel10) *CalcLevel11 {
	c := &CalcLevel11{
		Lvl2:  lvl2,
		Lvl5:  lvl5,
		Lvl6:  lvl6,
		Lvl7:  lvl7,
		Lvl8:  lvl8,
		Lvl10: lvl10,
	}
	c.Run()
	return c
}

// Run executes all calculations for level 11
func (c *CalcLevel11) Run() {
	c.HTransmission = c.calcHTransmission()
}

// calcHTransmission calculates indicator for energy quality of building envelope (compactness + insulation)
// Excel Formula: IFERROR(SUM(H_Transmission_Roof_1:H_Transmission_ThermalBridging)/A_C_Ref,0)
func (c *CalcLevel11) calcHTransmission() float64 {
	totalTransmission := c.Lvl6.HTransmissionFloor1 +
		c.Lvl7.HTransmissionFloor2 +
		c.Lvl6.HTransmissionRoof1 +
		c.Lvl6.HTransmissionRoof2 +
		c.Lvl6.HTransmissionWindow1 +
		c.Lvl6.HTransmissionWindow2 +
		c.Lvl8.HTransmissionWall1 +
		c.Lvl7.HTransmissionWall2 +
		c.Lvl8.HTransmissionWall3 +
		c.Lvl5.H_Transmission_Door_1 +
		c.Lvl10.HTransmissionThermalBridging

	if c.Lvl2.A_C_Ref == 0 {
		return 0
	}
	return totalTransmission / c.Lvl2.A_C_Ref
}
