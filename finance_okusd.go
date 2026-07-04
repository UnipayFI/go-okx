package okx

import (
	"context"
	"time"

	"github.com/UnipayFI/go-okx/request"
	"github.com/shopspring/decimal"
)

// OKUSDRedeemType selects how an OKUSD redemption settles.
type OKUSDRedeemType string

const (
	// OKUSDRedeemTypeFast settles in real time.
	OKUSDRedeemTypeFast OKUSDRedeemType = "1"
	// OKUSDRedeemTypeStandard settles in D+5 business days.
	OKUSDRedeemTypeStandard OKUSDRedeemType = "2"
)

// GetOKUSDLimitsService -- GET /api/v5/finance/okusd/limits (Read)
//
// Returns the account's remaining daily OKUSD subscription and redemption
// limits (personal VIP-tier limits and platform-wide limits, plus redemption
// fee rates).
type GetOKUSDLimitsService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetOKUSDLimitsService() *GetOKUSDLimitsService {
	return &GetOKUSDLimitsService{c: c, params: map[string]string{}}
}

func (s *GetOKUSDLimitsService) Do(ctx context.Context) (*OKUSDLimits, error) {
	req := request.Get(ctx, s.c, "/api/v5/finance/okusd/limits", s.params).WithSign()
	return request.DoOne[OKUSDLimits](req)
}

// OKUSDLimits is the account's OKUSD subscription/redemption limit snapshot.
type OKUSDLimits struct {
	SubLimit        OKUSDSubLimit    `json:"subLimit"`
	FastRedeemLimit OKUSDRedeemLimit `json:"fastRedeemLimit"`
	StdRedeemLimit  OKUSDRedeemLimit `json:"stdRedeemLimit"`
	Timestamp       time.Time        `json:"ts"`
}

// OKUSDSubLimit holds the subscription limit figures (USDT).
type OKUSDSubLimit struct {
	MaxSubAmount       decimal.Decimal `json:"maxSubAmt"`
	PersonalDailyLimit decimal.Decimal `json:"personalDailyLimit"`
	PersonalUsedAmount decimal.Decimal `json:"personalUsedAmt"`
	PlatformDailyLimit decimal.Decimal `json:"platformDailyLimit"`
	PlatformUsedAmount decimal.Decimal `json:"platformUsedAmt"`
}

// OKUSDRedeemLimit holds a redemption limit tier's figures (OKUSD) and its fee
// rate.
type OKUSDRedeemLimit struct {
	PersonalDailyLimit decimal.Decimal `json:"personalDailyLimit"`
	PersonalUsedAmount decimal.Decimal `json:"personalUsedAmt"`
	PlatformDailyLimit decimal.Decimal `json:"platformDailyLimit"`
	PlatformUsedAmount decimal.Decimal `json:"platformUsedAmt"`
	FeeRate            decimal.Decimal `json:"feeRate"`
}

// SetOKUSDSubscribeService -- POST /api/v5/finance/okusd/subscribe (Trade)
//
// Subscribes USDT to receive OKUSD at a 1:1 rate.
//
// State-changing: NOT exercised by the test suite.
type SetOKUSDSubscribeService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewSetOKUSDSubscribeService(amt decimal.Decimal) *SetOKUSDSubscribeService {
	return &SetOKUSDSubscribeService{c: c, body: map[string]any{
		"amt": amt.String(),
	}}
}

func (s *SetOKUSDSubscribeService) Do(ctx context.Context) (*OKUSDSubscription, error) {
	req := request.Post(ctx, s.c, "/api/v5/finance/okusd/subscribe", s.body).WithSign()
	return request.DoOne[OKUSDSubscription](req)
}

// OKUSDSubscription is the ack of an OKUSD subscription.
type OKUSDSubscription struct {
	OrderID     string          `json:"ordId"`
	Currency    string          `json:"ccy"`
	Amount      decimal.Decimal `json:"amt"`
	OKUSDAmount decimal.Decimal `json:"okusdAmt"`
	State       string          `json:"state"`
	Timestamp   time.Time       `json:"ts"`
}

// SetOKUSDRedeemService -- POST /api/v5/finance/okusd/redeem (Trade)
//
// Redeems OKUSD back to USDT, either fast (real-time) or standard (D+5).
//
// State-changing: NOT exercised by the test suite.
type SetOKUSDRedeemService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewSetOKUSDRedeemService(amt decimal.Decimal, redeemType OKUSDRedeemType) *SetOKUSDRedeemService {
	return &SetOKUSDRedeemService{c: c, body: map[string]any{
		"amt":        amt.String(),
		"redeemType": string(redeemType),
	}}
}

func (s *SetOKUSDRedeemService) Do(ctx context.Context) (*OKUSDRedemption, error) {
	req := request.Post(ctx, s.c, "/api/v5/finance/okusd/redeem", s.body).WithSign()
	return request.DoOne[OKUSDRedemption](req)
}

// OKUSDRedemption is the ack of an OKUSD redemption.
type OKUSDRedemption struct {
	OrderID                 string          `json:"ordId"`
	Currency                string          `json:"ccy"`
	Amount                  decimal.Decimal `json:"amt"`
	Fee                     decimal.Decimal `json:"fee"`
	USDTAmount              decimal.Decimal `json:"usdtAmt"`
	RedeemType              OKUSDRedeemType `json:"redeemType"`
	State                   string          `json:"state"`
	EstimatedSettlementTime time.Time       `json:"estSettlementTime"`
	Timestamp               time.Time       `json:"ts"`
}
