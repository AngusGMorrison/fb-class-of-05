// migrate wraps golang-migrate, allowing environment variables to be
// loaded from a config file before running migrations.
package main

// import (
// 	"errors"
// 	"flag"
// 	"fmt"
// 	"io"
// 	"os"
// 	"strconv"
// 	"strings"

// 	"github.com/golang-migrate/migrate/v4"
// 	_ "github.com/golang-migrate/migrate/v4/database/postgres"
// 	_ "github.com/golang-migrate/migrate/v4/source/file"
// 	"github.com/rs/zerolog/log"
// 	"github.com/spf13/viper"
// )

// // env is the current environment, e.g. development, prod, etc.
// var env string

// // envVars stores environment variables loaded by viper. The variables
// // loaded are specific to a single environment, e.g. development.
// var envVars *viper.Viper

// var out io.Writer = os.Stdout

// func main() {
// 	// Flags to control environment variable loading.
// 	envKey := flag.String(
// 		"envKey",
// 		"FB05_ENV",
// 		"The key used to look up the current environment, e.g., development, prod, etc.",
// 	)
// 	configName := flag.String(
// 		"configName",
// 		"environment",
// 		"The name of the config file (without extension) containing env vars required for the migration.",
// 	)
// 	configType := flag.String(
// 		"configType",
// 		"yaml",
// 		"The config file format, e.g. \"yaml\".",
// 	)
// 	configPath := flag.String(
// 		"configPath",
// 		".",
// 		"The path to the config file containing env vars required for the migration.",
// 	)
// 	migrationDir := flag.String(
// 		"migrationDir",
// 		"db/migrations",
// 		"The directory where migration files are stored.",
// 	)
// 	flag.Parse()

// 	// Ensure a command has been provided.
// 	args := flag.Args()
// 	if len(args) == 0 {
// 		fmt.Println("migrate must be used with a migration command")
// 		fmt.Println("Usage: migrate down | drop | up | version | force number | step number | toVersion number [-configName string] [-configPath string]")
// 		os.Exit(1)
// 	}

// 	// Load environment variables.
// 	var err error
// 	env = os.Getenv(*envKey)
// 	envConfig := envloader.NewConfig(*configName, *configType, *configPath, env)
// 	envVars, err = envloader.Load(envConfig)
// 	if err != nil {
// 		log.Error().Err(err)
// 		os.Exit(1)
// 	}

// 	// Create the migrator.
// 	migrationPath := fmt.Sprintf("file://%s", *migrationDir)
// 	m, err := migrate.New(migrationPath, databaseURL())
// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, migrationError(args[0], err))
// 		os.Exit(1)
// 	}

// 	// Run the specified command.
// 	err = runMigration(m, args[0], args[1:])
// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, err)
// 		os.Exit(1)
// 	}

// 	fmt.Println("Success!")
// }

// // migrator describes the interface of the migrate package to
// // faciliate testing runMigration.
// type migrator interface {
// 	Down() error
// 	Drop() error
// 	Force(version int) error
// 	Migrate(n uint) error
// 	Steps(n int) error
// 	Up() error
// 	Version() (uint, bool, error)
// }

// func runMigration(m migrator, command string, args []string) error {
// 	fmt.Printf("Running migrate %s %s...\n", command, strings.Join(args, " "))

// 	// Ensure numeric arguments can be parsed.
// 	var n int
// 	var err error
// 	if len(args) > 0 {
// 		n, err = strconv.Atoi(args[0])
// 		if err != nil {
// 			return migrationError(command, err)
// 		}
// 	}

// 	switch command {
// 	case "down":
// 		err = m.Down()
// 	case "drop":
// 		err = m.Drop()
// 	case "force":
// 		if len(args) == 0 {
// 			return migrationError(command, errors.New("a migration version number is required"))
// 		}
// 		err = m.Force(n)
// 	case "steps":
// 		if len(args) == 0 {
// 			return migrationError(command, errors.New("the number of steps to migrate is required"))
// 		}
// 		err = m.Steps(n)
// 	case "toVersion":
// 		if len(args) == 0 {
// 			return migrationError(command, errors.New("a migration version number is required"))
// 		}
// 		err = m.Migrate(uint(n))
// 	case "up":
// 		err = m.Up()
// 	case "version":
// 		version, dirty, err := m.Version()
// 		if err != nil {
// 			return migrationError(command, err)
// 		}
// 		fmt.Fprintf(out, "\tVersion: %d\n\tDirty: %t\n", version, dirty)
// 	default:
// 		return migrationError(command, errors.New("unknown command"))
// 	}

// 	if err != nil {
// 		return migrationError(command, err)
// 	}
// 	return nil
// }

// func migrationError(command string, err error) error {
// 	return fmt.Errorf("migrate %s: %v", command, err)
// }

// func databaseURL() string {
// 	return fmt.Sprintf(
// 		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
// 		envVars.Get("FB05_DB_USER"),
// 		envVars.Get("FB05_DB_PASSWORD"),
// 		envVars.Get("FB05_DB_HOST"),
// 		envVars.Get("FB05_DB_PORT"),
// 		envVars.Get("FB05_DB_NAME"),
// 	)
// }
