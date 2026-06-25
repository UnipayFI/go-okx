package okx

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/UnipayFI/go-okx/request"
	"github.com/shopspring/decimal"
)

// MaxClOrdIDLen is the maximum length OKX accepts for a client order id
// (clOrdId). Longer ids are rejected with code 51000 ("Parameter clOrdId
// error"). OKX also requires the id to be alphanumeric (letters and digits
// only); see ValidateClOrdID.
const MaxClOrdIDLen = 32

// MaxBatchOrders is the maximum number of orders OKX accepts in a single
// batch-orders / cancel-batch-orders / amend-batch-orders request (REST) and in
// the WS batch-* ops. Larger batches are rejected; callers should chunk by this.
const MaxBatchOrders = 20

// ValidateClOrdID reports whether clOrdId satisfies OKX's client-order-id rule:
// at most MaxClOrdIDLen characters, letters and digits only (case-sensitive). It
// lets callers fail fast instead of round-tripping to a 51000 rejection. An
// empty clOrdId is allowed (the field is optional) and reported valid.
func ValidateClOrdID(clOrdId string) error {
	if clOrdId == "" {
		return nil
	}
	if len(clOrdId) > MaxClOrdIDLen {
		return fmt.Errorf("clOrdId %q too long: %d > %d", clOrdId, len(clOrdId), MaxClOrdIDLen)
	}
	for i := 0; i < len(clOrdId); i++ {
		c := clOrdId[i]
		if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z') {
			return fmt.Errorf("clOrdId %q has invalid character %q (must be alphanumeric)", clOrdId, string(c))
		}
	}
	return nil
}

// AttachAlgoOrd is an attached take-profit / stop-loss algo order that can be
// supplied when placing or amending an order. It is used both in the request
// body (attachAlgoOrds) and as a sub-object of Order in responses.
type AttachAlgoOrd struct {
	AttachAlgoID               string          `json:"attachAlgoId,omitempty"`
	AttachAlgoClientOrderID    string          `json:"attachAlgoClOrdId,omitempty"`
	TakeProfitTriggerPrice     decimal.Decimal `json:"tpTriggerPx,omitzero"`
	TakeProfitTriggerPriceType string          `json:"tpTriggerPxType,omitempty"`
	TakeProfitOrderPrice       decimal.Decimal `json:"tpOrdPx,omitzero"`
	TakeProfitOrderKind        string          `json:"tpOrdKind,omitempty"`
	StopLossTriggerPrice       decimal.Decimal `json:"slTriggerPx,omitzero"`
	StopLossTriggerPriceType   string          `json:"slTriggerPxType,omitempty"`
	StopLossOrderPrice         decimal.Decimal `json:"slOrdPx,omitzero"`
	Size                       decimal.Decimal `json:"sz,omitzero"`
	AmendPriceOnTriggerType    string          `json:"amendPxOnTriggerType,omitempty"`
}

// OrderArg is one order leg of a batch place-orders request body. It mirrors the
// single place-order body and is used by NewBatchOrdersService and
// NewPlaceOrderService.
type OrderArg struct {
	InstrumentID               string          `json:"instId"`
	TradeMode                  TdMode          `json:"tdMode"`
	Side                       Side            `json:"side"`
	OrderType                  OrdType         `json:"ordType"`
	Size                       decimal.Decimal `json:"sz"`
	Currency                   string          `json:"ccy,omitempty"`
	ClientOrderID              string          `json:"clOrdId,omitempty"`
	Tag                        string          `json:"tag,omitempty"`
	PositionSide               PosSide         `json:"posSide,omitempty"`
	Price                      decimal.Decimal `json:"px,omitzero"`
	PriceUSD                   decimal.Decimal `json:"pxUsd,omitzero"`
	PriceVolatility            decimal.Decimal `json:"pxVol,omitzero"`
	ReduceOnly                 bool            `json:"reduceOnly,omitempty"`
	TargetCurrency             TgtCcy          `json:"tgtCcy,omitempty"`
	BanAmend                   bool            `json:"banAmend,omitempty"`
	TakeProfitTriggerPrice     decimal.Decimal `json:"tpTriggerPx,omitzero"`
	TakeProfitOrderPrice       decimal.Decimal `json:"tpOrdPx,omitzero"`
	StopLossTriggerPrice       decimal.Decimal `json:"slTriggerPx,omitzero"`
	StopLossOrderPrice         decimal.Decimal `json:"slOrdPx,omitzero"`
	TakeProfitTriggerPriceType string          `json:"tpTriggerPxType,omitempty"`
	StopLossTriggerPriceType   string          `json:"slTriggerPxType,omitempty"`
	QuickMarginType            string          `json:"quickMgnType,omitempty"`
	STPID                      string          `json:"stpId,omitempty"`
	STPMode                    string          `json:"stpMode,omitempty"`
	AttachAlgoOrders           []AttachAlgoOrd `json:"attachAlgoOrds,omitempty"`
}

