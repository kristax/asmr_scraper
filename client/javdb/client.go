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
	var index = -1
	box.Find(".video-title").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		if strings.Contains(strings.ToLower(selection.Text()), strings.ToLower(avCode)) {
			index = i
			return false
		}
		return true
	})

	if index == -1 {
		return nil, fmt.Errorf("未找到该项目: %s", avCode)
	}

	var result = &ListItem{
		Code: avCode,
	}

	href, ok := box.Attr("href")
	if !ok {
		return nil, fmt.Errorf("获取Key路径失败: %s", avCode)
	}
	result.Key = href

	fmt.Println(box.Children().First().Children().Attr("src"))
	result.CoverImg, _ = box.Children().First().Children().Attr("src")
	result.Rating, result.RateCount = RateSplit(box.Find(".score .value").Text())
	result.ReleaseDate = box.Find(".meta").Text()
	return result, nil
}

func (c *client) getDetail(ctx context.Context, item *ListItem, lang string) (*Detail, error) {
	mappingHandler, ok := infoMapping[lang]
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
	if !resp.IsSuccess() {
		return &Detail{
			Code:         item.Code,
			Title:        item.Title,
			OriginTitle:  item.Title,
			CoverImg:     item.CoverImg,
			ReleasedDate: item.ReleaseDate,
			Rating:       item.Rating,
			RateCount:    item.RateCount,
		}, nil
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
		if f, ok := mappingHandler[key]; ok {
			f(result, selection.Find(".value"))
		}
	})
	return result, nil
}

func (c *client) GetProjectInfo(ctx context.Context, code string) (*model.ProjectInfo, error) {
	detail, err := c.Get(ctx, code, c.Lang)
	if err != nil {
		return nil, err
	}

	releaseDate, err := time.Parse(time.DateOnly, detail.ReleasedDate)
	if err != nil {
		releaseDate = time.Now()
	}
	return &model.ProjectInfo{
		ItemId:          "",
		Code:            detail.Code,
		Path:            "",
		Name:            detail.Title,
		Name2:           detail.OriginTitle,
		Tags:            detail.Tags,
		ReleaseDate:     releaseDate,
		CreateDate:      time.Now(),
		Artists:         detail.Actors,
		Rating:          detail.Rating,
		Group:           detail.Maker,
		Nsfw:            true,
		Price:           0,
		Sales:           0,
		Overview:        fmt.Sprintf(overviewTemplate, detail.ReleasedDate, detail.Duration, detail.Maker, detail.Publisher, detail.Series),
		PrimaryImageUrl: detail.CoverImg,
	}, nil
}
