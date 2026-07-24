package okx

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-okx/request"
	"github.com/shopspring/decimal"
)

// AlgoOrdType is the type of an algo (trigger / advanced) order.
type AlgoOrdType string

const (
	AlgoOrdTypeConditional   AlgoOrdType = "conditional"
	AlgoOrdTypeOCO           AlgoOrdType = "oco"
	AlgoOrdTypeTrigger       AlgoOrdType = "trigger"
	AlgoOrdTypeMoveOrderStop AlgoOrdType = "move_order_stop"
	AlgoOrdTypeTwap          AlgoOrdType = "twap"
	AlgoOrdTypeIceberg       AlgoOrdType = "iceberg"
	AlgoOrdTypeChase         AlgoOrdType = "chase"
)

// AlgoState is the lifecycle state of an algo order. "live"/"pause" appear on the
// pending endpoint; "effective"/"canceled"/"order_failed"/"partially_failed"
// appear on the history endpoint.
type AlgoState string

const (
	AlgoStateLive            AlgoState = "live"
	AlgoStatePause           AlgoState = "pause"
	AlgoStatePartiallyEff    AlgoState = "partially_effective"
	AlgoStateEffective       AlgoState = "effective"
	AlgoStateCanceled        AlgoState = "canceled"
	AlgoStateOrderFailed     AlgoState = "order_failed"
	AlgoStatePartiallyFailed AlgoState = "partially_failed"
)

// --- State-changing endpoints (Trade): implemented but NEVER exercised by the
// test suite. Bodies and acks are modeled from the OKX docs. ---

// PlaceAlgoOrderService -- POST /api/v5/trade/order-algo (Trade)
//
// Places a single algo (trigger / advanced) order: conditional, oco, trigger,
// move_order_stop (trailing), twap, iceberg or chase. The required body fields
// (instId, tdMode, side, ordType, sz) are constructor args; the many type-specific
// fields are optional setters.
type PlaceAlgoOrderService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewPlaceAlgoOrderService(instId string, tdMode TdMode, side Side, ordType AlgoOrdType, sz decimal.Decimal) *PlaceAlgoOrderService {
	return &PlaceAlgoOrderService{c: c, body: map[string]any{
		"instId":  instId,
		"tdMode":  string(tdMode),
		"side":    string(side),
		"ordType": string(ordType),
		"sz":      sz.String(),
	}}
}

// SetCcy sets the margin currency (single-currency margin only).
func (s *PlaceAlgoOrderService) SetCcy(ccy string) *PlaceAlgoOrderService {
	s.body["ccy"] = ccy
	return s
}

// SetPosSide sets the position side (long/short/net).
func (s *PlaceAlgoOrderService) SetPosSide(posSide PosSide) *PlaceAlgoOrderService {
	s.body["posSide"] = string(posSide)
	return s
}

// SetReduceOnly toggles reduce-only (MARGIN/FUTURES/SWAP).
func (s *PlaceAlgoOrderService) SetReduceOnly(reduceOnly bool) *PlaceAlgoOrderService {
	s.body["reduceOnly"] = reduceOnly
	return s
}

// SetTgtCcy sets the spot market size unit (base_ccy / quote_ccy).
func (s *PlaceAlgoOrderService) SetTgtCcy(tgtCcy TgtCcy) *PlaceAlgoOrderService {
	s.body["tgtCcy"] = string(tgtCcy)
	return s
}

// SetAlgoClOrdId sets a client-supplied algo order id.
func (s *PlaceAlgoOrderService) SetAlgoClOrdId(algoClOrdId string) *PlaceAlgoOrderService {
	s.body["algoClOrdId"] = algoClOrdId
	return s
}

// SetClOrdId sets a client-supplied order id for the resulting order.
func (s *PlaceAlgoOrderService) SetClOrdId(clOrdId string) *PlaceAlgoOrderService {
	s.body["clOrdId"] = clOrdId
	return s
}

