package okx

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-okx/common"
	"github.com/UnipayFI/go-okx/request"
	"github.com/shopspring/decimal"
)

// This file wraps the OKX v5 WebSocket BUSINESS-gateway channels
// (request.GatewayBusiness): the candlestick streams (candle / mark-price-candle
// / index-candle / sprd-candle), the spread market channels (sprd-tickers /
// sprd-public-trades / sprd-books5), the algo & grid order channels (orders-algo
// / algo-advance / grid-orders-* / grid-positions / grid-sub-orders) and the
// economic-calendar channel. Candle channels are public; algo/grid/
// economic-calendar channels require login (private=true).

// --- candlestick channels (array-of-arrays payload) ---------------------------

// WsCandle is one OHLCV candlestick pushed by the "candle{bar}" channel. The raw
// data frame is an array-of-arrays with 9 columns:
// [ts, o, h, l, c, vol, volCcy, volCcyQuote, confirm].
type WsCandle struct {
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

// WsCandleHandler is invoked for every "candle{bar}" push (or error). The push's
// data array is parsed into typed WsCandle rows.
type WsCandleHandler func([]WsCandle, error)

// WsIndexCandle is one OHLC candlestick pushed by the "mark-price-candle{bar}"
// and "index-candle{bar}" channels. The raw data frame is an array-of-arrays
// with 6 columns: [ts, o, h, l, c, confirm] (no volume).
type WsIndexCandle struct {
	Timestamp time.Time
	Open      decimal.Decimal
	High      decimal.Decimal
	Low       decimal.Decimal
	Close     decimal.Decimal
	Confirm   string
}

// WsIndexCandleHandler is invoked for every mark-price-candle / index-candle
// push (or error).
type WsIndexCandleHandler func([]WsIndexCandle, error)

// WsSprdCandle is one OHLCV candlestick pushed by the "sprd-candle{bar}" channel.
// The raw data frame is an array-of-arrays with 7 columns:
// [ts, o, h, l, c, vol, confirm].
type WsSprdCandle struct {
	Timestamp time.Time
	Open      decimal.Decimal
	High      decimal.Decimal
	Low       decimal.Decimal
	Close     decimal.Decimal
	Volume    decimal.Decimal
	Confirm   string
}

// WsSprdCandleHandler is invoked for every "sprd-candle{bar}" push (or error).
type WsSprdCandleHandler func([]WsSprdCandle, error)

// candleCol parses one decimal column, tolerating a short row.
func candleCol(row []string, i int) decimal.Decimal {
	if len(row) > i {
		d, _ := decimal.NewFromString(row[i])
		return d
	}
	return decimal.Zero
}

// candleTs parses the millisecond ts column (index 0).
func candleTs(row []string) time.Time {
	if len(row) > 0 {
		if ms, err := strconv.ParseInt(row[0], 10, 64); err == nil {
			return time.UnixMilli(ms)
		}
	}
	return time.Time{}
}

// parseWsCandles maps [ts,o,h,l,c,vol,volCcy,volCcyQuote,confirm] rows into
// typed WsCandle values. Rows shorter than 9 columns are tolerated.
func parseWsCandles(rows [][]string) []WsCandle {
	out := make([]WsCandle, 0, len(rows))
	for _, row := range rows {
		c := WsCandle{
			Timestamp:           candleTs(row),
			Open:                candleCol(row, 1),
			High:                candleCol(row, 2),
			Low:                 candleCol(row, 3),
			Close:               candleCol(row, 4),
			Volume:              candleCol(row, 5),
			VolumeCurrency:      candleCol(row, 6),
			VolumeCurrencyQuote: candleCol(row, 7),
		}
		if len(row) > 8 {
			c.Confirm = row[8]
		}
		out = append(out, c)
	}
	return out
}

// parseWsIndexCandles maps [ts,o,h,l,c,confirm] rows into typed WsIndexCandle
// values. Rows shorter than 6 columns are tolerated.
func parseWsIndexCandles(rows [][]string) []WsIndexCandle {
	out := make([]WsIndexCandle, 0, len(rows))
	for _, row := range rows {
		c := WsIndexCandle{
			Timestamp: candleTs(row),
			Open:      candleCol(row, 1),
			High:      candleCol(row, 2),
			Low:       candleCol(row, 3),
			Close:     candleCol(row, 4),
		}
		if len(row) > 5 {
			c.Confirm = row[5]
		}
		out = append(out, c)
	}
	return out
}

// parseWsSprdCandles maps [ts,o,h,l,c,vol,confirm] rows into typed WsSprdCandle
// values. Rows shorter than 7 columns are tolerated.
func parseWsSprdCandles(rows [][]string) []WsSprdCandle {
	out := make([]WsSprdCandle, 0, len(rows))
	for _, row := range rows {
		c := WsSprdCandle{
			Timestamp: candleTs(row),
			Open:      candleCol(row, 1),
			High:      candleCol(row, 2),
			Low:       candleCol(row, 3),
			Close:     candleCol(row, 4),
			Volume:    candleCol(row, 5),
		}
		if len(row) > 6 {
			c.Confirm = row[6]
		}
		out = append(out, c)
	}
	return out
}

// wsCandleFrame is the raw data array of a candle push (array-of-arrays).
type wsCandleFrame struct {
	Data [][]string `json:"data"`
}

// SubscribeCandleService -- "candle{bar}" channel (business; public).
//
// Streams the instrument's candlesticks at the given bar (e.g. "1m", "1H", "1D",
// "1Dutc"). The channel name is "candle"+bar; the payload is an array-of-arrays,
// delivered as parsed WsCandle rows.
type SubscribeCandleService struct {
	c      *WebSocketClient
	instId string
	bar    string
}

func (c *WebSocketClient) NewSubscribeCandleService(instId, bar string) *SubscribeCandleService {
	return &SubscribeCandleService{c: c, instId: instId, bar: bar}
}

func (s *SubscribeCandleService) Do(ctx context.Context, cb WsCandleHandler) (chan<- struct{}, <-chan struct{}, error) {
	return request.SubscribeRaw(ctx, s.c, request.GatewayBusiness, false,
		request.WsArg{Channel: "candle" + s.bar, InstrumentID: s.instId},
		func(message []byte, err error) {
			if err != nil {
				cb(nil, err)
				return
			}
			var f wsCandleFrame
			if err := common.JSONUnmarshal(message, &f); err != nil {
				cb(nil, err)
				return
			}
			cb(parseWsCandles(f.Data), nil)
		})
}

// SubscribeMarkPriceCandleService -- "mark-price-candle{bar}" channel (business; public).
//
// Streams the instrument's mark-price candlesticks at the given bar. The channel
// name is "mark-price-candle"+bar; the payload is a 6-column array-of-arrays,
// delivered as parsed WsIndexCandle rows.
type SubscribeMarkPriceCandleService struct {
	c      *WebSocketClient
	instId string
	bar    string
}

func (c *WebSocketClient) NewSubscribeMarkPriceCandleService(instId, bar string) *SubscribeMarkPriceCandleService {
	return &SubscribeMarkPriceCandleService{c: c, instId: instId, bar: bar}
}

func (s *SubscribeMarkPriceCandleService) Do(ctx context.Context, cb WsIndexCandleHandler) (chan<- struct{}, <-chan struct{}, error) {
	return request.SubscribeRaw(ctx, s.c, request.GatewayBusiness, false,
		request.WsArg{Channel: "mark-price-candle" + s.bar, InstrumentID: s.instId},
		func(message []byte, err error) {
			if err != nil {
				cb(nil, err)
				return
			}
			var f wsCandleFrame
			if err := common.JSONUnmarshal(message, &f); err != nil {
				cb(nil, err)
				return
			}
			cb(parseWsIndexCandles(f.Data), nil)
		})
}

// SubscribeIndexCandleService -- "index-candle{bar}" channel (business; public).
//
// Streams the index's candlesticks at the given bar. The channel name is
// "index-candle"+bar; the payload is a 6-column array-of-arrays, delivered as
// parsed WsIndexCandle rows.
type SubscribeIndexCandleService struct {
	c      *WebSocketClient
	instId string
	bar    string
}

func (c *WebSocketClient) NewSubscribeIndexCandleService(instId, bar string) *SubscribeIndexCandleService {
	return &SubscribeIndexCandleService{c: c, instId: instId, bar: bar}
}

func (s *SubscribeIndexCandleService) Do(ctx context.Context, cb WsIndexCandleHandler) (chan<- struct{}, <-chan struct{}, error) {
	return request.SubscribeRaw(ctx, s.c, request.GatewayBusiness, false,
		request.WsArg{Channel: "index-candle" + s.bar, InstrumentID: s.instId},
		func(message []byte, err error) {
			if err != nil {
				cb(nil, err)
				return
			}
			var f wsCandleFrame
			if err := common.JSONUnmarshal(message, &f); err != nil {
				cb(nil, err)
				return
			}
			cb(parseWsIndexCandles(f.Data), nil)
		})
}

// SubscribeSprdCandleService -- "sprd-candle{bar}" channel (business; public).
//
// Streams a spread's candlesticks at the given bar. The channel name is
// "sprd-candle"+bar; the payload is a 7-column array-of-arrays, delivered as
// parsed WsSprdCandle rows.
type SubscribeSprdCandleService struct {
	c      *WebSocketClient
	sprdId string
	bar    string
}

func (c *WebSocketClient) NewSubscribeSprdCandleService(sprdId, bar string) *SubscribeSprdCandleService {
	return &SubscribeSprdCandleService{c: c, sprdId: sprdId, bar: bar}
}

func (s *SubscribeSprdCandleService) Do(ctx context.Context, cb WsSprdCandleHandler) (chan<- struct{}, <-chan struct{}, error) {
	return request.SubscribeRaw(ctx, s.c, request.GatewayBusiness, false,
		request.WsArg{Channel: "sprd-candle" + s.bar, SpreadID: s.sprdId},
		func(message []byte, err error) {
			if err != nil {
				cb(nil, err)
				return
			}
			var f wsCandleFrame
			if err := common.JSONUnmarshal(message, &f); err != nil {
				cb(nil, err)
				return
			}
			cb(parseWsSprdCandles(f.Data), nil)
		})
}

// --- spread market channels (object payload) ----------------------------------

// WsSprdTicker is the latest market snapshot for a spread pushed by the
// "sprd-tickers" channel.
type WsSprdTicker struct {
	SpreadID  string          `json:"sprdId"`
	Last      decimal.Decimal `json:"last"`
	LastSize  decimal.Decimal `json:"lastSz"`
	AskPrice  decimal.Decimal `json:"askPx"`
	AskSize   decimal.Decimal `json:"askSz"`
	BidPrice  decimal.Decimal `json:"bidPx"`
	BidSize   decimal.Decimal `json:"bidSz"`
	Open24h   decimal.Decimal `json:"open24h"`
	High24h   decimal.Decimal `json:"high24h"`
	Low24h    decimal.Decimal `json:"low24h"`
	Volume24h decimal.Decimal `json:"vol24h"`
	Timestamp time.Time       `json:"ts"`
}

// SubscribeSprdTickersService -- "sprd-tickers" channel (business; public).
type SubscribeSprdTickersService struct {
	c      *WebSocketClient
	sprdId string
}

func (c *WebSocketClient) NewSubscribeSprdTickersService(sprdId string) *SubscribeSprdTickersService {
	return &SubscribeSprdTickersService{c: c, sprdId: sprdId}
}

func (s *SubscribeSprdTickersService) Do(ctx context.Context, cb WsHandler[WsSprdTicker]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsSprdTicker](ctx, s.c, request.GatewayBusiness, false,
		request.WsArg{Channel: "sprd-tickers", SpreadID: s.sprdId}, cb)
}

