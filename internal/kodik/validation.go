package kodik

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// OpenCodeConfig represents the structure of opencode.json
type OpenCodeConfig struct {
	Schema string                 `json:"$schema,omitempty"`
	Mode   map[string]interface{} `json:"mode"`
}

// RooModesConfig represents the structure of .roomodes
type RooModesConfig struct {
	Modes map[string]interface{} `yaml:"modes"`
}

// ValidateOpenCodeConfig validates JSON syntax and checks for required kodik modes
func ValidateOpenCodeConfig(configPath string) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil // Config doesn't exist, that's okay
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to read OpenCode config: %v", err))
	}

	var config OpenCodeConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return NewKodikError(ErrUser, fmt.Sprintf("Invalid JSON in OpenCode config: %v", err))
	}

	// Check for required kodik modes
	requiredModes := []string{"kodik-spec", "kodik-design", "kodik-tasks", "kodik-agent", "kodik-review-planner"}
	for _, mode := range requiredModes {
		if _, exists := config.Mode[mode]; !exists {
			Output.Warning("Missing kodik mode '%s' in OpenCode config", mode)
		}
	}

	Output.Success("OpenCode config validation passed (%d modes found)", len(config.Mode))
	return nil
}

// ValidateRooModesConfig validates YAML syntax and checks for kodik modes
func ValidateRooModesConfig(configPath string) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil // Config doesn't exist, that's okay
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return NewKodikError(ErrFilesystem, fmt.Sprintf("Failed to read .roomodes config: %v", err))
	}

	var config RooModesConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return NewKodikError(ErrUser, fmt.Sprintf("Invalid YAML in .roomodes config: %v", err))
	}

	// Check for kodik modes
	kodikModeCount := 0
	for modeName := range config.Modes {
		if strings.HasPrefix(modeName, "kodik-") {
			kodikModeCount++
		}
	}

	if kodikModeCount == 0 {
		Output.Warning("No kodik modes found in .roomodes config")
	} else {
		Output.Success(".roomodes config validation passed (%d kodik modes found)", kodikModeCount)
	}

	return nil
}
