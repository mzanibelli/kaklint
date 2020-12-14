package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

// Default is the default configuration.
var Default *Config

// Initialize a default configuration without checking for errors.
func init() {
	confDir := os.Getenv("XDG_CONFIG_DIR")
	if confDir == "" {
		confDir = path.Join(os.Getenv("HOME"), ".config")
	}
	Default, _ = New(path.Join(confDir, "kaklint.json"))
}

// Config holds the configuration. It maps a file type to a configuration
// entry.
type Config struct {
	linters map[string]Linter
}

// New decodes JSON configuration from the given files and creates a new Config object.
func New(filename string) (*Config, error) {
	cfg := new(Config)
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return cfg, err
	}
	if err := json.Unmarshal(content, &cfg.linters); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// Get returns the configuration for a given linter.
func (cfg *Config) Get(linter string) ([]string, []string, bool, error) {
	if entry, ok := cfg.linters[linter]; ok {
		return entry.Cmd, entry.Efm, entry.Global, nil
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
	Cmd    []string `json:"cmd"`
	Efm    []string `json:"efm"`
	Global bool     `json:"global"`
}
