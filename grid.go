package okx

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-okx/request"
	"github.com/shopspring/decimal"
)

// Grid trading endpoints. NOTE: despite OKX documenting these under the
// "Order Book Trading > Grid Trading" section, the REAL live REST base path is
// /api/v5/tradingBot/grid/* (verified by live sign-curl: /api/v5/trade/grid/*
// returns 404, /api/v5/tradingBot/grid/* returns code "0"). The RSI back-testing
// helper lives at /api/v5/tradingBot/public/rsi-back-testing.

// GridAlgoOrdType is the kind of grid bot: spot/margin grid ("grid") or
// contract (perpetual/futures) grid ("contract_grid").
type GridAlgoOrdType string

const (
	GridAlgoOrdTypeGrid         GridAlgoOrdType = "grid"
	GridAlgoOrdTypeContractGrid GridAlgoOrdType = "contract_grid"
)

// GridRunType is the grid spacing mode: "1" arithmetic, "2" geometric.
type GridRunType string

const (
	GridRunTypeArithmetic GridRunType = "1"
	GridRunTypeGeometric  GridRunType = "2"
)

// GridDirection is the direction of a contract grid (long/short/neutral).
type GridDirection string

const (
	GridDirectionLong    GridDirection = "long"
	GridDirectionShort   GridDirection = "short"
	GridDirectionNeutral GridDirection = "neutral"
)

// GridSubOrderType selects which sub-orders to return: "live" (open) or
// "filled".
type GridSubOrderType string

const (
	GridSubOrderTypeLive   GridSubOrderType = "live"
	GridSubOrderTypeFilled GridSubOrderType = "filled"
)

// GridResult is the per-item ack of a grid algo place/amend/stop operation. The
// real reason for a failure is carried in sCode/sMsg even when the top-level
// envelope code is "1".
type GridResult struct {
	AlgoID            string `json:"algoId"`
	AlgoClientOrderID string `json:"algoClOrdId"`
	SCode             string `json:"sCode"`
	SMsg              string `json:"sMsg"`
	Tag               string `json:"tag"`
}

// PlaceGridAlgoOrderService -- POST /api/v5/tradingBot/grid/order-algo (Trade)
//
// Places a new spot/margin grid ("grid") or contract grid ("contract_grid")
// algo order. IMPLEMENT-ONLY: never executed during validation (real account).
type PlaceGridAlgoOrderService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewPlaceGridAlgoOrderService(instId string, algoOrdType GridAlgoOrdType, maxPx, minPx decimal.Decimal, gridNum int) *PlaceGridAlgoOrderService {
	return &PlaceGridAlgoOrderService{c: c, body: map[string]any{
		"instId":      instId,
		"algoOrdType": string(algoOrdType),
		"maxPx":       maxPx.String(),
		"minPx":       minPx.String(),
		"gridNum":     strconv.Itoa(gridNum),
	}}
}

// SetRunType sets the grid spacing mode ("1" arithmetic, "2" geometric).
func (s *PlaceGridAlgoOrderService) SetRunType(runType GridRunType) *PlaceGridAlgoOrderService {
	s.body["runType"] = string(runType)
	return s
}

// SetTpTriggerPx sets the take-profit trigger price (spot grid).
func (s *PlaceGridAlgoOrderService) SetTpTriggerPx(px decimal.Decimal) *PlaceGridAlgoOrderService {
	s.body["tpTriggerPx"] = px.String()
	return s
}

// SetSlTriggerPx sets the stop-loss trigger price (spot grid).
func (s *PlaceGridAlgoOrderService) SetSlTriggerPx(px decimal.Decimal) *PlaceGridAlgoOrderService {
	s.body["slTriggerPx"] = px.String()
	return s
}

// SetAlgoClOrdId sets a client-supplied algo order id.
func (s *PlaceGridAlgoOrderService) SetAlgoClOrdId(algoClOrdId string) *PlaceGridAlgoOrderService {
	s.body["algoClOrdId"] = algoClOrdId
	return s
}

// SetTag sets an order tag.
func (s *PlaceGridAlgoOrderService) SetTag(tag string) *PlaceGridAlgoOrderService {
	s.body["tag"] = tag
	return s
}

// SetProfitSharingRatio sets the profit-sharing ratio (copy-trading leaders).
func (s *PlaceGridAlgoOrderService) SetProfitSharingRatio(ratio decimal.Decimal) *PlaceGridAlgoOrderService {
	s.body["profitSharingRatio"] = ratio.String()
	return s
}

