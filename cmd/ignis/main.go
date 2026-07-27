package main

import (
	"github.com/thd-spatial-ai/ignis/internal/api/server"
	"github.com/thd-spatial-ai/ignis/internal/config"
	"github.com/thd-spatial-ai/ignis/internal/utils"
)

// Setup app server and routes
func main() {
	utils.InitLogger()
	cfg := config.LoadConfig()
	app, cleanup := server.SetupServer()
	defer cleanup()
	app.Run(":" + cfg.App.Port)
}
