package assignment

//go:generate mockgen -destination mocks/scheduler.go -package mocks . Fetcher,ScheduleUpdater,SchedulerClient
import (
	"fmt"

	"github.com/testrelay/testrelay/backend/internal/core"
	intTime "github.com/testrelay/testrelay/backend/internal/time"
)

// Fetcher defines an interface for a type that uses an id to retrieve a
// full assignment from underlying storage.
type Fetcher interface {
	GetAssignment(id int) (WithTestDetails, error)
}

// ScheduleUpdater defines an interface for a type that updates an assignment
// with the id for a future scheduled run.
type ScheduleUpdater interface {
	UpdateAssignmentWithDetails(id int, runID string, url string) error
}

// StartInput holds information needed to schedule an assignment in the future.
type StartInput struct {
	Type string

	ID         int64       `json:"id"`
	ScheduleAt string      `json:"schedule_at"`
	Duration   int         `json:"duration"`
	Data       interface{} `json:"data"`
}

// SchedulerClient defines a client that calls an entity that schedules
// assignments for a future date. See scheduler package for implementations.
type SchedulerClient interface {
	Stop(id string) error
	Start(input StartInput) (string, error)
}

// Scheduler orchestrates future assignment execution.
type Scheduler struct {
	Fetcher         Fetcher
	SchedulerClient SchedulerClient
	VCSCreator      core.VCSCreator
	Updater         ScheduleUpdater
}

// Stop terminates a previously started assignment using the assignmentID.
func (s Scheduler) Stop(assignmentID int) error {
	assignment, err := s.Fetcher.GetAssignment(assignmentID)
	if err != nil {
		return fmt.Errorf("could not fetch assignment id %d %w", assignmentID, err)
	}

	err = s.SchedulerClient.Stop(assignment.SchedulerID)
	if err != nil {
		return fmt.Errorf("could not stop previously scheduled assignment %w", err)
	}

	return nil
}

// Start schedules an assignment to execute at a date in the future.
func (s Scheduler) Start(assignmentID int) error {
	assignment, err := s.Fetcher.GetAssignment(assignmentID)
	if err != nil {
		return fmt.Errorf("could not fetch assignment id %d %w", assignmentID, err)
	}

	if assignment.SchedulerID != "" {
		err := s.SchedulerClient.Stop(assignment.SchedulerID)
		if err != nil {
			return fmt.Errorf("could not stop previously scheduled assignment %w", err)
		}
	}

	timeInput := intTime.AssignmentChoices{
		DayChosen:  assignment.TestDayChosen,
		TimeChosen: assignment.TestTimeChosen,
		Timezone:   assignment.TestTimezoneChosen,
	}
	t, err := intTime.Parse(timeInput)
	if err != nil {
		return fmt.Errorf("error formatting assignment schedule time %w", err)
	}

	githubRepoURL := assignment.GithubRepoURL
	if assignment.GithubRepoURL == "" {
		githubRepoURL, err = s.VCSCreator.CreateRepo(
			assignment.Test.Business.Name,
			assignment.Candidate.GithubUsername,
			assignment.ID,
		)
		if err != nil {
			return fmt.Errorf("could not generate repo for assignment %w", err)
		}
	}

	assignment.GithubRepoURL = githubRepoURL
	schedulerID, err := s.SchedulerClient.Start(StartInput{
		Type:       "start",
		ID:         int64(assignment.ID),
		ScheduleAt: t.SendNotificationAt,
		Duration:   int(assignment.TimeLimit) - 600,
		Data:       assignment,
	})
	if err != nil {
		return fmt.Errorf("could not schedule assignment to start %w", err)
	}

	err = s.Updater.UpdateAssignmentWithDetails(int(assignment.ID), schedulerID, githubRepoURL)
	if err != nil {
		return fmt.Errorf("could not update assignment with schedule details %w", err)
	}

	return nil
}
