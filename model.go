package main

import (
	"database/sql"

	"github.com/lib/pq"
)

// Article -
type Article struct {
	ID      string   `json:"id,omitempty"`
	Title   string   `json:"title"`
	PubDate string   `json:"date,omitempty"`
	Body    string   `json:"body"`
	Tags    []string `json:"tags,omitempty"`
}

func (a *Article) getArticle(db *sql.DB) error {
	return db.QueryRow("SELECT title, pub_date, body, tags FROM article WHERE id=$1",
		a.ID).Scan(&a.Title, &a.PubDate, &a.Body, pq.Array(&a.Tags))
}

func (a *Article) createArticle(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO article(title, pub_date, body, tags) VALUES($1, $2, $3, $4) RETURNING id",
		a.Title, a.PubDate, a.Body, pq.Array(&a.Tags)).Scan(&a.ID)

	if err != nil {
		return err
	}

	return nil
}
