package okx

// Trading-statistics ("rubik") endpoints under /api/v5/rubik/stat/*. These are
// all PUBLIC (no signing). Most return an array of positional string arrays
// ([[ts, value, ...], ...]); each is modeled as a typed struct plus a row
// parser. The raw rows are fetched as [][]string via request.DoList[[]string]
// and mapped to the typed struct, mirroring how candle endpoints are handled.

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-okx/request"
	"github.com/shopspring/decimal"
)

// RubikStatPeriod is the bar/aggregation window accepted by the rubik
// trading-statistics endpoints (e.g. "5m", "1H", "1D"). Each endpoint documents
// the subset it supports; values are passed through verbatim.
type RubikStatPeriod string

const (
	RubikStatPeriod5m RubikStatPeriod = "5m"
	RubikStatPeriod1H RubikStatPeriod = "1H"
	RubikStatPeriod1D RubikStatPeriod = "1D"
)

// parseRubikTs converts a millisecond timestamp string (the first column of
// every rubik row) into a time.Time, returning the zero value on a blank or
// unparsable token (matching the global codec's sentinels).
func parseRubikTs(s string) time.Time {
	switch s {
	case "", "0", "-1":
		return time.Time{}
	}
	ms, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Time{}
	}
	return time.UnixMilli(ms)
}

// parseRubikDec converts a decimal string column into a decimal.Decimal,
// returning zero on a blank or unparsable token.
func parseRubikDec(s string) decimal.Decimal {
	if s == "" {
		return decimal.Zero
	}
	v, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.Zero
	}
	return v
}

// GetSupportCoinService -- GET /api/v5/rubik/stat/trading-data/support-coin (public)
//
// Returns the currencies supported by the trading-statistics endpoints, grouped
// by product line.
type GetSupportCoinService struct {
	c *Client
}

func (c *Client) NewGetSupportCoinService() *GetSupportCoinService {
	return &GetSupportCoinService{c: c}
}

func (s *GetSupportCoinService) Do(ctx context.Context) (*RubikSupportCoin, error) {
	req := request.Get(ctx, s.c, "/api/v5/rubik/stat/trading-data/support-coin")
	return request.DoObject[RubikSupportCoin](req)
}

// RubikSupportCoin lists the currencies covered by the trading-statistics
// endpoints, split by product line.
type RubikSupportCoin struct {
	Contract []string `json:"contract"`
	Option   []string `json:"option"`
	Spot     []string `json:"spot"`
}

// GetTakerVolumeService -- GET /api/v5/rubik/stat/taker-volume (public)
//
// Returns the taker buy/sell trading volume per time bar for a currency and
// product line.
type GetTakerVolumeService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetTakerVolumeService(ccy string, instType InstType) *GetTakerVolumeService {
	return &GetTakerVolumeService{c: c, params: map[string]string{
		"ccy":      ccy,
		"instType": string(instType),
	}}
}

