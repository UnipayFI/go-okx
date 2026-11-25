package okx

import (
	"testing"
	"time"

	"github.com/UnipayFI/go-okx/client"
	"github.com/UnipayFI/go-okx/common"
)

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

// TestGLPV2 exercises the signed GLP v2 READ endpoints (overview, details,
// history/performance) live, asserting that the typed structs cover every key
// the real responses return. A non-enrolled account tolerates the
// permission/param error codes.
func TestGLPV2(t *testing.T) {
	c := testClient(t)
	_ = c.SyncServerTime(ctx(t))
	cx := ctx(t)

	// --- GET /api/v5/users/glp/v2/overview (Read) ---
	{
		const label = "users/glp/v2/overview"
		params := map[string]string{}
		resp, err := c.NewGetGLPV2OverviewService().Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50030", "51035", "50014", "51010") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if resp == nil || len(resp.Programs) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/users/glp/v2/overview", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/users/glp/v2/details (Read) ---
	{
		const label = "users/glp/v2/details"
		params := map[string]string{"bizId": string(GLPBusinessIDPerpMaker)}
		resp, err := c.NewGetGLPV2DetailsService(GLPBusinessIDPerpMaker).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50030", "51035", "50014", "51010") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if resp == nil || len(resp.Performance) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/users/glp/v2/details", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/users/glp/v2/history/performance (Read) ---
	{
		const label = "users/glp/v2/history/performance"
		params := map[string]string{"bizId": string(GLPBusinessIDPerpMaker)}
		resp, err := c.NewGetGLPV2PerformanceHistoryService(GLPBusinessIDPerpMaker).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50030", "51035", "50014", "51010") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if resp == nil || len(resp.Performance) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/users/glp/v2/history/performance", params, true)
			assertCovers(t, label, raw, resp)
		}
	}
}

// TestGLPC exercises the signed GLP-C READ endpoints (registered-pairs,
// performance, history/performance) live, asserting that the typed structs
// cover every key the real responses return. POST /users/glp/c/register-pairs
// is state-changing and is not executed against a live account. A non-enrolled
// account tolerates the permission/param error codes.
func TestGLPC(t *testing.T) {
	c := testClient(t)
	_ = c.SyncServerTime(ctx(t))
	cx := ctx(t)

	// --- GET /api/v5/users/glp/c/registered-pairs (Read) ---
	{
		const label = "users/glp/c/registered-pairs"
		params := map[string]string{}
		resp, err := c.NewGetGLPCRegisteredPairsService().Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50030", "51035", "50014", "51010") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/users/glp/c/registered-pairs", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/users/glp/c/performance (Read) ---
	{
		const label = "users/glp/c/performance"
		params := map[string]string{"limit": "10"}
		resp, err := c.NewGetGLPCPerformanceService().SetLimit(10).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50030", "51035", "50014", "51010") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/users/glp/c/performance", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/users/glp/c/history/performance (Read) ---
	{
		const label = "users/glp/c/history/performance"
		params := map[string]string{"limit": "10"}
		resp, err := c.NewGetGLPCPerformanceHistoryService().SetLimit(10).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50030", "51035", "50014", "51010") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/users/glp/c/history/performance", params, true)
			assertCovers(t, label, raw, resp)
		}
	}
}

// TestGLPCJSON verifies JSON encoding/decoding and key coverage of the GLP-C
// response structs against OKX's documented payloads.
func TestGLPCJSON(t *testing.T) {
	const registerRaw = `[{"instId":"DOGE-USDT","sCode":"0","sMsg":""},{"instId":"BAD-USDT","sCode":"51001","sMsg":"Instrument ID does not exist"}]`
	var reg []GLPCRegisterPairResult
	if err := common.JSONUnmarshal([]byte(registerRaw), &reg); err != nil {
		t.Fatalf("unmarshal register-pairs: %v", err)
	}
	assertCovers(t, "users/glp/c/register-pairs (fixture)", []byte(registerRaw), reg)
	if len(reg) != 2 || reg[0].InstrumentID != "DOGE-USDT" || reg[0].Err() != nil {
		t.Fatalf("reg[0] = %+v, err = %v", reg[0], reg[0].Err())
	}
	if !client.IsCode(reg[1].Err(), "51001") {
		t.Fatalf("reg[1].Err() = %v, want code 51001", reg[1].Err())
	}

	const pairsRaw = `[{"instId":"DOGE-USDT","instType":"SPOT","registerTime":"1732300800000"},{"instId":"PEPE-USDT-SWAP","instType":"SWAP","registerTime":"1732300800000"}]`
	var pairs []GLPCRegisteredPair
	if err := common.JSONUnmarshal([]byte(pairsRaw), &pairs); err != nil {
		t.Fatalf("unmarshal registered-pairs: %v", err)
	}
	assertCovers(t, "users/glp/c/registered-pairs (fixture)", []byte(pairsRaw), pairs)
	if len(pairs) != 2 || pairs[0].InstrumentID != "DOGE-USDT" || pairs[0].InstrumentType != InstTypeSpot {
		t.Fatalf("pairs = %+v", pairs)
	}
	if want := time.UnixMilli(1732300800000).UTC(); !pairs[0].RegisterTime.Equal(want) {
		t.Errorf("registerTime = %v, want %v", pairs[0].RegisterTime, want)
	}

	const perfRaw = `[{"instId":"DOGE-USDT","instType":"SPOT","ts":"1732363200000","qualified":"true","bidDepth":"50000","askDepth":"48000","reqDepth":"30000"}]`
	var perf []GLPCPerformance
	if err := common.JSONUnmarshal([]byte(perfRaw), &perf); err != nil {
		t.Fatalf("unmarshal performance: %v", err)
	}
	assertCovers(t, "users/glp/c/performance (fixture)", []byte(perfRaw), perf)
	if len(perf) != 1 || perf[0].Qualified != "true" || perf[0].BidDepth.String() != "50000" || perf[0].AskDepth.String() != "48000" || perf[0].RequiredDepth.String() != "30000" {
		t.Fatalf("perf = %+v", perf)
	}
	if want := time.UnixMilli(1732363200000).UTC(); !perf[0].Timestamp.Equal(want) {
		t.Errorf("ts = %v, want %v", perf[0].Timestamp, want)
	}

	const histRaw = `[{"instId":"DOGE-USDT","instType":"SPOT","assessStartTime":"1732276800000","assessEndTime":"1732363200000","qualifiedHours":"20","totalHours":"24","qualifiedRate":"0.8333"}]`
	var hist []GLPCPerformanceHistory
	if err := common.JSONUnmarshal([]byte(histRaw), &hist); err != nil {
		t.Fatalf("unmarshal history/performance: %v", err)
	}
	assertCovers(t, "users/glp/c/history/performance (fixture)", []byte(histRaw), hist)
	if len(hist) != 1 || hist[0].QualifiedHours.String() != "20" || hist[0].TotalHours.String() != "24" || hist[0].QualifiedRate.String() != "0.8333" {
		t.Fatalf("hist = %+v", hist)
	}
	if wantStart := time.UnixMilli(1732276800000).UTC(); !hist[0].AssessmentStartTime.Equal(wantStart) {
		t.Errorf("assessStartTime = %v, want %v", hist[0].AssessmentStartTime, wantStart)
	}
	if wantEnd := time.UnixMilli(1732363200000).UTC(); !hist[0].AssessmentEndTime.Equal(wantEnd) {
		t.Errorf("assessEndTime = %v, want %v", hist[0].AssessmentEndTime, wantEnd)
	}
}