// WsSprdPublicTrade is a single public trade (taker fill) on a spread pushed by
// the "sprd-public-trades" channel.
type WsSprdPublicTrade struct {
	SpreadID  string          `json:"sprdId"`
	TradeID   string          `json:"tradeId"`
	Price     decimal.Decimal `json:"px"`
	Size      decimal.Decimal `json:"sz"`
	Side      Side            `json:"side"`
	Timestamp time.Time       `json:"ts"`
}

// SubscribeSprdPublicTradesService -- "sprd-public-trades" channel (business; public).
type SubscribeSprdPublicTradesService struct {
	c      *WebSocketClient
	sprdId string
}

func (c *WebSocketClient) NewSubscribeSprdPublicTradesService(sprdId string) *SubscribeSprdPublicTradesService {
	return &SubscribeSprdPublicTradesService{c: c, sprdId: sprdId}
}

func (s *SubscribeSprdPublicTradesService) Do(ctx context.Context, cb WsHandler[WsSprdPublicTrade]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsSprdPublicTrade](ctx, s.c, request.GatewayBusiness, false,
		request.WsArg{Channel: "sprd-public-trades", SpreadID: s.sprdId}, cb)
}

// WsSprdBooks is a spread's 5-level order book snapshot/update pushed by the
// "sprd-books5" channel. Each ask/bid level is [price, size, numOrders].
type WsSprdBooks struct {
	Asks       [][]string `json:"asks"`
	Bids       [][]string `json:"bids"`
	Timestamp  time.Time  `json:"ts"`
	SequenceID int64      `json:"seqId"`
}

