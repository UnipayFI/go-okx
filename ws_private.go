package okx

import (
	"context"
	"time"

	"github.com/UnipayFI/go-okx/request"
	"github.com/shopspring/decimal"
)

// This file wraps OKX v5 private-gateway WebSocket channels (account / position /
// order state). All require login (request.GatewayPrivate, private=true). The
// account, positions and balance_and_position channels push a full snapshot right
// after login; orders, liquidation-warning and account-greeks are activity-driven.
// The Ws* structs are modeled from the LIVE pushes (a subset/variant of the REST
// shapes) and the OKX field tables for channels with no live data on the
// validating account.

// SubscribeAccountService -- "account" channel (private; login).
//
// Pushes the trading-account equity snapshot on login, then on balance change.
type SubscribeAccountService struct {
	c           *WebSocketClient
	ccy         string
	extraParams string
}

func (c *WebSocketClient) NewSubscribeAccountService() *SubscribeAccountService {
	return &SubscribeAccountService{c: c}
}

// SetCcy limits the push to a single settlement currency.
func (s *SubscribeAccountService) SetCcy(ccy string) *SubscribeAccountService {
	s.ccy = ccy
	return s
}

// SetExtraParams sets the channel's extraParams JSON (e.g. `{"updateInterval":"0"}`
// to push on every change rather than on the default interval).
func (s *SubscribeAccountService) SetExtraParams(extraParams string) *SubscribeAccountService {
	s.extraParams = extraParams
	return s
}

func (s *SubscribeAccountService) Do(ctx context.Context, cb WsHandler[WsAccount]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsAccount](ctx, s.c, request.GatewayPrivate, true,
		request.WsArg{Channel: "account", Currency: s.ccy, ExtraParams: s.extraParams}, cb)
}

// WsAccount is one "account" channel push: account-level equity plus per-currency
// details. It mirrors the REST Balance but is the WebSocket variant (it adds
// coinUsdPrice/frpType per detail and omits some REST-only fields).
type WsAccount struct {
	UpdateTime            time.Time         `json:"uTime"`
	TotalEquity           decimal.Decimal   `json:"totalEq"`
	IsolatedEquity        decimal.Decimal   `json:"isoEq"`
	AdjustedEquity        decimal.Decimal   `json:"adjEq"`
	AvailableEquity       decimal.Decimal   `json:"availEq"`
	OrderFrozen           decimal.Decimal   `json:"ordFroz"`
	IMR                   decimal.Decimal   `json:"imr"`
	MMR                   decimal.Decimal   `json:"mmr"`
	BorrowFrozen          decimal.Decimal   `json:"borrowFroz"`
	MarginRatio           decimal.Decimal   `json:"mgnRatio"`
	NotionalUSD           decimal.Decimal   `json:"notionalUsd"`
	NotionalUSDForBorrow  decimal.Decimal   `json:"notionalUsdForBorrow"`
	NotionalUSDForSwap    decimal.Decimal   `json:"notionalUsdForSwap"`
	NotionalUSDForFutures decimal.Decimal   `json:"notionalUsdForFutures"`
	NotionalUSDForOption  decimal.Decimal   `json:"notionalUsdForOption"`
	UPL                   decimal.Decimal   `json:"upl"`
	Delta                 decimal.Decimal   `json:"delta"`
	DeltaLeverage         decimal.Decimal   `json:"deltaLever"`
	DeltaNeutralStatus    string            `json:"deltaNeutralStatus"`
	Details               []WsAccountDetail `json:"details"`
}

