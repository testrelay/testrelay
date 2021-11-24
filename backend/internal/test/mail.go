package test

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
)

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

func AssertEmail(t *testing.T, data MailhogQueryResponse, from string, subject string, expectedBody string) {
	actual := data.Items[0]
	assert.Equal(t, "<"+from+">", actual.Content.Headers.From[0])
	assert.Equal(t, subject, actual.Content.Headers.Subject[0])

	r := quotedprintable.NewReader(strings.NewReader(actual.Content.Body))
	b, _ := io.ReadAll(r)

	assert.Contains(t, strings.ReplaceAll(strings.ReplaceAll(string(b), "\r\n", ""), "\t", ""), expectedBody)
}

func GetEmail(t *testing.T, email string) MailhogQueryResponse {
	res, err := http.Get("http://localhost:8025/api/v2/search?kind=to&query=" + email)
	require.NoError(t, err)

	var data MailhogQueryResponse
	err = json.NewDecoder(res.Body).Decode(&data)
	require.NoError(t, err)

	assert.NotZero(t, data.Items)
	return data
}

func DeleteEmails(t *testing.T, data MailhogQueryResponse) {
	for _, item := range data.Items {
		req, err := http.NewRequest(http.MethodDelete, "http://localhost:8025/api/v1/messages/"+item.ID, nil)
		assert.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, res.StatusCode)
	}
}

