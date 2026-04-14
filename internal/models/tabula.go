package models

// --- Constants & Enumerations ---
var (
	CodeConditionStatus  = []string{"-", "N", "P", "C", "NI", "PI"}  // Attic/Cellar conditioning
	NeighbourStatus      = []string{"B_Alone", "B_N1", "B_N2"}       // Attached neighbours
	FootprintShapeStatus = []string{"Simple", "Standard", "Complex"} // Complexity of building footprint
	RoofShapeStatus      = []string{"Simple", "Standard", "Complex"} // Complexity of roof shape
)

const (
	// Default room height in meters for calculating conditioned volume if not provided in dataset
	RoomHeight = 2.5
)

// --- Structs ---

// TabulaBuildingParameters represents the building parameters from Tabula
type TabulaBuildingParameters struct {
	BasicParameters    *BasicParameters    // Basic building parameters (thematic and envelope parameters)
	AdvancedParameters *AdvancedParameters // Advanced building parameters (air infiltration, climate conditions, U-values, insulation thicknesses, solar gains, thermal bridging, heat losses, etc.)
}

// BasicParameters holds the building's thematic and envelope parameters
type BasicParameters struct {
	BuildingAppearance *BuildingThematic // Thematic building parameters (e.g. building variant code, type variant code, number of storeys, etc.)
	Envelope           *Envelope         // Building envelope parameters (e.g. areas of different building elements, conditioned volume, etc.)
}

// AdvancedParameters holds advanced building parameters
type AdvancedParameters struct {
	AirInfiltration       *AirInfiltration              // Air infiltration and usage parameters
	ClimateConditions     *ClimateConditions            // Climate-related parameters
	Uvalues               *Uvalues                      // U-values for building elements
	Insulation            *InsulationThicknesses        // Insulation thickness parameters
	SolarGains            *SolarGains                   // Solar gain parameters
	ThermalBridges        *ThermalBridgeParameters      // Thermal bridging parameters
	HeatLosses            *TransmissionHeatLoss         // Transmission heat loss coefficients
	ThermalResistances    *ThermalResistances           // Thermal resistances for predefined measures
	InsulationMeasures    *InsulationPredefinedMeasures // Insulation thicknesses for predefined measures
	ActualInsulation      *ActualInsulationThicknesses  // Actual insulation thicknesses for input measures
	HeatTransfer          *HeatTransferCoefficients     // Heat transfer coefficients
	PredefinedCodes       *PredefinedCodes              // Pre-defined codes and configuration parameters
	MeasureTypes          *MeasureTypeCodes             // Measure type codes for building elements
	SolarTransmittance    *SolarEnergyTransmittance     // Solar energy transmittance parameters
	MeasureFractions      *MeasureAreaFractions         // Area fractions for building element measures
	AdditionalResistances *AdditionalThermalResistance  // Additional thermal resistance due to unheated space
}

// BuildingThematic holds thematic building parameters
type BuildingThematic struct {
	Code_BuildingVariant    string  `json:"Code_BuildingVariant"`    // e.g. DE.SFH
	Code_TypeVariant        string  `json:"Code_TypeVariant"`        // e.g. DE.SFH.1960
	Code_AttachedNeighbours string  `json:"Code_AttachedNeighbours"` // Attachment status of neighbours (e.g. "B_Alone", "B_N1", "B_N2")
	Code_ComplexFootprint   string  `json:"Code_ComplexFootprint"`   // Complexity of building footprint (e.g. "Simple", "Standard", "Complex")
	Code_ComplexRoof        string  `json:"Code_ComplexRoof"`        // Complexity of roof shape (e.g. "Simple", "Standard", "Complex")
	N_Storey                int     `json:"n_Storey"`                // Number of storeys (used for calculating conditioned volume if not provided in dataset)
	H_room                  float64 `json:"h_room"`                  // Room height in meters (used for calculating conditioned volume if not provided in dataset)
	Code_AtticCond          string  `json:"Code_AtticCond"`          // Attic condition code
	Code_CellarCond         string  `json:"Code_CellarCond"`         // Cellar condition code
}