// --- spot/margin grid ("grid") only ---

// SetQuoteSz sets the invested amount in the quote currency (spot grid).
func (s *PlaceGridAlgoOrderService) SetQuoteSz(sz decimal.Decimal) *PlaceGridAlgoOrderService {
	s.body["quoteSz"] = sz.String()
	return s
}

// SetBaseSz sets the invested amount in the base currency (spot grid).
func (s *PlaceGridAlgoOrderService) SetBaseSz(sz decimal.Decimal) *PlaceGridAlgoOrderService {
	s.body["baseSz"] = sz.String()
	return s
}

// SetTdMode sets the trade mode (cross/isolated for margin grid; cash for spot).
func (s *PlaceGridAlgoOrderService) SetTdMode(tdMode TdMode) *PlaceGridAlgoOrderService {
	s.body["tdMode"] = string(tdMode)
	return s
}

// SetCcy sets the margin currency (margin grid).
func (s *PlaceGridAlgoOrderService) SetCcy(ccy string) *PlaceGridAlgoOrderService {
	s.body["ccy"] = ccy
	return s
}

// --- contract grid ("contract_grid") only ---

// SetSz sets the invested margin amount (contract grid).
func (s *PlaceGridAlgoOrderService) SetSz(sz decimal.Decimal) *PlaceGridAlgoOrderService {
	s.body["sz"] = sz.String()
	return s
}

// SetDirection sets the contract grid direction (long/short/neutral).
func (s *PlaceGridAlgoOrderService) SetDirection(direction GridDirection) *PlaceGridAlgoOrderService {
	s.body["direction"] = string(direction)
	return s
}

// SetLever sets the leverage (contract grid).
func (s *PlaceGridAlgoOrderService) SetLever(lever decimal.Decimal) *PlaceGridAlgoOrderService {
	s.body["lever"] = lever.String()
	return s
}

// SetBasePos toggles whether to open a base position immediately (contract grid).
func (s *PlaceGridAlgoOrderService) SetBasePos(basePos bool) *PlaceGridAlgoOrderService {
	s.body["basePos"] = basePos
	return s
}

// SetTpRatio sets the take-profit ratio (contract grid).
func (s *PlaceGridAlgoOrderService) SetTpRatio(ratio decimal.Decimal) *PlaceGridAlgoOrderService {
	s.body["tpRatio"] = ratio.String()
	return s
}

// SetSlRatio sets the stop-loss ratio (contract grid).
func (s *PlaceGridAlgoOrderService) SetSlRatio(ratio decimal.Decimal) *PlaceGridAlgoOrderService {
	s.body["slRatio"] = ratio.String()
	return s
}

func (s *PlaceGridAlgoOrderService) Do(ctx context.Context) (*GridResult, error) {
	req := request.Post(ctx, s.c, "/api/v5/tradingBot/grid/order-algo", s.body).WithSign()
	list, err := request.DoListPartial[GridResult](req)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return &list[0], nil
}

// AmendGridAlgoOrderService -- POST /api/v5/tradingBot/grid/amend-order-algo (Trade)
//
// Amends the take-profit / stop-loss settings of a running grid algo order.
// IMPLEMENT-ONLY.
type AmendGridAlgoOrderService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewAmendGridAlgoOrderService(algoId, instId string) *AmendGridAlgoOrderService {
	return &AmendGridAlgoOrderService{c: c, body: map[string]any{
		"algoId": algoId,
		"instId": instId,
	}}
}

// SetSlTriggerPx sets the stop-loss trigger price.
func (s *AmendGridAlgoOrderService) SetSlTriggerPx(px decimal.Decimal) *AmendGridAlgoOrderService {
	s.body["slTriggerPx"] = px.String()
	return s
}

// SetTpTriggerPx sets the take-profit trigger price.
func (s *AmendGridAlgoOrderService) SetTpTriggerPx(px decimal.Decimal) *AmendGridAlgoOrderService {
	s.body["tpTriggerPx"] = px.String()
	return s
}

// SetTpRatio sets the take-profit ratio (contract grid).
func (s *AmendGridAlgoOrderService) SetTpRatio(ratio decimal.Decimal) *AmendGridAlgoOrderService {
	s.body["tpRatio"] = ratio.String()
	return s
}

// SetSlRatio sets the stop-loss ratio (contract grid).
func (s *AmendGridAlgoOrderService) SetSlRatio(ratio decimal.Decimal) *AmendGridAlgoOrderService {
	s.body["slRatio"] = ratio.String()
	return s
}

