package okx

import (
	"testing"
)

// TestPublicData exercises every /api/v5/public/* market-reference endpoint
// implemented in public_data.go against the live OKX REST API. For each endpoint
// it (a) calls the typed Service to prove the request works and deserializes,
// then (b) re-fetches the raw "data" array and asserts the struct covers every
// real JSON key.
func TestPublicData(t *testing.T) {
	c := testPublicClient()
	cx := ctx(t)

	// instruments (SPOT)
	{
		params := map[string]string{"instType": string(InstTypeSpot)}
		resp, err := c.NewGetInstrumentsService(InstTypeSpot).Do(cx)
		if err != nil {
			t.Fatalf("instruments: %v", err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/public/instruments", params, false)
		assertCovers(t, "instruments", raw, resp)
	}

	// estimated-price — fetch a live FUTURES instId, tolerate 51001 (not yet
	// within the estimation window / not found).
	{
		futs, err := c.NewGetInstrumentsService(InstTypeFutures).Do(cx)
		if err != nil {
			t.Fatalf("instruments FUTURES for estimated-price: %v", err)
		}
		if len(futs) == 0 {
			t.Fatal("estimated-price: no FUTURES instruments available")
		}
		instId := futs[0].InstrumentID
		params := map[string]string{"instId": instId}
		resp, err := c.NewGetEstimatedPriceService(instId).Do(cx)
		if err != nil && !tolerable(t, "estimated-price", err, "51001") {
			t.Fatalf("estimated-price: %v", err)
		} else if err == nil {
			raw := fetchRawGet(t, c, cx, "/api/v5/public/estimated-price", params, false)
			assertCovers(t, "estimated-price", raw, resp)
		}
	}

	// delivery-exercise-history (FUTURES, uly=BTC-USD)
	{
		params := map[string]string{"instType": string(InstTypeFutures), "uly": "BTC-USD"}
		resp, err := c.NewGetDeliveryExerciseHistoryService(InstTypeFutures).SetUly("BTC-USD").Do(cx)
		if err != nil {
			t.Fatalf("delivery-exercise-history: %v", err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/public/delivery-exercise-history", params, false)
		assertCovers(t, "delivery-exercise-history", raw, resp)
	}

	// settlement-history (FUTURES, instFamily=BTC-USD)
	{
		params := map[string]string{"instType": string(InstTypeFutures), "instFamily": "BTC-USD"}
		resp, err := c.NewGetSettlementHistoryService(InstTypeFutures, "BTC-USD").Do(cx)
		if err != nil {
			t.Fatalf("settlement-history: %v", err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/public/settlement-history", params, false)
		assertCovers(t, "settlement-history", raw, resp)
	}

	// funding-rate (SWAP)
	{
		params := map[string]string{"instId": "BTC-USDT-SWAP"}
		resp, err := c.NewGetFundingRateService("BTC-USDT-SWAP").Do(cx)
		if err != nil {
			t.Fatalf("funding-rate: %v", err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/public/funding-rate", params, false)
		assertCovers(t, "funding-rate", raw, resp)
	}

	// funding-rate-history (SWAP)
	{
		params := map[string]string{"instId": "BTC-USDT-SWAP"}
		resp, err := c.NewGetFundingRateHistoryService("BTC-USDT-SWAP").Do(cx)
		if err != nil {
			t.Fatalf("funding-rate-history: %v", err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/public/funding-rate-history", params, false)
		assertCovers(t, "funding-rate-history", raw, resp)
	}

	// open-interest (SWAP, BTC-USDT-SWAP)
	{
		params := map[string]string{"instType": string(InstTypeSwap), "instId": "BTC-USDT-SWAP"}
		resp, err := c.NewGetOpenInterestService(InstTypeSwap).SetInstId("BTC-USDT-SWAP").Do(cx)
		if err != nil {
			t.Fatalf("open-interest: %v", err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/public/open-interest", params, false)
		assertCovers(t, "open-interest", raw, resp)
	}

	// price-limit
	{
		params := map[string]string{"instId": "BTC-USDT-SWAP"}
		resp, err := c.NewGetPriceLimitService("BTC-USDT-SWAP").Do(cx)
		if err != nil {
			t.Fatalf("price-limit: %v", err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/public/price-limit", params, false)
		assertCovers(t, "price-limit", raw, resp)
	}

	// opt-summary (uly=BTC-USD)
	{
		params := map[string]string{"uly": "BTC-USD"}
		resp, err := c.NewGetOptSummaryService().SetUly("BTC-USD").Do(cx)
		if err != nil {
			t.Fatalf("opt-summary: %v", err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/public/opt-summary", params, false)
		assertCovers(t, "opt-summary", raw, resp)
	}

	// discount-rate-interest-free-quota (no params)
	{
		params := map[string]string{}
		resp, err := c.NewGetDiscountRateInterestFreeQuotaService().Do(cx)
		if err != nil {
			t.Fatalf("discount-rate-interest-free-quota: %v", err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/public/discount-rate-interest-free-quota", params, false)
		assertCovers(t, "discount-rate-interest-free-quota", raw, resp)
	}

	// mark-price (SWAP)
	{
		params := map[string]string{"instType": string(InstTypeSwap)}
		resp, err := c.NewGetMarkPriceService(InstTypeSwap).Do(cx)
		if err != nil {
			t.Fatalf("mark-price: %v", err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/public/mark-price", params, false)
		assertCovers(t, "mark-price", raw, resp)
	}

	// position-tiers (SWAP, cross, uly=BTC-USDT)
	{
		params := map[string]string{"instType": string(InstTypeSwap), "tdMode": string(TdModeCross), "uly": "BTC-USDT"}
		resp, err := c.NewGetPositionTiersService(InstTypeSwap, TdModeCross).SetUly("BTC-USDT").Do(cx)
		if err != nil {
			t.Fatalf("position-tiers: %v", err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/public/position-tiers", params, false)
		assertCovers(t, "position-tiers", raw, resp)
	}

	// interest-rate-loan-quota (no params)
	{
		resp, err := c.NewGetInterestRateLoanQuotaService().Do(cx)
		if err != nil {
			t.Fatalf("interest-rate-loan-quota: %v", err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/public/interest-rate-loan-quota", nil, false)
		assertCovers(t, "interest-rate-loan-quota", raw, resp)
	}

	// underlying (SWAP) — data is [][]string; assert non-empty.
	{
		resp, err := c.NewGetUnderlyingService(InstTypeSwap).Do(cx)
		if err != nil {
			t.Fatalf("underlying: %v", err)
		}
		if len(resp) == 0 || len(resp[0]) == 0 {
			t.Fatalf("underlying: empty result %v", resp)
		}
		t.Logf("underlying: OK, %d underlyings (e.g. %s)", len(resp[0]), resp[0][0])
	}

	// insurance-fund (SWAP, uly=BTC-USDT)
	{
		params := map[string]string{"instType": string(InstTypeSwap), "uly": "BTC-USDT"}
		resp, err := c.NewGetInsuranceFundService(InstTypeSwap).SetUly("BTC-USDT").Do(cx)
		if err != nil {
			t.Fatalf("insurance-fund: %v", err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/public/insurance-fund", params, false)
		assertCovers(t, "insurance-fund", raw, resp)
	}

	// convert-contract-coin (type=1, BTC-USD-SWAP, sz=1) — inverse contract
	// requires px, so supply one.
	{
		params := map[string]string{"type": "1", "instId": "BTC-USD-SWAP", "sz": "1", "px": "60000"}
		resp, err := c.NewGetConvertContractCoinService("1", "BTC-USD-SWAP", "1").SetPx("60000").Do(cx)
		if err != nil {
			t.Fatalf("convert-contract-coin: %v", err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/public/convert-contract-coin", params, false)
		assertCovers(t, "convert-contract-coin", raw, resp)
	}

	// instrument-tick-bands (OPTION)
	{
		params := map[string]string{"instType": string(InstTypeOption)}
		resp, err := c.NewGetInstrumentTickBandsService(InstTypeOption).Do(cx)
		if err != nil {
			t.Fatalf("instrument-tick-bands: %v", err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/public/instrument-tick-bands", params, false)
		assertCovers(t, "instrument-tick-bands", raw, resp)
	}

	// premium-history
	{
		params := map[string]string{"instId": "BTC-USDT-SWAP"}
		resp, err := c.NewGetPremiumHistoryService("BTC-USDT-SWAP").Do(cx)
		if err != nil {
			t.Fatalf("premium-history: %v", err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/public/premium-history", params, false)
		assertCovers(t, "premium-history", raw, resp)
	}

	// mm-instrument-types (SWAP)
	{
		params := map[string]string{"instType": string(InstTypeSwap)}
		resp, err := c.NewGetMMInstrumentTypesService().SetInstType(InstTypeSwap).Do(cx)
		if err != nil {
			t.Fatalf("mm-instrument-types: %v", err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/public/mm-instrument-types", params, false)
		assertCovers(t, "mm-instrument-types", raw, resp)
	}
}
