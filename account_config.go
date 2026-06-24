package okx

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-okx/request"
	"github.com/shopspring/decimal"
)

// PosMode is the account-wide position mode (long/short hedging vs net).
type PosMode string

const (
	PosModeLongShort PosMode = "long_short_mode"
	PosModeNet       PosMode = "net_mode"
)

// GreeksType is the display type of option greeks: "PA" (in coins) or "BS"
// (Black-Scholes, in dollars).
type GreeksType string

const (
	GreeksTypePA GreeksType = "PA"
	GreeksTypeBS GreeksType = "BS"
)

// IsoMode is the isolated-margin transfer behaviour.
type IsoMode string

const (
	IsoModeAutomatic        IsoMode = "automatic"
	IsoModeAutonomy         IsoMode = "autonomy"
	IsoModeAutoTransfersCcy IsoMode = "auto_transfers_ccy"
)

// IsoModeType is the instrument scope of an isolated-margin setting.
type IsoModeType string

const (
	IsoModeTypeMargin    IsoModeType = "MARGIN"
	IsoModeTypeContracts IsoModeType = "CONTRACTS"
)

// MarginBalanceType selects whether a position-margin adjustment adds or reduces
// margin.
type MarginBalanceType string

const (
	MarginBalanceTypeAdd    MarginBalanceType = "add"
	MarginBalanceTypeReduce MarginBalanceType = "reduce"
)

// CollateralAssetType selects whether a set-collateral-assets request applies to
// all currencies or to a custom currency list.
type CollateralAssetType string

const (
	CollateralAssetTypeAll    CollateralAssetType = "all"
	CollateralAssetTypeCustom CollateralAssetType = "custom"
)

// SetPositionModeService -- POST /api/v5/account/set-position-mode (Trade)
//
// Sets the account-wide position mode (long/short hedging or net) for
// derivatives. Single-currency margin accounts only.
type SetPositionModeService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewSetPositionModeService(posMode PosMode) *SetPositionModeService {
	return &SetPositionModeService{c: c, body: map[string]any{"posMode": string(posMode)}}
}

func (s *SetPositionModeService) Do(ctx context.Context) (*PositionMode, error) {
	req := request.Post(ctx, s.c, "/api/v5/account/set-position-mode", s.body).WithSign()
	return request.DoOne[PositionMode](req)
}

// PositionMode is the set-position-mode acknowledgement.
type PositionMode struct {
	PositionMode PosMode `json:"posMode"`
}

// SetLeverageService -- POST /api/v5/account/set-leverage (Trade)
//
// Sets the leverage of an instrument or currency. Either instId or ccy is
// required (instId for an instrument's leverage; ccy for cross-margin currency
// leverage in Futures mode).
type SetLeverageService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewSetLeverageService(lever string, mgnMode MgnMode) *SetLeverageService {
	return &SetLeverageService{c: c, body: map[string]any{
		"lever":   lever,
		"mgnMode": string(mgnMode),
	}}
}

// SetInstId sets the instrument whose leverage is changed.
func (s *SetLeverageService) SetInstId(instId string) *SetLeverageService {
	s.body["instId"] = instId
	return s
}

// SetCcy sets the currency whose cross-margin leverage is changed (Futures mode).
func (s *SetLeverageService) SetCcy(ccy string) *SetLeverageService {
	s.body["ccy"] = ccy
	return s
}

// SetPosSide sets the position side (long/short) for isolated long/short mode.
func (s *SetLeverageService) SetPosSide(posSide PosSide) *SetLeverageService {
	s.body["posSide"] = string(posSide)
	return s
}

func (s *SetLeverageService) Do(ctx context.Context) (*LeverageSetting, error) {
	req := request.Post(ctx, s.c, "/api/v5/account/set-leverage", s.body).WithSign()
	return request.DoOne[LeverageSetting](req)
}