// SetTriggerParams sets the trigger-parameter list (advanced trigger config).
func (s *AmendGridAlgoOrderService) SetTriggerParams(params []GridTriggerParam) *AmendGridAlgoOrderService {
	s.body["triggerParams"] = params
	return s
}

func (s *AmendGridAlgoOrderService) Do(ctx context.Context) (*GridResult, error) {
	req := request.Post(ctx, s.c, "/api/v5/tradingBot/grid/amend-order-algo", s.body).WithSign()
	list, err := request.DoListPartial[GridResult](req)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return &list[0], nil
}

// GridStopArg is one entry of a batch grid-stop request.
type GridStopArg struct {
	AlgoID        string          `json:"algoId"`
	InstrumentID  string          `json:"instId"`
	AlgoOrderType GridAlgoOrdType `json:"algoOrdType"`
	StopType      string          `json:"stopType"`
}

// StopGridAlgoOrderService -- POST /api/v5/tradingBot/grid/stop-order-algo (Trade)
//
// Stops one or more running grid algo orders. The request body is an ARRAY.
// IMPLEMENT-ONLY.
type StopGridAlgoOrderService struct {
	c     *Client
	items []GridStopArg
}

func (c *Client) NewStopGridAlgoOrderService(items []GridStopArg) *StopGridAlgoOrderService {
	return &StopGridAlgoOrderService{c: c, items: items}
}

func (s *StopGridAlgoOrderService) Do(ctx context.Context) ([]GridResult, error) {
	req := request.Post(ctx, s.c, "/api/v5/tradingBot/grid/stop-order-algo").SetBody(s.items).WithSign()
	return request.DoListPartial[GridResult](req)
}

// GridClosePositionResult is the ack of a contract-grid close-position request.
type GridClosePositionResult struct {
	AlgoID            string `json:"algoId"`
	OrderID           string `json:"ordId"`
	AlgoClientOrderID string `json:"algoClOrdId"`
	Tag               string `json:"tag"`
}

// CloseGridPositionService -- POST /api/v5/tradingBot/grid/close-position (Trade)
//
// Closes the position of a stopped contract grid (market or limit). IMPLEMENT-ONLY.
type CloseGridPositionService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewCloseGridPositionService(algoId string, mktClose bool) *CloseGridPositionService {
	return &CloseGridPositionService{c: c, body: map[string]any{
		"algoId":   algoId,
		"mktClose": mktClose,
	}}
}

// SetSz sets the quantity to close (required when mktClose is false).
func (s *CloseGridPositionService) SetSz(sz decimal.Decimal) *CloseGridPositionService {
	s.body["sz"] = sz.String()
	return s
}

// SetPx sets the close price (required when mktClose is false).
func (s *CloseGridPositionService) SetPx(px decimal.Decimal) *CloseGridPositionService {
	s.body["px"] = px.String()
	return s
}

func (s *CloseGridPositionService) Do(ctx context.Context) (*GridClosePositionResult, error) {
	req := request.Post(ctx, s.c, "/api/v5/tradingBot/grid/close-position", s.body).WithSign()
	return request.DoOne[GridClosePositionResult](req)
}

// GridCancelCloseOrderResult is the ack of a cancel-close-order request.
type GridCancelCloseOrderResult struct {
	AlgoID  string `json:"algoId"`
	OrderID string `json:"ordId"`
}

// CancelGridCloseOrderService -- POST /api/v5/tradingBot/grid/cancel-close-order (Trade)
//
// Cancels an outstanding close-position order of a stopped contract grid.
// IMPLEMENT-ONLY.
type CancelGridCloseOrderService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewCancelGridCloseOrderService(algoId, ordId string) *CancelGridCloseOrderService {
	return &CancelGridCloseOrderService{c: c, body: map[string]any{
		"algoId": algoId,
		"ordId":  ordId,
	}}
}

func (s *CancelGridCloseOrderService) Do(ctx context.Context) (*GridCancelCloseOrderResult, error) {
	req := request.Post(ctx, s.c, "/api/v5/tradingBot/grid/cancel-close-order", s.body).WithSign()
	return request.DoOne[GridCancelCloseOrderResult](req)
}

// GridInstantTriggerResult is the ack of an order-instant-trigger request.
type GridInstantTriggerResult struct {
	AlgoID string `json:"algoId"`
}

