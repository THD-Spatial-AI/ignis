package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/thd-spatial-ai/ignis/internal/api/server"
	"github.com/thd-spatial-ai/ignis/internal/config"
	"github.com/thd-spatial-ai/ignis/internal/utils"
	"github.com/thd-spatial-ai/ignis/internal/version"
)

// Setup app server and routes
func main() {
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.BoolVar(showVersion, "v", false, "print version and exit (shorthand)")
	flag.Parse()
	if *showVersion {
		fmt.Printf("%s (commit %s, built %s)\n", version.Version, version.Commit, version.Date)
		os.Exit(0)
	}

	utils.InitLogger()
	cfg := config.LoadConfig()
	app, cleanup := server.SetupServer()
	defer cleanup()
	app.Run(":" + cfg.App.Port)
}
