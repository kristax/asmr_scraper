package downloader

import (
	"context"
	"github.com/go-resty/resty/v2"
)

type client struct {
	client *resty.Client
}

func NewClient() Client {
	return new(client)
}

func (c *client) Init() error {
	c.client = resty.New()
	return nil
}

func (c *client) Download(ctx context.Context, url string) ([]byte, error) {
	response, err := c.client.R().SetContext(ctx).Get(url)
	if err != nil {
		return nil, err
	}
	return response.Body(), nil
}
