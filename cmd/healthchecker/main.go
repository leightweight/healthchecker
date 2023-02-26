package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/leightweight/healthchecker/internal/cli/http"
	"github.com/leightweight/healthchecker/internal/cli/wait"
	"github.com/urfave/cli/v2"
)

var (
	Version   = "DEV"
	BuildTime = time.Now().UTC().Format(time.RFC3339)
)

func main() {
	cli.HelpFlag = &cli.BoolFlag{
		Name:               "help",
		Usage:              "Show this help message",
		DisableDefaultText: true,
	}

	cli.VersionFlag = &cli.BoolFlag{
		Name:               "version",
		Usage:              "Print the version",
		DisableDefaultText: true,
		Aliases:            []string{"v"},
	}

	app := &cli.App{
		Name:    "healthchecker",
		Usage:   "Check the health of external services",
		Version: fmt.Sprintf("%s (built %s)", Version, BuildTime),

		EnableBashCompletion: true,

		Commands: []*cli.Command{
			http.Command(),
			wait.Command(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
