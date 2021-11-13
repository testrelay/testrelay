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

// HasuraAssignmentScheduler is hasura implementation of the assignment.SchedulerClient.
// It uses hasura's built in one off event scheduling to schedule messages to be delivered
// to the backend for future processing.
//
// You can read more about hasura one off events here:
// 		https://hasura.io/docs/latest/graphql/core/scheduled-triggers/create-one-off-scheduled-event.html
//
// and the metadata api that is used in this client here:
// 		https://hasura.io/docs/latest/graphql/core/api-reference/metadata-api/scheduled-triggers.html
type HasuraAssignmentScheduler struct {
	client      *http.Client
	webhookURL  string
	baseURL     string
	accessToken string
}

// NewHasuraAssignmentScheduler returns a  HasuraAssignmentScheduler with the underlying http
// transport inserting the necessary hasura auth headers.
func NewHasuraAssignmentScheduler(baseURL, hasuraToken, accessToken, webhookURL string) HasuraAssignmentScheduler {
	return HasuraAssignmentScheduler{
		client: &http.Client{
			Transport: &httputil.KeyTransport{
				Value: hasuraToken,
				Key:   "x-hasura-admin-secret",
			},
			Timeout: time.Second * 5,
		},
		webhookURL:  webhookURL,
		baseURL:     baseURL,
		accessToken: accessToken,
	}
}

// Stop cancels the scheduled_event in using the hasura metadata api.
// It returns an error if there is a non 200 status code.
func (h HasuraAssignmentScheduler) Stop(id string) error {
	body := bytes.NewBufferString(fmt.Sprintf(`{
    "type" : "delete_scheduled_event",
    "args" : {
        "type": "one_off",
        "event_id": %q
    }
}`, id))

	err := h.post(body, nil)
	if err != nil {
		return fmt.Errorf("could not delete scheduled event using hasura metadata api %w", err)
	}

	return nil
}

// Start schedules a one off event to be delivered to the application at a later date.
// It returns the id of the scheduled event which can be used to cancel the event.
func (h HasuraAssignmentScheduler) Start(input assignment.StartInput) (string, error) {
	data := hasuraSchedulePayload{
		Type: "create_scheduled_event",
		Args: hasuraScheduleData{
			Webhook:    h.webhookURL + "/assignments/process",
			ScheduleAt: input.ScheduleAt,
			Payload: stepPayload{
				Step: input.Type,
				Data: input.Data,
			},
			Headers: []header{
				{
					Name:  "Authorization",
					Value: h.accessToken,
				},
			},
		},
	}

	buf := bytes.NewBuffer([]byte{})
	err := json.NewEncoder(buf).Encode(data)
	if err != nil {
		return "", fmt.Errorf("could not marshal schedule event data %w", err)
	}

	var res hasuraScheduleResponse
	err = h.post(buf, &res)
	if err != nil {
		return "", fmt.Errorf("could not create scheduled event using hasura metadata api %w", err)
	}

	return res.EventId, nil
}

func (h HasuraAssignmentScheduler) post(body *bytes.Buffer, out interface{}) error {
	res, err := h.client.Post(h.baseURL+"/v1/metadata", "application/json", body)
	if err != nil {
		return fmt.Errorf("could send data to %s/v1/metadata %w", h.baseURL, err)
	}

	if res.StatusCode > 299 {
		rb, _ := io.ReadAll(res.Body)
		return fmt.Errorf("non 200 status code from POST schedule, code: %d body: %s", res.StatusCode, rb)
	}

	if out != nil {
		err := json.NewDecoder(res.Body).Decode(out)
		if err != nil {
			return fmt.Errorf("could not decode schedule response to out struct %w", err)
		}
	}

	return nil
}

type hasuraScheduleResponse struct {
	Message string `json:"message"`
	EventId string `json:"event_id"`
}

type hasuraSchedulePayload struct {
	Type string             `json:"type"`
	Args hasuraScheduleData `json:"args"`
}

type hasuraScheduleData struct {
	Webhook    string      `json:"webhook"`
	ScheduleAt string      `json:"schedule_at"`
	Payload    interface{} `json:"payload"`
	Headers    []header    `json:"headers"`
}

type header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type stepPayload struct {
	Step string      `json:"step"`
	Data interface{} `json:"data"`
}
