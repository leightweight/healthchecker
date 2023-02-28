package serve

import (
	"time"

	"github.com/leightweight/healthchecker/internal/cli/serve/http"
	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "serve",
		Usage: "Waits for a check request and then performs the configured health check",

		Flags: []cli.Flag{
			&cli.DurationFlag{
				Name:    "timeout",
				Usage:   "The amount of time to wait for a response",
				Value:   30 * time.Second,
				EnvVars: []string{"CHECK_TIMEOUT"},
			},
		},

		Subcommands: []*cli.Command{
			http.Command(),
		},
	}
}