// SetTag sets an order tag.
func (s *PlaceAlgoOrderService) SetTag(tag string) *PlaceAlgoOrderService {
	s.body["tag"] = tag
	return s
}

// --- conditional / oco take-profit & stop-loss fields ---

// SetTpTriggerPx sets the take-profit trigger price (conditional/oco).
func (s *PlaceAlgoOrderService) SetTpTriggerPx(px decimal.Decimal) *PlaceAlgoOrderService {
	s.body["tpTriggerPx"] = px.String()
	return s
}

// SetTpTriggerPxType sets the take-profit trigger price type (last/index/mark).
func (s *PlaceAlgoOrderService) SetTpTriggerPxType(pxType string) *PlaceAlgoOrderService {
	s.body["tpTriggerPxType"] = pxType
	return s
}

// SetTpOrdPx sets the take-profit order price ("-1" for market).
func (s *PlaceAlgoOrderService) SetTpOrdPx(px decimal.Decimal) *PlaceAlgoOrderService {
	s.body["tpOrdPx"] = px.String()
	return s
}

// SetTpOrdKind sets the take-profit order kind (condition / limit).
func (s *PlaceAlgoOrderService) SetTpOrdKind(kind string) *PlaceAlgoOrderService {
	s.body["tpOrdKind"] = kind
	return s
}

// SetSlTriggerPx sets the stop-loss trigger price (conditional/oco).
func (s *PlaceAlgoOrderService) SetSlTriggerPx(px decimal.Decimal) *PlaceAlgoOrderService {
	s.body["slTriggerPx"] = px.String()
	return s
}

// SetSlTriggerPxType sets the stop-loss trigger price type (last/index/mark).
func (s *PlaceAlgoOrderService) SetSlTriggerPxType(pxType string) *PlaceAlgoOrderService {
	s.body["slTriggerPxType"] = pxType
	return s
}

// SetSlOrdPx sets the stop-loss order price ("-1" for market).
func (s *PlaceAlgoOrderService) SetSlOrdPx(px decimal.Decimal) *PlaceAlgoOrderService {
	s.body["slOrdPx"] = px.String()
	return s
}

// SetCxlOnClosePos toggles whether the algo order is canceled when the position
// is fully closed (only for TP/SL orders bound to a position).
func (s *PlaceAlgoOrderService) SetCxlOnClosePos(cxl bool) *PlaceAlgoOrderService {
	s.body["cxlOnClosePos"] = cxl
	return s
}

// --- trigger fields ---

// SetTriggerPx sets the trigger price (trigger ordType).
func (s *PlaceAlgoOrderService) SetTriggerPx(px decimal.Decimal) *PlaceAlgoOrderService {
	s.body["triggerPx"] = px.String()
	return s
}

// SetTriggerPxType sets the trigger price type (last/index/mark).
func (s *PlaceAlgoOrderService) SetTriggerPxType(pxType string) *PlaceAlgoOrderService {
	s.body["triggerPxType"] = pxType
	return s
}

// SetOrdPx sets the order price for trigger/iceberg/twap ("-1" for market).
func (s *PlaceAlgoOrderService) SetOrdPx(px decimal.Decimal) *PlaceAlgoOrderService {
	s.body["ordPx"] = px.String()
	return s
}

// SetAttachAlgoOrds attaches a list of TP/SL orders to a trigger order.
func (s *PlaceAlgoOrderService) SetAttachAlgoOrds(ords []AlgoAttachOrd) *PlaceAlgoOrderService {
	s.body["attachAlgoOrds"] = ords
	return s
}

// --- move_order_stop (trailing) fields ---

// SetCallbackRatio sets the trailing callback ratio (move_order_stop).
func (s *PlaceAlgoOrderService) SetCallbackRatio(ratio decimal.Decimal) *PlaceAlgoOrderService {
	s.body["callbackRatio"] = ratio.String()
	return s
}

// SetCallbackSpread sets the trailing callback spread (move_order_stop).
func (s *PlaceAlgoOrderService) SetCallbackSpread(spread decimal.Decimal) *PlaceAlgoOrderService {
	s.body["callbackSpread"] = spread.String()
	return s
}

