package deploy

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// RailwayCLIConfig holds values read from the Railway CLI's config file (~/.railway/config.json).
type RailwayCLIConfig struct {
	Token       string
	Project     string
	Service     string
	Environment string
}

type railwayConfigFile struct {
	Projects map[string]railwayProjectEntry `json:"projects"`
	User     struct {
		Token string `json:"token"`
	} `json:"user"`
}

type railwayProjectEntry struct {
	Project     string `json:"project"`
	Service     string `json:"service"`
	Environment string `json:"environment"`
}

// ReadRailwayCLIConfig reads the Railway CLI configuration for the given working directory.
// It looks for ~/.railway/config.json and extracts the token and any project-specific settings.
func ReadRailwayCLIConfig(workdir string) (*RailwayCLIConfig, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return readRailwayCLIConfigFrom(filepath.Join(home, ".railway", "config.json"), workdir)
}

func readRailwayCLIConfigFrom(configPath, workdir string) (*RailwayCLIConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg railwayConfigFile
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	result := &RailwayCLIConfig{
		Token: cfg.User.Token,
	}

	// Look up project entry by workdir
	if entry, ok := cfg.Projects[workdir]; ok {
		result.Project = entry.Project
		result.Service = entry.Service
		result.Environment = entry.Environment
	}

	return result, nil
}
