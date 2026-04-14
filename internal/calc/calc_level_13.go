package calc

// CalcLevel13 represents the thirteenth calculation level with all dependencies
type CalcLevel13 struct {
	Lvl1  *CalcLevel1
	Lvl11 *CalcLevel11
	Lvl12 *CalcLevel12

	// Calculated attributes
	QHtTr float64 `json:"q_ht_tr"`
	QHtVe float64 `json:"q_ht_ve"`
	AH    float64 `json:"a_H"`
}

// NewCalcLevel13 creates a new CalcLevel13 instance and runs calculations
func NewCalcLevel13(lvl1 *CalcLevel1, lvl11 *CalcLevel11, lvl12 *CalcLevel12) *CalcLevel13 {
	c := &CalcLevel13{
		Lvl1:  lvl1,
		Lvl11: lvl11,
		Lvl12: lvl12,
	}
	c.Run()
	return c
}

// Run executes all calculations for level 13
func (c *CalcLevel13) Run() {
	c.QHtTr = c.calcQHtTr()
	c.QHtVe = c.calcQHtVe()
	c.AH = c.calcAH()
}

// calcQHtTr calculates the heat transfer due to thermal transmission
// Excel Formula: h_Transmission * 0.024 * Sum_DeltaT_for_HeatingDays * F_red_temp
func (c *CalcLevel13) calcQHtTr() float64 {
	return c.Lvl11.HTransmission * 0.024 * c.Lvl1.Sum_DeltaT_for_HeatingDays * c.Lvl12.FRedTemp
}

// calcQHtVe calculates the heat transfer due to ventilation
// Excel Formula: h_Ventilation * 0.024 * Sum_DeltaT_for_HeatingDays * F_red_temp
func (c *CalcLevel13) calcQHtVe() float64 {
	return c.Lvl1.H_Ventilation * 0.024 * c.Lvl1.Sum_DeltaT_for_HeatingDays * c.Lvl12.FRedTemp
}

// calcAH calculates the heat adaptation factor
// Excel Formula: 0.8 + tau / 30
func (c *CalcLevel13) calcAH() float64 {
	return 0.8 + c.Lvl12.Tau/30
}
