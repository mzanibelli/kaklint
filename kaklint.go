package kaklint

import (
	"context"
	"io"
	"kaklint/internal/config"
	"kaklint/internal/errfmt"
	"os"
	"os/exec"
	"text/template"
	"time"
)

// Default is the default instance.
var Default *KakLint

func init() {
	Default = New(config.Default, os.Stdout)
}

// KakLint lints files and reshapes error messages.
type KakLint struct {
	config *config.Config
	output io.Writer
}

// New returns a new instance.
func New(config *config.Config, output io.Writer) *KakLint {
	return &KakLint{config, output}
}

const commands = `{{if len . -}}
set-option buffer lint_flags %val{timestamp}{{range .}} {{.Flag}}{{end}}
set-option buffer lint_messages %val{timestamp}{{range .}} {{.Mess}}{{end}}
lint-show-diagnostics
{{else -}}
lint-hide-diagnostics
{{end}}`

// Lint runs the linter and formats results into Kakoune's format.
func (kl KakLint) Lint(linter, target string) error {
	cmd, efm, pkg, err := kl.config.Get(linter)
	if err != nil {
		return err
	}

	// If pkg is set, run the linter without arguments.
	// TODO: move to git top-level if any?
	if !pkg {
		cmd = append(cmd, target)
	}

	// Do not allow the linter to run forever.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Store error for later use since a linter failure generally
	// means something is there for us to parse.
	//nolint:gosec
	output, lintErr := exec.CommandContext(ctx, cmd[0], cmd[1:]...).CombinedOutput()

	messages, err := errfmt.Parse(output, efm, target)
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
	tpl, err := template.New("").Parse(commands)
	if err != nil {
		return err
	}

	return tpl.Execute(kl.output, &messages)
}
