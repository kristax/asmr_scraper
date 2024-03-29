package scraper

import (
	"asmr_scraper/client/downloader"
	"asmr_scraper/client/jellyfin"
	"asmr_scraper/model"
	"context"
	_ "embed"
	"fmt"
	"github.com/go-kid/ioc/util/reflectx"
	"github.com/go-resty/resty/v2"
	"github.com/samber/lo"
	"log"
	"sort"
	"strconv"
)

type client struct {
	JellyfinClient jellyfin.Client   `wire:""`
	Clients        []SourceClient    `wire:""`
	Downloader     downloader.Client `wire:""`
	Config         *Config

	clientsMap map[string][]SourceClient
}

func NewClient() Client {
	return new(client)
}

func (c *client) Init() error {
	c.clientsMap = lo.GroupBy(c.Clients, func(item SourceClient) string {
		return item.TargetName()
	})
	for key, _ := range c.clientsMap {
		sort.Slice(c.clientsMap[key], func(i, j int) bool {
			return c.clientsMap[key][i].Order() < c.clientsMap[key][j].Order()
		})
	}
	return nil
}

func (c *client) RefreshInfo(ctx context.Context, target *model.Target) (*RefreshInfoResult, error) {
	clients, ok := c.clientsMap[target.Name]
	if !ok {
		return nil, fmt.Errorf("none source clients for %s", target.Name)
	}
	id, err := c.JellyfinClient.GetViewIdByName(ctx, target.Name)
	if err != nil {
		return nil, err
	}
	itemsResponse, err := c.JellyfinClient.GetItems(ctx, id, target.Type, func(r *resty.Request) {
		if c.Config.Query.StartIndex == 0 && c.Config.Query.Limit == 0 {
			return
		}
		r.SetQueryParam("StartIndex", strconv.Itoa(c.Config.Query.StartIndex))
		r.SetQueryParam("Limit", strconv.Itoa(c.Config.Query.Limit))
	})
	if err != nil {
		return nil, err
	}
	log.Printf("start refresh media: %s (%s)", target.Name, id)
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	total := len(itemsResponse.Items)
	if total < 1 {
		return nil, nil
	}

	var wrapperChan = make(chan *projectWrapper)

	c.getProjectsInfo(ctx, itemsResponse.Items, clients, wrapperChan)
	c.updateInfos(ctx, wrapperChan)

	return &RefreshInfoResult{}, nil
}

type projectWrapper struct {
	projectInfo  *model.ProjectInfo
	missingImage bool
	missingInfo  bool
}

func (c *client) getProjectsInfo(ctx context.Context, items []*jellyfin.Items, clients []SourceClient, wrapperChan chan *projectWrapper) {
	go func() {
		for _, respItem := range items {
			item, err := c.JellyfinClient.GetItem(ctx, respItem.Id)
			if err != nil {
				log.Printf("jelly client get item error: %v %#v\n", err, item.Path)
				continue
			}
			var (
				projectInfo  *model.ProjectInfo
				missingImage bool
				missingInfo  bool
			)
			for _, cli := range clients {
				clientId := reflectx.Id(cli)
				missingImage = cli.ImageMissing(item)
				missingInfo = cli.InfoMissing(item)
				if !missingImage && !missingInfo {
					continue
				}
				log.Printf("client %s parse code start\n", clientId)
				code, err := cli.ParseCode(ctx, item)
				if err != nil {
					log.Printf("client %s parse code failed: %v\n", clientId, err)
					continue
				}
				log.Printf("client %s parse code success: %s\n", clientId, code)
				log.Printf("client %s get project info start\n", clientId)
				projectInfo, err = cli.GetProjectInfo(ctx, code)
				if err != nil {
					log.Printf("client %s get %s project info failed: %v\n", code, clientId, err)
					continue
				}
				log.Printf("client %s get project info success: %s\n", clientId, projectInfo.Name)
				projectInfo.ItemId = item.Id
				projectInfo.Path = item.Path
				projectInfo.Code = code
				break
			}
			if projectInfo == nil && (missingInfo || missingImage) {
				log.Printf("all clients failed to build project info for %s", item.Path)
				continue
			}
			wrapperChan <- &projectWrapper{
				projectInfo:  projectInfo,
				missingImage: missingImage,
				missingInfo:  missingInfo,
			}
		}
		close(wrapperChan)
	}()
}

func (c *client) updateInfos(ctx context.Context, wrapperChan chan *projectWrapper) {
	for wrapper := range wrapperChan {
		if wrapper.missingImage {
			cover, err := c.Downloader.Download(ctx, wrapper.projectInfo.PrimaryImageUrl)
			if err != nil {
				log.Printf("downloader error: %s %v\n", wrapper.projectInfo.Code, err)
				goto STAGE2
			}
			err = c.JellyfinClient.UploadPrimaryImage(ctx, wrapper.projectInfo.ItemId, cover)
			if err != nil {
				log.Printf("jellyfin upload image error: %s %v\n", wrapper.projectInfo.Code, err)
				goto STAGE2
			}
			log.Println("upload image for", wrapper.projectInfo.Code, "success")
		}
	STAGE2:
		if wrapper.missingInfo {
			err := c.JellyfinClient.UpdateItem(ctx, wrapper.projectInfo.ToJellyfinUpdateItemReq())
			if err != nil {
				log.Printf("jellyfin update item error: %s %v\n", wrapper.projectInfo.Code, err)
				return
			}
			log.Println("update info for", wrapper.projectInfo.Code, "success")
		}
	}
}
