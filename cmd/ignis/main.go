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

// parseVersionFlag reports whether args request the version string, so the
// flag-handling decision can be tested without starting a server or exiting.
func parseVersionFlag(args []string) bool {
	fs := flag.NewFlagSet("ignis", flag.ExitOnError)
	showVersion := fs.Bool("version", false, "print version and exit")
	fs.BoolVar(showVersion, "v", false, "print version and exit (shorthand)")
	fs.Parse(args)
	return *showVersion
}

// Setup app server and routes
func main() {
	if parseVersionFlag(os.Args[1:]) {
		fmt.Println(versionString())
		os.Exit(0)
	}

	utils.InitLogger()
	cfg := config.LoadConfig()
	app, cleanup := server.SetupServer()
	defer cleanup()
	app.Run(":" + cfg.App.Port)
}
