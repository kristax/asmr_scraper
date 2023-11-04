package scraper

import (
	"asmr_scraper/model"
	"context"
)

type Client interface {
	RefreshInfo(ctx context.Context, target *model.Target) (*RefreshInfoResult, error)
}

type SourceClient interface {
	TargetName() string
	Order() int
	ParseCodeFromPath(ctx context.Context, path string) (string, error)
	GetProjectInfo(ctx context.Context, code string) (*model.ProjectInfo, error)
}
