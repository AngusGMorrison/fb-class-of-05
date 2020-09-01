// app.go contains the main function required to run the application.
package main

import (
	"angusgmorrison/fb05/internal/app/middleware"
	"angusgmorrison/fb05/internal/app/middleware/httplog"
	"angusgmorrison/fb05/internal/app/middleware/stacklog"
	"angusgmorrison/fb05/internal/app/routing"
	"angusgmorrison/fb05/pkg/envloader"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/spf13/viper"
)

const defaultEnv = "development"

var (
	// The key used to look up the current environment, e.g. development
	// prod, etc.
	envKey string
	// env is the current environment, e.g. development, prod, etc.
	env string
	// envVars holds environment variables for the current environment.
	envVars *viper.Viper
)

func main() {
	graceSeconds := flag.Duration(
		"graceSeconds",
		time.Second*15,
		"How long the server waits for existing connects to finish before shutting down.",
	)
	flag.StringVar(&envKey, "envKey", "FB05_ENV", "The key used to look up the current environment.")
	configName := flag.String(
		"configName",
		"environment",
		"The name of the config file (without extension) containing env vars required for the migration.",
	)
	configType := flag.String(
		"configType",
		"yaml",
		`The config file format, e.g. "yaml".`,
	)
	configPath := flag.String(
		"configPath",
		".",
		"The path to the config file containing env vars required for the migration.",
	)
	debug := flag.Bool("debug", false, "Toggle debug-level logging.")
	flag.Parse()

	if targetEnv := os.Getenv(envKey); targetEnv == "" {
		log.Warn().
			Msg(fmt.Sprintf("No value for env found at %q; defaulting to %q", envKey, defaultEnv))
		env = defaultEnv
	} else {
		env = targetEnv
	}

	configureLogger(*debug)

	// Load environment variables.
	var err error
	envConfig := envloader.NewConfig(*configName, *configType, *configPath, env)
	envVars, err = envloader.Load(envConfig)
	if err != nil {
		log.Error().Err(err)
		os.Exit(1)
	}

	// Configure and launch server.
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", envVars.Get("FB05_HOST"), envVars.Get("FB05_PORT")),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      configureRouter(),
	}

	go func() {
		log.Info().Str("addr", srv.Addr).Msg("Starting server")
		if err := srv.ListenAndServe(); err != nil {
			log.Error().Err(err)
		}
	}()

	// Trigger graceful shutdown when quit with SIGINT (Ctrl+C).
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt

	log.Info().Msg("Shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), *graceSeconds)
	defer cancel()
	srv.Shutdown(ctx)
	os.Exit(0)
}

func configureLogger(debug bool) {
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Enable stack trace logging using Stack().
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	if env == "development" {
		// Pretty print logs.
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

func configureRouter() *mux.Router {
	// Configure middleware
	globalLogger := &log.Logger
	mw := []mux.MiddlewareFunc{
		middleware.Logging(httplog.NewLogger(globalLogger, env)),
		middleware.Recovery(stacklog.NewLogger(globalLogger, env)),
	}
	return routing.Router(mw)
}