// OrderResult is the ack returned by the place / batch-place / cancel /
// batch-cancel order endpoints. The real per-item status lives in sCode/sMsg.
type OrderResult struct {
	OrderID       string    `json:"ordId"`
	ClientOrderID string    `json:"clOrdId"`
	Tag           string    `json:"tag"`
	Timestamp     time.Time `json:"ts"`
	SCode         string    `json:"sCode"`
	SMsg          string    `json:"sMsg"`
}

// AmendResult is the ack returned by the amend / batch-amend order endpoints.
type AmendResult struct {
	OrderID       string    `json:"ordId"`
	ClientOrderID string    `json:"clOrdId"`
	Timestamp     time.Time `json:"ts"`
	RequestID     string    `json:"reqId"`
	SCode         string    `json:"sCode"`
	SMsg          string    `json:"sMsg"`
}

// ClosePositionResult is the ack returned by the close-position endpoint.
type ClosePositionResult struct {
	InstrumentID  string  `json:"instId"`
	PositionSide  PosSide `json:"posSide"`
	ClientOrderID string  `json:"clOrdId"`
	Tag           string  `json:"tag"`
}

// PlaceOrderService -- POST /api/v5/trade/order (Trade)
//
// Places a single order. OKX may report failure with a top-level code "1" and
// the real reason in data[0].sCode, so the ack is read via DoListPartial and the
// first (only) element is returned, preserving sCode/sMsg.
type PlaceOrderService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewPlaceOrderService(instId string, tdMode TdMode, side Side, ordType OrdType, sz decimal.Decimal) *PlaceOrderService {
	return &PlaceOrderService{c: c, body: map[string]any{
		"instId":  instId,
		"tdMode":  string(tdMode),
		"side":    string(side),
		"ordType": string(ordType),
		"sz":      sz.String(),
	}}
}

// SetCcy sets the margin currency (single-currency margin mode only).
func (s *PlaceOrderService) SetCcy(ccy string) *PlaceOrderService {
	s.body["ccy"] = ccy
	return s
}

// SetClOrdId sets the client-supplied order id.
func (s *PlaceOrderService) SetClOrdId(clOrdId string) *PlaceOrderService {
	s.body["clOrdId"] = clOrdId
	return s
}

// SetTag sets the order tag.
func (s *PlaceOrderService) SetTag(tag string) *PlaceOrderService {
	s.body["tag"] = tag
	return s
}

// SetPosSide sets the position side (long/short/net).
func (s *PlaceOrderService) SetPosSide(posSide PosSide) *PlaceOrderService {
	s.body["posSide"] = string(posSide)
	return s
}

// SetPx sets the limit price (required for limit-type orders).
func (s *PlaceOrderService) SetPx(px decimal.Decimal) *PlaceOrderService {
	s.body["px"] = px.String()
	return s
}

// SetReduceOnly toggles reduce-only (cross MARGIN / FUTURES / SWAP).
func (s *PlaceOrderService) SetReduceOnly(reduceOnly bool) *PlaceOrderService {
	s.body["reduceOnly"] = reduceOnly
	return s
}

// SetTgtCcy selects the unit of a spot market order's size (base_ccy/quote_ccy).
func (s *PlaceOrderService) SetTgtCcy(tgtCcy TgtCcy) *PlaceOrderService {
	s.body["tgtCcy"] = string(tgtCcy)
	return s
}

// SetBanAmend forbids the order from being amended by the system on self-fill.
func (s *PlaceOrderService) SetBanAmend(banAmend bool) *PlaceOrderService {
	s.body["banAmend"] = banAmend
	return s
}

// SetTpTriggerPx sets the take-profit trigger price.
func (s *PlaceOrderService) SetTpTriggerPx(px decimal.Decimal) *PlaceOrderService {
	s.body["tpTriggerPx"] = px.String()
	return s
}

// SetTpOrdPx sets the take-profit order price ("-1" for market).
func (s *PlaceOrderService) SetTpOrdPx(px decimal.Decimal) *PlaceOrderService {
	s.body["tpOrdPx"] = px.String()
	return s
}

// SetSlTriggerPx sets the stop-loss trigger price.
func (s *PlaceOrderService) SetSlTriggerPx(px decimal.Decimal) *PlaceOrderService {
	s.body["slTriggerPx"] = px.String()
	return s
}

// SetSlOrdPx sets the stop-loss order price ("-1" for market).
func (s *PlaceOrderService) SetSlOrdPx(px decimal.Decimal) *PlaceOrderService {
	s.body["slOrdPx"] = px.String()
	return s
}

// SetTpTriggerPxType sets the take-profit trigger price type (last/index/mark).
func (s *PlaceOrderService) SetTpTriggerPxType(pxType string) *PlaceOrderService {
	s.body["tpTriggerPxType"] = pxType
	return s
}

// SetSlTriggerPxType sets the stop-loss trigger price type (last/index/mark).
func (s *PlaceOrderService) SetSlTriggerPxType(pxType string) *PlaceOrderService {
	s.body["slTriggerPxType"] = pxType
	return s
}

