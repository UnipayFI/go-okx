package okx

import (
	"context"
	"time"

	"github.com/UnipayFI/go-okx/request"
	"github.com/shopspring/decimal"
)

// AffiliatePeriodType selects the stats window for an affiliate endpoint's
// period-scoped fields (volPeriod).
type AffiliatePeriodType string

const (
	AffiliatePeriodLast7D    AffiliatePeriodType = "last_7d"
	AffiliatePeriodLast30D   AffiliatePeriodType = "last_30d"
	AffiliatePeriodThisMonth AffiliatePeriodType = "this_month"
	AffiliatePeriodLastMonth AffiliatePeriodType = "last_month"
	AffiliatePeriodTotal     AffiliatePeriodType = "total"
	AffiliatePeriodToday     AffiliatePeriodType = "today"
	AffiliatePeriodThisWeek  AffiliatePeriodType = "this_week"
	// AffiliatePeriodCustom pairs with begin/end and is accepted by the invitee
	// list only; the invitee detail rejects it with code 51000.
	AffiliatePeriodCustom AffiliatePeriodType = "custom"
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

// SetPeriodType selects the stats window for the VolumePeriod response field.
// When unset, VolumePeriod is not returned. AffiliatePeriodCustom is not
// supported here — it (or any unknown value) returns code 51000.
func (s *GetAffiliateInviteeDetailService) SetPeriodType(periodType AffiliatePeriodType) *GetAffiliateInviteeDetailService {
	s.params["periodType"] = string(periodType)
	return s
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
	// VolumePeriod is the trading volume inside the SetPeriodType window, in
	// USDT. Only returned when periodType is supplied; 0 when the invitee did
	// not trade in the window.
	VolumePeriod decimal.Decimal `json:"volPeriod"`
}
