package okx

import (
	"testing"
)

// TestRFQ exercises the Block-trading (RFQ) read endpoints live, asserting that
// the typed structs cover every key the real responses return. Private endpoints
// are signed; the public block-trade feeds and block-trade tickers are not.
// State-changing (create/cancel/execute/set) endpoints are implement-only and are
// NOT exercised here.
func TestRFQ(t *testing.T) {
	c := testClient(t)
	_ = c.SyncServerTime(ctx(t))
	cx := ctx(t)

	// --- GET /api/v5/rfq/counterparties (private) ---
	{
		const label = "rfq/counterparties"
		resp, err := c.NewGetRfqCounterpartiesService().Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50030", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/rfq/counterparties", map[string]string{}, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/rfq/rfqs (private) ---
	{
		const label = "rfq/rfqs"
		resp, err := c.NewGetRfqsService().SetLimit(10).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50030", "70000", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/rfq/rfqs", map[string]string{"limit": "10"}, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/rfq/quotes (private) ---
	{
		const label = "rfq/quotes"
		resp, err := c.NewGetRfqQuotesService().SetLimit(10).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50030", "70000", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/rfq/quotes", map[string]string{"limit": "10"}, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/rfq/trades (private) ---
	{
		const label = "rfq/trades"
		resp, err := c.NewGetRfqTradesService().SetLimit(10).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50030", "70000", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/rfq/trades", map[string]string{"limit": "10"}, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/rfq/public-trades (public) ---
	{
		const label = "rfq/public-trades"
		c2 := testPublicClient()
		resp, err := c2.NewGetRfqPublicTradesService().SetLimit(5).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "51000", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c2, cx, "/api/v5/rfq/public-trades", map[string]string{"limit": "5"}, false)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/market/block-tickers (public) ---
	// OKX documents this under /api/v5/rfq/block-tickers, but that path returns
	// HTTP 404; the live endpoint is /api/v5/market/block-tickers.
	{
		const label = "rfq/block-tickers"
		c2 := testPublicClient()
		resp, err := c2.NewGetRfqBlockTickersService(InstTypeSwap).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "51000", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c2, cx, "/api/v5/market/block-tickers", map[string]string{"instType": "SWAP"}, false)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/market/block-ticker (public) ---
	// OKX documents this under /api/v5/rfq/block-ticker, but that path returns
	// HTTP 404; the live endpoint is /api/v5/market/block-ticker.
	{
		const label = "rfq/block-ticker"
		c2 := testPublicClient()
		resp, err := c2.NewGetRfqBlockTickerService("BTC-USDT").Do(cx)
		if err != nil {
			if tolerable(t, label, err, "51000", "51001", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if resp == nil {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c2, cx, "/api/v5/market/block-ticker", map[string]string{"instId": "BTC-USDT"}, false)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/rfq/maker-instrument-settings (private) ---
	{
		const label = "rfq/maker-instrument-settings"
		resp, err := c.NewGetRfqMakerInstrumentSettingsService().Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50030", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/rfq/maker-instrument-settings", map[string]string{}, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/rfq/mmp-config (private) ---
	{
		const label = "rfq/mmp-config"
		resp, err := c.NewGetRfqMmpConfigService().Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50030", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if resp == nil {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/rfq/mmp-config", map[string]string{}, true)
			assertCovers(t, label, raw, resp)
		}
	}
}
