//go:build e2e
// +build e2e

package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

//var db *sql.DB

func TestMain(m *testing.M) {
	err := godotenv.Overload("test_assets/e2e.env")
	if err != nil {
		log.Fatal("error loading e2e.env file, please specify")
	}

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

func bootContainers()  {
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