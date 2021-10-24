//go:build e2e
// +build e2e

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/bxcodec/faker/v3"
	"github.com/google/go-github/v39/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/testrelay/testrelay/backend/internal/core/assignment"
	http2 "github.com/testrelay/testrelay/backend/internal/events/http"
	"github.com/testrelay/testrelay/backend/internal/test"
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
			//var candidateRepo *github.Repository
			tr := test.NewRunner(t)
			defer tr.Clean()

			// setup
			testRepo := generateTestRepository(tr)
			fbRecruiter := createRecruiterFirebaseUser(tr)
			trBusinessWithUser := createRecruiterAndBusiness(tr, fbRecruiter)

			// insert assignment which triggers events
			assignmentInsertData := insertAssignment(tr, testRepo, trBusinessWithUser)

			// assertions
			assertEmailSent(tr, assignmentInsertData, trBusinessWithUser)
			cRec := assertCandidateCreatedInFirebase(tr, assignmentInsertData)
			vad := assertAssignmentUpdated(tr, cRec, assignmentInsertData, trBusinessWithUser)
			assertCandidateClaims(tr, trBusinessWithUser, cRec, vad)
			t.Run("insert assignment_events event", func(t *testing.T) {
				candidate := updateInsertedCandidateWithGithubUsername(tr, assignmentInsertData)
				now := time.Now()
				updateAssignmentWithTimeChosen(tr, assignmentInsertData, now.Format("15:04"), now.AddDate(0, 0, 4).Format("2006-01-02"))
				insertAssignmentEvent(tr, assignmentInsertData, candidate, "scheduled")

				assignmentDetails := waitForAssignmentDetails(tr, assignmentInsertData)
				assertGithubRepoCreated(tr, assignmentDetails, trBusinessWithUser)
				assertEventScheduled(tr, assignmentDetails)
			})

			t.Run("call process handler", func(t *testing.T) {
				fullAssignment, err := hasuraClient.GetAssignment(assignmentInsertData.Insert.ID)
				require.NoError(t, err)
				t.Run("start", func(t *testing.T) {
					step := "start"
					sendStepPayload(t, step, fullAssignment)
					assertWarningEmailSent(t, fullAssignment, tr, trBusinessWithUser)
				})

				t.Run("init", func(t *testing.T) {
					// call one off handler with send at now and payload init
					// check github uploaded
					// cehck assignment_events
				})

				t.Run("end", func(t *testing.T) {
					// call one off handler with send at now and payload end
					// check that email has been sent
				})

				t.Run("cleanup", func(t *testing.T) {
					// call one off handler with send at now and payload init
					// check github removed user
					// check emails sent to recruiter and candidate
					// check event created with missed
				})
			})
		})
	})
}

func assertWarningEmailSent(t *testing.T, fullAssignment assignment.WithTestDetails, tr *test.Runner, trBusinessWithUser insertUserWithBusinessMuData) {
	emails, ok := waitForEmail(t, fullAssignment.Candidate.Email)
	require.True(t, ok)
	assert.Equal(tr.T, "<candidates@testrelay.io>", emails.Items[0].Content.Headers.From[0])
	assert.Equal(tr.T, "5 minute reminder for your "+trBusinessWithUser.Insert.Name+" assignment", emails.Items[0].Content.Headers.Subject[0])
	assert.NotContains(tr.T, emails.Items[0].Content.Body, "{{")
	assert.Contains(tr.T, emails.Items[0].Content.Body, fullAssignment.GithubRepoURL)
}

func sendStepPayload(t *testing.T, step string, fullAssignment assignment.WithTestDetails) {
	body := newStepPayload(t, step, fullAssignment)
	res, err := http.Post("http://localhost:8000/assignments/process", "application/json", body)

	require.NoError(t, err)
	require.Equal(t, 200, res.StatusCode)
}

