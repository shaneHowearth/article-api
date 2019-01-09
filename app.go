// app.go

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// App -
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// Initialise -
func (a *App) Initialise(user, password, dbname string) {
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s", user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initialiseRoutes()
}

// Run -
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8000", a.Router))
}

func (a *App) initialiseRoutes() {
	a.Router.HandleFunc("/articles/{id:[0-9]+}", a.getArticle).Methods("GET")
	a.Router.HandleFunc("/articles", a.createArticle).Methods("POST")
	a.Router.HandleFunc("/articles/", a.createArticle).Methods("POST")
	a.Router.HandleFunc("/tag/{tagName:[a-zA-Z0-9]+}/{date:[0-9]{8}}", a.getTags).Methods("GET")
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) getArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Article ID")
		return
	}

	art := Article{ID: vars["id"]}
	if err := art.getArticle(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Article not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, art)
}

func (a *App) createArticle(w http.ResponseWriter, r *http.Request) {
	var art Article
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&art); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := art.createArticle(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, art)
}

func (a *App) getTags(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_, err := strconv.Atoi(vars["date"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid date")
		return
	}
	t := Tag{ArticleDate: vars["date"], TagName: vars["tagName"]}
	if err := t.getTagInfo(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Tag with date not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, t)
}
