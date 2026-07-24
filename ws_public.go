package okx

// This file wraps OKX v5 public WebSocket market channels. Every channel here is
// served on the public gateway (request.GatewayPublic). Most need no login; the
// two VIP-gated "*-l2-tbt" depth channels log in first (private=true) and then
// fail with code 64003 unless the account's fee tier qualifies. Each Ws* struct
// is modeled from a live push (json tags are the exact pushed keys); a struct may
// carry a superset of keys across order-book variants since the verification only
// requires the pushed keys to be covered.

import (
	"context"
	"time"

	"github.com/UnipayFI/go-okx/request"
	"github.com/shopspring/decimal"
)

// WsTicker is the latest market snapshot pushed on the "tickers" channel. It is a
// subset of the REST Ticker (no SodUtc fields differ, but it omits the REST-only
// extras).
type WsTicker struct {
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

// SubscribeTickersService -- "tickers" channel (public; no login).
type SubscribeTickersService struct {
	c      *WebSocketClient
	instId string
}

func (c *WebSocketClient) NewSubscribeTickersService(instId string) *SubscribeTickersService {
	return &SubscribeTickersService{c: c, instId: instId}
}

func (s *SubscribeTickersService) Do(ctx context.Context, cb WsHandler[WsTicker]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsTicker](ctx, s.c, request.GatewayPublic, false,
		request.WsArg{Channel: "tickers", InstrumentID: s.instId}, cb)
}

// WsTrade is a single public taker fill pushed on the "trades" channel.
type WsTrade struct {
	InstrumentID string          `json:"instId"`
	TradeID      string          `json:"tradeId"`
	Price        decimal.Decimal `json:"px"`
	Size         decimal.Decimal `json:"sz"`
	Side         Side            `json:"side"`
	Count        string          `json:"count"`
	// Source is the trade source ("1": RPI order, previously ELP order; see
	// Trade.Source).
	Source     string    `json:"source"`
	SequenceID int64     `json:"seqId"`
	Timestamp  time.Time `json:"ts"`
}

// SubscribeTradesService -- "trades" channel (public; no login).
type SubscribeTradesService struct {
	c      *WebSocketClient
	instId string
}

func (c *WebSocketClient) NewSubscribeTradesService(instId string) *SubscribeTradesService {
	return &SubscribeTradesService{c: c, instId: instId}
}

func (s *SubscribeTradesService) Do(ctx context.Context, cb WsHandler[WsTrade]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsTrade](ctx, s.c, request.GatewayPublic, false,
		request.WsArg{Channel: "trades", InstrumentID: s.instId}, cb)
}

// WsOrderBook is an order-book frame. It is the superset shape across the
// order-book channels: "books" carries asks/bids/ts/checksum/seqId/prevSeqId
// (with Action "snapshot"/"update" on the push); "books5" carries
// asks/bids/instId/ts/seqId; "bbo-tbt" carries asks/bids/ts/seqId; the VIP
// "*-l2-tbt" channels carry the same as "books". Each level is
// [price, size, deprecated("0"), numOrders] as strings.
type WsOrderBook struct {
	Asks               [][]string `json:"asks"`
	Bids               [][]string `json:"bids"`
	InstrumentID       string     `json:"instId,omitempty"`
	Timestamp          time.Time  `json:"ts"`
	Checksum           int64      `json:"checksum,omitempty"`
	PreviousSequenceID int64      `json:"prevSeqId,omitempty"`
	SequenceID         int64      `json:"seqId"`
}

// SubscribeBooksService -- "books" channel (public; no login). 400-level book:
// the first push is a "snapshot" (push.Action), subsequent pushes are "update".
type SubscribeBooksService struct {
	c      *WebSocketClient
	instId string
}

func (c *WebSocketClient) NewSubscribeBooksService(instId string) *SubscribeBooksService {
	return &SubscribeBooksService{c: c, instId: instId}
}

func (s *SubscribeBooksService) Do(ctx context.Context, cb WsHandler[WsOrderBook]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsOrderBook](ctx, s.c, request.GatewayPublic, false,
		request.WsArg{Channel: "books", InstrumentID: s.instId}, cb)
}

