package okx

import (
	"testing"

	"github.com/shopspring/decimal"
)

// TestGrid exercises the [Read] and [Pub] grid endpoints. The state-changing
// [Trade] endpoints (place/amend/stop/close/etc.) are implement-only and never
// executed here. The parent owns any gated order-lifecycle test.
//
// NOTE: the grid endpoints live under /api/v5/tradingBot/grid/* (and
// /api/v5/tradingBot/public/rsi-back-testing), NOT /api/v5/trade/grid/* which
// 404s. See grid.go.
func TestGrid(t *testing.T) {
	c := testClient(t)
	_ = c.SyncServerTime(ctx(t))
	cx := ctx(t)

	// [Read] orders-algo-pending (grid). The validating account has no grid
	// orders, so the data array is typically empty: do the Do call to prove the
	// path + signing + deserialization, and assertCovers only when non-empty.
	{
		label := "grid orders-algo-pending"
		resp, err := c.NewGetGridAlgoPendingService(GridAlgoOrdTypeGrid).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "51603", "51000", "51001", "50011", "51010", "51291") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) > 0 {
			raw := fetchRawGet(t, c, cx, "/api/v5/tradingBot/grid/orders-algo-pending",
				map[string]string{"algoOrdType": string(GridAlgoOrdTypeGrid)}, true)
			assertCovers(t, label, raw, resp)
		} else {
			t.Logf("%s: empty (no grid orders on this account) — path+signing OK", label)
		}
	}

	// [Read] orders-algo-history (grid).
	{
		label := "grid orders-algo-history"
		resp, err := c.NewGetGridAlgoHistoryService(GridAlgoOrdTypeGrid).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "51603", "51000", "51001", "50011", "51010", "51291") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) > 0 {
			raw := fetchRawGet(t, c, cx, "/api/v5/tradingBot/grid/orders-algo-history",
				map[string]string{"algoOrdType": string(GridAlgoOrdTypeGrid)}, true)
			assertCovers(t, label, raw, resp)
		} else {
			t.Logf("%s: empty (no grid history on this account) — path+signing OK", label)
		}
	}

	// [Read] orders-algo-details. Without a real algoId OKX returns 51291 ("bot
	// doesn't exist"); tolerate that (and empty data) — the path+signing are proven.
	{
		label := "grid orders-algo-details"
		resp, err := c.NewGetGridAlgoDetailsService(GridAlgoOrdTypeGrid, "0").Do(cx)
		if err != nil {
			if tolerable(t, label, err, "51603", "51000", "51001", "50011", "51010", "51291") {
				t.Logf("%s: tolerated (no such bot) — path+signing OK", label)
			} else {
				t.Fatalf("%s: %v", label, err)
			}
		} else if resp != nil {
			raw := fetchRawGet(t, c, cx, "/api/v5/tradingBot/grid/orders-algo-details",
				map[string]string{"algoOrdType": string(GridAlgoOrdTypeGrid), "algoId": "0"}, true)
			assertCovers(t, label, raw, resp)
		} else {
			t.Logf("%s: empty — path+signing OK", label)
		}
	}

	// [Read] sub-orders. Requires algoOrdType + algoId + type; without a real
	// algoId OKX returns 51291. Tolerate.
	{
		label := "grid sub-orders"
		resp, err := c.NewGetGridSubOrdersService("0", GridAlgoOrdTypeGrid, GridSubOrderTypeLive).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "51603", "51000", "51001", "50011", "51010", "51291", "50014") {
				t.Logf("%s: tolerated (no such bot) — path+signing OK", label)
			} else {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) > 0 {
			raw := fetchRawGet(t, c, cx, "/api/v5/tradingBot/grid/sub-orders",
				map[string]string{"algoOrdType": string(GridAlgoOrdTypeGrid), "algoId": "0", "type": string(GridSubOrderTypeLive)}, true)
			assertCovers(t, label, raw, resp)
		} else {
			t.Logf("%s: empty — path+signing OK", label)
		}
	}

	// [Read] positions. Without a real algoId OKX returns 51291. Tolerate.
	{
		label := "grid positions"
		resp, err := c.NewGetGridPositionsService(GridAlgoOrdTypeContractGrid, "0").Do(cx)
		if err != nil {
			if tolerable(t, label, err, "51603", "51000", "51001", "50011", "51010", "51291", "50014") {
				t.Logf("%s: tolerated (no such bot) — path+signing OK", label)
			} else {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) > 0 {
			raw := fetchRawGet(t, c, cx, "/api/v5/tradingBot/grid/positions",
				map[string]string{"algoOrdType": string(GridAlgoOrdTypeContractGrid), "algoId": "0"}, true)
			assertCovers(t, label, raw, resp)
		} else {
			t.Logf("%s: empty — path+signing OK", label)
		}
	}

	// [Pub] ai-param (BTC-USDT, grid). Public — no signing.
	{
		label := "grid ai-param"
		pc := testPublicClient()
		resp, err := pc.NewGetGridAIParamService(GridAlgoOrdTypeGrid, "BTC-USDT").Do(cx)
		if err != nil {
			if tolerable(t, label, err, "51000", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) > 0 {
			raw := fetchRawGet(t, pc, cx, "/api/v5/tradingBot/grid/ai-param",
				map[string]string{"algoOrdType": string(GridAlgoOrdTypeGrid), "instId": "BTC-USDT"}, false)
			assertCovers(t, label, raw, resp)
		} else {
			t.Logf("%s: empty — path OK", label)
		}
	}

	// [Pub] min-investment (public stateless calculator) — POST, no signing.
	{
		label := "grid min-investment"
		pc := testPublicClient()
		resp, err := pc.NewGetGridMinInvestmentService("BTC-USDT", GridAlgoOrdTypeGrid,
			decimal.RequireFromString("5000"), decimal.RequireFromString("3000"), 50, GridRunTypeArithmetic).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "51000", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if resp != nil {
			raw := fetchRawPost(t, pc, cx, "/api/v5/tradingBot/grid/min-investment", map[string]any{
				"instId":      "BTC-USDT",
				"algoOrdType": string(GridAlgoOrdTypeGrid),
				"maxPx":       "5000",
				"minPx":       "3000",
				"gridNum":     "50",
				"runType":     string(GridRunTypeArithmetic),
			}, false)
			assertCovers(t, label, raw, resp)
		} else {
			t.Logf("%s: empty — path OK", label)
		}
	}

	// [Pub] rsi-back-testing — public GET.
	{
		label := "grid rsi-back-testing"
		pc := testPublicClient()
		resp, err := pc.NewGetGridRSIBackTestingService("BTC-USDT", "1H", "50", 14, "cross_up", "1M").Do(cx)
		if err != nil {
			if tolerable(t, label, err, "51000", "50011", "50014") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if resp != nil {
			raw := fetchRawGet(t, pc, cx, "/api/v5/tradingBot/public/rsi-back-testing", map[string]string{
				"instId":      "BTC-USDT",
				"timeframe":   "1H",
				"thold":       "50",
				"timePeriod":  "14",
				"triggerCond": "cross_up",
				"duration":    "1M",
			}, false)
			assertCovers(t, label, raw, resp)
		} else {
			t.Logf("%s: empty — path OK", label)
		}
	}
}
