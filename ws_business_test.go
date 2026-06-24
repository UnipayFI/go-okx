package okx

import (
	"errors"
	"testing"
	"time"

	"github.com/UnipayFI/go-okx/common"
	"github.com/UnipayFI/go-okx/request"
	"github.com/go-json-experiment/json/jsontext"
)

// TestWsBusiness exercises every BUSINESS-gateway WebSocket channel: it
// subscribes (logging in for the private algo/grid/economic-calendar channels),
// captures the first live data push, and verifies the typed struct covers the
// real response shape. Candle channels (array-of-arrays) are validated by parsing
// instead of assertCovers. Activity-driven private channels that push nothing in
// time are tolerated (subscribe+login were still validated).
func TestWsBusiness(t *testing.T) {
	const sprdId = "BTC-USDT_BTC-USDT-SWAP"

	// --- public candle channels (array-of-arrays) ---

	t.Run("candle1m", func(t *testing.T) {
		c := testWsPublicClient()
		_ = c.NewSubscribeCandleService("BTC-USDT", "1m")
		arg := request.WsArg{Channel: "candle1m", InstrumentID: "BTC-USDT"}
		raw := wsFirstDataArray(t, c, request.GatewayBusiness, false, arg, 15*time.Second)
		if raw == nil {
			return
		}
		var rows [][]string
		if err := common.JSONUnmarshal(raw, &rows); err != nil {
			t.Fatalf("decode: %v", err)
		}
		candles := parseWsCandles(rows)
		if len(candles) == 0 || candles[0].Timestamp.IsZero() {
			t.Fatalf("candle1m: expected a parsed row with non-zero Ts, got %+v", candles)
		}
	})

	t.Run("mark-price-candle1m", func(t *testing.T) {
		c := testWsPublicClient()
		_ = c.NewSubscribeMarkPriceCandleService("BTC-USDT-SWAP", "1m")
		arg := request.WsArg{Channel: "mark-price-candle1m", InstrumentID: "BTC-USDT-SWAP"}
		raw := wsFirstDataArray(t, c, request.GatewayBusiness, false, arg, 15*time.Second)
		if raw == nil {
			return
		}
		var rows [][]string
		if err := common.JSONUnmarshal(raw, &rows); err != nil {
			t.Fatalf("decode: %v", err)
		}
		candles := parseWsIndexCandles(rows)
		if len(candles) == 0 || candles[0].Timestamp.IsZero() {
			t.Fatalf("mark-price-candle1m: expected a parsed row with non-zero Ts, got %+v", candles)
		}
	})

	t.Run("index-candle1m", func(t *testing.T) {
		c := testWsPublicClient()
		_ = c.NewSubscribeIndexCandleService("BTC-USDT", "1m")
		arg := request.WsArg{Channel: "index-candle1m", InstrumentID: "BTC-USDT"}
		raw := wsFirstDataArray(t, c, request.GatewayBusiness, false, arg, 15*time.Second)
		if raw == nil {
			return
		}
		var rows [][]string
		if err := common.JSONUnmarshal(raw, &rows); err != nil {
			t.Fatalf("decode: %v", err)
		}
		candles := parseWsIndexCandles(rows)
		if len(candles) == 0 || candles[0].Timestamp.IsZero() {
			t.Fatalf("index-candle1m: expected a parsed row with non-zero Ts, got %+v", candles)
		}
	})

	t.Run("sprd-candle1m", func(t *testing.T) {
		c := testWsPublicClient()
		_ = c.NewSubscribeSprdCandleService(sprdId, "1m")
		arg := request.WsArg{Channel: "sprd-candle1m", SpreadID: sprdId}
		raw := wsFirstDataArray(t, c, request.GatewayBusiness, false, arg, 15*time.Second)
		if raw == nil {
			return
		}
		var rows [][]string
		if err := common.JSONUnmarshal(raw, &rows); err != nil {
			t.Fatalf("decode: %v", err)
		}
		candles := parseWsSprdCandles(rows)
		if len(candles) == 0 || candles[0].Timestamp.IsZero() {
			t.Fatalf("sprd-candle1m: expected a parsed row with non-zero Ts, got %+v", candles)
		}
	})

	// --- public spread market channels (object payload) ---

	t.Run("sprd-tickers", func(t *testing.T) {
		c := testWsPublicClient()
		_ = c.NewSubscribeSprdTickersService(sprdId)
		arg := request.WsArg{Channel: "sprd-tickers", SpreadID: sprdId}
		raw := wsFirstDataArray(t, c, request.GatewayBusiness, false, arg, 15*time.Second)
		if raw == nil {
			return
		}
		var data []WsSprdTicker
		if err := common.JSONUnmarshal(raw, &data); err != nil {
			t.Fatalf("decode: %v", err)
		}
		assertCovers(t, "ws/sprd-tickers", raw, data)
	})

	t.Run("sprd-public-trades", func(t *testing.T) {
		c := testWsPublicClient()
		_ = c.NewSubscribeSprdPublicTradesService(sprdId)
		arg := request.WsArg{Channel: "sprd-public-trades", SpreadID: sprdId}
		raw := wsFirstDataArray(t, c, request.GatewayBusiness, false, arg, 15*time.Second)
		if raw == nil {
			return
		}
		var data []WsSprdPublicTrade
		if err := common.JSONUnmarshal(raw, &data); err != nil {
			t.Fatalf("decode: %v", err)
		}
		assertCovers(t, "ws/sprd-public-trades", raw, data)
	})

	t.Run("sprd-books5", func(t *testing.T) {
		c := testWsPublicClient()
		_ = c.NewSubscribeSprdBooks5Service(sprdId)
		arg := request.WsArg{Channel: "sprd-books5", SpreadID: sprdId}
		raw := wsFirstDataArray(t, c, request.GatewayBusiness, false, arg, 15*time.Second)
		if raw == nil {
			return
		}
		var data []WsSprdBooks
		if err := common.JSONUnmarshal(raw, &data); err != nil {
			t.Fatalf("decode: %v", err)
		}
		assertCovers(t, "ws/sprd-books5", raw, data)
	})

	// --- private (login) algo / grid / economic-calendar channels ---

	t.Run("orders-algo", func(t *testing.T) {
		c := testWsClient(t)
		_ = c.NewSubscribeOrdersAlgoService(InstTypeSpot)
		arg := request.WsArg{Channel: "orders-algo", InstrumentType: string(InstTypeSpot)}
		raw := wsFirstDataArray(t, c, request.GatewayBusiness, true, arg, 15*time.Second)
		if raw == nil {
			return
		}
		var data []WsAlgoOrder
		if err := common.JSONUnmarshal(raw, &data); err != nil {
			t.Fatalf("decode: %v", err)
		}
		assertCovers(t, "ws/orders-algo", raw, data)
	})

	t.Run("algo-advance", func(t *testing.T) {
		c := testWsClient(t)
		_ = c.NewSubscribeAlgoAdvanceService(InstTypeSpot)
		arg := request.WsArg{Channel: "algo-advance", InstrumentType: string(InstTypeSpot)}
		raw := wsFirstDataArray(t, c, request.GatewayBusiness, true, arg, 15*time.Second)
		if raw == nil {
			return
		}
		var data []WsAdvanceAlgoOrder
		if err := common.JSONUnmarshal(raw, &data); err != nil {
			t.Fatalf("decode: %v", err)
		}
		assertCovers(t, "ws/algo-advance", raw, data)
	})

	t.Run("grid-orders-spot", func(t *testing.T) {
		c := testWsClient(t)
		_ = c.NewSubscribeGridOrdersSpotService(InstTypeSpot)
		arg := request.WsArg{Channel: "grid-orders-spot", InstrumentType: string(InstTypeSpot)}
		raw := wsFirstDataArray(t, c, request.GatewayBusiness, true, arg, 15*time.Second)
		if raw == nil {
			return
		}
		var data []WsGridOrder
		if err := common.JSONUnmarshal(raw, &data); err != nil {
			t.Fatalf("decode: %v", err)
		}
		assertCovers(t, "ws/grid-orders-spot", raw, data)
	})

	t.Run("grid-orders-contract", func(t *testing.T) {
		c := testWsClient(t)
		_ = c.NewSubscribeGridOrdersContractService(InstTypeSwap)
		arg := request.WsArg{Channel: "grid-orders-contract", InstrumentType: string(InstTypeSwap)}
		raw := wsFirstDataArray(t, c, request.GatewayBusiness, true, arg, 15*time.Second)
		if raw == nil {
			return
		}
		var data []WsGridOrder
		if err := common.JSONUnmarshal(raw, &data); err != nil {
			t.Fatalf("decode: %v", err)
		}
		assertCovers(t, "ws/grid-orders-contract", raw, data)
	})

	t.Run("grid-orders-moon", func(t *testing.T) {
		c := testWsClient(t)
		_ = c.NewSubscribeGridOrdersMoonService(InstTypeSpot)
		arg := request.WsArg{Channel: "grid-orders-moon", InstrumentType: string(InstTypeSpot)}
		raw := wsFirstDataArray(t, c, request.GatewayBusiness, true, arg, 15*time.Second)
		if raw == nil {
			return
		}
		var data []WsGridOrder
		if err := common.JSONUnmarshal(raw, &data); err != nil {
			t.Fatalf("decode: %v", err)
		}
		assertCovers(t, "ws/grid-orders-moon", raw, data)
	})

	t.Run("grid-positions", func(t *testing.T) {
		c := testWsClient(t)
		_ = c.NewSubscribeGridPositionsService("123456")
		arg := request.WsArg{Channel: "grid-positions", AlgoID: "123456"}
		raw := wsFirstDataArray(t, c, request.GatewayBusiness, true, arg, 15*time.Second)
		if raw == nil {
			return
		}
		var data []WsGridPosition
		if err := common.JSONUnmarshal(raw, &data); err != nil {
			t.Fatalf("decode: %v", err)
		}
		assertCovers(t, "ws/grid-positions", raw, data)
	})

	t.Run("grid-sub-orders", func(t *testing.T) {
		c := testWsClient(t)
		_ = c.NewSubscribeGridSubOrdersService("123456")
		arg := request.WsArg{Channel: "grid-sub-orders", AlgoID: "123456"}
		raw := wsFirstDataArray(t, c, request.GatewayBusiness, true, arg, 15*time.Second)
		if raw == nil {
			return
		}
		var data []WsGridSubOrder
		if err := common.JSONUnmarshal(raw, &data); err != nil {
			t.Fatalf("decode: %v", err)
		}
		assertCovers(t, "ws/grid-sub-orders", raw, data)
	})

	t.Run("economic-calendar", func(t *testing.T) {
		c := testWsClient(t)
		_ = c.NewSubscribeEconomicCalendarService()
		arg := request.WsArg{Channel: "economic-calendar"}
		// Access to this channel is gated by trading-fee tier: OKX rejects lower
		// tiers with code 64003. Subscribe directly so that capability rejection
		// can be tolerated (login+subscribe path was still validated).
		ctx := wsCtx(t, 15*time.Second)
		out := make(chan []byte, 1)
		done, _, err := c.Subscribe(ctx, request.GatewayBusiness, true, arg, func(message []byte, e error) {
			if e != nil {
				var wsErr *request.WsError
				if errors.As(e, &wsErr) && wsErr.Code == "64003" {
					t.Logf("economic-calendar: account fee tier too low (code=64003) — login+subscribe OK")
					return
				}
				t.Errorf("economic-calendar push error: %v", e)
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
			t.Fatalf("economic-calendar subscribe: %v", err)
		}
		defer close(done)
		select {
		case raw := <-out:
			var data []WsEconomicCalendar
			if err := common.JSONUnmarshal(raw, &data); err != nil {
				t.Fatalf("decode: %v", err)
			}
			assertCovers(t, "ws/economic-calendar", raw, data)
		case <-ctx.Done():
			t.Logf("economic-calendar: no data push within timeout — login+subscribe OK")
		}
	})
}
