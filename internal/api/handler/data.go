package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/THD-Spatial-AI/hdcp-go/internal/db/repository"
	"github.com/THD-Spatial-AI/hdcp-go/internal/utils"

	"github.com/gin-gonic/gin"
)

const requestTimeout = 5 * time.Second

var tabulaCountryHelper = utils.NewTabulaCountryHelper()

// GetVariants lists all available building variants for a given country.
func (h *Handler) GetVariants(c *gin.Context) {
	isoCode := strings.ToUpper(strings.TrimSpace(c.Param("country_iso2")))
	tableName, err := tableNameFromISO(isoCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	variants, err := h.repo.ListVariants(ctx, tableName)
	if err != nil {
		utils.Error.Printf("failed to load variants for %s: %v", tableName, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query variants"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"country": tableName,
		"data":    variants,
	})
}

// GetVariantData retrieves TABULA data for a specific building variant.
func (h *Handler) GetVariantData(c *gin.Context) {
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

	utils.Info.Printf("Fetching TABULA data for variant %s in table %s", variantCode, tableName)

	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	building, buildingID, expectedQHND, err := h.repo.GetVariant(ctx, tableName, variantCode)
	if err != nil {
		if errors.Is(err, repository.ErrVariantNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Variant not found"})
			return
		}
		utils.Error.Printf("failed to load TABULA data for %s: %v", variantCode, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load TABULA data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"country":         tableName,
		"variant_code":    buildingID,
		"tabula_data":     building,
		"expected_q_h_nd": expectedQHND,
	})
}

// tableNameFromISO converts an ISO 3166-1 alpha-2 code to the TABULA table name.
// Returns an error if the code is unknown — CodeToCountry returns a lowercase fallback
// for unknown codes, so we verify the round-trip to detect them.
func tableNameFromISO(isoCode string) (string, error) {
	if len(isoCode) != 2 {
		return "", fmt.Errorf("invalid ISO2 code: %s", isoCode)
	}

	table := tabulaCountryHelper.CodeToCountry(isoCode)
	if tabulaCountryHelper.CountryToCode(table) != strings.ToUpper(isoCode) {
		return "", fmt.Errorf("no TABULA dataset configured for %s", isoCode)
	}

	return table, nil
}

// isoFromVariantCode extracts and validates the ISO 3166-1 alpha-2 prefix from a TABULA variant code.
// Valid codes follow the pattern "CC.something" (e.g. "DE.N.SFH.01.Gen").
func isoFromVariantCode(variantCode string) (string, error) {
	if len(variantCode) < 4 || variantCode[2] != '.' {
		return "", fmt.Errorf("invalid variant code %q: expected format CC.xxx", variantCode)
	}
	prefix := strings.ToUpper(variantCode[:2])
	for _, ch := range prefix {
		if ch < 'A' || ch > 'Z' {
			return "", fmt.Errorf("invalid variant code %q: country prefix must be two letters", variantCode)
		}
	}
	return prefix, nil
}
