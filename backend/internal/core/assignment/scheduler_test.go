package assignment_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/testrelay/testrelay/backend/internal/core/assignment"
	"github.com/testrelay/testrelay/backend/internal/core/assignment/mocks"
	coreMocks "github.com/testrelay/testrelay/backend/internal/core/mocks"
)

func TestScheduler(t *testing.T) {
	t.Run("Stop", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		f := mocks.NewMockFetcher(ctrl)
		su := mocks.NewMockScheduleUpdater(ctrl)
		sc := mocks.NewMockSchedulerClient(ctrl)
		vc := coreMocks.NewMockVCSCreator(ctrl)

		s := assignment.Scheduler{
			Fetcher:         f,
			SchedulerClient: sc,
			VCSCreator:      vc,
			Updater:         su,
		}

		assignmentID := 123
		schedulerID := "test-id"
		f.EXPECT().GetAssignment(assignmentID).Return(assignment.WithTestDetails{
			SchedulerID:        schedulerID,
		}, nil)

		sc.EXPECT().Stop(schedulerID).Return(nil)

		err := s.Stop(assignmentID)
		assert.NoError(t, err)
	})
}
