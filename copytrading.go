package okx

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-okx/request"
	"github.com/shopspring/decimal"
)

// CopySubPosType is a copy-trading sub-position side selector.
type CopySubPosType string

const (
	CopySubPosTypeLong  CopySubPosType = "long"
	CopySubPosTypeShort CopySubPosType = "short"
)

// CopySortType is the lead-trader ranking sort key (overview/pnl/aum/winRatio/...).
type CopySortType string

const (
	CopySortTypeOverview CopySortType = "overview"
	CopySortTypePnl      CopySortType = "pnl"
	CopySortTypePnlRatio CopySortType = "pnlRatio"
	CopySortTypeAum      CopySortType = "aum"
	CopySortTypeWinRatio CopySortType = "winRatio"
)

// CopyState is the copy/lead state flag ("0" not copying / "1" copying).
type CopyState string

const (
	CopyStateNo  CopyState = "0"
	CopyStateYes CopyState = "1"
)

// GetCopyCurrentSubpositionsService -- GET /api/v5/copytrading/current-subpositions (Read)
//
// Returns the lead trader's currently open sub-positions (the positions copiers
// are mirroring).
type GetCopyCurrentSubpositionsService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetCopyCurrentSubpositionsService() *GetCopyCurrentSubpositionsService {
	return &GetCopyCurrentSubpositionsService{c: c, params: map[string]string{}}
}

// SetInstType filters by product line (SWAP).
func (s *GetCopyCurrentSubpositionsService) SetInstType(instType InstType) *GetCopyCurrentSubpositionsService {
	s.params["instType"] = string(instType)
	return s
}

// SetInstId filters by a single instrument id.
func (s *GetCopyCurrentSubpositionsService) SetInstId(instId string) *GetCopyCurrentSubpositionsService {
	s.params["instId"] = instId
	return s
}

// SetSubPosType filters by sub-position side (long/short).
func (s *GetCopyCurrentSubpositionsService) SetSubPosType(subPosType CopySubPosType) *GetCopyCurrentSubpositionsService {
	s.params["subPosType"] = string(subPosType)
	return s
}

// SetAfter paginates to records with subPosId earlier than the given id.
func (s *GetCopyCurrentSubpositionsService) SetAfter(after string) *GetCopyCurrentSubpositionsService {
	s.params["after"] = after
	return s
}

// SetBefore paginates to records with subPosId newer than the given id.
func (s *GetCopyCurrentSubpositionsService) SetBefore(before string) *GetCopyCurrentSubpositionsService {
	s.params["before"] = before
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetCopyCurrentSubpositionsService) SetLimit(limit int) *GetCopyCurrentSubpositionsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetCopyCurrentSubpositionsService) Do(ctx context.Context) ([]CopySubPosition, error) {
	req := request.Get(ctx, s.c, "/api/v5/copytrading/current-subpositions", s.params).WithSign()
	return request.DoList[CopySubPosition](req)
}

// CopySubPosition is a single lead-trader sub-position. The validating account
// has not led any trades (the endpoint returns an empty array), so the field set
// is modeled from the OKX doc field table.
type CopySubPosition struct {
	InstrumentType         InstType        `json:"instType"`
	InstrumentID           string          `json:"instId"`
	SubPositionID          string          `json:"subPosId"`
	PositionSide           PosSide         `json:"posSide"`
	MarginMode             MgnMode         `json:"mgnMode"`
	Leverage               decimal.Decimal `json:"lever"`
	OpenOrderID            string          `json:"openOrdId"`
	OpenAveragePrice       decimal.Decimal `json:"openAvgPx"`
	OpenTime               time.Time       `json:"openTime"`
	SubPosition            decimal.Decimal `json:"subPos"`
	Currency               string          `json:"ccy"`
	MarkPrice              decimal.Decimal `json:"markPx"`
	UPL                    decimal.Decimal `json:"upl"`
	UPLRatio               decimal.Decimal `json:"uplRatio"`
	StopLossTriggerPrice   decimal.Decimal `json:"slTriggerPx"`
	TakeProfitTriggerPrice decimal.Decimal `json:"tpTriggerPx"`
	Margin                 decimal.Decimal `json:"margin"`
	UniqueCode             string          `json:"uniqueCode"`
	AvailableSubPosition   decimal.Decimal `json:"availSubPos"`
}

// GetCopySubpositionsHistoryService -- GET /api/v5/copytrading/subpositions-history (Read)
//
// Returns the lead trader's closed sub-position history.
type GetCopySubpositionsHistoryService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetCopySubpositionsHistoryService() *GetCopySubpositionsHistoryService {
	return &GetCopySubpositionsHistoryService{c: c, params: map[string]string{}}
}

// SetInstType filters by product line (SWAP).
func (s *GetCopySubpositionsHistoryService) SetInstType(instType InstType) *GetCopySubpositionsHistoryService {
	s.params["instType"] = string(instType)
	return s
}

// SetInstId filters by a single instrument id.
func (s *GetCopySubpositionsHistoryService) SetInstId(instId string) *GetCopySubpositionsHistoryService {
	s.params["instId"] = instId
	return s
}

// SetAfter paginates to records with subPosId earlier than the given id.
func (s *GetCopySubpositionsHistoryService) SetAfter(after string) *GetCopySubpositionsHistoryService {
	s.params["after"] = after
	return s
}

// SetBefore paginates to records with subPosId newer than the given id.
func (s *GetCopySubpositionsHistoryService) SetBefore(before string) *GetCopySubpositionsHistoryService {
	s.params["before"] = before
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetCopySubpositionsHistoryService) SetLimit(limit int) *GetCopySubpositionsHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetCopySubpositionsHistoryService) Do(ctx context.Context) ([]CopySubPositionHistory, error) {
	req := request.Get(ctx, s.c, "/api/v5/copytrading/subpositions-history", s.params).WithSign()
	return request.DoList[CopySubPositionHistory](req)
}

