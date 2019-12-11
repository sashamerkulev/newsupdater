package models

import (
	"time"
)

type Article struct {
	SourceName  string
	Title       string
	Link        string
	Description string
	PubDate     time.Time
	Category    string
	PictureUrl  string
}

type SaverArticlesFunc func(articles []Article)
