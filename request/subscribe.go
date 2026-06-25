package request

import (
	"context"
	"errors"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/UnipayFI/go-okx/common"
	"github.com/UnipayFI/go-okx/pkg/log"
	"github.com/gorilla/websocket"
)

// WsClient is what the subscribe framework needs from a *client.WebSocketClient.
type WsClient interface {
	GetPublicURL() string
	GetPrivateURL() string
	GetBusinessURL() string
	GetAPIKey() string
	GetAPISecret() string
	GetPassphrase() string
	GetSignFn() SignFn
	GetLogger() log.Logger
	GetDialer() *websocket.Dialer
	// GetReadTimeout bounds the idle time a steady-state read may block with no
	// inbound frame before the connection is treated as dead; zero disables it.
	GetReadTimeout() time.Duration
	// GetWriteTimeout bounds a single WS write (subscribe op / keepalive ping);
	// zero disables it.
	GetWriteTimeout() time.Duration
}

// Gateway selects which OKX v5 WebSocket endpoint a subscription uses. OKX splits
// channels across three gateways; the channel determines which one (and whether
// login is required).
type Gateway int

const (
	GatewayPublic Gateway = iota
	GatewayPrivate
	GatewayBusiness
)

func gatewayURL(c WsClient, g Gateway) string {
	switch g {
	case GatewayPrivate:
		return c.GetPrivateURL()
	case GatewayBusiness:
		return c.GetBusinessURL()
	default:
		return c.GetPublicURL()
	}
}

// WsArg identifies a channel subscription. channel is the channel name; the
// remaining fields narrow it and vary per channel (instId for a symbol channel,
// instType/instFamily/uly for product channels, ccy/algoId for others). Unset
// fields are omitted. The same shape is echoed back in each push's "arg".
type WsArg struct {
	Channel          string `json:"channel"`
	InstrumentType   string `json:"instType,omitempty"`
	InstrumentFamily string `json:"instFamily,omitempty"`
	InstrumentID     string `json:"instId,omitempty"`
	Underlying       string `json:"uly,omitempty"`
	Currency         string `json:"ccy,omitempty"`
	AlgoID           string `json:"algoId,omitempty"`
	SpreadID         string `json:"sprdId,omitempty"`
	ExtraParams      string `json:"extraParams,omitempty"`
}

// WsPush is the envelope OKX pushes for a data event. Action is set only on the
// order-book channels ("snapshot" / "update"). The typed Data field is the
// channel's payload (almost always an array).
type WsPush[T any] struct {
	Arg    WsArg  `json:"arg"`
	Action string `json:"action,omitempty"`
	Data   T      `json:"data"`
}

type wsOp struct {
	Operation string `json:"op"`
	Args      []any  `json:"args"`
}

type wsLoginOp struct {
	Operation string       `json:"op"`
	Args      []wsLoginArg `json:"args"`
}

type wsLoginArg struct {
	APIKey     string `json:"apiKey"`
	Passphrase string `json:"passphrase"`
	Timestamp  string `json:"timestamp"`
	Sign       string `json:"sign"`
}

// wsHeader is a lightweight view used to classify an inbound frame before
// committing to a typed decode. OKX encodes WS event codes as quoted strings.
type wsHeader struct {
	Event   string `json:"event"`
	Code    string `json:"code"`
	Message string `json:"msg"`
	// Arg is present on both data pushes and subscribe/error acks; Data only on
	// data pushes.
	Arg  *WsArg `json:"arg"`
	Data any    `json:"data"`
}

