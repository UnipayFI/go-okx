package okx

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-okx/request"
	"github.com/shopspring/decimal"
)

// StakingProtocolType is the on-chain earn protocol category (e.g. "defi").
type StakingProtocolType string

const (
	StakingProtocolTypeDefi    StakingProtocolType = "defi"
	StakingProtocolTypeStaking StakingProtocolType = "staking"
)

// StakingOfferState is the purchasability state of an on-chain earn offer.
type StakingOfferState string

const (
	StakingOfferStatePurchasable StakingOfferState = "purchasable"
	StakingOfferStateSoldOut     StakingOfferState = "sold_out"
	StakingOfferStateStop        StakingOfferState = "stop"
)

// StakingOrderState is the lifecycle state of an on-chain earn order.
type StakingOrderState string

const (
	StakingOrderState8  StakingOrderState = "8"  // pending
	StakingOrderState13 StakingOrderState = "13" // cancelling
	StakingOrderState9  StakingOrderState = "9"  // onchain
	StakingOrderState1  StakingOrderState = "1"  // earning
	StakingOrderState2  StakingOrderState = "2"  // redeeming
	StakingOrderState3  StakingOrderState = "3"  // expired
)

// GetStakingOffersService -- GET /api/v5/finance/staking-defi/offers (Read)
//
// Returns the available on-chain earn (staking/defi) offers, including the
// per-currency invest constraints and the assets earned.
type GetStakingOffersService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetStakingOffersService() *GetStakingOffersService {
	return &GetStakingOffersService{c: c, params: map[string]string{}}
}

// SetProductId filters by a single product id.
func (s *GetStakingOffersService) SetProductId(productId string) *GetStakingOffersService {
	s.params["productId"] = productId
	return s
}

// SetProtocolType filters by protocol category (defi/staking).
func (s *GetStakingOffersService) SetProtocolType(protocolType StakingProtocolType) *GetStakingOffersService {
	s.params["protocolType"] = string(protocolType)
	return s
}

// SetCcy filters by investment currency.
func (s *GetStakingOffersService) SetCcy(ccy string) *GetStakingOffersService {
	s.params["ccy"] = ccy
	return s
}

func (s *GetStakingOffersService) Do(ctx context.Context) ([]StakingOffer, error) {
	req := request.Get(ctx, s.c, "/api/v5/finance/staking-defi/offers", s.params).WithSign()
	return request.DoList[StakingOffer](req)
}

// StakingOffer is a single on-chain earn offer.
type StakingOffer struct {
	Currency                 string               `json:"ccy"`
	ProductID                string               `json:"productId"`
	Protocol                 string               `json:"protocol"`
	ProtocolType             StakingProtocolType  `json:"protocolType"`
	Term                     string               `json:"term"`
	APY                      decimal.Decimal      `json:"apy"`
	EarlyRedeem              bool                 `json:"earlyRedeem"`
	State                    StakingOfferState    `json:"state"`
	InvestmentData           []StakingInvestData  `json:"investData"`
	EarningData              []StakingEarningData `json:"earningData"`
	FastRedemptionDailyLimit decimal.Decimal      `json:"fastRedemptionDailyLimit"`
	ProjectDisplayName       string               `json:"projectDisplayName"`
	RedeemPeriod             []string             `json:"redeemPeriod"`
}

// StakingInvestData is one investable currency's balance and amount limits
// within an offer.
type StakingInvestData struct {
	Balance   decimal.Decimal `json:"bal"`
	Currency  string          `json:"ccy"`
	MaxAmount decimal.Decimal `json:"maxAmt"`
	MinAmount decimal.Decimal `json:"minAmt"`
}

// StakingEarningData is one asset earned by an offer and how it is earned.
type StakingEarningData struct {
	Currency    string `json:"ccy"`
	EarningType string `json:"earningType"`
}

// PurchaseStakingService -- POST /api/v5/finance/staking-defi/purchase (Trade)
//
// Subscribes to an on-chain earn offer. State-changing: implement-only, never
// executed by the SDK's test suite.
type PurchaseStakingService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewPurchaseStakingService(productId string, investData []StakingPurchaseInvest) *PurchaseStakingService {
	return &PurchaseStakingService{c: c, body: map[string]any{
		"productId":  productId,
		"investData": investData,
	}}
}

// SetTerm sets the protocol term (required for fixed-term protocols).
func (s *PurchaseStakingService) SetTerm(term string) *PurchaseStakingService {
	s.body["term"] = term
	return s
}

