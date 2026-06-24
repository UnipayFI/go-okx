package okx

import (
	"errors"
	"testing"

	"github.com/UnipayFI/go-okx/client"
)

// TestTradeConvert exercises the easy-convert and one-click-repay read endpoints
// under /api/v5/trade/* live, asserting the typed structs cover every key the
// real responses return. The state-changing POST endpoints (easy-convert,
// one-click-repay, one-click-repay-v2) are implemented but intentionally NOT
// exercised against the real account.
func TestTradeConvert(t *testing.T) {
	c := testClient(t)
	_ = c.SyncServerTime(ctx(t))
	cx := ctx(t)

	// --- GET /api/v5/trade/easy-convert-currency-list ---
	{
		const label = "trade/easy-convert-currency-list"
		params := map[string]string{}
		resp, err := c.NewGetEasyConvertCurrencyListService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "51010", "51000", "51001", "50011") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if resp == nil {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/trade/easy-convert-currency-list", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/trade/easy-convert-history ---
	{
		const label = "trade/easy-convert-history"
		params := map[string]string{}
		resp, err := c.NewGetEasyConvertHistoryService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "51010", "51000", "51001", "50011") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/trade/easy-convert-history", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/trade/one-click-repay-currency-list ---
	{
		const label = "trade/one-click-repay-currency-list"
		params := map[string]string{}
		resp, err := c.NewGetOneClickRepayCurrencyListService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "51010", "51000", "51001", "50011") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/trade/one-click-repay-currency-list", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/trade/one-click-repay-history ---
	{
		const label = "trade/one-click-repay-history"
		params := map[string]string{}
		resp, err := c.NewGetOneClickRepayHistoryService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "51010", "51000", "51001", "50011") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/trade/one-click-repay-history", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/trade/one-click-repay-currency-list-v2 (drop on 404) ---
	{
		const label = "trade/one-click-repay-currency-list-v2"
		params := map[string]string{}
		resp, err := c.NewGetOneClickRepayCurrencyListV2Service().Do(cx)
		if err != nil {
			if tradeConvertNotFound(err) {
				t.Logf("%s: endpoint not available (404) — skipped", label)
			} else if !tolerable(t, label, err, "51010", "51000", "51001", "50011") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/trade/one-click-repay-currency-list-v2", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/trade/one-click-repay-history-v2 (drop on 404) ---
	{
		const label = "trade/one-click-repay-history-v2"
		params := map[string]string{}
		resp, err := c.NewGetOneClickRepayHistoryV2Service().Do(cx)
		if err != nil {
			if tradeConvertNotFound(err) {
				t.Logf("%s: endpoint not available (404) — skipped", label)
			} else if !tolerable(t, label, err, "51010", "51000", "51001", "50011") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/trade/one-click-repay-history-v2", params, true)
			assertCovers(t, label, raw, resp)
		}
	}
}

// tradeConvertNotFound reports whether err is OKX's wrong-path 404 (code "404").
func tradeConvertNotFound(err error) bool {
	var apiErr *client.APIError
	return errors.As(err, &apiErr) && apiErr.Code == "404"
}
