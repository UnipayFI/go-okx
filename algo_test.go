package okx

import (
	"testing"
)

// TestAlgo exercises the algo-order read endpoints under /api/v5/trade/* live,
// asserting that the typed AlgoOrder struct covers every key the real responses
// return. All endpoints are private (signed). The state-changing endpoints
// (order-algo / cancel-algos / amend-algos / cancel-advance-algos) are
// implement-only and never exercised here. This account has no algo orders, so
// the reads typically return empty or "not found"; the tests still verify the
// path + signing.
func TestAlgo(t *testing.T) {
	c := testClient(t)
	_ = c.SyncServerTime(ctx(t))
	cx := ctx(t)

	// --- GET /api/v5/trade/order-algo ---
	{
		const label = "trade/order-algo"
		params := map[string]string{"algoId": "1"}
		resp, err := c.NewGetAlgoOrderService().SetAlgoId("1").Do(cx)
		if err != nil {
			if tolerable(t, label, err, "51603", "51000", "51001", "50011", "51010") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if resp == nil {
			t.Logf("%s: no such algo order — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/trade/order-algo", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/trade/orders-algo-pending ---
	{
		const label = "trade/orders-algo-pending"
		params := map[string]string{"ordType": "conditional"}
		resp, err := c.NewGetAlgoOrdersPendingService(AlgoOrdTypeConditional).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "51603", "51000", "51001", "50011", "51010") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: no pending algo orders — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/trade/orders-algo-pending", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/trade/orders-algo-history ---
	{
		const label = "trade/orders-algo-history"
		params := map[string]string{"ordType": "conditional", "state": "canceled"}
		resp, err := c.NewGetAlgoOrdersHistoryService(AlgoOrdTypeConditional).SetState(AlgoStateCanceled).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "51603", "51000", "51001", "50011", "51010") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: no algo-order history — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/trade/orders-algo-history", params, true)
			assertCovers(t, label, raw, resp)
		}
	}
}
