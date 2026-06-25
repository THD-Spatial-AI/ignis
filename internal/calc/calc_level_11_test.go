package calc

import "testing"

func TestCalcLevel11_calcHTransmission_nonZeroRef(t *testing.T) {
	lvl2 := &CalcLevel2{A_C_Ref: 100}
	lvl5 := &CalcLevel5{H_Transmission_Door_1: 5}
	lvl6 := &CalcLevel6{
		HTransmissionRoof1: 10, HTransmissionRoof2: 8,
		HTransmissionFloor1: 12, HTransmissionWindow1: 6, HTransmissionWindow2: 4,
	}
	lvl7 := &CalcLevel7{HTransmissionFloor2: 7, HTransmissionWall2: 9}
	lvl8 := &CalcLevel8{HTransmissionWall1: 15, HTransmissionWall3: 6}
	lvl10 := &CalcLevel10{HTransmissionThermalBridging: 3}
	c := &CalcLevel11{
		Lvl2: lvl2, Lvl5: lvl5, Lvl6: lvl6,
		Lvl7: lvl7, Lvl8: lvl8, Lvl10: lvl10,
	}
	// 12+7+10+8+6+4+15+9+6+5+3 = 85, / 100 = 0.85
	if got := c.calcHTransmission(); !approxEqual(got, 0.85) {
		t.Errorf("got %.4f, want 0.85", got)
	}
}

func TestCalcLevel11_calcHTransmission_zeroRef(t *testing.T) {
	lvl2 := &CalcLevel2{A_C_Ref: 0}
	c := &CalcLevel11{
		Lvl2: lvl2, Lvl5: &CalcLevel5{}, Lvl6: &CalcLevel6{},
		Lvl7: &CalcLevel7{}, Lvl8: &CalcLevel8{}, Lvl10: &CalcLevel10{},
	}
	if got := c.calcHTransmission(); got != 0 {
		t.Errorf("got %.4f, want 0", got)
	}
}
