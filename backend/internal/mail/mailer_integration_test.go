//go:build integration
// +build integration

package mail_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/testrelay/testrelay/backend/internal/core"
	"github.com/testrelay/testrelay/backend/internal/mail"
	"github.com/testrelay/testrelay/backend/internal/test"
)

type stubEmailSendingData struct {
	CandidateName string
	Test          struct {
		Business struct {
			Name string
		}
	}
}

type stubEmailRecruiterInviteData struct {
	Link         string
	BusinessName string
}

func TestSMTPMailer(t *testing.T) {
	t.Run("Send", func(t *testing.T) {
		domain := "@testdomain.io"
		mailer := mail.NewSMTPMailer(core.SMTPConfig{
			SendingDomain: domain,
			Host:          "localhost",
			Port:          1025,
		})

		email := "test@example.com"
		from := "testfrom"
		subject := "my funky subject"

		t.Run("submitted", func(t *testing.T) {
			err := mailer.Send(core.MailConfig{
				TemplateName: "submitted",
				Subject:      subject,
				From:         from,
				To:           email,
			}, stubEmailSendingData{
				CandidateName: "test name",
				Test: struct {
					Business struct {
						Name string
					}
				}{
					Business: struct {
						Name string
					}{
						Name: "test biz",
					},
				},
			})
			require.NoError(t, err)

			data := test.GetEmail(t, email)
			defer test.DeleteEmails(t, data)

			expectedBody := `<h3>Hello test name,</h3><p>Thanks for submitting your assignment. test biz will now review your code and get back to you with feedback.</p>`
			test.AssertEmail(t, data, from+domain, subject, expectedBody)
		})

		t.Run("recruiter-invite", func(t *testing.T) {
			err := mailer.Send(core.MailConfig{
				TemplateName: "recruiter-invite",
				Subject:      subject,
				From:         from,
				To:           email,
			}, stubEmailRecruiterInviteData{
				Link:         "http://mylink",
				BusinessName: "testee",
			})
			require.NoError(t, err)

			data := test.GetEmail(t, email)
			defer test.DeleteEmails(t, data)

			expectedBody := `You've received an invite to join testee on TestRelay. Click the link <a href="http://mylink">here</a> to login into the testee dashboard.</p>`
			test.AssertEmail(t, data, from+domain, subject, expectedBody)
		})

		t.Run("recruiter-invite-new", func(t *testing.T) {
			err := mailer.Send(core.MailConfig{
				TemplateName: "recruiter-invite-new",
				Subject:      subject,
				From:         from,
				To:           email,
			}, stubEmailRecruiterInviteData{
				Link:         "http://mylink",
				BusinessName: "testee",
			})
			require.NoError(t, err)

			data := test.GetEmail(t, email)
			defer test.DeleteEmails(t, data)

			expectedBody := `<p>You've received an invite to join testee on TestRelay. Click the link <a href="http://mylink">here</a> to reset your password and login to the dashboard.</p>`
			test.AssertEmail(t, data, from+domain, subject, expectedBody)
		})
	})
}