// SetQuickMgnType sets the quick-margin type (manual/auto_borrow/auto_borrow_repay).
func (s *PlaceOrderService) SetQuickMgnType(quickMgnType string) *PlaceOrderService {
	s.body["quickMgnType"] = quickMgnType
	return s
}

// SetStpId sets the self-trade-prevention id.
func (s *PlaceOrderService) SetStpId(stpId string) *PlaceOrderService {
	s.body["stpId"] = stpId
	return s
}

// SetStpMode sets the self-trade-prevention mode (cancel_maker/cancel_taker/cancel_both).
func (s *PlaceOrderService) SetStpMode(stpMode string) *PlaceOrderService {
	s.body["stpMode"] = stpMode
	return s
}

// SetAttachAlgoOrds attaches take-profit / stop-loss algo orders.
func (s *PlaceOrderService) SetAttachAlgoOrds(orders []AttachAlgoOrd) *PlaceOrderService {
	s.body["attachAlgoOrds"] = orders
	return s
}

func (s *PlaceOrderService) Do(ctx context.Context) (*OrderResult, error) {
	req := request.Post(ctx, s.c, "/api/v5/trade/order", s.body).WithSign()
	list, err := request.DoListPartial[OrderResult](req)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return &list[0], nil
}

// BatchOrdersService -- POST /api/v5/trade/batch-orders (Trade)
//
// Places up to 20 orders in a single request (array body). Each result carries
// its own sCode/sMsg.
type BatchOrdersService struct {
	c     *Client
	items []OrderArg
}

func (c *Client) NewBatchOrdersService(orders []OrderArg) *BatchOrdersService {
	return &BatchOrdersService{c: c, items: orders}
}

func (s *BatchOrdersService) Do(ctx context.Context) ([]OrderResult, error) {
	req := request.Post(ctx, s.c, "/api/v5/trade/batch-orders").SetBody(s.items).WithSign()
	return request.DoListPartial[OrderResult](req)
}

// CancelOrderService -- POST /api/v5/trade/cancel-order (Trade)
//
// Cancels a single pending order by ordId or clOrdId.
type CancelOrderService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewCancelOrderService(instId string) *CancelOrderService {
	return &CancelOrderService{c: c, body: map[string]any{"instId": instId}}
}

// SetOrdId sets the order id (either ordId or clOrdId is required).
func (s *CancelOrderService) SetOrdId(ordId string) *CancelOrderService {
	s.body["ordId"] = ordId
	return s
}

// SetClOrdId sets the client-supplied order id (either ordId or clOrdId is required).
func (s *CancelOrderService) SetClOrdId(clOrdId string) *CancelOrderService {
	s.body["clOrdId"] = clOrdId
	return s
}

func (s *CancelOrderService) Do(ctx context.Context) (*OrderResult, error) {
	req := request.Post(ctx, s.c, "/api/v5/trade/cancel-order", s.body).WithSign()
	list, err := request.DoListPartial[OrderResult](req)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return &list[0], nil
}

// CancelOrderArg is one leg of a cancel-batch-orders request body.
type CancelOrderArg struct {
	InstrumentID  string `json:"instId"`
	OrderID       string `json:"ordId,omitempty"`
	ClientOrderID string `json:"clOrdId,omitempty"`
}

// CancelBatchOrdersService -- POST /api/v5/trade/cancel-batch-orders (Trade)
//
// Cancels up to 20 pending orders in a single request (array body).
type CancelBatchOrdersService struct {
	c     *Client
	items []CancelOrderArg
}

func (c *Client) NewCancelBatchOrdersService(orders []CancelOrderArg) *CancelBatchOrdersService {
	return &CancelBatchOrdersService{c: c, items: orders}
}

func (s *CancelBatchOrdersService) Do(ctx context.Context) ([]OrderResult, error) {
	req := request.Post(ctx, s.c, "/api/v5/trade/cancel-batch-orders").SetBody(s.items).WithSign()
	return request.DoListPartial[OrderResult](req)
}

// AmendOrderService -- POST /api/v5/trade/amend-order (Trade)
//
// Amends a single pending order's size and/or price by ordId or clOrdId.
type AmendOrderService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewAmendOrderService(instId string) *AmendOrderService {
	return &AmendOrderService{c: c, body: map[string]any{"instId": instId}}
}

// SetOrdId sets the order id (either ordId or clOrdId is required).
func (s *AmendOrderService) SetOrdId(ordId string) *AmendOrderService {
	s.body["ordId"] = ordId
	return s
}

// SetClOrdId sets the client-supplied order id (either ordId or clOrdId is required).
func (s *AmendOrderService) SetClOrdId(clOrdId string) *AmendOrderService {
	s.body["clOrdId"] = clOrdId
	return s
}

