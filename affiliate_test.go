package okx

import "testing"

// TestAffiliate exercises the signed affiliate READ endpoints live, asserting
// that the typed structs cover every key the real responses return.
//
// Both paths are curl-verified to exist (a request without credentials returns
// 50103 rather than HTTP 404) and are affiliate-role gated: a non-affiliate
// account gets code 51620 ("Only affiliates can perform this action"). The
// validating account is not an affiliate, so the calls are expected to be
// tolerated.
func TestAffiliate(t *testing.T) {
	c := testClient(t)
	_ = c.SyncServerTime(ctx(t))
	cx := ctx(t)

	// --- GET /api/v5/affiliate/invitee/detail (Read) ---
	{
		const label = "affiliate/invitee/detail"
		params := map[string]string{"uid": "123456", "periodType": "last_30d"}
		resp, err := c.NewGetAffiliateInviteeDetailService("123456").
			SetPeriodType(AffiliatePeriodLast30D).
			Do(cx)
		if err != nil {
			// 51620: not an affiliate; 59509/18004: affiliate role/permission gating;
			// 50014: param probe (placeholder uid not a real invitee).
			if !tolerable(t, label, err, "51620", "59509", "18004", "50014") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if resp == nil {
			t.Logf("%s: empty data (no such invitee) — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/affiliate/invitee/detail", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/affiliate/invitee/list (Read) ---
	{
		const label = "affiliate/invitee/list"
		params := map[string]string{"limit": "1"}
		resp, err := c.NewGetAffiliateInviteeListService().SetLimit(1).Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "51620", "59509", "18004", "50014") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data (no invitees) — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/affiliate/invitee/list", params, true)
			assertCovers(t, label, raw, resp)
		}
	}
}
