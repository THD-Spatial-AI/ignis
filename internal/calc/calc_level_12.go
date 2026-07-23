package calc

import (
	"github.com/thd-spatial-ai/ignis/internal/models"
)

// CalcLevel12 represents the twelfth calculation level with all dependencies
type CalcLevel12 struct {
	Lvl0  *models.TabulaBuildingParameters
	Lvl1  *CalcLevel1
	Lvl11 *CalcLevel11

	// Calculated attributes
	FRedTemp float64 `json:"F_red_temp"`
	Tau      float64 `json:"tau"`
}

// NewCalcLevel12 creates a new CalcLevel12 instance and runs calculations
func NewCalcLevel12(lvl0 *models.TabulaBuildingParameters, lvl1 *CalcLevel1, lvl11 *CalcLevel11) *CalcLevel12 {
	c := &CalcLevel12{
		Lvl0:  lvl0,
		Lvl1:  lvl1,
		Lvl11: lvl11,
	}
	c.Run()
	return c
}

// Run executes all calculations for level 12
func (c *CalcLevel12) Run() {
	c.FRedTemp = c.calcFRedTemp()
	c.Tau = c.calcTau()
}

// determines the reduced temperature factor based on h_Transmission
// Excel Formula: IF(h_Transmission<=1,F_red_htr1+(1-h_Transmission)/0.5*(1-F_red_htr1),IF(h_Transmission>=4,F_red_htr4,F_red_htr1+(h_Transmission-1)*(F_red_htr4-F_red_htr1)/(4-1)))
func (c *CalcLevel12) calcFRedTemp() float64 {
	if c.Lvl11.HTransmission <= 1 {
		return c.Lvl0.AdvancedParameters.HeatTransfer.F_red_htr1 +
			(1-c.Lvl11.HTransmission)/0.5*(1-c.Lvl0.AdvancedParameters.HeatTransfer.F_red_htr1)
	} else if c.Lvl11.HTransmission >= 4 {
		return c.Lvl0.AdvancedParameters.HeatTransfer.F_red_htr4
	} else {
		return c.Lvl0.AdvancedParameters.HeatTransfer.F_red_htr1 +
			(c.Lvl11.HTransmission-1)*(c.Lvl0.AdvancedParameters.HeatTransfer.F_red_htr4-c.Lvl0.AdvancedParameters.HeatTransfer.F_red_htr1)/(4-1)
	}
}

// time constant relevant for seasonal method
// Excel Formula: c_m/(h_Transmission+h_Ventilation)
func (c *CalcLevel12) calcTau() float64 {
	return c.Lvl0.AdvancedParameters.HeatTransfer.C_m / (c.Lvl11.HTransmission + c.Lvl1.H_Ventilation)
}