// LeverageSetting is the set-leverage acknowledgement.
type LeverageSetting struct {
	Leverage     decimal.Decimal `json:"lever"`
	MarginMode   MgnMode         `json:"mgnMode"`
	InstrumentID string          `json:"instId"`
	PositionSide PosSide         `json:"posSide"`
}

// SetGreeksService -- POST /api/v5/account/set-greeks (Trade)
//
// Sets the display type of option greeks ("PA" in coins or "BS" in dollars).
type SetGreeksService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewSetGreeksService(greeksType GreeksType) *SetGreeksService {
	return &SetGreeksService{c: c, body: map[string]any{"greeksType": string(greeksType)}}
}

func (s *SetGreeksService) Do(ctx context.Context) (*GreeksSetting, error) {
	req := request.Post(ctx, s.c, "/api/v5/account/set-greeks", s.body).WithSign()
	return request.DoOne[GreeksSetting](req)
}

// GreeksSetting is the set-greeks acknowledgement.
type GreeksSetting struct {
	GreeksType GreeksType `json:"greeksType"`
}

// SetIsolatedModeService -- POST /api/v5/account/set-isolated-mode (Trade)
//
// Sets the isolated-margin transfer behaviour (automatic vs autonomy) for the
// given instrument type (MARGIN or CONTRACTS).
type SetIsolatedModeService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewSetIsolatedModeService(isoMode IsoMode, typ IsoModeType) *SetIsolatedModeService {
	return &SetIsolatedModeService{c: c, body: map[string]any{
		"isoMode": string(isoMode),
		"type":    string(typ),
	}}
}

func (s *SetIsolatedModeService) Do(ctx context.Context) (*IsolatedModeSetting, error) {
	req := request.Post(ctx, s.c, "/api/v5/account/set-isolated-mode", s.body).WithSign()
	return request.DoOne[IsolatedModeSetting](req)
}

// IsolatedModeSetting is the set-isolated-mode acknowledgement.
type IsolatedModeSetting struct {
	IsolatedMode IsoMode `json:"isoMode"`
}

// SetCollateralAssetsService -- POST /api/v5/account/set-collateral-assets (Trade)
//
// Enables or disables assets as collateral, either for all currencies (type
// "all") or for a custom currency list (type "custom", then set CcyList).
type SetCollateralAssetsService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewSetCollateralAssetsService(typ CollateralAssetType, collateralEnabled bool) *SetCollateralAssetsService {
	return &SetCollateralAssetsService{c: c, body: map[string]any{
		"type":              string(typ),
		"collateralEnabled": collateralEnabled,
	}}
}

// SetCcyList sets the currency list (required when type is "custom").
func (s *SetCollateralAssetsService) SetCcyList(ccyList []string) *SetCollateralAssetsService {
	s.body["ccyList"] = ccyList
	return s
}

func (s *SetCollateralAssetsService) Do(ctx context.Context) (*CollateralAssetsSetting, error) {
	req := request.Post(ctx, s.c, "/api/v5/account/set-collateral-assets", s.body).WithSign()
	return request.DoOne[CollateralAssetsSetting](req)
}

// CollateralAssetsSetting is the set-collateral-assets acknowledgement.
type CollateralAssetsSetting struct {
	Type              CollateralAssetType `json:"type"`
	CollateralEnabled bool                `json:"collateralEnabled"`
	CurrencyList      []string            `json:"ccyList"`
}

// GetCollateralAssetsService -- GET /api/v5/account/collateral-assets (Read)
//
// Returns the per-currency collateral-enabled state of the account.
type GetCollateralAssetsService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetCollateralAssetsService() *GetCollateralAssetsService {
	return &GetCollateralAssetsService{c: c, params: map[string]string{}}
}

// SetCcy filters by a single currency.
func (s *GetCollateralAssetsService) SetCcy(ccy string) *GetCollateralAssetsService {
	s.params["ccy"] = ccy
	return s
}

