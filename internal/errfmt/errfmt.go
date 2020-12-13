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
	return Entry{e}
}

const (
	// Info is a kind of entry that represents a informational message.
	Info string = "info"
	// Warning is a kind of entry that represents a warning message.
	Warning string = "warning"
	// Error is a kind of entry that represents an error message.
	Error string = "error"
)

// Entry is a wrapper to errorformat.Entry.
type Entry struct{ *errorformat.Entry }

// String implements fmt.Stringer.
func (e Entry) String() string {
	return fmt.Sprintf("%s:%d:%d: %s: %s", e.Filename, e.Lnum, e.Col, e.Kind(), e.Text)
}

// Kind returns the error kind. Kakoune only support Warning or Error.
// See: https://github.com/reviewdog/errorformat/blob/55531c7dabdfad07a928152b1c6eb9dcd2eb3bdb/errorformat.go#L138
func (e Entry) Kind() string {
	switch kind := e.Types(); {
	case strings.Index(kind, Info) == 0:
		return Warning
	case strings.Index(kind, Warning) == 0:
		return Warning
	default:
		return Error
	}
}
