package assignment

import (
	"errors"
	"fmt"

	"github.com/testrelay/testrelay/backend/internal/core"
	"github.com/testrelay/testrelay/backend/internal/core/business"
	"github.com/testrelay/testrelay/backend/internal/core/user"
)

type UserRepo interface {
	CreateUser(u *user.U) error
}

type BusinessRepo interface {
	GetTestBusiness(testID int) (business.Short, error)
}

type Repo interface {
	UpdateAssignmentToSent(a SentDetails) error
}

type AuthRepo interface {
	GetUserByEmail(email string) (user.AuthInfo, error)
	CreateUserFromAssignment(a Full) (user.AuthInfo, error)
	SetCustomUserClaims(claims user.AuthClaims) error
	GetPasswordResetLink(email, redirectLink string) (string, error)
}

type CandidateEmailData struct {
	EmailLink    string
	BusinessName string
	Assignment   Full
}

type Inviter struct {
	BusinessRepo   BusinessRepo
	Mailer         core.Mailer
	AssignmentRepo Repo
	UserRepo       UserRepo
	Auth           AuthRepo
	AppURL         string
}

func (i Inviter) Invite(data Full) error {
	b, err := i.BusinessRepo.GetTestBusiness(data.TestId)
	if err != nil {
		return fmt.Errorf("could not fetch business to invite user %w", err)
	}

	link := fmt.Sprintf("%s/assignments/%d/view", i.AppURL, data.Id)
	candidate, err := i.Auth.GetUserByEmail(data.CandidateEmail)
	if err == nil {
		err = i.Auth.SetCustomUserClaims(user.AuthClaims{
			AuthUID:     candidate.UID,
			BusinessIDs: []int64{int64(b.ID)},
		})
		if err != nil {
			return fmt.Errorf("could not update existing user claims %w", err)
		}
	}

	if err != nil && errors.Is(err, user.ErrorNotFound) {
		candidate, err = i.createUser(data, int64(b.ID))
		if err != nil {
			return fmt.Errorf("could not create user %w\n", err)
		}

		link, err = i.Auth.GetPasswordResetLink(candidate.Email, link)
		if err != nil {
			return fmt.Errorf("could not generate password reset link %w", err)
		}
	}

	if err != nil {
		return fmt.Errorf("error getting user from email for invite event %w\n", err)
	}

	err = i.Mailer.Send(core.MailConfig{
		TemplateName: "candidate-invite",
		Subject:      b.Name + " has invited you to a technical test",
		From:         "candidates",
		To:           candidate.Email,
	}, CandidateEmailData{
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

func (i Inviter) createUser(data Full, businessID int64) (user.AuthInfo, error) {
	candidate, err := i.Auth.CreateUserFromAssignment(data)
	if err != nil {
		return candidate, fmt.Errorf("could not create candidate user %w", err)
	}

	u := user.U{
		UID:   candidate.UID,
		Email: candidate.Email,
	}
	err = i.UserRepo.CreateUser(&u)
	if err != nil {
		return candidate, fmt.Errorf("could not create graphql candidate %w", err)
	}

	err = i.Auth.SetCustomUserClaims(user.AuthClaims{
		ID:           u.ID,
		AuthUID:      candidate.UID,
		Interviewing: []int64{businessID},
	})
	if err != nil {
		return candidate, fmt.Errorf("could not create custom claims for user %w", err)
	}

	return candidate, err
}
