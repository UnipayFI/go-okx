package okx

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-okx/request"
	"github.com/shopspring/decimal"
)

// Recurring buy (DCA) endpoints live under /api/v5/tradingBot/recurring/*.
//
// NOTE: although these are conceptually "trade" endpoints, OKX serves the
// recurring-buy strategy bots under the tradingBot path (the /api/v5/trade/
// prefix returns 404). All of them are private and therefore always signed.

// RecurringPeriod is the cadence at which a recurring-buy strategy invests.
type RecurringPeriod string

const (
	RecurringPeriodMonthly RecurringPeriod = "monthly"
	RecurringPeriodWeekly  RecurringPeriod = "weekly"
	RecurringPeriodDaily   RecurringPeriod = "daily"
	RecurringPeriodHourly  RecurringPeriod = "hourly"
)

// RecurringState is the lifecycle state of a recurring-buy strategy bot.
type RecurringState string

const (
	RecurringStateStarting   RecurringState = "starting"
	RecurringStateRunning    RecurringState = "running"
	RecurringStateStopping   RecurringState = "stopping"
	RecurringStateNotStarted RecurringState = "not_started"
	RecurringStateStopped    RecurringState = "stopped"
	RecurringStatePause      RecurringState = "pause"
)

// RecurringItem is one currency allocation within a recurring-buy strategy
// (used in both the create body and the strategy detail response).
type RecurringItem struct {
	Currency string          `json:"ccy"`
	Ratio    decimal.Decimal `json:"ratio"`
}

// PlaceRecurringOrderService -- POST /api/v5/tradingBot/recurring/order-algo (Trade)
//
// Creates a recurring-buy (DCA) strategy that invests amt of investmentCcy
// across recurringList on the configured period/day/time.
type PlaceRecurringOrderService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewPlaceRecurringOrderService(stgyName string, recurringList []RecurringItem, period RecurringPeriod, recurringTime string, timeZone string, amt decimal.Decimal, investmentCcy string, tdMode TdMode) *PlaceRecurringOrderService {
	return &PlaceRecurringOrderService{c: c, body: map[string]any{
		"stgyName":      stgyName,
		"recurringList": recurringList,
		"period":        string(period),
		"recurringTime": recurringTime,
		"timeZone":      timeZone,
		"amt":           amt.String(),
		"investmentCcy": investmentCcy,
		"tdMode":        string(tdMode),
	}}
}

// SetRecurringDay sets the day of the recurring buy. Required when period is
// "weekly" (1-7) or "monthly" (1-28).
func (s *PlaceRecurringOrderService) SetRecurringDay(recurringDay string) *PlaceRecurringOrderService {
	s.body["recurringDay"] = recurringDay
	return s
}

// SetRecurringHour sets the hour of the recurring buy ("0"/"6"/"12"/"18"),
// required when period is "hourly".
func (s *PlaceRecurringOrderService) SetRecurringHour(recurringHour string) *PlaceRecurringOrderService {
	s.body["recurringHour"] = recurringHour
	return s
}

// SetAlgoClOrdId sets a client-supplied strategy id.
func (s *PlaceRecurringOrderService) SetAlgoClOrdId(algoClOrdId string) *PlaceRecurringOrderService {
	s.body["algoClOrdId"] = algoClOrdId
	return s
}

// SetTag sets an order tag for the strategy.
func (s *PlaceRecurringOrderService) SetTag(tag string) *PlaceRecurringOrderService {
	s.body["tag"] = tag
	return s
}

func (s *PlaceRecurringOrderService) Do(ctx context.Context) (*RecurringResult, error) {
	req := request.Post(ctx, s.c, "/api/v5/tradingBot/recurring/order-algo", s.body).WithSign()
	list, err := request.DoListPartial[RecurringResult](req)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return &list[0], nil
}

// RecurringResult is the per-item ack of a recurring-buy create/stop op. OKX may
// set the top-level code to "1" with the real reason in sCode/sMsg.
type RecurringResult struct {
	AlgoID            string `json:"algoId"`
	AlgoClientOrderID string `json:"algoClOrdId"`
	SCode             string `json:"sCode"`
	SMsg              string `json:"sMsg"`
	Tag               string `json:"tag"`
}

// AmendRecurringOrderService -- POST /api/v5/tradingBot/recurring/amend-order-algo (Trade)
//
// Amends the display name of a running recurring-buy strategy. Implement-only;
// not exercised against the live account.
type AmendRecurringOrderService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewAmendRecurringOrderService(algoId string) *AmendRecurringOrderService {
	return &AmendRecurringOrderService{c: c, body: map[string]any{
		"algoId": algoId,
	}}
}

