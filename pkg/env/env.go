// Package env facilitates the loading of environment variables
// from configuration files.
package env

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type LoadConfig struct {
	configName, configType, configPath, targetEnv string
}

func NewConfig(configName, configType, configPath, targetEnv string) *LoadConfig {
	return &LoadConfig{configName, configType, configPath, targetEnv}
}

type store struct {
	sync.RWMutex
	vars *viper.Viper
}

var env *store = new(store)

// Load looks up the current environment using targetEnvKey and loads
// the corresponding environment variables from the configuration
// file specified by l, returning them as a *viper.Viper.
func Load(lc *LoadConfig) error {
	// Configure Viper and read in the configuration file.
	viper.SetConfigName(lc.configName)
	viper.SetConfigType(lc.configType)
	viper.AddConfigPath(lc.configPath)
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("loading env: %v", err)
	}

	// Return only the variables for the target environment.
	env.Lock()
	defer env.Unlock()
	env.vars = viper.Sub(lc.targetEnv)
	if env.vars == nil {
		return ErrNoEnvKeys
	}

	return nil
}

// ErrNoEnvKeys signifies that no variables were found within the
// environment specified by the user.
var ErrNoEnvKeys error = errors.New("environment loaded but it contains no env vars")

func Get(key string) interface{} {
	env.RLock()
	defer env.RUnlock()
	if env.vars == nil {
		return nil
	}
	return env.vars.Get(key)
}

func Reset() {
	env.Lock()
	env.vars = nil
	env.Unlock()
}