// SubscribeSprdBooks5Service -- "sprd-books5" channel (business; public).
type SubscribeSprdBooks5Service struct {
	c      *WebSocketClient
	sprdId string
}

func (c *WebSocketClient) NewSubscribeSprdBooks5Service(sprdId string) *SubscribeSprdBooks5Service {
	return &SubscribeSprdBooks5Service{c: c, sprdId: sprdId}
}

func (s *SubscribeSprdBooks5Service) Do(ctx context.Context, cb WsHandler[WsSprdBooks]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsSprdBooks](ctx, s.c, request.GatewayBusiness, false,
		request.WsArg{Channel: "sprd-books5", SpreadID: s.sprdId}, cb)
}

// --- algo order channels (private; login) -------------------------------------

// WsAlgoOrder is a single algo (conditional / oco / trigger / move_order_stop)
// order pushed by the "orders-algo" channel. The validating account had no algo
// orders, so the field set is a union modeled from the OKX WS doc field table for
// the orders-algo channel.
type WsAlgoOrder struct {
	InstrumentType             InstType        `json:"instType"`
	InstrumentID               string          `json:"instId"`
	Currency                   string          `json:"ccy"`
	OrderID                    string          `json:"ordId"`
	AlgoID                     string          `json:"algoId"`
	ClientOrderID              string          `json:"clOrdId"`
	AlgoClientOrderID          string          `json:"algoClOrdId"`
	Size                       decimal.Decimal `json:"sz"`
	OrderType                  AlgoOrdType     `json:"ordType"`
	Side                       Side            `json:"side"`
	PositionSide               PosSide         `json:"posSide"`
	TradeMode                  TdMode          `json:"tdMode"`
	TargetCurrency             TgtCcy          `json:"tgtCcy"`
	NotionalUSD                decimal.Decimal `json:"notionalUsd"`
	OrderPrice                 decimal.Decimal `json:"ordPx"`
	Price                      decimal.Decimal `json:"px"`
	State                      AlgoState       `json:"state"`
	Leverage                   decimal.Decimal `json:"lever"`
	TakeProfitTriggerPrice     decimal.Decimal `json:"tpTriggerPx"`
	TakeProfitTriggerPriceType string          `json:"tpTriggerPxType"`
	TakeProfitOrderPrice       decimal.Decimal `json:"tpOrdPx"`
	StopLossTriggerPrice       decimal.Decimal `json:"slTriggerPx"`
	StopLossTriggerPriceType   string          `json:"slTriggerPxType"`
	StopLossOrderPrice         decimal.Decimal `json:"slOrdPx"`
	TriggerPrice               decimal.Decimal `json:"triggerPx"`
	TriggerPriceType           string          `json:"triggerPxType"`
	OrderPriceType             string          `json:"ordPxType"`
	ActualSize                 decimal.Decimal `json:"actualSz"`
	ActualPrice                decimal.Decimal `json:"actualPx"`
	ActualSide                 string          `json:"actualSide"`
	TriggerTime                time.Time       `json:"triggerTime"`
	Tag                        string          `json:"tag"`
	ReduceOnly                 string          `json:"reduceOnly"`
	Last                       decimal.Decimal `json:"last"`
	FailCode                   string          `json:"failCode"`
	AmendResult                string          `json:"amendResult"`
	RequestID                  string          `json:"reqId"`
	AmendPriceOnTriggerType    string          `json:"amendPxOnTriggerType"`
	CreationTime               time.Time       `json:"cTime"`
	UpdateTime                 time.Time       `json:"uTime"`
}

