package model

import (
	"context"
	"github.com/go-resty/resty/v2"
	"regexp"
)

type ClientBase struct {
	TargetName_    string `mapstructure:"targetName"`
	Order_         int    `mapstructure:"order"`
	Host           string `mapstructure:"host"`
	Debug          bool   `mapstructure:"debug"`
	ParsePathRegex string `mapstructure:"parsePathRegex"`
	Cli            *resty.Client
	reg            *regexp.Regexp
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

func (c *ClientBase) ParseCodeFromPath(ctx context.Context, path string) (string, error) {
	return c.reg.FindString(path), nil
}

func (c *ClientBase) GetProjectInfo(ctx context.Context, code string) (*ProjectInfo, error) {
	panic("implement me")
}
