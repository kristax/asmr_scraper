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
	"sync"
	"time"
)

type client struct {
	//AsmrClient     asmr_one.Client   `wire:""`
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
	itemsResponse, err := c.JellyfinClient.GetItems(ctx, target.Id, target.Type, func(r *resty.Request) {
		if c.Config.Query.StartIndex == 0 && c.Config.Query.Limit == 0 {
			return
		}
		r.SetQueryParam("StartIndex", strconv.Itoa(c.Config.Query.StartIndex))
		r.SetQueryParam("Limit", strconv.Itoa(c.Config.Query.Limit))
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
	if total < 1 {
		return nil, nil
	}
	if target.Async {
		wg := sync.WaitGroup{}
		wg.Add(total)
		for i, item := range itemsResponse.Items {
			go func(i int, item *jellyfin.Items) {
				defer wg.Done()
				getItem, err := c.JellyfinClient.GetItem(ctx, item.Id)
				if err != nil {
					log.Fatal(err)
				}

				if getItem.LockData && !c.Config.ForceUploadImage && !c.Config.ForceUpdateInfo {
					return
				}
				err = c.UpdateInfo(ctx, getItem, clients)
				if err != nil {
					fmt.Printf("update item failed: %v\n", err)
				}
			}(i, item)
		}
		wg.Wait()
	} else {
		for _, item := range itemsResponse.Items {
			getItem, err := c.JellyfinClient.GetItem(ctx, item.Id)
			if err != nil {
				return nil, err
			}

			if getItem.LockData && !c.Config.ForceUploadImage && !c.Config.ForceUpdateInfo {
				continue
			}
			err = c.UpdateInfo(ctx, getItem, clients)
			if err != nil {
				fmt.Printf("update item failed: %v\n", err)
			}
			time.Sleep(time.Second * time.Duration(target.Jitter))
		}
	}

	return &RefreshInfoResult{}, nil
}

func (c *client) UpdateInfo(ctx context.Context, getItem *jellyfin.ItemInfoResponse, clients []SourceClient) error {
	var projectInfo *model.ProjectInfo
	for _, cli := range clients {
		clientId := reflectx.Id(cli)
		code, err := cli.ParseCodeFromPath(ctx, getItem.Path)
		if err != nil || code == "" {
			fmt.Printf("client %s parse code from path %s failed: %v\n", clientId, getItem.Path, err)
			continue
		}
		projectInfo, err = cli.GetProjectInfo(ctx, code)
		if err != nil || projectInfo == nil {
			fmt.Printf("client %s get %s project info failed: %v\n", code, clientId, err)
			continue
		}
		projectInfo.ItemId = getItem.Id
		projectInfo.Path = getItem.Path
		projectInfo.Code = code
		break
	}
	if projectInfo == nil {
		return fmt.Errorf("all clients failed to build project info")
	}

	if c.Config.ForceUploadImage || !getItem.LockData {
		cover, err := c.Downloader.Download(ctx, projectInfo.PrimaryImageUrl)
		if err != nil {
			fmt.Printf("downloader error: %s %v\n", projectInfo.Code, err)
			goto STAGE2
		}
		err = c.JellyfinClient.UploadPrimaryImage(ctx, projectInfo.ItemId, cover)
		if err != nil {
			fmt.Printf("jellyfin upload image error: %s %v\n", projectInfo.Code, err)
			goto STAGE2
		}
		fmt.Println("upload image for", projectInfo.Code, "success")
	}
STAGE2:
	if c.Config.ForceUpdateInfo || !getItem.LockData {
		err := c.JellyfinClient.UpdateItem(ctx, projectInfo.ToJellyfinUpdateItemReq())
		if err != nil {
			fmt.Printf("jellyfin update item error: %s %v\n", projectInfo.Code, err)
			return nil
		}
		fmt.Println("update info for", projectInfo.Code, "success")
	}
	return nil
}