// SetActivePx sets the activation price (move_order_stop).
func (s *PlaceAlgoOrderService) SetActivePx(px decimal.Decimal) *PlaceAlgoOrderService {
	s.body["activePx"] = px.String()
	return s
}

// --- iceberg / twap shared fields ---

// SetPxVar sets the price ratio versus the latest price (iceberg/twap, % offset).
func (s *PlaceAlgoOrderService) SetPxVar(pxVar decimal.Decimal) *PlaceAlgoOrderService {
	s.body["pxVar"] = pxVar.String()
	return s
}

// SetPxSpread sets the price variance per child order (iceberg/twap, absolute).
func (s *PlaceAlgoOrderService) SetPxSpread(pxSpread decimal.Decimal) *PlaceAlgoOrderService {
	s.body["pxSpread"] = pxSpread.String()
	return s
}

// SetSzLimit sets the average amount per child order (iceberg/twap).
func (s *PlaceAlgoOrderService) SetSzLimit(szLimit decimal.Decimal) *PlaceAlgoOrderService {
	s.body["szLimit"] = szLimit.String()
	return s
}

// SetPxLimit sets the price limit per child order (iceberg/twap).
func (s *PlaceAlgoOrderService) SetPxLimit(pxLimit decimal.Decimal) *PlaceAlgoOrderService {
	s.body["pxLimit"] = pxLimit.String()
	return s
}

// SetTimeInterval sets the interval in seconds between child orders (twap).
func (s *PlaceAlgoOrderService) SetTimeInterval(seconds int) *PlaceAlgoOrderService {
	s.body["timeInterval"] = strconv.Itoa(seconds)
	return s
}

// --- chase fields ---

// SetChaseType sets the chase distance type (distance / ratio).
func (s *PlaceAlgoOrderService) SetChaseType(chaseType string) *PlaceAlgoOrderService {
	s.body["chaseType"] = chaseType
	return s
}

// SetChaseVal sets the chase distance value (chase).
func (s *PlaceAlgoOrderService) SetChaseVal(val decimal.Decimal) *PlaceAlgoOrderService {
	s.body["chaseVal"] = val.String()
	return s
}

// SetMaxChaseType sets the max-chase distance type (distance / ratio).
func (s *PlaceAlgoOrderService) SetMaxChaseType(maxChaseType string) *PlaceAlgoOrderService {
	s.body["maxChaseType"] = maxChaseType
	return s
}

// SetMaxChaseVal sets the max-chase distance value (chase).
func (s *PlaceAlgoOrderService) SetMaxChaseVal(val decimal.Decimal) *PlaceAlgoOrderService {
	s.body["maxChaseVal"] = val.String()
	return s
}

// --- advance order type (chase trigger) fields ---

// SetAdvanceOrdType sets the order type spawned when a FUTURES/SWAP trigger order
// fires. Value "chase" makes the triggered order a chase limit order; ordPx is
// then not required and the chase parameters come from SetAdvChaseParams.
func (s *PlaceAlgoOrderService) SetAdvanceOrdType(advanceOrdType string) *PlaceAlgoOrderService {
	s.body["advanceOrdType"] = advanceOrdType
	return s
}

// SetAdvChaseParams sets the chase parameters carried by a trigger order whose
// advanceOrdType is "chase" (FUTURES/SWAP).
func (s *PlaceAlgoOrderService) SetAdvChaseParams(params []AdvChaseParam) *PlaceAlgoOrderService {
	s.body["advChaseParams"] = params
	return s
}

// SetQuickMgnType sets the quick-margin borrow type (manual/auto_borrow/auto_borrow_repay).
func (s *PlaceAlgoOrderService) SetQuickMgnType(quickMgnType string) *PlaceAlgoOrderService {
	s.body["quickMgnType"] = quickMgnType
	return s
}

