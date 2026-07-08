package models

// FieldMetadata describes one TABULA input field for human consumption:
// where to find it in a GET /api/v1/data/:code response, its unit, and two
// descriptions at different levels of expertise. This is static, generated
// once from this file — the same for every TABULA country, since the DB
// schema itself is identical across countries (see build_db).
//
// Client applications (e.g. Building Configurator) use this to power field
// labels/tooltips today, and are expected to use it to drive an interactive
// building-description questionnaire in the future, where SimpleDescription
// becomes the question text and the field's own value (from the matched
// TABULA variant) becomes the suggested default.
type FieldMetadata struct {
	Key               string `json:"key"`                // matches the json tag in tabula.go and the leaf key in tabula_data
	Group             string `json:"group"`              // matches the struct name in tabula.go, e.g. "ClimateConditions"
	Path              string `json:"path"`               // full dotted path into tabula_data, e.g. "AdvancedParameters.ClimateConditions.HeatingDays"
	Unit              string `json:"unit,omitempty"`     // e.g. "days/year", "°C", "m²", "W/m²K"
	Label             string `json:"label"`              // short human label, e.g. "Heating days"
	SimpleDescription string `json:"simple_description"` // plain-language description, for a future non-expert questionnaire
	ExpertDescription string `json:"expert_description"` // technical description, matches the tabula.go field comment
}

