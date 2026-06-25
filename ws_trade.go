package okx

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/UnipayFI/go-okx/common"
	"github.com/UnipayFI/go-okx/request"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/gorilla/websocket"
)

// WsTradeConn is a persistent, logged-in WebSocket connection for order entry —
// a low-latency alternative to the REST trade endpoints. OKX correlates each
// request with its response by a client-supplied "id"; this type assigns ids,
// matches responses, and exposes one method per order operation.
//
// Obtain one via (*WebSocketClient).DialTrade and Close it when done. All
// methods are safe for concurrent use.
type WsTradeConn struct {
	conn   *websocket.Conn
	logger interface{ Debugf(string, ...any) }

	readTimeout  time.Duration
	writeTimeout time.Duration

	mu      sync.Mutex
	pending map[string]chan wsTradeResp
	seq     int64
	writeMu sync.Mutex
	closed  chan struct{}
	closeMu sync.Once
}

// wsTradeResp is the op-response envelope OKX returns over the trade connection.
type wsTradeResp struct {
	ID        string         `json:"id"`
	Operation string         `json:"op"`
	Code      string         `json:"code"`
	Message   string         `json:"msg"`
	Data      jsontext.Value `json:"data"`
	InTime    string         `json:"inTime"`
	OutTime   string         `json:"outTime"`
}

type wsTradeReq struct {
	ID        string `json:"id"`
	Operation string `json:"op"`
	Args      []any  `json:"args"`
}

// DialTrade dials the private gateway, logs in, and returns a ready order-entry
// connection. The caller must Close it.
func (c *WebSocketClient) DialTrade(ctx context.Context) (*WsTradeConn, error) {
	conn, err := request.DialLoggedIn(ctx, c, request.GatewayPrivate)
	if err != nil {
		return nil, err
	}
	tc := &WsTradeConn{
		conn:         conn,
		logger:       c.GetLogger(),
		readTimeout:  c.GetReadTimeout(),
		writeTimeout: c.GetWriteTimeout(),
		pending:      make(map[string]chan wsTradeResp),
		closed:       make(chan struct{}),
	}
	go tc.readLoop()
	go tc.keepAlive()
	return tc, nil
}

// Close terminates the connection and fails any in-flight requests.
func (tc *WsTradeConn) Close() error {
	tc.closeMu.Do(func() { close(tc.closed) })
	return tc.conn.Close()
}

func (tc *WsTradeConn) nextID() string {
	return strconv.FormatInt(atomic.AddInt64(&tc.seq, 1), 10)
}

func (tc *WsTradeConn) readLoop() {
	for {
		// Bound the read so a half-open connection surfaces an error instead of
		// blocking forever (which would hang every in-flight request until its
		// ctx). A healthy idle connection still sees a keepalive "pong" within
		// this window, which resets the deadline each iteration.
		if tc.readTimeout > 0 {
			_ = tc.conn.SetReadDeadline(time.Now().Add(tc.readTimeout))
		}
		_, message, err := tc.conn.ReadMessage()
		if err != nil {
			tc.failAll(err)
			_ = tc.Close() // mark dead: close `closed` (stops keepAlive) + conn
			return
		}
		if common.BytesToString(message) == "pong" {
			continue
		}
		tc.logger.Debugf("ws trade recv: %s", common.BytesToString(message))
		var resp wsTradeResp
		if err := common.JSONUnmarshal(message, &resp); err != nil || resp.ID == "" {
			continue // control frame (e.g. login ack) or unparseable
		}
		tc.mu.Lock()
		ch := tc.pending[resp.ID]
		delete(tc.pending, resp.ID)
		tc.mu.Unlock()
		if ch != nil {
			ch <- resp
		}
	}
}

func (tc *WsTradeConn) failAll(err error) {
	tc.mu.Lock()
	for id, ch := range tc.pending {
		ch <- wsTradeResp{Code: "-1", Message: "connection closed: " + err.Error()}
		delete(tc.pending, id)
	}
	tc.mu.Unlock()
}

func (tc *WsTradeConn) keepAlive() {
	ticker := time.NewTicker(common.DEFAULT_KEEP_ALIVE_INTERVAL)
	defer ticker.Stop()
	for {
		select {
		case <-tc.closed:
			return
		case <-ticker.C:
			tc.writeMu.Lock()
			if tc.writeTimeout > 0 {
				_ = tc.conn.SetWriteDeadline(time.Now().Add(tc.writeTimeout))
			}
			err := tc.conn.WriteMessage(websocket.TextMessage, []byte("ping"))
			tc.writeMu.Unlock()
			if err != nil {
				// A failed keepalive write means the socket is gone (or half-open):
				// tear the connection down so the reader unblocks and pending
				// requests fail, rather than leaving a dead connection live.
				_ = tc.Close()
				return
			}
		}
	}
}

