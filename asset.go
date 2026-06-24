package okx

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-okx/request"
	"github.com/shopspring/decimal"
)

// AssetAcctType is a funding/transfer account-type code used by the asset
// (funding) endpoints. 6 is the funding account, 18 the unified trading account.
type AssetAcctType string

const (
	AssetAcctTypeFunding AssetAcctType = "6"
	AssetAcctTypeTrading AssetAcctType = "18"
)

// AssetTransferType selects the kind of transfer for /asset/transfer and
// /asset/transfer-state: "0" within own account (default), "1" master->sub,
// "2" sub->master (via master key), "3" sub->master (via sub key),
// "4" sub->sub.
type AssetTransferType string

const (
	AssetTransferTypeWithinAccount       AssetTransferType = "0"
	AssetTransferTypeMasterToSub         AssetTransferType = "1"
	AssetTransferTypeSubToMaster         AssetTransferType = "2"
	AssetTransferTypeSubToMasterBySubKey AssetTransferType = "3"
	AssetTransferTypeSubToSub            AssetTransferType = "4"
)

// GetCurrenciesService -- GET /api/v5/asset/currencies (private)
//
// Returns the list of all currencies (per chain) available to the account,
// including deposit/withdrawal capability flags, fee schedule and quotas.
type GetCurrenciesService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetCurrenciesService() *GetCurrenciesService {
	return &GetCurrenciesService{c: c, params: map[string]string{}}
}

// SetCcy filters by a single (or comma-separated) currency.
func (s *GetCurrenciesService) SetCcy(ccy string) *GetCurrenciesService {
	s.params["ccy"] = ccy
	return s
}

func (s *GetCurrenciesService) Do(ctx context.Context) ([]Currency, error) {
	req := request.Get(ctx, s.c, "/api/v5/asset/currencies", s.params).WithSign()
	return request.DoList[Currency](req)
}

// Currency is a single currency-on-chain entry and its deposit/withdrawal rules.
type Currency struct {
	Currency                    string          `json:"ccy"`
	Name                        string          `json:"name"`
	LogoLink                    string          `json:"logoLink"`
	Chain                       string          `json:"chain"`
	CanDeposit                  bool            `json:"canDep"`
	CanWithdrawal               bool            `json:"canWd"`
	CanInternal                 bool            `json:"canInternal"`
	MinDeposit                  decimal.Decimal `json:"minDep"`
	MinWithdrawal               decimal.Decimal `json:"minWd"`
	MaxWithdrawal               decimal.Decimal `json:"maxWd"`
	WithdrawalTickSize          decimal.Decimal `json:"wdTickSz"`
	WithdrawalQuota             decimal.Decimal `json:"wdQuota"`
	UsedWithdrawalQuota         decimal.Decimal `json:"usedWdQuota"`
	Fee                         decimal.Decimal `json:"fee"`
	MinFee                      decimal.Decimal `json:"minFee"`
	MaxFee                      decimal.Decimal `json:"maxFee"`
	MinFeeForContractAddress    decimal.Decimal `json:"minFeeForCtAddr"`
	MaxFeeForContractAddress    decimal.Decimal `json:"maxFeeForCtAddr"`
	BurningFeeRate              decimal.Decimal `json:"burningFeeRate"`
	MinInternal                 decimal.Decimal `json:"minInternal"`
	MinDepositArrivalConfirm    decimal.Decimal `json:"minDepArrivalConfirm"`
	MinWithdrawalUnlockConfirm  decimal.Decimal `json:"minWdUnlockConfirm"`
	MainNet                     bool            `json:"mainNet"`
	NeedTag                     bool            `json:"needTag"`
	ContractAddress             string          `json:"ctAddr"`
	DepositEstimatedOpenTime    time.Time       `json:"depEstOpenTime"`
	WithdrawalEstimatedOpenTime time.Time       `json:"wdEstOpenTime"`
	DepositQuotaFixed           decimal.Decimal `json:"depQuotaFixed"`
	UsedDepositQuotaFixed       decimal.Decimal `json:"usedDepQuotaFixed"`
	DepositQuoteDailyLayer2     decimal.Decimal `json:"depQuoteDailyLayer2"`
	StablecoinDailyQuota        decimal.Decimal `json:"stablecoinDailyQuota"`
	StablecoinMonthlyQuota      decimal.Decimal `json:"stablecoinMonthlyQuota"`
	UsedStablecoinDailyQuota    decimal.Decimal `json:"usedStablecoinDailyQuota"`
	UsedStablecoinMonthlyQuota  decimal.Decimal `json:"usedStablecoinMonthlyQuota"`
}

