package scraper

import (
	"asmr_scraper/client/asmr_one"
	"asmr_scraper/client/downloader"
	"asmr_scraper/client/jellyfin"
	"asmr_scraper/util/fas"
	"context"
	_ "embed"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/samber/lo"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
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
				projectInfo := GetProjectInfo("asmr_one", rjCode, workInfo)
				if projectInfo == nil {
					fmt.Printf("adaptor not found: %s", "asmr_one")
					return
				}
				projectInfo.ItemId = getItem.Id
				projectInfo.Path = getItem.Path
				request := buildUpdateRequest(projectInfo)
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

func buildUpdateRequest(project *ProjectInfo) *jellyfin.UpdateItemRequest {
	base := filepath.Base(project.Path)
	tags := project.Tags
	releaseDate := project.ReleaseDate
	createDate := project.CreateDate
	return &jellyfin.UpdateItemRequest{
		Id: project.ItemId,
		Name: func() string {
			builder := strings.Builder{}
			builder.WriteString(project.RJCode)
			builder.WriteString(project.Name)
			if base != project.RJCode {
				builder.WriteString(fmt.Sprintf(" [%s]", base))
			}
			builder.WriteString(fmt.Sprintf(" CV: %s", strings.Join(project.Artists, ",")))
			return builder.String()
		}(),
		OriginalTitle:           project.Path,
		ForcedSortName:          project.RJCode,
		CommunityRating:         fmt.Sprintf("%.1f", project.Rating),
		CriticRating:            "",
		IndexNumber:             nil,
		AirsBeforeSeasonNumber:  "",
		AirsAfterSeasonNumber:   "",
		AirsBeforeEpisodeNumber: "",
		ParentIndexNumber:       nil,
		DisplayOrder:            "",
		Album:                   project.RJCode,
		AlbumArtists: lo.Map[string, *jellyfin.Subject](project.Artists, func(item string, _ int) *jellyfin.Subject {
			return &jellyfin.Subject{Name: item}
		}),
		ArtistItems: []*jellyfin.Subject{{Name: project.Group}},
		Overview:    fmt.Sprintf(overviewTemplate, project.Group, project.Price, project.Sales, project.RJCode),
		Status:      "",
		AirDays:     []any{},
		AirTime:     "",
		Genres:      tags,
		Tags:        []string{fas.TernaryOp(project.Nsfw, "R18", "全年龄")},
		Studios: []*jellyfin.Subject{
			{
				Name: project.Group,
			},
		},
		PremiereDate:                 releaseDate,
		DateCreated:                  createDate,
		EndDate:                      nil,
		ProductionYear:               fmt.Sprintf("%d", releaseDate.Year()),
		AspectRatio:                  "",
		Video3DFormat:                "",
		OfficialRating:               fas.TernaryOp(project.Nsfw, "XXX", "APPROVED"),
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
