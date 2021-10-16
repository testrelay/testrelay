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
	graphql2 "github.com/hasura/go-graphql-client"
	"github.com/mailgun/mailgun-go/v4"
	"google.golang.org/api/option"

	"github.com/testrelay/testrelay/backend/internal/event"
	"github.com/testrelay/testrelay/backend/internal/github"
	"github.com/testrelay/testrelay/backend/internal/graphql"
	http2 "github.com/testrelay/testrelay/backend/internal/http"
	"github.com/testrelay/testrelay/backend/internal/mail"
)

var (
	client       *graphql.HasuraClient
	githubClient *github.Client
	sfnClient    *sfn.SFN
	mailer       mail.Mailer
	processor    event.Processor
	gh           *graphql.HttpHandler
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

	collector, err := github.NewGithubRepoCollectorFromENV()
	if err != nil {
		log.Fatal(err)
	}

	gh, err = graphql.NewHttpHandler(
		os.Getenv("HASURA_URL"),
		&graphql.GraphResolver{
			HasuraURL: os.Getenv("HASURA_URL"),
			Collector: collector,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	r := mux.NewRouter()

	a := r.PathPrefix("assignments").Subrouter()
	a.Methods(http.MethodPost).Path("events").HandlerFunc()
	a.Methods(http.MethodPost).Path("process").HandlerFunc()

	re := r.PathPrefix("reviewers").Subrouter()
	re.Methods(http.MethodPost).Path("events").HandlerFunc()

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
