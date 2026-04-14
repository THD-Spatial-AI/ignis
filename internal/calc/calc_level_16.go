package calc

import (
	"math"
)

// CalcLevel16 represents the sixteenth calculation level with all dependencies
type CalcLevel16 struct {
	Lvl13 *CalcLevel13
	Lvl15 *CalcLevel15

	// Calculated attributes
	EtaHGn float64 `json:"eta_h_gn"`
}

// NewCalcLevel16 creates a new CalcLevel16 instance and runs calculations
func NewCalcLevel16(lvl13 *CalcLevel13, lvl15 *CalcLevel15) *CalcLevel16 {
	c := &CalcLevel16{
		Lvl13: lvl13,
		Lvl15: lvl15,
	}
	c.Run()
	return c
}

// Run executes all calculations for level 16
func (c *CalcLevel16) Run() {
	c.EtaHGn = c.calcEtaHGn()
}

// calcEtaHGn calculates the gain utilization factor
// Excel Formula: (1-gamma_h_gn^a_H)/(1-gamma_h_gn^(a_H+1))
func (c *CalcLevel16) calcEtaHGn() float64 {
	numerator := 1 - math.Pow(c.Lvl15.GammaHGn, c.Lvl13.AH)
	denominator := 1 - math.Pow(c.Lvl15.GammaHGn, c.Lvl13.AH+1)

	if denominator == 0 {
		return 0
	}
	return numerator / denominator
}