// SubscribeBooks5Service -- "books5" channel (public; no login). 5-level
// snapshots pushed at ~100ms.
type SubscribeBooks5Service struct {
	c      *WebSocketClient
	instId string
}

func (c *WebSocketClient) NewSubscribeBooks5Service(instId string) *SubscribeBooks5Service {
	return &SubscribeBooks5Service{c: c, instId: instId}
}

func (s *SubscribeBooks5Service) Do(ctx context.Context, cb WsHandler[WsOrderBook]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsOrderBook](ctx, s.c, request.GatewayPublic, false,
		request.WsArg{Channel: "books5", InstrumentID: s.instId}, cb)
}

// SubscribeBboTbtService -- "bbo-tbt" channel (public; no login). Tick-by-tick
// best bid/offer (1-level book), pushed at ~10ms.
type SubscribeBboTbtService struct {
	c      *WebSocketClient
	instId string
}

func (c *WebSocketClient) NewSubscribeBboTbtService(instId string) *SubscribeBboTbtService {
	return &SubscribeBboTbtService{c: c, instId: instId}
}

func (s *SubscribeBboTbtService) Do(ctx context.Context, cb WsHandler[WsOrderBook]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsOrderBook](ctx, s.c, request.GatewayPublic, false,
		request.WsArg{Channel: "bbo-tbt", InstrumentID: s.instId}, cb)
}

// SubscribeBooks50L2TbtService -- "books50-l2-tbt" channel (public gateway;
// login required). Tick-by-tick 50-level book, gated behind a VIP fee tier;
// without it the subscribe returns code 64003. Action is "snapshot"/"update".
type SubscribeBooks50L2TbtService struct {
	c      *WebSocketClient
	instId string
}

func (c *WebSocketClient) NewSubscribeBooks50L2TbtService(instId string) *SubscribeBooks50L2TbtService {
	return &SubscribeBooks50L2TbtService{c: c, instId: instId}
}

func (s *SubscribeBooks50L2TbtService) Do(ctx context.Context, cb WsHandler[WsOrderBook]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsOrderBook](ctx, s.c, request.GatewayPublic, true,
		request.WsArg{Channel: "books50-l2-tbt", InstrumentID: s.instId}, cb)
}

// SubscribeBooksL2TbtService -- "books-l2-tbt" channel (public gateway; login
// required). Tick-by-tick 400-level book, gated behind a VIP fee tier; without it
// the subscribe returns code 64003. Action is "snapshot"/"update".
type SubscribeBooksL2TbtService struct {
	c      *WebSocketClient
	instId string
}

func (c *WebSocketClient) NewSubscribeBooksL2TbtService(instId string) *SubscribeBooksL2TbtService {
	return &SubscribeBooksL2TbtService{c: c, instId: instId}
}

func (s *SubscribeBooksL2TbtService) Do(ctx context.Context, cb WsHandler[WsOrderBook]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsOrderBook](ctx, s.c, request.GatewayPublic, true,
		request.WsArg{Channel: "books-l2-tbt", InstrumentID: s.instId}, cb)
}

