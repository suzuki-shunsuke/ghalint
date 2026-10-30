package controller_test

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

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
	logLockfileError = "lockfile_error"
	logReadLockfile  = "read the lockfile"
	logStepViolates  = "the step violates policies"
)

func TestController_Run_lockfile(t *testing.T) { //nolint:funlen
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
			// Required fields are validated only for entries the workflows reference.
			name: "lockfile has a broken entry that no workflow uses",
			ref:  "v6",
			lockfile: testLockfile + `  'actions/setup-go@v6':
    ref: 'v6'
    owner_id: 44036562
    repo_id: 192624594
`,
			noLogs: []string{logReadLockfile, logStepViolates},
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
			noLogs: []string{logReadLockfile, logLockfileError},
		},
		{
			// Actions are validated as if the lockfile didn't exist.
			name:     "lockfile is invalid and action isn't pinned",
			ref:      "v6",
			lockfile: testInvalidLockfile,
			isErr:    true,
			logs:     []string{logReadLockfile, logStepViolates, logLockfileError},
		},
		{
			// It's just a warning if all actions are pinned by SHA.
			name:     "lockfile is invalid and action is pinned by SHA",
			ref:      "3d3c42e5aac5ba805825da76410c181273ba90b1",
			lockfile: testInvalidLockfile,
			logs:     []string{logReadLockfile},
			noLogs:   []string{logStepViolates, logLockfileError},
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Chdir(t.TempDir())
			writeFile(t, ".github/workflows/test.yaml", fmt.Sprintf(testWorkflow, d.ref))
			if d.lockfile != "" {
				writeFile(t, lockfile.Path, d.lockfile)
			}
			buf := &bytes.Buffer{}
			logger := slog.New(slog.NewTextHandler(buf, nil))
			err := controller.New().Run(t.Context(), logger, "")
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

// writeFile writes a file relative to the current directory, creating its parent directories.
func writeFile(t *testing.T, p, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil { //nolint:gosec
		t.Fatal(err)
	}
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

func TestController_Run_noWorkflow(t *testing.T) {
	t.Chdir(t.TempDir())
	// The lockfile has a broken entry, but it isn't read because there is no workflow.
	writeFile(t, lockfile.Path, testInvalidLockfile)
	buf := &bytes.Buffer{}
	logger := slog.New(slog.NewTextHandler(buf, nil))
	if err := controller.New().Run(t.Context(), logger, ""); err != nil {
		t.Fatal(err)
	}
	assertLogs(t, buf.String(), nil, []string{logReadLockfile})
}
