package okx

import (
	"testing"
)

// TestTradeFills exercises the transaction-fill (/api/v5/trade/fills*) read
// endpoints live, asserting that the typed Fill struct covers every key the real
// responses return. Both endpoints are private (signed). No order is placed.
func TestTradeFills(t *testing.T) {
	c := testClient(t)
	_ = c.SyncServerTime(ctx(t))
	cx := ctx(t)

	// --- GET /api/v5/trade/fills (last 3 days, no params) ---
	{
		const label = "trade/fills"
		params := map[string]string{}
		resp, err := c.NewGetFillsService().Do(cx)
		if err != nil {
			if tolerable(t, label, err, "51603", "51000", "51001", "50011", "51010") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: no fills in the last 3 days — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/trade/fills", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/trade/fills-history (instType=SPOT) ---
	{
		const label = "trade/fills-history"
		params := map[string]string{"instType": "SPOT"}
		resp, err := c.NewGetFillsHistoryService(InstTypeSpot).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "51603", "51000", "51001", "50011", "51010") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty SPOT fill history — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/trade/fills-history", params, true)
			assertCovers(t, label, raw, resp)
		}
	}
}
