package okx

import (
	"testing"
)

// TestRubik exercises every public trading-statistics ("rubik") endpoint live.
// Most endpoints return positional string arrays (array-of-arrays), so
// assertCovers (which diffs object key sets) does not apply: instead we assert
// the call deserializes and, when non-empty, a parsed row has a non-zero Ts.
// support-coin is a real object, so its shape is checked with assertCovers.
func TestRubik(t *testing.T) {
	c := testPublicClient()
	cx := ctx(t)

	// support-coin (object {contract,option,spot}).
	{
		const path = "/api/v5/rubik/stat/trading-data/support-coin"
		resp, err := c.NewGetSupportCoinService().Do(cx)
		if err != nil {
			t.Fatalf("support-coin: %v", err)
		}
		if resp == nil || len(resp.Contract) == 0 || len(resp.Spot) == 0 {
			t.Fatalf("support-coin: empty lists: %+v", resp)
		}
		raw := fetchRawGet(t, c, cx, path, nil, false)
		assertCovers(t, "support-coin", raw, resp)
	}

	// taker-volume (ccy, instType): [ts, sellVol, buyVol].
	{
		const path = "/api/v5/rubik/stat/taker-volume"
		params := map[string]string{"ccy": "BTC", "instType": string(InstTypeSpot)}
		rows, err := c.NewGetTakerVolumeService("BTC", InstTypeSpot).Do(cx)
		if err != nil {
			t.Fatalf("taker-volume: %v", err)
		}
		if len(rows) > 0 && rows[0].Timestamp.IsZero() {
			t.Errorf("taker-volume: parsed row has zero Ts: %+v", rows[0])
		}
		_ = fetchRawGet(t, c, cx, path, params, false)
		t.Logf("taker-volume: %d rows", len(rows))
	}

	// margin/loan-ratio (ccy): [ts, ratio].
	{
		const path = "/api/v5/rubik/stat/margin/loan-ratio"
		params := map[string]string{"ccy": "BTC"}
		rows, err := c.NewGetLoanRatioService("BTC").Do(cx)
		if err != nil {
			t.Fatalf("loan-ratio: %v", err)
		}
		if len(rows) > 0 && rows[0].Timestamp.IsZero() {
			t.Errorf("loan-ratio: parsed row has zero Ts: %+v", rows[0])
		}
		_ = fetchRawGet(t, c, cx, path, params, false)
		t.Logf("loan-ratio: %d rows", len(rows))
	}

	// contracts/long-short-account-ratio (ccy): [ts, ratio].
	{
		const path = "/api/v5/rubik/stat/contracts/long-short-account-ratio"
		params := map[string]string{"ccy": "BTC"}
		rows, err := c.NewGetLongShortAccountRatioService("BTC").Do(cx)
		if err != nil {
			t.Fatalf("long-short-account-ratio: %v", err)
		}
		if len(rows) > 0 && rows[0].Timestamp.IsZero() {
			t.Errorf("long-short-account-ratio: parsed row has zero Ts: %+v", rows[0])
		}
		_ = fetchRawGet(t, c, cx, path, params, false)
		t.Logf("long-short-account-ratio: %d rows", len(rows))
	}

	// contracts/long-short-account-ratio-contract (instId): [ts, ratio].
	{
		const path = "/api/v5/rubik/stat/contracts/long-short-account-ratio-contract"
		params := map[string]string{"instId": "BTC-USDT-SWAP"}
		rows, err := c.NewGetLongShortAccountRatioContractService("BTC-USDT-SWAP").Do(cx)
		if err != nil {
			t.Fatalf("long-short-account-ratio-contract: %v", err)
		}
		if len(rows) > 0 && rows[0].Timestamp.IsZero() {
			t.Errorf("long-short-account-ratio-contract: parsed row has zero Ts: %+v", rows[0])
		}
		_ = fetchRawGet(t, c, cx, path, params, false)
		t.Logf("long-short-account-ratio-contract: %d rows", len(rows))
	}

	// contracts/long-short-account-ratio-contract-top-trader (instId): [ts, ratio].
	{
		const path = "/api/v5/rubik/stat/contracts/long-short-account-ratio-contract-top-trader"
		params := map[string]string{"instId": "BTC-USDT-SWAP"}
		rows, err := c.NewGetLongShortAccountRatioContractTopTraderService("BTC-USDT-SWAP").Do(cx)
		if err != nil {
			t.Fatalf("long-short-account-ratio-contract-top-trader: %v", err)
		}
		if len(rows) > 0 && rows[0].Timestamp.IsZero() {
			t.Errorf("long-short-account-ratio-contract-top-trader: parsed row has zero Ts: %+v", rows[0])
		}
		_ = fetchRawGet(t, c, cx, path, params, false)
		t.Logf("long-short-account-ratio-contract-top-trader: %d rows", len(rows))
	}

	// contracts/long-short-position-ratio-contract-top-trader (instId): [ts, ratio].
	{
		const path = "/api/v5/rubik/stat/contracts/long-short-position-ratio-contract-top-trader"
		params := map[string]string{"instId": "BTC-USDT-SWAP"}
		rows, err := c.NewGetLongShortPositionRatioContractTopTraderService("BTC-USDT-SWAP").Do(cx)
		if err != nil {
			t.Fatalf("long-short-position-ratio-contract-top-trader: %v", err)
		}
		if len(rows) > 0 && rows[0].Timestamp.IsZero() {
			t.Errorf("long-short-position-ratio-contract-top-trader: parsed row has zero Ts: %+v", rows[0])
		}
		_ = fetchRawGet(t, c, cx, path, params, false)
		t.Logf("long-short-position-ratio-contract-top-trader: %d rows", len(rows))
	}

	// contracts/open-interest-history (instId): [ts, oi, oiCcy, oiUsd].
	{
		const path = "/api/v5/rubik/stat/contracts/open-interest-history"
		params := map[string]string{"instId": "BTC-USDT-SWAP"}
		rows, err := c.NewGetOpenInterestHistoryService("BTC-USDT-SWAP").Do(cx)
		if err != nil {
			t.Fatalf("open-interest-history: %v", err)
		}
		if len(rows) > 0 && rows[0].Timestamp.IsZero() {
			t.Errorf("open-interest-history: parsed row has zero Ts: %+v", rows[0])
		}
		_ = fetchRawGet(t, c, cx, path, params, false)
		t.Logf("open-interest-history: %d rows", len(rows))
	}

	// contracts/open-interest-volume (ccy): [ts, oi, vol].
	{
		const path = "/api/v5/rubik/stat/contracts/open-interest-volume"
		params := map[string]string{"ccy": "BTC"}
		rows, err := c.NewGetContractsOpenInterestVolumeService("BTC").Do(cx)
		if err != nil {
			t.Fatalf("contracts open-interest-volume: %v", err)
		}
		if len(rows) > 0 && rows[0].Timestamp.IsZero() {
			t.Errorf("contracts open-interest-volume: parsed row has zero Ts: %+v", rows[0])
		}
		_ = fetchRawGet(t, c, cx, path, params, false)
		t.Logf("contracts open-interest-volume: %d rows", len(rows))
	}

	// option/open-interest-volume (ccy): [ts, oi, vol].
	{
		const path = "/api/v5/rubik/stat/option/open-interest-volume"
		params := map[string]string{"ccy": "BTC"}
		rows, err := c.NewGetOptionOpenInterestVolumeService("BTC").Do(cx)
		if err != nil {
			t.Fatalf("option open-interest-volume: %v", err)
		}
		if len(rows) > 0 && rows[0].Timestamp.IsZero() {
			t.Errorf("option open-interest-volume: parsed row has zero Ts: %+v", rows[0])
		}
		_ = fetchRawGet(t, c, cx, path, params, false)
		t.Logf("option open-interest-volume: %d rows", len(rows))
	}

	// option/open-interest-volume-expiry (ccy): [ts, expTime, callOI, putOI, callVol, putVol].
	var nearExpiry string
	{
		const path = "/api/v5/rubik/stat/option/open-interest-volume-expiry"
		params := map[string]string{"ccy": "BTC"}
		rows, err := c.NewGetOptionOpenInterestVolumeExpiryService("BTC").Do(cx)
		if err != nil {
			t.Fatalf("option open-interest-volume-expiry: %v", err)
		}
		if len(rows) > 0 {
			if rows[0].Timestamp.IsZero() {
				t.Errorf("option open-interest-volume-expiry: parsed row has zero Ts: %+v", rows[0])
			}
			nearExpiry = rows[0].ExpiryTime
		}
		_ = fetchRawGet(t, c, cx, path, params, false)
		t.Logf("option open-interest-volume-expiry: %d rows (nearExpiry=%s)", len(rows), nearExpiry)
	}

	// option/open-interest-volume-strike (ccy, expTime): [ts, strike, callOI, putOI, callVol, putVol].
	{
		const path = "/api/v5/rubik/stat/option/open-interest-volume-strike"
		expTime := nearExpiry
		if expTime == "" {
			expTime = "20260626"
		}
		params := map[string]string{"ccy": "BTC", "expTime": expTime}
		rows, err := c.NewGetOptionOpenInterestVolumeStrikeService("BTC", expTime).Do(cx)
		if err != nil {
			t.Fatalf("option open-interest-volume-strike: %v", err)
		}
		if len(rows) > 0 && rows[0].Timestamp.IsZero() {
			t.Errorf("option open-interest-volume-strike: parsed row has zero Ts: %+v", rows[0])
		}
		_ = fetchRawGet(t, c, cx, path, params, false)
		t.Logf("option open-interest-volume-strike: %d rows", len(rows))
	}

	// option/taker-block-volume (ccy): single flat array snapshot.
	{
		const path = "/api/v5/rubik/stat/option/taker-block-volume"
		params := map[string]string{"ccy": "BTC"}
		resp, err := c.NewGetOptionTakerBlockVolumeService("BTC").Do(cx)
		if err != nil {
			t.Fatalf("option taker-block-volume: %v", err)
		}
		if resp != nil && resp.Timestamp.IsZero() {
			t.Errorf("option taker-block-volume: parsed snapshot has zero Ts: %+v", resp)
		}
		_ = fetchRawGet(t, c, cx, path, params, false)
		t.Logf("option taker-block-volume: %+v", resp)
	}
}
