package schema

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

var fullCommitSHAPattern = regexp.MustCompile(`^[0-9a-f]{40}$`)

// remoteRef is a remote action or reusable workflow `<owner>/<repo>[/<path>]@<ref>`.
type remoteRef struct {
	Owner string
	Repo  string
	Path  string
	Ref   string
}

// parseRemoteRef parses `<owner>/<repo>[/<path>]@<ref>`.
// It rejects empty, `.`, and `..` path elements so that the cache path can't escape the cache directory.
func parseRemoteRef(uses string) (*remoteRef, error) {
	fullPath, ref, ok := strings.Cut(uses, "@")
	if !ok || ref == "" {
		return nil, slogerr.With(errors.New("the format must be <owner>/<repo>[/<path>]@<ref>"), "uses", uses) //nolint:wrapcheck
	}
	elems := strings.Split(fullPath, "/")
	if len(elems) < 2 { //nolint:mnd
		return nil, slogerr.With(errors.New("the format must be <owner>/<repo>[/<path>]@<ref>"), "uses", uses) //nolint:wrapcheck
	}
	for _, elem := range elems {
		if elem == "" || elem == "." || elem == ".." || strings.Contains(elem, `\`) {
			return nil, slogerr.With(errors.New("the path includes an invalid element"), "uses", uses) //nolint:wrapcheck
		}
	}
	return &remoteRef{
		Owner: elems[0],
		Repo:  elems[1],
		Path:  strings.Join(elems[2:], "/"),
		Ref:   ref,
	}, nil
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
