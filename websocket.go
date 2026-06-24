package okx

import (
	"context"

	"github.com/UnipayFI/go-okx/client"
	"github.com/UnipayFI/go-okx/request"
)

var _ request.WsClient = (*WebSocketClient)(nil)

// WebSocketClient streams OKX v5 public, private and business channels. Public
// channels need no credentials; private and (private) business channels require
// WithWebSocketAuth and log in automatically.
type WebSocketClient struct {
	*client.WebSocketClient
}

// NewWebSocketClient constructs an OKX v5 WebSocket client.
func NewWebSocketClient(options ...client.WebSocketOptions) *WebSocketClient {
	return &WebSocketClient{client.NewWebSocketClient(options...)}
}

// WsHandler is invoked for every push (or error) on a subscription. The push's
// Data field is already decoded into the channel's typed slice.
type WsHandler[T any] func(*request.WsPush[[]T], error)

// Subscribe is the low-level escape hatch: it subscribes to an arbitrary channel
// on the given gateway (request.GatewayPublic / GatewayPrivate / GatewayBusiness)
// and delivers each data push's raw bytes. Prefer the typed NewSubscribe*
// services; use this for channels the SDK does not yet wrap.
func (c *WebSocketClient) Subscribe(ctx context.Context, gateway request.Gateway, private bool, arg request.WsArg, cb func([]byte, error)) (chan<- struct{}, <-chan struct{}, error) {
	return request.SubscribeRaw(ctx, c, gateway, private, arg, cb)
}
