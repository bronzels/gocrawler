package web

import (
	"time"
)

type News struct {
	EntryURL string

	URL   string
	URLid string

	QuoteeURL   string
	QuoteeURLid string

	Title       string
	PublishedAt time.Time
	Contents    string
	CrawledAt   time.Time
}
