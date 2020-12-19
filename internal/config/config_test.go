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
	cfg := config.New()

	if err := cfg.Load("/dev/null"); err == nil {
		t.Error("config should require valid JSON")
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	file := path.Join(cwd, "..", "..", "testdata", "kaklint.json")

	if err := cfg.Load(file); err != nil {
		t.Error(err)
	}

	var missingConfigErr config.ErrMissingConfiguration
	if _, _, _, err := cfg.Get("unknown"); err == nil || !errors.As(err, &missingConfigErr) {
		t.Error("config should return ErrMissingConfiguration")
	}

	cmd, _, _, err := cfg.Get("php")
	if err != nil {
		t.Error(err)
	}

	want := []string{"php", "-l"}
	if !reflect.DeepEqual(want, cmd) {
		t.Errorf("want: %v, got: %v", want, cmd)
	}
}