func (s *PlaceAlgoOrderService) Do(ctx context.Context) (*AlgoResult, error) {
	req := request.Post(ctx, s.c, "/api/v5/trade/order-algo", s.body).WithSign()
	list, err := request.DoListPartial[AlgoResult](req)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return &list[0], nil
}

// AlgoAttachOrd is a take-profit / stop-loss order attached to a trigger order.
type AlgoAttachOrd struct {
	AttachAlgoClientOrderID    string          `json:"attachAlgoClOrdId,omitempty"`
	TakeProfitTriggerPrice     decimal.Decimal `json:"tpTriggerPx,omitzero"`
	TakeProfitTriggerPriceType string          `json:"tpTriggerPxType,omitempty"`
	TakeProfitOrderPrice       decimal.Decimal `json:"tpOrdPx,omitzero"`
	StopLossTriggerPrice       decimal.Decimal `json:"slTriggerPx,omitzero"`
	StopLossTriggerPriceType   string          `json:"slTriggerPxType,omitempty"`
	StopLossOrderPrice         decimal.Decimal `json:"slOrdPx,omitzero"`
	Size                       decimal.Decimal `json:"sz,omitzero"`
	AmendPriceOnTriggerType    string          `json:"amendPxOnTriggerType,omitempty"`
}

// AdvChaseParam carries the chase execution parameters of a FUTURES/SWAP trigger
// order whose advanceOrdType is "chase": when the trigger fires the spawned order
// chases the order book. Mirrors the standalone chase order fields.
type AdvChaseParam struct {
	ChaseType     string          `json:"chaseType,omitempty"`
	ChaseValue    decimal.Decimal `json:"chaseVal,omitzero"`
	MaxChaseType  string          `json:"maxChaseType,omitempty"`
	MaxChaseValue decimal.Decimal `json:"maxChaseVal,omitzero"`
}

// AlgoResult is the per-item ack returned by the algo place / amend / cancel
// endpoints. On a top-level code of "1" the real reason is in sCode/sMsg.
type AlgoResult struct {
	AlgoID            string `json:"algoId"`
	AlgoClientOrderID string `json:"algoClOrdId"`
	ClientOrderID     string `json:"clOrdId"`
	Tag               string `json:"tag"`
	RequestID         string `json:"reqId"`
	SCode             string `json:"sCode"`
	SMsg              string `json:"sMsg"`
}

// CancelAlgoOrdersService -- POST /api/v5/trade/cancel-algos (Trade)
//
// Cancels up to 10 pending algo orders in one request. The body is a JSON ARRAY
// of {algoId, instId} items.
type CancelAlgoOrdersService struct {
	c     *Client
	items []AlgoCancelArg
}

func (c *Client) NewCancelAlgoOrdersService(items []AlgoCancelArg) *CancelAlgoOrdersService {
	return &CancelAlgoOrdersService{c: c, items: items}
}

func (s *CancelAlgoOrdersService) Do(ctx context.Context) ([]AlgoResult, error) {
	req := request.Post(ctx, s.c, "/api/v5/trade/cancel-algos").SetBody(s.items).WithSign()
	return request.DoListPartial[AlgoResult](req)
}

// AlgoCancelArg identifies one algo order to cancel.
type AlgoCancelArg struct {
	AlgoID       string `json:"algoId"`
	InstrumentID string `json:"instId"`
}

// AmendAlgoOrderService -- POST /api/v5/trade/amend-algos (Trade)
//
// Amends a single pending algo order. One of algoId / algoClOrdId is required.
type AmendAlgoOrderService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewAmendAlgoOrderService(instId string) *AmendAlgoOrderService {
	return &AmendAlgoOrderService{c: c, body: map[string]any{"instId": instId}}
}

// SetAlgoId targets the algo order by id (one of algoId / algoClOrdId required).
func (s *AmendAlgoOrderService) SetAlgoId(algoId string) *AmendAlgoOrderService {
	s.body["algoId"] = algoId
	return s
}

// SetAlgoClOrdId targets the algo order by client-supplied id.
func (s *AmendAlgoOrderService) SetAlgoClOrdId(algoClOrdId string) *AmendAlgoOrderService {
	s.body["algoClOrdId"] = algoClOrdId
	return s
}