// WsAccountDetail is one currency's balance within an "account" push.
type WsAccountDetail struct {
	Currency                       string          `json:"ccy"`
	Equity                         decimal.Decimal `json:"eq"`
	CashBalance                    decimal.Decimal `json:"cashBal"`
	UpdateTime                     time.Time       `json:"uTime"`
	IsolatedEquity                 decimal.Decimal `json:"isoEq"`
	AvailableEquity                decimal.Decimal `json:"availEq"`
	DiscountEquity                 decimal.Decimal `json:"disEq"`
	FixedBalance                   decimal.Decimal `json:"fixedBal"`
	AvailableBalance               decimal.Decimal `json:"availBal"`
	FrozenBalance                  decimal.Decimal `json:"frozenBal"`
	OrderFrozen                    decimal.Decimal `json:"ordFrozen"`
	Liability                      decimal.Decimal `json:"liab"`
	UPL                            decimal.Decimal `json:"upl"`
	UPLLiability                   decimal.Decimal `json:"uplLiab"`
	CrossLiability                 decimal.Decimal `json:"crossLiab"`
	IsolatedLiability              decimal.Decimal `json:"isoLiab"`
	RewardBalance                  decimal.Decimal `json:"rewardBal"`
	MarginRatio                    decimal.Decimal `json:"mgnRatio"`
	IMR                            decimal.Decimal `json:"imr"`
	MMR                            decimal.Decimal `json:"mmr"`
	Interest                       decimal.Decimal `json:"interest"`
	TWAP                           decimal.Decimal `json:"twap"`
	MaxLoan                        decimal.Decimal `json:"maxLoan"`
	EquityUSD                      decimal.Decimal `json:"eqUsd"`
	CoinUSDPrice                   decimal.Decimal `json:"coinUsdPrice"`
	BorrowFrozen                   decimal.Decimal `json:"borrowFroz"`
	NotionalLeverage               decimal.Decimal `json:"notionalLever"`
	StrategyEquity                 decimal.Decimal `json:"stgyEq"`
	IsolatedUPL                    decimal.Decimal `json:"isoUpl"`
	SpotInUseAmount                decimal.Decimal `json:"spotInUseAmt"`
	ClientSpotInUseAmount          decimal.Decimal `json:"clSpotInUseAmt"`
	MaxSpotInUseAmount             decimal.Decimal `json:"maxSpotInUseAmt"`
	SpotIsolatedBalance            decimal.Decimal `json:"spotIsoBal"`
	SmtSyncEquity                  decimal.Decimal `json:"smtSyncEq"`
	SpotCopyTradingEquity          decimal.Decimal `json:"spotCopyTradingEq"`
	SpotBalance                    decimal.Decimal `json:"spotBal"`
	OpenAveragePrice               decimal.Decimal `json:"openAvgPx"`
	AccumulatedAveragePrice        decimal.Decimal `json:"accAvgPx"`
	SpotUPL                        decimal.Decimal `json:"spotUpl"`
	SpotUPLRatio                   decimal.Decimal `json:"spotUplRatio"`
	TotalPnl                       decimal.Decimal `json:"totalPnl"`
	TotalPnlRatio                  decimal.Decimal `json:"totalPnlRatio"`
	CollateralEnabled              bool            `json:"collateralEnabled"`
	CollateralRestrict             bool            `json:"collateralRestrict"`
	CollateralBorrowAutoConversion decimal.Decimal `json:"colBorrAutoConversion"`
	ColRes                         string          `json:"colRes"`
	AutoLendStatus                 string          `json:"autoLendStatus"`
	AutoLendAmount                 decimal.Decimal `json:"autoLendAmt"`
	AutoLendMatchedAmount          decimal.Decimal `json:"autoLendMtAmt"`
	AutoStakingStatus              string          `json:"autoStakingStatus"`
	FrpType                        string          `json:"frpType"`
}

// SubscribePositionsService -- "positions" channel (private; login).
//
// Pushes the open-positions snapshot on login (empty array when none), then on
// every position change. instType is required; pass InstTypeAny for all.
type SubscribePositionsService struct {
	c          *WebSocketClient
	instType   InstType
	instFamily string
	instId     string
}

func (c *WebSocketClient) NewSubscribePositionsService(instType InstType) *SubscribePositionsService {
	return &SubscribePositionsService{c: c, instType: instType}
}

// SetInstFamily narrows the push to one instrument family.
func (s *SubscribePositionsService) SetInstFamily(instFamily string) *SubscribePositionsService {
	s.instFamily = instFamily
	return s
}

// SetInstId narrows the push to one instrument.
func (s *SubscribePositionsService) SetInstId(instId string) *SubscribePositionsService {
	s.instId = instId
	return s
}

