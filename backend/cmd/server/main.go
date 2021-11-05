package main

import (
	"bytes"
	"context"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	firebase "firebase.google.com/go/v4"
	firebaseAuth "firebase.google.com/go/v4/auth"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"google.golang.org/api/option"

	"github.com/testrelay/testrelay/backend/internal/api"
	"github.com/testrelay/testrelay/backend/internal/auth"
	"github.com/testrelay/testrelay/backend/internal/core"
	"github.com/testrelay/testrelay/backend/internal/core/assignment"
	"github.com/testrelay/testrelay/backend/internal/core/assignmentuser"
	eventsHttp "github.com/testrelay/testrelay/backend/internal/events/http"
	"github.com/testrelay/testrelay/backend/internal/mail"
	"github.com/testrelay/testrelay/backend/internal/options"
	"github.com/testrelay/testrelay/backend/internal/scheduler"
	"github.com/testrelay/testrelay/backend/internal/store/graphql"
	"github.com/testrelay/testrelay/backend/internal/vcs"
)

func newFirebaseAuth(config options.Config) *firebaseAuth.Client {
	app, err := firebase.NewApp(
		context.Background(),
		nil,
		option.WithCredentialsFile(config.GoogleServiceAccountLocation),
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
	collector, err := vcs.NewGithubRepoCollector(config.GithubPrivateKeyLocation, config.GithubAppID)
	if err != nil {
		log.Fatalf("could not init github repository collector %s", err)
	}

	gh, err := api.NewGraphQLQueryHandler(
		config.HasuraURL+"/v1/graphql",
		&api.GraphResolver{
			HasuraURL: config.HasuraURL + "/v1/graphql",
			Collector: collector,
		},
	)
	if err != nil {
		log.Fatalf("could not init graphql api handler %s", err)
	}

	return gh
}

func newMailer(config options.Config) mail.SMTPMailer {
	m, err := mail.NewSMTPMailer(core.SMTPConfig{
		SendingDomain: config.MailFromDomain,
		Host:          config.SMTPHost,
		Port:          int(config.SMTPPort),
		Username:      config.SMTPUsername,
		Password:      config.SMTPPassword,
	})
	if err != nil {
		log.Fatalf("could not init mailer %s", err)
	}

	return m
}

func newLogger(config options.Config) *zap.SugaredLogger {
	c := zap.NewDevelopmentConfig()
	c.DisableStacktrace = true
	zlog, _ := c.Build()

	if config.APPEnv == "production" {
		zlog, _ = zap.NewProduction()
	}
	logger := zlog.Sugar()
	return logger
}

func run() {
	config, err := options.ConfigFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	logger := newLogger(config)

	hasuraClient := graphql.NewClient(config.HasuraURL+"/v1/graphql", config.HasuraToken)
	githubClient := vcs.NewClient(config.GithubAccessToken)

	mailer := newMailer(config)

	scheduleClient := scheduler.NewHasuraAssignmentScheduler(
		config.HasuraURL,
		config.HasuraToken,
		config.BackendURL,
	)
	ah := eventsHttp.AssignmentHandler{
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
			SchedulerClient:   scheduleClient,
			Time:              time.Now,
			StartDelay:        time.Minute * 5,
			WarningBeforeEnd:  time.Minute * 10,
		},
		Scheduler: assignment.Scheduler{
			Fetcher:         hasuraClient,
			SchedulerClient: scheduleClient,
			VCSCreator:      githubClient,
			Updater:         hasuraClient,
		},
	}

	rh := eventsHttp.ReviewerHandler{
		Logger: logger,
		Assigner: assignmentuser.Assigner{
			ReviewerRepository: hasuraClient,
			VCSClient:          githubClient,
			Mailer:             mailer,
		},
	}

	gh := newGraphQLQueryHandler(config)

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	r := mux.NewRouter()
	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/healthz" {
				b, _ := io.ReadAll(r.Body)
				r.Body = io.NopCloser(bytes.NewBuffer(b))
				logger.Info(r.URL.Path, string(b))
			}

			h.ServeHTTP(w, r)
		})
	})

	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	a := r.PathPrefix("/assignments").Subrouter()
	a.Methods(http.MethodPost).Path("/events").HandlerFunc(ah.EventHandler)
	a.Methods(http.MethodPost).Path("/process").HandlerFunc(ah.ProcessHandler)

	re := r.PathPrefix("/reviewers").Subrouter()
	re.Methods(http.MethodPost).Path("/events").HandlerFunc(rh.EventsHandler)

	r.Methods(http.MethodPost).Path("/graphql").HandlerFunc(gh.Query)

	srv := &http.Server{
		Addr:         "0.0.0.0:8000",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		log.Println("starting server")
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

func main() {
	run()
}
