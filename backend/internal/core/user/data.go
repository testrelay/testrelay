package user

import "errors"

var (
	ErrorNotFound = errors.New("user not found")
)

type U struct {
	ID    int64
	UID   string
	Email string
}

type Short struct {
	Email          string
	GithubUsername string
}

type AuthClaims struct {
	ID           int64
	AuthUID      string
	BusinessIDs  []int64
	Interviewing []int64
}

type AuthInfo struct {
	UID           string
	DisplayName   string
	Email         string
	PhoneNumber   string
	PhotoURL      string
	ProviderID    string
	CustomClaims  map[string]interface{}
	Disabled      bool
	EmailVerified bool
}