func (s *SubscribePositionsService) Do(ctx context.Context, cb WsHandler[WsPosition]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsPosition](ctx, s.c, request.GatewayPrivate, true,
		request.WsArg{Channel: "positions", InstrumentType: string(s.instType), InstrumentFamily: s.instFamily, InstrumentID: s.instId}, cb)
}

// WsPosition is one open position from the "positions" channel. The validating
// account had no open positions, so the field set is modeled from the OKX
// "positions" channel field table (it mirrors the REST Position shape).
type WsPosition struct {
	InstrumentType                  InstType                 `json:"instType"`
	MarginMode                      MgnMode                  `json:"mgnMode"`
	PositionID                      string                   `json:"posId"`
	PositionSide                    PosSide                  `json:"posSide"`
	Position                        decimal.Decimal          `json:"pos"`
	BaseBalance                     decimal.Decimal          `json:"baseBal"`
	QuoteBalance                    decimal.Decimal          `json:"quoteBal"`
	BaseBorrowed                    decimal.Decimal          `json:"baseBorrowed"`
	BaseInterest                    decimal.Decimal          `json:"baseInterest"`
	QuoteBorrowed                   decimal.Decimal          `json:"quoteBorrowed"`
	QuoteInterest                   decimal.Decimal          `json:"quoteInterest"`
	PositionCurrency                string                   `json:"posCcy"`
	AvailablePosition               decimal.Decimal          `json:"availPos"`
	AveragePrice                    decimal.Decimal          `json:"avgPx"`
	NonSettleAveragePrice           decimal.Decimal          `json:"nonSettleAvgPx"`
	MarkPrice                       decimal.Decimal          `json:"markPx"`
	UPL                             decimal.Decimal          `json:"upl"`
	UPLRatio                        decimal.Decimal          `json:"uplRatio"`
	UPLLastPrice                    decimal.Decimal          `json:"uplLastPx"`
	UPLRatioLastPrice               decimal.Decimal          `json:"uplRatioLastPx"`
	InstrumentID                    string                   `json:"instId"`
	Leverage                        decimal.Decimal          `json:"lever"`
	LiquidationPrice                decimal.Decimal          `json:"liqPx"`
	IMR                             decimal.Decimal          `json:"imr"`
	Margin                          decimal.Decimal          `json:"margin"`
	MarginRatio                     decimal.Decimal          `json:"mgnRatio"`
	MMR                             decimal.Decimal          `json:"mmr"`
	Liability                       decimal.Decimal          `json:"liab"`
	LiabilityCurrency               string                   `json:"liabCcy"`
	Interest                        decimal.Decimal          `json:"interest"`
	TradeID                         string                   `json:"tradeId"`
	OptionValue                     decimal.Decimal          `json:"optVal"`
	PendingCloseOrderLiabilityValue decimal.Decimal          `json:"pendingCloseOrdLiabVal"`
	NotionalUSD                     decimal.Decimal          `json:"notionalUsd"`
	ADL                             decimal.Decimal          `json:"adl"`
	Currency                        string                   `json:"ccy"`
	Last                            decimal.Decimal          `json:"last"`
	IndexPrice                      decimal.Decimal          `json:"idxPx"`
	USDPrice                        decimal.Decimal          `json:"usdPx"`
	BreakEvenPrice                  decimal.Decimal          `json:"bePx"`
	DeltaBS                         decimal.Decimal          `json:"deltaBS"`
	DeltaPA                         decimal.Decimal          `json:"deltaPA"`
	GammaBS                         decimal.Decimal          `json:"gammaBS"`
	GammaPA                         decimal.Decimal          `json:"gammaPA"`
	ThetaBS                         decimal.Decimal          `json:"thetaBS"`
	ThetaPA                         decimal.Decimal          `json:"thetaPA"`
	VegaBS                          decimal.Decimal          `json:"vegaBS"`
	VegaPA                          decimal.Decimal          `json:"vegaPA"`
	SpotInUseAmount                 decimal.Decimal          `json:"spotInUseAmt"`
	SpotInUseCurrency               string                   `json:"spotInUseCcy"`
	ClientSpotInUseAmount           decimal.Decimal          `json:"clSpotInUseAmt"`
	MaxSpotInUseAmount              decimal.Decimal          `json:"maxSpotInUseAmt"`
	RealizedPnl                     decimal.Decimal          `json:"realizedPnl"`
	SettledPnl                      decimal.Decimal          `json:"settledPnl"`
	Pnl                             decimal.Decimal          `json:"pnl"`
	Fee                             decimal.Decimal          `json:"fee"`
	FundingFee                      decimal.Decimal          `json:"fundingFee"`
	LiquidationPenalty              decimal.Decimal          `json:"liqPenalty"`
	CloseOrderAlgo                  []PositionCloseOrderAlgo `json:"closeOrderAlgo"`
	CreationTime                    time.Time                `json:"cTime"`
	UpdateTime                      time.Time                `json:"uTime"`
	PushTime                        time.Time                `json:"pTime"`
	BusinessReferenceID             string                   `json:"bizRefId"`
	BusinessReferenceType           string                   `json:"bizRefType"`
}