// Envelope holds building envelope parameters like areas and volumes of different building elements.
// It also includes the conditioned volume and floor area, which are essential for calculating heat losses and energy demand.
type Envelope struct {
	V_C           float64 `json:"V_C"`           // m³ conditioned volume
	A_C_ExtDim    float64 `json:"A_C_ExtDim"`    // m² conditioned floor area (external dims)
	A_C_Ref_Estim float64 `json:"A_C_Ref_Estim"` // m² estimated reference floor area (used for calculating area fractions for measures)
	A_C_Ref_Input float64 `json:"A_C_Ref_Input"` // m² reference floor area from input dataset (used for calculating area fractions for measures)
	A_C_IntDim    float64 `json:"A_C_IntDim"`    // m² conditioned floor area (internal dims)
	A_C_Use       float64 `json:"A_C_Use"`       // m² useful conditioned floor area
	A_C_Living    float64 `json:"A_C_Living"`    // m² conditioned living area
	A_C_RefInput  float64 `json:"A_C_Ref_Input"` // m² reference floor area

	// Roofs
	A_Roof_1 float64 `json:"A_Roof_1"` // m² roof area 1
	A_Roof_2 float64 `json:"A_Roof_2"` // m² roof area 2

	A_Wall_1 float64 `json:"A_Wall_1"` // m² wall area 1
	A_Wall_2 float64 `json:"A_Wall_2"` // m² wall area 2
	A_Wall_3 float64 `json:"A_Wall_3"` // m² wall area 3

	A_Floor_1 float64 `json:"A_Floor_1"` // m² floor area 1
	A_Floor_2 float64 `json:"A_Floor_2"` // m² floor area 2

	A_Window_1          float64 `json:"A_Window_1"`          // m² window area 1
	A_Window_2          float64 `json:"A_Window_2"`          // m² window area 2
	A_Window_Horizontal float64 `json:"A_Window_Horizontal"` // m² horizontal window area (for solar gains calculation)
	A_Window_East       float64 `json:"A_Window_East"`       // m² east-facing window area (for solar gains calculation)
	A_Window_South      float64 `json:"A_Window_South"`      // m² south-facing window area (for solar gains calculation)
	A_Window_West       float64 `json:"A_Window_West"`       // m² west-facing window area (for solar gains calculation)
	A_Window_North      float64 `json:"A_Window_North"`      // m² north-facing window area (for solar gains calculation)

	// Doors
	A_Door_1 float64 `json:"A_Door_1"` // m² door area 1
}

// AirInfiltration holds parameters related to air infiltration and usage
type AirInfiltration struct {
	N_air_infiltration float64 `json:"n_air_infiltration"` // Air change rate due to infiltration (1/h)
	N_air_use          float64 `json:"n_air_use"`          // Air change rate due to building usage (1/h)
}

// ClimateConditions holds climate-related parameters
type ClimateConditions struct {
	HeatingDays int     `json:"HeatingDays"` // Number of heating days per year
	Theta_e     float64 `json:"Theta_e"`     // °C External design temperature (used for calculating temperature difference and heat losses)
	Theta_i     float64 `json:"theta_i"`     // °C Internal design temperature (used for calculating temperature difference and heat losses) - lowercase to match Excel
}

