package okx

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-okx/request"
	"github.com/shopspring/decimal"
)

// GLPProgram selects a GLP (Global Liquidity Program) business line for the
// historicalperformance endpoint.
type GLPProgram string

const (
	GLPProgramSpot   GLPProgram = "SPOT"
	GLPProgramPerp   GLPProgram = "PERP"
	GLPProgramFutNto GLPProgram = "FUT_NTO"
)

// GetGLPTodayPerformanceService -- GET /api/v5/users/glp/todayperformance (Read)
//
// Returns the enrolled GLP (Global Liquidity Program) account's current-day and
// month-to-date assessment snapshots across all programs. Only available to
// accounts enrolled in a GLP market-maker program; the account is resolved from
// the API key, so there are no request parameters.
type GetGLPTodayPerformanceService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetGLPTodayPerformanceService() *GetGLPTodayPerformanceService {
	return &GetGLPTodayPerformanceService{c: c, params: map[string]string{}}
}

func (s *GetGLPTodayPerformanceService) Do(ctx context.Context) (*GLPTodayPerformance, error) {
	req := request.Get(ctx, s.c, "/api/v5/users/glp/todayperformance", s.params).WithSign()
	return request.DoOne[GLPTodayPerformance](req)
}

// GetGLPHistoricalPerformanceService -- GET /api/v5/users/glp/historicalperformance (Read)
//
// Returns the enrolled GLP account's daily performance history for one program,
// sorted newest-first.
type GetGLPHistoricalPerformanceService struct {
	c      *Client
	params map[string]string
}

// NewGetGLPHistoricalPerformanceService requires the GLP program (SPOT/PERP/FUT_NTO).
func (c *Client) NewGetGLPHistoricalPerformanceService(program GLPProgram) *GetGLPHistoricalPerformanceService {
	return &GetGLPHistoricalPerformanceService{c: c, params: map[string]string{"program": string(program)}}
}

// SetBegin filters to records at or after the given ms timestamp.
func (s *GetGLPHistoricalPerformanceService) SetBegin(ms int64) *GetGLPHistoricalPerformanceService {
	s.params["begin"] = strconv.FormatInt(ms, 10)
	return s
}

// SetEnd filters to records at or before the given ms timestamp.
func (s *GetGLPHistoricalPerformanceService) SetEnd(ms int64) *GetGLPHistoricalPerformanceService {
	s.params["end"] = strconv.FormatInt(ms, 10)
	return s
}

// SetLimit caps the number of daily records returned (default 31, max 100).
func (s *GetGLPHistoricalPerformanceService) SetLimit(limit int) *GetGLPHistoricalPerformanceService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetGLPHistoricalPerformanceService) Do(ctx context.Context) ([]GLPHistoricalPerformance, error) {
	req := request.Get(ctx, s.c, "/api/v5/users/glp/historicalperformance", s.params).WithSign()
	return request.DoList[GLPHistoricalPerformance](req)
}

// GLPTodayPerformance is the today / month-to-date GLP snapshot. The validating
// account is not enrolled in any GLP program, so the field set is modeled from
// the OKX doc field tables.
type GLPTodayPerformance struct {
	DataReady bool             `json:"dataReady"`
	DataDate  string           `json:"dataDate"`
	Account   GLPAccount       `json:"account"`
	Programs  []GLPProgramPerf `json:"programs"`
}

// GLPAccount identifies the master account and any combined accounts assessed
// together for GLP.
type GLPAccount struct {
	MasterAccountID    string   `json:"masterAccountId"`
	CombinedAccountIDs []string `json:"combinedAccountIds"`
}

// GLPProgramPerf is one enrolled program's daily and month-to-date performance.
type GLPProgramPerf struct {
	Program               string             `json:"program"`
	MarketMakerBusinessID string             `json:"marketMakerBusinessId"`
	EnrollmentStatus      string             `json:"enrollmentStatus"`
	MarketMakerLevelID    string             `json:"marketMakerLevelId"`
	EnrolledTierDisplay   string             `json:"enrolledTierDisplay"`
	QualifyingPool        string             `json:"qualifyingPool"`
	QualifyingRows        []string           `json:"qualifyingRows"` // row keys to highlight, e.g. "TOTAL"
	Daily                 GLPPeriodPerf      `json:"daily"`
	MonthToDate           GLPMonthToDatePerf `json:"mtd"`
}