// SetCollateralEnabled filters by collateral-enabled state.
func (s *GetCollateralAssetsService) SetCollateralEnabled(enabled bool) *GetCollateralAssetsService {
	s.params["collateralEnabled"] = strconv.FormatBool(enabled)
	return s
}

func (s *GetCollateralAssetsService) Do(ctx context.Context) ([]CollateralAsset, error) {
	req := request.Get(ctx, s.c, "/api/v5/account/collateral-assets", s.params).WithSign()
	return request.DoList[CollateralAsset](req)
}

// CollateralAsset is a currency's collateral-enabled state.
type CollateralAsset struct {
	Currency          string `json:"ccy"`
	CollateralEnabled bool   `json:"collateralEnabled"`
}

// SetAccountLevelService -- POST /api/v5/account/set-account-level (Trade)
//
// Sets the account mode (acctLv: "1" spot, "2" spot&futures, "3" multi-currency
// margin, "4" portfolio margin). DANGEROUS: this changes the whole account's
// trading mode.
type SetAccountLevelService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewSetAccountLevelService(acctLv string) *SetAccountLevelService {
	return &SetAccountLevelService{c: c, body: map[string]any{"acctLv": acctLv}}
}

func (s *SetAccountLevelService) Do(ctx context.Context) (*AccountLevelSetting, error) {
	req := request.Post(ctx, s.c, "/api/v5/account/set-account-level", s.body).WithSign()
	return request.DoOne[AccountLevelSetting](req)
}

// AccountLevelSetting is the set-account-level acknowledgement.
type AccountLevelSetting struct {
	AccountLevel string `json:"acctLv"`
}

// ActivateOptionService -- POST /api/v5/account/activate-option (Trade)
//
// Activates option trading for the account. Returns the activation timestamp.
type ActivateOptionService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewActivateOptionService() *ActivateOptionService {
	return &ActivateOptionService{c: c, body: map[string]any{}}
}

func (s *ActivateOptionService) Do(ctx context.Context) (*ActivateOption, error) {
	req := request.Post(ctx, s.c, "/api/v5/account/activate-option", s.body).WithSign()
	return request.DoOne[ActivateOption](req)
}

// ActivateOption is the activate-option acknowledgement. Some account states
// also return a "result" boolean.
type ActivateOption struct {
	Result    bool      `json:"result"`
	Timestamp time.Time `json:"ts"`
}

// SetAutoLoanService -- POST /api/v5/account/set-auto-loan (Trade)
//
// Enables or disables automatic margin borrowing for the account.
type SetAutoLoanService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewSetAutoLoanService(autoLoan bool) *SetAutoLoanService {
	return &SetAutoLoanService{c: c, body: map[string]any{"autoLoan": autoLoan}}
}

func (s *SetAutoLoanService) Do(ctx context.Context) (*AutoLoanSetting, error) {
	req := request.Post(ctx, s.c, "/api/v5/account/set-auto-loan", s.body).WithSign()
	return request.DoOne[AutoLoanSetting](req)
}

// AutoLoanSetting is the set-auto-loan acknowledgement.
type AutoLoanSetting struct {
	AutoLoan bool `json:"autoLoan"`
}

// SetMarginBalanceService -- POST /api/v5/account/position/margin-balance (Trade)
//
// Increases or decreases the margin of an isolated position.
type SetMarginBalanceService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewSetMarginBalanceService(instId string, posSide PosSide, typ MarginBalanceType, amt string) *SetMarginBalanceService {
	return &SetMarginBalanceService{c: c, body: map[string]any{
		"instId":  instId,
		"posSide": string(posSide),
		"type":    string(typ),
		"amt":     amt,
	}}
}

// SetCcy sets the currency (applicable to isolated MARGIN positions).
func (s *SetMarginBalanceService) SetCcy(ccy string) *SetMarginBalanceService {
	s.body["ccy"] = ccy
	return s
}