// CopySubPositionHistory is a single closed lead-trader sub-position. The
// validating account has not led any trades (the endpoint returns an empty
// array), so the field set is modeled from the OKX doc field table.
type CopySubPositionHistory struct {
	InstrumentType    InstType        `json:"instType"`
	InstrumentID      string          `json:"instId"`
	SubPositionID     string          `json:"subPosId"`
	PositionSide      PosSide         `json:"posSide"`
	MarginMode        MgnMode         `json:"mgnMode"`
	Leverage          decimal.Decimal `json:"lever"`
	OpenAveragePrice  decimal.Decimal `json:"openAvgPx"`
	OpenTime          time.Time       `json:"openTime"`
	CloseAveragePrice decimal.Decimal `json:"closeAvgPx"`
	CloseTime         time.Time       `json:"closeTime"`
	SubPosition       decimal.Decimal `json:"subPos"`
	Currency          string          `json:"ccy"`
	Pnl               decimal.Decimal `json:"pnl"`
	PnlRatio          decimal.Decimal `json:"pnlRatio"`
	UniqueCode        string          `json:"uniqueCode"`
	Type              string          `json:"type"`
}

// GetCopyInstrumentsService -- GET /api/v5/copytrading/instruments (Read)
//
// Returns the instruments the current lead trader can lead-trade and whether
// each is enabled for leading.
type GetCopyInstrumentsService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetCopyInstrumentsService() *GetCopyInstrumentsService {
	return &GetCopyInstrumentsService{c: c, params: map[string]string{}}
}

// SetInstType filters by product line (SWAP).
func (s *GetCopyInstrumentsService) SetInstType(instType InstType) *GetCopyInstrumentsService {
	s.params["instType"] = string(instType)
	return s
}

func (s *GetCopyInstrumentsService) Do(ctx context.Context) ([]CopyInstrument, error) {
	req := request.Get(ctx, s.c, "/api/v5/copytrading/instruments", s.params).WithSign()
	return request.DoList[CopyInstrument](req)
}

// CopyInstrument is one lead-tradable instrument and its enabled flag.
type CopyInstrument struct {
	InstrumentID string `json:"instId"`
	Enabled      bool   `json:"enabled"`
}

// GetCopyProfitSharingDetailsService -- GET /api/v5/copytrading/profit-sharing-details (Read)
//
// Returns the lead trader's profit-sharing detail records (the share collected
// from each copier).
type GetCopyProfitSharingDetailsService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetCopyProfitSharingDetailsService() *GetCopyProfitSharingDetailsService {
	return &GetCopyProfitSharingDetailsService{c: c, params: map[string]string{}}
}

// SetInstType filters by product line (SWAP).
func (s *GetCopyProfitSharingDetailsService) SetInstType(instType InstType) *GetCopyProfitSharingDetailsService {
	s.params["instType"] = string(instType)
	return s
}

// SetAfter paginates to records with profitSharingId earlier than the given id.
func (s *GetCopyProfitSharingDetailsService) SetAfter(after string) *GetCopyProfitSharingDetailsService {
	s.params["after"] = after
	return s
}

// SetBefore paginates to records with profitSharingId newer than the given id.
func (s *GetCopyProfitSharingDetailsService) SetBefore(before string) *GetCopyProfitSharingDetailsService {
	s.params["before"] = before
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetCopyProfitSharingDetailsService) SetLimit(limit int) *GetCopyProfitSharingDetailsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetCopyProfitSharingDetailsService) Do(ctx context.Context) ([]CopyProfitSharingDetail, error) {
	req := request.Get(ctx, s.c, "/api/v5/copytrading/profit-sharing-details", s.params).WithSign()
	return request.DoList[CopyProfitSharingDetail](req)
}

// CopyProfitSharingDetail is one profit-sharing record. The validating account
// has not led any trades (the endpoint returns an empty array), so the field set
// is modeled from the OKX doc field table.
type CopyProfitSharingDetail struct {
	InstrumentType      InstType        `json:"instType"`
	Currency            string          `json:"ccy"`
	NickName            string          `json:"nickName"`
	ProfitSharingAmount decimal.Decimal `json:"profitSharingAmt"`
	ProfitSharingID     string          `json:"profitSharingId"`
	Timestamp           time.Time       `json:"ts"`
}

// GetCopyTotalProfitSharingService -- GET /api/v5/copytrading/total-profit-sharing (Read)
//
// Returns the lead trader's accumulated profit-sharing amount per currency.
type GetCopyTotalProfitSharingService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetCopyTotalProfitSharingService() *GetCopyTotalProfitSharingService {
	return &GetCopyTotalProfitSharingService{c: c, params: map[string]string{}}
}

// SetInstType filters by product line (SWAP).
func (s *GetCopyTotalProfitSharingService) SetInstType(instType InstType) *GetCopyTotalProfitSharingService {
	s.params["instType"] = string(instType)
	return s
}

func (s *GetCopyTotalProfitSharingService) Do(ctx context.Context) ([]CopyTotalProfitSharing, error) {
	req := request.Get(ctx, s.c, "/api/v5/copytrading/total-profit-sharing", s.params).WithSign()
	return request.DoList[CopyTotalProfitSharing](req)
}

// CopyTotalProfitSharing is the lead trader's total profit-sharing for a
// currency.
type CopyTotalProfitSharing struct {
	InstrumentType           InstType        `json:"instType"`
	Currency                 string          `json:"ccy"`
	TotalProfitSharingAmount decimal.Decimal `json:"totalProfitSharingAmt"`
}

