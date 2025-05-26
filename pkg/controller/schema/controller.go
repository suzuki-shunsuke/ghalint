package schema

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ghalint/pkg/github"
)

type Controller struct {
	fs      afero.Fs
	logE    *logrus.Entry
	gh      GitHub
	rootDir string
}

func New(fs afero.Fs, logE *logrus.Entry, gh GitHub, rootDir string) *Controller {
	return &Controller{
		fs:      fs,
		logE:    logE,
		gh:      gh,
		rootDir: rootDir,
	}
}

type GitHub interface {
	GetCommitSHA1(ctx context.Context, owner, repo, ref, lastSHA string) (string, *github.Response, error)
	GetContents(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentGetOptions) (fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent, resp *github.Response, err error)
}
