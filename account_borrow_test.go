package okx

import "testing"

// TestAccountBorrow exercises the interest / borrow-repay / VIP-loan read
// endpoints under /api/v5/account/* live, asserting that the typed structs cover
// every key the real responses return. The state-changing POST endpoints
// (spot-manual-borrow-repay, set-auto-repay, borrow-repay) are implemented but
// intentionally NOT exercised against the real account.
func TestAccountBorrow(t *testing.T) {
	c := testClient(t)
	_ = c.SyncServerTime(ctx(t))
	cx := ctx(t)

	// --- GET /api/v5/account/interest-accrued ---
	{
		const label = "account/interest-accrued"
		params := map[string]string{}
		resp, err := c.NewGetInterestAccruedService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "51010", "50014", "51000") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/account/interest-accrued", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/account/interest-rate ---
	{
		const label = "account/interest-rate"
		params := map[string]string{}
		resp, err := c.NewGetInterestRateService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "51010", "50014") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/account/interest-rate", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/account/interest-limits ---
	{
		const label = "account/interest-limits"
		params := map[string]string{}
		resp, err := c.NewGetInterestLimitsService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "59307", "51010") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if resp == nil {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/account/interest-limits", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/account/spot-borrow-repay-history ---
	{
		const label = "account/spot-borrow-repay-history"
		params := map[string]string{}
		resp, err := c.NewGetSpotBorrowRepayHistoryService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "51010", "50014") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/account/spot-borrow-repay-history", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/account/borrow-repay-history ---
	{
		const label = "account/borrow-repay-history"
		params := map[string]string{}
		resp, err := c.NewGetBorrowRepayHistoryService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "51010", "50014") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/account/borrow-repay-history", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/account/vip-loan-order-list ---
	{
		const label = "account/vip-loan-order-list"
		params := map[string]string{}
		resp, err := c.NewGetVipLoanOrderListService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "51010", "50014") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/account/vip-loan-order-list", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/account/vip-loan-order-detail (requires ordId) ---
	{
		const label = "account/vip-loan-order-detail"
		params := map[string]string{"ordId": "0"}
		resp, err := c.NewGetVipLoanOrderDetailService("0").Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "51000", "51010", "50014") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data (no such order) — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/account/vip-loan-order-detail", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/account/vip-interest-accrued ---
	{
		const label = "account/vip-interest-accrued"
		params := map[string]string{}
		resp, err := c.NewGetVipInterestAccruedService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "51010", "50014") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/account/vip-interest-accrued", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/account/vip-interest-deducted ---
	{
		const label = "account/vip-interest-deducted"
		params := map[string]string{}
		resp, err := c.NewGetVipInterestDeductedService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "51010", "50014") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/account/vip-interest-deducted", params, true)
			assertCovers(t, label, raw, resp)
		}
	}
}