// SubscribeOrdersAlgoService -- "orders-algo" channel (business; login).
//
// Streams the account's conditional / oco / trigger / move_order_stop algo
// orders for a product line. InstFamily/InstId narrow the subscription.
type SubscribeOrdersAlgoService struct {
	c          *WebSocketClient
	instType   InstType
	instFamily string
	instId     string
}

func (c *WebSocketClient) NewSubscribeOrdersAlgoService(instType InstType) *SubscribeOrdersAlgoService {
	return &SubscribeOrdersAlgoService{c: c, instType: instType}
}

// SetInstFamily narrows the subscription to an instrument family.
func (s *SubscribeOrdersAlgoService) SetInstFamily(instFamily string) *SubscribeOrdersAlgoService {
	s.instFamily = instFamily
	return s
}

// SetInstId narrows the subscription to a single instrument.
func (s *SubscribeOrdersAlgoService) SetInstId(instId string) *SubscribeOrdersAlgoService {
	s.instId = instId
	return s
}

func (s *SubscribeOrdersAlgoService) Do(ctx context.Context, cb WsHandler[WsAlgoOrder]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsAlgoOrder](ctx, s.c, request.GatewayBusiness, true,
		request.WsArg{Channel: "orders-algo", InstrumentType: string(s.instType), InstrumentFamily: s.instFamily, InstrumentID: s.instId}, cb)
}

