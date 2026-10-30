package schema

import (
	"context"
	"log/slog"
	"path/filepath"
	"testing"

	"github.com/suzuki-shunsuke/ghalint/pkg/github"
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

func Test_invalidReusableWorkflow(t *testing.T) {
	t.Parallel()
	// The workflow path is required. gh is nil, so the test panics if the GitHub API is called.
	vj := &validateJob{
		job:     &workflow.Job{Uses: "owner/repo@" + testSHA},
		logger:  slog.New(slog.DiscardHandler),
		rootDir: t.TempDir(),
	}
	if err := vj.read(t.Context(), &ReusableWorkflow{}); err == nil {
		t.Fatal("error must be returned")
	}
}

func Test_validateStep_validate_docker(t *testing.T) {
	t.Parallel()
	// Docker container actions are skipped. gh is nil, so the test panics if the GitHub API is called.
	vs := &validateStep{
		step:    &workflow.Step{Uses: "docker://alpine@sha256:abc"},
		logger:  slog.New(slog.DiscardHandler),
		rootDir: t.TempDir(),
	}
	if err := vs.validate(t.Context()); err != nil {
		t.Fatal(err)
	}
}

// directoryGitHub returns a directory listing for any path.
type directoryGitHub struct{}

func (directoryGitHub) GetCommitSHA1(_ context.Context, _, _, _, _ string) (string, *github.Response, error) {
	return testSHA, nil, nil
}

func (directoryGitHub) GetContents(_ context.Context, _, _, _ string, _ *github.RepositoryContentGetOptions) (*github.RepositoryContent, []*github.RepositoryContent, *github.Response, error) {
	return nil, []*github.RepositoryContent{}, nil, nil
}

// It must return an error instead of panicking if the path is a directory.
func Test_downloadDirectory(t *testing.T) {
	t.Parallel()
	vs := &validateStep{
		step:    &workflow.Step{Uses: "owner/repo@v1"},
		logger:  slog.New(slog.DiscardHandler),
		gh:      directoryGitHub{},
		rootDir: t.TempDir(),
	}
	if err := vs.readAction(t.Context(), &workflow.Action{}); err == nil {
		t.Fatal("error must be returned")
	}
	vj := &validateJob{
		job:     &workflow.Job{Uses: "owner/repo/.github/workflows@v1"},
		logger:  slog.New(slog.DiscardHandler),
		gh:      directoryGitHub{},
		rootDir: t.TempDir(),
	}
	if err := vj.read(t.Context(), &ReusableWorkflow{}); err == nil {
		t.Fatal("error must be returned")
	}
}
