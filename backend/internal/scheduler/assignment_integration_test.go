//go:build integration
// +build integration

package scheduler_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/testrelay/testrelay/backend/internal/core/assignment"
	. "github.com/testrelay/testrelay/backend/internal/scheduler"
)

func TestHasuraAssignmentScheduler(t *testing.T) {
	type hdbScheduledEvent struct {
		ID            string     `db:"id"`
		ScheduledTime *time.Time `db:"scheduled_time"`
		Payload       string     `db:"payload"`
		WebhookConf   string     `db:"webhook_conf"`
	}

	db, err := sqlx.Connect("postgres", "user=postgres dbname=postgres password=postgrespassword sslmode=disable")
	require.NoError(t, err)
	defer db.Close()

	webhookURL := "http://webhook.url"
	h := NewHasuraAssignmentScheduler(
		"http://localhost:8080",
		"myadminsecretkey",
		webhookURL,
	)

	t.Run("Start", func(t *testing.T) {
		l, err := time.LoadLocation("Asia/Qatar")
		assert.NoError(t, err)

		type fakeData struct {
			Name string
		}

		now := time.Now().Add(time.Hour).In(l)
		id, err := h.Start(assignment.StartInput{
			Type:       "sent",
			ScheduleAt: now.Format(time.RFC3339),
			Data:       fakeData{Name: "testname"},
		})
		assert.NoError(t, err)
		defer func() {
			_, err := db.Exec("DELETE FROM hdb_catalog.hdb_scheduled_events WHERE id = $1", id)
			assert.NoError(t, err)
		}()

		var event hdbScheduledEvent
		err = db.Get(&event, "SELECT id, payload, scheduled_time, webhook_conf FROM hdb_catalog.hdb_scheduled_events WHERE id = $1", id)
		assert.NoError(t, err)

		assert.JSONEq(t, `{"data":{"Name":"testname"},"step":"sent"}`, event.Payload)
		assert.Equal(t, time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, now.Location()), event.ScheduledTime.In(l))
		assert.Equal(t, fmt.Sprintf("%q", webhookURL+"/assignments/process"), event.WebhookConf)
	})

	t.Run("Stop", func(t *testing.T) {
		l, err := time.LoadLocation("Asia/Qatar")
		assert.NoError(t, err)

		type fakeData struct {
			Name string
		}

		now := time.Now().Add(time.Hour).In(l)
		id, err := h.Start(assignment.StartInput{
			Type:       "sent",
			ScheduleAt: now.Format(time.RFC3339),
			Data:       fakeData{Name: "testname"},
		})
		assert.NoError(t, err)

		err = h.Stop(id)
		assert.NoError(t, err)

		var count int
		err = db.Get(&count, "SELECT count(id) FROM hdb_catalog.hdb_scheduled_events WHERE id = $1", id)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})
}
