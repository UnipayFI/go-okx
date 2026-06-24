package okx

import "testing"

// TestCopyTrading exercises the signed copy-trading READ endpoints (plus the two
// public endpoints) live, asserting the typed structs cover every key the real
// responses return. The validating account has neither led nor copied any trades
// and is not on the copy-trading allowlist, so most signed reads return empty or
// a capability error (59285 not-a-lead/copy-trader, 59263 allowlist-only, 50030
// no-permission); those are tolerated — the request path and signing are still
// verified. The state-changing POST endpoints in copytrading.go are
// implement-only and are NOT exercised here.
func TestCopyTrading(t *testing.T) {
	c := testClient(t)
	_ = c.SyncServerTime(ctx(t))
	cx := ctx(t)

	// --- GET /api/v5/copytrading/public-config (public) ---
	c2 := testPublicClient()
	{
		const label = "copytrading/public-config"
		params := map[string]string{"instType": "SWAP"}
		resp, err := c2.NewGetCopyPublicConfigService().SetInstType(InstTypeSwap).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if resp == nil {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c2, cx, "/api/v5/copytrading/public-config", params, false)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/copytrading/public-lead-traders (public) ---
	// Also captures a real uniqueCode to reuse in the signed per-trader reads.
	var uniqueCode string
	{
		const label = "copytrading/public-lead-traders"
		params := map[string]string{"instType": "SWAP", "limit": "5"}
		resp, err := c2.NewGetCopyPublicLeadTradersService().SetInstType(InstTypeSwap).SetLimit(5).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if resp == nil || len(resp.Ranks) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			uniqueCode = resp.Ranks[0].UniqueCode
			raw := fetchRawGet(t, c2, cx, "/api/v5/copytrading/public-lead-traders", params, false)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/copytrading/instruments (Read) ---
	{
		const label = "copytrading/instruments"
		params := map[string]string{"instType": "SWAP"}
		resp, err := c.NewGetCopyInstrumentsService().SetInstType(InstTypeSwap).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "59285", "59263", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/copytrading/instruments", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/copytrading/current-subpositions (Read) ---
	{
		const label = "copytrading/current-subpositions"
		params := map[string]string{"instType": "SWAP"}
		resp, err := c.NewGetCopyCurrentSubpositionsService().SetInstType(InstTypeSwap).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "59285", "59263", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/copytrading/current-subpositions", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/copytrading/subpositions-history (Read) ---
	{
		const label = "copytrading/subpositions-history"
		params := map[string]string{"instType": "SWAP"}
		resp, err := c.NewGetCopySubpositionsHistoryService().SetInstType(InstTypeSwap).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "59285", "59263", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/copytrading/subpositions-history", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/copytrading/profit-sharing-details (Read) ---
	{
		const label = "copytrading/profit-sharing-details"
		params := map[string]string{"instType": "SWAP"}
		resp, err := c.NewGetCopyProfitSharingDetailsService().SetInstType(InstTypeSwap).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "59285", "59263", "50030", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/copytrading/profit-sharing-details", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/copytrading/total-profit-sharing (Read) ---
	{
		const label = "copytrading/total-profit-sharing"
		params := map[string]string{"instType": "SWAP"}
		resp, err := c.NewGetCopyTotalProfitSharingService().SetInstType(InstTypeSwap).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "59285", "59263", "50030", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/copytrading/total-profit-sharing", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/copytrading/unrealized-profit-sharing-details (Read) ---
	{
		const label = "copytrading/unrealized-profit-sharing-details"
		params := map[string]string{"instType": "SWAP"}
		resp, err := c.NewGetCopyUnrealizedProfitSharingDetailsService().SetInstType(InstTypeSwap).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "59285", "59263", "50030", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/copytrading/unrealized-profit-sharing-details", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/copytrading/total-unrealized-profit-sharing (Read) ---
	{
		const label = "copytrading/total-unrealized-profit-sharing"
		params := map[string]string{"instType": "SWAP"}
		resp, err := c.NewGetCopyTotalUnrealizedProfitSharingService().SetInstType(InstTypeSwap).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "59285", "59263", "50030", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/copytrading/total-unrealized-profit-sharing", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/copytrading/config (Read) ---
	{
		const label = "copytrading/config"
		resp, err := c.NewGetCopyConfigService().Do(cx)
		if err != nil {
			if tolerable(t, label, err, "59285", "59263", "50030", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if resp == nil {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/copytrading/config", map[string]string{}, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/copytrading/lead-traders (Read) ---
	{
		const label = "copytrading/lead-traders"
		params := map[string]string{"instType": "SWAP", "limit": "5"}
		resp, err := c.NewGetCopyLeadTradersService().SetInstType(InstTypeSwap).SetLimit(5).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "59285", "59263", "50030", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if resp == nil || len(resp.Ranks) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			if uniqueCode == "" {
				uniqueCode = resp.Ranks[0].UniqueCode
			}
			raw := fetchRawGet(t, c, cx, "/api/v5/copytrading/lead-traders", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/copytrading/copy-settings (Read) ---
	{
		const label = "copytrading/copy-settings"
		uc := uniqueCode
		if uc == "" {
			uc = "6308333277A08132"
		}
		params := map[string]string{"instType": "SWAP", "uniqueCode": uc}
		resp, err := c.NewGetCopySettingsService(InstTypeSwap, uc).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "59285", "59263", "50030", "51000", "50014", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if resp == nil {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/copytrading/copy-settings", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/copytrading/batch-leverage-info (Read) ---
	{
		const label = "copytrading/batch-leverage-info"
		uc := uniqueCode
		if uc == "" {
			uc = "6308333277A08132"
		}
		params := map[string]string{"mgnMode": "cross", "uniqueCode": uc, "instId": "BTC-USDT-SWAP"}
		resp, err := c.NewGetCopyBatchLeverageInfoService(MgnModeCross, uc).SetInstId("BTC-USDT-SWAP").Do(cx)
		if err != nil {
			if tolerable(t, label, err, "59285", "59263", "50030", "51000", "50014", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/copytrading/batch-leverage-info", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/copytrading/copy-traders (Read) ---
	{
		const label = "copytrading/copy-traders"
		uc := uniqueCode
		if uc == "" {
			uc = "6308333277A08132"
		}
		params := map[string]string{"instType": "SWAP", "uniqueCode": uc, "limit": "5"}
		resp, err := c.NewGetCopyTradersService(uc).SetInstType(InstTypeSwap).SetLimit(5).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "59285", "59263", "50030", "50014", "51000", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if resp == nil || len(resp.CopyTraders) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/copytrading/copy-traders", params, true)
			assertCovers(t, label, raw, resp)
		}
	}
}
