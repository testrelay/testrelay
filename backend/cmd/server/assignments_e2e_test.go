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
			tr := &testRunner{mu: &sync.Mutex{}, t: t}
			defer tr.clean()

			// setup
			testRepo := generateTestRepository(tr)
			fbRecruiter := createRecruiterFirebaseUser(tr)
			trBusinessWithUser := createRecruiterAndBusiness(tr, fbRecruiter)

			res := insertAssignment(tr, testRepo, trBusinessWithUser)

			// assert the invite email sent which adds the candidate e.t.c
			assertEmailSent(tr, res, trBusinessWithUser)

			candidate := updateInsertedCandidateWithGithubUsername(tr, res)
			now := time.Now()
			updateAssignmentWithTimeChosen(tr, res, now.Format("15:04"), now.AddDate(0, 0, 4).Format("2006-01-02"))
			insertAssignmentEvent(tr, res, candidate, "scheduled")

			assignmentDetails := waitForAssignmentDetails(tr, res)
			assertGithubRepoCreated(tr, assignmentDetails, trBusinessWithUser)
			assertEventScheduled(tr, assignmentDetails)
		})
	})

	t.Run("/process", func(t *testing.T) {
		t.Run("/init", func(t *testing.T) {

		})
	})
}

func assertEventScheduled(tr *testRunner, details assignmentTestDetailsData) {
	// todo
}

func assertGithubRepoCreated(tr *testRunner, details assignmentTestDetailsData, user insertUserWithBusinessMuData) {
	r, _, err := githubClient.Repositories.Get(
		context.Background(),
		githubTestOwner,
		strings.ToLower(fmt.Sprintf("%s-%s-test-%d", testUserGithubUsername, user.Insert.Name, details.AssignmentsByPK.ID)),
	)
	require.NoError(tr.t, err)

	owner := r.GetOwner().GetLogin()
	repoName := r.GetName()
	tr.addCleanupStep(func() error {
		_, err := githubClient.Repositories.Delete(context.Background(), owner, repoName)
		return err
	})

	invites, _, err := githubClient.Repositories.ListInvitations(context.Background(), owner, repoName, nil)
	require.NoError(tr.t, err)

	var found bool
	for _, u := range invites {
		if u.Invitee.GetLogin() == testUserGithubUsername {
			found = true
			break
		}
	}

	assert.True(tr.t, found, "could not find test user %s as invitee on generated repo %+v", testUserGithubUsername, invites)
}

func assertHasCommits(tr *testRunner, owner string, repoName string) {
	commits, _, err := githubClient.Repositories.ListCommits(context.Background(), owner, repoName, nil)
	require.NoError(tr.t, err)
	require.Len(tr.t, commits, 1)

	c := commits[0].GetCommit()
	assert.Equal(tr.t, "start test", c.GetMessage())

	tree := c.GetTree()
	assert.Len(tr.t, tree.Entries, 2)

	filenames := make([]string, 0, 2)
	for _, e := range tree.Entries {
		if e != nil {
			filenames = append(filenames, e.GetPath())
		}
	}

	assert.Contains(tr.t, filenames, "test/index.txt")
	assert.Contains(tr.t, filenames, "test.txt")
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

func assertAssignmentUpdatedWithRepoDetails(tr *testRunner, res insertAssignmentMuData) assignmentTestDetailsData {
	var d assignmentTestDetailsData
	_, err := hasuraClient.do(fetchAssignmentTestDetails, map[string]interface{}{
		"id": res.Insert.ID,
	}, &d)
	require.NoError(tr.t, err)

	return d
}

func waitForAssignmentDetails(tr *testRunner, res insertAssignmentMuData) assignmentTestDetailsData {
	for i := 0; i < 5; i++ {
		d := assertAssignmentUpdatedWithRepoDetails(tr, res)
		if d.AssignmentsByPK.GithubRepoURL != "" && d.AssignmentsByPK.StepARN != "" {
			return d
		}

		time.Sleep(time.Second)
	}

	tr.t.Fatalf("could not find updated repo details for assignment %d", res.Insert.ID)
	return assignmentTestDetailsData{}
}

var updateAssignmentWithTimeMu = `
mutation ($id: Int!, $test_time_chosen: time!, $test_timezone_chosen: String!, $test_day_chosen: date!) {
  update_assignments_by_pk(pk_columns: {id: $id}, _set: {test_time_chosen: $test_time_chosen, test_timezone_chosen: $test_timezone_chosen, test_day_chosen: $test_day_chosen}) {
    id
  }
}
`

func updateAssignmentWithTimeChosen(tr *testRunner, res insertAssignmentMuData, timeChosen string, dayChosen string) {
	_, err := hasuraClient.do(updateAssignmentWithTimeMu, map[string]interface{}{
		"id":                   res.Insert.ID,
		"test_time_chosen":     timeChosen,
		"test_day_chosen":      dayChosen,
		"test_timezone_chosen": "Europe/London",
	}, nil)
	require.NoError(tr.t, err)
}

var insertAssignmentEventMu = `
mutation InsertAssignmentEvent($assignment_id: Int!, $user_id: Int!, $event_type: assignment_status_enum!) {
  insert_assignment_events_one(object: {assignment_id: $assignment_id, user_id: $user_id, event_type: $event_type}) {
    id
  }
}
`

func insertAssignmentEvent(tr *testRunner, res insertAssignmentMuData, candidate userQueryData, eventType string) {
	_, err := hasuraClient.do(insertAssignmentEventMu, map[string]interface{}{
		"assignment_id": res.Insert.ID,
		"user_id":       candidate.ID,
		"event_type":    eventType,
	}, nil)
	require.NoError(tr.t, err)
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

func updateInsertedCandidateWithGithubUsername(tr *testRunner, a insertAssignmentMuData) userQueryData {
	var d userUpdateData

	_, err := hasuraClient.do(updateUserWithGithubUsername, map[string]interface{}{
		"email":           strings.ToLower(a.Insert.CandidateEmail),
		"github_username": strings.ToLower(testUserGithubUsername),
	}, &d)
	require.NoError(tr.t, err)
	require.Len(tr.t, d.UpdateUsers.Returning, 1)

	return d.UpdateUsers.Returning[0]
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
		RecruiterID:    trBusinessWithUser.Insert.Creator.ID,
		BusinessID:     trBusinessWithUser.Insert.ID,
		Email:          candidateEmail,
		TestGithubRepo: testRepo.GetCloneURL(),
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
