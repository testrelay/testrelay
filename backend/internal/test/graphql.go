package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"firebase.google.com/go/v4/auth"
	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/require"
)

type graphErrors []struct {
	Message   string
	Locations []struct {
		Line   int
		Column int
	}
}

func (e graphErrors) Error() string {
	b := strings.Builder{}
	for _, err := range e {
		b.WriteString(fmt.Sprintf("Message: %s, Locations: %+v", err.Message, err.Locations))
	}
	return b.String()
}

type GraphQLClient struct {
	Client  *http.Client
	BaseURL string
}

func (c GraphQLClient) Do(query string, variables map[string]interface{}, v interface{}) (string, error) {
	in := struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables,omitempty"`
	}{
		Query:     query,
		Variables: variables,
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(in)
	if err != nil {
		return "", err
	}

	resp, err := c.Client.Post(c.BaseURL, "application/json", &buf)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("non-200 OK status code: %v body: %q", resp.Status, body)
	}
	var out struct {
		Data   *json.RawMessage
		Errors graphErrors
	}

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &out)
	if err != nil {
		return "", fmt.Errorf("could not unmarshall body: %s %w", body, err)
	}

	if out.Data != nil && v != nil {
		err := json.Unmarshal(*out.Data, &v)
		if err != nil {
			return string(body), fmt.Errorf("could not unmarshall data: %s %w", body, err)
		}
	}

	if len(out.Errors) > 0 {
		b, _ := json.Marshal(out.Errors)
		return string(b), out.Errors
	}

	return string(body), nil
}

var insertUserWithBusiness = `
mutation ($auth_id: String!, $email: String!, $business_name: String!, $github_installation_id: String!) {
  insert_businesses_one(
	object: {
		name: $business_name, 
		github_installation_id: $github_installation_id,
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

type InsertUserWithBusinessMuData struct {
	Insert struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Creator struct {
			ID int `json:"id"`
		} `json:"creator"`
	} `json:"insert_businesses_one"`
}

type insertUserWithBusinessVars struct {
	AuthID               string `json:"auth_id" faker:"uuid_hyphenated"`
	Email                string `json:"email" faker:"-"`
	BusinessName         string `json:"business_name" faker:"username"`
	GithubInstallationID string `json:"github_installation_id" faker:"-"`
}

func (c GraphQLClient) CreateRecruiterAndBusiness(t *testing.T, recruiterUser *auth.UserRecord) (InsertUserWithBusinessMuData, func() error) {
	vb := insertUserWithBusinessVars{
		Email:                recruiterUser.Email,
		AuthID:               recruiterUser.UID,
		GithubInstallationID: os.Getenv("TEST_GITHUB_INSTALLATION"),
	}

	var res InsertUserWithBusinessMuData
	_, err := c.Do(insertUserWithBusiness, toQueryVars(t, &vb), &res)
	require.NoError(t, err)

	return res, func() error {
		_, err := c.Do(deleteBusinessMu, map[string]interface{}{
			"id":      res.Insert.ID,
			"user_id": res.Insert.Creator.ID,
		}, nil)
		return err
	}
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
