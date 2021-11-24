package user

//go:generate mockgen -destination mocks/create.go -package mocks . AuthClient,Repo,Creator
import (
	"errors"
	"fmt"
)

// AuthClient defines an interface for a client to an authorization platform.
type AuthClient interface {
	GetUserByEmail(email string) (AuthInfo, error)
	CreateUser(name, email string) (AuthInfo, error)
	SetCustomUserClaims(claims AuthClaims) error
	GetPasswordResetLink(email, redirectLink string) (string, error)
}

type Repo interface {
	CreateUser(u *U) error
}

type Creator interface {
	FirstOrCreate(data CreateParams) (AuthInfo, error)
}

// AuthCreator orchestrates creating users in the testrelay system and a 3rd party auth system.
// It implements the Creator interface.
type AuthCreator struct {
	Auth AuthClient
	Repo Repo
}

// CreateParams defines input params for the FirstOrCreate method.
type CreateParams struct {
	// Name defines the name of the user to display. Can be blank.
	Name       string
	Email      string
	BusinessId int64

	// RedirectLink defines the full url of the link to redirect the user to after
	// a password reset - if the user is created as part of FirstOrCreate.
	RedirectLink string
	// Type is either oneof candidate|recruiter. If empty defaults to candidate.
	Type string
}

// FirstOrCreate creates a user in the testrelay system and links it to the provided BusinessId in data.
// If the user already exists Invite will simply update the claims on the user.
// If a new user is generated, a temporary password will be given and
//  a password reset link populated in AuthInfo.ResetLink.
func (c AuthCreator) FirstOrCreate(data CreateParams) (AuthInfo, error) {
	u, err := c.Auth.GetUserByEmail(data.Email)
	if err != nil {
		if !errors.Is(err, ErrorNotFound) {
			return u, fmt.Errorf("error getting user from email to invite %w\n", err)
		}

		u, err = c.createUser(data)
		if err != nil {
			return u, fmt.Errorf("could not create user for invite %w\n", err)
		}

		resetLink, err := c.Auth.GetPasswordResetLink(u.Email, data.RedirectLink)
		if err != nil {
			return u, fmt.Errorf("could not generate password reset link for invite %w", err)
		}

		u.ResetLink = resetLink
		return u, nil
	}

	interviewing := []int64{data.BusinessId}
	var businesses []int64
	if data.Type == "recruiter" {
		businesses = interviewing
		interviewing = nil
	}

	err = c.Auth.SetCustomUserClaims(AuthClaims{
		AuthUID:      u.UID,
		Interviewing: interviewing,
		BusinessIDs:  businesses,
	})
	if err != nil {
		return u, fmt.Errorf("could not update existing user claims %w", err)
	}

	return u, nil
}

func (c AuthCreator) createUser(data CreateParams) (AuthInfo, error) {
	au, err := c.Auth.CreateUser(data.Name, data.Email)
	if err != nil {
		return au, fmt.Errorf("could not create user %w", err)
	}

	u := U{
		UID:   au.UID,
		Email: au.Email,
	}
	err = c.Repo.CreateUser(&u)
	if err != nil {
		return au, fmt.Errorf("could not create graphql user %w", err)
	}

	interviewing := []int64{data.BusinessId}
	var businesses []int64
	if data.Type == "recruiter" {
		businesses = interviewing
		interviewing = nil
	}

	err = c.Auth.SetCustomUserClaims(AuthClaims{
		ID:           u.ID,
		AuthUID:      au.UID,
		Interviewing: interviewing,
		BusinessIDs:  businesses,
	})
	if err != nil {
		return au, fmt.Errorf("could not set custom claims for new user %w", err)
	}

	au.New = true
	return au, err
}
