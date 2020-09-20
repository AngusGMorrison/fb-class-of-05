// app.go contains the main function required to run the application.
package main

import (
	"angusgmorrison/fb05/internal/app/middleware"
	"angusgmorrison/fb05/internal/app/routing"
	"angusgmorrison/fb05/internal/app/templates"
	"angusgmorrison/fb05/pkg/env"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

const defaultEnv = "development"

var (
	// The key used to look up the current environment, e.g. development
	// prod, etc.
	envKey string
	// env is the current environment, e.g. development, prod, etc.
	currentEnv string
	configName string
	configType string
	configPath string
)

func main() {
	graceSeconds := flag.Duration(
		"graceSeconds",
		time.Second*15,
		"How long the server waits for existing connects to finish before shutting down.",
	)
	flag.StringVar(&envKey, "envKey", "FB05_ENV", "The key used to look up the current environment.")
	flag.StringVar(
		&configName,
		"configName",
		"environment",
		"The name of the config file (without extension) containing env vars required for the migration.",
	)
	flag.StringVar(
		&configType,
		"configType",
		"yaml",
		`The config file format, e.g. "yaml".`,
	)
	flag.StringVar(
		&configPath,
		"configPath",
		".",
		"The path to the config file containing env vars required for the migration.",
	)
	debug := flag.Bool("debug", false, "Toggle debug-level logging.")
	flag.Parse()

	if targetEnv := os.Getenv(envKey); targetEnv == "" {
		log.Printf("%-8s %s", "INFO",
			fmt.Sprintf("No value for env found at %q; defaulting to %q", envKey, defaultEnv))
		currentEnv = defaultEnv
	} else {
		currentEnv = targetEnv
	}

	done := make(chan struct{})
	defer close(done)

	err := loadPrerequisites(done)
	if err != nil {
		log.Fatalf("%-8s %v", "FATAL", err)
	}

	err = runServer(*debug, *graceSeconds)
	if err != nil {
		log.Printf("%-8s %v", "FATAL", err)
		os.Exit(1)
	}
}

func loadPrerequisites(done <-chan struct{}) error {
	fatalErrors := make(chan error, 1)
	wgDone := make(chan interface{})
	var wg sync.WaitGroup // start-up concurrently
	wg.Add(2)

	// Load environment.
	go func(configName, configType, configPath string) {
		envConfig := env.NewConfig(configName, configType, configPath, currentEnv)
		err := env.Load(envConfig)
		if err != nil {
			fatalErrors <- err
		}
		wg.Done()
	}(configName, configType, configPath)

	// Parse templates.
	go func() {
		templates.InitCache(done, filepath.Join("web", "templates"))
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(wgDone)
	}()

	select {
	case err := <-fatalErrors:
		return err
	case <-wgDone:
		return nil
	}
}

func runServer(debug bool, graceSeconds time.Duration) error {
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", env.Get("FB05_HOST"), env.Get("FB05_PORT")),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      configureRouter(debug),
	}

	fatalErrors := make(chan error, 1)
	go func() {
		log.Printf("%-8s Starting server at %s", "INFO", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			fatalErrors <- err
		}
	}()

	// Trigger graceful shutdown when quit with SIGINT (Ctrl+C).
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	select {
	case <-interrupt:
		log.Printf("%-8s %s", "INFO", "Shutting down gracefully...")
		ctx, cancel := context.WithTimeout(context.Background(), graceSeconds)
		defer cancel()
		srv.Shutdown(ctx)
	case err := <-fatalErrors:
		return err
	}

	return nil
}

func configureRouter(debug bool) *mux.Router {
	// Configure middleware
	l := log.New(os.Stderr, "FB05 ", log.Ldate|log.Ltime)
	dl := middleware.NewDebuggableLog(l, debug)
	mw := []func(http.Handler) http.Handler{
		middleware.Logging(dl),
		middleware.Recovery(dl),
	}
	return routing.Router(l, mw)
}
