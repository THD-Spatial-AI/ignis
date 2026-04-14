package calc

// CalcLevel17 represents the seventeenth and final calculation level
type CalcLevel17 struct {
	Lvl1  *CalcLevel1
	Lvl8  *CalcLevel8
	Lvl14 *CalcLevel14
	Lvl16 *CalcLevel16

	// Calculated attributes - FINAL OUTPUT
	QHNd float64 `json:"q_h_nd"` // Final heating demand per unit area
}

// NewCalcLevel17 creates a new CalcLevel17 instance and runs calculations
func NewCalcLevel17(lvl1 *CalcLevel1, lvl8 *CalcLevel8, lvl14 *CalcLevel14, lvl16 *CalcLevel16) *CalcLevel17 {
	c := &CalcLevel17{
		Lvl1:  lvl1,
		Lvl8:  lvl8,
		Lvl14: lvl14,
		Lvl16: lvl16,
	}
	c.Run()
	return c
}

// Run executes all calculations for level 17 and returns the final output
func (c *CalcLevel17) Run() float64 {
	c.QHNd = c.calcQHNd()
	return c.QHNd
}

// calcQHNd calculates the final net heating demand per unit area
// This is the FINAL OUTPUT of the entire calculation pipeline
// Excel Formula: q_ht - eta_h_gn * (q_sol + q_int)
func (c *CalcLevel17) calcQHNd() float64 {
	return c.Lvl14.QHt - c.Lvl16.EtaHGn*(c.Lvl8.QSol+c.Lvl1.Q_int)
}
