package schema

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/suzuki-shunsuke/ghalint/pkg/lockfile"
)

// testLockedSHA is the commit actions/checkout@v6 is locked to.
// It differs from the SHA countingGitHub returns, so the tests can tell which one is used.
const testLockedSHA = "0123456789abcdef0123456789abcdef01234567"

const testWorkflowLocked = `on: push
jobs:
  foo:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6
      - uses: actions/setup-go@v5
`

const testLockfileCheckout = `version: 'v0.0.3'
workflows:
  '.github/workflows/test.yaml':
    - 'actions/checkout@v6'
dependencies:
  'actions/checkout@v6':
    ref: 'v6'
    commit: 'sha1-` + testLockedSHA + `'
    owner_id: 44036562
    repo_id: 197814629
`

// testCompositeAction is an action file using actions/checkout@v6.
// The lockfile records actions/checkout@v6 for the workflow, but the commit is used for action files too.
const testCompositeAction = `name: test
runs:
  using: composite
  steps:
    - uses: actions/checkout@v6
`

// writeFiles creates the files in the current directory.
// The keys are slash-separated paths.
func writeFiles(t *testing.T, files map[string]string) {
	t.Helper()
	for p, content := range files {
		p = filepath.FromSlash(p)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil { //nolint:gosec
			t.Fatal(err)
		}
	}
}

func TestController_Run_lockfile(t *testing.T) {
	t.Chdir(t.TempDir())
	writeFiles(t, map[string]string{
		".github/workflows/test.yaml": testWorkflowLocked,
		lockfile.Path:                 testLockfileCheckout,
		"action.yaml":                 testCompositeAction,
	})
	gh := newCountingGitHub()
	rootDir := t.TempDir()
	if err := New(slog.New(slog.DiscardHandler), gh, rootDir).Run(t.Context()); err != nil {
		t.Fatal(err)
	}
	// actions/checkout@v6 is recorded in the lockfile, so its SHA is read from the lockfile
	// for both the workflow and the action file.
	if n := gh.commitSHA1["actions/checkout@v6"]; n != 0 {
		t.Errorf("actions/checkout@v6 must not be resolved through the API, but resolved %d times", n)
	}
	// actions/setup-go@v5 isn't locked, so its SHA is resolved through the API.
	if n := gh.commitSHA1["actions/setup-go@v5"]; n != 1 {
		t.Errorf("actions/setup-go@v5 must be resolved once, but resolved %d times", n)
	}
	// The action file of the locked commit is downloaded and cached.
	if _, err := os.Stat(filepath.Join(rootDir, "actions", "actions", "checkout", testLockedSHA, "action.yaml")); err != nil {
		t.Errorf("the action file of the locked commit must be cached: %v", err)
	}
}

func TestController_Run_invalidLockfile(t *testing.T) {
	// The lockfile can't be read, so SHAs are resolved through the API.
	t.Chdir(t.TempDir())
	writeFiles(t, map[string]string{
		".github/workflows/test.yaml": testWorkflowLocked,
		lockfile.Path:                 "version: 'v0.0.3'\nfoo: bar\n",
	})
	gh := newCountingGitHub()
	if err := New(slog.New(slog.DiscardHandler), gh, t.TempDir()).Run(t.Context()); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"actions/checkout@v6", "actions/setup-go@v5"} {
		if n := gh.commitSHA1[key]; n != 1 {
			t.Errorf("%s must be resolved once, but resolved %d times", key, n)
		}
	}
}
