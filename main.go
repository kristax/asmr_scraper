package main

import (
	"asmr_scraper/client/scraper"
	"asmr_scraper/model"
	"bufio"
	"context"
	"github.com/go-kid/ioc"
	ap "github.com/go-kid/ioc/app"
	"github.com/gosuri/uiprogress"
	"log"
	"os"
	"sync"
)

type App struct {
	Targets []*model.Target `prop:"App.targets"`
	Scraper scraper.Client  `wire:""`
}

func main() {
	uiprogress.Start()
	var app = &App{}
	_, err := ioc.Run(ap.SetConfig("config.yaml"), ap.SetComponents(app))
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(len(app.Targets))
	for _, parent := range app.Targets {
		go func(parent *model.Target) {
			defer wg.Done()
			if parent.Disable {
				return
			}
			_, err = app.Scraper.RefreshInfo(context.Background(), parent)
			if err != nil {
				log.Printf("refresh target %s failed: %v\n", parent.Name, err)
				return
			}
		}(parent)
	}

	wg.Wait()
	uiprogress.Stop()

	log.Println("refresh finished, press Enter to exit")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
}
