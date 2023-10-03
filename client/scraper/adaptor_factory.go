package scraper

import (
	"asmr_scraper/client/asmr_one"
	"github.com/samber/lo"
	"time"
)

func GetProjectInfo(source, rjCode string, workInfo any) *ProjectInfo {
	switch source {
	case "asmr_one":
		return asmrOneAdaptor(rjCode, workInfo.(*asmr_one.WorkInfoResponse))
	}
	return nil
}

func asmrOneAdaptor(rjCode string, workInfo *asmr_one.WorkInfoResponse) *ProjectInfo {
	tags := lo.Map[*asmr_one.Tag, string](workInfo.Tags, func(item *asmr_one.Tag, _ int) string {
		return item.Name
	})
	releaseDate, err := time.Parse(time.DateOnly, workInfo.Release)
	if err != nil {
		releaseDate = time.Now()
	}
	createDate, err := time.Parse(time.DateOnly, workInfo.CreateDate)
	if err != nil {
		createDate = time.Now()
	}
	artist := lo.Map[*asmr_one.Vas, string](workInfo.Vas, func(item *asmr_one.Vas, _ int) string {
		return item.Name
	})
	return &ProjectInfo{
		RJCode:      rjCode,
		Name:        workInfo.Title,
		Tags:        tags,
		ReleaseDate: releaseDate,
		CreateDate:  createDate,
		Artists:     artist,
		Rating:      workInfo.RateAverage2Dp,
		Group:       workInfo.Circle.Name,
		Nsfw:        workInfo.Nsfw,
		Price:       workInfo.Price,
		Sales:       workInfo.DlCount,
	}
}
