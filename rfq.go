package okx

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-okx/request"
	"github.com/shopspring/decimal"
)

// This file covers OKX's Block-trading (RFQ) endpoints under /api/v5/rfq/*.
// Most endpoints are private (signed); the block-trade quote/trade public feeds
// omit signing. The block-trade ticker endpoints documented here actually live
// under /api/v5/market/* (the /api/v5/rfq/block-ticker[s] paths return HTTP 404);
// the services below target the real /market paths.
//
// All exported identifiers are Rfq*-prefixed to avoid collisions with the shared
// types owned by market.go (BlockTicker etc.).

// RfqState is the lifecycle state of an RFQ.
type RfqState string

const (
	RfqStateActive   RfqState = "active"
	RfqStateCanceled RfqState = "canceled"
	RfqStatePendFill RfqState = "pending_fill"
	RfqStateFilled   RfqState = "filled"
	RfqStateExpired  RfqState = "expired"
	RfqStateTraded   RfqState = "traded_away"
	RfqStateFailed   RfqState = "failed"
)

// RfqQuoteState is the lifecycle state of a quote.
type RfqQuoteState string

const (
	RfqQuoteStateActive   RfqQuoteState = "active"
	RfqQuoteStateCanceled RfqQuoteState = "canceled"
	RfqQuoteStatePendFill RfqQuoteState = "pending_fill"
	RfqQuoteStateFilled   RfqQuoteState = "filled"
	RfqQuoteStateExpired  RfqQuoteState = "expired"
	RfqQuoteStateFailed   RfqQuoteState = "failed"
)

// RfqQuoteSide is the side a maker quotes an RFQ on.
type RfqQuoteSide string

const (
	RfqQuoteSideBuy  RfqQuoteSide = "buy"
	RfqQuoteSideSell RfqQuoteSide = "sell"
)

// GetRfqCounterpartiesService -- GET /api/v5/rfq/counterparties (Read)
//
// Returns the list of counterparties (makers) the account may direct RFQs to.
type GetRfqCounterpartiesService struct {
	c *Client
}

func (c *Client) NewGetRfqCounterpartiesService() *GetRfqCounterpartiesService {
	return &GetRfqCounterpartiesService{c: c}
}

func (s *GetRfqCounterpartiesService) Do(ctx context.Context) ([]RfqCounterparty, error) {
	req := request.Get(ctx, s.c, "/api/v5/rfq/counterparties").WithSign()
	return request.DoList[RfqCounterparty](req)
}

// RfqCounterparty is one available RFQ counterparty (maker).
type RfqCounterparty struct {
	TraderName string `json:"traderName"`
	TraderCode string `json:"traderCode"`
	Type       string `json:"type"`
}

// GetRfqsService -- GET /api/v5/rfq/rfqs (Read)
//
// Returns the account's RFQ records, both created (taker) and received (maker).
type GetRfqsService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetRfqsService() *GetRfqsService {
	return &GetRfqsService{c: c, params: map[string]string{}}
}

// SetRfqId filters by RFQ id.
func (s *GetRfqsService) SetRfqId(rfqId string) *GetRfqsService {
	s.params["rfqId"] = rfqId
	return s
}

// SetClRfqId filters by client-supplied RFQ id.
func (s *GetRfqsService) SetClRfqId(clRfqId string) *GetRfqsService {
	s.params["clRfqId"] = clRfqId
	return s
}

// SetState filters by RFQ state.
func (s *GetRfqsService) SetState(state RfqState) *GetRfqsService {
	s.params["state"] = string(state)
	return s
}

// SetBeginId pages from a starting RFQ id (records newer than the id).
func (s *GetRfqsService) SetBeginId(beginId string) *GetRfqsService {
	s.params["beginId"] = beginId
	return s
}

// SetEndId pages up to an ending RFQ id (records older than the id).
func (s *GetRfqsService) SetEndId(endId string) *GetRfqsService {
	s.params["endId"] = endId
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetRfqsService) SetLimit(limit int) *GetRfqsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetRfqsService) Do(ctx context.Context) ([]Rfq, error) {
	req := request.Get(ctx, s.c, "/api/v5/rfq/rfqs", s.params).WithSign()
	return request.DoList[Rfq](req)
}

// Rfq is one RFQ record. The validating account has no RFQ history, so the field
// set is modeled from the OKX doc field table.
type Rfq struct {
	CreateTime     time.Time `json:"cTime"`
	UpdateTime     time.Time `json:"uTime"`
	TraderCode     string    `json:"traderCode"`
	RFQID          string    `json:"rfqId"`
	ClientRFQID    string    `json:"clRfqId"`
	State          RfqState  `json:"state"`
	ValidUntil     time.Time `json:"validUntil"`
	Counterparties []string  `json:"counterparties"`
	Legs           []RfqLeg  `json:"legs"`
	AllowPartial   bool      `json:"allowPartialExecution"`
}

