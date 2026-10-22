package okx

import "testing"

// TestGLPPerformance exercises the signed GLP (Global Liquidity Program) READ
// endpoints (todayperformance, historicalperformance) live, asserting that the
// typed structs cover every key the real responses return. These are market-maker
// only; a non-enrolled account tolerates the permission/param error codes.
func TestGLPPerformance(t *testing.T) {
	c := testClient(t)
	_ = c.SyncServerTime(ctx(t))
	cx := ctx(t)

	// --- GET /api/v5/users/glp/todayperformance (Read) ---
	{
		const label = "users/glp/todayperformance"
		params := map[string]string{}
		resp, err := c.NewGetGLPTodayPerformanceService().Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50030", "51035", "50014", "51010") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if resp == nil {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/users/glp/todayperformance", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/users/glp/historicalperformance (Read) ---
	{
		const label = "users/glp/historicalperformance"
		params := map[string]string{"program": string(GLPProgramSpot)}
		resp, err := c.NewGetGLPHistoricalPerformanceService(GLPProgramSpot).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50030", "51035", "50014", "51010") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/users/glp/historicalperformance", params, true)
			assertCovers(t, label, raw, resp)
		}
	}
}