// SetCxlOnFail cancels the order if the amendment fails.
func (s *AmendOrderService) SetCxlOnFail(cxlOnFail bool) *AmendOrderService {
	s.body["cxlOnFail"] = cxlOnFail
	return s
}

// SetReqId sets the client-supplied request id for the amendment.
func (s *AmendOrderService) SetReqId(reqId string) *AmendOrderService {
	s.body["reqId"] = reqId
	return s
}

// SetNewSz sets the new order size.
func (s *AmendOrderService) SetNewSz(newSz decimal.Decimal) *AmendOrderService {
	s.body["newSz"] = newSz.String()
	return s
}

// SetNewPx sets the new order price.
func (s *AmendOrderService) SetNewPx(newPx decimal.Decimal) *AmendOrderService {
	s.body["newPx"] = newPx.String()
	return s
}

// SetNewPxUsd sets the new options price in USD (OPTION only).
func (s *AmendOrderService) SetNewPxUsd(newPxUsd decimal.Decimal) *AmendOrderService {
	s.body["newPxUsd"] = newPxUsd.String()
	return s
}

// SetNewPxVol sets the new implied volatility (OPTION only).
func (s *AmendOrderService) SetNewPxVol(newPxVol decimal.Decimal) *AmendOrderService {
	s.body["newPxVol"] = newPxVol.String()
	return s
}

// SetNewTpTriggerPx sets the new take-profit trigger price.
func (s *AmendOrderService) SetNewTpTriggerPx(px decimal.Decimal) *AmendOrderService {
	s.body["newTpTriggerPx"] = px.String()
	return s
}

// SetNewTpOrdPx sets the new take-profit order price.
func (s *AmendOrderService) SetNewTpOrdPx(px decimal.Decimal) *AmendOrderService {
	s.body["newTpOrdPx"] = px.String()
	return s
}

// SetNewSlTriggerPx sets the new stop-loss trigger price.
func (s *AmendOrderService) SetNewSlTriggerPx(px decimal.Decimal) *AmendOrderService {
	s.body["newSlTriggerPx"] = px.String()
	return s
}

// SetNewSlOrdPx sets the new stop-loss order price.
func (s *AmendOrderService) SetNewSlOrdPx(px decimal.Decimal) *AmendOrderService {
	s.body["newSlOrdPx"] = px.String()
	return s
}

// SetNewTpTriggerPxType sets the new take-profit trigger price type.
func (s *AmendOrderService) SetNewTpTriggerPxType(pxType string) *AmendOrderService {
	s.body["newTpTriggerPxType"] = pxType
	return s
}

// SetNewSlTriggerPxType sets the new stop-loss trigger price type.
func (s *AmendOrderService) SetNewSlTriggerPxType(pxType string) *AmendOrderService {
	s.body["newSlTriggerPxType"] = pxType
	return s
}

// SetAttachAlgoOrds amends the attached take-profit / stop-loss algo orders.
func (s *AmendOrderService) SetAttachAlgoOrds(orders []AttachAlgoOrd) *AmendOrderService {
	s.body["attachAlgoOrds"] = orders
	return s
}

func (s *AmendOrderService) Do(ctx context.Context) (*AmendResult, error) {
	req := request.Post(ctx, s.c, "/api/v5/trade/amend-order", s.body).WithSign()
	list, err := request.DoListPartial[AmendResult](req)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return &list[0], nil
}

// AmendOrderArg is one leg of an amend-batch-orders request body.
type AmendOrderArg struct {
	InstrumentID                  string          `json:"instId"`
	OrderID                       string          `json:"ordId,omitempty"`
	ClientOrderID                 string          `json:"clOrdId,omitempty"`
	CancelOnFail                  bool            `json:"cxlOnFail,omitempty"`
	RequestID                     string          `json:"reqId,omitempty"`
	NewSize                       decimal.Decimal `json:"newSz,omitzero"`
	NewPrice                      decimal.Decimal `json:"newPx,omitzero"`
	NewPriceUSD                   decimal.Decimal `json:"newPxUsd,omitzero"`
	NewPriceVolatility            decimal.Decimal `json:"newPxVol,omitzero"`
	NewTakeProfitTriggerPrice     decimal.Decimal `json:"newTpTriggerPx,omitzero"`
	NewTakeProfitOrderPrice       decimal.Decimal `json:"newTpOrdPx,omitzero"`
	NewStopLossTriggerPrice       decimal.Decimal `json:"newSlTriggerPx,omitzero"`
	NewStopLossOrderPrice         decimal.Decimal `json:"newSlOrdPx,omitzero"`
	NewTakeProfitTriggerPriceType string          `json:"newTpTriggerPxType,omitempty"`
	NewStopLossTriggerPriceType   string          `json:"newSlTriggerPxType,omitempty"`
	AttachAlgoOrders              []AttachAlgoOrd `json:"attachAlgoOrds,omitempty"`
}

// AmendBatchOrdersService -- POST /api/v5/trade/amend-batch-orders (Trade)
//
// Amends up to 20 pending orders in a single request (array body).
type AmendBatchOrdersService struct {
	c     *Client
	items []AmendOrderArg
}

