package scraper

import (
	"asmr_scraper/client/asmr_one"
	"asmr_scraper/client/jellyfin"
	"context"
	"github.com/go-kid/ioc"
	"github.com/go-kid/ioc/app"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_client_RefreshInfo(t *testing.T) {
	c := NewClient()
	ioc.RunTest(t, app.SetComponents(c, asmr_one.NewClient(), jellyfin.NewClient()))
	_, err := c.RefreshInfo(context.Background(), "2aa1d857635177546f8785032805c532")
	assert.NoError(t, err)
}
