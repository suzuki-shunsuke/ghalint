package schema

import (
	"log/slog"
	"path/filepath"
	"testing"

	"github.com/suzuki-shunsuke/ghalint/pkg/workflow"
)

const testSHA = "3d3c42e5aac5ba805825da76410c181273ba90b1"

func Test_parseRemoteRef(t *testing.T) {
	t.Parallel()
	data := []struct {
		name  string
		uses  string
		exp   *remoteRef
		isErr bool
	}{
		{
			name: "action",
			uses: "actions/checkout@v6",
			exp:  &remoteRef{Owner: "actions", Repo: "checkout", Ref: "v6"},
		},
		{
			name: "path",
			uses: "actions/cache/save@v4.2.0",
			exp:  &remoteRef{Owner: "actions", Repo: "cache", Path: "save", Ref: "v4.2.0"},
		},
		{
			name: "reusable workflow",
			uses: "suzuki-shunsuke/go-release-workflow/.github/workflows/release.yaml@" + testSHA,
			exp:  &remoteRef{Owner: "suzuki-shunsuke", Repo: "go-release-workflow", Path: ".github/workflows/release.yaml", Ref: testSHA},
		},
		{name: "no slash", uses: "foo@v1", isErr: true},
		{name: "no ref", uses: "actions/checkout", isErr: true},
		{name: "empty ref", uses: "actions/checkout@", isErr: true},
		{name: "empty owner", uses: "/checkout@v6", isErr: true},
		{name: "empty element", uses: "actions//checkout@v6", isErr: true},
		{name: "trailing slash", uses: "actions/checkout/@v6", isErr: true},
		{name: "parent directory", uses: "owner/repo/../../../../etc@" + testSHA, isErr: true},
		{name: "parent directory as repo", uses: "owner/..@" + testSHA, isErr: true},
		{name: "current directory", uses: "owner/repo/./foo@v1", isErr: true},
		{name: "backslash", uses: `owner/repo/..\..\foo@v1`, isErr: true},
		{name: "docker", uses: "docker://alpine@sha256:abc", isErr: true},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			r, err := parseRemoteRef(d.uses)
			if err != nil {
				if d.isErr {
					return
				}
				t.Fatal(err)
			}
			if d.isErr {
				t.Fatalf("error must be returned: %+v", r)
			}
			if *r != *d.exp {
				t.Fatalf("got %+v, wanted %+v", r, d.exp)
			}
		})
	}
}

func Test_remoteRef_cachePath(t *testing.T) {
	t.Parallel()
	rootDir := filepath.FromSlash("/root")
	data := []struct {
		name  string
		ref   *remoteRef
		sha   string
		elems []string
		exp   string
		isErr bool
	}{
		{
			name:  "action",
			ref:   &remoteRef{Owner: "actions", Repo: "cache", Path: "save"},
			sha:   testSHA,
			elems: []string{"action.yaml"},
			exp:   filepath.FromSlash("/root/actions/actions/cache/" + testSHA + "/save/action.yaml"),
		},
		{
			name: "reusable workflow",
			ref:  &remoteRef{Owner: "owner", Repo: "repo", Path: ".github/workflows/test.yaml"},
			sha:  testSHA,
			exp:  filepath.FromSlash("/root/actions/owner/repo/" + testSHA + "/.github/workflows/test.yaml"),
		},
		{
			name:  "sha includes parent directories",
			ref:   &remoteRef{Owner: "owner", Repo: "repo"},
			sha:   "../../../../x/" + testSHA,
			isErr: true,
		},
		{
			name:  "sha has a prefix",
			ref:   &remoteRef{Owner: "owner", Repo: "repo"},
			sha:   "v1-" + testSHA,
			isErr: true,
		},
		{
			name:  "short sha",
			ref:   &remoteRef{Owner: "owner", Repo: "repo"},
			sha:   "3d3c42e",
			isErr: true,
		},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			p, err := d.ref.cachePath(rootDir, d.sha, d.elems...)
			if err != nil {
				if d.isErr {
					return
				}
				t.Fatal(err)
			}
			if d.isErr {
				t.Fatalf("error must be returned: %s", p)
			}
			if p != d.exp {
				t.Fatalf("got %s, wanted %s", p, d.exp)
			}
		})
	}
}

// Invalid uses must be rejected before calling the GitHub API or touching the cache.
// gh is nil, so the test panics if the GitHub API is called.
func Test_invalidUses(t *testing.T) {
	t.Parallel()
	for _, uses := range []string{"foo@v1", "owner/repo/../../../../etc@" + testSHA} {
		t.Run(uses, func(t *testing.T) {
			t.Parallel()
			rootDir := t.TempDir()
			vs := &validateStep{
				step:    &workflow.Step{Uses: uses},
				logger:  slog.New(slog.DiscardHandler),
				rootDir: rootDir,
			}
			if err := vs.readAction(t.Context(), &workflow.Action{}); err == nil {
				t.Fatal("error must be returned")
			}
			vj := &validateJob{
				job:     &workflow.Job{Uses: uses},
				logger:  slog.New(slog.DiscardHandler),
				rootDir: rootDir,
			}
			if err := vj.read(t.Context(), &ReusableWorkflow{}); err == nil {
				t.Fatal("error must be returned")
			}
		})
	}
}
