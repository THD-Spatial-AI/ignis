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

func versionString() string {
	return fmt.Sprintf("%s (commit %s, built %s)", version.Version, version.Commit, version.Date)
}

// Setup app server and routes
func main() {
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.BoolVar(showVersion, "v", false, "print version and exit (shorthand)")
	flag.Parse()
	if *showVersion {
		fmt.Println(versionString())
		os.Exit(0)
	}

	utils.InitLogger()
	cfg := config.LoadConfig()
	app, cleanup := server.SetupServer()
	defer cleanup()
	app.Run(":" + cfg.App.Port)
}