// SetAuto enables automatic loan transfer when reducing margin.
func (s *SetMarginBalanceService) SetAuto(auto bool) *SetMarginBalanceService {
	s.body["auto"] = auto
	return s
}

func (s *SetMarginBalanceService) Do(ctx context.Context) (*MarginBalance, error) {
	req := request.Post(ctx, s.c, "/api/v5/account/position/margin-balance", s.body).WithSign()
	return request.DoOne[MarginBalance](req)
}

// MarginBalance is the position-margin adjustment acknowledgement.
type MarginBalance struct {
	InstrumentID string            `json:"instId"`
	PositionSide PosSide           `json:"posSide"`
	Amount       decimal.Decimal   `json:"amt"`
	Type         MarginBalanceType `json:"type"`
	Leverage     decimal.Decimal   `json:"leverage"`
	Currency     string            `json:"ccy"`
}

// GetMMPConfigService -- GET /api/v5/account/mmp-config (Read)
//
// Returns the Market Maker Protection configuration of an option instrument
// family. Requires market-maker permission.
type GetMMPConfigService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetMMPConfigService(instFamily string) *GetMMPConfigService {
	return &GetMMPConfigService{c: c, params: map[string]string{"instFamily": instFamily}}
}

func (s *GetMMPConfigService) Do(ctx context.Context) ([]MMPConfig, error) {
	req := request.Get(ctx, s.c, "/api/v5/account/mmp-config", s.params).WithSign()
	return request.DoList[MMPConfig](req)
}

// MMPConfig is a Market Maker Protection configuration and current state.
type MMPConfig struct {
	InstrumentFamily string          `json:"instFamily"`
	MMPFrozen        bool            `json:"mmpFrozen"`
	MMPFrozenUntil   time.Time       `json:"mmpFrozenUntil"`
	TimeInterval     decimal.Decimal `json:"timeInterval"`
	FrozenInterval   decimal.Decimal `json:"frozenInterval"`
	QtyLimit         decimal.Decimal `json:"qtyLimit"`
}

// SetMMPConfigService -- POST /api/v5/account/mmp-config (Trade)
//
// Sets the Market Maker Protection configuration of an option instrument family.
type SetMMPConfigService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewSetMMPConfigService(instFamily, timeInterval, frozenInterval, qtyLimit string) *SetMMPConfigService {
	return &SetMMPConfigService{c: c, body: map[string]any{
		"instFamily":     instFamily,
		"timeInterval":   timeInterval,
		"frozenInterval": frozenInterval,
		"qtyLimit":       qtyLimit,
	}}
}

func (s *SetMMPConfigService) Do(ctx context.Context) (*MMPConfigSetting, error) {
	req := request.Post(ctx, s.c, "/api/v5/account/mmp-config", s.body).WithSign()
	return request.DoOne[MMPConfigSetting](req)
}

// MMPConfigSetting is the set-MMP acknowledgement.
type MMPConfigSetting struct {
	InstrumentFamily string          `json:"instFamily"`
	TimeInterval     decimal.Decimal `json:"timeInterval"`
	FrozenInterval   decimal.Decimal `json:"frozenInterval"`
	QtyLimit         decimal.Decimal `json:"qtyLimit"`
}

// ResetMMPService -- POST /api/v5/account/mmp-reset (Trade)
//
// Resets the Market Maker Protection frozen state of an option instrument
// family, allowing quoting to resume.
type ResetMMPService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewResetMMPService(instFamily string) *ResetMMPService {
	return &ResetMMPService{c: c, body: map[string]any{"instFamily": instFamily}}
}

// SetInstType sets the instrument type (optional).
func (s *ResetMMPService) SetInstType(instType InstType) *ResetMMPService {
	s.body["instType"] = string(instType)
	return s
}

