//go:build integration
// +build integration

package mail_test

import (
	"encoding/json"
	"io"
	"mime/quotedprintable"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/testrelay/testrelay/backend/internal/core"
	"github.com/testrelay/testrelay/backend/internal/mail"
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

type mailhogQueryResponse struct {
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

			data := getEmail(t, email)
			defer deleteEmails(t, data)

			expectedBody := `<h3>Hello test name,</h3><p>Thanks for submitting your assignment. test biz will now review your code and get back to you with feedback.</p>`
			assertEmail(t, data, from+domain, subject, expectedBody)
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

			data := getEmail(t, email)
			defer deleteEmails(t, data)

			expectedBody := `You've received an invite to join testee on TestRelay. Click the link <a href="http://mylink">here</a> to login into the testee dashboard.</p>`
			assertEmail(t, data, from+domain, subject, expectedBody)
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

			data := getEmail(t, email)
			defer deleteEmails(t, data)

			expectedBody := `<p>You've received an invite to join testee on TestRelay. Click the link <a href="http://mylink">here</a> to reset your password and login to the dashboard.</p>`
			assertEmail(t, data, from+domain, subject, expectedBody)
		})
	})
}

func assertEmail(t *testing.T, data mailhogQueryResponse, from string, subject string, expectedBody string) {
	actual := data.Items[0]
	assert.Equal(t, "<"+from+">", actual.Content.Headers.From[0])
	assert.Equal(t, subject, actual.Content.Headers.Subject[0])

	r := quotedprintable.NewReader(strings.NewReader(actual.Content.Body))
	b, _ := io.ReadAll(r)

	assert.Contains(t, strings.ReplaceAll(strings.ReplaceAll(string(b), "\r\n", ""), "\t", ""), expectedBody)
}

func getEmail(t *testing.T, email string) mailhogQueryResponse {
	res, err := http.Get("http://localhost:8025/api/v2/search?kind=to&query=" + email)
	require.NoError(t, err)

	var data mailhogQueryResponse
	err = json.NewDecoder(res.Body).Decode(&data)
	require.NoError(t, err)

	assert.NotZero(t, data.Items)
	return data
}

func deleteEmails(t *testing.T, data mailhogQueryResponse) {
	for _, item := range data.Items {
		req, err := http.NewRequest(http.MethodDelete, "http://localhost:8025/api/v1/messages/"+item.ID, nil)
		assert.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, res.StatusCode)
	}
}
