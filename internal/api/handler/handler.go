package handler

import (
	"github.com/THD-Spatial-AI/hdcp-go/internal/db/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Handler holds shared dependencies for all API handlers.
// A single instance is created at startup and shared across all requests.
type Handler struct {
	repo *repository.TabulaRepository
}

// New creates a Handler using the shared connection pool.
// pool must remain open for the lifetime of the process.
func New(pool *pgxpool.Pool, schema string) *Handler {
	return &Handler{repo: repository.NewTabulaRepository(pool, schema)}
}
