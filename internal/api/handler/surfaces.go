package handler

import (
	"fmt"
	"math"

	"github.com/thd-spatial-ai/ignis/internal/models"
)

// Surface describes one physical building element (a wall, a window, ...) as
// produced by a real building's geometry, in contrast to TABULA's fixed
// per-category slots (A_Wall_1/2/3, U_Wall_1/2/3, ...). Its shape mirrors
// BuEM's envelope.elements[] so a caller can build both models' requests
// from the same source data.
type Surface struct {
	ID      string   `json:"id"`
	Type    string   `json:"type"`              // "roof" | "wall" | "floor" | "window" | "door"
	Area    float64  `json:"area"`              // m²
	UValue  float64  `json:"u_value"`           // W/m²K
	Azimuth *float64 `json:"azimuth,omitempty"` // degrees, 0=N/90=E/180=S/270=W; windows only, defaults to South (180)
	Tilt    *float64 `json:"tilt,omitempty"`    // degrees; windows only, ~0 = horizontal (skylight)
}

// horizontalTiltThreshold is how close Tilt must be to 0 for a window to be
// bucketed as horizontal (skylight) rather than by azimuth.
const horizontalTiltThreshold = 5.0

// aggregateSurfaces collapses an arbitrary list of physical surfaces into the
// area-weighted category totals ignis's pipeline already consumes
// (A_<Type>_1, U_<Type>_1, and the window orientation split), mutating
// building in place. An empty or nil list is a no-op, leaving the
// TABULA-sourced defaults on building untouched.
//
// The aggregation is exact, not approximate: U-values are thermal
// conductances (W/m²K), and for surfaces in the same category facing the
// same boundary, conductances add in parallel weighted by area, i.e.
// sum(area_i * U_i) == area_total * U_weighted_average. This is the same
// quantity TABULA's own methodology already represents with a single
// U_<Type>_1 per category, just computed from finer-grained input instead of
// being pre-computed once into the database row.
func aggregateSurfaces(surfaces []Surface, building *models.TabulaBuildingParameters) error {
	if len(surfaces) == 0 {
		return nil
	}

	totals := map[string]*weightedTotal{
		"roof":   {},
		"wall":   {},
		"floor":  {},
		"door":   {},
		"window": {},
	}
	windowOrientationArea := map[string]float64{
		"North": 0, "East": 0, "South": 0, "West": 0, "Horizontal": 0,
	}

	for _, s := range surfaces {
		total, ok := totals[s.Type]
		if !ok {
			return fmt.Errorf("surface %q: unknown type %q (must be roof, wall, floor, window, or door)", s.ID, s.Type)
		}
		if s.Area <= 0 {
			return fmt.Errorf("surface %q: area must be a positive number", s.ID)
		}
		if s.UValue <= 0 {
			return fmt.Errorf("surface %q: u_value must be a positive number", s.ID)
		}
		total.add(s.Area, s.UValue)

		if s.Type == "window" {
			windowOrientationArea[windowOrientationBucket(s)] += s.Area
		}
	}

	env := building.BasicParameters.Envelope
	uv := building.AdvancedParameters.Uvalues

	env.A_Roof_1, uv.U_Roof_1 = totals["roof"].resolve(env.A_Roof_1, uv.U_Roof_1)
	env.A_Wall_1, uv.U_Wall_1 = totals["wall"].resolve(env.A_Wall_1, uv.U_Wall_1)
	env.A_Floor_1, uv.U_Floor_1 = totals["floor"].resolve(env.A_Floor_1, uv.U_Floor_1)
	env.A_Door_1, uv.U_Door_1 = totals["door"].resolve(env.A_Door_1, uv.U_Door_1)
	env.A_Window_1, uv.U_Window_1 = totals["window"].resolve(env.A_Window_1, uv.U_Window_1)

	if totals["window"].totalArea > 0 {
		env.A_Window_North = windowOrientationArea["North"]
		env.A_Window_East = windowOrientationArea["East"]
		env.A_Window_South = windowOrientationArea["South"]
		env.A_Window_West = windowOrientationArea["West"]
		env.A_Window_Horizontal = windowOrientationArea["Horizontal"]
	}

	return nil
}

// windowOrientationBucket maps a window surface to one of ignis's five
// existing solar-gain buckets (North/East/South/West/Horizontal). Tilt near
// 0 wins (skylight); otherwise azimuth rounds to the nearest cardinal.
// A window with no azimuth given defaults to South.
func windowOrientationBucket(s Surface) string {
	if s.Tilt != nil && math.Abs(*s.Tilt) <= horizontalTiltThreshold {
		return "Horizontal"
	}
	if s.Azimuth == nil {
		return "South"
	}
	switch nearestCardinal(*s.Azimuth) {
	case 90:
		return "East"
	case 180:
		return "South"
	case 270:
		return "West"
	default:
		return "North"
	}
}

// nearestCardinal rounds an azimuth in degrees to the nearest of 0/90/180/270.
func nearestCardinal(azimuth float64) int {
	normalized := math.Mod(azimuth, 360)
	if normalized < 0 {
		normalized += 360
	}
	return int(math.Round(normalized/90)) % 4 * 90
}

// weightedTotal accumulates area and area-weighted U-value across the
// surfaces of one category.
type weightedTotal struct {
	totalArea  float64
	areaTimesU float64
}

func (w *weightedTotal) add(area, uValue float64) {
	w.totalArea += area
	w.areaTimesU += area * uValue
}

// resolve returns the aggregated (area, U-value) pair, or the given
// fallback values unchanged if no surface of this category was provided.
func (w *weightedTotal) resolve(fallbackArea, fallbackU float64) (float64, float64) {
	if w.totalArea == 0 {
		return fallbackArea, fallbackU
	}
	return w.totalArea, w.areaTimesU / w.totalArea
}
