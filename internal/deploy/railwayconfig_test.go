package deploy

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadRailwayCLIConfig(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	workdir := "/Users/test/myapp"

	config := `{
		"projects": {
			"/Users/test/myapp": {
				"project": "proj-123",
				"service": "svc-456",
				"environment": "env-789"
			},
			"/Users/test/other": {
				"project": "proj-other",
				"service": "svc-other",
				"environment": "env-other"
			}
		},
		"user": {
			"token": "rw_test_token_abc"
		}
	}`

	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := readRailwayCLIConfigFrom(configPath, workdir)
	if err != nil {
		t.Fatalf("readRailwayCLIConfigFrom() error: %v", err)
	}

	if cfg.Token != "rw_test_token_abc" {
		t.Errorf("Token = %q, want %q", cfg.Token, "rw_test_token_abc")
	}
	if cfg.Project != "proj-123" {
		t.Errorf("Project = %q, want %q", cfg.Project, "proj-123")
	}
	if cfg.Service != "svc-456" {
		t.Errorf("Service = %q, want %q", cfg.Service, "svc-456")
	}
	if cfg.Environment != "env-789" {
		t.Errorf("Environment = %q, want %q", cfg.Environment, "env-789")
	}
}

func TestReadRailwayCLIConfigUnknownWorkdir(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")

	config := `{
		"projects": {
			"/Users/test/myapp": {
				"project": "proj-123"
			}
		},
		"user": {
			"token": "rw_token"
		}
	}`

	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := readRailwayCLIConfigFrom(configPath, "/some/other/path")
	if err != nil {
		t.Fatalf("readRailwayCLIConfigFrom() error: %v", err)
	}

	// Token should still be available
	if cfg.Token != "rw_token" {
		t.Errorf("Token = %q, want %q", cfg.Token, "rw_token")
	}
	// But project should be empty since workdir doesn't match
	if cfg.Project != "" {
		t.Errorf("Project = %q, want empty", cfg.Project)
	}
}

func TestReadRailwayCLIConfigMissing(t *testing.T) {
	_, err := readRailwayCLIConfigFrom("/nonexistent/config.json", "/some/path")
	if err == nil {
		t.Fatal("expected error for missing config file")
	}
}

func TestReadRailwayCLIConfigInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")

	if err := os.WriteFile(configPath, []byte("not json"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := readRailwayCLIConfigFrom(configPath, "/some/path")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
