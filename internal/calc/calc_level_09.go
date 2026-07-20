package calc

import (
	"fmt"
	"github.com/thd-spatial-ai/ignis/internal/models"
)

// CalcLevel9 represents the ninth calculation level with all dependencies
type CalcLevel9 struct {
	Lvl0 *models.TabulaBuildingParameters
	Lvl8 *CalcLevel8

	// Calculated attributes
	CheckEnvSumExactToEstim   int     `json:"Check_EnvSum_ExactToEstim"`
	TypeThermalBridgingActual string  `json:"Type_ThermalBridging_Actual"`
	DeltaUThermalBridging     float64 `json:"delta_U_ThermalBridging"`
}

// NewCalcLevel9 creates a new CalcLevel9 instance and runs calculations
func NewCalcLevel9(lvl0 *models.TabulaBuildingParameters, lvl8 *CalcLevel8) *CalcLevel9 {
	c := &CalcLevel9{
		Lvl0: lvl0,
		Lvl8: lvl8,
	}
	c.Run()
	return c
}

// Run executes all calculations for level 9
func (c *CalcLevel9) Run() {
	c.CheckEnvSumExactToEstim = c.calcCheckEnvSumExactToEstim()
	c.TypeThermalBridgingActual = c.calcTypeThermalBridgingActual()
	c.DeltaUThermalBridging = c.calcDeltaUThermalBridging()
}

// checks if the ratio of exact to estimated envelope sum falls within plausible limits
// Excel Formula: IF(AND(r_EnvTotal_ExactToEstim>=f_PlausiCrit_EnvSum_LowerLimit,r_EnvTotal_ExactToEstim<=f_PlausiCrit_EnvSum_UpperLimit),1,0)
func (c *CalcLevel9) calcCheckEnvSumExactToEstim() int {
	if c.Lvl8.REnvTotalExactToEstim >= c.Lvl0.AdvancedParameters.PredefinedCodes.F_PlausiCrit_EnvSum_LowerLimit &&
		c.Lvl8.REnvTotalExactToEstim <= c.Lvl0.AdvancedParameters.PredefinedCodes.F_PlausiCrit_EnvSum_UpperLimit {
		return 1
	}
	return 0
}

// Excel Formula: IF(OR(Code_ThermalBridging_Refurbished="",Code_ThermalBridging_Original=Code_ThermalBridging_Refurbished),Code_ThermalBridging_Original,IF(Code_TypeVariant="Variation",Code_ThermalBridging_Refurbished,IF(Fraction_EnvelopeRefurbished=0,Code_ThermalBridging_Original,IF(Fraction_EnvelopeRefurbished=1,Code_ThermalBridging_Refurbished,Code_ThermalBridging_Original&"("&TEXT(1-Fraction_EnvelopeRefurbished,"##0%")&")."&Code_ThermalBridging_Refurbished&"("&TEXT(Fraction_EnvelopeRefurbished,"##0%")&")"))))
func (c *CalcLevel9) calcTypeThermalBridgingActual() string {
	if c.Lvl0.AdvancedParameters.ThermalBridges.Code_ThermalBridging_Refurbished == "" ||
		c.Lvl0.AdvancedParameters.ThermalBridges.Code_ThermalBridging_Original == c.Lvl0.AdvancedParameters.ThermalBridges.Code_ThermalBridging_Refurbished {
		return c.Lvl0.AdvancedParameters.ThermalBridges.Code_ThermalBridging_Original
	} else if c.Lvl0.BasicParameters.BuildingAppearance.Code_TypeVariant == "Variation" {
		return c.Lvl0.AdvancedParameters.ThermalBridges.Code_ThermalBridging_Refurbished
	} else if c.Lvl8.FractionEnvelopeRefurbished == 0 {
		return c.Lvl0.AdvancedParameters.ThermalBridges.Code_ThermalBridging_Original
	} else if c.Lvl8.FractionEnvelopeRefurbished == 1 {
		return c.Lvl0.AdvancedParameters.ThermalBridges.Code_ThermalBridging_Refurbished
	} else {
		originalPercent := (1 - c.Lvl8.FractionEnvelopeRefurbished) * 100
		refurbishedPercent := c.Lvl8.FractionEnvelopeRefurbished * 100
		return fmt.Sprintf("%s (%.0f%%). %s (%.0f%%)",
			c.Lvl0.AdvancedParameters.ThermalBridges.Code_ThermalBridging_Original,
			originalPercent,
			c.Lvl0.AdvancedParameters.ThermalBridges.Code_ThermalBridging_Refurbished,
			refurbishedPercent)
	}
}

// Excel Formula: IFERROR(IF(Code_ThermalBridging_Refurbished<>"",IF(Code_TypeVariant="Variation",delta_U_ThermalBridging_Refurbished,(1-Fraction_EnvelopeRefurbished)*delta_U_ThermalBridging_Original+Fraction_EnvelopeRefurbished*delta_U_ThermalBridging_Refurbished),delta_U_ThermalBridging_Original),0)
func (c *CalcLevel9) calcDeltaUThermalBridging() float64 {
	if c.Lvl0.AdvancedParameters.ThermalBridges.Code_ThermalBridging_Refurbished != "" {
		if c.Lvl0.BasicParameters.BuildingAppearance.Code_TypeVariant == "Variation" {
			return c.Lvl0.AdvancedParameters.ThermalBridges.Delta_U_ThermalBridging_Refurbished
		} else {
			return (1-c.Lvl8.FractionEnvelopeRefurbished)*c.Lvl0.AdvancedParameters.ThermalBridges.Delta_U_ThermalBridging_Original +
				c.Lvl8.FractionEnvelopeRefurbished*c.Lvl0.AdvancedParameters.ThermalBridges.Delta_U_ThermalBridging_Refurbished
		}
	}
	return c.Lvl0.AdvancedParameters.ThermalBridges.Delta_U_ThermalBridging_Original
}
