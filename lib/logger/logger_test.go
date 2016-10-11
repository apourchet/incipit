package logger

import (
	"testing"

	"golang.org/x/net/context"

	"github.com/apourchet/incipit/lib/utils"
	protos "github.com/apourchet/incipit/protos/go"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestZeroOriginal(t *testing.T) {
	if utils.InKubernetes() {
		return
	}
	c, err := GetDefaultClient()
	assert.Nil(t, err)
	k := uuid.NewV4().String()
	res, err := c.LogLogin(context.Background(), &protos.LogLoginReq{k})
	assert.Nil(t, err)
	assert.Equal(t, res.LastLogin, 0)
}

func TestReplaceOriginal(t *testing.T) {
	if utils.InKubernetes() {
		return
	}
	c, err := GetDefaultClient()
	assert.Nil(t, err)
	k := uuid.NewV4().String()
	_, err = c.LogLogin(context.Background(), &protos.LogLoginReq{k})
	assert.Nil(t, err)
	res, err := c.LogLogin(context.Background(), &protos.LogLoginReq{k})
	assert.Nil(t, err)
	assert.NotEqual(t, res.LastLogin, 0)
}