// GLPPeriodPerf holds a period's volume and share breakdowns.
type GLPPeriodPerf struct {
	Volume GLPCategoryBreakdown `json:"volume"`
	Share  GLPCategoryBreakdown `json:"share"`
}

// GLPMonthToDatePerf is the month-to-date breakdown plus the MTD qualification
// outcome. MonthToDateStatus is QUALIFIED, UPGRADE or DOWNGRADE.
type GLPMonthToDatePerf struct {
	Volume            GLPCategoryBreakdown `json:"volume"`
	Share             GLPCategoryBreakdown `json:"share"`
	MonthToDateStatus string               `json:"mtdStatus"`
	QualifyingShare   GLPMakerTaker        `json:"qualifyingShare"`
}

// GLPCategoryBreakdown splits a metric (volume or share) by GLP pair-type
// category, each carrying a maker/taker pair.
type GLPCategoryBreakdown struct {
	TypeA      GLPMakerTaker `json:"typeA"`
	TypeBTotal GLPMakerTaker `json:"typeBTotal"`
	TypeBAdj   GLPMakerTaker `json:"typeBAdj"`
	TradFiX2   GLPMakerTaker `json:"tradfiX2"`
	Total      GLPMakerTaker `json:"total"`
}

// GLPMakerTaker is a maker/taker value pair (decimal-encoded strings).
type GLPMakerTaker struct {
	Maker decimal.Decimal `json:"maker"`
	Taker decimal.Decimal `json:"taker"`
}

// GLPHistoricalPerformance is one day of a program's GLP performance history.
type GLPHistoricalPerformance struct {
	Date   string               `json:"date"`
	Volume GLPCategoryBreakdown `json:"volume"`
	Share  GLPCategoryBreakdown `json:"share"`
}

// GLPBusinessID selects a GLP v2 product and trading side (bizId).
type GLPBusinessID string

const (
	GLPBusinessIDSpotMaker   GLPBusinessID = "14"
	GLPBusinessIDFutNtoMaker GLPBusinessID = "15"
	GLPBusinessIDPerpMaker   GLPBusinessID = "16"
	GLPBusinessIDSpotTaker   GLPBusinessID = "17"
	GLPBusinessIDFutNtoTaker GLPBusinessID = "18"
	GLPBusinessIDPerpTaker   GLPBusinessID = "19"
)

// GetGLPV2OverviewService -- GET /api/v5/users/glp/v2/overview (Read)
//
// Returns the account's active GLP enrollments with current and estimated
// tiers, tier requirements and week-to-date performance. A sub-account API key
// returns the master account's GLP data.
type GetGLPV2OverviewService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetGLPV2OverviewService() *GetGLPV2OverviewService {
	return &GetGLPV2OverviewService{c: c, params: map[string]string{}}
}

func (s *GetGLPV2OverviewService) Do(ctx context.Context) (*GLPV2Overview, error) {
	req := request.Get(ctx, s.c, "/api/v5/users/glp/v2/overview", s.params).WithSign()
	return request.DoOne[GLPV2Overview](req)
}

// GetGLPV2DetailsService -- GET /api/v5/users/glp/v2/details (Read)
//
// Returns daily and week-to-date performance for one GLP product and trading
// side within the current assessment period.
type GetGLPV2DetailsService struct {
	c      *Client
	params map[string]string
}

// NewGetGLPV2DetailsService requires the GLP product and trading side (bizId).
func (c *Client) NewGetGLPV2DetailsService(bizID GLPBusinessID) *GetGLPV2DetailsService {
	return &GetGLPV2DetailsService{c: c, params: map[string]string{"bizId": string(bizID)}}
}

func (s *GetGLPV2DetailsService) Do(ctx context.Context) (*GLPV2Details, error) {
	req := request.Get(ctx, s.c, "/api/v5/users/glp/v2/details", s.params).WithSign()
	return request.DoOne[GLPV2Details](req)
}