// GetCopyUnrealizedProfitSharingDetailsService -- GET /api/v5/copytrading/unrealized-profit-sharing-details (Read)
//
// Returns the lead trader's pending (unrealized) profit-sharing per copier.
type GetCopyUnrealizedProfitSharingDetailsService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetCopyUnrealizedProfitSharingDetailsService() *GetCopyUnrealizedProfitSharingDetailsService {
	return &GetCopyUnrealizedProfitSharingDetailsService{c: c, params: map[string]string{}}
}

// SetInstType filters by product line (SWAP).
func (s *GetCopyUnrealizedProfitSharingDetailsService) SetInstType(instType InstType) *GetCopyUnrealizedProfitSharingDetailsService {
	s.params["instType"] = string(instType)
	return s
}

func (s *GetCopyUnrealizedProfitSharingDetailsService) Do(ctx context.Context) ([]CopyUnrealizedProfitSharingDetail, error) {
	req := request.Get(ctx, s.c, "/api/v5/copytrading/unrealized-profit-sharing-details", s.params).WithSign()
	return request.DoList[CopyUnrealizedProfitSharingDetail](req)
}

// CopyUnrealizedProfitSharingDetail is one copier's pending profit-sharing. The
// validating account lacks permission for this endpoint (code 50030), so the
// field set is modeled from the OKX doc field table.
type CopyUnrealizedProfitSharingDetail struct {
	InstrumentType                InstType        `json:"instType"`
	Currency                      string          `json:"ccy"`
	NickName                      string          `json:"nickName"`
	UnrealizedProfitSharingAmount decimal.Decimal `json:"unrealizedProfitSharingAmt"`
}

// GetCopyTotalUnrealizedProfitSharingService -- GET /api/v5/copytrading/total-unrealized-profit-sharing (Read)
//
// Returns the lead trader's total pending (unrealized) profit-sharing per
// currency.
type GetCopyTotalUnrealizedProfitSharingService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetCopyTotalUnrealizedProfitSharingService() *GetCopyTotalUnrealizedProfitSharingService {
	return &GetCopyTotalUnrealizedProfitSharingService{c: c, params: map[string]string{}}
}

// SetInstType filters by product line (SWAP).
func (s *GetCopyTotalUnrealizedProfitSharingService) SetInstType(instType InstType) *GetCopyTotalUnrealizedProfitSharingService {
	s.params["instType"] = string(instType)
	return s
}

func (s *GetCopyTotalUnrealizedProfitSharingService) Do(ctx context.Context) ([]CopyTotalUnrealizedProfitSharing, error) {
	req := request.Get(ctx, s.c, "/api/v5/copytrading/total-unrealized-profit-sharing", s.params).WithSign()
	return request.DoList[CopyTotalUnrealizedProfitSharing](req)
}

// CopyTotalUnrealizedProfitSharing is the lead trader's total pending
// profit-sharing for a currency. The validating account lacks permission for
// this endpoint (code 50030), so the field set is modeled from the OKX doc field
// table.
type CopyTotalUnrealizedProfitSharing struct {
	InstrumentType                     InstType        `json:"instType"`
	Currency                           string          `json:"ccy"`
	TotalUnrealizedProfitSharingAmount decimal.Decimal `json:"totalUnrealizedProfitSharingAmt"`
}

// GetCopyConfigService -- GET /api/v5/copytrading/config (Read)
//
// Returns the current account's copy-trading role configuration (lead / copy
// state and per-role nick names / unique codes).
type GetCopyConfigService struct {
	c *Client
}

func (c *Client) NewGetCopyConfigService() *GetCopyConfigService {
	return &GetCopyConfigService{c: c}
}

func (s *GetCopyConfigService) Do(ctx context.Context) (*CopyConfig, error) {
	req := request.Get(ctx, s.c, "/api/v5/copytrading/config").WithSign()
	return request.DoOne[CopyConfig](req)
}

// CopyConfig is the account's copy-trading role configuration. The validating
// account has neither led nor copied any trades (the endpoint returns code
// 59285), so the field set is modeled from the OKX doc field table.
type CopyConfig struct {
	NickName   string    `json:"nickName"`
	PortLink   string    `json:"portLink"`
	Roles      []string  `json:"roles"`
	CopyState  CopyState `json:"copyState"`
	LeadState  CopyState `json:"leadState"`
	UniqueCode string    `json:"uniqueCode"`
}

// GetCopyPublicConfigService -- GET /api/v5/copytrading/public-config (public)
//
// Returns the platform-wide copy-trading limits (min/max copy amounts, ratios
// and SL/TP caps). Public; no signing required.
type GetCopyPublicConfigService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetCopyPublicConfigService() *GetCopyPublicConfigService {
	return &GetCopyPublicConfigService{c: c, params: map[string]string{}}
}

// SetInstType filters by product line (SWAP).
func (s *GetCopyPublicConfigService) SetInstType(instType InstType) *GetCopyPublicConfigService {
	s.params["instType"] = string(instType)
	return s
}

func (s *GetCopyPublicConfigService) Do(ctx context.Context) (*CopyPublicConfig, error) {
	req := request.Get(ctx, s.c, "/api/v5/copytrading/public-config", s.params)
	return request.DoOne[CopyPublicConfig](req)
}

// CopyPublicConfig is the platform-wide copy-trading configuration.
type CopyPublicConfig struct {
	MaxCopyAmount      decimal.Decimal `json:"maxCopyAmt"`
	MinCopyAmount      decimal.Decimal `json:"minCopyAmt"`
	MaxCopyTotalAmount decimal.Decimal `json:"maxCopyTotalAmt"`
	MaxCopyRatio       decimal.Decimal `json:"maxCopyRatio"`
	MinCopyRatio       decimal.Decimal `json:"minCopyRatio"`
	MaxTakeProfitRatio decimal.Decimal `json:"maxTpRatio"`
	MaxStopLossRatio   decimal.Decimal `json:"maxSlRatio"`
}

