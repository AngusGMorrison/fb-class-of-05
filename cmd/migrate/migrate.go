// migrate wraps golang-migrate, allowing environment variables to be
// loaded from a config file before running migrations.
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/viper"
)

// env stores environment variables loaded by viper. The variables
// loaded are specific to a single environment, e.g. development.
var env *viper.Viper

func main() {
	// Flags to control environment variable loading.
	configName := flag.String(
		"configName",
		"environment",
		"The name of the config file (without extension) containing env vars required for the migration.",
	)
	configPath := flag.String(
		"configPath",
		".",
		"The path to the config file containing env vars required for the migration.",
	)
	flag.Parse()

	// Ensure a command has been provided.
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("migrate must be used with a migration command")
		fmt.Println("Usage: migrate down | drop | up | version | force number | step number | toVersion number [-configName string] [-configPath string]")
	}

	// Load environment variables
	var err error
	env, err = loadEnvVars(*configName, *configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "migrate: %v\n", err)
		os.Exit(1)
	}

	// Run the specified command
	err = runMigration(args[0], args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	} else {
		fmt.Println("Success!")
	}
}

func loadEnvVars(configName, configPath string) (*viper.Viper, error) {
	targetEnv := os.Getenv("FB05_ENV")
	if targetEnv == "" {
		return nil, fmt.Errorf("loadEnvVars: target environment unknown, FB_ENV is blank")
	}

	fmt.Printf("Loading environment %q...\n", targetEnv)
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("loadEnvVars: %v", err)
	}
	// Return only the variables for the target environment.
	return viper.Sub(targetEnv), nil
}

func runMigration(command string, args []string) error {
	fmt.Printf("Running migrate %s %s...\n", command, strings.Join(args, " "))
	m, err := migrate.New("file://db/migrations", databaseURL())
	if err != nil {
		return migrationError(command, err)
	}

	var n int
	if len(args) > 0 {
		n, err = strconv.Atoi(args[0])
		if err != nil {
			return migrationError(command, err)
		}
	}

	switch command {
	case "down":
		err = m.Down()
	case "drop":
		err = m.Drop()
	case "force":
		if len(args) == 0 {
			return migrationError(command, errors.New("a migration version number is required"))
		}
		err = m.Force(n)
	case "steps":
		if len(args) == 0 {
			return migrationError(command, errors.New("the number of steps to migrate is required"))
		}
		err = m.Steps(n)
	case "toVersion":
		if len(args) == 0 {
			return migrationError(command, errors.New("a migration version number is required"))
		}
		err = m.Migrate(uint(n))
	case "up":
		err = m.Up()
	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			return migrationError(command, err)
		}
		fmt.Printf("\tVersion: %d\n\tDirty: %t\n", version, dirty)
	default:
		return migrationError(command, errors.New("unknown command"))
	}

	if err != nil {
		return migrationError(command, err)
	}
	return nil
}

func migrationError(command string, err error) error {
	return fmt.Errorf("migrate %s: %v", command, err)
}

func databaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		env.Get("FB05_DB_USER"),
		env.Get("FB05_DB_PASSWORD"),
		env.Get("FB05_DB_HOST"),
		env.Get("FB05_DB_PORT"),
		env.Get("FB05_DB_NAME"),
	)
}
