package okx

import "testing"

// TestFinanceSavings exercises the Simple Earn Flexible (savings) read endpoints
// under /api/v5/finance/savings/* live, asserting that the typed structs cover
// every key the real responses return. The validating account holds no savings,
// so balance and lending-history return empty data arrays (coverage skipped, Do
// call kept). The state-changing POST endpoints (purchase-redemption,
// set-lending-rate) are implemented but intentionally NOT exercised.
func TestFinanceSavings(t *testing.T) {
	c := testClient(t)
	_ = c.SyncServerTime(ctx(t))
	cx := ctx(t)

	// --- GET /api/v5/finance/savings/balance ---
	{
		const label = "finance/savings/balance"
		params := map[string]string{}
		resp, err := c.NewGetSavingsBalanceService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "50030", "51010", "51000", "50014", "50011") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/finance/savings/balance", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/finance/savings/lending-history ---
	{
		const label = "finance/savings/lending-history"
		params := map[string]string{}
		resp, err := c.NewGetSavingsLendingHistoryService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "50030", "51010", "51000", "50014", "50011") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/finance/savings/lending-history", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/finance/savings/lending-rate-summary ---
	{
		const label = "finance/savings/lending-rate-summary"
		params := map[string]string{}
		resp, err := c.NewGetSavingsLendingRateSummaryService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "50030", "51010", "51000", "50014", "50011") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/finance/savings/lending-rate-summary", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/finance/savings/lending-rate-history ---
	{
		const label = "finance/savings/lending-rate-history"
		params := map[string]string{}
		resp, err := c.NewGetSavingsLendingRateHistoryService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "50030", "51010", "51000", "50014", "50011") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/finance/savings/lending-rate-history", params, true)
			assertCovers(t, label, raw, resp)
		}
	}
}