// WsAdvanceAlgoOrder is a single advanced algo (iceberg / twap / grid) order
// pushed by the "algo-advance" channel. The validating account had no advanced
// algo orders, so the field set is modeled from the OKX WS doc field table for
// the algo-advance channel.
type WsAdvanceAlgoOrder struct {
	InstrumentType    InstType        `json:"instType"`
	InstrumentID      string          `json:"instId"`
	Currency          string          `json:"ccy"`
	OrderID           string          `json:"ordId"`
	AlgoID            string          `json:"algoId"`
	ClientOrderID     string          `json:"clOrdId"`
	AlgoClientOrderID string          `json:"algoClOrdId"`
	Size              decimal.Decimal `json:"sz"`
	OrderType         AlgoOrdType     `json:"ordType"`
	Side              Side            `json:"side"`
	PositionSide      PosSide         `json:"posSide"`
	TradeMode         TdMode          `json:"tdMode"`
	TargetCurrency    TgtCcy          `json:"tgtCcy"`
	State             AlgoState       `json:"state"`
	Leverage          decimal.Decimal `json:"lever"`
	PriceVariation    decimal.Decimal `json:"pxVar"`
	PriceSpread       decimal.Decimal `json:"pxSpread"`
	PriceLimit        decimal.Decimal `json:"pxLimit"`
	SizeLimit         decimal.Decimal `json:"szLimit"`
	TimeInterval      string          `json:"timeInterval"`
	TriggerPrice      decimal.Decimal `json:"triggerPx"`
	OrderPrice        decimal.Decimal `json:"ordPx"`
	ActualSize        decimal.Decimal `json:"actualSz"`
	ActualPrice       decimal.Decimal `json:"actualPx"`
	ActualSide        string          `json:"actualSide"`
	NotionalUSD       decimal.Decimal `json:"notionalUsd"`
	TriggerTime       time.Time       `json:"triggerTime"`
	Tag               string          `json:"tag"`
	Count             string          `json:"count"`
	Last              decimal.Decimal `json:"last"`
	FailCode          string          `json:"failCode"`
	AmendResult       string          `json:"amendResult"`
	RequestID         string          `json:"reqId"`
	CreationTime      time.Time       `json:"cTime"`
	UpdateTime        time.Time       `json:"uTime"`
}

// SubscribeAlgoAdvanceService -- "algo-advance" channel (business; login).
//
// Streams the account's advanced algo (iceberg / twap) orders for a product line.
// InstId narrows the subscription.
type SubscribeAlgoAdvanceService struct {
	c        *WebSocketClient
	instType InstType
	instId   string
}

func (c *WebSocketClient) NewSubscribeAlgoAdvanceService(instType InstType) *SubscribeAlgoAdvanceService {
	return &SubscribeAlgoAdvanceService{c: c, instType: instType}
}

// SetInstId narrows the subscription to a single instrument.
func (s *SubscribeAlgoAdvanceService) SetInstId(instId string) *SubscribeAlgoAdvanceService {
	s.instId = instId
	return s
}

func (s *SubscribeAlgoAdvanceService) Do(ctx context.Context, cb WsHandler[WsAdvanceAlgoOrder]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsAdvanceAlgoOrder](ctx, s.c, request.GatewayBusiness, true,
		request.WsArg{Channel: "algo-advance", InstrumentType: string(s.instType), InstrumentID: s.instId}, cb)
}

// --- grid order channels (private; login) -------------------------------------

