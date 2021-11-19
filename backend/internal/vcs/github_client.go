package vcs

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	gitHttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"

	"github.com/testrelay/testrelay/backend/internal/core"
)

var (
	letters = []rune("abcdefghijklmnopqrstuvwxyz")
)

// GithubClient handles communicating with the github api to orchestrate github repository management.
// GithubClient communicates with the github API using two main authentication methods: an installationID
// and an accessToken.
//
// Access token interactions are scoped to assignment repository generation and management.
// These repositories are fully maintained by testrelay and thus use a single github user to perform actions.
//
// Installation interactions are scoped to reading the contents for business test repositories that have been
// given access through an app installation.
type GithubClient struct {
	client          *github.Client
	intervConf      GithubInterviewerConfig
	newInstallation InstallationFunc
}

// GithubInterviewerConfig represents fields required to generate repos using a personal access token.
type GithubInterviewerConfig struct {
	AccessToken string
	Username    string
	Email       string
}

// NewGithubClient returns a GithubClient with all the underlying github api client initialized.
// repoCreatorAccessToken must be an github access token for a user that has full permissions to manage
// github repos. This includes deletion and collaborator management. appPrivKeyLoc must be the file path
// of a github private key that is the same app as appID.
func NewGithubClient(intervConf GithubInterviewerConfig, appPrivKeyLoc string, appID int64) (*GithubClient, error) {
	b, err := os.ReadFile(appPrivKeyLoc)
	if err != nil {
		return nil, fmt.Errorf("could not read github priv key file %w", err)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: intervConf.AccessToken},
	)

	tc := oauth2.NewClient(context.Background(), ts)
	return &GithubClient{
		client:          github.NewClient(tc),
		newInstallation: NewGithubAppInstallationFunc(appID, b),
		intervConf:      intervConf,
	}, nil
}

func (c GithubClient) CreateRepo(bName, username string, id int) (string, error) {
	name := makeRepoName(bName, username, id)
	r := &github.Repository{
		Name:         github.String(name),
		Private:      github.Bool(true),
		Description:  github.String(username + " code assignment for " + bName),
		MasterBranch: github.String("master"),
	}

	repo, _, err := c.client.Repositories.Create(context.Background(), "", r)
	if err != nil {
		return "", fmt.Errorf("could not create repo %w", err)
	}

	login := repo.GetOwner().GetLogin()
	repoName := repo.GetName()
	err = c.addCollaborator(login, repoName, username)
	if err != nil {
		return "", err
	}

	return repo.GetCloneURL(), nil
}

var (
	repl  = regexp.MustCompile("https://github.com/")
	grepl = regexp.MustCompile("\\.git")

	ErrorAlreadyCollaborator = errors.New("already collaborator")
)

func (c GithubClient) AddCollaborator(repo string, username string) error {
	repo = repl.ReplaceAllString(repo, "")
	repo = grepl.ReplaceAllString(repo, "")

	pieces := strings.Split(repo, "/")
	owner := pieces[0]
	name := pieces[1]

	var colabs []*github.User
	var err error
	for i := 0; i < 3; i++ {
		colabs, _, err = c.client.Repositories.ListCollaborators(context.Background(), owner, name, nil)
		if err == nil {
			break
		}
		log.Printf("could not list colaborator %s sleeping", err)
		time.Sleep(time.Second)
		i++
	}

	if err != nil {
		return fmt.Errorf("could not list colaborators for repo %s %s %w", owner, name, err)
	}

	invites, _, err := c.client.Repositories.ListInvitations(context.Background(), owner, name, nil)
	if err != nil {
		return fmt.Errorf("could not list invites for repo %s %s %w", owner, name, err)
	}

	alreadyActive := make(map[string]struct{})
	for _, invite := range invites {
		alreadyActive[invite.Invitee.GetLogin()] = struct{}{}
	}
	for _, u := range colabs {
		alreadyActive[u.GetLogin()] = struct{}{}
	}

	if _, ok := alreadyActive[username]; ok {
		return ErrorAlreadyCollaborator
	}

	_, _, err = c.client.Repositories.AddCollaborator(context.Background(), owner, name, username, nil)
	if err != nil {
		return fmt.Errorf("could not add %s to generated repository %s %w", username, repo, err)
	}

	return nil
}

