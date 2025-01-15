package javdb

import (
	"asmr_scraper/client/scraper"
	"context"
)

type Client interface {
	scraper.SourceClient
	Get(ctx context.Context, avCode, lang string) (*Detail, error)
}
