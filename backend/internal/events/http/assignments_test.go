package http_test

import (
	"bytes"
	"fmt"
	http2 "net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/testrelay/testrelay/backend/internal/events/http"
	"github.com/testrelay/testrelay/backend/internal/events/http/mocks"
)

func TestAssignmentHandler(t *testing.T) {
	logger := zap.NewNop().Sugar()

	t.Run("EventHandler", func(t *testing.T) {
		t.Run("assignment_events", func(t *testing.T) {
			t.Run("cancelled", func(t *testing.T) {
				ctrl := gomock.NewController(t)
				s := mocks.NewMockAssignmentScheduler(ctrl)

				h := http.AssignmentHandler{
					Logger:    logger,
					Scheduler: s,
				}
				assignmentID := 47
				body := bytes.NewBuffer([]byte(fmt.Sprintf(`{
        "event": {
            "session_variables": {
                "x-hasura-role": "candidate",
                "x-hasura-user-pk": "157",
                "x-hasura-interviewing-ids": "{97}",
                "x-hasura-user-id": "IfUgofPYv2ZiT4jCDhzwC1c7E9h1",
                "x-hasura-business-ids": "{}"
            },
            "op": "INSERT",
            "data": {
                "old": null,
                "new": {
                    "event_type": "cancelled",
                    "assignment_id": %d,
                    "created_at": "2021-11-10T18:14:59.537331+00:00",
                    "id": 117,
                    "meta": {},
                    "user_id": 157
                }
            },
            "trace_context": {
                "trace_id": "ddfaf60f7b429ec0",
                "span_id": "9fd916050accf512"
            }
        },
        "created_at": "2021-11-10T18:14:59.537331Z",
        "id": "8b6f4978-6f99-43d8-b23a-8c6541abda9b",
        "delivery_info": {
            "max_retries": 0,
            "current_retry": 0
        },
        "trigger": {
            "name": "assignment_events"
        },
        "table": {
            "schema": "public",
            "name": "assignment_events"
        }
    }
}`, assignmentID)))

				w := httptest.NewRecorder()
				r := httptest.NewRequest("POST", "/", body)

				s.EXPECT().Stop(assignmentID).Return(nil)
				h.EventHandler(w, r)

				assert.Equal(t, http2.StatusOK, w.Code)
			})
		})
	})
}
