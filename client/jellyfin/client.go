package jellyfin

import (
	"asmr_scraper/util/restyop"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/go-resty/resty/v2"
)

var (
	auth = `MediaBrowser Client="Asmr Scraper", Device="Golang", Version="0.0.1", Token="%s"`
	//userId = "8e7118d28d104dd6acdad76c7ac58647"
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

func (c *client) GetItems(ctx context.Context, parentId string, options ...restyop.Option) (*ItemsResponse, error) {
	var result = &ItemsResponse{}
	r := c.client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"SortBy":           "DateCreated%2CSortName",
			"SortOrder":        "Descending",
			"IncludeItemTypes": "MusicAlbum",
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
