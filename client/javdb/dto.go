package javdb

import (
	"asmr_scraper/model"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-kid/ioc/util/fas"
	"github.com/samber/lo"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ListItem struct {
	Code        string
	Key         string
	Title       string
	CoverImg    string
	Rating      float64
	RateCount   int
	ReleaseDate string
}

type Detail struct {
	Code         string
	Title        string
	OriginTitle  string
	CoverImg     string
	ReleasedDate string
	Duration     string
	Maker        string
	Publisher    string
	Series       string
	Rating       float64
	RateCount    int
	Tags         []string
	Actors       []*Actor
	Director     *Actor
}

func (d *Detail) ToProjectInfo(code, path string) (*model.ProjectInfo, error) {
	base := filepath.Base(path)
	releaseDate, err := time.Parse(time.DateOnly, d.ReleasedDate)
	if err != nil {
		releaseDate = time.Now()
	}
	return &model.ProjectInfo{
		ItemId: "",
		Code:   d.Code,
		Name: func() string {
			builder := strings.Builder{}
			if base != code {
				builder.WriteString(fmt.Sprintf("「%s」 ", base))
			}
			builder.WriteString(d.Title)
			return builder.String()
		}(),
		Name2:       d.OriginTitle,
		Tags:        d.Tags,
		ReleaseDate: releaseDate,
		CreateDate:  time.Now(),
		People: func() []*model.People {
			var people []*model.People
			if len(d.Actors) > 0 {
				people = append(people, lo.Map(d.Actors, func(item *Actor, index int) *model.People {
					return &model.People{
						Name:     item.Name,
						Type:     model.TypeActor,
						Role:     fas.TernaryOp(item.Gender == "male", "男演员", "女优"),
						Gender:   item.Gender,
						HomePage: item.HomePage,
					}
				})...)
			}
			if d.Director != nil {
				people = append(people, &model.People{
					Name:     d.Director.Name,
					Type:     model.TypeDirector,
					HomePage: d.Director.HomePage,
				})
			}
			return people
		}(),
		Rating:          d.Rating,
		Group:           d.Maker,
		Nsfw:            true,
		Price:           0,
		Sales:           0,
		Overview:        fmt.Sprintf(overviewTemplate, d.ReleasedDate, d.Duration, d.Maker, d.Publisher, d.Series),
		PrimaryImageUrl: d.CoverImg,
	}, nil
}

type Actor struct {
	Name     string
	Gender   string
	HomePage string
}

func (d *Detail) String() string {
	bytes, _ := json.MarshalIndent(d, "", "  ")
	return string(bytes)
}

type Category string

const (
	CategoryCode        Category = "code"
	CategoryReleaseDate Category = "releaseDate"
	CategoryDuration    Category = "duration"
	CategoryMaker       Category = "maker"
	CategoryPublisher   Category = "publisher"
	CategorySeries      Category = "series"
	CategoryRating      Category = "rating"
	CategoryTags        Category = "tags"
	CategoryActors      Category = "actors"
	CategoryDirector    Category = "director"
)

var infoMapping = map[string]map[string]Category{
	"en": {
		"ID:":            CategoryCode,
		"Released Date:": CategoryReleaseDate,
		"Duration:":      CategoryDuration,
		"Maker:":         CategoryMaker,
		"Publisher:":     CategoryPublisher,
		"Series:":        CategorySeries,
		"Rating:":        CategoryRating,
		"Tags:":          CategoryTags,
		"Actor(s):":      CategoryActors,
		"Director:":      CategoryDirector,
	},
	"zh": {
		"番號:": CategoryCode,
		"日期:": CategoryReleaseDate,
		"時長:": CategoryDuration,
		"片商:": CategoryMaker,
		"發行:": CategoryPublisher,
		"系列:": CategorySeries,
		"評分:": CategoryRating,
		"類別:": CategoryTags,
		"演員:": CategoryActors,
		"導演:": CategoryDirector,
	},
}

var handlerMapping = map[Category]func(detail *Detail, value *goquery.Selection){
	CategoryCode: func(detail *Detail, val *goquery.Selection) {
		detail.Code = val.Text()
	},
	CategoryReleaseDate: func(detail *Detail, val *goquery.Selection) {
		detail.ReleasedDate = val.Text()
	},
	CategoryDuration: func(detail *Detail, val *goquery.Selection) {
		detail.Duration = val.Text()
	},
	CategoryMaker: func(detail *Detail, val *goquery.Selection) {
		detail.Maker = val.Text()
	},
	CategoryPublisher: func(detail *Detail, val *goquery.Selection) {
		detail.Publisher = val.Text()
	},
	CategorySeries: func(detail *Detail, val *goquery.Selection) {
		detail.Series = val.Text()
	},
	CategoryRating: func(detail *Detail, val *goquery.Selection) {
		detail.Rating, detail.RateCount = RateSplit(val.Text())
	},
	CategoryTags: func(detail *Detail, val *goquery.Selection) {
		tags := strings.Split(val.Text(), ",")
		detail.Tags = lo.Map(tags, func(item string, _ int) string {
			return strings.TrimSpace(item)
		})
	},
	CategoryActors: func(detail *Detail, val *goquery.Selection) {
		var actors []*Actor
		val.Children().Each(func(i int, selection *goquery.Selection) {
			if i%2 == 0 {
				url, _ := selection.Attr("href")
				actors = append(actors, &Actor{
					Name:     selection.Text(),
					Gender:   "",
					HomePage: url,
				})
			} else {
				switch selection.Text() {
				case "♀":
					actors[len(actors)-1].Gender = "female"
				case "♂":
					actors[len(actors)-1].Gender = "male"
				}
			}
		})
		detail.Actors = actors
	},
	CategoryDirector: func(detail *Detail, value *goquery.Selection) {
		url, _ := value.Children().Attr("href")
		detail.Director = &Actor{
			Name:     value.Text(),
			Gender:   "",
			HomePage: url,
		}
	},
}

var rateRegex, _ = regexp.Compile("(-?\\d+)(\\.\\d+)?")

func RateSplit(val string) (rating float64, rateCount int) {
	rates := rateRegex.FindAllString(val, -1)
	if len(rates) > 0 {
		rating, _ = strconv.ParseFloat(rates[0], 64)
	}
	if len(rates) > 1 {
		rateCount, _ = strconv.Atoi(rates[1])
	}
	return
}
