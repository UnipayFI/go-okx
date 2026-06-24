package okx

import "testing"

// TestMarket exercises every /api/v5/market/* endpoint live (all public): it
// calls each Service.Do to prove deserialization works, then diffs the typed
// struct against the raw response to guarantee field coverage. Candle endpoints
// (array-of-arrays) are asserted on parsed fields instead of assertCovers.
func TestMarket(t *testing.T) {
	c := testPublicClient()
	cx := ctx(t)

	// GET /api/v5/market/tickers
	{
		params := map[string]string{"instType": "SPOT"}
		resp, err := c.NewGetTickersService(InstTypeSpot).Do(cx)
		if err != nil {
			t.Fatalf("tickers: %v", err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/market/tickers", params, false)
		assertCovers(t, "tickers", raw, resp)
	}

	// GET /api/v5/market/ticker
	{
		params := map[string]string{"instId": "BTC-USDT"}
		resp, err := c.NewGetTickerService("BTC-USDT").Do(cx)
		if err != nil {
			t.Fatalf("ticker: %v", err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/market/ticker", params, false)
		assertCovers(t, "ticker", raw, resp)
	}

	// GET /api/v5/market/books
	{
		params := map[string]string{"instId": "BTC-USDT", "sz": "5"}
		resp, err := c.NewGetBooksService("BTC-USDT").SetSz(5).Do(cx)
		if err != nil {
			t.Fatalf("books: %v", err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/market/books", params, false)
		assertCovers(t, "books", raw, resp)
	}

	// GET /api/v5/market/books-full
	{
		params := map[string]string{"instId": "BTC-USDT", "sz": "5"}
		resp, err := c.NewGetBooksFullService("BTC-USDT").SetSz(5).Do(cx)
		if err != nil {
			t.Fatalf("books-full: %v", err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/market/books-full", params, false)
		assertCovers(t, "books-full", raw, resp)
	}

	// GET /api/v5/market/candles (array-of-arrays)
	{
		candles, err := c.NewGetCandlesService("BTC-USDT").SetLimit(5).Do(cx)
		if err != nil {
			t.Fatalf("candles: %v", err)
		}
		if len(candles) == 0 {
			t.Fatal("candles: empty result")
		}
		if candles[0].Timestamp.IsZero() || candles[0].Close.IsZero() {
			t.Fatalf("candles: parsed row has zero ts or close: %+v", candles[0])
		}
	}

	// GET /api/v5/market/history-candles (array-of-arrays)
	{
		candles, err := c.NewGetHistoryCandlesService("BTC-USDT").SetLimit(5).Do(cx)
		if err != nil {
			t.Fatalf("history-candles: %v", err)
		}
		if len(candles) == 0 {
			t.Fatal("history-candles: empty result")
		}
		if candles[0].Timestamp.IsZero() || candles[0].Close.IsZero() {
			t.Fatalf("history-candles: parsed row has zero ts or close: %+v", candles[0])
		}
	}

	// GET /api/v5/market/trades
	{
		params := map[string]string{"instId": "BTC-USDT", "limit": "5"}
		resp, err := c.NewGetTradesService("BTC-USDT").SetLimit(5).Do(cx)
		if err != nil {
			t.Fatalf("trades: %v", err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/market/trades", params, false)
		assertCovers(t, "trades", raw, resp)
	}

	// GET /api/v5/market/history-trades
	{
		params := map[string]string{"instId": "BTC-USDT", "limit": "5"}
		resp, err := c.NewGetHistoryTradesService("BTC-USDT").SetLimit(5).Do(cx)
		if err != nil {
			t.Fatalf("history-trades: %v", err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/market/history-trades", params, false)
		assertCovers(t, "history-trades", raw, resp)
	}

	// GET /api/v5/market/option/instrument-family-trades
	{
		params := map[string]string{"instFamily": "BTC-USD"}
		resp, err := c.NewGetOptionInstrumentFamilyTradesService("BTC-USD").Do(cx)
		if err != nil {
			t.Fatalf("option family trades: %v", err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/market/option/instrument-family-trades", params, false)
		assertCovers(t, "option family trades", raw, resp)
	}

	// GET /api/v5/market/platform-24-volume
	{
		resp, err := c.NewGetPlatform24VolumeService().Do(cx)
		if err != nil {
			t.Fatalf("platform-24-volume: %v", err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/market/platform-24-volume", nil, false)
		assertCovers(t, "platform-24-volume", raw, resp)
	}

	// GET /api/v5/market/block-tickers
	{
		params := map[string]string{"instType": "SPOT"}
		resp, err := c.NewGetBlockTickersService(InstTypeSpot).Do(cx)
		if err != nil {
			t.Fatalf("block-tickers: %v", err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/market/block-tickers", params, false)
		assertCovers(t, "block-tickers", raw, resp)
	}
}
