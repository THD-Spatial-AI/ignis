package utils

import (
	"strings"
)

// TabulaCountryHelper handles conversion between country codes and names
// This is a Go implementation based on the Python version using pycountry
type TabulaCountryHelper struct{}

// NewTabulaCountryHelper creates a new instance of TabulaCountryHelper
func NewTabulaCountryHelper() *TabulaCountryHelper {
	return &TabulaCountryHelper{}
}

// countryMapping provides ISO2 code to country name mapping
// This maps ISO 3166-1 alpha-2 codes to standardized country names
var countryMapping = map[string]string{
	"AD": "andorra",
	"AE": "united_arab_emirates",
	"AF": "afghanistan",
	"AG": "antigua_and_barbuda",
	"AI": "anguilla",
	"AL": "albania",
	"AM": "armenia",
	"AO": "angola",
	"AQ": "antarctica",
	"AR": "argentina",
	"AS": "american_samoa",
	"AT": "austria",
	"AU": "australia",
	"AW": "aruba",
	"AX": "aland_islands",
	"AZ": "azerbaijan",
	"BA": "bosnia_and_herzegovina",
	"BB": "barbados",
	"BD": "bangladesh",
	"BE": "belgium",
	"BF": "burkina_faso",
	"BG": "bulgaria",
	"BH": "bahrain",
	"BI": "burundi",
	"BJ": "benin",
	"BL": "saint_barthelemy",
	"BM": "bermuda",
	"BN": "brunei",
	"BO": "bolivia",
	"BQ": "bonaire_sint_eustatius_and_saba",
	"BR": "brazil",
	"BS": "bahamas",
	"BT": "bhutan",
	"BV": "bouvet_island",
	"BW": "botswana",
	"BY": "belarus",
	"BZ": "belize",
	"CA": "canada",
	"CC": "cocos_keeling_islands",
	"CD": "congo_democratic_republic",
	"CF": "central_african_republic",
	"CG": "congo",
	"CH": "switzerland",
	"CI": "cote_divoire",
	"CK": "cook_islands",
	"CL": "chile",
	"CM": "cameroon",
	"CN": "china",
	"CO": "colombia",
	"CR": "costa_rica",
	"CU": "cuba",
	"CV": "cape_verde",
	"CW": "curacao",
	"CX": "christmas_island",
	"CY": "cyprus",
	"CZ": "czechia",
	"DE": "germany",
	"DJ": "djibouti",
	"DK": "denmark",
	"DM": "dominica",
	"DO": "dominican_republic",
	"DZ": "algeria",
	"EC": "ecuador",
	"EE": "estonia",
	"EG": "egypt",
	"EH": "western_sahara",
	"ER": "eritrea",
	"ES": "spain",
	"ET": "ethiopia",
	"FI": "finland",
	"FJ": "fiji",
	"FK": "falkland_islands",
	"FM": "micronesia",
	"FO": "faroe_islands",
	"FR": "france",
	"GA": "gabon",
	"GB": "united_kingdom",
	"GD": "grenada",
	"GE": "georgia",
	"GF": "french_guiana",
	"GG": "guernsey",
	"GH": "ghana",
	"GI": "gibraltar",
	"GL": "greenland",
	"GM": "gambia",
	"GN": "guinea",
	"GP": "guadeloupe",
	"GQ": "equatorial_guinea",
	"GR": "greece",
	"GS": "south_georgia_and_south_sandwich_islands",
	"GT": "guatemala",
	"GU": "guam",
	"GW": "guinea_bissau",
	"GY": "guyana",
	"HK": "hong_kong",
	"HM": "heard_island_and_mcdonald_islands",
	"HN": "honduras",
	"HR": "croatia",
	"HT": "haiti",
	"HU": "hungary",
	"ID": "indonesia",
	"IE": "ireland",
	"IL": "israel",
	"IM": "isle_of_man",
	"IN": "india",
	"IO": "british_indian_ocean_territory",
	"IQ": "iraq",
	"IR": "iran",
	"IS": "iceland",
	"IT": "italy",
	"JE": "jersey",
	"JM": "jamaica",
	"JO": "jordan",
	"JP": "japan",
	"KE": "kenya",
	"KG": "kyrgyzstan",
	"KH": "cambodia",
	"KI": "kiribati",
	"KM": "comoros",
	"KN": "saint_kitts_and_nevis",
	"KP": "north_korea",
	"KR": "south_korea",
	"KW": "kuwait",
	"KY": "cayman_islands",
	"KZ": "kazakhstan",
	"LA": "laos",
	"LB": "lebanon",
	"LC": "saint_lucia",
	"LI": "liechtenstein",
	"LK": "sri_lanka",
	"LR": "liberia",
	"LS": "lesotho",
	"LT": "lithuania",
	"LU": "luxembourg",
	"LV": "latvia",
	"LY": "libya",
	"MA": "morocco",
	"MC": "monaco",
	"MD": "moldova",
	"ME": "montenegro",
	"MF": "saint_martin",
	"MG": "madagascar",
	"MH": "marshall_islands",
	"MK": "north_macedonia",
	"ML": "mali",
	"MM": "myanmar",
	"MN": "mongolia",
	"MO": "macao",
	"MP": "northern_mariana_islands",
	"MQ": "martinique",
	"MR": "mauritania",
	"MS": "montserrat",
	"MT": "malta",
	"MU": "mauritius",
	"MV": "maldives",
	"MW": "malawi",
	"MX": "mexico",
	"MY": "malaysia",
	"MZ": "mozambique",
	"NA": "namibia",
	"NC": "new_caledonia",
	"NE": "niger",
	"NF": "norfolk_island",
	"NG": "nigeria",
	"NI": "nicaragua",
	"NL": "netherlands",
	"NO": "norway",
	"NP": "nepal",
	"NR": "nauru",
	"NU": "niue",
	"NZ": "new_zealand",
	"OM": "oman",
	"PA": "panama",
	"PE": "peru",
	"PF": "french_polynesia",
	"PG": "papua_new_guinea",
	"PH": "philippines",
	"PK": "pakistan",
	"PL": "poland",
	"PM": "saint_pierre_and_miquelon",
	"PN": "pitcairn",
	"PR": "puerto_rico",
	"PS": "palestine",
	"PT": "portugal",
	"PW": "palau",
	"PY": "paraguay",
	"QA": "qatar",
	"RE": "reunion",
	"RO": "romania",
	"RS": "serbia",
	"RU": "russia",
	"RW": "rwanda",
	"SA": "saudi_arabia",
	"SB": "solomon_islands",
	"SC": "seychelles",
	"SD": "sudan",
	"SE": "sweden",
	"SG": "singapore",
	"SH": "saint_helena_ascension_and_tristan_da_cunha",
	"SI": "slovenia",
	"SJ": "svalbard_and_jan_mayen",
	"SK": "slovakia",
	"SL": "sierra_leone",
	"SM": "san_marino",
	"SN": "senegal",
	"SO": "somalia",
	"SR": "suriname",
	"SS": "south_sudan",
	"ST": "sao_tome_and_principe",
	"SV": "el_salvador",
	"SX": "sint_maarten",
	"SY": "syria",
	"SZ": "eswatini",
	"TC": "turks_and_caicos_islands",
	"TD": "chad",
	"TF": "french_southern_territories",
	"TG": "togo",
	"TH": "thailand",
	"TJ": "tajikistan",
	"TK": "tokelau",
	"TL": "timor_leste",
	"TM": "turkmenistan",
	"TN": "tunisia",
	"TO": "tonga",
	"TR": "turkey",
	"TT": "trinidad_and_tobago",
	"TV": "tuvalu",
	"TW": "taiwan",
	"TZ": "tanzania",
	"UA": "ukraine",
	"UG": "uganda",
	"UM": "united_states_minor_outlying_islands",
	"US": "united_states",
	"UY": "uruguay",
	"UZ": "uzbekistan",
	"VA": "vatican_city",
	"VC": "saint_vincent_and_the_grenadines",
	"VE": "venezuela",
	"VG": "virgin_islands_british",
	"VI": "virgin_islands_us",
	"VN": "vietnam",
	"VU": "vanuatu",
	"WF": "wallis_and_futuna",
	"WS": "samoa",
	"YE": "yemen",
	"YT": "mayotte",
	"ZA": "south_africa",
	"ZM": "zambia",
	"ZW": "zimbabwe",
}