// GetCopyPublicLeadTradersService -- GET /api/v5/copytrading/public-lead-traders (public)
//
// Returns the public lead-trader ranking. Public; no signing required.
type GetCopyPublicLeadTradersService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetCopyPublicLeadTradersService() *GetCopyPublicLeadTradersService {
	return &GetCopyPublicLeadTradersService{c: c, params: map[string]string{}}
}

// SetInstType filters by product line (SWAP).
func (s *GetCopyPublicLeadTradersService) SetInstType(instType InstType) *GetCopyPublicLeadTradersService {
	s.params["instType"] = string(instType)
	return s
}

// SetSortType sets the ranking sort key (overview/pnl/pnlRatio/aum/winRatio).
func (s *GetCopyPublicLeadTradersService) SetSortType(sortType CopySortType) *GetCopyPublicLeadTradersService {
	s.params["sortType"] = string(sortType)
	return s
}

// SetState filters by copy state ("0" not copying / "1" copying).
func (s *GetCopyPublicLeadTradersService) SetState(state CopyState) *GetCopyPublicLeadTradersService {
	s.params["state"] = string(state)
	return s
}

// SetMinLeadDays filters by minimum lead days (1/2/3/4 buckets).
func (s *GetCopyPublicLeadTradersService) SetMinLeadDays(minLeadDays string) *GetCopyPublicLeadTradersService {
	s.params["minLeadDays"] = minLeadDays
	return s
}

// SetMinAssets filters by minimum trading assets.
func (s *GetCopyPublicLeadTradersService) SetMinAssets(minAssets decimal.Decimal) *GetCopyPublicLeadTradersService {
	s.params["minAssets"] = minAssets.String()
	return s
}

// SetMaxAssets filters by maximum trading assets.
func (s *GetCopyPublicLeadTradersService) SetMaxAssets(maxAssets decimal.Decimal) *GetCopyPublicLeadTradersService {
	s.params["maxAssets"] = maxAssets.String()
	return s
}

// SetMinAum filters by minimum assets under management.
func (s *GetCopyPublicLeadTradersService) SetMinAum(minAum decimal.Decimal) *GetCopyPublicLeadTradersService {
	s.params["minAum"] = minAum.String()
	return s
}

// SetMaxAum filters by maximum assets under management.
func (s *GetCopyPublicLeadTradersService) SetMaxAum(maxAum decimal.Decimal) *GetCopyPublicLeadTradersService {
	s.params["maxAum"] = maxAum.String()
	return s
}

// SetDataVer sets the data version (snapshot id from a prior call).
func (s *GetCopyPublicLeadTradersService) SetDataVer(dataVer string) *GetCopyPublicLeadTradersService {
	s.params["dataVer"] = dataVer
	return s
}

// SetPage sets the page number.
func (s *GetCopyPublicLeadTradersService) SetPage(page int) *GetCopyPublicLeadTradersService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of records returned (max 20, default 10).
func (s *GetCopyPublicLeadTradersService) SetLimit(limit int) *GetCopyPublicLeadTradersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetCopyPublicLeadTradersService) Do(ctx context.Context) (*CopyLeadTradersPage, error) {
	req := request.Get(ctx, s.c, "/api/v5/copytrading/public-lead-traders", s.params)
	return request.DoOne[CopyLeadTradersPage](req)
}

// CopyLeadTradersPage is one page of the lead-trader ranking.
type CopyLeadTradersPage struct {
	DataVersion string           `json:"dataVer"`
	TotalPage   decimal.Decimal  `json:"totalPage"`
	Ranks       []CopyLeadTrader `json:"ranks"`
}

// CopyLeadTrader is one lead trader's ranking row. chanType is only present in
// the signed lead-traders response.
type CopyLeadTrader struct {
	UniqueCode                  string                   `json:"uniqueCode"`
	NickName                    string                   `json:"nickName"`
	PortLink                    string                   `json:"portLink"`
	Currency                    string                   `json:"ccy"`
	ChannelType                 string                   `json:"chanType"`
	CopyState                   CopyState                `json:"copyState"`
	AUM                         decimal.Decimal          `json:"aum"`
	Pnl                         decimal.Decimal          `json:"pnl"`
	PnlRatio                    decimal.Decimal          `json:"pnlRatio"`
	WinRatio                    decimal.Decimal          `json:"winRatio"`
	LeadDays                    decimal.Decimal          `json:"leadDays"`
	CopyTraderNumber            decimal.Decimal          `json:"copyTraderNum"`
	MaxCopyTraderNumber         decimal.Decimal          `json:"maxCopyTraderNum"`
	AccumulatedCopyTraderNumber decimal.Decimal          `json:"accCopyTraderNum"`
	TraderInstruments           []string                 `json:"traderInsts"`
	PnlRatios                   []CopyLeadTraderPnlRatio `json:"pnlRatios"`
}

// CopyLeadTraderPnlRatio is a single dated pnl-ratio point in a lead trader's
// history curve.
type CopyLeadTraderPnlRatio struct {
	BeginTimestamp time.Time       `json:"beginTs"`
	PnlRatio       decimal.Decimal `json:"pnlRatio"`
}

// GetCopyLeadTradersService -- GET /api/v5/copytrading/lead-traders (Read)
//
// Returns the lead-trader ranking visible to the current (signed) account. Shape
// matches the public ranking but the rows additionally carry chanType.
type GetCopyLeadTradersService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetCopyLeadTradersService() *GetCopyLeadTradersService {
	return &GetCopyLeadTradersService{c: c, params: map[string]string{}}
}

