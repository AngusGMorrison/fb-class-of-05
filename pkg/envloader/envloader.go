// Package envloader facilitates the loading of environment variables
// from configuration files.
package envloader

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// LoadConfig details the configuration file to be loaded.
type LoadConfig struct {
	configName, configType, configPath, targetEnv string
}

// NewConfig is a convenience function for constructing LoadConfigs.
func NewConfig(configName, configType, configPath, targetEnv string) *LoadConfig {
	return &LoadConfig{configName, configType, configPath, targetEnv}
}

// Load looks up the current environment using targetEnvKey and loads
// the corresponding environment variables from the configuration
// file specified by l, returning them as a *viper.Viper.
func Load(lc *LoadConfig) (*viper.Viper, error) {
	var err error
	// Configure Viper and read in the configuration file.
	fmt.Printf("Loading environment %q...\n", lc.targetEnv)
	viper.SetConfigName(lc.configName)
	viper.SetConfigType(lc.configType)
	viper.AddConfigPath(lc.configPath)
	if err = viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("loading env: %v", err)
	}

	// Return only the variables for the target environment.
	sub := viper.Sub(lc.targetEnv)
	if sub == nil {
		err = ErrNoEnvKeys
	}
	return sub, err
}

// ErrNoEnvKeys signifies that no variables were found within the
// environment specified by the user.
var ErrNoEnvKeys error = errors.New("environment loaded but it contains no env vars")
