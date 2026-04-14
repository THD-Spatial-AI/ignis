package calc

import (
	"github.com/THD-Spatial-AI/hdcp-go/internal/models"
)

// CalcLevel1 holds calculated values for level 1
type CalcLevel1 struct {
	Lvl0 *models.TabulaBuildingParameters

	// Calculated attributes
	A_C_Ref_Estim                                     float64 `json:"A_C_Ref_Estim"`
	P_Roof                                            float64 `json:"p_Roof"`
	Q_Roof                                            float64 `json:"q_Roof"`
	P_Ceiling                                         float64 `json:"p_Ceiling"`
	Q_Ceiling                                         float64 `json:"q_Ceiling"`
	F_AtticCond                                       float64 `json:"f_AtticCond"`
	F_CellarCond                                      float64 `json:"f_CellarCond"`
	F_ComplexRoof                                     float64 `json:"f_ComplexRoof"`
	F_ComplexFootprint                                float64 `json:"f_ComplexFootprint"`
	Code_Estim_ConstructionBorder_Wall_ToCellarOrSoil string  `json:"Code_Estim_ConstructionBorder_Wall_ToCellarOrSoil"`
	Code_Estim_ConstructionBorder_Floor               string  `json:"Code_Estim_ConstructionBorder_Floor"`
	A_Exact_Env_Sum                                   float64 `json:"A_Exact_Env_Sum"`
	A_Calc_Window_2                                   float64 `json:"A_Calc_Window_2"`
	A_Calc_Window_Horizontal                          float64 `json:"A_Calc_Window_Horizontal"`
	A_Calc_Window_South                               float64 `json:"A_Calc_Window_South"`
	A_Calc_Window_North                               float64 `json:"A_Calc_Window_North"`
	D_Insulation_Measure_Roof_1                       float64 `json:"d_Insulation_Measure_Roof_1"`
	D_Insulation_Measure_Roof_2                       float64 `json:"d_Insulation_Measure_Roof_2"`
	D_Insulation_Measure_Wall_1                       float64 `json:"d_Insulation_Measure_Wall_1"`
	D_Insulation_Measure_Wall_2                       float64 `json:"d_Insulation_Measure_Wall_2"`
	D_Insulation_Measure_Wall_3                       float64 `json:"d_Insulation_Measure_Wall_3"`
	D_Insulation_Measure_Floor_1                      float64 `json:"d_Insulation_Measure_Floor_1"`
	D_Insulation_Measure_Floor_2                      float64 `json:"d_Insulation_Measure_Floor_2"`
	R_Measure_Window_1                                float64 `json:"R_Measure_Window_1"`
	R_Measure_Window_2                                float64 `json:"R_Measure_Window_2"`
	R_Measure_Door_1                                  float64 `json:"R_Measure_Door_1"`
	G_gl_n_Measure_Window_1                           float64 `json:"g_gl_n_Measure_Window_1"`
	G_gl_n_Measure_Window_2                           float64 `json:"g_gl_n_Measure_Window_2"`
	R_Before_Window_1                                 float64 `json:"R_Before_Window_1"`
	R_Before_Window_2                                 float64 `json:"R_Before_Window_2"`
	R_Before_Door_1                                   float64 `json:"R_Before_Door_1"`
	H_Ventilation                                     float64 `json:"h_Ventilation"`
	Sum_DeltaT_for_HeatingDays                        float64 `json:"Sum_DeltaT_for_HeatingDays"`
	Q_int                                             float64 `json:"q_int"`
}

// NewCalcLevel1 creates a new CalcLevel1 instance and runs all calculations
func NewCalcLevel1(lvl0 *models.TabulaBuildingParameters) *CalcLevel1 {
	calc := &CalcLevel1{
		Lvl0: lvl0,
	}
	calc.Run()
	return calc
}

