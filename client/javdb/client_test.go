package javdb

import (
	"context"
	"github.com/go-kid/ioc"
	"github.com/go-kid/ioc/app"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_client_Search(t *testing.T) {
	c := NewClient()
	ioc.RunTest(t, app.SetComponents(c), app.SetConfig("../../config.yaml"))
	detail, err := c.Get(context.Background(), "STAR-907", "zh")
	//detail, err := c.Get(context.Background(), "JUKF-045", "zh")
	assert.NoError(t, err)
	log.Println(detail)
}
