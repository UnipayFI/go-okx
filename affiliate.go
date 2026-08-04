package okx

import (
	"context"
	"strconv"
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
// code 51620 (path REAL, capability-gated).
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

// AffiliateCommissionCategory is the commission calculation bucket an affiliate
// endpoint reports or filters on.
type AffiliateCommissionCategory string

const (
	AffiliateCommissionSpot       AffiliateCommissionCategory = "SPOT"
	AffiliateCommissionDerivative AffiliateCommissionCategory = "DERIVATIVE"
	AffiliateCommissionBSC        AffiliateCommissionCategory = "BSC"
)

// GetAffiliateInviteeListService -- GET /api/v5/affiliate/invitee/list (Read)
//
// Returns the calling affiliate's invitees, one page at a time, with per-invitee
// trading stats and KYC info. Like the invitee detail, it is affiliate-role
// gated: other accounts get OKX code 51620.
//
// The envelope carries a totalPage field alongside data; it is not surfaced by
// the typed helpers, so paginate with SetPage until a short page comes back.
type GetAffiliateInviteeListService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetAffiliateInviteeListService() *GetAffiliateInviteeListService {
	return &GetAffiliateInviteeListService{c: c, params: map[string]string{}}
}

// SetPage selects the 1-indexed page (default 1).
func (s *GetAffiliateInviteeListService) SetPage(page int) *GetAffiliateInviteeListService {
	s.params["page"] = strconv.Itoa(page)
	return s
}

// SetLimit sets the items per page, clamped server-side to [1, 100] (default 100).
func (s *GetAffiliateInviteeListService) SetLimit(limit int) *GetAffiliateInviteeListService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetPeriodType selects the stats window. AffiliatePeriodCustom requires both
// SetBegin and SetEnd; for every other value a server-defined window is used and
// any begin/end passed alongside is ignored.
func (s *GetAffiliateInviteeListService) SetPeriodType(periodType AffiliatePeriodType) *GetAffiliateInviteeListService {
	s.params["periodType"] = string(periodType)
	return s
}

// SetBegin sets the custom stats-window start (ms, inclusive). Required together
// with SetEnd when periodType is custom; the span may not exceed 90 days and may
// not start earlier than 180 days ago.
func (s *GetAffiliateInviteeListService) SetBegin(ms int64) *GetAffiliateInviteeListService {
	s.params["begin"] = strconv.FormatInt(ms, 10)
	return s
}

// SetEnd sets the custom stats-window end (ms, inclusive).
func (s *GetAffiliateInviteeListService) SetEnd(ms int64) *GetAffiliateInviteeListService {
	s.params["end"] = strconv.FormatInt(ms, 10)
	return s
}

// SetKeyword searches by UID or channel name.
func (s *GetAffiliateInviteeListService) SetKeyword(keyword string) *GetAffiliateInviteeListService {
	s.params["keyword"] = keyword
	return s
}

// SetCommissionCategory filters by commission calculation category.
func (s *GetAffiliateInviteeListService) SetCommissionCategory(category AffiliateCommissionCategory) *GetAffiliateInviteeListService {
	s.params["commissionCategory"] = string(category)
	return s
}

// SetOrderBy sets the sort field (cTime/depAmt/vol/fee/rebate, default cTime).
func (s *GetAffiliateInviteeListService) SetOrderBy(orderBy string) *GetAffiliateInviteeListService {
	s.params["orderBy"] = orderBy
	return s
}

// SetOrderDir sets the sort direction (asc/desc, default desc).
func (s *GetAffiliateInviteeListService) SetOrderDir(orderDir string) *GetAffiliateInviteeListService {
	s.params["orderDir"] = orderDir
	return s
}

// SetKYCStatus filters by KYC status (unverified/verified).
func (s *GetAffiliateInviteeListService) SetKYCStatus(status string) *GetAffiliateInviteeListService {
	s.params["kycStatus"] = status
	return s
}

// SetSubAffiliateUID restricts the page to invitees under one sub-affiliate.
func (s *GetAffiliateInviteeListService) SetSubAffiliateUID(uid string) *GetAffiliateInviteeListService {
	s.params["subAffiliateUid"] = uid
	return s
}

// SetUID matches invitees by external UID exactly: one UID or up to 100
// comma-separated. Unknown UIDs are skipped; if none resolve the page is empty
// (never the full list).
func (s *GetAffiliateInviteeListService) SetUID(uid string) *GetAffiliateInviteeListService {
	s.params["uid"] = uid
	return s
}

// SetJoinTimeBegin sets the inclusive lower bound on joinTime (ms). Must be sent
// together with SetJoinTimeEnd; the span may not exceed 90 days and may not
// start earlier than 180 days ago. Independent of the stats window.
func (s *GetAffiliateInviteeListService) SetJoinTimeBegin(ms int64) *GetAffiliateInviteeListService {
	s.params["joinTimeBegin"] = strconv.FormatInt(ms, 10)
	return s
}

// SetJoinTimeEnd sets the inclusive upper bound on joinTime (ms). Must be sent
// together with SetJoinTimeBegin; equal bounds are a valid single-point range.
func (s *GetAffiliateInviteeListService) SetJoinTimeEnd(ms int64) *GetAffiliateInviteeListService {
	s.params["joinTimeEnd"] = strconv.FormatInt(ms, 10)
	return s
}

func (s *GetAffiliateInviteeListService) Do(ctx context.Context) ([]AffiliateInvitee, error) {
	req := request.Get(ctx, s.c, "/api/v5/affiliate/invitee/list", s.params).WithSign()
	return request.DoList[AffiliateInvitee](req)
}

// AffiliateInvitee is one row of the affiliate's invitee list. The validating
// account does not have the affiliate role (the endpoint returns code 51620), so
// the field set is modeled from the OKX affiliate doc field table.
type AffiliateInvitee struct {
	UID             string          `json:"uid"`
	Country         string          `json:"country"`
	JoinTime        time.Time       `json:"joinTime"`
	FirstTradeTime  time.Time       `json:"firstTradeTime"`
	ChannelName     string          `json:"channelName"`
	RebateRate      decimal.Decimal `json:"rebateRate"`
	FeeTierRank     string          `json:"feeTierRank"`
	KYCStatus       string          `json:"kycStatus"`
	KYCTime         time.Time       `json:"kycTime"`
	DepositAmount   decimal.Decimal `json:"depAmt"`
	TotalVolume     decimal.Decimal `json:"totalVol"`
	TotalFee        decimal.Decimal `json:"totalFee"`
	TotalCommission decimal.Decimal `json:"totalCommission"`
	// IsCompliant is false when the invitee is restricted by KYC entity or
	// jurisdiction.
	IsCompliant bool `json:"isCompliant"`
}