// WsInstrument is an instrument definition pushed on the "instruments" channel
// (on a listing/rule change). It mirrors the REST Instrument's pushed subset.
type WsInstrument struct {
	InstrumentType            InstType        `json:"instType"`
	InstrumentID              string          `json:"instId"`
	InstrumentIDCode          int64           `json:"instIdCode"`
	Underlying                string          `json:"uly"`
	InstrumentFamily          string          `json:"instFamily"`
	InstrumentCategory        string          `json:"instCategory"`
	BaseCurrency              string          `json:"baseCcy"`
	QuoteCurrency             string          `json:"quoteCcy"`
	SettleCurrency            string          `json:"settleCcy"`
	ContractValue             decimal.Decimal `json:"ctVal"`
	ContractMultiplier        decimal.Decimal `json:"ctMult"`
	ContractValueCurrency     string          `json:"ctValCcy"`
	OptionType                OptType         `json:"optType"`
	Strike                    decimal.Decimal `json:"stk"`
	ListTime                  time.Time       `json:"listTime"`
	AuctionEndTime            time.Time       `json:"auctionEndTime"`
	ContinuousTradeSwitchTime time.Time       `json:"contTdSwTime"`
	OpenType                  string          `json:"openType"`
	ExpiryTime                time.Time       `json:"expTime"`
	Leverage                  decimal.Decimal `json:"lever"`
	TickSize                  decimal.Decimal `json:"tickSz"`
	LotSize                   decimal.Decimal `json:"lotSz"`
	MinSize                   decimal.Decimal `json:"minSz"`
	ContractType              CtType          `json:"ctType"`
	Alias                     string          `json:"alias"`
	State                     InstState       `json:"state"`
	// RuleType is the trading rule type: "normal", "pre_market" (including
	// pre-market X-Perp FUTURES), "rebase_contract" and "xperp" (a pre-market
	// X-Perp changes from "pre_market" to "xperp" after converting to a normal
	// X-Perp).
	RuleType                         string          `json:"ruleType"`
	MaxLimitSize                     decimal.Decimal `json:"maxLmtSz"`
	MaxMarketSize                    decimal.Decimal `json:"maxMktSz"`
	MaxLimitAmount                   decimal.Decimal `json:"maxLmtAmt"`
	MaxMarketAmount                  decimal.Decimal `json:"maxMktAmt"`
	MaxTWAPSize                      decimal.Decimal `json:"maxTwapSz"`
	MaxIcebergSize                   decimal.Decimal `json:"maxIcebergSz"`
	MaxTriggerSize                   decimal.Decimal `json:"maxTriggerSz"`
	MaxStopSize                      decimal.Decimal `json:"maxStopSz"`
	MaxPlatformOpenInterestLimit     decimal.Decimal `json:"maxPlatOILmt"`
	MaxPlatformOpenInterestCoinLimit decimal.Decimal `json:"maxPlatOICoinLmt"`
	FutureSettlement                 bool            `json:"futureSettlement"`
	TradeQuoteCurrencyList           []string        `json:"tradeQuoteCcyList"`
	GroupID                          string          `json:"groupId"`
	PositionLimitAmount              decimal.Decimal `json:"posLmtAmt"`
	PositionLimitPercent             decimal.Decimal `json:"posLmtPct"`
	// PreMarketSwitchTime is the time a pre-market instrument switched to normal
	// trading. Applicable to pre-market SWAP and pre-market X-Perp FUTURES.
	PreMarketSwitchTime time.Time `json:"preMktSwTime"`
	// Elp is the ELP (Enhanced Liquidity Program) maker permission (values
	// "0"/"1"/"2"; see Instrument.Elp). OKX is rebranding ELP to RPI (Retail
	// Price Improvement); the json key stays "elp" until the old names retire on
	// 2026-10-31, after which it becomes "rpi".
	Elp string `json:"elp"`
}

// SubscribeInstrumentsService -- "instruments" channel (public; no login).
type SubscribeInstrumentsService struct {
	c        *WebSocketClient
	instType InstType
}

func (c *WebSocketClient) NewSubscribeInstrumentsService(instType InstType) *SubscribeInstrumentsService {
	return &SubscribeInstrumentsService{c: c, instType: instType}
}

func (s *SubscribeInstrumentsService) Do(ctx context.Context, cb WsHandler[WsInstrument]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsInstrument](ctx, s.c, request.GatewayPublic, false,
		request.WsArg{Channel: "instruments", InstrumentType: string(s.instType)}, cb)
}

// WsOpenInterest is an instrument's open interest pushed on the "open-interest"
// channel.
type WsOpenInterest struct {
	InstrumentType       InstType        `json:"instType"`
	InstrumentID         string          `json:"instId"`
	OpenInterest         decimal.Decimal `json:"oi"`
	OpenInterestCurrency decimal.Decimal `json:"oiCcy"`
	OpenInterestUSD      decimal.Decimal `json:"oiUsd"`
	Timestamp            time.Time       `json:"ts"`
}

// SubscribeOpenInterestService -- "open-interest" channel (public; no login).
type SubscribeOpenInterestService struct {
	c      *WebSocketClient
	instId string
}

func (c *WebSocketClient) NewSubscribeOpenInterestService(instId string) *SubscribeOpenInterestService {
	return &SubscribeOpenInterestService{c: c, instId: instId}
}

