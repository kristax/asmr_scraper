package repository

import (
	"asmr_scraper/model"
	"context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type repository struct {
	db *gorm.DB
}

func New() Repository {
	return &repository{}
}

//func (r *repository) Init() error {
//	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
//		Logger: logger.Default.LogMode(logger.Info),
//	})
//	if err != nil {
//		return err
//	}
//	r.db = db
//	return nil
//}

func (r *repository) Init() error {
	dsn := "root:QIANXIAOfanhua123@tcp(kristas.top:13306)/jellyfin_scraper?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}
	r.db = db
	return nil
}

func (r *repository) Do(f func(db *gorm.DB) any) error {
	result := f(r.db)
	switch rt := result.(type) {
	case error:
		return rt
	case *gorm.DB:
		return rt.Error
	}
	return nil
}

func (r *repository) GetDataCacheByCode(ctx context.Context, target, code string) (*model.DataCache, error) {
	var result = &model.DataCache{}
	err := r.db.WithContext(ctx).Where("`target` = ? and `code` = ?", target, code).First(result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *repository) SaveDataCache(ctx context.Context, data *model.DataCache) error {
	return r.db.WithContext(ctx).Save(data).Error
}
