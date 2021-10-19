//go:build e2e
// +build e2e

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"

	firebase "firebase.google.com/go/v4"
	firebaseAuth "firebase.google.com/go/v4/auth"
	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"

	"github.com/testrelay/testrelay/backend/internal/httputil"
)

type graphErrors []struct {
	Message   string
	Locations []struct {
		Line   int
		Column int
	}
}

// Error implements error interface.
func (e graphErrors) Error() string {
	b := strings.Builder{}
	for _, err := range e {
		b.WriteString(fmt.Sprintf("Message: %s, Locations: %+v", err.Message, err.Locations))
	}
	return b.String()
}

type graphQLClient struct {
	client  *http.Client
	baseURL string
}

func (c graphQLClient) do(query string, variables map[string]interface{}, v interface{}) error {
	in := struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables,omitempty"`
	}{
		Query:     query,
		Variables: variables,
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(in)
	if err != nil {
		return err
	}

	resp, err := c.client.Post(c.baseURL, "application/json", &buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("non-200 OK status code: %v body: %q", resp.Status, body)
	}
	var out struct {
		Data   *json.RawMessage
		Errors graphErrors
	}

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &out)
	if err != nil {
		return err
	}

	if out.Data != nil {
		err := json.Unmarshal(*out.Data, &v)
		if err != nil {
			return err
		}
	}

	if len(out.Errors) > 0 {
		return out.Errors
	}

	return nil
}

var (
	githubClient   *github.Client
	hasuraClient   graphQLClient
	firebaseClient *firebaseAuth.Client
)

func TestMain(m *testing.M) {
	err := godotenv.Overload("test_assets/e2e.env")
	if err != nil {
		log.Fatal("error loading e2e.env file, please specify")
	}

	initHasuraClient()
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
	//if err := pool.Purge(pg); err != nil {
	//	log.Fatalf("Could not purge resource: %s", err)
	//}

	os.Exit(code)
}

func initFirebaseAuth(){
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

func initHasuraClient() {
	hasuraClient = graphQLClient{
		baseURL: os.Getenv("HASURA_URL"),
		client: &http.Client{
			Transport: &httputil.KeyTransport{Key: "x-hasura-admin-secret", Value: os.Getenv("HASURA_TOKEN")},
		},
	}
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

func bootContainers() {
	//pool, err := dockertest.NewPool("")
	//if err != nil {
	//	log.Fatalf("Could not connect to docker: %s", err)
	//}
	//
	//pg, err := pool.Run("postgres:12", "latest", []string{"POSTGRES_PASSWORD=postgrespassword"})
	//if err != nil {
	//	log.Fatalf("Could not start resource pg: %s", err)
	//}
	//
	//// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	//if err := pool.Retry(func() error {
	//	var err error
	//	connStr := "user=postgres dbname=postgres password=postgrespassword ssl-mode=skip-verify"
	//	db, err = sql.Open("postgres", connStr)
	//	if err != nil {
	//		return err
	//	}
	//	return db.Ping()
	//}); err != nil {
	//	log.Fatalf("Could not connect to pg: %s", err)
	//}
	//
	//// add host.docker.internal as the application URL
	//hasura, err := pool.Run("hasura/graphql-engine:v2.0.9.cli-migrations-v3", "latest", []string{"POSTGRES_PASSWORD=postgrespassword"})
	//if err != nil {
	//	log.Fatalf("Could not start resource pg: %s", err)
	//}
	//
	//// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	//if err := pool.Retry(func() error {
	//	// hasura retry functionality
	//}); err != nil {
	//	log.Fatalf("Could not connect to pg: %s", err)
	//}
}