func (s *SubscribeOpenInterestService) Do(ctx context.Context, cb WsHandler[WsOpenInterest]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsOpenInterest](ctx, s.c, request.GatewayPublic, false,
		request.WsArg{Channel: "open-interest", InstrumentID: s.instId}, cb)
}

// WsFundingRate is a perpetual swap's funding rate pushed on the "funding-rate"
// channel.
type WsFundingRate struct {
	InstrumentType        InstType        `json:"instType"`
	InstrumentID          string          `json:"instId"`
	Method                string          `json:"method"`
	FormulaType           string          `json:"formulaType"`
	FundingRate           decimal.Decimal `json:"fundingRate"`
	NextFundingRate       decimal.Decimal `json:"nextFundingRate"`
	FundingTime           time.Time       `json:"fundingTime"`
	NextFundingTime       time.Time       `json:"nextFundingTime"`
	MinFundingRate        decimal.Decimal `json:"minFundingRate"`
	MaxFundingRate        decimal.Decimal `json:"maxFundingRate"`
	InterestRate          decimal.Decimal `json:"interestRate"`
	ImpactValue           decimal.Decimal `json:"impactValue"`
	SettlementState       string          `json:"settState"`
	SettlementFundingRate decimal.Decimal `json:"settFundingRate"`
	Premium               decimal.Decimal `json:"premium"`
	PreviousFundingTime   time.Time       `json:"prevFundingTime"`
	Timestamp             time.Time       `json:"ts"`
}

// SubscribeFundingRateService -- "funding-rate" channel (public; no login).
type SubscribeFundingRateService struct {
	c      *WebSocketClient
	instId string
}

func (c *WebSocketClient) NewSubscribeFundingRateService(instId string) *SubscribeFundingRateService {
	return &SubscribeFundingRateService{c: c, instId: instId}
}

func (s *SubscribeFundingRateService) Do(ctx context.Context, cb WsHandler[WsFundingRate]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsFundingRate](ctx, s.c, request.GatewayPublic, false,
		request.WsArg{Channel: "funding-rate", InstrumentID: s.instId}, cb)
}

// WsPriceLimit is an instrument's buy/sell price limits pushed on the
// "price-limit" channel.
type WsPriceLimit struct {
	InstrumentType InstType        `json:"instType"`
	InstrumentID   string          `json:"instId"`
	BuyLimit       decimal.Decimal `json:"buyLmt"`
	SellLimit      decimal.Decimal `json:"sellLmt"`
	Enabled        bool            `json:"enabled"`
	Timestamp      time.Time       `json:"ts"`
}

// SubscribePriceLimitService -- "price-limit" channel (public; no login).
type SubscribePriceLimitService struct {
	c      *WebSocketClient
	instId string
}

func (c *WebSocketClient) NewSubscribePriceLimitService(instId string) *SubscribePriceLimitService {
	return &SubscribePriceLimitService{c: c, instId: instId}
}

func (s *SubscribePriceLimitService) Do(ctx context.Context, cb WsHandler[WsPriceLimit]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsPriceLimit](ctx, s.c, request.GatewayPublic, false,
		request.WsArg{Channel: "price-limit", InstrumentID: s.instId}, cb)
}

// WsMarkPrice is an instrument's mark price pushed on the "mark-price" channel.
type WsMarkPrice struct {
	InstrumentType InstType        `json:"instType"`
	InstrumentID   string          `json:"instId"`
	MarkPrice      decimal.Decimal `json:"markPx"`
	Timestamp      time.Time       `json:"ts"`
}

// SubscribeMarkPriceService -- "mark-price" channel (public; no login).
type SubscribeMarkPriceService struct {
	c      *WebSocketClient
	instId string
}

func (c *WebSocketClient) NewSubscribeMarkPriceService(instId string) *SubscribeMarkPriceService {
	return &SubscribeMarkPriceService{c: c, instId: instId}
}

func (s *SubscribeMarkPriceService) Do(ctx context.Context, cb WsHandler[WsMarkPrice]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsMarkPrice](ctx, s.c, request.GatewayPublic, false,
		request.WsArg{Channel: "mark-price", InstrumentID: s.instId}, cb)
}

