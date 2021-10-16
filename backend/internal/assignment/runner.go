package assignment

import (
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/testrelay/testrelay/backend/internal"
	"github.com/testrelay/testrelay/backend/internal/github"
	"github.com/testrelay/testrelay/backend/internal/graphql"
	"github.com/testrelay/testrelay/backend/internal/mail"
)

type RunData struct {
	ID           int                     `json:"id"`
	TestStart    time.Time               `json:"testStart"`
	TestDuration int                     `json:"testDuration"`
	Data         internal.FullAssignment `json:"data"`
}

type Runner struct {
	GHClient      *github.Client
	GraphQLClient *graphql.HasuraClient
	Mailer        mail.Mailer
	Logger        *zap.SugaredLogger
}

func (r Runner) Run(step string, data RunData) error {
	assignment := data.Data

	switch step {
	case "start":
		err := r.Mailer.Send(mail.Config{
			TemplateName: "warning",
			Subject:      "5 minute reminder for your " + assignment.Test.Business.Name + " assignment",
			From:         "candidates@testrelay.io",
			To:           assignment.CandidateEmail,
		}, assignment)
		if err != nil {
			return fmt.Errorf("could not send reminder email to candidate %s %w", assignment.CandidateEmail, err)
		}
	case "init":
		err := r.GHClient.Upload(assignment)
		if err != nil {
			return fmt.Errorf("could not upload assignment to github %w", err)
		}

		err = r.GraphQLClient.NewAssignmentEvent(assignment.CandidateID, assignment.ID, "inprogress")
		if err != nil {
			return fmt.Errorf("could not insert event 'inprogress' %w", err)
		}
	case "end":
		err := r.Mailer.Send(mail.Config{
			TemplateName: "end",
			Subject:      "Your test is about to finish",
			From:         "candidates@testrelay.io",
			To:           assignment.CandidateEmail,
		}, assignment)
		if err != nil {
			return fmt.Errorf("could not send finish email to candidate %s %w", assignment.CandidateEmail, err)
		}
	case "cleanup":
		reviewers, err := r.GraphQLClient.Reviewers(assignment.ID)
		if err != nil {
			return fmt.Errorf("could not get reviewers for assignemnt %d %w", assignment.ID, err)
		}

		err = r.GHClient.Cleanup(assignment, reviewers)
		if err != nil {
			return fmt.Errorf("could not cleanup github repo for assignemnt %d %w", assignment.ID, err)
		}

		ok, err := r.GHClient.IsSubmitted(assignment)
		if err != nil {
			return fmt.Errorf("could not check github repo is submitted assignemnt %d %w", assignment.ID, err)
		}

		status := "submitted"
		if !ok {
			status = "missed"
		}
		err = r.GraphQLClient.NewAssignmentEvent(assignment.CandidateID, assignment.ID, status)
		if err != nil {
			return fmt.Errorf("could not insert event '%s' %w", status, err)
		}

		err = r.Mailer.SendEnd(status, assignment)
		if err != nil {
			return fmt.Errorf("could not send end emails %w", err)
		}
	default:
		r.Logger.Info("assignment step does not exist", "step", step)
	}

	return nil
}