// RfqLeg is one leg of an RFQ.
type RfqLeg struct {
	InstrumentID   string          `json:"instId"`
	Size           decimal.Decimal `json:"sz"`
	Side           Side            `json:"side"`
	TargetCurrency TgtCcy          `json:"tgtCcy"`
	PositionSide   PosSide         `json:"posSide"`
}

// GetRfqQuotesService -- GET /api/v5/rfq/quotes (Read)
//
// Returns the account's quote records (quotes it made, or quotes it received as
// taker).
type GetRfqQuotesService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetRfqQuotesService() *GetRfqQuotesService {
	return &GetRfqQuotesService{c: c, params: map[string]string{}}
}

// SetRfqId filters by RFQ id.
func (s *GetRfqQuotesService) SetRfqId(rfqId string) *GetRfqQuotesService {
	s.params["rfqId"] = rfqId
	return s
}

// SetClRfqId filters by client-supplied RFQ id.
func (s *GetRfqQuotesService) SetClRfqId(clRfqId string) *GetRfqQuotesService {
	s.params["clRfqId"] = clRfqId
	return s
}

// SetQuoteId filters by quote id.
func (s *GetRfqQuotesService) SetQuoteId(quoteId string) *GetRfqQuotesService {
	s.params["quoteId"] = quoteId
	return s
}

// SetClQuoteId filters by client-supplied quote id.
func (s *GetRfqQuotesService) SetClQuoteId(clQuoteId string) *GetRfqQuotesService {
	s.params["clQuoteId"] = clQuoteId
	return s
}

// SetState filters by quote state.
func (s *GetRfqQuotesService) SetState(state RfqQuoteState) *GetRfqQuotesService {
	s.params["state"] = string(state)
	return s
}

// SetBeginId pages from a starting quote id (records newer than the id).
func (s *GetRfqQuotesService) SetBeginId(beginId string) *GetRfqQuotesService {
	s.params["beginId"] = beginId
	return s
}

// SetEndId pages up to an ending quote id (records older than the id).
func (s *GetRfqQuotesService) SetEndId(endId string) *GetRfqQuotesService {
	s.params["endId"] = endId
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetRfqQuotesService) SetLimit(limit int) *GetRfqQuotesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetRfqQuotesService) Do(ctx context.Context) ([]RfqQuote, error) {
	req := request.Get(ctx, s.c, "/api/v5/rfq/quotes", s.params).WithSign()
	return request.DoList[RfqQuote](req)
}

// RfqQuote is one quote record. The validating account has no quote history, so
// the field set is modeled from the OKX doc field table.
type RfqQuote struct {
	CreateTime    time.Time     `json:"cTime"`
	UpdateTime    time.Time     `json:"uTime"`
	TraderCode    string        `json:"traderCode"`
	RFQID         string        `json:"rfqId"`
	ClientRFQID   string        `json:"clRfqId"`
	QuoteID       string        `json:"quoteId"`
	ClientQuoteID string        `json:"clQuoteId"`
	State         RfqQuoteState `json:"state"`
	ValidUntil    time.Time     `json:"validUntil"`
	QuoteSide     RfqQuoteSide  `json:"quoteSide"`
	Legs          []RfqQuoteLeg `json:"legs"`
}

// RfqQuoteLeg is one leg of a quote.
type RfqQuoteLeg struct {
	InstrumentID   string          `json:"instId"`
	Size           decimal.Decimal `json:"sz"`
	Price          decimal.Decimal `json:"px"`
	Side           Side            `json:"side"`
	TargetCurrency TgtCcy          `json:"tgtCcy"`
	PositionSide   PosSide         `json:"posSide"`
}

// GetRfqTradesService -- GET /api/v5/rfq/trades (Read)
//
// Returns the account's executed block trades (RFQ trades).
type GetRfqTradesService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetRfqTradesService() *GetRfqTradesService {
	return &GetRfqTradesService{c: c, params: map[string]string{}}
}

// SetRfqId filters by RFQ id.
func (s *GetRfqTradesService) SetRfqId(rfqId string) *GetRfqTradesService {
	s.params["rfqId"] = rfqId
	return s
}

// SetClRfqId filters by client-supplied RFQ id.
func (s *GetRfqTradesService) SetClRfqId(clRfqId string) *GetRfqTradesService {
	s.params["clRfqId"] = clRfqId
	return s
}

