package kaklint_test

import (
	"bytes"
	"io/ioutil"
	"kaklint"
	"kaklint/internal/config"
	"os"
	"path"
	"strings"
	"testing"
)

func TestKakLint(t *testing.T) {
	if os.Getenv("KAKLINT_ENV") != "docker" {
		t.Skip("this test is meant to run inside docker")
	}

	out := bytes.NewBuffer(nil)

	cfg := config.New()
	if err := cfg.Load(path.Join("testdata", "kaklint.json")); err != nil {
		t.Fatal(err)
	}

	kl := kaklint.New(cfg, out)

	tests, err := snapshots()
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range tests {
		t.Run(test.name, check(test, cfg, kl, out))
	}
}

func check(
	test snapshot,
	cfg *config.Config,
	kl *kaklint.KakLint,
	out *bytes.Buffer,
) func(t *testing.T) {
	return func(t *testing.T) {
		defer out.Reset()

		err := kl.Lint(test.linter, test.input)
		if err != nil {
			t.Error(err)
		}

		if string(test.want) != out.String() {
			t.Errorf("\n\n%s\n\n%s\n\n", test.want, out)
		}
	}
}

type snapshot struct {
	name   string
	linter string
	input  string
	want   []byte
}

func snapshots() ([]snapshot, error) {
	root := path.Join("testdata", "snapshots")

	dirs, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, err
	}

	res := make([]snapshot, len(dirs), len(dirs))

	for i, fi := range dirs {
		snap, err := parse(root, fi)
		if err != nil {
			return nil, err
		}
		res[i] = snap
	}

	return res, nil
}

func parse(root string, fi os.FileInfo) (snapshot, error) {
	input := path.Join(root, fi.Name(), "input")

	ftPath := path.Join(root, fi.Name(), "linter")

	linter, err := ioutil.ReadFile(ftPath)
	if err != nil {
		return snapshot{}, err
	}

	wantPath := path.Join(root, fi.Name(), "want")

	want, err := ioutil.ReadFile(wantPath)
	if err != nil {
		return snapshot{}, err
	}

	ft := strings.TrimSpace(string(linter))

	return snapshot{fi.Name(), ft, input, want}, nil
}