// WsGridOrder is a single grid algo order pushed by the "grid-orders-spot",
// "grid-orders-contract" and "grid-orders-moon" channels. The validating account
// had no grid orders, so the field set is a union modeled from the OKX WS doc
// field tables for those channels (spot / contract / moon grid).
type WsGridOrder struct {
	AlgoID                 string             `json:"algoId"`
	AlgoClientOrderID      string             `json:"algoClOrdId"`
	AlgoOrderType          GridAlgoOrdType    `json:"algoOrdType"`
	InstrumentType         InstType           `json:"instType"`
	InstrumentID           string             `json:"instId"`
	CancelType             string             `json:"cancelType"`
	State                  string             `json:"state"`
	RunType                GridRunType        `json:"runType"`
	GridNumber             decimal.Decimal    `json:"gridNum"`
	MaxPrice               decimal.Decimal    `json:"maxPx"`
	MinPrice               decimal.Decimal    `json:"minPx"`
	GridProfit             decimal.Decimal    `json:"gridProfit"`
	TotalPnl               decimal.Decimal    `json:"totalPnl"`
	PnlRatio               decimal.Decimal    `json:"pnlRatio"`
	FloatProfit            decimal.Decimal    `json:"floatProfit"`
	TotalAnnualizedRate    decimal.Decimal    `json:"totalAnnualizedRate"`
	AnnualizedRate         decimal.Decimal    `json:"annualizedRate"`
	Investment             decimal.Decimal    `json:"investment"`
	TakeProfitTriggerPrice decimal.Decimal    `json:"tpTriggerPx"`
	StopLossTriggerPrice   decimal.Decimal    `json:"slTriggerPx"`
	TriggerPrice           decimal.Decimal    `json:"triggerPx"`
	StopType               string             `json:"stopType"`
	StopResult             string             `json:"stopResult"`
	ActiveOrderNumber      decimal.Decimal    `json:"activeOrdNum"`
	Tag                    string             `json:"tag"`
	ProfitSharingRatio     decimal.Decimal    `json:"profitSharingRatio"`
	CopyType               string             `json:"copyType"`
	Fee                    decimal.Decimal    `json:"fee"`
	FundingFee             decimal.Decimal    `json:"fundingFee"`
	RebateTransfer         []GridRebateTrans  `json:"rebateTrans"`
	TriggerParams          []GridTriggerParam `json:"triggerParams"`
	TriggerTime            time.Time          `json:"triggerTime"`
	CreationTime           time.Time          `json:"cTime"`
	UpdateTime             time.Time          `json:"uTime"`

	// --- spot/moon grid ("grid"/"moon_grid") ---
	BaseSize                decimal.Decimal `json:"baseSz"`
	QuoteSize               decimal.Decimal `json:"quoteSz"`
	BaseCurrency            string          `json:"baseCcy"`
	QuoteCurrency           string          `json:"quoteCcy"`
	TradeMode               TdMode          `json:"tdMode"`
	ProfitAndLoss           decimal.Decimal `json:"profitAndLoss"`
	GridArithmeticGeometric string          `json:"gridArithGeo"`
	MinTradeFeeRate         decimal.Decimal `json:"minTradeFeeRate"`

	// --- contract grid ("contract_grid") ---
	Direction          GridDirection   `json:"direction"`
	BasePosition       bool            `json:"basePos"`
	Size               decimal.Decimal `json:"sz"`
	Currency           string          `json:"ccy"`
	Equity             decimal.Decimal `json:"eq"`
	Underlying         string          `json:"uly"`
	InstrumentFamily   string          `json:"instFamily"`
	Leverage           decimal.Decimal `json:"lever"`
	TakeProfitRatio    decimal.Decimal `json:"tpRatio"`
	StopLossRatio      decimal.Decimal `json:"slRatio"`
	AvailableEquity    decimal.Decimal `json:"availEq"`
	LiquidationPrice   decimal.Decimal `json:"liqPx"`
	UPLRatio           decimal.Decimal `json:"uplRatio"`
	UPL                decimal.Decimal `json:"upl"`
	TotalInvestment    decimal.Decimal `json:"totalInvestment"`
	GridInvestment     decimal.Decimal `json:"gridInvestment"`
	MarginRatio        decimal.Decimal `json:"marginRatio"`
	Arbitrage          decimal.Decimal `json:"arbitrage"`
	SingleAmount       decimal.Decimal `json:"singleAmt"`
	PerMaxProfitRate   decimal.Decimal `json:"perMaxProfitRate"`
	PerMinProfitRate   decimal.Decimal `json:"perMinProfitRate"`
	OrderFrozen        decimal.Decimal `json:"ordFrozen"`
	ActualLeverage     decimal.Decimal `json:"actualLever"`
	InvestmentCurrency string          `json:"investmentCcy"`
}

// SubscribeGridOrdersSpotService -- "grid-orders-spot" channel (business; login).
//
// Streams the account's spot grid algo orders for a product line. InstId narrows
// the subscription.
type SubscribeGridOrdersSpotService struct {
	c        *WebSocketClient
	instType InstType
	instId   string
}

func (c *WebSocketClient) NewSubscribeGridOrdersSpotService(instType InstType) *SubscribeGridOrdersSpotService {
	return &SubscribeGridOrdersSpotService{c: c, instType: instType}
}

