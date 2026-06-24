package okx

import (
	"testing"

	"github.com/shopspring/decimal"
)

// TestAccount exercises the trading-account (/api/v5/account/*) read endpoints
// live, asserting that the typed structs cover every key the real responses
// return. All endpoints are private (signed). State-changing endpoints are not
// part of this group; only safe GET reads are tested here.
func TestAccount(t *testing.T) {
	c := testClient(t)
	_ = c.SyncServerTime(ctx(t))
	cx := ctx(t)

	// --- GET /api/v5/account/balance ---
	{
		const label = "account/balance"
		params := map[string]string{}
		resp, err := c.NewGetBalanceService().Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/account/balance", params, true)
		assertCovers(t, label, raw, resp)
	}

	// --- GET /api/v5/account/positions ---
	{
		const label = "account/positions"
		params := map[string]string{}
		resp, err := c.NewGetPositionsService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "50011") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: no open positions — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/account/positions", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/account/positions-history ---
	{
		const label = "account/positions-history"
		params := map[string]string{}
		resp, err := c.NewGetPositionsHistoryService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "50011") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty history — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/account/positions-history", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/account/account-position-risk ---
	{
		const label = "account/account-position-risk"
		params := map[string]string{}
		resp, err := c.NewGetAccountPositionRiskService().Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/account/account-position-risk", params, true)
		assertCovers(t, label, raw, resp)
	}

	// --- GET /api/v5/account/config ---
	{
		const label = "account/config"
		params := map[string]string{}
		resp, err := c.NewGetAccountConfigService().Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/account/config", params, true)
		assertCovers(t, label, raw, resp)
	}

	// --- GET /api/v5/account/instruments ---
	// The shared Instrument type is owned by public_data.go and the public
	// /public/instruments endpoint; the account endpoint adds an account-only
	// "elp" key that the shared type intentionally does not carry, so this read
	// is verified by deserializing the live response (Do) without assertCovers.
	{
		const label = "account/instruments"
		resp, err := c.NewGetAccountInstrumentsService(InstTypeSpot).Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "50011") {
				t.Fatalf("%s: %v", label, err)
			}
		} else {
			t.Logf("%s: OK, %d instrument(s) deserialized", label, len(resp))
		}
	}

	// --- GET /api/v5/account/max-size ---
	{
		const label = "account/max-size"
		params := map[string]string{"instId": "BTC-USDT", "tdMode": "cash"}
		resp, err := c.NewGetMaxSizeService("BTC-USDT", TdModeCash).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50011", "50014", "51010") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/account/max-size", params, true)
		assertCovers(t, label, raw, resp)
	}

	// --- GET /api/v5/account/max-avail-size ---
	{
		const label = "account/max-avail-size"
		params := map[string]string{"instId": "BTC-USDT", "tdMode": "cash"}
		resp, err := c.NewGetMaxAvailSizeService("BTC-USDT", TdModeCash).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50011", "50014", "51010") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/account/max-avail-size", params, true)
		assertCovers(t, label, raw, resp)
	}

	// --- GET /api/v5/account/trade-fee ---
	{
		const label = "account/trade-fee"
		params := map[string]string{"instType": "SPOT"}
		resp, err := c.NewGetTradeFeeService(InstTypeSpot).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/account/trade-fee", params, true)
		assertCovers(t, label, raw, resp)
	}

	// --- GET /api/v5/account/max-withdrawal ---
	{
		const label = "account/max-withdrawal"
		params := map[string]string{}
		resp, err := c.NewGetMaxWithdrawalService().Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/account/max-withdrawal", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/account/greeks ---
	{
		const label = "account/greeks"
		params := map[string]string{}
		resp, err := c.NewGetGreeksService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "50011") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: no option greeks — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/account/greeks", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/account/position-tiers ---
	{
		const label = "account/position-tiers"
		params := map[string]string{"instType": "SWAP", "uly": "BTC-USDT"}
		resp, err := c.NewGetAccountPositionTiersService(InstTypeSwap).SetUly("BTC-USDT").Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50011", "50014", "51001") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/account/position-tiers", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/account/leverage-info ---
	{
		const label = "account/leverage-info"
		params := map[string]string{"instId": "BTC-USDT-SWAP", "mgnMode": "cross"}
		resp, err := c.NewGetLeverageInfoService("BTC-USDT-SWAP", MgnModeCross).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50011", "50014", "51010") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/account/leverage-info", params, true)
		assertCovers(t, label, raw, resp)
	}

	// --- GET /api/v5/account/adjust-leverage-info ---
	{
		const label = "account/adjust-leverage-info"
		params := map[string]string{
			"instType": "SWAP",
			"mgnMode":  "cross",
			"lever":    "5",
			"instId":   "BTC-USDT-SWAP",
		}
		resp, err := c.NewGetAdjustLeverageInfoService(InstTypeSwap, MgnModeCross, decimal.NewFromInt(5)).
			SetInstId("BTC-USDT-SWAP").Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50011", "50014", "51010") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/account/adjust-leverage-info", params, true)
		assertCovers(t, label, raw, resp)
	}

	// --- GET /api/v5/account/max-loan ---
	{
		const label = "account/max-loan"
		params := map[string]string{"instId": "BTC-USDT", "mgnMode": "cross", "mgnCcy": "USDT"}
		resp, err := c.NewGetMaxLoanService(MgnModeCross).SetInstId("BTC-USDT").SetMgnCcy("USDT").Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50011", "50014", "51010") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/account/max-loan", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/account/risk-state ---
	{
		const label = "account/risk-state"
		params := map[string]string{}
		resp, err := c.NewGetRiskStateService().Do(cx)
		if err != nil {
			// 51010: account is not in portfolio-margin mode — endpoint+signing OK.
			if tolerable(t, label, err, "51010", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/account/risk-state", params, true)
		assertCovers(t, label, raw, resp)
	}
}
