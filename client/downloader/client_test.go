package downloader

import (
	"context"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/go-kid/ioc"
	"github.com/go-kid/ioc/app"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_client_Download(t *testing.T) {
	c := NewClient()
	ioc.RunTest(t, app.SetComponents(c))
	bytes, err := c.Download(context.Background(), "https://api.asmr-200.com/api/cover/RJ01096697.jpg?type=main")
	assert.NoError(t, err)
	detect := mimetype.Detect(bytes)
	fmt.Println(detect.String())
}
