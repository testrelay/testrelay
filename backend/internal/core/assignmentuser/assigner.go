package assignmentuser

import (
	"errors"
	"fmt"

	"github.com/testrelay/testrelay/backend/internal/core"
	"github.com/testrelay/testrelay/backend/internal/vcs"
)

type ReviewerRepository interface {
	GetReviewer(id int) (ReviewerDetail, error)
}

type Assigner struct {
	ReviewerRepository ReviewerRepository
	VCSClient          core.VCSCollaboratorAdder
	Mailer             core.Mailer
	APPURL             string
}

func (a Assigner) Assign(r RawReviewer) error {
	rd, err := a.ReviewerRepository.GetReviewer(r.ID)
	if err != nil {
		return fmt.Errorf("could not fetch reviewer id: %d err %w", r.ID, err)
	}

	if rd.Assignment.GithubRepoUrl != "" && rd.User.GithubUsername != "" {
		err := a.VCSClient.AddCollaborator(rd.Assignment.GithubRepoUrl, rd.User.GithubUsername)
		if err != nil {
			// return as nothing to do here
			if errors.Is(err, vcs.ErrorAlreadyCollaborator) {
				return nil
			}

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
	}, assignmentData{APPURL: a.APPURL, CandidateName: rd.Assignment.CandidateName})
	if err != nil {
		return fmt.Errorf("could not send reviewer invite to %s %w", rd.User.Email, err)
	}

	return nil
}

type assignmentData struct {
	CandidateName string
	APPURL        string
}
