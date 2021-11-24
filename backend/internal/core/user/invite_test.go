package user_test

import (
	"fmt"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/testrelay/testrelay/backend/internal/core"
	"github.com/testrelay/testrelay/backend/internal/core/business"
	coreMocks "github.com/testrelay/testrelay/backend/internal/core/mocks"
	"github.com/testrelay/testrelay/backend/internal/core/user"
	"github.com/testrelay/testrelay/backend/internal/core/user/mocks"
)

func TestInviter(t *testing.T) {
	t.Run("Invite", func(t *testing.T) {
		t.Run("should create user and send mail using new template", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			fetcher := mocks.NewMockBusinessFetcher(ctrl)
			linker := mocks.NewMockBusinessLinker(ctrl)
			creator := mocks.NewMockCreator(ctrl)
			mailer := coreMocks.NewMockMailer(ctrl)

			i := user.Inviter{
				BusinessFetcher: fetcher,
				BusinessLinker:  linker,
				UserCreator:     creator,
				Mailer:          mailer,
			}

			email := faker.Email()
			link := faker.URL()
			var businessID int64 = 736
			var userID int64 = 66
			businessName := faker.Name()

			short := business.Short{
				Name: businessName,
				ID:   int(businessID),
			}
			fetcher.EXPECT().GetBusiness(businessID).Return(short, nil)

			authEmail := faker.Email()
			resetLink := faker.URL()
			ai := user.AuthInfo{
				Email: authEmail,
				CustomClaims: map[string]interface{}{
					user.CustomClaimKey: map[string]interface{}{
						"x-hasura-user-pk": fmt.Sprintf("%d", userID),
					},
				},
				ResetLink: resetLink,
				New:       true,
			}
			creator.EXPECT().FirstOrCreate(user.CreateParams{
				Email:        email,
				BusinessId:   int64(businessID),
				RedirectLink: link,
				Type:         "recruiter",
			}).Return(ai, nil)

			linker.EXPECT().LinkUser(userID, businessID, "recruiter").Return(nil)

			mailer.EXPECT().Send(core.MailConfig{
				TemplateName: "recruiter-invite-new",
				Subject:      "You've been invited to join " + short.Name + " on TestRelay",
				From:         "info",
				To:           authEmail,
			}, user.RecruiterInviteParams{
				Link:         resetLink,
				BusinessName: short.Name,
			})

			a, err := i.Invite(email, link, businessID)
			require.NoError(t, err)
			assert.Equal(t, &ai, a)
		})
	})
}
