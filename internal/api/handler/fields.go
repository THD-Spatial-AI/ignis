package handler

import (
	"net/http"

	"github.com/thd-spatial-ai/ignis/internal/models"

	"github.com/gin-gonic/gin"
)

// GetFieldMetadata returns the static description of every TABULA input field
// ignis's clients care about: where to find it in a GET /api/v1/data/:code
// response, its unit, and a simple/expert description pair. The list is
// identical for every country, since the underlying DB schema is uniform.
func (h *Handler) GetFieldMetadata(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": models.AllFieldMetadata,
	})
}
