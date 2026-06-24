package okx

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-okx/request"
	"github.com/shopspring/decimal"
)

// This file covers the flexible-loan (/api/v5/finance/flexible-loan/*) and
// fixed-loan (/api/v5/finance/fixed-loan/*) endpoints. The fixed-loan endpoints
// are restricted on the validating key (every path returns HTTP 403 Forbidden),
// so their structs are modeled from the OKX doc field tables; the request paths
// themselves are verified (403, not 404).

// FlexLoanAdjustType selects whether collateral is added or reduced when
// adjusting a flexible-loan position.
type FlexLoanAdjustType string

const (
	FlexLoanAdjustTypeAdd    FlexLoanAdjustType = "add"
	FlexLoanAdjustTypeReduce FlexLoanAdjustType = "reduce"
)

// GetFlexLoanBorrowCurrenciesService -- GET /api/v5/finance/flexible-loan/borrow-currencies (Read)
//
// Returns the currencies that can be borrowed via flexible loan.
type GetFlexLoanBorrowCurrenciesService struct {
	c *Client
}

func (c *Client) NewGetFlexLoanBorrowCurrenciesService() *GetFlexLoanBorrowCurrenciesService {
	return &GetFlexLoanBorrowCurrenciesService{c: c}
}

func (s *GetFlexLoanBorrowCurrenciesService) Do(ctx context.Context) ([]FlexLoanBorrowCurrency, error) {
	req := request.Get(ctx, s.c, "/api/v5/finance/flexible-loan/borrow-currencies").WithSign()
	return request.DoList[FlexLoanBorrowCurrency](req)
}

// FlexLoanBorrowCurrency is a single borrowable flexible-loan currency.
type FlexLoanBorrowCurrency struct {
	BorrowCurrency string `json:"borrowCcy"`
}

// GetFlexLoanCollateralAssetsService -- GET /api/v5/finance/flexible-loan/collateral-assets (Read)
//
// Returns the assets eligible as flexible-loan collateral. The data is a single
// object wrapping an "assets" array.
type GetFlexLoanCollateralAssetsService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetFlexLoanCollateralAssetsService() *GetFlexLoanCollateralAssetsService {
	return &GetFlexLoanCollateralAssetsService{c: c, params: map[string]string{}}
}

// SetCcy filters by collateral currency.
func (s *GetFlexLoanCollateralAssetsService) SetCcy(ccy string) *GetFlexLoanCollateralAssetsService {
	s.params["ccy"] = ccy
	return s
}

func (s *GetFlexLoanCollateralAssetsService) Do(ctx context.Context) (*FlexLoanCollateralAssets, error) {
	req := request.Get(ctx, s.c, "/api/v5/finance/flexible-loan/collateral-assets", s.params).WithSign()
	return request.DoOne[FlexLoanCollateralAssets](req)
}

// FlexLoanCollateralAssets wraps the list of collateral assets. The validating
// account had no eligible collateral (assets is []), so FlexLoanCollateralAsset
// is modeled from the OKX doc field table.
type FlexLoanCollateralAssets struct {
	Assets []FlexLoanCollateralAsset `json:"assets"`
}

// FlexLoanCollateralAsset is one asset usable as flexible-loan collateral.
type FlexLoanCollateralAsset struct {
	Currency string          `json:"ccy"`
	Amount   decimal.Decimal `json:"amt"`
}

// GetFlexLoanMaxLoanService -- POST /api/v5/finance/flexible-loan/max-loan (Read)
//
// Computes the maximum loanable amount for a borrow currency given the supplied
// (or current) collateral. This is a non-state-changing calculation served over
// POST (it returns 405 to GET).
type GetFlexLoanMaxLoanService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewGetFlexLoanMaxLoanService(borrowCcy string) *GetFlexLoanMaxLoanService {
	return &GetFlexLoanMaxLoanService{c: c, body: map[string]any{"borrowCcy": borrowCcy}}
}

// SetSupCollateral sets the supplementary collateral used in the calculation
// (each item a {ccy, amt} pair).
func (s *GetFlexLoanMaxLoanService) SetSupCollateral(supCollateral []FlexLoanSupCollateral) *GetFlexLoanMaxLoanService {
	s.body["supCollateral"] = supCollateral
	return s
}

