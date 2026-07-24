package okx

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-okx/request"
	"github.com/shopspring/decimal"
)

// GetTickersService -- GET /api/v5/market/tickers (public)
//
// Returns the latest tickers for every instrument of a product line.
type GetTickersService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetTickersService(instType InstType) *GetTickersService {
	return &GetTickersService{c: c, params: map[string]string{"instType": string(instType)}}
}

// SetUly filters tickers by underlying (FUTURES/SWAP/OPTION).
func (s *GetTickersService) SetUly(uly string) *GetTickersService {
	s.params["uly"] = uly
	return s
}

// SetInstFamily filters tickers by instrument family (FUTURES/SWAP/OPTION).
func (s *GetTickersService) SetInstFamily(instFamily string) *GetTickersService {
	s.params["instFamily"] = instFamily
	return s
}

func (s *GetTickersService) Do(ctx context.Context) ([]Ticker, error) {
	req := request.Get(ctx, s.c, "/api/v5/market/tickers", s.params)
	return request.DoList[Ticker](req)
}

// GetTickerService -- GET /api/v5/market/ticker (public)
//
// Returns the latest ticker for a single instrument.
type GetTickerService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetTickerService(instId string) *GetTickerService {
	return &GetTickerService{c: c, params: map[string]string{"instId": instId}}
}

func (s *GetTickerService) Do(ctx context.Context) (*Ticker, error) {
	req := request.Get(ctx, s.c, "/api/v5/market/ticker", s.params)
	return request.DoOne[Ticker](req)
}

// Ticker is the latest market snapshot for an instrument.
type Ticker struct {
	InstrumentType    InstType        `json:"instType"`
	InstrumentID      string          `json:"instId"`
	Last              decimal.Decimal `json:"last"`
	LastSize          decimal.Decimal `json:"lastSz"`
	AskPrice          decimal.Decimal `json:"askPx"`
	AskSize           decimal.Decimal `json:"askSz"`
	BidPrice          decimal.Decimal `json:"bidPx"`
	BidSize           decimal.Decimal `json:"bidSz"`
	Open24h           decimal.Decimal `json:"open24h"`
	High24h           decimal.Decimal `json:"high24h"`
	Low24h            decimal.Decimal `json:"low24h"`
	VolumeCurrency24h decimal.Decimal `json:"volCcy24h"`
	Volume24h         decimal.Decimal `json:"vol24h"`
	StartOfDayUTC0    decimal.Decimal `json:"sodUtc0"`
	StartOfDayUTC8    decimal.Decimal `json:"sodUtc8"`
	Timestamp         time.Time       `json:"ts"`
}

// GetBooksService -- GET /api/v5/market/books (public)
//
// Returns the order book (up to 400 levels) for an instrument.
type GetBooksService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetBooksService(instId string) *GetBooksService {
	return &GetBooksService{c: c, params: map[string]string{"instId": instId}}
}

// SetSz sets the order book depth (number of price levels).
func (s *GetBooksService) SetSz(sz int) *GetBooksService {
	s.params["sz"] = strconv.Itoa(sz)
	return s
}

func (s *GetBooksService) Do(ctx context.Context) (*OrderBook, error) {
	req := request.Get(ctx, s.c, "/api/v5/market/books", s.params)
	return request.DoOne[OrderBook](req)
}

// OrderBook is an order book snapshot. Each ask/bid level is
// [price, size, deprecated("0"), numOrders] for /books and
// [price, size, numOrders] for /books-full.
type OrderBook struct {
	Asks       [][]string `json:"asks"`
	Bids       [][]string `json:"bids"`
	Timestamp  time.Time  `json:"ts"`
	SequenceID int64      `json:"seqId"`
}

// GetBooksFullService -- GET /api/v5/market/books-full (public)
//
// Returns the full-depth order book (up to 5000 levels) for an instrument.
type GetBooksFullService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetBooksFullService(instId string) *GetBooksFullService {
	return &GetBooksFullService{c: c, params: map[string]string{"instId": instId}}
}

// SetSz sets the order book depth (number of price levels, max 5000).
func (s *GetBooksFullService) SetSz(sz int) *GetBooksFullService {
	s.params["sz"] = strconv.Itoa(sz)
	return s
}

func (s *GetBooksFullService) Do(ctx context.Context) (*OrderBook, error) {
	req := request.Get(ctx, s.c, "/api/v5/market/books-full", s.params)
	return request.DoOne[OrderBook](req)
}

