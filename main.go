package main

import (
	"asmr_scraper/client/scraper"
	"context"
	"github.com/go-kid/ioc"
	ap "github.com/go-kid/ioc/app"
)

type app struct {
	Scraper scraper.Client `wire:""`
}

var App = &app{}

func init() {
	ioc.Register(App)
}

func main() {
	err := ioc.Run(ap.SetDefaultConfigure(), ap.SetConfig("config.yaml"))
	if err != nil {
		panic(err)
	}
	_, err = App.Scraper.RefreshInfo(context.Background(), "2aa1d857635177546f8785032805c532")
	if err != nil {
		panic(err)
	}
}