func (c *Client) NewAmendBatchOrdersService(orders []AmendOrderArg) *AmendBatchOrdersService {
	return &AmendBatchOrdersService{c: c, items: orders}
}

func (s *AmendBatchOrdersService) Do(ctx context.Context) ([]AmendResult, error) {
	req := request.Post(ctx, s.c, "/api/v5/trade/amend-batch-orders").SetBody(s.items).WithSign()
	return request.DoListPartial[AmendResult](req)
}

// ClosePositionService -- POST /api/v5/trade/close-position (Trade)
//
// Market-closes the position of an instrument under the given margin mode.
type ClosePositionService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewClosePositionService(instId string, mgnMode MgnMode) *ClosePositionService {
	return &ClosePositionService{c: c, body: map[string]any{
		"instId":  instId,
		"mgnMode": string(mgnMode),
	}}
}

// SetPosSide sets the position side (long/short; required in long/short mode).
func (s *ClosePositionService) SetPosSide(posSide PosSide) *ClosePositionService {
	s.body["posSide"] = string(posSide)
	return s
}

// SetCcy sets the margin currency (single-currency margin cross only).
func (s *ClosePositionService) SetCcy(ccy string) *ClosePositionService {
	s.body["ccy"] = ccy
	return s
}

// SetAutoCxl cancels pending orders that would prevent the close.
func (s *ClosePositionService) SetAutoCxl(autoCxl bool) *ClosePositionService {
	s.body["autoCxl"] = autoCxl
	return s
}

// SetClOrdId sets the client-supplied order id.
func (s *ClosePositionService) SetClOrdId(clOrdId string) *ClosePositionService {
	s.body["clOrdId"] = clOrdId
	return s
}

// SetTag sets the order tag.
func (s *ClosePositionService) SetTag(tag string) *ClosePositionService {
	s.body["tag"] = tag
	return s
}

func (s *ClosePositionService) Do(ctx context.Context) (*ClosePositionResult, error) {
	req := request.Post(ctx, s.c, "/api/v5/trade/close-position", s.body).WithSign()
	return request.DoOne[ClosePositionResult](req)
}

// GetOrderService -- GET /api/v5/trade/order (Read)
//
// Returns the details of a single order by ordId or clOrdId.
type GetOrderService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetOrderService(instId string) *GetOrderService {
	return &GetOrderService{c: c, params: map[string]string{"instId": instId}}
}

// SetOrdId sets the order id (either ordId or clOrdId is required).
func (s *GetOrderService) SetOrdId(ordId string) *GetOrderService {
	s.params["ordId"] = ordId
	return s
}

// SetClOrdId sets the client-supplied order id (either ordId or clOrdId is required).
func (s *GetOrderService) SetClOrdId(clOrdId string) *GetOrderService {
	s.params["clOrdId"] = clOrdId
	return s
}

func (s *GetOrderService) Do(ctx context.Context) (*Order, error) {
	req := request.Get(ctx, s.c, "/api/v5/trade/order", s.params).WithSign()
	return request.DoOne[Order](req)
}

// GetOrdersPendingService -- GET /api/v5/trade/orders-pending (Read)
//
// Returns the account's currently live (incomplete) orders.
type GetOrdersPendingService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetOrdersPendingService() *GetOrdersPendingService {
	return &GetOrdersPendingService{c: c, params: map[string]string{}}
}

// SetInstType filters by product line (SPOT/MARGIN/SWAP/FUTURES/OPTION).
func (s *GetOrdersPendingService) SetInstType(instType InstType) *GetOrdersPendingService {
	s.params["instType"] = string(instType)
	return s
}

// SetUly filters by underlying.
func (s *GetOrdersPendingService) SetUly(uly string) *GetOrdersPendingService {
	s.params["uly"] = uly
	return s
}

// SetInstFamily filters by instrument family.
func (s *GetOrdersPendingService) SetInstFamily(instFamily string) *GetOrdersPendingService {
	s.params["instFamily"] = instFamily
	return s
}

// SetInstId filters by instrument id.
func (s *GetOrdersPendingService) SetInstId(instId string) *GetOrdersPendingService {
	s.params["instId"] = instId
	return s
}

// SetOrdType filters by order type (market/limit/post_only/fok/ioc/...).
func (s *GetOrdersPendingService) SetOrdType(ordType OrdType) *GetOrdersPendingService {
	s.params["ordType"] = string(ordType)
	return s
}

// SetState filters by order state (live/partially_filled).
func (s *GetOrdersPendingService) SetState(state OrdState) *GetOrdersPendingService {
	s.params["state"] = string(state)
	return s
}

// SetAfter pages to orders with an ordId earlier than this value.
func (s *GetOrdersPendingService) SetAfter(ordId string) *GetOrdersPendingService {
	s.params["after"] = ordId
	return s
}