// SetQuoteId filters by quote id.
func (s *GetRfqTradesService) SetQuoteId(quoteId string) *GetRfqTradesService {
	s.params["quoteId"] = quoteId
	return s
}

// SetClQuoteId filters by client-supplied quote id.
func (s *GetRfqTradesService) SetClQuoteId(clQuoteId string) *GetRfqTradesService {
	s.params["clQuoteId"] = clQuoteId
	return s
}

// SetState filters by trade state.
func (s *GetRfqTradesService) SetState(state string) *GetRfqTradesService {
	s.params["state"] = state
	return s
}

// SetBeginId pages from a starting block-trade id (records newer than the id).
func (s *GetRfqTradesService) SetBeginId(beginId string) *GetRfqTradesService {
	s.params["beginId"] = beginId
	return s
}

// SetEndId pages up to an ending block-trade id (records older than the id).
func (s *GetRfqTradesService) SetEndId(endId string) *GetRfqTradesService {
	s.params["endId"] = endId
	return s
}

// SetBeginTs filters to trades at or after the given time (ms).
func (s *GetRfqTradesService) SetBeginTs(t time.Time) *GetRfqTradesService {
	s.params["beginTs"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTs filters to trades at or before the given time (ms).
func (s *GetRfqTradesService) SetEndTs(t time.Time) *GetRfqTradesService {
	s.params["endTs"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetRfqTradesService) SetLimit(limit int) *GetRfqTradesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetRfqTradesService) Do(ctx context.Context) ([]RfqTrade, error) {
	req := request.Get(ctx, s.c, "/api/v5/rfq/trades", s.params).WithSign()
	return request.DoList[RfqTrade](req)
}

// RfqTrade is one executed block trade. The validating account has no block-trade
// history, so the field set is modeled from the OKX doc field table.
type RfqTrade struct {
	CreateTime      time.Time     `json:"cTime"`
	RFQID           string        `json:"rfqId"`
	ClientRFQID     string        `json:"clRfqId"`
	QuoteID         string        `json:"quoteId"`
	ClientQuoteID   string        `json:"clQuoteId"`
	BlockTradeID    string        `json:"blockTdId"`
	TakerTraderCode string        `json:"tTraderCode"`
	MakerTraderCode string        `json:"mTraderCode"`
	Tag             string        `json:"tag"`
	Legs            []RfqTradeLeg `json:"legs"`
}

// RfqTradeLeg is one leg of an executed block trade.
type RfqTradeLeg struct {
	InstrumentID   string          `json:"instId"`
	Side           Side            `json:"side"`
	Size           decimal.Decimal `json:"sz"`
	Price          decimal.Decimal `json:"px"`
	TradeID        string          `json:"tradeId"`
	Fee            decimal.Decimal `json:"fee"`
	FeeCurrency    string          `json:"feeCcy"`
	TargetCurrency TgtCcy          `json:"tgtCcy"`
	PositionSide   PosSide         `json:"posSide"`
}

// GetRfqPublicTradesService -- GET /api/v5/rfq/public-trades (public)
//
// Returns the most recent public block trades across all counterparties.
type GetRfqPublicTradesService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetRfqPublicTradesService() *GetRfqPublicTradesService {
	return &GetRfqPublicTradesService{c: c, params: map[string]string{}}
}

// SetBeginId pages from a starting block-trade id (records newer than the id).
func (s *GetRfqPublicTradesService) SetBeginId(beginId string) *GetRfqPublicTradesService {
	s.params["beginId"] = beginId
	return s
}

// SetEndId pages up to an ending block-trade id (records older than the id).
func (s *GetRfqPublicTradesService) SetEndId(endId string) *GetRfqPublicTradesService {
	s.params["endId"] = endId
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetRfqPublicTradesService) SetLimit(limit int) *GetRfqPublicTradesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetRfqPublicTradesService) Do(ctx context.Context) ([]RfqPublicTrade, error) {
	req := request.Get(ctx, s.c, "/api/v5/rfq/public-trades", s.params)
	return request.DoList[RfqPublicTrade](req)
}

// RfqPublicTrade is one public block trade.
type RfqPublicTrade struct {
	BlockTradeID string              `json:"blockTdId"`
	GroupID      string              `json:"groupId"`
	Strategy     string              `json:"strategy"`
	Inverse      bool                `json:"inverse"`
	CreateTime   time.Time           `json:"cTime"`
	Legs         []RfqPublicTradeLeg `json:"legs"`
}

// RfqPublicTradeLeg is one leg of a public block trade.
type RfqPublicTradeLeg struct {
	InstrumentID string          `json:"instId"`
	Side         Side            `json:"side"`
	Size         decimal.Decimal `json:"sz"`
	Price        decimal.Decimal `json:"px"`
	TradeID      string          `json:"tradeId"`
}

// GetRfqBlockTickersService -- GET /api/v5/market/block-tickers (public)
//
// Returns the 24h block-trade volume tickers for every instrument of a product
// line. NOTE: OKX documents this under /api/v5/rfq/block-tickers but that path
// returns HTTP 404; the live endpoint is /api/v5/market/block-tickers.
type GetRfqBlockTickersService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetRfqBlockTickersService(instType InstType) *GetRfqBlockTickersService {
	return &GetRfqBlockTickersService{c: c, params: map[string]string{"instType": string(instType)}}
}

// SetUly filters by underlying (FUTURES/SWAP/OPTION).
func (s *GetRfqBlockTickersService) SetUly(uly string) *GetRfqBlockTickersService {
	s.params["uly"] = uly
	return s
}

// SetInstFamily filters by instrument family (FUTURES/SWAP/OPTION).
func (s *GetRfqBlockTickersService) SetInstFamily(instFamily string) *GetRfqBlockTickersService {
	s.params["instFamily"] = instFamily
	return s
}

func (s *GetRfqBlockTickersService) Do(ctx context.Context) ([]RfqBlockTicker, error) {
	req := request.Get(ctx, s.c, "/api/v5/market/block-tickers", s.params)
	return request.DoList[RfqBlockTicker](req)
}

// RfqBlockTicker is the 24h block-trade volume snapshot for an instrument.
type RfqBlockTicker struct {
	InstrumentType    InstType        `json:"instType"`
	InstrumentID      string          `json:"instId"`
	VolumeCurrency24h decimal.Decimal `json:"volCcy24h"`
	Volume24h         decimal.Decimal `json:"vol24h"`
	Timestamp         time.Time       `json:"ts"`
}

// GetRfqBlockTickerService -- GET /api/v5/market/block-ticker (public)
//
// Returns the 24h block-trade volume ticker for a single instrument. NOTE: OKX
// documents this under /api/v5/rfq/block-ticker but that path returns HTTP 404;
// the live endpoint is /api/v5/market/block-ticker.
type GetRfqBlockTickerService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetRfqBlockTickerService(instId string) *GetRfqBlockTickerService {
	return &GetRfqBlockTickerService{c: c, params: map[string]string{"instId": instId}}
}

func (s *GetRfqBlockTickerService) Do(ctx context.Context) (*RfqBlockTicker, error) {
	req := request.Get(ctx, s.c, "/api/v5/market/block-ticker", s.params)
	return request.DoOne[RfqBlockTicker](req)
}

// GetRfqMakerInstrumentSettingsService -- GET /api/v5/rfq/maker-instrument-settings (Read)
//
// Returns the maker's per-instrument-type quoting settings (the instruments a
// maker is configured to quote and their parameters).
type GetRfqMakerInstrumentSettingsService struct {
	c *Client
}

func (c *Client) NewGetRfqMakerInstrumentSettingsService() *GetRfqMakerInstrumentSettingsService {
	return &GetRfqMakerInstrumentSettingsService{c: c}
}

func (s *GetRfqMakerInstrumentSettingsService) Do(ctx context.Context) ([]RfqMakerInstrumentSetting, error) {
	req := request.Get(ctx, s.c, "/api/v5/rfq/maker-instrument-settings").WithSign()
	return request.DoList[RfqMakerInstrumentSetting](req)
}

// RfqMakerInstrumentSetting is the maker's quoting configuration for one
// instrument type. The validating account lacks maker permission (50030), so the
// field set is modeled from the OKX doc field table.
type RfqMakerInstrumentSetting struct {
	InstrumentType InstType                        `json:"instType"`
	IncludeALL     bool                            `json:"includeALL"`
	Data           []RfqMakerInstrumentSettingData `json:"data"`
}

// RfqMakerInstrumentSettingData is one instrument-family/underlying entry within
// a maker's quoting configuration.
type RfqMakerInstrumentSettingData struct {
	InstrumentFamily string          `json:"instFamily"`
	InstrumentID     string          `json:"instId"`
	MaxBlockSize     decimal.Decimal `json:"maxBlockSz"`
	MakerPriceBand   decimal.Decimal `json:"makerPxBand"`
}

// GetRfqMmpConfigService -- GET /api/v5/rfq/mmp-config (Read)
//
// Returns the maker's MMP (market-maker protection) configuration.
type GetRfqMmpConfigService struct {
	c *Client
}

func (c *Client) NewGetRfqMmpConfigService() *GetRfqMmpConfigService {
	return &GetRfqMmpConfigService{c: c}
}

func (s *GetRfqMmpConfigService) Do(ctx context.Context) (*RfqMmpConfig, error) {
	req := request.Get(ctx, s.c, "/api/v5/rfq/mmp-config").WithSign()
	return request.DoOne[RfqMmpConfig](req)
}

// RfqMmpConfig is the maker's MMP configuration. The validating account lacks
// maker permission (50030), so the field set is modeled from the OKX doc field
// table.
type RfqMmpConfig struct {
	TimeInterval   decimal.Decimal `json:"timeInterval"`
	FrozenInterval decimal.Decimal `json:"frozenInterval"`
	CountLimit     decimal.Decimal `json:"countLimit"`
	MMPFrozen      bool            `json:"mmpFrozen"`
	MMPFrozenUntil time.Time       `json:"mmpFrozenUntil"`
}

// --- State-changing endpoints (Trade): implemented but NEVER exercised by the
// test suite. Bodies and acks are modeled from the OKX docs. ---

// CreateRfqService -- POST /api/v5/rfq/create-rfq (Trade)
//
// Creates an RFQ as a taker and broadcasts it to the chosen counterparties.
type CreateRfqService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewCreateRfqService(counterparties []string, legs []RfqCreateLeg) *CreateRfqService {
	return &CreateRfqService{c: c, body: map[string]any{
		"counterparties": counterparties,
		"legs":           legs,
	}}
}

// SetClRfqId sets a client-supplied RFQ id.
func (s *CreateRfqService) SetClRfqId(clRfqId string) *CreateRfqService {
	s.body["clRfqId"] = clRfqId
	return s
}

// SetTag sets an order tag for the RFQ.
func (s *CreateRfqService) SetTag(tag string) *CreateRfqService {
	s.body["tag"] = tag
	return s
}

// SetAllowPartialExecution toggles whether partial execution is allowed.
func (s *CreateRfqService) SetAllowPartialExecution(allow bool) *CreateRfqService {
	s.body["allowPartialExecution"] = allow
	return s
}

func (s *CreateRfqService) Do(ctx context.Context) (*Rfq, error) {
	req := request.Post(ctx, s.c, "/api/v5/rfq/create-rfq", s.body).WithSign()
	return request.DoOne[Rfq](req)
}

// RfqCreateLeg is one leg supplied when creating an RFQ.
type RfqCreateLeg struct {
	InstrumentID   string          `json:"instId"`
	Size           decimal.Decimal `json:"sz"`
	Side           Side            `json:"side"`
	TargetCurrency TgtCcy          `json:"tgtCcy,omitempty"`
	PositionSide   PosSide         `json:"posSide,omitempty"`
}

// CancelRfqService -- POST /api/v5/rfq/cancel-rfq (Trade)
//
// Cancels a single active RFQ.
type CancelRfqService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewCancelRfqService() *CancelRfqService {
	return &CancelRfqService{c: c, body: map[string]any{}}
}

// SetRfqId targets the RFQ by id (one of rfqId / clRfqId is required).
func (s *CancelRfqService) SetRfqId(rfqId string) *CancelRfqService {
	s.body["rfqId"] = rfqId
	return s
}

// SetClRfqId targets the RFQ by client-supplied id.
func (s *CancelRfqService) SetClRfqId(clRfqId string) *CancelRfqService {
	s.body["clRfqId"] = clRfqId
	return s
}

func (s *CancelRfqService) Do(ctx context.Context) (*RfqCancelAck, error) {
	req := request.Post(ctx, s.c, "/api/v5/rfq/cancel-rfq", s.body).WithSign()
	return request.DoOne[RfqCancelAck](req)
}

// RfqCancelAck is the per-item ack of an RFQ cancel.
type RfqCancelAck struct {
	RFQID       string `json:"rfqId"`
	ClientRFQID string `json:"clRfqId"`
	SCode       string `json:"sCode"`
	SMsg        string `json:"sMsg"`
}

// CancelBatchRfqsService -- POST /api/v5/rfq/cancel-batch-rfqs (Trade)
//
// Cancels multiple active RFQs in one request.
type CancelBatchRfqsService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewCancelBatchRfqsService() *CancelBatchRfqsService {
	return &CancelBatchRfqsService{c: c, body: map[string]any{}}
}

// SetRfqIds targets the RFQs by id.
func (s *CancelBatchRfqsService) SetRfqIds(rfqIds []string) *CancelBatchRfqsService {
	s.body["rfqIds"] = rfqIds
	return s
}

// SetClRfqIds targets the RFQs by client-supplied id.
func (s *CancelBatchRfqsService) SetClRfqIds(clRfqIds []string) *CancelBatchRfqsService {
	s.body["clRfqIds"] = clRfqIds
	return s
}

func (s *CancelBatchRfqsService) Do(ctx context.Context) ([]RfqCancelAck, error) {
	req := request.Post(ctx, s.c, "/api/v5/rfq/cancel-batch-rfqs", s.body).WithSign()
	return request.DoListPartial[RfqCancelAck](req)
}

// CancelAllRfqsService -- POST /api/v5/rfq/cancel-all-rfqs (Trade)
//
// Cancels all of the account's active RFQs.
type CancelAllRfqsService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewCancelAllRfqsService() *CancelAllRfqsService {
	return &CancelAllRfqsService{c: c, body: map[string]any{}}
}

func (s *CancelAllRfqsService) Do(ctx context.Context) (*RfqTimestampAck, error) {
	req := request.Post(ctx, s.c, "/api/v5/rfq/cancel-all-rfqs", s.body).WithSign()
	return request.DoOne[RfqTimestampAck](req)
}

// RfqTimestampAck is the ack of a bulk cancel/reset action, carrying the server
// timestamp at which the action took effect.
type RfqTimestampAck struct {
	Timestamp time.Time `json:"ts"`
}

// ExecuteQuoteService -- POST /api/v5/rfq/execute-quote (Trade)
//
// Executes (accepts) a maker's quote as the taker, producing a block trade.
type ExecuteQuoteService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewExecuteQuoteService(rfqId, quoteId string) *ExecuteQuoteService {
	return &ExecuteQuoteService{c: c, body: map[string]any{
		"rfqId":   rfqId,
		"quoteId": quoteId,
	}}
}

// SetLegs restricts execution to a subset of legs (partial execution).
func (s *ExecuteQuoteService) SetLegs(legs []RfqExecuteLeg) *ExecuteQuoteService {
	s.body["legs"] = legs
	return s
}

func (s *ExecuteQuoteService) Do(ctx context.Context) (*RfqTrade, error) {
	req := request.Post(ctx, s.c, "/api/v5/rfq/execute-quote", s.body).WithSign()
	return request.DoOne[RfqTrade](req)
}

// RfqExecuteLeg is one leg supplied when executing a quote.
type RfqExecuteLeg struct {
	InstrumentID string          `json:"instId"`
	Size         decimal.Decimal `json:"sz"`
}

// CreateRfqQuoteService -- POST /api/v5/rfq/create-quote (Trade)
//
// Creates a quote against an RFQ as a maker.
type CreateRfqQuoteService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewCreateRfqQuoteService(rfqId string, quoteSide RfqQuoteSide, legs []RfqQuoteCreateLeg) *CreateRfqQuoteService {
	return &CreateRfqQuoteService{c: c, body: map[string]any{
		"rfqId":     rfqId,
		"quoteSide": string(quoteSide),
		"legs":      legs,
	}}
}

// SetClQuoteId sets a client-supplied quote id.
func (s *CreateRfqQuoteService) SetClQuoteId(clQuoteId string) *CreateRfqQuoteService {
	s.body["clQuoteId"] = clQuoteId
	return s
}

// SetTag sets an order tag for the quote.
func (s *CreateRfqQuoteService) SetTag(tag string) *CreateRfqQuoteService {
	s.body["tag"] = tag
	return s
}

// SetExpiresIn sets the quote validity window in seconds.
func (s *CreateRfqQuoteService) SetExpiresIn(seconds int) *CreateRfqQuoteService {
	s.body["expiresIn"] = strconv.Itoa(seconds)
	return s
}

// SetAnonymous toggles whether the quote is anonymous.
func (s *CreateRfqQuoteService) SetAnonymous(anonymous bool) *CreateRfqQuoteService {
	s.body["anonymous"] = anonymous
	return s
}

func (s *CreateRfqQuoteService) Do(ctx context.Context) (*RfqQuote, error) {
	req := request.Post(ctx, s.c, "/api/v5/rfq/create-quote", s.body).WithSign()
	return request.DoOne[RfqQuote](req)
}

// RfqQuoteCreateLeg is one leg supplied when creating a quote.
type RfqQuoteCreateLeg struct {
	InstrumentID   string          `json:"instId"`
	Size           decimal.Decimal `json:"sz"`
	Price          decimal.Decimal `json:"px"`
	Side           Side            `json:"side"`
	TargetCurrency TgtCcy          `json:"tgtCcy,omitempty"`
	PositionSide   PosSide         `json:"posSide,omitempty"`
}

// CancelRfqQuoteService -- POST /api/v5/rfq/cancel-quote (Trade)
//
// Cancels a single active quote.
type CancelRfqQuoteService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewCancelRfqQuoteService() *CancelRfqQuoteService {
	return &CancelRfqQuoteService{c: c, body: map[string]any{}}
}

// SetQuoteId targets the quote by id (one of quoteId / clQuoteId is required).
func (s *CancelRfqQuoteService) SetQuoteId(quoteId string) *CancelRfqQuoteService {
	s.body["quoteId"] = quoteId
	return s
}

// SetClQuoteId targets the quote by client-supplied id.
func (s *CancelRfqQuoteService) SetClQuoteId(clQuoteId string) *CancelRfqQuoteService {
	s.body["clQuoteId"] = clQuoteId
	return s
}

// SetRfqId restricts cancellation to quotes on the given RFQ.
func (s *CancelRfqQuoteService) SetRfqId(rfqId string) *CancelRfqQuoteService {
	s.body["rfqId"] = rfqId
	return s
}

// SetClRfqId restricts cancellation to quotes on the given client RFQ id.
func (s *CancelRfqQuoteService) SetClRfqId(clRfqId string) *CancelRfqQuoteService {
	s.body["clRfqId"] = clRfqId
	return s
}

func (s *CancelRfqQuoteService) Do(ctx context.Context) (*RfqQuoteCancelAck, error) {
	req := request.Post(ctx, s.c, "/api/v5/rfq/cancel-quote", s.body).WithSign()
	return request.DoOne[RfqQuoteCancelAck](req)
}

// RfqQuoteCancelAck is the per-item ack of a quote cancel.
type RfqQuoteCancelAck struct {
	QuoteID       string `json:"quoteId"`
	ClientQuoteID string `json:"clQuoteId"`
	SCode         string `json:"sCode"`
	SMsg          string `json:"sMsg"`
}

// CancelBatchRfqQuotesService -- POST /api/v5/rfq/cancel-batch-quotes (Trade)
//
// Cancels multiple active quotes in one request.
type CancelBatchRfqQuotesService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewCancelBatchRfqQuotesService() *CancelBatchRfqQuotesService {
	return &CancelBatchRfqQuotesService{c: c, body: map[string]any{}}
}

// SetQuoteIds targets the quotes by id.
func (s *CancelBatchRfqQuotesService) SetQuoteIds(quoteIds []string) *CancelBatchRfqQuotesService {
	s.body["quoteIds"] = quoteIds
	return s
}

// SetClQuoteIds targets the quotes by client-supplied id.
func (s *CancelBatchRfqQuotesService) SetClQuoteIds(clQuoteIds []string) *CancelBatchRfqQuotesService {
	s.body["clQuoteIds"] = clQuoteIds
	return s
}

func (s *CancelBatchRfqQuotesService) Do(ctx context.Context) ([]RfqQuoteCancelAck, error) {
	req := request.Post(ctx, s.c, "/api/v5/rfq/cancel-batch-quotes", s.body).WithSign()
	return request.DoListPartial[RfqQuoteCancelAck](req)
}

// CancelAllRfqQuotesService -- POST /api/v5/rfq/cancel-all-quotes (Trade)
//
// Cancels all of the account's active quotes.
type CancelAllRfqQuotesService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewCancelAllRfqQuotesService() *CancelAllRfqQuotesService {
	return &CancelAllRfqQuotesService{c: c, body: map[string]any{}}
}

func (s *CancelAllRfqQuotesService) Do(ctx context.Context) (*RfqTimestampAck, error) {
	req := request.Post(ctx, s.c, "/api/v5/rfq/cancel-all-quotes", s.body).WithSign()
	return request.DoOne[RfqTimestampAck](req)
}

// CancelAllAfterRfqService -- POST /api/v5/rfq/cancel-all-after (Trade)
//
// Arms a dead-man's-switch: cancels all quotes after the given number of seconds
// unless re-armed (0 disarms).
type CancelAllAfterRfqService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewCancelAllAfterRfqService(timeOut int) *CancelAllAfterRfqService {
	return &CancelAllAfterRfqService{c: c, body: map[string]any{
		"timeOut": strconv.Itoa(timeOut),
	}}
}

func (s *CancelAllAfterRfqService) Do(ctx context.Context) (*RfqCancelAllAfter, error) {
	req := request.Post(ctx, s.c, "/api/v5/rfq/cancel-all-after", s.body).WithSign()
	return request.DoOne[RfqCancelAllAfter](req)
}

// RfqCancelAllAfter is the ack of a cancel-all-after action.
type RfqCancelAllAfter struct {
	TriggerTime time.Time `json:"triggerTime"`
	Timestamp   time.Time `json:"ts"`
}

// SetRfqMmpService -- POST /api/v5/rfq/set-mmp (Trade)
//
// Sets the maker's MMP (market-maker protection) configuration.
type SetRfqMmpService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewSetRfqMmpService(timeInterval, frozenInterval, countLimit int) *SetRfqMmpService {
	return &SetRfqMmpService{c: c, body: map[string]any{
		"timeInterval":   strconv.Itoa(timeInterval),
		"frozenInterval": strconv.Itoa(frozenInterval),
		"countLimit":     strconv.Itoa(countLimit),
	}}
}

func (s *SetRfqMmpService) Do(ctx context.Context) (*RfqMmpConfig, error) {
	req := request.Post(ctx, s.c, "/api/v5/rfq/set-mmp", s.body).WithSign()
	return request.DoOne[RfqMmpConfig](req)
}

// ResetRfqMmpService -- POST /api/v5/rfq/mmp-reset (Trade)
//
// Resets (unfreezes) the maker's MMP state.
type ResetRfqMmpService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewResetRfqMmpService() *ResetRfqMmpService {
	return &ResetRfqMmpService{c: c, body: map[string]any{}}
}

func (s *ResetRfqMmpService) Do(ctx context.Context) (*RfqMmpReset, error) {
	req := request.Post(ctx, s.c, "/api/v5/rfq/mmp-reset", s.body).WithSign()
	return request.DoOne[RfqMmpReset](req)
}

// RfqMmpReset is the ack of an MMP reset.
type RfqMmpReset struct {
	Result bool `json:"result"`
}

// SetRfqMakerInstrumentSettingsService -- POST /api/v5/rfq/maker-instrument-settings (Trade)
//
// Sets the maker's per-instrument-type quoting settings.
type SetRfqMakerInstrumentSettingsService struct {
	c    *Client
	body []RfqMakerInstrumentSetting
}

func (c *Client) NewSetRfqMakerInstrumentSettingsService(settings []RfqMakerInstrumentSetting) *SetRfqMakerInstrumentSettingsService {
	return &SetRfqMakerInstrumentSettingsService{c: c, body: settings}
}

func (s *SetRfqMakerInstrumentSettingsService) Do(ctx context.Context) (*RfqMmpReset, error) {
	req := request.Post(ctx, s.c, "/api/v5/rfq/maker-instrument-settings").SetBody(s.body).WithSign()
	return request.DoOne[RfqMmpReset](req)
}

// GetRfqQuoteProductsService -- GET /api/v5/rfq/quote-products (deprecated)
//
// Returns the maker's configured quote products. DEPRECATED: this legacy path now
// returns HTTP 404; OKX folded quote-product configuration into
// /api/v5/rfq/maker-instrument-settings. Kept implement-only and NOT exercised by
// the test suite. Prefer GetRfqMakerInstrumentSettingsService.
type GetRfqQuoteProductsService struct {
	c *Client
}

func (c *Client) NewGetRfqQuoteProductsService() *GetRfqQuoteProductsService {
	return &GetRfqQuoteProductsService{c: c}
}

func (s *GetRfqQuoteProductsService) Do(ctx context.Context) ([]RfqQuoteProduct, error) {
	req := request.Get(ctx, s.c, "/api/v5/rfq/quote-products").WithSign()
	return request.DoList[RfqQuoteProduct](req)
}

// RfqQuoteProduct is one quote-product configuration. The validating account
// lacks maker permission, so the field set is modeled from the OKX doc field
// table.
type RfqQuoteProduct struct {
	InstrumentType InstType              `json:"instType"`
	IncludeALL     bool                  `json:"includeALL"`
	Data           []RfqQuoteProductData `json:"data"`
}

// RfqQuoteProductData is one instrument entry within a quote product.
type RfqQuoteProductData struct {
	Underlying       string          `json:"underlying"`
	InstrumentFamily string          `json:"instFamily"`
	InstrumentID     string          `json:"instId"`
	MaxBlockSize     decimal.Decimal `json:"maxBlockSz"`
	MakerPriceBand   decimal.Decimal `json:"makerPxBand"`
}

// SetRfqQuoteProductsService -- POST /api/v5/rfq/set-quote-products (Trade)
//
// Sets the maker's quote products.
type SetRfqQuoteProductsService struct {
	c    *Client
	body []RfqQuoteProduct
}

func (c *Client) NewSetRfqQuoteProductsService(products []RfqQuoteProduct) *SetRfqQuoteProductsService {
	return &SetRfqQuoteProductsService{c: c, body: products}
}

func (s *SetRfqQuoteProductsService) Do(ctx context.Context) (*RfqMmpReset, error) {
	req := request.Post(ctx, s.c, "/api/v5/rfq/set-quote-products").SetBody(s.body).WithSign()
	return request.DoOne[RfqMmpReset](req)
}