// SetCxlOnFail cancels the order if amendment fails.
func (s *AmendAlgoOrderService) SetCxlOnFail(cxl bool) *AmendAlgoOrderService {
	s.body["cxlOnFail"] = cxl
	return s
}

// SetReqId sets a client-supplied request id for the amendment.
func (s *AmendAlgoOrderService) SetReqId(reqId string) *AmendAlgoOrderService {
	s.body["reqId"] = reqId
	return s
}

// SetNewSz sets the new size.
func (s *AmendAlgoOrderService) SetNewSz(sz decimal.Decimal) *AmendAlgoOrderService {
	s.body["newSz"] = sz.String()
	return s
}

// SetNewTpTriggerPx sets the new take-profit trigger price.
func (s *AmendAlgoOrderService) SetNewTpTriggerPx(px decimal.Decimal) *AmendAlgoOrderService {
	s.body["newTpTriggerPx"] = px.String()
	return s
}

// SetNewTpTriggerPxType sets the new take-profit trigger price type.
func (s *AmendAlgoOrderService) SetNewTpTriggerPxType(pxType string) *AmendAlgoOrderService {
	s.body["newTpTriggerPxType"] = pxType
	return s
}

// SetNewTpOrdPx sets the new take-profit order price ("-1" for market).
func (s *AmendAlgoOrderService) SetNewTpOrdPx(px decimal.Decimal) *AmendAlgoOrderService {
	s.body["newTpOrdPx"] = px.String()
	return s
}

// SetNewTpOrdKind sets the new take-profit order kind (condition / limit).
func (s *AmendAlgoOrderService) SetNewTpOrdKind(kind string) *AmendAlgoOrderService {
	s.body["newTpOrdKind"] = kind
	return s
}

// SetNewSlTriggerPx sets the new stop-loss trigger price.
func (s *AmendAlgoOrderService) SetNewSlTriggerPx(px decimal.Decimal) *AmendAlgoOrderService {
	s.body["newSlTriggerPx"] = px.String()
	return s
}

// SetNewSlTriggerPxType sets the new stop-loss trigger price type.
func (s *AmendAlgoOrderService) SetNewSlTriggerPxType(pxType string) *AmendAlgoOrderService {
	s.body["newSlTriggerPxType"] = pxType
	return s
}

// SetNewSlOrdPx sets the new stop-loss order price ("-1" for market).
func (s *AmendAlgoOrderService) SetNewSlOrdPx(px decimal.Decimal) *AmendAlgoOrderService {
	s.body["newSlOrdPx"] = px.String()
	return s
}

// SetNewTriggerPx sets the new trigger price.
func (s *AmendAlgoOrderService) SetNewTriggerPx(px decimal.Decimal) *AmendAlgoOrderService {
	s.body["newTriggerPx"] = px.String()
	return s
}

// SetNewOrdPx sets the new order price ("-1" for market).
func (s *AmendAlgoOrderService) SetNewOrdPx(px decimal.Decimal) *AmendAlgoOrderService {
	s.body["newOrdPx"] = px.String()
	return s
}

// SetNewTriggerPxType sets the new trigger price type.
func (s *AmendAlgoOrderService) SetNewTriggerPxType(pxType string) *AmendAlgoOrderService {
	s.body["newTriggerPxType"] = pxType
	return s
}

// SetAttachAlgoOrds amends the attached TP/SL orders.
func (s *AmendAlgoOrderService) SetAttachAlgoOrds(ords []AlgoAttachOrd) *AmendAlgoOrderService {
	s.body["attachAlgoOrds"] = ords
	return s
}

func (s *AmendAlgoOrderService) Do(ctx context.Context) (*AlgoResult, error) {
	req := request.Post(ctx, s.c, "/api/v5/trade/amend-algos", s.body).WithSign()
	list, err := request.DoListPartial[AlgoResult](req)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return &list[0], nil
}

