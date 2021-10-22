package core

type UploadDetails struct {
	ID             int64
	VCSRepoURL     string
	TestVCSRepoURL string
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