// SetTag sets an order tag (brokerId-issued label).
func (s *PurchaseStakingService) SetTag(tag string) *PurchaseStakingService {
	s.body["tag"] = tag
	return s
}

func (s *PurchaseStakingService) Do(ctx context.Context) (*StakingOrderAck, error) {
	req := request.Post(ctx, s.c, "/api/v5/finance/staking-defi/purchase", s.body).WithSign()
	return request.DoOne[StakingOrderAck](req)
}

// StakingPurchaseInvest is one currency+amount line of a purchase request.
type StakingPurchaseInvest struct {
	Currency string `json:"ccy"`
	Amount   string `json:"amt"`
}

// StakingOrderAck is the acknowledgement returned by purchase/redeem/cancel.
type StakingOrderAck struct {
	OrderID string `json:"ordId"`
	Tag     string `json:"tag"`
}

// RedeemStakingService -- POST /api/v5/finance/staking-defi/redeem (Trade)
//
// Redeems an active on-chain earn order. State-changing: implement-only.
type RedeemStakingService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewRedeemStakingService(ordId string, protocolType StakingProtocolType) *RedeemStakingService {
	return &RedeemStakingService{c: c, body: map[string]any{
		"ordId":        ordId,
		"protocolType": string(protocolType),
	}}
}

// SetAllowEarlyRedeem permits early redemption (incurring any applicable cost).
func (s *RedeemStakingService) SetAllowEarlyRedeem(allow bool) *RedeemStakingService {
	s.body["allowEarlyRedeem"] = allow
	return s
}

func (s *RedeemStakingService) Do(ctx context.Context) (*StakingOrderAck, error) {
	req := request.Post(ctx, s.c, "/api/v5/finance/staking-defi/redeem", s.body).WithSign()
	return request.DoOne[StakingOrderAck](req)
}

// CancelStakingService -- POST /api/v5/finance/staking-defi/cancel (Trade)
//
// Cancels a pending purchase/redemption of an on-chain earn order.
// State-changing: implement-only.
type CancelStakingService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewCancelStakingService(ordId string, protocolType StakingProtocolType) *CancelStakingService {
	return &CancelStakingService{c: c, body: map[string]any{
		"ordId":        ordId,
		"protocolType": string(protocolType),
	}}
}

func (s *CancelStakingService) Do(ctx context.Context) (*StakingOrderAck, error) {
	req := request.Post(ctx, s.c, "/api/v5/finance/staking-defi/cancel", s.body).WithSign()
	return request.DoOne[StakingOrderAck](req)
}

// GetStakingActiveOrdersService -- GET /api/v5/finance/staking-defi/orders-active (Read)
//
// Returns the account's active (not yet fully redeemed) on-chain earn orders.
type GetStakingActiveOrdersService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetStakingActiveOrdersService() *GetStakingActiveOrdersService {
	return &GetStakingActiveOrdersService{c: c, params: map[string]string{}}
}

// SetProductId filters by product id.
func (s *GetStakingActiveOrdersService) SetProductId(productId string) *GetStakingActiveOrdersService {
	s.params["productId"] = productId
	return s
}

// SetProtocolType filters by protocol category (defi/staking).
func (s *GetStakingActiveOrdersService) SetProtocolType(protocolType StakingProtocolType) *GetStakingActiveOrdersService {
	s.params["protocolType"] = string(protocolType)
	return s
}

// SetCcy filters by investment currency.
func (s *GetStakingActiveOrdersService) SetCcy(ccy string) *GetStakingActiveOrdersService {
	s.params["ccy"] = ccy
	return s
}

// SetState filters by order state.
func (s *GetStakingActiveOrdersService) SetState(state StakingOrderState) *GetStakingActiveOrdersService {
	s.params["state"] = string(state)
	return s
}

func (s *GetStakingActiveOrdersService) Do(ctx context.Context) ([]StakingActiveOrder, error) {
	req := request.Get(ctx, s.c, "/api/v5/finance/staking-defi/orders-active", s.params).WithSign()
	return request.DoList[StakingActiveOrder](req)
}

