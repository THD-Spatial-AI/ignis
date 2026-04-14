package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/THD-Spatial-AI/hdcp-go/internal/db/repository"
	"github.com/THD-Spatial-AI/hdcp-go/internal/service"
	"github.com/THD-Spatial-AI/hdcp-go/internal/utils"

	"github.com/gin-gonic/gin"
)

// CalculateHeatDemand handles POST /api/v1/calculate/:code.
//
// :code is the TABULA Code_BuildingVariant, e.g. "DE.N.SFH.01.Gen".
// The first two characters of the code identify the country (ISO 3166-1 alpha-2).
//
// Optional JSON request body:
//
//	{ "A_ref": 150.0 }
//
// A_ref overrides the reference floor area stored in the TABULA record
// (A_C_Ref_Input). Omit the body to use the TABULA default.
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
		utils.Error.Printf("hdcp: failed to load TABULA data for %s: %v", variantCode, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load TABULA data"})
		return
	}

	// Apply optional body overrides. ShouldBindJSON tolerates an empty
	// JSON object {} — any present field replaces the TABULA default.
	var overrides struct {
		ARef *float64 `json:"A_ref"`
	}
	if err := c.ShouldBindJSON(&overrides); err != nil && err.Error() != "EOF" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}
	if overrides.ARef != nil {
		building.BasicParameters.Envelope.A_C_Ref_Input = *overrides.ARef
	}

	svc := service.NewHDCPService()
	qHND, err := svc.CalculateHeatingDemand(building)
	if err != nil {
		utils.Error.Printf("hdcp: pipeline failed for %s: %v", variantCode, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Pipeline execution failed"})
		return
	}

	utils.Info.Printf("hdcp: variant=%s q_h_nd=%.2f kWh/(m2.a)", variantCode, qHND)
	c.JSON(http.StatusOK, gin.H{
		"variant_code": variantCode,
		"q_h_nd":       qHND,
		"unit":         "kWh/(m2.a)",
	})
}