// Run executes all calculation methods in CalcLevel1 and stores output in corresponding attributes
func (c *CalcLevel1) Run() {
	c.A_C_Ref_Estim = c.calcACRefEstim()
	c.P_Roof = c.calcPRoof()
	c.Q_Roof = c.calcQRoof()
	c.P_Ceiling = c.calcPCeiling()
	c.Q_Ceiling = c.calcQCeiling()
	c.F_AtticCond = c.calcFAtticCond()
	c.F_CellarCond = c.calcFCellarCond()
	c.F_ComplexRoof = c.calcFComplexRoof()
	c.F_ComplexFootprint = c.calcFComplexFootprint()
	c.Code_Estim_ConstructionBorder_Wall_ToCellarOrSoil = c.calcCodeEstimConstructionBorderWallToCellarOrSoil()
	c.Code_Estim_ConstructionBorder_Floor = c.calcCodeEstimConstructionBorderFloor()
	c.A_Exact_Env_Sum = c.calcAExactEnvSum()
	c.A_Calc_Window_2 = c.calcACalcWindow2()
	c.A_Calc_Window_Horizontal = c.calcACalcWindowHorizontal()
	c.A_Calc_Window_South = c.calcACalcWindowSouth()
	c.A_Calc_Window_North = c.calcACalcWindowNorth()
	c.D_Insulation_Measure_Roof_1 = c.calcDInsulationMeasureRoof1()
	c.D_Insulation_Measure_Roof_2 = c.calcDInsulationMeasureRoof2()
	c.D_Insulation_Measure_Wall_1 = c.calcDInsulationMeasureWall1()
	c.D_Insulation_Measure_Wall_2 = c.calcDInsulationMeasureWall2()
	c.D_Insulation_Measure_Wall_3 = c.calcDInsulationMeasureWall3()
	c.D_Insulation_Measure_Floor_1 = c.calcDInsulationMeasureFloor1()
	c.D_Insulation_Measure_Floor_2 = c.calcDInsulationMeasureFloor2()
	c.R_Measure_Window_1 = c.calcRMeasureWindow1()
	c.R_Measure_Window_2 = c.calcRMeasureWindow2()
	c.R_Measure_Door_1 = c.calcRMeasureDoor1()
	c.G_gl_n_Measure_Window_1 = c.calcGGlNMeasureWindow1()
	c.G_gl_n_Measure_Window_2 = c.calcGGlNMeasureWindow2()
	c.R_Before_Window_1 = c.calcRBeforeWindow1()
	c.R_Before_Window_2 = c.calcRBeforeWindow2()
	c.R_Before_Door_1 = c.calcRBeforeDoor1()
	c.H_Ventilation = c.calcHVentilation()
	c.Sum_DeltaT_for_HeatingDays = c.calcSumDeltaTForHeatingDays()
	c.Q_int = c.calcQInt()
}

// calcACRefEstim calculates A_C_Ref_Estim
// Estimated by use of conversion factors
// Unit: m²
// Formula: IF(A_C_IntDim>0, A_C_IntDim, IF(A_C_ExtDim>0, 0.85*A_C_ExtDim, IF(A_C_Living>0, 1.1*A_C_Living, IF(A_C_Use>0, 1.4*A_C_Use, 0.85/3*V_C))))
func (c *CalcLevel1) calcACRefEstim() float64 {
	if c.Lvl0.BasicParameters.Envelope.A_C_IntDim > 0 {
		return c.Lvl0.BasicParameters.Envelope.A_C_IntDim
	} else if c.Lvl0.BasicParameters.Envelope.A_C_ExtDim > 0 {
		return 0.85 * c.Lvl0.BasicParameters.Envelope.A_C_ExtDim
	} else if c.Lvl0.BasicParameters.Envelope.A_C_Living > 0 {
		return 1.1 * c.Lvl0.BasicParameters.Envelope.A_C_Living
	} else if c.Lvl0.BasicParameters.Envelope.A_C_Use > 0 {
		return 1.4 * c.Lvl0.BasicParameters.Envelope.A_C_Use
	} else {
		return 0.85 / 3 * c.Lvl0.BasicParameters.Envelope.V_C
	}
}

// calcPRoof calculates p_Roof
// Unit: m²/m²
// Formula: IF(OR($Code_AtticCond="C",$Code_AtticCond="PI",$Code_AtticCond="NI"),1.6,IF($Code_AtticCond="P",0.8,IF($Code_AtticCond="-",1.2,0)))
func (c *CalcLevel1) calcPRoof() float64 {
	code := c.Lvl0.BasicParameters.BuildingAppearance.Code_AtticCond
	if code == "C" || code == "PI" || code == "NI" {
		return 1.6
	} else if code == "P" {
		return 0.8
	} else if code == "-" {
		return 1.2
	}
	return 0
}