// SetBefore pages to orders with an ordId later than this value.
func (s *GetOrdersPendingService) SetBefore(ordId string) *GetOrdersPendingService {
	s.params["before"] = ordId
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetOrdersPendingService) SetLimit(limit int) *GetOrdersPendingService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetOrdersPendingService) Do(ctx context.Context) ([]Order, error) {
	req := request.Get(ctx, s.c, "/api/v5/trade/orders-pending", s.params).WithSign()
	return request.DoList[Order](req)
}

// GetOrdersHistoryService -- GET /api/v5/trade/orders-history (Read)
//
// Returns the account's completed orders from the last seven days.
type GetOrdersHistoryService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetOrdersHistoryService(instType InstType) *GetOrdersHistoryService {
	return &GetOrdersHistoryService{c: c, params: map[string]string{"instType": string(instType)}}
}

// SetUly filters by underlying.
func (s *GetOrdersHistoryService) SetUly(uly string) *GetOrdersHistoryService {
	s.params["uly"] = uly
	return s
}

// SetInstFamily filters by instrument family.
func (s *GetOrdersHistoryService) SetInstFamily(instFamily string) *GetOrdersHistoryService {
	s.params["instFamily"] = instFamily
	return s
}

// SetInstId filters by instrument id.
func (s *GetOrdersHistoryService) SetInstId(instId string) *GetOrdersHistoryService {
	s.params["instId"] = instId
	return s
}

// SetOrdType filters by order type.
func (s *GetOrdersHistoryService) SetOrdType(ordType OrdType) *GetOrdersHistoryService {
	s.params["ordType"] = string(ordType)
	return s
}

// SetState filters by order state (canceled/filled/mmp_canceled).
func (s *GetOrdersHistoryService) SetState(state OrdState) *GetOrdersHistoryService {
	s.params["state"] = string(state)
	return s
}

// SetCategory filters by order category (twap/adl/full_liquidation/...).
func (s *GetOrdersHistoryService) SetCategory(category string) *GetOrdersHistoryService {
	s.params["category"] = category
	return s
}

// SetAfter pages to orders with an ordId earlier than this value.
func (s *GetOrdersHistoryService) SetAfter(ordId string) *GetOrdersHistoryService {
	s.params["after"] = ordId
	return s
}

// SetBefore pages to orders with an ordId later than this value.
func (s *GetOrdersHistoryService) SetBefore(ordId string) *GetOrdersHistoryService {
	s.params["before"] = ordId
	return s
}

