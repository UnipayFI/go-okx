package okx

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-okx/request"
	"github.com/shopspring/decimal"
)

// GetBillsService -- GET /api/v5/account/bills (Read)
//
// Returns the account's bill (balance-change) records for the last 7 days, most
// recent first.
type GetBillsService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetBillsService() *GetBillsService {
	return &GetBillsService{c: c, params: map[string]string{}}
}

// SetInstType filters bills by instrument type (SPOT/MARGIN/SWAP/FUTURES/OPTION).
func (s *GetBillsService) SetInstType(instType InstType) *GetBillsService {
	s.params["instType"] = string(instType)
	return s
}

// SetInstId filters bills by a single instrument id.
func (s *GetBillsService) SetInstId(instId string) *GetBillsService {
	s.params["instId"] = instId
	return s
}

// SetCcy filters bills by currency.
func (s *GetBillsService) SetCcy(ccy string) *GetBillsService {
	s.params["ccy"] = ccy
	return s
}

// SetMgnMode filters bills by margin mode (isolated/cross).
func (s *GetBillsService) SetMgnMode(mgnMode MgnMode) *GetBillsService {
	s.params["mgnMode"] = string(mgnMode)
	return s
}

// SetCtType filters bills by contract type (linear/inverse).
func (s *GetBillsService) SetCtType(ctType CtType) *GetBillsService {
	s.params["ctType"] = string(ctType)
	return s
}

// SetType filters bills by bill type (e.g. "1" transfer, "2" trade, "8" funding fee).
func (s *GetBillsService) SetType(typ string) *GetBillsService {
	s.params["type"] = typ
	return s
}

// SetSubType filters bills by bill sub-type (e.g. "1" buy, "2" sell).
func (s *GetBillsService) SetSubType(subType string) *GetBillsService {
	s.params["subType"] = subType
	return s
}

// SetAfter paginates to bills with a billId smaller than the given value (older).
func (s *GetBillsService) SetAfter(billId string) *GetBillsService {
	s.params["after"] = billId
	return s
}

// SetBefore paginates to bills with a billId larger than the given value (newer).
func (s *GetBillsService) SetBefore(billId string) *GetBillsService {
	s.params["before"] = billId
	return s
}

