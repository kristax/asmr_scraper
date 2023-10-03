package asmr_one

import (
	"context"
	"github.com/go-kid/ioc"
	"github.com/go-kid/ioc/app"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_client_GetWorkInfo(t *testing.T) {
	var c = NewClient()
	ioc.RunTest(t, app.SetComponents(c), app.SetDefaultConfigure(), app.SetConfig("../../config.yaml"))
	workInfo, err := c.GetWorkInfo(context.Background(), "RJ438754")
	assert.NoError(t, err)
	assert.Equal(t, 438754, workInfo.Id)
}
