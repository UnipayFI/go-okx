package okx

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-okx/request"
	"github.com/shopspring/decimal"
)

// convertState is the fill state of a convert trade reported by the convert
// history endpoint.
type convertState string

const (
	convertStateFullyFilled convertState = "fullyFilled"
	convertStateRejected    convertState = "rejected"
)

// GetConvertCurrenciesService -- GET /api/v5/asset/convert/currencies (Read)
//
// Returns the list of currencies that can be used in small-asset conversions.
type GetConvertCurrenciesService struct {
	c *Client
}

func (c *Client) NewGetConvertCurrenciesService() *GetConvertCurrenciesService {
	return &GetConvertCurrenciesService{c: c}
}

func (s *GetConvertCurrenciesService) Do(ctx context.Context) ([]ConvertCurrency, error) {
	req := request.Get(ctx, s.c, "/api/v5/asset/convert/currencies").WithSign()
	return request.DoList[ConvertCurrency](req)
}

// ConvertCurrency is a single convertible currency and its conversion bounds.
type ConvertCurrency struct {
	Currency string          `json:"ccy"`
	Min      decimal.Decimal `json:"min"`
	Max      decimal.Decimal `json:"max"`
}

// GetConvertCurrencyPairService -- GET /api/v5/asset/convert/currency-pair (Read)
//
// Returns the tradable convert pair and its per-currency size bounds for a
// from/to currency combination.
type GetConvertCurrencyPairService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetConvertCurrencyPairService(fromCcy, toCcy string) *GetConvertCurrencyPairService {
	return &GetConvertCurrencyPairService{c: c, params: map[string]string{
		"fromCcy": fromCcy,
		"toCcy":   toCcy,
	}}
}

func (s *GetConvertCurrencyPairService) Do(ctx context.Context) (*ConvertCurrencyPair, error) {
	req := request.Get(ctx, s.c, "/api/v5/asset/convert/currency-pair", s.params).WithSign()
	return request.DoOne[ConvertCurrencyPair](req)
}

// ConvertCurrencyPair is a convert pair and its base/quote size limits.
type ConvertCurrencyPair struct {
	InstrumentID     string          `json:"instId"`
	BaseCurrency     string          `json:"baseCcy"`
	BaseCurrencyMax  decimal.Decimal `json:"baseCcyMax"`
	BaseCurrencyMin  decimal.Decimal `json:"baseCcyMin"`
	QuoteCurrency    string          `json:"quoteCcy"`
	QuoteCurrencyMax decimal.Decimal `json:"quoteCcyMax"`
	QuoteCurrencyMin decimal.Decimal `json:"quoteCcyMin"`
}

// ConvertEstimateQuoteService -- POST /api/v5/asset/convert/estimate-quote (Trade)
//
// Requests a binding conversion quote (returns a quoteId valid for a short
// window) for a base/quote currency and side.
//
// State-changing: NOT exercised by the test suite.
type ConvertEstimateQuoteService struct {
	c    *Client
	body map[string]any
}

// NewConvertEstimateQuoteService builds an estimate-quote request. side is the
// trade direction relative to baseCcy; rfqSz is the requested size denominated
// in rfqSzCcy (which must be either baseCcy or quoteCcy).
func (c *Client) NewConvertEstimateQuoteService(baseCcy, quoteCcy string, side Side, rfqSz decimal.Decimal, rfqSzCcy string) *ConvertEstimateQuoteService {
	return &ConvertEstimateQuoteService{c: c, body: map[string]any{
		"baseCcy":  baseCcy,
		"quoteCcy": quoteCcy,
		"side":     string(side),
		"rfqSz":    rfqSz.String(),
		"rfqSzCcy": rfqSzCcy,
	}}
}

// SetClQReqId sets the client-supplied quote-request id (used for idempotency).
func (s *ConvertEstimateQuoteService) SetClQReqId(clQReqId string) *ConvertEstimateQuoteService {
	s.body["clQReqId"] = clQReqId
	return s
}

// SetTag sets the order tag (broker id).
func (s *ConvertEstimateQuoteService) SetTag(tag string) *ConvertEstimateQuoteService {
	s.body["tag"] = tag
	return s
}

func (s *ConvertEstimateQuoteService) Do(ctx context.Context) (*ConvertEstimateQuote, error) {
	req := request.Post(ctx, s.c, "/api/v5/asset/convert/estimate-quote", s.body).WithSign()
	return request.DoOne[ConvertEstimateQuote](req)
}

