package calc

import "testing"

func TestCalcLevel17_calcQHNd(t *testing.T) {
	lvl1 := &CalcLevel1{Q_int: 10}
	lvl8 := &CalcLevel8{QSol: 20}
	lvl14 := &CalcLevel14{QHt: 162}
	lvl16 := &CalcLevel16{EtaHGn: 0.8}
	c := &CalcLevel17{Lvl1: lvl1, Lvl8: lvl8, Lvl14: lvl14, Lvl16: lvl16}
	// 162 - 0.8*(20+10) = 162 - 24 = 138
	if got := c.calcQHNd(); !approxEqual(got, 138) {
		t.Errorf("got %.4f, want 138", got)
	}
}

func TestCalcLevel17_Run_returnsQHNd(t *testing.T) {
	lvl1 := &CalcLevel1{Q_int: 10}
	lvl8 := &CalcLevel8{QSol: 20}
	lvl14 := &CalcLevel14{QHt: 162}
	lvl16 := &CalcLevel16{EtaHGn: 0.8}
	c := NewCalcLevel17(lvl1, lvl8, lvl14, lvl16)
	if c.QHNd != 138 {
		t.Errorf("QHNd field: got %.4f, want 138", c.QHNd)
	}
}
