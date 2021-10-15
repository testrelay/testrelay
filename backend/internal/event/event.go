package event

import (
	"encoding/json"
	"time"
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

type assignment_status_enum string

func newStatus(s string) *assignment_status_enum {
	a := assignment_status_enum(s)
	return &a
}
