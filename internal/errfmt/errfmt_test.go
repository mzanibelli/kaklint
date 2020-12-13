package errfmt_test

import (
	"fmt"
	"kaklint/internal/errfmt"
	"testing"
)

const phpError = `
PHP Parse error:  syntax error, unexpected 'die' (T_EXIT), expecting ';' or ',' in testdata/badfile.php on line 5
Errors parsing testdata/badfile.php
`

func TestErrfmt(t *testing.T) {
	phpEfm := []string{"%m in %f on line %l", "%-G%.%#"}

	entries, err := errfmt.Parse([]byte(phpError), phpEfm)
	if err != nil {
		t.Error(err)
	}

	if len(entries) != 1 {
		t.Errorf("want: 1 entry, got: %d", len(entries))
	}

	want := "testdata/badfile.php:5:1: error: PHP Parse error:  syntax error, unexpected 'die' (T_EXIT), expecting ';' or ','"
	got := fmt.Sprint(entries[0])
	if want != got {
		t.Errorf("want: %s, got: %s", want, got)
	}
}
