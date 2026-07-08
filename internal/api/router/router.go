package router

import (
	"github.com/thd-spatial-ai/ignis/internal/api/handler"
	"github.com/thd-spatial-ai/ignis/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes initializes the router with all routes.
// h holds the shared database connection and is injected at startup.
func RegisterRoutes(r *gin.Engine, h *handler.Handler) {
	utils.Info.Println("Setting up routes...")
	r.GET("/favicon.ico", func(c *gin.Context) { c.Status(http.StatusNoContent) })
	r.GET("/health", HealthCheck)

	api := r.Group("/api")
	{
		// Version 1 routes
		v1 := api.Group("/v1")
		{
			// GET
			v1.GET("/data/:code", h.GetVariantData)
			v1.GET("/variants/:country_iso2", h.GetVariants)
			v1.GET("/variants/:country_iso2/match", h.MatchVariants)
			v1.GET("/fields", h.GetFieldMetadata)

			// POST
			v1.POST("/calculate/:code", h.CalculateHeatDemand)
		}
	}
}

// HealthCheck handles health check requests
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}
