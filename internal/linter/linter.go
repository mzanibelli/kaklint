package linter

import "os/exec"

// Default is the default linter.
var Default Linter

// Linter is an executable command-line tool.
type Linter struct{}

// Run runs the external program.
func (l Linter) Run(args []string) ([]byte, error) {
	return exec.Command(args[0], args[1:]...).CombinedOutput()
}
