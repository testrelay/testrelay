package api_test

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/testrelay/testrelay/backend/internal/api"
	"github.com/testrelay/testrelay/backend/internal/api/mocks"
	"github.com/testrelay/testrelay/backend/internal/core/user"
)

func TestUserResolver(t *testing.T) {
	t.Run("InviteUser", func(t *testing.T) {
		t.Run("should call invite with graphql params", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			inviter := mocks.NewMockInviter(ctrl)

			r := api.UserResolver{
				Inviter: inviter,
				Logger:  nil,
			}

			email := "test@test.com"
			link := "http://mylink"
			businessID := 12
			p := graphql.ResolveParams{
				Args: map[string]interface{}{
					"email":         email,
					"redirect_link": link,
					"business_id":   businessID,
				},
			}

			userID := 776
			info := &user.AuthInfo{
				UID:   uuid.New().String(),
				Email: email,
				CustomClaims: map[string]interface{}{
					user.CustomClaimKey: map[string]interface{}{
						"x-hasura-user-pk": fmt.Sprintf("%d", userID),
					},
				},
			}
			inviter.EXPECT().Invite(email, link, int64(businessID)).Return(info, nil)
			actual, err := r.InviteUser(p)
			require.NoError(t, err)

			expected := api.InviteUserResponse{
				ID:    int64(userID),
				Email: email,
			}
			assert.Equal(t, expected, actual)
		})
	})
}
