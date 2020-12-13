package config_test

import (
	"errors"
	"kaklint/internal/config"
	"os"
	"path"
	"reflect"
	"testing"
)

func TestConfig(t *testing.T) {
	if _, err := config.New("/dev/null"); err == nil {
		t.Error("config should require valid JSON")
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	cfg, err := config.New(path.Join(cwd, "..", "..", "testdata", "config.json"))
	if err != nil {
		t.Error(err)
	}

	var missingConfigErr config.ErrMissingConfiguration
	if _, _, err := cfg.Get("unknown"); err == nil || !errors.As(err, &missingConfigErr) {
		t.Error("config should return ErrMissingConfiguration")
	}

	cmd, _, err := cfg.Get("php")
	if err != nil {
		t.Error(err)
	}

	want := []string{"php", "-l"}
	if !reflect.DeepEqual(want, cmd) {
		t.Errorf("want: %v, got: %v", want, cmd)
	}
}
