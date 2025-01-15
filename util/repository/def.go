package repository

import (
	"asmr_scraper/model"
	"context"
	"gorm.io/gorm"
)

type Repository interface {
	Do(f func(db *gorm.DB) any) error
	GetDataCacheByCode(ctx context.Context, target, code string) (*model.DataCache, error)
	SaveDataCache(ctx context.Context, data *model.DataCache) error
}
