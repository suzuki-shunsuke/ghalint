package schema

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

var (
	fullCommitSHAPattern = regexp.MustCompile(`^[0-9a-f]{40}$`)
	// ownerRepoPattern is the characters GitHub allows in owner and repository names.
	// owner and repo are put into GitHub API URLs without escaping.
	ownerRepoPattern = regexp.MustCompile(`^[A-Za-z0-9_.-]+$`)
	// reusableWorkflowPathPattern is the path of a reusable workflow, which must be in .github/workflows.
	reusableWorkflowPathPattern = regexp.MustCompile(`^\.github/workflows/[^/]+\.ya?ml$`)
)

const (
	remoteActionFormat     = "<owner>/<repo>[/<path>]@<ref>"
	reusableWorkflowFormat = "<owner>/<repo>/.github/workflows/<file>.y(a)ml@<ref>"
)

func formatError(format, uses string) error {
	return slogerr.With(errors.New("the format must be "+format), "uses", uses) //nolint:wrapcheck
}

// remoteRef is a remote action or reusable workflow `<owner>/<repo>[/<path>]@<ref>`.
type remoteRef struct {
	Owner string
	Repo  string
	Path  string
	Ref   string
}

// parseRemoteRef parses `<owner>/<repo>[/<path>]@<ref>`.
// It rejects empty, `.`, and `..` path elements so that the cache path can't escape the cache directory,
// and owner and repo with characters GitHub doesn't allow.
func parseRemoteRef(uses string) (*remoteRef, error) {
	fullPath, ref, ok := strings.Cut(uses, "@")
	if !ok || ref == "" {
		return nil, formatError(remoteActionFormat, uses)
	}
	elems := strings.Split(fullPath, "/")
	if len(elems) < 2 { //nolint:mnd
		return nil, formatError(remoteActionFormat, uses)
	}
	for _, elem := range elems {
		if !validPathElement(elem) {
			return nil, slogerr.With(errors.New("the path includes an invalid element"), "uses", uses) //nolint:wrapcheck
		}
	}
	if !ownerRepoPattern.MatchString(elems[0]) || !ownerRepoPattern.MatchString(elems[1]) {
		return nil, slogerr.With(errors.New("owner or repo includes an invalid character"), "uses", uses) //nolint:wrapcheck
	}
	return &remoteRef{
		Owner: elems[0],
		Repo:  elems[1],
		Path:  strings.Join(elems[2:], "/"),
		Ref:   ref,
	}, nil
}

// validPathElement reports whether elem can be a path element of the cache path.
// A backslash is a path separator on Windows.
func validPathElement(elem string) bool {
	return elem != "" && elem != "." && elem != ".." && !strings.Contains(elem, `\`)
}

// resolveSHA returns the full length commit SHA of the ref.
// It calls the GitHub API only if the ref isn't a full length commit SHA.
func (r *remoteRef) resolveSHA(ctx context.Context, gh GitHub) (string, error) {
	if fullCommitSHAPattern.MatchString(r.Ref) {
		return r.Ref, nil
	}
	sha, _, err := gh.GetCommitSHA1(ctx, r.Owner, r.Repo, r.Ref, "")
	if err != nil {
		return "", fmt.Errorf("get commit SHA1: %w", err)
	}
	return sha, nil
}

// cachePath returns `<rootDir>/actions/<owner>/<repo>/<sha>/<path>/<elems...>`.
// It returns an error if sha isn't a full length commit SHA or the path would escape `<rootDir>/actions`.
func (r *remoteRef) cachePath(rootDir, sha string, elems ...string) (string, error) {
	if !fullCommitSHAPattern.MatchString(sha) {
		return "", slogerr.With(errors.New("the commit SHA is invalid"), "sha", sha) //nolint:wrapcheck
	}
	rel := filepath.Join(append([]string{r.Owner, r.Repo, sha, filepath.FromSlash(r.Path)}, elems...)...)
	if !filepath.IsLocal(rel) {
		return "", fmt.Errorf("the cache path escapes the cache directory: %s", rel)
	}
	return filepath.Join(rootDir, "actions", rel), nil
}