// SetStgyName sets the new display name for the strategy.
func (s *AmendRecurringOrderService) SetStgyName(stgyName string) *AmendRecurringOrderService {
	s.body["stgyName"] = stgyName
	return s
}

func (s *AmendRecurringOrderService) Do(ctx context.Context) (*RecurringResult, error) {
	req := request.Post(ctx, s.c, "/api/v5/tradingBot/recurring/amend-order-algo", s.body).WithSign()
	list, err := request.DoListPartial[RecurringResult](req)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return &list[0], nil
}

// StopRecurringOrderService -- POST /api/v5/tradingBot/recurring/stop-order-algo (Trade)
//
// Stops one or more recurring-buy strategies. The request body is an ARRAY of
// {algoId} objects (max 10).
type StopRecurringOrderService struct {
	c     *Client
	items []map[string]any
}

func (c *Client) NewStopRecurringOrderService(algoIds []string) *StopRecurringOrderService {
	items := make([]map[string]any, 0, len(algoIds))
	for _, id := range algoIds {
		items = append(items, map[string]any{"algoId": id})
	}
	return &StopRecurringOrderService{c: c, items: items}
}

func (s *StopRecurringOrderService) Do(ctx context.Context) ([]RecurringResult, error) {
	req := request.Post(ctx, s.c, "/api/v5/tradingBot/recurring/stop-order-algo").SetBody(s.items).WithSign()
	return request.DoListPartial[RecurringResult](req)
}

// GetRecurringOrdersPendingService -- GET /api/v5/tradingBot/recurring/orders-algo-pending (Read)
//
// Returns the account's active (not-yet-stopped) recurring-buy strategies.
type GetRecurringOrdersPendingService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetRecurringOrdersPendingService() *GetRecurringOrdersPendingService {
	return &GetRecurringOrdersPendingService{c: c, params: map[string]string{}}
}

// SetAlgoId filters to a single strategy id.
func (s *GetRecurringOrdersPendingService) SetAlgoId(algoId string) *GetRecurringOrdersPendingService {
	s.params["algoId"] = algoId
	return s
}

// SetAfter pages backward from the given algoId (records older than it).
func (s *GetRecurringOrdersPendingService) SetAfter(after string) *GetRecurringOrdersPendingService {
	s.params["after"] = after
	return s
}

// SetBefore pages forward from the given algoId (records newer than it).
func (s *GetRecurringOrdersPendingService) SetBefore(before string) *GetRecurringOrdersPendingService {
	s.params["before"] = before
	return s
}

// SetLimit caps the number of returned strategies (max 100, default 100).
func (s *GetRecurringOrdersPendingService) SetLimit(limit int) *GetRecurringOrdersPendingService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetRecurringOrdersPendingService) Do(ctx context.Context) ([]RecurringOrder, error) {
	req := request.Get(ctx, s.c, "/api/v5/tradingBot/recurring/orders-algo-pending", s.params).WithSign()
	return request.DoList[RecurringOrder](req)
}

// GetRecurringOrdersHistoryService -- GET /api/v5/tradingBot/recurring/orders-algo-history (Read)
//
// Returns the account's stopped recurring-buy strategies.
type GetRecurringOrdersHistoryService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetRecurringOrdersHistoryService() *GetRecurringOrdersHistoryService {
	return &GetRecurringOrdersHistoryService{c: c, params: map[string]string{}}
}

// SetAlgoId filters to a single strategy id.
func (s *GetRecurringOrdersHistoryService) SetAlgoId(algoId string) *GetRecurringOrdersHistoryService {
	s.params["algoId"] = algoId
	return s
}

// SetAfter pages backward from the given algoId (records older than it).
func (s *GetRecurringOrdersHistoryService) SetAfter(after string) *GetRecurringOrdersHistoryService {
	s.params["after"] = after
	return s
}

// SetBefore pages forward from the given algoId (records newer than it).
func (s *GetRecurringOrdersHistoryService) SetBefore(before string) *GetRecurringOrdersHistoryService {
	s.params["before"] = before
	return s
}

// SetLimit caps the number of returned strategies (max 100, default 100).
func (s *GetRecurringOrdersHistoryService) SetLimit(limit int) *GetRecurringOrdersHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetRecurringOrdersHistoryService) Do(ctx context.Context) ([]RecurringOrder, error) {
	req := request.Get(ctx, s.c, "/api/v5/tradingBot/recurring/orders-algo-history", s.params).WithSign()
	return request.DoList[RecurringOrder](req)
}

