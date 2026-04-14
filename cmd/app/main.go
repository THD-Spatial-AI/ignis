package main

import (
	"github.com/THD-Spatial-AI/hdcp-go/internal/api/server"
	"github.com/THD-Spatial-AI/hdcp-go/internal/utils"
)

// Setup app server and routes
func main() {
	utils.InitLogger()
	app, cleanup := server.SetupServer()
	defer cleanup()
	app.Run(":8080")
}
