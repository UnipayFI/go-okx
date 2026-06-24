package okx

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-okx/request"
	"github.com/shopspring/decimal"
)

// SprdType is a spread's pricing type (linear vs inverse legs).
type SprdType string

const (
	SprdTypeLinear   SprdType = "linear"
	SprdTypeInverse  SprdType = "inverse"
	SprdTypeHybrid   SprdType = "hybrid"
	SprdTypeFundRate SprdType = "fund_rate"
)

// SprdState is the listing state of a spread.
type SprdState string

const (
	SprdStateLive    SprdState = "live"
	SprdStateExpired SprdState = "expired"
	SprdStateSuspend SprdState = "suspend"
	SprdStatePreOpen SprdState = "preopen"
)

// SprdOrdType is a spread order type.
type SprdOrdType string

const (
	SprdOrdTypeLimit    SprdOrdType = "limit"
	SprdOrdTypeMarket   SprdOrdType = "market"
	SprdOrdTypePostOnly SprdOrdType = "post_only"
	SprdOrdTypeIOC      SprdOrdType = "ioc"
)

// SprdOrdState is the lifecycle state of a spread order.
type SprdOrdState string

const (
	SprdOrdStateLive            SprdOrdState = "live"
	SprdOrdStatePartiallyFilled SprdOrdState = "partially_filled"
	SprdOrdStateFilled          SprdOrdState = "filled"
	SprdOrdStateCanceled        SprdOrdState = "canceled"
)

// GetSprdSpreadsService -- GET /api/v5/sprd/spreads (public)
//
// Returns the tradable spreads, their legs and lot/tick sizes.
type GetSprdSpreadsService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetSprdSpreadsService() *GetSprdSpreadsService {
	return &GetSprdSpreadsService{c: c, params: map[string]string{}}
}

// SetBaseCcy filters by base currency.
func (s *GetSprdSpreadsService) SetBaseCcy(baseCcy string) *GetSprdSpreadsService {
	s.params["baseCcy"] = baseCcy
	return s
}

// SetInstId filters to spreads that include the given instrument as a leg.
func (s *GetSprdSpreadsService) SetInstId(instId string) *GetSprdSpreadsService {
	s.params["instId"] = instId
	return s
}

// SetSprdId filters by a single spread id.
func (s *GetSprdSpreadsService) SetSprdId(sprdId string) *GetSprdSpreadsService {
	s.params["sprdId"] = sprdId
	return s
}

// SetState filters by spread state.
func (s *GetSprdSpreadsService) SetState(state SprdState) *GetSprdSpreadsService {
	s.params["state"] = string(state)
	return s
}

func (s *GetSprdSpreadsService) Do(ctx context.Context) ([]SprdSpread, error) {
	req := request.Get(ctx, s.c, "/api/v5/sprd/spreads", s.params)
	return request.DoList[SprdSpread](req)
}

// SprdSpread is a single tradable spread and its legs.
type SprdSpread struct {
	SpreadID      string          `json:"sprdId"`
	SpreadType    SprdType        `json:"sprdType"`
	State         SprdState       `json:"state"`
	BaseCurrency  string          `json:"baseCcy"`
	SizeCurrency  string          `json:"szCcy"`
	QuoteCurrency string          `json:"quoteCcy"`
	TickSize      decimal.Decimal `json:"tickSz"`
	MinSize       decimal.Decimal `json:"minSz"`
	LotSize       decimal.Decimal `json:"lotSz"`
	ListTime      time.Time       `json:"listTime"`
	Legs          []SprdLeg       `json:"legs"`
	ExpiryTime    time.Time       `json:"expTime"`
	UpdateTime    time.Time       `json:"uTime"`
}

// SprdLeg is one leg of a spread.
type SprdLeg struct {
	InstrumentID string `json:"instId"`
	Side         Side   `json:"side"`
}

// GetSprdBooksService -- GET /api/v5/sprd/books (public)
//
// Returns the order book of a spread.
type GetSprdBooksService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetSprdBooksService(sprdId string) *GetSprdBooksService {
	return &GetSprdBooksService{c: c, params: map[string]string{"sprdId": sprdId}}
}

// SetSz sets the order book depth (number of price levels, max 400).
func (s *GetSprdBooksService) SetSz(sz int) *GetSprdBooksService {
	s.params["sz"] = strconv.Itoa(sz)
	return s
}

func (s *GetSprdBooksService) Do(ctx context.Context) (*SprdOrderBook, error) {
	req := request.Get(ctx, s.c, "/api/v5/sprd/books", s.params)
	return request.DoOne[SprdOrderBook](req)
}

// SprdOrderBook is a spread's order book snapshot. Each ask/bid level is
// [price, size, numOrders].
type SprdOrderBook struct {
	Asks      [][]string `json:"asks"`
	Bids      [][]string `json:"bids"`
	Timestamp time.Time  `json:"ts"`
}

// GetSprdTickerService -- GET /api/v5/sprd/ticker (public)
//
// Returns the latest ticker for a spread.
type GetSprdTickerService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetSprdTickerService(sprdId string) *GetSprdTickerService {
	return &GetSprdTickerService{c: c, params: map[string]string{"sprdId": sprdId}}
}