// GetFundingBalanceService -- GET /api/v5/asset/balances (private)
//
// Returns the balances in the funding account.
type GetFundingBalanceService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetFundingBalanceService() *GetFundingBalanceService {
	return &GetFundingBalanceService{c: c, params: map[string]string{}}
}

// SetCcy filters by a single (or comma-separated) currency.
func (s *GetFundingBalanceService) SetCcy(ccy string) *GetFundingBalanceService {
	s.params["ccy"] = ccy
	return s
}

func (s *GetFundingBalanceService) Do(ctx context.Context) ([]FundingBalance, error) {
	req := request.Get(ctx, s.c, "/api/v5/asset/balances", s.params).WithSign()
	return request.DoList[FundingBalance](req)
}

// FundingBalance is a currency's balance in the funding account.
type FundingBalance struct {
	Currency         string          `json:"ccy"`
	Balance          decimal.Decimal `json:"bal"`
	FrozenBalance    decimal.Decimal `json:"frozenBal"`
	AvailableBalance decimal.Decimal `json:"availBal"`
}

// GetNonTradableAssetsService -- GET /api/v5/asset/non-tradable-assets (private)
//
// Returns the non-tradable assets held in the funding account (assets that can
// only be withdrawn, not traded).
type GetNonTradableAssetsService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetNonTradableAssetsService() *GetNonTradableAssetsService {
	return &GetNonTradableAssetsService{c: c, params: map[string]string{}}
}

// SetCcy filters by a single (or comma-separated) currency.
func (s *GetNonTradableAssetsService) SetCcy(ccy string) *GetNonTradableAssetsService {
	s.params["ccy"] = ccy
	return s
}

func (s *GetNonTradableAssetsService) Do(ctx context.Context) ([]NonTradableAsset, error) {
	req := request.Get(ctx, s.c, "/api/v5/asset/non-tradable-assets", s.params).WithSign()
	return request.DoList[NonTradableAsset](req)
}

// NonTradableAsset is a single non-tradable asset and its withdrawal rules.
type NonTradableAsset struct {
	Currency           string          `json:"ccy"`
	Name               string          `json:"name"`
	LogoLink           string          `json:"logoLink"`
	Balance            decimal.Decimal `json:"bal"`
	CanWithdrawal      bool            `json:"canWd"`
	Chain              string          `json:"chain"`
	ContractAddress    string          `json:"ctAddr"`
	Fee                decimal.Decimal `json:"fee"`
	FeeCurrency        string          `json:"feeCcy"`
	MinWithdrawal      decimal.Decimal `json:"minWd"`
	WithdrawalTickSize decimal.Decimal `json:"wdTickSz"`
	WithdrawalAll      bool            `json:"wdAll"`
	NeedTag            bool            `json:"needTag"`
	BurningFeeRate     decimal.Decimal `json:"burningFeeRate"`
}

// GetAssetValuationService -- GET /api/v5/asset/asset-valuation (private)
//
// Returns the total valuation of all assets across the account's sub-accounts
// (funding, trading, classic, earn), denominated in a chosen currency.
type GetAssetValuationService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetAssetValuationService() *GetAssetValuationService {
	return &GetAssetValuationService{c: c, params: map[string]string{}}
}

// SetCcy sets the valuation currency (e.g. BTC, USDT; default BTC).
func (s *GetAssetValuationService) SetCcy(ccy string) *GetAssetValuationService {
	s.params["ccy"] = ccy
	return s
}

func (s *GetAssetValuationService) Do(ctx context.Context) (*AssetValuation, error) {
	req := request.Get(ctx, s.c, "/api/v5/asset/asset-valuation", s.params).WithSign()
	return request.DoOne[AssetValuation](req)
}

// AssetValuation is the account-wide asset valuation and its per-account-type
// breakdown.
type AssetValuation struct {
	TotalBalance decimal.Decimal       `json:"totalBal"`
	Timestamp    time.Time             `json:"ts"`
	Details      AssetValuationDetails `json:"details"`
}

// AssetValuationDetails breaks the valuation down by account type.
type AssetValuationDetails struct {
	Funding decimal.Decimal `json:"funding"`
	Trading decimal.Decimal `json:"trading"`
	Classic decimal.Decimal `json:"classic"`
	Earn    decimal.Decimal `json:"earn"`
}

