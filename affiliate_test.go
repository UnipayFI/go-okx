package okx

import "testing"

// TestAffiliate exercises the signed affiliate READ endpoint live, asserting
// that the typed struct covers every key the real response returns.
//
// Only /api/v5/affiliate/invitee/detail is part of the public OKX v5 REST API;
// it is curl-verified to return code 51620 ("Only affiliates can perform this
// action") for a non-affiliate account (path REAL, capability-gated). The other
// documented-sounding affiliate paths (performance-summary, invitee-list,
// co-inviter-link-list, link-list, sub-affiliate-list) all return HTTP 404 and
// were dropped. The validating account is not an affiliate, so the call is
// expected to be tolerated.
func TestAffiliate(t *testing.T) {
	c := testClient(t)
	_ = c.SyncServerTime(ctx(t))
	cx := ctx(t)

	// --- GET /api/v5/affiliate/invitee/detail (Read) ---
	{
		const label = "affiliate/invitee/detail"
		params := map[string]string{"uid": "123456"}
		resp, err := c.NewGetAffiliateInviteeDetailService("123456").Do(cx)
		if err != nil {
			// 51620: not an affiliate; 59509/18004: affiliate role/permission gating;
			// 50014: param probe (placeholder uid not a real invitee).
			if tolerable(t, label, err, "51620", "59509", "18004", "50014") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if resp == nil {
			t.Logf("%s: empty data (no such invitee) — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/affiliate/invitee/detail", params, true)
			assertCovers(t, label, raw, resp)
		}
	}
}