// ConvertEstimateQuote is the binding quote returned by estimate-quote.
type ConvertEstimateQuote struct {
	QuoteID              string          `json:"quoteId"`
	ClientQuoteRequestID string          `json:"clQReqId"`
	BaseCurrency         string          `json:"baseCcy"`
	QuoteCurrency        string          `json:"quoteCcy"`
	Side                 Side            `json:"side"`
	OriginalRFQSize      decimal.Decimal `json:"origRfqSz"`
	RFQSize              decimal.Decimal `json:"rfqSz"`
	RFQSizeCurrency      string          `json:"rfqSzCcy"`
	ConvertedPrice       decimal.Decimal `json:"cnvtPx"`
	BaseSize             decimal.Decimal `json:"baseSz"`
	QuoteSize            decimal.Decimal `json:"quoteSz"`
	TTLMilliseconds      decimal.Decimal `json:"ttlMs"`
	QuoteTime            time.Time       `json:"quoteTime"`
}

// ConvertTradeService -- POST /api/v5/asset/convert/trade (Trade)
//
// Executes a conversion against a quoteId previously obtained from
// estimate-quote.
//
// State-changing: NOT exercised by the test suite.
type ConvertTradeService struct {
	c    *Client
	body map[string]any
}

// NewConvertTradeService builds a convert-trade request. quoteId is the id from
// estimate-quote; sz is the size denominated in szCcy (which must be either
// baseCcy or quoteCcy).
func (c *Client) NewConvertTradeService(quoteId, baseCcy, quoteCcy string, side Side, sz decimal.Decimal, szCcy string) *ConvertTradeService {
	return &ConvertTradeService{c: c, body: map[string]any{
		"quoteId":  quoteId,
		"baseCcy":  baseCcy,
		"quoteCcy": quoteCcy,
		"side":     string(side),
		"sz":       sz.String(),
		"szCcy":    szCcy,
	}}
}

// SetClTReqId sets the client-supplied trade-request id (used for idempotency).
func (s *ConvertTradeService) SetClTReqId(clTReqId string) *ConvertTradeService {
	s.body["clTReqId"] = clTReqId
	return s
}

// SetTag sets the order tag (broker id).
func (s *ConvertTradeService) SetTag(tag string) *ConvertTradeService {
	s.body["tag"] = tag
	return s
}

func (s *ConvertTradeService) Do(ctx context.Context) (*ConvertTrade, error) {
	req := request.Post(ctx, s.c, "/api/v5/asset/convert/trade", s.body).WithSign()
	return request.DoOne[ConvertTrade](req)
}

// ConvertTrade is the result of an executed conversion.
type ConvertTrade struct {
	TradeID              string          `json:"tradeId"`
	QuoteID              string          `json:"quoteId"`
	ClientTradeRequestID string          `json:"clTReqId"`
	State                convertState    `json:"state"`
	InstrumentID         string          `json:"instId"`
	BaseCurrency         string          `json:"baseCcy"`
	QuoteCurrency        string          `json:"quoteCcy"`
	Side                 Side            `json:"side"`
	FillPrice            decimal.Decimal `json:"fillPx"`
	FillBaseSize         decimal.Decimal `json:"fillBaseSz"`
	FillQuoteSize        decimal.Decimal `json:"fillQuoteSz"`
	Timestamp            time.Time       `json:"ts"`
}

// GetConvertHistoryService -- GET /api/v5/asset/convert/history (Read)
//
// Returns the account's past conversions (last three months).
type GetConvertHistoryService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetConvertHistoryService() *GetConvertHistoryService {
	return &GetConvertHistoryService{c: c, params: map[string]string{}}
}

// SetAfter paginates to records earlier than the given time (older).
func (s *GetConvertHistoryService) SetAfter(t time.Time) *GetConvertHistoryService {
	s.params["after"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetBefore paginates to records later than the given time (newer).
func (s *GetConvertHistoryService) SetBefore(t time.Time) *GetConvertHistoryService {
	s.params["before"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetConvertHistoryService) SetLimit(limit int) *GetConvertHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetTag filters by order tag (broker id).
func (s *GetConvertHistoryService) SetTag(tag string) *GetConvertHistoryService {
	s.params["tag"] = tag
	return s
}

func (s *GetConvertHistoryService) Do(ctx context.Context) ([]ConvertHistory, error) {
	req := request.Get(ctx, s.c, "/api/v5/asset/convert/history", s.params).WithSign()
	return request.DoList[ConvertHistory](req)
}

// ConvertHistory is a single past conversion record.
type ConvertHistory struct {
	InstrumentID  string          `json:"instId"`
	Side          Side            `json:"side"`
	FillPrice     decimal.Decimal `json:"fillPx"`
	BaseCurrency  string          `json:"baseCcy"`
	QuoteCurrency string          `json:"quoteCcy"`
	FillBaseSize  decimal.Decimal `json:"fillBaseSz"`
	FillQuoteSize decimal.Decimal `json:"fillQuoteSz"`
	State         convertState    `json:"state"`
	TradeID       string          `json:"tradeId"`
	Timestamp     time.Time       `json:"ts"`
}
