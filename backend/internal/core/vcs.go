package core

//go:generate mockgen -destination mocks/vcs.go -package mocks . VCSCollaboratorAdder,VCSUploader,VCSCleaner,VCSSubmissionChecker,VCSCreator

type UploadDetails struct {
	ID             int64
	VCSRepoURL     string
	TestVCSRepoURL string
	InstallationID int64
}

type CleanDetails struct {
	ID                 int64
	VCSRepoURL         string
	CandidateUsername  string
	ReviewersUsernames []string
}

type VCSCollaboratorAdder interface {
	AddCollaborator(repo string, username string) error
}

type VCSUploader interface {
	Upload(data UploadDetails) error
}

type VCSCleaner interface {
	Cleanup(details CleanDetails) error
}

type VCSSubmissionChecker interface {
	IsSubmitted(vcsURL, username string) (bool, error)
}

type VCSCreator interface {
	CreateRepo(businessName, username string, id int) (string, error)
}

type Repo struct {
	ID       int64  `json:"id"`
	FullName string `json:"full_name"`
}

type RepoCollector interface {
	CollectRepos(installationID int64) ([]Repo, error)
}