// GetGLPV2PerformanceHistoryService -- GET /api/v5/users/glp/v2/history/performance (Read)
//
// Returns weekly assessment results and tier outcomes for one GLP product and
// trading side.
type GetGLPV2PerformanceHistoryService struct {
	c      *Client
	params map[string]string
}

// NewGetGLPV2PerformanceHistoryService requires the GLP product and trading side (bizId).
func (c *Client) NewGetGLPV2PerformanceHistoryService(bizID GLPBusinessID) *GetGLPV2PerformanceHistoryService {
	return &GetGLPV2PerformanceHistoryService{c: c, params: map[string]string{"bizId": string(bizID)}}
}

func (s *GetGLPV2PerformanceHistoryService) Do(ctx context.Context) (*GLPV2PerformanceHistory, error) {
	req := request.Get(ctx, s.c, "/api/v5/users/glp/v2/history/performance", s.params).WithSign()
	return request.DoOne[GLPV2PerformanceHistory](req)
}

// GLPV2Overview lists the account's active GLP enrollments (empty when none).
// The validating account is not enrolled in any GLP program, so the GLP v2
// structs are modeled from the OKX doc field tables.
type GLPV2Overview struct {
	Programs []GLPV2Program `json:"programs"`
}

// GLPV2Program is one active GLP enrollment. Program is GLP_A or GLP_B, Side is
// MAKER or TAKER, and Instrument is SPOT, PERP or FUT_NTO. CurrentTier,
// EstimatedNextTier and EstimatedNextStatus (disqualify/downgrade/qualify/
// upgrade) are empty when unavailable.
type GLPV2Program struct {
	BusinessID          GLPBusinessID          `json:"bizId"`
	Program             string                 `json:"program"`
	Side                string                 `json:"side"`
	Instrument          string                 `json:"instrument"`
	CurrentTier         string                 `json:"currentTier"`
	EstimatedNextTier   string                 `json:"estNextTier"`
	EstimatedNextStatus string                 `json:"estNextStatus"`
	GracePeriodEndTime  time.Time              `json:"gracePeriodEndTime,format:unixmilli"` // earliest date a downgrade can take effect; zero when not applicable
	Requirement         []GLPV2TierRequirement `json:"requirement"`
	Performance         []GLPV2WeekToDate      `json:"performance"`
}

// GLPV2TierRequirement is the market-share threshold for one GLP tier, as a
// decimal fraction ("0.007" = 0.7%).
type GLPV2TierRequirement struct {
	Tier  string          `json:"tier"`
	Share decimal.Decimal `json:"share"`
}

// GLPV2WeekToDate is the week-to-date share and volume as of one trading date.
// AssessmentPeriod is PREVIOUS or CURRENT.
type GLPV2WeekToDate struct {
	Date             time.Time       `json:"date,format:unixmilli"`
	AssessmentPeriod string          `json:"assessPeriod"`
	Share            decimal.Decimal `json:"share"`
	Volume           decimal.Decimal `json:"vol"`
}

// GLPV2Details is the current assessment period's performance for one GLP
// product and trading side.
type GLPV2Details struct {
	AssessmentStartTime time.Time        `json:"assessStartTime,format:unixmilli"`
	AssessmentEndTime   time.Time        `json:"assessEndTime,format:unixmilli"`
	Performance         []GLPV2DailyPerf `json:"performance"`
}

// GLPV2DailyPerf is one trading date's volume and market share (decimal
// fractions), both for the day and accumulated week-to-date, before weighting
// (Original) and adjusted.
type GLPV2DailyPerf struct {
	Date                      time.Time       `json:"date,format:unixmilli"`
	OriginalVolume            decimal.Decimal `json:"origVol"`
	OriginalShare             decimal.Decimal `json:"origShare"`
	AdjustedVolume            decimal.Decimal `json:"adjVol"`
	AdjustedShare             decimal.Decimal `json:"adjShare"`
	AccumulatedOriginalVolume decimal.Decimal `json:"accOrigVol"`
	AccumulatedOriginalShare  decimal.Decimal `json:"accOrigShare"`
	AccumulatedAdjustedVolume decimal.Decimal `json:"accAdjVol"`
	AccumulatedAdjustedShare  decimal.Decimal `json:"accAdjShare"`
}

