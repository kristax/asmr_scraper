package scraper

import "context"

type Client interface {
	RefreshInfo(ctx context.Context, parentId string) (*RefreshInfoResult, error)
}
