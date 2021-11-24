//go:build integration
// +build integration

package graphql_test

import (
	"net/http"
	"os"
	"testing"

	"firebase.google.com/go/v4/auth"
	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/testrelay/testrelay/backend/internal/core/business"
	"github.com/testrelay/testrelay/backend/internal/httputil"
	"github.com/testrelay/testrelay/backend/internal/store/graphql"
	"github.com/testrelay/testrelay/backend/internal/test"
)

var fetchBusinessUsersQ = `
query ($user_id: Int!, $business_id: Int!, $user_type: String!) {
  business_users(where: {business_id: {_eq: $business_id}, user_id: {_eq: $user_id}, user_type: {_eq: $user_type}}) {
    id
  }
}
`

type bURes struct {
	BusinessUsers []struct {
		ID int64 `json:"id"`
	} `json:"business_users"`
}

var insertUserMu = `
mutation ($auth_id: String!, $email: String!) {
  insert_users_one(object: {auth_id: $auth_id, email: $email}) {
    id
  }
}
`

type uRes struct {
	InsertUsersOne struct {
		ID int64 `json:"id"`
	} `json:"insert_users_one"`
}

var deleteUserMu = `
mutation ($id: Int!) {
  delete_users_by_pk(id: $id) {
    id
  }
}
`

func TestHasuraClient(t *testing.T) {
	client := graphql.NewHasuraClient(
		os.Getenv("HASURA_URL")+"/v1/graphql",
		os.Getenv("HASURA_TOKEN"),
	)

	rawClient := test.GraphQLClient{
		BaseURL: os.Getenv("HASURA_URL") + "/v1/graphql",
		Client: &http.Client{
			Transport: &httputil.KeyTransport{Key: "x-hasura-admin-secret", Value: os.Getenv("HASURA_TOKEN")},
		},
	}

	t.Run("GetBusiness", func(t *testing.T) {
		d, cleanup := rawClient.CreateRecruiterAndBusiness(t, &auth.UserRecord{
			UserInfo: &auth.UserInfo{
				Email: faker.Email(),
				UID:   uuid.New().String(),
			},
		})
		defer func() {
			err := cleanup()
			assert.NoError(t, err)
		}()

		short, err := client.GetBusiness(d.Insert.ID)
		require.NoError(t, err)

		require.Equal(t, business.Short{
			Name: d.Insert.Name,
			ID:   d.Insert.ID,
		}, short)
	})

	t.Run("LinkUser", func(t *testing.T) {
		d, cleanup := rawClient.CreateRecruiterAndBusiness(t, &auth.UserRecord{
			UserInfo: &auth.UserInfo{
				Email: faker.Email(),
				UID:   uuid.New().String(),
			},
		})
		defer func() {
			err := cleanup()
			assert.NoError(t, err)
		}()

		var r uRes
		_, err := rawClient.Do(insertUserMu, map[string]interface{}{
			"email":   faker.Email(),
			"auth_id": faker.UUIDHyphenated(),
		}, &r)
		require.NoError(t, err)
		defer func() {
			_, err := rawClient.Do(deleteUserMu, map[string]interface{}{
				"id": r.InsertUsersOne.ID,
			}, nil)
			assert.NoError(t, err)
		}()

		err = client.LinkUser(r.InsertUsersOne.ID, int64(d.Insert.ID), "recruiter")
		require.NoError(t, err)

		var bu bURes
		_, err = rawClient.Do(fetchBusinessUsersQ, map[string]interface{}{
			"business_id": d.Insert.ID,
			"user_id":     r.InsertUsersOne.ID,
			"user_type":   "recruiter",
		}, &bu)
		require.NoError(t, err)

		assert.Len(t, bu.BusinessUsers, 1)
	})
}
