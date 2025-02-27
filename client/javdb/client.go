package javdb

import (
	"asmr_scraper/model"
	"bytes"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"time"
)

type client struct {
	model.ClientBase `prop:"Clients.JavDBConfig"`
	Lang             string `prop:"Clients.JavDBConfig.lang"`
}

func (c *client) ClientID() string {
	return "JavDB"
}

func (c *client) DataModel() model.ProjectInfoData {
	return &Detail{}
}

func NewClient() Client {
	return &client{}
}

func (c *client) Init() error {
	return c.InitClient()
}

func (c *client) Get(ctx context.Context, avCode, lang string) (*Detail, error) {
	search, err := c.search(ctx, avCode)
	if err != nil {
		return nil, err
	}
	return c.getDetail(ctx, search, lang)
}

func (c *client) search(ctx context.Context, avCode string) (*ListItem, error) {
	resp, err := c.Cli.R().
		SetContext(ctx).
		SetQueryParam("q", avCode).
		SetQueryParam("f", "all").
		Get("/search")
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(resp.Body()))
	if err != nil {
		return nil, err
	}
	box := doc.Find(".movie-list .item .box")
	matchedItem := box.FilterFunction(func(i int, selection *goquery.Selection) bool {
		return strings.Contains(strings.ToLower(selection.Find(".video-title").Text()), strings.ToLower(avCode))
	}).First()

	if matchedItem == nil || matchedItem.Text() == "" {
		return nil, fmt.Errorf("未找到该项目: %s", avCode)
	}

	var result = &ListItem{
		Code: avCode,
	}

	href, ok := matchedItem.Attr("href")
	if !ok {
		return nil, fmt.Errorf("获取Key路径失败: %s", avCode)
	}
	result.Key = href
	result.Title, _ = matchedItem.Attr("title")

	matchedItem.Find(".cover").Children().EachWithBreak(func(i int, selection *goquery.Selection) bool {
		img, exists := selection.Attr("src")
		if exists {
			result.CoverImg = img
			return false
		}
		return true
	})

	result.Rating, result.RateCount = RateSplit(matchedItem.Find(".score .value").Text())
	result.ReleaseDate = strings.ReplaceAll(strings.ReplaceAll(matchedItem.Find(".meta").First().Text(), " ", ""), "\n", "")
	return result, nil
}

func (c *client) getDetail(ctx context.Context, item *ListItem, lang string) (*Detail, error) {
	if strings.HasPrefix(strings.ToUpper(item.Code), "FC2") {
		return &Detail{
			Code:         item.Code,
			Title:        item.Title,
			OriginTitle:  item.Title,
			CoverImg:     item.CoverImg,
			ReleasedDate: item.ReleaseDate,
			Duration:     "",
			Maker:        "FC2",
			Publisher:    "FC2",
			Series:       "",
			Rating:       item.Rating,
			RateCount:    item.RateCount,
			Tags:         []string{"素人"},
			Actors:       []*Actor{},
		}, nil
	}
	categoryMapping, ok := infoMapping[lang]
	if !ok {
		return nil, fmt.Errorf("language %s not support", lang)
	}
	resp, err := c.Cli.R().
		SetContext(ctx).
		SetQueryParam("locale", lang).
		Get(item.Key)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(resp.Body()))
	if err != nil {
		return nil, err
	}

	var result = &Detail{}
	videoDetail := doc.Find(".video-detail")

	title := videoDetail.Find(".title").First()
	result.Title = title.Find(".current-title").Text()
	result.OriginTitle = title.Find(".origin-title").Text()

	meta := videoDetail.Find(".video-meta-panel")
	cover, _ := meta.Find(".video-cover").Attr("src")
	result.CoverImg = cover

	meta.Find(".movie-panel-info .panel-block").Each(func(i int, selection *goquery.Selection) {
		key := selection.Children().First().Text()
		if category, ok := categoryMapping[key]; ok {
			handlerMapping[category](result, selection.Find(".value"))
		}
	})
	return result, nil
}

func (c *client) GetData(ctx context.Context, code string) (model.ProjectInfoData, error) {
	detail, err := c.Get(ctx, code, c.Lang)
	if err != nil {
		return nil, err
	}
	time.Sleep(time.Second * 5)
	return detail, nil
}
