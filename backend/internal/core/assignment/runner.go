package assignment

import (
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/testrelay/testrelay/backend/internal/core"
)

type EventCreator interface {
	NewAssignmentEvent(userID int, assignmentID int, status string) error
}

type ReviewerCollector interface {
	Reviewers(assignmentID int) ([]string, error)
}
type RunData struct {
	Data WithTestDetails `json:"data"`
}

type Time func() time.Time

type Runner struct {
	Uploader          core.VCSUploader
	Cleaner           core.VCSCleaner
	SubmissionChecker core.VCSSubmissionChecker
	ReviewerCollector ReviewerCollector
	EventCreator      EventCreator
	Mailer            core.Mailer
	Logger            *zap.SugaredLogger
	SchedulerClient   SchedulerClient
	Time              Time

	StartDelay       time.Duration
	WarningBeforeEnd time.Duration
}

func (r Runner) Run(step string, data RunData) error {
	assignment := data.Data

	switch step {
	case "start":
		return r.start(assignment)
	case "init":
		return r.init(assignment)
	case "end":
		return r.end(assignment)
	case "cleanup":
		return r.cleanup(assignment)
	default:
		r.Logger.Info("assignment step does not exist", "step", step)
	}

	return nil
}

func (r Runner) cleanup(assignment WithTestDetails) error {
	reviewers, err := r.ReviewerCollector.Reviewers(assignment.ID)
	if err != nil {
		return fmt.Errorf("could not get reviewers for assignemnt %d %w", assignment.ID, err)
	}

	err = r.Cleaner.Cleanup(core.CleanDetails{
		ID:                 int64(assignment.ID),
		VCSRepoURL:         assignment.GithubRepoURL,
		CandidateUsername:  assignment.Candidate.GithubUsername,
		ReviewersUsernames: reviewers,
	})
	if err != nil {
		return fmt.Errorf("could not cleanup github repo for assignemnt %d %w", assignment.ID, err)
	}

	ok, err := r.SubmissionChecker.IsSubmitted(assignment.GithubRepoURL, assignment.Candidate.GithubUsername)
	if err != nil {
		return fmt.Errorf("could not check github repo is submitted assignemnt %d %w", assignment.ID, err)
	}

	status := "submitted"
	if !ok {
		status = "missed"
	}
	err = r.EventCreator.NewAssignmentEvent(assignment.CandidateID, assignment.ID, status)
	if err != nil {
		return fmt.Errorf("could not insert event '%s' %w", status, err)
	}

	err = r.sendEnd(status, assignment)
	if err != nil {
		return fmt.Errorf("could not send end emails %w", err)
	}
	return nil
}

func (r Runner) end(assignment WithTestDetails) error {
	err := r.Mailer.Send(core.MailConfig{
		TemplateName: "end",
		Subject:      "Your test is about to finish",
		From:         "candidates",
		To:           assignment.CandidateEmail,
	}, assignment)
	if err != nil {
		return fmt.Errorf("could not send finish email to candidate %s %w", assignment.CandidateEmail, err)
	}

	_, err = r.SchedulerClient.Start(StartInput{
		Type:       "cleanup",
		ID:         int64(assignment.ID),
		ScheduleAt: r.Time().Add(r.WarningBeforeEnd).Format(time.RFC3339),
		Data:       assignment,
	})
	if err != nil {
		return fmt.Errorf("could not schedule assignment to cleanup %w", err)
	}

	return nil
}

func (r Runner) init(assignment WithTestDetails) error {
	err := r.Uploader.Upload(core.UploadDetails{
		ID:             int64(assignment.ID),
		VCSRepoURL:     assignment.GithubRepoURL,
		TestVCSRepoURL: assignment.Test.GithubRepo,
	})
	if err != nil {
		return fmt.Errorf("could not upload assignment to github %w", err)
	}

	err = r.EventCreator.NewAssignmentEvent(assignment.CandidateID, assignment.ID, "inprogress")
	if err != nil {
		return fmt.Errorf("could not insert event 'inprogress' %w", err)
	}

	_, err = r.SchedulerClient.Start(StartInput{
		Type:       "end",
		ID:         int64(assignment.ID),
		ScheduleAt: r.Time().Add(-r.WarningBeforeEnd).Add(time.Second * time.Duration(assignment.TimeLimit)).Format(time.RFC3339),
		Data:       assignment,
	})
	if err != nil {
		return fmt.Errorf("could not schedule assignment to end %w", err)
	}

	return nil
}

func (r Runner) start(assignment WithTestDetails) error {
	err := r.Mailer.Send(core.MailConfig{
		TemplateName: "warning",
		Subject:      "5 minute reminder for your " + assignment.Test.Business.Name + " assignment",
		From:         "candidates",
		To:           assignment.CandidateEmail,
	}, assignment)
	if err != nil {
		return fmt.Errorf("could not send reminder email to candidate %s %w", assignment.CandidateEmail, err)
	}

	_, err = r.SchedulerClient.Start(StartInput{
		Type:       "init",
		ID:         int64(assignment.ID),
		ScheduleAt: r.Time().Add(r.StartDelay).Format(time.RFC3339),
		Data:       assignment,
	})
	if err != nil {
		return fmt.Errorf("could not schedule assignment to init %w", err)
	}

	return nil
}

func (r Runner) sendEnd(status string, data WithTestDetails) error {
	subject := "Thanks for submitting your test for " + data.Test.Business.Name
	if status != "submitted" {
		subject = "You missed the deadline for submitting your test"
	}

	err := r.Mailer.Send(core.MailConfig{
		TemplateName: status,
		Subject:      subject,
		From:         "candidates",
		To:           data.CandidateEmail,
	}, data)
	if err != nil {
		return fmt.Errorf("could not send email to candidate %w", err)
	}

	subject = data.CandidateName + " has submitted their assignment"
	if status != "submitted" {
		subject = data.CandidateName + " missed the deadline to submit their assignment"
	}

	err = r.Mailer.Send(core.MailConfig{
		TemplateName: status + "-recruiter",
		Subject:      subject,
		From:         "candidates",
		To:           data.Recruiter.Email,
	}, data)
	if err != nil {
		return fmt.Errorf("could not send email to recruiter %w", err)
	}

	return err
}
