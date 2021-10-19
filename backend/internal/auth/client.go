package auth

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/errorutils"

	"github.com/testrelay/testrelay/backend/internal/core/assignment"
	"github.com/testrelay/testrelay/backend/internal/core/user"
)

var (
	bracketReg = regexp.MustCompile(`{}`)
	commaReg   = regexp.MustCompile(`,`)
)

type FirebaseClient struct {
	Auth            *auth.Client
	CustomClaimName string
}

func (f FirebaseClient) GetUserByEmail(email string) (user.AuthInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	u, err := f.Auth.GetUserByEmail(ctx, email)
	if err != nil {
		if errorutils.IsNotFound(err) {
			return user.AuthInfo{}, user.ErrorNotFound
		}

		return user.AuthInfo{}, fmt.Errorf("error occured fetchin user %s from firebase %w", email, err)
	}

	return user.AuthInfo{
		UID:           u.UID,
		DisplayName:   u.DisplayName,
		Email:         u.Email,
		PhoneNumber:   u.PhoneNumber,
		PhotoURL:      u.PhotoURL,
		ProviderID:    u.ProviderID,
		CustomClaims:  u.CustomClaims,
		Disabled:      u.Disabled,
		EmailVerified: u.EmailVerified,
	}, nil
}

func (f FirebaseClient) CreateUserFromAssignment(a assignment.Full) (user.AuthInfo, error) {
	toCreate := &auth.UserToCreate{}
	toCreate.DisplayName(a.CandidateName).Email(a.CandidateEmail).Password(randomPassword(8))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	u, err := f.Auth.CreateUser(ctx, toCreate)
	if err != nil {
		return user.AuthInfo{}, fmt.Errorf("could not create user from assignment %w", err)
	}

	return user.AuthInfo{
		UID:           u.UID,
		DisplayName:   u.DisplayName,
		Email:         u.Email,
		PhoneNumber:   u.PhoneNumber,
		PhotoURL:      u.PhotoURL,
		ProviderID:    u.ProviderID,
		CustomClaims:  u.CustomClaims,
		Disabled:      u.Disabled,
		EmailVerified: u.EmailVerified,
	}, nil
}

func (f FirebaseClient) SetCustomUserClaims(claims user.AuthClaims) error {
	if claims.AuthUID == "" {
		return errors.New("auth id cannot be nil when setting claims")
	}

	rec, err := f.Auth.GetUser(context.Background(), claims.AuthUID)
	if err != nil {
		return fmt.Errorf("could not fetch user claims to refresh, for auth id %s %w", claims.AuthUID, err)
	}

	custom := map[string]interface{}{
		"x-hasura-allowed-roles": []string{"user", "candidate"},
		"x-hasura-default-role":  "user",
		"x-hasura-user-id":       claims.AuthUID,
	}

	existing := make(map[string]interface{})
	if v, ok := rec.CustomClaims[f.CustomClaimName]; ok {
		existing, _ = v.(map[string]interface{})
	}

	if v, ok := existing["x-hasura-user-pk"]; ok {
		claims.ID, err = strconv.ParseInt(v.(string), 10, 64)
		if err != nil {
			return fmt.Errorf("could not parse user claims existing pk %w", err)
		}
	}

	if claims.ID == 0 {
		return fmt.Errorf("user claims must contain a valid pk")
	}

	custom["x-hasura-user-pk"] = fmt.Sprintf("%d", claims.ID)

	businessIds := claims.BusinessIDs
	if v, ok := existing["x-hasura-business-ids"]; ok {
		businessIds = appendToExisting(v.(string), businessIds)
	}

	interviewingIds := claims.Interviewing
	if v, ok := existing["x-hasura-interviewing-ids"]; ok {
		interviewingIds = appendToExisting(v.(string), interviewingIds)
	}

	custom["x-hasura-business-ids"] = intSliceToString(businessIds)
	custom["x-hasura-interviewing-ids"] = intSliceToString(interviewingIds)

	err = f.Auth.SetCustomUserClaims(context.Background(), claims.AuthUID, map[string]interface{}{
		f.CustomClaimName: custom,
	})
	if err != nil {
		return fmt.Errorf("could not set custom claims %+v %w", custom, err)
	}

	return nil
}

func (f FirebaseClient) PasswordResetLink(email, redirectLink string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	link, err := f.Auth.PasswordResetLinkWithSettings(ctx, email, &auth.ActionCodeSettings{
		URL: redirectLink,
	})
	if err != nil {
		return "", fmt.Errorf("could not generate password reset for email %s %w", email, err)
	}

	return link, nil
}

func appendToExisting(existing string, toAdd []int64) []int64 {
	pieces := commaReg.Split(
		bracketReg.ReplaceAllString(existing, ""),
		-1,
	)

	var ids []int64
	for _, piece := range pieces {
		id, err := strconv.ParseInt(strings.TrimSpace(piece), 10, 64)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}

	lookup := make(map[int64]struct{})
	for _, id := range ids {
		lookup[id] = struct{}{}
	}

	for _, id := range toAdd {
		if _, ok := lookup[id]; !ok {
			ids = append(ids, id)
		}
	}

	return ids
}

func intSliceToString(ids []int64) string {
	s := "{"
	for _, id := range ids {
		s += fmt.Sprintf("%d,", id)
	}

	return strings.TrimRight(s, ",") + "}"
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomPassword(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