// GLPV2PerformanceHistory lists weekly GLP assessment results (empty when
// there is no history).
type GLPV2PerformanceHistory struct {
	Performance []GLPV2Assessment `json:"performance"`
}

// GLPV2Assessment is one weekly assessment. Tier is the tier assessed against,
// NextTier the tier it determined (empty when unavailable), and
// AssessmentResult is disqualify, downgrade, qualify or upgrade (empty when
// unavailable).
type GLPV2Assessment struct {
	AssessmentStartTime       time.Time       `json:"assessStartTime,format:unixmilli"`
	AssessmentEndTime         time.Time       `json:"assessEndTime,format:unixmilli"`
	Tier                      string          `json:"tier"`
	NextTier                  string          `json:"nextTier"`
	AccumulatedOriginalVolume decimal.Decimal `json:"accOrigVol"`
	AccumulatedOriginalShare  decimal.Decimal `json:"accOrigShare"`
	AccumulatedAdjustedVolume decimal.Decimal `json:"accAdjVol"`
	AccumulatedAdjustedShare  decimal.Decimal `json:"accAdjShare"`
	AssessmentResult          string          `json:"assessResult"`
}

// GLPCInstType selects the instrument category for GLP-C (Global Liquidity
// Program - C) endpoints. Requests use SPOT or PERP (including X-Perp); the
// hourly assessment summary returns SWAP for perpetuals.
type GLPCInstType string

const (
	GLPCInstTypeSpot GLPCInstType = "SPOT"
	GLPCInstTypePerp GLPCInstType = "PERP"
	GLPCInstTypeSwap GLPCInstType = "SWAP"
)

// GLPCFilterType filters eligible GLP-C instruments by enrollment state.
type GLPCFilterType int

const (
	GLPCFilterAll        GLPCFilterType = 1 // all eligible instruments
	GLPCFilterRegistered GLPCFilterType = 2 // registered
	GLPCFilterOpen       GLPCFilterType = 3 // open for registration
)

// GetGLPCInstrumentsService -- GET /api/v5/users/glp-c/instruments (Read)
//
// Retrieves GLP-C instruments eligible for the current account, with
// registration counts and qualification requirements. Only instruments tradable
// by the current account are returned. Note: OKX returns "data" as a JSON
// object rather than an array, so the service decodes via DoObject.
// Rate limit: 5 requests per second per User ID.
type GetGLPCInstrumentsService struct {
	c      *Client
	params map[string]string
}

// NewGetGLPCInstrumentsService requires the instrument category (SPOT or PERP),
// 1-based page number, page size (max 100), and enrollment filter type.
func (c *Client) NewGetGLPCInstrumentsService(instType GLPCInstType, page, pageSize int, filterType GLPCFilterType) *GetGLPCInstrumentsService {
	return &GetGLPCInstrumentsService{
		c: c,
		params: map[string]string{
			"instType":   string(instType),
			"page":       strconv.Itoa(page),
			"pageSize":   strconv.Itoa(pageSize),
			"filterType": strconv.Itoa(int(filterType)),
		},
	}
}

// SetInstId filters by a single instrument ID (e.g. "BTC-USDT").
func (s *GetGLPCInstrumentsService) SetInstId(instId string) *GetGLPCInstrumentsService {
	s.params["instId"] = instId
	return s
}

func (s *GetGLPCInstrumentsService) Do(ctx context.Context) (*GLPCInstrumentsPage, error) {
	req := request.Get(ctx, s.c, "/api/v5/users/glp-c/instruments", s.params).WithSign()
	return request.DoObject[GLPCInstrumentsPage](req)
}

// GLPCInstrumentsPage is one paginated page of eligible GLP-C instruments.
// The validating account is not whitelisted for GLP-C, so the GLP-C structs are
// modeled from the OKX doc field tables.
type GLPCInstrumentsPage struct {
	Page        int              `json:"page"`
	PageSize    int              `json:"pageSize"`
	Total       int              `json:"total"`
	StartIndex  int              `json:"startIndex"`
	HasNextPage bool             `json:"hasNextPage"`
	Rows        []GLPCInstrument `json:"rows"`
}

