package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/vention/booking_api/internal/config"
	"github.com/vention/booking_api/pkg/server"
)

func main() {
	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	var cfgPath string
	flag.StringVar(&cfgPath, "cfg", "", "path where cfg is located")

	flag.Parse()

	if os.Getenv("CONFIG_PATH") != "" {
		cfgPath = os.Getenv("CONFIG_PATH")
	}

	// Initialize logger
	logger := log.With().Str("component", "api").Logger()

	cfg, err := config.New(cfgPath)
	if err != nil {
		logger.Fatal().Err(err).Msg("Error loading configuration")
	}

	go func() {
		if err := server.Start(ctx, cfg, logger); err != nil {
			logger.Fatal().Err(err).Msg("Error starting server")
		}
	}()

	// Handle shutdown signals
	signs := make(chan os.Signal, 1)
	signal.Notify(signs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	sign := <-signs
	logger.Info().Str("signal", sign.String()).Msg("shutting server down, since received signal to exit")

	// Cancel context and close active connections
	ctxCancel()
}
