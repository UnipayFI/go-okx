package okx

import "testing"

// TestAsset exercises the funding-account (asset) read endpoints live and
// asserts that the typed structs cover every key the real responses return. The
// state-changing endpoints (transfer / withdrawal / cancel-withdrawal /
// monthly-statement apply) are implemented but never executed: they move real
// funds or create export jobs on a real account.
func TestAsset(t *testing.T) {
	c := testClient(t)
	_ = c.SyncServerTime(ctx(t))
	cx := ctx(t)

	// --- GET /api/v5/asset/currencies ---
	{
		const label = "asset/currencies"
		params := map[string]string{}
		resp, err := c.NewGetCurrenciesService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "50011") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/asset/currencies", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/asset/balances ---
	{
		const label = "asset/balances"
		params := map[string]string{}
		resp, err := c.NewGetFundingBalanceService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "50011") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/asset/balances", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/asset/non-tradable-assets ---
	{
		const label = "asset/non-tradable-assets"
		params := map[string]string{}
		resp, err := c.NewGetNonTradableAssetsService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "50011") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/asset/non-tradable-assets", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/asset/asset-valuation ---
	{
		const label = "asset/asset-valuation"
		params := map[string]string{}
		resp, err := c.NewGetAssetValuationService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "50011") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if resp == nil {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/asset/asset-valuation", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/asset/transfer-state (transId=1&type=0) ---
	// A bogus transId returns 58129; the path + signing are still exercised.
	{
		const label = "asset/transfer-state"
		params := map[string]string{"transId": "1", "type": "0"}
		resp, err := c.NewGetTransferStateService().SetTransId("1").SetType(AssetTransferTypeWithinAccount).Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "58129", "50011") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/asset/transfer-state", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/asset/bills ---
	{
		const label = "asset/bills"
		params := map[string]string{}
		resp, err := c.NewGetAssetBillsService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "50011") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/asset/bills", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/asset/deposit-address (ccy=USDT) ---
	{
		const label = "asset/deposit-address"
		params := map[string]string{"ccy": "USDT"}
		resp, err := c.NewGetDepositAddressService("USDT").Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "58006", "50011") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/asset/deposit-address", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/asset/deposit-history ---
	{
		const label = "asset/deposit-history"
		params := map[string]string{}
		resp, err := c.NewGetDepositHistoryService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "50011") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/asset/deposit-history", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/asset/withdrawal-history ---
	{
		const label = "asset/withdrawal-history"
		params := map[string]string{}
		resp, err := c.NewGetWithdrawalHistoryService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "50011") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/asset/withdrawal-history", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/asset/deposit-withdraw-status (wdId=1) ---
	// A bogus wdId returns 50026/50015/51000; path + signing are still exercised.
	{
		const label = "asset/deposit-withdraw-status"
		params := map[string]string{"wdId": "1"}
		resp, err := c.NewGetDepositWithdrawStatusService().SetWdId("1").Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "50026", "50015", "51000", "50011") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/asset/deposit-withdraw-status", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/asset/exchange-list ---
	{
		const label = "asset/exchange-list"
		params := map[string]string{}
		resp, err := c.NewGetExchangeListService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "50011") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/asset/exchange-list", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/asset/monthly-statement (month=Jan) ---
	// Requires a prior POST apply (code 51604) before the link exists; tolerate.
	{
		const label = "asset/monthly-statement"
		params := map[string]string{"month": "Jan"}
		resp, err := c.NewGetMonthlyStatementService("Jan").Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "51604", "50011") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if resp == nil {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/asset/monthly-statement", params, true)
			assertCovers(t, label, raw, resp)
		}
	}
}
