package user

//go:generate mockgen -destination mocks/invite.go -package mocks . BusinessFetcher,BusinessLinker
import (
	"fmt"

	"github.com/testrelay/testrelay/backend/internal/core"
	"github.com/testrelay/testrelay/backend/internal/core/business"
)

type BusinessFetcher interface {
	GetBusiness(businessID int64) (business.Short, error)
}

type BusinessLinker interface {
	LinkUser(userID, businessID int64, userType string) error
}

// Inviter takes care of inviting users to a businesses.
type Inviter struct {
	BusinessFetcher BusinessFetcher
	BusinessLinker  BusinessLinker
	UserCreator     Creator
	Mailer          core.Mailer
}

type RecruiterInviteParams struct {
	Link         string
	BusinessName string
}

// Invite invites a given email to a business. If the user email does not exist in the system
// Invite creates the user with a temp account. It generates a password reset link which is sent to
// the user via invite email.
func (i Inviter) Invite(email, redirectLink string, businessID int64) (*AuthInfo, error) {
	short, err := i.BusinessFetcher.GetBusiness(businessID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch business %d %w", businessID, err)
	}

	a, err := i.UserCreator.FirstOrCreate(CreateParams{
		Email:        email,
		BusinessId:   businessID,
		RedirectLink: redirectLink,
		Type: "recruiter",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user for invite request %w", err)
	}

	err = i.BusinessLinker.LinkUser(a.PK(), businessID, "recruiter")
	if err != nil {
		return nil, fmt.Errorf("failed to link user %d with business id %d %w", a.PK(), businessID, err)
	}

	template := "recruiter-invite"
	if a.New {
		template = "recruiter-invite-new"
	}

	err = i.Mailer.Send(core.MailConfig{
		TemplateName: template,
		Subject:      "You've been invited to join " + short.Name + " on TestRelay",
		From:         "info",
		To:           a.Email,
	}, RecruiterInviteParams{
		Link:         a.ResetLink,
		BusinessName: short.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("could not invite recruiter, email failed template: %s err: %w", template, err)
	}

	return &a, nil
}