// SetInstType filters by product line (SWAP).
func (s *GetCopyLeadTradersService) SetInstType(instType InstType) *GetCopyLeadTradersService {
	s.params["instType"] = string(instType)
	return s
}

// SetSortType sets the ranking sort key (overview/pnl/pnlRatio/aum/winRatio).
func (s *GetCopyLeadTradersService) SetSortType(sortType CopySortType) *GetCopyLeadTradersService {
	s.params["sortType"] = string(sortType)
	return s
}

// SetState filters by copy state ("0" not copying / "1" copying).
func (s *GetCopyLeadTradersService) SetState(state CopyState) *GetCopyLeadTradersService {
	s.params["state"] = string(state)
	return s
}

// SetMinLeadDays filters by minimum lead days (1/2/3/4 buckets).
func (s *GetCopyLeadTradersService) SetMinLeadDays(minLeadDays string) *GetCopyLeadTradersService {
	s.params["minLeadDays"] = minLeadDays
	return s
}

// SetMinAssets filters by minimum trading assets.
func (s *GetCopyLeadTradersService) SetMinAssets(minAssets decimal.Decimal) *GetCopyLeadTradersService {
	s.params["minAssets"] = minAssets.String()
	return s
}

// SetMaxAssets filters by maximum trading assets.
func (s *GetCopyLeadTradersService) SetMaxAssets(maxAssets decimal.Decimal) *GetCopyLeadTradersService {
	s.params["maxAssets"] = maxAssets.String()
	return s
}

// SetMinAum filters by minimum assets under management.
func (s *GetCopyLeadTradersService) SetMinAum(minAum decimal.Decimal) *GetCopyLeadTradersService {
	s.params["minAum"] = minAum.String()
	return s
}

// SetMaxAum filters by maximum assets under management.
func (s *GetCopyLeadTradersService) SetMaxAum(maxAum decimal.Decimal) *GetCopyLeadTradersService {
	s.params["maxAum"] = maxAum.String()
	return s
}

// SetDataVer sets the data version (snapshot id from a prior call).
func (s *GetCopyLeadTradersService) SetDataVer(dataVer string) *GetCopyLeadTradersService {
	s.params["dataVer"] = dataVer
	return s
}

// SetPage sets the page number.
func (s *GetCopyLeadTradersService) SetPage(page int) *GetCopyLeadTradersService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit caps the number of records returned (max 20, default 10).
func (s *GetCopyLeadTradersService) SetLimit(limit int) *GetCopyLeadTradersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetCopyLeadTradersService) Do(ctx context.Context) (*CopyLeadTradersPage, error) {
	req := request.Get(ctx, s.c, "/api/v5/copytrading/lead-traders", s.params).WithSign()
	return request.DoOne[CopyLeadTradersPage](req)
}

// GetCopySettingsService -- GET /api/v5/copytrading/copy-settings (Read)
//
// Returns the current account's copy settings for a specific lead trader.
type GetCopySettingsService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetCopySettingsService(instType InstType, uniqueCode string) *GetCopySettingsService {
	return &GetCopySettingsService{c: c, params: map[string]string{
		"instType":   string(instType),
		"uniqueCode": uniqueCode,
	}}
}

func (s *GetCopySettingsService) Do(ctx context.Context) (*CopySettings, error) {
	req := request.Get(ctx, s.c, "/api/v5/copytrading/copy-settings", s.params).WithSign()
	return request.DoOne[CopySettings](req)
}

// CopySettings is the current account's copy configuration for a lead trader.
type CopySettings struct {
	Currency             string          `json:"ccy"`
	CopyAmount           decimal.Decimal `json:"copyAmt"`
	CopyTotalAmount      decimal.Decimal `json:"copyTotalAmt"`
	CopyInstrumentIDType string          `json:"copyInstIdType"`
	CopyMarginMode       string          `json:"copyMgnMode"`
	CopyMode             string          `json:"copyMode"`
	CopyRatio            decimal.Decimal `json:"copyRatio"`
	CopyState            CopyState       `json:"copyState"`
	InstrumentIDs        []string        `json:"instIds"`
	StopLossRatio        decimal.Decimal `json:"slRatio"`
	StopLossTotalAmount  decimal.Decimal `json:"slTotalAmt"`
	SubPositionCloseType string          `json:"subPosCloseType"`
	Tag                  string          `json:"tag"`
	TakeProfitRatio      decimal.Decimal `json:"tpRatio"`
}

// GetCopyBatchLeverageInfoService -- GET /api/v5/copytrading/batch-leverage-info (Read)
//
// Returns the leverage a lead trader uses for a set of instruments under the
// given margin mode.
type GetCopyBatchLeverageInfoService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetCopyBatchLeverageInfoService(mgnMode MgnMode, uniqueCode string) *GetCopyBatchLeverageInfoService {
	return &GetCopyBatchLeverageInfoService{c: c, params: map[string]string{
		"mgnMode":    string(mgnMode),
		"uniqueCode": uniqueCode,
	}}
}

// SetInstId filters by instrument id (single or comma-separated).
func (s *GetCopyBatchLeverageInfoService) SetInstId(instId string) *GetCopyBatchLeverageInfoService {
	s.params["instId"] = instId
	return s
}

func (s *GetCopyBatchLeverageInfoService) Do(ctx context.Context) ([]CopyBatchLeverageInfo, error) {
	req := request.Get(ctx, s.c, "/api/v5/copytrading/batch-leverage-info", s.params).WithSign()
	return request.DoList[CopyBatchLeverageInfo](req)
}

// CopyBatchLeverageInfo is a lead trader's leverage for an instrument. The
// validating account is not on the copy-trading allowlist (code 59263), so the
// field set is modeled from the OKX doc field table.
type CopyBatchLeverageInfo struct {
	InstrumentID string          `json:"instId"`
	MarginMode   MgnMode         `json:"mgnMode"`
	Leverage     decimal.Decimal `json:"lever"`
	PositionSide PosSide         `json:"posSide"`
}

