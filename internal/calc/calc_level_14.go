package calc

// CalcLevel14 represents the fourteenth calculation level with all dependencies
type CalcLevel14 struct {
	Lvl13 *CalcLevel13

	// Calculated attributes
	QHt float64 `json:"q_ht"`
}

// NewCalcLevel14 creates a new CalcLevel14 instance and runs calculations
func NewCalcLevel14(lvl13 *CalcLevel13) *CalcLevel14 {
	c := &CalcLevel14{
		Lvl13: lvl13,
	}
	c.Run()
	return c
}

// Run executes all calculations for level 14
func (c *CalcLevel14) Run() {
	c.QHt = c.calcQHt()
}

// calcQHt calculates total heat transfer
// Excel Formula: q_ht_tr+q_ht_ve
func (c *CalcLevel14) calcQHt() float64 {
	return c.Lvl13.QHtTr + c.Lvl13.QHtVe
}
