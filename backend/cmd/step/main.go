package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mailgun/mailgun-go/v4"

	"github.com/testrelay/testrelay/backend/internal"
	"github.com/testrelay/testrelay/backend/internal/github"
	"github.com/testrelay/testrelay/backend/internal/graphql"
	"github.com/testrelay/testrelay/backend/internal/mail"
)

var (
	mailer        mail.Mailer
	ghClient      *github.Client
	graphQLClient *graphql.HasuraClient
)

func init() {
	mg, err := mailgun.NewMailgunFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	mailer = &mail.MailgunMailer{MG: mg}
	ghClient = github.NewClient(os.Getenv("GITHUB_ACCESS_TOKEN"))

	graphQLClient = graphql.NewClient(os.Getenv("HASURA_URL"), os.Getenv("HASURA_TOKEN"))

}

func Handler(request internal.StepPayload) (internal.Data, error) {
	log.Printf("%+v", request)

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	assignment := request.Data.Assignment
	switch request.Step {
	case "start":
		err := mailer.Send(mail.Config{
			TemplateName: "warning",
			Subject:      "5 minute reminder for your " + assignment.Test.Business.Name + " assignment",
			From:         "candidates@testrelay.io",
			To:           assignment.CandidateEmail,
		}, assignment)
		if err != nil {
			return internal.Data{}, err
		}
	case "init":
		err := ghClient.Upload(assignment)
		if err != nil {
			return internal.Data{}, err
		}

		err = graphQLClient.NewAssignmentEvent(assignment.CandidateID, assignment.ID, "inprogress")
		if err != nil {
			return internal.Data{}, fmt.Errorf("could not add assignment status %s %w", "inprogress", err)
		}
	case "end":
		err := mailer.Send(mail.Config{
			TemplateName: "end",
			Subject:      "Your test is about to finish",
			From:         "candidates@testrelay.io",
			To:           assignment.CandidateEmail,
		}, assignment)
		if err != nil {
			return internal.Data{}, err
		}
	case "cleanup":
		reviewers, err := graphQLClient.Reviewers(request.Data.AssignmentID)
		if err != nil {
			return internal.Data{}, err
		}

		err = ghClient.Cleanup(assignment, reviewers)
		if err != nil {
			return internal.Data{}, err
		}

		ok, err := ghClient.IsSubmitted(assignment)
		if err != nil {
			return internal.Data{}, err
		}

		status := "submitted"
		if !ok {
			status = "missed"
		}

		err = graphQLClient.NewAssignmentEvent(assignment.CandidateID, assignment.ID, status)
		if err != nil {
			return internal.Data{}, fmt.Errorf("could not add assignment status %s %w", status, err)
		}

		err = mailer.SendEnd(status, assignment)
		if err != nil {
			return internal.Data{}, fmt.Errorf("could not send email to candidate %w", err)
		}

	}

	return request.Data, nil
}

func main() {
	lambda.Start(Handler)
}
