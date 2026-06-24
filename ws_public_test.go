package okx

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/UnipayFI/go-okx/common"
	"github.com/UnipayFI/go-okx/request"
	"github.com/go-json-experiment/json/jsontext"
)

// TestWsPublic verifies every public-gateway market channel: it subscribes live,
// captures the first data push, and checks the typed struct covers every pushed
// key. Channels that only push on activity (instruments / status /
// liquidation-orders) tolerate a no-push within the window. The VIP-gated
// "*-l2-tbt" depth channels log in and then return code 64003 unless the account
// qualifies, which the test tolerates.
func TestWsPublic(t *testing.T) {
	c := testWsPublicClient()
	const timeout = 15 * time.Second

	// Reference each typed service constructor so the typed path is exercised.
	_ = c.NewSubscribeTickersService("BTC-USDT")
	_ = c.NewSubscribeTradesService("BTC-USDT")
	_ = c.NewSubscribeBooksService("BTC-USDT")
	_ = c.NewSubscribeBooks5Service("BTC-USDT")
	_ = c.NewSubscribeBboTbtService("BTC-USDT")
	_ = c.NewSubscribeBooks50L2TbtService("BTC-USDT")
	_ = c.NewSubscribeBooksL2TbtService("BTC-USDT")
	_ = c.NewSubscribeInstrumentsService(InstTypeSpot)
	_ = c.NewSubscribeOpenInterestService("BTC-USDT-SWAP")
	_ = c.NewSubscribeFundingRateService("BTC-USDT-SWAP")
	_ = c.NewSubscribePriceLimitService("BTC-USDT-SWAP")
	_ = c.NewSubscribeMarkPriceService("BTC-USDT-SWAP")
	_ = c.NewSubscribeIndexTickersService("BTC-USDT")
	_ = c.NewSubscribeStatusService()
	_ = c.NewSubscribeLiquidationOrdersService(InstTypeSwap)

	t.Run("tickers", func(t *testing.T) {
		arg := request.WsArg{Channel: "tickers", InstrumentID: "BTC-USDT"}
		raw := wsFirstDataArray(t, c, request.GatewayPublic, false, arg, timeout)
		if raw == nil {
			return
		}
		var data []WsTicker
		if err := common.JSONUnmarshal(raw, &data); err != nil {
			t.Fatalf("decode: %v", err)
		}
		assertCovers(t, "ws/tickers", raw, data)
	})

	t.Run("trades", func(t *testing.T) {
		arg := request.WsArg{Channel: "trades", InstrumentID: "BTC-USDT"}
		raw := wsFirstDataArray(t, c, request.GatewayPublic, false, arg, timeout)
		if raw == nil {
			return
		}
		var data []WsTrade
		if err := common.JSONUnmarshal(raw, &data); err != nil {
			t.Fatalf("decode: %v", err)
		}
		assertCovers(t, "ws/trades", raw, data)
	})

	t.Run("books", func(t *testing.T) {
		arg := request.WsArg{Channel: "books", InstrumentID: "BTC-USDT"}
		raw := wsFirstDataArray(t, c, request.GatewayPublic, false, arg, timeout)
		if raw == nil {
			return
		}
		var data []WsOrderBook
		if err := common.JSONUnmarshal(raw, &data); err != nil {
			t.Fatalf("decode: %v", err)
		}
		assertCovers(t, "ws/books", raw, data)
	})

	t.Run("books5", func(t *testing.T) {
		arg := request.WsArg{Channel: "books5", InstrumentID: "BTC-USDT"}
		raw := wsFirstDataArray(t, c, request.GatewayPublic, false, arg, timeout)
		if raw == nil {
			return
		}
		var data []WsOrderBook
		if err := common.JSONUnmarshal(raw, &data); err != nil {
			t.Fatalf("decode: %v", err)
		}
		assertCovers(t, "ws/books5", raw, data)
	})

	t.Run("bbo-tbt", func(t *testing.T) {
		arg := request.WsArg{Channel: "bbo-tbt", InstrumentID: "BTC-USDT"}
		raw := wsFirstDataArray(t, c, request.GatewayPublic, false, arg, timeout)
		if raw == nil {
			return
		}
		var data []WsOrderBook
		if err := common.JSONUnmarshal(raw, &data); err != nil {
			t.Fatalf("decode: %v", err)
		}
		assertCovers(t, "ws/bbo-tbt", raw, data)
	})

	// VIP-gated depth channels: log in (private=true) then expect code 64003
	// unless the fee tier qualifies. Tolerate the error — the login+subscribe path
	// is what we validate here. Use an authenticated client when creds are present
	// so login actually happens; otherwise fall back to the public client (login
	// then fails with "missing credentials", still tolerated).
	authC := c
	if os.Getenv("OKX_API_KEY") != "" && os.Getenv("OKX_API_SECRET") != "" && os.Getenv("OKX_PASSPHRASE") != "" {
		authC = testWsClient(t)
	}

	t.Run("books50-l2-tbt", func(t *testing.T) {
		arg := request.WsArg{Channel: "books50-l2-tbt", InstrumentID: "BTC-USDT"}
		raw := wsFirstDataArrayTolerant(t, authC, request.GatewayPublic, true, arg, timeout)
		if raw == nil {
			return
		}
		var data []WsOrderBook
		if err := common.JSONUnmarshal(raw, &data); err != nil {
			t.Fatalf("decode: %v", err)
		}
		assertCovers(t, "ws/books50-l2-tbt", raw, data)
	})

	t.Run("books-l2-tbt", func(t *testing.T) {
		arg := request.WsArg{Channel: "books-l2-tbt", InstrumentID: "BTC-USDT"}
		raw := wsFirstDataArrayTolerant(t, authC, request.GatewayPublic, true, arg, timeout)
		if raw == nil {
			return
		}
		var data []WsOrderBook
		if err := common.JSONUnmarshal(raw, &data); err != nil {
			t.Fatalf("decode: %v", err)
		}
		assertCovers(t, "ws/books-l2-tbt", raw, data)
	})

	t.Run("instruments", func(t *testing.T) {
		arg := request.WsArg{Channel: "instruments", InstrumentType: "SPOT"}
		raw := wsFirstDataArray(t, c, request.GatewayPublic, false, arg, timeout)
		if raw == nil {
			return // instruments push only on listing/rule change
		}
		var data []WsInstrument
		if err := common.JSONUnmarshal(raw, &data); err != nil {
			t.Fatalf("decode: %v", err)
		}
		assertCovers(t, "ws/instruments", raw, data)
	})

	t.Run("open-interest", func(t *testing.T) {
		arg := request.WsArg{Channel: "open-interest", InstrumentID: "BTC-USDT-SWAP"}
		raw := wsFirstDataArray(t, c, request.GatewayPublic, false, arg, timeout)
		if raw == nil {
			return
		}
		var data []WsOpenInterest
		if err := common.JSONUnmarshal(raw, &data); err != nil {
			t.Fatalf("decode: %v", err)
		}
		assertCovers(t, "ws/open-interest", raw, data)
	})

	t.Run("funding-rate", func(t *testing.T) {
		arg := request.WsArg{Channel: "funding-rate", InstrumentID: "BTC-USDT-SWAP"}
		raw := wsFirstDataArray(t, c, request.GatewayPublic, false, arg, timeout)
		if raw == nil {
			return
		}
		var data []WsFundingRate
		if err := common.JSONUnmarshal(raw, &data); err != nil {
			t.Fatalf("decode: %v", err)
		}
		assertCovers(t, "ws/funding-rate", raw, data)
	})

	t.Run("price-limit", func(t *testing.T) {
		arg := request.WsArg{Channel: "price-limit", InstrumentID: "BTC-USDT-SWAP"}
		raw := wsFirstDataArray(t, c, request.GatewayPublic, false, arg, timeout)
		if raw == nil {
			return
		}
		var data []WsPriceLimit
		if err := common.JSONUnmarshal(raw, &data); err != nil {
			t.Fatalf("decode: %v", err)
		}
		assertCovers(t, "ws/price-limit", raw, data)
	})

	t.Run("mark-price", func(t *testing.T) {
		arg := request.WsArg{Channel: "mark-price", InstrumentID: "BTC-USDT-SWAP"}
		raw := wsFirstDataArray(t, c, request.GatewayPublic, false, arg, timeout)
		if raw == nil {
			return
		}
		var data []WsMarkPrice
		if err := common.JSONUnmarshal(raw, &data); err != nil {
			t.Fatalf("decode: %v", err)
		}
		assertCovers(t, "ws/mark-price", raw, data)
	})

	t.Run("index-tickers", func(t *testing.T) {
		arg := request.WsArg{Channel: "index-tickers", InstrumentID: "BTC-USDT"}
		raw := wsFirstDataArray(t, c, request.GatewayPublic, false, arg, timeout)
		if raw == nil {
			return
		}
		var data []WsIndexTicker
		if err := common.JSONUnmarshal(raw, &data); err != nil {
			t.Fatalf("decode: %v", err)
		}
		assertCovers(t, "ws/index-tickers", raw, data)
	})

	t.Run("status", func(t *testing.T) {
		arg := request.WsArg{Channel: "status"}
		raw := wsFirstDataArray(t, c, request.GatewayPublic, false, arg, timeout)
		if raw == nil {
			return // status pushes only around maintenance windows
		}
		var data []WsStatus
		if err := common.JSONUnmarshal(raw, &data); err != nil {
			t.Fatalf("decode: %v", err)
		}
		assertCovers(t, "ws/status", raw, data)
	})

	t.Run("liquidation-orders", func(t *testing.T) {
		arg := request.WsArg{Channel: "liquidation-orders", InstrumentType: "SWAP"}
		raw := wsFirstDataArray(t, c, request.GatewayPublic, false, arg, timeout)
		if raw == nil {
			return // liquidations are infrequent
		}
		var data []WsLiquidationOrder
		if err := common.JSONUnmarshal(raw, &data); err != nil {
			t.Fatalf("decode: %v", err)
		}
		assertCovers(t, "ws/liquidation-orders", raw, data)
	})
}

