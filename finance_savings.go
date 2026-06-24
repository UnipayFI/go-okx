package okx

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-okx/request"
	"github.com/shopspring/decimal"
)

// SavingsSide selects whether a Simple Earn Flexible (savings) purchase /
// redemption order is a subscribe (purchase) or a redeem (redempt).
type SavingsSide string

const (
	SavingsSidePurchase SavingsSide = "purchase"
	SavingsSideRedempt  SavingsSide = "redempt"
)

// GetSavingsBalanceService -- GET /api/v5/finance/savings/balance (Read)
//
// Returns the Simple Earn Flexible (savings) balance of each currency the
// account is subscribed to.
type GetSavingsBalanceService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetSavingsBalanceService() *GetSavingsBalanceService {
	return &GetSavingsBalanceService{c: c, params: map[string]string{}}
}

// SetCcy filters by currency.
func (s *GetSavingsBalanceService) SetCcy(ccy string) *GetSavingsBalanceService {
	s.params["ccy"] = ccy
	return s
}

func (s *GetSavingsBalanceService) Do(ctx context.Context) ([]SavingsBalance, error) {
	req := request.Get(ctx, s.c, "/api/v5/finance/savings/balance", s.params).WithSign()
	return request.DoList[SavingsBalance](req)
}

// SavingsBalance is one currency's savings balance. The validating account holds
// no savings (the endpoint returns an empty data array), so the field set is
// modeled from the OKX doc field table.
type SavingsBalance struct {
	Currency         string          `json:"ccy"`
	Amount           decimal.Decimal `json:"amt"`
	Earnings         decimal.Decimal `json:"earnings"`
	Rate             decimal.Decimal `json:"rate"`
	LoanAmount       decimal.Decimal `json:"loanAmt"`
	PendingAmount    decimal.Decimal `json:"pendingAmt"`
	RedemptionAmount decimal.Decimal `json:"redemptAmt"`
}

// SetSavingsPurchaseRedemptionService -- POST /api/v5/finance/savings/purchase-redemption (Trade)
//
// Subscribes to (purchase) or redeems (redempt) a currency in Simple Earn
// Flexible savings.
//
// State-changing: NOT exercised by the test suite.
type SetSavingsPurchaseRedemptionService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewSetSavingsPurchaseRedemptionService(ccy string, amt decimal.Decimal, side SavingsSide) *SetSavingsPurchaseRedemptionService {
	return &SetSavingsPurchaseRedemptionService{c: c, body: map[string]any{
		"ccy":  ccy,
		"amt":  amt.String(),
		"side": string(side),
	}}
}

// SetRate sets the lending rate (between 0.01 and 3.65, only applicable when
// subscribing).
func (s *SetSavingsPurchaseRedemptionService) SetRate(rate decimal.Decimal) *SetSavingsPurchaseRedemptionService {
	s.body["rate"] = rate.String()
	return s
}

func (s *SetSavingsPurchaseRedemptionService) Do(ctx context.Context) (*SavingsPurchaseRedemption, error) {
	req := request.Post(ctx, s.c, "/api/v5/finance/savings/purchase-redemption", s.body).WithSign()
	return request.DoOne[SavingsPurchaseRedemption](req)
}

// SavingsPurchaseRedemption is the ack of a savings purchase/redemption action.
type SavingsPurchaseRedemption struct {
	Currency string          `json:"ccy"`
	Amount   decimal.Decimal `json:"amt"`
	Side     SavingsSide     `json:"side"`
	Rate     decimal.Decimal `json:"rate"`
}

// SetSavingsLendingRateService -- POST /api/v5/finance/savings/set-lending-rate (Trade)
//
// Sets the minimum lending rate the account is willing to accept for a savings
// currency.
//
// State-changing: NOT exercised by the test suite.
type SetSavingsLendingRateService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewSetSavingsLendingRateService(ccy string, rate decimal.Decimal) *SetSavingsLendingRateService {
	return &SetSavingsLendingRateService{c: c, body: map[string]any{
		"ccy":  ccy,
		"rate": rate.String(),
	}}
}

func (s *SetSavingsLendingRateService) Do(ctx context.Context) (*SavingsLendingRate, error) {
	req := request.Post(ctx, s.c, "/api/v5/finance/savings/set-lending-rate", s.body).WithSign()
	return request.DoOne[SavingsLendingRate](req)
}

// SavingsLendingRate is the ack of a set-lending-rate action.
type SavingsLendingRate struct {
	Currency string          `json:"ccy"`
	Rate     decimal.Decimal `json:"rate"`
}

// GetSavingsLendingHistoryService -- GET /api/v5/finance/savings/lending-history (Read)
//
// Returns the account's savings lending history.
type GetSavingsLendingHistoryService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetSavingsLendingHistoryService() *GetSavingsLendingHistoryService {
	return &GetSavingsLendingHistoryService{c: c, params: map[string]string{}}
}

