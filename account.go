package okx

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-okx/request"
	"github.com/shopspring/decimal"
)

// GetBalanceService -- GET /api/v5/account/balance (Read)
//
// Returns the trading account's balance, total equity and per-currency details.
type GetBalanceService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetBalanceService() *GetBalanceService {
	return &GetBalanceService{c: c, params: map[string]string{}}
}

// SetCcy filters the details to a single currency (comma-separated for several).
func (s *GetBalanceService) SetCcy(ccy string) *GetBalanceService {
	s.params["ccy"] = ccy
	return s
}

func (s *GetBalanceService) Do(ctx context.Context) (*Balance, error) {
	req := request.Get(ctx, s.c, "/api/v5/account/balance", s.params).WithSign()
	return request.DoOne[Balance](req)
}

// Balance is the account-level balance summary plus per-currency details.
type Balance struct {
	UpdateTime            time.Time       `json:"uTime"`
	TotalEquity           decimal.Decimal `json:"totalEq"`
	IsolatedEquity        decimal.Decimal `json:"isoEq"`
	AdjustedEquity        decimal.Decimal `json:"adjEq"`
	AvailableEquity       decimal.Decimal `json:"availEq"`
	OrderFrozen           decimal.Decimal `json:"ordFroz"`
	IMR                   decimal.Decimal `json:"imr"`
	MMR                   decimal.Decimal `json:"mmr"`
	BorrowFrozen          decimal.Decimal `json:"borrowFroz"`
	MarginRatio           decimal.Decimal `json:"mgnRatio"`
	NotionalUSD           decimal.Decimal `json:"notionalUsd"`
	NotionalUSDForBorrow  decimal.Decimal `json:"notionalUsdForBorrow"`
	NotionalUSDForSwap    decimal.Decimal `json:"notionalUsdForSwap"`
	NotionalUSDForFutures decimal.Decimal `json:"notionalUsdForFutures"`
	NotionalUSDForOption  decimal.Decimal `json:"notionalUsdForOption"`
	UPL                   decimal.Decimal `json:"upl"`
	Delta                 decimal.Decimal `json:"delta"`
	DeltaLeverage         decimal.Decimal `json:"deltaLever"`
	DeltaNeutralStatus    string          `json:"deltaNeutralStatus"`
	Details               []BalanceDetail `json:"details"`
}