// MarketBar is a candlestick time granularity (e.g. "1m", "1H", "1D").
type MarketBar string

// OKX candle bar granularities. Note the case convention OKX requires: minute
// bars are lower-case "m" while hour/day/week/month bars are upper-case
// (H/D/W/M). The ...UTC variants align the bar boundary to 00:00 UTC instead of
// the default Hong Kong time (UTC+8).
const (
	MarketBar1m  MarketBar = "1m"
	MarketBar3m  MarketBar = "3m"
	MarketBar5m  MarketBar = "5m"
	MarketBar15m MarketBar = "15m"
	MarketBar30m MarketBar = "30m"
	MarketBar1H  MarketBar = "1H"
	MarketBar2H  MarketBar = "2H"
	MarketBar4H  MarketBar = "4H"
	MarketBar6H  MarketBar = "6H"
	MarketBar12H MarketBar = "12H"
	MarketBar1D  MarketBar = "1D"
	MarketBar2D  MarketBar = "2D"
	MarketBar3D  MarketBar = "3D"
	MarketBar1W  MarketBar = "1W"
	MarketBar1M  MarketBar = "1M"
	MarketBar3M  MarketBar = "3M"

	MarketBar6Hutc  MarketBar = "6Hutc"
	MarketBar12Hutc MarketBar = "12Hutc"
	MarketBar1Dutc  MarketBar = "1Dutc"
	MarketBar2Dutc  MarketBar = "2Dutc"
	MarketBar3Dutc  MarketBar = "3Dutc"
	MarketBar1Wutc  MarketBar = "1Wutc"
	MarketBar1Mutc  MarketBar = "1Mutc"
	MarketBar3Mutc  MarketBar = "3Mutc"
)

// Candle is one OHLCV candlestick from the market candle endpoints. The raw
// response is an array-of-arrays with 9 columns:
// [ts, o, h, l, c, vol, volCcy, volCcyQuote, confirm].
type Candle struct {
	Timestamp           time.Time
	Open                decimal.Decimal
	High                decimal.Decimal
	Low                 decimal.Decimal
	Close               decimal.Decimal
	Volume              decimal.Decimal
	VolumeCurrency      decimal.Decimal
	VolumeCurrencyQuote decimal.Decimal
	Confirm             string
}

// parseCandles maps the raw [ts,o,h,l,c,vol,volCcy,volCcyQuote,confirm] rows
// into typed Candle values. Rows shorter than 9 columns are tolerated.
func parseCandles(rows [][]string) []Candle {
	out := make([]Candle, 0, len(rows))
	for _, row := range rows {
		var c Candle
		if len(row) > 0 {
			if ms, err := strconv.ParseInt(row[0], 10, 64); err == nil {
				c.Timestamp = time.UnixMilli(ms)
			}
		}
		if len(row) > 1 {
			c.Open, _ = decimal.NewFromString(row[1])
		}
		if len(row) > 2 {
			c.High, _ = decimal.NewFromString(row[2])
		}
		if len(row) > 3 {
			c.Low, _ = decimal.NewFromString(row[3])
		}
		if len(row) > 4 {
			c.Close, _ = decimal.NewFromString(row[4])
		}
		if len(row) > 5 {
			c.Volume, _ = decimal.NewFromString(row[5])
		}
		if len(row) > 6 {
			c.VolumeCurrency, _ = decimal.NewFromString(row[6])
		}
		if len(row) > 7 {
			c.VolumeCurrencyQuote, _ = decimal.NewFromString(row[7])
		}
		if len(row) > 8 {
			c.Confirm = row[8]
		}
		out = append(out, c)
	}
	return out
}

// GetCandlesService -- GET /api/v5/market/candles (public)
//
// Returns recent candlestick data for an instrument.
type GetCandlesService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetCandlesService(instId string) *GetCandlesService {
	return &GetCandlesService{c: c, params: map[string]string{"instId": instId}}
}

// SetBar sets the candle time granularity (default 1m).
func (s *GetCandlesService) SetBar(bar MarketBar) *GetCandlesService {
	s.params["bar"] = string(bar)
	return s
}

// SetAfter requests candles before the given timestamp (paging backward).
func (s *GetCandlesService) SetAfter(after time.Time) *GetCandlesService {
	s.params["after"] = strconv.FormatInt(after.UnixMilli(), 10)
	return s
}

