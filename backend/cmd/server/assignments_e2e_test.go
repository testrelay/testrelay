//go:build e2e
// +build e2e

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/bxcodec/faker/v3"
	"github.com/google/go-github/v39/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var insertAssignmentMutation = `
mutation (
	$email: String!,
	$name: String!,
	$choose_until: date!,
	$time_limit: Int!,
	$recruiter_id: Int!,
	$business_id: Int!,
	$test_github_repo: String!,
	$test_name: String!,
	$test_window: Int,
	$test_time_limit: Int!,
	$status: assignment_status_enum = sending
) { 
insert_assignments_one (
	object: {
		candidate_email: $email, 
		candidate_name: $name, 
		choose_until: $choose_until, 
		time_limit: $time_limit, 
		recruiter_id: $recruiter_id
		status: $status
		test: {
			data: {
				user_id: $recruiter_id,
				business_id: $business_id,
				github_repo: $test_github_repo, 
				name: $test_name, 
				test_window: $test_window, 
				time_limit: $test_time_limit
			}
		}
	}
) {
	id
	candidate_email
}
}`

var insertUserWithBusiness = `
mutation ($auth_id: String!, $email: String!, $business_name: String!) {
  insert_businesses_one(
	object: {
		name: $business_name, 
		creator: {
			data: {
				email: $email, 
				auth_id: $auth_id
			}
		}
	}
) {
    id
	name
    creator {
      id
    }
  }
}

