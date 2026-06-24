package okx

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-okx/request"
	"github.com/shopspring/decimal"
)

// ConvertStatus is the lifecycle state of an easy-convert / one-click-repay job.
type ConvertStatus string

const (
	ConvertStatusRunning ConvertStatus = "running"
	ConvertStatusFilled  ConvertStatus = "filled"
	ConvertStatusFailed  ConvertStatus = "failed"
)

// DebtType selects whether a one-click-repay currency list reports cross or
// isolated debt.
type DebtType string

const (
	DebtTypeCross    DebtType = "cross"
	DebtTypeIsolated DebtType = "isolated"
)

// ConvertSource selects the funding source for an easy-convert order: "1"
// (trading account) or "2" (funding account).
type ConvertSource string

const (
	ConvertSourceTrading ConvertSource = "1"
	ConvertSourceFunding ConvertSource = "2"
)

// GetEasyConvertCurrencyListService -- GET /api/v5/trade/easy-convert-currency-list (Read)
//
// Returns the small-balance ("dust") currencies eligible for one-click
// easy-convert and the mainstream currencies they can be converted to.
type GetEasyConvertCurrencyListService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetEasyConvertCurrencyListService() *GetEasyConvertCurrencyListService {
	return &GetEasyConvertCurrencyListService{c: c, params: map[string]string{}}
}

// SetSource sets the funding source: "1" (trading account, default) or "2"
// (funding account).
func (s *GetEasyConvertCurrencyListService) SetSource(source ConvertSource) *GetEasyConvertCurrencyListService {
	s.params["source"] = string(source)
	return s
}

func (s *GetEasyConvertCurrencyListService) Do(ctx context.Context) (*EasyConvertCurrencyList, error) {
	req := request.Get(ctx, s.c, "/api/v5/trade/easy-convert-currency-list", s.params).WithSign()
	return request.DoOne[EasyConvertCurrencyList](req)
}

// EasyConvertCurrencyList is the set of convertible dust currencies and the
// mainstream currencies they can be converted into.
type EasyConvertCurrencyList struct {
	FromData   []EasyConvertFromCcy `json:"fromData"`
	ToCurrency []string             `json:"toCcy"`
}

// EasyConvertFromCcy is a single dust currency and its convertible balance.
type EasyConvertFromCcy struct {
	FromCurrency string          `json:"fromCcy"`
	FromAmount   decimal.Decimal `json:"fromAmt"`
}

// PlaceEasyConvertService -- POST /api/v5/trade/easy-convert (Trade)
//
// Converts up to five small-balance ("dust") currencies into a single
// mainstream currency. State-changing: implemented but never executed by tests.
type PlaceEasyConvertService struct {
	c    *Client
	body map[string]any
}

// NewPlaceEasyConvertService starts an easy-convert order. fromCcy is the list
// of dust currencies (max 5); toCcy is the single mainstream currency to receive.
func (c *Client) NewPlaceEasyConvertService(fromCcy []string, toCcy string) *PlaceEasyConvertService {
	return &PlaceEasyConvertService{c: c, body: map[string]any{
		"fromCcy": fromCcy,
		"toCcy":   toCcy,
	}}
}

// SetSource sets the funding source: "1" (trading account, default) or "2"
// (funding account).
func (s *PlaceEasyConvertService) SetSource(source ConvertSource) *PlaceEasyConvertService {
	s.body["source"] = string(source)
	return s
}

func (s *PlaceEasyConvertService) Do(ctx context.Context) ([]EasyConvertResult, error) {
	req := request.Post(ctx, s.c, "/api/v5/trade/easy-convert", s.body).WithSign()
	return request.DoListPartial[EasyConvertResult](req)
}

// EasyConvertResult is the acknowledgement for an easy-convert order.
type EasyConvertResult struct {
	Status       ConvertStatus   `json:"status"`
	FromCurrency string          `json:"fromCcy"`
	ToCurrency   string          `json:"toCcy"`
	FillFromSize decimal.Decimal `json:"fillFromSz"`
	FillToSize   decimal.Decimal `json:"fillToSz"`
	UpdateTime   time.Time       `json:"uTime"`
}