// SetBefore requests candles after the given timestamp (paging forward).
func (s *GetCandlesService) SetBefore(before time.Time) *GetCandlesService {
	s.params["before"] = strconv.FormatInt(before.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of returned candles (max 300, default 100).
func (s *GetCandlesService) SetLimit(limit int) *GetCandlesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetCandlesService) Do(ctx context.Context) ([]Candle, error) {
	req := request.Get(ctx, s.c, "/api/v5/market/candles", s.params)
	rows, err := request.DoList[[]string](req)
	if err != nil {
		return nil, err
	}
	return parseCandles(rows), nil
}

// GetHistoryCandlesService -- GET /api/v5/market/history-candles (public)
//
// Returns historical candlestick data for an instrument.
type GetHistoryCandlesService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetHistoryCandlesService(instId string) *GetHistoryCandlesService {
	return &GetHistoryCandlesService{c: c, params: map[string]string{"instId": instId}}
}

// SetBar sets the candle time granularity (default 1m).
func (s *GetHistoryCandlesService) SetBar(bar MarketBar) *GetHistoryCandlesService {
	s.params["bar"] = string(bar)
	return s
}

// SetAfter requests candles before the given timestamp (paging backward).
func (s *GetHistoryCandlesService) SetAfter(after time.Time) *GetHistoryCandlesService {
	s.params["after"] = strconv.FormatInt(after.UnixMilli(), 10)
	return s
}

// SetBefore requests candles after the given timestamp (paging forward).
func (s *GetHistoryCandlesService) SetBefore(before time.Time) *GetHistoryCandlesService {
	s.params["before"] = strconv.FormatInt(before.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of returned candles (max 100, default 100).
func (s *GetHistoryCandlesService) SetLimit(limit int) *GetHistoryCandlesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetHistoryCandlesService) Do(ctx context.Context) ([]Candle, error) {
	req := request.Get(ctx, s.c, "/api/v5/market/history-candles", s.params)
	rows, err := request.DoList[[]string](req)
	if err != nil {
		return nil, err
	}
	return parseCandles(rows), nil
}

// GetTradesService -- GET /api/v5/market/trades (public)
//
// Returns the most recent public trades for an instrument.
type GetTradesService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetTradesService(instId string) *GetTradesService {
	return &GetTradesService{c: c, params: map[string]string{"instId": instId}}
}

// SetLimit caps the number of returned trades (max 500, default 100).
func (s *GetTradesService) SetLimit(limit int) *GetTradesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetTradesService) Do(ctx context.Context) ([]Trade, error) {
	req := request.Get(ctx, s.c, "/api/v5/market/trades", s.params)
	return request.DoList[Trade](req)
}

// Trade is a single public trade (taker fill) on an instrument.
type Trade struct {
	InstrumentID string          `json:"instId"`
	TradeID      string          `json:"tradeId"`
	Price        decimal.Decimal `json:"px"`
	Size         decimal.Decimal `json:"sz"`
	Side         Side            `json:"side"`
	// Source is the trade source. "1": RPI (Retail Price Improvement) order —
	// previously documented as ELP order; the returned value itself is
	// unchanged by the ELP->RPI rebranding.
	Source    string    `json:"source"`
	Timestamp time.Time `json:"ts"`
}

// MarketHistoryTradeType selects the paging field for history-trades:
// "1" by tradeId (default), "2" by ts.
type MarketHistoryTradeType string

const (
	MarketHistoryTradeTypeTradeId MarketHistoryTradeType = "1"
	MarketHistoryTradeTypeTs      MarketHistoryTradeType = "2"
)

// GetHistoryTradesService -- GET /api/v5/market/history-trades (public)
//
// Returns historical public trades for an instrument (paged).
type GetHistoryTradesService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetHistoryTradesService(instId string) *GetHistoryTradesService {
	return &GetHistoryTradesService{c: c, params: map[string]string{"instId": instId}}
}

// SetType selects the paging field: "1" by tradeId (default), "2" by ts.
func (s *GetHistoryTradesService) SetType(typ MarketHistoryTradeType) *GetHistoryTradesService {
	s.params["type"] = string(typ)
	return s
}

// SetAfter pages backward from the given tradeId or timestamp (per type).
func (s *GetHistoryTradesService) SetAfter(after string) *GetHistoryTradesService {
	s.params["after"] = after
	return s
}

// SetBefore pages forward from the given tradeId (type=1 only).
func (s *GetHistoryTradesService) SetBefore(before string) *GetHistoryTradesService {
	s.params["before"] = before
	return s
}

// SetLimit caps the number of returned trades (max 100, default 100).
func (s *GetHistoryTradesService) SetLimit(limit int) *GetHistoryTradesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetHistoryTradesService) Do(ctx context.Context) ([]Trade, error) {
	req := request.Get(ctx, s.c, "/api/v5/market/history-trades", s.params)
	return request.DoList[Trade](req)
}