// U-values holds thermal transmittance parameters for building elements, including original,
// measure, and actual U-values for different building components. These are essential for calculating heat losses and energy demand,
// as well as for evaluating the impact of insulation measures.
type Uvalues struct {
	U_Roof_1   float64 `json:"U_Roof_1"`   // Original U-value for roof 1 (W/m²K)
	U_Roof_2   float64 `json:"U_Roof_2"`   // Original U-value for roof 2 (W/m²K)
	U_Wall_1   float64 `json:"U_Wall_1"`   // Original U-value for wall 1 (W/m²K)
	U_Wall_2   float64 `json:"U_Wall_2"`   // Original U-value for wall 2 (W/m²K)
	U_Wall_3   float64 `json:"U_Wall_3"`   // Original U-value for wall 3 (W/m²K)
	U_Floor_1  float64 `json:"U_Floor_1"`  // Original U-value for floor 1 (W/m²K)
	U_Floor_2  float64 `json:"U_Floor_2"`  // Original U-value for floor 2 (W/m²K)
	U_Window_1 float64 `json:"U_Window_1"` // Original U-value for window 1 (W/m²K)
	U_Window_2 float64 `json:"U_Window_2"` // Original U-value for window 2 (W/m²K)
	U_Door_1   float64 `json:"U_Door_1"`   // Original U-value for door 1 (W/m²K)

	// Measure U-values (calculated in Level 4)
	U_Measure_Roof_1   float64 `json:"U_Measure_Roof_1"`   // U-value for roof 1 after applying insulation measure (W/m²K)
	U_Measure_Roof_2   float64 `json:"U_Measure_Roof_2"`   // U-value for roof 2 after applying insulation measure (W/m²K)
	U_Measure_Wall_1   float64 `json:"U_Measure_Wall_1"`   // U-value for wall 1 after applying insulation measure (W/m²K)
	U_Measure_Wall_2   float64 `json:"U_Measure_Wall_2"`   // U-value for wall 2 after applying insulation measure (W/m²K)
	U_Measure_Wall_3   float64 `json:"U_Measure_Wall_3"`   // U-value for wall 3 after applying insulation measure (W/m²K)
	U_Measure_Floor_1  float64 `json:"U_Measure_Floor_1"`  // U-value for floor 1 after applying insulation measure (W/m²K)
	U_Measure_Floor_2  float64 `json:"U_Measure_Floor_2"`  // U-value for floor 2 after applying insulation measure (W/m²K)
	U_Measure_Window_1 float64 `json:"U_Measure_Window_1"` // U-value for window 1 after applying insulation measure (W/m²K)
	U_Measure_Window_2 float64 `json:"U_Measure_Window_2"` // U-value for window 2 after applying insulation measure (W/m²K)
	U_Measure_Door_1   float64 `json:"U_Measure_Door_1"`   // U-value for door 1 after applying insulation measure (W/m²K)

	// Actual U-values (calculated in Level 5)
	U_Actual_Roof_1   float64 `json:"U_Actual_Roof_1"`   // Actual U-value for roof 1 after applying insulation measure and considering thermal bridging (W/m²K)
	U_Actual_Roof_2   float64 `json:"U_Actual_Roof_2"`   // Actual U-value for roof 2 after applying insulation measure and considering thermal bridging (W/m²K)
	U_Actual_Wall_1   float64 `json:"U_Actual_Wall_1"`   // Actual U-value for wall 1 after applying insulation measure and considering thermal bridging (W/m²K)
	U_Actual_Wall_2   float64 `json:"U_Actual_Wall_2"`   // Actual U-value for wall 2 after applying insulation measure and considering thermal bridging (W/m²K)
	U_Actual_Wall_3   float64 `json:"U_Actual_Wall_3"`   // Actual U-value for wall 3 after applying insulation measure and considering thermal bridging (W/m²K)
	U_Actual_Floor_1  float64 `json:"U_Actual_Floor_1"`  // Actual U-value for floor 1 after applying insulation measure and considering thermal bridging (W/m²K)
	U_Actual_Floor_2  float64 `json:"U_Actual_Floor_2"`  // Actual U-value for floor 2 after applying insulation measure and considering thermal bridging (W/m²K)
	U_Actual_Window_1 float64 `json:"U_Actual_Window_1"` // Actual U-value for window 1 after applying insulation measure and considering thermal bridging (W/m²K)
	U_Actual_Window_2 float64 `json:"U_Actual_Window_2"` // Actual U-value for window 2 after applying insulation measure and considering thermal bridging (W/m²K)
	U_Actual_Door_1   float64 `json:"U_Actual_Door_1"`   // Actual U-value for door 1 after applying insulation measure and considering thermal bridging (W/m²K)
}