// FundsTransferService -- POST /api/v5/asset/transfer (private)
//
// Transfers funds between the funding account and the trading account, or
// between a master account and its sub-accounts. IMPLEMENT-ONLY: this moves real
// funds and must never be executed by the test suite.
type FundsTransferService struct {
	c    *Client
	body map[string]any
}

// NewFundsTransferService builds a transfer. ccy is the currency, amt the
// amount, from/to the source and destination account types.
func (c *Client) NewFundsTransferService(ccy string, amt decimal.Decimal, from, to AssetAcctType) *FundsTransferService {
	return &FundsTransferService{c: c, body: map[string]any{
		"ccy":  ccy,
		"amt":  amt.String(),
		"from": string(from),
		"to":   string(to),
	}}
}

// SetType sets the transfer type (default "0" within own account).
func (s *FundsTransferService) SetType(typ AssetTransferType) *FundsTransferService {
	s.body["type"] = string(typ)
	return s
}

// SetSubAcct sets the sub-account name (required for master<->sub transfers).
func (s *FundsTransferService) SetSubAcct(subAcct string) *FundsTransferService {
	s.body["subAcct"] = subAcct
	return s
}

// SetLoanTrans enables transfer of borrowed funds under Spot mode / Multi-currency
// margin / Portfolio margin.
func (s *FundsTransferService) SetLoanTrans(loanTrans bool) *FundsTransferService {
	s.body["loanTrans"] = loanTrans
	return s
}

// SetOmitPosRisk skips the position-risk check (Portfolio margin only).
func (s *FundsTransferService) SetOmitPosRisk(omitPosRisk string) *FundsTransferService {
	s.body["omitPosRisk"] = omitPosRisk
	return s
}

// SetClientId sets a client-supplied transfer id.
func (s *FundsTransferService) SetClientId(clientId string) *FundsTransferService {
	s.body["clientId"] = clientId
	return s
}

func (s *FundsTransferService) Do(ctx context.Context) (*FundsTransfer, error) {
	req := request.Post(ctx, s.c, "/api/v5/asset/transfer", s.body).WithSign()
	return request.DoOne[FundsTransfer](req)
}

// FundsTransfer is the acknowledgement of a funds transfer.
type FundsTransfer struct {
	TransferID string          `json:"transId"`
	ClientID   string          `json:"clientId"`
	Currency   string          `json:"ccy"`
	Amount     decimal.Decimal `json:"amt"`
	From       AssetAcctType   `json:"from"`
	To         AssetAcctType   `json:"to"`
}

// GetTransferStateService -- GET /api/v5/asset/transfer-state (private)
//
// Returns the status of a funds transfer by transfer id or client id.
type GetTransferStateService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetTransferStateService() *GetTransferStateService {
	return &GetTransferStateService{c: c, params: map[string]string{}}
}

// SetTransId filters by transfer id (mutually exclusive with clientId).
func (s *GetTransferStateService) SetTransId(transId string) *GetTransferStateService {
	s.params["transId"] = transId
	return s
}

// SetClientId filters by client-supplied transfer id.
func (s *GetTransferStateService) SetClientId(clientId string) *GetTransferStateService {
	s.params["clientId"] = clientId
	return s
}

// SetType sets the transfer type (default "0").
func (s *GetTransferStateService) SetType(typ AssetTransferType) *GetTransferStateService {
	s.params["type"] = string(typ)
	return s
}

func (s *GetTransferStateService) Do(ctx context.Context) ([]TransferState, error) {
	req := request.Get(ctx, s.c, "/api/v5/asset/transfer-state", s.params).WithSign()
	return request.DoList[TransferState](req)
}

// TransferState is the status of a funds transfer.
type TransferState struct {
	TransferID     string            `json:"transId"`
	ClientID       string            `json:"clientId"`
	Currency       string            `json:"ccy"`
	Amount         decimal.Decimal   `json:"amt"`
	Type           AssetTransferType `json:"type"`
	From           AssetAcctType     `json:"from"`
	To             AssetAcctType     `json:"to"`
	SubAccount     string            `json:"subAcct"`
	InstrumentID   string            `json:"instId"`
	ToInstrumentID string            `json:"toInstId"`
	State          string            `json:"state"`
}

