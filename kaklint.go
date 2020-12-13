package kaklint

import (
	"fmt"
	"io"
	"kaklint/internal/config"
	"kaklint/internal/errfmt"
	"kaklint/internal/linter"
	"os"
)

// Default is the default instance.
var Default *KakLint

func init() {
	Default = New(config.Default, linter.Default, os.Stdout)
}

// Linter is a runnable external program that produces output.
type Linter interface {
	Run(args []string) ([]byte, error)
}

// Config reads configuration for a given file type.
type Config interface {
	Get(filetype string) (cmd, efm []string, global bool, err error)
}

// KakLint lints files and reshapes error messages.
type KakLint struct {
	config Config
	linter Linter
	output io.Writer
}

// New returns a new instance.
func New(config Config, linter Linter, output io.Writer) *KakLint {
	return &KakLint{config, linter, output}
}

// Lint runs the linter and formats results into Kakoune's format.
func (kl KakLint) Lint(filetype, target string) error {
	cmd, efm, global, err := kl.config.Get(filetype)
	if err != nil {
		return err
	}

	// If global is set, run the linter without arguments.
	// TODO: move to git top-level if any?
	if !global {
		cmd = append(cmd, target)
	}

	// Store error for later use since a linter failure generally
	// means something is there for us to parse.
	output, lintErr := kl.linter.Run(cmd)

	messages, err := errfmt.Parse(output, efm)
	if err != nil {
		return err
	}

	// If there was an error executing the linter but no messages
	// were parsed, this means we made a configuration mistake or
	// something unexpected happened.
	if len(messages) == 0 && lintErr != nil {
		return lintErr
	}

	for _, mess := range messages {
		if _, err := fmt.Fprintln(kl.output, mess); err != nil {
			return err
		}
	}

	return nil
}