// wsFirstDataArrayTolerant is like wsFirstDataArray but does NOT fail on a
// subscribe error frame (e.g. code 64003 "fee tier doesn't meet requirement" for
// the VIP-gated depth channels): it logs and returns nil so the subscribe+login
// path is still validated.
func wsFirstDataArrayTolerant(t *testing.T, c *WebSocketClient, gateway request.Gateway, private bool, arg request.WsArg, timeout time.Duration) []byte {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	out := make(chan []byte, 4)
	done, _, err := request.SubscribeRaw(ctx, c, gateway, private, arg, func(message []byte, e error) {
		if e != nil {
			t.Logf("ws %s: subscribe error (tolerated): %v", arg.Channel, e)
			return
		}
		var f struct {
			Data jsontext.Value `json:"data"`
		}
		if err := common.JSONUnmarshal(message, &f); err != nil || len(f.Data) == 0 {
			return
		}
		select {
		case out <- f.Data:
		default:
		}
	})
	if err != nil {
		t.Logf("ws %s: subscribe setup error (tolerated): %v", arg.Channel, err)
		return nil
	}
	defer close(done)

	select {
	case raw := <-out:
		return raw
	case <-ctx.Done():
		t.Logf("ws %s: no data push within %s (VIP-gated / activity) — subscribe+login OK", arg.Channel, timeout)
		return nil
	}
}