// CancelAdvanceAlgoOrdersService -- POST /api/v5/trade/cancel-advance-algos (Trade)
//
// Cancels up to 10 pending advanced algo orders (iceberg / twap / move_order_stop)
// in one request. The body is a JSON ARRAY of {algoId, instId} items.
type CancelAdvanceAlgoOrdersService struct {
	c     *Client
	items []AlgoCancelArg
}

func (c *Client) NewCancelAdvanceAlgoOrdersService(items []AlgoCancelArg) *CancelAdvanceAlgoOrdersService {
	return &CancelAdvanceAlgoOrdersService{c: c, items: items}
}

func (s *CancelAdvanceAlgoOrdersService) Do(ctx context.Context) ([]AlgoResult, error) {
	req := request.Post(ctx, s.c, "/api/v5/trade/cancel-advance-algos").SetBody(s.items).WithSign()
	return request.DoListPartial[AlgoResult](req)
}

// --- Read endpoints ---

// GetAlgoOrderService -- GET /api/v5/trade/order-algo (Read)
//
// Returns the details of a single algo order. One of algoId / algoClOrdId is
// required.
type GetAlgoOrderService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetAlgoOrderService() *GetAlgoOrderService {
	return &GetAlgoOrderService{c: c, params: map[string]string{}}
}

// SetAlgoId targets the algo order by id.
func (s *GetAlgoOrderService) SetAlgoId(algoId string) *GetAlgoOrderService {
	s.params["algoId"] = algoId
	return s
}

// SetAlgoClOrdId targets the algo order by client-supplied id.
func (s *GetAlgoOrderService) SetAlgoClOrdId(algoClOrdId string) *GetAlgoOrderService {
	s.params["algoClOrdId"] = algoClOrdId
	return s
}

func (s *GetAlgoOrderService) Do(ctx context.Context) (*AlgoOrder, error) {
	req := request.Get(ctx, s.c, "/api/v5/trade/order-algo", s.params).WithSign()
	return request.DoOne[AlgoOrder](req)
}

// GetAlgoOrdersPendingService -- GET /api/v5/trade/orders-algo-pending (Read)
//
// Returns the account's incomplete algo orders for a given ordType.
type GetAlgoOrdersPendingService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetAlgoOrdersPendingService(ordType AlgoOrdType) *GetAlgoOrdersPendingService {
	return &GetAlgoOrdersPendingService{c: c, params: map[string]string{"ordType": string(ordType)}}
}

// SetAlgoId filters by algo order id.
func (s *GetAlgoOrdersPendingService) SetAlgoId(algoId string) *GetAlgoOrdersPendingService {
	s.params["algoId"] = algoId
	return s
}

// SetAlgoClOrdId filters by client-supplied algo order id.
func (s *GetAlgoOrdersPendingService) SetAlgoClOrdId(algoClOrdId string) *GetAlgoOrdersPendingService {
	s.params["algoClOrdId"] = algoClOrdId
	return s
}

// SetInstType filters by product line (SPOT/MARGIN/SWAP/FUTURES).
func (s *GetAlgoOrdersPendingService) SetInstType(instType InstType) *GetAlgoOrdersPendingService {
	s.params["instType"] = string(instType)
	return s
}

// SetInstId filters by instrument id.
func (s *GetAlgoOrdersPendingService) SetInstId(instId string) *GetAlgoOrdersPendingService {
	s.params["instId"] = instId
	return s
}

// SetAfter pages to records with an algoId earlier than the given one.
func (s *GetAlgoOrdersPendingService) SetAfter(algoId string) *GetAlgoOrdersPendingService {
	s.params["after"] = algoId
	return s
}

// SetBefore pages to records with an algoId later than the given one.
func (s *GetAlgoOrdersPendingService) SetBefore(algoId string) *GetAlgoOrdersPendingService {
	s.params["before"] = algoId
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetAlgoOrdersPendingService) SetLimit(limit int) *GetAlgoOrdersPendingService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetAlgoOrdersPendingService) Do(ctx context.Context) ([]AlgoOrder, error) {
	req := request.Get(ctx, s.c, "/api/v5/trade/orders-algo-pending", s.params).WithSign()
	return request.DoList[AlgoOrder](req)
}

