package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vention/booking_api/internal/config"
	"github.com/vention/booking_api/internal/util"
	"github.com/vention/booking_api/pkg/server"
)

func main() {
	// Configure logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Initialize timezone
	if err := util.InitTimezone(); err != nil {
		log.Warn().Err(err).Msg("Failed to load timezone, falling back to UTC")
		// Debug: List available timezones
		log.Debug().Msg("Available timezones in container:")
		if data, err := os.ReadFile("/usr/share/zoneinfo/zone.tab"); err == nil {
			lines := strings.Split(string(data), "\n")
			for i, line := range lines {
				if i < 10 && strings.Contains(line, "Europe") { // Show first 10 Europe timezones
					log.Debug().Str("timezone", line).Msg("Available timezone")
				}
			}
		}
	} else {
		log.Info().Str("timezone", util.GetAppTimezone().String()).Msg("Timezone initialized")
	}

	// Create context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Info().Str("signal", sig.String()).Msg("Received shutdown signal")
		cancel()
	}()

	// Start server
	log.Info().Msg("Starting booking API server...")
	if err := server.Start(ctx, cfg, log.Logger); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}

	log.Info().Msg("Server stopped")
}
