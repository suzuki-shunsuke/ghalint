package lockfile_test

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghalint/pkg/lockfile"
)

func TestRead(t *testing.T) {
	t.Parallel()
	afs := afero.NewMemMapFs()
	lf, err := lockfile.Read(afs)
	if err != nil {
		t.Fatal(err)
	}
	if lf != nil {
		t.Fatal("lockfile should be nil if it doesn't exist")
	}
	if err := afero.WriteFile(afs, lockfile.Path, []byte(`version: 'v0.0.3'
workflows:
  '.github/workflows/ci.yml':
    - 'actions/checkout@v6'
dependencies:
  'actions/checkout@v6':
    ref: 'v6'
    commit: '1af3b93b6815bc44a9784bd300feb67ff0d1eeb3'
    owner_id: 44036562
    repo_id: 197814629
`), 0o644); err != nil {
		t.Fatal(err)
	}
	lf, err = lockfile.Read(afs)
	if err != nil {
		t.Fatal(err)
	}
	if !lf.Locked(".github/workflows/ci.yml", "actions/checkout@v6") {
		t.Fatal("actions/checkout@v6 should be locked")
	}
	if lf.Locked(".github/workflows/ci.yml", "actions/checkout@v5") {
		t.Fatal("actions/checkout@v5 shouldn't be locked")
	}
	if lf.Locked(".github/workflows/test.yml", "actions/checkout@v6") {
		t.Fatal("actions/checkout@v6 shouldn't be locked in test.yml")
	}
}
