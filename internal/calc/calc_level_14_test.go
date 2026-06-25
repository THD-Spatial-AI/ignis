package calc

import "testing"

func TestCalcLevel14_calcQHt(t *testing.T) {
	lvl13 := &CalcLevel13{QHtTr: 108, QHtVe: 54}
	c := &CalcLevel14{Lvl13: lvl13}
	if got := c.calcQHt(); !approxEqual(got, 162) {
		t.Errorf("got %.4f, want 162", got)
	}
}