// InsulationThicknesses holds insulation thickness parameters
type InsulationThicknesses struct {
	D_Insulation_Roof_1  float64 `json:"d_Insulation_Roof_1"`  // Insulation thickness for roof 1 (m)
	D_Insulation_Roof_2  float64 `json:"d_Insulation_Roof_2"`  // Insulation thickness for roof 2 (m)
	D_Insulation_Wall_1  float64 `json:"d_Insulation_Wall_1"`  // Insulation thickness for wall 1 (m)
	D_Insulation_Wall_2  float64 `json:"d_Insulation_Wall_2"`  // Insulation thickness for wall 2 (m)
	D_Insulation_Wall_3  float64 `json:"d_Insulation_Wall_3"`  // Insulation thickness for wall 3 (m)
	D_Insulation_Floor_1 float64 `json:"d_Insulation_Floor_1"` // Insulation thickness for floor 1 (m)
	D_Insulation_Floor_2 float64 `json:"d_Insulation_Floor_2"` // Insulation thickness for floor 2 (m)
}

// SolarGains holds solar gain parameters
type SolarGains struct {
	I_Sol_Horizontal float64 `json:"I_Sol_Hor"`   // Solar radiation on horizontal surface (W/m²)
	I_Sol_East       float64 `json:"I_Sol_East"`  // Solar radiation on east-facing vertical surface (W/m²)
	I_Sol_South      float64 `json:"I_Sol_South"` // Solar radiation on south-facing vertical surface (W/m²)
	I_Sol_West       float64 `json:"I_Sol_West"`  // Solar radiation on west-facing vertical surface (W/m²)
	I_Sol_North      float64 `json:"I_Sol_North"` // Solar radiation on north-facing vertical surface (W/m²)
}

type ThermalBridgeParameters struct {
	Code_ThermalBridging_Original       string  `json:"Code_ThermalBridging_Original"`       // Thermal bridging code for original building elements (e.g. "Original", "Refurbished", "None")
	Code_ThermalBridging_Refurbished    string  `json:"Code_ThermalBridging_Refurbished"`    // Thermal bridging code for refurbished building elements (e.g. "Original", "Refurbished", "None")
	Delta_U_ThermalBridging_Original    float64 `json:"delta_U_ThermalBridging_Original"`    // Additional U-value due to thermal bridging for original building elements (W/m²K)
	Delta_U_ThermalBridging_Refurbished float64 `json:"delta_U_ThermalBridging_Refurbished"` // Additional U-value due to thermal bridging for refurbished building elements (W/m²K)
}

type TransmissionHeatLoss struct {
	B_Transmission_Roof_1  float64 `json:"b_Transmission_Roof_1"`  // Transmission heat loss coefficient for roof 1
	B_Transmission_Roof_2  float64 `json:"b_Transmission_Roof_2"`  // Transmission heat loss coefficient for roof 2
	B_Transmission_Wall_1  float64 `json:"b_Transmission_Wall_1"`  // Transmission heat loss coefficient for wall 1
	B_Transmission_Wall_2  float64 `json:"b_Transmission_Wall_2"`  // Transmission heat loss coefficient for wall 2
	B_Transmission_Wall_3  float64 `json:"b_Transmission_Wall_3"`  // Transmission heat loss coefficient for wall 3
	B_Transmission_Floor_1 float64 `json:"b_Transmission_Floor_1"` // Transmission heat loss coefficient for floor 1
	B_Transmission_Floor_2 float64 `json:"b_Transmission_Floor_2"` // Transmission heat loss coefficient for floor 2
}

