package server

import (
	"context"

	"github.com/thd-spatial-ai/ignis/internal/api/handler"
	"github.com/thd-spatial-ai/ignis/internal/api/middleware"
	"github.com/thd-spatial-ai/ignis/internal/api/router"
	"github.com/thd-spatial-ai/ignis/internal/config"
	"github.com/thd-spatial-ai/ignis/internal/utils"

	"github.com/gin-gonic/gin"
)

// SetupServer creates the connection pool, wires handlers, and returns the
// configured engine together with a cleanup function that closes the pool.
// The caller must invoke cleanup() — typically via defer — before the process exits.
//
// Pool defaults: pgxpool uses max(4, runtime.NumCPU()) connections.
// Override with the DB_POOL_MAX_CONNS environment variable if needed.
func SetupServer() (*gin.Engine, func()) {
	utils.InitLogger()
	utils.Info.Println("Setting up server...")

	cfg := config.LoadConfig()
	if cfg.DB == nil {
		utils.Error.Fatal("database configuration missing — set DB_HOST, DB_NAME, DB_USER, DB_PASSWORD")
	}

	connString := utils.BuildConnectionString(cfg)
	pool, err := utils.ConnectPool(context.Background(), connString)
	if err != nil {
		utils.Error.Fatalf("failed to connect to database: %v", err)
	}
	utils.Info.Println("Database pool established")

	schema := ""
	if cfg.DB.Schemas != nil {
		schema = cfg.DB.Schemas.Tabula
	}

	h := handler.New(pool, schema)

	r := gin.New()
	r.Use(gin.Recovery(), middleware.CORS(), middleware.RequestBodyLimit(), middleware.RequestLogger())
	r.SetTrustedProxies(nil)

	router.RegisterRoutes(r, h)

	return r, pool.Close
}
