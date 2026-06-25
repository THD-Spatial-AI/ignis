package calc

import "testing"

func TestCalcLevel9_calcCheckEnvSumExactToEstim_within(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.PredefinedCodes.F_PlausiCrit_EnvSum_LowerLimit = 0.7
	p.AdvancedParameters.PredefinedCodes.F_PlausiCrit_EnvSum_UpperLimit = 1.3
	lvl8 := &CalcLevel8{REnvTotalExactToEstim: 1.0}
	c := &CalcLevel9{Lvl0: p, Lvl8: lvl8}
	if got := c.calcCheckEnvSumExactToEstim(); got != 1 {
		t.Errorf("got %d, want 1", got)
	}
}

func TestCalcLevel9_calcCheckEnvSumExactToEstim_outside(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.PredefinedCodes.F_PlausiCrit_EnvSum_LowerLimit = 0.7
	p.AdvancedParameters.PredefinedCodes.F_PlausiCrit_EnvSum_UpperLimit = 1.3
	lvl8 := &CalcLevel8{REnvTotalExactToEstim: 2.0}
	c := &CalcLevel9{Lvl0: p, Lvl8: lvl8}
	if got := c.calcCheckEnvSumExactToEstim(); got != 0 {
		t.Errorf("got %d, want 0", got)
	}
}

func TestCalcLevel9_calcTypeThermalBridgingActual(t *testing.T) {
	tests := []struct {
		name         string
		original     string
		refurbished  string
		typeVariant  string
		fraction     float64
		want         string
	}{
		{
			name:        "refurbished empty → return original",
			original:    "Orig", refurbished: "",
			want: "Orig",
		},
		{
			name:        "original == refurbished → return original",
			original:    "Same", refurbished: "Same",
			want: "Same",
		},
		{
			name:        "Variation → return refurbished",
			original:    "Orig", refurbished: "Refurb", typeVariant: "Variation",
			want: "Refurb",
		},
		{
			name:        "fraction == 0 → return original",
			original:    "Orig", refurbished: "Refurb", typeVariant: "Standard",
			fraction: 0, want: "Orig",
		},
		{
			name:        "fraction == 1 → return refurbished",
			original:    "Orig", refurbished: "Refurb", typeVariant: "Standard",
			fraction: 1, want: "Refurb",
		},
		{
			name:        "partial refurbishment → formatted string",
			original:    "Orig", refurbished: "Refurb", typeVariant: "Standard",
			fraction: 0.5, want: "Orig (50%). Refurb (50%)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := newTestParams()
			p.AdvancedParameters.ThermalBridges.Code_ThermalBridging_Original = tt.original
			p.AdvancedParameters.ThermalBridges.Code_ThermalBridging_Refurbished = tt.refurbished
			p.BasicParameters.BuildingAppearance.Code_TypeVariant = tt.typeVariant
			lvl8 := &CalcLevel8{FractionEnvelopeRefurbished: tt.fraction}
			c := &CalcLevel9{Lvl0: p, Lvl8: lvl8}
			if got := c.calcTypeThermalBridgingActual(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCalcLevel9_calcDeltaUThermalBridging(t *testing.T) {
	tests := []struct {
		name        string
		refurbished string
		typeVariant string
		fraction    float64
		deltaOrig   float64
		deltaRefurb float64
		want        float64
	}{
		{
			name: "no refurbished code → return original",
			refurbished: "", deltaOrig: 0.1, want: 0.1,
		},
		{
			name: "Variation → return refurbished delta",
			refurbished: "Refurb", typeVariant: "Variation",
			deltaOrig: 0.1, deltaRefurb: 0.05, want: 0.05,
		},
		{
			name: "mixed fraction",
			refurbished: "Refurb", typeVariant: "Standard",
			fraction: 0.5, deltaOrig: 0.2, deltaRefurb: 0.1,
			want: 0.5*0.2 + 0.5*0.1, // 0.15
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := newTestParams()
			p.AdvancedParameters.ThermalBridges.Code_ThermalBridging_Refurbished = tt.refurbished
			p.AdvancedParameters.ThermalBridges.Delta_U_ThermalBridging_Original = tt.deltaOrig
			p.AdvancedParameters.ThermalBridges.Delta_U_ThermalBridging_Refurbished = tt.deltaRefurb
			p.BasicParameters.BuildingAppearance.Code_TypeVariant = tt.typeVariant
			lvl8 := &CalcLevel8{FractionEnvelopeRefurbished: tt.fraction}
			c := &CalcLevel9{Lvl0: p, Lvl8: lvl8}
			if got := c.calcDeltaUThermalBridging(); !approxEqual(got, tt.want) {
				t.Errorf("got %.6f, want %.6f", got, tt.want)
			}
		})
	}
}