// StakingActiveOrder is one active on-chain earn order. The validating account
// holds no on-chain earn orders, so the field set is modeled from the OKX doc
// field table.
type StakingActiveOrder struct {
	OrderID                  string                  `json:"ordId"`
	Currency                 string                  `json:"ccy"`
	ProductID                string                  `json:"productId"`
	State                    StakingOrderState       `json:"state"`
	Protocol                 string                  `json:"protocol"`
	ProtocolType             StakingProtocolType     `json:"protocolType"`
	Term                     string                  `json:"term"`
	APY                      decimal.Decimal         `json:"apy"`
	InvestmentData           []StakingOrderInvest    `json:"investData"`
	EarningData              []StakingOrderEarn      `json:"earningData"`
	PurchasedTime            time.Time               `json:"purchasedTime"`
	EstimatedSettlementTime  time.Time               `json:"estSettlementTime"`
	CancelRedemptionDeadline time.Time               `json:"cancelRedemptionDeadline"`
	FastRedemptionData       []StakingFastRedemption `json:"fastRedemptionData"`
	Tag                      string                  `json:"tag"`
}

// StakingOrderInvest is one invested currency line of an on-chain earn order.
type StakingOrderInvest struct {
	Currency string          `json:"ccy"`
	Amount   decimal.Decimal `json:"amt"`
}

// StakingOrderEarn is one earned-asset line of an on-chain earn order.
type StakingOrderEarn struct {
	Currency         string          `json:"ccy"`
	Earnings         decimal.Decimal `json:"earnings"`
	EarningType      string          `json:"earningType"`
	RealizedEarnings decimal.Decimal `json:"realizedEarnings"`
}

// StakingFastRedemption is the fast-redemption availability of an order.
type StakingFastRedemption struct {
	Currency                string          `json:"ccy"`
	RedeemingAmount         decimal.Decimal `json:"redeemingAmt"`
	FastRedemptionAvailable decimal.Decimal `json:"fastRedemptionAvail"`
}

// GetStakingOrdersHistoryService -- GET /api/v5/finance/staking-defi/orders-history (Read)
//
// Returns the account's completed/redeemed on-chain earn orders.
type GetStakingOrdersHistoryService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetStakingOrdersHistoryService() *GetStakingOrdersHistoryService {
	return &GetStakingOrdersHistoryService{c: c, params: map[string]string{}}
}

// SetProductId filters by product id.
func (s *GetStakingOrdersHistoryService) SetProductId(productId string) *GetStakingOrdersHistoryService {
	s.params["productId"] = productId
	return s
}

// SetProtocolType filters by protocol category (defi/staking).
func (s *GetStakingOrdersHistoryService) SetProtocolType(protocolType StakingProtocolType) *GetStakingOrdersHistoryService {
	s.params["protocolType"] = string(protocolType)
	return s
}

// SetCcy filters by investment currency.
func (s *GetStakingOrdersHistoryService) SetCcy(ccy string) *GetStakingOrdersHistoryService {
	s.params["ccy"] = ccy
	return s
}

