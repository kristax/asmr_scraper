package asmr_one

import (
	"asmr_scraper/model"
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"time"
)

type RateCountDetail struct {
	ReviewPoint int `json:"review_point"`
	Count       int `json:"count"`
	Ratio       int `json:"ratio"`
}

type Rank struct {
	Term     string `json:"term"`
	Category string `json:"category"`
	Rank     int    `json:"rank"`
	RankDate string `json:"rank_date"`
}

type Vas struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Tag struct {
	Id   int              `json:"id"`
	I18N map[string]*I18N `json:"i18n"`
	Name string           `json:"name"`
}

type I18N struct {
	Name string `json:"name"`
}

type Circle struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	SourceId   string `json:"source_id"`
	SourceType string `json:"source_type"`
}

type TranslationInfo struct {
	Lang                     string          `json:"lang"`
	IsChild                  bool            `json:"is_child"`
	IsParent                 bool            `json:"is_parent"`
	IsOriginal               bool            `json:"is_original"`
	IsVolunteer              bool            `json:"is_volunteer"`
	ChildWorknos             json.RawMessage `json:"child_worknos"`
	ParentWorkno             string          `json:"parent_workno"`
	OriginalWorkno           string          `json:"original_workno"`
	IsTranslationAgree       bool            `json:"is_translation_agree"`
	TranslationBonusLangs    json.RawMessage `json:"translation_bonus_langs"`
	IsTranslationBonusChild  bool            `json:"is_translation_bonus_child"`
	ProductionTradePriceRate int             `json:"production_trade_price_rate"`
}

type WorkInfoResponse struct {
	Id                        int                `json:"id"`
	Title                     string             `json:"title"`
	CircleId                  int                `json:"circle_id"`
	Name                      string             `json:"name"`
	Nsfw                      bool               `json:"nsfw"`
	Release                   string             `json:"release"`
	DlCount                   int                `json:"dl_count"`
	Price                     int                `json:"price"`
	ReviewCount               int                `json:"review_count"`
	RateCount                 int                `json:"rate_count"`
	RateAverage2Dp            float64            `json:"rate_average_2dp"`
	RateCountDetail           []*RateCountDetail `json:"rate_count_detail"`
	Rank                      []*Rank            `json:"rank"`
	HasSubtitle               bool               `json:"has_subtitle"`
	CreateDate                string             `json:"create_date"`
	Vas                       []*Vas             `json:"vas"`
	Tags                      []*Tag             `json:"tags"`
	LanguageEditions          json.RawMessage    `json:"language_editions"`
	OriginalWorkno            string             `json:"original_workno"`
	OtherLanguageEditionsInDb json.RawMessage    `json:"other_language_editions_in_db"`
	TranslationInfo           *TranslationInfo   `json:"translation_info"`
	WorkAttributes            string             `json:"work_attributes"`
	AgeCategoryString         string             `json:"age_category_string"`
	Duration                  int                `json:"duration"`
	SourceType                string             `json:"source_type"`
	SourceId                  string             `json:"source_id"`
	SourceUrl                 string             `json:"source_url"`
	Circle                    *Circle            `json:"circle"`
	SamCoverUrl               string             `json:"samCoverUrl"`
	ThumbnailCoverUrl         string             `json:"thumbnailCoverUrl"`
	MainCoverUrl              string             `json:"mainCoverUrl"`
}

func (workInfo *WorkInfoResponse) ToProjectInfo(code string) (*model.ProjectInfo, error) {
	tags := lo.Map[*Tag, string](workInfo.Tags, func(item *Tag, _ int) string {
		return item.Name
	})
	releaseDate, err := time.Parse(time.DateOnly, workInfo.Release)
	if err != nil {
		releaseDate = time.Now()
	}
	createDate, err := time.Parse(time.DateOnly, workInfo.CreateDate)
	if err != nil {
		createDate = time.Now()
	}
	artist := lo.Map(workInfo.Vas, func(item *Vas, _ int) string {
		return item.Name
	})
	if len(artist) < 1 {
		artist = append(artist, "Unknown")
	}
	return &model.ProjectInfo{
		ItemId:          "",
		Code:            code,
		Path:            "",
		Name:            workInfo.Title,
		Name2:           "",
		Tags:            tags,
		ReleaseDate:     releaseDate,
		CreateDate:      createDate,
		Artist:          artist,
		Rating:          workInfo.RateAverage2Dp,
		Group:           workInfo.Circle.Name,
		Nsfw:            workInfo.Nsfw,
		Price:           workInfo.Price,
		Sales:           workInfo.DlCount,
		Overview:        fmt.Sprintf(overviewTemplate, workInfo.Circle.Name, workInfo.Price, workInfo.DlCount, code, code),
		PrimaryImageUrl: workInfo.MainCoverUrl,
	}, nil
}

func asmrOneAdaptor(rjCode string, workInfo *WorkInfoResponse) *model.ProjectInfo {
	tags := lo.Map[*Tag, string](workInfo.Tags, func(item *Tag, _ int) string {
		return item.Name
	})
	releaseDate, err := time.Parse(time.DateOnly, workInfo.Release)
	if err != nil {
		releaseDate = time.Now()
	}
	createDate, err := time.Parse(time.DateOnly, workInfo.CreateDate)
	if err != nil {
		createDate = time.Now()
	}
	artist := lo.Map(workInfo.Vas, func(item *Vas, _ int) string {
		return item.Name
	})
	if len(artist) < 1 {
		artist = append(artist, "Unknown")
	}
	return &model.ProjectInfo{
		ItemId:          "",
		Code:            rjCode,
		Path:            "",
		Name:            workInfo.Title,
		Name2:           "",
		Tags:            tags,
		ReleaseDate:     releaseDate,
		CreateDate:      createDate,
		Artist:          artist,
		Rating:          workInfo.RateAverage2Dp,
		Group:           workInfo.Circle.Name,
		Nsfw:            workInfo.Nsfw,
		Price:           workInfo.Price,
		Sales:           workInfo.DlCount,
		Overview:        fmt.Sprintf(overviewTemplate, workInfo.Circle.Name, workInfo.Price, workInfo.DlCount, rjCode, rjCode),
		PrimaryImageUrl: workInfo.MainCoverUrl,
	}
}
