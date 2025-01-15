package main

import (
	"asmr_scraper/client/asmr_one"
	"asmr_scraper/client/javdb"
	"asmr_scraper/util/repository"
	"github.com/go-kid/ioc"
)

func init() {
	ioc.Register(
		asmr_one.NewClient(),
		javdb.NewClient(),
		repository.New(),
	)
}