// SetAfter paginates to records earlier than the given order-creation time (older).
func (s *GetStakingOrdersHistoryService) SetAfter(t time.Time) *GetStakingOrdersHistoryService {
	s.params["after"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetBefore paginates to records later than the given order-creation time (newer).
func (s *GetStakingOrdersHistoryService) SetBefore(t time.Time) *GetStakingOrdersHistoryService {
	s.params["before"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetStakingOrdersHistoryService) SetLimit(limit int) *GetStakingOrdersHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetStakingOrdersHistoryService) Do(ctx context.Context) ([]StakingHistoryOrder, error) {
	req := request.Get(ctx, s.c, "/api/v5/finance/staking-defi/orders-history", s.params).WithSign()
	return request.DoList[StakingHistoryOrder](req)
}

// StakingHistoryOrder is one completed on-chain earn order. The validating
// account holds no order history, so the field set is modeled from the OKX doc
// field table.
type StakingHistoryOrder struct {
	OrderID        string               `json:"ordId"`
	Currency       string               `json:"ccy"`
	ProductID      string               `json:"productId"`
	State          StakingOrderState    `json:"state"`
	Protocol       string               `json:"protocol"`
	ProtocolType   StakingProtocolType  `json:"protocolType"`
	Term           string               `json:"term"`
	APY            decimal.Decimal      `json:"apy"`
	InvestmentData []StakingOrderInvest `json:"investData"`
	EarningData    []StakingOrderEarn   `json:"earningData"`
	PurchasedTime  time.Time            `json:"purchasedTime"`
	RedeemedTime   time.Time            `json:"redeemedTime"`
	Tag            string               `json:"tag"`
}

// GetEthStakingProductInfoService -- GET /api/v5/finance/staking-defi/eth/product-info (Read)
//
// Returns the ETH staking product parameters (rate, min amount, redemption
// window and fast-redemption daily limit).
type GetEthStakingProductInfoService struct {
	c *Client
}

func (c *Client) NewGetEthStakingProductInfoService() *GetEthStakingProductInfoService {
	return &GetEthStakingProductInfoService{c: c}
}

func (s *GetEthStakingProductInfoService) Do(ctx context.Context) (*EthStakingProductInfo, error) {
	req := request.Get(ctx, s.c, "/api/v5/finance/staking-defi/eth/product-info").WithSign()
	return request.DoOne[EthStakingProductInfo](req)
}

// EthStakingProductInfo is the ETH staking product configuration.
type EthStakingProductInfo struct {
	FastRedemptionDailyLimit decimal.Decimal `json:"fastRedemptionDailyLimit"`
	MinAmount                decimal.Decimal `json:"minAmt"`
	Rate                     decimal.Decimal `json:"rate"`
	RedemptionDays           decimal.Decimal `json:"redemptDays"`
}

// PurchaseEthStakingService -- POST /api/v5/finance/staking-defi/eth/purchase (Trade)
//
// Stakes ETH (receiving BETH). State-changing: implement-only.
type PurchaseEthStakingService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewPurchaseEthStakingService(amt decimal.Decimal) *PurchaseEthStakingService {
	return &PurchaseEthStakingService{c: c, body: map[string]any{"amt": amt.String()}}
}

func (s *PurchaseEthStakingService) Do(ctx context.Context) (*EthStakingAck, error) {
	req := request.Post(ctx, s.c, "/api/v5/finance/staking-defi/eth/purchase", s.body).WithSign()
	return request.DoOne[EthStakingAck](req)
}

// EthStakingAck is the (empty) acknowledgement returned by ETH purchase/redeem.
type EthStakingAck struct{}

// RedeemEthStakingService -- POST /api/v5/finance/staking-defi/eth/redeem (Trade)
//
// Redeems BETH back to ETH. State-changing: implement-only.
type RedeemEthStakingService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewRedeemEthStakingService(amt decimal.Decimal) *RedeemEthStakingService {
	return &RedeemEthStakingService{c: c, body: map[string]any{"amt": amt.String()}}
}

func (s *RedeemEthStakingService) Do(ctx context.Context) (*EthStakingAck, error) {
	req := request.Post(ctx, s.c, "/api/v5/finance/staking-defi/eth/redeem", s.body).WithSign()
	return request.DoOne[EthStakingAck](req)
}

// GetEthStakingBalanceService -- GET /api/v5/finance/staking-defi/eth/balance (Read)
//
// Returns the account's BETH balance and accrued staking interest.
type GetEthStakingBalanceService struct {
	c *Client
}

func (c *Client) NewGetEthStakingBalanceService() *GetEthStakingBalanceService {
	return &GetEthStakingBalanceService{c: c}
}

func (s *GetEthStakingBalanceService) Do(ctx context.Context) ([]EthStakingBalance, error) {
	req := request.Get(ctx, s.c, "/api/v5/finance/staking-defi/eth/balance").WithSign()
	return request.DoList[EthStakingBalance](req)
}

// EthStakingBalance is the BETH balance and accrued interest.
type EthStakingBalance struct {
	Currency              string          `json:"ccy"`
	Amount                decimal.Decimal `json:"amt"`
	LatestInterestAccrual decimal.Decimal `json:"latestInterestAccrual"`
	TotalInterestAccrual  decimal.Decimal `json:"totalInterestAccrual"`
	Timestamp             time.Time       `json:"ts"`
}

// GetEthStakingHistoryService -- GET /api/v5/finance/staking-defi/eth/purchase-redeem-history (Read)
//
// Returns the account's ETH staking purchase and redemption history.
type GetEthStakingHistoryService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetEthStakingHistoryService() *GetEthStakingHistoryService {
	return &GetEthStakingHistoryService{c: c, params: map[string]string{}}
}

// SetType filters by record type (purchase/redeem).
func (s *GetEthStakingHistoryService) SetType(typ string) *GetEthStakingHistoryService {
	s.params["type"] = typ
	return s
}

// SetStatus filters by record status (pending/success/failed).
func (s *GetEthStakingHistoryService) SetStatus(status string) *GetEthStakingHistoryService {
	s.params["status"] = status
	return s
}

// SetAfter paginates to records earlier than the given request time (older).
func (s *GetEthStakingHistoryService) SetAfter(t time.Time) *GetEthStakingHistoryService {
	s.params["after"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetBefore paginates to records later than the given request time (newer).
func (s *GetEthStakingHistoryService) SetBefore(t time.Time) *GetEthStakingHistoryService {
	s.params["before"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetEthStakingHistoryService) SetLimit(limit int) *GetEthStakingHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetEthStakingHistoryService) Do(ctx context.Context) ([]EthStakingHistory, error) {
	req := request.Get(ctx, s.c, "/api/v5/finance/staking-defi/eth/purchase-redeem-history", s.params).WithSign()
	return request.DoList[EthStakingHistory](req)
}

// EthStakingHistory is one ETH staking purchase/redemption record. The
// validating account has no ETH staking history, so the field set is modeled
// from the OKX doc field table.
type EthStakingHistory struct {
	Type                   string          `json:"type"`
	Amount                 decimal.Decimal `json:"amt"`
	Status                 string          `json:"status"`
	RequestTime            time.Time       `json:"requestTime"`
	CompletedTime          time.Time       `json:"completedTime"`
	EstimatedCompletedTime time.Time       `json:"estCompletedTime"`
}

// GetEthStakingApyHistoryService -- GET /api/v5/finance/staking-defi/eth/apy-history (Read)
//
// Returns the ETH staking APY history for the trailing number of days.
type GetEthStakingApyHistoryService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetEthStakingApyHistoryService(days int) *GetEthStakingApyHistoryService {
	return &GetEthStakingApyHistoryService{c: c, params: map[string]string{"days": strconv.Itoa(days)}}
}

func (s *GetEthStakingApyHistoryService) Do(ctx context.Context) ([]StakingApyHistory, error) {
	req := request.Get(ctx, s.c, "/api/v5/finance/staking-defi/eth/apy-history", s.params).WithSign()
	return request.DoList[StakingApyHistory](req)
}

// StakingApyHistory is one daily APY point of the ETH/SOL staking APY history.
type StakingApyHistory struct {
	Rate      decimal.Decimal `json:"rate"`
	Timestamp time.Time       `json:"ts"`
}

// GetSolStakingProductInfoService -- GET /api/v5/finance/staking-defi/sol/product-info (Read)
//
// Returns the SOL staking product parameters. Unlike the ETH variant, OKX
// returns this endpoint's "data" as a single JSON object, so it decodes via
// DoObject.
type GetSolStakingProductInfoService struct {
	c *Client
}

func (c *Client) NewGetSolStakingProductInfoService() *GetSolStakingProductInfoService {
	return &GetSolStakingProductInfoService{c: c}
}

func (s *GetSolStakingProductInfoService) Do(ctx context.Context) (*SolStakingProductInfo, error) {
	req := request.Get(ctx, s.c, "/api/v5/finance/staking-defi/sol/product-info").WithSign()
	return request.DoObject[SolStakingProductInfo](req)
}

// SolStakingProductInfo is the SOL staking product configuration.
type SolStakingProductInfo struct {
	FastRedemptionAvailable  decimal.Decimal `json:"fastRedemptionAvail"`
	FastRedemptionDailyLimit decimal.Decimal `json:"fastRedemptionDailyLimit"`
	MinAmount                decimal.Decimal `json:"minAmt"`
	Rate                     decimal.Decimal `json:"rate"`
	RedemptionDays           decimal.Decimal `json:"redemptDays"`
}

// PurchaseSolStakingService -- POST /api/v5/finance/staking-defi/sol/purchase (Trade)
//
// Stakes SOL (receiving OKSOL). State-changing: implement-only.
type PurchaseSolStakingService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewPurchaseSolStakingService(amt decimal.Decimal) *PurchaseSolStakingService {
	return &PurchaseSolStakingService{c: c, body: map[string]any{"amt": amt.String()}}
}

func (s *PurchaseSolStakingService) Do(ctx context.Context) (*SolStakingAck, error) {
	req := request.Post(ctx, s.c, "/api/v5/finance/staking-defi/sol/purchase", s.body).WithSign()
	return request.DoOne[SolStakingAck](req)
}

// SolStakingAck is the (empty) acknowledgement returned by SOL purchase/redeem.
type SolStakingAck struct{}

// RedeemSolStakingService -- POST /api/v5/finance/staking-defi/sol/redeem (Trade)
//
// Redeems OKSOL back to SOL. State-changing: implement-only.
type RedeemSolStakingService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewRedeemSolStakingService(amt decimal.Decimal) *RedeemSolStakingService {
	return &RedeemSolStakingService{c: c, body: map[string]any{"amt": amt.String()}}
}

func (s *RedeemSolStakingService) Do(ctx context.Context) (*SolStakingAck, error) {
	req := request.Post(ctx, s.c, "/api/v5/finance/staking-defi/sol/redeem", s.body).WithSign()
	return request.DoOne[SolStakingAck](req)
}

// GetSolStakingBalanceService -- GET /api/v5/finance/staking-defi/sol/balance (Read)
//
// Returns the account's OKSOL balance and accrued staking interest.
type GetSolStakingBalanceService struct {
	c *Client
}

func (c *Client) NewGetSolStakingBalanceService() *GetSolStakingBalanceService {
	return &GetSolStakingBalanceService{c: c}
}

func (s *GetSolStakingBalanceService) Do(ctx context.Context) ([]SolStakingBalance, error) {
	req := request.Get(ctx, s.c, "/api/v5/finance/staking-defi/sol/balance").WithSign()
	return request.DoList[SolStakingBalance](req)
}

// SolStakingBalance is the OKSOL balance and accrued interest.
type SolStakingBalance struct {
	Currency              string          `json:"ccy"`
	Amount                decimal.Decimal `json:"amt"`
	LatestInterestAccrual decimal.Decimal `json:"latestInterestAccrual"`
	TotalInterestAccrual  decimal.Decimal `json:"totalInterestAccrual"`
	Timestamp             time.Time       `json:"ts"`
}

// GetSolStakingHistoryService -- GET /api/v5/finance/staking-defi/sol/purchase-redeem-history (Read)
//
// Returns the account's SOL staking purchase and redemption history.
type GetSolStakingHistoryService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetSolStakingHistoryService() *GetSolStakingHistoryService {
	return &GetSolStakingHistoryService{c: c, params: map[string]string{}}
}

// SetType filters by record type (purchase/redeem).
func (s *GetSolStakingHistoryService) SetType(typ string) *GetSolStakingHistoryService {
	s.params["type"] = typ
	return s
}

// SetStatus filters by record status (pending/success/failed).
func (s *GetSolStakingHistoryService) SetStatus(status string) *GetSolStakingHistoryService {
	s.params["status"] = status
	return s
}

// SetAfter paginates to records earlier than the given request time (older).
func (s *GetSolStakingHistoryService) SetAfter(t time.Time) *GetSolStakingHistoryService {
	s.params["after"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetBefore paginates to records later than the given request time (newer).
func (s *GetSolStakingHistoryService) SetBefore(t time.Time) *GetSolStakingHistoryService {
	s.params["before"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetSolStakingHistoryService) SetLimit(limit int) *GetSolStakingHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetSolStakingHistoryService) Do(ctx context.Context) ([]SolStakingHistory, error) {
	req := request.Get(ctx, s.c, "/api/v5/finance/staking-defi/sol/purchase-redeem-history", s.params).WithSign()
	return request.DoList[SolStakingHistory](req)
}

// SolStakingHistory is one SOL staking purchase/redemption record. The
// validating account has no SOL staking history, so the field set is modeled
// from the OKX doc field table.
type SolStakingHistory struct {
	Type                   string          `json:"type"`
	Amount                 decimal.Decimal `json:"amt"`
	Status                 string          `json:"status"`
	RequestTime            time.Time       `json:"requestTime"`
	CompletedTime          time.Time       `json:"completedTime"`
	EstimatedCompletedTime time.Time       `json:"estCompletedTime"`
}

// GetSolStakingApyHistoryService -- GET /api/v5/finance/staking-defi/sol/apy-history (Read)
//
// Returns the SOL staking APY history for the trailing number of days.
type GetSolStakingApyHistoryService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetSolStakingApyHistoryService(days int) *GetSolStakingApyHistoryService {
	return &GetSolStakingApyHistoryService{c: c, params: map[string]string{"days": strconv.Itoa(days)}}
}

func (s *GetSolStakingApyHistoryService) Do(ctx context.Context) ([]StakingApyHistory, error) {
	req := request.Get(ctx, s.c, "/api/v5/finance/staking-defi/sol/apy-history", s.params).WithSign()
	return request.DoList[StakingApyHistory](req)
}
