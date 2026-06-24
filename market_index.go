package okx

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-okx/common"
	"github.com/UnipayFI/go-okx/request"
	"github.com/shopspring/decimal"
)

// GetIndexTickersService -- GET /api/v5/market/index-tickers (public)
//
// Retrieves index tickers. Either quoteCcy or instId is required.
type GetIndexTickersService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetIndexTickersService() *GetIndexTickersService {
	return &GetIndexTickersService{c: c, params: map[string]string{}}
}

// SetQuoteCcy sets the quote currency filter (e.g. USD, USDT, BTC). Currently
// only USD, USDT and BTC are supported.
func (s *GetIndexTickersService) SetQuoteCcy(v string) *GetIndexTickersService {
	s.params["quoteCcy"] = v
	return s
}

// SetInstId sets the index instrument ID (e.g. BTC-USDT).
func (s *GetIndexTickersService) SetInstId(v string) *GetIndexTickersService {
	s.params["instId"] = v
	return s
}

func (s *GetIndexTickersService) Do(ctx context.Context) ([]IndexTicker, error) {
	req := request.Get(ctx, s.c, "/api/v5/market/index-tickers", s.params)
	return request.DoList[IndexTicker](req)
}

// IndexTicker is the latest index price snapshot for an index instrument.
type IndexTicker struct {
	InstrumentID   string          `json:"instId"`
	IndexPrice     decimal.Decimal `json:"idxPx"`
	High24h        decimal.Decimal `json:"high24h"`
	StartOfDayUTC0 decimal.Decimal `json:"sodUtc0"`
	Open24h        decimal.Decimal `json:"open24h"`
	Low24h         decimal.Decimal `json:"low24h"`
	StartOfDayUTC8 decimal.Decimal `json:"sodUtc8"`
	Timestamp      time.Time       `json:"ts"`
}

// GetIndexCandlesService -- GET /api/v5/market/index-candles (public)
//
// Retrieves the candlestick charts of an index. Data is sorted newest-first and
// is returned for up to the most recent ~1440 bars.
type GetIndexCandlesService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetIndexCandlesService(instId string) *GetIndexCandlesService {
	return &GetIndexCandlesService{c: c, params: map[string]string{"instId": instId}}
}

// SetBar sets the bar size (e.g. 1m, 5m, 1H, 1D). Default is 1m.
func (s *GetIndexCandlesService) SetBar(v string) *GetIndexCandlesService {
	s.params["bar"] = v
	return s
}

// SetAfter requests records earlier than the given timestamp (pagination).
func (s *GetIndexCandlesService) SetAfter(v time.Time) *GetIndexCandlesService {
	s.params["after"] = strconv.FormatInt(v.UnixMilli(), 10)
	return s
}

// SetBefore requests records newer than the given timestamp (pagination).
func (s *GetIndexCandlesService) SetBefore(v time.Time) *GetIndexCandlesService {
	s.params["before"] = strconv.FormatInt(v.UnixMilli(), 10)
	return s
}

// SetLimit sets the number of results per request (max 100, default 100).
func (s *GetIndexCandlesService) SetLimit(v int) *GetIndexCandlesService {
	s.params["limit"] = strconv.Itoa(v)
	return s
}

func (s *GetIndexCandlesService) Do(ctx context.Context) ([]IndexCandle, error) {
	req := request.Get(ctx, s.c, "/api/v5/market/index-candles", s.params)
	rows, err := request.DoList[[]string](req)
	if err != nil {
		return nil, err
	}
	return parseIndexCandles(rows), nil
}

// GetHistoryIndexCandlesService -- GET /api/v5/market/history-index-candles (public)
//
// Retrieves the older candlestick charts of an index from the last 3 months.
type GetHistoryIndexCandlesService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetHistoryIndexCandlesService(instId string) *GetHistoryIndexCandlesService {
	return &GetHistoryIndexCandlesService{c: c, params: map[string]string{"instId": instId}}
}

// SetBar sets the bar size (e.g. 1m, 5m, 1H, 1D). Default is 1m.
func (s *GetHistoryIndexCandlesService) SetBar(v string) *GetHistoryIndexCandlesService {
	s.params["bar"] = v
	return s
}

// SetAfter requests records earlier than the given timestamp (pagination).
func (s *GetHistoryIndexCandlesService) SetAfter(v time.Time) *GetHistoryIndexCandlesService {
	s.params["after"] = strconv.FormatInt(v.UnixMilli(), 10)
	return s
}