`

var validateAssignmentQuery = `
query ($email: String!, $id: Int!) {
  users(where: {email: {_eq: $email}}) {
    id
    auth_id
    business_users {
      business_id
      user_type
    }
  }
  assignments_by_pk(id: $id) {
    status
    assignment_events {
      event_type
    }
  }
}
`

var deleteBusinessMu = `
mutation ($id: Int!, $user_id: Int!) {
  delete_businesses_by_pk(id: $id) {
    id
  }
  delete_users_by_pk(id: $user_id) {
    id
  }
}
`

var deleteUserMu = `
mutation ($id: Int!) {
  delete_users_by_pk(id: $id) {
    id
  }
}
`

type deleteUserVars struct {
	Id int `json:"id" faker:"-"`
}

type deleteBusinessVars struct {
	Id     int `json:"id" faker:"-"`
	UserId int `json:"user_id" faker:"-"`
}

type validateAssignmentVars struct {
	Email string `json:"email" faker:"-"`
	Id    int    `json:"id" faker:"-"`
}

type validateAssignmentQueryData struct {
	AssignmentsByPk struct {
		Status           string `json:"status"`
		AssignmentEvents []struct {
			EventType string `json:"event_type"`
		} `json:"assignment_events"`
	} `json:"assignments_by_pk"`
	Users []struct {
		Id            int    `json:"id"`
		AuthId        string `json:"auth_id"`
		BusinessUsers []struct {
			BusinessId int    `json:"business_id"`
			UserType   string `json:"user_type"`
		} `json:"business_users"`
	} `json:"users"`
}

type insertAssignmentMuData struct {
	Insert struct {
		ID             int    `json:"id"`
		CandidateEmail string `json:"candidate_email"`
	} `json:"insert_assignments_one"`
}

type insertUserWithBusinessMuData struct {
	Insert struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Creator struct {
			ID int `json:"id"`
		} `json:"creator"`
	} `json:"insert_businesses_one"`
}

type insertUserWithBusinessVars struct {
	AuthID       string `json:"auth_id" faker:"uuid_hyphenated"`
	Email        string `json:"email" faker:"-"`
	BusinessName string `json:"business_name" faker:"username"`
}

type insertAssignmentVars struct {
	RecruiterID    int    `json:"recruiter_id" faker:"-"`
	BusinessID     int    `json:"business_id" faker:"-"`
	Email          string `json:"email" faker:"-"`
	Name           string `json:"name" faker:"name"`
	ChooseUntil    string `json:"choose_until" faker:"date"`
	TimeLimit      int    `json:"time_limit" faker:"oneof: 14400, 28800, 129000"`
	TestGithubRepo string `json:"test_github_repo" faker:"-"`
	TestName       string `json:"test_name" faker:"username"`
	TestWindow     int    `json:"test_window" faker:"oneof: 100,200,600"`
	TestTimeLimit  int    `json:"test_time_limit" faker:"oneof: 100,200,600"`
}

type MailhogQueryResponse struct {
	Total int `json:"total"`
	Count int `json:"count"`
	Start int `json:"start"`
	Items []struct {
		ID   string `json:"ID"`
		From struct {
			Relays  interface{} `json:"Relays"`
			Mailbox string      `json:"Mailbox"`
			Domain  string      `json:"Domain"`
			Params  string      `json:"Params"`
		} `json:"From"`
		To []struct {
			Relays  interface{} `json:"Relays"`
			Mailbox string      `json:"Mailbox"`
			Domain  string      `json:"Domain"`
			Params  string      `json:"Params"`
		} `json:"To"`
		Content struct {
			Headers struct {
				ContentTransferEncoding []string `json:"Content-Transfer-Encoding"`
				ContentType             []string `json:"Content-Type"`
				Date                    []string `json:"Date"`
				From                    []string `json:"From"`
				MessageID               []string `json:"Message-ID"`
				MimeVersion             []string `json:"Mime-Version"`
				Received                []string `json:"Received"`
				ReturnPath              []string `json:"Return-Path"`
				Subject                 []string `json:"Subject"`
				To                      []string `json:"To"`
			} `json:"Headers"`
			Body string      `json:"Body"`
			Size int         `json:"Size"`
			MIME interface{} `json:"MIME"`
		} `json:"Content"`
		Created time.Time   `json:"Created"`
		MIME    interface{} `json:"MIME"`
		Raw     struct {
			From string   `json:"From"`
			To   []string `json:"To"`
			Data string   `json:"Data"`
			Helo string   `json:"Helo"`
		} `json:"Raw"`
	} `json:"items"`
}

func TestAssignments(t *testing.T) {
	t.Run("/events", func(t *testing.T) {
		t.Run("insert assignment event", func(t *testing.T) {
			tr := &testRunner{mu: &sync.Mutex{}, t: t}
			defer tr.clean()

			// setup
			testRepo := generateTestRepository(tr)
			fbRecruiter := createRecruiterFirebaseUser(tr)
			trBusinessWithUser := createRecruiterAndBusiness(tr, fbRecruiter)

			// insert assignment which triggers events
			res := insertAssignment(tr, testRepo, trBusinessWithUser)

			// assertions
			assertEmailSent(tr, res, trBusinessWithUser)
			cRec := assertCandidateCreatedInFirebase(tr, res)
			vad := assertAssignmentUpdated(tr, cRec, res, trBusinessWithUser)
			assertCandidateClaims(tr, trBusinessWithUser, cRec, vad)
		})

		t.Run("insert assignment_events event", func(t *testing.T) {
		})
	})

	t.Run("/process", func(t *testing.T) {
		t.Run("/init", func(t *testing.T) {

		})
	})
}

func assertCandidateClaims(tr *testRunner, trBusinessWithUser insertUserWithBusinessMuData, cRec *auth.UserRecord, vad validateAssignmentQueryData) bool {
	if len(vad.Users) == 0 {
		return false
	}

	return assert.Equal(tr.t, map[string]interface{}{
		"https://hasura.io/jwt/claims": map[string]interface{}{
			"x-hasura-allowed-roles":    []interface{}{"user", "candidate"},
			"x-hasura-business-ids":     "{}",
			"x-hasura-default-role":     "user",
			"x-hasura-interviewing-ids": fmt.Sprintf("{%d}", trBusinessWithUser.Insert.ID),
			"x-hasura-user-id":          cRec.UID,
			"x-hasura-user-pk":          fmt.Sprintf("%d", vad.Users[0].Id),
		},
	}, cRec.CustomClaims)
}

func assertCandidateCreatedInFirebase(tr *testRunner, res insertAssignmentMuData) *auth.UserRecord {
	cRec, err := firebaseClient.GetUserByEmail(context.Background(), res.Insert.CandidateEmail)
	require.NoError(tr.t, err, "firebase user not generated")

	tr.addCleanupStep(func() error {
		return firebaseClient.DeleteUser(context.Background(), cRec.UID)
	})

	return cRec
}

func assertEmailSent(tr *testRunner, res insertAssignmentMuData, trBusinessWithUser insertUserWithBusinessMuData) {
	qr, ok := waitForEmail(tr.t, res.Insert.CandidateEmail)
	if ok {
		assert.Equal(tr.t, "<candidates@testrelay.io>", qr.Items[0].Content.Headers.From[0])
		assert.Equal(tr.t, trBusinessWithUser.Insert.Name+" has invited you to a technical test", qr.Items[0].Content.Headers.Subject[0])
		assert.NotContains(tr.t, qr.Items[0].Content.Body, "{{")
	}

	tr.addCleanupStep(func() error {
		// TODO delete email
		return nil
	})
}

func assertAssignmentUpdated(tr *testRunner, cRec *auth.UserRecord, res insertAssignmentMuData, trBusinessWithUser insertUserWithBusinessMuData) validateAssignmentQueryData {
	vav := validateAssignmentVars{
		Email: strings.ToLower(res.Insert.CandidateEmail),
		Id:    res.Insert.ID,
	}
	var vad validateAssignmentQueryData

	actual, err := hasuraClient.do(validateAssignmentQuery, toQueryVars(tr.t, &vav), &vad)
	if assert.Len(tr.t, vad.Users, 1) {
		assert.NoError(tr.t, err)
		assert.JSONEq(
			tr.t,
			fmt.Sprintf(
				`{"data":{"users":[{"id":%d,"auth_id":"%s","business_users":[{"business_id":%d,"user_type":"candidate"}]}],"assignments_by_pk":{"status":"sent","assignment_events":[{"event_type":"sent"}]}}}`,
				vad.Users[0].Id,
				cRec.UID,
				trBusinessWithUser.Insert.ID,
			),
			actual)
	}

	tr.addCleanupStep(func() error {
		if len(vad.Users) > 0 {
			_, err := hasuraClient.do(deleteUserMu, toQueryVars(tr.t, &deleteUserVars{Id: vad.Users[0].Id}), nil)
			return err
		}

		return nil
	})

	return vad
}

func insertAssignment(tr *testRunner, testRepo *github.Repository, trBusinessWithUser insertUserWithBusinessMuData) insertAssignmentMuData {
	candidateEmail := faker.Email()
	v := insertAssignmentVars{
		TestGithubRepo: testRepo.GetCloneURL(),
		Email:          candidateEmail,
		BusinessID:     trBusinessWithUser.Insert.ID,
		RecruiterID:    trBusinessWithUser.Insert.Creator.ID,
	}

	var res insertAssignmentMuData
	_, err := hasuraClient.do(insertAssignmentMutation, toQueryVars(tr.t, &v), &res)
	require.NoError(tr.t, err)
	return res
}

func createRecruiterAndBusiness(tr *testRunner, recruiterUser *auth.UserRecord) insertUserWithBusinessMuData {
	vb := insertUserWithBusinessVars{
		Email:  recruiterUser.Email,
		AuthID: recruiterUser.UID,
	}

	var res insertUserWithBusinessMuData
	_, err := hasuraClient.do(insertUserWithBusiness, toQueryVars(tr.t, &vb), &res)
	require.NoError(tr.t, err)

	tr.addCleanupStep(func() error {
		vars := deleteBusinessVars{Id: res.Insert.ID, UserId: res.Insert.Creator.ID}
		_, err := hasuraClient.do(deleteBusinessMu, toQueryVars(tr.t, &vars), nil)
		return err
	})

	return res
}

func createRecruiterFirebaseUser(tr *testRunner) *auth.UserRecord {
	user := &auth.UserToCreate{}
	user.Email(faker.Email()).Password("mypassword1234").DisplayName(faker.Name())
	rec, err := firebaseClient.CreateUser(context.Background(), user)
	require.NoError(tr.t, err)

	tr.addCleanupStep(func() error {
		return firebaseClient.DeleteUser(context.Background(), rec.UID)
	})

	return rec
}

func generateTestRepository(tr *testRunner) *github.Repository {
	nowUnix := time.Now().Unix()
	repoName := fmt.Sprintf("%s-%d", strings.ToLower(faker.Username()), nowUnix)

	repo, _, err := githubClient.Repositories.CreateFromTemplate(context.Background(), "the-foreman", "test-template", &github.TemplateRepoRequest{
		Name:        github.String(repoName),
		Private:     github.Bool(true),
		Description: github.String(repoName + " generated by e2e test runner"),
	})
	require.NoError(tr.t, err)

	tr.addCleanupStep(func() error {
		_, err := githubClient.Repositories.Delete(context.Background(), repo.GetOwner().GetLogin(), repo.GetName())
		return err
	})

	return repo
}

func waitForEmail(t *testing.T, email string) (MailhogQueryResponse, bool) {
	t.Helper()

	for i := 0; i < 5; i++ {
		res, err := http.Get("http://localhost:8025/api/v2/search?kind=to&query=" + email)
		require.NoError(t, err)

		var data MailhogQueryResponse
		err = json.NewDecoder(res.Body).Decode(&data)
		require.NoError(t, err)

		if len(data.Items) == 0 {
			time.Sleep(time.Second)
			continue
		}

		return data, true
	}

	t.Errorf("no email for user %s found ", email)
	return MailhogQueryResponse{}, false
}

func toQueryVars(t *testing.T, a interface{}) map[string]interface{} {
	t.Helper()

	err := faker.FakeData(a)
	require.NoError(t, err)

	b, err := json.Marshal(a)
	require.NoError(t, err)

	var v map[string]interface{}
	err = json.Unmarshal(b, &v)
	require.NoError(t, err)

	return v
}