// GLPCInstrument is one GLP-C eligible instrument and its quoting requirements.
// InstrumentGroup is G1, G2 or G3; RequiredMarketShare and RequiredUptime are
// decimal fractions (e.g. "0.05" = 5%, "0.7" = 70%; "0.5" during public beta).
type GLPCInstrument struct {
	InstrumentType        GLPCInstType    `json:"instType"`
	InstrumentID          string          `json:"instId"`
	InstrumentGroup       string          `json:"instGroup"`
	RegisteredClientCount int             `json:"registeredClientCount"`
	DailyAvgPlatTV7d      decimal.Decimal `json:"dailyAvgPlatTv7d"`    // 7-day daily average platform trading volume
	B3PriceLevel          decimal.Decimal `json:"b3PriceLevel"`        // 7-day avg spread between 3rd-best ask and bid, in bp
	B3PriceQuoteSize      decimal.Decimal `json:"b3PriceQuoteSize"`    // 7-day avg combined bid+ask quote size across levels 1-3, in USD
	RequiredMarketShare   decimal.Decimal `json:"requiredMarketShare"` // required minute market share
	RequiredUptime        decimal.Decimal `json:"requiredUptime"`      // required hourly uptime
}

// GLPCInstrumentRef identifies one instrument by category + instrument ID in
// GLP-C enrollment requests and responses.
type GLPCInstrumentRef struct {
	InstrumentType GLPCInstType `json:"instType"`
	InstrumentID   string       `json:"instId"`
}

// EnrollGLPCInstrumentsService -- POST /api/v5/users/glp-c/enroll-instruments (Trade)
//
// Enrolls one or more instruments in GLP-C. Enrollment applies to the whole
// participant and is assessed from the current minute. Enrolling an already
// enrolled instrument returns it in EnrolledInstruments; withdrawal is not
// supported in this phase. IMPLEMENT-ONLY: mutates market-maker enrollments.
// Rate limit: 5 requests per second per User ID.
type EnrollGLPCInstrumentsService struct {
	c    *Client
	body map[string]any
}

// NewEnrollGLPCInstrumentsService requires the list of instruments to enroll.
func (c *Client) NewEnrollGLPCInstrumentsService(insts []GLPCInstrumentRef) *EnrollGLPCInstrumentsService {
	return &EnrollGLPCInstrumentsService{c: c, body: map[string]any{"instsToEnroll": insts}}
}

func (s *EnrollGLPCInstrumentsService) Do(ctx context.Context) (*GLPCEnrollResult, error) {
	req := request.Post(ctx, s.c, "/api/v5/users/glp-c/enroll-instruments", s.body).WithSign()
	return request.DoObject[GLPCEnrollResult](req)
}

// GLPCEnrollResult splits the enrollment attempt into successfully enrolled and
// failed (invalid or ineligible) instruments.
type GLPCEnrollResult struct {
	EnrolledInstruments []GLPCInstrumentRef `json:"enrolledInsts"`
	FailedInstruments   []GLPCInstrumentRef `json:"failedInsts"`
}

// GetGLPCEnrollStatsService -- GET /api/v5/users/glp-c/enroll/stats (Read)
//
// Retrieves the total, registered and open-for-registration GLP-C instrument
// counts. Note: OKX returns "data" as a JSON object rather than an array, so
// the service decodes via DoObject. Rate limit: 5 requests per second per User ID.
type GetGLPCEnrollStatsService struct {
	c      *Client
	params map[string]string
}

// NewGetGLPCEnrollStatsService requires the instrument category (SPOT or PERP).
func (c *Client) NewGetGLPCEnrollStatsService(instType GLPCInstType) *GetGLPCEnrollStatsService {
	return &GetGLPCEnrollStatsService{c: c, params: map[string]string{"instType": string(instType)}}
}

