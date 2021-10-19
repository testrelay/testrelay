package vcs

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/github"
)

type Repo struct {
	ID       int64  `json:"id"`
	FullName string `json:"full_name"`
}

type GithubRepoCollector struct {
	GithubPK    []byte
	GithubAppID int64
}

func NewGithubRepoCollector(privateKey string, appID int64) (GithubRepoCollector, error) {
	b, err := os.ReadFile(privateKey)
	if err != nil {
		return GithubRepoCollector{}, fmt.Errorf("could not read github priv key file %w", err)
	}

	return GithubRepoCollector{
		GithubAppID: appID,
		GithubPK:    b,
	}, nil
}

func (g GithubRepoCollector) CollectRepos(installationID int64) ([]Repo, error) {
	itr, err := ghinstallation.New(
		http.DefaultTransport,
		g.GithubAppID,
		installationID,
		g.GithubPK,
	)
	if err != nil {
		return nil, fmt.Errorf("could not init new github installation %w", err)
	}

	appClient := github.NewClient(
		&http.Client{Transport: itr},
	)
	repos, _, err := appClient.Apps.ListRepos(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list repos %w", err)
	}

	var qrepos = make([]Repo, len(repos))
	for i, repo := range repos {
		qrepos[i] = Repo{
			ID:       *repo.ID,
			FullName: *repo.FullName,
		}
	}

	return qrepos, nil
}