// Subscribe opens a dedicated connection to the chosen gateway, logs in when
// private, subscribes to arg, and invokes cb for every data push. Returns a done
// channel (close to stop) and a stop channel (closed when the reader exits). The
// typed Data field of the push is decoded into T.
func Subscribe[T any](ctx context.Context, client WsClient, gateway Gateway, private bool, arg WsArg, cb func(*WsPush[T], error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	return subscribeBytes(ctx, client, gateway, private, arg, func(message []byte, e error) {
		if e != nil {
			cb(nil, e)
			return
		}
		var push WsPush[T]
		if err := common.JSONUnmarshal(message, &push); err != nil {
			cb(nil, err)
			return
		}
		cb(&push, nil)
	})
}

// SubscribeRaw is like Subscribe but delivers each data frame's raw bytes, for
// channels whose payload shape the caller wants to decode itself.
func SubscribeRaw(ctx context.Context, client WsClient, gateway Gateway, private bool, arg WsArg, cb func(message []byte, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	return subscribeBytes(ctx, client, gateway, private, arg, cb)
}

// SubscribeManyRaw opens one connection to the gateway, logs in when private,
// subscribes to every arg in a single "subscribe" op (OKX accepts an args array
// on one connection), and delivers each data frame's raw bytes. The caller
// routes each frame to its channel by inspecting the push's "arg" (decode into
// WsPush[T] or a wsHeader-shaped struct and read Arg.Channel).
//
// This is the multi-channel counterpart to SubscribeRaw: a single private login
// can serve several channels (e.g. orders + positions + account) on one
// connection instead of one connection — and one login — per channel.
//
// Unlike Subscribe/SubscribeRaw it returns BIDIRECTIONAL channels, shaped for a
// reconnect supervisor that both tears the stream down and waits on its death:
// close done to tear the stream down; stop is closed when the reader exits (on a
// connection error or after a requested teardown completes).
func SubscribeManyRaw(ctx context.Context, client WsClient, gateway Gateway, private bool, args []WsArg, cb func(message []byte, err error)) (done, stop chan struct{}, err error) {
	anyArgs := make([]any, len(args))
	for i := range args {
		anyArgs[i] = args[i]
	}
	return subscribeManyBytes(ctx, client, gateway, private, anyArgs, cb)
}

func subscribeBytes(ctx context.Context, client WsClient, gateway Gateway, private bool, arg any, cb func(message []byte, err error)) (done, stop chan struct{}, err error) {
	return subscribeManyBytes(ctx, client, gateway, private, []any{arg}, cb)
}

func subscribeManyBytes(ctx context.Context, client WsClient, gateway Gateway, private bool, args []any, cb func(message []byte, err error)) (done, stop chan struct{}, err error) {
	conn, _, err := client.GetDialer().DialContext(ctx, gatewayURL(client, gateway), nil)
	if err != nil {
		return nil, nil, err
	}
	conn.SetReadLimit(10 << 20)
	readTimeout := client.GetReadTimeout()
	writeTimeout := client.GetWriteTimeout()

	if private {
		if err := wsLogin(client, conn); err != nil {
			conn.Close()
			return nil, nil, err
		}
	}

	sub := wsOp{Operation: "subscribe", Args: args}
	data, _ := common.JSONMarshal(sub)
	if writeTimeout > 0 {
		_ = conn.SetWriteDeadline(time.Now().Add(writeTimeout))
	}
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		conn.Close()
		return nil, nil, err
	}

	doneC := make(chan struct{})
	stopC := make(chan struct{})
	var silent atomic.Bool

	go keepAlive(conn, common.DEFAULT_KEEP_ALIVE_INTERVAL, writeTimeout)
	go func() {
		select {
		case <-stopC:
		case <-doneC:
		}
		// Either path is an intentional teardown: silence the reader so the
		// close-induced read error is not delivered to cb as a real error.
		silent.Store(true)
		unsub := wsOp{Operation: "unsubscribe", Args: args}
		if b, e := common.JSONMarshal(unsub); e == nil {
			if writeTimeout > 0 {
				_ = conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			}
			_ = conn.WriteMessage(websocket.TextMessage, b)
		}
		conn.Close()
	}()
	go func() {
		for {
			// Bound the read so a half-open (silently dropped) connection surfaces
			// an error and the caller can reconnect, instead of blocking forever.
			// A healthy stream always sees a frame (data or keepalive pong) within
			// this window, which resets each iteration.
			if readTimeout > 0 {
				_ = conn.SetReadDeadline(time.Now().Add(readTimeout))
			}
			_, message, err := conn.ReadMessage()
			if err != nil {
				if !silent.Load() {
					cb(nil, err)
				}
				close(stopC)
				return
			}
			if common.BytesToString(message) == "pong" {
				continue
			}
			client.GetLogger().Debugf("ws recv: %s", common.BytesToString(message))

			var hdr wsHeader
			if err := common.JSONUnmarshal(message, &hdr); err != nil {
				cb(nil, err)
				continue
			}
			switch {
			case hdr.Event == "error":
				cb(nil, &WsError{Code: hdr.Code, Message: hdr.Message})
			case hdr.Event != "":
				// subscribe / unsubscribe / login / channel-conn-count acks.
			case hdr.Data != nil:
				cb(message, nil)
			default:
				// other control frames.
			}
		}
	}()
	return doneC, stopC, nil
}

// DialLoggedIn dials the given gateway and completes the login handshake,
// returning a ready connection. WebSocket order entry (op:"order", ...) builds
// on this. The caller owns and must Close the returned connection.
func DialLoggedIn(ctx context.Context, client WsClient, gateway Gateway) (*websocket.Conn, error) {
	conn, _, err := client.GetDialer().DialContext(ctx, gatewayURL(client, gateway), nil)
	if err != nil {
		return nil, err
	}
	conn.SetReadLimit(10 << 20)
	if err := wsLogin(client, conn); err != nil {
		conn.Close()
		return nil, err
	}
	return conn, nil
}

// wsLogin performs the login handshake and blocks until the server acknowledges
// (or rejects) it. OKX signs ts + "GET" + "/users/self/verify" with the ts as
// Unix epoch seconds.
func wsLogin(client WsClient, conn *websocket.Conn) error {
	apiKey := client.GetAPIKey()
	secret := client.GetAPISecret()
	passphrase := client.GetPassphrase()
	if apiKey == "" || secret == "" || passphrase == "" {
		return errors.New("ws login: missing credentials (WithWebSocketAuth)")
	}
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	prehash := ts + "GET" + "/users/self/verify"
	var (
		sign string
		err  error
	)
	if fn := client.GetSignFn(); fn != nil {
		sign, err = fn(secret, prehash)
	} else {
		sign, err = HMACSign(secret, prehash)
	}
	if err != nil {
		return err
	}

	login := wsLoginOp{Operation: "login", Args: []wsLoginArg{{
		APIKey:     apiKey,
		Passphrase: passphrase,
		Timestamp:  ts,
		Sign:       sign,
	}}}
	data, _ := common.JSONMarshal(login)
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return err
	}

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	defer conn.SetReadDeadline(time.Time{})
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return err
		}
		if common.BytesToString(message) == "pong" {
			continue
		}
		var hdr wsHeader
		if err := common.JSONUnmarshal(message, &hdr); err != nil {
			return err
		}
		switch hdr.Event {
		case "login":
			if hdr.Code != "" && hdr.Code != "0" {
				return &WsError{Code: hdr.Code, Message: hdr.Message}
			}
			return nil
		case "error":
			return &WsError{Code: hdr.Code, Message: hdr.Message}
		}
	}
}

// keepAlive sends OKX's literal "ping" text frame on an interval; the server
// replies "pong" (handled in the read loop). A failed write means the socket is
// gone (or half-open), so it closes the connection to unblock the reader, which
// then surfaces the error and lets the caller reconnect — rather than leaving
// the reader hanging until its read deadline.
func keepAlive(conn *websocket.Conn, interval, writeTimeout time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		if writeTimeout > 0 {
			_ = conn.SetWriteDeadline(time.Now().Add(writeTimeout))
		}
		if err := conn.WriteMessage(websocket.TextMessage, []byte("ping")); err != nil {
			conn.Close()
			return
		}
	}
}

// WsError is an OKX WebSocket control-frame error.
type WsError struct {
	Code    string
	Message string
}

func (e *WsError) Error() string {
	return "<WsError> code=" + e.Code + ", msg=" + e.Message
}
