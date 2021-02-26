package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

// Default is the default configuration.
var Default *Config

// Initialize a default configuration without checking for errors.
func init() {
	Default = New()
	Default.Load(path.Join(os.Getenv("HOME"), ".config", "kaklint.json"))
	Default.Load(".kaklint.json") // Allow project-level overrides.
}

// Config holds the configuration. It maps a file type to a configuration
// entry.
type Config struct{ linters map[string]Linter }

// New returns a new configuration object.
func New() *Config { return &Config{make(map[string]Linter)} }

// Load decodes JSON configuration from the given file.
func (cfg *Config) Load(filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(content, &cfg.linters); err != nil {
		return err
	}
	return nil
}

// Get returns the configuration for a given linter.
func (cfg *Config) Get(linter string) ([]string, []string, bool, error) {
	if entry, ok := cfg.linters[linter]; ok {
		return entry.Cmd, entry.Efm, entry.Pkg, nil
	}
	return nil, nil, false, ErrMissingConfiguration{linter}
}

// ErrMissingConfiguration is returned upon missing configuration for filetype.
type ErrMissingConfiguration struct{ linter string }

// Error implements standard error for ErrMissingConfiguration.
func (err ErrMissingConfiguration) Error() string {
	return fmt.Sprintf("missing configuration for linter: %s", err.linter)
}

// Linter is a configuration entry.
type Linter struct {
	Cmd []string `json:"cmd"` // Which command to run.
	Efm []string `json:"efm"` // Which shape do error messages have.
	Pkg bool     `json:"pkg"` // Should the command run at package level?
}
