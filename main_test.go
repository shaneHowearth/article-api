package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/lib/pq"
)

var a App

func TestMain(m *testing.M) {
	a = App{}
	a.Initialise(
		os.Getenv("TEST_DB_USERNAME"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"))
	ensureTablesExist()

	code := m.Run()

	clearTables()

	os.Exit(code)
}

func ensureTablesExist() {
	if _, err := a.DB.Exec(articleTableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTables() {
	a.DB.Exec("DELETE FROM article")
	a.DB.Exec("ALTER SEQUENCE article_id_seq RESTART WITH 1")
}

const articleTableCreationQuery = `CREATE TABLE IF NOT EXISTS article
(id SERIAL PRIMARY KEY,
 title	TEXT,
 pub_date	TIMESTAMP,
 body	TEXT,
 tags	TEXT[]
)`

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func addArticles(count int, date time.Time) {
	// Always create at least one
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		strI := strconv.Itoa(i)
		tags := []string{"go", "goroutines", "queues"}

		_, err := a.DB.Exec("INSERT INTO article(title, pub_date, body, tags) VALUES($1, $2, $3, $4)", "Title "+strI, date, "Body "+strI, pq.Array(tags))
		if err != nil {
			log.Fatal(err)
		}
	}
}
