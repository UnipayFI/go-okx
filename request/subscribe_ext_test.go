package request

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/UnipayFI/go-okx/client"
	"github.com/UnipayFI/go-okx/common"
	"github.com/gorilla/websocket"
)

func wsURL(s *httptest.Server) string {
	return "ws" + strings.TrimPrefix(s.URL, "http")
}

// TestSubscribeManyRaw verifies that the multi-arg primitive subscribes every
// channel in a single op on one connection and delivers each data frame raw.
func TestSubscribeManyRaw(t *testing.T) {
	upgrader := websocket.Upgrader{}
	argCount := make(chan int, 1)
	srvDone := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		var op struct {
			Op   string  `json:"op"`
			Args []WsArg `json:"args"`
		}
		_ = common.JSONUnmarshal(msg, &op)
		argCount <- len(op.Args)
		_ = conn.WriteMessage(websocket.TextMessage,
			[]byte(`{"arg":{"channel":"books","instId":"BTC-USDT"},"data":[{"x":"1"}]}`))
		<-srvDone
	}))
	defer srv.Close()
	defer close(srvDone)

	c := client.NewWebSocketClient(client.WithWebSocketPublicURL(wsURL(srv)))
	got := make(chan []byte, 1)
	done, _, err := SubscribeManyRaw(context.Background(), c, GatewayPublic, false,
		[]WsArg{{Channel: "books", InstrumentID: "BTC-USDT"}, {Channel: "trades", InstrumentID: "BTC-USDT"}},
		func(msg []byte, err error) {
			if err == nil && msg != nil {
				select {
				case got <- msg:
				default:
				}
			}
		})
	if err != nil {
		t.Fatalf("SubscribeManyRaw: %v", err)
	}
	defer close(done)

	select {
	case n := <-argCount:
		if n != 2 {
			t.Errorf("subscribe op carried %d args, want 2 (single-connection multi-channel)", n)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("server never received the subscribe op")
	}
	select {
	case msg := <-got:
		if !strings.Contains(string(msg), `"channel":"books"`) {
			t.Errorf("unexpected data frame: %s", msg)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("cb never received the data frame")
	}
}

// TestSubscribeReadTimeout verifies that a steady-state read deadline fires on a
// silent (half-open) connection and surfaces a timeout error to the callback.
func TestSubscribeReadTimeout(t *testing.T) {
	upgrader := websocket.Upgrader{}
	srvDone := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		// Consume the subscribe op, then stay silent (keep the conn open) so the
		// client's read deadline — not a server close — is what ends the read.
		_, _, _ = conn.ReadMessage()
		<-srvDone
	}))
	defer srv.Close()
	defer close(srvDone)

	c := client.NewWebSocketClient(
		client.WithWebSocketPublicURL(wsURL(srv)),
		client.WithWebSocketReadTimeout(150*time.Millisecond),
	)
	errC := make(chan error, 1)
	done, _, err := SubscribeRaw(context.Background(), c, GatewayPublic, false,
		WsArg{Channel: "books", InstrumentID: "BTC-USDT"},
		func(msg []byte, err error) {
			if err != nil {
				select {
				case errC <- err:
				default:
				}
			}
		})
	if err != nil {
		t.Fatalf("SubscribeRaw: %v", err)
	}
	defer close(done)

	select {
	case got := <-errC:
		if !os.IsTimeout(got) {
			t.Errorf("read error = %v, want a timeout (read deadline)", got)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("read deadline never fired")
	}
}
