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
//	  "delta_U_ThermalBridging_Refurbished": 0.05
//	}
//
// A_ref overrides the reference floor area (A_C_Ref_Input); the rest override
// the matching AdvancedParameters fields used by the climate/solar-gain/
// thermal-bridging calc levels. Omit the body entirely to use TABULA defaults
// throughout.
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
		ARef                             *float64 `json:"A_ref"`
		HeatingDays                      *int     `json:"HeatingDays"`
		ThetaE                           *float64 `json:"Theta_e"`
		ThetaI                           *float64 `json:"theta_i"`
		ISolSouth                        *float64 `json:"I_Sol_South"`
		ISolEast                         *float64 `json:"I_Sol_East"`
		ISolWest                         *float64 `json:"I_Sol_West"`
		ISolNorth                        *float64 `json:"I_Sol_North"`
		ISolHorizontal                   *float64 `json:"I_Sol_Hor"`
		DeltaUThermalBridgingOriginal    *float64 `json:"delta_U_ThermalBridging_Original"`
		DeltaUThermalBridgingRefurbished *float64 `json:"delta_U_ThermalBridging_Refurbished"`
	}
	if err := c.ShouldBindJSON(&overrides); err != nil && !errors.Is(err, io.EOF) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
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