func (s *GetGLPCEnrollStatsService) Do(ctx context.Context) (*GLPCEnrollStats, error) {
	req := request.Get(ctx, s.c, "/api/v5/users/glp-c/enroll/stats", s.params).WithSign()
	return request.DoObject[GLPCEnrollStats](req)
}

// GLPCEnrollStats holds GLP-C instrument counts. AvailableInstrumentCount is
// not necessarily equal to TotalInstrumentCount - RegisteredInstrumentCount.
type GLPCEnrollStats struct {
	TotalInstrumentCount      int `json:"totalInstCount"`
	RegisteredInstrumentCount int `json:"registeredInstCount"`
	AvailableInstrumentCount  int `json:"availableInstCount"`
}

// GetGLPCMinuteAssessmentService -- GET /api/v5/users/glp-c/quote-assessment/minute (Read)
//
// Retrieves minute-level GLP-C quote assessment results of one instrument,
// sorted ascending by minuteTs. Either SetInstId or SetInstIdCode is required
// (if both are passed, instIdCode takes precedence). Assessment data is
// available for the last 72 hours. Rate limit: 5 requests per second per User ID.
type GetGLPCMinuteAssessmentService struct {
	c      *Client
	params map[string]string
}

// NewGetGLPCMinuteAssessmentService requires the inclusive begin and end
// timestamps (within the last 72 hours, end not later than request time).
func (c *Client) NewGetGLPCMinuteAssessmentService(beginTs, endTs time.Time) *GetGLPCMinuteAssessmentService {
	return &GetGLPCMinuteAssessmentService{
		c: c,
		params: map[string]string{
			"beginTs": strconv.FormatInt(beginTs.UnixMilli(), 10),
			"endTs":   strconv.FormatInt(endTs.UnixMilli(), 10),
		},
	}
}

// SetInstType filters by instrument category (SPOT or PERP, including X-Perp).
func (s *GetGLPCMinuteAssessmentService) SetInstType(instType GLPCInstType) *GetGLPCMinuteAssessmentService {
	s.params["instType"] = string(instType)
	return s
}

// SetInstId selects the instrument by ID (e.g. "BTC-USDT-SWAP"). Either
// SetInstId or SetInstIdCode is required.
func (s *GetGLPCMinuteAssessmentService) SetInstId(instId string) *GetGLPCMinuteAssessmentService {
	s.params["instId"] = instId
	return s
}

// SetInstIdCode selects the instrument by ID code (e.g. "10459"). Takes
// precedence if both instId and instIdCode are passed.
func (s *GetGLPCMinuteAssessmentService) SetInstIdCode(instIdCode string) *GetGLPCMinuteAssessmentService {
	s.params["instIdCode"] = instIdCode
	return s
}

func (s *GetGLPCMinuteAssessmentService) Do(ctx context.Context) ([]GLPCMinuteAssessment, error) {
	req := request.Get(ctx, s.c, "/api/v5/users/glp-c/quote-assessment/minute", s.params).WithSign()
	return request.DoList[GLPCMinuteAssessment](req)
}

// GLPCMinuteAssessment is one minute-level GLP-C quote assessment record.
// MinuteStatus is QUALIFIED or NOT_QUALIFIED. ReasonCode is one of QUALIFIED,
// BELOW_THRESHOLD, NO_CUSTOMER_QUOTE, MARKET_DENOMINATOR_ZERO,
// SOURCE_DATA_MISSING or SOURCE_DATA_INVALID.
type GLPCMinuteAssessment struct {
	InstrumentIDCode                   string                 `json:"instIdCode"`
	InstrumentType                     GLPCInstType           `json:"instType"`
	InstrumentID                       string                 `json:"instId"`
	MinuteTimestamp                    time.Time              `json:"minuteTs,format:unixmilli"`
	BilateralWeightedShare             decimal.Decimal        `json:"bilateralWeightedShare"`
	Threshold                          decimal.Decimal        `json:"threshold"`
	MinuteStatus                       string                 `json:"minuteStatus"`
	ReasonCode                         string                 `json:"reasonCode"`
	ParticipantBidOriginalQuoteSizeUSD decimal.Decimal        `json:"participantBidOriginalQuoteSizeUsd"`
	ParticipantAskOriginalQuoteSizeUSD decimal.Decimal        `json:"participantAskOriginalQuoteSizeUsd"`
	ParticipantBidWeightedQuoteSizeUSD decimal.Decimal        `json:"participantBidWeightedQuoteSizeUsd"`
	ParticipantAskWeightedQuoteSizeUSD decimal.Decimal        `json:"participantAskWeightedQuoteSizeUsd"`
	Levels                             []GLPCMinuteLevelQuote `json:"levels"`
}