func (c GithubClient) addCollaborator(login string, repoName string, username string) error {
	var i int
	var err error
	for i < 3 {
		_, _, err = c.client.Repositories.AddCollaborator(context.Background(), login, repoName, username, nil)
		if err == nil {
			return nil
		}

		log.Printf("could not add colaborator %s sleeping", err)
		time.Sleep(time.Second)
		i++
	}

	return fmt.Errorf("could not add %s to generated repository %s %w", username, repoName, err)
}

var space = regexp.MustCompile("/s+")

func makeRepoName(bName, username string, id int) string {
	rand.Seed(time.Now().UnixNano())

	return strings.ToLower(
		fmt.Sprintf(
			"%s-%s-test-%d",
			username,
			space.ReplaceAllString(bName, "-"),
			id,
		),
	)
}

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Upload uploads the test code into the assignment repository.
// It first clones the test github repo specified for the assignment into the tmp directory.
// After cloning, it bundles all the test files into a single commit. Signing it with a start test message.
// The Upload method expects that the repository provided in data have the correct permissions to be
// able to init a test. This means that the test repository needs to have access given to the github app
// as part of an installation. The assignment repository needs to be also created by the user whom
// c.accessToken stems from.
//
// Upload returns an error if there is any problem in execution of the upload. It cleans the temp directory
// of the cloned repository.
func (c GithubClient) Upload(data core.UploadDetails) error {
	i, err := c.newInstallation(data.InstallationID)
	if err != nil {
		return fmt.Errorf("failed to generate installation with id %d %w", data.InstallationID, err)
	}

	buf, err := i.DownloadRepo(context.Background(), data.TestVCSRepoURL)
	if err != nil {
		return fmt.Errorf("could not download repo contents %w", err)
	}

	zipPath := os.TempDir()
	clonePath := path.Join(zipPath, fmt.Sprintf("%d_%d", data.ID, time.Now().Unix()))
	err = os.MkdirAll(clonePath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not create repo clone dir %s %w", clonePath, err)
	}

	defer func() {
		err := removeContents(clonePath)
		if err != nil {
			log.Printf("could not remove clone path dir %s\n", err)
		}
	}()

	zipFile := path.Join(zipPath, fmt.Sprintf("%d.zip", time.Now().Unix()))
	out, err := os.Create(zipFile)
	if err != nil {
		return fmt.Errorf("could not create zipFile %s %w", zipFile, err)
	}
	defer out.Close()

	// Write the body to file
	_, err = buf.WriteTo(out)
	if err != nil {
		return fmt.Errorf("failed to write to zipFile %w", err)
	}

	r, err := git.PlainInit(clonePath, false)
	if err != nil {
		return fmt.Errorf("could not plain init dir %s %w", clonePath, err)
	}

	_, err = r.CreateRemote(&config.RemoteConfig{Name: "origin", URLs: []string{data.VCSRepoURL}})
	if err != nil {
		return fmt.Errorf("could not create remote %s %w", data.VCSRepoURL, err)
	}

	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("failed to init worktree %w", err)
	}

	_, err = unzip(out.Name(), clonePath)
	if err != nil {
		return fmt.Errorf("could not untip to %s %w", clonePath, err)
	}

	files, err := ioutil.ReadDir(clonePath)
	if err != nil {
		return fmt.Errorf("could not read clone dir files %w", err)
	}

	// copy the top level directory as the zipFile creates a throwaway dir
	var dirname string
	for _, file := range files {
		if file.IsDir() && file.Name() != ".git" {
			dirname = file.Name()
			break
		}
	}

	abs := path.Join(clonePath, dirname)
	err = copyDirectory(abs, clonePath)
	if err != nil {
		return fmt.Errorf("could not copy ziped dir %s to %s %w", abs, clonePath, err)
	}

	err = removeContents(abs)
	if err != nil {
		return fmt.Errorf("could not remove dir %s %w", abs, err)
	}

	_, err = w.Add(".")
	if err != nil {
		return fmt.Errorf("could not add all files %w", err)
	}

	commit, err := w.Commit("start test", &git.CommitOptions{
		Author: &object.Signature{
			Name:  c.intervConf.Username,
			Email: c.intervConf.Email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("could not commit changes %w", err)
	}

	_, err = r.CommitObject(commit)
	if err != nil {
		return fmt.Errorf("failed to get commit obj %w", err)
	}

	err = r.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth: &gitHttp.BasicAuth{
			Username: c.intervConf.Username,
			Password: c.intervConf.AccessToken,
		},
		Force: true,
	})
	if err != nil {
		return fmt.Errorf("could not push to remote %w", err)
	}

	return nil
}