// WsIndexTicker is an index's latest snapshot pushed on the "index-tickers"
// channel. idxPx is the index price; there is no last-trade or volume.
type WsIndexTicker struct {
	InstrumentID   string          `json:"instId"`
	IndexPrice     decimal.Decimal `json:"idxPx"`
	Open24h        decimal.Decimal `json:"open24h"`
	High24h        decimal.Decimal `json:"high24h"`
	Low24h         decimal.Decimal `json:"low24h"`
	StartOfDayUTC0 decimal.Decimal `json:"sodUtc0"`
	StartOfDayUTC8 decimal.Decimal `json:"sodUtc8"`
	Timestamp      time.Time       `json:"ts"`
}

// SubscribeIndexTickersService -- "index-tickers" channel (public; no login).
type SubscribeIndexTickersService struct {
	c      *WebSocketClient
	instId string
}

func (c *WebSocketClient) NewSubscribeIndexTickersService(instId string) *SubscribeIndexTickersService {
	return &SubscribeIndexTickersService{c: c, instId: instId}
}

func (s *SubscribeIndexTickersService) Do(ctx context.Context, cb WsHandler[WsIndexTicker]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsIndexTicker](ctx, s.c, request.GatewayPublic, false,
		request.WsArg{Channel: "index-tickers", InstrumentID: s.instId}, cb)
}

// WsStatus is a system maintenance event pushed on the "status" channel.
type WsStatus struct {
	Title               string    `json:"title"`
	State               string    `json:"state"`
	Begin               time.Time `json:"begin"`
	End                 time.Time `json:"end"`
	PreOpenBegin        time.Time `json:"preOpenBegin"`
	Href                string    `json:"href"`
	ServiceType         string    `json:"serviceType"`
	System              string    `json:"system"`
	ScheduleDescription string    `json:"scheDesc"`
	MaintenanceType     string    `json:"maintType"`
	Env                 string    `json:"env"`
	Timestamp           time.Time `json:"ts"`
}

// SubscribeStatusService -- "status" channel (public; no login). System
// maintenance announcements; pushes only around maintenance windows.
type SubscribeStatusService struct {
	c *WebSocketClient
}

func (c *WebSocketClient) NewSubscribeStatusService() *SubscribeStatusService {
	return &SubscribeStatusService{c: c}
}

func (s *SubscribeStatusService) Do(ctx context.Context, cb WsHandler[WsStatus]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsStatus](ctx, s.c, request.GatewayPublic, false,
		request.WsArg{Channel: "status"}, cb)
}

// WsLiquidationOrder groups an instrument's recent forced-liquidation fills
// pushed on the "liquidation-orders" channel.
type WsLiquidationOrder struct {
	InstrumentType   InstType                   `json:"instType"`
	InstrumentID     string                     `json:"instId"`
	InstrumentFamily string                     `json:"instFamily"`
	Underlying       string                     `json:"uly"`
	Details          []WsLiquidationOrderDetail `json:"details"`
}

// WsLiquidationOrderDetail is one forced-liquidation fill.
type WsLiquidationOrderDetail struct {
	Side            Side            `json:"side"`
	PositionSide    PosSide         `json:"posSide"`
	BankruptcyPrice decimal.Decimal `json:"bkPx"`
	Size            decimal.Decimal `json:"sz"`
	BankruptcyLoss  decimal.Decimal `json:"bkLoss"`
	Currency        string          `json:"ccy"`
	Timestamp       time.Time       `json:"ts"`
}

// SubscribeLiquidationOrdersService -- "liquidation-orders" channel (public; no
// login). Forced-liquidation fills for a product line; infrequent.
type SubscribeLiquidationOrdersService struct {
	c        *WebSocketClient
	instType InstType
}

func (c *WebSocketClient) NewSubscribeLiquidationOrdersService(instType InstType) *SubscribeLiquidationOrdersService {
	return &SubscribeLiquidationOrdersService{c: c, instType: instType}
}

func (s *SubscribeLiquidationOrdersService) Do(ctx context.Context, cb WsHandler[WsLiquidationOrder]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsLiquidationOrder](ctx, s.c, request.GatewayPublic, false,
		request.WsArg{Channel: "liquidation-orders", InstrumentType: string(s.instType)}, cb)
}
