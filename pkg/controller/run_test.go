package controller_test

import (
	"log/slog"
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
      - uses: actions/checkout@v6
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

func TestController_Run_lockfile(t *testing.T) {
	t.Parallel()
	data := []struct {
		name     string
		lockfile string
		isErr    bool
	}{
		{
			name:     "action is locked",
			lockfile: testLockfile,
		},
		{
			name:  "lockfile doesn't exist",
			isErr: true,
		},
		{
			name:     "lockfile is invalid",
			lockfile: "version: 'v0.0.3'\nfoo: bar\n",
			isErr:    true,
		},
	}
	logger := slog.New(slog.DiscardHandler)
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			fs := afero.NewMemMapFs()
			if err := afero.WriteFile(fs, ".github/workflows/test.yaml", []byte(testWorkflow), 0o644); err != nil {
				t.Fatal(err)
			}
			if d.lockfile != "" {
				if err := afero.WriteFile(fs, lockfile.Path, []byte(d.lockfile), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			if err := controller.New(fs).Run(t.Context(), logger, ""); err != nil {
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