// SetInstId narrows the subscription to a single instrument.
func (s *SubscribeGridOrdersSpotService) SetInstId(instId string) *SubscribeGridOrdersSpotService {
	s.instId = instId
	return s
}

func (s *SubscribeGridOrdersSpotService) Do(ctx context.Context, cb WsHandler[WsGridOrder]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsGridOrder](ctx, s.c, request.GatewayBusiness, true,
		request.WsArg{Channel: "grid-orders-spot", InstrumentType: string(s.instType), InstrumentID: s.instId}, cb)
}

// SubscribeGridOrdersContractService -- "grid-orders-contract" channel (business; login).
//
// Streams the account's contract grid algo orders for a product line. InstId
// narrows the subscription.
type SubscribeGridOrdersContractService struct {
	c        *WebSocketClient
	instType InstType
	instId   string
}

func (c *WebSocketClient) NewSubscribeGridOrdersContractService(instType InstType) *SubscribeGridOrdersContractService {
	return &SubscribeGridOrdersContractService{c: c, instType: instType}
}

// SetInstId narrows the subscription to a single instrument.
func (s *SubscribeGridOrdersContractService) SetInstId(instId string) *SubscribeGridOrdersContractService {
	s.instId = instId
	return s
}

func (s *SubscribeGridOrdersContractService) Do(ctx context.Context, cb WsHandler[WsGridOrder]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsGridOrder](ctx, s.c, request.GatewayBusiness, true,
		request.WsArg{Channel: "grid-orders-contract", InstrumentType: string(s.instType), InstrumentID: s.instId}, cb)
}

// SubscribeGridOrdersMoonService -- "grid-orders-moon" channel (business; login).
//
// Streams the account's moon grid algo orders for a product line. InstId narrows
// the subscription.
type SubscribeGridOrdersMoonService struct {
	c        *WebSocketClient
	instType InstType
	instId   string
}

func (c *WebSocketClient) NewSubscribeGridOrdersMoonService(instType InstType) *SubscribeGridOrdersMoonService {
	return &SubscribeGridOrdersMoonService{c: c, instType: instType}
}

// SetInstId narrows the subscription to a single instrument.
func (s *SubscribeGridOrdersMoonService) SetInstId(instId string) *SubscribeGridOrdersMoonService {
	s.instId = instId
	return s
}

func (s *SubscribeGridOrdersMoonService) Do(ctx context.Context, cb WsHandler[WsGridOrder]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsGridOrder](ctx, s.c, request.GatewayBusiness, true,
		request.WsArg{Channel: "grid-orders-moon", InstrumentType: string(s.instType), InstrumentID: s.instId}, cb)
}

// WsGridPosition is the position held by a contract grid algo order pushed by the
// "grid-positions" channel. The validating account had no grid orders, so the
// field set is modeled from the OKX WS doc field table.
type WsGridPosition struct {
	AlgoID            string          `json:"algoId"`
	AlgoClientOrderID string          `json:"algoClOrdId"`
	InstrumentType    InstType        `json:"instType"`
	InstrumentID      string          `json:"instId"`
	Currency          string          `json:"ccy"`
	PositionSide      PosSide         `json:"posSide"`
	MarginMode        MgnMode         `json:"mgnMode"`
	Position          decimal.Decimal `json:"pos"`
	AveragePrice      decimal.Decimal `json:"avgPx"`
	LiquidationPrice  decimal.Decimal `json:"liqPx"`
	MarkPrice         decimal.Decimal `json:"markPx"`
	Leverage          decimal.Decimal `json:"lever"`
	IMR               decimal.Decimal `json:"imr"`
	MMR               decimal.Decimal `json:"mmr"`
	MarginRatio       decimal.Decimal `json:"mgnRatio"`
	Margin            decimal.Decimal `json:"margin"`
	NotionalUSD       decimal.Decimal `json:"notionalUsd"`
	Last              decimal.Decimal `json:"last"`
	UPL               decimal.Decimal `json:"upl"`
	UPLRatio          decimal.Decimal `json:"uplRatio"`
	CreationTime      time.Time       `json:"cTime"`
	UpdateTime        time.Time       `json:"uTime"`
}

// SubscribeGridPositionsService -- "grid-positions" channel (business; login).
//
// Streams the open positions of a contract grid algo order, keyed by algoId.
type SubscribeGridPositionsService struct {
	c      *WebSocketClient
	algoId string
}

func (c *WebSocketClient) NewSubscribeGridPositionsService(algoId string) *SubscribeGridPositionsService {
	return &SubscribeGridPositionsService{c: c, algoId: algoId}
}

