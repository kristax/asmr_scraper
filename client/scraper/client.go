package scraper

import (
	"asmr_scraper/client/asmr_one"
	"asmr_scraper/client/downloader"
	"asmr_scraper/client/jellyfin"
	"asmr_scraper/util/fas"
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/samber/lo"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

type client struct {
	AsmrClient     asmr_one.Client   `wire:""`
	JellyfinClient jellyfin.Client   `wire:""`
	Downloader     downloader.Client `wire:""`
	Config         *Config
	rjExp          *regexp.Regexp
}

func NewClient() Client {
	return new(client)
}

func (c *client) Init() error {
	rjExp, err := regexp.Compile("RJ\\d+")
	if err != nil {
		return err
	}
	c.rjExp = rjExp
	return nil
}

func (c *client) RefreshInfo(ctx context.Context, parentId string) (*RefreshInfoResult, error) {
	itemsResponse, err := c.JellyfinClient.GetItems(ctx, parentId, func(r *resty.Request) {
		//r.SetQueryParam("StartIndex", "0")
		//r.SetQueryParam("Limit", "1")
	})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	total := len(itemsResponse.Items)
	wg := sync.WaitGroup{}
	wg.Add(total)
	for i, item := range itemsResponse.Items {
		go func(i int, item *jellyfin.Items) {
			defer wg.Done()
			//timerStart := time.Now()
			getItem, err := c.JellyfinClient.GetItem(ctx, item.Id)
			if err != nil {
				panic(err)
			}
			rjCode := c.rjExp.FindString(getItem.Path)
			workInfo, err := c.AsmrClient.GetWorkInfo(ctx, rjCode)
			if err != nil {
				fmt.Printf("asmr one error: %s %v\n", rjCode, err)
				return
			}
			if c.Config.ForceUploadImage || !getItem.LockData {
				cover, err := c.Downloader.Download(ctx, workInfo.MainCoverUrl)
				if err != nil {
					fmt.Printf("downloader error: %s %v\n", rjCode, err)
					return
				}
				err = c.JellyfinClient.UploadPrimaryImage(ctx, item.Id, cover)
				if err != nil {
					fmt.Printf("jellyfin upload image error: %s %v\n", rjCode, err)
					return
				}
				fmt.Println("upload image for", rjCode, "success")
			}
			if c.Config.ForceUpdateInfo || !getItem.LockData {
				request := buildUpdateRequest(item.Id, rjCode, getItem.Path, workInfo)
				err = c.JellyfinClient.UpdateItem(ctx, request)
				if err != nil {
					fmt.Printf("jellyfin update item error: %s %v\n", rjCode, err)
					return
				}
				fmt.Println("update info for", rjCode, "success")
			}
			//fmt.Printf("[%d/%d] refresh %s succeed, title: %s, cost: %v\n", i+1, total, rjCode, workInfo.Title, time.Now().Sub(timerStart))
		}(i, item)
	}
	wg.Wait()
	return &RefreshInfoResult{}, nil
}

func buildUpdateRequest(itemId, rjCode, path string, workInfo *asmr_one.WorkInfoResponse) *jellyfin.UpdateItemRequest {
	base := filepath.Base(path)
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
	artist := lo.Map[*asmr_one.Vas, *jellyfin.Subject](workInfo.Vas, func(item *asmr_one.Vas, _ int) *jellyfin.Subject {
		return &jellyfin.Subject{
			Name: item.Name,
		}
	})
	var overviewTemplate = `<br>
<span>%s<span><br><br>
<span style="color: #f44336!important">%d JPY</span> &nbsp&nbsp&nbsp
<span style="color: #ffffff!important">销量: %d</span> &nbsp&nbsp&nbsp
<a href="https://www.dlsite.com/home/work/=/product_id/%s.html/?locale=zh_CN" target="_blank" style="color: #4992F2!important">DLsite</a>`
	return &jellyfin.UpdateItemRequest{
		Id: itemId,
		Name: func() string {
			builder := strings.Builder{}
			builder.WriteString(rjCode)
			builder.WriteString(workInfo.Title)
			if base != rjCode {
				builder.WriteString(fmt.Sprintf(" [%s]", base))
			}
			builder.WriteString(fmt.Sprintf(" CV: %s", strings.Join(lo.Map(workInfo.Vas, func(item *asmr_one.Vas, _ int) string {
				return item.Name
			}), ",")))
			return builder.String()
		}(),
		OriginalTitle:           path,
		ForcedSortName:          rjCode,
		CommunityRating:         fmt.Sprintf("%.1f", workInfo.RateAverage2Dp),
		CriticRating:            "",
		IndexNumber:             nil,
		AirsBeforeSeasonNumber:  "",
		AirsAfterSeasonNumber:   "",
		AirsBeforeEpisodeNumber: "",
		ParentIndexNumber:       nil,
		DisplayOrder:            "",
		Album:                   rjCode,
		AlbumArtists:            artist,
		ArtistItems:             []*jellyfin.Subject{{Name: workInfo.Circle.Name}},
		Overview:                fmt.Sprintf(overviewTemplate, workInfo.Circle.Name, workInfo.Price, workInfo.DlCount, rjCode),
		Status:                  "",
		AirDays:                 []any{},
		AirTime:                 "",
		Genres:                  tags,
		Tags:                    []string{fas.TernaryOp(workInfo.Nsfw, "R18", "全年龄")},
		Studios: []*jellyfin.Subject{
			{
				Name: workInfo.Circle.Name,
			},
		},
		PremiereDate:                 releaseDate,
		DateCreated:                  createDate,
		EndDate:                      nil,
		ProductionYear:               fmt.Sprintf("%d", releaseDate.Year()),
		AspectRatio:                  "",
		Video3DFormat:                "",
		OfficialRating:               fas.TernaryOp(workInfo.Nsfw, "XXX", "APPROVED"),
		CustomRating:                 "",
		People:                       []any{},
		LockData:                     true,
		LockedFields:                 []any{},
		ProviderIds:                  &jellyfin.ProviderIds{},
		PreferredMetadataLanguage:    "",
		PreferredMetadataCountryCode: "",
		Taglines:                     []any{},
	}
}
