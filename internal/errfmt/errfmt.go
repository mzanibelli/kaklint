package errfmt

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/reviewdog/errorformat"
)

// Parse uses the third-party library to parse a given input using a
// given efm. It returns entries that can be understood by Kakoune.
func Parse(input []byte, shape []string) ([]Entry, error) {
	res := make([]Entry, 0)

	efm, err := errorformat.NewErrorformat(shape)
	if err != nil {
		return res, err
	}

	scanner := efm.NewScanner(bytes.NewBuffer(input))
	for scanner.Scan() {
		res = append(res, newEntry(scanner.Entry()))
	}

	return res, nil
}

// Kakoune linter bugs if lines or columns are not greater than 0.
func newEntry(e *errorformat.Entry) Entry {
	if e.Lnum == 0 {
		e.Lnum = 1
	}
	if e.Col == 0 {
		e.Col = 1
	}
	return Entry{flag(e), mess(e)}
}

const (
	// Info is a kind of entry that represents a informational message.
	Info string = "info"
	// Note is a kind of entry that represents a simple note.
	Note string = "note"
	// Warning is a kind of entry that represents a warning message.
	Warning string = "warning"
	// Error is a kind of entry that represents an error message.
	Error string = "error"
)

// Entry is a wrapper to errorformat.Entry.
type Entry struct {
	Flag string
	Mess string
}

// See: https://github.com/reviewdog/errorformat/blob/55531c7dabdfad07a928152b1c6eb9dcd2eb3bdb/errorformat.go#L138
func flag(e *errorformat.Entry) string {
	var icon string

	switch kind := e.Types(); {

	// Third-party lib doesn't seem to support notes. If they do
	// one day, some tests should break (shellcheck-note).
	case strings.Index(kind, Note) == 0:
		icon = "?"

	case strings.Index(kind, Info) == 0:
		icon = "?"
	case strings.Index(kind, Warning) == 0:
		icon = "!"
	default:
		icon = "x"
	}

	return spec(e.Lnum, icon)
}

// Handle pipe in flag message: there is currently no other way than
// escaping pipes with backslashes, but this displays a literal backslash
// in the resulting info message.
func mess(e *errorformat.Entry) string {
	return spec(e.Lnum, strings.Replace(e.Text, `|`, `\|`, -1))
}

// Quote string with "...". Any "" must be doubled.
// See: :doc commands-parsing
func spec(line int, text string) string {
	return fmt.Sprintf(`"%d|%s"`, line, strings.Replace(text, `"`, `""`, -1))
}
