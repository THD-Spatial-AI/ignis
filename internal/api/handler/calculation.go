package handler

import (
	"context"
	"errors"
	"io"
	"math"
	"net/http"
	"strings"

	"github.com/thd-spatial-ai/ignis/internal/db/repository"
	"github.com/thd-spatial-ai/ignis/internal/service"
	"github.com/thd-spatial-ai/ignis/internal/utils"

	"github.com/gin-gonic/gin"
)

// CalculateHeatDemand handles POST /api/v1/calculate/:code.
//
// :code is the TABULA Code_BuildingVariant, e.g. "DE.N.SFH.01.Gen".
// The first two characters of the code identify the country (ISO 3166-1 alpha-2).
//
// Optional JSON request body — every field overrides the corresponding TABULA
// record default, omit any field to keep the TABULA value:
//
//	{
//	  "A_ref": 150.0,
//	  "HeatingDays": 200,
//	  "Theta_e": -5.0,
//	  "theta_i": 20.0,
//	  "I_Sol_South": 400.0,
//	  "I_Sol_East": 150.0,
//	  "I_Sol_West": 150.0,
//	  "I_Sol_North": 80.0,
//	  "I_Sol_Hor": 500.0,
//	  "delta_U_ThermalBridging_Original": 0.1,
//	  "delta_U_ThermalBridging_Refurbished": 0.05,
//	  "h_room": 2.8,
//	  "n_Storey": 3,
//	  "surfaces": [
//	    {"id": "wall-1", "type": "wall", "area": 60.0, "u_value": 0.8, "azimuth": 180},
//	    {"id": "win-1", "type": "window", "area": 8.0, "u_value": 1.2, "azimuth": 90}
//	  ]
//	}
//
// A_ref overrides the reference floor area (A_C_Ref_Input); the rest override
// the matching AdvancedParameters fields used by the climate/solar-gain/
// thermal-bridging calc levels. Omit the body entirely to use TABULA defaults
// throughout.
//
// h_room overrides the archetype's assumed room height (BuildingAppearance.H_room,
// meters) — feeds the ventilation heat transfer coefficient directly (calc_level_01.go),
// so this is often a large lever: TABULA's generic default is frequently a poor
// match for a specific real building (a lecture hall or workshop can easily run
// 4-5m against a 2.5m archetype default). n_Storey overrides the assumed storey
// count (BuildingAppearance.N_Storey), which feeds the envelope-area *estimation*
// path (calc_level_02.go/03.go/06.go) — has less effect once real surfaces are
// also given, since those bypass estimation for the categories they cover.
//
// surfaces replaces TABULA's fixed 2-3 slots per element category with an
// arbitrary list of individual physical surfaces (as a real building's
// geometry has, rather than a generic archetype). Surfaces are grouped by
// type ("roof", "wall", "floor", "window", or "door") and collapsed into an
// area-weighted equivalent (summed area, area-weighted average U-value)
// before the pipeline runs — mathematically exact, since U-values are
// conductances and conductances in parallel add weighted by area. Window
// surfaces are additionally bucketed by azimuth into ignis's five existing
// solar-gain orientations (nearest of North/East/South/West, or Horizontal
// if tilt is near 0); a window with no azimuth given defaults to South.
// Any category with no surfaces given keeps its TABULA default.
//
// Caveat: for a variant with a nonzero measure fraction (i.e. not "Existing
// state" — check GET /api/v1/data/:code), only the raw U-value term is
// affected by surfaces or U_<Type>_1 overrides; the separate refurbishment
// blend term is untouched, so the override only partially moves the result.
//
// Response:
//
//	{ "variant_code": "DE.N.SFH.01.Gen", "q_h_nd": 123.45, "unit": "kWh/(m2.a)" }
func (h *Handler) CalculateHeatDemand(c *gin.Context) {
	variantCode := strings.TrimSpace(c.Param("code"))

	isoCode, err := isoFromVariantCode(variantCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tableName, err := tableNameFromISO(isoCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	building, _, _, err := h.repo.GetVariant(ctx, tableName, variantCode)
	if err != nil {
		if errors.Is(err, repository.ErrVariantNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "variant not found: " + variantCode})
			return
		}
		utils.Error.Printf("ignis: failed to load TABULA data for %s: %v", variantCode, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load TABULA data"})
		return
	}

	// Apply optional body overrides. ShouldBindJSON tolerates an empty
	// JSON object {} — any present field replaces the TABULA default.
	var overrides struct {
		ARef                             *float64  `json:"A_ref"`
		HeatingDays                      *int      `json:"HeatingDays"`
		ThetaE                           *float64  `json:"Theta_e"`
		ThetaI                           *float64  `json:"theta_i"`
		ISolSouth                        *float64  `json:"I_Sol_South"`
		ISolEast                         *float64  `json:"I_Sol_East"`
		ISolWest                         *float64  `json:"I_Sol_West"`
		ISolNorth                        *float64  `json:"I_Sol_North"`
		ISolHorizontal                   *float64  `json:"I_Sol_Hor"`
		DeltaUThermalBridgingOriginal    *float64  `json:"delta_U_ThermalBridging_Original"`
		DeltaUThermalBridgingRefurbished *float64  `json:"delta_U_ThermalBridging_Refurbished"`
		HRoom                            *float64  `json:"h_room"`
		NStorey                          *int      `json:"n_Storey"`
		Surfaces                         []Surface `json:"surfaces"`
	}
	if err := c.ShouldBindJSON(&overrides); err != nil && !errors.Is(err, io.EOF) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	if err := aggregateSurfaces(overrides.Surfaces, building); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if overrides.ARef != nil {
		if *overrides.ARef <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "A_ref must be a positive number"})
			return
		}
		building.BasicParameters.Envelope.A_C_Ref_Input = *overrides.ARef
	}
	if overrides.HeatingDays != nil {
		if *overrides.HeatingDays < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "HeatingDays must not be negative"})
			return
		}
		building.AdvancedParameters.ClimateConditions.HeatingDays = *overrides.HeatingDays
	}
	if overrides.ThetaE != nil {
		building.AdvancedParameters.ClimateConditions.Theta_e = *overrides.ThetaE
	}
	if overrides.ThetaI != nil {
		building.AdvancedParameters.ClimateConditions.Theta_i = *overrides.ThetaI
	}
	if overrides.ISolSouth != nil {
		if *overrides.ISolSouth < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "I_Sol_South must not be negative"})
			return
		}
		building.AdvancedParameters.SolarGains.I_Sol_South = *overrides.ISolSouth
	}
	if overrides.ISolEast != nil {
		if *overrides.ISolEast < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "I_Sol_East must not be negative"})
			return
		}
		building.AdvancedParameters.SolarGains.I_Sol_East = *overrides.ISolEast
	}
	if overrides.ISolWest != nil {
		if *overrides.ISolWest < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "I_Sol_West must not be negative"})
			return
		}
		building.AdvancedParameters.SolarGains.I_Sol_West = *overrides.ISolWest
	}
	if overrides.ISolNorth != nil {
		if *overrides.ISolNorth < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "I_Sol_North must not be negative"})
			return
		}
		building.AdvancedParameters.SolarGains.I_Sol_North = *overrides.ISolNorth
	}
	if overrides.ISolHorizontal != nil {
		if *overrides.ISolHorizontal < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "I_Sol_Hor must not be negative"})
			return
		}
		building.AdvancedParameters.SolarGains.I_Sol_Horizontal = *overrides.ISolHorizontal
	}
	if overrides.DeltaUThermalBridgingOriginal != nil {
		if *overrides.DeltaUThermalBridgingOriginal < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "delta_U_ThermalBridging_Original must not be negative"})
			return
		}
		building.AdvancedParameters.ThermalBridges.Delta_U_ThermalBridging_Original = *overrides.DeltaUThermalBridgingOriginal
	}
	if overrides.DeltaUThermalBridgingRefurbished != nil {
		if *overrides.DeltaUThermalBridgingRefurbished < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "delta_U_ThermalBridging_Refurbished must not be negative"})
			return
		}
		building.AdvancedParameters.ThermalBridges.Delta_U_ThermalBridging_Refurbished = *overrides.DeltaUThermalBridgingRefurbished
	}
	if overrides.HRoom != nil {
		if *overrides.HRoom <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "h_room must be a positive number"})
			return
		}
		building.BasicParameters.BuildingAppearance.H_room = *overrides.HRoom
	}
	if overrides.NStorey != nil {
		if *overrides.NStorey <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "n_Storey must be a positive number"})
			return
		}
		building.BasicParameters.BuildingAppearance.N_Storey = *overrides.NStorey
	}

	svc := service.NewIgnisService()
	qHND, err := svc.CalculateHeatingDemand(building)
	if err != nil {
		utils.Error.Printf("ignis: pipeline failed for %s: %v", variantCode, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Pipeline execution failed"})
		return
	}

	if math.IsNaN(qHND) || math.IsInf(qHND, 0) {
		utils.Error.Printf("ignis: pipeline returned non-finite value for %s: %v", variantCode, qHND)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Pipeline returned invalid result"})
		return
	}

	// Rounded here to two decimal places for cleaner output; the underlying calculation uses full precision.
	qHND = math.Round(qHND*100) / 100

	utils.Info.Printf("ignis: variant=%s q_h_nd=%.2f kWh/(m2.a)", variantCode, qHND)
	c.JSON(http.StatusOK, gin.H{
		"variant_code": variantCode,
		"q_h_nd":       qHND,
		"unit":         "kWh/(m2.a)",
	})
}