// reverseCountryMapping provides country name to ISO2 code mapping
var reverseCountryMapping map[string]string

// init initializes the reverse mapping
func init() {
	reverseCountryMapping = make(map[string]string)
	for code, country := range countryMapping {
		reverseCountryMapping[country] = code
	}
}

// CodeToCountry converts a country code to its corresponding country name
// This method takes a 2-letter ISO country code and returns the corresponding country name
// in lowercase and replaces spaces in the country name with an underscore ("_").
//
// Output from this method can be passed into CountryToCode method to get back the country code.
//
// Parameters:
//   - code: A 2-letter country code (e.g., 'US' for United States).
//
// Returns:
//   - The name of the country in lowercase (e.g., 'united_states') or empty string if not found.
func (h *TabulaCountryHelper) CodeToCountry(code string) string {
	code = strings.ToUpper(strings.TrimSpace(code))
	if country, exists := countryMapping[code]; exists {
		return country
	}

	// Log warning if logger is available
	Warn.Printf("Country code '%s' not found", code)
	return strings.ToLower(code) // Default to lowercase version of the code
}

// CountryToCode converts a country name to its corresponding country code
// This method takes a country name and returns the corresponding 2-letter ISO country code
// in uppercase. If the country name is invalid or not found, it returns empty string.
//
// Output from this method can be passed into CodeToCountry method to get back the country name.
//
// Parameters:
//   - country: The name of the country (e.g., 'united_states').
//
// Returns:
//   - The 2-letter country code (e.g., 'US') or empty string if not found.
func (h *TabulaCountryHelper) CountryToCode(country string) string {
	country = strings.ToLower(strings.TrimSpace(country))
	country = strings.ReplaceAll(country, " ", "_")

	if code, exists := reverseCountryMapping[country]; exists {
		return code
	}

	// Log warning if logger is available
	Warn.Printf("Country name '%s' not found", country)
	return ""
}

// GetAllCountryCodes returns all available ISO2 country codes
func (h *TabulaCountryHelper) GetAllCountryCodes() []string {
	codes := make([]string, 0, len(countryMapping))
	for code := range countryMapping {
		codes = append(codes, code)
	}
	return codes
}

// GetAllCountryNames returns all available country names
func (h *TabulaCountryHelper) GetAllCountryNames() []string {
	names := make([]string, 0, len(countryMapping))
	for _, name := range countryMapping {
		names = append(names, name)
	}
	return names
}

// IsValidCountryCode checks if the given code is a valid ISO2 country code
func (h *TabulaCountryHelper) IsValidCountryCode(code string) bool {
	code = strings.ToUpper(strings.TrimSpace(code))
	_, exists := countryMapping[code]
	return exists
}

// IsValidCountryName checks if the given name is a valid country name
func (h *TabulaCountryHelper) IsValidCountryName(country string) bool {
	country = strings.ToLower(strings.TrimSpace(country))
	country = strings.ReplaceAll(country, " ", "_")
	_, exists := reverseCountryMapping[country]
	return exists
}
