package errfmt

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/reviewdog/errorformat"
)

// Parse uses the third-party library to parse a given input using a
// given efm. It returns entries that can be understood by Kakoune.
func Parse(input []byte, shape []string, target string) ([]Entry, error) {
	res := make([]Entry, 0)

	efm, err := errorformat.NewErrorformat(shape)
	if err != nil {
		return res, err
	}

	scanner := efm.NewScanner(bytes.NewBuffer(input))
	for scanner.Scan() {
		ent := scanner.Entry()
		ok, err := sameFile(ent.Filename, target)
		if err != nil {
			return res, err
		}
		if ok {
			res = append(res, newEntry(ent))
		}
	}

	return res, nil
}

func sameFile(entry, target string) (bool, error) {
	a, err := filepath.Abs(entry)
	if err != nil {
		return false, err
	}

	b, err := filepath.Abs(target)
	if err != nil {
		return false, err
	}

	return a == b, nil
}

// Entry is a wrapper to errorformat.Entry.
type Entry struct {
	Flag string
	Mess string
}

// Kakoune linter bugs if lines or columns are not greater than 0.
// TODO: check if line and column are too big...
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
	info    string = "info"
	note    string = "note"
	warning string = "warning"
)

// Third-party lib doesn't seem to support notes. If they do one day,
// some tests should break (shellcheck-note).
// See: https://github.com/reviewdog/errorformat/blob/55531c7dabdfad07a928152b1c6eb9dcd2eb3bdb/errorformat.go#L138
func flag(e *errorformat.Entry) string {
	var icon string

	switch kind := e.Types(); {
	case strings.Index(kind, note) == 0:
		icon = "?"
	case strings.Index(kind, info) == 0:
		icon = "?"
	case strings.Index(kind, warning) == 0:
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
	txt := e.Text
	txt = strings.ReplaceAll(txt, `|`, `<pipe>`)
	txt = strings.ReplaceAll(txt, `%`, `<percent>`)
	return spec(e.Lnum, txt)
}

// Quote string with "...". Any "" must be doubled.
// See: :doc commands-parsing.
func spec(line int, text string) string {
	return fmt.Sprintf(`"%d|%s"`, line, strings.ReplaceAll(text, `"`, `""`))
}
