package model

import (
	"asmr_scraper/client/jellyfin"
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"regexp"
)

type ClientBase struct {
	TargetName_       string `mapstructure:"targetName"`
	Order_            int    `mapstructure:"order"`
	Host              string `mapstructure:"host"`
	Debug             bool   `mapstructure:"debug"`
	ParsePathRegex    string `mapstructure:"parsePathRegex"`
	ForceMissingInfo  bool   `mapstructure:"forceMissingInfo"`
	ForceMissingImage bool   `mapstructure:"forceMissingImage"`
	Cli               *resty.Client
	reg               *regexp.Regexp
}

func (c *ClientBase) InitClient() error {
	c.Cli = resty.New().
		SetBaseURL(c.Host).
		SetDebug(c.Debug)
	reg, err := regexp.Compile(c.ParsePathRegex)
	if err != nil {
		return err
	}
	c.reg = reg
	return nil
}

func (c *ClientBase) TargetName() string {
	return c.TargetName_
}

func (c *ClientBase) Order() int {
	return c.Order_
}

func (c *ClientBase) ParseCode(ctx context.Context, item *jellyfin.ItemInfoResponse) (string, error) {
	var sources = []string{
		item.Name, item.Path, item.SortName, item.ForcedSortName, item.OriginalTitle,
	}
	for _, source := range sources {
		code := c.reg.FindString(source)
		if code != "" {
			return code, nil
		}
	}
	return "", fmt.Errorf("parse code failed from sources: %v", sources)
}

func (c *ClientBase) GetProjectInfo(ctx context.Context, code string) (*ProjectInfo, error) {
	panic("implement me")
}

func (c *ClientBase) ImageMissing(item *jellyfin.ItemInfoResponse) bool {
	return c.ForceMissingImage || item.ImageTags.Primary == ""
}

func (c *ClientBase) InfoMissing(item *jellyfin.ItemInfoResponse) bool {
	return c.ForceMissingInfo || item.OriginalTitle == ""
}
