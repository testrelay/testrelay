package assignment

import (
	"fmt"

	"github.com/testrelay/testrelay/backend/internal/core"
	intTime "github.com/testrelay/testrelay/backend/internal/time"
)

type Fetcher interface {
	GetAssignment(id int) (WithTestDetails, error)
}

type UpdaterForScheduler interface {
	UpdateAssignmentWithDetails(id int, arn string, url string) error
}

type StartInput struct {
	Type string

	ID         int64       `json:"id"`
	ScheduleAt string      `json:"schedule_at"`
	Duration   int         `json:"duration"`
	Data       interface{} `json:"data"`
}

type SchedulerClient interface {
	Stop(id string) error
	Start(input StartInput) (string, error)
}

type Scheduler struct {
	Fetcher         Fetcher
	SchedulerClient SchedulerClient
	VCSCreator      core.VCSCreator
	Updater         UpdaterForScheduler
}

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