// Thermal resistance of building elements
type ThermalResistances struct {
	R_PredefinedMeasure_Roof_1   float64 `json:"R_PredefinedMeasure_Roof_1"`   // Thermal resistance of predefined measure for roof 1
	R_PredefinedMeasure_Roof_2   float64 `json:"R_PredefinedMeasure_Roof_2"`   // Thermal resistance of predefined measure for roof 2
	R_PredefinedMeasure_Wall_1   float64 `json:"R_PredefinedMeasure_Wall_1"`   // Thermal resistance of predefined measure for wall 1
	R_PredefinedMeasure_Wall_2   float64 `json:"R_PredefinedMeasure_Wall_2"`   // Thermal resistance of predefined measure for wall 2
	R_PredefinedMeasure_Wall_3   float64 `json:"R_PredefinedMeasure_Wall_3"`   // Thermal resistance of predefined measure for wall 3
	R_PredefinedMeasure_Floor_1  float64 `json:"R_PredefinedMeasure_Floor_1"`  // Thermal resistance of predefined measure for floor 1
	R_PredefinedMeasure_Floor_2  float64 `json:"R_PredefinedMeasure_Floor_2"`  // Thermal resistance of predefined measure for floor 2
	R_PredefinedMeasure_Window_1 float64 `json:"R_PredefinedMeasure_Window_1"` // Thermal resistance of predefined measure for window 1
	R_PredefinedMeasure_Window_2 float64 `json:"R_PredefinedMeasure_Window_2"` // Thermal resistance of predefined measure for window 2
	R_PredefinedMeasure_Door_1   float64 `json:"R_PredefinedMeasure_Door_1"`   // Thermal resistance of predefined measure for door 1
}

type InsulationPredefinedMeasures struct {
	D_Insulation_PredefinedMeasure_Roof_1  float64 `json:"d_Insulation_PredefinedMeasure_Roof_1"`  // Insulation thickness of predefined measure for roof 1 (m)
	D_Insulation_PredefinedMeasure_Roof_2  float64 `json:"d_Insulation_PredefinedMeasure_Roof_2"`  // Insulation thickness of predefined measure for roof 2 (m)
	D_Insulation_PredefinedMeasure_Wall_1  float64 `json:"d_Insulation_PredefinedMeasure_Wall_1"`  // Insulation thickness of predefined measure for wall 1 (m)
	D_Insulation_PredefinedMeasure_Wall_2  float64 `json:"d_Insulation_PredefinedMeasure_Wall_2"`  // Insulation thickness of predefined measure for wall 2 (m)
	D_Insulation_PredefinedMeasure_Wall_3  float64 `json:"d_Insulation_PredefinedMeasure_Wall_3"`  // Insulation thickness of predefined measure for wall 3 (m)
	D_Insulation_PredefinedMeasure_Floor_1 float64 `json:"d_Insulation_PredefinedMeasure_Floor_1"` // Insulation thickness of predefined measure for floor 1 (m)
	D_Insulation_PredefinedMeasure_Floor_2 float64 `json:"d_Insulation_PredefinedMeasure_Floor_2"` // Insulation thickness of predefined measure for floor 2 (m)
}

type ActualInsulationThicknesses struct {
	D_Insulation_Input_Measure_Roof_1  float64 `json:"d_Insulation_Input_Measure_Roof_1"`  // Actual insulation thickness for roof 1 after applying insulation measure (m)
	D_Insulation_Input_Measure_Roof_2  float64 `json:"d_Insulation_Input_Measure_Roof_2"`  // Actual insulation thickness for roof 2 after applying insulation measure (m)
	D_Insulation_Input_Measure_Wall_1  float64 `json:"d_Insulation_Input_Measure_Wall_1"`  // Actual insulation thickness for wall 1 after applying insulation measure (m)
	D_Insulation_Input_Measure_Wall_2  float64 `json:"d_Insulation_Input_Measure_Wall_2"`  // Actual insulation thickness for wall 2 after applying insulation measure (m)
	D_Insulation_Input_Measure_Wall_3  float64 `json:"d_Insulation_Input_Measure_Wall_3"`  // Actual insulation thickness for wall 3 after applying insulation measure (m)
	D_Insulation_Input_Measure_Floor_1 float64 `json:"d_Insulation_Input_Measure_Floor_1"` // Actual insulation thickness for floor 1 after applying insulation measure (m)
	D_Insulation_Input_Measure_Floor_2 float64 `json:"d_Insulation_Input_Measure_Floor_2"` // Actual insulation thickness for floor 2 after applying insulation measure (m)
}

