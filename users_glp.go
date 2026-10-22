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
