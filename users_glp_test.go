package okx

import (
	"strconv"
	"testing"
	"time"

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

// TestGLPC exercises the signed GLP-C READ endpoints (instruments,
// enroll/stats, quote-assessment/minute, quote-assessment/hour/instrument,
// quote-assessment/hour/summary) live, asserting that the typed structs cover
// every key the real responses return. POST /users/glp-c/enroll-instruments is
// state-changing and is not executed against a live account. A non-whitelisted
// account tolerates the permission/param error codes.
func TestGLPC(t *testing.T) {
	c := testClient(t)
	_ = c.SyncServerTime(ctx(t))
	cx := ctx(t)

	end := time.Now().UTC().Truncate(time.Hour)
	begin := end.Add(-2 * time.Hour)
	beginMs := strconv.FormatInt(begin.UnixMilli(), 10)
	endMs := strconv.FormatInt(end.UnixMilli(), 10)

	// --- GET /api/v5/users/glp-c/instruments (Read) ---
	{
		const label = "users/glp-c/instruments"
		params := map[string]string{
			"instType":   string(GLPCInstTypeSpot),
			"page":       "1",
			"pageSize":   "10",
			"filterType": strconv.Itoa(int(GLPCFilterAll)),
		}
		resp, err := c.NewGetGLPCInstrumentsService(GLPCInstTypeSpot, 1, 10, GLPCFilterAll).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50030", "51035", "50014", "51010") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if resp == nil {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/users/glp-c/instruments", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/users/glp-c/enroll/stats (Read) ---
	{
		const label = "users/glp-c/enroll/stats"
		params := map[string]string{"instType": string(GLPCInstTypeSpot)}
		resp, err := c.NewGetGLPCEnrollStatsService(GLPCInstTypeSpot).Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50030", "51035", "50014", "51010") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if resp == nil {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/users/glp-c/enroll/stats", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/users/glp-c/quote-assessment/minute (Read) ---
	{
		const label = "users/glp-c/quote-assessment/minute"
		params := map[string]string{
			"instType": string(GLPCInstTypePerp),
			"instId":   "BTC-USDT-SWAP",
			"beginTs":  beginMs,
			"endTs":    endMs,
		}
		resp, err := c.NewGetGLPCMinuteAssessmentService(begin, end).
			SetInstType(GLPCInstTypePerp).
			SetInstId("BTC-USDT-SWAP").
			Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50030", "51035", "50014", "51010") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/users/glp-c/quote-assessment/minute", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/users/glp-c/quote-assessment/hour/instrument (Read) ---
	{
		const label = "users/glp-c/quote-assessment/hour/instrument"
		params := map[string]string{
			"instType": string(GLPCInstTypePerp),
			"instId":   "BTC-USDT-SWAP",
			"beginTs":  beginMs,
			"endTs":    endMs,
		}
		resp, err := c.NewGetGLPCHourlyInstrumentAssessmentService(begin, end).
			SetInstType(GLPCInstTypePerp).
			SetInstId("BTC-USDT-SWAP").
			Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50030", "51035", "50014", "51010") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/users/glp-c/quote-assessment/hour/instrument", params, true)
			assertCovers(t, label, raw, resp)
		}
	}

	// --- GET /api/v5/users/glp-c/quote-assessment/hour/summary (Read) ---
	{
		const label = "users/glp-c/quote-assessment/hour/summary"
		params := map[string]string{
			"instType": string(GLPCInstTypePerp),
			"beginTs":  beginMs,
			"endTs":    endMs,
		}
		resp, err := c.NewGetGLPCHourlyAssessmentSummaryService(begin, end).
			SetInstType(GLPCInstTypePerp).
			Do(cx)
		if err != nil {
			if tolerable(t, label, err, "50030", "51035", "50014", "51010") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		if len(resp) == 0 {
			t.Logf("%s: empty data — coverage check skipped", label)
		} else {
			raw := fetchRawGet(t, c, cx, "/api/v5/users/glp-c/quote-assessment/hour/summary", params, true)
			assertCovers(t, label, raw, resp)
		}
	}
}

// TestGLPCJSON verifies JSON encoding/decoding and key coverage of all six
// GLP-C response structs against OKX's documented payloads.
func TestGLPCJSON(t *testing.T) {
	const instsRaw = `{"page":1,"pageSize":10,"total":1,"startIndex":0,"hasNextPage":false,"rows":[{"instType":"SPOT","instId":"BTC-USDT","instGroup":"G1","registeredClientCount":0,"dailyAvgPlatTv7d":"0","b3PriceLevel":"5","b3PriceQuoteSize":"25","requiredMarketShare":"0.05","requiredUptime":"0.7"}]}`
	var page GLPCInstrumentsPage
	if err := common.JSONUnmarshal([]byte(instsRaw), &page); err != nil {
		t.Fatalf("unmarshal instruments: %v", err)
	}
	assertCovers(t, "users/glp-c/instruments (fixture)", []byte(instsRaw), &page)
	if page.Page != 1 || page.PageSize != 10 || page.Total != 1 || page.HasNextPage || len(page.Rows) != 1 {
		t.Fatalf("page = %+v", page)
	}
	if r := page.Rows[0]; r.InstrumentType != GLPCInstTypeSpot || r.InstrumentID != "BTC-USDT" || r.InstrumentGroup != "G1" || r.RequiredMarketShare.String() != "0.05" || r.RequiredUptime.String() != "0.7" {
		t.Fatalf("row = %+v", r)
	}

	const enrollRaw = `{"enrolledInsts":[{"instType":"SPOT","instId":"BTC-USDT"}],"failedInsts":[{"instType":"PERP","instId":"BTC-USDT"}]}`
	var enroll GLPCEnrollResult
	if err := common.JSONUnmarshal([]byte(enrollRaw), &enroll); err != nil {
		t.Fatalf("unmarshal enroll-instruments: %v", err)
	}
	assertCovers(t, "users/glp-c/enroll-instruments (fixture)", []byte(enrollRaw), &enroll)
	if len(enroll.EnrolledInstruments) != 1 || enroll.EnrolledInstruments[0].InstrumentType != GLPCInstTypeSpot || len(enroll.FailedInstruments) != 1 || enroll.FailedInstruments[0].InstrumentType != GLPCInstTypePerp {
		t.Fatalf("enroll = %+v", enroll)
	}

	const statsRaw = `{"totalInstCount":500,"registeredInstCount":40,"availableInstCount":460}`
	var stats GLPCEnrollStats
	if err := common.JSONUnmarshal([]byte(statsRaw), &stats); err != nil {
		t.Fatalf("unmarshal enroll/stats: %v", err)
	}
	assertCovers(t, "users/glp-c/enroll/stats (fixture)", []byte(statsRaw), &stats)
	if stats.TotalInstrumentCount != 500 || stats.RegisteredInstrumentCount != 40 || stats.AvailableInstrumentCount != 460 {
		t.Fatalf("stats = %+v", stats)
	}

	const minuteRaw = `[{"instIdCode":"10459","instType":"PERP","instId":"BTC-USDT-SWAP","minuteTs":"1795154460000","bilateralWeightedShare":"0.0600000000","threshold":"0.05","minuteStatus":"QUALIFIED","reasonCode":"QUALIFIED","participantBidOriginalQuoteSizeUsd":"2300","participantAskOriginalQuoteSizeUsd":"2300","participantBidWeightedQuoteSizeUsd":"3050","participantAskWeightedQuoteSizeUsd":"3250","levels":[{"level":"1","multiplier":"2","bidQuoteSizeUsd":"1000","askQuoteSizeUsd":"1200","bidMultipliedSizeUsd":"2000","askMultipliedSizeUsd":"2400"},{"level":"2","multiplier":"1","bidQuoteSizeUsd":"800","askQuoteSizeUsd":"600","bidMultipliedSizeUsd":"800","askMultipliedSizeUsd":"600"},{"level":"3","multiplier":"0.5","bidQuoteSizeUsd":"500","askQuoteSizeUsd":"500","bidMultipliedSizeUsd":"250","askMultipliedSizeUsd":"250"}]}]`
	var minutes []GLPCMinuteAssessment
	if err := common.JSONUnmarshal([]byte(minuteRaw), &minutes); err != nil {
		t.Fatalf("unmarshal quote-assessment/minute: %v", err)
	}
	assertCovers(t, "users/glp-c/quote-assessment/minute (fixture)", []byte(minuteRaw), minutes)
	if len(minutes) != 1 || minutes[0].InstrumentIDCode != "10459" || minutes[0].InstrumentType != GLPCInstTypePerp || len(minutes[0].Levels) != 3 {
		t.Fatalf("minutes = %+v", minutes)
	}
	if want := time.UnixMilli(1795154460000).UTC(); !minutes[0].MinuteTimestamp.Equal(want) {
		t.Errorf("minuteTs = %v, want %v", minutes[0].MinuteTimestamp, want)
	}
	if minutes[0].BilateralWeightedShare.String() != "0.06" || minutes[0].Levels[2].Multiplier.String() != "0.5" {
		t.Fatalf("minutes[0] values = %+v", minutes[0])
	}

	const hourInstRaw = `[{"instIdCode":"10459","instType":"PERP","instId":"BTC-USDT-SWAP","assessmentHour":"1795154400000","instGroup":"G1","threshold":"0.05","requiredHourlyUptime":"0.7","qualifiedMinutes":"1","notQualifiedMinutes":"1","hourlyUptime":"0.0166666667","hourStatus":"DELAYED","intendedEffectiveHour":"1795158000000","finalizedAt":"","minuteStatusDetail":[{"minuteTs":"1795154400000","minuteStatus":"NOT_QUALIFIED","reasonCode":"BELOW_THRESHOLD"},{"minuteTs":"1795154460000","minuteStatus":"QUALIFIED","reasonCode":"QUALIFIED"}]}]`
	var hourInst []GLPCHourlyInstrumentAssessment
	if err := common.JSONUnmarshal([]byte(hourInstRaw), &hourInst); err != nil {
		t.Fatalf("unmarshal quote-assessment/hour/instrument: %v", err)
	}
	assertCovers(t, "users/glp-c/quote-assessment/hour/instrument (fixture)", []byte(hourInstRaw), hourInst)
	if len(hourInst) != 1 || hourInst[0].HourStatus != "DELAYED" || !hourInst[0].FinalizedAt.IsZero() || len(hourInst[0].MinuteStatusDetail) != 2 {
		t.Fatalf("hourInst = %+v", hourInst)
	}
	if want := time.UnixMilli(1795154400000).UTC(); !hourInst[0].AssessmentHour.Equal(want) {
		t.Errorf("assessmentHour = %v, want %v", hourInst[0].AssessmentHour, want)
	}
	if want := time.UnixMilli(1795158000000).UTC(); !hourInst[0].IntendedEffectiveHour.Equal(want) {
		t.Errorf("intendedEffectiveHour = %v, want %v", hourInst[0].IntendedEffectiveHour, want)
	}
	if want := time.UnixMilli(1795154460000).UTC(); !hourInst[0].MinuteStatusDetail[1].MinuteTimestamp.Equal(want) {
		t.Errorf("minuteTs[1] = %v, want %v", hourInst[0].MinuteStatusDetail[1].MinuteTimestamp, want)
	}

	const summaryRaw = `[{"assessmentHour":"1795150800000","instType":"SWAP","state":"finalized","registeredInstCount":"47","qualifiedInstCount":"23","qualifiedInstRatio":"0.4893617021","tier":"3","intendedEffectiveHour":"1795154400000"},{"assessmentHour":"1795158000000","instType":"SWAP","state":"in_progress","registeredInstCount":"47","qualifiedInstCount":"0","qualifiedInstRatio":"0.0000000000","tier":"","intendedEffectiveHour":"1795161600000"}]`
	var summary []GLPCHourlyAssessmentSummary
	if err := common.JSONUnmarshal([]byte(summaryRaw), &summary); err != nil {
		t.Fatalf("unmarshal quote-assessment/hour/summary: %v", err)
	}
	assertCovers(t, "users/glp-c/quote-assessment/hour/summary (fixture)", []byte(summaryRaw), summary)
	if len(summary) != 2 || summary[0].InstrumentType != GLPCInstTypeSwap || summary[0].Tier != "3" || summary[1].State != "in_progress" || summary[1].Tier != "" {
		t.Fatalf("summary = %+v", summary)
	}
	if summary[0].RegisteredInstrumentCount.String() != "47" || summary[0].QualifiedInstrumentCount.String() != "23" || summary[0].QualifiedInstrumentRatio.String() != "0.4893617021" {
		t.Fatalf("summary[0] = %+v", summary[0])
	}
	if want := time.UnixMilli(1795150800000).UTC(); !summary[0].AssessmentHour.Equal(want) {
		t.Errorf("assessmentHour = %v, want %v", summary[0].AssessmentHour, want)
	}
}
