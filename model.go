package main

import (
	"database/sql"
	"log"
	"time"

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

// Tag -
type Tag struct {
	TagName     string   `json:"tag"`
	Count       int      `json:"count"`
	Articles    []int    `json:"articles"`
	Related     []string `json:"related_tags"`
	ArticleDate string   `json:"-"`
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

func (t *Tag) getTagInfo(db *sql.DB) error {
	d, _ := time.Parse("20060102", t.ArticleDate)
	rows, err := db.Query("SELECT id, tags FROM article WHERE $1=ANY(tags) AND pub_date=$2", t.TagName, d)
	tmpTags := make(map[string]int)
	for rows.Next() {
		var (
			id   int
			tags []string
		)
		if err := rows.Scan(&id, pq.Array(&tags)); err != nil {
			log.Fatal(err)
		}
		t.Articles = append(t.Articles, id)
		for _, tag := range tags {
			tmpTags[tag] = 1
		}
	}
	for k, _ := range tmpTags {
		t.Related = append(t.Related, k)
	}
	t.Count = len(t.Articles)

	return err
}