func newStepPayload(t *testing.T, step string, fullAssignment assignment.WithTestDetails) *bytes.Buffer {
	body := http2.StepPayload{
		Payload: struct {
			Data assignment.WithTestDetails `json:"data"`
			Step string                     `json:"step"`
		}{
			Data: fullAssignment,
			Step: step,
		},
	}
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(body)
	require.NoError(t, err)

	return buf
}

func assertEventScheduled(tr *test.Runner, details assignmentTestDetailsData) {
	// todo check hasura event scheduled
}

func assertGithubRepoCreated(tr *test.Runner, details assignmentTestDetailsData, user insertUserWithBusinessMuData) *github.Repository {
	r, _, err := githubClient.Repositories.Get(
		context.Background(),
		githubTestOwner,
		strings.ToLower(fmt.Sprintf("%s-%s-test-%d", testUserGithubUsername, user.Insert.Name, details.AssignmentsByPK.ID)),
	)
	require.NoError(tr.T, err)

	owner := r.GetOwner().GetLogin()
	repoName := r.GetName()
	tr.AddCleanupStep(func() error {
		_, err := githubClient.Repositories.Delete(context.Background(), owner, repoName)
		return err
	})

	invites, _, err := githubClient.Repositories.ListInvitations(context.Background(), owner, repoName, nil)
	require.NoError(tr.T, err)

	var found bool
	for _, u := range invites {
		if u.Invitee.GetLogin() == testUserGithubUsername {
			found = true
			break
		}
	}

	assert.True(tr.T, found, "could not find test user %s as invitee on generated repo %+v", testUserGithubUsername, invites)
	return r
}

func assertHasCommits(tr *test.Runner, owner string, repoName string) {
	commits, _, err := githubClient.Repositories.ListCommits(context.Background(), owner, repoName, nil)
	require.NoError(tr.T, err)
	require.Len(tr.T, commits, 1)

	c := commits[0].GetCommit()
	assert.Equal(tr.T, "start test", c.GetMessage())

	tree := c.GetTree()
	assert.Len(tr.T, tree.Entries, 2)

	filenames := make([]string, 0, 2)
	for _, e := range tree.Entries {
		if e != nil {
			filenames = append(filenames, e.GetPath())
		}
	}

	assert.Contains(tr.T, filenames, "test/index.txt")
	assert.Contains(tr.T, filenames, "test.txt")
}

var fetchAssignmentTestDetails = `
query ($id: Int!) {
  assignments_by_pk(id: $id) {
	id
    github_repo_url
    step_arn
  }
}
`

type assignmentTestDetailsData struct {
	AssignmentsByPK struct {
		ID            int    `json:"id"`
		GithubRepoURL string `json:"github_repo_url"`
		StepARN       string `json:"step_arn"`
	} `json:"assignments_by_pk"`
}

func assertAssignmentUpdatedWithRepoDetails(tr *test.Runner, res insertAssignmentMuData) assignmentTestDetailsData {
	var d assignmentTestDetailsData
	_, err := rawGraphlClient.Do(fetchAssignmentTestDetails, map[string]interface{}{
		"id": res.Insert.ID,
	}, &d)
	require.NoError(tr.T, err)

	return d
}

func waitForAssignmentDetails(tr *test.Runner, res insertAssignmentMuData) assignmentTestDetailsData {
	for i := 0; i < 5; i++ {
		d := assertAssignmentUpdatedWithRepoDetails(tr, res)
		if d.AssignmentsByPK.GithubRepoURL != "" && d.AssignmentsByPK.StepARN != "" {
			return d
		}

		time.Sleep(time.Second)
	}

	tr.T.Fatalf("could not find updated repo details for assignment %d", res.Insert.ID)
	return assignmentTestDetailsData{}
}

var updateAssignmentWithTimeMu = `
mutation ($id: Int!, $test_time_chosen: time!, $test_timezone_chosen: String!, $test_day_chosen: date!) {
  update_assignments_by_pk(pk_columns: {id: $id}, _set: {test_time_chosen: $test_time_chosen, test_timezone_chosen: $test_timezone_chosen, test_day_chosen: $test_day_chosen}) {
    id
  }
}
`

