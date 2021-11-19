package vcs

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v39/github"

	"github.com/testrelay/testrelay/backend/internal/core"
)

// GithubRepoCollector implements a vcs.RepoCollector interface for github.
// It uses github app installation to list repositories that the app installation
// has access to.
type GithubRepoCollector struct {
	newInstallation InstallationFunc
}

// NewGithubRepoCollector returns a GithubRepoCollector first parsing a privKeyLoc.
// The privKeyLoc mus be a valid path to a github app private key. This can be downloaded
// settings overview of the github app.
func NewGithubRepoCollector(privKeyLoc string, appID int64) (GithubRepoCollector, error) {
	b, err := os.ReadFile(privKeyLoc)
	if err != nil {
		return GithubRepoCollector{}, fmt.Errorf("could not read github priv key file %w", err)
	}

	return GithubRepoCollector{
		newInstallation: NewGithubAppInstallationFunc(appID, b),
	}, nil
}

// CollectRepos fetches a list of the github repos scoped to the passed installationID.
// The github installation must have read access to the repositories otherwise CollectRepos will fail.
func (g GithubRepoCollector) CollectRepos(installationID int64) ([]core.Repo, error) {
	c, err := g.newInstallation(installationID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate installation with id %d %w", installationID, err)
	}

	repos, _, err := c.ListRepos(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list repos %w", err)
	}

	var qrepos = make([]core.Repo, len(repos.Repositories))
	for i, repo := range repos.Repositories {
		qrepos[i] = core.Repo{
			ID:       *repo.ID,
			FullName: *repo.FullName,
		}
	}

	return qrepos, nil
}

// InstallationFunc represents a function that returns a new client for the given installationID.
// In most cases this represents a given github APP installation.
type InstallationFunc func(installationID int64) (GithubInstallationClient, error)

// NewGithubAppInstallationFunc returns a InstallationFunc for a given github app.
// NewGithubAppInstallationFunc uses ghinstallation to init a new github client with
// modified http transport. ghinstallation pkg takes care of generating/refreshing a valid
// access token for the github installation per http request.
//
// pk must be a valid primary key for the appID given. This can be found/generated under the
// github app developer settings page.
func NewGithubAppInstallationFunc(appID int64, pk []byte) InstallationFunc {
	return func(installationID int64) (GithubInstallationClient, error) {
		itr, err := ghinstallation.New(
			http.DefaultTransport,
			appID,
			installationID,
			pk,
		)
		if err != nil {
			return GithubInstallationClient{}, fmt.Errorf("could not init new github installation %w", err)
		}

		return GithubInstallationClient{
			InstallationClient: GithubInstallationWrapper{
				client: github.NewClient(
					&http.Client{Transport: itr},
				),
			},
		}, nil
	}
}

// InstallationClient defines an interface around a vcs installation that is scoped for read-only access to repos.
type InstallationClient interface {
	ListRepos(ctx context.Context, opts *github.ListOptions) (*github.ListRepositories, *github.Response, error)
	DownloadRepo(ctx context.Context, url string) (*bytes.Buffer, error)
}

// GithubInstallationClient wraps an InstallationClient interface with a hard type. This is done for extra interface/
// parameter flexibility down the line.
type GithubInstallationClient struct {
	InstallationClient
}

// GithubInstallationWrapper adapts the github.client into a InstallationClient interface.
type GithubInstallationWrapper struct {
	client *github.Client
}

// DownloadRepo downloads a zip of the given repo the provided url and returns it as a bytes.Buffer.
func (g GithubInstallationWrapper) DownloadRepo(ctx context.Context, url string) (*bytes.Buffer, error) {
	owner, repo := getRepoName(url)

	u, _, err := g.client.Repositories.GetArchiveLink(ctx, owner, repo, github.Zipball, nil, true)
	if err != nil {
		return nil, fmt.Errorf("could not get archive link for repo %s/%s %w", owner, repo, err)
	}

	req, _ := g.client.NewRequest("GET", u.String(), nil)
	buf := bytes.NewBuffer([]byte{})
	_, err = g.client.Do(context.Background(), req, buf)
	if err != nil {
		return nil, fmt.Errorf("problem downloading zipFile for repo %s/%s %w", owner, repo, err)
	}

	return buf, nil
}

// ListRepos is a slim wrapper around the Apps.ListRepos. In future this could be adapted with more
// generic params defined on the InstallationClient interface.
func (g GithubInstallationWrapper) ListRepos(ctx context.Context, opts *github.ListOptions) (*github.ListRepositories, *github.Response, error) {
	return g.client.Apps.ListRepos(ctx, opts)
}
