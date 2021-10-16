package graphql

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hasura/go-graphql-client"

	"github.com/testrelay/testrelay/backend/internal"
)

// KeyTransport adds a keyed header to the request
type KeyTransport struct {
	Key   string
	Value string
}

// RoundTrip implements the roundtripper interface adding a key value to the request.
func (t *KeyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add(t.Key, t.Value)

	return http.DefaultTransport.RoundTrip(req)
}

func NewClient(url string, token string) *HasuraClient {
	return &HasuraClient{
		client: graphql.NewClient(
			url,
			&http.Client{
				Transport: &KeyTransport{Key: "x-hasura-admin-secret", Value: token},
			},
		),
	}
}

type AssignmentUsers struct {
	Assignment ShortAssignment `graphql:"assignment" json:"assignment"`
	User       User            `graphql:"user" json:"user"`
}

type ShortAssignment struct {
	CandidateName graphql.String `graphql:"candidate_name" json:"candidate_name"`
	GithubRepoUrl graphql.String `graphql:"github_repo_url" json:"github_repo_url"`
}

type User struct {
	Email          graphql.String `graphql:"email" json:"email"`
	GithubUsername graphql.String `graphql:"github_username" json:"github_username"`
}

type Assignment struct {
	Status             graphql.String `graphql:"status" json:"status"`
	TestTimeChosen     graphql.String `graphql:"test_time_chosen" json:"test_time_chosen"`
	ChooseUntil        graphql.String `graphql:"choose_until" json:"choose_until"`
	TestDayChosen      graphql.String `graphql:"test_day_chosen" json:"test_day_chosen"`
	TestId             graphql.Int    `graphql:"test_id" json:"test_id"`
	TimeLimit          graphql.Int    `graphql:"time_limit" json:"time_limit"`
	CandidateId        graphql.Int    `graphql:"candidate_id" json:"candidate_id"`
	ID                 graphql.Int    `graphql:"id" json:"id"`
	CandidateName      graphql.String `graphql:"candidate_name" json:"candidate_name"`
	RecruiterId        graphql.Int    `graphql:"recruiter_id" json:"recruiter_id"`
	InviteCode         graphql.String `graphql:"invite_code" json:"invite_code"`
	GithubRepoUrl      graphql.String `graphql:"github_repo_url" json:"github_repo_url"`
	CandidateEmail     graphql.String `graphql:"candidate_email" json:"candidate_email"`
	TestTimezoneChosen graphql.String `graphql:"test_timezone_chosen" json:"test_timezone_chosen"`
	StepArn            graphql.String `graphql:"step_arn" json:"step_arn"`
	Candidate          Candidate      `graphql:"candidate" json:"candidate"`
	Recruiter          Recruiter      `graphql:"recruiter" json:"recruiter"`
	Test               Test           `graphql:"test" json:"test"`
}

type Test struct {
	Business   Business       `graphql:"business" json:"business"`
	Name       string         `graphql:"name" json:"name"`
	GithubRepo graphql.String `graphql:"github_repo" json:"github_repo"`
}

type Business struct {
	Name graphql.String `graphql:"name" json:"name"`
}

type Recruiter struct {
	Email graphql.String `graphql:"email" json:"email"`
}

type Candidate struct {
	Email             graphql.String `graphql:"email" json:"email"`
	GithubUsername    graphql.String `graphql:"github_username" json:"github_username"`
	GithubAccessToken graphql.String `graphql:"github_access_token" json:"github_access_token"`
}

type assignmentQ struct {
	AssignmentsByPK Assignment `graphql:"assignments_by_pk(id:$id)"`
}

type assignmentUQ struct {
	AssignmentUsersByPK AssignmentUsers `graphql:"assignment_users_by_pk(id:$id)"`
}

type assignmentMu struct {
	UpdateAssignmentsByPk struct {
		ID graphql.Int `graphql:"id"`
	} `graphql:"update_assignments_by_pk(pk_columns: {id: $id}, _set: {step_arn: $step_arn, github_repo_url: $github_repo_url})"`
}

type HasuraClient struct {
	client *graphql.Client
}

func (h HasuraClient) GetAssignmentUser(id int) (*AssignmentUsers, error) {
	var q assignmentUQ
	err := h.client.Query(context.Background(), &q, map[string]interface{}{
		"id": graphql.Int(id),
	})
	if err != nil {
		return nil, fmt.Errorf("could not fetch graphql assignment %w", err)
	}

	return &q.AssignmentUsersByPK, nil
}

func (h HasuraClient) GetAssignment(id int) (*Assignment, error) {
	var q assignmentQ
	err := h.client.Query(context.Background(), &q, map[string]interface{}{
		"id": graphql.Int(id),
	})
	if err != nil {
		return nil, fmt.Errorf("could not fetch graphql assignment %w", err)
	}

	return &q.AssignmentsByPK, nil
}

func (h HasuraClient) UpdateAssignmentWithDetails(id int, arn string, url string) error {
	var mu assignmentMu

	err := h.client.Mutate(context.Background(), &mu, map[string]interface{}{
		"id":              graphql.Int(id),
		"step_arn":        graphql.String(arn),
		"github_repo_url": graphql.String(url),
	})
	if err != nil {
		return fmt.Errorf("could not update graphql assignment %w", err)
	}

	return nil
}

func (h HasuraClient) Reviewers(id int) ([]internal.Reviewer, error) {
	var q AssignmentReviewers

	err := h.client.Mutate(context.Background(), &q, map[string]interface{}{
		"id": graphql.Int(id),
	})

	return q.AssignmentUsers.Reviewers, err
}

func (h HasuraClient) NewAssignmentEvent(userID int, assignmentID int, status string) error {
	var mu InsertAssignmentEvent
	return h.client.Mutate(context.Background(), &mu, map[string]interface{}{
		"user_id": graphql.Int(userID),
		"id":      graphql.Int(assignmentID),
		"status":  newStatus(status),
	})
}
