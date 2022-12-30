package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var db *sql.DB

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "11",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=user_name",
			"POSTGRES_DB=dbname",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://user_name:secret@%s/dbname?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url: ", databaseUrl)

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		db, err = sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	//Run tests
	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestPosgres(t *testing.T) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS mytable (id SERIAL PRIMARY KEY, some_text TEXT NOT NULL)")
	if err != nil {
		t.Fatal(err)
	}

	// Create
	var id int
	err = db.QueryRow("INSERT INTO mytable(some_text) VALUES ($1) RETURNING id", "hello world").Scan(&id)
	if err != nil {
		t.Fatal(err)
	}

	// Read
	var someText string
	row := db.QueryRow("SELECT some_text FROM mytable WHERE id = $1 LIMIT 1", id)
	if err := row.Scan(&someText); err != nil {
		t.Fatal(err)
	}
	fmt.Println(someText)

	// Update
	_, err = db.Exec("UPDATE mytable SET some_text = $1 WHERE id = $2", "Hello, æøå", id)
	if err != nil {
		t.Fatal(err)
	}

	// Delete
	_, err = db.Exec("DELETE FROM mytable WHERE id = $1", id)
	if err != nil {
		t.Fatal(err)
	}
}
