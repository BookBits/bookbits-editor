package models

import (
	"context"

	"github.com/google/go-github/v61/github"
	"golang.org/x/oauth2"
)

type GitClient struct {
	Client *github.Client
	Repo string
	Owner string
	Context context.Context
}

func NewGitClient(vars EnvVars) GitClient {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: vars.GitToken})

	tc := oauth2.NewClient(ctx, ts)

	return GitClient{
		Client: github.NewClient(tc),
		Context: ctx,
		Repo: vars.GitRepo,
		Owner: vars.GitOwner,
	}
}