// GLPCMinuteLevelQuote is the per-price-level bid/ask quote breakdown in USD.
type GLPCMinuteLevelQuote struct {
	Level                string          `json:"level"`
	Multiplier           decimal.Decimal `json:"multiplier"`
	BidQuoteSizeUSD      decimal.Decimal `json:"bidQuoteSizeUsd"`
	AskQuoteSizeUSD      decimal.Decimal `json:"askQuoteSizeUsd"`
	BidMultipliedSizeUSD decimal.Decimal `json:"bidMultipliedSizeUsd"`
	AskMultipliedSizeUSD decimal.Decimal `json:"askMultipliedSizeUsd"`
}

// GetGLPCHourlyInstrumentAssessmentService -- GET /api/v5/users/glp-c/quote-assessment/hour/instrument (Read)
//
// Retrieves hourly GLP-C qualification results of one instrument for ended
// hours (the current hour is not returned). Either SetInstId or SetInstIdCode
// is required (if both are passed, instIdCode takes precedence).
// Rate limit: 5 requests per second per User ID.
type GetGLPCHourlyInstrumentAssessmentService struct {
	c      *Client
	params map[string]string
}

// NewGetGLPCHourlyInstrumentAssessmentService requires the inclusive begin and
// end timestamps (within the last 72 hours, end not later than request time).
func (c *Client) NewGetGLPCHourlyInstrumentAssessmentService(beginTs, endTs time.Time) *GetGLPCHourlyInstrumentAssessmentService {
	return &GetGLPCHourlyInstrumentAssessmentService{
		c: c,
		params: map[string]string{
			"beginTs": strconv.FormatInt(beginTs.UnixMilli(), 10),
			"endTs":   strconv.FormatInt(endTs.UnixMilli(), 10),
		},
	}
}

// SetInstType filters by instrument category (SPOT or PERP, including X-Perp).
func (s *GetGLPCHourlyInstrumentAssessmentService) SetInstType(instType GLPCInstType) *GetGLPCHourlyInstrumentAssessmentService {
	s.params["instType"] = string(instType)
	return s
}

// SetInstId selects the instrument by ID (e.g. "BTC-USDT-SWAP"). Either
// SetInstId or SetInstIdCode is required.
func (s *GetGLPCHourlyInstrumentAssessmentService) SetInstId(instId string) *GetGLPCHourlyInstrumentAssessmentService {
	s.params["instId"] = instId
	return s
}

// SetInstIdCode selects the instrument by ID code (e.g. "10459"). Takes
// precedence if both instId and instIdCode are passed.
func (s *GetGLPCHourlyInstrumentAssessmentService) SetInstIdCode(instIdCode string) *GetGLPCHourlyInstrumentAssessmentService {
	s.params["instIdCode"] = instIdCode
	return s
}

func (s *GetGLPCHourlyInstrumentAssessmentService) Do(ctx context.Context) ([]GLPCHourlyInstrumentAssessment, error) {
	req := request.Get(ctx, s.c, "/api/v5/users/glp-c/quote-assessment/hour/instrument", s.params).WithSign()
	return request.DoList[GLPCHourlyInstrumentAssessment](req)
}