// SetBegin filters to orders created at or after the given time.
func (s *GetOrdersHistoryService) SetBegin(t time.Time) *GetOrdersHistoryService {
	s.params["begin"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEnd filters to orders created at or before the given time.
func (s *GetOrdersHistoryService) SetEnd(t time.Time) *GetOrdersHistoryService {
	s.params["end"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetOrdersHistoryService) SetLimit(limit int) *GetOrdersHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetOrdersHistoryService) Do(ctx context.Context) ([]Order, error) {
	req := request.Get(ctx, s.c, "/api/v5/trade/orders-history", s.params).WithSign()
	return request.DoList[Order](req)
}

// GetOrdersHistoryArchiveService -- GET /api/v5/trade/orders-history-archive (Read)
//
// Returns the account's completed orders from the last three months.
type GetOrdersHistoryArchiveService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetOrdersHistoryArchiveService(instType InstType) *GetOrdersHistoryArchiveService {
	return &GetOrdersHistoryArchiveService{c: c, params: map[string]string{"instType": string(instType)}}
}

// SetUly filters by underlying.
func (s *GetOrdersHistoryArchiveService) SetUly(uly string) *GetOrdersHistoryArchiveService {
	s.params["uly"] = uly
	return s
}

// SetInstFamily filters by instrument family.
func (s *GetOrdersHistoryArchiveService) SetInstFamily(instFamily string) *GetOrdersHistoryArchiveService {
	s.params["instFamily"] = instFamily
	return s
}

// SetInstId filters by instrument id.
func (s *GetOrdersHistoryArchiveService) SetInstId(instId string) *GetOrdersHistoryArchiveService {
	s.params["instId"] = instId
	return s
}

// SetOrdType filters by order type.
func (s *GetOrdersHistoryArchiveService) SetOrdType(ordType OrdType) *GetOrdersHistoryArchiveService {
	s.params["ordType"] = string(ordType)
	return s
}

// SetState filters by order state (canceled/filled/mmp_canceled).
func (s *GetOrdersHistoryArchiveService) SetState(state OrdState) *GetOrdersHistoryArchiveService {
	s.params["state"] = string(state)
	return s
}

// SetCategory filters by order category.
func (s *GetOrdersHistoryArchiveService) SetCategory(category string) *GetOrdersHistoryArchiveService {
	s.params["category"] = category
	return s
}

// SetAfter pages to orders with an ordId earlier than this value.
func (s *GetOrdersHistoryArchiveService) SetAfter(ordId string) *GetOrdersHistoryArchiveService {
	s.params["after"] = ordId
	return s
}

// SetBefore pages to orders with an ordId later than this value.
func (s *GetOrdersHistoryArchiveService) SetBefore(ordId string) *GetOrdersHistoryArchiveService {
	s.params["before"] = ordId
	return s
}

// SetBegin filters to orders created at or after the given time.
func (s *GetOrdersHistoryArchiveService) SetBegin(t time.Time) *GetOrdersHistoryArchiveService {
	s.params["begin"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEnd filters to orders created at or before the given time.
func (s *GetOrdersHistoryArchiveService) SetEnd(t time.Time) *GetOrdersHistoryArchiveService {
	s.params["end"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetOrdersHistoryArchiveService) SetLimit(limit int) *GetOrdersHistoryArchiveService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetOrdersHistoryArchiveService) Do(ctx context.Context) ([]Order, error) {
	req := request.Get(ctx, s.c, "/api/v5/trade/orders-history-archive", s.params).WithSign()
	return request.DoList[Order](req)
}

// Order is a single order record returned by the order / orders-pending /
// orders-history(-archive) endpoints. The validating account had no order
// history, so the field set is modeled from the OKX doc field table.
type Order struct {
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
	Pnl                        decimal.Decimal `json:"pnl"`
	OrderType                  OrdType         `json:"ordType"`
	Side                       Side            `json:"side"`
	PositionSide               PosSide         `json:"posSide"`
	TradeMode                  TdMode          `json:"tdMode"`
	AccumulatedFillSize        decimal.Decimal `json:"accFillSz"`
	FillPrice                  decimal.Decimal `json:"fillPx"`
	TradeID                    string          `json:"tradeId"`
	FillSize                   decimal.Decimal `json:"fillSz"`
	FillTime                   time.Time       `json:"fillTime"`
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
	Source                     string          `json:"source"`
	Category                   string          `json:"category"`
	ReduceOnly                 string          `json:"reduceOnly"` // OKX sends a quoted "true"/"false"
	CancelSource               string          `json:"cancelSource"`
	CancelSourceReason         string          `json:"cancelSourceReason"`
	QuickMarginType            string          `json:"quickMgnType"`
	AlgoClientOrderID          string          `json:"algoClOrdId"`
	AlgoID                     string          `json:"algoId"`
	IsTakeProfitLimit          string          `json:"isTpLimit"`
	Outcome                    string          `json:"outcome"`
	LinkedAlgoOrder            OrderLinkedAlgo `json:"linkedAlgoOrd"`
	UpdateTime                 time.Time       `json:"uTime"`
	CreationTime               time.Time       `json:"cTime"`
	TradeQuoteCurrency         string          `json:"tradeQuoteCcy"`
}

// OrderLinkedAlgo is the algo order linked to a regular order (returned by
// order-info as a nested object).
type OrderLinkedAlgo struct {
	AlgoID string `json:"algoId"`
}

// GetAccountRateLimitService -- GET /api/v5/trade/account-rate-limit (Read)
//
// Returns the account's order-placement rate limit and current fill ratio.
type GetAccountRateLimitService struct {
	c *Client
}

func (c *Client) NewGetAccountRateLimitService() *GetAccountRateLimitService {
	return &GetAccountRateLimitService{c: c}
}

func (s *GetAccountRateLimitService) Do(ctx context.Context) (*AccountRateLimit, error) {
	req := request.Get(ctx, s.c, "/api/v5/trade/account-rate-limit").WithSign()
	return request.DoOne[AccountRateLimit](req)
}

// AccountRateLimit is the account's order rate-limit snapshot.
type AccountRateLimit struct {
	AccountRateLimit     decimal.Decimal `json:"accRateLimit"`
	FillRatio            decimal.Decimal `json:"fillRatio"`
	MainFillRatio        decimal.Decimal `json:"mainFillRatio"`
	NextAccountRateLimit decimal.Decimal `json:"nextAccRateLimit"`
	Timestamp            time.Time       `json:"ts"`
}

// OrderPrecheckService -- POST /api/v5/trade/order-precheck (Trade)
//
// Pre-checks the account-level impact (estimated margin, available balance,
// liquidation price) of placing an order without actually placing it.
type OrderPrecheckService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewOrderPrecheckService(instId string, tdMode TdMode, side Side, ordType OrdType, sz decimal.Decimal) *OrderPrecheckService {
	return &OrderPrecheckService{c: c, body: map[string]any{
		"instId":  instId,
		"tdMode":  string(tdMode),
		"side":    string(side),
		"ordType": string(ordType),
		"sz":      sz.String(),
	}}
}

// SetPosSide sets the position side (long/short/net).
func (s *OrderPrecheckService) SetPosSide(posSide PosSide) *OrderPrecheckService {
	s.body["posSide"] = string(posSide)
	return s
}

// SetPx sets the order price (required for limit-type orders).
func (s *OrderPrecheckService) SetPx(px decimal.Decimal) *OrderPrecheckService {
	s.body["px"] = px.String()
	return s
}

// SetReduceOnly toggles reduce-only.
func (s *OrderPrecheckService) SetReduceOnly(reduceOnly bool) *OrderPrecheckService {
	s.body["reduceOnly"] = reduceOnly
	return s
}

// SetTgtCcy selects the unit of a spot market order's size (base_ccy/quote_ccy).
func (s *OrderPrecheckService) SetTgtCcy(tgtCcy TgtCcy) *OrderPrecheckService {
	s.body["tgtCcy"] = string(tgtCcy)
	return s
}

// SetAttachAlgoOrds attaches take-profit / stop-loss algo orders for the check.
func (s *OrderPrecheckService) SetAttachAlgoOrds(orders []AttachAlgoOrd) *OrderPrecheckService {
	s.body["attachAlgoOrds"] = orders
	return s
}

func (s *OrderPrecheckService) Do(ctx context.Context) (*OrderPrecheck, error) {
	req := request.Post(ctx, s.c, "/api/v5/trade/order-precheck", s.body).WithSign()
	return request.DoOne[OrderPrecheck](req)
}

// OrderPrecheck is the estimated account impact of a hypothetical order.
type OrderPrecheck struct {
	AdjustedEquity            decimal.Decimal `json:"adjEq"`
	AdjustedEquityChange      decimal.Decimal `json:"adjEqChg"`
	AvailableBalance          decimal.Decimal `json:"availBal"`
	AvailableBalanceChange    decimal.Decimal `json:"availBalChg"`
	IMR                       decimal.Decimal `json:"imr"`
	IMRChange                 decimal.Decimal `json:"imrChg"`
	Borrowed                  decimal.Decimal `json:"borrowed"`
	Liability                 decimal.Decimal `json:"liab"`
	LiabilityChange           decimal.Decimal `json:"liabChg"`
	LiabilityChangeCurrency   string          `json:"liabChgCcy"`
	LiquidationPrice          decimal.Decimal `json:"liqPx"`
	LiquidationPriceDiff      decimal.Decimal `json:"liqPxDiff"`
	LiquidationPriceDiffRatio decimal.Decimal `json:"liqPxDiffRatio"`
	Margin                    decimal.Decimal `json:"mgn"`
	MarginRatio               decimal.Decimal `json:"mgnRatio"`
	MarginRatioChange         decimal.Decimal `json:"mgnRatioChg"`
	MMR                       decimal.Decimal `json:"mmr"`
	MMRChange                 decimal.Decimal `json:"mmrChg"`
	PositionBalance           decimal.Decimal `json:"posBal"`
	PositionBalanceChange     decimal.Decimal `json:"posBalChg"`
	Type                      string          `json:"type"`
}

// MassCancelService -- POST /api/v5/trade/mass-cancel (Trade)
//
// Cancels all pending MMP orders for an instrument family (market-maker
// protection).
type MassCancelService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewMassCancelService(instType InstType, instFamily string) *MassCancelService {
	return &MassCancelService{c: c, body: map[string]any{
		"instType":   string(instType),
		"instFamily": instFamily,
	}}
}

// SetLockInterval sets the auto-cancel lock interval in milliseconds (0-10000).
func (s *MassCancelService) SetLockInterval(ms int) *MassCancelService {
	s.body["lockInterval"] = strconv.Itoa(ms)
	return s
}

func (s *MassCancelService) Do(ctx context.Context) (*MassCancelResult, error) {
	req := request.Post(ctx, s.c, "/api/v5/trade/mass-cancel", s.body).WithSign()
	return request.DoOne[MassCancelResult](req)
}

// MassCancelResult is the ack returned by the mass-cancel endpoint.
type MassCancelResult struct {
	Result bool `json:"result"`
}

// CancelAllAfterService -- POST /api/v5/trade/cancel-all-after (Trade)
//
// Arms (or disarms with timeOut 0) a dead-man's-switch that cancels all pending
// orders after the given number of seconds if not refreshed.
type CancelAllAfterService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewCancelAllAfterService(timeOut int) *CancelAllAfterService {
	return &CancelAllAfterService{c: c, body: map[string]any{"timeOut": strconv.Itoa(timeOut)}}
}

// SetTag sets the order tag scope for the timer.
func (s *CancelAllAfterService) SetTag(tag string) *CancelAllAfterService {
	s.body["tag"] = tag
	return s
}

func (s *CancelAllAfterService) Do(ctx context.Context) (*CancelAllAfterResult, error) {
	req := request.Post(ctx, s.c, "/api/v5/trade/cancel-all-after", s.body).WithSign()
	return request.DoOne[CancelAllAfterResult](req)
}

// CancelAllAfterResult is the ack returned by the cancel-all-after endpoint.
type CancelAllAfterResult struct {
	TriggerTime time.Time `json:"triggerTime"`
	Tag         string    `json:"tag"`
	Timestamp   time.Time `json:"ts"`
}
