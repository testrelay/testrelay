package graphql

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hasura/go-graphql-client"

	"github.com/testrelay/testrelay/backend/internal/core/assignment"
	"github.com/testrelay/testrelay/backend/internal/core/assignmentuser"
	"github.com/testrelay/testrelay/backend/internal/core/business"
	"github.com/testrelay/testrelay/backend/internal/core/user"
	"github.com/testrelay/testrelay/backend/internal/httputil"
)

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
	SchedulerID        graphql.String `graphql:"step_arn" json:"step_arn"`
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
	Name                 graphql.String `graphql:"name" json:"name"`
	GithubInstallationID graphql.String `graphql:"github_installation_id" json:"github_installation_id"`
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

// HasuraClient wraps a hasura graphql client with methods to CRUD data and return it in domain types.
type HasuraClient struct {
	client *graphql.Client
}

// NewHasuraClient returns a HasuraClient setting the base graphql client with a custom http.RoundTripper that
// sets the provided token as an Authorization header. The url provided for client interactions must be
// the full url to the root graphql query. e.g. /v1/graphql
func NewHasuraClient(url string, token string) *HasuraClient {
	return &HasuraClient{
		client: graphql.NewClient(
			url,
			&http.Client{
				Transport: &httputil.KeyTransport{Key: "x-hasura-admin-secret", Value: token},
			},
		),
	}
}

// GetBusiness returns a short business from the given businessID.
func (h HasuraClient) GetBusiness(businessID int64) (business.Short, error) {
	var q getBusinessQuery
	err := h.client.Query(context.Background(), &q, map[string]interface{}{
		"id": graphql.Int(businessID),
	})
	if err != nil {
		return business.Short{}, fmt.Errorf("couldn't retrieve business from id %d %w", businessID, err)
	}

	return business.Short{
		ID:                   int(q.Business.ID),
		Name:                 string(q.Business.Name),
		GithubInstallationID: string(q.Business.GithubInstallationID),
	}, nil
}

// LinkUser adds an entry to business_users with the given user type. LinkUser ignores the insert
// if there is already an entry with the same details.
func (h HasuraClient) LinkUser(userID, businessID int64, userType string) error {
	var mu linkUserMutation
	err := h.client.Mutate(context.Background(), &mu, map[string]interface{}{
		"user_id":     graphql.Int(userID),
		"business_id": graphql.Int(businessID),
		"user_type":   graphql.String(userType),
	})
	if err != nil {
		return fmt.Errorf("error linking user in hasura %w", err)
	}

	return nil
}

func (h HasuraClient) GetTestBusiness(testID int) (business.Short, error) {
	var q getBusinessForTestQuery
	err := h.client.Query(context.Background(), &q, map[string]interface{}{
		"test_id": graphql.Int(testID),
	})
	if err != nil {
		return business.Short{}, fmt.Errorf("couldn't retrieve business from test_id %d %s", testID, err)
	}

	return business.Short{
		ID:                   int(q.TestByPk.Business.ID),
		Name:                 string(q.TestByPk.Business.Name),
		GithubInstallationID: string(q.TestByPk.Business.GithubInstallationID),
	}, nil
}

func (h HasuraClient) GetReviewer(id int) (assignmentuser.ReviewerDetail, error) {
	var q assignmentUQ
	err := h.client.Query(context.Background(), &q, map[string]interface{}{
		"id": graphql.Int(id),
	})
	if err != nil {
		return assignmentuser.ReviewerDetail{}, fmt.Errorf("could not fetch graphql assignment %w", err)
	}

	return assignmentuser.ReviewerDetail{
		User: user.Short{
			Email:          string(q.AssignmentUsersByPK.User.Email),
			GithubUsername: string(q.AssignmentUsersByPK.User.GithubUsername),
		},
		Assignment: assignment.Short{
			CandidateName: string(q.AssignmentUsersByPK.Assignment.CandidateName),
			GithubRepoUrl: string(q.AssignmentUsersByPK.Assignment.GithubRepoUrl),
		},
	}, nil
}