// calcQRoof calculates q_Roof
// Unit: m²
// Formula: IF(OR($Code_AtticCond="C",$Code_AtticCond="PI",$Code_AtticCond="NI"),15,IF($Code_AtticCond="P",7,IF($Code_AtticCond="-",5,0)))
func (c *CalcLevel1) calcQRoof() float64 {
	code := c.Lvl0.BasicParameters.BuildingAppearance.Code_AtticCond
	if code == "C" || code == "PI" || code == "NI" {
		return 15
	} else if code == "P" {
		return 7
	} else if code == "-" {
		return 5
	}
	return 0
}

// calcPCeiling calculates p_Ceiling
// Unit: m²/m²
// Formula: IF(OR($Code_AtticCond="C",$Code_AtticCond="PI",$Code_AtticCond="NI"),0,IF($Code_AtticCond="P",0.6,IF($Code_AtticCond="-",0,1.2)))
func (c *CalcLevel1) calcPCeiling() float64 {
	code := c.Lvl0.BasicParameters.BuildingAppearance.Code_AtticCond
	if code == "C" || code == "PI" || code == "NI" {
		return 0
	} else if code == "P" {
		return 0.6
	} else if code == "-" {
		return 0
	}
	return 1.2
}

// calcQCeiling calculates q_Ceiling
// Unit: m²
// Formula: IF(OR($Code_AtticCond="C",$Code_AtticCond="PI",$Code_AtticCond="NI"),0,IF($Code_AtticCond="P",3,IF($Code_AtticCond="-",0,5)))
func (c *CalcLevel1) calcQCeiling() float64 {
	code := c.Lvl0.BasicParameters.BuildingAppearance.Code_AtticCond
	if code == "C" || code == "PI" || code == "NI" {
		return 0
	} else if code == "P" {
		return 3
	} else if code == "-" {
		return 0
	}
	return 5
}

// calcFAtticCond calculates f_AtticCond
// Formula: IF(LEFT($Code_AtticCond,1)="C",1,IF(LEFT($Code_AtticCond,1)="P",0.5,0))
func (c *CalcLevel1) calcFAtticCond() float64 {
	code := c.Lvl0.BasicParameters.BuildingAppearance.Code_AtticCond
	if len(code) > 0 {
		firstChar := string(code[0])
		if firstChar == "C" {
			return 1
		} else if firstChar == "P" {
			return 0.5
		}
	}
	return 0
}

// calcFCellarCond calculates f_CellarCond
// Formula: IF(LEFT($Code_CellarCond,1)="C",1,IF(LEFT($Code_CellarCond,1)="P",0.5,0))
func (c *CalcLevel1) calcFCellarCond() float64 {
	code := c.Lvl0.BasicParameters.BuildingAppearance.Code_CellarCond
	if len(code) > 0 {
		firstChar := string(code[0])
		if firstChar == "C" {
			return 1
		} else if firstChar == "P" {
			return 0.5
		}
	}
	return 0
}

// calcFComplexRoof calculates f_ComplexRoof
// Formula: IF(Code_ComplexRoof="Simple",0.9,IF(Code_ComplexRoof="Complex",1.3,1))
func (c *CalcLevel1) calcFComplexRoof() float64 {
	code := c.Lvl0.BasicParameters.BuildingAppearance.Code_ComplexRoof
	if code == "Simple" {
		return 0.9
	} else if code == "Complex" {
		return 1.3
	}
	return 1
}

// calcFComplexFootprint calculates f_ComplexFootprint
// Formula: IF(Code_ComplexFootprint="Simple",0.9,IF(Code_ComplexFootprint="Complex",1.2,1))
func (c *CalcLevel1) calcFComplexFootprint() float64 {
	code := c.Lvl0.BasicParameters.BuildingAppearance.Code_ComplexFootprint
	if code == "Simple" {
		return 0.9
	} else if code == "Complex" {
		return 1.2
	}
	return 1
}

// calcCodeEstimConstructionBorderWallToCellarOrSoil calculates Code_Estim_ConstructionBorder_Wall_ToCellarOrSoil
// wall bordering at soil or unheated cellar
// Formula: IF($Code_CellarCond="P","Unh","Soil")
func (c *CalcLevel1) calcCodeEstimConstructionBorderWallToCellarOrSoil() string {
	if c.Lvl0.BasicParameters.BuildingAppearance.Code_CellarCond == "P" {
		return "Unh"
	}
	return "Soil"
}

// calcCodeEstimConstructionBorderFloor calculates Code_Estim_ConstructionBorder_Floor
// floor above cellar or soil
// Formula: IF($Code_CellarCond="-","Soil","Cellar")
func (c *CalcLevel1) calcCodeEstimConstructionBorderFloor() string {
	if c.Lvl0.BasicParameters.BuildingAppearance.Code_CellarCond == "-" {
		return "Soil"
	}
	return "Cellar"
}

