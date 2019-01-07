package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestEmptyArticleID(t *testing.T) {
	clearTables()

	req, _ := http.NewRequest("GET", "/articles", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusMethodNotAllowed, response.Code)
}

func TestGetNonExistentArticle(t *testing.T) {
	clearTables()

	req, _ := http.NewRequest("GET", "/articles/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Article not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Article not found'. Got '%s'", m["error"])
	}
}

func TestCreateArticle(t *testing.T) {
	clearTables()

	payload := []byte(`{"title":"test","body":"article", "date": "2016-09-22", "tags" : ["health", "fitness", "science"] }`)

	req, _ := http.NewRequest("POST", "/articles", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["title"] != "test" {
		t.Errorf("Expected article title to be 'test'. Got '%v'", m["title"])
	}

	if m["body"] != "article" {
		t.Errorf("Expected article body to be 'article'. Got '%v'", m["article"])
	}

	if m["date"] != "2016-09-22" {
		t.Errorf("Expected article date to be '2016-09-22'. Got '%v'", m["date"])
	}

	/*
	* Tags is a slice of strings,
	* In order to check that they match the expected slice
	* the m["tags"] needs a type assertion to a slice of interfaces
	* which can then be looped over.
	* Each tag within that slice of interfaces is type asserted to be a string,
	* and then tested against the corresponding string in the expected slice of strings.
	 */
	expected := []string{"health", "fitness", "science"}
	for idx, tag := range m["tags"].([]interface{}) {
		if tag.(string) != expected[idx] {
			t.Errorf("Expected article tags to be '%s'. Got '%v'", expected[idx], tag)
		}
	}

	if m["id"] != "1" {
		t.Errorf("Expected article ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetArticle(t *testing.T) {
	clearTables()
	addArticles(1, time.Now())

	req, _ := http.NewRequest("GET", "/articles/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetTags(t *testing.T) {
	clearTables()

	addArticles(2, time.Now())
	addArticles(2, time.Now().AddDate(0, 0, 1))

	d := time.Now()
	year, month, day := d.Date()
	uri := fmt.Sprintf("/tag/queues/%d%02d%02d", year, int(month), day)
	req, _ := http.NewRequest("GET", uri, nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}
