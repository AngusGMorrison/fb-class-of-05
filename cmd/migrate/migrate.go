// migrate wraps golang-migrate, allowing environment variables to be
// loaded from a config file before running migrations.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/viper"
)

// env stores environment variables loaded by viper. The variables
// loaded are specific to a single environment, e.g. development.
var env *viper.Viper

func main() {
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
	steps := flag.Int("steps", 0, "Optional number of steps to migrate (default: all).")
	version := flag.Int("version", 0, "Optional version to migrate to.")
	flag.Parse()

	// Ensure a command has been provided.
	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("migrate must be used with a migration command")
		fmt.Println("Usage: migrate <command> <flags>")
		fmt.Println("Valid commands are down, drop, force, toVersion, up, and version.")
	}

	// Load environment variables
	var err error
	env, err = loadEnvVars(*configName, *configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "migrate: %v\n", err)
		os.Exit(1)
	}

	err = runMigration(args[0], *steps, *version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "migrate: %v\n", err)
	}

	fmt.Println("Success!")
}

func loadEnvVars(configName, configPath string) (*viper.Viper, error) {
	targetEnv := os.Getenv("FB_ENV")
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

func runMigration(command string, steps, version int) error {
	fmt.Printf("Running migrate %s...\n", command)
	m, err := migrate.New("file://db/migrations", databaseURL())
	if err != nil {
		return fmt.Errorf("runMigration: %v", err)
	}

	switch command {
	case "down":
		if steps > 0 {
			return m.Steps(-steps)
		} else {
			return m.Down()
		}
	case "drop":
		return m.Drop()
	case "force":
		return m.Force(version)
	case "toVersion":
		return m.Migrate(uint(version))
	case "up":
		if steps > 0 {
			return m.Steps(steps)
		} else {
			return m.Up()
		}
	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			return err
		}
		fmt.Printf("version: %d\ndirty: %t\n", version, dirty)
	default:
		return fmt.Errorf("runMigration: unknown command %s", command)
	}
	return nil
}

func databaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		env.Get("FB_DB_USERNAME"),
		env.Get("FB_DB_PASSWORD"),
		env.Get("FB_DB_HOST"),
		env.Get("FB_DB_PORT"),
		env.Get("FB_DB_NAME"),
	)
}
