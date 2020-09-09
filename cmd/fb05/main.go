// app.go contains the main function required to run the application.
package main

import (
	"angusgmorrison/fb05/internal/app/middleware"
	"angusgmorrison/fb05/internal/app/routing"
	"angusgmorrison/fb05/internal/app/templates"
	"angusgmorrison/fb05/pkg/envloader"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
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
		log.Printf("%-8s %s", "INFO",
			fmt.Sprintf("No value for env found at %q; defaulting to %q", envKey, defaultEnv))
		env = defaultEnv
	} else {
		env = targetEnv
	}

	// Load environment variables.
	var err error
	envConfig := envloader.NewConfig(*configName, *configType, *configPath, env)
	envVars, err = envloader.Load(envConfig)
	if err != nil {
		log.Fatalf("%-8s %v", "FATAL", err)
	}

	// Parse templates.
	err = templates.Initialize(filepath.Join("internal", "app", "templates"))
	if err != nil {
		log.Fatalf("%-8s %v", "FATAL", err)
	}

	// Configure and launch server.
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", envVars.Get("FB05_HOST"), envVars.Get("FB05_PORT")),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      configureRouter(*debug),
	}

	srvErrors := make(chan error, 1)
	go func() {
		log.Printf("%-8s Starting server at %s", "INFO", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			srvErrors <- err
		}
	}()

	// Trigger graceful shutdown when quit with SIGINT (Ctrl+C).
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	select {
	case <-interrupt:
		log.Printf("%-8s %s", "INFO", "Shutting down gracefully...")
		ctx, cancel := context.WithTimeout(context.Background(), *graceSeconds)
		defer cancel()
		srv.Shutdown(ctx)
		os.Exit(0)
	case err := <-srvErrors:
		log.Printf("%-8s %v", "FATAL", err)
		os.Exit(1)
	}
}

func configureRouter(debug bool) *mux.Router {
	// Configure middleware
	l := log.New(os.Stderr, "FB05 ", log.Ldate|log.Ltime)
	dl := middleware.NewDebuggableLog(l, debug)
	mw := []mux.MiddlewareFunc{
		middleware.Logging(dl),
		middleware.Recovery(dl),
	}
	return routing.Router(l, mw)
}