// GetCopyTradersService -- GET /api/v5/copytrading/copy-traders (Read)
//
// Returns the copiers currently following a given lead trader.
type GetCopyTradersService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetCopyTradersService(uniqueCode string) *GetCopyTradersService {
	return &GetCopyTradersService{c: c, params: map[string]string{"uniqueCode": uniqueCode}}
}

// SetInstType filters by product line (SWAP).
func (s *GetCopyTradersService) SetInstType(instType InstType) *GetCopyTradersService {
	s.params["instType"] = string(instType)
	return s
}

// SetLimit caps the number of records returned (max 10, default 10).
func (s *GetCopyTradersService) SetLimit(limit int) *GetCopyTradersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetCopyTradersService) Do(ctx context.Context) (*CopyTradersPage, error) {
	req := request.Get(ctx, s.c, "/api/v5/copytrading/copy-traders", s.params).WithSign()
	return request.DoOne[CopyTradersPage](req)
}

// CopyTradersPage is the summary plus per-copier list following a lead trader.
// The validating account is not on the copy-trading allowlist (code 59263), so
// the field set is modeled from the OKX doc field table.
type CopyTradersPage struct {
	InstrumentType   InstType        `json:"instType"`
	CopyTotalAmount  decimal.Decimal `json:"copyTotalAmt"`
	CopyTotalPnl     decimal.Decimal `json:"copyTotalPnl"`
	Currency         string          `json:"ccy"`
	CopyTraderNumber decimal.Decimal `json:"copyTraderNum"`
	CopyTraders      []CopyTrader    `json:"copyTraders"`
}

// CopyTrader is one copier following a lead trader.
type CopyTrader struct {
	NickName      string          `json:"nickName"`
	PortLink      string          `json:"portLink"`
	BeginCopyTime time.Time       `json:"beginCopyTime"`
	Pnl           decimal.Decimal `json:"pnl"`
}

// =====================================================================
// State-changing endpoints below are IMPLEMENT-ONLY. They place/close
// lead orders, change copy settings, or toggle lead/copy trading on a
// REAL account and must NOT be exercised by tests. Bodies and ack shapes
// are modeled from the OKX docs.
// =====================================================================

// PlaceCopyLeadStopOrderService -- POST /api/v5/copytrading/place-lead-stop-order (Trade)
//
// Places a take-profit / stop-loss order against a lead trader's sub-position.
type PlaceCopyLeadStopOrderService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewPlaceCopyLeadStopOrderService(subPosId string) *PlaceCopyLeadStopOrderService {
	return &PlaceCopyLeadStopOrderService{c: c, body: map[string]any{"subPosId": subPosId}}
}

// SetTpTriggerPx sets the take-profit trigger price.
func (s *PlaceCopyLeadStopOrderService) SetTpTriggerPx(px decimal.Decimal) *PlaceCopyLeadStopOrderService {
	s.body["tpTriggerPx"] = px.String()
	return s
}

// SetSlTriggerPx sets the stop-loss trigger price.
func (s *PlaceCopyLeadStopOrderService) SetSlTriggerPx(px decimal.Decimal) *PlaceCopyLeadStopOrderService {
	s.body["slTriggerPx"] = px.String()
	return s
}

// SetTpTriggerPxType sets the take-profit trigger price type (last/index/mark).
func (s *PlaceCopyLeadStopOrderService) SetTpTriggerPxType(pxType string) *PlaceCopyLeadStopOrderService {
	s.body["tpTriggerPxType"] = pxType
	return s
}

// SetSlTriggerPxType sets the stop-loss trigger price type (last/index/mark).
func (s *PlaceCopyLeadStopOrderService) SetSlTriggerPxType(pxType string) *PlaceCopyLeadStopOrderService {
	s.body["slTriggerPxType"] = pxType
	return s
}

func (s *PlaceCopyLeadStopOrderService) Do(ctx context.Context) (*CopySubPosAck, error) {
	req := request.Post(ctx, s.c, "/api/v5/copytrading/place-lead-stop-order", s.body).WithSign()
	return request.DoOne[CopySubPosAck](req)
}

// CopySubPosAck is the acknowledgement for a sub-position lead order/close.
type CopySubPosAck struct {
	SubPositionID string `json:"subPosId"`
	Tag           string `json:"tag"`
}

// CloseCopySubpositionService -- POST /api/v5/copytrading/close-subposition (Trade)
//
// Market-closes a lead trader's open sub-position.
type CloseCopySubpositionService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewCloseCopySubpositionService(subPosId string) *CloseCopySubpositionService {
	return &CloseCopySubpositionService{c: c, body: map[string]any{"subPosId": subPosId}}
}

func (s *CloseCopySubpositionService) Do(ctx context.Context) (*CopySubPosAck, error) {
	req := request.Post(ctx, s.c, "/api/v5/copytrading/close-subposition", s.body).WithSign()
	return request.DoOne[CopySubPosAck](req)
}

// AmendCopyLeadingInstrumentsService -- POST /api/v5/copytrading/amend-leading-instruments (Trade)
//
// Updates the set of instruments the current lead trader leads on.
type AmendCopyLeadingInstrumentsService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewAmendCopyLeadingInstrumentsService(instId string) *AmendCopyLeadingInstrumentsService {
	return &AmendCopyLeadingInstrumentsService{c: c, body: map[string]any{"instId": instId}}
}

// SetInstType sets the product line (SWAP).
func (s *AmendCopyLeadingInstrumentsService) SetInstType(instType InstType) *AmendCopyLeadingInstrumentsService {
	s.body["instType"] = string(instType)
	return s
}

