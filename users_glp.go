package okx

import (
	"context"
	"strconv"

	"github.com/UnipayFI/go-okx/request"
	"github.com/shopspring/decimal"
)

// GLPProgram selects a GLP (Global Liquidity Program) business line for the
// historical-performance endpoint.
type GLPProgram string

const (
	GLPProgramSpot   GLPProgram = "SPOT"
	GLPProgramPerp   GLPProgram = "PERP"
	GLPProgramFutNto GLPProgram = "FUT_NTO"
)

// GetGLPTodayPerformanceService -- GET /api/v5/users/glp/today-performance (Read)
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
	req := request.Get(ctx, s.c, "/api/v5/users/glp/today-performance", s.params).WithSign()
	return request.DoOne[GLPTodayPerformance](req)
}

// GetGLPHistoricalPerformanceService -- GET /api/v5/users/glp/historical-performance (Read)
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
	req := request.Get(ctx, s.c, "/api/v5/users/glp/historical-performance", s.params).WithSign()
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
	QualifyingRows        []GLPQualifyingRow `json:"qualifyingRows"`
	Daily                 GLPPeriodPerf      `json:"daily"`
	MonthToDate           GLPPeriodPerf      `json:"mtd"`
}

// GLPQualifyingRow is one row of a program's qualifying-pool breakdown. Modeled
// from the OKX docs; unknown keys decode away.
type GLPQualifyingRow struct {
	PairType string          `json:"pairType"`
	Volume   decimal.Decimal `json:"volume"`
	Share    decimal.Decimal `json:"share"`
}

// GLPPeriodPerf holds a period's volume and share breakdowns.
type GLPPeriodPerf struct {
	Volume GLPCategoryBreakdown `json:"volume"`
	Share  GLPCategoryBreakdown `json:"share"`
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
