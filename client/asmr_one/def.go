package asmr_one

import (
	"asmr_scraper/client/scraper"
	"context"
)

type Client interface {
	scraper.SourceClient
	GetWorkInfo(ctx context.Context, rj string) (*WorkInfoResponse, error)
}
