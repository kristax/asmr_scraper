package asmr_one

import (
	"asmr_scraper/model"
	"context"
	"errors"
	"fmt"
	"github.com/samber/lo"
	"strings"
	"time"
)

type client struct {
	model.ClientBase `prop:"Clients.AsmrOneConfig"`
}

func NewClient() Client {
	return new(client)
}

func (c *client) Init() error {
	err := c.InitClient()
	if err != nil {
		return err
	}
	c.Cli.SetHeaders(map[string]string{
		"accept":          "application/json, text/plain, */*",
		"accept-language": "zh,zh-TW;q=0.9,en-US;q=0.8,en;q=0.7,zh-CN;q=0.6,zh-HK;q=0.5,ja;q=0.4",
		"user-agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36",
	})
	return nil
}

func (c *client) GetWorkInfo(ctx context.Context, rj string) (*WorkInfoResponse, error) {
	var result = &WorkInfoResponse{}
	resp, err := c.Cli.R().
		SetContext(ctx).
		SetPathParam("rjCode", strings.TrimPrefix(rj, "RJ")).
		SetResult(result).
		Get("/api/workInfo/{rjCode}")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}
	return result, nil
}

func (c *client) GetProjectInfo(ctx context.Context, code string) (*model.ProjectInfo, error) {
	workInfo, err := c.GetWorkInfo(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("asmr one error: %s %v", code, err)
	}
	return asmrOneAdaptor(code, workInfo), nil
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
	artist := lo.Map(workInfo.Vas, func(item *Vas, _ int) *model.People {
		return &model.People{
			Name:     item.Name,
			Type:     model.TypeActor,
			Role:     "CV",
			Gender:   "female",
			HomePage: "",
		}
	})
	return &model.ProjectInfo{
		ItemId:          "",
		Code:            rjCode,
		Path:            "",
		Name:            workInfo.Title,
		Tags:            tags,
		ReleaseDate:     releaseDate,
		CreateDate:      createDate,
		People:          artist,
		Rating:          workInfo.RateAverage2Dp,
		Group:           workInfo.Circle.Name,
		Nsfw:            workInfo.Nsfw,
		Price:           workInfo.Price,
		Sales:           workInfo.DlCount,
		Overview:        fmt.Sprintf(overviewTemplate, workInfo.Circle.Name, workInfo.Price, workInfo.DlCount, rjCode),
		PrimaryImageUrl: workInfo.MainCoverUrl,
	}
}
