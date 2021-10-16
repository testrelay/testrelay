package assignmentuser

import (
	"fmt"

	"github.com/testrelay/testrelay/backend/internal/core"
)

type ReviewerRepository interface {
	GetReviewer(id int) (ReviewerDetail, error)
}

type VCSClient interface {
	AddCollaborator(repo string, username string) error
}

type Assigner struct {
	ReviewerRepository ReviewerRepository
	VCSClient          VCSClient
	Mailer             core.Mailer
}

func (a Assigner) Assign(r RawReviewer) error {
	rd, err := a.ReviewerRepository.GetReviewer(r.ID)
	if err != nil {
		return fmt.Errorf("could not fetch reviewer id: %d err %w", r.ID, err)
	}

	if rd.Assignment.GithubRepoUrl != "" && rd.User.GithubUsername != "" {
		err := a.VCSClient.AddCollaborator(rd.Assignment.GithubRepoUrl, rd.User.GithubUsername)
		if err != nil {
			return fmt.Errorf(
				"could not vcs collaborator: %s to repo: %s %w",
				rd.User.GithubUsername,
				rd.Assignment.GithubRepoUrl,
				err,
			)
		}
	}

	err = a.Mailer.Send(core.MailConfig{
		TemplateName: "reviewer-invite",
		Subject:      "You've been invited you to review " + rd.Assignment.CandidateName + "'s technical assignment",
		To:           rd.User.Email,
	}, nil)
	if err != nil {
		return fmt.Errorf("could not send reviewer invite to %s %w", rd.User.Email, err)
	}

	return nil
}