// GetEasyConvertHistoryService -- GET /api/v5/trade/easy-convert-history (Read)
//
// Returns the easy-convert order history (last 7 days).
type GetEasyConvertHistoryService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetEasyConvertHistoryService() *GetEasyConvertHistoryService {
	return &GetEasyConvertHistoryService{c: c, params: map[string]string{}}
}

// SetAfter returns records earlier than the given time (paginates by uTime).
func (s *GetEasyConvertHistoryService) SetAfter(t time.Time) *GetEasyConvertHistoryService {
	s.params["after"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetBefore returns records newer than the given time (paginates by uTime).
func (s *GetEasyConvertHistoryService) SetBefore(t time.Time) *GetEasyConvertHistoryService {
	s.params["before"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit sets the number of results per request (max 100, default 100).
func (s *GetEasyConvertHistoryService) SetLimit(limit int) *GetEasyConvertHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetEasyConvertHistoryService) Do(ctx context.Context) ([]EasyConvertHistory, error) {
	req := request.Get(ctx, s.c, "/api/v5/trade/easy-convert-history", s.params).WithSign()
	return request.DoList[EasyConvertHistory](req)
}

// EasyConvertHistory is a single past easy-convert order.
type EasyConvertHistory struct {
	FromCurrency string          `json:"fromCcy"`
	FillFromSize decimal.Decimal `json:"fillFromSz"`
	ToCurrency   string          `json:"toCcy"`
	FillToSize   decimal.Decimal `json:"fillToSz"`
	Account      string          `json:"acct"`
	Status       ConvertStatus   `json:"status"`
	UpdateTime   time.Time       `json:"uTime"`
}

// GetOneClickRepayCurrencyListService -- GET /api/v5/trade/one-click-repay-currency-list (Read)
//
// Returns the cross/isolated debt currencies eligible for one-click repay and
// the currencies that can be used to repay them.
type GetOneClickRepayCurrencyListService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetOneClickRepayCurrencyListService() *GetOneClickRepayCurrencyListService {
	return &GetOneClickRepayCurrencyListService{c: c, params: map[string]string{}}
}

// SetDebtType filters by debt type: "cross" or "isolated".
func (s *GetOneClickRepayCurrencyListService) SetDebtType(debtType DebtType) *GetOneClickRepayCurrencyListService {
	s.params["debtType"] = string(debtType)
	return s
}

func (s *GetOneClickRepayCurrencyListService) Do(ctx context.Context) ([]OneClickRepayCurrencyList, error) {
	req := request.Get(ctx, s.c, "/api/v5/trade/one-click-repay-currency-list", s.params).WithSign()
	return request.DoList[OneClickRepayCurrencyList](req)
}

// OneClickRepayCurrencyList is the set of repayable debt currencies and the
// currencies available to repay them, grouped by debt type.
type OneClickRepayCurrencyList struct {
	DebtType  DebtType             `json:"debtType"`
	DebtData  []OneClickRepayDebt  `json:"debtData"`
	RepayData []OneClickRepayRepay `json:"repayData"`
}

// OneClickRepayDebt is a single debt currency and its outstanding amount.
type OneClickRepayDebt struct {
	DebtCurrency string          `json:"debtCcy"`
	DebtAmount   decimal.Decimal `json:"debtAmt"`
}

// OneClickRepayRepay is a single currency available to repay debt and its
// available balance.
type OneClickRepayRepay struct {
	RepayCurrency string          `json:"repayCcy"`
	RepayAmount   decimal.Decimal `json:"repayAmt"`
}

// TradeOneClickRepayService -- POST /api/v5/trade/one-click-repay (Trade)
//
// Repays up to five debt currencies using a single repay currency.
// State-changing: implemented but never executed by tests.
type TradeOneClickRepayService struct {
	c    *Client
	body map[string]any
}

// NewTradeOneClickRepayService starts a one-click repay. debtCcy is the list of
// debt currencies (max 5); repayCcy is the single currency used to repay.
func (c *Client) NewTradeOneClickRepayService(debtCcy []string, repayCcy string) *TradeOneClickRepayService {
	return &TradeOneClickRepayService{c: c, body: map[string]any{
		"debtCcy":  debtCcy,
		"repayCcy": repayCcy,
	}}
}

func (s *TradeOneClickRepayService) Do(ctx context.Context) ([]OneClickRepayResult, error) {
	req := request.Post(ctx, s.c, "/api/v5/trade/one-click-repay", s.body).WithSign()
	return request.DoListPartial[OneClickRepayResult](req)
}

// OneClickRepayResult is the acknowledgement for a one-click repay order.
type OneClickRepayResult struct {
	Status        ConvertStatus   `json:"status"`
	DebtCurrency  string          `json:"debtCcy"`
	RepayCurrency string          `json:"repayCcy"`
	FillDebtSize  decimal.Decimal `json:"fillDebtSz"`
	FillRepaySize decimal.Decimal `json:"fillRepaySz"`
	UpdateTime    time.Time       `json:"uTime"`
}

// GetOneClickRepayHistoryService -- GET /api/v5/trade/one-click-repay-history (Read)
//
// Returns the one-click repay order history (last 7 days).
type GetOneClickRepayHistoryService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetOneClickRepayHistoryService() *GetOneClickRepayHistoryService {
	return &GetOneClickRepayHistoryService{c: c, params: map[string]string{}}
}

// SetAfter returns records earlier than the given time (paginates by uTime).
func (s *GetOneClickRepayHistoryService) SetAfter(t time.Time) *GetOneClickRepayHistoryService {
	s.params["after"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetBefore returns records newer than the given time (paginates by uTime).
func (s *GetOneClickRepayHistoryService) SetBefore(t time.Time) *GetOneClickRepayHistoryService {
	s.params["before"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit sets the number of results per request (max 100, default 100).
func (s *GetOneClickRepayHistoryService) SetLimit(limit int) *GetOneClickRepayHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetOneClickRepayHistoryService) Do(ctx context.Context) ([]OneClickRepayHistory, error) {
	req := request.Get(ctx, s.c, "/api/v5/trade/one-click-repay-history", s.params).WithSign()
	return request.DoList[OneClickRepayHistory](req)
}

// OneClickRepayHistory is a single past one-click repay order.
type OneClickRepayHistory struct {
	DebtCurrency  string          `json:"debtCcy"`
	FillDebtSize  decimal.Decimal `json:"fillDebtSz"`
	RepayCurrency string          `json:"repayCcy"`
	FillRepaySize decimal.Decimal `json:"fillRepaySz"`
	Status        ConvertStatus   `json:"status"`
	UpdateTime    time.Time       `json:"uTime"`
}

// GetOneClickRepayCurrencyListV2Service -- GET /api/v5/trade/one-click-repay-currency-list-v2 (Read)
//
// Returns the debt currencies eligible for one-click repay (v2) and the
// currencies that can be used to repay them.
type GetOneClickRepayCurrencyListV2Service struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetOneClickRepayCurrencyListV2Service() *GetOneClickRepayCurrencyListV2Service {
	return &GetOneClickRepayCurrencyListV2Service{c: c, params: map[string]string{}}
}

func (s *GetOneClickRepayCurrencyListV2Service) Do(ctx context.Context) ([]OneClickRepayCurrencyListV2, error) {
	req := request.Get(ctx, s.c, "/api/v5/trade/one-click-repay-currency-list-v2", s.params).WithSign()
	return request.DoList[OneClickRepayCurrencyListV2](req)
}

// OneClickRepayCurrencyListV2 is the set of repayable debt currencies and the
// currencies available to repay them (v2).
type OneClickRepayCurrencyListV2 struct {
	DebtData  []OneClickRepayDebt  `json:"debtData"`
	RepayData []OneClickRepayRepay `json:"repayData"`
}

// TradeOneClickRepayV2Service -- POST /api/v5/trade/one-click-repay-v2 (Trade)
//
// Repays a single debt currency using a prioritized list of repay currencies.
// State-changing: implemented but never executed by tests.
type TradeOneClickRepayV2Service struct {
	c    *Client
	body map[string]any
}

// NewTradeOneClickRepayV2Service starts a one-click repay (v2). debtCcy is the
// single debt currency; repayCcyList is the prioritized list of currencies used
// to repay it.
func (c *Client) NewTradeOneClickRepayV2Service(debtCcy string, repayCcyList []string) *TradeOneClickRepayV2Service {
	return &TradeOneClickRepayV2Service{c: c, body: map[string]any{
		"debtCcy":      debtCcy,
		"repayCcyList": repayCcyList,
	}}
}

func (s *TradeOneClickRepayV2Service) Do(ctx context.Context) ([]OneClickRepayResultV2, error) {
	req := request.Post(ctx, s.c, "/api/v5/trade/one-click-repay-v2", s.body).WithSign()
	return request.DoListPartial[OneClickRepayResultV2](req)
}

// OneClickRepayResultV2 is the acknowledgement for a one-click repay (v2) order.
type OneClickRepayResultV2 struct {
	DebtCurrency      string    `json:"debtCcy"`
	RepayCurrencyList []string  `json:"repayCcyList"`
	Timestamp         time.Time `json:"ts"`
}

// GetOneClickRepayHistoryV2Service -- GET /api/v5/trade/one-click-repay-history-v2 (Read)
//
// Returns the one-click repay (v2) order history.
type GetOneClickRepayHistoryV2Service struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetOneClickRepayHistoryV2Service() *GetOneClickRepayHistoryV2Service {
	return &GetOneClickRepayHistoryV2Service{c: c, params: map[string]string{}}
}

// SetAfter returns records earlier than (included) the given ts.
func (s *GetOneClickRepayHistoryV2Service) SetAfter(t time.Time) *GetOneClickRepayHistoryV2Service {
	s.params["after"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetBefore returns records newer than (included) the given ts.
func (s *GetOneClickRepayHistoryV2Service) SetBefore(t time.Time) *GetOneClickRepayHistoryV2Service {
	s.params["before"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit sets the number of results per request (max 100, default 100).
func (s *GetOneClickRepayHistoryV2Service) SetLimit(limit int) *GetOneClickRepayHistoryV2Service {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetOneClickRepayHistoryV2Service) Do(ctx context.Context) ([]OneClickRepayHistoryV2, error) {
	req := request.Get(ctx, s.c, "/api/v5/trade/one-click-repay-history-v2", s.params).WithSign()
	return request.DoList[OneClickRepayHistoryV2](req)
}

// OneClickRepayHistoryV2 is a single past one-click repay (v2) order.
type OneClickRepayHistoryV2 struct {
	DebtCurrency      string                  `json:"debtCcy"`
	RepayCurrencyList []string                `json:"repayCcyList"`
	FillDebtSize      decimal.Decimal         `json:"fillDebtSz"`
	Status            ConvertStatus           `json:"status"`
	OrderIDInfo       []OneClickRepayOrderRef `json:"ordIdInfo"`
	Timestamp         time.Time               `json:"ts"`
}

// OneClickRepayOrderRef is a single underlying order placed to execute a
// one-click repay (v2).
type OneClickRepayOrderRef struct {
	OrderID      string          `json:"ordId"`
	InstrumentID string          `json:"instId"`
	OrderType    OrdType         `json:"ordType"`
	Side         Side            `json:"side"`
	Price        decimal.Decimal `json:"px"`
	Size         decimal.Decimal `json:"sz"`
	FillPrice    decimal.Decimal `json:"fillPx"`
	FillSize     decimal.Decimal `json:"fillSz"`
	State        OrdState        `json:"state"`
	CreationTime time.Time       `json:"cTime"`
}
