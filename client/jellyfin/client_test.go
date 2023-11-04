package jellyfin

import (
	"context"
	"fmt"
	"github.com/go-kid/ioc"
	"github.com/go-kid/ioc/app"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func initClient(t *testing.T) Client {
	c := NewClient()
	ioc.RunTest(t, app.SetComponents(c), app.SetConfig("../../config.yaml"))
	return c
}

func Test_client_GetItems(t *testing.T) {
	c := initClient(t)
	itemsResponse, err := c.GetItems(context.Background(), "2aa1d857635177546f8785032805c532")
	assert.NoError(t, err)
	assert.Equal(t, itemsResponse.TotalRecordCount, len(itemsResponse.Items))
}

func Test_client_GetItemsLimit(t *testing.T) {
	c := initClient(t)
	itemsResponse, err := c.GetItems(context.Background(), "2aa1d857635177546f8785032805c532", func(r *resty.Request) {
		r.SetQueryParam("StartIndex", "0")
		r.SetQueryParam("Limit", "1")
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(itemsResponse.Items))
}

func Test_client_GetItem(t *testing.T) {
	c := initClient(t)
	itemInfoResponse, err := c.GetItem(context.Background(), "b808921b42453e40001fe006204aa830")
	assert.NoError(t, err)
	fmt.Println(itemInfoResponse.Path)
}
