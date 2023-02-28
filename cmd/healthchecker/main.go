package main

import (
	"fmt"
	"os"
	"time"

	"github.com/leightweight/healthchecker/internal/cli/check"
	"github.com/leightweight/healthchecker/internal/cli/serve"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var (
	Version   = "DEV"
	BuildTime = time.Now().UTC().Format(time.RFC3339)
)

func main() {
	initLogger()
	app := initCli()

	err := app.Run(os.Args)
	if err != nil {
		os.Exit(1)
	}
}

func initLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}

func initCli() *cli.App {
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

	return &cli.App{
		Name:    "healthchecker",
		Usage:   "Check the health of external services",
		Version: fmt.Sprintf("%s (built %s)", Version, BuildTime),

		Suggest:              true,
		EnableBashCompletion: true,

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "socket",
				Usage:   "The socket file for healthchecker to communicate over",
				Value:   "/tmp/healthchecker.sock",
				EnvVars: []string{"HEALTHCHECKER_SOCKET"},
			},
		},

		Commands: []*cli.Command{
			serve.Command(),
			check.Command(),
		},
	}
}
