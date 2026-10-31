package lockfile_test

import (
	"os"
	"path/filepath"
	"testing"

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

// writeLockfile writes the lockfile to the temporary directory, which is the current directory.
func writeLockfile(t *testing.T, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(lockfile.Path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(lockfile.Path, []byte(content), 0o644); err != nil { //nolint:gosec
		t.Fatal(err)
	}
}

func TestRead(t *testing.T) {
	t.Chdir(t.TempDir())
	workflows := []string{".github/workflows/ci.yml"}
	lf, err := lockfile.Read(workflows)
	if err != nil {
		t.Fatal(err)
	}
	if lf != nil {
		t.Fatal("lockfile should be nil if it doesn't exist")
	}
	writeLockfile(t, testLockfile)
	lf, err = lockfile.Read(workflows)
	if err != nil {
		t.Fatal(err)
	}
	if lf == nil {
		t.Fatal("lockfile should be returned")
	}
}

func TestRead_noWorkflow(t *testing.T) {
	t.Chdir(t.TempDir())
	// The lockfile is broken, but it isn't read because no workflow is given.
	writeLockfile(t, "version: 'v0.0.3'\nfoo: bar\n")
	lf, err := lockfile.Read(nil)
	if err != nil {
		t.Fatal(err)
	}
	if lf != nil {
		t.Fatal("lockfile should be nil if no workflow is given")
	}
}

func TestParse_workflowFilePaths(t *testing.T) {
	t.Parallel()
	// actions/setup-go@v6 lacks commit, but ci.yml doesn't use it.
	content := testLockfile + `  'actions/setup-go@v6':
    ref: 'v6'
    owner_id: 44036562
    repo_id: 192624594
`
	if _, err := lockfile.Parse([]byte(content)); err == nil {
		t.Fatal("the broken entry should be rejected if workflow file paths aren't given")
	}
	lf, err := lockfile.Parse([]byte(content), ".github/workflows/ci.yml")
	if err != nil {
		t.Fatal(err)
	}
	if !lf.Locked(".github/workflows/ci.yml", "actions/checkout@v6") {
		t.Fatal("actions/checkout@v6 should be locked")
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

func TestLockfile_Commit(t *testing.T) {
	t.Parallel()
	// actions/setup-go@v6 isn't referenced by any workflow, and actions/setup-node@v6 has a broken commit.
	// Only entries ci.yml references are validated, so the broken commit is accepted by Parse.
	content := testLockfile + `  'actions/setup-go@v6':
    ref: 'v6'
    commit: 'sha1-44694675825211faa026b3c33043df3e48a5fa00'
    owner_id: 44036562
    repo_id: 192624594
  'actions/setup-node@v6':
    ref: 'v6'
    commit: 'sha1-xyz'
    owner_id: 44036562
    repo_id: 223024751
`
	lf, err := lockfile.Parse([]byte(content), ".github/workflows/ci.yml")
	if err != nil {
		t.Fatal(err)
	}
	data := []struct {
		name     string
		lockfile *lockfile.Lockfile
		owner    string
		repo     string
		ref      string
		exp      string
	}{
		{name: "locked", lockfile: lf, owner: "actions", repo: "checkout", ref: "v6", exp: "3d3c42e5aac5ba805825da76410c181273ba90b1"},
		{name: "owner and repo are case insensitive", lockfile: lf, owner: "Actions", repo: "Checkout", ref: "v6", exp: "3d3c42e5aac5ba805825da76410c181273ba90b1"},
		{name: "no workflow references the entry", lockfile: lf, owner: "actions", repo: "setup-go", ref: "v6", exp: "44694675825211faa026b3c33043df3e48a5fa00"},
		{name: "broken commit", lockfile: lf, owner: "actions", repo: "setup-node", ref: "v6"},
		{name: "different ref", lockfile: lf, owner: "actions", repo: "checkout", ref: "v5"},
		{name: "ref is case sensitive", lockfile: lf, owner: "actions", repo: "checkout", ref: "V6"},
		{name: "not recorded", lockfile: lf, owner: "suzuki-shunsuke", repo: "ghalint", ref: "v1"},
		{name: "nil lockfile", owner: "actions", repo: "checkout", ref: "v6"},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			sha, ok := d.lockfile.Commit(d.owner, d.repo, d.ref)
			if ok != (d.exp != "") {
				t.Fatalf("got ok %v, wanted %v", ok, d.exp != "")
			}
			if sha != d.exp {
				t.Fatalf("got %q, wanted %q", sha, d.exp)
			}
		})
	}
}