type HeatTransferCoefficients struct {
	F_red_htr1 float64 `json:"F_red_htr1"`
	F_red_htr4 float64 `json:"F_red_htr4"`
	Phi_int    float64 `json:"phi_int"`   // Internal heat Source per m² - FIXED: lowercase to match Excel
	F_sh_hor   float64 `json:"F_sh_hor"`  // Horizontal shading factor
	F_sh_vert  float64 `json:"F_sh_vert"` // Vertical shading factor
	F_f        float64 `json:"F_f"`       // Window frame area fraction
	F_w        float64 `json:"F_w"`       // Window wall area fraction
	C_m        float64 `json:"c_m"`       // Internal heat capacity per m² of useful floor area
}

// Pre-defined Codes and Configuration Parameters
type PredefinedCodes struct {
	Code_TypeIntake_EnvelopeArea       string  `json:"Code_TypeIntake_EnvelopeArea"`       // Envelope Area Calculation Method
	F_Corr_CeilingHeight               float64 `json:"f_Corr_CeilingHeight"`               // Correction Factor for Ceiling Height
	F_PlausiCrit_EnvSum_LowerLimit     float64 `json:"f_PlausiCrit_EnvSum_LowerLimit"`     // Lower Limit for Envelope Sum Ratio
	F_PlausiCrit_EnvSum_UpperLimit     float64 `json:"f_PlausiCrit_EnvSum_UpperLimit"`     // Upper Limit for Envelope Sum Ratio
	F_PlausiCrit_FloorArea_LowerLimit  float64 `json:"f_PlausiCrit_FloorArea_LowerLimit"`  // Lower Limit for Floor Area
	F_PlausiCrit_FloorArea_UpperLimit  float64 `json:"f_PlausiCrit_FloorArea_UpperLimit"`  // Upper Limit for Floor Area
	F_PlausiCrit_WindowArea_LowerLimit float64 `json:"f_PlausiCrit_WindowArea_LowerLimit"` // Lower Limit for Window Area
	F_PlausiCrit_WindowArea_UpperLimit float64 `json:"f_PlausiCrit_WindowArea_UpperLimit"` // Upper Limit for Window Area
}

// Measure Type Codes for Building Elements
type MeasureTypeCodes struct {
	Code_MeasureType_Roof_1   string `json:"Code_MeasureType_Roof_1"`   // Measure Type Code Roof 1
	Code_MeasureType_Roof_2   string `json:"Code_MeasureType_Roof_2"`   // Measure Type Code Roof 2
	Code_MeasureType_Wall_1   string `json:"Code_MeasureType_Wall_1"`   // Measure Type Code Wall 1
	Code_MeasureType_Wall_2   string `json:"Code_MeasureType_Wall_2"`   // Measure Type Code Wall 2
	Code_MeasureType_Wall_3   string `json:"Code_MeasureType_Wall_3"`   // Measure Type Code Wall 3
	Code_MeasureType_Floor_1  string `json:"Code_MeasureType_Floor_1"`  // Measure Type Code Floor 1
	Code_MeasureType_Floor_2  string `json:"Code_MeasureType_Floor_2"`  // Measure Type Code Floor 2
	Code_Measure_Window_1     string `json:"Code_Measure_Window_1"`     // Measure Dataset ID for Window 1
	Code_Measure_Window_2     string `json:"Code_Measure_Window_2"`     // Measure Dataset ID for Window 2
	Code_MeasureType_Window_1 string `json:"Code_MeasureType_Window_1"` // Measure Type Code Window 1
	Code_MeasureType_Window_2 string `json:"Code_MeasureType_Window_2"` // Measure Type Code Window 2
	Code_MeasureType_Door_1   string `json:"Code_MeasureType_Door_1"`   // Measure Type Code Door 1
}

