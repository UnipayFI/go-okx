package okx

import "testing"

// TestRecurring exercises the recurring-buy (DCA) read endpoints under
// /api/v5/tradingBot/recurring/* live, asserting that the typed structs cover
// every key the real responses return. The validating account holds no
// recurring-buy strategies, so the list endpoints return empty data (coverage
// skipped, Do call kept) and the by-id detail/sub-order endpoints return 51291
// ("the bot doesn't exist or has already stopped"). The state-changing POST
// endpoints (order-algo, amend-order-algo, stop-order-algo) are implemented but
// intentionally NOT exercised.
func TestRecurring(t *testing.T) {
	c := testClient(t)
	_ = c.SyncServerTime(ctx(t))
	cx := ctx(t)

	// --- GET /api/v5/tradingBot/recurring/orders-algo-pending ---
	{
		const label = "tradingBot/recurring/orders-algo-pending"
		params := map[string]string{}
		resp, err := c.NewGetRecurringOrdersPendingService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "51603", "51000", "51001", "50011", "51010", "51291") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/tradingBot/recurring/orders-algo-pending", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/tradingBot/recurring/orders-algo-history ---
	{
		const label = "tradingBot/recurring/orders-algo-history"
		params := map[string]string{}
		resp, err := c.NewGetRecurringOrdersHistoryService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "51603", "51000", "51001", "50011", "51010", "51291") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/tradingBot/recurring/orders-algo-history", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/tradingBot/recurring/orders-algo-details ---
	{
		const label = "tradingBot/recurring/orders-algo-details"
		// Resolve a real algoId from the pending list when one exists; otherwise
		// the by-id lookup returns 51291 (no such bot), which we tolerate.
		algoId := "0"
		if pending, perr := c.NewGetRecurringOrdersPendingService().Do(cx); perr == nil && len(pending) > 0 {
			algoId = pending[0].AlgoID
		}
		params := map[string]string{"algoId": algoId}
		resp, err := c.NewGetRecurringOrderDetailsService(algoId).Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "51603", "51000", "51001", "50011", "51010", "51291") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if resp == nil {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/tradingBot/recurring/orders-algo-details", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/tradingBot/recurring/sub-orders ---
	{
		const label = "tradingBot/recurring/sub-orders"
		algoId := "0"
		if pending, perr := c.NewGetRecurringOrdersPendingService().Do(cx); perr == nil && len(pending) > 0 {
			algoId = pending[0].AlgoID
		}
		params := map[string]string{"algoId": algoId}
		resp, err := c.NewGetRecurringSubOrdersService(algoId).Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "51603", "51000", "51001", "50011", "51010", "51291") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/tradingBot/recurring/sub-orders", params, true)
			assertCovers(t, label, raw, resp)
		}
	}
}