// calcAExactEnvSum calculates A_Exact_Env_Sum
// Calculates the sum of various building element areas
// Unit: m²
// Formula: SUM(A_Roof_1:A_Door_1)
func (c *CalcLevel1) calcAExactEnvSum() float64 {
	envelope := c.Lvl0.BasicParameters.Envelope
	return envelope.A_Roof_1 +
		envelope.A_Roof_2 +
		envelope.A_Wall_1 +
		envelope.A_Wall_2 +
		envelope.A_Wall_3 +
		envelope.A_Floor_1 +
		envelope.A_Floor_2 +
		envelope.A_Window_1 +
		envelope.A_Window_2 +
		envelope.A_Door_1
}

// calcACalcWindow2 calculates A_Calc_Window_2
// element type window 2
// Unit: m²
// Formula: IF($Code_TypeIntake_EnvelopeArea="Estimation",0,A_Window_2)
func (c *CalcLevel1) calcACalcWindow2() float64 {
	if c.Lvl0.AdvancedParameters.PredefinedCodes.Code_TypeIntake_EnvelopeArea == "Estimation" {
		return 0
	}
	return c.Lvl0.BasicParameters.Envelope.A_Window_2
}

// calcACalcWindowHorizontal calculates A_Calc_Window_Horizontal
// tilted below 30°, otherwise classified as vertical
// Unit: m²
// Formula: IF($Code_TypeIntake_EnvelopeArea="Estimation",0,A_Window_Horizontal)
func (c *CalcLevel1) calcACalcWindowHorizontal() float64 {
	if c.Lvl0.AdvancedParameters.PredefinedCodes.Code_TypeIntake_EnvelopeArea == "Estimation" {
		return 0
	}
	return c.Lvl0.BasicParameters.Envelope.A_Window_Horizontal
}

// calcACalcWindowSouth calculates A_Calc_Window_South
// deviation from orientation: +/- 45°
// Unit: m²
// Formula: IF($Code_TypeIntake_EnvelopeArea="Estimation",0,A_Window_South)
func (c *CalcLevel1) calcACalcWindowSouth() float64 {
	if c.Lvl0.AdvancedParameters.PredefinedCodes.Code_TypeIntake_EnvelopeArea == "Estimation" {
		return 0
	}
	return c.Lvl0.BasicParameters.Envelope.A_Window_South
}

// calcACalcWindowNorth calculates A_Calc_Window_North
// deviation from orientation: +/- 45°
// Unit: m²
// Formula: IF($Code_TypeIntake_EnvelopeArea="Estimation",0,A_Window_North)
func (c *CalcLevel1) calcACalcWindowNorth() float64 {
	if c.Lvl0.AdvancedParameters.PredefinedCodes.Code_TypeIntake_EnvelopeArea == "Estimation" {
		return 0
	}
	return c.Lvl0.BasicParameters.Envelope.A_Window_North
}

// calcDInsulationMeasureRoof1 calculates d_Insulation_Measure_Roof_1
// Element type roof 1
// Unit: m
// Formula: IF(d_Insulation_Input_Measure_Roof_1<>0, d_Insulation_Input_Measure_Roof_1, d_Insulation_PredefinedMeasure_Roof_1)
func (c *CalcLevel1) calcDInsulationMeasureRoof1() float64 {
	if c.Lvl0.AdvancedParameters.ActualInsulation.D_Insulation_Input_Measure_Roof_1 != 0 {
		return c.Lvl0.AdvancedParameters.ActualInsulation.D_Insulation_Input_Measure_Roof_1
	}
	return c.Lvl0.AdvancedParameters.InsulationMeasures.D_Insulation_PredefinedMeasure_Roof_1
}

// calcDInsulationMeasureRoof2 calculates d_Insulation_Measure_Roof_2
// Element type roof 2
// Unit: m
// Formula: IF(d_Insulation_Input_Measure_Roof_2<>0, d_Insulation_Input_Measure_Roof_2, d_Insulation_PredefinedMeasure_Roof_2)
func (c *CalcLevel1) calcDInsulationMeasureRoof2() float64 {
	if c.Lvl0.AdvancedParameters.ActualInsulation.D_Insulation_Input_Measure_Roof_2 != 0 {
		return c.Lvl0.AdvancedParameters.ActualInsulation.D_Insulation_Input_Measure_Roof_2
	}
	return c.Lvl0.AdvancedParameters.InsulationMeasures.D_Insulation_PredefinedMeasure_Roof_2
}

