//go:build e2e
// +build e2e

package vcs_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/google/go-github/v39/github"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"

	"github.com/testrelay/testrelay/backend/internal/core"
	"github.com/testrelay/testrelay/backend/internal/vcs"
)

func TestGithubClient(t *testing.T) {
	require.NoError(t, godotenv.Overload("./test_assets/e2e.env"))

	at := os.Getenv("GITHUB_ACCESS_TOKEN")
	kl := os.Getenv("GITHUB_PRIVATE_KEY_LOCATION")
	appID, _ := strconv.ParseInt(os.Getenv("GITHUB_APP_ID"), 10, 64)
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_ACCESS_TOKEN")},
	)

	tc := oauth2.NewClient(context.Background(), ts)
	rawClient := github.NewClient(tc)

	githubClient, err := vcs.NewGithubClient(vcs.GithubInterviewerConfig{
		AccessToken: at,
		Username:    os.Getenv("GITHUB_USERNAME"),
		Email:       os.Getenv("GITHUB_EMAIl"),
	}, kl, appID)
	require.NoError(t, err)

	t.Run("Upload", func(t *testing.T) {
		unix := time.Now().Unix()
		name := fmt.Sprintf("e2e-assignment-%d", unix)
		r := &github.Repository{
			Name:         github.String(name),
			Private:      github.Bool(true),
			MasterBranch: github.String("master"),
		}
		repo, _, err := rawClient.Repositories.Create(context.Background(), "", r)
		require.NoError(t, err)
		defer func() {
			_, err := rawClient.Repositories.Delete(context.Background(), repo.GetOwner().GetLogin(), repo.GetName())
			assert.NoError(t, err)
		}()

		installationID, _ := strconv.ParseInt(os.Getenv("TEST_GITHUB_INSTALLATION"), 10, 64)
		err = githubClient.Upload(core.UploadDetails{
			ID:             unix,
			VCSRepoURL:     repo.GetCloneURL(),
			TestVCSRepoURL: os.Getenv("TEST_GITHUB_REPO_URL"),
			InstallationID: installationID,
		})
		require.NoError(t, err)

		owner := repo.GetOwner().GetLogin()
		repoName := repo.GetName()
		commits, _, err := rawClient.Repositories.ListCommits(context.Background(), owner, repoName, nil)
		require.NoError(t, err)

		require.NoError(t, err)
		require.Len(t, commits, 1)

		c := commits[0].GetCommit()
		assert.Equal(t, "start test", c.GetMessage())

		tree, _, err := rawClient.Git.GetTree(context.Background(), owner, repoName, c.GetTree().GetSHA(), true)
		assert.NoError(t, err)

		filenames := make([]string, 0, 2)
		for _, e := range tree.Entries {
			if e != nil {
				filenames = append(filenames, e.GetPath())
			}
		}

		assert.Contains(t, filenames, "test/index.txt")
		assert.Contains(t, filenames, "echo.txt")
	})
}
