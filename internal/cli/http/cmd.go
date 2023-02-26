package http

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:   "http",
		Usage:  "Calls an HTTP service to get its status",
		Action: command,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "url",
				Usage:    "The URL of the health check endpoint",
				Required: true,
				Aliases:  []string{"u"},
				EnvVars:  []string{"CHECK_URL"},
			},
			&cli.StringFlag{
				Name:    "method",
				Usage:   "The HTTP method of the health check endpoint",
				Value:   "GET",
				Aliases: []string{"m"},
				EnvVars: []string{"CHECK_METHOD"},
			},
			&cli.StringSliceFlag{
				Name:  "header",
				Usage: "A header to add to the HTTP request",
				//Aliases: []string{"h"},
				EnvVars: []string{"CHECK_HEADERS"},
			},
			&cli.BoolFlag{
				Name:    "allow-redirects",
				Usage:   "Allow HTTP response redirects",
				EnvVars: []string{"CHECK_ALLOW_REDIRECTS"},
			},
			&cli.DurationFlag{
				Name:    "timeout",
				Usage:   "The amount of time to wait for a response",
				Value:   30 * time.Second,
				EnvVars: []string{"CHECK_TIMEOUT"},
			},
			&cli.StringFlag{
				Name:    "status-code",
				Usage:   "Regular expression for the expected response status code",
				Value:   "^200$",
				EnvVars: []string{"CHECK_STATUS_CODE"},
			},
			&cli.StringFlag{
				Name:    "response",
				Usage:   "Regular expression for the expected response content",
				Value:   ".*",
				EnvVars: []string{"CHECK_RESPONSE"},
			},
		},
	}
}

func command(ctx *cli.Context) error {
	method := ctx.String("method")
	url := ctx.String("url")
	headers := ctx.StringSlice("header")
	statusCodeRegex := ctx.String("status-code")
	responseRegex := ctx.String("response")

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if ctx.Bool("allow-redirects") {
				return nil
			}
			return errors.New("tried to redirect")
		},
		Timeout: ctx.Duration("timeout"),
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	for _, header := range headers {
		split := strings.SplitN(header, ":", 2)
		if len(split) != 2 {
			return fmt.Errorf("invalid request header '%s'", header)
		}

		req.Header.Add(split[0], strings.TrimLeft(split[1], " "))
	}

	log.Printf("Executing HTTP request: %s %s", method, url)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	log.Printf("Check status code %d matches /%s/", res.StatusCode, statusCodeRegex)
	matched, err := regexp.MatchString(statusCodeRegex, strconv.Itoa(res.StatusCode))
	if err != nil {
		return err
	}
	if !matched {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	responseContent, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	log.Printf("Check response content matches /%s/", responseRegex)
	matched, err = regexp.Match(responseRegex, responseContent)
	if err != nil {
		return err
	}
	if !matched {
		return fmt.Errorf("unexpected response: %s", responseContent)
	}

	log.Printf("All checks succeeded")
	return nil
}
