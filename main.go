package main

import (
	"asmr_scraper/client/scraper"
	"asmr_scraper/model"
	"bufio"
	"context"
	"fmt"
	"github.com/go-kid/ioc"
	ap "github.com/go-kid/ioc/app"
	"os"
)

type App struct {
	Targets []*model.Target `prop:"App.targets"`
	Scraper scraper.Client  `wire:""`
}

func main() {
	var app = &App{}
	_, err := ioc.Run(ap.SetConfig("config.yaml"), ap.SetComponents(app))
	if err != nil {
		panic(err)
	}

	for _, parent := range app.Targets {
		if parent.Disable {
			continue
		}
		_, err = app.Scraper.RefreshInfo(context.Background(), parent)
		if err != nil {
			fmt.Printf("refresh target %s:%s failed %v\n", parent.Id, parent.Name, err)
			continue
		}
	}
	fmt.Println("refresh finished, press any key to exit")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
}
