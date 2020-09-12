package main

import (
	"angusgmorrison/fb05/pkg/env"
	"bytes"
	"errors"
	"fmt"
	"testing"
)

type mockMigrator struct {
	lastFuncCalled string
}

func (m *mockMigrator) Down() error {
	m.lastFuncCalled = "Down"
	return nil
}

func (m *mockMigrator) Drop() error {
	m.lastFuncCalled = "Drop"
	return nil
}

func (m *mockMigrator) Migrate(_ uint) error {
	m.lastFuncCalled = "Migrate"
	return nil
}

func (m *mockMigrator) Force(_ int) error {
	m.lastFuncCalled = "Force"
	return nil
}

func (m *mockMigrator) Steps(_ int) error {
	m.lastFuncCalled = "Steps"
	return nil
}

func (m *mockMigrator) Up() error {
	m.lastFuncCalled = "Up"
	return nil
}

func (m *mockMigrator) Version() (uint, bool, error) {
	m.lastFuncCalled = "Version"
	return 1, false, nil
}

func TestRunMigration(t *testing.T) {
	testCases := []struct {
		command        string
		args           []string
		wantLastCalled string
		wantErr        bool
		wantOutput     string
	}{
		{"down", nil, "Down", false, ""},
		{"drop", nil, "Drop", false, ""},
		{"force", []string{"1"}, "Force", false, ""},
		{"force", []string{"@"}, "", true, ""},
		{"steps", []string{"1"}, "Steps", false, ""},
		{"toVersion", []string{"1"}, "Migrate", false, ""},
		{"up", nil, "Up", false, ""},
		{"version", nil, "Version", false, "\tVersion: 1\n\tDirty: false\n"},
		{"unknown", nil, "", true, ""},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("runMigration(m, %q, %v)", tc.command, tc.args), func(t *testing.T) {
			if tc.wantOutput != "" {
				oldOut := out
				out = new(bytes.Buffer)
				defer func() {
					out = oldOut
				}()
			}

			var m mockMigrator
			err := runMigration(&m, tc.command, tc.args)
			if gotErr := err != nil; gotErr != tc.wantErr {
				t.Fatalf("wantErr == %t, but got %v", tc.wantErr, err)
			}
			if m.lastFuncCalled != tc.wantLastCalled {
				t.Fatalf("want lastFuncCalled to be %q, got %q",
					tc.wantLastCalled, m.lastFuncCalled)
			}
			if tc.wantOutput != "" {
				if gotOutput := out.(*bytes.Buffer).String(); gotOutput != tc.wantOutput {
					t.Errorf("want output %q, got %q", tc.wantOutput, gotOutput)
				}
			}
		})

	}
}

func TestMigrationError(t *testing.T) {
	command := "test"
	err := errors.New("error message")
	wantErrorMsg := fmt.Sprintf("migrate %s: %v", command, err)
	gotErrorMsg := migrationError(command, err).Error()
	if gotErrorMsg != wantErrorMsg {
		t.Errorf("want error message %q, got %q", wantErrorMsg, gotErrorMsg)
	}
}

func TestDatabaseURL(t *testing.T) {
	envConfig := env.NewConfig("environment", "yaml", "fixtures", "test")
	err := env.Load(envConfig)
	if err != nil {
		t.Fatalf("failed to load environment: %v", err)
	}
	defer func() { env.Reset() }()

	wantDBURL := "postgres://test_user:password@test_host:1234/test_db?sslmode=disable"
	if gotDBURL := databaseURL(); gotDBURL != wantDBURL {
		t.Errorf("want DB URL %q, got %q", wantDBURL, gotDBURL)
	}
}
