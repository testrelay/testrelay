package http

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/testrelay/testrelay/backend/internal/core/assignment"
	"github.com/testrelay/testrelay/backend/internal/httputil"
	"github.com/testrelay/testrelay/backend/internal/store/graphql"
	"github.com/testrelay/testrelay/backend/internal/vcs"
)

type HasuraEvent struct {
	Event        Event        `json:"event"`
	CreatedAt    time.Time    `json:"created_at"`
	ID           string       `json:"id"`
	DeliveryInfo DeliveryInfo `json:"delivery_info"`
	Trigger      Trigger      `json:"trigger"`
	Table        Table        `json:"table"`
}

type SessionVariables struct {
	XHasuraBusinessID string `json:"x-hasura-business-id"`
	XHasuraRole       string `json:"x-hasura-role"`
	XHasuraUserPk     string `json:"x-hasura-user-pk"`
	XHasuraUserID     string `json:"x-hasura-user-id"`
}

type Data struct {
	Old json.RawMessage `json:"old"`
	New json.RawMessage `json:"new"`
}

type TraceContext struct {
	TraceID string `json:"trace_id"`
	SpanID  string `json:"span_id"`
}

type Event struct {
	SessionVariables SessionVariables `json:"session_variables"`
	Op               string           `json:"op"`
	Data             Data             `json:"data"`
	TraceContext     TraceContext     `json:"trace_context"`
}

type DeliveryInfo struct {
	MaxRetries   int `json:"max_retries"`
	CurrentRetry int `json:"current_retry"`
}

type Trigger struct {
	Name string `json:"name"`
}

type Table struct {
	Schema string `json:"schema"`
	Name   string `json:"name"`
}

type AssignmentEvent struct {
	ID           int             `json:"id"`
	UserID       int             `json:"user_id"`
	AssignmentID int             `json:"assignment_id"`
	Meta         json.RawMessage `json:"meta"`
	EventType    string          `json:"event_type"`
	CreatedAt    time.Time       `json:"created_at"`
}

type AssignmentHandler struct {
	HasuraClient *graphql.HasuraClient
	GithubClient *vcs.GithubClient
	Inviter      assignment.Inviter
	Logger       *zap.SugaredLogger
	Scheduler    assignment.Scheduler
	Runner       assignment.Runner
}

func (a AssignmentHandler) EventHandler(w http.ResponseWriter, r *http.Request) {
	var data HasuraEvent
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		body, _ := ioutil.ReadAll(r.Body)
		a.Logger.Error(
			"could not decode event data",
			"error", err,
			"body", body,
		)

		httputil.BadRequest(w)
		return
	}

	switch data.Table.Name {
	case "assignments":
		if data.Event.Op == "INSERT" {
			var body assignment.Full
			err := json.Unmarshal(data.Event.Data.New, &body)
			if err != nil {
				a.Logger.Error(
					"could not unmarshall assignments insert event data",
					"error", err,
					"body", string(data.Event.Data.New),
				)

				httputil.BadRequest(w)
				return
			}

			err = a.Inviter.Invite(body)
			if err != nil {
				a.Logger.Error(
					"could not process event data",
					"error", err,
					"data", data,
				)

				httputil.BadRequest(w)
			}
		}

	case "assignment_events":
		var body AssignmentEvent
		if err := json.Unmarshal(data.Event.Data.New, &body); err != nil {
			httputil.BadRequest(w)
			return
		}

		if data.Event.Op == "INSERT" && body.EventType == "scheduled" {
			err = a.Scheduler.Start(body.AssignmentID)
			if err != nil {
				a.Logger.Error(
					"could not start assignment",
					"assignment_id", body.AssignmentID,
					"error", err,
				)
				httputil.BadRequest(w)
				return
			}
		}
	}

	httputil.Success(w)
}

type StepPayload struct {
	Step string             `json:"step"`
	Data assignment.RunData `json:"data"`
}

func (a AssignmentHandler) ProcessHandler(w http.ResponseWriter, r *http.Request) {
	var data StepPayload
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		b, _ := io.ReadAll(r.Body)
		a.Logger.Error(
			"could not decode assignment process payload",
			"body", string(b),
			"error", err,
		)

		httputil.BadRequest(w)
		return
	}

	err = a.Runner.Run(data.Step, data.Data)
	if err != nil {
		a.Logger.Error(
			"run step errored",
			"step", data.Step,
			"data", data.Data,
			"error", err,
		)

		httputil.BadRequest(w)
		return
	}

	httputil.Success(w)
}
