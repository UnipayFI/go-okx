package okx

import "testing"

// TestTradeOrder exercises the [Read] order endpoints under /api/v5/trade/*. It
// never places, amends or cancels an order — the [Trade] services are
// implement-only and a tiny-order lifecycle is verified separately by the parent.
func TestTradeOrder(t *testing.T) {
	c := testClient(t)
	_ = c.SyncServerTime(ctx(t))
	cx := ctx(t)

	// account-rate-limit (single object, always populated).
	{
		const label = "account-rate-limit"
		resp, err := c.NewGetAccountRateLimitService().Do(cx)
		if err != nil {
			if tolerable(t, label, err, "51603", "51000", "51001", "50011", "51010") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/trade/account-rate-limit", nil, true)
		if resp != nil {
			assertCovers(t, label, raw, resp)
		}
	}

	// orders-pending (list; account may have none).
	{
		const label = "orders-pending"
		resp, err := c.NewGetOrdersPendingService().Do(cx)
		if err != nil {
			if tolerable(t, label, err, "51603", "51000", "51001", "50011", "51010") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/trade/orders-pending", nil, true)
		if len(resp) > 0 {
			assertCovers(t, label, raw, resp)
		}
	}

	// orders-history SPOT (list; account may have none).
	{
		const label = "orders-history"
		params := map[string]string{"instType": string(InstTypeSpot)}
		resp, err := c.NewGetOrdersHistoryService(InstTypeSpot).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "51603", "51000", "51001", "50011", "51010") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/trade/orders-history", params, true)
		if len(resp) > 0 {
			assertCovers(t, label, raw, resp)
		}
	}

	// orders-history-archive SPOT (list; account may have none).
	{
		const label = "orders-history-archive"
		params := map[string]string{"instType": string(InstTypeSpot)}
		resp, err := c.NewGetOrdersHistoryArchiveService(InstTypeSpot).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "51603", "51000", "51001", "50011", "51010") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/trade/orders-history-archive", params, true)
		if len(resp) > 0 {
			assertCovers(t, label, raw, resp)
		}
	}

	// order lookup by id (no such order on the validating account -> 51603).
	{
		const label = "order"
		_, err := c.NewGetOrderService("BTC-USDT").SetOrdId("1").Do(cx)
		if err != nil {
			if tolerable(t, label, err, "51603", "51000", "51001", "50011", "51010") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
	}
}
