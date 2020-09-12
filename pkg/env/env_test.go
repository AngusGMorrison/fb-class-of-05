package env

import (
	"fmt"
	"os"
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
	dir, _ := os.Getwd()
	fmt.Println(dir)
	testCases := []struct {
		name, fileType string
	}{
		{"jsonenv", "json"},
		{"yamlenv", "yaml"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("loads env vars from %s", tc.fileType), func(t *testing.T) {
			config := NewConfig(tc.name, tc.fileType, dir+"/fixtures", currentEnv)
			envVars, err := Load(config)
			if err != nil {
				t.Fatal(err)
			}
			if gotVar := envVars.Get(envKey); gotVar != currentEnv {
				t.Fatalf("want %s=%s, got %s", envKey, currentEnv, gotVar)
			}

			prodOnlyKey := "PROD_ONLY"
			if prodVar := envVars.Get(prodOnlyKey); prodVar != nil {
				t.Fatalf("want env var %s to be inaccessible, got %s=%s",
					prodOnlyKey, prodOnlyKey, prodVar)
			}
		})
	}
}
