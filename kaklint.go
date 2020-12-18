package kaklint

import (
	"io"
	"kaklint/internal/config"
	"kaklint/internal/errfmt"
	"kaklint/internal/linter"
	"os"
	"text/template"
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

// Config reads configuration for a given linter.
type Config interface {
	Get(linter string) (cmd, efm []string, global bool, err error)
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

const tpl = `{{if len . -}}
set-option buffer lint_flags %val{timestamp}{{range .}} {{.Flag}}{{end}}
set-option buffer lint_messages %val{timestamp}{{range .}} {{.Mess}}{{end}}
lint-show-diagnostics
{{else -}}
lint-hide-diagnostics
{{end}}`

// Lint runs the linter and formats results into Kakoune's format.
func (kl KakLint) Lint(linter, target string) error {
	cmd, efm, global, err := kl.config.Get(linter)
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

	// Parsed errors must be translated into Kakoune instructions.
	templ, err := template.New("").Parse(tpl)
	if err != nil {
		return err
	}

	return templ.Execute(kl.output, &messages)
}
