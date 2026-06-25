package calc

import "testing"

func TestCalcLevel15_calcGammaHGn_nonZeroQHt(t *testing.T) {
	lvl1 := &CalcLevel1{Q_int: 10}
	lvl8 := &CalcLevel8{QSol: 20}
	lvl14 := &CalcLevel14{QHt: 100}
	c := &CalcLevel15{Lvl1: lvl1, Lvl8: lvl8, Lvl14: lvl14}
	// (20+10)/100 = 0.3
	if got := c.calcGammaHGn(); !approxEqual(got, 0.3) {
		t.Errorf("got %.4f, want 0.3", got)
	}
}

func TestCalcLevel15_calcGammaHGn_zeroQHt(t *testing.T) {
	lvl1 := &CalcLevel1{Q_int: 10}
	lvl8 := &CalcLevel8{QSol: 20}
	lvl14 := &CalcLevel14{QHt: 0}
	c := &CalcLevel15{Lvl1: lvl1, Lvl8: lvl8, Lvl14: lvl14}
	if got := c.calcGammaHGn(); got != 0 {
		t.Errorf("got %.4f, want 0", got)
	}
}
