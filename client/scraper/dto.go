package scraper

import "time"

type RefreshInfoResult struct {
}

type ProjectInfo struct {
	ItemId      string
	RJCode      string
	Path        string
	Name        string
	Tags        []string
	ReleaseDate time.Time
	CreateDate  time.Time
	Artists     []string
	Rating      float64
	Group       string
	Nsfw        bool
	Price       int
	Sales       int
}
