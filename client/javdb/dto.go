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
	Actors       []string
}

func (d *Detail) String() string {
	bytes, _ := json.MarshalIndent(d, "", "  ")
	return string(bytes)
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

var infoMapping = map[string]map[string]func(detail *Detail, value *goquery.Selection){
	"en": {
		"ID:": func(detail *Detail, val *goquery.Selection) {
			detail.Code = val.Text()
		},
		"Released Date:": func(detail *Detail, val *goquery.Selection) {
			detail.ReleasedDate = val.Text()
		},
		"Duration:": func(detail *Detail, val *goquery.Selection) {
			detail.Duration = val.Text()
		},
		"Maker:": func(detail *Detail, val *goquery.Selection) {
			detail.Maker = val.Text()
		},
		"Publisher:": func(detail *Detail, val *goquery.Selection) {
			detail.Publisher = val.Text()
		},
		"Series:": func(detail *Detail, val *goquery.Selection) {
			detail.Series = val.Text()
		},
		"Rating:": func(detail *Detail, val *goquery.Selection) {
			detail.Rating, detail.RateCount = RateSplit(val.Text())
		},
		"Tags:": func(detail *Detail, val *goquery.Selection) {
			tags := strings.Split(val.Text(), ",")
			detail.Tags = lo.Map(tags, func(item string, _ int) string {
				return strings.TrimSpace(item)
			})
		},
		"Actor(s):": func(detail *Detail, val *goquery.Selection) {
			actors := val.Children().FilterFunction(func(i int, selection *goquery.Selection) bool {
				return i%2 == 0
			}).Map(func(i int, selection *goquery.Selection) string {
				return selection.Text()
			})
			detail.Actors = actors
		},
	},
	"zh": {
		"番號:": func(detail *Detail, val *goquery.Selection) {
			detail.Code = val.Text()
		},
		"日期:": func(detail *Detail, val *goquery.Selection) {
			detail.ReleasedDate = val.Text()
		},
		"時長:": func(detail *Detail, val *goquery.Selection) {
			detail.Duration = val.Text()
		},
		"片商:": func(detail *Detail, val *goquery.Selection) {
			detail.Maker = val.Text()
		},
		"發行:": func(detail *Detail, val *goquery.Selection) {
			detail.Publisher = val.Text()
		},
		"系列:": func(detail *Detail, val *goquery.Selection) {
			detail.Series = val.Text()
		},
		"評分:": func(detail *Detail, val *goquery.Selection) {
			rates := rateRegex.FindAllString(val.Text(), -1)
			if len(rates) > 0 {
				detail.Rating, _ = strconv.ParseFloat(rates[0], 64)
			}
			if len(rates) > 1 {
				detail.RateCount, _ = strconv.Atoi(rates[1])
			}
		},
		"類別:": func(detail *Detail, val *goquery.Selection) {
			tags := strings.Split(val.Text(), ",")
			detail.Tags = lo.Map(tags, func(item string, _ int) string {
				return strings.TrimSpace(item)
			})
		},
		"演員:": func(detail *Detail, val *goquery.Selection) {
			actors := val.Children().FilterFunction(func(i int, selection *goquery.Selection) bool {
				return i%2 == 0
			}).Map(func(i int, selection *goquery.Selection) string {
				return selection.Text()
			})
			detail.Actors = actors
		},
	},
}
