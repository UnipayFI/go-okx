package okx

import "testing"

func TestMarketIndex(t *testing.T) {
	c := testPublicClient()
	cx := ctx(t)

	// index-tickers (instId=BTC-USDT) -> []IndexTicker
	{
		params := map[string]string{"instId": "BTC-USDT"}
		resp, err := c.NewGetIndexTickersService().SetInstId("BTC-USDT").Do(cx)
		if err != nil {
			t.Fatalf("index-tickers: %v", err)
		}
		if len(resp) == 0 {
			t.Fatal("index-tickers: empty response")
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/market/index-tickers", params, false)
		assertCovers(t, "index-tickers", raw, resp)
	}

	// index-candles (instId=BTC-USDT&limit=5) -> []IndexCandle (6 cols)
	{
		resp, err := c.NewGetIndexCandlesService("BTC-USDT").SetLimit(5).Do(cx)
		if err != nil {
			t.Fatalf("index-candles: %v", err)
		}
		if len(resp) == 0 {
			t.Fatal("index-candles: empty response")
		}
		if resp[0].Timestamp.IsZero() || resp[0].Close.IsZero() {
			t.Errorf("index-candles: parsed row missing Ts/Close: %+v", resp[0])
		}
	}

	// history-index-candles (instId=BTC-USDT&limit=5) -> []IndexCandle (6 cols)
	{
		resp, err := c.NewGetHistoryIndexCandlesService("BTC-USDT").SetLimit(5).Do(cx)
		if err != nil {
			t.Fatalf("history-index-candles: %v", err)
		}
		if len(resp) == 0 {
			t.Fatal("history-index-candles: empty response")
		}
		if resp[0].Timestamp.IsZero() || resp[0].Close.IsZero() {
			t.Errorf("history-index-candles: parsed row missing Ts/Close: %+v", resp[0])
		}
	}

	// mark-price-candles (instId=BTC-USDT-SWAP&limit=5) -> []IndexCandle (6 cols)
	{
		resp, err := c.NewGetMarkPriceCandlesService("BTC-USDT-SWAP").SetLimit(5).Do(cx)
		if err != nil {
			t.Fatalf("mark-price-candles: %v", err)
		}
		if len(resp) == 0 {
			t.Fatal("mark-price-candles: empty response")
		}
		if resp[0].Timestamp.IsZero() || resp[0].Close.IsZero() {
			t.Errorf("mark-price-candles: parsed row missing Ts/Close: %+v", resp[0])
		}
	}

	// history-mark-price-candles (instId=BTC-USDT-SWAP&limit=5) -> []IndexCandle (6 cols)
	{
		resp, err := c.NewGetHistoryMarkPriceCandlesService("BTC-USDT-SWAP").SetLimit(5).Do(cx)
		if err != nil {
			t.Fatalf("history-mark-price-candles: %v", err)
		}
		if len(resp) == 0 {
			t.Fatal("history-mark-price-candles: empty response")
		}
		if resp[0].Timestamp.IsZero() || resp[0].Close.IsZero() {
			t.Errorf("history-mark-price-candles: parsed row missing Ts/Close: %+v", resp[0])
		}
	}

	// exchange-rate (no params) -> *ExchangeRate
	{
		resp, err := c.NewGetExchangeRateService().Do(cx)
		if err != nil {
			t.Fatalf("exchange-rate: %v", err)
		}
		if resp == nil {
			t.Fatal("exchange-rate: nil response")
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/market/exchange-rate", nil, false)
		assertCovers(t, "exchange-rate", raw, resp)
	}

	// index-components (index=BTC-USDT) -> *IndexComponents (data is an OBJECT)
	{
		params := map[string]string{"index": "BTC-USDT"}
		resp, err := c.NewGetIndexComponentsService("BTC-USDT").Do(cx)
		if err != nil {
			t.Fatalf("index-components: %v", err)
		}
		if resp == nil {
			t.Fatal("index-components: nil response")
		}
		// "data" here is a single JSON object, not an array; assertCovers expects
		// the envelope's array shape, so wrap the raw object in a one-element
		// array to line both sides up.
		raw := fetchRawGet(t, c, cx, "/api/v5/market/index-components", params, false)
		wrapped := append(append([]byte{'['}, raw...), ']')
		assertCovers(t, "index-components", wrapped, resp)
	}
}
