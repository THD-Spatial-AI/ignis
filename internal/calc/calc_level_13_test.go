package calc

import "testing"

func TestCalcLevel13_calcQHtTr(t *testing.T) {
	lvl1 := &CalcLevel1{Sum_DeltaT_for_HeatingDays: 5000}
	lvl11 := &CalcLevel11{HTransmission: 1.0}
	lvl12 := &CalcLevel12{FRedTemp: 0.9}
	c := &CalcLevel13{Lvl1: lvl1, Lvl11: lvl11, Lvl12: lvl12}
	// 1.0 * 0.024 * 5000 * 0.9 = 108
	if got := c.calcQHtTr(); !approxEqual(got, 108) {
		t.Errorf("got %.4f, want 108", got)
	}
}

func TestCalcLevel13_calcQHtVe(t *testing.T) {
	lvl1 := &CalcLevel1{H_Ventilation: 0.5, Sum_DeltaT_for_HeatingDays: 5000}
	lvl12 := &CalcLevel12{FRedTemp: 0.9}
	c := &CalcLevel13{Lvl1: lvl1, Lvl12: lvl12}
	// 0.5 * 0.024 * 5000 * 0.9 = 54
	if got := c.calcQHtVe(); !approxEqual(got, 54) {
		t.Errorf("got %.4f, want 54", got)
	}
}

func TestCalcLevel13_calcAH(t *testing.T) {
	lvl12 := &CalcLevel12{Tau: 60}
	c := &CalcLevel13{Lvl12: lvl12}
	// 0.8 + 60/30 = 2.8
	if got := c.calcAH(); !approxEqual(got, 2.8) {
		t.Errorf("got %.4f, want 2.8", got)
	}
}
