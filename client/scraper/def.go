package scraper

import (
	"asmr_scraper/client/jellyfin"
	"asmr_scraper/model"
	"context"
)

type Client interface {
	RefreshInfo(ctx context.Context, target *model.Target) (*RefreshInfoResult, error)
}

type SourceClient interface {
	TargetName() string
	Order() int
	ParseCode(ctx context.Context, item *jellyfin.ItemInfoResponse) (string, error)
	GetProjectInfo(ctx context.Context, code string) (*model.ProjectInfo, error)
	ImageMissing(item *jellyfin.ItemInfoResponse) bool
	InfoMissing(item *jellyfin.ItemInfoResponse) bool
}