// InstantTriggerGridService -- POST /api/v5/tradingBot/grid/order-instant-trigger (Trade)
//
// Manually triggers a pending grid algo order immediately. IMPLEMENT-ONLY.
type InstantTriggerGridService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewInstantTriggerGridService(algoId string, algoOrdType GridAlgoOrdType) *InstantTriggerGridService {
	return &InstantTriggerGridService{c: c, body: map[string]any{
		"algoId":      algoId,
		"algoOrdType": string(algoOrdType),
	}}
}

func (s *InstantTriggerGridService) Do(ctx context.Context) (*GridInstantTriggerResult, error) {
	req := request.Post(ctx, s.c, "/api/v5/tradingBot/grid/order-instant-trigger", s.body).WithSign()
	return request.DoOne[GridInstantTriggerResult](req)
}

// GetGridAlgoPendingService -- GET /api/v5/tradingBot/grid/orders-algo-pending (Read)
//
// Returns the account's currently running (pending) grid algo orders.
type GetGridAlgoPendingService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetGridAlgoPendingService(algoOrdType GridAlgoOrdType) *GetGridAlgoPendingService {
	return &GetGridAlgoPendingService{c: c, params: map[string]string{"algoOrdType": string(algoOrdType)}}
}

// SetInstType filters by product line.
func (s *GetGridAlgoPendingService) SetInstType(instType InstType) *GetGridAlgoPendingService {
	s.params["instType"] = string(instType)
	return s
}

// SetInstId filters by instrument id.
func (s *GetGridAlgoPendingService) SetInstId(instId string) *GetGridAlgoPendingService {
	s.params["instId"] = instId
	return s
}

// SetAlgoId filters by algo order id.
func (s *GetGridAlgoPendingService) SetAlgoId(algoId string) *GetGridAlgoPendingService {
	s.params["algoId"] = algoId
	return s
}

// SetAfter paginates to records earlier than the given algoId (older).
func (s *GetGridAlgoPendingService) SetAfter(algoId string) *GetGridAlgoPendingService {
	s.params["after"] = algoId
	return s
}

// SetBefore paginates to records later than the given algoId (newer).
func (s *GetGridAlgoPendingService) SetBefore(algoId string) *GetGridAlgoPendingService {
	s.params["before"] = algoId
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetGridAlgoPendingService) SetLimit(limit int) *GetGridAlgoPendingService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetGridAlgoPendingService) Do(ctx context.Context) ([]GridAlgoOrder, error) {
	req := request.Get(ctx, s.c, "/api/v5/tradingBot/grid/orders-algo-pending", s.params).WithSign()
	return request.DoList[GridAlgoOrder](req)
}

// GetGridAlgoHistoryService -- GET /api/v5/tradingBot/grid/orders-algo-history (Read)
//
// Returns the account's stopped (historical) grid algo orders.
type GetGridAlgoHistoryService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetGridAlgoHistoryService(algoOrdType GridAlgoOrdType) *GetGridAlgoHistoryService {
	return &GetGridAlgoHistoryService{c: c, params: map[string]string{"algoOrdType": string(algoOrdType)}}
}

// SetInstType filters by product line.
func (s *GetGridAlgoHistoryService) SetInstType(instType InstType) *GetGridAlgoHistoryService {
	s.params["instType"] = string(instType)
	return s
}

// SetInstId filters by instrument id.
func (s *GetGridAlgoHistoryService) SetInstId(instId string) *GetGridAlgoHistoryService {
	s.params["instId"] = instId
	return s
}

// SetAlgoId filters by algo order id.
func (s *GetGridAlgoHistoryService) SetAlgoId(algoId string) *GetGridAlgoHistoryService {
	s.params["algoId"] = algoId
	return s
}

// SetAfter paginates to records earlier than the given algoId (older).
func (s *GetGridAlgoHistoryService) SetAfter(algoId string) *GetGridAlgoHistoryService {
	s.params["after"] = algoId
	return s
}

// SetBefore paginates to records later than the given algoId (newer).
func (s *GetGridAlgoHistoryService) SetBefore(algoId string) *GetGridAlgoHistoryService {
	s.params["before"] = algoId
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetGridAlgoHistoryService) SetLimit(limit int) *GetGridAlgoHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetGridAlgoHistoryService) Do(ctx context.Context) ([]GridAlgoOrder, error) {
	req := request.Get(ctx, s.c, "/api/v5/tradingBot/grid/orders-algo-history", s.params).WithSign()
	return request.DoList[GridAlgoOrder](req)
}