// GetAlgoOrdersHistoryService -- GET /api/v5/trade/orders-algo-history (Read)
//
// Returns the account's completed algo orders (canceled or effective) for a given
// ordType over the last three months. Either state or algoId is required in
// addition to ordType.
type GetAlgoOrdersHistoryService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetAlgoOrdersHistoryService(ordType AlgoOrdType) *GetAlgoOrdersHistoryService {
	return &GetAlgoOrdersHistoryService{c: c, params: map[string]string{"ordType": string(ordType)}}
}

// SetState filters by terminal state (effective / canceled / order_failed).
func (s *GetAlgoOrdersHistoryService) SetState(state AlgoState) *GetAlgoOrdersHistoryService {
	s.params["state"] = string(state)
	return s
}

// SetAlgoId filters by algo order id.
func (s *GetAlgoOrdersHistoryService) SetAlgoId(algoId string) *GetAlgoOrdersHistoryService {
	s.params["algoId"] = algoId
	return s
}

// SetInstType filters by product line (SPOT/MARGIN/SWAP/FUTURES).
func (s *GetAlgoOrdersHistoryService) SetInstType(instType InstType) *GetAlgoOrdersHistoryService {
	s.params["instType"] = string(instType)
	return s
}

// SetInstId filters by instrument id.
func (s *GetAlgoOrdersHistoryService) SetInstId(instId string) *GetAlgoOrdersHistoryService {
	s.params["instId"] = instId
	return s
}

// SetAfter pages to records with an algoId earlier than the given one.
func (s *GetAlgoOrdersHistoryService) SetAfter(algoId string) *GetAlgoOrdersHistoryService {
	s.params["after"] = algoId
	return s
}

// SetBefore pages to records with an algoId later than the given one.
func (s *GetAlgoOrdersHistoryService) SetBefore(algoId string) *GetAlgoOrdersHistoryService {
	s.params["before"] = algoId
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetAlgoOrdersHistoryService) SetLimit(limit int) *GetAlgoOrdersHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetAlgoOrdersHistoryService) Do(ctx context.Context) ([]AlgoOrder, error) {
	req := request.Get(ctx, s.c, "/api/v5/trade/orders-algo-history", s.params).WithSign()
	return request.DoList[AlgoOrder](req)
}

