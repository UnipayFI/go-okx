package okx

import (
	"testing"
)

// TestSubAccount exercises the sub-account read endpoints
// (/api/v5/users/*, /api/v5/account/subaccount/*, /api/v5/asset/subaccount/*)
// live, asserting that the typed structs cover every key the real responses
// return. All endpoints are private (signed). State-changing endpoints
// (create/modify/delete apikey, transfer, set-transfer-out) are implement-only
// and not exercised here. This account has no sub-accounts, so most reads return
// empty or "not found"; the tests still verify the path + signing.
func TestSubAccount(t *testing.T) {
	c := testClient(t)
	_ = c.SyncServerTime(ctx(t))
	cx := ctx(t)

	// --- GET /api/v5/users/subaccount/list ---
	{
		const label = "users/subaccount/list"
		params := map[string]string{}
		resp, err := c.NewGetSubAccountListService().Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50011", "58115", "59510") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: no sub-accounts — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/users/subaccount/list", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/users/subaccount/apikey ---
	{
		const label = "users/subaccount/apikey"
		params := map[string]string{"subAcct": "test"}
		resp, err := c.NewGetSubAccountApiKeyService("test").Do(cx)
		if err != nil {
			if tolerable(t, label, err, "51000", "58115", "59510", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: no api keys — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/users/subaccount/apikey", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/account/subaccount/balances ---
	{
		const label = "account/subaccount/balances"
		params := map[string]string{"subAcct": "test"}
		resp, err := c.NewGetSubAccountTradingBalancesService("test").Do(cx)
		if err != nil {
			if tolerable(t, label, err, "58115", "59510", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/account/subaccount/balances", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/asset/subaccount/balances ---
	{
		const label = "asset/subaccount/balances"
		params := map[string]string{"subAcct": "test"}
		resp, err := c.NewGetSubAccountFundingBalancesService("test").Do(cx)
		if err != nil {
			if tolerable(t, label, err, "59510", "58115", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/asset/subaccount/balances", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/account/subaccount/max-withdrawal ---
	{
		const label = "account/subaccount/max-withdrawal"
		params := map[string]string{"subAcct": "test"}
		resp, err := c.NewGetSubAccountMaxWithdrawalService("test").Do(cx)
		if err != nil {
			if tolerable(t, label, err, "58115", "59510", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/account/subaccount/max-withdrawal", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/asset/subaccount/bills ---
	{
		const label = "asset/subaccount/bills"
		params := map[string]string{}
		resp, err := c.NewGetSubAccountBillsService().Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50011", "58115", "59510") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/asset/subaccount/bills", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/users/entrust-subaccount-list ---
	{
		const label = "users/entrust-subaccount-list"
		params := map[string]string{}
		resp, err := c.NewGetEntrustSubAccountListService().Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50011", "58016", "58115", "59510") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/users/entrust-subaccount-list", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/asset/subaccount/managed-subaccount-bills ---
	// The /api/v5/account/subaccount/managed-subaccount-bills path returns HTTP
	// 404; the real route is under /api/v5/asset/. 58016 is returned when the
	// caller is not a trading-team master account.
	{
		const label = "asset/subaccount/managed-subaccount-bills"
		params := map[string]string{}
		resp, err := c.NewGetManagedSubAccountBillsService().Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50011", "58016", "58115", "59510") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/asset/subaccount/managed-subaccount-bills", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// NOTE: GET /api/v5/users/custody-trading-subaccount-list returns HTTP 404
	// under every base-path convention (/users, /account/subaccount,
	// /asset/subaccount); it does not exist, so it is intentionally not
	// implemented. The custody (entrusted) trading sub-account list is served by
	// /api/v5/users/entrust-subaccount-list above.
}
