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
