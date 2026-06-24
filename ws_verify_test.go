package okx

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/UnipayFI/go-okx/client"
	"github.com/UnipayFI/go-okx/common"
	"github.com/UnipayFI/go-okx/request"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/gorilla/websocket"
)

// testWsPublicClient builds a WebSocket client for public channels (no auth),
// honoring OKX_PROXY.
func testWsPublicClient() *WebSocketClient {
	opts := []client.WebSocketOptions{}
	if proxy := os.Getenv("OKX_PROXY"); proxy != "" {
		opts = append(opts, client.WithWebSocketProxy(proxy))
	}
	return NewWebSocketClient(opts...)
}

// testWsClient builds an authenticated WebSocket client (for private/business
// login channels), skipping the test when credentials are absent.
func testWsClient(t *testing.T) *WebSocketClient {
	t.Helper()
	apiKey := os.Getenv("OKX_API_KEY")
	apiSecret := os.Getenv("OKX_API_SECRET")
	passphrase := os.Getenv("OKX_PASSPHRASE")
	if apiKey == "" || apiSecret == "" || passphrase == "" {
		t.Skip("OKX_API_KEY/SECRET/PASSPHRASE not set; skipping private ws test")
	}
	opts := []client.WebSocketOptions{client.WithWebSocketAuth(apiKey, apiSecret, passphrase)}
	if proxy := os.Getenv("OKX_PROXY"); proxy != "" {
		opts = append(opts, client.WithWebSocketProxy(proxy))
	}
	return NewWebSocketClient(opts...)
}

// wsFirstDataArray subscribes (raw) to arg on the given gateway and returns the
// raw JSON bytes of the FIRST data frame's "data" array — the WebSocket analogue
// of fetchRawGet, so the same assertCovers can diff a push against the typed
// struct. Fails the test on subscribe/login error; returns nil (with a logged
// skip) when no data push arrives before the timeout (some private channels only
// push on activity).
func wsFirstDataArray(t *testing.T, c *WebSocketClient, gateway request.Gateway, private bool, arg request.WsArg, timeout time.Duration) []byte {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	type frame struct {
		Data jsontext.Value `json:"data"`
	}
	out := make(chan []byte, 4)
	errCh := make(chan error, 4)
	done, _, err := request.SubscribeRaw(ctx, c, gateway, private, arg, func(message []byte, e error) {
		if e != nil {
			select {
			case errCh <- e:
			default:
			}
			return
		}
		var f frame
		if err := common.JSONUnmarshal(message, &f); err != nil || len(f.Data) == 0 {
			return
		}
		select {
		case out <- f.Data:
		default:
		}
	})
	if err != nil {
		t.Fatalf("ws subscribe %s: %v", arg.Channel, err)
		return nil
	}
	defer close(done)

	select {
	case raw := <-out:
		return raw
	case e := <-errCh:
		t.Fatalf("ws %s push error: %v", arg.Channel, e)
		return nil
	case <-ctx.Done():
		t.Logf("ws %s: no data push within %s (activity-driven channel?) — subscribe+login OK", arg.Channel, timeout)
		return nil
	}
}

// wsCtx returns a context with the given timeout, cleaned up via t.Cleanup.
func wsCtx(t *testing.T, d time.Duration) context.Context {
	t.Helper()
	c, cancel := context.WithTimeout(context.Background(), d)
	t.Cleanup(cancel)
	return c
}

// ensure gorilla/websocket stays a direct dependency even if only the request
// package references it transitively.
var _ = websocket.TextMessage
