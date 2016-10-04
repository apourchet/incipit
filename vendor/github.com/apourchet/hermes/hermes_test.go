package hermes

import (
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

// Service definition
type MyService struct{}

func (s *MyService) Host() string {
	return "localhost:9000"
}

func (s *MyService) Endpoints() EndpointMap {
	return EndpointMap{
		Endpoint{"RpcCall", "GET", "/test", NewInbound, NewOutbound},
		Endpoint{"RpcCall", "POST", "/test", NewInbound, NewOutbound},
	}
}

// Endpoint definitions
type Inbound struct {
	Message string
}

type Outbound struct {
	Ok bool
}

func NewInbound() interface{}  { return &Inbound{} }
func NewOutbound() interface{} { return &Outbound{} }

func (s *MyService) RpcCall(c *gin.Context, in *Inbound, out *Outbound) (int, error) {
	if in.Message == "secret" {
		out.Ok = true
		return http.StatusOK, nil
	}
	out.Ok = false
	return http.StatusBadRequest, nil
}

// Mocked Service
type MockedService struct{}

func (s *MockedService) Endpoints() EndpointMap {
	return EndpointMap{
		Endpoint{"RpcCall", "GET", "/test", NewInbound, NewOutbound},
		Endpoint{"RpcCall", "POST", "/test", NewInbound, NewOutbound},
	}
}

func (s *MockedService) RpcCall(in *Inbound, out *Outbound) (int, error) {
	if in.Message == "secret" {
		out.Ok = true
		return http.StatusOK, nil
	}
	out.Ok = false
	return http.StatusBadRequest, nil
}

// Tests
func TestMain(m *testing.M) {
	si := InitService(&MyService{})
	engine := gin.New()
	si.Serve(engine)
	go engine.Run(":9000")
	time.Sleep(500 * time.Millisecond)
	m.Run()
}

func TestMock(t *testing.T) {
	si := InitMockService(&MockedService{})
	out := &Outbound{false}
	code, err := si.Call("RpcCall", &Inbound{"secret"}, out)
	if code != 200 {
		t.Fail()
	}
	if err != nil {
		t.Fail()
	}
	if !out.Ok {
		t.Fail()
	}
}

func TestCallSuccess(t *testing.T) {
	si := InitService(&MyService{})
	out := &Outbound{false}
	code, err := si.Call("RpcCall", &Inbound{"secret"}, out)
	if code != 200 {
		t.Fail()
	}
	if err != nil {
		t.Fail()
	}
	if !out.Ok {
		t.Fail()
	}
}

func TestCallWrongSecret(t *testing.T) {
	si := InitService(&MyService{})
	out := &Outbound{true}
	code, err := si.Call("RpcCall", &Inbound{"wrong secret"}, out)
	if code != 400 {
		t.Fail()
	}
	if err != nil {
		t.Fail()
	}
	if out.Ok {
		t.Fail()
	}
}

func TestCallNotFound(t *testing.T) {
	si := InitService(&MyService{})
	code, err := si.Call("NotAnEndpoint", &Inbound{}, &Outbound{})
	if code != 404 || err == nil {
		t.Fail()
	}
}
