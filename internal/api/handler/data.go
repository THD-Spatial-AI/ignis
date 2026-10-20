package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/thd-spatial-ai/ignis/internal/db/repository"
	"github.com/thd-spatial-ai/ignis/internal/utils"

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

// refurbishmentLabels are TABULA's names for its standard variants, keyed by Number_BuildingVariant.
var refurbishmentLabels = map[int]string{
	1: "Existing state",
	2: "Usual refurbishment",
	3: "Advanced refurbishment",
}

// defaultSubtype is TABULA's generic subtype, present for all but a few building types.
const defaultSubtype = "Gen"

// refurbishmentLabel returns TABULA's description of a variant, or TABULA's standard name for
// its number where the record has no description. In recent periods the standard variants are
// construction standards (e.g. "passive house"), which the standard names would misstate.
func refurbishmentLabel(v repository.RefurbishmentVariant) string {
	if v.Description != "" {
		return v.Description
	}
	return refurbishmentLabels[v.Number]
}

// MatchVariants returns TABULA's standard refurbishment variants (existing state, usual and
// advanced refurbishment) of one building subtype, ordered by variant number.
// Query params: type (e.g. SFH) is required; exactly one of period (e.g. 01) or
// year (e.g. 1975) must be given; subtype (e.g. Semi) is optional and defaults to Gen.
// year must be a non-negative integer no later than the current calendar year.
// With year, ignis resolves the period whose band contains it for that country and type;
// a year outside every defined band for that type clamps to the oldest or newest period
// rather than failing to match. A subtype without standard variants returns an empty list.
func (h *Handler) MatchVariants(c *gin.Context) {
	isoCode := strings.ToUpper(strings.TrimSpace(c.Param("country_iso2")))
	tableName, err := tableNameFromISO(isoCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	buildingType := strings.ToUpper(strings.TrimSpace(c.Query("type")))
	period := strings.TrimSpace(c.Query("period"))
	yearParam := strings.TrimSpace(c.Query("year"))
	subtype := strings.TrimSpace(c.Query("subtype"))
	if subtype == "" {
		subtype = defaultSubtype
	}

	if buildingType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query param 'type' is required"})
		return
	}
	if (period == "") == (yearParam == "") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "exactly one of 'period' or 'year' is required"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	// TABULA codes follow CC.N.TYPE.PERIOD.VariantSuffix — N is the national dataset identifier.
	typePrefix := fmt.Sprintf("%s.N.%s", isoCode, buildingType)

	if yearParam != "" {
		year, convErr := strconv.Atoi(yearParam)
		if convErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("query param 'year' must be an integer, got %q", yearParam)})
			return
		}
		if year < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("query param 'year' must not be negative, got %d", year)})
			return
		}
		if currentYear := time.Now().Year(); year > currentYear {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("query param 'year' must not be in the future, got %d", year)})
			return
		}
		period, err = h.repo.ResolvePeriodByYear(ctx, tableName, typePrefix, year)
		if err != nil {
			utils.Error.Printf("failed to resolve period for %s year %d: %v", typePrefix, year, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resolve construction period"})
			return
		}
		if period == "" {
			// No archetype of this type covers the year — an empty match list, not an error.
			c.JSON(http.StatusOK, gin.H{"country": tableName, "prefix": typePrefix, "subtype": subtype, "data": []any{}})
			return
		}
	}

	prefix := fmt.Sprintf("%s.%s.%s", typePrefix, period, subtype)

	variants, err := h.repo.MatchVariants(ctx, tableName, prefix)
	if err != nil {
		utils.Error.Printf("failed to match variants for %s: %v", prefix, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query variants"})
		return
	}

	type variantEntry struct {
		Code        string `json:"code"`
		Number      int    `json:"number"`
		Label       string `json:"label"`
		Description string `json:"description"`
	}
	entries := make([]variantEntry, len(variants))
	for i, v := range variants {
		entries[i] = variantEntry{Code: v.Code, Number: v.Number, Label: refurbishmentLabel(v), Description: v.Description}
	}

	c.JSON(http.StatusOK, gin.H{
		"country": tableName,
		"prefix":  prefix,
		"subtype": subtype,
		"data":    entries,
	})
}

// ListPeriods returns the country's construction-year bands, oldest first.
// year_from 0 means open-ended (oldest band); year_to 9999 means open-ended (newest band).
func (h *Handler) ListPeriods(c *gin.Context) {
	isoCode := strings.ToUpper(strings.TrimSpace(c.Param("country_iso2")))
	tableName, err := tableNameFromISO(isoCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), requestTimeout)
	defer cancel()

	periods, err := h.repo.ListPeriods(ctx, tableName)
	if err != nil {
		utils.Error.Printf("failed to load construction periods for %s: %v", tableName, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query construction periods"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"country": tableName,
		"data":    periods,
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
