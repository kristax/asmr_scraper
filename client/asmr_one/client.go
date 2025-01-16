package asmr_one

import (
	"asmr_scraper/client/jellyfin"
	"asmr_scraper/model"
	"context"
	"errors"
	"fmt"
	"strings"
)

type client struct {
	model.ClientBase `prop:"Clients.AsmrOneConfig"`
}

func (c *client) ClientID() string {
	return "AsmrOne"
}

func (c *client) DataModel() model.ProjectInfoData {
	return &WorkInfoResponse{}
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
func (c *client) InfoMissing(item *jellyfin.ItemInfoResponse) bool {
	return c.ClientBase.InfoMissing(item) || item.AlbumArtist == "" || !strings.HasPrefix(item.Overview, "<div>")
}

func (c *client) GetData(ctx context.Context, code string) (model.ProjectInfoData, error) {
	workInfo, err := c.GetWorkInfo(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("asmr one error: %s %v", code, err)
	}
	return workInfo, nil
}