// BalanceDetail is a single currency's balance within the trading account.
type BalanceDetail struct {
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
	BorrowFrozen                   decimal.Decimal `json:"borrowFroz"`
	NotionalLeverage               decimal.Decimal `json:"notionalLever"`
	StrategyEquity                 decimal.Decimal `json:"stgyEq"`
	IsolatedUPL                    decimal.Decimal `json:"isoUpl"`
	SpotInUseAmount                decimal.Decimal `json:"spotInUseAmt"`
	ClientSpotInUseAmount          decimal.Decimal `json:"clSpotInUseAmt"`
	MaxSpotInUse                   decimal.Decimal `json:"maxSpotInUse"`
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

// GetPositionsService -- GET /api/v5/account/positions (Read)
//
// Returns the account's currently open positions.
type GetPositionsService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetPositionsService() *GetPositionsService {
	return &GetPositionsService{c: c, params: map[string]string{}}
}

// SetInstType filters by product line (MARGIN/SWAP/FUTURES/OPTION).
func (s *GetPositionsService) SetInstType(instType InstType) *GetPositionsService {
	s.params["instType"] = string(instType)
	return s
}

// SetInstId filters by instrument id (comma-separated for several).
func (s *GetPositionsService) SetInstId(instId string) *GetPositionsService {
	s.params["instId"] = instId
	return s
}

// SetPosId filters by position id (comma-separated for several).
func (s *GetPositionsService) SetPosId(posId string) *GetPositionsService {
	s.params["posId"] = posId
	return s
}

func (s *GetPositionsService) Do(ctx context.Context) ([]Position, error) {
	req := request.Get(ctx, s.c, "/api/v5/account/positions", s.params).WithSign()
	return request.DoList[Position](req)
}

// Position is a single open position. The account used to validate this SDK had
// no open positions, so the field set is modeled from the OKX doc field table.
type Position struct {
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
	BusinessReferenceID             string                   `json:"bizRefId"`
	BusinessReferenceType           string                   `json:"bizRefType"`
}

// PositionCloseOrderAlgo is a take-profit / stop-loss algo order attached to a
// position.
type PositionCloseOrderAlgo struct {
	AlgoID                     string          `json:"algoId"`
	StopLossTriggerPrice       decimal.Decimal `json:"slTriggerPx"`
	StopLossTriggerPriceType   string          `json:"slTriggerPxType"`
	TakeProfitTriggerPrice     decimal.Decimal `json:"tpTriggerPx"`
	TakeProfitTriggerPriceType string          `json:"tpTriggerPxType"`
	CloseFraction              decimal.Decimal `json:"closeFraction"`
}

// GetPositionsHistoryService -- GET /api/v5/account/positions-history (Read)
//
// Returns the account's closed-position history over the last three months.
type GetPositionsHistoryService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetPositionsHistoryService() *GetPositionsHistoryService {
	return &GetPositionsHistoryService{c: c, params: map[string]string{}}
}

// SetInstType filters by product line (MARGIN/SWAP/FUTURES/OPTION).
func (s *GetPositionsHistoryService) SetInstType(instType InstType) *GetPositionsHistoryService {
	s.params["instType"] = string(instType)
	return s
}

// SetInstId filters by instrument id.
func (s *GetPositionsHistoryService) SetInstId(instId string) *GetPositionsHistoryService {
	s.params["instId"] = instId
	return s
}

// SetMgnMode filters by margin mode (cross/isolated).
func (s *GetPositionsHistoryService) SetMgnMode(mgnMode MgnMode) *GetPositionsHistoryService {
	s.params["mgnMode"] = string(mgnMode)
	return s
}

// SetType filters by close type (1 partial close ... 5 partial liquidation, ...).
func (s *GetPositionsHistoryService) SetType(typ string) *GetPositionsHistoryService {
	s.params["type"] = typ
	return s
}

// SetPosId filters by position id.
func (s *GetPositionsHistoryService) SetPosId(posId string) *GetPositionsHistoryService {
	s.params["posId"] = posId
	return s
}

// SetAfter paginates to records earlier than the given time (older).
func (s *GetPositionsHistoryService) SetAfter(t time.Time) *GetPositionsHistoryService {
	s.params["after"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetBefore paginates to records later than the given time (newer).
func (s *GetPositionsHistoryService) SetBefore(t time.Time) *GetPositionsHistoryService {
	s.params["before"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetPositionsHistoryService) SetLimit(limit int) *GetPositionsHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetPositionsHistoryService) Do(ctx context.Context) ([]PositionHistory, error) {
	req := request.Get(ctx, s.c, "/api/v5/account/positions-history", s.params).WithSign()
	return request.DoList[PositionHistory](req)
}

// PositionHistory is one closed-position record.
type PositionHistory struct {
	InstrumentType        InstType        `json:"instType"`
	InstrumentID          string          `json:"instId"`
	MarginMode            MgnMode         `json:"mgnMode"`
	Type                  string          `json:"type"`
	Currency              string          `json:"ccy"`
	PositionID            string          `json:"posId"`
	PositionSide          PosSide         `json:"posSide"`
	Direction             string          `json:"direction"`
	Leverage              decimal.Decimal `json:"lever"`
	OpenAveragePrice      decimal.Decimal `json:"openAvgPx"`
	CloseAveragePrice     decimal.Decimal `json:"closeAvgPx"`
	NonSettleAveragePrice decimal.Decimal `json:"nonSettleAvgPx"`
	OpenMaxPosition       decimal.Decimal `json:"openMaxPos"`
	CloseTotalPosition    decimal.Decimal `json:"closeTotalPos"`
	Pnl                   decimal.Decimal `json:"pnl"`
	PnlRatio              decimal.Decimal `json:"pnlRatio"`
	RealizedPnl           decimal.Decimal `json:"realizedPnl"`
	SettledPnl            decimal.Decimal `json:"settledPnl"`
	Fee                   decimal.Decimal `json:"fee"`
	FundingFee            decimal.Decimal `json:"fundingFee"`
	LiquidationPenalty    decimal.Decimal `json:"liqPenalty"`
	TriggerPrice          decimal.Decimal `json:"triggerPx"`
	Underlying            string          `json:"uly"`
	CreationTime          time.Time       `json:"cTime"`
	UpdateTime            time.Time       `json:"uTime"`
}

// GetAccountPositionRiskService -- GET /api/v5/account/account-position-risk (Read)
//
// Returns the account- and position-level risk snapshot (equity per currency and
// per-position notional/quantity) at a single point in time.
type GetAccountPositionRiskService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetAccountPositionRiskService() *GetAccountPositionRiskService {
	return &GetAccountPositionRiskService{c: c, params: map[string]string{}}
}

// SetInstType filters by product line (MARGIN/SWAP/FUTURES/OPTION).
func (s *GetAccountPositionRiskService) SetInstType(instType InstType) *GetAccountPositionRiskService {
	s.params["instType"] = string(instType)
	return s
}

func (s *GetAccountPositionRiskService) Do(ctx context.Context) (*AccountPositionRisk, error) {
	req := request.Get(ctx, s.c, "/api/v5/account/account-position-risk", s.params).WithSign()
	return request.DoOne[AccountPositionRisk](req)
}

// AccountPositionRisk is the account/position risk snapshot.
type AccountPositionRisk struct {
	Timestamp      time.Time                    `json:"ts"`
	AdjustedEquity decimal.Decimal              `json:"adjEq"`
	BalanceData    []AccountPositionRiskBalance `json:"balData"`
	PositionData   []AccountPositionRiskPos     `json:"posData"`
}

// AccountPositionRiskBalance is one currency's equity within the risk snapshot.
type AccountPositionRiskBalance struct {
	Currency       string          `json:"ccy"`
	Equity         decimal.Decimal `json:"eq"`
	DiscountEquity decimal.Decimal `json:"disEq"`
}

// AccountPositionRiskPos is one position's risk data within the snapshot. The
// account used to validate this SDK had no open positions, so this field set is
// modeled from the OKX doc field table.
type AccountPositionRiskPos struct {
	InstrumentType   InstType        `json:"instType"`
	InstrumentID     string          `json:"instId"`
	MarginMode       MgnMode         `json:"mgnMode"`
	PositionID       string          `json:"posId"`
	PositionSide     PosSide         `json:"posSide"`
	Position         decimal.Decimal `json:"pos"`
	BaseBalance      decimal.Decimal `json:"baseBal"`
	QuoteBalance     decimal.Decimal `json:"quoteBal"`
	PositionCurrency string          `json:"posCcy"`
	Currency         string          `json:"ccy"`
	NotionalCurrency decimal.Decimal `json:"notionalCcy"`
	NotionalUSD      decimal.Decimal `json:"notionalUsd"`
}

// GetAccountConfigService -- GET /api/v5/account/config (Read)
//
// Returns the account's configuration: account level, position mode, fee/greeks
// settings and related flags.
type GetAccountConfigService struct {
	c *Client
}

func (c *Client) NewGetAccountConfigService() *GetAccountConfigService {
	return &GetAccountConfigService{c: c}
}

func (s *GetAccountConfigService) Do(ctx context.Context) (*AccountConfig, error) {
	req := request.Get(ctx, s.c, "/api/v5/account/config").WithSign()
	return request.DoOne[AccountConfig](req)
}

// AccountConfig is the account's configuration.
type AccountConfig struct {
	UID                   string   `json:"uid"`
	MainUID               string   `json:"mainUid"`
	AccountLevel          string   `json:"acctLv"`
	AccountSTPMode        string   `json:"acctStpMode"`
	PositionMode          string   `json:"posMode"`
	AutoLoan              bool     `json:"autoLoan"`
	GreeksType            string   `json:"greeksType"`
	Level                 string   `json:"level"`
	LevelTemporary        string   `json:"levelTmp"`
	ContractIsolatedMode  string   `json:"ctIsoMode"`
	MarginIsolatedMode    string   `json:"mgnIsoMode"`
	SpotOffsetType        string   `json:"spotOffsetType"`
	RoleType              string   `json:"roleType"`
	TraderInstruments     []string `json:"traderInsts"`
	SpotTraderInstruments []string `json:"spotTraderInsts"`
	OpAuth                string   `json:"opAuth"`
	KYCLevel              string   `json:"kycLv"`
	Label                 string   `json:"label"`
	IP                    string   `json:"ip"`
	Perm                  string   `json:"perm"`
	LiquidationGear       string   `json:"liquidationGear"`
	EnableSpotBorrow      bool     `json:"enableSpotBorrow"`
	SpotBorrowAutoRepay   bool     `json:"spotBorrowAutoRepay"`
	Type                  string   `json:"type"`
	SpotRoleType          string   `json:"spotRoleType"`
	StrategyType          string   `json:"stgyType"`
	SettleCurrency        string   `json:"settleCcy"`
	SettleCurrencyList    []string `json:"settleCcyList"`
	FeeType               string   `json:"feeType"`
}

// GetAccountInstrumentsService -- GET /api/v5/account/instruments (Read)
//
// Returns the instruments tradable by the current account for a product line.
// The shape matches the public instruments endpoint, so it reuses Instrument.
type GetAccountInstrumentsService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetAccountInstrumentsService(instType InstType) *GetAccountInstrumentsService {
	return &GetAccountInstrumentsService{c: c, params: map[string]string{"instType": string(instType)}}
}

// SetUly filters by underlying (FUTURES/SWAP/OPTION).
func (s *GetAccountInstrumentsService) SetUly(uly string) *GetAccountInstrumentsService {
	s.params["uly"] = uly
	return s
}

// SetInstFamily filters by instrument family (FUTURES/SWAP/OPTION).
func (s *GetAccountInstrumentsService) SetInstFamily(instFamily string) *GetAccountInstrumentsService {
	s.params["instFamily"] = instFamily
	return s
}

// SetInstId filters by a single instrument id.
func (s *GetAccountInstrumentsService) SetInstId(instId string) *GetAccountInstrumentsService {
	s.params["instId"] = instId
	return s
}

func (s *GetAccountInstrumentsService) Do(ctx context.Context) ([]Instrument, error) {
	req := request.Get(ctx, s.c, "/api/v5/account/instruments", s.params).WithSign()
	return request.DoList[Instrument](req)
}

// GetMaxSizeService -- GET /api/v5/account/max-size (Read)
//
// Returns the maximum buyable/sellable (or openable) size of an instrument under
// the given trade mode.
type GetMaxSizeService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetMaxSizeService(instId string, tdMode TdMode) *GetMaxSizeService {
	return &GetMaxSizeService{c: c, params: map[string]string{
		"instId": instId,
		"tdMode": string(tdMode),
	}}
}

// SetCcy sets the margin currency (MARGIN cross only).
func (s *GetMaxSizeService) SetCcy(ccy string) *GetMaxSizeService {
	s.params["ccy"] = ccy
	return s
}

// SetPx sets the order price used to estimate the max size.
func (s *GetMaxSizeService) SetPx(px decimal.Decimal) *GetMaxSizeService {
	s.params["px"] = px.String()
	return s
}

// SetLeverage sets the leverage used to estimate the max size.
func (s *GetMaxSizeService) SetLeverage(leverage decimal.Decimal) *GetMaxSizeService {
	s.params["leverage"] = leverage.String()
	return s
}

// SetUnSpotOffset toggles whether spot-derivatives offset is excluded.
func (s *GetMaxSizeService) SetUnSpotOffset(unSpotOffset bool) *GetMaxSizeService {
	s.params["unSpotOffset"] = strconv.FormatBool(unSpotOffset)
	return s
}

func (s *GetMaxSizeService) Do(ctx context.Context) ([]MaxSize, error) {
	req := request.Get(ctx, s.c, "/api/v5/account/max-size", s.params).WithSign()
	return request.DoList[MaxSize](req)
}

// MaxSize is the maximum buy/sell size of an instrument.
type MaxSize struct {
	InstrumentID       string          `json:"instId"`
	Currency           string          `json:"ccy"`
	MaxBuy             decimal.Decimal `json:"maxBuy"`
	MaxSell            decimal.Decimal `json:"maxSell"`
	TradeQuoteCurrency string          `json:"tradeQuoteCcy"`
}

// GetMaxAvailSizeService -- GET /api/v5/account/max-avail-size (Read)
//
// Returns the maximum available tradable amount of an instrument under the given
// trade mode (accounting for current balance and positions).
type GetMaxAvailSizeService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetMaxAvailSizeService(instId string, tdMode TdMode) *GetMaxAvailSizeService {
	return &GetMaxAvailSizeService{c: c, params: map[string]string{
		"instId": instId,
		"tdMode": string(tdMode),
	}}
}

// SetCcy sets the margin currency (MARGIN cross only).
func (s *GetMaxAvailSizeService) SetCcy(ccy string) *GetMaxAvailSizeService {
	s.params["ccy"] = ccy
	return s
}

// SetReduceOnly restricts the estimate to reduce-only sizing.
func (s *GetMaxAvailSizeService) SetReduceOnly(reduceOnly bool) *GetMaxAvailSizeService {
	s.params["reduceOnly"] = strconv.FormatBool(reduceOnly)
	return s
}

// SetPx sets the order price used to estimate the available size.
func (s *GetMaxAvailSizeService) SetPx(px decimal.Decimal) *GetMaxAvailSizeService {
	s.params["px"] = px.String()
	return s
}

// SetUnSpotOffset toggles whether spot-derivatives offset is excluded.
func (s *GetMaxAvailSizeService) SetUnSpotOffset(unSpotOffset bool) *GetMaxAvailSizeService {
	s.params["unSpotOffset"] = strconv.FormatBool(unSpotOffset)
	return s
}

// SetQuickMgnType sets the quick-margin borrow type (manual/auto_borrow/auto_borrow_repay).
func (s *GetMaxAvailSizeService) SetQuickMgnType(quickMgnType string) *GetMaxAvailSizeService {
	s.params["quickMgnType"] = quickMgnType
	return s
}

func (s *GetMaxAvailSizeService) Do(ctx context.Context) ([]MaxAvailSize, error) {
	req := request.Get(ctx, s.c, "/api/v5/account/max-avail-size", s.params).WithSign()
	return request.DoList[MaxAvailSize](req)
}

// MaxAvailSize is the maximum available buy/sell amount of an instrument.
type MaxAvailSize struct {
	InstrumentID       string          `json:"instId"`
	AvailableBuy       decimal.Decimal `json:"availBuy"`
	AvailableSell      decimal.Decimal `json:"availSell"`
	TradeQuoteCurrency string          `json:"tradeQuoteCcy"`
}

// GetTradeFeeService -- GET /api/v5/account/trade-fee (Read)
//
// Returns the account's maker/taker trading fee rates for a product line.
type GetTradeFeeService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetTradeFeeService(instType InstType) *GetTradeFeeService {
	return &GetTradeFeeService{c: c, params: map[string]string{"instType": string(instType)}}
}

// SetInstId filters by a single instrument id (SPOT/MARGIN).
func (s *GetTradeFeeService) SetInstId(instId string) *GetTradeFeeService {
	s.params["instId"] = instId
	return s
}

// SetUly filters by underlying (FUTURES/SWAP/OPTION).
func (s *GetTradeFeeService) SetUly(uly string) *GetTradeFeeService {
	s.params["uly"] = uly
	return s
}

// SetInstFamily filters by instrument family (FUTURES/SWAP/OPTION).
func (s *GetTradeFeeService) SetInstFamily(instFamily string) *GetTradeFeeService {
	s.params["instFamily"] = instFamily
	return s
}

// SetRuleType filters by rule type (normal/pre_market).
func (s *GetTradeFeeService) SetRuleType(ruleType string) *GetTradeFeeService {
	s.params["ruleType"] = ruleType
	return s
}

func (s *GetTradeFeeService) Do(ctx context.Context) (*TradeFee, error) {
	req := request.Get(ctx, s.c, "/api/v5/account/trade-fee", s.params).WithSign()
	return request.DoOne[TradeFee](req)
}

// TradeFee is the account's fee schedule for a product line.
type TradeFee struct {
	InstrumentType InstType        `json:"instType"`
	Level          string          `json:"level"`
	Taker          decimal.Decimal `json:"taker"`
	Maker          decimal.Decimal `json:"maker"`
	TakerUSDT      decimal.Decimal `json:"takerU"`
	MakerUSDT      decimal.Decimal `json:"makerU"`
	TakerUSDC      decimal.Decimal `json:"takerUSDC"`
	MakerUSDC      decimal.Decimal `json:"makerUSDC"`
	Delivery       decimal.Decimal `json:"delivery"`
	Exercise       decimal.Decimal `json:"exercise"`
	Settle         decimal.Decimal `json:"settle"`
	Category       string          `json:"category"`
	RuleType       string          `json:"ruleType"`
	FeeGroup       []TradeFeeGroup `json:"feeGroup"`
	Fiat           []TradeFeeFiat  `json:"fiat"`
	Timestamp      time.Time       `json:"ts"`
}

// TradeFeeGroup is one fee group's maker/taker rates.
type TradeFeeGroup struct {
	GroupID  string          `json:"groupId"`
	Maker    decimal.Decimal `json:"maker"`
	Taker    decimal.Decimal `json:"taker"`
	ElpMaker decimal.Decimal `json:"elpMaker"`
}

// TradeFeeFiat is one fiat currency's maker/taker rates.
type TradeFeeFiat struct {
	Currency string          `json:"ccy"`
	Maker    decimal.Decimal `json:"maker"`
	Taker    decimal.Decimal `json:"taker"`
}

// GetMaxWithdrawalService -- GET /api/v5/account/max-withdrawal (Read)
//
// Returns the maximum amount that can be transferred out of the trading account
// per currency.
type GetMaxWithdrawalService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetMaxWithdrawalService() *GetMaxWithdrawalService {
	return &GetMaxWithdrawalService{c: c, params: map[string]string{}}
}

// SetCcy filters by currency (comma-separated for several).
func (s *GetMaxWithdrawalService) SetCcy(ccy string) *GetMaxWithdrawalService {
	s.params["ccy"] = ccy
	return s
}

func (s *GetMaxWithdrawalService) Do(ctx context.Context) ([]MaxWithdrawal, error) {
	req := request.Get(ctx, s.c, "/api/v5/account/max-withdrawal", s.params).WithSign()
	return request.DoList[MaxWithdrawal](req)
}

// MaxWithdrawal is a currency's maximum transferable-out amount.
type MaxWithdrawal struct {
	Currency                string          `json:"ccy"`
	MaxWithdrawal           decimal.Decimal `json:"maxWd"`
	MaxWdEx                 decimal.Decimal `json:"maxWdEx"`
	SpotOffsetMaxWithdrawal decimal.Decimal `json:"spotOffsetMaxWd"`
	SpotOffsetMaxWdEx       decimal.Decimal `json:"spotOffsetMaxWdEx"`
}

// GetGreeksService -- GET /api/v5/account/greeks (Read)
//
// Returns the account's per-currency option greeks.
type GetGreeksService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetGreeksService() *GetGreeksService {
	return &GetGreeksService{c: c, params: map[string]string{}}
}

// SetCcy filters by currency.
func (s *GetGreeksService) SetCcy(ccy string) *GetGreeksService {
	s.params["ccy"] = ccy
	return s
}

func (s *GetGreeksService) Do(ctx context.Context) ([]Greeks, error) {
	req := request.Get(ctx, s.c, "/api/v5/account/greeks", s.params).WithSign()
	return request.DoList[Greeks](req)
}

// Greeks is a currency's aggregated option greeks. The account used to validate
// this SDK had no option positions, so the field set is modeled from the OKX doc
// field table.
type Greeks struct {
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

// GetAccountPositionTiersService -- GET /api/v5/account/position-tiers (Read)
//
// Returns the account's maximum position size per underlying / instrument family
// (the per-account view of the public position-tier schedule).
type GetAccountPositionTiersService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetAccountPositionTiersService(instType InstType) *GetAccountPositionTiersService {
	return &GetAccountPositionTiersService{c: c, params: map[string]string{"instType": string(instType)}}
}

// SetUly filters by underlying (comma-separated for several).
func (s *GetAccountPositionTiersService) SetUly(uly string) *GetAccountPositionTiersService {
	s.params["uly"] = uly
	return s
}

// SetInstFamily filters by instrument family (comma-separated for several).
func (s *GetAccountPositionTiersService) SetInstFamily(instFamily string) *GetAccountPositionTiersService {
	s.params["instFamily"] = instFamily
	return s
}

func (s *GetAccountPositionTiersService) Do(ctx context.Context) ([]AccountPositionTier, error) {
	req := request.Get(ctx, s.c, "/api/v5/account/position-tiers", s.params).WithSign()
	return request.DoList[AccountPositionTier](req)
}

// AccountPositionTier is the account's maximum position size for an underlying or
// instrument family.
type AccountPositionTier struct {
	Underlying       string          `json:"uly"`
	InstrumentFamily string          `json:"instFamily"`
	MaxSize          decimal.Decimal `json:"maxSz"`
	PositionType     string          `json:"posType"`
}

// GetLeverageInfoService -- GET /api/v5/account/leverage-info (Read)
//
// Returns the account's current leverage setting for an instrument under the
// given margin mode.
type GetLeverageInfoService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetLeverageInfoService(instId string, mgnMode MgnMode) *GetLeverageInfoService {
	return &GetLeverageInfoService{c: c, params: map[string]string{
		"instId":  instId,
		"mgnMode": string(mgnMode),
	}}
}

func (s *GetLeverageInfoService) Do(ctx context.Context) ([]LeverageInfo, error) {
	req := request.Get(ctx, s.c, "/api/v5/account/leverage-info", s.params).WithSign()
	return request.DoList[LeverageInfo](req)
}

// LeverageInfo is the account's leverage setting for an instrument/side.
type LeverageInfo struct {
	InstrumentID string          `json:"instId"`
	Currency     string          `json:"ccy"`
	MarginMode   MgnMode         `json:"mgnMode"`
	PositionSide PosSide         `json:"posSide"`
	Leverage     decimal.Decimal `json:"lever"`
}

// GetAdjustLeverageInfoService -- GET /api/v5/account/adjust-leverage-info (Read)
//
// Returns the estimated information (available transfer, liq price, max amount)
// for adjusting leverage to a target value.
type GetAdjustLeverageInfoService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetAdjustLeverageInfoService(instType InstType, mgnMode MgnMode, lever decimal.Decimal) *GetAdjustLeverageInfoService {
	return &GetAdjustLeverageInfoService{c: c, params: map[string]string{
		"instType": string(instType),
		"mgnMode":  string(mgnMode),
		"lever":    lever.String(),
	}}
}

// SetInstId sets the instrument id (required unless ccy is set).
func (s *GetAdjustLeverageInfoService) SetInstId(instId string) *GetAdjustLeverageInfoService {
	s.params["instId"] = instId
	return s
}

// SetCcy sets the margin currency (required for cross MARGIN unless instId is set).
func (s *GetAdjustLeverageInfoService) SetCcy(ccy string) *GetAdjustLeverageInfoService {
	s.params["ccy"] = ccy
	return s
}

// SetPosSide sets the position side (long/short/net).
func (s *GetAdjustLeverageInfoService) SetPosSide(posSide PosSide) *GetAdjustLeverageInfoService {
	s.params["posSide"] = string(posSide)
	return s
}

func (s *GetAdjustLeverageInfoService) Do(ctx context.Context) (*AdjustLeverageInfo, error) {
	req := request.Get(ctx, s.c, "/api/v5/account/adjust-leverage-info", s.params).WithSign()
	return request.DoOne[AdjustLeverageInfo](req)
}

// AdjustLeverageInfo is the estimated impact of adjusting leverage.
type AdjustLeverageInfo struct {
	EstimatedAvailableQuoteTransfer decimal.Decimal `json:"estAvailQuoteTrans"`
	EstimatedAvailableTransfer      decimal.Decimal `json:"estAvailTrans"`
	EstimatedLiquidationPrice       decimal.Decimal `json:"estLiqPx"`
	EstimatedMargin                 decimal.Decimal `json:"estMgn"`
	EstimatedQuoteMargin            decimal.Decimal `json:"estQuoteMgn"`
	EstimatedMaxAmount              decimal.Decimal `json:"estMaxAmt"`
	EstimatedQuoteMaxAmount         decimal.Decimal `json:"estQuoteMaxAmt"`
	ExistOrder                      bool            `json:"existOrd"`
	MaxLeverage                     decimal.Decimal `json:"maxLever"`
	MinLeverage                     decimal.Decimal `json:"minLever"`
}

// GetMaxLoanService -- GET /api/v5/account/max-loan (Read)
//
// Returns the account's maximum loanable amount for an instrument/currency under
// the given margin mode.
type GetMaxLoanService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetMaxLoanService(mgnMode MgnMode) *GetMaxLoanService {
	return &GetMaxLoanService{c: c, params: map[string]string{"mgnMode": string(mgnMode)}}
}

// SetInstId sets the instrument id (single or comma-separated).
func (s *GetMaxLoanService) SetInstId(instId string) *GetMaxLoanService {
	s.params["instId"] = instId
	return s
}

// SetCcy sets the borrow currency (cross MARGIN).
func (s *GetMaxLoanService) SetCcy(ccy string) *GetMaxLoanService {
	s.params["ccy"] = ccy
	return s
}

// SetMgnCcy sets the margin currency (cross MARGIN).
func (s *GetMaxLoanService) SetMgnCcy(mgnCcy string) *GetMaxLoanService {
	s.params["mgnCcy"] = mgnCcy
	return s
}

func (s *GetMaxLoanService) Do(ctx context.Context) ([]MaxLoan, error) {
	req := request.Get(ctx, s.c, "/api/v5/account/max-loan", s.params).WithSign()
	return request.DoList[MaxLoan](req)
}

// MaxLoan is the account's maximum loanable amount for an instrument side.
type MaxLoan struct {
	InstrumentID   string          `json:"instId"`
	MarginMode     MgnMode         `json:"mgnMode"`
	MarginCurrency string          `json:"mgnCcy"`
	MaxLoan        decimal.Decimal `json:"maxLoan"`
	Currency       string          `json:"ccy"`
	Side           Side            `json:"side"`
}

// GetRiskStateService -- GET /api/v5/account/risk-state (Read)
//
// Returns the portfolio-margin account's auto-borrow/auto-repay risk state. Only
// available when the account is in portfolio-margin mode.
type GetRiskStateService struct {
	c *Client
}

func (c *Client) NewGetRiskStateService() *GetRiskStateService {
	return &GetRiskStateService{c: c}
}

func (s *GetRiskStateService) Do(ctx context.Context) (*RiskState, error) {
	req := request.Get(ctx, s.c, "/api/v5/account/risk-state").WithSign()
	return request.DoOne[RiskState](req)
}

// RiskState is the portfolio-margin account's risk-offset state. The validating
// account is not in portfolio-margin mode (the endpoint returns 51010), so the
// field set is modeled from the OKX doc field table.
type RiskState struct {
	AtRisk       bool      `json:"atRisk"`
	AtRiskIndex  []string  `json:"atRiskIdx"`
	AtRiskMargin []string  `json:"atRiskMgn"`
	Timestamp    time.Time `json:"ts"`
}