func (s *GetTakerVolumeService) SetBegin(t time.Time) *GetTakerVolumeService {
	s.params["begin"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetTakerVolumeService) SetEnd(t time.Time) *GetTakerVolumeService {
	s.params["end"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetTakerVolumeService) SetPeriod(p RubikStatPeriod) *GetTakerVolumeService {
	s.params["period"] = string(p)
	return s
}

func (s *GetTakerVolumeService) Do(ctx context.Context) ([]RubikTakerVolume, error) {
	req := request.Get(ctx, s.c, "/api/v5/rubik/stat/taker-volume", s.params)
	rows, err := request.DoList[[]string](req)
	if err != nil {
		return nil, err
	}
	out := make([]RubikTakerVolume, 0, len(rows))
	for _, r := range rows {
		out = append(out, parseRubikTakerVolume(r))
	}
	return out, nil
}

// RubikTakerVolume is one taker-volume bar: [ts, sellVol, buyVol].
type RubikTakerVolume struct {
	Timestamp  time.Time       `json:"ts"`
	SellVolume decimal.Decimal `json:"sellVol"`
	BuyVolume  decimal.Decimal `json:"buyVol"`
}

func parseRubikTakerVolume(r []string) RubikTakerVolume {
	var v RubikTakerVolume
	if len(r) > 0 {
		v.Timestamp = parseRubikTs(r[0])
	}
	if len(r) > 1 {
		v.SellVolume = parseRubikDec(r[1])
	}
	if len(r) > 2 {
		v.BuyVolume = parseRubikDec(r[2])
	}
	return v
}

// GetLoanRatioService -- GET /api/v5/rubik/stat/margin/loan-ratio (public)
//
// Returns the margin lending ratio (margin long vs short) per time bar for a
// currency.
type GetLoanRatioService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetLoanRatioService(ccy string) *GetLoanRatioService {
	return &GetLoanRatioService{c: c, params: map[string]string{"ccy": ccy}}
}

func (s *GetLoanRatioService) SetBegin(t time.Time) *GetLoanRatioService {
	s.params["begin"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetLoanRatioService) SetEnd(t time.Time) *GetLoanRatioService {
	s.params["end"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetLoanRatioService) SetPeriod(p RubikStatPeriod) *GetLoanRatioService {
	s.params["period"] = string(p)
	return s
}

func (s *GetLoanRatioService) Do(ctx context.Context) ([]RubikRatio, error) {
	req := request.Get(ctx, s.c, "/api/v5/rubik/stat/margin/loan-ratio", s.params)
	return doRubikRatio(req)
}

// GetLongShortAccountRatioService -- GET /api/v5/rubik/stat/contracts/long-short-account-ratio (public)
//
// Returns the ratio of accounts net-long vs net-short on derivatives per time
// bar for a currency.
type GetLongShortAccountRatioService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetLongShortAccountRatioService(ccy string) *GetLongShortAccountRatioService {
	return &GetLongShortAccountRatioService{c: c, params: map[string]string{"ccy": ccy}}
}

func (s *GetLongShortAccountRatioService) SetBegin(t time.Time) *GetLongShortAccountRatioService {
	s.params["begin"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetLongShortAccountRatioService) SetEnd(t time.Time) *GetLongShortAccountRatioService {
	s.params["end"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetLongShortAccountRatioService) SetPeriod(p RubikStatPeriod) *GetLongShortAccountRatioService {
	s.params["period"] = string(p)
	return s
}

func (s *GetLongShortAccountRatioService) Do(ctx context.Context) ([]RubikRatio, error) {
	req := request.Get(ctx, s.c, "/api/v5/rubik/stat/contracts/long-short-account-ratio", s.params)
	return doRubikRatio(req)
}

// GetLongShortAccountRatioContractService -- GET /api/v5/rubik/stat/contracts/long-short-account-ratio-contract (public)
//
// Returns the long/short account ratio per time bar for a single contract
// instrument.
type GetLongShortAccountRatioContractService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetLongShortAccountRatioContractService(instId string) *GetLongShortAccountRatioContractService {
	return &GetLongShortAccountRatioContractService{c: c, params: map[string]string{"instId": instId}}
}

func (s *GetLongShortAccountRatioContractService) SetBegin(t time.Time) *GetLongShortAccountRatioContractService {
	s.params["begin"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetLongShortAccountRatioContractService) SetEnd(t time.Time) *GetLongShortAccountRatioContractService {
	s.params["end"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetLongShortAccountRatioContractService) SetPeriod(p RubikStatPeriod) *GetLongShortAccountRatioContractService {
	s.params["period"] = string(p)
	return s
}

func (s *GetLongShortAccountRatioContractService) Do(ctx context.Context) ([]RubikRatio, error) {
	req := request.Get(ctx, s.c, "/api/v5/rubik/stat/contracts/long-short-account-ratio-contract", s.params)
	return doRubikRatio(req)
}

// GetLongShortAccountRatioContractTopTraderService -- GET /api/v5/rubik/stat/contracts/long-short-account-ratio-contract-top-trader (public)
//
// Returns the long/short account ratio of top traders per time bar for a single
// contract instrument.
type GetLongShortAccountRatioContractTopTraderService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetLongShortAccountRatioContractTopTraderService(instId string) *GetLongShortAccountRatioContractTopTraderService {
	return &GetLongShortAccountRatioContractTopTraderService{c: c, params: map[string]string{"instId": instId}}
}

func (s *GetLongShortAccountRatioContractTopTraderService) SetBegin(t time.Time) *GetLongShortAccountRatioContractTopTraderService {
	s.params["begin"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetLongShortAccountRatioContractTopTraderService) SetEnd(t time.Time) *GetLongShortAccountRatioContractTopTraderService {
	s.params["end"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetLongShortAccountRatioContractTopTraderService) SetPeriod(p RubikStatPeriod) *GetLongShortAccountRatioContractTopTraderService {
	s.params["period"] = string(p)
	return s
}

func (s *GetLongShortAccountRatioContractTopTraderService) Do(ctx context.Context) ([]RubikRatio, error) {
	req := request.Get(ctx, s.c, "/api/v5/rubik/stat/contracts/long-short-account-ratio-contract-top-trader", s.params)
	return doRubikRatio(req)
}

// GetLongShortPositionRatioContractTopTraderService -- GET /api/v5/rubik/stat/contracts/long-short-position-ratio-contract-top-trader (public)
//
// Returns the long/short position ratio of top traders per time bar for a single
// contract instrument.
type GetLongShortPositionRatioContractTopTraderService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetLongShortPositionRatioContractTopTraderService(instId string) *GetLongShortPositionRatioContractTopTraderService {
	return &GetLongShortPositionRatioContractTopTraderService{c: c, params: map[string]string{"instId": instId}}
}

func (s *GetLongShortPositionRatioContractTopTraderService) SetBegin(t time.Time) *GetLongShortPositionRatioContractTopTraderService {
	s.params["begin"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetLongShortPositionRatioContractTopTraderService) SetEnd(t time.Time) *GetLongShortPositionRatioContractTopTraderService {
	s.params["end"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetLongShortPositionRatioContractTopTraderService) SetPeriod(p RubikStatPeriod) *GetLongShortPositionRatioContractTopTraderService {
	s.params["period"] = string(p)
	return s
}

func (s *GetLongShortPositionRatioContractTopTraderService) Do(ctx context.Context) ([]RubikRatio, error) {
	req := request.Get(ctx, s.c, "/api/v5/rubik/stat/contracts/long-short-position-ratio-contract-top-trader", s.params)
	return doRubikRatio(req)
}

// RubikRatio is one [ts, ratio] bar shared by the loan-ratio and the
// long/short account/position ratio endpoints.
type RubikRatio struct {
	Timestamp time.Time       `json:"ts"`
	Ratio     decimal.Decimal `json:"ratio"`
}

func parseRubikRatio(r []string) RubikRatio {
	var v RubikRatio
	if len(r) > 0 {
		v.Timestamp = parseRubikTs(r[0])
	}
	if len(r) > 1 {
		v.Ratio = parseRubikDec(r[1])
	}
	return v
}

// doRubikRatio runs a request whose rows are [ts, ratio] and maps them.
func doRubikRatio(req *request.Request) ([]RubikRatio, error) {
	rows, err := request.DoList[[]string](req)
	if err != nil {
		return nil, err
	}
	out := make([]RubikRatio, 0, len(rows))
	for _, r := range rows {
		out = append(out, parseRubikRatio(r))
	}
	return out, nil
}

// GetOpenInterestHistoryService -- GET /api/v5/rubik/stat/contracts/open-interest-history (public)
//
// Returns the historical open interest per time bar for a single contract
// instrument, in contracts, base currency and USD.
type GetOpenInterestHistoryService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetOpenInterestHistoryService(instId string) *GetOpenInterestHistoryService {
	return &GetOpenInterestHistoryService{c: c, params: map[string]string{"instId": instId}}
}

func (s *GetOpenInterestHistoryService) SetPeriod(p RubikStatPeriod) *GetOpenInterestHistoryService {
	s.params["period"] = string(p)
	return s
}

func (s *GetOpenInterestHistoryService) SetBegin(t time.Time) *GetOpenInterestHistoryService {
	s.params["begin"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetOpenInterestHistoryService) SetEnd(t time.Time) *GetOpenInterestHistoryService {
	s.params["end"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetOpenInterestHistoryService) SetLimit(limit int) *GetOpenInterestHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetOpenInterestHistoryService) Do(ctx context.Context) ([]RubikOpenInterestHistory, error) {
	req := request.Get(ctx, s.c, "/api/v5/rubik/stat/contracts/open-interest-history", s.params)
	rows, err := request.DoList[[]string](req)
	if err != nil {
		return nil, err
	}
	out := make([]RubikOpenInterestHistory, 0, len(rows))
	for _, r := range rows {
		out = append(out, parseRubikOpenInterestHistory(r))
	}
	return out, nil
}

// RubikOpenInterestHistory is one open-interest-history bar:
// [ts, oi (contracts), oiCcy (base ccy), oiUsd].
type RubikOpenInterestHistory struct {
	Timestamp            time.Time       `json:"ts"`
	OpenInterest         decimal.Decimal `json:"oi"`
	OpenInterestCurrency decimal.Decimal `json:"oiCcy"`
	OpenInterestUSD      decimal.Decimal `json:"oiUsd"`
}

func parseRubikOpenInterestHistory(r []string) RubikOpenInterestHistory {
	var v RubikOpenInterestHistory
	if len(r) > 0 {
		v.Timestamp = parseRubikTs(r[0])
	}
	if len(r) > 1 {
		v.OpenInterest = parseRubikDec(r[1])
	}
	if len(r) > 2 {
		v.OpenInterestCurrency = parseRubikDec(r[2])
	}
	if len(r) > 3 {
		v.OpenInterestUSD = parseRubikDec(r[3])
	}
	return v
}

// GetContractsOpenInterestVolumeService -- GET /api/v5/rubik/stat/contracts/open-interest-volume (public)
//
// Returns the aggregated contract open interest and trading volume per time bar
// for a currency.
type GetContractsOpenInterestVolumeService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetContractsOpenInterestVolumeService(ccy string) *GetContractsOpenInterestVolumeService {
	return &GetContractsOpenInterestVolumeService{c: c, params: map[string]string{"ccy": ccy}}
}

func (s *GetContractsOpenInterestVolumeService) SetBegin(t time.Time) *GetContractsOpenInterestVolumeService {
	s.params["begin"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetContractsOpenInterestVolumeService) SetEnd(t time.Time) *GetContractsOpenInterestVolumeService {
	s.params["end"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetContractsOpenInterestVolumeService) SetPeriod(p RubikStatPeriod) *GetContractsOpenInterestVolumeService {
	s.params["period"] = string(p)
	return s
}

func (s *GetContractsOpenInterestVolumeService) Do(ctx context.Context) ([]RubikOpenInterestVolume, error) {
	req := request.Get(ctx, s.c, "/api/v5/rubik/stat/contracts/open-interest-volume", s.params)
	return doRubikOpenInterestVolume(req)
}

// GetOptionOpenInterestVolumeService -- GET /api/v5/rubik/stat/option/open-interest-volume (public)
//
// Returns the aggregated option open interest and trading volume per time bar
// for a currency.
type GetOptionOpenInterestVolumeService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetOptionOpenInterestVolumeService(ccy string) *GetOptionOpenInterestVolumeService {
	return &GetOptionOpenInterestVolumeService{c: c, params: map[string]string{"ccy": ccy}}
}

func (s *GetOptionOpenInterestVolumeService) SetPeriod(p RubikStatPeriod) *GetOptionOpenInterestVolumeService {
	s.params["period"] = string(p)
	return s
}

func (s *GetOptionOpenInterestVolumeService) Do(ctx context.Context) ([]RubikOpenInterestVolume, error) {
	req := request.Get(ctx, s.c, "/api/v5/rubik/stat/option/open-interest-volume", s.params)
	return doRubikOpenInterestVolume(req)
}

// RubikOpenInterestVolume is one [ts, oi, vol] bar shared by the contract and
// option open-interest-volume endpoints.
type RubikOpenInterestVolume struct {
	Timestamp    time.Time       `json:"ts"`
	OpenInterest decimal.Decimal `json:"oi"`
	Volume       decimal.Decimal `json:"vol"`
}

func parseRubikOpenInterestVolume(r []string) RubikOpenInterestVolume {
	var v RubikOpenInterestVolume
	if len(r) > 0 {
		v.Timestamp = parseRubikTs(r[0])
	}
	if len(r) > 1 {
		v.OpenInterest = parseRubikDec(r[1])
	}
	if len(r) > 2 {
		v.Volume = parseRubikDec(r[2])
	}
	return v
}

// doRubikOpenInterestVolume runs a request whose rows are [ts, oi, vol].
func doRubikOpenInterestVolume(req *request.Request) ([]RubikOpenInterestVolume, error) {
	rows, err := request.DoList[[]string](req)
	if err != nil {
		return nil, err
	}
	out := make([]RubikOpenInterestVolume, 0, len(rows))
	for _, r := range rows {
		out = append(out, parseRubikOpenInterestVolume(r))
	}
	return out, nil
}

// GetOptionOpenInterestVolumeExpiryService -- GET /api/v5/rubik/stat/option/open-interest-volume-expiry (public)
//
// Returns option open interest and volume broken down by expiry per time bar for
// a currency.
type GetOptionOpenInterestVolumeExpiryService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetOptionOpenInterestVolumeExpiryService(ccy string) *GetOptionOpenInterestVolumeExpiryService {
	return &GetOptionOpenInterestVolumeExpiryService{c: c, params: map[string]string{"ccy": ccy}}
}

func (s *GetOptionOpenInterestVolumeExpiryService) SetPeriod(p RubikStatPeriod) *GetOptionOpenInterestVolumeExpiryService {
	s.params["period"] = string(p)
	return s
}

func (s *GetOptionOpenInterestVolumeExpiryService) Do(ctx context.Context) ([]RubikOpenInterestVolumeExpiry, error) {
	req := request.Get(ctx, s.c, "/api/v5/rubik/stat/option/open-interest-volume-expiry", s.params)
	rows, err := request.DoList[[]string](req)
	if err != nil {
		return nil, err
	}
	out := make([]RubikOpenInterestVolumeExpiry, 0, len(rows))
	for _, r := range rows {
		out = append(out, parseRubikOpenInterestVolumeExpiry(r))
	}
	return out, nil
}

// RubikOpenInterestVolumeExpiry is one option-by-expiry bar:
// [ts, expTime, callOI, putOI, callVol, putVol]. ExpTime is the expiry date in
// YYYYMMDD form.
type RubikOpenInterestVolumeExpiry struct {
	Timestamp        time.Time       `json:"ts"`
	ExpiryTime       string          `json:"expTime"`
	CallOpenInterest decimal.Decimal `json:"callOI"`
	PutOpenInterest  decimal.Decimal `json:"putOI"`
	CallVolume       decimal.Decimal `json:"callVol"`
	PutVolume        decimal.Decimal `json:"putVol"`
}

func parseRubikOpenInterestVolumeExpiry(r []string) RubikOpenInterestVolumeExpiry {
	var v RubikOpenInterestVolumeExpiry
	if len(r) > 0 {
		v.Timestamp = parseRubikTs(r[0])
	}
	if len(r) > 1 {
		v.ExpiryTime = r[1]
	}
	if len(r) > 2 {
		v.CallOpenInterest = parseRubikDec(r[2])
	}
	if len(r) > 3 {
		v.PutOpenInterest = parseRubikDec(r[3])
	}
	if len(r) > 4 {
		v.CallVolume = parseRubikDec(r[4])
	}
	if len(r) > 5 {
		v.PutVolume = parseRubikDec(r[5])
	}
	return v
}

// GetOptionOpenInterestVolumeStrikeService -- GET /api/v5/rubik/stat/option/open-interest-volume-strike (public)
//
// Returns option open interest and volume broken down by strike price for a
// given expiry per time bar for a currency.
type GetOptionOpenInterestVolumeStrikeService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetOptionOpenInterestVolumeStrikeService(ccy, expTime string) *GetOptionOpenInterestVolumeStrikeService {
	return &GetOptionOpenInterestVolumeStrikeService{c: c, params: map[string]string{
		"ccy":     ccy,
		"expTime": expTime,
	}}
}

func (s *GetOptionOpenInterestVolumeStrikeService) SetPeriod(p RubikStatPeriod) *GetOptionOpenInterestVolumeStrikeService {
	s.params["period"] = string(p)
	return s
}

func (s *GetOptionOpenInterestVolumeStrikeService) Do(ctx context.Context) ([]RubikOpenInterestVolumeStrike, error) {
	req := request.Get(ctx, s.c, "/api/v5/rubik/stat/option/open-interest-volume-strike", s.params)
	rows, err := request.DoList[[]string](req)
	if err != nil {
		return nil, err
	}
	out := make([]RubikOpenInterestVolumeStrike, 0, len(rows))
	for _, r := range rows {
		out = append(out, parseRubikOpenInterestVolumeStrike(r))
	}
	return out, nil
}

// RubikOpenInterestVolumeStrike is one option-by-strike bar:
// [ts, strike, callOI, putOI, callVol, putVol].
type RubikOpenInterestVolumeStrike struct {
	Timestamp        time.Time       `json:"ts"`
	Strike           decimal.Decimal `json:"strike"`
	CallOpenInterest decimal.Decimal `json:"callOI"`
	PutOpenInterest  decimal.Decimal `json:"putOI"`
	CallVolume       decimal.Decimal `json:"callVol"`
	PutVolume        decimal.Decimal `json:"putVol"`
}

func parseRubikOpenInterestVolumeStrike(r []string) RubikOpenInterestVolumeStrike {
	var v RubikOpenInterestVolumeStrike
	if len(r) > 0 {
		v.Timestamp = parseRubikTs(r[0])
	}
	if len(r) > 1 {
		v.Strike = parseRubikDec(r[1])
	}
	if len(r) > 2 {
		v.CallOpenInterest = parseRubikDec(r[2])
	}
	if len(r) > 3 {
		v.PutOpenInterest = parseRubikDec(r[3])
	}
	if len(r) > 4 {
		v.CallVolume = parseRubikDec(r[4])
	}
	if len(r) > 5 {
		v.PutVolume = parseRubikDec(r[5])
	}
	return v
}

// GetOptionTakerBlockVolumeService -- GET /api/v5/rubik/stat/option/taker-block-volume (public)
//
// Returns the latest option taker buy/sell and block-trade volumes for calls and
// puts of a currency. Unlike the other rubik endpoints, the data array is a
// single flat positional array (one snapshot), not an array of bars.
type GetOptionTakerBlockVolumeService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetOptionTakerBlockVolumeService(ccy string) *GetOptionTakerBlockVolumeService {
	return &GetOptionTakerBlockVolumeService{c: c, params: map[string]string{"ccy": ccy}}
}

func (s *GetOptionTakerBlockVolumeService) SetPeriod(p RubikStatPeriod) *GetOptionTakerBlockVolumeService {
	s.params["period"] = string(p)
	return s
}

func (s *GetOptionTakerBlockVolumeService) Do(ctx context.Context) (*RubikOptionTakerBlockVolume, error) {
	req := request.Get(ctx, s.c, "/api/v5/rubik/stat/option/taker-block-volume", s.params)
	row, err := request.DoList[string](req)
	if err != nil {
		return nil, err
	}
	if len(row) == 0 {
		return nil, nil
	}
	v := parseRubikOptionTakerBlockVolume(row)
	return &v, nil
}

// RubikOptionTakerBlockVolume is the option taker/block-volume snapshot:
// [ts, callBuyVol, callSellVol, putBuyVol, putSellVol, callBlockVol, putBlockVol].
type RubikOptionTakerBlockVolume struct {
	Timestamp       time.Time       `json:"ts"`
	CallBuyVolume   decimal.Decimal `json:"callBuyVol"`
	CallSellVolume  decimal.Decimal `json:"callSellVol"`
	PutBuyVolume    decimal.Decimal `json:"putBuyVol"`
	PutSellVolume   decimal.Decimal `json:"putSellVol"`
	CallBlockVolume decimal.Decimal `json:"callBlockVol"`
	PutBlockVolume  decimal.Decimal `json:"putBlockVol"`
}

func parseRubikOptionTakerBlockVolume(r []string) RubikOptionTakerBlockVolume {
	var v RubikOptionTakerBlockVolume
	if len(r) > 0 {
		v.Timestamp = parseRubikTs(r[0])
	}
	if len(r) > 1 {
		v.CallBuyVolume = parseRubikDec(r[1])
	}
	if len(r) > 2 {
		v.CallSellVolume = parseRubikDec(r[2])
	}
	if len(r) > 3 {
		v.PutBuyVolume = parseRubikDec(r[3])
	}
	if len(r) > 4 {
		v.PutSellVolume = parseRubikDec(r[4])
	}
	if len(r) > 5 {
		v.CallBlockVolume = parseRubikDec(r[5])
	}
	if len(r) > 6 {
		v.PutBlockVolume = parseRubikDec(r[6])
	}
	return v
}
