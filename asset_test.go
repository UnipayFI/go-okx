package okx

import (
	"testing"
	"time"

	"github.com/UnipayFI/go-okx/common"
)

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

// TestDepositWithdrawStatusEstCompleteTime pins estCompleteTime's wire form:
// UTC+8 wall-clock text (the layout OKX documents and real responses show), or
// "" when there is no estimate.
func TestDepositWithdrawStatusEstCompleteTime(t *testing.T) {
	const payload = `[{"wdId":"298527196","txId":"","state":"Pending withdrawal","estCompleteTime":"04/26/2025, 8:04:27 AM"},` +
		`{"wdId":"200045249","txId":"","state":"Pending withdrawal","estCompleteTime":"01/09/2023, 8:10:48 PM"},` +
		`{"wdId":"94075211","txId":"0x38","state":"Withdrawal complete","estCompleteTime":""}]`
	var got []DepositWithdrawStatus
	if err := common.JSONUnmarshal([]byte(payload), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	want := []time.Time{
		time.Date(2025, 4, 26, 0, 4, 27, 0, time.UTC),
		time.Date(2023, 1, 9, 12, 10, 48, 0, time.UTC),
		{},
	}
	for i, w := range want {
		if got[i].EstimatedCompleteTime != w {
			t.Errorf("[%d] estCompleteTime = %v, want %v", i, got[i].EstimatedCompleteTime, w)
		}
	}
	if got[0].WithdrawalID != "298527196" || got[2].TransactionID != "0x38" || got[2].State != "Withdrawal complete" {
		t.Errorf("other fields = %+v", got)
	}
	out, err := common.JSONMarshal(got)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(out) != payload {
		t.Errorf("marshal = %s, want %s", out, payload)
	}

	// Decoding into a reused value without the key keeps the instant, and a
	// caller's struct embedding the type keeps its own fields.
	s := got[0]
	if err := common.JSONUnmarshal([]byte(`{"state":"done"}`), &s); err != nil {
		t.Fatalf("re-unmarshal: %v", err)
	}
	if s.EstimatedCompleteTime != want[0] || s.State != "done" {
		t.Errorf("re-unmarshal = %+v", s)
	}
	var outer struct {
		DepositWithdrawStatus
		Note string `json:"note"`
	}
	const embedded = `{"wdId":"1","txId":"","state":"","estCompleteTime":"04/26/2025, 8:04:27 AM","note":"n"}`
	if err := common.JSONUnmarshal([]byte(embedded), &outer); err != nil {
		t.Fatalf("embedded unmarshal: %v", err)
	}
	if outer.Note != "n" || outer.EstimatedCompleteTime != want[0] {
		t.Errorf("embedded unmarshal = %+v", outer)
	}
	if out, err := common.JSONMarshal(outer); err != nil || string(out) != embedded {
		t.Errorf("embedded marshal = %s (%v), want %s", out, err, embedded)
	}
}
