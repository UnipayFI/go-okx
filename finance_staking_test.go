package okx

import (
	"testing"
)

// TestFinanceStaking exercises the on-chain earn and ETH/SOL staking read
// endpoints (/api/v5/finance/staking-defi/*) live, asserting that the typed
// structs cover every key the real responses return. All endpoints are private
// (signed). State-changing purchase/redeem/cancel endpoints are implement-only
// and are not exercised here.
func TestFinanceStaking(t *testing.T) {
	c := testClient(t)
	_ = c.SyncServerTime(ctx(t))
	cx := ctx(t)

	// --- GET /api/v5/finance/staking-defi/offers ---
	{
		const label = "staking-defi/offers"
		params := map[string]string{}
		resp, err := c.NewGetStakingOffersService().Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50030", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/finance/staking-defi/offers", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/finance/staking-defi/orders-active ---
	{
		const label = "staking-defi/orders-active"
		params := map[string]string{}
		resp, err := c.NewGetStakingActiveOrdersService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "50030", "50011") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: no active orders — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/finance/staking-defi/orders-active", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/finance/staking-defi/orders-history ---
	{
		const label = "staking-defi/orders-history"
		params := map[string]string{}
		resp, err := c.NewGetStakingOrdersHistoryService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "50030", "50011") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty history — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/finance/staking-defi/orders-history", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/finance/staking-defi/eth/product-info ---
	{
		const label = "staking-defi/eth/product-info"
		params := map[string]string{}
		resp, err := c.NewGetEthStakingProductInfoService().Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50030", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/finance/staking-defi/eth/product-info", params, true)
		assertCovers(t, label, raw, resp)
	}

	// --- GET /api/v5/finance/staking-defi/eth/balance ---
	{
		const label = "staking-defi/eth/balance"
		params := map[string]string{}
		resp, err := c.NewGetEthStakingBalanceService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "50030", "50011") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/finance/staking-defi/eth/balance", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/finance/staking-defi/eth/purchase-redeem-history ---
	{
		const label = "staking-defi/eth/purchase-redeem-history"
		params := map[string]string{}
		resp, err := c.NewGetEthStakingHistoryService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "50030", "50011") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty history — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/finance/staking-defi/eth/purchase-redeem-history", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/finance/staking-defi/eth/apy-history ---
	{
		const label = "staking-defi/eth/apy-history"
		params := map[string]string{"days": "7"}
		resp, err := c.NewGetEthStakingApyHistoryService(7).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50030", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/finance/staking-defi/eth/apy-history", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/finance/staking-defi/sol/product-info ---
	// OKX returns this endpoint's "data" as a single JSON object (not an array),
	// so the service decodes via DoObject.
	{
		const label = "staking-defi/sol/product-info"
		params := map[string]string{}
		resp, err := c.NewGetSolStakingProductInfoService().Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50030", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := fetchRawGet(t, c, cx, "/api/v5/finance/staking-defi/sol/product-info", params, true)
		assertCovers(t, label, raw, resp)
	}

	// --- GET /api/v5/finance/staking-defi/sol/balance ---
	{
		const label = "staking-defi/sol/balance"
		params := map[string]string{}
		resp, err := c.NewGetSolStakingBalanceService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "50030", "50011") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/finance/staking-defi/sol/balance", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/finance/staking-defi/sol/purchase-redeem-history ---
	{
		const label = "staking-defi/sol/purchase-redeem-history"
		params := map[string]string{}
		resp, err := c.NewGetSolStakingHistoryService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "50030", "50011") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty history — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/finance/staking-defi/sol/purchase-redeem-history", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/finance/staking-defi/sol/apy-history ---
	{
		const label = "staking-defi/sol/apy-history"
		params := map[string]string{"days": "7"}
		resp, err := c.NewGetSolStakingApyHistoryService(7).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50030", "50011") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/finance/staking-defi/sol/apy-history", params, true)
			assertCovers(t, label, raw, resp)
		}
	}
}
