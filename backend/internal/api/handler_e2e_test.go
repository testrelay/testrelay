//go:build e2e
// +build e2e

package api_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	firebase "firebase.google.com/go/v4"
	firebaseAuth "firebase.google.com/go/v4/auth"
	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/api/option"

	"github.com/testrelay/testrelay/backend/internal/api"
	"github.com/testrelay/testrelay/backend/internal/api/mocks"
	"github.com/testrelay/testrelay/backend/internal/auth"
	"github.com/testrelay/testrelay/backend/internal/core"
	"github.com/testrelay/testrelay/backend/internal/core/user"
	"github.com/testrelay/testrelay/backend/internal/mail"
	"github.com/testrelay/testrelay/backend/internal/store/graphql"
	"github.com/testrelay/testrelay/backend/internal/test"
)

var fetchUserQ = `
query ($email: String!) {
  users(where: {email: {_eq: $email}}) {
    id
  }
}
`

type usersRes struct {
	Users []struct {
		Id int `json:"id"`
	} `json:"users"`
}

func TestGraphQLQueryHandler(t *testing.T) {
	t.Run("ServeHTTP", func(t *testing.T) {
		t.Run("inviteUser", func(t *testing.T) {
			ctrl := gomock.NewController(t)

			hasuraClient := graphql.NewHasuraClient(os.Getenv("HASURA_URL")+"/v1/graphql", os.Getenv("HASURA_TOKEN"))
			firebaseClient := newFirebaseAuth(t)
			uCreator := user.AuthCreator{
				Auth: auth.FirebaseClient{
					Auth:            firebaseClient,
					CustomClaimName: "https://hasura.io/jwt/claims",
				},
				Repo: hasuraClient,
			}

			mailer := mail.NewSMTPMailer(core.SMTPConfig{
				SendingDomain: "@testrelay.io",
				Host:          os.Getenv("SMTP_HOST"),
				Port:          1025,
			})

			l, _ := zap.NewDevelopment()
			ur := api.UserResolver{
				Inviter: user.Inviter{
					BusinessFetcher: hasuraClient,
					BusinessLinker:  hasuraClient,
					UserCreator:     uCreator,
					Mailer:          mailer,
				},
				Logger: l.Sugar(),
			}

			verifier := mocks.NewMockVerifier(ctrl)
			h, err := api.NewGraphQLQueryHandler(
				os.Getenv("HASURA_URL")+"/v1/graphql",
				verifier,
				ur,
				&api.RepositoryResolver{},
			)
			require.NoError(t, err)

			client := test.NewGraphQLClientFromOS()
			d, clean := client.CreateRecruiterAndBusiness(t, &firebaseAuth.UserRecord{
				UserInfo: &firebaseAuth.UserInfo{
					Email: faker.Email(),
					UID:   uuid.New().String(),
				},
			})
			defer func() {
				err := clean()
				assert.NoError(t, err)
			}()

			w := httptest.NewRecorder()
			link := "https://app.testrelay.io"
			email := faker.Email()
			r, err := http.NewRequest(http.MethodPost, "/graphql", strings.NewReader(fmt.Sprintf(`{
	"query": "mutation {inviteUser(business_id:%d, email:\"%s\", redirect_link:\"%s\") { id }}"
}`, d.Insert.ID, email, link)))
			require.NoError(t, err)
			token := faker.UUIDHyphenated()
			r.Header.Set("Authorization", "Bearer "+token)

			verifier.EXPECT().Parse(token).Return(nil)
			h.ServeHTTP(w, r)

			u, err := firebaseClient.GetUserByEmail(context.Background(), email)
			assert.NoError(t, err)
			defer func() {
				err = firebaseClient.DeleteUser(context.Background(), u.UID)
				assert.NoError(t, err)
			}()

			var users usersRes
			_, err = client.Do(fetchUserQ, map[string]interface{}{
				"email": strings.ToLower(email),
			}, &users)
			assert.NoError(t, err)
			require.Len(t, users.Users, 1)
			defer func() {
				_, err = client.Do(`mutation ($email: String!) {
  delete_users(where: {email: {_eq: $email}}) {
    affected_rows
  }
}`,
					map[string]interface{}{
						"email": email,
					}, nil)
				assert.NoError(t, err)
			}()

			assert.Equal(t, map[string]interface{}{
				"https://hasura.io/jwt/claims": map[string]interface{}{
					"x-hasura-allowed-roles":    []interface{}{"user", "candidate"},
					"x-hasura-business-ids":     fmt.Sprintf("{%d}", d.Insert.ID),
					"x-hasura-default-role":     "user",
					"x-hasura-interviewing-ids": "{}",
					"x-hasura-user-id":          fmt.Sprintf("%s", u.UID),
					"x-hasura-user-pk":          fmt.Sprintf("%d", users.Users[0].Id),
				},
			}, u.CustomClaims)

			data := test.GetEmail(t, email)
			defer test.DeleteEmails(t, data)

			expectedBody := `<p>You've received an invite to join ` + d.Insert.Name + ` on TestRelay. Click the link <a href="https://testrelay-sandbox.firebaseapp.com`
			test.AssertEmail(t, data, "info@testrelay.io", "You've been invited to join "+d.Insert.Name+" on TestRelay", expectedBody)
		})
	})
}

func newFirebaseAuth(t *testing.T) *firebaseAuth.Client {
	t.Helper()
	app, err := firebase.NewApp(
		context.Background(),
		nil,
		option.WithCredentialsFile(os.Getenv("GOOGLE_SERVICE_ACC_LOCATION")),
	)
	require.NoError(t, err)

	a, err := app.Auth(context.Background())
	require.NoError(t, err)

	return a
}
