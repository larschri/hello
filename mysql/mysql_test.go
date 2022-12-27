package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest/v3"
)

var db *sql.DB

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("mysql", "5.7", []string{"MYSQL_ROOT_PASSWORD=secret"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		db, err = sql.Open("mysql", fmt.Sprintf("root:secret@(localhost:%s)/mysql", resource.GetPort("3306/tcp")))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestMySQL(t *testing.T) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS mytable (id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY, some_text TEXT NOT NULL)")
	if err != nil {
		t.Fatal(err)
	}

	// Create
	res, err := db.Exec("INSERT INTO mytable (some_text) VALUES (?)", "hello world")
	if err != nil {
		t.Fatal(err)
	}

	// get the id of the newly inserted record
	id, err := res.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	// Read
	var someText string
	row := db.QueryRow("SELECT some_text FROM mytable WHERE id = ? LIMIT 1", id)
	if err := row.Scan(&someText); err != nil {
		t.Fatal(err)
	}
	fmt.Println(someText)

	// Update
	_, err = db.Exec("UPDATE mytable SET some_text = ? WHERE id = ?", "Hello, æøå", id)
	if err != nil {
		t.Fatal(err)
	}

	// Delete
	_, err = db.Exec("DELETE FROM mytable WHERE id = ?", id)
	if err != nil {
		t.Fatal(err)
	}
}
