package jellyfin

import (
	"asmr_scraper/util/restyop"
	"context"
)

type Client interface {
	GetItems(ctx context.Context, parentId string, options ...restyop.Option) (*ItemsResponse, error)
	GetItem(ctx context.Context, itemId string) (*ItemInfoResponse, error)
	UpdateItem(ctx context.Context, req *UpdateItemRequest) error
	UploadPrimaryImage(ctx context.Context, itemId string, data []byte) error
}
