package assignment

import (
	"fmt"

	"github.com/testrelay/testrelay/backend/internal/core"
	"github.com/testrelay/testrelay/backend/internal/core/business"
	"github.com/testrelay/testrelay/backend/internal/core/user"
)

type BusinessRepo interface {
	GetTestBusiness(testID int) (business.Short, error)
}

type Repo interface {
	UpdateAssignmentToSent(a SentDetails) error
}

type UserCreator interface {
	FirstOrCreate(data user.CreateParams) (user.AuthInfo, error)
}

type candidateEmailData struct {
	EmailLink    string
	BusinessName string
	Assignment   Full
}

// Inviter invites users for a given assignment.
// It handles state management and user notification via email.
type Inviter struct {
	BusinessRepo   BusinessRepo
	Mailer         core.Mailer
	AssignmentRepo Repo
	UserCreator    UserCreator
	CandidatesURL  string
}

// Invite invites a user from the provided Full assignment. It uses
// the candidate email and name from the assignment data to generate a new user and link it to the parent business.
func (i Inviter) Invite(data Full) error {
	b, err := i.BusinessRepo.GetTestBusiness(data.TestId)
	if err != nil {
		return fmt.Errorf("could not fetch business to invite user %w", err)
	}

	link := fmt.Sprintf("%s/assignments/%d/view", i.CandidatesURL, data.Id)
	candidate, err := i.UserCreator.FirstOrCreate(user.CreateParams{
		Name:         data.CandidateName,
		Email:        data.CandidateEmail,
		BusinessId:   int64(b.ID),
		RedirectLink: link,
	})
	if err != nil {
		return fmt.Errorf("error inviting user from assignment create event %w\n", err)
	}

	err = i.Mailer.Send(core.MailConfig{
		TemplateName: "candidate-invite",
		Subject:      b.Name + " has invited you to a technical test",
		From:         "candidates",
		To:           candidate.Email,
	}, candidateEmailData{
		EmailLink:    link,
		BusinessName: b.Name,
		Assignment:   data,
	})
	if err != nil {
		return fmt.Errorf("couldn't send candidate email %w", err)
	}

	err = i.AssignmentRepo.UpdateAssignmentToSent(SentDetails{
		ID:           int64(data.Id),
		RecruiterID:  int64(data.RecruiterId),
		CandidateUID: candidate.UID,
		BusinessID:   int64(b.ID),
	})
	if err != nil {
		return fmt.Errorf("could not update assignment to sent status after invite %w", err)
	}

	return nil
}