func (c GithubClient) IsSubmitted(vcsURL, username string) (bool, error) {
	owner, name := getRepoName(vcsURL)
	prs, _, err := c.client.PullRequests.List(context.Background(), owner, name, nil)
	if err != nil {
		return false, fmt.Errorf("could not list prs %w", err)
	}

	for _, pr := range prs {
		if pr.GetUser().GetLogin() == username {
			return true, nil
		}
	}

	return false, nil
}

func (c GithubClient) Cleanup(details core.CleanDetails) error {
	owner, name := getRepoName(details.VCSRepoURL)
	_, err := c.client.Repositories.RemoveCollaborator(context.Background(), owner, name, details.CandidateUsername)
	if err != nil {
		return fmt.Errorf("could not remove collaborator from test repo %s %s %w", owner, name, err)
	}

	invites, _, err := c.client.Repositories.ListInvitations(context.Background(), owner, name, nil)
	if err != nil {
		return fmt.Errorf("could not list invitations for repo %s %s %s", owner, name, err)
	}

	for _, invite := range invites {
		if invite.GetInvitee().GetLogin() == details.CandidateUsername {
			_, err := c.client.Repositories.DeleteInvitation(context.Background(), owner, name, invite.GetID())
			if err != nil {
				return fmt.Errorf("could not remove collaborator invitation from test repo %s %s %w", owner, name, err)
			}
		}
	}

	for _, reviewer := range details.ReviewersUsernames {
		err := c.addCollaborator(owner, name, reviewer)
		if err != nil {
			return fmt.Errorf("could not add %s to repo %w", reviewer, err)
		}
	}

	return nil
}

func getRepoName(githubRepoURL string) (string, string) {
	pieces := strings.Split(githubRepoURL, "/")
	end := pieces[len(pieces)-1]

	return pieces[len(pieces)-2], strings.Replace(end, ".git", "", 1)
}

func removeContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return os.Remove(dir)
}

func copyDirectory(scrDir, dest string) error {
	entries, err := ioutil.ReadDir(scrDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		sourcePath := filepath.Join(scrDir, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
			return err
		}

		stat, ok := fileInfo.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("failed to get raw syscall.Stat_t data for '%s'", sourcePath)
		}

		switch fileInfo.Mode() & os.ModeType {
		case os.ModeDir:
			if err := createIfNotExists(destPath, 0755); err != nil {
				return err
			}
			if err := copyDirectory(sourcePath, destPath); err != nil {
				return err
			}
		case os.ModeSymlink:
			if err := copySymLink(sourcePath, destPath); err != nil {
				return err
			}
		default:
			if err := copy(sourcePath, destPath); err != nil {
				return err
			}
		}

		if err := os.Lchown(destPath, int(stat.Uid), int(stat.Gid)); err != nil {
			return err
		}

		isSymlink := entry.Mode()&os.ModeSymlink != 0
		if !isSymlink {
			if err := os.Chmod(destPath, entry.Mode()); err != nil {
				return err
			}
		}
	}
	return nil
}

func copy(srcFile, dstFile string) error {
	out, err := os.Create(dstFile)
	if err != nil {
		return err
	}

	defer out.Close()

	in, err := os.Open(srcFile)
	defer in.Close()
	if err != nil {
		return err
	}

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}

func exists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

func createIfNotExists(dir string, perm os.FileMode) error {
	if exists(dir) {
		return nil
	}

	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
	}

	return nil
}

func copySymLink(source, dest string) error {
	link, err := os.Readlink(source)
	if err != nil {
		return err
	}
	return os.Symlink(link, dest)
}

func unzip(src string, dest string) ([]string, error) {
	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}
