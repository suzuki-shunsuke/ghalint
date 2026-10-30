package controller_test

import (
	"bytes"
	"fmt"
	"log/slog"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghalint/pkg/controller"
	"github.com/suzuki-shunsuke/ghalint/pkg/lockfile"
)

const testWorkflow = `on: push
permissions: {}
jobs:
  test:
    runs-on: ubuntu-latest
    timeout-minutes: 10
    permissions: {}
    steps:
      - uses: actions/checkout@%s
        with:
          persist-credentials: false
`

const testLockfile = `version: 'v0.0.3'
workflows:
  '.github/workflows/test.yaml':
    - 'actions/checkout@v6'
dependencies:
  'actions/checkout@v6':
    ref: 'v6'
    commit: 'sha1-3d3c42e5aac5ba805825da76410c181273ba90b1'
    owner_id: 44036562
    repo_id: 197814629
`

const testInvalidLockfile = "version: 'v0.0.3'\nfoo: bar\n"

const (
	logReadLockfile = "read the lockfile"
	logStepViolates = "the step violates policies"
)

func TestController_Run_lockfile(t *testing.T) { //nolint:funlen
	t.Parallel()
	data := []struct {
		name     string
		ref      string
		lockfile string
		isErr    bool
		// logs must be in the log output, and noLogs must not.
		logs   []string
		noLogs []string
	}{
		{
			name:     "action is locked",
			ref:      "v6",
			lockfile: testLockfile,
			noLogs:   []string{logReadLockfile, logStepViolates},
		},
		{
			name:     "lockfile exists but action isn't locked",
			ref:      "v5",
			lockfile: testLockfile,
			isErr:    true,
			logs:     []string{logStepViolates},
			noLogs:   []string{logReadLockfile},
		},
		{
			name:   "lockfile doesn't exist and action is pinned by SHA",
			ref:    "3d3c42e5aac5ba805825da76410c181273ba90b1",
			noLogs: []string{logReadLockfile, logStepViolates},
		},
		{
			name:   "lockfile doesn't exist and action isn't pinned",
			ref:    "v6",
			isErr:  true,
			logs:   []string{logStepViolates},
			noLogs: []string{logReadLockfile},
		},
		{
			// Actions are validated as if the lockfile didn't exist.
			name:     "lockfile is invalid and action isn't pinned",
			ref:      "v6",
			lockfile: testInvalidLockfile,
			isErr:    true,
			logs:     []string{logReadLockfile, logStepViolates},
		},
		{
			// It's just a warning if all actions are pinned by SHA.
			name:     "lockfile is invalid and action is pinned by SHA",
			ref:      "3d3c42e5aac5ba805825da76410c181273ba90b1",
			lockfile: testInvalidLockfile,
			logs:     []string{logReadLockfile},
			noLogs:   []string{logStepViolates},
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			fs := newFs(t, d.ref, d.lockfile)
			buf := &bytes.Buffer{}
			logger := slog.New(slog.NewTextHandler(buf, nil))
			err := controller.New(fs).Run(t.Context(), logger, "")
			if d.isErr && err == nil {
				t.Fatal("error must be returned")
			}
			if !d.isErr && err != nil {
				t.Fatalf("unexpected error: %v\n%s", err, buf.String())
			}
			assertLogs(t, buf.String(), d.logs, d.noLogs)
		})
	}
}

// newFs returns a file system with a workflow using actions/checkout@<ref> and the lockfile.
// The lockfile isn't created if it's empty.
func newFs(t *testing.T, ref, lf string) afero.Fs {
	t.Helper()
	fs := afero.NewMemMapFs()
	if err := afero.WriteFile(fs, ".github/workflows/test.yaml", []byte(fmt.Sprintf(testWorkflow, ref)), 0o644); err != nil {
		t.Fatal(err)
	}
	if lf == "" {
		return fs
	}
	if err := afero.WriteFile(fs, lockfile.Path, []byte(lf), 0o644); err != nil {
		t.Fatal(err)
	}
	return fs
}

func assertLogs(t *testing.T, out string, logs, noLogs []string) {
	t.Helper()
	for _, s := range logs {
		if !strings.Contains(out, s) {
			t.Errorf("log %q must be output\n%s", s, out)
		}
	}
	for _, s := range noLogs {
		if strings.Contains(out, s) {
			t.Errorf("log %q must not be output\n%s", s, out)
		}
	}
}
