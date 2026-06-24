package okx

import "testing"

// TestAccountBills exercises the account bills read endpoints live, asserting
// that the typed structs cover every key the real responses return. The
// bills-history-archive apply (POST) is implemented but never executed (it
// creates a downloadable export on a real account); only the safe GET reads are
// tested here.
func TestAccountBills(t *testing.T) {
	c := testClient(t)
	_ = c.SyncServerTime(ctx(t))
	cx := ctx(t)

	// --- GET /api/v5/account/bills (last 7 days) ---
	{
		const label = "account/bills"
		params := map[string]string{}
		resp, err := c.NewGetBillsService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "51010", "50014", "51000") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/account/bills", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/account/bills-archive (last 3 months) ---
	{
		const label = "account/bills-archive"
		params := map[string]string{}
		resp, err := c.NewGetBillsArchiveService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "51010", "50014", "51000") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/account/bills-archive", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/account/bills-history-archive (year=2024 quarter=Q1) ---
	// Requires a prior POST apply (code 51604) before the link exists; tolerate
	// that and the missing/invalid-param codes.
	{
		const label = "account/bills-history-archive"
		params := map[string]string{"year": "2024", "quarter": "Q1"}
		resp, err := c.NewGetBillsHistoryArchiveService("2024", BillsHistoryArchiveQuarterQ1).Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "51604", "51000", "50014") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if resp == nil {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/account/bills-history-archive", params, true)
			assertCovers(t, label, raw, resp)
		}
	}
}
