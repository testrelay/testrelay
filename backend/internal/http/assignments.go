package http

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	graphql2 "github.com/hasura/go-graphql-client"
	"go.uber.org/zap"

	"github.com/testrelay/testrelay/backend/internal"
	assignment2 "github.com/testrelay/testrelay/backend/internal/core/assignment"
	"github.com/testrelay/testrelay/backend/internal/scheduler"
	"github.com/testrelay/testrelay/backend/internal/store/graphql"
	intTime "github.com/testrelay/testrelay/backend/internal/time"
	"github.com/testrelay/testrelay/backend/internal/vcs"
)

type AssignmentHandler struct {
	HasuraClient *graphql.HasuraClient
	GithubClient *vcs.GithubClient
	Processor    assignment2.Inviter
	Logger       *zap.SugaredLogger
	Scheduler    scheduler.Scheduler
	Runner       assignment2.Runner
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

		BadRequest(w)
		return
	}

	switch data.Table.Name {
	case "assignments":
		if data.Event.Op == "INSERT" {
		}
		err = a.Processor.Process(data.Event)
		if err != nil {
			a.Logger.Error(
				"could not process event data",
				"error", err,
				"data", data,
			)

			BadRequest(w)
		}

	case "assignment_events":
		var body internal.AssignmentEvent
		if err := json.Unmarshal(data.Event.Data.New, &body); err != nil {
			BadRequest(w)
			return
		}

		if data.Event.Op == "INSERT" && body.EventType == "scheduled" {
			a.handleAssignmentScheduled(w, body)
			return
		}
	}

	Success(w)
}

func (a AssignmentHandler) handleAssignmentScheduled(w http.ResponseWriter, data internal.AssignmentEvent) {
	assignment, err := a.HasuraClient.GetAssignment(data.AssignmentID)
	if err != nil {
		a.Logger.Error(
			"could not retrieve assignment",
			"assignment", data.AssignmentID,
			"error", err,
		)
		BadRequest(w)
		return
	}

	if assignment.StepArn != "" {
		err := a.Scheduler.Stop(string(assignment.StepArn))
		if err != nil {
			a.Logger.Error(
				"could not stop assignment execution",
				"error", err,
			)

			BadRequest(w)
			return
		}
	}

	timeInput := intTime.AssignmentChoices{
		DayChosen:  string(assignment.TestDayChosen),
		TimeChosen: string(assignment.TestTimeChosen),
		Timezone:   string(assignment.TestTimezoneChosen),
	}
	t, err := intTime.Parse(timeInput)
	if err != nil {
		a.Logger.Error(
			"formating assignment time failed",
			"time_input", timeInput,
			"error", err,
		)

		BadRequest(w)
		return
	}

	githubRepoURL := string(assignment.GithubRepoUrl)
	if githubRepoURL == "" {
		githubRepoURL, err = a.GithubClient.CreateRepo(
			string(assignment.Test.Business.Name),
			string(assignment.Candidate.GithubUsername),
		)
		if err != nil {
			a.Logger.Error(
				"could not get github repo url",
				"business_name", assignment.Test.Business.Name,
				"github_username", assignment.Candidate.GithubUsername,
				"error", err,
			)

			BadRequest(w)
			return
		}
	}

	assignment.GithubRepoUrl = graphql2.String(githubRepoURL)
	startID, err := a.Scheduler.Start(scheduler.StartInput{
		ID:           int64(assignment.ID),
		TestStart:    t.SendNotificationAt,
		TestDuration: int(assignment.TimeLimit) - 600,
		Data:         assignment,
	})
	if err != nil {
		a.Logger.Error(
			"could not start assignment execution",
			"error", err,
		)

		BadRequest(w)
		return
	}

	err = a.HasuraClient.UpdateAssignmentWithDetails(int(assignment.ID), startID, githubRepoURL)
	if err != nil {
		a.Logger.Error(
			"could not update assignment after assignment trigger",
			"assignment", assignment.ID,
			"arn", startID,
			"repo_url", githubRepoURL,
			"error", err,
		)
		BadRequest(w)
		return
	}

	Success(w)
}

type StepPayload struct {
	Step string              `json:"step"`
	Data assignment2.RunData `json:"data"`
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

		BadRequest(w)
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

		BadRequest(w)
		return
	}

	Success(w)
}
