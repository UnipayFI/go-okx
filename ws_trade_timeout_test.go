package okx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/UnipayFI/go-okx/client"
	"github.com/gorilla/websocket"
	"github.com/shopspring/decimal"
)

func tradeWsURL(s *httptest.Server) string {
	return "ws" + strings.TrimPrefix(s.URL, "http")
}

// TestWsTradeReadTimeoutTeardown verifies that a silent (half-open) order-entry
// connection is torn down by the steady-state read deadline, so an in-flight
// request fails promptly instead of blocking until its context deadline.
func TestWsTradeReadTimeoutTeardown(t *testing.T) {
	upgrader := websocket.Upgrader{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		// login handshake: read the login op, ack success.
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
		_ = conn.WriteMessage(websocket.TextMessage, []byte(`{"event":"login","code":"0","msg":""}`))
		// then go silent: keep reading (to consume the order op) but never reply,
		// simulating a half-open connection that never produces a frame.
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}))
	defer srv.Close()

	c := NewWebSocketClient(
		client.WithWebSocketPrivateURL(tradeWsURL(srv)),
		client.WithWebSocketAuth("k", "s", "p"),
		client.WithWebSocketReadTimeout(200*time.Millisecond),
		client.WithWebSocketWriteTimeout(time.Second),
	)
	tc, err := c.DialTrade(context.Background())
	if err != nil {
		t.Fatalf("DialTrade: %v", err)
	}
	defer tc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	start := time.Now()
	_, err = tc.PlaceOrder(ctx, OrderArg{
		InstrumentID: "BTC-USDT",
		TradeMode:    TdMode("cash"),
		Side:         SideBuy,
		OrderType:    OrdTypeLimit,
		Size:         decimal.RequireFromString("1"),
	})
	elapsed := time.Since(start)
	if err == nil {
		t.Fatal("PlaceOrder against a silent server: expected error, got nil")
	}
	if elapsed > 2*time.Second {
		t.Fatalf("PlaceOrder took %v — read deadline did not tear down the half-open connection", elapsed)
	}
	t.Logf("PlaceOrder failed in %v as expected: %v", elapsed, err)
}