func (s *AmendCopyLeadingInstrumentsService) Do(ctx context.Context) ([]CopyInstrument, error) {
	req := request.Post(ctx, s.c, "/api/v5/copytrading/amend-leading-instruments", s.body).WithSign()
	return request.DoList[CopyInstrument](req)
}

// FirstCopyCopySettingsService -- POST /api/v5/copytrading/first-copy-settings (Trade)
//
// Starts copy-trading a lead trader for the first time, with the given copy
// configuration.
type FirstCopyCopySettingsService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewFirstCopyCopySettingsService(uniqueCode, copyMgnMode, copyInstIdType, copyMode, subPosCloseType string) *FirstCopyCopySettingsService {
	return &FirstCopyCopySettingsService{c: c, body: map[string]any{
		"uniqueCode":      uniqueCode,
		"copyMgnMode":     copyMgnMode,
		"copyInstIdType":  copyInstIdType,
		"copyMode":        copyMode,
		"subPosCloseType": subPosCloseType,
	}}
}

// SetInstType sets the product line (SWAP).
func (s *FirstCopyCopySettingsService) SetInstType(instType InstType) *FirstCopyCopySettingsService {
	s.body["instType"] = string(instType)
	return s
}

// SetInstId sets the comma-separated instruments to copy (copyInstIdType=custom).
func (s *FirstCopyCopySettingsService) SetInstId(instId string) *FirstCopyCopySettingsService {
	s.body["instId"] = instId
	return s
}

// SetCopyAmt sets the fixed copy amount per order (copyMode=fixed_amount).
func (s *FirstCopyCopySettingsService) SetCopyAmt(copyAmt decimal.Decimal) *FirstCopyCopySettingsService {
	s.body["copyAmt"] = copyAmt.String()
	return s
}

// SetCopyRatio sets the copy ratio (copyMode=ratio_copy).
func (s *FirstCopyCopySettingsService) SetCopyRatio(copyRatio decimal.Decimal) *FirstCopyCopySettingsService {
	s.body["copyRatio"] = copyRatio.String()
	return s
}

// SetCopyTotalAmt sets the total amount budgeted for copying.
func (s *FirstCopyCopySettingsService) SetCopyTotalAmt(copyTotalAmt decimal.Decimal) *FirstCopyCopySettingsService {
	s.body["copyTotalAmt"] = copyTotalAmt.String()
	return s
}

// SetTpRatio sets the take-profit ratio.
func (s *FirstCopyCopySettingsService) SetTpRatio(tpRatio decimal.Decimal) *FirstCopyCopySettingsService {
	s.body["tpRatio"] = tpRatio.String()
	return s
}

// SetSlRatio sets the stop-loss ratio.
func (s *FirstCopyCopySettingsService) SetSlRatio(slRatio decimal.Decimal) *FirstCopyCopySettingsService {
	s.body["slRatio"] = slRatio.String()
	return s
}

// SetSlTotalAmt sets the total stop-loss amount.
func (s *FirstCopyCopySettingsService) SetSlTotalAmt(slTotalAmt decimal.Decimal) *FirstCopyCopySettingsService {
	s.body["slTotalAmt"] = slTotalAmt.String()
	return s
}

func (s *FirstCopyCopySettingsService) Do(ctx context.Context) (*CopyResult, error) {
	req := request.Post(ctx, s.c, "/api/v5/copytrading/first-copy-settings", s.body).WithSign()
	return request.DoOne[CopyResult](req)
}

// CopyResult is the generic ack for the copy/lead settings POST endpoints.
type CopyResult struct {
	Result bool `json:"result"`
}

// AmendCopyCopySettingsService -- POST /api/v5/copytrading/amend-copy-settings (Trade)
//
// Updates the copy configuration for a lead trader already being copied.
type AmendCopyCopySettingsService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewAmendCopyCopySettingsService(uniqueCode, copyMgnMode, copyInstIdType, copyMode, subPosCloseType string) *AmendCopyCopySettingsService {
	return &AmendCopyCopySettingsService{c: c, body: map[string]any{
		"uniqueCode":      uniqueCode,
		"copyMgnMode":     copyMgnMode,
		"copyInstIdType":  copyInstIdType,
		"copyMode":        copyMode,
		"subPosCloseType": subPosCloseType,
	}}
}

// SetInstType sets the product line (SWAP).
func (s *AmendCopyCopySettingsService) SetInstType(instType InstType) *AmendCopyCopySettingsService {
	s.body["instType"] = string(instType)
	return s
}

// SetInstId sets the comma-separated instruments to copy (copyInstIdType=custom).
func (s *AmendCopyCopySettingsService) SetInstId(instId string) *AmendCopyCopySettingsService {
	s.body["instId"] = instId
	return s
}

// SetCopyAmt sets the fixed copy amount per order (copyMode=fixed_amount).
func (s *AmendCopyCopySettingsService) SetCopyAmt(copyAmt decimal.Decimal) *AmendCopyCopySettingsService {
	s.body["copyAmt"] = copyAmt.String()
	return s
}

// SetCopyRatio sets the copy ratio (copyMode=ratio_copy).
func (s *AmendCopyCopySettingsService) SetCopyRatio(copyRatio decimal.Decimal) *AmendCopyCopySettingsService {
	s.body["copyRatio"] = copyRatio.String()
	return s
}

// SetCopyTotalAmt sets the total amount budgeted for copying.
func (s *AmendCopyCopySettingsService) SetCopyTotalAmt(copyTotalAmt decimal.Decimal) *AmendCopyCopySettingsService {
	s.body["copyTotalAmt"] = copyTotalAmt.String()
	return s
}

