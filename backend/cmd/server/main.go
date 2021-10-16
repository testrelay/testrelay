package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	firebase "firebase.google.com/go/v4"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sfn"
	"github.com/gorilla/mux"
	"github.com/mailgun/mailgun-go/v4"
	"go.uber.org/zap"
	"google.golang.org/api/option"

	"github.com/testrelay/testrelay/backend/internal/api"
	"github.com/testrelay/testrelay/backend/internal/auth"
	"github.com/testrelay/testrelay/backend/internal/core"
	"github.com/testrelay/testrelay/backend/internal/core/assignment"
	"github.com/testrelay/testrelay/backend/internal/core/assignmentuser"
	eventsHttp "github.com/testrelay/testrelay/backend/internal/events/http"
	"github.com/testrelay/testrelay/backend/internal/mail"
	"github.com/testrelay/testrelay/backend/internal/scheduler"
	"github.com/testrelay/testrelay/backend/internal/store/graphql"
	"github.com/testrelay/testrelay/backend/internal/vcs"
)

var (
	client       *graphql.HasuraClient
	githubClient *vcs.GithubClient
	sfnClient    *sfn.SFN
	mailer       core.Mailer
	inviter      assignment.Inviter
	gh           *api.GraphQLQueryHandler
	ah           eventsHttp.AssignmentHandler
	rh           eventsHttp.ReviewerHandler
	logger       *zap.SugaredLogger
)

func init() {
	client = graphql.NewClient(os.Getenv("HASURA_URL"), os.Getenv("HASURA_TOKEN"))

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("eu-west-2"),
	}))

	sfnClient = sfn.New(sess, &aws.Config{Region: aws.String("eu-west-2")})

	githubClient = vcs.NewClient(os.Getenv("GITHUB_ACCESS_TOKEN"))

	mg, err := mailgun.NewMailgunFromEnv()
	if err != nil {
		log.Fatalf("could not generate mailgun err: %s\n", err)
	}

	mailer = &mail.MailgunMailer{MG: mg}

	options := option.WithCredentialsJSON([]byte(os.Getenv("GOOGLE_SERVICE_ACC")))

	app, err := firebase.NewApp(context.Background(), nil, options)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	a, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("could not generate auth client err: %s\n", err)
	}

	inviter = assignment.Inviter{
		BusinessRepo:   client,
		Mailer:         mailer,
		AssignmentRepo: client,
		UserRepo:       client,
		Auth: auth.FirebaseClient{
			Auth:            a,
			CustomClaimName: "https://hasura.io/jwt/claims",
		},
		AppURL: os.Getenv("APP_URL"),
	}

	collector, err := vcs.NewGithubRepoCollectorFromENV()
	if err != nil {
		log.Fatal(err)
	}

	gh, err = api.NewGraphQLQueryHandler(
		os.Getenv("HASURA_URL"),
		&api.GraphResolver{
			HasuraURL: os.Getenv("HASURA_URL"),
			Collector: collector,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	zlog, _ := zap.NewDevelopment()
	if os.Getenv("APP_ENV") == "production" {
		zlog, _ = zap.NewProduction()
	}
	logger = zlog.Sugar()

	ah = eventsHttp.AssignmentHandler{
		HasuraClient: client,
		GithubClient: githubClient,
		Inviter:      inviter,
		Logger:       logger,
		Runner: assignment.Runner{
			Uploader:          githubClient,
			Cleaner:           githubClient,
			SubmissionChecker: githubClient,
			ReviewerCollector: client,
			EventCreator:      client,
			Mailer:            mailer,
			Logger:            logger,
		},
		Scheduler: assignment.Scheduler{
			Fetcher: client,
			SchedulerClient: scheduler.StepFunctionAssignmentScheduler{
				StateMachineArn: os.Getenv("ASSIGNMENT_SCHEDULER_ARN"),
				SFNClient:       sfn.New(sess, &aws.Config{Region: aws.String("eu-west-2")}),
			},
			VCSCreator: githubClient,
			Updater:    client,
		},
	}

	rh = eventsHttp.ReviewerHandler{
		Logger: logger,
		Assigner: assignmentuser.Assigner{
			ReviewerRepository: client,
			VCSClient:          githubClient,
			Mailer:             mailer,
		},
	}
}

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	r := mux.NewRouter()

	a := r.PathPrefix("assignments").Subrouter()
	a.Methods(http.MethodPost).Path("events").HandlerFunc(ah.EventHandler)
	a.Methods(http.MethodPost).Path("process").HandlerFunc(ah.ProcessHandler)

	re := r.PathPrefix("reviewers").Subrouter()
	re.Methods(http.MethodPost).Path("events").HandlerFunc(rh.EventsHandler)

	r.Methods(http.MethodPost).Path("graphql").HandlerFunc(gh.Query)

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}