// calcDInsulationMeasureWall1 calculates d_Insulation_Measure_Wall_1
// Element type wall 1
// Unit: m
// Formula: IF(d_Insulation_Input_Measure_Wall_1<>0, d_Insulation_Input_Measure_Wall_1, d_Insulation_PredefinedMeasure_Wall_1)
func (c *CalcLevel1) calcDInsulationMeasureWall1() float64 {
	if c.Lvl0.AdvancedParameters.ActualInsulation.D_Insulation_Input_Measure_Wall_1 != 0 {
		return c.Lvl0.AdvancedParameters.ActualInsulation.D_Insulation_Input_Measure_Wall_1
	}
	return c.Lvl0.AdvancedParameters.InsulationMeasures.D_Insulation_PredefinedMeasure_Wall_1
}

// calcDInsulationMeasureWall2 calculates d_Insulation_Measure_Wall_2
// Element type wall 2
// Unit: m
// Formula: IF(d_Insulation_Input_Measure_Wall_2<>0, d_Insulation_Input_Measure_Wall_2, d_Insulation_PredefinedMeasure_Wall_2)
func (c *CalcLevel1) calcDInsulationMeasureWall2() float64 {
	if c.Lvl0.AdvancedParameters.ActualInsulation.D_Insulation_Input_Measure_Wall_2 != 0 {
		return c.Lvl0.AdvancedParameters.ActualInsulation.D_Insulation_Input_Measure_Wall_2
	}
	return c.Lvl0.AdvancedParameters.InsulationMeasures.D_Insulation_PredefinedMeasure_Wall_2
}

// calcDInsulationMeasureWall3 calculates d_Insulation_Measure_Wall_3
// Element type wall 3
// Unit: m
// Formula: IF(d_Insulation_Input_Measure_Wall_3<>0, d_Insulation_Input_Measure_Wall_3, d_Insulation_PredefinedMeasure_Wall_3)
func (c *CalcLevel1) calcDInsulationMeasureWall3() float64 {
	if c.Lvl0.AdvancedParameters.ActualInsulation.D_Insulation_Input_Measure_Wall_3 != 0 {
		return c.Lvl0.AdvancedParameters.ActualInsulation.D_Insulation_Input_Measure_Wall_3
	}
	return c.Lvl0.AdvancedParameters.InsulationMeasures.D_Insulation_PredefinedMeasure_Wall_3
}

// calcDInsulationMeasureFloor1 calculates d_Insulation_Measure_Floor_1
// Element type floor 1
// Unit: m
// Formula: IF(d_Insulation_Input_Measure_Floor_1<>0, d_Insulation_Input_Measure_Floor_1, d_Insulation_PredefinedMeasure_Floor_1)
func (c *CalcLevel1) calcDInsulationMeasureFloor1() float64 {
	if c.Lvl0.AdvancedParameters.ActualInsulation.D_Insulation_Input_Measure_Floor_1 != 0 {
		return c.Lvl0.AdvancedParameters.ActualInsulation.D_Insulation_Input_Measure_Floor_1
	}
	return c.Lvl0.AdvancedParameters.InsulationMeasures.D_Insulation_PredefinedMeasure_Floor_1
}

// calcDInsulationMeasureFloor2 calculates d_Insulation_Measure_Floor_2
// Element type floor 2
// Unit: m
// Formula: IF(d_Insulation_Input_Measure_Floor_2<>0, d_Insulation_Input_Measure_Floor_2, d_Insulation_PredefinedMeasure_Floor_2)
func (c *CalcLevel1) calcDInsulationMeasureFloor2() float64 {
	if c.Lvl0.AdvancedParameters.ActualInsulation.D_Insulation_Input_Measure_Floor_2 != 0 {
		return c.Lvl0.AdvancedParameters.ActualInsulation.D_Insulation_Input_Measure_Floor_2
	}
	return c.Lvl0.AdvancedParameters.InsulationMeasures.D_Insulation_PredefinedMeasure_Floor_2
}

// calcRMeasureWindow1 calculates R_Measure_Window_1
// Element type window 1
// Unit: m²*K/W
// Formula: R_PredefinedMeasure_Window_1
func (c *CalcLevel1) calcRMeasureWindow1() float64 {
	return c.Lvl0.AdvancedParameters.ThermalResistances.R_PredefinedMeasure_Window_1
}

