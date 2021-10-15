package event

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"

	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/errorutils"
	hGraph "github.com/hasura/go-graphql-client"

	"github.com/testrelay/testrelay/backend/internal"
	"github.com/testrelay/testrelay/backend/internal/graphql"
	"github.com/testrelay/testrelay/backend/internal/mail"
)

type Processor interface {
	Process(event Event) error
}

type GraphqlClient interface {
	Query(ctx context.Context, q interface{}, variables map[string]interface{}) error
	Mutate(ctx context.Context, m interface{}, variables map[string]interface{}) error
}

type AWSProcessor struct {
	GraphqlClient GraphqlClient
	Mailer        mail.Mailer
	Auth          *auth.Client
	AppURL        string
}

func (p AWSProcessor) Process(event Event) error {
	bg := context.Background()

	var data internal.Assignment
	err := json.Unmarshal(event.Data.New, &data)
	if err != nil {
		return err
	}

	if event.Op == "INSERT" && data.Status == "sending" {
		return p.processSendingEvent(data, bg)
	}

	return nil
}

func (p AWSProcessor) processSendingEvent(data internal.Assignment, bg context.Context) error {
	var q graphql.BusinessQuery
	id := data.TestId
	err := p.GraphqlClient.Query(bg, &q, map[string]interface{}{
		"test_id": hGraph.Int(id),
	})
	if err != nil {
		return fmt.Errorf("couldn't retrieve business from id %d %s\n", id, err)
	}

	candidate, err := p.Auth.GetUserByEmail(bg, data.CandidateEmail)
	if err != nil && !errorutils.IsNotFound(err) {
		return fmt.Errorf("error getting user from email %w\n", err)
	}

	link := fmt.Sprintf("%s/assignments/%d/view", p.AppURL, data.Id)
	if candidate == nil {
		candidate, link, err = p.createUser(data, link, int(q.TestByPk.Business.ID))
		if err != nil {
			return fmt.Errorf("could not create user %w\n", err)
		}
	} else {
		// @todo update claims of user
	}

	err = p.Mailer.SendCandidateInviteEmail(mail.CandidateEmailData{
		Sender:       "candidates@testrelay.io",
		EmailLink:    link,
		BusinessName: string(q.TestByPk.Business.Name),
		Assignment:   data,
	})
	if err != nil {
		return fmt.Errorf("couldn't send candidate email %w\n", err)
	}

	var u graphql.UserQuery
	err = p.GraphqlClient.Query(bg, &u, map[string]interface{}{
		"user_id": hGraph.String(candidate.UID),
	})
	if err != nil {
		return fmt.Errorf("couldn't retrieve candidate for err %w\n", err)
	}

	if len(u.Users) == 0 {
		return fmt.Errorf("couldn't retrieve candidate for id %s none exist\n", candidate.UID)
	}

	var m graphql.UpdateAssignmentMutation
	err = p.GraphqlClient.Mutate(bg, &m, map[string]interface{}{
		"id":           hGraph.Int(data.Id),
		"status":       newStatus("sent"),
		"user_id":      hGraph.Int(data.RecruiterId),
		"candidate_id": hGraph.Int(u.Users[0].ID),
		"user_type":    hGraph.String("candidate"),
		"business_id":  hGraph.Int(q.TestByPk.Business.ID),
	})
	if err != nil {
		return fmt.Errorf("could not update candidate state to sent err: %w\n", err)
	}

	return nil
}

func (p AWSProcessor) createUser(data internal.Assignment, baseLink string, businessID int) (*auth.UserRecord, string, error) {
	bg := context.Background()

	user := auth.UserToCreate{}
	user.DisplayName(data.CandidateName).Email(data.CandidateEmail).Password(randomPassword(8))
	candidate, err := p.Auth.CreateUser(bg, &user)
	if err != nil {
		return nil, "", fmt.Errorf("could not create candidate user %w", err)
	}

	// create candidate
	var mu graphql.InsertUserMutation
	err = p.GraphqlClient.Mutate(bg, &mu, map[string]interface{}{
		"auth_id": hGraph.String(candidate.UID),
		"email":   hGraph.String(candidate.Email),
	})
	if err != nil {
		return nil, "", fmt.Errorf("could not create hasura candidate %w", err)
	}

	err = p.Auth.SetCustomUserClaims(bg, candidate.UID, map[string]interface{}{
		"https://hasura.io/jwt/claims": map[string]interface{}{
			"x-hasura-allowed-roles":    []string{"user", "candidate"},
			"x-hasura-default-role":     "user",
			"x-hasura-user-id":          candidate.UID,
			"x-hasura-user-pk":          fmt.Sprintf("%d", mu.InsertUsersOne.ID),
			"x-hasura-business-ids":     "{}",
			"x-hasura-interviewing-ids": fmt.Sprintf("{%d}", businessID),
		},
	})
	if err != nil {
		return nil, "", fmt.Errorf("could not create custom claims for user %w", err)
	}

	link, err := p.Auth.PasswordResetLinkWithSettings(bg, candidate.Email, &auth.ActionCodeSettings{
		URL: baseLink,
	})
	if err != nil {
		return nil, "", fmt.Errorf("could not generate password reset link %w", err)
	}

	return candidate, link, err
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomPassword(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