// GetOptionInstrumentFamilyTradesService -- GET /api/v5/market/option/instrument-family-trades (public)
//
// Returns the most recent trades per strike for an option instrument family.
type GetOptionInstrumentFamilyTradesService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetOptionInstrumentFamilyTradesService(instFamily string) *GetOptionInstrumentFamilyTradesService {
	return &GetOptionInstrumentFamilyTradesService{c: c, params: map[string]string{"instFamily": instFamily}}
}

func (s *GetOptionInstrumentFamilyTradesService) Do(ctx context.Context) ([]OptionFamilyTrades, error) {
	req := request.Get(ctx, s.c, "/api/v5/market/option/instrument-family-trades", s.params)
	return request.DoList[OptionFamilyTrades](req)
}

// OptionFamilyTrades groups recent option trades by option type and underlying
// 24h volume.
type OptionFamilyTrades struct {
	OptionType OptType            `json:"optType"`
	Volume24h  decimal.Decimal    `json:"vol24h"`
	TradeInfo  []OptionFamilyInfo `json:"tradeInfo"`
}

// OptionFamilyInfo is one trade within an option instrument family.
type OptionFamilyInfo struct {
	InstrumentID string          `json:"instId"`
	TradeID      string          `json:"tradeId"`
	Price        decimal.Decimal `json:"px"`
	Size         decimal.Decimal `json:"sz"`
	Side         Side            `json:"side"`
	Timestamp    time.Time       `json:"ts"`
}

// GetPlatform24VolumeService -- GET /api/v5/market/platform-24-volume (public)
//
// Returns the platform-wide 24h trading volume.
type GetPlatform24VolumeService struct {
	c *Client
}

func (c *Client) NewGetPlatform24VolumeService() *GetPlatform24VolumeService {
	return &GetPlatform24VolumeService{c: c}
}

func (s *GetPlatform24VolumeService) Do(ctx context.Context) (*Platform24Volume, error) {
	req := request.Get(ctx, s.c, "/api/v5/market/platform-24-volume")
	return request.DoOne[Platform24Volume](req)
}

// Platform24Volume is the platform-wide 24h trading volume in CNY and USD.
type Platform24Volume struct {
	VolumeUSD decimal.Decimal `json:"volUsd"`
	VolumeCNY decimal.Decimal `json:"volCny"`
	Timestamp time.Time       `json:"ts"`
}

// GetBlockTickersService -- GET /api/v5/market/block-tickers (public)
//
// Returns the block-trade tickers for every instrument of a product line.
type GetBlockTickersService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetBlockTickersService(instType InstType) *GetBlockTickersService {
	return &GetBlockTickersService{c: c, params: map[string]string{"instType": string(instType)}}
}

// SetUly filters block tickers by underlying (FUTURES/SWAP/OPTION).
func (s *GetBlockTickersService) SetUly(uly string) *GetBlockTickersService {
	s.params["uly"] = uly
	return s
}

// SetInstFamily filters block tickers by instrument family (FUTURES/SWAP/OPTION).
func (s *GetBlockTickersService) SetInstFamily(instFamily string) *GetBlockTickersService {
	s.params["instFamily"] = instFamily
	return s
}

func (s *GetBlockTickersService) Do(ctx context.Context) ([]BlockTicker, error) {
	req := request.Get(ctx, s.c, "/api/v5/market/block-tickers", s.params)
	return request.DoList[BlockTicker](req)
}

// BlockTicker is the 24h block-trade volume snapshot for an instrument.
type BlockTicker struct {
	InstrumentType    InstType        `json:"instType"`
	InstrumentID      string          `json:"instId"`
	VolumeCurrency24h decimal.Decimal `json:"volCcy24h"`
	Volume24h         decimal.Decimal `json:"vol24h"`
	Timestamp         time.Time       `json:"ts"`
}
