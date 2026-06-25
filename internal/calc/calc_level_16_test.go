package calc

import (
	"math"
	"testing"
)

func TestCalcLevel16_calcEtaHGn_normalCase(t *testing.T) {
	// gamma = 0.5, aH = 2.0
	// numerator = 1 - 0.5^2 = 0.75
	// denominator = 1 - 0.5^3 = 0.875
	// result = 0.75/0.875 ≈ 0.857142...
	lvl13 := &CalcLevel13{AH: 2.0}
	lvl15 := &CalcLevel15{GammaHGn: 0.5}
	c := &CalcLevel16{Lvl13: lvl13, Lvl15: lvl15}
	want := (1 - math.Pow(0.5, 2.0)) / (1 - math.Pow(0.5, 3.0))
	if got := c.calcEtaHGn(); !approxEqual(got, want) {
		t.Errorf("got %.6f, want %.6f", got, want)
	}
}

func TestCalcLevel16_calcEtaHGn_zeroDenominator(t *testing.T) {
	// gamma = 1.0 → denominator = 1 - 1^(aH+1) = 0
	lvl13 := &CalcLevel13{AH: 2.0}
	lvl15 := &CalcLevel15{GammaHGn: 1.0}
	c := &CalcLevel16{Lvl13: lvl13, Lvl15: lvl15}
	if got := c.calcEtaHGn(); got != 0 {
		t.Errorf("got %.4f, want 0", got)
	}
}