// SetBegin filters to bills at or after the given time (by ts).
func (s *GetBillsService) SetBegin(t time.Time) *GetBillsService {
	s.params["begin"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEnd filters to bills at or before the given time (by ts).
func (s *GetBillsService) SetEnd(t time.Time) *GetBillsService {
	s.params["end"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100, default 100).
func (s *GetBillsService) SetLimit(limit int) *GetBillsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetBillsService) Do(ctx context.Context) ([]Bill, error) {
	req := request.Get(ctx, s.c, "/api/v5/account/bills", s.params).WithSign()
	return request.DoList[Bill](req)
}

// GetBillsArchiveService -- GET /api/v5/account/bills-archive (Read)
//
// Returns the account's bill (balance-change) records for the last 3 months,
// most recent first. Same shape and filters as /account/bills.
type GetBillsArchiveService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetBillsArchiveService() *GetBillsArchiveService {
	return &GetBillsArchiveService{c: c, params: map[string]string{}}
}

// SetInstType filters bills by instrument type (SPOT/MARGIN/SWAP/FUTURES/OPTION).
func (s *GetBillsArchiveService) SetInstType(instType InstType) *GetBillsArchiveService {
	s.params["instType"] = string(instType)
	return s
}

// SetInstId filters bills by a single instrument id.
func (s *GetBillsArchiveService) SetInstId(instId string) *GetBillsArchiveService {
	s.params["instId"] = instId
	return s
}

// SetCcy filters bills by currency.
func (s *GetBillsArchiveService) SetCcy(ccy string) *GetBillsArchiveService {
	s.params["ccy"] = ccy
	return s
}

// SetMgnMode filters bills by margin mode (isolated/cross).
func (s *GetBillsArchiveService) SetMgnMode(mgnMode MgnMode) *GetBillsArchiveService {
	s.params["mgnMode"] = string(mgnMode)
	return s
}

// SetCtType filters bills by contract type (linear/inverse).
func (s *GetBillsArchiveService) SetCtType(ctType CtType) *GetBillsArchiveService {
	s.params["ctType"] = string(ctType)
	return s
}

// SetType filters bills by bill type (e.g. "1" transfer, "2" trade, "8" funding fee).
func (s *GetBillsArchiveService) SetType(typ string) *GetBillsArchiveService {
	s.params["type"] = typ
	return s
}

// SetSubType filters bills by bill sub-type (e.g. "1" buy, "2" sell).
func (s *GetBillsArchiveService) SetSubType(subType string) *GetBillsArchiveService {
	s.params["subType"] = subType
	return s
}

// SetAfter paginates to bills with a billId smaller than the given value (older).
func (s *GetBillsArchiveService) SetAfter(billId string) *GetBillsArchiveService {
	s.params["after"] = billId
	return s
}

// SetBefore paginates to bills with a billId larger than the given value (newer).
func (s *GetBillsArchiveService) SetBefore(billId string) *GetBillsArchiveService {
	s.params["before"] = billId
	return s
}

// SetBegin filters to bills at or after the given time (by ts).
func (s *GetBillsArchiveService) SetBegin(t time.Time) *GetBillsArchiveService {
	s.params["begin"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEnd filters to bills at or before the given time (by ts).
func (s *GetBillsArchiveService) SetEnd(t time.Time) *GetBillsArchiveService {
	s.params["end"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100, default 100).
func (s *GetBillsArchiveService) SetLimit(limit int) *GetBillsArchiveService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetBillsArchiveService) Do(ctx context.Context) ([]Bill, error) {
	req := request.Get(ctx, s.c, "/api/v5/account/bills-archive", s.params).WithSign()
	return request.DoList[Bill](req)
}

// Bill is a single account balance-change record, shared by /account/bills and
// /account/bills-archive.
type Bill struct {
	InstrumentType        InstType        `json:"instType"`
	InstrumentID          string          `json:"instId"`
	BillID                string          `json:"billId"`
	Type                  string          `json:"type"`
	SubType               string          `json:"subType"`
	Timestamp             time.Time       `json:"ts"`
	BalanceChange         decimal.Decimal `json:"balChg"`
	PositionBalanceChange decimal.Decimal `json:"posBalChg"`
	Balance               decimal.Decimal `json:"bal"`
	PositionBalance       decimal.Decimal `json:"posBal"`
	Size                  decimal.Decimal `json:"sz"`
	Price                 decimal.Decimal `json:"px"`
	Pnl                   decimal.Decimal `json:"pnl"`
	Fee                   decimal.Decimal `json:"fee"`
	Currency              string          `json:"ccy"`
	OrderID               string          `json:"ordId"`
	ClientOrderID         string          `json:"clOrdId"`
	TradeID               string          `json:"tradeId"`
	ExecutionType         ExecType        `json:"execType"`
	MarginMode            MgnMode         `json:"mgnMode"`
	ContractType          CtType          `json:"ctType"`
	From                  string          `json:"from"`
	To                    string          `json:"to"`
	Notes                 string          `json:"notes"`
	Interest              decimal.Decimal `json:"interest"`
	Tag                   string          `json:"tag"`
	EarnAmount            decimal.Decimal `json:"earnAmt"`
	EarnAPR               decimal.Decimal `json:"earnApr"`
	FillTime              time.Time       `json:"fillTime"`
	FillForwardPrice      decimal.Decimal `json:"fillFwdPx"`
	FillIndexPrice        decimal.Decimal `json:"fillIdxPx"`
	FillMarkPrice         decimal.Decimal `json:"fillMarkPx"`
	FillMarkVolatility    decimal.Decimal `json:"fillMarkVol"`
	FillPriceVolatility   decimal.Decimal `json:"fillPxVol"`
	FillPriceUSD          decimal.Decimal `json:"fillPxUsd"`
}

// BillsHistoryArchiveQuarter selects the quarter of a yearly bills-history
// archive request.
type BillsHistoryArchiveQuarter string

const (
	BillsHistoryArchiveQuarterQ1 BillsHistoryArchiveQuarter = "Q1"
	BillsHistoryArchiveQuarterQ2 BillsHistoryArchiveQuarter = "Q2"
	BillsHistoryArchiveQuarterQ3 BillsHistoryArchiveQuarter = "Q3"
	BillsHistoryArchiveQuarterQ4 BillsHistoryArchiveQuarter = "Q4"
)

// ApplyBillsHistoryArchiveService -- POST /api/v5/account/bills-history-archive (Trade)
//
// Applies for a downloadable archive of one quarter's bill records. After the
// request is processed the download link is fetched via
// GetBillsHistoryArchiveService.
type ApplyBillsHistoryArchiveService struct {
	c    *Client
	body map[string]any
}

// NewApplyBillsHistoryArchiveService builds the apply request. year is a 4-digit
// year (e.g. "2024"); quarter is Q1..Q4.
func (c *Client) NewApplyBillsHistoryArchiveService(year string, quarter BillsHistoryArchiveQuarter) *ApplyBillsHistoryArchiveService {
	return &ApplyBillsHistoryArchiveService{c: c, body: map[string]any{
		"year":    year,
		"quarter": string(quarter),
	}}
}

func (s *ApplyBillsHistoryArchiveService) Do(ctx context.Context) (*BillsHistoryArchive, error) {
	req := request.Post(ctx, s.c, "/api/v5/account/bills-history-archive", s.body).WithSign()
	return request.DoOne[BillsHistoryArchive](req)
}

// GetBillsHistoryArchiveService -- GET /api/v5/account/bills-history-archive (Read)
//
// Returns the download link and state of a previously requested quarterly bills
// archive (apply first via ApplyBillsHistoryArchiveService).
type GetBillsHistoryArchiveService struct {
	c      *Client
	params map[string]string
}

// NewGetBillsHistoryArchiveService builds the request. year is a 4-digit year
// (e.g. "2024"); quarter is Q1..Q4.
func (c *Client) NewGetBillsHistoryArchiveService(year string, quarter BillsHistoryArchiveQuarter) *GetBillsHistoryArchiveService {
	return &GetBillsHistoryArchiveService{c: c, params: map[string]string{
		"year":    year,
		"quarter": string(quarter),
	}}
}

func (s *GetBillsHistoryArchiveService) Do(ctx context.Context) (*BillsHistoryArchive, error) {
	req := request.Get(ctx, s.c, "/api/v5/account/bills-history-archive", s.params).WithSign()
	return request.DoOne[BillsHistoryArchive](req)
}

// BillsHistoryArchive is the state and download link of a quarterly bills
// archive request.
type BillsHistoryArchive struct {
	Timestamp time.Time `json:"ts"`
	FileHref  string    `json:"fileHref"`
	State     string    `json:"state"`
}