// AllFieldMetadata lists every TABULA input field consumed by ignis's clients
// (envelope areas, climate conditions, U-values, air infiltration, heat
// transfer coefficients, solar gains, and thermal bridging). Internal-only
// TABULA fields (measure codes, plausibility thresholds, etc.) are omitted.
var AllFieldMetadata = []FieldMetadata{
	// Envelope areas
	{
		Key:               "A_C_Ref_Input",
		Group:             "Envelope",
		Path:              "BasicParameters.Envelope.A_C_Ref_Input",
		Unit:              "m²",
		Label:             "Reference floor area",
		SimpleDescription: "The total heated floor area of the building.",
		ExpertDescription: "Reference floor area from input dataset, used for calculating area fractions for measures."
	},
	{
		Key:               "A_Roof_1",
		Group:             "Envelope",
		Path:              "BasicParameters.Envelope.A_Roof_1",
		Unit:              "m²",
		Label:             "Roof area",
		SimpleDescription: "How much roof surface the building has.",
		ExpertDescription: "Roof area 1."
	},
	{
		Key:               "A_Roof_2",
		Group:             "Envelope",
		Path:              "BasicParameters.Envelope.A_Roof_2",
		Unit:              "m²",
		Label:             "Second roof area",
		SimpleDescription: "The area of a second, differently-shaped roof section, if the building has one.",
		ExpertDescription: "Roof area 2."
	},
	{
		Key:               "A_Wall_1",
		Group:             "Envelope",
		Path:              "BasicParameters.Envelope.A_Wall_1",
		Unit:              "m²",
		Label:             "Wall area",
		SimpleDescription: "How much exterior wall surface the building has.",
		ExpertDescription: "Wall area 1."
	},
	{
		Key:               "A_Wall_2",
		Group:             "Envelope",
		Path:              "BasicParameters.Envelope.A_Wall_2",
		Unit:              "m²",
		Label:             "Second wall area",
		SimpleDescription: "The area of a second type of exterior wall construction, if present.",
		ExpertDescription: "Wall area 2."
	},
	{
		Key:               "A_Wall_3",
		Group:             "Envelope",
		Path:              "BasicParameters.Envelope.A_Wall_3",
		Unit:              "m²",
		Label:             "Third wall area",
		SimpleDescription: "The area of a third type of exterior wall construction, if present.",
		ExpertDescription: "Wall area 3."
	},
	{
		Key:               "A_Floor_1",
		Group:             "Envelope",
		Path:              "BasicParameters.Envelope.A_Floor_1",
		Unit:              "m²",
		Label:             "Floor area (against ground/basement)",
		SimpleDescription: "How much of the building's floor sits directly on the ground or above an unheated basement.",
		ExpertDescription: "Floor area 1."
	},
	{
		Key:               "A_Floor_2",
		Group:             "Envelope",
		Path:              "BasicParameters.Envelope.A_Floor_2",
		Unit:              "m²",
		Label:             "Second floor type area",
		SimpleDescription: "The area of a second type of ground-facing floor construction, if present.",
		ExpertDescription: "Floor area 2."
	},
	{
		Key:               "A_Window_1",
		Group:             "Envelope",
		Path:              "BasicParameters.Envelope.A_Window_1",
		Unit:              "m²",
		Label:             "Window area",
		SimpleDescription: "The total glazed window area of the building.",
		ExpertDescription: "Window area 1."
	},
	{
		Key:               "A_Window_2",
		Group:             "Envelope",
		Path:              "BasicParameters.Envelope.A_Window_2",
		Unit:              "m²",
		Label:             "Second window type area",
		SimpleDescription: "The area of a second type of window (e.g. a different glazing), if present.",
		ExpertDescription: "Window area 2."
	},
	{
		Key:               "A_Window_South",
		Group:             "Envelope",
		Path:              "BasicParameters.Envelope.A_Window_South",
		Unit:              "m²",
		Label:             "South-facing window area",
		SimpleDescription: "How much window area faces south.",
		ExpertDescription: "South-facing window area, used for solar gains calculation."},
	{
		Key:               "A_Window_East",
		Group:             "Envelope",
		Path:              "BasicParameters.Envelope.A_Window_East",
		Unit:              "m²",
		Label:             "East-facing window area",
		SimpleDescription: "How much window area faces east.",
		ExpertDescription: "East-facing window area, used for solar gains calculation."},
	{
		Key:               "A_Window_West",
		Group:             "Envelope",
		Path:              "BasicParameters.Envelope.A_Window_West",
		Unit:              "m²",
		Label:             "West-facing window area",
		SimpleDescription: "How much window area faces west.",
		ExpertDescription: "West-facing window area, used for solar gains calculation."},
	{
		Key:               "A_Window_North",
		Group:             "Envelope",
		Path:              "BasicParameters.Envelope.A_Window_North",
		Unit:              "m²",
		Label:             "North-facing window area",
		SimpleDescription: "How much window area faces north.",
		ExpertDescription: "North-facing window area, used for solar gains calculation."},
	{
		Key:               "A_Door_1",
		Group:             "Envelope",
		Path:              "BasicParameters.Envelope.A_Door_1",
		Unit:              "m²",
		Label:             "Door area",
		SimpleDescription: "The total area of exterior doors.",
		ExpertDescription: "Door area 1."},

	// Climate conditions
	{
		Key:               "HeatingDays",
		Group:             "ClimateConditions",
		Path:              "AdvancedParameters.ClimateConditions.HeatingDays",
		Unit:              "days/year",
		Label:             "Heating days",
		SimpleDescription: "How many days a year the building typically needs heating, based on local climate.",
		ExpertDescription: "Number of heating days per year."},
	{
		Key:               "Theta_e",
		Group:             "ClimateConditions",
		Path:              "AdvancedParameters.ClimateConditions.Theta_e",
		Unit:              "°C",
		Label:             "Outside design temperature",
		SimpleDescription: "How cold it typically gets outside in winter in your area.",
		ExpertDescription: "External design temperature, used for calculating temperature difference and heat losses."},
	{
		Key:               "theta_i",
		Group:             "ClimateConditions",
		Path:              "AdvancedParameters.ClimateConditions.theta_i",
		Unit:              "°C",
		Label:             "Inside design temperature",
		SimpleDescription: "The indoor temperature you keep the building at.",
		ExpertDescription: "Internal design temperature, used for calculating temperature difference and heat losses."},

	// U-values
	{
		Key:               "U_Roof_1",
		Group:             "Uvalues",
		Path:              "AdvancedParameters.Uvalues.U_Roof_1",
		Unit:              "W/m²K",
		Label:             "Roof insulation value",
		SimpleDescription: "How well the roof keeps heat in — lower means better insulated.",
		ExpertDescription: "Original U-value for roof 1."},
	{
		Key:               "U_Wall_1",
		Group:             "Uvalues",
		Path:              "AdvancedParameters.Uvalues.U_Wall_1",
		Unit:              "W/m²K",
		Label:             "Wall insulation value",
		SimpleDescription: "How well the exterior walls keep heat in — lower means better insulated.",
		ExpertDescription: "Original U-value for wall 1."},
	{
		Key:               "U_Floor_1",
		Group:             "Uvalues",
		Path:              "AdvancedParameters.Uvalues.U_Floor_1",
		Unit:              "W/m²K",
		Label:             "Floor insulation value",
		SimpleDescription: "How well the ground-facing floor keeps heat in — lower means better insulated.",
		ExpertDescription: "Original U-value for floor 1."},
	{
		Key:               "U_Window_1",
		Group:             "Uvalues",
		Path:              "AdvancedParameters.Uvalues.U_Window_1",
		Unit:              "W/m²K",
		Label:             "Window insulation value",
		SimpleDescription: "How well the windows keep heat in — lower means better insulated (e.g. double/triple glazing).",
		ExpertDescription: "Original U-value for window 1."},
	{
		Key:               "U_Door_1",
		Group:             "Uvalues",
		Path:              "AdvancedParameters.Uvalues.U_Door_1",
		Unit:              "W/m²K",
		Label:             "Door insulation value",
		SimpleDescription: "How well the exterior doors keep heat in — lower means better insulated.",
		ExpertDescription: "Original U-value for door 1."},

	// Air infiltration
	{
		Key:               "n_air_infiltration",
		Group:             "AirInfiltration",
		Path:              "AdvancedParameters.AirInfiltration.n_air_infiltration",
		Unit:              "1/h",
		Label:             "Draughtiness",
		SimpleDescription: "How much outside air leaks in through gaps and cracks, independent of ventilation.",
		ExpertDescription: "Air change rate due to infiltration."},
	{
		Key:               "n_air_use",
		Group:             "AirInfiltration",
		Path:              "AdvancedParameters.AirInfiltration.n_air_use",
		Unit:              "1/h",
		Label:             "Ventilation rate",
		SimpleDescription: "How often the indoor air is exchanged through normal use (opening windows, mechanical ventilation).",
		ExpertDescription: "Air change rate due to building usage."},

	// Heat transfer coefficients
	{
		Key:               "F_sh_hor",
		Group:             "HeatTransfer",
		Path:              "AdvancedParameters.HeatTransfer.F_sh_hor",
		Unit:              "",
		Label:             "Horizontal shading",
		SimpleDescription: "How much overhangs or nearby buildings block sunlight from above.",
		ExpertDescription: "Horizontal shading factor."},
	{
		Key:               "F_sh_vert",
		Group:             "HeatTransfer",
		Path:              "AdvancedParameters.HeatTransfer.F_sh_vert",
		Unit:              "",
		Label:             "Vertical shading",
		SimpleDescription: "How much nearby buildings or trees block sunlight from the side.",
		ExpertDescription: "Vertical shading factor."},
	{
		Key:               "F_f",
		Group:             "HeatTransfer",
		Path:              "AdvancedParameters.HeatTransfer.F_f",
		Unit:              "",
		Label:             "Window frame fraction",
		SimpleDescription: "How much of the window area is frame rather than glass.",
		ExpertDescription: "Window frame area fraction."},
	{
		Key:               "F_w",
		Group:             "HeatTransfer",
		Path:              "AdvancedParameters.HeatTransfer.F_w",
		Unit:              "",
		Label:             "Window wall fraction",
		SimpleDescription: "How much of the exterior wall is taken up by windows.",
		ExpertDescription: "Window wall area fraction."},
	{
		Key:               "phi_int",
		Group:             "HeatTransfer",
		Path:              "AdvancedParameters.HeatTransfer.phi_int",
		Unit:              "W/m²",
		Label:             "Internal heat gains",
		SimpleDescription: "Heat given off inside the building by people, appliances and lighting.",
		ExpertDescription: "Internal heat source per m² of useful floor area."},
	{
		Key:               "c_m",
		Group:             "HeatTransfer",
		Path:              "AdvancedParameters.HeatTransfer.c_m",
		Unit:              "J/(m²K)",
		Label:             "Thermal mass",
		SimpleDescription: "How much heat the building's structure (walls, floors) can store — heavier, denser buildings store more.",
		ExpertDescription: "Internal heat capacity per m² of useful floor area."},

	// Solar gains
	{
		Key:               "I_Sol_South",
		Group:             "SolarGains",
		Path:              "AdvancedParameters.SolarGains.I_Sol_South",
		Unit:              "W/m²",
		Label:             "Solar radiation (south)",
		SimpleDescription: "How much sunlight energy typically hits a south-facing surface in your area.",
		ExpertDescription: "Solar radiation on south-facing vertical surface."},
	{
		Key:               "I_Sol_East",
		Group:             "SolarGains",
		Path:              "AdvancedParameters.SolarGains.I_Sol_East",
		Unit:              "W/m²",
		Label:             "Solar radiation (east)",
		SimpleDescription: "How much sunlight energy typically hits an east-facing surface in your area.",
		ExpertDescription: "Solar radiation on east-facing vertical surface."},
	{
		Key:               "I_Sol_West",
		Group:             "SolarGains",
		Path:              "AdvancedParameters.SolarGains.I_Sol_West",
		Unit:              "W/m²",
		Label:             "Solar radiation (west)",
		SimpleDescription: "How much sunlight energy typically hits a west-facing surface in your area.",
		ExpertDescription: "Solar radiation on west-facing vertical surface."},
	{
		Key:               "I_Sol_North",
		Group:             "SolarGains",
		Path:              "AdvancedParameters.SolarGains.I_Sol_North",
		Unit:              "W/m²",
		Label:             "Solar radiation (north)",
		SimpleDescription: "How much sunlight energy typically hits a north-facing surface in your area.",
		ExpertDescription: "Solar radiation on north-facing vertical surface."},
	{
		Key:               "I_Sol_Hor",
		Group:             "SolarGains",
		Path:              "AdvancedParameters.SolarGains.I_Sol_Hor",
		Unit:              "W/m²",
		Label:             "Solar radiation (horizontal)",
		SimpleDescription: "How much sunlight energy typically hits a flat, upward-facing surface (like a flat roof) in your area.",
		ExpertDescription: "Solar radiation on horizontal surface."},

	// Thermal bridging
	{
		Key:               "delta_U_ThermalBridging_Original",
		Group:             "ThermalBridges",
		Path:              "AdvancedParameters.ThermalBridges.delta_U_ThermalBridging_Original",
		Unit:              "W/m²K",
		Label:             "Thermal bridging (existing state)",
		SimpleDescription: "Extra heat loss through junctions and corners (e.g. where a wall meets a balcony), for the building as it stands today.",
		ExpertDescription: "Additional U-value due to thermal bridging for original building elements."},
	{
		Key:               "delta_U_ThermalBridging_Refurbished",
		Group:             "ThermalBridges",
		Path:              "AdvancedParameters.ThermalBridges.delta_U_ThermalBridging_Refurbished",
		Unit:              "W/m²K",
		Label:             "Thermal bridging (after refurbishment)",
		SimpleDescription: "Extra heat loss through junctions and corners, assuming refurbishment measures have been applied.",
		ExpertDescription: "Additional U-value due to thermal bridging for refurbished building elements."
	},
}