// SubscribeBalanceAndPositionService -- "balance_and_position" channel (private; login).
//
// Pushes a combined balance + position snapshot on login, then incremental
// updates (with the originating trades) on every balance/position change.
type SubscribeBalanceAndPositionService struct {
	c *WebSocketClient
}

func (c *WebSocketClient) NewSubscribeBalanceAndPositionService() *SubscribeBalanceAndPositionService {
	return &SubscribeBalanceAndPositionService{c: c}
}

func (s *SubscribeBalanceAndPositionService) Do(ctx context.Context, cb WsHandler[WsBalanceAndPosition]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsBalanceAndPosition](ctx, s.c, request.GatewayPrivate, true,
		request.WsArg{Channel: "balance_and_position"}, cb)
}

// WsBalanceAndPosition is one "balance_and_position" push: the publish time, the
// event type ("snapshot"/"delivered"/"exercised"/"transferred"/"filled"/...),
// the changed balances, the changed positions and the triggering trades.
type WsBalanceAndPosition struct {
	PushTime     time.Time       `json:"pTime"`
	EventType    string          `json:"eventType"`
	BalanceData  []WsBalData     `json:"balData"`
	PositionData []WsPosition    `json:"posData"`
	Trades       []WsBalPosTrade `json:"trades"`
}

// WsBalData is one currency's balance change within a "balance_and_position" push.
type WsBalData struct {
	Currency    string          `json:"ccy"`
	CashBalance decimal.Decimal `json:"cashBal"`
	UpdateTime  time.Time       `json:"uTime"`
}

// WsBalPosTrade is one trade that triggered a "balance_and_position" update. The
// validating account had no live trades, so the field set is modeled from the
// OKX channel field table.
type WsBalPosTrade struct {
	InstrumentID string `json:"instId"`
	TradeID      string `json:"tradeId"`
}

// SubscribeOrdersService -- "orders" channel (private; login).
//
// Pushes on order create / amend / fill / cancel. Activity-driven (no snapshot).
type SubscribeOrdersService struct {
	c          *WebSocketClient
	instType   InstType
	instFamily string
	instId     string
}

func (c *WebSocketClient) NewSubscribeOrdersService(instType InstType) *SubscribeOrdersService {
	return &SubscribeOrdersService{c: c, instType: instType}
}

// SetInstFamily narrows the push to one instrument family.
func (s *SubscribeOrdersService) SetInstFamily(instFamily string) *SubscribeOrdersService {
	s.instFamily = instFamily
	return s
}

// SetInstId narrows the push to one instrument.
func (s *SubscribeOrdersService) SetInstId(instId string) *SubscribeOrdersService {
	s.instId = instId
	return s
}

func (s *SubscribeOrdersService) Do(ctx context.Context, cb WsHandler[WsOrder]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsOrder](ctx, s.c, request.GatewayPrivate, true,
		request.WsArg{Channel: "orders", InstrumentType: string(s.instType), InstrumentFamily: s.instFamily, InstrumentID: s.instId}, cb)
}

