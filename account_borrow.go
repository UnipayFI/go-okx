package okx

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-okx/request"
	"github.com/shopspring/decimal"
)

// BorrowRepayType selects whether a borrow/repay record (or action) is a borrow
// or a repayment.
type BorrowRepayType string

const (
	BorrowRepayTypeBorrow BorrowRepayType = "borrow"
	BorrowRepayTypeRepay  BorrowRepayType = "repay"
)

// VipLoanState is the lifecycle state of a VIP loan order.
type VipLoanState string

const (
	VipLoanStateBorrowing       VipLoanState = "1"
	VipLoanStateBorrowed        VipLoanState = "2"
	VipLoanStatePartiallyRepaid VipLoanState = "3"
	VipLoanStateRepaid          VipLoanState = "4"
	VipLoanStateRepaying        VipLoanState = "5"
)

// GetInterestAccruedService -- GET /api/v5/account/interest-accrued (Read)
//
// Returns the interest accrued on margin and quick-margin borrowings.
type GetInterestAccruedService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetInterestAccruedService() *GetInterestAccruedService {
	return &GetInterestAccruedService{c: c, params: map[string]string{}}
}

// SetType filters by loan type ("1" VIP loan, "2" market loan).
func (s *GetInterestAccruedService) SetType(typ string) *GetInterestAccruedService {
	s.params["type"] = typ
	return s
}

// SetCcy filters by currency.
func (s *GetInterestAccruedService) SetCcy(ccy string) *GetInterestAccruedService {
	s.params["ccy"] = ccy
	return s
}

// SetInstId filters by instrument id (only applicable to isolated margin).
func (s *GetInterestAccruedService) SetInstId(instId string) *GetInterestAccruedService {
	s.params["instId"] = instId
	return s
}

// SetMgnMode filters by margin mode (cross / isolated).
func (s *GetInterestAccruedService) SetMgnMode(mgnMode MgnMode) *GetInterestAccruedService {
	s.params["mgnMode"] = string(mgnMode)
	return s
}

