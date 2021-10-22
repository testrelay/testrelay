package scheduler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/testrelay/testrelay/backend/internal/core/assignment"
	"github.com/testrelay/testrelay/backend/internal/httputil"
)

type HasuraScheduleResponse struct {
	Message string `json:"message"`
	EventId string `json:"event_id"`
}

type HasuraSchedulePayload struct {
	Type string             `json:"type"`
	Args HasuraScheduleData `json:"args"`
}

type HasuraScheduleData struct {
	Webhook    string      `json:"webhook"`
	ScheduleAt string      `json:"schedule_at"`
	Payload    interface{} `json:"payload"`
	Headers    []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"headers"`
	Comment string `json:"comment"`
}

type StepPayload struct {
	Step string      `json:"step"`
	Data interface{} `json:"data"`
}

type HasuraAssignmentScheduler struct {
	client     *http.Client
	WebhookURL string
	baseURL    string
}

func NewHasuraAssignmentScheduler(baseURL, token, webhookURL string) HasuraAssignmentScheduler {
	return HasuraAssignmentScheduler{
		client: &http.Client{
			Transport: &httputil.KeyTransport{
				Value: token,
				Key:   "x-hasura-admin-secret",
			},
			Timeout: time.Second * 5,
		},
		WebhookURL: webhookURL,
		baseURL:    baseURL,
	}
}

func (h HasuraAssignmentScheduler) Stop(id string) error {
	body := bytes.NewBufferString(fmt.Sprintf(`{
    "type": "delete",
    "args": {
      "table": {
        "name": "hdb_scheduled_events",
        "schema": "hdb_catalog",
      },
      "where": {
        "comment": {
          "$eq": %q,
        },
      },
    },
  }`, id))

	res, err := h.client.Post(
		h.baseURL+"/v1/metadata",
		"application/json",
		body,
	)
	if err != nil {
		return fmt.Errorf("could not post to %s %w", h.baseURL, err)
	}

	if res.StatusCode > 299 {
		rb, _ := io.ReadAll(res.Body)
		return fmt.Errorf("non 200 status code for delete schedule: %s", rb)
	}

	return nil
}

func (h HasuraAssignmentScheduler) Start(input assignment.StartInput) (string, error) {
	id := fmt.Sprintf("assignment-%d-%s", input.ID, input.Type)

	data := HasuraSchedulePayload{
		Type: "create_scheduled_event",
		Args: HasuraScheduleData{
			Webhook:    h.WebhookURL + "/assignments/process",
			ScheduleAt: input.ScheduleAt,
			Payload: StepPayload{
				Step: input.Type,
				Data: input.Data,
			},
			Comment: id,
		},
	}

	b, err := json.Marshal(data)
	if err != nil {
		return id, fmt.Errorf("could not marshal schedule event data %w", err)
	}

	buffer := bytes.NewBuffer(b)
	res, err := h.client.Post(
		h.baseURL+"/v1/metadata",
		"application/json",
		buffer,
	)
	if err != nil {
		return id, fmt.Errorf("could not post to %s %w", h.baseURL, err)
	}

	if res.StatusCode > 299 {
		rb, _ := io.ReadAll(res.Body)
		return id, fmt.Errorf("non 200 status code for schedule body: %s", rb)
	}

	// TODO one 2.1 hasura becomes stable we should use the event_id returned
	return id, nil
}
