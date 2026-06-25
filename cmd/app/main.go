package main

import (
	"github.com/thd-spatial-ai/ignis/internal/api/server"
	"github.com/thd-spatial-ai/ignis/internal/utils"
)

// Setup app server and routes
func main() {
	utils.InitLogger()
	app, cleanup := server.SetupServer()
	defer cleanup()
	app.Run(":8080")
}