// WsOrder is one order update from the "orders" channel. The validating account
// had no order activity, so the field set is modeled from the OKX "orders"
// channel field table (it mirrors the REST Order shape plus the WS-only fill /
// notify fields).
type WsOrder struct {
	InstrumentType             InstType        `json:"instType"`
	InstrumentID               string          `json:"instId"`
	TargetCurrency             TgtCcy          `json:"tgtCcy"`
	Currency                   string          `json:"ccy"`
	OrderID                    string          `json:"ordId"`
	ClientOrderID              string          `json:"clOrdId"`
	Tag                        string          `json:"tag"`
	Price                      decimal.Decimal `json:"px"`
	PriceUSD                   decimal.Decimal `json:"pxUsd"`
	PriceVolatility            decimal.Decimal `json:"pxVol"`
	PriceType                  string          `json:"pxType"`
	Size                       decimal.Decimal `json:"sz"`
	NotionalUSD                decimal.Decimal `json:"notionalUsd"`
	OrderType                  OrdType         `json:"ordType"`
	Side                       Side            `json:"side"`
	PositionSide               PosSide         `json:"posSide"`
	TradeMode                  TdMode          `json:"tdMode"`
	FillPrice                  decimal.Decimal `json:"fillPx"`
	TradeID                    string          `json:"tradeId"`
	FillSize                   decimal.Decimal `json:"fillSz"`
	FillPnl                    decimal.Decimal `json:"fillPnl"`
	FillTime                   time.Time       `json:"fillTime"`
	FillFee                    decimal.Decimal `json:"fillFee"`
	FillFeeCurrency            string          `json:"fillFeeCcy"`
	FillPriceVolatility        decimal.Decimal `json:"fillPxVol"`
	FillPriceUSD               decimal.Decimal `json:"fillPxUsd"`
	FillMarkVolatility         decimal.Decimal `json:"fillMarkVol"`
	FillForwardPrice           decimal.Decimal `json:"fillFwdPx"`
	FillMarkPrice              decimal.Decimal `json:"fillMarkPx"`
	ExecutionType              ExecType        `json:"execType"`
	AccumulatedFillSize        decimal.Decimal `json:"accFillSz"`
	FillNotionalUSD            decimal.Decimal `json:"fillNotionalUsd"`
	AveragePrice               decimal.Decimal `json:"avgPx"`
	State                      OrdState        `json:"state"`
	Leverage                   decimal.Decimal `json:"lever"`
	AttachAlgoClientOrderID    string          `json:"attachAlgoClOrdId"`
	TakeProfitTriggerPrice     decimal.Decimal `json:"tpTriggerPx"`
	TakeProfitTriggerPriceType string          `json:"tpTriggerPxType"`
	TakeProfitOrderPrice       decimal.Decimal `json:"tpOrdPx"`
	StopLossTriggerPrice       decimal.Decimal `json:"slTriggerPx"`
	StopLossTriggerPriceType   string          `json:"slTriggerPxType"`
	StopLossOrderPrice         decimal.Decimal `json:"slOrdPx"`
	AttachAlgoOrders           []AttachAlgoOrd `json:"attachAlgoOrds"`
	STPID                      string          `json:"stpId"`
	STPMode                    string          `json:"stpMode"`
	FeeCurrency                string          `json:"feeCcy"`
	Fee                        decimal.Decimal `json:"fee"`
	RebateCurrency             string          `json:"rebateCcy"`
	Rebate                     decimal.Decimal `json:"rebate"`
	Pnl                        decimal.Decimal `json:"pnl"`
	Source                     string          `json:"source"`
	Category                   string          `json:"category"`
	ReduceOnly                 string          `json:"reduceOnly"` // OKX sends a quoted "true"/"false"
	CancelSource               string          `json:"cancelSource"`
	CancelSourceReason         string          `json:"cancelSourceReason"`
	QuickMarginType            string          `json:"quickMgnType"`
	AlgoClientOrderID          string          `json:"algoClOrdId"`
	AlgoID                     string          `json:"algoId"`
	IsTakeProfitLimit          string          `json:"isTpLimit"`
	LastPrice                  decimal.Decimal `json:"lastPx"`
	UpdateTime                 time.Time       `json:"uTime"`
	CreationTime               time.Time       `json:"cTime"`
	RequestID                  string          `json:"reqId"`
	AmendResult                string          `json:"amendResult"`
	// AmendSource gains value "6" with the ELP->RPI rebranding rollout: order
	// price adjusted (rounded) by the system to satisfy the RPI maker spacing
	// rule.
	AmendSource        string `json:"amendSource"`
	Code               string `json:"code"`
	Message            string `json:"msg"`
	TradeQuoteCurrency string `json:"tradeQuoteCcy"`
}