// GetRecurringOrderDetailsService -- GET /api/v5/tradingBot/recurring/orders-algo-details (Read)
//
// Returns the full detail of a single recurring-buy strategy.
type GetRecurringOrderDetailsService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetRecurringOrderDetailsService(algoId string) *GetRecurringOrderDetailsService {
	return &GetRecurringOrderDetailsService{c: c, params: map[string]string{"algoId": algoId}}
}

func (s *GetRecurringOrderDetailsService) Do(ctx context.Context) (*RecurringOrder, error) {
	req := request.Get(ctx, s.c, "/api/v5/tradingBot/recurring/orders-algo-details", s.params).WithSign()
	return request.DoOne[RecurringOrder](req)
}

// RecurringOrder is a recurring-buy (DCA) strategy bot.
type RecurringOrder struct {
	AlgoID               string          `json:"algoId"`
	AlgoClientOrderID    string          `json:"algoClOrdId"`
	InstrumentType       InstType        `json:"instType"`
	Cycles               string          `json:"cycles"`
	StrategyName         string          `json:"stgyName"`
	State                RecurringState  `json:"state"`
	Period               RecurringPeriod `json:"period"`
	RecurringDay         string          `json:"recurringDay"`
	RecurringHour        string          `json:"recurringHour"`
	RecurringTime        string          `json:"recurringTime"`
	TimeZone             string          `json:"timeZone"`
	Amount               decimal.Decimal `json:"amt"`
	InvestmentCurrency   string          `json:"investmentCcy"`
	InvestmentAmount     decimal.Decimal `json:"investmentAmt"`
	TotalPnl             decimal.Decimal `json:"totalPnl"`
	TotalAnnualRate      decimal.Decimal `json:"totalAnnRate"`
	PnlRatio             decimal.Decimal `json:"pnlRatio"`
	MarketCapitalization decimal.Decimal `json:"mktCap"`
	TradeMode            TdMode          `json:"tdMode"`
	Tag                  string          `json:"tag"`
	RecurringList        []RecurringItem `json:"recurringList"`
	CreationTime         time.Time       `json:"cTime"`
	UpdateTime           time.Time       `json:"uTime"`
}

// GetRecurringSubOrdersService -- GET /api/v5/tradingBot/recurring/sub-orders (Read)
//
// Returns the individual sub-orders placed by a recurring-buy strategy.
type GetRecurringSubOrdersService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetRecurringSubOrdersService(algoId string) *GetRecurringSubOrdersService {
	return &GetRecurringSubOrdersService{c: c, params: map[string]string{"algoId": algoId}}
}

// SetOrdId filters to a single sub-order id.
func (s *GetRecurringSubOrdersService) SetOrdId(ordId string) *GetRecurringSubOrdersService {
	s.params["ordId"] = ordId
	return s
}

// SetAfter pages backward from the given ordId (records older than it).
func (s *GetRecurringSubOrdersService) SetAfter(after string) *GetRecurringSubOrdersService {
	s.params["after"] = after
	return s
}

// SetBefore pages forward from the given ordId (records newer than it).
func (s *GetRecurringSubOrdersService) SetBefore(before string) *GetRecurringSubOrdersService {
	s.params["before"] = before
	return s
}

// SetLimit caps the number of returned sub-orders (max 100, default 100).
func (s *GetRecurringSubOrdersService) SetLimit(limit int) *GetRecurringSubOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetRecurringSubOrdersService) Do(ctx context.Context) ([]RecurringSubOrder, error) {
	req := request.Get(ctx, s.c, "/api/v5/tradingBot/recurring/sub-orders", s.params).WithSign()
	return request.DoList[RecurringSubOrder](req)
}

// RecurringSubOrder is a single underlying order executed by a recurring-buy
// strategy.
type RecurringSubOrder struct {
	AlgoID              string          `json:"algoId"`
	AlgoClientOrderID   string          `json:"algoClOrdId"`
	InstrumentType      InstType        `json:"instType"`
	BotID               string          `json:"botId"`
	InstrumentID        string          `json:"instId"`
	OrderID             string          `json:"ordId"`
	ClientOrderID       string          `json:"clOrdId"`
	OrderType           OrdType         `json:"ordType"`
	Side                Side            `json:"side"`
	Price               decimal.Decimal `json:"px"`
	Size                decimal.Decimal `json:"sz"`
	State               OrdState        `json:"state"`
	AccumulatedFillSize decimal.Decimal `json:"accFillSz"`
	AveragePrice        decimal.Decimal `json:"avgPx"`
	Fee                 decimal.Decimal `json:"fee"`
	FeeCurrency         string          `json:"feeCcy"`
	Tag                 string          `json:"tag"`
	CreationTime        time.Time       `json:"cTime"`
	UpdateTime          time.Time       `json:"uTime"`
}
