package okx

import (
	"context"
	"time"

	"github.com/UnipayFI/go-okx/request"
	"github.com/shopspring/decimal"
)

// GetAffiliateInviteeDetailService -- GET /api/v5/affiliate/invitee/detail (Read)
//
// Returns the affiliate-relationship detail for one of the calling affiliate's
// invitees, identified by the invitee's user id (uid). Only accounts with the
// affiliate/agent role may call this; other accounts get OKX code 51620 ("Only
// affiliates can perform this action").
//
// Path is curl-verified live: an account without the affiliate role receives
// code 51620 (path REAL, capability-gated). The other affiliate paths probed
// (performance-summary, invitee-list, co-inviter-link-list, link-list,
// sub-affiliate-list) return HTTP 404 — they are not part of the public OKX v5
// REST API and were therefore dropped.
type GetAffiliateInviteeDetailService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetAffiliateInviteeDetailService(uid string) *GetAffiliateInviteeDetailService {
	return &GetAffiliateInviteeDetailService{c: c, params: map[string]string{"uid": uid}}
}

func (s *GetAffiliateInviteeDetailService) Do(ctx context.Context) (*AffiliateInviteeDetail, error) {
	req := request.Get(ctx, s.c, "/api/v5/affiliate/invitee/detail", s.params).WithSign()
	return request.DoOne[AffiliateInviteeDetail](req)
}

// AffiliateInviteeDetail is one invitee's affiliate-relationship detail. The
// validating account does not have the affiliate role (the endpoint returns code
// 51620), so the field set is modeled from the OKX affiliate doc field table.
type AffiliateInviteeDetail struct {
	InviteeLevel                string          `json:"inviteeLv"`
	JoinTime                    time.Time       `json:"joinTime"`
	InviteeRebateRate           decimal.Decimal `json:"inviteeRebateRate"`
	TotalCommission             decimal.Decimal `json:"totalCommission"`
	FirstTradeTime              time.Time       `json:"firstTradeTime"`
	Level                       string          `json:"level"`
	DepositAmount               decimal.Decimal `json:"depAmt"`
	Volume                      decimal.Decimal `json:"vol"`
	KYCTime                     time.Time       `json:"kycTime"`
	Region                      string          `json:"region"`
	AffiliateCode               string          `json:"affiliateCode"`
	InvitedTradeVolumeThirtyDay decimal.Decimal `json:"invitedTradeVolThirtyD"`
}
