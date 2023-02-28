package check

import (
	"errors"
	"io"
	"net"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:   "check",
		Usage:  "Performs the configured health check",
		Action: command,
	}
}

func command(ctx *cli.Context) error {
	socket := ctx.String("socket")

	log.Debug().Msgf("Dialing socket: unix://%s", socket)

	conn, err := net.Dial("unix", socket)
	if err != nil {
		log.Fatal().Err(err).Msgf("Error dialing socket")
		return err
	}
	defer conn.Close()

	log.Info().Msgf("Reading health check response")

	response, err := io.ReadAll(conn)
	if err != nil {
		log.Fatal().Err(err).Msgf("Error reading health check response")
		return err
	}
	if len(response) != 1 {
		log.Fatal().Msgf("Received incorrect health check response: %v", response)
		return errors.New("unexpected response")
	}

	if response[0] == 0 {
		log.Info().Msgf("Health check returned: healthy")
		return nil
	}

	log.Warn().Msgf("Health check returned: unhealthy")
	return errors.New("unhealthy")
}
