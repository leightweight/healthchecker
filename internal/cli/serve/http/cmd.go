package http

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/leightweight/healthchecker/internal/healthcheck"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:   "http",
		Usage:  "Checks service status with an HTTP request",
		Action: healthcheck.Build(check),

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
				Name:    "header",
				Usage:   "A header to add to the HTTP request",
				EnvVars: []string{"CHECK_HEADERS"},
			},
			&cli.BoolFlag{
				Name:    "allow-redirects",
				Usage:   "Allow HTTP response redirects",
				EnvVars: []string{"CHECK_ALLOW_REDIRECTS"},
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

func check(ctx *cli.Context) error {
	logger := log.With().Str("check", "http").Logger()
	url := ctx.String("url")
	method := ctx.String("method")
	headers := ctx.StringSlice("header")
	allowRedirects := ctx.Bool("allow-redirects")
	statusCodeRegex := ctx.String("status-code")
	responseRegex := ctx.String("response")
	timeout := ctx.Duration("timeout")

	logger.
		Trace().
		Bool("allowRedirects", allowRedirects).
		Dur("timeout", timeout).
		Msgf("Building HTTP client")

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if allowRedirects {
				return nil
			}
			return errors.New("tried to redirect")
		},
		Timeout: timeout,
	}

	logger.
		Trace().
		Str("method", method).
		Str("url", url).
		Msgf("Building HTTP request")

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		logger.Error().Err(err).Msgf("Error building HTTP request")
		return err
	}

	logger.Trace().Msgf("Adding %d request headers", len(headers))

	for _, header := range headers {
		split := strings.SplitN(header, ":", 2)
		if len(split) != 2 {
			logger.Error().Str("header", header).Msgf("Invalid HTTP header")
			return fmt.Errorf("invalid http header '%s'", header)
		}

		key := split[0]
		value := strings.TrimLeft(split[1], " ")

		logger.
			Trace().
			Str("key", key).
			Str("value", value).
			Msgf("Adding request header")

		req.Header.Add(key, value)
	}

	logger.Trace().Msgf("Executing HTTP request")

	res, err := client.Do(req)
	if err != nil {
		logger.Error().Err(err).Msgf("Error executing HTTP request")
		return err
	}
	defer res.Body.Close()

	logger.
		Debug().
		Int("statusCode", res.StatusCode).
		Str("regex", statusCodeRegex).
		Msgf("Checking status code")

	matched, err := regexp.MatchString(statusCodeRegex, strconv.Itoa(res.StatusCode))
	if err != nil {
		logger.Error().Err(err).Msgf("Error checking status code")
		return err
	}
	if !matched {
		logger.Warn().Msgf("Status code did not match")
		return fmt.Errorf("incorrect status code")
	}

	logger.Trace().Msgf("Reading response body")

	responseContent, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error().Err(err).Msgf("Error reading response body")
		return err
	}

	truncate := len(responseContent)
	if truncate > 50 {
		truncate = 50
	}

	logger.
		Debug().
		Bytes("body", responseContent[:truncate]).
		Str("regex", responseRegex).
		Msgf("Checking response body")

	matched, err = regexp.Match(responseRegex, responseContent)
	if err != nil {
		logger.Error().Err(err).Msgf("Error checking response body")
		return err
	}
	if !matched {
		logger.Warn().Msgf("Response body did not match")
		return fmt.Errorf("incorrect response body")
	}

	logger.Info().Msgf("Health check succeeded")
	return nil
}
