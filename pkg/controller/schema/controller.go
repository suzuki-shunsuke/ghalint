package schema

import (
	"context"
	"log/slog"

	"github.com/suzuki-shunsuke/ghalint/pkg/github"
)

type Controller struct {
	logger  *slog.Logger
	gh      GitHub
	rootDir string
}

func New(logger *slog.Logger, gh GitHub, rootDir string) *Controller {
	return &Controller{
		logger:  logger,
		gh:      gh,
		rootDir: rootDir,
	}
}

type GitHub interface {
	GetCommitSHA1(ctx context.Context, owner, repo, ref, lastSHA string) (string, *github.Response, error)
	GetContents(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentGetOptions) (fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent, resp *github.Response, err error)
}
