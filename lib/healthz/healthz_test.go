package healthz

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSpawn(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	SpawnHealthCheck(12345)
	time.Sleep(5 * time.Millisecond)
	testRequest(t, "http://localhost:12345/healthz")
}

func testRequest(t *testing.T, url string) {
	resp, err := http.Get(url)
	defer resp.Body.Close()
	assert.NoError(t, err)

	body, ioerr := ioutil.ReadAll(resp.Body)
	assert.NoError(t, ioerr)
	assert.Equal(t, "", string(body), "resp body should be empty")
	assert.Equal(t, "200 OK", resp.Status, "should get a 200")
}
