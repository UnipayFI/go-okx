package okx

import (
	"testing"
)

// TestSpread exercises the spread-trading endpoints. The public market-data
// sub-endpoints are tested unsigned via testPublicClient(); the private order
// reads are signed and tolerate 50030 (this key lacks the sprd permission). The
// spread candle endpoints live under /api/v5/market/sprd-(history-)candles, not
// /api/v5/sprd/, and the spread algo-order read endpoints return HTTP 404 on the
// live host (not currently deployed), so they are implemented but not live-tested.
func TestSpread(t *testing.T) {
	c2 := testPublicClient()
	cx := ctx(t)

	// pick a known-active spread (BCH cash/swap) so books/ticker/candles are
	// non-empty; fall back to the first live spread if absent.
	sprdId := "BCH-USDT_BCH-USDT-SWAP"

	// --- [Pub] GET /api/v5/sprd/spreads ---
	{
		label := "sprd/spreads"
		resp, err := c2.NewGetSprdSpreadsService().Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := fetchRawGet(t, c2, cx, "/api/v5/sprd/spreads", nil, false)
		if len(resp) > 0 {
			assertCovers(t, label, raw, resp)
		}
	}

	// --- [Pub] GET /api/v5/sprd/books ---
	{
		label := "sprd/books"
		resp, err := c2.NewGetSprdBooksService(sprdId).SetSz(5).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50011", "51000", "50014") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := fetchRawGet(t, c2, cx, "/api/v5/sprd/books", map[string]string{"sprdId": sprdId, "sz": "5"}, false)
		if resp != nil {
			assertCovers(t, label, raw, resp)
		}
	}

	// --- [Pub] GET /api/v5/sprd/ticker ---
	{
		label := "sprd/ticker"
		resp, err := c2.NewGetSprdTickerService(sprdId).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50011", "51000", "50014") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := fetchRawGet(t, c2, cx, "/api/v5/sprd/ticker", map[string]string{"sprdId": sprdId}, false)
		if resp != nil {
			assertCovers(t, label, raw, resp)
		}
	}

	// --- [Pub] GET /api/v5/sprd/public-trades ---
	{
		label := "sprd/public-trades"
		resp, err := c2.NewGetSprdPublicTradesService().Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := fetchRawGet(t, c2, cx, "/api/v5/sprd/public-trades", nil, false)
		if len(resp) > 0 {
			assertCovers(t, label, raw, resp)
		}
	}

	// --- [Pub] GET /api/v5/market/sprd-candles ---
	{
		label := "market/sprd-candles"
		resp, err := c2.NewGetSprdCandlesService(sprdId).SetLimit(10).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50011", "51000", "50014") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Errorf("%s: expected at least one candle, got 0", label)
		}
	}

	// --- [Pub] GET /api/v5/market/sprd-history-candles ---
	{
		label := "market/sprd-history-candles"
		resp, err := c2.NewGetSprdHistoryCandlesService(sprdId).SetLimit(10).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50011", "51000", "50014") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Errorf("%s: expected at least one history candle, got 0", label)
		}
	}

	// --- private order reads (signed; this key lacks sprd permission -> 50030) ---
	c := testClient(t)
	_ = c.SyncServerTime(cx)

	// --- [Read] GET /api/v5/sprd/order ---
	{
		label := "sprd/order"
		_, err := c.NewGetSprdOrderService().SetOrdId("1").Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50030", "51000", "50014", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
	}

	// --- [Read] GET /api/v5/sprd/orders-pending ---
	{
		label := "sprd/orders-pending"
		resp, err := c.NewGetSprdOrdersPendingService().Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50030", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/sprd/orders-pending", nil, true)
		if len(resp) > 0 {
			assertCovers(t, label, raw, resp)
		}
	}

	// --- [Read] GET /api/v5/sprd/orders-history ---
	{
		label := "sprd/orders-history"
		resp, err := c.NewGetSprdOrdersHistoryService().Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50030", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/sprd/orders-history", nil, true)
		if len(resp) > 0 {
			assertCovers(t, label, raw, resp)
		}
	}

	// --- [Read] GET /api/v5/sprd/orders-history-archive ---
	{
		label := "sprd/orders-history-archive"
		resp, err := c.NewGetSprdOrdersHistoryArchiveService().Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50030", "404", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/sprd/orders-history-archive", nil, true)
		if len(resp) > 0 {
			assertCovers(t, label, raw, resp)
		}
	}

	// --- [Read] GET /api/v5/sprd/trades ---
	{
		label := "sprd/trades"
		resp, err := c.NewGetSprdTradesService().Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50030", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/sprd/trades", nil, true)
		if len(resp) > 0 {
			assertCovers(t, label, raw, resp)
		}
	}
}
