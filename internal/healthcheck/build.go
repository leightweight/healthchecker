package healthcheck

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func Build(check cli.ActionFunc) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		var listener net.Listener
		var interrupted bool
		socket := ctx.String("socket")

		log.Trace().Msgf("Making signal channel")

		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

		go func() {
			log.Trace().Msgf("Waiting for signals")
			log.Trace().Msgf("Received signal: %v", <-ch)
			interrupted = true

			err := listener.Close()
			if err != nil {
				log.Panic().Err(err).Msgf("Error closing socket")
			}
		}()

		log.Debug().Msgf("Opening socket: %s", socket)

		listener, err := net.Listen("unix", socket)
		if err != nil {
			log.Fatal().Err(err).Msgf("Error opening socket")
			return err
		}
		defer func(listener net.Listener) {
			_ = listener.Close()
		}(listener)

		log.Info().Msgf("Waiting for health check connections")

		for {
			conn, err := listener.Accept()
			if err != nil && interrupted {
				log.Info().Msgf("Received interrupt, exiting")
				return nil
			} else if err != nil {
				log.Warn().Err(err).Msgf("Error accepting health check connection")
				continue
			}

			log.Info().Msgf("Executing health check")

			var response byte
			err = check(ctx)

			if err != nil {
				log.Error().Err(err).Msgf("Health check failed: %s", err)
				response = 1
			} else {
				log.Info().Msgf("Health check succeeded")
				response = 0
			}

			log.Trace().Msgf("Writing response to check client")

			_, err = conn.Write([]byte{response})
			if err != nil {
				log.Error().Err(err).Msgf("Error writing response back to checker")
			}

			_ = conn.Close()
		}
	}
}
