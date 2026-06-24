package okx

import "testing"

// TestStatus exercises the system-status (public) and the announcements /
// announcement-types / economic-calendar (signed) endpoints live, asserting that
// the typed structs cover every key the real responses return.
func TestStatus(t *testing.T) {
	cx := ctx(t)

	// --- GET /api/v5/system/status (public) ---
	pub := testPublicClient()
	{
		const label = "system/status"
		params := map[string]string{}
		resp, err := pub.NewGetSystemStatusService().Do(cx)
		if err != nil {
			t.Fatalf("%s: %v", label, err)
		}
		t.Logf("%s: %d maintenance window(s)", label, len(resp))
		if len(resp) > 0 {
			raw := fetchRawGet(t, pub, cx, "/api/v5/system/status", params, false)
			assertCovers(t, label, raw, resp)
		} else {
			t.Logf("%s: empty (no maintenance) — coverage check skipped", label)
		}
	}

	// --- signed endpoints ---
	c := testClient(t)
	_ = c.SyncServerTime(cx)

	// --- GET /api/v5/support/announcements (signed) ---
	{
		const label = "support/announcements"
		params := map[string]string{}
		resp, err := c.NewGetAnnouncementsService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "50105", "50103", "50111", "50030") {
				t.Fatalf("%s: %v", label, err)
			}
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/support/announcements", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/support/announcement-types (signed) ---
	{
		const label = "support/announcement-types"
		params := map[string]string{}
		resp, err := c.NewGetAnnouncementTypesService().Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "50105", "50103", "50111", "50030") {
				t.Fatalf("%s: %v", label, err)
			}
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/support/announcement-types", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/public/economic-calendar (signed) ---
	{
		const label = "public/economic-calendar"
		params := map[string]string{"limit": "5"}
		resp, err := c.NewGetEconomicCalendarService().SetLimit(5).Do(cx)
		if err != nil {
			if !tolerable(t, label, err, "50105", "50103", "50111", "50030") {
				t.Fatalf("%s: %v", label, err)
			}
		} else if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/public/economic-calendar", params, true)
			assertCovers(t, label, raw, resp)
		}
	}
}
