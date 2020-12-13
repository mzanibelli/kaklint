package kaklint_test

import (
	"bytes"
	"kaklint"
	"testing"
)

func TestKakLint(t *testing.T) {
	kl := kaklint.New(stubConfig{}, stubLinter{}, bytes.NewBuffer(nil))
	if err := kl.Lint("php", "testdata/badfile.php"); err != nil {
		t.Error(err)
	}
}

type stubConfig struct{}

func (cfg stubConfig) Get(filetype string) (cmd, efm []string, err error) {
	return []string{"php", "-l"}, []string{"%m in %f on line %l", "%-G%.%#"}, nil
}

const phpError = `
PHP Parse error:  syntax error, unexpected 'die' (T_EXIT), expecting ';' or ',' in testdata/badfile.php on line 5
Errors parsing testdata/badfile.php
`

type stubLinter struct{}

func (l stubLinter) Run(args []string) ([]byte, error) {
	return []byte(phpError), nil
}
