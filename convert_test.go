package okx

import "testing"

// TestConvert exercises the convert read endpoints under
// /api/v5/asset/convert/* live, asserting that the typed structs cover every key
// the real responses return. The state-changing POST endpoints (estimate-quote,
// trade) are implemented but intentionally NOT exercised against the real
// account.
func TestConvert(t *testing.T) {
	c := testClient(t)
	_ = c.SyncServerTime(ctx(t))
	cx := ctx(t)

	// --- GET /api/v5/asset/convert/currencies ---
	{
		const label = "asset/convert/currencies"
		params := map[string]string{}
		resp, err := c.NewGetConvertCurrenciesService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "50026", "50011", "51000") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/asset/convert/currencies", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/asset/convert/currency-pair ---
	{
		const label = "asset/convert/currency-pair"
		params := map[string]string{"fromCcy": "USDT", "toCcy": "BTC"}
		resp, err := c.NewGetConvertCurrencyPairService("USDT", "BTC").Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "50026", "50011", "51000") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if resp == nil {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/asset/convert/currency-pair", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/asset/convert/history ---
	{
		const label = "asset/convert/history"
		params := map[string]string{}
		resp, err := c.NewGetConvertHistoryService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "50026", "50011", "51000") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/asset/convert/history", params, true)
			assertCovers(t, label, raw, resp)
		}
	}
}