func (s *GetSprdTickerService) Do(ctx context.Context) (*SprdTicker, error) {
	req := request.Get(ctx, s.c, "/api/v5/sprd/ticker", s.params)
	return request.DoOne[SprdTicker](req)
}

// SprdTicker is the latest market snapshot for a spread.
type SprdTicker struct {
	SpreadID  string          `json:"sprdId"`
	Last      decimal.Decimal `json:"last"`
	LastSize  decimal.Decimal `json:"lastSz"`
	AskPrice  decimal.Decimal `json:"askPx"`
	AskSize   decimal.Decimal `json:"askSz"`
	BidPrice  decimal.Decimal `json:"bidPx"`
	BidSize   decimal.Decimal `json:"bidSz"`
	Timestamp time.Time       `json:"ts"`
}

// GetSprdPublicTradesService -- GET /api/v5/sprd/public-trades (public)
//
// Returns the most recent public trades for spreads.
type GetSprdPublicTradesService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetSprdPublicTradesService() *GetSprdPublicTradesService {
	return &GetSprdPublicTradesService{c: c, params: map[string]string{}}
}

// SetSprdId filters trades to a single spread.
func (s *GetSprdPublicTradesService) SetSprdId(sprdId string) *GetSprdPublicTradesService {
	s.params["sprdId"] = sprdId
	return s
}

func (s *GetSprdPublicTradesService) Do(ctx context.Context) ([]SprdPublicTrade, error) {
	req := request.Get(ctx, s.c, "/api/v5/sprd/public-trades", s.params)
	return request.DoList[SprdPublicTrade](req)
}

// SprdPublicTrade is a single public trade (taker fill) on a spread.
type SprdPublicTrade struct {
	SpreadID  string          `json:"sprdId"`
	TradeID   string          `json:"tradeId"`
	Price     decimal.Decimal `json:"px"`
	Size      decimal.Decimal `json:"sz"`
	Side      Side            `json:"side"`
	Timestamp time.Time       `json:"ts"`
}

// SprdCandle is one OHLCV candlestick for a spread. The raw response is an
// array-of-arrays with 7 columns: [ts, o, h, l, c, vol, confirm].
type SprdCandle struct {
	Timestamp time.Time
	Open      decimal.Decimal
	High      decimal.Decimal
	Low       decimal.Decimal
	Close     decimal.Decimal
	Volume    decimal.Decimal
	Confirm   string
}

// parseSprdCandles maps the raw [ts,o,h,l,c,vol,confirm] rows into typed
// SprdCandle values. Rows shorter than 7 columns are tolerated.
func parseSprdCandles(rows [][]string) []SprdCandle {
	out := make([]SprdCandle, 0, len(rows))
	for _, row := range rows {
		var c SprdCandle
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
			c.Confirm = row[6]
		}
		out = append(out, c)
	}
	return out
}

// GetSprdCandlesService -- GET /api/v5/market/sprd-candles (public)
//
// Returns recent candlestick data for a spread. Note the path lives under
// /market/, not /sprd/.
type GetSprdCandlesService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetSprdCandlesService(sprdId string) *GetSprdCandlesService {
	return &GetSprdCandlesService{c: c, params: map[string]string{"sprdId": sprdId}}
}

// SetBar sets the candle time granularity (default 1m).
func (s *GetSprdCandlesService) SetBar(bar MarketBar) *GetSprdCandlesService {
	s.params["bar"] = string(bar)
	return s
}

// SetAfter requests candles before the given timestamp (paging backward).
func (s *GetSprdCandlesService) SetAfter(after time.Time) *GetSprdCandlesService {
	s.params["after"] = strconv.FormatInt(after.UnixMilli(), 10)
	return s
}