// Solar Energy Transmittance Parameters
type SolarEnergyTransmittance struct {
	G_gl_n_Window_1                   float64 `json:"g_gl_n_Window_1"`                   // Solar Energy Transmittance for Window 1 (Perpendicular Radiation)
	G_gl_n_Window_2                   float64 `json:"g_gl_n_Window_2"`                   // Solar Energy Transmittance for Window 2 (Perpendicular Radiation)
	G_gl_n_PredefinedMeasure_Window_1 float64 `json:"g_gl_n_PredefinedMeasure_Window_1"` // Solar Energy Transmittance for Refurbished Window 1 (Perpendicular Radiation)
	G_gl_n_PredefinedMeasure_Window_2 float64 `json:"g_gl_n_PredefinedMeasure_Window_2"` // Solar Energy Transmittance for Refurbished Window 2 (Perpendicular Radiation)
}

// Area Fractions for Building Element Measures
// NOTE: These fields should be loaded from database if available, otherwise need to be calculated or set as constants
type MeasureAreaFractions struct {
	F_Measure_Roof_1   float64 `json:"f_Measure_Roof_1"`   // Area Fraction of Roof Measure 1
	F_Measure_Roof_2   float64 `json:"f_Measure_Roof_2"`   // Area Fraction of Roof Measure 2
	F_Measure_Wall_1   float64 `json:"f_Measure_Wall_1"`   // Area Fraction of Wall Measure 1
	F_Measure_Wall_2   float64 `json:"f_Measure_Wall_2"`   // Area Fraction of Wall Measure 2
	F_Measure_Wall_3   float64 `json:"f_Measure_Wall_3"`   // Area Fraction of Wall Measure 3
	F_Measure_Floor_1  float64 `json:"f_Measure_Floor_1"`  // Area Fraction of Floor Measure 1
	F_Measure_Floor_2  float64 `json:"f_Measure_Floor_2"`  // Area Fraction of Floor Measure 2
	F_Measure_Window_1 float64 `json:"f_Measure_Window_1"` // Area Fraction of Window Measure 1
	F_Measure_Window_2 float64 `json:"f_Measure_Window_2"` // Area Fraction of Window Measure 2
	F_Measure_Door_1   float64 `json:"f_Measure_Door_1"`   // Area Fraction of Door Measure 1
}

// Additional Thermal Resistance due to Unheated Space
type AdditionalThermalResistance struct {
	R_Add_UnheatedSpace_Roof_1  float64 `json:"R_Add_UnheatedSpace_Roof_1"`  // Additional Thermal Resistance Roof 1
	R_Add_UnheatedSpace_Roof_2  float64 `json:"R_Add_UnheatedSpace_Roof_2"`  // Additional Thermal Resistance Roof 2
	R_Add_UnheatedSpace_Wall_1  float64 `json:"R_Add_UnheatedSpace_Wall_1"`  // Additional Thermal Resistance Wall 1
	R_Add_UnheatedSpace_Wall_2  float64 `json:"R_Add_UnheatedSpace_Wall_2"`  // Additional Thermal Resistance Wall 2
	R_Add_UnheatedSpace_Wall_3  float64 `json:"R_Add_UnheatedSpace_Wall_3"`  // Additional Thermal Resistance Wall 3
	R_Add_UnheatedSpace_Floor_1 float64 `json:"R_Add_UnheatedSpace_Floor_1"` // Additional Thermal Resistance Floor 1
	R_Add_UnheatedSpace_Floor_2 float64 `json:"R_Add_UnheatedSpace_Floor_2"` // Additional Thermal Resistance Floor 2
}