// SetBefore requests records newer than the given timestamp (pagination).
func (s *GetHistoryIndexCandlesService) SetBefore(v time.Time) *GetHistoryIndexCandlesService {
	s.params["before"] = strconv.FormatInt(v.UnixMilli(), 10)
	return s
}

// SetLimit sets the number of results per request (max 100, default 100).
func (s *GetHistoryIndexCandlesService) SetLimit(v int) *GetHistoryIndexCandlesService {
	s.params["limit"] = strconv.Itoa(v)
	return s
}

func (s *GetHistoryIndexCandlesService) Do(ctx context.Context) ([]IndexCandle, error) {
	req := request.Get(ctx, s.c, "/api/v5/market/history-index-candles", s.params)
	rows, err := request.DoList[[]string](req)
	if err != nil {
		return nil, err
	}
	return parseIndexCandles(rows), nil
}

// GetMarkPriceCandlesService -- GET /api/v5/market/mark-price-candles (public)
//
// Retrieves the candlestick charts of the mark price. Data is sorted
// newest-first for up to the most recent ~1440 bars.
type GetMarkPriceCandlesService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetMarkPriceCandlesService(instId string) *GetMarkPriceCandlesService {
	return &GetMarkPriceCandlesService{c: c, params: map[string]string{"instId": instId}}
}

// SetBar sets the bar size (e.g. 1m, 5m, 1H, 1D). Default is 1m.
func (s *GetMarkPriceCandlesService) SetBar(v string) *GetMarkPriceCandlesService {
	s.params["bar"] = v
	return s
}

// SetAfter requests records earlier than the given timestamp (pagination).
func (s *GetMarkPriceCandlesService) SetAfter(v time.Time) *GetMarkPriceCandlesService {
	s.params["after"] = strconv.FormatInt(v.UnixMilli(), 10)
	return s
}

// SetBefore requests records newer than the given timestamp (pagination).
func (s *GetMarkPriceCandlesService) SetBefore(v time.Time) *GetMarkPriceCandlesService {
	s.params["before"] = strconv.FormatInt(v.UnixMilli(), 10)
	return s
}

// SetLimit sets the number of results per request (max 100, default 100).
func (s *GetMarkPriceCandlesService) SetLimit(v int) *GetMarkPriceCandlesService {
	s.params["limit"] = strconv.Itoa(v)
	return s
}

func (s *GetMarkPriceCandlesService) Do(ctx context.Context) ([]IndexCandle, error) {
	req := request.Get(ctx, s.c, "/api/v5/market/mark-price-candles", s.params)
	rows, err := request.DoList[[]string](req)
	if err != nil {
		return nil, err
	}
	return parseIndexCandles(rows), nil
}

// GetHistoryMarkPriceCandlesService -- GET /api/v5/market/history-mark-price-candles (public)
//
// Retrieves the older candlestick charts of the mark price from the last 3
// months.
type GetHistoryMarkPriceCandlesService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetHistoryMarkPriceCandlesService(instId string) *GetHistoryMarkPriceCandlesService {
	return &GetHistoryMarkPriceCandlesService{c: c, params: map[string]string{"instId": instId}}
}

// SetBar sets the bar size (e.g. 1m, 5m, 1H, 1D). Default is 1m.
func (s *GetHistoryMarkPriceCandlesService) SetBar(v string) *GetHistoryMarkPriceCandlesService {
	s.params["bar"] = v
	return s
}

// SetAfter requests records earlier than the given timestamp (pagination).
func (s *GetHistoryMarkPriceCandlesService) SetAfter(v time.Time) *GetHistoryMarkPriceCandlesService {
	s.params["after"] = strconv.FormatInt(v.UnixMilli(), 10)
	return s
}

// SetBefore requests records newer than the given timestamp (pagination).
func (s *GetHistoryMarkPriceCandlesService) SetBefore(v time.Time) *GetHistoryMarkPriceCandlesService {
	s.params["before"] = strconv.FormatInt(v.UnixMilli(), 10)
	return s
}

// SetLimit sets the number of results per request (max 100, default 100).
func (s *GetHistoryMarkPriceCandlesService) SetLimit(v int) *GetHistoryMarkPriceCandlesService {
	s.params["limit"] = strconv.Itoa(v)
	return s
}