// GetGridAlgoDetailsService -- GET /api/v5/tradingBot/grid/orders-algo-details (Read)
//
// Returns the full detail of a single grid algo order.
type GetGridAlgoDetailsService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetGridAlgoDetailsService(algoOrdType GridAlgoOrdType, algoId string) *GetGridAlgoDetailsService {
	return &GetGridAlgoDetailsService{c: c, params: map[string]string{
		"algoOrdType": string(algoOrdType),
		"algoId":      algoId,
	}}
}

func (s *GetGridAlgoDetailsService) Do(ctx context.Context) (*GridAlgoOrder, error) {
	req := request.Get(ctx, s.c, "/api/v5/tradingBot/grid/orders-algo-details", s.params).WithSign()
	return request.DoOne[GridAlgoOrder](req)
}

// GridAlgoOrder is a single grid algo order (spot/margin or contract grid). The
// validating account had no grid orders, so the field set is modeled from the
// OKX doc field table (union of spot-grid and contract-grid fields). Fields that
// apply to only one grid kind are simply empty for the other.
type GridAlgoOrder struct {
	AlgoID                 string             `json:"algoId"`
	AlgoClientOrderID      string             `json:"algoClOrdId"`
	AlgoOrderType          GridAlgoOrdType    `json:"algoOrdType"`
	InstrumentType         InstType           `json:"instType"`
	InstrumentID           string             `json:"instId"`
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
	CancelType             string             `json:"cancelType"`
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

	// --- spot/margin grid ("grid") ---
	BaseSize                decimal.Decimal `json:"baseSz"`
	QuoteSize               decimal.Decimal `json:"quoteSz"`
	BaseCurrency            string          `json:"baseCcy"`
	QuoteCurrency           string          `json:"quoteCcy"`
	TradeMode               TdMode          `json:"tdMode"`
	Leverage                decimal.Decimal `json:"lever"`
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

// GridRebateTrans is one rebate-transfer record attached to a grid algo order.
type GridRebateTrans struct {
	Rebate         decimal.Decimal `json:"rebate"`
	RebateCurrency string          `json:"rebateCcy"`
}

// GridTriggerParam is one advanced trigger configuration of a grid algo order
// (used both in the order detail and in the amend request body).
type GridTriggerParam struct {
	TriggerAction   string          `json:"triggerAction"`
	TriggerStrategy string          `json:"triggerStrategy"`
	DelaySeconds    string          `json:"delaySeconds"`
	TriggerType     string          `json:"triggerType"`
	Timeframe       string          `json:"timeframe"`
	Thold           string          `json:"thold"`
	TriggerCond     string          `json:"triggerCond"`
	TimePeriod      string          `json:"timePeriod"`
	TriggerPrice    decimal.Decimal `json:"triggerPx"`
	StopType        string          `json:"stopType"`
	TriggerTime     time.Time       `json:"triggerTime"`
}

// GetGridSubOrdersService -- GET /api/v5/tradingBot/grid/sub-orders (Read)
//
// Returns the live or filled sub-orders of a grid algo order.
type GetGridSubOrdersService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetGridSubOrdersService(algoId string, algoOrdType GridAlgoOrdType, typ GridSubOrderType) *GetGridSubOrdersService {
	return &GetGridSubOrdersService{c: c, params: map[string]string{
		"algoId":      algoId,
		"algoOrdType": string(algoOrdType),
		"type":        string(typ),
	}}
}

// SetGroupId filters by grid group id.
func (s *GetGridSubOrdersService) SetGroupId(groupId string) *GetGridSubOrdersService {
	s.params["groupId"] = groupId
	return s
}

// SetInstType filters by product line.
func (s *GetGridSubOrdersService) SetInstType(instType InstType) *GetGridSubOrdersService {
	s.params["instType"] = string(instType)
	return s
}

// SetAfter paginates to records earlier than the given ordId (older).
func (s *GetGridSubOrdersService) SetAfter(ordId string) *GetGridSubOrdersService {
	s.params["after"] = ordId
	return s
}

// SetBefore paginates to records later than the given ordId (newer).
func (s *GetGridSubOrdersService) SetBefore(ordId string) *GetGridSubOrdersService {
	s.params["before"] = ordId
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetGridSubOrdersService) SetLimit(limit int) *GetGridSubOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetGridSubOrdersService) Do(ctx context.Context) ([]GridSubOrder, error) {
	req := request.Get(ctx, s.c, "/api/v5/tradingBot/grid/sub-orders", s.params).WithSign()
	return request.DoList[GridSubOrder](req)
}