func (s *GetFlexLoanMaxLoanService) Do(ctx context.Context) ([]FlexLoanMaxLoan, error) {
	req := request.Post(ctx, s.c, "/api/v5/finance/flexible-loan/max-loan", s.body).WithSign()
	return request.DoList[FlexLoanMaxLoan](req)
}

// FlexLoanSupCollateral is one supplementary-collateral item passed to the
// max-loan calculation.
type FlexLoanSupCollateral struct {
	Currency string          `json:"ccy"`
	Amount   decimal.Decimal `json:"amt"`
}

// FlexLoanMaxLoan is the computed max-loan result for a borrow currency.
type FlexLoanMaxLoan struct {
	BorrowCurrency string          `json:"borrowCcy"`
	MaxLoan        decimal.Decimal `json:"maxLoan"`
	NotionalUSD    decimal.Decimal `json:"notionalUsd"`
	RemainingQuota decimal.Decimal `json:"remainingQuota"`
}

// AdjustFlexLoanCollateralService -- POST /api/v5/finance/flexible-loan/adjust-collateral (Trade)
//
// Adds or reduces collateral on a flexible-loan position. State-changing:
// implemented but never executed by the test suite. The acknowledgement shape is
// modeled from the OKX doc field table.
type AdjustFlexLoanCollateralService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewAdjustFlexLoanCollateralService(typ FlexLoanAdjustType, collateralCcy string, collateralAmt decimal.Decimal) *AdjustFlexLoanCollateralService {
	return &AdjustFlexLoanCollateralService{c: c, body: map[string]any{
		"type":          string(typ),
		"collateralCcy": collateralCcy,
		"collateralAmt": collateralAmt.String(),
	}}
}

func (s *AdjustFlexLoanCollateralService) Do(ctx context.Context) (*FlexLoanAdjustCollateral, error) {
	req := request.Post(ctx, s.c, "/api/v5/finance/flexible-loan/adjust-collateral", s.body).WithSign()
	return request.DoOne[FlexLoanAdjustCollateral](req)
}

// FlexLoanAdjustCollateral is the acknowledgement of a collateral adjustment.
type FlexLoanAdjustCollateral struct {
	Type               FlexLoanAdjustType `json:"type"`
	CollateralCurrency string             `json:"collateralCcy"`
	CollateralAmount   decimal.Decimal    `json:"collateralAmt"`
}

// GetFlexLoanLoanInfoService -- GET /api/v5/finance/flexible-loan/loan-info (Read)
//
// Returns the account's current flexible-loan position (borrowed amounts,
// collateral and loan-to-value ratios). The validating account has no flexible
// loan (empty data), so FlexLoanInfo is modeled from the OKX doc field table.
type GetFlexLoanLoanInfoService struct {
	c *Client
}

func (c *Client) NewGetFlexLoanLoanInfoService() *GetFlexLoanLoanInfoService {
	return &GetFlexLoanLoanInfoService{c: c}
}

func (s *GetFlexLoanLoanInfoService) Do(ctx context.Context) (*FlexLoanInfo, error) {
	req := request.Get(ctx, s.c, "/api/v5/finance/flexible-loan/loan-info").WithSign()
	return request.DoOne[FlexLoanInfo](req)
}

// FlexLoanInfo is the account's current flexible-loan position.
type FlexLoanInfo struct {
	LiquidationLTV          decimal.Decimal          `json:"liquidationLTV"`
	CurrentLTV              decimal.Decimal          `json:"curLTV"`
	RiskWarningLTV          decimal.Decimal          `json:"riskWarningLTV"`
	InitialLTV              decimal.Decimal          `json:"initLTV"`
	TotalBorrowValueUSD     decimal.Decimal          `json:"totalBorrowValueUsd"`
	TotalCollateralValueUSD decimal.Decimal          `json:"totalCollateralValueUsd"`
	LoanData                []FlexLoanInfoLoan       `json:"loanData"`
	CollateralData          []FlexLoanInfoCollateral `json:"collateralData"`
}

// FlexLoanInfoLoan is one borrowed currency within a flexible-loan position.
type FlexLoanInfoLoan struct {
	BorrowCurrency string          `json:"borrowCcy"`
	Amount         decimal.Decimal `json:"amt"`
	NotionalUSD    decimal.Decimal `json:"notionalUsd"`
}

