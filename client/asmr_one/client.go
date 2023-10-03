package asmr_one

import (
	"context"
	"errors"
	"github.com/go-resty/resty/v2"
	"strings"
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
			"accept":          "application/json, text/plain, */*",
			"accept-language": "zh,zh-TW;q=0.9,en-US;q=0.8,en;q=0.7,zh-CN;q=0.6,zh-HK;q=0.5,ja;q=0.4",
			"user-agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36",
		})
	return nil
}

func (c *client) GetWorkInfo(ctx context.Context, rj string) (*WorkInfoResponse, error) {
	var result = &WorkInfoResponse{}
	resp, err := c.client.R().
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