// GridSubOrder is one sub-order (working leg) of a grid algo order. The
// validating account had no grid orders, so the field set is modeled from the
// OKX doc field table.
type GridSubOrder struct {
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

// GetGridPositionsService -- GET /api/v5/tradingBot/grid/positions (Read)
//
// Returns the open positions held by a contract grid algo order.
type GetGridPositionsService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetGridPositionsService(algoOrdType GridAlgoOrdType, algoId string) *GetGridPositionsService {
	return &GetGridPositionsService{c: c, params: map[string]string{
		"algoOrdType": string(algoOrdType),
		"algoId":      algoId,
	}}
}

func (s *GetGridPositionsService) Do(ctx context.Context) ([]GridPosition, error) {
	req := request.Get(ctx, s.c, "/api/v5/tradingBot/grid/positions", s.params).WithSign()
	return request.DoList[GridPosition](req)
}

// GridPosition is the position held by a contract grid algo order. The
// validating account had no grid orders, so the field set is modeled from the
// OKX doc field table.
type GridPosition struct {
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

// WithdrawGridIncomeService -- POST /api/v5/tradingBot/grid/withdraw-income (Trade)
//
// Withdraws the realized grid profit of a spot grid to the trading account.
// IMPLEMENT-ONLY.
type WithdrawGridIncomeService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewWithdrawGridIncomeService(algoId string) *WithdrawGridIncomeService {
	return &WithdrawGridIncomeService{c: c, body: map[string]any{"algoId": algoId}}
}

func (s *WithdrawGridIncomeService) Do(ctx context.Context) (*GridWithdrawIncomeResult, error) {
	req := request.Post(ctx, s.c, "/api/v5/tradingBot/grid/withdraw-income", s.body).WithSign()
	return request.DoOne[GridWithdrawIncomeResult](req)
}

// GridWithdrawIncomeResult is the ack of a grid income withdrawal.
type GridWithdrawIncomeResult struct {
	AlgoID string          `json:"algoId"`
	Profit decimal.Decimal `json:"profit"`
}

// ComputeGridMarginBalanceService -- POST /api/v5/tradingBot/grid/compute-margin-balance (Trade)
//
// Estimates the impact of adding/reducing the margin of a contract grid.
// IMPLEMENT-ONLY.
type ComputeGridMarginBalanceService struct {
	c    *Client
	body map[string]any
}

// NewComputeGridMarginBalanceService builds the estimate. typ is "add" or
// "reduce".
func (c *Client) NewComputeGridMarginBalanceService(algoId, typ string) *ComputeGridMarginBalanceService {
	return &ComputeGridMarginBalanceService{c: c, body: map[string]any{
		"algoId": algoId,
		"type":   typ,
	}}
}

// SetAmt sets the margin amount to add/reduce.
func (s *ComputeGridMarginBalanceService) SetAmt(amt decimal.Decimal) *ComputeGridMarginBalanceService {
	s.body["amt"] = amt.String()
	return s
}

func (s *ComputeGridMarginBalanceService) Do(ctx context.Context) (*GridComputeMarginBalance, error) {
	req := request.Post(ctx, s.c, "/api/v5/tradingBot/grid/compute-margin-balance", s.body).WithSign()
	return request.DoOne[GridComputeMarginBalance](req)
}

// GridComputeMarginBalance is the estimated result of a margin-balance change.
type GridComputeMarginBalance struct {
	Leverage  decimal.Decimal `json:"lever"`
	MaxAmount decimal.Decimal `json:"maxAmt"`
}

// AdjustGridMarginBalanceService -- POST /api/v5/tradingBot/grid/margin-balance (Trade)
//
// Adds or reduces the margin of a running contract grid. IMPLEMENT-ONLY.
type AdjustGridMarginBalanceService struct {
	c    *Client
	body map[string]any
}

// NewAdjustGridMarginBalanceService builds the request. typ is "add" or "reduce".
func (c *Client) NewAdjustGridMarginBalanceService(algoId, typ string) *AdjustGridMarginBalanceService {
	return &AdjustGridMarginBalanceService{c: c, body: map[string]any{
		"algoId": algoId,
		"type":   typ,
	}}
}

// SetAmt sets the margin amount to add/reduce (mutually exclusive with percent).
func (s *AdjustGridMarginBalanceService) SetAmt(amt decimal.Decimal) *AdjustGridMarginBalanceService {
	s.body["amt"] = amt.String()
	return s
}

// SetPercent sets the percentage of margin to add/reduce.
func (s *AdjustGridMarginBalanceService) SetPercent(percent decimal.Decimal) *AdjustGridMarginBalanceService {
	s.body["percent"] = percent.String()
	return s
}

func (s *AdjustGridMarginBalanceService) Do(ctx context.Context) (*GridMarginBalanceResult, error) {
	req := request.Post(ctx, s.c, "/api/v5/tradingBot/grid/margin-balance", s.body).WithSign()
	return request.DoOne[GridMarginBalanceResult](req)
}

// GridMarginBalanceResult is the ack of a margin-balance adjustment.
type GridMarginBalanceResult struct {
	AlgoID string `json:"algoId"`
}

// AdjustGridInvestmentService -- POST /api/v5/tradingBot/grid/adjust-investment (Trade)
//
// Adds investment to a running spot grid. IMPLEMENT-ONLY.
type AdjustGridInvestmentService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewAdjustGridInvestmentService(algoId string, amt decimal.Decimal) *AdjustGridInvestmentService {
	return &AdjustGridInvestmentService{c: c, body: map[string]any{
		"algoId": algoId,
		"amt":    amt.String(),
	}}
}

func (s *AdjustGridInvestmentService) Do(ctx context.Context) (*GridAdjustInvestmentResult, error) {
	req := request.Post(ctx, s.c, "/api/v5/tradingBot/grid/adjust-investment", s.body).WithSign()
	return request.DoOne[GridAdjustInvestmentResult](req)
}

// GridAdjustInvestmentResult is the ack of a grid investment top-up.
type GridAdjustInvestmentResult struct {
	AlgoID     string          `json:"algoId"`
	Investment decimal.Decimal `json:"investment"`
}

// GetGridAIParamService -- GET /api/v5/tradingBot/grid/ai-param (public)
//
// Returns OKX's AI-recommended grid parameters for an instrument. Public (no
// signing required).
type GetGridAIParamService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetGridAIParamService(algoOrdType GridAlgoOrdType, instId string) *GetGridAIParamService {
	return &GetGridAIParamService{c: c, params: map[string]string{
		"algoOrdType": string(algoOrdType),
		"instId":      instId,
	}}
}

// SetDirection sets the contract grid direction (contract_grid only).
func (s *GetGridAIParamService) SetDirection(direction GridDirection) *GetGridAIParamService {
	s.params["direction"] = string(direction)
	return s
}

// SetDuration sets the back-test duration (e.g. "7D", "30D", "180D").
func (s *GetGridAIParamService) SetDuration(duration string) *GetGridAIParamService {
	s.params["duration"] = duration
	return s
}

func (s *GetGridAIParamService) Do(ctx context.Context) ([]GridAIParam, error) {
	req := request.Get(ctx, s.c, "/api/v5/tradingBot/grid/ai-param", s.params)
	return request.DoList[GridAIParam](req)
}

// GridAIParam is OKX's AI-recommended grid parameters for an instrument.
type GridAIParam struct {
	AlgoOrderType      GridAlgoOrdType `json:"algoOrdType"`
	AnnualizedRate     decimal.Decimal `json:"annualizedRate"`
	Currency           string          `json:"ccy"`
	Direction          GridDirection   `json:"direction"`
	Duration           string          `json:"duration"`
	GridNumber         decimal.Decimal `json:"gridNum"`
	InstrumentID       string          `json:"instId"`
	Leverage           decimal.Decimal `json:"lever"`
	MaxPrice           decimal.Decimal `json:"maxPx"`
	MinPrice           decimal.Decimal `json:"minPx"`
	MinInvestment      decimal.Decimal `json:"minInvestment"`
	PerMaxProfitRate   decimal.Decimal `json:"perMaxProfitRate"`
	PerMinProfitRate   decimal.Decimal `json:"perMinProfitRate"`
	PerGridProfitRatio decimal.Decimal `json:"perGridProfitRatio"`
	RunType            GridRunType     `json:"runType"`
	SourceCurrency     string          `json:"sourceCcy"`
}

// GetGridMinInvestmentService -- POST /api/v5/tradingBot/grid/min-investment (public calc)
//
// Computes the minimum investment required for a given grid configuration. This
// is a public stateless calculator (no signing required), but its body carries
// several conditional fields, so it is provided implement-only and is not
// exercised by the live test. Returns the minimum investable amount per
// currency.
type GetGridMinInvestmentService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewGetGridMinInvestmentService(instId string, algoOrdType GridAlgoOrdType, maxPx, minPx decimal.Decimal, gridNum int, runType GridRunType) *GetGridMinInvestmentService {
	return &GetGridMinInvestmentService{c: c, body: map[string]any{
		"instId":      instId,
		"algoOrdType": string(algoOrdType),
		"maxPx":       maxPx.String(),
		"minPx":       minPx.String(),
		"gridNum":     strconv.Itoa(gridNum),
		"runType":     string(runType),
	}}
}

// SetAmt sets the planned investment amount.
func (s *GetGridMinInvestmentService) SetAmt(amt decimal.Decimal) *GetGridMinInvestmentService {
	s.body["amt"] = amt.String()
	return s
}

// SetInvestmentType sets the investment type ("quote", "base").
func (s *GetGridMinInvestmentService) SetInvestmentType(investmentType string) *GetGridMinInvestmentService {
	s.body["investmentType"] = investmentType
	return s
}

// SetDirection sets the contract grid direction (contract_grid only).
func (s *GetGridMinInvestmentService) SetDirection(direction GridDirection) *GetGridMinInvestmentService {
	s.body["direction"] = string(direction)
	return s
}

// SetLever sets the leverage (contract_grid only).
func (s *GetGridMinInvestmentService) SetLever(lever decimal.Decimal) *GetGridMinInvestmentService {
	s.body["lever"] = lever.String()
	return s
}

// SetBasePos toggles whether to open a base position (contract_grid only).
func (s *GetGridMinInvestmentService) SetBasePos(basePos bool) *GetGridMinInvestmentService {
	s.body["basePos"] = basePos
	return s
}

// SetInvestmentData sets per-currency planned investment entries.
func (s *GetGridMinInvestmentService) SetInvestmentData(data []GridInvestmentDataInput) *GetGridMinInvestmentService {
	s.body["investmentData"] = data
	return s
}

func (s *GetGridMinInvestmentService) Do(ctx context.Context) (*GridMinInvestment, error) {
	req := request.Post(ctx, s.c, "/api/v5/tradingBot/grid/min-investment", s.body)
	return request.DoOne[GridMinInvestment](req)
}

// GridInvestmentDataInput is one per-currency planned investment entry of a
// min-investment request.
type GridInvestmentDataInput struct {
	Amount   decimal.Decimal `json:"amt"`
	Currency string          `json:"ccy"`
}

// GridMinInvestment is the minimum-investment computation result.
type GridMinInvestment struct {
	MinInvestmentData []GridMinInvestmentData `json:"minInvestmentData"`
	SingleAmount      decimal.Decimal         `json:"singleAmt"`
}

// GridMinInvestmentData is one currency's minimum investable amount.
type GridMinInvestmentData struct {
	Amount   decimal.Decimal `json:"amt"`
	Currency string          `json:"ccy"`
}

// GetGridRSIBackTestingService -- GET /api/v5/tradingBot/public/rsi-back-testing (public)
//
// Back-tests how many times an RSI trigger condition would have fired for an
// instrument over a duration. Public (no signing required).
type GetGridRSIBackTestingService struct {
	c      *Client
	params map[string]string
}

// NewGetGridRSIBackTestingService builds the back-test. triggerCond is e.g.
// "cross_up"/"cross_down"/"above"/"below"; timeframe is e.g. "3m"/"15m"/"1H".
func (c *Client) NewGetGridRSIBackTestingService(instId, timeframe, thold string, timePeriod int, triggerCond, duration string) *GetGridRSIBackTestingService {
	return &GetGridRSIBackTestingService{c: c, params: map[string]string{
		"instId":      instId,
		"timeframe":   timeframe,
		"thold":       thold,
		"timePeriod":  strconv.Itoa(timePeriod),
		"triggerCond": triggerCond,
		"duration":    duration,
	}}
}

func (s *GetGridRSIBackTestingService) Do(ctx context.Context) (*GridRSIBackTesting, error) {
	req := request.Get(ctx, s.c, "/api/v5/tradingBot/public/rsi-back-testing", s.params)
	return request.DoOne[GridRSIBackTesting](req)
}

// GridRSIBackTesting is the RSI back-test result (how many times the condition
// fired over the duration).
type GridRSIBackTesting struct {
	TriggerNumber string `json:"triggerNum"`
}