// SetBefore requests candles after the given timestamp (paging forward).
func (s *GetSprdCandlesService) SetBefore(before time.Time) *GetSprdCandlesService {
	s.params["before"] = strconv.FormatInt(before.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of returned candles (max 100, default 100).
func (s *GetSprdCandlesService) SetLimit(limit int) *GetSprdCandlesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetSprdCandlesService) Do(ctx context.Context) ([]SprdCandle, error) {
	req := request.Get(ctx, s.c, "/api/v5/market/sprd-candles", s.params)
	rows, err := request.DoList[[]string](req)
	if err != nil {
		return nil, err
	}
	return parseSprdCandles(rows), nil
}

// GetSprdHistoryCandlesService -- GET /api/v5/market/sprd-history-candles (public)
//
// Returns historical candlestick data for a spread. Note the path lives under
// /market/, not /sprd/.
type GetSprdHistoryCandlesService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetSprdHistoryCandlesService(sprdId string) *GetSprdHistoryCandlesService {
	return &GetSprdHistoryCandlesService{c: c, params: map[string]string{"sprdId": sprdId}}
}

// SetBar sets the candle time granularity (default 1m).
func (s *GetSprdHistoryCandlesService) SetBar(bar MarketBar) *GetSprdHistoryCandlesService {
	s.params["bar"] = string(bar)
	return s
}

// SetAfter requests candles before the given timestamp (paging backward).
func (s *GetSprdHistoryCandlesService) SetAfter(after time.Time) *GetSprdHistoryCandlesService {
	s.params["after"] = strconv.FormatInt(after.UnixMilli(), 10)
	return s
}

// SetBefore requests candles after the given timestamp (paging forward).
func (s *GetSprdHistoryCandlesService) SetBefore(before time.Time) *GetSprdHistoryCandlesService {
	s.params["before"] = strconv.FormatInt(before.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of returned candles (max 100, default 100).
func (s *GetSprdHistoryCandlesService) SetLimit(limit int) *GetSprdHistoryCandlesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetSprdHistoryCandlesService) Do(ctx context.Context) ([]SprdCandle, error) {
	req := request.Get(ctx, s.c, "/api/v5/market/sprd-history-candles", s.params)
	rows, err := request.DoList[[]string](req)
	if err != nil {
		return nil, err
	}
	return parseSprdCandles(rows), nil
}

// GetSprdOrderService -- GET /api/v5/sprd/order (Read)
//
// Returns the details of a single spread order. Either ordId or clOrdId must be
// supplied.
type GetSprdOrderService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetSprdOrderService() *GetSprdOrderService {
	return &GetSprdOrderService{c: c, params: map[string]string{}}
}

// SetOrdId looks the order up by its OKX order id.
func (s *GetSprdOrderService) SetOrdId(ordId string) *GetSprdOrderService {
	s.params["ordId"] = ordId
	return s
}

// SetClOrdId looks the order up by its client-supplied order id.
func (s *GetSprdOrderService) SetClOrdId(clOrdId string) *GetSprdOrderService {
	s.params["clOrdId"] = clOrdId
	return s
}

func (s *GetSprdOrderService) Do(ctx context.Context) (*SprdOrder, error) {
	req := request.Get(ctx, s.c, "/api/v5/sprd/order", s.params).WithSign()
	return request.DoOne[SprdOrder](req)
}

// SprdOrder is a single spread order's full state.
type SprdOrder struct {
	SpreadID            string          `json:"sprdId"`
	OrderID             string          `json:"ordId"`
	ClientOrderID       string          `json:"clOrdId"`
	Tag                 string          `json:"tag"`
	Price               decimal.Decimal `json:"px"`
	Size                decimal.Decimal `json:"sz"`
	OrderType           SprdOrdType     `json:"ordType"`
	Side                Side            `json:"side"`
	FillSize            decimal.Decimal `json:"fillSz"`
	FillPrice           decimal.Decimal `json:"fillPx"`
	TradeID             string          `json:"tradeId"`
	AccumulatedFillSize decimal.Decimal `json:"accFillSz"`
	PendingFillSize     decimal.Decimal `json:"pendingFillSz"`
	PendingSettleSize   decimal.Decimal `json:"pendingSettleSz"`
	CanceledSize        decimal.Decimal `json:"canceledSz"`
	State               SprdOrdState    `json:"state"`
	AveragePrice        decimal.Decimal `json:"avgPx"`
	CancelSource        string          `json:"cancelSource"`
	UpdateTime          time.Time       `json:"uTime"`
	CreationTime        time.Time       `json:"cTime"`
}

// GetSprdOrdersPendingService -- GET /api/v5/sprd/orders-pending (Read)
//
// Returns the account's currently open (incomplete) spread orders.
type GetSprdOrdersPendingService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetSprdOrdersPendingService() *GetSprdOrdersPendingService {
	return &GetSprdOrdersPendingService{c: c, params: map[string]string{}}
}

// SetSprdId filters by spread id.
func (s *GetSprdOrdersPendingService) SetSprdId(sprdId string) *GetSprdOrdersPendingService {
	s.params["sprdId"] = sprdId
	return s
}

// SetOrdType filters by order type.
func (s *GetSprdOrdersPendingService) SetOrdType(ordType SprdOrdType) *GetSprdOrdersPendingService {
	s.params["ordType"] = string(ordType)
	return s
}

// SetState filters by order state (live / partially_filled).
func (s *GetSprdOrdersPendingService) SetState(state SprdOrdState) *GetSprdOrdersPendingService {
	s.params["state"] = string(state)
	return s
}

// SetBeginId pages from orders with an id newer than beginId.
func (s *GetSprdOrdersPendingService) SetBeginId(beginId string) *GetSprdOrdersPendingService {
	s.params["beginId"] = beginId
	return s
}

// SetEndId pages to orders with an id older than endId.
func (s *GetSprdOrdersPendingService) SetEndId(endId string) *GetSprdOrdersPendingService {
	s.params["endId"] = endId
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetSprdOrdersPendingService) SetLimit(limit int) *GetSprdOrdersPendingService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetSprdOrdersPendingService) Do(ctx context.Context) ([]SprdOrder, error) {
	req := request.Get(ctx, s.c, "/api/v5/sprd/orders-pending", s.params).WithSign()
	return request.DoList[SprdOrder](req)
}

// GetSprdOrdersHistoryService -- GET /api/v5/sprd/orders-history (Read)
//
// Returns the account's completed spread orders from the last 21 days.
type GetSprdOrdersHistoryService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetSprdOrdersHistoryService() *GetSprdOrdersHistoryService {
	return &GetSprdOrdersHistoryService{c: c, params: map[string]string{}}
}

// SetSprdId filters by spread id.
func (s *GetSprdOrdersHistoryService) SetSprdId(sprdId string) *GetSprdOrdersHistoryService {
	s.params["sprdId"] = sprdId
	return s
}

// SetOrdType filters by order type.
func (s *GetSprdOrdersHistoryService) SetOrdType(ordType SprdOrdType) *GetSprdOrdersHistoryService {
	s.params["ordType"] = string(ordType)
	return s
}

// SetState filters by order state (canceled / filled).
func (s *GetSprdOrdersHistoryService) SetState(state SprdOrdState) *GetSprdOrdersHistoryService {
	s.params["state"] = string(state)
	return s
}

// SetBeginId pages from orders with an id newer than beginId.
func (s *GetSprdOrdersHistoryService) SetBeginId(beginId string) *GetSprdOrdersHistoryService {
	s.params["beginId"] = beginId
	return s
}

// SetEndId pages to orders with an id older than endId.
func (s *GetSprdOrdersHistoryService) SetEndId(endId string) *GetSprdOrdersHistoryService {
	s.params["endId"] = endId
	return s
}

// SetBegin filters to orders created at or after the given time.
func (s *GetSprdOrdersHistoryService) SetBegin(t time.Time) *GetSprdOrdersHistoryService {
	s.params["begin"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEnd filters to orders created at or before the given time.
func (s *GetSprdOrdersHistoryService) SetEnd(t time.Time) *GetSprdOrdersHistoryService {
	s.params["end"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetSprdOrdersHistoryService) SetLimit(limit int) *GetSprdOrdersHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetSprdOrdersHistoryService) Do(ctx context.Context) ([]SprdOrder, error) {
	req := request.Get(ctx, s.c, "/api/v5/sprd/orders-history", s.params).WithSign()
	return request.DoList[SprdOrder](req)
}

// GetSprdOrdersHistoryArchiveService -- GET /api/v5/sprd/orders-history-archive (Read)
//
// Returns the account's archived (older than 21 days) completed spread orders.
type GetSprdOrdersHistoryArchiveService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetSprdOrdersHistoryArchiveService() *GetSprdOrdersHistoryArchiveService {
	return &GetSprdOrdersHistoryArchiveService{c: c, params: map[string]string{}}
}

// SetSprdId filters by spread id.
func (s *GetSprdOrdersHistoryArchiveService) SetSprdId(sprdId string) *GetSprdOrdersHistoryArchiveService {
	s.params["sprdId"] = sprdId
	return s
}

// SetOrdType filters by order type.
func (s *GetSprdOrdersHistoryArchiveService) SetOrdType(ordType SprdOrdType) *GetSprdOrdersHistoryArchiveService {
	s.params["ordType"] = string(ordType)
	return s
}

// SetState filters by order state (canceled / filled).
func (s *GetSprdOrdersHistoryArchiveService) SetState(state SprdOrdState) *GetSprdOrdersHistoryArchiveService {
	s.params["state"] = string(state)
	return s
}

// SetBeginId pages from orders with an id newer than beginId.
func (s *GetSprdOrdersHistoryArchiveService) SetBeginId(beginId string) *GetSprdOrdersHistoryArchiveService {
	s.params["beginId"] = beginId
	return s
}

// SetEndId pages to orders with an id older than endId.
func (s *GetSprdOrdersHistoryArchiveService) SetEndId(endId string) *GetSprdOrdersHistoryArchiveService {
	s.params["endId"] = endId
	return s
}

// SetBegin filters to orders created at or after the given time.
func (s *GetSprdOrdersHistoryArchiveService) SetBegin(t time.Time) *GetSprdOrdersHistoryArchiveService {
	s.params["begin"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEnd filters to orders created at or before the given time.
func (s *GetSprdOrdersHistoryArchiveService) SetEnd(t time.Time) *GetSprdOrdersHistoryArchiveService {
	s.params["end"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetSprdOrdersHistoryArchiveService) SetLimit(limit int) *GetSprdOrdersHistoryArchiveService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetSprdOrdersHistoryArchiveService) Do(ctx context.Context) ([]SprdOrder, error) {
	req := request.Get(ctx, s.c, "/api/v5/sprd/orders-history-archive", s.params).WithSign()
	return request.DoList[SprdOrder](req)
}

// GetSprdTradesService -- GET /api/v5/sprd/trades (Read)
//
// Returns the account's spread fills from the last 21 days.
type GetSprdTradesService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetSprdTradesService() *GetSprdTradesService {
	return &GetSprdTradesService{c: c, params: map[string]string{}}
}

// SetSprdId filters by spread id.
func (s *GetSprdTradesService) SetSprdId(sprdId string) *GetSprdTradesService {
	s.params["sprdId"] = sprdId
	return s
}

// SetTradeId filters by a single trade id.
func (s *GetSprdTradesService) SetTradeId(tradeId string) *GetSprdTradesService {
	s.params["tradeId"] = tradeId
	return s
}

// SetOrdId filters fills by their parent order id.
func (s *GetSprdTradesService) SetOrdId(ordId string) *GetSprdTradesService {
	s.params["ordId"] = ordId
	return s
}

// SetBeginId pages from fills with an id newer than beginId.
func (s *GetSprdTradesService) SetBeginId(beginId string) *GetSprdTradesService {
	s.params["beginId"] = beginId
	return s
}

// SetEndId pages to fills with an id older than endId.
func (s *GetSprdTradesService) SetEndId(endId string) *GetSprdTradesService {
	s.params["endId"] = endId
	return s
}

// SetBegin filters to fills at or after the given time.
func (s *GetSprdTradesService) SetBegin(t time.Time) *GetSprdTradesService {
	s.params["begin"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEnd filters to fills at or before the given time.
func (s *GetSprdTradesService) SetEnd(t time.Time) *GetSprdTradesService {
	s.params["end"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetSprdTradesService) SetLimit(limit int) *GetSprdTradesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetSprdTradesService) Do(ctx context.Context) ([]SprdTrade, error) {
	req := request.Get(ctx, s.c, "/api/v5/sprd/trades", s.params).WithSign()
	return request.DoList[SprdTrade](req)
}

// SprdTrade is a single account spread fill, including its per-leg fills.
type SprdTrade struct {
	SpreadID      string          `json:"sprdId"`
	TradeID       string          `json:"tradeId"`
	OrderID       string          `json:"ordId"`
	ClientOrderID string          `json:"clOrdId"`
	Tag           string          `json:"tag"`
	FillPrice     decimal.Decimal `json:"fillPx"`
	FillSize      decimal.Decimal `json:"fillSz"`
	Side          Side            `json:"side"`
	State         string          `json:"state"`
	ExecutionType ExecType        `json:"execType"`
	Timestamp     time.Time       `json:"ts"`
	Legs          []SprdTradeLeg  `json:"legs"`
	Code          string          `json:"code"`
	Message       string          `json:"msg"`
}

// SprdTradeLeg is one leg fill within a spread trade.
type SprdTradeLeg struct {
	InstrumentID string          `json:"instId"`
	Price        decimal.Decimal `json:"px"`
	Size         decimal.Decimal `json:"sz"`
	SizeContract decimal.Decimal `json:"szCont"`
	Side         Side            `json:"side"`
	Fee          decimal.Decimal `json:"fee"`
	FeeCurrency  string          `json:"feeCcy"`
	TradeID      string          `json:"tradeId"`
}

// GetSprdOrderAlgoService -- GET /api/v5/sprd/order-algo (Read)
//
// Returns the details of a single spread algo order. Either algoId or
// algoClOrdId must be supplied.
//
// NOTE: live verification (2026-06) returns HTTP 404 for this path on
// www.okx.com; the spread algo-order endpoints are not currently deployed.
// The service is modeled from the OKX docs so it works if/when OKX enables it.
type GetSprdOrderAlgoService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetSprdOrderAlgoService() *GetSprdOrderAlgoService {
	return &GetSprdOrderAlgoService{c: c, params: map[string]string{}}
}

// SetAlgoId looks the algo order up by its OKX algo id.
func (s *GetSprdOrderAlgoService) SetAlgoId(algoId string) *GetSprdOrderAlgoService {
	s.params["algoId"] = algoId
	return s
}

// SetAlgoClOrdId looks the algo order up by its client-supplied id.
func (s *GetSprdOrderAlgoService) SetAlgoClOrdId(algoClOrdId string) *GetSprdOrderAlgoService {
	s.params["algoClOrdId"] = algoClOrdId
	return s
}

func (s *GetSprdOrderAlgoService) Do(ctx context.Context) (*SprdAlgoOrder, error) {
	req := request.Get(ctx, s.c, "/api/v5/sprd/order-algo", s.params).WithSign()
	return request.DoOne[SprdAlgoOrder](req)
}

// SprdAlgoOrder is a single spread algo order's state.
type SprdAlgoOrder struct {
	SpreadID          string          `json:"sprdId"`
	AlgoID            string          `json:"algoId"`
	AlgoClientOrderID string          `json:"algoClOrdId"`
	OrderType         string          `json:"ordType"`
	Side              Side            `json:"side"`
	Size              decimal.Decimal `json:"sz"`
	SizeLimit         decimal.Decimal `json:"szLimit"`
	PriceVariation    decimal.Decimal `json:"pxVar"`
	PriceSpread       decimal.Decimal `json:"pxSpread"`
	PriceLimit        decimal.Decimal `json:"pxLimit"`
	TimeInterval      string          `json:"timeInterval"`
	State             string          `json:"state"`
	Tag               string          `json:"tag"`
	UpdateTime        time.Time       `json:"uTime"`
	CreationTime      time.Time       `json:"cTime"`
}

// GetSprdOrdersAlgoPendingService -- GET /api/v5/sprd/orders-algo-pending (Read)
//
// Returns the account's currently open spread algo orders.
//
// NOTE: live verification (2026-06) returns HTTP 404 for this path on
// www.okx.com; the spread algo-order endpoints are not currently deployed.
type GetSprdOrdersAlgoPendingService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetSprdOrdersAlgoPendingService() *GetSprdOrdersAlgoPendingService {
	return &GetSprdOrdersAlgoPendingService{c: c, params: map[string]string{}}
}

// SetSprdId filters by spread id.
func (s *GetSprdOrdersAlgoPendingService) SetSprdId(sprdId string) *GetSprdOrdersAlgoPendingService {
	s.params["sprdId"] = sprdId
	return s
}

// SetAlgoId filters by a single algo id.
func (s *GetSprdOrdersAlgoPendingService) SetAlgoId(algoId string) *GetSprdOrdersAlgoPendingService {
	s.params["algoId"] = algoId
	return s
}

// SetOrdType filters by algo order type.
func (s *GetSprdOrdersAlgoPendingService) SetOrdType(ordType string) *GetSprdOrdersAlgoPendingService {
	s.params["ordType"] = ordType
	return s
}

// SetState filters by algo order state.
func (s *GetSprdOrdersAlgoPendingService) SetState(state string) *GetSprdOrdersAlgoPendingService {
	s.params["state"] = state
	return s
}

func (s *GetSprdOrdersAlgoPendingService) Do(ctx context.Context) ([]SprdAlgoOrder, error) {
	req := request.Get(ctx, s.c, "/api/v5/sprd/orders-algo-pending", s.params).WithSign()
	return request.DoList[SprdAlgoOrder](req)
}

// GetSprdOrdersAlgoHistoryService -- GET /api/v5/sprd/orders-algo-history (Read)
//
// Returns the account's completed spread algo orders.
//
// NOTE: live verification (2026-06) returns HTTP 404 for this path on
// www.okx.com; the spread algo-order endpoints are not currently deployed.
type GetSprdOrdersAlgoHistoryService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetSprdOrdersAlgoHistoryService() *GetSprdOrdersAlgoHistoryService {
	return &GetSprdOrdersAlgoHistoryService{c: c, params: map[string]string{}}
}

// SetSprdId filters by spread id.
func (s *GetSprdOrdersAlgoHistoryService) SetSprdId(sprdId string) *GetSprdOrdersAlgoHistoryService {
	s.params["sprdId"] = sprdId
	return s
}

// SetAlgoId filters by a single algo id.
func (s *GetSprdOrdersAlgoHistoryService) SetAlgoId(algoId string) *GetSprdOrdersAlgoHistoryService {
	s.params["algoId"] = algoId
	return s
}

// SetOrdType filters by algo order type.
func (s *GetSprdOrdersAlgoHistoryService) SetOrdType(ordType string) *GetSprdOrdersAlgoHistoryService {
	s.params["ordType"] = ordType
	return s
}

// SetState filters by algo order state.
func (s *GetSprdOrdersAlgoHistoryService) SetState(state string) *GetSprdOrdersAlgoHistoryService {
	s.params["state"] = state
	return s
}

// SetBegin filters to algo orders created at or after the given time.
func (s *GetSprdOrdersAlgoHistoryService) SetBegin(t time.Time) *GetSprdOrdersAlgoHistoryService {
	s.params["begin"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEnd filters to algo orders created at or before the given time.
func (s *GetSprdOrdersAlgoHistoryService) SetEnd(t time.Time) *GetSprdOrdersAlgoHistoryService {
	s.params["end"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetSprdOrdersAlgoHistoryService) SetLimit(limit int) *GetSprdOrdersAlgoHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetSprdOrdersAlgoHistoryService) Do(ctx context.Context) ([]SprdAlgoOrder, error) {
	req := request.Get(ctx, s.c, "/api/v5/sprd/orders-algo-history", s.params).WithSign()
	return request.DoList[SprdAlgoOrder](req)
}

// --- State-changing endpoints (Trade): implement-only, modeled from OKX docs.
// These place / amend / cancel spread (algo) orders and MUST NOT be executed
// during verification. ---

// PlaceSprdOrderService -- POST /api/v5/sprd/order (Trade)
//
// Places a spread order. sprdId, side, ordType and sz are required; px is
// required for limit / post_only / ioc orders.
type PlaceSprdOrderService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewPlaceSprdOrderService(sprdId string, side Side, ordType SprdOrdType, sz string) *PlaceSprdOrderService {
	return &PlaceSprdOrderService{c: c, body: map[string]any{
		"sprdId":  sprdId,
		"side":    string(side),
		"ordType": string(ordType),
		"sz":      sz,
	}}
}

// SetPx sets the order price (required for limit / post_only / ioc).
func (s *PlaceSprdOrderService) SetPx(px string) *PlaceSprdOrderService {
	s.body["px"] = px
	return s
}

// SetClOrdId sets a client-supplied order id.
func (s *PlaceSprdOrderService) SetClOrdId(clOrdId string) *PlaceSprdOrderService {
	s.body["clOrdId"] = clOrdId
	return s
}

// SetTag sets an order tag.
func (s *PlaceSprdOrderService) SetTag(tag string) *PlaceSprdOrderService {
	s.body["tag"] = tag
	return s
}

func (s *PlaceSprdOrderService) Do(ctx context.Context) (*SprdOrderAck, error) {
	req := request.Post(ctx, s.c, "/api/v5/sprd/order", s.body).WithSign()
	return request.DoOne[SprdOrderAck](req)
}

// SprdOrderAck is the acknowledgement of a place / amend / cancel spread order.
type SprdOrderAck struct {
	OrderID       string `json:"ordId"`
	ClientOrderID string `json:"clOrdId"`
	Tag           string `json:"tag"`
	SCode         string `json:"sCode"`
	SMsg          string `json:"sMsg"`
}

// AmendSprdOrderService -- POST /api/v5/sprd/amend-order (Trade)
//
// Amends an existing spread order's price and/or size. Either ordId or clOrdId
// is required, plus at least one of newSz / newPx.
type AmendSprdOrderService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewAmendSprdOrderService() *AmendSprdOrderService {
	return &AmendSprdOrderService{c: c, body: map[string]any{}}
}

// SetOrdId identifies the order to amend by its OKX order id.
func (s *AmendSprdOrderService) SetOrdId(ordId string) *AmendSprdOrderService {
	s.body["ordId"] = ordId
	return s
}

// SetClOrdId identifies the order to amend by its client-supplied id.
func (s *AmendSprdOrderService) SetClOrdId(clOrdId string) *AmendSprdOrderService {
	s.body["clOrdId"] = clOrdId
	return s
}

// SetReqId sets a client request id echoed back in the response.
func (s *AmendSprdOrderService) SetReqId(reqId string) *AmendSprdOrderService {
	s.body["reqId"] = reqId
	return s
}

// SetNewSz sets the new order size.
func (s *AmendSprdOrderService) SetNewSz(newSz string) *AmendSprdOrderService {
	s.body["newSz"] = newSz
	return s
}

// SetNewPx sets the new order price.
func (s *AmendSprdOrderService) SetNewPx(newPx string) *AmendSprdOrderService {
	s.body["newPx"] = newPx
	return s
}

func (s *AmendSprdOrderService) Do(ctx context.Context) (*SprdAmendAck, error) {
	req := request.Post(ctx, s.c, "/api/v5/sprd/amend-order", s.body).WithSign()
	return request.DoOne[SprdAmendAck](req)
}

// SprdAmendAck is the acknowledgement of an amend spread order request.
type SprdAmendAck struct {
	OrderID       string `json:"ordId"`
	ClientOrderID string `json:"clOrdId"`
	RequestID     string `json:"reqId"`
	SCode         string `json:"sCode"`
	SMsg          string `json:"sMsg"`
}

// CancelSprdOrderService -- POST /api/v5/sprd/cancel-order (Trade)
//
// Cancels a single spread order by ordId or clOrdId.
type CancelSprdOrderService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewCancelSprdOrderService() *CancelSprdOrderService {
	return &CancelSprdOrderService{c: c, body: map[string]any{}}
}

// SetOrdId cancels by OKX order id.
func (s *CancelSprdOrderService) SetOrdId(ordId string) *CancelSprdOrderService {
	s.body["ordId"] = ordId
	return s
}

// SetClOrdId cancels by client-supplied order id.
func (s *CancelSprdOrderService) SetClOrdId(clOrdId string) *CancelSprdOrderService {
	s.body["clOrdId"] = clOrdId
	return s
}

func (s *CancelSprdOrderService) Do(ctx context.Context) (*SprdOrderAck, error) {
	req := request.Post(ctx, s.c, "/api/v5/sprd/cancel-order", s.body).WithSign()
	return request.DoOne[SprdOrderAck](req)
}

// CancelAllSprdOrdersService -- POST /api/v5/sprd/cancel-all-orders (Trade)
//
// Cancels all pending spread orders, optionally scoped to one spread.
type CancelAllSprdOrdersService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewCancelAllSprdOrdersService() *CancelAllSprdOrdersService {
	return &CancelAllSprdOrdersService{c: c, body: map[string]any{}}
}

// SetSprdId scopes the cancellation to a single spread.
func (s *CancelAllSprdOrdersService) SetSprdId(sprdId string) *CancelAllSprdOrdersService {
	s.body["sprdId"] = sprdId
	return s
}

func (s *CancelAllSprdOrdersService) Do(ctx context.Context) (*SprdCancelAllAck, error) {
	req := request.Post(ctx, s.c, "/api/v5/sprd/cancel-all-orders", s.body).WithSign()
	return request.DoOne[SprdCancelAllAck](req)
}

// SprdCancelAllAck is the acknowledgement of a cancel-all spread orders request.
type SprdCancelAllAck struct {
	TriggerTime time.Time `json:"triggerTime"`
	Timestamp   time.Time `json:"ts"`
}

// SprdMassCancelOrdersService -- POST /api/v5/sprd/mass-cancel (Trade)
//
// Cancels all pending spread orders (mass cancel), optionally scoped to one
// spread.
type SprdMassCancelOrdersService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewSprdMassCancelOrdersService() *SprdMassCancelOrdersService {
	return &SprdMassCancelOrdersService{c: c, body: map[string]any{}}
}

// SetSprdId scopes the mass cancel to a single spread.
func (s *SprdMassCancelOrdersService) SetSprdId(sprdId string) *SprdMassCancelOrdersService {
	s.body["sprdId"] = sprdId
	return s
}

func (s *SprdMassCancelOrdersService) Do(ctx context.Context) (*SprdMassCancelAck, error) {
	req := request.Post(ctx, s.c, "/api/v5/sprd/mass-cancel", s.body).WithSign()
	return request.DoOne[SprdMassCancelAck](req)
}

// SprdMassCancelAck is the acknowledgement of a mass-cancel spread orders
// request.
type SprdMassCancelAck struct {
	Result bool `json:"result"`
}

// SprdCancelAllAfterService -- POST /api/v5/sprd/cancel-all-after (Trade)
//
// Arms or disarms the dead-man's-switch: all pending spread orders are canceled
// timeOut seconds after the last heartbeat. timeOut of 0 cancels the timer.
type SprdCancelAllAfterService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewSprdCancelAllAfterService(timeOut int) *SprdCancelAllAfterService {
	return &SprdCancelAllAfterService{c: c, body: map[string]any{
		"timeOut": strconv.Itoa(timeOut),
	}}
}

func (s *SprdCancelAllAfterService) Do(ctx context.Context) (*SprdCancelAllAfterAck, error) {
	req := request.Post(ctx, s.c, "/api/v5/sprd/cancel-all-after", s.body).WithSign()
	return request.DoOne[SprdCancelAllAfterAck](req)
}

// SprdCancelAllAfterAck is the acknowledgement of a cancel-all-after request.
type SprdCancelAllAfterAck struct {
	TriggerTime time.Time `json:"triggerTime"`
	Timestamp   time.Time `json:"ts"`
}

// PlaceSprdOrderAlgoService -- POST /api/v5/sprd/order-algo (Trade)
//
// Places a spread algo (e.g. TWAP) order. sprdId, side, ordType and sz are
// required; the remaining params depend on the algo order type.
type PlaceSprdOrderAlgoService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewPlaceSprdOrderAlgoService(sprdId string, side Side, ordType, sz string) *PlaceSprdOrderAlgoService {
	return &PlaceSprdOrderAlgoService{c: c, body: map[string]any{
		"sprdId":  sprdId,
		"side":    string(side),
		"ordType": ordType,
		"sz":      sz,
	}}
}

// SetAlgoClOrdId sets a client-supplied algo order id.
func (s *PlaceSprdOrderAlgoService) SetAlgoClOrdId(algoClOrdId string) *PlaceSprdOrderAlgoService {
	s.body["algoClOrdId"] = algoClOrdId
	return s
}

// SetSzLimit sets the per-interval size limit (TWAP).
func (s *PlaceSprdOrderAlgoService) SetSzLimit(szLimit string) *PlaceSprdOrderAlgoService {
	s.body["szLimit"] = szLimit
	return s
}

// SetPxVar sets the price variance (TWAP).
func (s *PlaceSprdOrderAlgoService) SetPxVar(pxVar string) *PlaceSprdOrderAlgoService {
	s.body["pxVar"] = pxVar
	return s
}

// SetPxSpread sets the price spread (TWAP).
func (s *PlaceSprdOrderAlgoService) SetPxSpread(pxSpread string) *PlaceSprdOrderAlgoService {
	s.body["pxSpread"] = pxSpread
	return s
}

// SetPxLimit sets the limit price (TWAP).
func (s *PlaceSprdOrderAlgoService) SetPxLimit(pxLimit string) *PlaceSprdOrderAlgoService {
	s.body["pxLimit"] = pxLimit
	return s
}

// SetTimeInterval sets the order interval in seconds (TWAP).
func (s *PlaceSprdOrderAlgoService) SetTimeInterval(timeInterval string) *PlaceSprdOrderAlgoService {
	s.body["timeInterval"] = timeInterval
	return s
}

// SetTag sets an order tag.
func (s *PlaceSprdOrderAlgoService) SetTag(tag string) *PlaceSprdOrderAlgoService {
	s.body["tag"] = tag
	return s
}

func (s *PlaceSprdOrderAlgoService) Do(ctx context.Context) (*SprdAlgoAck, error) {
	req := request.Post(ctx, s.c, "/api/v5/sprd/order-algo", s.body).WithSign()
	return request.DoOne[SprdAlgoAck](req)
}

// SprdAlgoAck is the acknowledgement of a place / cancel spread algo order.
type SprdAlgoAck struct {
	AlgoID            string `json:"algoId"`
	AlgoClientOrderID string `json:"algoClOrdId"`
	SCode             string `json:"sCode"`
	SMsg              string `json:"sMsg"`
}

// CancelSprdOrderAlgoService -- POST /api/v5/sprd/cancel-order-algo (Trade)
//
// Cancels a single spread algo order by algoId or algoClOrdId.
type CancelSprdOrderAlgoService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewCancelSprdOrderAlgoService() *CancelSprdOrderAlgoService {
	return &CancelSprdOrderAlgoService{c: c, body: map[string]any{}}
}

// SetAlgoId cancels by OKX algo id.
func (s *CancelSprdOrderAlgoService) SetAlgoId(algoId string) *CancelSprdOrderAlgoService {
	s.body["algoId"] = algoId
	return s
}

// SetAlgoClOrdId cancels by client-supplied algo id.
func (s *CancelSprdOrderAlgoService) SetAlgoClOrdId(algoClOrdId string) *CancelSprdOrderAlgoService {
	s.body["algoClOrdId"] = algoClOrdId
	return s
}

func (s *CancelSprdOrderAlgoService) Do(ctx context.Context) (*SprdAlgoAck, error) {
	req := request.Post(ctx, s.c, "/api/v5/sprd/cancel-order-algo", s.body).WithSign()
	return request.DoOne[SprdAlgoAck](req)
}