// AlgoOrder is a single algo (trigger / advanced) order. The validating account
// had no algo orders, so the field set is modeled from the OKX doc field tables
// for order-algo / orders-algo-pending / orders-algo-history (a union covering all
// algo order types).
type AlgoOrder struct {
	InstrumentType             InstType        `json:"instType"`
	InstrumentID               string          `json:"instId"`
	Currency                   string          `json:"ccy"`
	OrderID                    string          `json:"ordId"`
	OrderIDList                []string        `json:"ordIdList"`
	AlgoID                     string          `json:"algoId"`
	AlgoClientOrderID          string          `json:"algoClOrdId"`
	ClientOrderID              string          `json:"clOrdId"`
	Size                       decimal.Decimal `json:"sz"`
	CloseFraction              decimal.Decimal `json:"closeFraction"`
	OrderType                  AlgoOrdType     `json:"ordType"`
	Side                       Side            `json:"side"`
	PositionSide               PosSide         `json:"posSide"`
	TradeMode                  TdMode          `json:"tdMode"`
	TargetCurrency             TgtCcy          `json:"tgtCcy"`
	State                      AlgoState       `json:"state"`
	Leverage                   decimal.Decimal `json:"lever"`
	TakeProfitTriggerPrice     decimal.Decimal `json:"tpTriggerPx"`
	TakeProfitTriggerPriceType string          `json:"tpTriggerPxType"`
	TakeProfitOrderPrice       decimal.Decimal `json:"tpOrdPx"`
	TakeProfitOrderKind        string          `json:"tpOrdKind"`
	StopLossTriggerPrice       decimal.Decimal `json:"slTriggerPx"`
	StopLossTriggerPriceType   string          `json:"slTriggerPxType"`
	StopLossOrderPrice         decimal.Decimal `json:"slOrdPx"`
	TriggerPrice               decimal.Decimal `json:"triggerPx"`
	TriggerPriceType           string          `json:"triggerPxType"`
	OrderPrice                 decimal.Decimal `json:"ordPx"`
	ActualSize                 decimal.Decimal `json:"actualSz"`
	ActualPrice                decimal.Decimal `json:"actualPx"`
	ActualSide                 string          `json:"actualSide"`
	TriggerTime                time.Time       `json:"triggerTime"`
	PriceVariation             decimal.Decimal `json:"pxVar"`
	PriceSpread                decimal.Decimal `json:"pxSpread"`
	SizeLimit                  decimal.Decimal `json:"szLimit"`
	PriceLimit                 decimal.Decimal `json:"pxLimit"`
	TimeInterval               string          `json:"timeInterval"`
	CallbackRatio              decimal.Decimal `json:"callbackRatio"`
	CallbackSpread             decimal.Decimal `json:"callbackSpread"`
	ActivePrice                decimal.Decimal `json:"activePx"`
	MoveTriggerPrice           decimal.Decimal `json:"moveTriggerPx"`
	ChaseType                  string          `json:"chaseType"`
	ChaseValue                 decimal.Decimal `json:"chaseVal"`
	MaxChaseType               string          `json:"maxChaseType"`
	MaxChaseValue              decimal.Decimal `json:"maxChaseVal"`
	// SubAlgoIDList holds the algoId(s) of the chase order(s) spawned when a
	// trigger order with advanceOrdType "chase" fires (FUTURES/SWAP).
	SubAlgoIDList           []string            `json:"subAlgoIdList"`
	ReduceOnly              string              `json:"reduceOnly"`
	QuickMarginType         string              `json:"quickMgnType"`
	Last                    decimal.Decimal     `json:"last"`
	FailCode                string              `json:"failCode"`
	AlgoClientOrderIDParent string              `json:"algoClOrdIdParent"`
	AmendPriceOnTriggerType string              `json:"amendPxOnTriggerType"`
	LinkedOrder             *AlgoLinkedOrd      `json:"linkedOrd"`
	AttachAlgoOrders        []AlgoAttachOrdInfo `json:"attachAlgoOrds"`
	Tag                     string              `json:"tag"`
	CancelOnClosePosition   string              `json:"cxlOnClosePos"`
	IsTradeBorrowMode       string              `json:"isTradeBorrowMode"`
	CreationTime            time.Time           `json:"cTime"`
	UpdateTime              time.Time           `json:"uTime"`
}

// AlgoLinkedOrd is the order linked to an algo order (OCO/conditional).
type AlgoLinkedOrd struct {
	OrderID string `json:"ordId"`
}

// AlgoAttachOrdInfo is a take-profit / stop-loss order attached to a trigger
// order as returned by the read endpoints.
type AlgoAttachOrdInfo struct {
	AttachAlgoID               string          `json:"attachAlgoId"`
	AttachAlgoClientOrderID    string          `json:"attachAlgoClOrdId"`
	TakeProfitTriggerPrice     decimal.Decimal `json:"tpTriggerPx"`
	TakeProfitTriggerPriceType string          `json:"tpTriggerPxType"`
	TakeProfitOrderPrice       decimal.Decimal `json:"tpOrdPx"`
	StopLossTriggerPrice       decimal.Decimal `json:"slTriggerPx"`
	StopLossTriggerPriceType   string          `json:"slTriggerPxType"`
	StopLossOrderPrice         decimal.Decimal `json:"slOrdPx"`
	Size                       decimal.Decimal `json:"sz"`
	AmendPriceOnTriggerType    string          `json:"amendPxOnTriggerType"`
	FailCode                   string          `json:"failCode"`
	FailReason                 string          `json:"failReason"`
}
