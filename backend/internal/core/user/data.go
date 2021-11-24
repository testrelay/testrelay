package user

import (
	"errors"
	"strconv"
)

var (
	ErrorNotFound = errors.New("user not found")

	CustomClaimKey = "https://hasura.io/jwt/claims"
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

	ResetLink string
	New       bool
}

// PK returns the user pk for the testrelay system from the custom claims.
// It returns a zero value if no pk is present or there is an error in the custom claims.
func (a AuthInfo) PK() int64 {
	existing := make(map[string]interface{})
	if v, ok := a.CustomClaims[CustomClaimKey]; ok {
		existing, _ = v.(map[string]interface{})
	}

	pk, ok := existing["x-hasura-user-pk"].(string)
	if !ok {
		return 0
	}

	i, _ := strconv.ParseInt(pk, 10, 64)
	return i
}
