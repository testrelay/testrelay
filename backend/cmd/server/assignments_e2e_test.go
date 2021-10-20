//go:build e2e
// +build e2e

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/bxcodec/faker/v3"
	"github.com/google/go-github/github"
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
		ID int `json:"id"`
	} `json:"insert_assignments_one"`
}

type insertUserWithBusinessMuData struct {
	Insert struct {
		ID      int `json:"id"`
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
	TimeLimit      int    `json:"time_limit" faker:"oneof: 600, 100, 200"`
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
			nowUnix := time.Now().Unix()
			repoName := fmt.Sprintf("%s-%d", strings.ToLower(faker.Username()), nowUnix)
			repo, _, err := githubClient.Repositories.Create(context.Background(), "", &github.Repository{
				Name:         github.String(repoName),
				Private:      github.Bool(true),
				Description:  github.String(repoName + " generated by e2e test runner"),
				MasterBranch: github.String("master"),
			})
			require.NoError(t, err)

			defer func() {
				_, err := githubClient.Repositories.Delete(context.Background(), repo.GetOwner().GetLogin(), repo.GetName())
				require.NoError(t, err)
			}()

			recruiterEmail := faker.Email()
			user := &auth.UserToCreate{}
			user.Email(recruiterEmail).Password("mypassword1234").DisplayName(faker.Name())
			rec, err := firebaseClient.CreateUser(context.Background(), user)
			require.NoError(t, err)
			defer func() {
				err := firebaseClient.DeleteUser(context.Background(), rec.UID)
				require.NoError(t, err)
			}()

			vb := insertUserWithBusinessVars{
				Email:  recruiterEmail,
				AuthID: rec.UID,
			}

			var uRes insertUserWithBusinessMuData
			vars := toQueryVars(t, &vb)
			_, err = hasuraClient.do(insertUserWithBusiness, vars, &uRes)
			require.NoError(t, err)

			candidateEmail := faker.Email()
			v := insertAssignmentVars{
				TestGithubRepo: repo.GetCloneURL(),
				Email:          candidateEmail,
				BusinessID:     uRes.Insert.ID,
				RecruiterID:    uRes.Insert.Creator.ID,
			}

			var res insertAssignmentMuData
			vars = toQueryVars(t, &v)
			_, err = hasuraClient.do(insertAssignmentMutation, vars, &res)
			require.NoError(t, err)

			qr, ok := waitForEmail(t, candidateEmail)
			if ok {
				assert.Equal(t, "<candidates@testrelay.io>", qr.Items[0].Content.Headers.From[0])
				assert.Equal(t, vb.BusinessName+" has invited you to a technical test", qr.Items[0].Content.Headers.Subject[0])
				assert.NotContains(t, qr.Items[0].Content.Body, "%")
			}

			cRec, err := firebaseClient.GetUserByEmail(context.Background(), candidateEmail)
			require.NoError(t, err, "firebase user not generated")

			vav := validateAssignmentVars{
				Email: cRec.Email,
				Id:    res.Insert.ID,
			}
			var vad validateAssignmentQueryData
			vars = toQueryVars(t, &vav)
			actual, err := hasuraClient.do(validateAssignmentQuery, vars, &vad)
			assert.NoError(t, err)
			assert.JSONEq(
				t,
				fmt.Sprintf(
					`{"data":{"users":[{"id":%d,"auth_id":"%s","business_users":[{"business_id":%d,"user_type":"candidate"}]}],"assignments_by_pk":{"status":"sent","assignment_events":[{"event_type":"sent"}]}}}`,
					vad.Users[0].Id,
					cRec.UID,
					uRes.Insert.ID,
				),
				actual)

			assert.Equal(t, map[string]interface{}{
				"https://hasura.io/jwt/claims": map[string]interface{}{
					"x-hasura-allowed-roles":    []interface{}{"user", "candidate"},
					"x-hasura-business-ids":     "{}",
					"x-hasura-default-role":     "user",
					"x-hasura-interviewing-ids": fmt.Sprintf("{%d}", uRes.Insert.ID),
					"x-hasura-user-id":          cRec.UID,
					"x-hasura-user-pk":          fmt.Sprintf("%d", vad.Users[0].Id),
				},
			}, cRec.CustomClaims)

			defer func() {
				err := firebaseClient.DeleteUser(context.Background(), cRec.UID)
				require.NoError(t, err)
			}()

			defer func() {
				// TODO delete the business & user associated with the test
				// make sure to alter hasura so that this cascades for all the added information
			}()
		})

		t.Run("insert assignment_events event", func(t *testing.T) {
		})
	})

	t.Run("/process", func(t *testing.T) {
		t.Run("/init", func(t *testing.T) {

		})
	})
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
