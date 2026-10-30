package lockfile_test

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghalint/pkg/lockfile"
)

const testLockfile = `version: 'v0.0.3'
workflows:
  '.github/workflows/ci.yml':
    - 'actions/cache@v4.2.0'
    - 'actions/checkout@v6'
dependencies:
  'actions/cache@v4.2.0':
    ref: 'v4.2.0'
    commit: 'sha1-1bd1e32a3bdc45362d1e726936510720a7c30a57'
    owner_id: 44036562
    repo_id: 207226787
  'actions/checkout@v6':
    ref: 'v6'
    commit: 'sha1-3d3c42e5aac5ba805825da76410c181273ba90b1'
    owner_id: 44036562
    repo_id: 197814629
`

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
	if err := afero.WriteFile(afs, lockfile.Path, []byte(testLockfile), 0o644); err != nil {
		t.Fatal(err)
	}
	lf, err = lockfile.Read(afs)
	if err != nil {
		t.Fatal(err)
	}
	if lf == nil {
		t.Fatal("lockfile should be returned")
	}
}

func TestParse(t *testing.T) {
	t.Parallel()
	data := []struct {
		name    string
		content string
		isErr   bool
	}{
		{
			name:    "valid",
			content: testLockfile,
		},
		{
			name: "commit without algorithm",
			content: `version: 'v0.0.3'
workflows:
  '.github/workflows/ci.yml':
    - 'actions/checkout@v6'
dependencies:
  'actions/checkout@v6':
    ref: 'v6'
    commit: '3d3c42e5aac5ba805825da76410c181273ba90b1'
    owner_id: 44036562
    repo_id: 197814629
`,
			isErr: true,
		},
		{
			name: "commit is missing",
			content: `version: 'v0.0.3'
workflows:
  '.github/workflows/ci.yml':
    - 'actions/checkout@v6'
dependencies:
  'actions/checkout@v6':
    ref: 'v6'
    owner_id: 44036562
    repo_id: 197814629
`,
			isErr: true,
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			if _, err := lockfile.Parse([]byte(d.content)); err != nil {
				if d.isErr {
					return
				}
				t.Fatal(err)
			}
			if d.isErr {
				t.Fatal("error must be returned")
			}
		})
	}
}

func TestLockfile_Locked(t *testing.T) {
	t.Parallel()
	lf, err := lockfile.Parse([]byte(testLockfile))
	if err != nil {
		t.Fatal(err)
	}
	data := []struct {
		name         string
		lockfile     *lockfile.Lockfile
		workflowPath string
		uses         string
		exp          bool
	}{
		{name: "locked", lockfile: lf, workflowPath: ".github/workflows/ci.yml", uses: "actions/checkout@v6", exp: true},
		{name: "sub-action path", lockfile: lf, workflowPath: ".github/workflows/ci.yml", uses: "actions/cache/save@v4.2.0", exp: true},
		{name: "owner and repo are case insensitive", lockfile: lf, workflowPath: ".github/workflows/ci.yml", uses: "Actions/Checkout@v6", exp: true},
		{name: "different ref", lockfile: lf, workflowPath: ".github/workflows/ci.yml", uses: "actions/checkout@v5"},
		{name: "different workflow", lockfile: lf, workflowPath: ".github/workflows/test.yml", uses: "actions/checkout@v6"},
		{name: "docker image", lockfile: lf, workflowPath: ".github/workflows/ci.yml", uses: "docker://actions/checkout:v6"},
		{name: "local action", lockfile: lf, workflowPath: ".github/workflows/ci.yml", uses: "./.github/actions/foo"},
		{name: "nil lockfile", workflowPath: ".github/workflows/ci.yml", uses: "actions/checkout@v6"},
		{name: "invalid lockfile", lockfile: lockfile.Invalid(), workflowPath: ".github/workflows/test.yml", uses: "actions/setup-go@v6", exp: true},
		{name: "docker image in invalid lockfile", lockfile: lockfile.Invalid(), workflowPath: ".github/workflows/test.yml", uses: "docker://rhysd/actionlint:latest"},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			if got := d.lockfile.Locked(d.workflowPath, d.uses); got != d.exp {
				t.Fatalf("got %v, wanted %v", got, d.exp)
			}
		})
	}
}