func (s *ResetMMPService) Do(ctx context.Context) (*MMPReset, error) {
	req := request.Post(ctx, s.c, "/api/v5/account/mmp-reset", s.body).WithSign()
	return request.DoOne[MMPReset](req)
}

// MMPReset is the mmp-reset acknowledgement.
type MMPReset struct {
	Result bool `json:"result"`
}

// MovePositionLeg is one position to move in a move-positions request. From
// describes the source-account leg and To the destination-account leg.
type MovePositionLeg struct {
	From MovePositionLegFromReq `json:"from"`
	To   MovePositionLegToReq   `json:"to"`
}

// MovePositionLegFromReq is the source leg of a move-positions request.
type MovePositionLegFromReq struct {
	PositionID string `json:"posId"`
	Size       string `json:"sz"`
	Side       Side   `json:"side"`
}

// MovePositionLegToReq is the destination leg of a move-positions request.
type MovePositionLegToReq struct {
	TradeMode    TdMode  `json:"tdMode,omitempty"`
	PositionSide PosSide `json:"posSide,omitempty"`
	Currency     string  `json:"ccy,omitempty"`
}

// MovePositionsService -- POST /api/v5/account/move-positions (Trade)
//
// Moves positions between the master account and a managed sub-account (block
// transfer). Each leg pairs a source-account position with its destination.
type MovePositionsService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewMovePositionsService(fromAcct, toAcct, clientId string, legs []MovePositionLeg) *MovePositionsService {
	return &MovePositionsService{c: c, body: map[string]any{
		"fromAcct": fromAcct,
		"toAcct":   toAcct,
		"clientId": clientId,
		"legs":     legs,
	}}
}

func (s *MovePositionsService) Do(ctx context.Context) (*MovePositions, error) {
	req := request.Post(ctx, s.c, "/api/v5/account/move-positions", s.body).WithSign()
	return request.DoOne[MovePositions](req)
}

// MovePositions is the move-positions acknowledgement.
type MovePositions struct {
	BlockTradeID string             `json:"blockTdId"`
	ClientID     string             `json:"clientId"`
	State        string             `json:"state"`
	FromAccount  string             `json:"fromAcct"`
	ToAccount    string             `json:"toAcct"`
	Legs         []MovePositionsLeg `json:"legs"`
	Timestamp    time.Time          `json:"ts"`
}

// MovePositionsLeg is one moved-position leg in a move-positions response.
type MovePositionsLeg struct {
	From MovePositionsLegFrom `json:"from"`
	To   MovePositionsLegTo   `json:"to"`
}

// MovePositionsLegFrom is the source-account side of a moved-position leg.
type MovePositionsLegFrom struct {
	InstrumentID string          `json:"instId"`
	PositionID   string          `json:"posId"`
	Price        decimal.Decimal `json:"px"`
	Side         Side            `json:"side"`
	Size         decimal.Decimal `json:"sz"`
	SCode        string          `json:"sCode"`
	SMsg         string          `json:"sMsg"`
}

// MovePositionsLegTo is the destination-account side of a moved-position leg.
type MovePositionsLegTo struct {
	InstrumentID string          `json:"instId"`
	Side         Side            `json:"side"`
	PositionSide PosSide         `json:"posSide"`
	TradeMode    TdMode          `json:"tdMode"`
	Price        decimal.Decimal `json:"px"`
	Currency     string          `json:"ccy"`
	SCode        string          `json:"sCode"`
	SMsg         string          `json:"sMsg"`
}

// GetMovePositionsHistoryService -- GET /api/v5/account/move-positions-history (Read)
//
// Returns the history of move-positions (block transfer) requests.
type GetMovePositionsHistoryService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetMovePositionsHistoryService() *GetMovePositionsHistoryService {
	return &GetMovePositionsHistoryService{c: c, params: map[string]string{}}
}

// SetBlockTdId filters by block trade id.
func (s *GetMovePositionsHistoryService) SetBlockTdId(blockTdId string) *GetMovePositionsHistoryService {
	s.params["blockTdId"] = blockTdId
	return s
}

