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
	firebaseAuth "firebase.google.com/go/v4/auth"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sfn"
	"github.com/gorilla/mux"
	"github.com/mailgun/mailgun-go/v4"
	"go.uber.org/zap"
	"google.golang.org/api/option"

	"github.com/testrelay/testrelay/backend/internal/api"
	"github.com/testrelay/testrelay/backend/internal/auth"
	"github.com/testrelay/testrelay/backend/internal/core/assignment"
	"github.com/testrelay/testrelay/backend/internal/core/assignmentuser"
	eventsHttp "github.com/testrelay/testrelay/backend/internal/events/http"
	"github.com/testrelay/testrelay/backend/internal/mail"
	"github.com/testrelay/testrelay/backend/internal/options"
	"github.com/testrelay/testrelay/backend/internal/scheduler"
	"github.com/testrelay/testrelay/backend/internal/store/graphql"
	"github.com/testrelay/testrelay/backend/internal/vcs"
)

var (
	gh *api.GraphQLQueryHandler
	ah eventsHttp.AssignmentHandler
	rh eventsHttp.ReviewerHandler
)

func init() {
	config, err := options.ConfigFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	logger := newLogger(config)

	hasuraClient := graphql.NewClient(config.HasuraURL, config.HasuraToken)
	githubClient := vcs.NewClient(config.GithubAccessToken)

	mailer := newMailer(config)

	ah = eventsHttp.AssignmentHandler{
		HasuraClient: hasuraClient,
		GithubClient: githubClient,
		Inviter: assignment.Inviter{
			BusinessRepo:   hasuraClient,
			Mailer:         mailer,
			AssignmentRepo: hasuraClient,
			UserRepo:       hasuraClient,
			Auth: auth.FirebaseClient{
				Auth:            newFirebaseAuth(config),
				CustomClaimName: "https://hasura.io/jwt/claims",
			},
			AppURL: config.AppURL,
		},
		Logger: logger,
		Runner: assignment.Runner{
			Uploader:          githubClient,
			Cleaner:           githubClient,
			SubmissionChecker: githubClient,
			ReviewerCollector: hasuraClient,
			EventCreator:      hasuraClient,
			Mailer:            mailer,
			Logger:            logger,
		},
		Scheduler: assignment.Scheduler{
			Fetcher: hasuraClient,
			SchedulerClient: scheduler.StepFunctionAssignmentScheduler{
				StateMachineArn: config.AssignmentSchedulerARN,
				SFNClient: sfn.New(session.Must(session.NewSession(&aws.Config{
					Region: aws.String(config.AWSRegion),
				})), &aws.Config{Region: aws.String(config.AWSRegion)}),
			},
			VCSCreator: githubClient,
			Updater:    hasuraClient,
		},
	}

	rh = eventsHttp.ReviewerHandler{
		Logger: logger,
		Assigner: assignmentuser.Assigner{
			ReviewerRepository: hasuraClient,
			VCSClient:          githubClient,
			Mailer:             mailer,
		},
	}

	gh = newGraphQLQueryHandler(config)
}

func newFirebaseAuth(config options.Config) *firebaseAuth.Client {
	app, err := firebase.NewApp(
		context.Background(),
		nil,
		option.WithCredentialsJSON([]byte(config.GoogleServiceAccount)),
	)
	if err != nil {
		log.Fatalf("error initializing firebase app: %v", err)
	}

	a, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("could not generate auth client err: %s", err)
	}

	return a
}

func newGraphQLQueryHandler(config options.Config) *api.GraphQLQueryHandler {
	collector, err := vcs.NewGithubRepoCollector(config.GithubPrivateKey, config.GithubAppID)
	if err != nil {
		log.Fatalf("could not init github repository collector %s", err)
	}

	gh, err := api.NewGraphQLQueryHandler(
		config.HasuraURL,
		&api.GraphResolver{
			HasuraURL: config.HasuraToken,
			Collector: collector,
		},
	)
	if err != nil {
		log.Fatalf("could not init graphql api handler %s", err)
	}

	return gh
}

func newMailer(config options.Config) *mail.MailgunMailer {
	mg := mailgun.NewMailgun(config.MGDomain, config.MGAPIKey)
	mg.SetAPIBase(config.MGURL)

	return &mail.MailgunMailer{MG: mg}
}

func newLogger(config options.Config) *zap.SugaredLogger {
	zlog, _ := zap.NewDevelopment()
	if config.APPEnv == "production" {
		zlog, _ = zap.NewProduction()
	}
	logger := zlog.Sugar()
	return logger
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
