package linter_test

import (
	"kaklint/internal/linter"
	"testing"
)

func TestDefaultLinter(t *testing.T) {
	_, err := linter.Default.Run([]string{"/bin/false"})
	if err == nil {
		t.Error("linter should fail if external program fails")
	}

	args := []string{"/bin/sh", "-c", "echo out; echo err >&2 2>&1"}
	output, err := linter.Default.Run(args)
	if err != nil {
		t.Error(err)
	}

	want := "out\nerr\n"
	got := string(output)
	if want != got {
		t.Errorf("want: %s, got: %s", want, got)
	}
}