// SetTpRatio sets the take-profit ratio.
func (s *AmendCopyCopySettingsService) SetTpRatio(tpRatio decimal.Decimal) *AmendCopyCopySettingsService {
	s.body["tpRatio"] = tpRatio.String()
	return s
}

// SetSlRatio sets the stop-loss ratio.
func (s *AmendCopyCopySettingsService) SetSlRatio(slRatio decimal.Decimal) *AmendCopyCopySettingsService {
	s.body["slRatio"] = slRatio.String()
	return s
}

// SetSlTotalAmt sets the total stop-loss amount.
func (s *AmendCopyCopySettingsService) SetSlTotalAmt(slTotalAmt decimal.Decimal) *AmendCopyCopySettingsService {
	s.body["slTotalAmt"] = slTotalAmt.String()
	return s
}

func (s *AmendCopyCopySettingsService) Do(ctx context.Context) (*CopyResult, error) {
	req := request.Post(ctx, s.c, "/api/v5/copytrading/amend-copy-settings", s.body).WithSign()
	return request.DoOne[CopyResult](req)
}

// StopCopyCopyTradingService -- POST /api/v5/copytrading/stop-copy-trading (Trade)
//
// Stops copying a lead trader and (per subPosCloseType) closes or keeps the
// mirrored positions.
type StopCopyCopyTradingService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewStopCopyCopyTradingService(uniqueCode, subPosCloseType string) *StopCopyCopyTradingService {
	return &StopCopyCopyTradingService{c: c, body: map[string]any{
		"uniqueCode":      uniqueCode,
		"subPosCloseType": subPosCloseType,
	}}
}

// SetInstType sets the product line (SWAP).
func (s *StopCopyCopyTradingService) SetInstType(instType InstType) *StopCopyCopyTradingService {
	s.body["instType"] = string(instType)
	return s
}

func (s *StopCopyCopyTradingService) Do(ctx context.Context) (*CopyResult, error) {
	req := request.Post(ctx, s.c, "/api/v5/copytrading/stop-copy-trading", s.body).WithSign()
	return request.DoOne[CopyResult](req)
}

// AmendCopyProfitSharingRatioService -- POST /api/v5/copytrading/amend-profit-sharing-ratio (Trade)
//
// Updates the lead trader's profit-sharing ratio.
type AmendCopyProfitSharingRatioService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewAmendCopyProfitSharingRatioService(profitSharingRatio decimal.Decimal) *AmendCopyProfitSharingRatioService {
	return &AmendCopyProfitSharingRatioService{c: c, body: map[string]any{
		"profitSharingRatio": profitSharingRatio.String(),
	}}
}

// SetInstType sets the product line (SWAP).
func (s *AmendCopyProfitSharingRatioService) SetInstType(instType InstType) *AmendCopyProfitSharingRatioService {
	s.body["instType"] = string(instType)
	return s
}

func (s *AmendCopyProfitSharingRatioService) Do(ctx context.Context) (*CopyResult, error) {
	req := request.Post(ctx, s.c, "/api/v5/copytrading/amend-profit-sharing-ratio", s.body).WithSign()
	return request.DoOne[CopyResult](req)
}

// SetCopyLeverageService -- POST /api/v5/copytrading/set-leverage (Trade)
//
// Sets the leverage a lead trader uses for one or more instruments.
type SetCopyLeverageService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewSetCopyLeverageService(mgnMode MgnMode, lever decimal.Decimal) *SetCopyLeverageService {
	return &SetCopyLeverageService{c: c, body: map[string]any{
		"mgnMode": string(mgnMode),
		"lever":   lever.String(),
	}}
}

// SetInstId sets the comma-separated instruments to set leverage on.
func (s *SetCopyLeverageService) SetInstId(instId string) *SetCopyLeverageService {
	s.body["instId"] = instId
	return s
}

func (s *SetCopyLeverageService) Do(ctx context.Context) ([]CopyBatchLeverageInfo, error) {
	req := request.Post(ctx, s.c, "/api/v5/copytrading/set-leverage", s.body).WithSign()
	return request.DoList[CopyBatchLeverageInfo](req)
}

// ApplyCopyLeadTradingService -- POST /api/v5/copytrading/apply-lead-trading (Trade)
//
// Applies to become a lead trader on the given instruments.
type ApplyCopyLeadTradingService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewApplyCopyLeadTradingService(instId string) *ApplyCopyLeadTradingService {
	return &ApplyCopyLeadTradingService{c: c, body: map[string]any{"instId": instId}}
}

// SetInstType sets the product line (SWAP).
func (s *ApplyCopyLeadTradingService) SetInstType(instType InstType) *ApplyCopyLeadTradingService {
	s.body["instType"] = string(instType)
	return s
}

func (s *ApplyCopyLeadTradingService) Do(ctx context.Context) ([]CopyInstrument, error) {
	req := request.Post(ctx, s.c, "/api/v5/copytrading/apply-lead-trading", s.body).WithSign()
	return request.DoList[CopyInstrument](req)
}

// StopCopyLeadTradingService -- POST /api/v5/copytrading/stop-lead-trading (Trade)
//
// Stops being a lead trader (closes lead positions per platform rules).
type StopCopyLeadTradingService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewStopCopyLeadTradingService() *StopCopyLeadTradingService {
	return &StopCopyLeadTradingService{c: c, body: map[string]any{}}
}

// SetInstType sets the product line (SWAP).
func (s *StopCopyLeadTradingService) SetInstType(instType InstType) *StopCopyLeadTradingService {
	s.body["instType"] = string(instType)
	return s
}

func (s *StopCopyLeadTradingService) Do(ctx context.Context) (*CopyResult, error) {
	req := request.Post(ctx, s.c, "/api/v5/copytrading/stop-lead-trading", s.body).WithSign()
	return request.DoOne[CopyResult](req)
}
