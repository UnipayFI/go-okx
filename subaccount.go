package okx

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-okx/request"
	"github.com/shopspring/decimal"
)

// subAccountType is a sub-account's classification (standard, managed-trading,
// custody-trading, ...), reported as a bare-number string by the sub-account
// endpoints.
type subAccountType string

const (
	subAccountTypeStandard       subAccountType = "1"
	subAccountTypeManagedTrading subAccountType = "2"
	subAccountTypeCustodyAnnex   subAccountType = "5"
)

// GetSubAccountListService -- GET /api/v5/users/subaccount/list (Read)
//
// Returns the list of sub-accounts under the (master) account, with their
// enable/trade-out flags and creation time.
type GetSubAccountListService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetSubAccountListService() *GetSubAccountListService {
	return &GetSubAccountListService{c: c, params: map[string]string{}}
}

// SetEnable filters by sub-account enabled state.
func (s *GetSubAccountListService) SetEnable(enable bool) *GetSubAccountListService {
	s.params["enable"] = strconv.FormatBool(enable)
	return s
}

// SetSubAcct filters by a single sub-account name.
func (s *GetSubAccountListService) SetSubAcct(subAcct string) *GetSubAccountListService {
	s.params["subAcct"] = subAcct
	return s
}

// SetAfter paginates to sub-accounts created earlier than the given time (older).
func (s *GetSubAccountListService) SetAfter(t time.Time) *GetSubAccountListService {
	s.params["after"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetBefore paginates to sub-accounts created later than the given time (newer).
func (s *GetSubAccountListService) SetBefore(t time.Time) *GetSubAccountListService {
	s.params["before"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetSubAccountListService) SetLimit(limit int) *GetSubAccountListService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetSubAccountListService) Do(ctx context.Context) ([]SubAccount, error) {
	req := request.Get(ctx, s.c, "/api/v5/users/subaccount/list", s.params).WithSign()
	return request.DoList[SubAccount](req)
}

// SubAccount is a single sub-account.
type SubAccount struct {
	Type           subAccountType `json:"type"`
	Enable         bool           `json:"enable"`
	SubAccount     string         `json:"subAcct"`
	UID            string         `json:"uid"`
	Label          string         `json:"label"`
	Mobile         string         `json:"mobile"`
	GoogleAuth     bool           `json:"gAuth"`
	Frozen         bool           `json:"frozen"`
	CanTransferOut bool           `json:"canTransOut"`
	FrozenFunc     []string       `json:"frozenFunc"`
	Timestamp      time.Time      `json:"ts"`
}

// ModifySubAccountApiKeyService -- POST /api/v5/users/subaccount/modify-apikey (Trade)
//
// Modifies (master account) a sub-account API key's permissions, label or IP
// allow-list. Implement-only: never executed by the live test suite.
type ModifySubAccountApiKeyService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewModifySubAccountApiKeyService(subAcct, apiKey string) *ModifySubAccountApiKeyService {
	return &ModifySubAccountApiKeyService{c: c, body: map[string]any{
		"subAcct": subAcct,
		"apiKey":  apiKey,
	}}
}

// SetLabel sets a new label for the API key.
func (s *ModifySubAccountApiKeyService) SetLabel(label string) *ModifySubAccountApiKeyService {
	s.body["label"] = label
	return s
}

// SetPerm sets the API key permissions (comma-separated: read_only, trade).
func (s *ModifySubAccountApiKeyService) SetPerm(perm string) *ModifySubAccountApiKeyService {
	s.body["perm"] = perm
	return s
}

// SetIP sets the IP allow-list (comma-separated, up to 20 addresses).
func (s *ModifySubAccountApiKeyService) SetIP(ip string) *ModifySubAccountApiKeyService {
	s.body["ip"] = ip
	return s
}

func (s *ModifySubAccountApiKeyService) Do(ctx context.Context) (*SubAccountApiKey, error) {
	req := request.Post(ctx, s.c, "/api/v5/users/subaccount/modify-apikey", s.body).WithSign()
	return request.DoOne[SubAccountApiKey](req)
}

// CreateSubAccountApiKeyService -- POST /api/v5/users/subaccount/apikey (Trade)
//
// Creates an API key for a sub-account (master account). Implement-only: never
// executed by the live test suite.
type CreateSubAccountApiKeyService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewCreateSubAccountApiKeyService(subAcct, label, passphrase string) *CreateSubAccountApiKeyService {
	return &CreateSubAccountApiKeyService{c: c, body: map[string]any{
		"subAcct":    subAcct,
		"label":      label,
		"passphrase": passphrase,
	}}
}

// SetPerm sets the API key permissions (comma-separated: read_only, trade).
func (s *CreateSubAccountApiKeyService) SetPerm(perm string) *CreateSubAccountApiKeyService {
	s.body["perm"] = perm
	return s
}

// SetIP sets the IP allow-list (comma-separated, up to 20 addresses).
func (s *CreateSubAccountApiKeyService) SetIP(ip string) *CreateSubAccountApiKeyService {
	s.body["ip"] = ip
	return s
}

func (s *CreateSubAccountApiKeyService) Do(ctx context.Context) (*SubAccountApiKey, error) {
	req := request.Post(ctx, s.c, "/api/v5/users/subaccount/apikey", s.body).WithSign()
	return request.DoOne[SubAccountApiKey](req)
}

// GetSubAccountApiKeyService -- GET /api/v5/users/subaccount/apikey (Read)
//
// Returns the API keys (excluding the secret) of a sub-account.
type GetSubAccountApiKeyService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetSubAccountApiKeyService(subAcct string) *GetSubAccountApiKeyService {
	return &GetSubAccountApiKeyService{c: c, params: map[string]string{"subAcct": subAcct}}
}

// SetApiKey filters by a single API key.
func (s *GetSubAccountApiKeyService) SetApiKey(apiKey string) *GetSubAccountApiKeyService {
	s.params["apiKey"] = apiKey
	return s
}

func (s *GetSubAccountApiKeyService) Do(ctx context.Context) ([]SubAccountApiKey, error) {
	req := request.Get(ctx, s.c, "/api/v5/users/subaccount/apikey", s.params).WithSign()
	return request.DoList[SubAccountApiKey](req)
}

// SubAccountApiKey is a sub-account API key and its permissions.
type SubAccountApiKey struct {
	SubAccount string    `json:"subAcct"`
	Label      string    `json:"label"`
	APIKey     string    `json:"apiKey"`
	SecretKey  string    `json:"secretKey"`
	Perm       string    `json:"perm"`
	IP         string    `json:"ip"`
	Timestamp  time.Time `json:"ts"`
}

// DeleteSubAccountApiKeyService -- POST /api/v5/users/subaccount/delete-apikey (Trade)
//
// Deletes a sub-account API key (master account). Implement-only: never executed
// by the live test suite.
type DeleteSubAccountApiKeyService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewDeleteSubAccountApiKeyService(subAcct, apiKey string) *DeleteSubAccountApiKeyService {
	return &DeleteSubAccountApiKeyService{c: c, body: map[string]any{
		"subAcct": subAcct,
		"apiKey":  apiKey,
	}}
}

func (s *DeleteSubAccountApiKeyService) Do(ctx context.Context) (*SubAccountApiKey, error) {
	req := request.Post(ctx, s.c, "/api/v5/users/subaccount/delete-apikey", s.body).WithSign()
	return request.DoOne[SubAccountApiKey](req)
}

// GetSubAccountTradingBalancesService -- GET /api/v5/account/subaccount/balances (Read)
//
// Returns a sub-account's trading-account balance details (master account).
type GetSubAccountTradingBalancesService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetSubAccountTradingBalancesService(subAcct string) *GetSubAccountTradingBalancesService {
	return &GetSubAccountTradingBalancesService{c: c, params: map[string]string{"subAcct": subAcct}}
}

func (s *GetSubAccountTradingBalancesService) Do(ctx context.Context) ([]SubAccountTradingBalance, error) {
	req := request.Get(ctx, s.c, "/api/v5/account/subaccount/balances", s.params).WithSign()
	return request.DoList[SubAccountTradingBalance](req)
}

// SubAccountTradingBalance is a sub-account's trading-account balance summary.
type SubAccountTradingBalance struct {
	AdjustedEquity decimal.Decimal               `json:"adjEq"`
	BorrowFrozen   decimal.Decimal               `json:"borrowFroz"`
	IMR            decimal.Decimal               `json:"imr"`
	IsolatedEquity decimal.Decimal               `json:"isoEq"`
	MarginRatio    decimal.Decimal               `json:"mgnRatio"`
	MMR            decimal.Decimal               `json:"mmr"`
	NotionalUSD    decimal.Decimal               `json:"notionalUsd"`
	OrderFrozen    decimal.Decimal               `json:"ordFroz"`
	TotalEquity    decimal.Decimal               `json:"totalEq"`
	UpdateTime     time.Time                     `json:"uTime"`
	Details        []SubAccountTradingBalanceCcy `json:"details"`
}

// SubAccountTradingBalanceCcy is one currency's balance within a sub-account's
// trading account.
type SubAccountTradingBalanceCcy struct {
	Currency            string          `json:"ccy"`
	Equity              decimal.Decimal `json:"eq"`
	CashBalance         decimal.Decimal `json:"cashBal"`
	UpdateTime          time.Time       `json:"uTime"`
	IsolatedEquity      decimal.Decimal `json:"isoEq"`
	AvailableEquity     decimal.Decimal `json:"availEq"`
	DiscountEquity      decimal.Decimal `json:"disEq"`
	AvailableBalance    decimal.Decimal `json:"availBal"`
	FrozenBalance       decimal.Decimal `json:"frozenBal"`
	OrderFrozen         decimal.Decimal `json:"ordFrozen"`
	Liability           decimal.Decimal `json:"liab"`
	UPL                 decimal.Decimal `json:"upl"`
	UPLLiability        decimal.Decimal `json:"uplLiab"`
	CrossLiability      decimal.Decimal `json:"crossLiab"`
	IsolatedLiability   decimal.Decimal `json:"isoLiab"`
	MarginRatio         decimal.Decimal `json:"mgnRatio"`
	Interest            decimal.Decimal `json:"interest"`
	TWAP                decimal.Decimal `json:"twap"`
	MaxLoan             decimal.Decimal `json:"maxLoan"`
	EquityUSD           decimal.Decimal `json:"eqUsd"`
	NotionalLeverage    decimal.Decimal `json:"notionalLever"`
	StrategyEquity      decimal.Decimal `json:"stgyEq"`
	IsolatedUPL         decimal.Decimal `json:"isoUpl"`
	SpotInUseAmount     decimal.Decimal `json:"spotInUseAmt"`
	SpotIsolatedBalance decimal.Decimal `json:"spotIsoBal"`
	IMR                 decimal.Decimal `json:"imr"`
	MMR                 decimal.Decimal `json:"mmr"`
	SmtSyncEquity       decimal.Decimal `json:"smtSyncEq"`
}

// GetSubAccountFundingBalancesService -- GET /api/v5/asset/subaccount/balances (Read)
//
// Returns a sub-account's funding-account balances (master account).
type GetSubAccountFundingBalancesService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetSubAccountFundingBalancesService(subAcct string) *GetSubAccountFundingBalancesService {
	return &GetSubAccountFundingBalancesService{c: c, params: map[string]string{"subAcct": subAcct}}
}

// SetCcy filters by a single currency.
func (s *GetSubAccountFundingBalancesService) SetCcy(ccy string) *GetSubAccountFundingBalancesService {
	s.params["ccy"] = ccy
	return s
}

func (s *GetSubAccountFundingBalancesService) Do(ctx context.Context) ([]SubAccountFundingBalance, error) {
	req := request.Get(ctx, s.c, "/api/v5/asset/subaccount/balances", s.params).WithSign()
	return request.DoList[SubAccountFundingBalance](req)
}

// SubAccountFundingBalance is one currency's funding-account balance within a
// sub-account.
type SubAccountFundingBalance struct {
	Currency         string          `json:"ccy"`
	Balance          decimal.Decimal `json:"bal"`
	FrozenBalance    decimal.Decimal `json:"frozenBal"`
	AvailableBalance decimal.Decimal `json:"availBal"`
}

// GetSubAccountMaxWithdrawalService -- GET /api/v5/account/subaccount/max-withdrawal (Read)
//
// Returns a sub-account's maximum withdrawal amount per currency (master
// account).
type GetSubAccountMaxWithdrawalService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetSubAccountMaxWithdrawalService(subAcct string) *GetSubAccountMaxWithdrawalService {
	return &GetSubAccountMaxWithdrawalService{c: c, params: map[string]string{"subAcct": subAcct}}
}

// SetCcy filters by a single currency.
func (s *GetSubAccountMaxWithdrawalService) SetCcy(ccy string) *GetSubAccountMaxWithdrawalService {
	s.params["ccy"] = ccy
	return s
}

func (s *GetSubAccountMaxWithdrawalService) Do(ctx context.Context) ([]SubAccountMaxWithdrawal, error) {
	req := request.Get(ctx, s.c, "/api/v5/account/subaccount/max-withdrawal", s.params).WithSign()
	return request.DoList[SubAccountMaxWithdrawal](req)
}

// SubAccountMaxWithdrawal is a currency's maximum withdrawal amounts for a
// sub-account.
type SubAccountMaxWithdrawal struct {
	Currency                string          `json:"ccy"`
	MaxWithdrawal           decimal.Decimal `json:"maxWd"`
	MaxWdEx                 decimal.Decimal `json:"maxWdEx"`
	SpotOffsetMaxWithdrawal decimal.Decimal `json:"spotOffsetMaxWd"`
	SpotOffsetMaxWdEx       decimal.Decimal `json:"spotOffsetMaxWdEx"`
}

// GetSubAccountBillsService -- GET /api/v5/asset/subaccount/bills (Read)
//
// Returns the master-account funding transfer records to/from sub-accounts over
// the last three months.
type GetSubAccountBillsService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetSubAccountBillsService() *GetSubAccountBillsService {
	return &GetSubAccountBillsService{c: c, params: map[string]string{}}
}

// SetCcy filters by currency.
func (s *GetSubAccountBillsService) SetCcy(ccy string) *GetSubAccountBillsService {
	s.params["ccy"] = ccy
	return s
}

// SetType filters by transfer type ("0" master->sub, "1" sub->master).
func (s *GetSubAccountBillsService) SetType(typ string) *GetSubAccountBillsService {
	s.params["type"] = typ
	return s
}

// SetSubAcct filters by sub-account name.
func (s *GetSubAccountBillsService) SetSubAcct(subAcct string) *GetSubAccountBillsService {
	s.params["subAcct"] = subAcct
	return s
}

// SetAfter paginates to records earlier than the given time (older).
func (s *GetSubAccountBillsService) SetAfter(t time.Time) *GetSubAccountBillsService {
	s.params["after"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetBefore paginates to records later than the given time (newer).
func (s *GetSubAccountBillsService) SetBefore(t time.Time) *GetSubAccountBillsService {
	s.params["before"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetSubAccountBillsService) SetLimit(limit int) *GetSubAccountBillsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetSubAccountBillsService) Do(ctx context.Context) ([]SubAccountBill, error) {
	req := request.Get(ctx, s.c, "/api/v5/asset/subaccount/bills", s.params).WithSign()
	return request.DoList[SubAccountBill](req)
}

// SubAccountBill is one master<->sub funding transfer record.
type SubAccountBill struct {
	BillID     string          `json:"billId"`
	Currency   string          `json:"ccy"`
	Amount     decimal.Decimal `json:"amt"`
	Type       string          `json:"type"`
	SubAccount string          `json:"subAcct"`
	Timestamp  time.Time       `json:"ts"`
}

// SubAccountTransferService -- POST /api/v5/asset/subaccount/transfer (Trade)
//
// Transfers funds between the master account and a sub-account, or between two
// sub-accounts (master account). Implement-only: never executed by the live test
// suite.
type SubAccountTransferService struct {
	c    *Client
	body map[string]any
}

// NewSubAccountTransferService builds the transfer. froms/to are the source/dest
// account types ("6" funding, "18" trading); fromSubAccount/toSubAccount are the
// sub-account names.
func (c *Client) NewSubAccountTransferService(ccy string, amt decimal.Decimal, froms, to, fromSubAccount, toSubAccount string) *SubAccountTransferService {
	return &SubAccountTransferService{c: c, body: map[string]any{
		"ccy":            ccy,
		"amt":            amt.String(),
		"from":           froms,
		"to":             to,
		"fromSubAccount": fromSubAccount,
		"toSubAccount":   toSubAccount,
	}}
}

// SetLoanTrans sets whether borrowing/repayment is allowed during the transfer.
func (s *SubAccountTransferService) SetLoanTrans(loanTrans bool) *SubAccountTransferService {
	s.body["loanTrans"] = loanTrans
	return s
}

// SetOmitPosRisk sets whether to ignore position risk ("true"/"false").
func (s *SubAccountTransferService) SetOmitPosRisk(omitPosRisk string) *SubAccountTransferService {
	s.body["omitPosRisk"] = omitPosRisk
	return s
}

func (s *SubAccountTransferService) Do(ctx context.Context) (*SubAccountTransfer, error) {
	req := request.Post(ctx, s.c, "/api/v5/asset/subaccount/transfer", s.body).WithSign()
	return request.DoOne[SubAccountTransfer](req)
}

// SubAccountTransfer is the sub-account transfer acknowledgement.
type SubAccountTransfer struct {
	TransferID string `json:"transId"`
}

// SetSubAccountTransferOutService -- POST /api/v5/users/subaccount/set-transfer-out (Trade)
//
// Sets whether a sub-account is permitted to transfer funds out (master
// account). Implement-only: never executed by the live test suite.
type SetSubAccountTransferOutService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewSetSubAccountTransferOutService(subAcct string, canTransOut bool) *SetSubAccountTransferOutService {
	return &SetSubAccountTransferOutService{c: c, body: map[string]any{
		"subAcct":     subAcct,
		"canTransOut": canTransOut,
	}}
}

func (s *SetSubAccountTransferOutService) Do(ctx context.Context) (*SubAccountTransferOut, error) {
	req := request.Post(ctx, s.c, "/api/v5/users/subaccount/set-transfer-out", s.body).WithSign()
	return request.DoOne[SubAccountTransferOut](req)
}

// SubAccountTransferOut is the set-transfer-out acknowledgement.
type SubAccountTransferOut struct {
	SubAccount     string `json:"subAcct"`
	CanTransferOut bool   `json:"canTransOut"`
}

// GetEntrustSubAccountListService -- GET /api/v5/users/entrust-subaccount-list (Read)
//
// Returns the custody (entrusted) trading sub-accounts the trading-team master
// account manages.
type GetEntrustSubAccountListService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetEntrustSubAccountListService() *GetEntrustSubAccountListService {
	return &GetEntrustSubAccountListService{c: c, params: map[string]string{}}
}

// SetSubAcct filters by sub-account name.
func (s *GetEntrustSubAccountListService) SetSubAcct(subAcct string) *GetEntrustSubAccountListService {
	s.params["subAcct"] = subAcct
	return s
}

func (s *GetEntrustSubAccountListService) Do(ctx context.Context) ([]EntrustSubAccount, error) {
	req := request.Get(ctx, s.c, "/api/v5/users/entrust-subaccount-list", s.params).WithSign()
	return request.DoList[EntrustSubAccount](req)
}

// EntrustSubAccount is one custody (entrusted) trading sub-account.
type EntrustSubAccount struct {
	SubAccount string `json:"subAcct"`
}

// GetManagedSubAccountBillsService -- GET /api/v5/asset/subaccount/managed-subaccount-bills (Read)
//
// Returns the asset-transfer history of managed (custody) sub-accounts, callable
// only by the trading-team master account.
//
// Note: /api/v5/account/subaccount/managed-subaccount-bills does NOT exist (HTTP
// 404); the real path is under /api/v5/asset/.
type GetManagedSubAccountBillsService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetManagedSubAccountBillsService() *GetManagedSubAccountBillsService {
	return &GetManagedSubAccountBillsService{c: c, params: map[string]string{}}
}

// SetCcy filters by currency.
func (s *GetManagedSubAccountBillsService) SetCcy(ccy string) *GetManagedSubAccountBillsService {
	s.params["ccy"] = ccy
	return s
}

// SetType filters by transfer type ("0" transfer in, "1" transfer out).
func (s *GetManagedSubAccountBillsService) SetType(typ string) *GetManagedSubAccountBillsService {
	s.params["type"] = typ
	return s
}

// SetSubAcct filters by managed sub-account name.
func (s *GetManagedSubAccountBillsService) SetSubAcct(subAcct string) *GetManagedSubAccountBillsService {
	s.params["subAcct"] = subAcct
	return s
}

// SetSubUid filters by managed sub-account uid.
func (s *GetManagedSubAccountBillsService) SetSubUid(subUid string) *GetManagedSubAccountBillsService {
	s.params["subUid"] = subUid
	return s
}

// SetAfter paginates to records earlier than the given time (older).
func (s *GetManagedSubAccountBillsService) SetAfter(t time.Time) *GetManagedSubAccountBillsService {
	s.params["after"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetBefore paginates to records later than the given time (newer).
func (s *GetManagedSubAccountBillsService) SetBefore(t time.Time) *GetManagedSubAccountBillsService {
	s.params["before"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetManagedSubAccountBillsService) SetLimit(limit int) *GetManagedSubAccountBillsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetManagedSubAccountBillsService) Do(ctx context.Context) ([]ManagedSubAccountBill, error) {
	req := request.Get(ctx, s.c, "/api/v5/asset/subaccount/managed-subaccount-bills", s.params).WithSign()
	return request.DoList[ManagedSubAccountBill](req)
}

// ManagedSubAccountBill is one managed-sub-account asset-transfer record.
type ManagedSubAccountBill struct {
	BillID     string          `json:"billId"`
	Type       string          `json:"type"`
	Currency   string          `json:"ccy"`
	Amount     decimal.Decimal `json:"amt"`
	SubAccount string          `json:"subAcct"`
	SubUID     string          `json:"subUid"`
	Timestamp  time.Time       `json:"ts"`
}