// SubscribeLiquidationWarningService -- "liquidation-warning" channel (private; login).
//
// Pushes when a position approaches liquidation. instType is required.
// Activity-driven (no snapshot unless a position is at risk).
type SubscribeLiquidationWarningService struct {
	c          *WebSocketClient
	instType   InstType
	instFamily string
	instId     string
}

func (c *WebSocketClient) NewSubscribeLiquidationWarningService(instType InstType) *SubscribeLiquidationWarningService {
	return &SubscribeLiquidationWarningService{c: c, instType: instType}
}

// SetInstFamily narrows the push to one instrument family.
func (s *SubscribeLiquidationWarningService) SetInstFamily(instFamily string) *SubscribeLiquidationWarningService {
	s.instFamily = instFamily
	return s
}

// SetInstId narrows the push to one instrument.
func (s *SubscribeLiquidationWarningService) SetInstId(instId string) *SubscribeLiquidationWarningService {
	s.instId = instId
	return s
}

func (s *SubscribeLiquidationWarningService) Do(ctx context.Context, cb WsHandler[WsLiquidationWarning]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsLiquidationWarning](ctx, s.c, request.GatewayPrivate, true,
		request.WsArg{Channel: "liquidation-warning", InstrumentType: string(s.instType), InstrumentFamily: s.instFamily, InstrumentID: s.instId}, cb)
}

// WsLiquidationWarning is one at-risk position from the "liquidation-warning"
// channel. It carries the same per-position fields as the positions channel; the
// validating account had no at-risk positions, so the field set is modeled from
// the OKX channel field table.
type WsLiquidationWarning = WsPosition

// SubscribeAccountGreeksService -- "account-greeks" channel (private; login).
//
// Pushes the account's per-currency option greeks on change. Activity-driven
// (empty when the account holds no options).
type SubscribeAccountGreeksService struct {
	c   *WebSocketClient
	ccy string
}

func (c *WebSocketClient) NewSubscribeAccountGreeksService() *SubscribeAccountGreeksService {
	return &SubscribeAccountGreeksService{c: c}
}

// SetCcy limits the push to a single currency.
func (s *SubscribeAccountGreeksService) SetCcy(ccy string) *SubscribeAccountGreeksService {
	s.ccy = ccy
	return s
}

func (s *SubscribeAccountGreeksService) Do(ctx context.Context, cb WsHandler[WsAccountGreeks]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsAccountGreeks](ctx, s.c, request.GatewayPrivate, true,
		request.WsArg{Channel: "account-greeks", Currency: s.ccy}, cb)
}

// WsAccountGreeks is one currency's aggregated option greeks from the
// "account-greeks" channel. The validating account held no options, so the field
// set is modeled from the OKX channel field table.
type WsAccountGreeks struct {
	Currency  string          `json:"ccy"`
	DeltaBS   decimal.Decimal `json:"deltaBS"`
	DeltaPA   decimal.Decimal `json:"deltaPA"`
	GammaBS   decimal.Decimal `json:"gammaBS"`
	GammaPA   decimal.Decimal `json:"gammaPA"`
	ThetaBS   decimal.Decimal `json:"thetaBS"`
	ThetaPA   decimal.Decimal `json:"thetaPA"`
	VegaBS    decimal.Decimal `json:"vegaBS"`
	VegaPA    decimal.Decimal `json:"vegaPA"`
	Timestamp time.Time       `json:"ts"`
}