func updateAssignmentWithTimeChosen(tr *test.Runner, res insertAssignmentMuData, timeChosen string, dayChosen string) {
	_, err := rawGraphlClient.Do(updateAssignmentWithTimeMu, map[string]interface{}{
		"id":                   res.Insert.ID,
		"test_time_chosen":     timeChosen,
		"test_day_chosen":      dayChosen,
		"test_timezone_chosen": "Europe/London",
	}, nil)
	require.NoError(tr.T, err)
}

var insertAssignmentEventMu = `
mutation InsertAssignmentEvent($assignment_id: Int!, $user_id: Int!, $event_type: assignment_status_enum!) {
  insert_assignment_events_one(object: {assignment_id: $assignment_id, user_id: $user_id, event_type: $event_type}) {
    id
  }
}
`

func insertAssignmentEvent(tr *test.Runner, res insertAssignmentMuData, candidate userQueryData, eventType string) {
	_, err := rawGraphlClient.Do(insertAssignmentEventMu, map[string]interface{}{
		"assignment_id": res.Insert.ID,
		"user_id":       candidate.ID,
		"event_type":    eventType,
	}, nil)
	require.NoError(tr.T, err)
}

var updateUserWithGithubUsername = `
mutation ($email: String!, $github_username: String!) {
  update_users(where: {email: {_eq: $email}}, _set: {github_username: $github_username}) {
    returning {
      id
      email
      auth_id
    }
  }
}

`

type userUpdateData struct {
	UpdateUsers struct {
		Returning []userQueryData `json:"returning"`
	} `json:"update_users"`
}

type userQueryData struct {
	ID     int    `json:"id"`
	AuthID string `json:"auth_id"`
	Email  string `json:"email"`
}

func updateInsertedCandidateWithGithubUsername(tr *test.Runner, a insertAssignmentMuData) userQueryData {
	var d userUpdateData

	_, err := rawGraphlClient.Do(updateUserWithGithubUsername, map[string]interface{}{
		"email":           strings.ToLower(a.Insert.CandidateEmail),
		"github_username": strings.ToLower(testUserGithubUsername),
	}, &d)
	require.NoError(tr.T, err)
	require.Len(tr.T, d.UpdateUsers.Returning, 1)

	return d.UpdateUsers.Returning[0]
}

