package calc

// CalcLevel15 represents the fifteenth calculation level with all dependencies
type CalcLevel15 struct {
	Lvl1  *CalcLevel1
	Lvl8  *CalcLevel8
	Lvl14 *CalcLevel14

	// Calculated attributes
	GammaHGn float64 `json:"gamma_h_gn"`
}

// NewCalcLevel15 creates a new CalcLevel15 instance and runs calculations
func NewCalcLevel15(lvl1 *CalcLevel1, lvl8 *CalcLevel8, lvl14 *CalcLevel14) *CalcLevel15 {
	c := &CalcLevel15{
		Lvl1:  lvl1,
		Lvl8:  lvl8,
		Lvl14: lvl14,
	}
	c.Run()
	return c
}

// Run executes all calculations for level 15
func (c *CalcLevel15) Run() {
	c.GammaHGn = c.calcGammaHGn()
}

// calcGammaHGn calculates the heat gain utilization ratio
// Excel Formula: IFERROR((q_sol+q_int)/q_ht,0)
func (c *CalcLevel15) calcGammaHGn() float64 {
	if c.Lvl14.QHt == 0 {
		return 0
	}
	return (c.Lvl8.QSol + c.Lvl1.Q_int) / c.Lvl14.QHt
}
