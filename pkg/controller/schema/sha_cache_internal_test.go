package schema

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/suzuki-shunsuke/ghalint/pkg/github"
)

const testSHA2 = "1bd1e32a3bdc45362d1e726936510720a7c30a57"

// countingGitHub counts API calls.
// GetCommitSHA1 returns testSHA, or an error if the ref is "error".
// GetContents returns an action without inputs.
type countingGitHub struct {
	commitSHA1 map[string]int
	contents   int
}

func newCountingGitHub() *countingGitHub {
	return &countingGitHub{commitSHA1: map[string]int{}}
}

func (g *countingGitHub) GetCommitSHA1(_ context.Context, owner, repo, ref, _ string) (string, *github.Response, error) {
	g.commitSHA1[owner+"/"+repo+"@"+ref]++
	switch ref {
	case "error":
		return "", nil, errors.New("not found")
	case "v5":
		return testSHA2, nil, nil
	}
	return testSHA, nil, nil
}

func (g *countingGitHub) GetContents(_ context.Context, _, _, _ string, _ *github.RepositoryContentGetOptions) (*github.RepositoryContent, []*github.RepositoryContent, *github.Response, error) {
	g.contents++
	content := "name: test\nruns:\n  using: node24\n  main: index.js\n"
	return &github.RepositoryContent{Content: &content}, nil, nil, nil
}

func Test_shaCache_GetCommitSHA1(t *testing.T) {
	t.Parallel()
	gh := newCountingGitHub()
	c := newSHACache(gh)
	for _, d := range []struct {
		owner, repo, ref, exp string
	}{
		{owner: "actions", repo: "checkout", ref: "v4", exp: testSHA},
		{owner: "actions", repo: "checkout", ref: "v4", exp: testSHA},
		// owner and repo are case insensitive.
		{owner: "Actions", repo: "Checkout", ref: "v4", exp: testSHA},
		{owner: "actions", repo: "checkout", ref: "v5", exp: testSHA2},
	} {
		sha, _, err := c.GetCommitSHA1(t.Context(), d.owner, d.repo, d.ref, "")
		if err != nil {
			t.Fatal(err)
		}
		if sha != d.exp {
			t.Fatalf("got %s, wanted %s", sha, d.exp)
		}
	}
	if n := gh.commitSHA1["actions/checkout@v4"]; n != 1 {
		t.Fatalf("actions/checkout@v4 must be resolved once, but resolved %d times", n)
	}
	if n := gh.commitSHA1["actions/checkout@v5"]; n != 1 {
		t.Fatalf("actions/checkout@v5 must be resolved once, but resolved %d times", n)
	}
	// Errors aren't cached.
	for range 2 {
		if _, _, err := c.GetCommitSHA1(t.Context(), "actions", "checkout", "error", ""); err == nil {
			t.Fatal("error must be returned")
		}
	}
	if n := gh.commitSHA1["actions/checkout@error"]; n != 2 {
		t.Fatalf("errors must not be cached, but the API was called %d times", n)
	}
}

const testWorkflowManySteps = `on: push
jobs:
  foo:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/checkout@v4
  bar:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
      - uses: actions/setup-go@v5
`

func TestController_Run_shaCache(t *testing.T) {
	t.Chdir(t.TempDir())
	if err := os.MkdirAll(filepath.Join(".github", "workflows"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(".github", "workflows", "test.yaml"), []byte(testWorkflowManySteps), 0o644); err != nil { //nolint:gosec
		t.Fatal(err)
	}
	gh := newCountingGitHub()
	ctrl := New(slog.New(slog.DiscardHandler), gh, t.TempDir())
	if err := ctrl.Run(t.Context()); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"actions/checkout@v4", "actions/setup-go@v5"} {
		if n := gh.commitSHA1[key]; n != 1 {
			t.Errorf("%s must be resolved once, but resolved %d times", key, n)
		}
	}
	// The action files are downloaded once per action and read from the cache after that.
	if gh.contents != 2 {
		t.Errorf("action files must be downloaded twice, but downloaded %d times", gh.contents)
	}
}
