package user_test

import (
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/testrelay/testrelay/backend/internal/core/user"
	"github.com/testrelay/testrelay/backend/internal/core/user/mocks"
)

func TestAuthCreator(t *testing.T) {
	t.Run("FirstOrCreate", func(t *testing.T) {
		t.Run("with type recruiter", func(t *testing.T) {
			t.Run("if user not exists should create user ", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				auth := mocks.NewMockAuthClient(ctrl)
				repo := mocks.NewMockRepo(ctrl)

				c := user.AuthCreator{
					Auth: auth,
					Repo: repo,
				}

				params := user.CreateParams{
					Name:         faker.Name(),
					Email:        faker.Email(),
					BusinessId:   123,
					RedirectLink: faker.URL(),
					Type:         "recruiter",
				}

				auth.EXPECT().GetUserByEmail(params.Email).Return(user.AuthInfo{}, user.ErrorNotFound)

				info := user.AuthInfo{
					UID:   faker.UUIDHyphenated(),
					Email: params.Email,
				}
				var id int64 = 81123
				auth.EXPECT().CreateUser(params.Name, params.Email).Return(info, nil)
				repo.EXPECT().CreateUser(gomock.Any()).DoAndReturn(func(u *user.U) error {
					t.Helper()

					assert.Equal(t, info.UID, u.UID)
					assert.Equal(t, info.Email, u.Email)

					u.ID = id

					return nil
				})

				auth.EXPECT().SetCustomUserClaims(user.AuthClaims{
					ID:          id,
					AuthUID:     info.UID,
					BusinessIDs: []int64{params.BusinessId},
				})

				resetLink := faker.URL()
				auth.EXPECT().GetPasswordResetLink(params.Email, params.RedirectLink).Return(resetLink, nil)

				a, err := c.FirstOrCreate(params)
				require.NoError(t, err)

				info.New = true
				info.ResetLink = resetLink
				assert.Equal(t, info, a)
			})

			t.Run("if user exists should update claims", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				auth := mocks.NewMockAuthClient(ctrl)
				repo := mocks.NewMockRepo(ctrl)

				c := user.AuthCreator{
					Auth: auth,
					Repo: repo,
				}

				params := user.CreateParams{
					Name:         faker.Name(),
					Email:        faker.Email(),
					BusinessId:   123,
					RedirectLink: faker.URL(),
					Type:         "recruiter",
				}

				info := user.AuthInfo{
					UID: faker.UUIDHyphenated(),
				}
				auth.EXPECT().GetUserByEmail(params.Email).Return(info, nil)

				auth.EXPECT().SetCustomUserClaims(user.AuthClaims{
					AuthUID:     info.UID,
					BusinessIDs: []int64{params.BusinessId},
				})

				a, err := c.FirstOrCreate(params)
				require.NoError(t, err)
				assert.Equal(t, info, a)

			})
		})

		t.Run("with default type", func(t *testing.T) {
			t.Run("if user not exists should create user ", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				auth := mocks.NewMockAuthClient(ctrl)
				repo := mocks.NewMockRepo(ctrl)

				c := user.AuthCreator{
					Auth: auth,
					Repo: repo,
				}

				params := user.CreateParams{
					Name:         faker.Name(),
					Email:        faker.Email(),
					BusinessId:   123,
					RedirectLink: faker.URL(),
				}

				auth.EXPECT().GetUserByEmail(params.Email).Return(user.AuthInfo{}, user.ErrorNotFound)

				info := user.AuthInfo{
					UID:   faker.UUIDHyphenated(),
					Email: params.Email,
				}
				var id int64 = 81123
				auth.EXPECT().CreateUser(params.Name, params.Email).Return(info, nil)
				repo.EXPECT().CreateUser(gomock.Any()).DoAndReturn(func(u *user.U) error {
					t.Helper()

					assert.Equal(t, info.UID, u.UID)
					assert.Equal(t, info.Email, u.Email)

					u.ID = id

					return nil
				})

				auth.EXPECT().SetCustomUserClaims(user.AuthClaims{
					ID:           id,
					AuthUID:      info.UID,
					Interviewing: []int64{params.BusinessId},
				})

				resetLink := faker.URL()
				auth.EXPECT().GetPasswordResetLink(params.Email, params.RedirectLink).Return(resetLink, nil)

				a, err := c.FirstOrCreate(params)
				require.NoError(t, err)

				info.New = true
				info.ResetLink = resetLink
				assert.Equal(t, info, a)
			})

			t.Run("if user exists should update claims", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				auth := mocks.NewMockAuthClient(ctrl)
				repo := mocks.NewMockRepo(ctrl)

				c := user.AuthCreator{
					Auth: auth,
					Repo: repo,
				}

				params := user.CreateParams{
					Name:         faker.Name(),
					Email:        faker.Email(),
					BusinessId:   123,
					RedirectLink: faker.URL(),
				}

				info := user.AuthInfo{
					UID: faker.UUIDHyphenated(),
				}
				auth.EXPECT().GetUserByEmail(params.Email).Return(info, nil)

				auth.EXPECT().SetCustomUserClaims(user.AuthClaims{
					AuthUID:      info.UID,
					Interviewing: []int64{params.BusinessId},
				})

				a, err := c.FirstOrCreate(params)
				require.NoError(t, err)
				assert.Equal(t, info, a)

			})
		})
	})

}
