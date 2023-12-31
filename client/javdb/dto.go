package javdb

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/samber/lo"
	"regexp"
	"strconv"
	"strings"
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