func (h HasuraClient) GetAssignment(id int) (assignment.WithTestDetails, error) {
	var q assignmentQ
	err := h.client.Query(context.Background(), &q, map[string]interface{}{
		"id": graphql.Int(id),
	})
	if err != nil {
		return assignment.WithTestDetails{}, fmt.Errorf("could not fetch graphql assignment %w", err)
	}

	installationID, _ := strconv.ParseInt(string(q.AssignmentsByPK.Test.Business.GithubInstallationID), 10, 64)
	return assignment.WithTestDetails{
		Status:             string(q.AssignmentsByPK.Status),
		TestTimeChosen:     string(q.AssignmentsByPK.TestTimeChosen),
		ChooseUntil:        string(q.AssignmentsByPK.ChooseUntil),
		TestDayChosen:      string(q.AssignmentsByPK.TestDayChosen),
		TestID:             int(q.AssignmentsByPK.TestId),
		TimeLimit:          int(q.AssignmentsByPK.TimeLimit),
		CandidateID:        int(q.AssignmentsByPK.CandidateId),
		ID:                 int(q.AssignmentsByPK.ID),
		CandidateName:      string(q.AssignmentsByPK.CandidateName),
		RecruiterID:        int(q.AssignmentsByPK.RecruiterId),
		InviteCode:         string(q.AssignmentsByPK.InviteCode),
		GithubRepoURL:      string(q.AssignmentsByPK.GithubRepoUrl),
		CandidateEmail:     string(q.AssignmentsByPK.CandidateEmail),
		TestTimezoneChosen: string(q.AssignmentsByPK.TestTimezoneChosen),
		SchedulerID:        string(q.AssignmentsByPK.SchedulerID),
		Candidate: assignment.Candidate{
			Email:             string(q.AssignmentsByPK.Candidate.Email),
			GithubUsername:    string(q.AssignmentsByPK.Candidate.GithubUsername),
			GithubAccessToken: string(q.AssignmentsByPK.Candidate.GithubAccessToken),
		},
		Recruiter: assignment.Recruiter{
			Email: string(q.AssignmentsByPK.Recruiter.Email),
		},
		Test: assignment.Test{
			Business: assignment.Business{
				Name:                 string(q.AssignmentsByPK.Test.Business.Name),
				GithubInstallationID: installationID,
			},
			Name:       string(q.AssignmentsByPK.Test.Name),
			GithubRepo: string(q.AssignmentsByPK.Test.GithubRepo),
		},
	}, nil
}

func newInt(i int) *int {
	return &i
}

func newString(s string) *string {
	return &s
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

func (h HasuraClient) UpdateAssignmentToSent(a assignment.SentDetails) error {
	var q UserQuery
	err := h.client.Query(context.Background(), &q, map[string]interface{}{
		"user_id": graphql.String(a.CandidateUID),
	})
	if err != nil {
		return fmt.Errorf("error fetching user for uid %s %w", a.CandidateUID, err)
	}

	if len(q.Users) == 0 {
		return fmt.Errorf("could not find user for uid %s %w", a.CandidateUID, err)
	}

	var m UpdateAssignmentMutation
	err = h.client.Mutate(context.Background(), &m, map[string]interface{}{
		"id":           graphql.Int(a.ID),
		"status":       newStatus("sent"),
		"user_id":      graphql.Int(a.RecruiterID),
		"candidate_id": graphql.Int(q.Users[0].ID),
		"user_type":    graphql.String("candidate"),
		"business_id":  graphql.Int(a.BusinessID),
	})
	if err != nil {
		return fmt.Errorf("could not update candidate state to sent err: %w\n", err)
	}

	return nil
}

func (h HasuraClient) Reviewers(id int) ([]string, error) {
	var q AssignmentReviewers

	err := h.client.Query(context.Background(), &q, map[string]interface{}{
		"id": graphql.Int(id),
	})

	reviewers := make([]string, len(q.AssignmentUsers.Reviewers))
	for i, reviewer := range q.AssignmentUsers.Reviewers {
		reviewers[i] = reviewer.User.GithubUsername
	}

	return reviewers, err
}

func (h HasuraClient) CreateUser(u *user.U) error {
	var mu InsertUserMutation
	err := h.client.Mutate(context.Background(), &mu, map[string]interface{}{
		"auth_id": graphql.String(u.UID),
		"email":   graphql.String(u.Email),
	})

	if err != nil {
		return fmt.Errorf("could not create graph user %w", err)
	}

	u.ID = int64(mu.InsertUsersOne.ID)
	return nil
}

func (h HasuraClient) NewAssignmentEvent(userID int, assignmentID int, status string) error {
	var mu InsertAssignmentEvent
	return h.client.Mutate(context.Background(), &mu, map[string]interface{}{
		"user_id": graphql.Int(userID),
		"id":      graphql.Int(assignmentID),
		"status":  newStatus(status),
	})
}
