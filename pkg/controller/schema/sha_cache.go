package schema

import (
	"context"
	"strings"

	"github.com/suzuki-shunsuke/ghalint/pkg/github"
)

// shaCache caches commit SHAs of refs in memory during a run,
// so that a ref used by many steps and jobs is resolved through the GitHub API only once.
// It isn't persisted because a tag or branch may be moved to another commit.
// It isn't safe for concurrent use.
type shaCache struct {
	GitHub

	shas map[string]string
}

func newSHACache(gh GitHub) *shaCache {
	return &shaCache{
		GitHub: gh,
		shas:   map[string]string{},
	}
}

// GetCommitSHA1 returns the cached SHA if the ref has been resolved.
// The returned response is nil in that case.
// Errors aren't cached.
func (c *shaCache) GetCommitSHA1(ctx context.Context, owner, repo, ref, lastSHA string) (string, *github.Response, error) {
	// owner and repo are case insensitive, but ref is case sensitive.
	key := strings.ToLower(owner+"/"+repo) + "@" + ref + "#" + lastSHA
	if sha, ok := c.shas[key]; ok {
		return sha, nil, nil
	}
	sha, resp, err := c.GitHub.GetCommitSHA1(ctx, owner, repo, ref, lastSHA)
	if err != nil {
		return "", resp, err //nolint:wrapcheck
	}
	c.shas[key] = sha
	return sha, resp, nil
}