func assertCandidateClaims(tr *test.Runner, trBusinessWithUser insertUserWithBusinessMuData, cRec *auth.UserRecord, vad validateAssignmentQueryData) bool {
	if len(vad.Users) == 0 {
		return false
	}

	return assert.Equal(tr.T, map[string]interface{}{
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

func assertCandidateCreatedInFirebase(tr *test.Runner, res insertAssignmentMuData) *auth.UserRecord {
	cRec, err := firebaseClient.GetUserByEmail(context.Background(), res.Insert.CandidateEmail)
	require.NoError(tr.T, err, "firebase user not generated")

	tr.AddCleanupStep(func() error {
		return firebaseClient.DeleteUser(context.Background(), cRec.UID)
	})

	return cRec
}

func assertEmailSent(tr *test.Runner, res insertAssignmentMuData, trBusinessWithUser insertUserWithBusinessMuData) {
	qr, ok := waitForEmail(tr.T, res.Insert.CandidateEmail)
	if ok {
		assert.Equal(tr.T, "<candidates@testrelay.io>", qr.Items[0].Content.Headers.From[0])
		assert.Equal(tr.T, trBusinessWithUser.Insert.Name+" has invited you to a technical test", qr.Items[0].Content.Headers.Subject[0])
		assert.NotContains(tr.T, qr.Items[0].Content.Body, "{{")
	}
}

func assertAssignmentUpdated(tr *test.Runner, cRec *auth.UserRecord, res insertAssignmentMuData, trBusinessWithUser insertUserWithBusinessMuData) validateAssignmentQueryData {
	vav := validateAssignmentVars{
		Email: strings.ToLower(res.Insert.CandidateEmail),
		Id:    res.Insert.ID,
	}
	var vad validateAssignmentQueryData

	actual, err := rawGraphlClient.Do(validateAssignmentQuery, toQueryVars(tr.T, &vav), &vad)
	if assert.Len(tr.T, vad.Users, 1) {
		assert.NoError(tr.T, err)
		assert.JSONEq(
			tr.T,
			fmt.Sprintf(
				`{"data":{"users":[{"id":%d,"auth_id":"%s","business_users":[{"business_id":%d,"user_type":"candidate"}]}],"assignments_by_pk":{"status":"sent","assignment_events":[{"event_type":"sent"}]}}}`,
				vad.Users[0].Id,
				cRec.UID,
				trBusinessWithUser.Insert.ID,
			),
			actual)
	}

	tr.AddCleanupStep(func() error {
		if len(vad.Users) > 0 {
			_, err := rawGraphlClient.Do(deleteUserMu, toQueryVars(tr.T, &deleteUserVars{Id: vad.Users[0].Id}), nil)
			return err
		}

		return nil
	})

	return vad
}

func insertAssignment(tr *test.Runner, testRepo *github.Repository, trBusinessWithUser insertUserWithBusinessMuData) insertAssignmentMuData {
	candidateEmail := faker.Email()
	v := insertAssignmentVars{
		RecruiterID:    trBusinessWithUser.Insert.Creator.ID,
		BusinessID:     trBusinessWithUser.Insert.ID,
		Email:          candidateEmail,
		TestGithubRepo: testRepo.GetCloneURL(),
	}

	var res insertAssignmentMuData
	_, err := rawGraphlClient.Do(insertAssignmentMutation, toQueryVars(tr.T, &v), &res)
	require.NoError(tr.T, err)
	return res
}

func createRecruiterAndBusiness(tr *test.Runner, recruiterUser *auth.UserRecord) insertUserWithBusinessMuData {
	vb := insertUserWithBusinessVars{
		Email:  recruiterUser.Email,
		AuthID: recruiterUser.UID,
	}

	var res insertUserWithBusinessMuData
	_, err := rawGraphlClient.Do(insertUserWithBusiness, toQueryVars(tr.T, &vb), &res)
	require.NoError(tr.T, err)

	tr.AddCleanupStep(func() error {
		vars := deleteBusinessVars{Id: res.Insert.ID, UserId: res.Insert.Creator.ID}
		_, err := rawGraphlClient.Do(deleteBusinessMu, toQueryVars(tr.T, &vars), nil)
		return err
	})

	return res
}

func createRecruiterFirebaseUser(tr *test.Runner) *auth.UserRecord {
	user := &auth.UserToCreate{}
	user.Email(faker.Email()).Password("mypassword1234").DisplayName(faker.Name())
	rec, err := firebaseClient.CreateUser(context.Background(), user)
	require.NoError(tr.T, err)

	tr.AddCleanupStep(func() error {
		return firebaseClient.DeleteUser(context.Background(), rec.UID)
	})

	return rec
}

func generateTestRepository(tr *test.Runner) *github.Repository {
	nowUnix := time.Now().Unix()
	repoName := fmt.Sprintf("%s-%d", strings.ToLower(faker.Username()), nowUnix)

	repo, _, err := githubClient.Repositories.CreateFromTemplate(context.Background(), "the-foreman", "test-template", &github.TemplateRepoRequest{
		Name:        github.String(repoName),
		Private:     github.Bool(true),
		Description: github.String(repoName + " generated by e2e test runner"),
	})
	require.NoError(tr.T, err)

	tr.AddCleanupStep(func() error {
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

		req, err := http.NewRequest(http.MethodDelete, "http://localhost:8025/api/v1/messages", nil)
		assert.NoError(t, err)

		res, err = http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, res.StatusCode)

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
