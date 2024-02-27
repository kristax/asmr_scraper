package jellyfin

import (
	"asmr_scraper/util/restyop"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/go-resty/resty/v2"
	"github.com/samber/lo"
)

var (
	auth = `MediaBrowser Client="Asmr Scraper", Device="Golang", Version="0.0.1", Token="%s"`
)

type client struct {
	Cfg    *Config
	client *resty.Client
}

func NewClient() Client {
	return new(client)
}

func (c *client) Init() error {
	c.client = resty.New().
		SetBaseURL(c.Cfg.Host).
		SetDebug(c.Cfg.Debug).
		SetHeaders(map[string]string{
			"Accept-Language":      "zh,zh-TW;q=0.9,en-US;q=0.8,en;q=0.7,zh-CN;q=0.6,zh-HK;q=0.5,ja;q=0.4",
			"Connection":           "keep-alive",
			"X-Emby-Authorization": fmt.Sprintf(auth, c.Cfg.ApiKey),
			"accept":               "application/json",
		})
	return nil
}

func (c *client) GetViews(ctx context.Context) ([]*ViewItem, error) {
	var result = &ViewsResponse{}
	resp, err := c.client.R().
		SetContext(ctx).
		SetPathParam("userId", c.Cfg.UserId).
		SetResult(result).
		Get("/Users/{userId}/Views")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}
	return result.Items, nil
}

func (c *client) GetViewIdByName(ctx context.Context, name string) (string, error) {
	views, err := c.GetViews(ctx)
	if err != nil {
		return "", err
	}
	if len(views) < 1 {
		return "", errors.New("no view found")
	}
	viewItem, ok := lo.Find(views, func(item *ViewItem) bool {
		return item.Name == name
	})
	if !ok {
		return "", errors.New("no view found for " + name)
	}
	return viewItem.Id, nil
}

func (c *client) GetItems(ctx context.Context, parentId, itemType string, options ...restyop.Option) (*ItemsResponse, error) {
	var result = &ItemsResponse{}
	r := c.client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"SortBy":           "DateCreated%2CSortName",
			"SortOrder":        "Descending",
			"IncludeItemTypes": itemType,
			"Recursive":        "true",
			"Fields":           "PrimaryImageAspectRatio%2CSortName%2CBasicSyncInfo",
			"ImageTypeLimit":   "1",
			"EnableImageTypes": "Primary%2CBackdrop%2CBanner%2CThumb",
			"ParentId":         parentId,
		}).
		SetPathParam("userId", c.Cfg.UserId)
	for _, option := range options {
		option(r)
	}
	resp, err := r.
		SetResult(result).
		Get("/Users/{userId}/Items")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}
	return result, nil
}

func (c *client) GetItem(ctx context.Context, itemId string) (*ItemInfoResponse, error) {
	var result = &ItemInfoResponse{}
	resp, err := c.client.R().
		SetContext(ctx).
		SetPathParam("itemId", itemId).
		SetPathParam("userId", c.Cfg.UserId).
		SetResult(result).
		Get("/Users/{userId}/Items/{itemId}")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}
	return result, nil
}

func (c *client) UpdateItem(ctx context.Context, req *UpdateItemRequest) error {
	resp, err := c.client.R().
		SetBody(req).
		SetPathParam("itemId", req.Id).
		Post("/Items/{itemId}")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return errors.New(resp.String())
	}
	return nil
}

func (c *client) UploadPrimaryImage(ctx context.Context, itemId string, data []byte) error {
	resp, err := c.client.R().SetPathParam("itemId", itemId).
		SetHeader("Accept", "*/*").
		SetHeader("Content-Type", mimetype.Detect(data).String()).
		SetBody(base64.StdEncoding.EncodeToString(data)).
		Post("/Items/{itemId}/Images/Primary")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return errors.New(resp.String())
	}
	return nil
}
