package vcs

import (
	"archive/zip"
	"bytes"
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
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"github.com/testrelay/testrelay/backend/internal/core"
)

var (
	letters = []rune("abcdefghijklmnopqrstuvwxyz")
)

type GithubClient struct {
	Client      *github.Client
	AccessToken string
}

func NewClient(accessToken string) *GithubClient {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)

	tc := oauth2.NewClient(context.Background(), ts)
	return &GithubClient{
		Client:      github.NewClient(tc),
		AccessToken: accessToken,
	}
}

func (c GithubClient) CreateRepo(bName, username string, id int) (string, error) {
	name := makeRepoName(bName, username, id)
	r := &github.Repository{
		Name:         github.String(name),
		Private:      github.Bool(true),
		Description:  github.String(username + " code assignment for " + bName),
		MasterBranch: github.String("master"),
	}

	log.Printf("creating repository %s\n", name)
	repo, _, err := c.Client.Repositories.Create(context.Background(), "", r)
	if err != nil {
		return "", fmt.Errorf("could not create repo %w", err)
	}

	log.Printf("repo generated: %+v\n", repo)
	log.Printf("repo ownder: %+v\n", repo.GetOwner())

	login := repo.GetOwner().GetLogin()
	repoName := repo.GetName()
	err = c.addCollaborator(login, repoName, username)
	if err != nil {
		return "", err
	}

	return repo.GetCloneURL(), nil
}

var (
	repl                     = regexp.MustCompile("https://github.com/")
	grepl                    = regexp.MustCompile("\\.git")

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
		colabs, _, err = c.Client.Repositories.ListCollaborators(context.Background(), owner, name, nil)
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

	invites, _, err := c.Client.Repositories.ListInvitations(context.Background(), owner, name, nil)
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

	_, err = c.Client.Repositories.AddCollaborator(context.Background(), owner, name, username, nil)
	if err != nil {
		return fmt.Errorf("could not add %s to generated repository %s %w", username, repo, err)
	}

	return nil
}

func (c GithubClient) addCollaborator(login string, repoName string, username string) error {
	var i int
	var err error
	for i < 3 {
		_, err = c.Client.Repositories.AddCollaborator(context.Background(), login, repoName, username, nil)
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

func (c GithubClient) Upload(data core.UploadDetails) error {
	id := data.ID
	from := data.TestVCSRepoURL
	to := data.VCSRepoURL

	owner, repo := getRepoName(from)
	u, _, err := c.Client.Repositories.GetArchiveLink(context.Background(), owner, repo, github.Zipball, nil)
	if err != nil {
		return fmt.Errorf("could not get archive link for repo %s %w", from, err)
	}

	req, _ := c.Client.NewRequest("GET", u.String(), nil)
	buf := bytes.NewBuffer([]byte{})
	_, err = c.Client.Do(context.Background(), req, buf)
	if err != nil {
		return fmt.Errorf("could not download zipFile %w", err)
	}

	// Create the file
	tmp := os.TempDir()
	zipPath := tmp
	clonePath := path.Join(zipPath, fmt.Sprintf("%d_%d", id, time.Now().Unix()))
	err = os.MkdirAll(clonePath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not create repo clone dir %s %w", clonePath, err)
	}
	defer func() {
		err := RemoveContents(clonePath)
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

	_, err = r.CreateRemote(&config.RemoteConfig{Name: "origin", URLs: []string{to}})
	if err != nil {
		return fmt.Errorf("could not create remote %s %w", to, err)
	}

	w, err := r.Worktree()
	if err != nil {
		return fmt.Errorf("failed to init worktree %w", err)
	}

	_, err = Unzip(out.Name(), clonePath)
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
	err = CopyDirectory(abs, clonePath)
	if err != nil {
		return fmt.Errorf("could not copy ziped dir %s to %s %w", abs, clonePath, err)
	}

	err = RemoveContents(abs)
	if err != nil {
		return fmt.Errorf("could not remove dir %s %w", abs, err)
	}

	_, err = w.Add(".")
	if err != nil {
		return fmt.Errorf("could not add all files %w", err)
	}

	commit, err := w.Commit("start test", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "testrelay",
			Email: "hugorut+2@gmail.com",
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
			Username: "test",
			Password: c.AccessToken,
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
	prs, _, err := c.Client.PullRequests.List(context.Background(), owner, name, nil)
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
	_, err := c.Client.Repositories.RemoveCollaborator(context.Background(), owner, name, details.CandidateUsername)
	if err != nil {
		return fmt.Errorf("could not remove collaborator from test repo %s %s %w", owner, name, err)
	}

	invites, _, err := c.Client.Repositories.ListInvitations(context.Background(), owner, name, nil)
	if err != nil {
		return fmt.Errorf("could not list invitations for repo %s %s %s", owner, name, err)
	}

	for _, invite := range invites {
		if invite.GetInvitee().GetLogin() == details.CandidateUsername {
			_, err := c.Client.Repositories.DeleteInvitation(context.Background(), owner, name, invite.GetID())
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

func RemoveContents(dir string) error {
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

func CopyDirectory(scrDir, dest string) error {
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
			if err := CreateIfNotExists(destPath, 0755); err != nil {
				return err
			}
			if err := CopyDirectory(sourcePath, destPath); err != nil {
				return err
			}
		case os.ModeSymlink:
			if err := CopySymLink(sourcePath, destPath); err != nil {
				return err
			}
		default:
			if err := Copy(sourcePath, destPath); err != nil {
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

func Copy(srcFile, dstFile string) error {
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

func Exists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

func CreateIfNotExists(dir string, perm os.FileMode) error {
	if Exists(dir) {
		return nil
	}

	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
	}

	return nil
}

func CopySymLink(source, dest string) error {
	link, err := os.Readlink(source)
	if err != nil {
		return err
	}
	return os.Symlink(link, dest)
}

func Unzip(src string, dest string) ([]string, error) {
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