func (s *SubscribeGridPositionsService) Do(ctx context.Context, cb WsHandler[WsGridPosition]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsGridPosition](ctx, s.c, request.GatewayBusiness, true,
		request.WsArg{Channel: "grid-positions", AlgoID: s.algoId}, cb)
}

// WsGridSubOrder is one sub-order (working leg) of a grid algo order pushed by
// the "grid-sub-orders" channel. The validating account had no grid orders, so
// the field set is modeled from the OKX WS doc field table.
type WsGridSubOrder struct {
	AlgoID              string          `json:"algoId"`
	AlgoClientOrderID   string          `json:"algoClOrdId"`
	AlgoOrderType       GridAlgoOrdType `json:"algoOrdType"`
	InstrumentType      InstType        `json:"instType"`
	InstrumentID        string          `json:"instId"`
	GroupID             string          `json:"groupId"`
	OrderID             string          `json:"ordId"`
	ClientOrderID       string          `json:"clOrdId"`
	Tag                 string          `json:"tag"`
	OrderType           OrdType         `json:"ordType"`
	Side                Side            `json:"side"`
	PositionSide        PosSide         `json:"posSide"`
	TradeMode           TdMode          `json:"tdMode"`
	Currency            string          `json:"ccy"`
	Price               decimal.Decimal `json:"px"`
	Size                decimal.Decimal `json:"sz"`
	State               OrdState        `json:"state"`
	AccumulatedFillSize decimal.Decimal `json:"accFillSz"`
	AveragePrice        decimal.Decimal `json:"avgPx"`
	Leverage            decimal.Decimal `json:"lever"`
	Fee                 decimal.Decimal `json:"fee"`
	FeeCurrency         string          `json:"feeCcy"`
	Rebate              decimal.Decimal `json:"rebate"`
	RebateCurrency      string          `json:"rebateCcy"`
	Pnl                 decimal.Decimal `json:"pnl"`
	CreationTime        time.Time       `json:"cTime"`
	UpdateTime          time.Time       `json:"uTime"`
}

// SubscribeGridSubOrdersService -- "grid-sub-orders" channel (business; login).
//
// Streams the sub-orders (working legs) of a grid algo order, keyed by algoId.
type SubscribeGridSubOrdersService struct {
	c      *WebSocketClient
	algoId string
}

func (c *WebSocketClient) NewSubscribeGridSubOrdersService(algoId string) *SubscribeGridSubOrdersService {
	return &SubscribeGridSubOrdersService{c: c, algoId: algoId}
}

func (s *SubscribeGridSubOrdersService) Do(ctx context.Context, cb WsHandler[WsGridSubOrder]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsGridSubOrder](ctx, s.c, request.GatewayBusiness, true,
		request.WsArg{Channel: "grid-sub-orders", AlgoID: s.algoId}, cb)
}

// --- economic-calendar channel (private; login) -------------------------------

// WsEconomicCalendar is one macro-economic calendar event pushed by the
// "economic-calendar" channel. Access requires a sufficient trading-fee tier; the
// field set matches the REST economic-calendar response.
type WsEconomicCalendar struct {
	CalendarID      string                     `json:"calendarId"`
	Date            time.Time                  `json:"date"`
	Region          string                     `json:"region"`
	Category        string                     `json:"category"`
	Event           string                     `json:"event"`
	ReferenceDate   time.Time                  `json:"refDate"`
	Actual          string                     `json:"actual"`
	Previous        string                     `json:"previous"`
	Forecast        string                     `json:"forecast"`
	DateSpan        string                     `json:"dateSpan"`
	Importance      EconomicCalendarImportance `json:"importance"`
	UpdateTime      time.Time                  `json:"uTime"`
	PreviousInitial string                     `json:"prevInitial"`
	Currency        string                     `json:"ccy"`
	Unit            string                     `json:"unit"`
}

// SubscribeEconomicCalendarService -- "economic-calendar" channel (business; login).
//
// Streams macro-economic calendar event updates. Access requires a sufficient
// trading-fee tier (OKX rejects lower tiers with code 64003).
type SubscribeEconomicCalendarService struct {
	c *WebSocketClient
}

func (c *WebSocketClient) NewSubscribeEconomicCalendarService() *SubscribeEconomicCalendarService {
	return &SubscribeEconomicCalendarService{c: c}
}

func (s *SubscribeEconomicCalendarService) Do(ctx context.Context, cb WsHandler[WsEconomicCalendar]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsEconomicCalendar](ctx, s.c, request.GatewayBusiness, true,
		request.WsArg{Channel: "economic-calendar"}, cb)
}
