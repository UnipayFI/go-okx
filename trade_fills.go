package okx

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-okx/request"
	"github.com/shopspring/decimal"
)

// GetFillsService -- GET /api/v5/trade/fills (Read)
//
// Returns the account's transaction (fill) details for the last three days.
type GetFillsService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetFillsService() *GetFillsService {
	return &GetFillsService{c: c, params: map[string]string{}}
}

// SetInstType filters by product line (SPOT/MARGIN/SWAP/FUTURES/OPTION).
func (s *GetFillsService) SetInstType(instType InstType) *GetFillsService {
	s.params["instType"] = string(instType)
	return s
}

// SetUly filters by underlying.
func (s *GetFillsService) SetUly(uly string) *GetFillsService {
	s.params["uly"] = uly
	return s
}

// SetInstFamily filters by instrument family.
func (s *GetFillsService) SetInstFamily(instFamily string) *GetFillsService {
	s.params["instFamily"] = instFamily
	return s
}

// SetInstId filters by instrument id.
func (s *GetFillsService) SetInstId(instId string) *GetFillsService {
	s.params["instId"] = instId
	return s
}

// SetOrdId filters by order id.
func (s *GetFillsService) SetOrdId(ordId string) *GetFillsService {
	s.params["ordId"] = ordId
	return s
}

// SetSubType filters by transaction type (e.g. 1 buy, 2 sell, ...).
func (s *GetFillsService) SetSubType(subType string) *GetFillsService {
	s.params["subType"] = subType
	return s
}

// SetAfter paginates to records with a billId earlier than the given one (older).
func (s *GetFillsService) SetAfter(after string) *GetFillsService {
	s.params["after"] = after
	return s
}

// SetBefore paginates to records with a billId later than the given one (newer).
func (s *GetFillsService) SetBefore(before string) *GetFillsService {
	s.params["before"] = before
	return s
}

// SetBegin filters to fills at or after the given time (by ts).
func (s *GetFillsService) SetBegin(t time.Time) *GetFillsService {
	s.params["begin"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEnd filters to fills at or before the given time (by ts).
func (s *GetFillsService) SetEnd(t time.Time) *GetFillsService {
	s.params["end"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100, default 100).
func (s *GetFillsService) SetLimit(limit int) *GetFillsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetFillsService) Do(ctx context.Context) ([]Fill, error) {
	req := request.Get(ctx, s.c, "/api/v5/trade/fills", s.params).WithSign()
	return request.DoList[Fill](req)
}

// GetFillsHistoryService -- GET /api/v5/trade/fills-history (Read)
//
// Returns the account's transaction (fill) details over the last three months.
// instType is required.
type GetFillsHistoryService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetFillsHistoryService(instType InstType) *GetFillsHistoryService {
	return &GetFillsHistoryService{c: c, params: map[string]string{"instType": string(instType)}}
}

// SetUly filters by underlying.
func (s *GetFillsHistoryService) SetUly(uly string) *GetFillsHistoryService {
	s.params["uly"] = uly
	return s
}

// SetInstFamily filters by instrument family.
func (s *GetFillsHistoryService) SetInstFamily(instFamily string) *GetFillsHistoryService {
	s.params["instFamily"] = instFamily
	return s
}

// SetInstId filters by instrument id.
func (s *GetFillsHistoryService) SetInstId(instId string) *GetFillsHistoryService {
	s.params["instId"] = instId
	return s
}

// SetOrdId filters by order id.
func (s *GetFillsHistoryService) SetOrdId(ordId string) *GetFillsHistoryService {
	s.params["ordId"] = ordId
	return s
}

// SetSubType filters by transaction type (e.g. 1 buy, 2 sell, ...).
func (s *GetFillsHistoryService) SetSubType(subType string) *GetFillsHistoryService {
	s.params["subType"] = subType
	return s
}

// SetAfter paginates to records with a billId earlier than the given one (older).
func (s *GetFillsHistoryService) SetAfter(after string) *GetFillsHistoryService {
	s.params["after"] = after
	return s
}

// SetBefore paginates to records with a billId later than the given one (newer).
func (s *GetFillsHistoryService) SetBefore(before string) *GetFillsHistoryService {
	s.params["before"] = before
	return s
}

// SetBegin filters to fills at or after the given time (by ts).
func (s *GetFillsHistoryService) SetBegin(t time.Time) *GetFillsHistoryService {
	s.params["begin"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEnd filters to fills at or before the given time (by ts).
func (s *GetFillsHistoryService) SetEnd(t time.Time) *GetFillsHistoryService {
	s.params["end"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100, default 100).
func (s *GetFillsHistoryService) SetLimit(limit int) *GetFillsHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetFillsHistoryService) Do(ctx context.Context) ([]Fill, error) {
	req := request.Get(ctx, s.c, "/api/v5/trade/fills-history", s.params).WithSign()
	return request.DoList[Fill](req)
}

// Fill is one transaction (fill) detail, shared by /api/v5/trade/fills and
// /api/v5/trade/fills-history.
type Fill struct {
	InstrumentType      InstType        `json:"instType"`
	InstrumentID        string          `json:"instId"`
	TradeID             string          `json:"tradeId"`
	OrderID             string          `json:"ordId"`
	ClientOrderID       string          `json:"clOrdId"`
	BillID              string          `json:"billId"`
	SubType             string          `json:"subType"`
	Tag                 string          `json:"tag"`
	FillPrice           decimal.Decimal `json:"fillPx"`
	FillSize            decimal.Decimal `json:"fillSz"`
	FillIndexPrice      decimal.Decimal `json:"fillIdxPx"`
	FillPnl             decimal.Decimal `json:"fillPnl"`
	FillPriceVolatility decimal.Decimal `json:"fillPxVol"`
	FillPriceUSD        decimal.Decimal `json:"fillPxUsd"`
	FillMarkVolatility  decimal.Decimal `json:"fillMarkVol"`
	FillForwardPrice    decimal.Decimal `json:"fillFwdPx"`
	FillMarkPrice       decimal.Decimal `json:"fillMarkPx"`
	Side                Side            `json:"side"`
	PositionSide        PosSide         `json:"posSide"`
	ExecutionType       ExecType        `json:"execType"`
	FeeCurrency         string          `json:"feeCcy"`
	Fee                 decimal.Decimal `json:"fee"`
	TradeQuoteCurrency  string          `json:"tradeQuoteCcy"`
	Timestamp           time.Time       `json:"ts"`
	FillTime            time.Time       `json:"fillTime"`
}
