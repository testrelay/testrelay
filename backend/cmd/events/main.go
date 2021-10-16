package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	firebase "firebase.google.com/go/v4"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sfn"
	graphql2 "github.com/hasura/go-graphql-client"
	"github.com/mailgun/mailgun-go/v4"
	"google.golang.org/api/option"

	"github.com/testrelay/testrelay/backend/internal"
	"github.com/testrelay/testrelay/backend/internal/event"
	"github.com/testrelay/testrelay/backend/internal/github"
	"github.com/testrelay/testrelay/backend/internal/graphql"
	http2 "github.com/testrelay/testrelay/backend/internal/http"
	"github.com/testrelay/testrelay/backend/internal/mail"
	intTime "github.com/testrelay/testrelay/backend/internal/time"
)

var (
	client       *graphql.HasuraClient
	githubClient *github.Client
	sfnClient    *sfn.SFN
	mailer       mail.Mailer
	processor    event.Processor
)

func init() {
	client = graphql.NewClient(os.Getenv("HASURA_URL"), os.Getenv("HASURA_TOKEN"))

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-west-2"),
	}))

	sfnClient = sfn.New(sess, &aws.Config{Region: aws.String("eu-west-2")})

	githubClient = github.NewClient(os.Getenv("GITHUB_ACCESS_TOKEN"))

	mg, err := mailgun.NewMailgunFromEnv()
	if err != nil {
		log.Fatalf("could not generate mailgun err: %s\n", err)
	}

	mailer = &mail.MailgunMailer{MG: mg}

	graphClient := graphql2.NewClient(
		"https://delicate-gator-74.hasura.app/v1/graphql",
		&http.Client{
			Transport: &http2.KeyTransport{Key: "x-hasura-admin-secret", Value: os.Getenv("HASURA_TOKEN")},
		},
	)

	options := option.WithCredentialsJSON([]byte(os.Getenv("GOOGLE_SERVICE_ACC")))

	app, err := firebase.NewApp(context.Background(), nil, options)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	a, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("could not generate auth client err: %s\n", err)
	}

	processor = event.AWSProcessor{
		GraphqlClient: graphClient,
		Mailer:        mailer,
		Auth:          a,
		AppURL:        os.Getenv("APP_URL"),
	}
}

type AssignmentSchedulerInput struct {
	AssignmentID int                `json:"assignmentId"`
	TestStart    string             `json:"testStart"`
	TestDuration int                `json:"testDuration"`
	Assignment   graphql.Assignment `json:"assignment"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("inbound assignment event %s\n", request.Body)

	var data event.HasuraEvent
	err := json.Unmarshal([]byte(request.Body), &data)
	if err != nil {
		log.Fatalf("could not unmarshal body of event error %s\n", err)
	}

	switch data.Table.Name {
	case "assignments":
		err := processor.Process(data.Event)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, err
		}
	case "assignment_events":
		var body internal.AssignmentEvent
		if err := json.Unmarshal(data.Event.Data.New, &body); err != nil {
			return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, err
		}

		if data.Event.Op == "INSERT" && body.EventType == "scheduled" {
			return handleAssignmentScheduled(body)
		}
	case "assignment_users":
		if data.Event.Op == "INSERT" {
			var body internal.AssignmentUser
			if err := json.Unmarshal(data.Event.Data.New, &body); err != nil {
				return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, err
			}

			au, err := client.GetAssignmentUser(body.ID)
			if err != nil {
				return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, err
			}

			if au.Assignment.GithubRepoUrl != "" && au.User.GithubUsername != "" {
				err := githubClient.AddCollaborator(string(au.Assignment.GithubRepoUrl), string(au.User.GithubUsername))
				if err != nil {
					return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, err
				}
			}

			err = mailer.SendReviewerInvite(mail.EmailData{
				Sender:        "info@testrelay.io",
				Email:         string(au.User.Email),
				CandidateName: string(au.Assignment.CandidateName),
			})
			if err != nil {
				return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, err
			}
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       `{"status": "ok"}`,
	}, nil
}

func handleAssignmentScheduled(data internal.AssignmentEvent) (events.APIGatewayProxyResponse, error) {
	assignment, err := client.GetAssignment(data.AssignmentID)
	if err != nil {
		return failAndLog(fmt.Sprintf("could not get assignment %s", err))
	}
	log.Printf("assignment returned %+v\n", assignment)

	if assignment.StepArn != "" {
		_, err := sfnClient.StopExecution(&sfn.StopExecutionInput{
			ExecutionArn: aws.String(string(assignment.StepArn)),
		})
		if err != nil {
			return failAndLog(fmt.Sprintf("could not stop arn %s executing err: %s", assignment.StepArn, err))
		}
	}

	t, err := intTime.Parse(intTime.AssignmentChoices{
		DayChosen:  string(assignment.TestDayChosen),
		TimeChosen: string(assignment.TestTimeChosen),
		Timezone:   string(assignment.TestTimezoneChosen),
	})
	if err != nil {
		return failAndLog(err.Error())
	}

	githubRepoURL := string(assignment.GithubRepoUrl)
	if githubRepoURL == "" {
		githubRepoURL, err = githubClient.CreateRepo(string(assignment.Test.Business.Name), string(assignment.Candidate.GithubUsername))
		if err != nil {
			return failAndLog(err.Error())
		}
	}

	assignment.GithubRepoUrl = graphql2.String(githubRepoURL)
	b, _ := json.Marshal(AssignmentSchedulerInput{
		AssignmentID: int(assignment.ID),
		TestStart:    t.SendNotificationAt,
		TestDuration: int(assignment.TimeLimit) - (600),
		Assignment:   *assignment,
	})

	out, err := sfnClient.StartExecution(&sfn.StartExecutionInput{
		Input:           aws.String(string(b)),
		Name:            aws.String(fmt.Sprintf("assignment-%d-%d", assignment.ID, time.Now().Unix())),
		StateMachineArn: aws.String(os.Getenv("ASSIGNMENT_SCHEDULER_ARN")),
	})
	if err != nil {
		return failAndLog(fmt.Sprintf("could not start step func exection %s", err))
	}

	err = client.UpdateAssignmentWithDetails(int(assignment.ID), *out.ExecutionArn, githubRepoURL)
	if err != nil {
		return failAndLog(
			fmt.Sprintf(
				"could not update assignment id %d with execution arn %s and github url %s err: %s",
				assignment.ID,
				*out.ExecutionArn,
				githubRepoURL,
				err,
			),
		)
	}

	return events.APIGatewayProxyResponse{}, nil
}

func main() {
	lambda.Start(Handler)
}

func failAndLog(msg string) (events.APIGatewayProxyResponse, error) {
	log.Println(msg)

	return events.APIGatewayProxyResponse{
		Body:       `{"message": "fail"}`,
		StatusCode: 400,
	}, nil
}