// calcRMeasureWindow2 calculates R_Measure_Window_2
// Element type window 2
// Unit: m²*K/W
// Formula: R_PredefinedMeasure_Window_2
func (c *CalcLevel1) calcRMeasureWindow2() float64 {
	return c.Lvl0.AdvancedParameters.ThermalResistances.R_PredefinedMeasure_Window_2
}

// calcRMeasureDoor1 calculates R_Measure_Door_1
// Element type door 1
// Unit: m²*K/W
// Formula: R_PredefinedMeasure_Door_1
func (c *CalcLevel1) calcRMeasureDoor1() float64 {
	return c.Lvl0.AdvancedParameters.ThermalResistances.R_PredefinedMeasure_Door_1
}

// calcGGlNMeasureWindow1 calculates g_gl_n_Measure_Window_1
// Element type window 1
// Unit: m²*K/W
// Formula: g_gl_n_PredefinedMeasure_Window_1
func (c *CalcLevel1) calcGGlNMeasureWindow1() float64 {
	return c.Lvl0.AdvancedParameters.SolarTransmittance.G_gl_n_PredefinedMeasure_Window_1
}

// calcGGlNMeasureWindow2 calculates g_gl_n_Measure_Window_2
// Element type window 2
// Unit: m²*K/W
// Formula: g_gl_n_PredefinedMeasure_Window_2
func (c *CalcLevel1) calcGGlNMeasureWindow2() float64 {
	return c.Lvl0.AdvancedParameters.SolarTransmittance.G_gl_n_PredefinedMeasure_Window_2
}

// calcRBeforeWindow1 calculates R_Before_Window_1
// Element type window 1
// Unit: m²*K/W
// Formula: IFERROR(IF(ISNUMBER(U_Window_1), 1/U_Window_1, 0), 0)
func (c *CalcLevel1) calcRBeforeWindow1() float64 {
	uValue := c.Lvl0.AdvancedParameters.Uvalues.U_Window_1
	if uValue != 0 {
		return 1 / uValue
	}
	return 0
}

// calcRBeforeWindow2 calculates R_Before_Window_2
// Element type window 2
// Unit: m²*K/W
// Formula: IFERROR(IF(ISNUMBER(U_Window_2), 1/U_Window_2, 0), 0)
func (c *CalcLevel1) calcRBeforeWindow2() float64 {
	uValue := c.Lvl0.AdvancedParameters.Uvalues.U_Window_2
	if uValue != 0 {
		return 1 / uValue
	}
	return 0
}

// calcRBeforeDoor1 calculates R_Before_Door_1
// Element type door 1
// Unit: m²*K/W
// Formula: IFERROR(IF(ISNUMBER(U_Door_1), 1/U_Door_1, 0), 0)
func (c *CalcLevel1) calcRBeforeDoor1() float64 {
	uValue := c.Lvl0.AdvancedParameters.Uvalues.U_Door_1
	if uValue != 0 {
		return 1 / uValue
	}
	return 0
}

// calcHVentilation calculates h_Ventilation
// Unit: W/(m²*K)
// Formula: 0.34*(n_air_use + n_air_infiltration)*h_room
func (c *CalcLevel1) calcHVentilation() float64 {
	return 0.34 * (c.Lvl0.AdvancedParameters.AirInfiltration.N_air_use + c.Lvl0.AdvancedParameters.AirInfiltration.N_air_infiltration) * c.Lvl0.BasicParameters.BuildingAppearance.H_room
}

// calcSumDeltaTForHeatingDays calculates Sum_DeltaT_for_HeatingDays
// Unit: kKh/a
// Formula: (theta_i - Theta_e) * HeatingDays
func (c *CalcLevel1) calcSumDeltaTForHeatingDays() float64 {
	return (c.Lvl0.AdvancedParameters.ClimateConditions.Theta_i - c.Lvl0.AdvancedParameters.ClimateConditions.Theta_e) * float64(c.Lvl0.AdvancedParameters.ClimateConditions.HeatingDays)
}

// calcQInt calculates q_int
// Unit: kWh/(m²*a)
// Formula: phi_int * 0.024 * HeatingDays
func (c *CalcLevel1) calcQInt() float64 {
	return c.Lvl0.AdvancedParameters.HeatTransfer.Phi_int * 0.024 * float64(c.Lvl0.AdvancedParameters.ClimateConditions.HeatingDays)
}