// FlexLoanInfoCollateral is one collateral currency within a flexible-loan
// position.
type FlexLoanInfoCollateral struct {
	CollateralCurrency string          `json:"collateralCcy"`
	Amount             decimal.Decimal `json:"amt"`
	NotionalUSD        decimal.Decimal `json:"notionalUsd"`
}

// GetFlexLoanLoanHistoryService -- GET /api/v5/finance/flexible-loan/loan-history (Read)
//
// Returns the account's flexible-loan event history (borrow / repay /
// liquidation). The validating account has no history (empty data), so
// FlexLoanHistory is modeled from the OKX doc field table.
type GetFlexLoanLoanHistoryService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetFlexLoanLoanHistoryService() *GetFlexLoanLoanHistoryService {
	return &GetFlexLoanLoanHistoryService{c: c, params: map[string]string{}}
}

// SetType filters by event type.
func (s *GetFlexLoanLoanHistoryService) SetType(typ string) *GetFlexLoanLoanHistoryService {
	s.params["type"] = typ
	return s
}

// SetAfter paginates to records earlier than the given time (older).
func (s *GetFlexLoanLoanHistoryService) SetAfter(t time.Time) *GetFlexLoanLoanHistoryService {
	s.params["after"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetBefore paginates to records later than the given time (newer).
func (s *GetFlexLoanLoanHistoryService) SetBefore(t time.Time) *GetFlexLoanLoanHistoryService {
	s.params["before"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetFlexLoanLoanHistoryService) SetLimit(limit int) *GetFlexLoanLoanHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetFlexLoanLoanHistoryService) Do(ctx context.Context) ([]FlexLoanHistory, error) {
	req := request.Get(ctx, s.c, "/api/v5/finance/flexible-loan/loan-history", s.params).WithSign()
	return request.DoList[FlexLoanHistory](req)
}

// FlexLoanHistory is one flexible-loan event record.
type FlexLoanHistory struct {
	ReferenceID string          `json:"refId"`
	Type        string          `json:"type"`
	Currency    string          `json:"ccy"`
	Amount      decimal.Decimal `json:"amt"`
	Timestamp   time.Time       `json:"ts"`
}

// GetFlexLoanInterestAccruedService -- GET /api/v5/finance/flexible-loan/interest-accrued (Read)
//
// Returns the account's accrued flexible-loan interest records. The validating
// account has no records (empty data), so FlexLoanInterestAccrued is modeled
// from the OKX doc field table.
type GetFlexLoanInterestAccruedService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetFlexLoanInterestAccruedService() *GetFlexLoanInterestAccruedService {
	return &GetFlexLoanInterestAccruedService{c: c, params: map[string]string{}}
}

// SetCcy filters by borrowed currency.
func (s *GetFlexLoanInterestAccruedService) SetCcy(ccy string) *GetFlexLoanInterestAccruedService {
	s.params["ccy"] = ccy
	return s
}

// SetAfter paginates to records earlier than the given time (older).
func (s *GetFlexLoanInterestAccruedService) SetAfter(t time.Time) *GetFlexLoanInterestAccruedService {
	s.params["after"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetBefore paginates to records later than the given time (newer).
func (s *GetFlexLoanInterestAccruedService) SetBefore(t time.Time) *GetFlexLoanInterestAccruedService {
	s.params["before"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetFlexLoanInterestAccruedService) SetLimit(limit int) *GetFlexLoanInterestAccruedService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetFlexLoanInterestAccruedService) Do(ctx context.Context) ([]FlexLoanInterestAccrued, error) {
	req := request.Get(ctx, s.c, "/api/v5/finance/flexible-loan/interest-accrued", s.params).WithSign()
	return request.DoList[FlexLoanInterestAccrued](req)
}

// FlexLoanInterestAccrued is one accrued-interest record.
type FlexLoanInterestAccrued struct {
	Currency     string          `json:"ccy"`
	Interest     decimal.Decimal `json:"interest"`
	InterestRate decimal.Decimal `json:"interestRate"`
	Timestamp    time.Time       `json:"ts"`
}

// ---------------------------------------------------------------------------
// Fixed loan (/api/v5/finance/fixed-loan/*) — restricted on this key (HTTP 403).
// Paths are verified (403, not 404); structs are modeled from the OKX docs.
// ---------------------------------------------------------------------------

// GetFixedLoanLendingOffersService -- GET /api/v5/finance/fixed-loan/lending-offers (Read)
//
// Returns the available fixed-loan lending offers (estimated APY and limits per
// currency/term).
type GetFixedLoanLendingOffersService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetFixedLoanLendingOffersService() *GetFixedLoanLendingOffersService {
	return &GetFixedLoanLendingOffersService{c: c, params: map[string]string{}}
}

// SetCcy filters by lending currency.
func (s *GetFixedLoanLendingOffersService) SetCcy(ccy string) *GetFixedLoanLendingOffersService {
	s.params["ccy"] = ccy
	return s
}

// SetTerm filters by lending term (e.g. 30D).
func (s *GetFixedLoanLendingOffersService) SetTerm(term string) *GetFixedLoanLendingOffersService {
	s.params["term"] = term
	return s
}

func (s *GetFixedLoanLendingOffersService) Do(ctx context.Context) ([]FixedLoanLendingOffer, error) {
	req := request.Get(ctx, s.c, "/api/v5/finance/fixed-loan/lending-offers", s.params).WithSign()
	return request.DoList[FixedLoanLendingOffer](req)
}

// FixedLoanLendingOffer is one available fixed-loan lending offer.
type FixedLoanLendingOffer struct {
	Currency        string                        `json:"ccy"`
	Term            string                        `json:"term"`
	Lend            decimal.Decimal               `json:"lend"`
	APY             decimal.Decimal               `json:"apy"`
	EarningCurrency []FixedLoanOfferEarning       `json:"earningCcy"`
	Details         []FixedLoanLendingOfferDetail `json:"details"`
}

// FixedLoanOfferEarning is one earning-currency entry of a lending offer.
type FixedLoanOfferEarning struct {
	EarningCurrency string `json:"earningCcy"`
}

// FixedLoanLendingOfferDetail is one min/max amount detail of a lending offer.
type FixedLoanLendingOfferDetail struct {
	MinLend decimal.Decimal `json:"minLend"`
	MaxLend decimal.Decimal `json:"maxLend"`
}

// GetFixedLoanLendingApyHistoryService -- GET /api/v5/finance/fixed-loan/lending-apy-history (Read)
//
// Returns the historical lending APY for a currency/term.
type GetFixedLoanLendingApyHistoryService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetFixedLoanLendingApyHistoryService(ccy, term string) *GetFixedLoanLendingApyHistoryService {
	return &GetFixedLoanLendingApyHistoryService{c: c, params: map[string]string{
		"ccy":  ccy,
		"term": term,
	}}
}

func (s *GetFixedLoanLendingApyHistoryService) Do(ctx context.Context) ([]FixedLoanLendingApyHistory, error) {
	req := request.Get(ctx, s.c, "/api/v5/finance/fixed-loan/lending-apy-history", s.params).WithSign()
	return request.DoList[FixedLoanLendingApyHistory](req)
}

// FixedLoanLendingApyHistory is one historical APY data point.
type FixedLoanLendingApyHistory struct {
	Rate      decimal.Decimal `json:"rate"`
	Timestamp time.Time       `json:"ts"`
}

// GetFixedLoanPendingLendingVolumeService -- GET /api/v5/finance/fixed-loan/pending-lending-volume (Read)
//
// Returns the pending (queued) lending volume for a currency/term.
type GetFixedLoanPendingLendingVolumeService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetFixedLoanPendingLendingVolumeService(ccy, term string) *GetFixedLoanPendingLendingVolumeService {
	return &GetFixedLoanPendingLendingVolumeService{c: c, params: map[string]string{
		"ccy":  ccy,
		"term": term,
	}}
}

func (s *GetFixedLoanPendingLendingVolumeService) Do(ctx context.Context) ([]FixedLoanPendingLendingVolume, error) {
	req := request.Get(ctx, s.c, "/api/v5/finance/fixed-loan/pending-lending-volume", s.params).WithSign()
	return request.DoList[FixedLoanPendingLendingVolume](req)
}

// FixedLoanPendingLendingVolume is the pending lending volume for a
// currency/term.
type FixedLoanPendingLendingVolume struct {
	Currency string          `json:"ccy"`
	Term     string          `json:"term"`
	Amount   decimal.Decimal `json:"amt"`
	MinRate  decimal.Decimal `json:"minRate"`
	MaxRate  decimal.Decimal `json:"maxRate"`
}

// PlaceFixedLoanLendingOrderService -- POST /api/v5/finance/fixed-loan/lending-order (Trade)
//
// Places a fixed-loan lending order. State-changing: implemented but never
// executed by the test suite. The acknowledgement shape is modeled from the OKX
// doc field table.
type PlaceFixedLoanLendingOrderService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewPlaceFixedLoanLendingOrderService(ccy string, amt, rate decimal.Decimal, term string) *PlaceFixedLoanLendingOrderService {
	return &PlaceFixedLoanLendingOrderService{c: c, body: map[string]any{
		"ccy":  ccy,
		"amt":  amt.String(),
		"rate": rate.String(),
		"term": term,
	}}
}

// SetAutoRenewal enables auto-renewal of the lending order on maturity.
func (s *PlaceFixedLoanLendingOrderService) SetAutoRenewal(autoRenewal bool) *PlaceFixedLoanLendingOrderService {
	s.body["autoRenewal"] = autoRenewal
	return s
}

func (s *PlaceFixedLoanLendingOrderService) Do(ctx context.Context) (*FixedLoanLendingOrderAck, error) {
	req := request.Post(ctx, s.c, "/api/v5/finance/fixed-loan/lending-order", s.body).WithSign()
	return request.DoOne[FixedLoanLendingOrderAck](req)
}

// FixedLoanLendingOrderAck is the acknowledgement of a placed lending order.
type FixedLoanLendingOrderAck struct {
	OrderID string `json:"ordId"`
}

// AmendFixedLoanLendingOrderService -- POST /api/v5/finance/fixed-loan/amend-lending-order (Trade)
//
// Amends an existing fixed-loan lending order. State-changing: implemented but
// never executed by the test suite. The acknowledgement shape is modeled from
// the OKX doc field table.
type AmendFixedLoanLendingOrderService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewAmendFixedLoanLendingOrderService(ordId string) *AmendFixedLoanLendingOrderService {
	return &AmendFixedLoanLendingOrderService{c: c, body: map[string]any{"ordId": ordId}}
}

// SetChangeAmt sets the amount to add to (or reduce from) the order.
func (s *AmendFixedLoanLendingOrderService) SetChangeAmt(changeAmt decimal.Decimal) *AmendFixedLoanLendingOrderService {
	s.body["changeAmt"] = changeAmt.String()
	return s
}

// SetRate sets a new lending rate for the order.
func (s *AmendFixedLoanLendingOrderService) SetRate(rate decimal.Decimal) *AmendFixedLoanLendingOrderService {
	s.body["rate"] = rate.String()
	return s
}

// SetAutoRenewal toggles auto-renewal of the order.
func (s *AmendFixedLoanLendingOrderService) SetAutoRenewal(autoRenewal bool) *AmendFixedLoanLendingOrderService {
	s.body["autoRenewal"] = autoRenewal
	return s
}

func (s *AmendFixedLoanLendingOrderService) Do(ctx context.Context) (*FixedLoanAmendLendingOrderAck, error) {
	req := request.Post(ctx, s.c, "/api/v5/finance/fixed-loan/amend-lending-order", s.body).WithSign()
	return request.DoOne[FixedLoanAmendLendingOrderAck](req)
}

// FixedLoanAmendLendingOrderAck is the acknowledgement of an amended lending
// order.
type FixedLoanAmendLendingOrderAck struct {
	OrderID string `json:"ordId"`
}

// GetFixedLoanLendingOrdersListService -- GET /api/v5/finance/fixed-loan/lending-orders-list (Read)
//
// Returns the account's fixed-loan lending orders.
type GetFixedLoanLendingOrdersListService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetFixedLoanLendingOrdersListService() *GetFixedLoanLendingOrdersListService {
	return &GetFixedLoanLendingOrdersListService{c: c, params: map[string]string{}}
}

// SetOrdId filters by order id.
func (s *GetFixedLoanLendingOrdersListService) SetOrdId(ordId string) *GetFixedLoanLendingOrdersListService {
	s.params["ordId"] = ordId
	return s
}

// SetCcy filters by lending currency.
func (s *GetFixedLoanLendingOrdersListService) SetCcy(ccy string) *GetFixedLoanLendingOrdersListService {
	s.params["ccy"] = ccy
	return s
}

// SetState filters by order state.
func (s *GetFixedLoanLendingOrdersListService) SetState(state string) *GetFixedLoanLendingOrdersListService {
	s.params["state"] = state
	return s
}

// SetTerm filters by lending term.
func (s *GetFixedLoanLendingOrdersListService) SetTerm(term string) *GetFixedLoanLendingOrdersListService {
	s.params["term"] = term
	return s
}

// SetAfter paginates to records earlier than the given time (older).
func (s *GetFixedLoanLendingOrdersListService) SetAfter(t time.Time) *GetFixedLoanLendingOrdersListService {
	s.params["after"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetBefore paginates to records later than the given time (newer).
func (s *GetFixedLoanLendingOrdersListService) SetBefore(t time.Time) *GetFixedLoanLendingOrdersListService {
	s.params["before"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetFixedLoanLendingOrdersListService) SetLimit(limit int) *GetFixedLoanLendingOrdersListService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetFixedLoanLendingOrdersListService) Do(ctx context.Context) ([]FixedLoanLendingOrder, error) {
	req := request.Get(ctx, s.c, "/api/v5/finance/fixed-loan/lending-orders-list", s.params).WithSign()
	return request.DoList[FixedLoanLendingOrder](req)
}

// FixedLoanLendingOrder is one fixed-loan lending order.
type FixedLoanLendingOrder struct {
	OrderID         string          `json:"ordId"`
	Currency        string          `json:"ccy"`
	Amount          decimal.Decimal `json:"amt"`
	Term            string          `json:"term"`
	Rate            decimal.Decimal `json:"rate"`
	EarningCurrency string          `json:"earningCcy"`
	Earnings        decimal.Decimal `json:"earnings"`
	State           string          `json:"state"`
	AutoRenewal     bool            `json:"autoRenewal"`
	SettledRate     decimal.Decimal `json:"settledRate"`
	PendingAmount   decimal.Decimal `json:"pendingAmt"`
	UsedAmount      decimal.Decimal `json:"usedAmt"`
	CreationTime    time.Time       `json:"cTime"`
	UpdateTime      time.Time       `json:"uTime"`
}

// GetFixedLoanLendingSubOrdersService -- GET /api/v5/finance/fixed-loan/lending-sub-orders (Read)
//
// Returns the sub-orders of a fixed-loan lending order.
type GetFixedLoanLendingSubOrdersService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetFixedLoanLendingSubOrdersService(ordId string) *GetFixedLoanLendingSubOrdersService {
	return &GetFixedLoanLendingSubOrdersService{c: c, params: map[string]string{"ordId": ordId}}
}

func (s *GetFixedLoanLendingSubOrdersService) Do(ctx context.Context) ([]FixedLoanLendingSubOrder, error) {
	req := request.Get(ctx, s.c, "/api/v5/finance/fixed-loan/lending-sub-orders", s.params).WithSign()
	return request.DoList[FixedLoanLendingSubOrder](req)
}

// FixedLoanLendingSubOrder is one sub-order of a fixed-loan lending order.
type FixedLoanLendingSubOrder struct {
	OrderID            string          `json:"ordId"`
	SubOrderID         string          `json:"subOrdId"`
	Currency           string          `json:"ccy"`
	Amount             decimal.Decimal `json:"amt"`
	Term               string          `json:"term"`
	Rate               decimal.Decimal `json:"rate"`
	EarningCurrency    string          `json:"earningCcy"`
	Earnings           decimal.Decimal `json:"earnings"`
	State              string          `json:"state"`
	SettledRate        decimal.Decimal `json:"settledRate"`
	CreationTime       time.Time       `json:"cTime"`
	UpdateTime         time.Time       `json:"uTime"`
	EffectiveTimestamp time.Time       `json:"effectiveTs"`
	ExpiryTimestamp    time.Time       `json:"expiryTs"`
}
