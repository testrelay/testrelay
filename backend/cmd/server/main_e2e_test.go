//go:build e2e
// +build e2e

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"

	firebase "firebase.google.com/go/v4"
	firebaseAuth "firebase.google.com/go/v4/auth"
	"github.com/google/go-github/v39/github"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"

	"github.com/testrelay/testrelay/backend/internal/httputil"
	"github.com/testrelay/testrelay/backend/internal/store/graphql"
	"github.com/testrelay/testrelay/backend/internal/test"
)

var (
	testUserGithubUsername = "testrelaycandidate"
	githubTestOwner        = "the-foreman"
)

var (
	githubClient    *github.Client
	rawGraphlClient test.GraphQLClient
	hasuraClient    *graphql.HasuraClient
	firebaseClient  *firebaseAuth.Client
)

func TestMain(m *testing.M) {
	err := godotenv.Overload("./test_assets/e2e.env")
	if err != nil {
		log.Fatal("error loading e2e.env file, please specify")
	}

	initGraphqlClients()
	initGithubClient()
	initFirebaseAuth()

	go run()

	err = waitForPort(8000)
	if err != nil {
		log.Fatal("backend server port was not ready after 3 tries")
	}

	res, err := http.Get("http://localhost:8000/healthz")
	if err != nil {
		log.Fatalf("backend server cannot be contacted %s", err)
	}

	if res.StatusCode != http.StatusOK {
		log.Fatal("backend server unhealthy")
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	os.Exit(code)
}

func initFirebaseAuth() {
	app, err := firebase.NewApp(
		context.Background(),
		nil,
		option.WithCredentialsFile(os.Getenv("GOOGLE_SERVICE_ACC")),
	)
	if err != nil {
		log.Fatalf("error initializing firebase app: %v", err)
	}

	a, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("could not generate auth client err: %s", err)
	}

	firebaseClient = a
}

func initGithubClient() {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_ACCESS_TOKEN")},
	)

	tc := oauth2.NewClient(context.Background(), ts)
	githubClient = github.NewClient(tc)
}

func initGraphqlClients() {
	rawGraphlClient = test.GraphQLClient{
		BaseURL: os.Getenv("HASURA_URL") + "/v1/graphql",
		Client: &http.Client{
			Transport: &httputil.KeyTransport{Key: "x-hasura-admin-secret", Value: os.Getenv("HASURA_TOKEN")},
		},
	}

	hasuraClient = graphql.NewClient(os.Getenv("HASURA_URL")+"/v1/graphql", os.Getenv("HASURA_TOKEN"))
}

func waitForPort(port int) error {
	for i := 0; i < 3; i++ {
		_, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d", port))
		if errors.Is(err, syscall.ECONNREFUSED) {
			time.Sleep(time.Second)
			continue
		}

		return nil
	}

	return fmt.Errorf("port %d was not open", port)
}