func (s *GetHistoryMarkPriceCandlesService) Do(ctx context.Context) ([]IndexCandle, error) {
	req := request.Get(ctx, s.c, "/api/v5/market/history-mark-price-candles", s.params)
	rows, err := request.DoList[[]string](req)
	if err != nil {
		return nil, err
	}
	return parseIndexCandles(rows), nil
}

// IndexCandle is one row of an index or mark-price candlestick. OKX returns
// these endpoints as arrays-of-arrays with 6 columns: [ts, o, h, l, c, confirm].
// Confirm is "0" for an in-progress bar and "1" once the bar is closed.
type IndexCandle struct {
	Timestamp time.Time       `json:"ts"`
	Open      decimal.Decimal `json:"o"`
	High      decimal.Decimal `json:"h"`
	Low       decimal.Decimal `json:"l"`
	Close     decimal.Decimal `json:"c"`
	Confirm   string          `json:"confirm"`
}

// parseIndexCandles maps the raw [ts, o, h, l, c, confirm] rows into typed
// IndexCandle values, skipping any row that does not have the expected width.
func parseIndexCandles(rows [][]string) []IndexCandle {
	out := make([]IndexCandle, 0, len(rows))
	for _, row := range rows {
		if len(row) < 6 {
			continue
		}
		var ts time.Time
		if ms, err := strconv.ParseInt(row[0], 10, 64); err == nil {
			ts = time.UnixMilli(ms)
		}
		out = append(out, IndexCandle{
			Timestamp: ts,
			Open:      decimalOrZero(row[1]),
			High:      decimalOrZero(row[2]),
			Low:       decimalOrZero(row[3]),
			Close:     decimalOrZero(row[4]),
			Confirm:   row[5],
		})
	}
	return out
}

// decimalOrZero parses a decimal string, returning zero when the value is empty
// or malformed.
func decimalOrZero(s string) decimal.Decimal {
	if s == "" {
		return decimal.Zero
	}
	v, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.Zero
	}
	return v
}

// GetExchangeRateService -- GET /api/v5/market/exchange-rate (public)
//
// Retrieves the average exchange rate for USD-CNY over the past 2 weeks.
type GetExchangeRateService struct {
	c *Client
}

func (c *Client) NewGetExchangeRateService() *GetExchangeRateService {
	return &GetExchangeRateService{c: c}
}

func (s *GetExchangeRateService) Do(ctx context.Context) (*ExchangeRate, error) {
	req := request.Get(ctx, s.c, "/api/v5/market/exchange-rate")
	return request.DoOne[ExchangeRate](req)
}

// ExchangeRate is the USD-CNY reference exchange rate.
type ExchangeRate struct {
	USDCNY decimal.Decimal `json:"usdCny"`
}

// GetIndexComponentsService -- GET /api/v5/market/index-components (public)
//
// Retrieves the constituent components of an index and their weights.
type GetIndexComponentsService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetIndexComponentsService(index string) *GetIndexComponentsService {
	return &GetIndexComponentsService{c: c, params: map[string]string{"index": index}}
}

func (s *GetIndexComponentsService) Do(ctx context.Context) (*IndexComponents, error) {
	req := request.Get(ctx, s.c, "/api/v5/market/index-components", s.params)
	// Unlike every other OKX endpoint, index-components returns its "data" as a
	// single JSON object rather than an array, so the array-based Do* helpers
	// cannot decode it. Read the raw "data" bytes and unmarshal the object.
	raw, err := request.DoRawData(req)
	if err != nil {
		return nil, err
	}
	if len(raw) == 0 || string(raw) == "null" {
		return nil, nil
	}
	var out IndexComponents
	if err := common.JSONUnmarshal(raw, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// IndexComponents describes an index, its latest value and the venues that make
// it up.
type IndexComponents struct {
	Index      string           `json:"index"`
	Last       decimal.Decimal  `json:"last"`
	Timestamp  time.Time        `json:"ts"`
	Components []IndexComponent `json:"components"`
}

// IndexComponent is one venue's contribution to an index.
type IndexComponent struct {
	Symbol       string          `json:"symbol"`
	SymbolPrice  decimal.Decimal `json:"symPx"`
	Weight       decimal.Decimal `json:"wgt"`
	ConvertPrice decimal.Decimal `json:"cnvPx"`
	Exch         string          `json:"exch"`
}