// SetCcy filters by currency.
func (s *GetSavingsLendingHistoryService) SetCcy(ccy string) *GetSavingsLendingHistoryService {
	s.params["ccy"] = ccy
	return s
}

// SetAfter paginates to records earlier than the given time (older).
func (s *GetSavingsLendingHistoryService) SetAfter(t time.Time) *GetSavingsLendingHistoryService {
	s.params["after"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetBefore paginates to records later than the given time (newer).
func (s *GetSavingsLendingHistoryService) SetBefore(t time.Time) *GetSavingsLendingHistoryService {
	s.params["before"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetSavingsLendingHistoryService) SetLimit(limit int) *GetSavingsLendingHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetSavingsLendingHistoryService) Do(ctx context.Context) ([]SavingsLendingHistory, error) {
	req := request.Get(ctx, s.c, "/api/v5/finance/savings/lending-history", s.params).WithSign()
	return request.DoList[SavingsLendingHistory](req)
}

// SavingsLendingHistory is one savings lending record. The validating account
// has no lending history (the endpoint returns an empty data array), so the
// field set is modeled from the OKX doc field table.
type SavingsLendingHistory struct {
	Currency  string          `json:"ccy"`
	Amount    decimal.Decimal `json:"amt"`
	Earnings  decimal.Decimal `json:"earnings"`
	Rate      decimal.Decimal `json:"rate"`
	Timestamp time.Time       `json:"ts"`
}

// GetSavingsLendingRateSummaryService -- GET /api/v5/finance/savings/lending-rate-summary (Read)
//
// Returns the current public lending-rate summary (estimated/previous/average
// rates) per savings currency.
type GetSavingsLendingRateSummaryService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetSavingsLendingRateSummaryService() *GetSavingsLendingRateSummaryService {
	return &GetSavingsLendingRateSummaryService{c: c, params: map[string]string{}}
}

// SetCcy filters by currency.
func (s *GetSavingsLendingRateSummaryService) SetCcy(ccy string) *GetSavingsLendingRateSummaryService {
	s.params["ccy"] = ccy
	return s
}

func (s *GetSavingsLendingRateSummaryService) Do(ctx context.Context) ([]SavingsLendingRateSummary, error) {
	req := request.Get(ctx, s.c, "/api/v5/finance/savings/lending-rate-summary", s.params).WithSign()
	return request.DoList[SavingsLendingRateSummary](req)
}

// SavingsLendingRateSummary is one currency's lending-rate summary.
type SavingsLendingRateSummary struct {
	Currency         string          `json:"ccy"`
	AverageAmount    decimal.Decimal `json:"avgAmt"`
	AverageAmountUSD decimal.Decimal `json:"avgAmtUsd"`
	AverageRate      decimal.Decimal `json:"avgRate"`
	PreRate          decimal.Decimal `json:"preRate"`
	EstimatedRate    decimal.Decimal `json:"estRate"`
}

// GetSavingsLendingRateHistoryService -- GET /api/v5/finance/savings/lending-rate-history (Read)
//
// Returns the historical public lending rates per savings currency.
type GetSavingsLendingRateHistoryService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetSavingsLendingRateHistoryService() *GetSavingsLendingRateHistoryService {
	return &GetSavingsLendingRateHistoryService{c: c, params: map[string]string{}}
}

// SetCcy filters by currency.
func (s *GetSavingsLendingRateHistoryService) SetCcy(ccy string) *GetSavingsLendingRateHistoryService {
	s.params["ccy"] = ccy
	return s
}

// SetAfter paginates to records earlier than the given time (older).
func (s *GetSavingsLendingRateHistoryService) SetAfter(t time.Time) *GetSavingsLendingRateHistoryService {
	s.params["after"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetBefore paginates to records later than the given time (newer).
func (s *GetSavingsLendingRateHistoryService) SetBefore(t time.Time) *GetSavingsLendingRateHistoryService {
	s.params["before"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetSavingsLendingRateHistoryService) SetLimit(limit int) *GetSavingsLendingRateHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetSavingsLendingRateHistoryService) Do(ctx context.Context) ([]SavingsLendingRateHistory, error) {
	req := request.Get(ctx, s.c, "/api/v5/finance/savings/lending-rate-history", s.params).WithSign()
	return request.DoList[SavingsLendingRateHistory](req)
}

// SavingsLendingRateHistory is one historical lending-rate record.
type SavingsLendingRateHistory struct {
	Currency    string          `json:"ccy"`
	Amount      decimal.Decimal `json:"amt"`
	LendingRate decimal.Decimal `json:"lendingRate"`
	Rate        decimal.Decimal `json:"rate"`
	Timestamp   time.Time       `json:"ts"`
}
