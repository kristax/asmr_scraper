package asmr_one

import (
	"asmr_scraper/client/jellyfin"
	"asmr_scraper/model"
	"asmr_scraper/util/guess_epsides"
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"log"
	"os"
	"path/filepath"
	"strings"
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

func (workInfo *WorkInfoResponse) ToProjectInfo(code, path string, item *jellyfin.ItemInfoResponse, subItems []*jellyfin.ItemInfoResponse) (*model.ProjectInfo, error) {
	base := filepath.Base(path)
	tags := lo.Map[*Tag, string](workInfo.Tags, func(item *Tag, _ int) string {
		return item.Name
	})
	releaseDate, err := time.Parse(time.DateOnly, workInfo.Release)
	if err != nil {
		releaseDate = time.Now()
	}
	artist := lo.Map(workInfo.Vas, func(item *Vas, _ int) string {
		return item.Name
	})
	if len(artist) < 1 {
		artist = append(artist, "Unknown")
	}
	var subItemsInfo []*model.ProjectInfo
	var extractName = func(file string) string {
		fileName := filepath.Base(file)
		fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))
		return fileName
	}
	fileNames := lo.Map(subItems, func(item *jellyfin.ItemInfoResponse, _ int) string {
		return extractName(item.Path)
	})
	mapOrders, err := guess_epsides.MapOrders(fileNames)
	if err != nil {
		log.Printf("guess epsides for %s %s failed: %s", code, path, fileNames)
		// 使用正确的打开方式：追加模式、不存在时创建、只写模式
		file, err := os.OpenFile("error_orders.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		if _, err := file.Write([]byte(code + " \n")); err != nil {
			log.Fatal(err)
		}
	}
	subItemsInfo = lo.Map(subItems, func(item *jellyfin.ItemInfoResponse, index int) *model.ProjectInfo {
		name := extractName(item.Path)
		var zero = 0
		var order = &zero
		if mapOrders != nil {
			order = mapOrders[name]
		}
		return &model.ProjectInfo{
			ItemId:   item.Id,
			Name:     name,
			Name2:    item.Name,
			SortName: name,
			Index: lo.IfF(order != nil, func() string {
				return fmt.Sprintf("%02d", *order)
			}).Else(fmt.Sprintf("%02d", index)),
			ParentIndex: lo.If(order != nil, "1").Else("0"),
			Artist:      artist,
			Nsfw:        workInfo.Nsfw,
			Tags:        append(item.Tags, tags...),
			ReleaseDate: releaseDate,
			CreateDate:  item.DateCreated,
			Rating:      workInfo.RateAverage2Dp,
			Group: lo.If(item.AlbumArtist != "",
				lo.If(item.AlbumArtist == workInfo.Circle.Name, []string{item.AlbumArtist}).
					Else([]string{item.AlbumArtist, workInfo.Circle.Name})).
				Else([]string{workInfo.Circle.Name}),
		}
	})

	return &model.ProjectInfo{
		ItemId: "",
		Code:   code,
		Name: func() string {
			builder := strings.Builder{}
			if base != code && !lo.Contains([]string{"本編", "本编", "mp3", "MP3"}, base) {
				builder.WriteString(fmt.Sprintf("「%s」 ", base))
			}
			var prefixes = []string{"【简体中文版】", "【簡体中文版】", "【繁体中文版】"}
			var name = workInfo.Title
			for _, prefix := range prefixes {
				name = strings.TrimPrefix(name, prefix)
			}
			builder.WriteString(name)
			return builder.String()
		}(),
		Name2:           path,
		SortName:        "",
		Index:           "",
		Tags:            tags,
		ReleaseDate:     releaseDate,
		CreateDate:      item.DateCreated,
		Artist:          artist,
		People:          nil,
		Rating:          workInfo.RateAverage2Dp,
		Group:           []string{workInfo.Circle.Name},
		Nsfw:            workInfo.Nsfw,
		Price:           workInfo.Price,
		Sales:           workInfo.DlCount,
		Overview:        fmt.Sprintf(overviewTemplate, workInfo.Circle.Name, workInfo.Price, workInfo.DlCount, code, code),
		PrimaryImageUrl: workInfo.MainCoverUrl,
		ItemsInfo:       subItemsInfo,
	}, nil
}