// GetAssetBillsService -- GET /api/v5/asset/bills (private)
//
// Returns the funding-account bill (balance-change) history of the last year.
type GetAssetBillsService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetAssetBillsService() *GetAssetBillsService {
	return &GetAssetBillsService{c: c, params: map[string]string{}}
}

// SetCcy filters by currency.
func (s *GetAssetBillsService) SetCcy(ccy string) *GetAssetBillsService {
	s.params["ccy"] = ccy
	return s
}

// SetType filters by bill type (see OKX docs for the numeric type codes).
func (s *GetAssetBillsService) SetType(typ string) *GetAssetBillsService {
	s.params["type"] = typ
	return s
}

// SetClientId filters by client-supplied id.
func (s *GetAssetBillsService) SetClientId(clientId string) *GetAssetBillsService {
	s.params["clientId"] = clientId
	return s
}

// SetAfter paginates to records earlier than the given time (older).
func (s *GetAssetBillsService) SetAfter(t time.Time) *GetAssetBillsService {
	s.params["after"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetBefore paginates to records later than the given time (newer).
func (s *GetAssetBillsService) SetBefore(t time.Time) *GetAssetBillsService {
	s.params["before"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetAssetBillsService) SetLimit(limit int) *GetAssetBillsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetAssetBillsService) Do(ctx context.Context) ([]AssetBill, error) {
	req := request.Get(ctx, s.c, "/api/v5/asset/bills", s.params).WithSign()
	return request.DoList[AssetBill](req)
}

// AssetBill is one funding-account balance-change record.
type AssetBill struct {
	BillID        string          `json:"billId"`
	Currency      string          `json:"ccy"`
	ClientID      string          `json:"clientId"`
	BalanceChange decimal.Decimal `json:"balChg"`
	Balance       decimal.Decimal `json:"bal"`
	Type          string          `json:"type"`
	Notes         string          `json:"notes"`
	Timestamp     time.Time       `json:"ts"`
}

// GetDepositAddressService -- GET /api/v5/asset/deposit-address (private)
//
// Returns the deposit addresses of a currency (one per chain / sub-address).
type GetDepositAddressService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetDepositAddressService(ccy string) *GetDepositAddressService {
	return &GetDepositAddressService{c: c, params: map[string]string{"ccy": ccy}}
}

func (s *GetDepositAddressService) Do(ctx context.Context) ([]DepositAddress, error) {
	req := request.Get(ctx, s.c, "/api/v5/asset/deposit-address", s.params).WithSign()
	return request.DoList[DepositAddress](req)
}

// DepositAddress is a single deposit address for a currency on a chain.
type DepositAddress struct {
	Address         string        `json:"addr"`
	Tag             string        `json:"tag"`
	Memo            string        `json:"memo"`
	PaymentID       string        `json:"pmtId"`
	AddrEx          any           `json:"addrEx"`
	Currency        string        `json:"ccy"`
	Chain           string        `json:"chain"`
	To              AssetAcctType `json:"to"`
	Selected        bool          `json:"selected"`
	ContractAddress string        `json:"ctAddr"`
	VerifiedName    string        `json:"verifiedName"`
}

// GetDepositHistoryService -- GET /api/v5/asset/deposit-history (private)
//
// Returns the deposit records of the account.
type GetDepositHistoryService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetDepositHistoryService() *GetDepositHistoryService {
	return &GetDepositHistoryService{c: c, params: map[string]string{}}
}

// SetCcy filters by currency.
func (s *GetDepositHistoryService) SetCcy(ccy string) *GetDepositHistoryService {
	s.params["ccy"] = ccy
	return s
}

// SetDepId filters by deposit id.
func (s *GetDepositHistoryService) SetDepId(depId string) *GetDepositHistoryService {
	s.params["depId"] = depId
	return s
}

// SetFromWdId filters by the internal-transfer withdrawal id of the sender.
func (s *GetDepositHistoryService) SetFromWdId(fromWdId string) *GetDepositHistoryService {
	s.params["fromWdId"] = fromWdId
	return s
}

// SetTxId filters by the on-chain transaction hash.
func (s *GetDepositHistoryService) SetTxId(txId string) *GetDepositHistoryService {
	s.params["txId"] = txId
	return s
}

// SetType filters by deposit type code.
func (s *GetDepositHistoryService) SetType(typ string) *GetDepositHistoryService {
	s.params["type"] = typ
	return s
}

// SetState filters by deposit state code.
func (s *GetDepositHistoryService) SetState(state string) *GetDepositHistoryService {
	s.params["state"] = state
	return s
}

// SetAfter paginates to records earlier than the given time (older).
func (s *GetDepositHistoryService) SetAfter(t time.Time) *GetDepositHistoryService {
	s.params["after"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetBefore paginates to records later than the given time (newer).
func (s *GetDepositHistoryService) SetBefore(t time.Time) *GetDepositHistoryService {
	s.params["before"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetDepositHistoryService) SetLimit(limit int) *GetDepositHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetDepositHistoryService) Do(ctx context.Context) ([]DepositHistory, error) {
	req := request.Get(ctx, s.c, "/api/v5/asset/deposit-history", s.params).WithSign()
	return request.DoList[DepositHistory](req)
}

// DepositHistory is one deposit record.
type DepositHistory struct {
	Currency                  string          `json:"ccy"`
	Chain                     string          `json:"chain"`
	Amount                    decimal.Decimal `json:"amt"`
	From                      string          `json:"from"`
	AreaCodeFrom              string          `json:"areaCodeFrom"`
	To                        string          `json:"to"`
	TransactionID             string          `json:"txId"`
	DepositID                 string          `json:"depId"`
	FromWithdrawalID          string          `json:"fromWdId"`
	State                     string          `json:"state"`
	ActualDepositBlockConfirm decimal.Decimal `json:"actualDepBlkConfirm"`
	Timestamp                 time.Time       `json:"ts"`
}

// WithdrawalService -- POST /api/v5/asset/withdrawal (private)
//
// Submits a withdrawal (on-chain or internal). IMPLEMENT-ONLY: this withdraws
// real funds and must NEVER be executed.
type WithdrawalService struct {
	c    *Client
	body map[string]any
}

// NewWithdrawalService builds a withdrawal. ccy is the currency, amt the amount,
// dest the destination ("3" internal, "4" on-chain), toAddr the recipient
// address (suffix ":tag" / ":memo" when required).
func (c *Client) NewWithdrawalService(ccy string, amt decimal.Decimal, dest, toAddr string) *WithdrawalService {
	return &WithdrawalService{c: c, body: map[string]any{
		"ccy":    ccy,
		"amt":    amt.String(),
		"dest":   dest,
		"toAddr": toAddr,
	}}
}

// SetChain sets the chain (e.g. "USDT-ERC20"); required for on-chain withdrawals
// of multi-chain currencies.
func (s *WithdrawalService) SetChain(chain string) *WithdrawalService {
	s.body["chain"] = chain
	return s
}

// SetFee sets the network fee (on-chain only).
func (s *WithdrawalService) SetFee(fee decimal.Decimal) *WithdrawalService {
	s.body["fee"] = fee.String()
	return s
}

// SetAreaCode sets the mobile area code (required for some internal withdrawals).
func (s *WithdrawalService) SetAreaCode(areaCode string) *WithdrawalService {
	s.body["areaCode"] = areaCode
	return s
}

// SetRcvrInfo sets the recipient information (required by some jurisdictions).
func (s *WithdrawalService) SetRcvrInfo(rcvrInfo any) *WithdrawalService {
	s.body["rcvrInfo"] = rcvrInfo
	return s
}

// SetClientId sets a client-supplied withdrawal id.
func (s *WithdrawalService) SetClientId(clientId string) *WithdrawalService {
	s.body["clientId"] = clientId
	return s
}

func (s *WithdrawalService) Do(ctx context.Context) (*Withdrawal, error) {
	req := request.Post(ctx, s.c, "/api/v5/asset/withdrawal", s.body).WithSign()
	return request.DoOne[Withdrawal](req)
}

// Withdrawal is the acknowledgement of a submitted withdrawal.
type Withdrawal struct {
	Currency     string          `json:"ccy"`
	Chain        string          `json:"chain"`
	Amount       decimal.Decimal `json:"amt"`
	WithdrawalID string          `json:"wdId"`
	ClientID     string          `json:"clientId"`
}

// CancelWithdrawalService -- POST /api/v5/asset/cancel-withdrawal (private)
//
// Cancels a pending withdrawal. IMPLEMENT-ONLY: acts on a real withdrawal and is
// not executed by the test suite.
type CancelWithdrawalService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewCancelWithdrawalService(wdId string) *CancelWithdrawalService {
	return &CancelWithdrawalService{c: c, body: map[string]any{"wdId": wdId}}
}

func (s *CancelWithdrawalService) Do(ctx context.Context) (*CancelWithdrawal, error) {
	req := request.Post(ctx, s.c, "/api/v5/asset/cancel-withdrawal", s.body).WithSign()
	return request.DoOne[CancelWithdrawal](req)
}

// CancelWithdrawal is the acknowledgement of a withdrawal cancellation.
type CancelWithdrawal struct {
	WithdrawalID string `json:"wdId"`
}

// GetWithdrawalHistoryService -- GET /api/v5/asset/withdrawal-history (private)
//
// Returns the withdrawal records of the account.
type GetWithdrawalHistoryService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetWithdrawalHistoryService() *GetWithdrawalHistoryService {
	return &GetWithdrawalHistoryService{c: c, params: map[string]string{}}
}

// SetCcy filters by currency.
func (s *GetWithdrawalHistoryService) SetCcy(ccy string) *GetWithdrawalHistoryService {
	s.params["ccy"] = ccy
	return s
}

// SetWdId filters by withdrawal id.
func (s *GetWithdrawalHistoryService) SetWdId(wdId string) *GetWithdrawalHistoryService {
	s.params["wdId"] = wdId
	return s
}

// SetClientId filters by client-supplied id.
func (s *GetWithdrawalHistoryService) SetClientId(clientId string) *GetWithdrawalHistoryService {
	s.params["clientId"] = clientId
	return s
}

// SetTxId filters by the on-chain transaction hash.
func (s *GetWithdrawalHistoryService) SetTxId(txId string) *GetWithdrawalHistoryService {
	s.params["txId"] = txId
	return s
}

// SetType filters by withdrawal type code.
func (s *GetWithdrawalHistoryService) SetType(typ string) *GetWithdrawalHistoryService {
	s.params["type"] = typ
	return s
}

// SetState filters by withdrawal state code.
func (s *GetWithdrawalHistoryService) SetState(state string) *GetWithdrawalHistoryService {
	s.params["state"] = state
	return s
}

// SetAfter paginates to records earlier than the given time (older).
func (s *GetWithdrawalHistoryService) SetAfter(t time.Time) *GetWithdrawalHistoryService {
	s.params["after"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetBefore paginates to records later than the given time (newer).
func (s *GetWithdrawalHistoryService) SetBefore(t time.Time) *GetWithdrawalHistoryService {
	s.params["before"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetWithdrawalHistoryService) SetLimit(limit int) *GetWithdrawalHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetWithdrawalHistoryService) Do(ctx context.Context) ([]WithdrawalHistory, error) {
	req := request.Get(ctx, s.c, "/api/v5/asset/withdrawal-history", s.params).WithSign()
	return request.DoList[WithdrawalHistory](req)
}

// WithdrawalHistory is one withdrawal record.
type WithdrawalHistory struct {
	Currency         string          `json:"ccy"`
	Chain            string          `json:"chain"`
	NonTradableAsset bool            `json:"nonTradableAsset"`
	Amount           decimal.Decimal `json:"amt"`
	Timestamp        time.Time       `json:"ts"`
	From             string          `json:"from"`
	AreaCodeFrom     string          `json:"areaCodeFrom"`
	To               string          `json:"to"`
	AreaCodeTo       string          `json:"areaCodeTo"`
	Memo             string          `json:"memo"`
	ToAddressType    string          `json:"toAddrType"`
	TransactionID    string          `json:"txId"`
	Fee              decimal.Decimal `json:"fee"`
	FeeCurrency      string          `json:"feeCcy"`
	State            string          `json:"state"`
	WithdrawalID     string          `json:"wdId"`
	ClientID         string          `json:"clientId"`
	Note             string          `json:"note"`
}

// GetDepositWithdrawStatusService -- GET /api/v5/asset/deposit-withdraw-status (private)
//
// Returns the on-chain status of a deposit or withdrawal. Query by withdrawal id
// (wdId) for a withdrawal, or by txId + ccy + to + chain for a deposit.
type GetDepositWithdrawStatusService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetDepositWithdrawStatusService() *GetDepositWithdrawStatusService {
	return &GetDepositWithdrawStatusService{c: c, params: map[string]string{}}
}

// SetWdId queries by withdrawal id.
func (s *GetDepositWithdrawStatusService) SetWdId(wdId string) *GetDepositWithdrawStatusService {
	s.params["wdId"] = wdId
	return s
}

// SetTxId queries by the on-chain transaction hash (deposit; with ccy/to/chain).
func (s *GetDepositWithdrawStatusService) SetTxId(txId string) *GetDepositWithdrawStatusService {
	s.params["txId"] = txId
	return s
}

// SetCcy sets the currency (deposit query).
func (s *GetDepositWithdrawStatusService) SetCcy(ccy string) *GetDepositWithdrawStatusService {
	s.params["ccy"] = ccy
	return s
}

// SetTo sets the receiving address (deposit query).
func (s *GetDepositWithdrawStatusService) SetTo(to string) *GetDepositWithdrawStatusService {
	s.params["to"] = to
	return s
}

// SetChain sets the chain (deposit query).
func (s *GetDepositWithdrawStatusService) SetChain(chain string) *GetDepositWithdrawStatusService {
	s.params["chain"] = chain
	return s
}

func (s *GetDepositWithdrawStatusService) Do(ctx context.Context) ([]DepositWithdrawStatus, error) {
	req := request.Get(ctx, s.c, "/api/v5/asset/deposit-withdraw-status", s.params).WithSign()
	return request.DoList[DepositWithdrawStatus](req)
}

// DepositWithdrawStatus is the on-chain progress of a deposit/withdrawal.
type DepositWithdrawStatus struct {
	WithdrawalID          string    `json:"wdId"`
	TransactionID         string    `json:"txId"`
	State                 string    `json:"state"`
	EstimatedCompleteTime time.Time `json:"estCompleteTime"`
}

// GetExchangeListService -- GET /api/v5/asset/exchange-list (private)
//
// Returns the list of supported exchanges and their ids (used to tag
// counterparty exchanges for travel-rule withdrawals).
type GetExchangeListService struct {
	c *Client
}

func (c *Client) NewGetExchangeListService() *GetExchangeListService {
	return &GetExchangeListService{c: c}
}

func (s *GetExchangeListService) Do(ctx context.Context) ([]Exchange, error) {
	req := request.Get(ctx, s.c, "/api/v5/asset/exchange-list").WithSign()
	return request.DoList[Exchange](req)
}

// Exchange is a supported counterparty exchange.
type Exchange struct {
	ExchID   string `json:"exchId"`
	ExchName string `json:"exchName"`
}

// ApplyMonthlyStatementService -- POST /api/v5/asset/monthly-statement (private)
//
// Applies for a downloadable monthly transaction statement. IMPLEMENT-ONLY: this
// creates an export job on a real account and is not executed by the test suite.
type ApplyMonthlyStatementService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewApplyMonthlyStatementService() *ApplyMonthlyStatementService {
	return &ApplyMonthlyStatementService{c: c, body: map[string]any{}}
}

// SetMonth sets the statement month (e.g. "Jan"); defaults to the previous month.
func (s *ApplyMonthlyStatementService) SetMonth(month string) *ApplyMonthlyStatementService {
	s.body["month"] = month
	return s
}

func (s *ApplyMonthlyStatementService) Do(ctx context.Context) (*MonthlyStatementApply, error) {
	req := request.Post(ctx, s.c, "/api/v5/asset/monthly-statement", s.body).WithSign()
	return request.DoOne[MonthlyStatementApply](req)
}

// MonthlyStatementApply is the acknowledgement of a monthly-statement request.
type MonthlyStatementApply struct {
	Timestamp time.Time `json:"ts"`
}

// GetMonthlyStatementService -- GET /api/v5/asset/monthly-statement (private)
//
// Returns the download link for a previously applied-for monthly statement.
type GetMonthlyStatementService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetMonthlyStatementService(month string) *GetMonthlyStatementService {
	return &GetMonthlyStatementService{c: c, params: map[string]string{"month": month}}
}

func (s *GetMonthlyStatementService) Do(ctx context.Context) (*MonthlyStatement, error) {
	req := request.Get(ctx, s.c, "/api/v5/asset/monthly-statement", s.params).WithSign()
	return request.DoOne[MonthlyStatement](req)
}

// MonthlyStatement is the download link and state of a monthly statement.
type MonthlyStatement struct {
	FileHref  string    `json:"fileHref"`
	State     string    `json:"state"`
	Timestamp time.Time `json:"ts"`
}