// SetAfter paginates to records earlier than the given time (older).
func (s *GetInterestAccruedService) SetAfter(t time.Time) *GetInterestAccruedService {
	s.params["after"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetBefore paginates to records later than the given time (newer).
func (s *GetInterestAccruedService) SetBefore(t time.Time) *GetInterestAccruedService {
	s.params["before"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetInterestAccruedService) SetLimit(limit int) *GetInterestAccruedService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetInterestAccruedService) Do(ctx context.Context) ([]InterestAccrued, error) {
	req := request.Get(ctx, s.c, "/api/v5/account/interest-accrued", s.params).WithSign()
	return request.DoList[InterestAccrued](req)
}

// InterestAccrued is one accrued-interest record.
type InterestAccrued struct {
	Type         string          `json:"type"`
	Currency     string          `json:"ccy"`
	InstrumentID string          `json:"instId"`
	MarginMode   MgnMode         `json:"mgnMode"`
	Interest     decimal.Decimal `json:"interest"`
	InterestRate decimal.Decimal `json:"interestRate"`
	Liability    decimal.Decimal `json:"liab"`
	Timestamp    time.Time       `json:"ts"`
}

// GetInterestRateService -- GET /api/v5/account/interest-rate (Read)
//
// Returns the current per-currency borrowing interest rate for the account.
type GetInterestRateService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetInterestRateService() *GetInterestRateService {
	return &GetInterestRateService{c: c, params: map[string]string{}}
}

// SetCcy filters by currency.
func (s *GetInterestRateService) SetCcy(ccy string) *GetInterestRateService {
	s.params["ccy"] = ccy
	return s
}

func (s *GetInterestRateService) Do(ctx context.Context) ([]InterestRate, error) {
	req := request.Get(ctx, s.c, "/api/v5/account/interest-rate", s.params).WithSign()
	return request.DoList[InterestRate](req)
}

// InterestRate is a currency's current account borrowing interest rate.
type InterestRate struct {
	Currency          string          `json:"ccy"`
	InterestRate      decimal.Decimal `json:"interestRate"`
	NextEstimatedRate decimal.Decimal `json:"nextEstRate"`
	Timestamp         time.Time       `json:"ts"`
}

// GetInterestLimitsService -- GET /api/v5/account/interest-limits (Read)
//
// Returns the borrowing-interest summary and the per-currency loan-quota /
// interest schedule for the account.
type GetInterestLimitsService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetInterestLimitsService() *GetInterestLimitsService {
	return &GetInterestLimitsService{c: c, params: map[string]string{}}
}

// SetType filters by loan type ("1" VIP loan, "2" market loan).
func (s *GetInterestLimitsService) SetType(typ string) *GetInterestLimitsService {
	s.params["type"] = typ
	return s
}

// SetCcy filters by currency.
func (s *GetInterestLimitsService) SetCcy(ccy string) *GetInterestLimitsService {
	s.params["ccy"] = ccy
	return s
}

func (s *GetInterestLimitsService) Do(ctx context.Context) (*InterestLimits, error) {
	req := request.Get(ctx, s.c, "/api/v5/account/interest-limits", s.params).WithSign()
	return request.DoOne[InterestLimits](req)
}

// InterestLimits is the account's interest summary and per-currency loan-quota
// schedule.
type InterestLimits struct {
	Debt             decimal.Decimal        `json:"debt"`
	Interest         decimal.Decimal        `json:"interest"`
	NextDiscountTime time.Time              `json:"nextDiscountTime"`
	NextInterestTime time.Time              `json:"nextInterestTime"`
	LoanAlloc        decimal.Decimal        `json:"loanAlloc"`
	Records          []InterestLimitsRecord `json:"records"`
}

// InterestLimitsRecord is one currency's loan-quota and interest detail within
// the account interest-limits schedule.
type InterestLimitsRecord struct {
	Currency                 string            `json:"ccy"`
	Rate                     decimal.Decimal   `json:"rate"`
	LoanQuota                decimal.Decimal   `json:"loanQuota"`
	SurplusLimit             decimal.Decimal   `json:"surplusLmt"`
	SurplusLimitDetails      SurplusLmtDetails `json:"surplusLmtDetails"`
	UsedLimit                decimal.Decimal   `json:"usedLmt"`
	Interest                 decimal.Decimal   `json:"interest"`
	InterestFreeLiability    decimal.Decimal   `json:"interestFreeLiab"`
	PositionLoan             decimal.Decimal   `json:"posLoan"`
	AvailableLoan            decimal.Decimal   `json:"availLoan"`
	UsedLoan                 decimal.Decimal   `json:"usedLoan"`
	AverageRate              decimal.Decimal   `json:"avgRate"`
	PotentialBorrowingAmount decimal.Decimal   `json:"potentialBorrowingAmt"`
}

// SurplusLmtDetails breaks the surplus borrowing limit down by its binding
// constraint. OKX returns it as an empty object when no limit is currently
// constraining the currency.
type SurplusLmtDetails struct {
	AllAccountRemainingQuota     decimal.Decimal `json:"allAcctRemainingQuota"`
	CurrentAccountRemainingQuota decimal.Decimal `json:"curAcctRemainingQuota"`
	PlatformRemainingQuota       decimal.Decimal `json:"platRemainingQuota"`
}

// GetSpotBorrowRepayHistoryService -- GET /api/v5/account/spot-borrow-repay-history (Read)
//
// Returns the history of spot (auto / manual) borrows and repayments.
type GetSpotBorrowRepayHistoryService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetSpotBorrowRepayHistoryService() *GetSpotBorrowRepayHistoryService {
	return &GetSpotBorrowRepayHistoryService{c: c, params: map[string]string{}}
}

// SetCcy filters by currency.
func (s *GetSpotBorrowRepayHistoryService) SetCcy(ccy string) *GetSpotBorrowRepayHistoryService {
	s.params["ccy"] = ccy
	return s
}

// SetType filters by event type (auto_borrow / auto_repay / manual_borrow /
// manual_repay).
func (s *GetSpotBorrowRepayHistoryService) SetType(typ string) *GetSpotBorrowRepayHistoryService {
	s.params["type"] = typ
	return s
}

// SetAfter paginates to records earlier than the given time (older).
func (s *GetSpotBorrowRepayHistoryService) SetAfter(t time.Time) *GetSpotBorrowRepayHistoryService {
	s.params["after"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetBefore paginates to records later than the given time (newer).
func (s *GetSpotBorrowRepayHistoryService) SetBefore(t time.Time) *GetSpotBorrowRepayHistoryService {
	s.params["before"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetSpotBorrowRepayHistoryService) SetLimit(limit int) *GetSpotBorrowRepayHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetSpotBorrowRepayHistoryService) Do(ctx context.Context) ([]SpotBorrowRepayHistory, error) {
	req := request.Get(ctx, s.c, "/api/v5/account/spot-borrow-repay-history", s.params).WithSign()
	return request.DoList[SpotBorrowRepayHistory](req)
}

// SpotBorrowRepayHistory is one spot borrow/repay event.
type SpotBorrowRepayHistory struct {
	Currency            string          `json:"ccy"`
	Type                string          `json:"type"`
	Amount              decimal.Decimal `json:"amt"`
	AccumulatedBorrowed decimal.Decimal `json:"accBorrowed"`
	Timestamp           time.Time       `json:"ts"`
}

// GetBorrowRepayHistoryService -- GET /api/v5/account/borrow-repay-history (Read)
//
// Returns the VIP-loan borrow and repay history.
type GetBorrowRepayHistoryService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetBorrowRepayHistoryService() *GetBorrowRepayHistoryService {
	return &GetBorrowRepayHistoryService{c: c, params: map[string]string{}}
}

// SetCcy filters by currency.
func (s *GetBorrowRepayHistoryService) SetCcy(ccy string) *GetBorrowRepayHistoryService {
	s.params["ccy"] = ccy
	return s
}

// SetAfter paginates to records earlier than the given time (older).
func (s *GetBorrowRepayHistoryService) SetAfter(t time.Time) *GetBorrowRepayHistoryService {
	s.params["after"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetBefore paginates to records later than the given time (newer).
func (s *GetBorrowRepayHistoryService) SetBefore(t time.Time) *GetBorrowRepayHistoryService {
	s.params["before"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetBorrowRepayHistoryService) SetLimit(limit int) *GetBorrowRepayHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetBorrowRepayHistoryService) Do(ctx context.Context) ([]BorrowRepayHistory, error) {
	req := request.Get(ctx, s.c, "/api/v5/account/borrow-repay-history", s.params).WithSign()
	return request.DoList[BorrowRepayHistory](req)
}

// BorrowRepayHistory is one VIP-loan borrow/repay event.
type BorrowRepayHistory struct {
	Currency            string          `json:"ccy"`
	Type                BorrowRepayType `json:"type"`
	Amount              decimal.Decimal `json:"amt"`
	AccumulatedBorrowed decimal.Decimal `json:"accBorrowed"`
	Timestamp           time.Time       `json:"ts"`
}

// GetVipLoanOrderListService -- GET /api/v5/account/vip-loan-order-list (Read)
//
// Returns the account's VIP loan orders.
type GetVipLoanOrderListService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetVipLoanOrderListService() *GetVipLoanOrderListService {
	return &GetVipLoanOrderListService{c: c, params: map[string]string{}}
}

// SetOrdId filters by VIP loan order id.
func (s *GetVipLoanOrderListService) SetOrdId(ordId string) *GetVipLoanOrderListService {
	s.params["ordId"] = ordId
	return s
}

// SetState filters by loan order state (1 borrowing / 2 borrowed / 3 partially
// repaid / 4 repaid / 5 repaying).
func (s *GetVipLoanOrderListService) SetState(state VipLoanState) *GetVipLoanOrderListService {
	s.params["state"] = string(state)
	return s
}

// SetCcy filters by currency.
func (s *GetVipLoanOrderListService) SetCcy(ccy string) *GetVipLoanOrderListService {
	s.params["ccy"] = ccy
	return s
}

// SetAfter paginates to orders earlier than the given time (older).
func (s *GetVipLoanOrderListService) SetAfter(t time.Time) *GetVipLoanOrderListService {
	s.params["after"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetBefore paginates to orders later than the given time (newer).
func (s *GetVipLoanOrderListService) SetBefore(t time.Time) *GetVipLoanOrderListService {
	s.params["before"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of orders returned (max 100).
func (s *GetVipLoanOrderListService) SetLimit(limit int) *GetVipLoanOrderListService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetVipLoanOrderListService) Do(ctx context.Context) ([]VipLoanOrder, error) {
	req := request.Get(ctx, s.c, "/api/v5/account/vip-loan-order-list", s.params).WithSign()
	return request.DoList[VipLoanOrder](req)
}

// VipLoanOrder is one VIP loan order summary.
type VipLoanOrder struct {
	OrderID      string          `json:"ordId"`
	State        VipLoanState    `json:"state"`
	Currency     string          `json:"ccy"`
	BorrowAmount decimal.Decimal `json:"borrowAmt"`
	CurrentRate  decimal.Decimal `json:"curRate"`
	DueAmount    decimal.Decimal `json:"dueAmt"`
	RepayAmount  decimal.Decimal `json:"repayAmt"`
	Interest     decimal.Decimal `json:"interest"`
	Timestamp    time.Time       `json:"ts"`
}

// GetVipLoanOrderDetailService -- GET /api/v5/account/vip-loan-order-detail (Read)
//
// Returns the detail (per-event borrow/repay records) of a single VIP loan
// order.
type GetVipLoanOrderDetailService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetVipLoanOrderDetailService(ordId string) *GetVipLoanOrderDetailService {
	return &GetVipLoanOrderDetailService{c: c, params: map[string]string{"ordId": ordId}}
}

// SetCcy filters by currency.
func (s *GetVipLoanOrderDetailService) SetCcy(ccy string) *GetVipLoanOrderDetailService {
	s.params["ccy"] = ccy
	return s
}

// SetAfter paginates to records earlier than the given time (older).
func (s *GetVipLoanOrderDetailService) SetAfter(t time.Time) *GetVipLoanOrderDetailService {
	s.params["after"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetBefore paginates to records later than the given time (newer).
func (s *GetVipLoanOrderDetailService) SetBefore(t time.Time) *GetVipLoanOrderDetailService {
	s.params["before"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetVipLoanOrderDetailService) SetLimit(limit int) *GetVipLoanOrderDetailService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetVipLoanOrderDetailService) Do(ctx context.Context) ([]VipLoanOrderDetail, error) {
	req := request.Get(ctx, s.c, "/api/v5/account/vip-loan-order-detail", s.params).WithSign()
	return request.DoList[VipLoanOrderDetail](req)
}

// VipLoanOrderDetail is the per-event detail of a VIP loan order.
type VipLoanOrderDetail struct {
	Currency         string                    `json:"ccy"`
	CurrentRate      decimal.Decimal           `json:"curRate"`
	DueAmount        decimal.Decimal           `json:"dueAmt"`
	TotalRepayAmount decimal.Decimal           `json:"totalRepayAmt"`
	TotalInterest    decimal.Decimal           `json:"totalInterest"`
	Timestamp        time.Time                 `json:"ts"`
	BorrowAmount     decimal.Decimal           `json:"borrowAmt"`
	RepayAmount      decimal.Decimal           `json:"repayAmt"`
	List             []VipLoanOrderDetailEvent `json:"list"`
}

// VipLoanOrderDetailEvent is one borrow/repay/interest event within a VIP loan
// order's history.
type VipLoanOrderDetailEvent struct {
	Type      BorrowRepayType `json:"type"`
	Amount    decimal.Decimal `json:"amt"`
	Currency  string          `json:"ccy"`
	Timestamp time.Time       `json:"ts"`
}

// GetVipInterestAccruedService -- GET /api/v5/account/vip-interest-accrued (Read)
//
// Returns the interest accrued on the account's VIP loans.
type GetVipInterestAccruedService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetVipInterestAccruedService() *GetVipInterestAccruedService {
	return &GetVipInterestAccruedService{c: c, params: map[string]string{}}
}

// SetCcy filters by currency.
func (s *GetVipInterestAccruedService) SetCcy(ccy string) *GetVipInterestAccruedService {
	s.params["ccy"] = ccy
	return s
}

// SetOrdId filters by VIP loan order id.
func (s *GetVipInterestAccruedService) SetOrdId(ordId string) *GetVipInterestAccruedService {
	s.params["ordId"] = ordId
	return s
}

// SetAfter paginates to records earlier than the given time (older).
func (s *GetVipInterestAccruedService) SetAfter(t time.Time) *GetVipInterestAccruedService {
	s.params["after"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetBefore paginates to records later than the given time (newer).
func (s *GetVipInterestAccruedService) SetBefore(t time.Time) *GetVipInterestAccruedService {
	s.params["before"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetVipInterestAccruedService) SetLimit(limit int) *GetVipInterestAccruedService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetVipInterestAccruedService) Do(ctx context.Context) ([]VipInterestAccrued, error) {
	req := request.Get(ctx, s.c, "/api/v5/account/vip-interest-accrued", s.params).WithSign()
	return request.DoList[VipInterestAccrued](req)
}

// VipInterestAccrued is one VIP-loan accrued-interest record.
type VipInterestAccrued struct {
	OrderID      string          `json:"ordId"`
	Currency     string          `json:"ccy"`
	Interest     decimal.Decimal `json:"interest"`
	InterestRate decimal.Decimal `json:"interestRate"`
	Liability    decimal.Decimal `json:"liab"`
	Timestamp    time.Time       `json:"ts"`
}

// GetVipInterestDeductedService -- GET /api/v5/account/vip-interest-deducted (Read)
//
// Returns the interest deducted from the account's VIP loans.
type GetVipInterestDeductedService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetVipInterestDeductedService() *GetVipInterestDeductedService {
	return &GetVipInterestDeductedService{c: c, params: map[string]string{}}
}

// SetCcy filters by currency.
func (s *GetVipInterestDeductedService) SetCcy(ccy string) *GetVipInterestDeductedService {
	s.params["ccy"] = ccy
	return s
}

// SetOrdId filters by VIP loan order id.
func (s *GetVipInterestDeductedService) SetOrdId(ordId string) *GetVipInterestDeductedService {
	s.params["ordId"] = ordId
	return s
}

// SetAfter paginates to records earlier than the given time (older).
func (s *GetVipInterestDeductedService) SetAfter(t time.Time) *GetVipInterestDeductedService {
	s.params["after"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetBefore paginates to records later than the given time (newer).
func (s *GetVipInterestDeductedService) SetBefore(t time.Time) *GetVipInterestDeductedService {
	s.params["before"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetVipInterestDeductedService) SetLimit(limit int) *GetVipInterestDeductedService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetVipInterestDeductedService) Do(ctx context.Context) ([]VipInterestDeducted, error) {
	req := request.Get(ctx, s.c, "/api/v5/account/vip-interest-deducted", s.params).WithSign()
	return request.DoList[VipInterestDeducted](req)
}

// VipInterestDeducted is one VIP-loan interest-deduction record.
type VipInterestDeducted struct {
	OrderID      string          `json:"ordId"`
	Currency     string          `json:"ccy"`
	Interest     decimal.Decimal `json:"interest"`
	InterestRate decimal.Decimal `json:"interestRate"`
	Liability    decimal.Decimal `json:"liab"`
	Timestamp    time.Time       `json:"ts"`
}

// SetSpotManualBorrowRepayService -- POST /api/v5/account/spot-manual-borrow-repay (Trade)
//
// Manually borrows or repays spot in the cross-margin account.
//
// State-changing: NOT exercised by the test suite.
type SetSpotManualBorrowRepayService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewSetSpotManualBorrowRepayService(ccy string, side BorrowRepayType, amt decimal.Decimal) *SetSpotManualBorrowRepayService {
	return &SetSpotManualBorrowRepayService{c: c, body: map[string]any{
		"ccy":  ccy,
		"side": string(side),
		"amt":  amt.String(),
	}}
}

func (s *SetSpotManualBorrowRepayService) Do(ctx context.Context) (*SpotManualBorrowRepay, error) {
	req := request.Post(ctx, s.c, "/api/v5/account/spot-manual-borrow-repay", s.body).WithSign()
	return request.DoOne[SpotManualBorrowRepay](req)
}

// SpotManualBorrowRepay is the ack of a manual spot borrow/repay action.
type SpotManualBorrowRepay struct {
	Currency string          `json:"ccy"`
	Side     BorrowRepayType `json:"side"`
	Amount   decimal.Decimal `json:"amt"`
}

// SetAutoRepayService -- POST /api/v5/account/set-auto-repay (Trade)
//
// Enables or disables automatic repayment of spot borrowings.
//
// State-changing: NOT exercised by the test suite.
type SetAutoRepayService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewSetAutoRepayService(autoRepay bool) *SetAutoRepayService {
	return &SetAutoRepayService{c: c, body: map[string]any{
		"autoRepay": autoRepay,
	}}
}

func (s *SetAutoRepayService) Do(ctx context.Context) (*AutoRepay, error) {
	req := request.Post(ctx, s.c, "/api/v5/account/set-auto-repay", s.body).WithSign()
	return request.DoOne[AutoRepay](req)
}

// AutoRepay is the ack of a set-auto-repay action.
type AutoRepay struct {
	AutoRepay bool `json:"autoRepay"`
}

// SetBorrowRepayService -- POST /api/v5/account/borrow-repay (Trade)
//
// Borrows or repays a VIP loan.
//
// State-changing: NOT exercised by the test suite.
type SetBorrowRepayService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewSetBorrowRepayService(ccy string, side BorrowRepayType, amt decimal.Decimal) *SetBorrowRepayService {
	return &SetBorrowRepayService{c: c, body: map[string]any{
		"ccy":  ccy,
		"side": string(side),
		"amt":  amt.String(),
	}}
}

// SetOrdId targets an existing VIP loan order (required when repaying).
func (s *SetBorrowRepayService) SetOrdId(ordId string) *SetBorrowRepayService {
	s.body["ordId"] = ordId
	return s
}

func (s *SetBorrowRepayService) Do(ctx context.Context) (*BorrowRepay, error) {
	req := request.Post(ctx, s.c, "/api/v5/account/borrow-repay", s.body).WithSign()
	return request.DoOne[BorrowRepay](req)
}

// BorrowRepay is the ack of a VIP-loan borrow/repay action.
type BorrowRepay struct {
	Currency string          `json:"ccy"`
	Side     BorrowRepayType `json:"side"`
	Amount   decimal.Decimal `json:"amt"`
	OrderID  string          `json:"ordId"`
	State    VipLoanState    `json:"state"`
}
