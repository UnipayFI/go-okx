package okx

import "testing"

// TestAccountConfig exercises the signed trading-account configuration READ
// endpoints (collateral-assets, mmp-config, move-positions-history) live,
// asserting that the typed structs cover every key the real responses return.
// The setter (POST) endpoints in account_config.go are state-changing and are
// implement-only; they are NOT exercised here.
func TestAccountConfig(t *testing.T) {
	c := testClient(t)
	_ = c.SyncServerTime(ctx(t))
	cx := ctx(t)

	// --- GET /api/v5/account/collateral-assets (Read) ---
	{
		const label = "account/collateral-assets"
		params := map[string]string{}
		resp, err := c.NewGetCollateralAssetsService().Do(cx)
		if err != nil {
			if tolerable(t, label, err, "51010", "50014") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/account/collateral-assets", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/account/mmp-config (Read) ---
	{
		const label = "account/mmp-config"
		params := map[string]string{"instFamily": "BTC-USD"}
		resp, err := c.NewGetMMPConfigService("BTC-USD").Do(cx)
		if err != nil {
			// 51035: no market-maker permission; 50014: missing param; 51010: acct mode.
			if tolerable(t, label, err, "51035", "50014", "51010") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/account/mmp-config", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/account/move-positions-history (Read) ---
	{
		const label = "account/move-positions-history"
		params := map[string]string{}
		resp, err := c.NewGetMovePositionsHistoryService().Do(cx)
		if err != nil {
			if tolerable(t, label, err, "51010", "50014") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/account/move-positions-history", params, true)
			assertCovers(t, label, raw, resp)
		}
	}
}
