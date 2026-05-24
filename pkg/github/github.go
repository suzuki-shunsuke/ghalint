package github

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/google/go-github/v88/github"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/keyring/ghtoken"
	"golang.org/x/oauth2"
)

type (
	ListOptions                 = github.ListOptions
	Reference                   = github.Reference
	Response                    = github.Response
	RepositoryTag               = github.RepositoryTag
	RepositoryRelease           = github.RepositoryRelease
	Client                      = github.Client
	GitObject                   = github.GitObject
	Commit                      = github.Commit
	RepositoryContentGetOptions = github.RepositoryContentGetOptions
	RepositoryContent           = github.RepositoryContent
)

func New(ctx context.Context, logger *slog.Logger) (*Client, error) {
	client, err := github.NewClient(github.WithHTTPClient(getHTTPClientForGitHub(ctx, logger, getGitHubToken())))
	if err != nil {
		return nil, fmt.Errorf("create a GitHub client: %w", err)
	}
	return client, nil
}

func getGitHubToken() string {
	return os.Getenv("GITHUB_TOKEN")
}

func checkKeyringEnabled() bool {
	return os.Getenv("GHALINT_KEYRING_ENABLED") == "true"
}

func getHTTPClientForGitHub(ctx context.Context, logger *slog.Logger, token string) *http.Client {
	if token == "" {
		if checkKeyringEnabled() {
			return oauth2.NewClient(ctx, ghtoken.NewTokenSource(logger, KeyService))
		}
		return http.DefaultClient
	}
	return oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))
}