// request sends an op request and waits for the matching response (or ctx
// cancellation). The caller decodes resp.Data.
func (tc *WsTradeConn) request(ctx context.Context, op string, args []any) (*wsTradeResp, error) {
	id := tc.nextID()
	ch := make(chan wsTradeResp, 1)
	tc.mu.Lock()
	tc.pending[id] = ch
	tc.mu.Unlock()

	data, err := common.JSONMarshal(wsTradeReq{ID: id, Operation: op, Args: args})
	if err != nil {
		tc.mu.Lock()
		delete(tc.pending, id)
		tc.mu.Unlock()
		return nil, err
	}
	tc.writeMu.Lock()
	if tc.writeTimeout > 0 {
		_ = tc.conn.SetWriteDeadline(time.Now().Add(tc.writeTimeout))
	}
	err = tc.conn.WriteMessage(websocket.TextMessage, data)
	tc.writeMu.Unlock()
	if err != nil {
		tc.mu.Lock()
		delete(tc.pending, id)
		tc.mu.Unlock()
		return nil, err
	}

	select {
	case <-ctx.Done():
		tc.mu.Lock()
		delete(tc.pending, id)
		tc.mu.Unlock()
		return nil, ctx.Err()
	case <-tc.closed:
		return nil, errors.New("ws trade: connection closed")
	case resp := <-ch:
		return &resp, nil
	}
}

// decodeResults decodes the op response's data array and returns it. A
// connection/system-level failure (code outside 0/1/2, or empty data on a
// non-zero code) is returned as a *client.WsError-style error via WsError.
func decodeResults[T any](resp *wsTradeResp) ([]T, error) {
	var results []T
	if len(resp.Data) > 0 {
		if err := common.JSONUnmarshal(resp.Data, &results); err != nil {
			return nil, err
		}
	}
	switch resp.Code {
	case "0", "1", "2":
		return results, nil
	default:
		return results, &request.WsError{Code: resp.Code, Message: resp.Message}
	}
}

// PlaceOrder places a single order over the trade connection. The per-order
// result carries sCode/sMsg for the order-level outcome.
func (tc *WsTradeConn) PlaceOrder(ctx context.Context, arg OrderArg) (*OrderResult, error) {
	resp, err := tc.request(ctx, "order", []any{arg})
	if err != nil {
		return nil, err
	}
	results, err := decodeResults[OrderResult](resp)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, &request.WsError{Code: resp.Code, Message: resp.Message}
	}
	return &results[0], nil
}

// BatchPlaceOrders places up to 20 orders in one request.
func (tc *WsTradeConn) BatchPlaceOrders(ctx context.Context, args []OrderArg) ([]OrderResult, error) {
	resp, err := tc.request(ctx, "batch-orders", toAnySlice(args))
	if err != nil {
		return nil, err
	}
	return decodeResults[OrderResult](resp)
}

// CancelOrder cancels a single order.
func (tc *WsTradeConn) CancelOrder(ctx context.Context, arg CancelOrderArg) (*OrderResult, error) {
	resp, err := tc.request(ctx, "cancel-order", []any{arg})
	if err != nil {
		return nil, err
	}
	results, err := decodeResults[OrderResult](resp)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, &request.WsError{Code: resp.Code, Message: resp.Message}
	}
	return &results[0], nil
}

// BatchCancelOrders cancels up to 20 orders in one request.
func (tc *WsTradeConn) BatchCancelOrders(ctx context.Context, args []CancelOrderArg) ([]OrderResult, error) {
	resp, err := tc.request(ctx, "batch-cancel-orders", toAnySlice(args))
	if err != nil {
		return nil, err
	}
	return decodeResults[OrderResult](resp)
}

// AmendOrder amends a single order.
func (tc *WsTradeConn) AmendOrder(ctx context.Context, arg AmendOrderArg) (*AmendResult, error) {
	resp, err := tc.request(ctx, "amend-order", []any{arg})
	if err != nil {
		return nil, err
	}
	results, err := decodeResults[AmendResult](resp)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, &request.WsError{Code: resp.Code, Message: resp.Message}
	}
	return &results[0], nil
}

// BatchAmendOrders amends up to 20 orders in one request.
func (tc *WsTradeConn) BatchAmendOrders(ctx context.Context, args []AmendOrderArg) ([]AmendResult, error) {
	resp, err := tc.request(ctx, "batch-amend-orders", toAnySlice(args))
	if err != nil {
		return nil, err
	}
	return decodeResults[AmendResult](resp)
}

// MassCancel cancels all pending orders for an option instrument family.
func (tc *WsTradeConn) MassCancel(ctx context.Context, instType InstType, instFamily string) (*MassCancelResult, error) {
	arg := map[string]any{"instType": string(instType), "instFamily": instFamily}
	resp, err := tc.request(ctx, "mass-cancel", []any{arg})
	if err != nil {
		return nil, err
	}
	results, err := decodeResults[MassCancelResult](resp)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, &request.WsError{Code: resp.Code, Message: resp.Message}
	}
	return &results[0], nil
}

func toAnySlice[T any](items []T) []any {
	out := make([]any, len(items))
	for i, it := range items {
		out[i] = it
	}
	return out
}
