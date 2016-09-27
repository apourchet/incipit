# Hermes
Hermes is a simple pseudo-RPC framework in golang. It's built on top of [gin](https://github.com/gin-gonic/gin). The main advantage is that it's 1) curlable and 2) has no code generation.

# Example
### Service Definition
```go
// Service definition
type MyService struct{}

func (s *MyService) Host() string {
	return "localhost:9000"
}

func (s *MyService) Endpoints() EndpointMap {
	return EndpointMap{
        // RpcCall will default to GET since it is first
		Endpoint{"RpcCall", "GET", "/test", NewInbound, NewOutbound}, 
		Endpoint{"RpcCall", "POST", "/test", NewInbound, NewOutbound},
		Endpoint{"OtherRpcCall", "POST", "/test", NewOtherInbound, NewOtherOutbound},
	}
}

// Endpoint definitions
// RpcCall
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

// OtherRpcCall
type OtherInbound struct {
    MyField int
}

type OtherOutbound struct {
    SomeFloat float64
}

func NewOtherInbound() interface{}  { return &OtherInbound{} }
func NewOtherOutbound() interface{} { return &OtherOutbound{} }

func (s *MyService) OtherRpcCall(c *gin.Context, in *OtherInbound, out *OtherOutbound) (int, error) {
    out.SomeFloat = 3.14 * in.MyField
    return http.StatusOK, nil
}
```

### Server Creation
```go
si := InitService(&MyService{})
engine := gin.New()
si.Serve(engine)
go engine.Run(":9000")
```

### Client RPC Call
```go
si := InitService(&MyService{})
out := &Outbound{false}
code, err := si.Call("RpcCall", &Inbound{"secret"}, out)
```
