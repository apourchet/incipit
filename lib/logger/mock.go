package logger

import (
	"fmt"
	"strconv"
	"time"

	"github.com/apourchet/incipit/lib/etcd"
	protos "github.com/apourchet/incipit/protos/go"
	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type mockLogger struct {
	store etcd.EtcdClient
}

func NewMockLoggerClient() protos.LoggerClient {
	m := mockLogger{}
	m.store = etcd.GetDefaultClient()
	return &m
}

func (m *mockLogger) LogLogin(ctx context.Context, req *protos.LogLoginReq, opts ...grpc.CallOption) (*protos.LogLoginRes, error) {
	key := req.UserId
	glog.Infof("LogLogin: %s", key)

	// Get the last login value
	var lastLogin int64
	val, found, err := m.store.Get(key)
	if err != nil {
		return nil, err
	}
	if found {
		lastLogin, err = strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, err
		}
	}

	res := &protos.LogLoginRes{}
	res.LastLogin = int64(lastLogin)

	// Store the new
	newLogin := time.Now().Unix()
	return res, m.store.Set(key, fmt.Sprintf("%d", newLogin))
}
