package env

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	wantConfig := &LoadConfig{"a", "b", "c", "d"}
	gotConfig := NewConfig("a", "b", "c", "d")
	if !reflect.DeepEqual(wantConfig, gotConfig) {
		t.Errorf("want *LoadConfig %+v, got %+v", wantConfig, gotConfig)
	}
}

func TestLoad(t *testing.T) {
	envKey := "CURRENT_ENV"
	currentEnv := "test"

	testCases := []struct {
		name, fileType string
	}{
		{"jsonenv", "json"},
		{"yamlenv", "yaml"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("loads env vars from %s", tc.fileType), func(t *testing.T) {
			config := NewConfig(tc.name, tc.fileType, fixturesDir(), currentEnv)
			err := Load(config)
			if err != nil {
				t.Fatal(err)
			}

			if gotVar := env.vars.Get(envKey); gotVar != currentEnv { // not thread-safe
				t.Fatalf("want %s=%s, got %s=%s", envKey, currentEnv, envKey, gotVar)
			}

			prodOnlyKey := "PROD_ONLY"
			if prodVar := env.vars.Get(prodOnlyKey); prodVar != nil { // not thread-safe
				t.Fatalf("want env var %s to be inaccessible, got %s=%s",
					prodOnlyKey, prodOnlyKey, prodVar)
			}
		})
	}
}

func TestGet(t *testing.T) {
	t.Run("calling Get before Load", func(t *testing.T) {
		got := Get("CURRENT_ENV")
		if got != nil {
			t.Errorf("want Get to return nil when called before Load, got %v", got)
		}
	})

	config := NewConfig("yamlenv", "yaml", fixturesDir(), "test")
	err := Load(config)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		key  string
		want interface{}
	}{
		{"CURRENT_ENV", "test"},
		{"NOT_PRESENT", nil},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Get(%q)", tc.key), func(t *testing.T) {
			got := Get(tc.key)
			if got != tc.want {
				t.Errorf("want %s=%v, got %s=%v", tc.key, tc.want, tc.key, got)
			}
		})
	}
}

func TestReset(t *testing.T) {
	config := NewConfig("yamlenv", "yaml", fixturesDir(), "test")
	err := Load(config)
	if err != nil {
		t.Fatal(err)
	}

	if env.vars == nil {
		t.Fatal("failed to load env vars")
	}
	Reset()
	if env.vars != nil {
		t.Fatalf("want env.vars to be nil, got %v", env.vars)
	}
}

func fixturesDir() string {
	dir, _ := os.Getwd()
	return filepath.Join(dir, "fixtures")
}
