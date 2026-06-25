package calc

import "testing"

func TestCalcLevel10_calcCheckEnvAreaExactToEstim(t *testing.T) {
	tests := []struct {
		name                  string
		envSum                int
		toBeAppliedFloorArea  int
		floorCheck            int
		windowCheck           int
		want                  int
	}{
		{
			name: "all checks pass → 1",
			envSum: 1, toBeAppliedFloorArea: 1, floorCheck: 1, windowCheck: 1,
			want: 1,
		},
		{
			name: "env sum fails → 0",
			envSum: 0, toBeAppliedFloorArea: 1, floorCheck: 1, windowCheck: 1,
			want: 0,
		},
		{
			name: "floor check skipped when toBeApplied==0",
			envSum: 1, toBeAppliedFloorArea: 0, floorCheck: 0, windowCheck: 1,
			want: 1, // floor not applied, window OK → 1*1 = 1
		},
		{
			name: "window check fails → 0",
			envSum: 1, toBeAppliedFloorArea: 0, floorCheck: 1, windowCheck: 0,
			want: 0,
		},
		{
			// Artificial: result > 1 → clamped to 1
			name: "result > 1 clamped to 1",
			envSum: 2, toBeAppliedFloorArea: 0, floorCheck: 0, windowCheck: 1,
			want: 1,
		},
		{
			// Artificial: result < 0 → clamped to -1
			name: "result < 0 clamped to -1",
			envSum: -1, toBeAppliedFloorArea: 0, floorCheck: 0, windowCheck: 1,
			want: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CalcLevel10{
				Lvl9: &CalcLevel9{CheckEnvSumExactToEstim: tt.envSum},
				Lvl2: &CalcLevel2{Check_ToBeApplied_FloorArea_ExactToEstim: tt.toBeAppliedFloorArea},
				Lvl6: &CalcLevel6{
					CheckFloorAreaExactToEstim:  tt.floorCheck,
					CheckWindowAreaExactToEstim: tt.windowCheck,
				},
				Lvl1: &CalcLevel1{}, Lvl4: &CalcLevel4{},
				Lvl5: &CalcLevel5{}, Lvl7: &CalcLevel7{},
			}
			if got := c.calcCheckEnvAreaExactToEstim(); got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestCalcLevel10_calcHTransmissionThermalBridging(t *testing.T) {
	lvl1 := &CalcLevel1{A_Calc_Window_2: 5}
	lvl5 := &CalcLevel5{A_Calc_Roof_1: 10, A_Calc_Roof_2: 10, A_Calc_Floor_1: 10, A_Calc_Window_1: 5}
	lvl6 := &CalcLevel6{ACalcFloor2: 10, ACalcWall2: 10}
	lvl7 := &CalcLevel7{ACalcWall1: 20, ACalcWall3: 10}
	lvl4 := &CalcLevel4{A_Calc_Door_1: 2}
	lvl9 := &CalcLevel9{DeltaUThermalBridging: 0.1}
	c := &CalcLevel10{
		Lvl1: lvl1, Lvl5: lvl5, Lvl6: lvl6,
		Lvl7: lvl7, Lvl4: lvl4, Lvl9: lvl9,
	}
	// Total area = 10+10+10+10+10+20+10+5+5+2 = 92, * 0.1 = 9.2
	if got := c.calcHTransmissionThermalBridging(); !approxEqual(got, 9.2) {
		t.Errorf("got %.4f, want 9.2", got)
	}
}
