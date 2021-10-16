package github

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

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

func NewGithubRepoCollectorFromENV() (GithubRepoCollector, error) {
	pkey := os.Getenv("GITHUB_PRIVATE_KEY")
	pkey = strings.ReplaceAll(pkey, `\n`, "\n")

	appID := os.Getenv("GITHUB_APP_ID")
	id, err := strconv.ParseInt(appID, 10, 64)
	if err != nil {
		return GithubRepoCollector{}, fmt.Errorf("could not parse github app id %w", err)
	}

	return GithubRepoCollector{
		GithubAppID: id,
		GithubPK:    []byte(pkey),
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