// GLPCHourlyInstrumentAssessment is one ended hour's GLP-C qualification result
// for an instrument. HourStatus is QUALIFIED, NOT_QUALIFIED or DELAYED (ended
// but not finalized, in which case FinalizedAt is zero).
type GLPCHourlyInstrumentAssessment struct {
	InstrumentIDCode      string                 `json:"instIdCode"`
	InstrumentType        GLPCInstType           `json:"instType"`
	InstrumentID          string                 `json:"instId"`
	AssessmentHour        time.Time              `json:"assessmentHour,format:unixmilli"`
	InstrumentGroup       string                 `json:"instGroup"`
	Threshold             decimal.Decimal        `json:"threshold"`
	RequiredHourlyUptime  decimal.Decimal        `json:"requiredHourlyUptime"`
	QualifiedMinutes      decimal.Decimal        `json:"qualifiedMinutes"`
	NotQualifiedMinutes   decimal.Decimal        `json:"notQualifiedMinutes"`
	HourlyUptime          decimal.Decimal        `json:"hourlyUptime"`
	HourStatus            string                 `json:"hourStatus"`
	IntendedEffectiveHour time.Time              `json:"intendedEffectiveHour,format:unixmilli"`
	FinalizedAt           time.Time              `json:"finalizedAt,format:unixmilli"`
	MinuteStatusDetail    []GLPCMinuteStatusItem `json:"minuteStatusDetail"`
}

// GLPCMinuteStatusItem is one minute's qualification status within an hourly
// instrument assessment.
type GLPCMinuteStatusItem struct {
	MinuteTimestamp time.Time `json:"minuteTs,format:unixmilli"`
	MinuteStatus    string    `json:"minuteStatus"`
	ReasonCode      string    `json:"reasonCode"`
}

// GetGLPCHourlyAssessmentSummaryService -- GET /api/v5/users/glp-c/quote-assessment/hour/summary (Read)
//
// Retrieves hourly qualified instrument counts and the calculated GLP-C tier,
// sorted ascending by assessmentHour (SPOT before SWAP in the same hour). The
// current hour can be included. Rate limit: 5 requests per second per User ID.
type GetGLPCHourlyAssessmentSummaryService struct {
	c      *Client
	params map[string]string
}

// NewGetGLPCHourlyAssessmentSummaryService requires the inclusive begin and end
// timestamps (within the last 72 hours, end not later than request time).
func (c *Client) NewGetGLPCHourlyAssessmentSummaryService(beginTs, endTs time.Time) *GetGLPCHourlyAssessmentSummaryService {
	return &GetGLPCHourlyAssessmentSummaryService{
		c: c,
		params: map[string]string{
			"beginTs": strconv.FormatInt(beginTs.UnixMilli(), 10),
			"endTs":   strconv.FormatInt(endTs.UnixMilli(), 10),
		},
	}
}

// SetInstType filters by instrument category (SPOT or PERP, including X-Perp).
// Returns both categories if omitted.
func (s *GetGLPCHourlyAssessmentSummaryService) SetInstType(instType GLPCInstType) *GetGLPCHourlyAssessmentSummaryService {
	s.params["instType"] = string(instType)
	return s
}

func (s *GetGLPCHourlyAssessmentSummaryService) Do(ctx context.Context) ([]GLPCHourlyAssessmentSummary, error) {
	req := request.Get(ctx, s.c, "/api/v5/users/glp-c/quote-assessment/hour/summary", s.params).WithSign()
	return request.DoList[GLPCHourlyAssessmentSummary](req)
}

// GLPCHourlyAssessmentSummary is one hourly GLP-C tier summary record.
// InstrumentType is SPOT or SWAP (including X-Perp). State is in_progress,
// finalized or delayed. Tier is "1" (>=40, VIP 9 Maker), "2" (30-39, VIP 8),
// "3" (20-29, VIP 7), "4" (10-19, VIP 6), or "" if fewer than 10 or unknown.
type GLPCHourlyAssessmentSummary struct {
	AssessmentHour            time.Time       `json:"assessmentHour,format:unixmilli"`
	InstrumentType            GLPCInstType    `json:"instType"`
	State                     string          `json:"state"`
	RegisteredInstrumentCount decimal.Decimal `json:"registeredInstCount"`
	QualifiedInstrumentCount  decimal.Decimal `json:"qualifiedInstCount"`
	QualifiedInstrumentRatio  decimal.Decimal `json:"qualifiedInstRatio"`
	Tier                      string          `json:"tier"`
	IntendedEffectiveHour     time.Time       `json:"intendedEffectiveHour,format:unixmilli"`
}
