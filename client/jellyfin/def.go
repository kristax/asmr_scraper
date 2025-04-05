package jellyfin

import (
	"asmr_scraper/util/restyop"
	"context"
)

type Client interface {
	GetViews(ctx context.Context) ([]*ViewItem, error)
	GetViewIdByName(ctx context.Context, name string) (string, error)
	GetItems(ctx context.Context, parentId string, queryParams map[string]string, options ...restyop.Option) (*ItemsResponse, error)
	GetItem(ctx context.Context, itemId string) (*ItemInfoResponse, error)
	UpdateItem(ctx context.Context, req *UpdateItemRequest) error
	UploadPrimaryImage(ctx context.Context, itemId string, data []byte) error
}
