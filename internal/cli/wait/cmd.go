package wait

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:   "wait",
		Usage:  "Waits for an interrupt and then quits",
		Action: command,
	}
}

func command(ctx *cli.Context) error {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

	log.Printf("Waiting for signals on PID %d...", os.Getpid())
	log.Printf("Received '%s'. Exiting.", <-ch)

	return nil
}