// SetClientId filters by client-supplied id.
func (s *GetMovePositionsHistoryService) SetClientId(clientId string) *GetMovePositionsHistoryService {
	s.params["clientId"] = clientId
	return s
}

// SetBeginTs filters to requests at or after the given time.
func (s *GetMovePositionsHistoryService) SetBeginTs(t time.Time) *GetMovePositionsHistoryService {
	s.params["beginTs"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTs filters to requests at or before the given time.
func (s *GetMovePositionsHistoryService) SetEndTs(t time.Time) *GetMovePositionsHistoryService {
	s.params["endTs"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetState filters by transfer state ("filled" / "pending").
func (s *GetMovePositionsHistoryService) SetState(state string) *GetMovePositionsHistoryService {
	s.params["state"] = state
	return s
}

// SetLimit caps the number of records returned.
func (s *GetMovePositionsHistoryService) SetLimit(limit int) *GetMovePositionsHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetMovePositionsHistoryService) Do(ctx context.Context) ([]MovePositionsHistory, error) {
	req := request.Get(ctx, s.c, "/api/v5/account/move-positions-history", s.params).WithSign()
	return request.DoList[MovePositionsHistory](req)
}

// MovePositionsHistory is one move-positions (block transfer) record.
type MovePositionsHistory struct {
	ClientID     string                    `json:"clientId"`
	BlockTradeID string                    `json:"blockTdId"`
	State        string                    `json:"state"`
	Timestamp    time.Time                 `json:"ts"`
	FromAccount  string                    `json:"fromAcct"`
	ToAccount    string                    `json:"toAcct"`
	Legs         []MovePositionsHistoryLeg `json:"legs"`
}

// MovePositionsHistoryLeg is one leg of a move-positions history record.
type MovePositionsHistoryLeg struct {
	From MovePositionsHistoryLegFrom `json:"from"`
	To   MovePositionsHistoryLegTo   `json:"to"`
}

// MovePositionsHistoryLegFrom is the source-account side of a history leg.
type MovePositionsHistoryLegFrom struct {
	InstrumentID string          `json:"instId"`
	PositionID   string          `json:"posId"`
	Price        decimal.Decimal `json:"px"`
	Side         Side            `json:"side"`
	Size         decimal.Decimal `json:"sz"`
}

// MovePositionsHistoryLegTo is the destination-account side of a history leg.
type MovePositionsHistoryLegTo struct {
	InstrumentID string          `json:"instId"`
	Price        decimal.Decimal `json:"px"`
	Side         Side            `json:"side"`
	Size         decimal.Decimal `json:"sz"`
	TradeMode    TdMode          `json:"tdMode"`
	PositionSide PosSide         `json:"posSide"`
	Currency     string          `json:"ccy"`
}

// PositionBuilderSimPos is a simulated position fed to the position builder.
type PositionBuilderSimPos struct {
	InstrumentID string `json:"instId"`
	Position     string `json:"pos"`
	AveragePrice string `json:"avgPx,omitempty"`
	Leverage     string `json:"lever,omitempty"`
}

// PositionBuilderSimAsset is a simulated asset fed to the position builder.
type PositionBuilderSimAsset struct {
	Currency string `json:"ccy"`
	Amount   string `json:"amt"`
}

// PositionBuilderService -- POST /api/v5/account/position-builder (Read)
//
// Simulates portfolio-margin risk for a hypothetical set of positions and
// assets, returning the resulting margin requirements and greeks. This is a
// read-only simulation that does not change account state.
type PositionBuilderService struct {
	c    *Client
	body map[string]any
}

func (c *Client) NewPositionBuilderService() *PositionBuilderService {
	return &PositionBuilderService{c: c, body: map[string]any{}}
}

// SetAcctLv sets the simulated account mode.
func (s *PositionBuilderService) SetAcctLv(acctLv string) *PositionBuilderService {
	s.body["acctLv"] = acctLv
	return s
}

// SetInclRealPosAndEq includes the account's real positions and equity in the
// simulation.
func (s *PositionBuilderService) SetInclRealPosAndEq(incl bool) *PositionBuilderService {
	s.body["inclRealPosAndEq"] = incl
	return s
}

// SetLever sets the simulated leverage.
func (s *PositionBuilderService) SetLever(lever string) *PositionBuilderService {
	s.body["lever"] = lever
	return s
}

// SetGreeksType sets the greeks display type.
func (s *PositionBuilderService) SetGreeksType(greeksType GreeksType) *PositionBuilderService {
	s.body["greeksType"] = string(greeksType)
	return s
}

// SetIdxVol sets the index volatility used for the stress-test scenarios.
func (s *PositionBuilderService) SetIdxVol(idxVol string) *PositionBuilderService {
	s.body["idxVol"] = idxVol
	return s
}

// SetSimPos sets the simulated positions.
func (s *PositionBuilderService) SetSimPos(simPos []PositionBuilderSimPos) *PositionBuilderService {
	s.body["simPos"] = simPos
	return s
}

// SetSimAsset sets the simulated assets.
func (s *PositionBuilderService) SetSimAsset(simAsset []PositionBuilderSimAsset) *PositionBuilderService {
	s.body["simAsset"] = simAsset
	return s
}

func (s *PositionBuilderService) Do(ctx context.Context) (*PositionBuilder, error) {
	req := request.Post(ctx, s.c, "/api/v5/account/position-builder", s.body).WithSign()
	return request.DoOne[PositionBuilder](req)
}

// PositionBuilder is the result of a portfolio-margin simulation.
type PositionBuilder struct {
	Equity          decimal.Decimal           `json:"eq"`
	TotalMMR        decimal.Decimal           `json:"totalMmr"`
	TotalIMR        decimal.Decimal           `json:"totalImr"`
	BorrowMMR       decimal.Decimal           `json:"borrowMmr"`
	DerivativesMMR  decimal.Decimal           `json:"derivMmr"`
	MarginRatio     decimal.Decimal           `json:"marginRatio"`
	UPL             decimal.Decimal           `json:"upl"`
	AccountLeverage decimal.Decimal           `json:"acctLever"`
	Timestamp       time.Time                 `json:"ts"`
	Assets          []PositionBuilderAsset    `json:"assets"`
	RiskUnitData    []PositionBuilderRiskUnit `json:"riskUnitData"`
	Positions       []PositionBuilderPosition `json:"positions"`
}

// PositionBuilderAsset is one asset row of a position-builder simulation.
type PositionBuilderAsset struct {
	Currency        string          `json:"ccy"`
	AvailableEquity decimal.Decimal `json:"availEq"`
	SpotInUse       decimal.Decimal `json:"spotInUse"`
	BorrowMMR       decimal.Decimal `json:"borrowMmr"`
	BorrowIMR       decimal.Decimal `json:"borrowImr"`
}

// PositionBuilderRiskUnit is one risk-unit row of a position-builder simulation.
type PositionBuilderRiskUnit struct {
	RiskUnit       string                        `json:"riskUnit"`
	MMRBefore      decimal.Decimal               `json:"mmrBf"`
	MMR            decimal.Decimal               `json:"mmr"`
	IMRBefore      decimal.Decimal               `json:"imrBf"`
	IMR            decimal.Decimal               `json:"imr"`
	UPL            decimal.Decimal               `json:"upl"`
	Mr1            decimal.Decimal               `json:"mr1"`
	Mr2            decimal.Decimal               `json:"mr2"`
	Mr3            decimal.Decimal               `json:"mr3"`
	Mr4            decimal.Decimal               `json:"mr4"`
	Mr5            decimal.Decimal               `json:"mr5"`
	Mr6            decimal.Decimal               `json:"mr6"`
	Mr7            decimal.Decimal               `json:"mr7"`
	Mr8            decimal.Decimal               `json:"mr8"`
	Mr9            decimal.Decimal               `json:"mr9"`
	Mr1Scenarios   PositionBuilderMr1Scenarios   `json:"mr1Scenarios"`
	Mr1FinalResult PositionBuilderMrFinalResult  `json:"mr1FinalResult"`
	Mr6FinalResult PositionBuilderMr6FinalResult `json:"mr6FinalResult"`
	Delta          decimal.Decimal               `json:"delta"`
	Gamma          decimal.Decimal               `json:"gamma"`
	Theta          decimal.Decimal               `json:"theta"`
	Vega           decimal.Decimal               `json:"vega"`
	Portfolios     []PositionBuilderPortfolio    `json:"portfolios"`
}

// PositionBuilderMr1Scenarios holds the MR1 stress-test P&L grids keyed by price
// volatility ratio (a "change" -> "value" map per volatility shock direction).
type PositionBuilderMr1Scenarios struct {
	VolatilityShockDown map[string]string `json:"volShockDown"`
	VolatilitySame      map[string]string `json:"volSame"`
	VolatilityShockUp   map[string]string `json:"volShockUp"`
}

// PositionBuilderMrFinalResult is the worst-case MR1 scenario.
type PositionBuilderMrFinalResult struct {
	Pnl             decimal.Decimal `json:"pnl"`
	SpotShock       decimal.Decimal `json:"spotShock"`
	VolatilityShock string          `json:"volShock"`
}

// PositionBuilderMr6FinalResult is the worst-case MR6 scenario.
type PositionBuilderMr6FinalResult struct {
	Pnl       decimal.Decimal `json:"pnl"`
	SpotShock decimal.Decimal `json:"spotShock"`
}

// PositionBuilderPortfolio is one instrument within a risk unit (portfolio
// margin only).
type PositionBuilderPortfolio struct {
	InstrumentID    string          `json:"instId"`
	InstrumentType  InstType        `json:"instType"`
	Amount          decimal.Decimal `json:"amt"`
	PositionSide    PosSide         `json:"posSide"`
	AveragePrice    decimal.Decimal `json:"avgPx"`
	MarkPriceBefore decimal.Decimal `json:"markPxBf"`
	MarkPrice       decimal.Decimal `json:"markPx"`
	FloatPnl        decimal.Decimal `json:"floatPnl"`
	NotionalUSD     decimal.Decimal `json:"notionalUsd"`
	Delta           decimal.Decimal `json:"delta"`
	Gamma           decimal.Decimal `json:"gamma"`
	Theta           decimal.Decimal `json:"theta"`
	Vega            decimal.Decimal `json:"vega"`
	IsRealPosition  bool            `json:"isRealPos"`
}

// PositionBuilderPosition is one position row (multi-currency margin only).
type PositionBuilderPosition struct {
	InstrumentID    string          `json:"instId"`
	InstrumentType  InstType        `json:"instType"`
	Amount          decimal.Decimal `json:"amt"`
	PositionSide    PosSide         `json:"posSide"`
	AveragePrice    decimal.Decimal `json:"avgPx"`
	MarkPriceBefore decimal.Decimal `json:"markPxBf"`
	MarkPrice       decimal.Decimal `json:"markPx"`
	FloatPnl        decimal.Decimal `json:"floatPnl"`
	IMRBefore       decimal.Decimal `json:"imrBf"`
	IMR             decimal.Decimal `json:"imr"`
	MarginRatio     decimal.Decimal `json:"mgnRatio"`
	Leverage        decimal.Decimal `json:"lever"`
	NotionalUSD     decimal.Decimal `json:"notionalUsd"`
	IsRealPosition  bool            `json:"isRealPos"`
}
