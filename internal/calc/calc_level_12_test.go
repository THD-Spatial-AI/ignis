package calc

import "testing"

func TestCalcLevel12_calcFRedTemp(t *testing.T) {
	tests := []struct {
		name         string
		hTransmission float64
		fRedHtr1     float64
		fRedHtr4     float64
		want         float64
	}{
		{
			name:          "h <= 1: lower branch",
			hTransmission: 0.5,
			fRedHtr1:      0.8, fRedHtr4: 0.5,
			// 0.8 + (1-0.5)/0.5*(1-0.8) = 0.8 + 1.0*0.2 = 1.0
			want: 1.0,
		},
		{
			name:          "h >= 4: upper branch",
			hTransmission: 5.0,
			fRedHtr1:      0.8, fRedHtr4: 0.5,
			want: 0.5,
		},
		{
			name:          "1 < h < 4: middle branch",
			hTransmission: 2.5,
			fRedHtr1:      0.8, fRedHtr4: 0.5,
			// 0.8 + (2.5-1)*(0.5-0.8)/(4-1) = 0.8 + 1.5*(-0.1) = 0.65
			want: 0.65,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := newTestParams()
			p.AdvancedParameters.HeatTransfer.F_red_htr1 = tt.fRedHtr1
			p.AdvancedParameters.HeatTransfer.F_red_htr4 = tt.fRedHtr4
			lvl1 := &CalcLevel1{}
			lvl11 := &CalcLevel11{HTransmission: tt.hTransmission}
			c := &CalcLevel12{Lvl0: p, Lvl1: lvl1, Lvl11: lvl11}
			if got := c.calcFRedTemp(); !approxEqual(got, tt.want) {
				t.Errorf("got %.6f, want %.6f", got, tt.want)
			}
		})
	}
}

func TestCalcLevel12_calcTau(t *testing.T) {
	p := newTestParams()
	p.AdvancedParameters.HeatTransfer.C_m = 165
	lvl1 := &CalcLevel1{H_Ventilation: 0.5}
	lvl11 := &CalcLevel11{HTransmission: 1.0}
	c := &CalcLevel12{Lvl0: p, Lvl1: lvl1, Lvl11: lvl11}
	// 165 / (1.0 + 0.5) = 110
	if got := c.calcTau(); !approxEqual(got, 110) {
		t.Errorf("got %.4f, want 110", got)
	}
}
