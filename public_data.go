package okx

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-okx/request"
	"github.com/shopspring/decimal"
)

// publicOpenType is an instrument's pre-market opening mechanism (auction vs
// fixed-price), reported by the instruments endpoint.
type publicOpenType string

const (
	publicOpenTypeFixPrice    publicOpenType = "fix_price"
	publicOpenTypePreQuote    publicOpenType = "pre_quote"
	publicOpenTypeCallAuction publicOpenType = "call_auction"
)

// publicSettState is the funding-rate settlement state reported by the
// funding-rate endpoint.
type publicSettState string

const (
	publicSettStateProcessing publicSettState = "processing"
	publicSettStateSettled    publicSettState = "settled"
)

// GetInstrumentsService -- GET /api/v5/public/instruments (public)
//
// Returns the tradable instruments for a product line, including their lot/tick
// sizes, contract specs and listing state.
type GetInstrumentsService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetInstrumentsService(instType InstType) *GetInstrumentsService {
	return &GetInstrumentsService{c: c, params: map[string]string{"instType": string(instType)}}
}

// SetUly filters by underlying (applicable to FUTURES/SWAP/OPTION).
func (s *GetInstrumentsService) SetUly(uly string) *GetInstrumentsService {
	s.params["uly"] = uly
	return s
}

// SetInstFamily filters by instrument family (applicable to FUTURES/SWAP/OPTION).
func (s *GetInstrumentsService) SetInstFamily(instFamily string) *GetInstrumentsService {
	s.params["instFamily"] = instFamily
	return s
}

// SetInstId filters by a single instrument id.
func (s *GetInstrumentsService) SetInstId(instId string) *GetInstrumentsService {
	s.params["instId"] = instId
	return s
}

func (s *GetInstrumentsService) Do(ctx context.Context) ([]Instrument, error) {
	req := request.Get(ctx, s.c, "/api/v5/public/instruments", s.params)
	return request.DoList[Instrument](req)
}

// Instrument is a single tradable instrument.
type Instrument struct {
	InstrumentType                   InstType           `json:"instType"`
	InstrumentID                     string             `json:"instId"`
	InstrumentIDCode                 int64              `json:"instIdCode"`
	Underlying                       string             `json:"uly"`
	InstrumentFamily                 string             `json:"instFamily"`
	Category                         string             `json:"category"`
	InstrumentCategory               string             `json:"instCategory"`
	BaseCurrency                     string             `json:"baseCcy"`
	QuoteCurrency                    string             `json:"quoteCcy"`
	SettleCurrency                   string             `json:"settleCcy"`
	ContractValue                    decimal.Decimal    `json:"ctVal"`
	ContractMultiplier               decimal.Decimal    `json:"ctMult"`
	ContractValueCurrency            string             `json:"ctValCcy"`
	OptionType                       OptType            `json:"optType"`
	Strike                           decimal.Decimal    `json:"stk"`
	ListTime                         time.Time          `json:"listTime"`
	AuctionEndTime                   time.Time          `json:"auctionEndTime"`
	ContinuousTradeSwitchTime        time.Time          `json:"contTdSwTime"`
	OpenType                         publicOpenType     `json:"openType"`
	ExpiryTime                       time.Time          `json:"expTime"`
	Leverage                         decimal.Decimal    `json:"lever"`
	TickSize                         decimal.Decimal    `json:"tickSz"`
	LotSize                          decimal.Decimal    `json:"lotSz"`
	MinSize                          decimal.Decimal    `json:"minSz"`
	ContractType                     CtType             `json:"ctType"`
	Alias                            string             `json:"alias"`
	State                            InstState          `json:"state"`
	RuleType                         string             `json:"ruleType"`
	MaxLimitSize                     decimal.Decimal    `json:"maxLmtSz"`
	MaxMarketSize                    decimal.Decimal    `json:"maxMktSz"`
	MaxLimitAmount                   decimal.Decimal    `json:"maxLmtAmt"`
	MaxMarketAmount                  decimal.Decimal    `json:"maxMktAmt"`
	MaxTWAPSize                      decimal.Decimal    `json:"maxTwapSz"`
	MaxIcebergSize                   decimal.Decimal    `json:"maxIcebergSz"`
	MaxTriggerSize                   decimal.Decimal    `json:"maxTriggerSz"`
	MaxStopSize                      decimal.Decimal    `json:"maxStopSz"`
	MaxPlatformOpenInterestLimit     decimal.Decimal    `json:"maxPlatOILmt"`
	MaxPlatformOpenInterestCoinLimit decimal.Decimal    `json:"maxPlatOICoinLmt"`
	FutureSettlement                 bool               `json:"futureSettlement"`
	TradeQuoteCurrencyList           []string           `json:"tradeQuoteCcyList"`
	UpcomingChange                   []InstrumentUpcChg `json:"upcChg"`
	Freq                             string             `json:"freq"`
	GroupID                          string             `json:"groupId"`
	SeriesID                         string             `json:"seriesId"`
	Method                           string             `json:"method"`
	LongPositionRemainingQuota       decimal.Decimal    `json:"longPosRemainingQuota"`
	ShortPositionRemainingQuota      decimal.Decimal    `json:"shortPosRemainingQuota"`
	PositionLimitAmount              decimal.Decimal    `json:"posLmtAmt"`
	PositionLimitPercent             decimal.Decimal    `json:"posLmtPct"`
	PreMarketSwitchTime              time.Time          `json:"preMktSwTime"`
	// InitialPriceLimitPercent is the initial price-limit band applied during the
	// first 10 minutes after contract listing. Empty for OPTION and EVENTS.
	InitialPriceLimitPercent decimal.Decimal `json:"initPxLmtPct"`
	// FloatingPriceLimitPercent is the floating price-limit band during normal
	// trading. Empty for OPTION and EVENTS.
	FloatingPriceLimitPercent decimal.Decimal `json:"floatPxLmtPct"`
	// MaxPriceLimitPercent is the maximum price-limit cap (hard ceiling). Empty
	// for OPTION and EVENTS.
	MaxPriceLimitPercent decimal.Decimal `json:"maxPxLmtPct"`
	// Elp (effective leverage profile) is returned only by the private
	// GET /api/v5/account/instruments variant; it is empty on the public endpoint.
	Elp string `json:"elp"`
}

// InstrumentUpcChg is a scheduled upcoming change to one of an instrument's
// trading rules (e.g. a tickSz adjustment).
type InstrumentUpcChg struct {
	EffectiveTime time.Time `json:"effTime"`
	NewValue      string    `json:"newValue"`
	Param         string    `json:"param"`
}

// GetEstimatedPriceService -- GET /api/v5/public/estimated-price (public)
//
// Returns the estimated delivery/exercise price of a FUTURES or OPTION
// instrument (only available within the hour before delivery/exercise).
type GetEstimatedPriceService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetEstimatedPriceService(instId string) *GetEstimatedPriceService {
	return &GetEstimatedPriceService{c: c, params: map[string]string{"instId": instId}}
}

func (s *GetEstimatedPriceService) Do(ctx context.Context) (*EstimatedPrice, error) {
	req := request.Get(ctx, s.c, "/api/v5/public/estimated-price", s.params)
	return request.DoOne[EstimatedPrice](req)
}

// EstimatedPrice is the estimated delivery/exercise price of an instrument.
type EstimatedPrice struct {
	InstrumentType InstType        `json:"instType"`
	InstrumentID   string          `json:"instId"`
	SettlePrice    decimal.Decimal `json:"settlePx"`
	Timestamp      time.Time       `json:"ts"`
}

// GetDeliveryExerciseHistoryService -- GET /api/v5/public/delivery-exercise-history (public)
//
// Returns the delivery (FUTURES) / exercise (OPTION) records of an underlying or
// instrument family over the last three months.
type GetDeliveryExerciseHistoryService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetDeliveryExerciseHistoryService(instType InstType) *GetDeliveryExerciseHistoryService {
	return &GetDeliveryExerciseHistoryService{c: c, params: map[string]string{"instType": string(instType)}}
}

// SetUly filters by underlying.
func (s *GetDeliveryExerciseHistoryService) SetUly(uly string) *GetDeliveryExerciseHistoryService {
	s.params["uly"] = uly
	return s
}

// SetInstFamily filters by instrument family.
func (s *GetDeliveryExerciseHistoryService) SetInstFamily(instFamily string) *GetDeliveryExerciseHistoryService {
	s.params["instFamily"] = instFamily
	return s
}

// SetAfter paginates to records earlier than the given time (older).
func (s *GetDeliveryExerciseHistoryService) SetAfter(t time.Time) *GetDeliveryExerciseHistoryService {
	s.params["after"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetBefore paginates to records later than the given time (newer).
func (s *GetDeliveryExerciseHistoryService) SetBefore(t time.Time) *GetDeliveryExerciseHistoryService {
	s.params["before"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetDeliveryExerciseHistoryService) SetLimit(limit int) *GetDeliveryExerciseHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetDeliveryExerciseHistoryService) Do(ctx context.Context) ([]DeliveryExerciseHistory, error) {
	req := request.Get(ctx, s.c, "/api/v5/public/delivery-exercise-history", s.params)
	return request.DoList[DeliveryExerciseHistory](req)
}

// DeliveryExerciseHistory is one delivery/exercise event and its per-instrument
// details.
type DeliveryExerciseHistory struct {
	Timestamp time.Time                       `json:"ts"`
	Details   []DeliveryExerciseHistoryDetail `json:"details"`
}

// DeliveryExerciseHistoryDetail is a single instrument's delivery/exercise
// record.
type DeliveryExerciseHistoryDetail struct {
	InstrumentID string          `json:"insId"`
	Price        decimal.Decimal `json:"px"`
	Type         string          `json:"type"`
}

// GetSettlementHistoryService -- GET /api/v5/public/settlement-history (public)
//
// Returns the futures settlement history of an instrument family.
type GetSettlementHistoryService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetSettlementHistoryService(instType InstType, instFamily string) *GetSettlementHistoryService {
	return &GetSettlementHistoryService{c: c, params: map[string]string{
		"instType":   string(instType),
		"instFamily": instFamily,
	}}
}

// SetAfter paginates to records earlier than the given time (older).
func (s *GetSettlementHistoryService) SetAfter(t time.Time) *GetSettlementHistoryService {
	s.params["after"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetBefore paginates to records later than the given time (newer).
func (s *GetSettlementHistoryService) SetBefore(t time.Time) *GetSettlementHistoryService {
	s.params["before"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetSettlementHistoryService) SetLimit(limit int) *GetSettlementHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetSettlementHistoryService) Do(ctx context.Context) ([]SettlementHistory, error) {
	req := request.Get(ctx, s.c, "/api/v5/public/settlement-history", s.params)
	return request.DoList[SettlementHistory](req)
}

// SettlementHistory is one settlement event and its per-instrument prices.
type SettlementHistory struct {
	Timestamp time.Time                 `json:"ts"`
	Details   []SettlementHistoryDetail `json:"details"`
}

// SettlementHistoryDetail is a single instrument's settlement price.
type SettlementHistoryDetail struct {
	InstrumentID string          `json:"instId"`
	SettlePrice  decimal.Decimal `json:"settlePx"`
}

// GetFundingRateService -- GET /api/v5/public/funding-rate (public)
//
// Returns the current funding rate of a perpetual swap.
type GetFundingRateService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetFundingRateService(instId string) *GetFundingRateService {
	return &GetFundingRateService{c: c, params: map[string]string{"instId": instId}}
}

func (s *GetFundingRateService) Do(ctx context.Context) (*FundingRate, error) {
	req := request.Get(ctx, s.c, "/api/v5/public/funding-rate", s.params)
	return request.DoOne[FundingRate](req)
}

// FundingRate is a perpetual swap's current funding rate.
type FundingRate struct {
	InstrumentType        InstType        `json:"instType"`
	InstrumentID          string          `json:"instId"`
	Method                string          `json:"method"`
	FormulaType           string          `json:"formulaType"`
	FundingRate           decimal.Decimal `json:"fundingRate"`
	NextFundingRate       decimal.Decimal `json:"nextFundingRate"`
	FundingTime           time.Time       `json:"fundingTime"`
	NextFundingTime       time.Time       `json:"nextFundingTime"`
	MinFundingRate        decimal.Decimal `json:"minFundingRate"`
	MaxFundingRate        decimal.Decimal `json:"maxFundingRate"`
	InterestRate          decimal.Decimal `json:"interestRate"`
	ImpactValue           decimal.Decimal `json:"impactValue"`
	SettlementState       publicSettState `json:"settState"`
	SettlementFundingRate decimal.Decimal `json:"settFundingRate"`
	Premium               decimal.Decimal `json:"premium"`
	PreviousFundingTime   time.Time       `json:"prevFundingTime"`
	Timestamp             time.Time       `json:"ts"`
}

// GetFundingRateHistoryService -- GET /api/v5/public/funding-rate-history (public)
//
// Returns the historical (realized) funding rates of a perpetual swap.
type GetFundingRateHistoryService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetFundingRateHistoryService(instId string) *GetFundingRateHistoryService {
	return &GetFundingRateHistoryService{c: c, params: map[string]string{"instId": instId}}
}

// SetAfter paginates to records earlier than the given time (older).
func (s *GetFundingRateHistoryService) SetAfter(t time.Time) *GetFundingRateHistoryService {
	s.params["after"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetBefore paginates to records later than the given time (newer).
func (s *GetFundingRateHistoryService) SetBefore(t time.Time) *GetFundingRateHistoryService {
	s.params["before"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetFundingRateHistoryService) SetLimit(limit int) *GetFundingRateHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetFundingRateHistoryService) Do(ctx context.Context) ([]FundingRateHistory, error) {
	req := request.Get(ctx, s.c, "/api/v5/public/funding-rate-history", s.params)
	return request.DoList[FundingRateHistory](req)
}

// FundingRateHistory is one realized funding-rate record.
type FundingRateHistory struct {
	InstrumentType InstType        `json:"instType"`
	InstrumentID   string          `json:"instId"`
	Method         string          `json:"method"`
	FormulaType    string          `json:"formulaType"`
	FundingRate    decimal.Decimal `json:"fundingRate"`
	RealizedRate   decimal.Decimal `json:"realizedRate"`
	FundingTime    time.Time       `json:"fundingTime"`
}

// GetOpenInterestService -- GET /api/v5/public/open-interest (public)
//
// Returns the open interest of derivatives (SWAP/FUTURES/OPTION).
type GetOpenInterestService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetOpenInterestService(instType InstType) *GetOpenInterestService {
	return &GetOpenInterestService{c: c, params: map[string]string{"instType": string(instType)}}
}

// SetUly filters by underlying.
func (s *GetOpenInterestService) SetUly(uly string) *GetOpenInterestService {
	s.params["uly"] = uly
	return s
}

// SetInstFamily filters by instrument family.
func (s *GetOpenInterestService) SetInstFamily(instFamily string) *GetOpenInterestService {
	s.params["instFamily"] = instFamily
	return s
}

// SetInstId filters by a single instrument id.
func (s *GetOpenInterestService) SetInstId(instId string) *GetOpenInterestService {
	s.params["instId"] = instId
	return s
}

func (s *GetOpenInterestService) Do(ctx context.Context) ([]OpenInterest, error) {
	req := request.Get(ctx, s.c, "/api/v5/public/open-interest", s.params)
	return request.DoList[OpenInterest](req)
}

// OpenInterest is an instrument's open interest.
type OpenInterest struct {
	InstrumentType       InstType        `json:"instType"`
	InstrumentID         string          `json:"instId"`
	OpenInterest         decimal.Decimal `json:"oi"`
	OpenInterestCurrency decimal.Decimal `json:"oiCcy"`
	OpenInterestUSD      decimal.Decimal `json:"oiUsd"`
	Timestamp            time.Time       `json:"ts"`
}

// GetPriceLimitService -- GET /api/v5/public/price-limit (public)
//
// Returns the highest buy and lowest sell limit prices of an instrument.
type GetPriceLimitService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetPriceLimitService(instId string) *GetPriceLimitService {
	return &GetPriceLimitService{c: c, params: map[string]string{"instId": instId}}
}

func (s *GetPriceLimitService) Do(ctx context.Context) (*PriceLimit, error) {
	req := request.Get(ctx, s.c, "/api/v5/public/price-limit", s.params)
	return request.DoOne[PriceLimit](req)
}

// PriceLimit is an instrument's buy/sell price limits.
type PriceLimit struct {
	InstrumentType InstType        `json:"instType"`
	InstrumentID   string          `json:"instId"`
	BuyLimit       decimal.Decimal `json:"buyLmt"`
	SellLimit      decimal.Decimal `json:"sellLmt"`
	Enabled        bool            `json:"enabled"`
	Timestamp      time.Time       `json:"ts"`
}

// GetOptSummaryService -- GET /api/v5/public/opt-summary (public)
//
// Returns the option market-data summary (greeks and implied volatilities) for
// an underlying or instrument family.
type GetOptSummaryService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetOptSummaryService() *GetOptSummaryService {
	return &GetOptSummaryService{c: c, params: map[string]string{}}
}

// SetUly filters by underlying.
func (s *GetOptSummaryService) SetUly(uly string) *GetOptSummaryService {
	s.params["uly"] = uly
	return s
}

// SetInstFamily filters by instrument family.
func (s *GetOptSummaryService) SetInstFamily(instFamily string) *GetOptSummaryService {
	s.params["instFamily"] = instFamily
	return s
}

// SetExpTime filters by an option expiry (YYYYMMdd).
func (s *GetOptSummaryService) SetExpTime(expTime string) *GetOptSummaryService {
	s.params["expTime"] = expTime
	return s
}

func (s *GetOptSummaryService) Do(ctx context.Context) ([]OptSummary, error) {
	req := request.Get(ctx, s.c, "/api/v5/public/opt-summary", s.params)
	return request.DoList[OptSummary](req)
}

// OptSummary is the greeks/IV summary of a single option instrument.
type OptSummary struct {
	InstrumentType     InstType        `json:"instType"`
	InstrumentID       string          `json:"instId"`
	Underlying         string          `json:"uly"`
	Delta              decimal.Decimal `json:"delta"`
	Gamma              decimal.Decimal `json:"gamma"`
	Theta              decimal.Decimal `json:"theta"`
	Vega               decimal.Decimal `json:"vega"`
	DeltaBS            decimal.Decimal `json:"deltaBS"`
	GammaBS            decimal.Decimal `json:"gammaBS"`
	ThetaBS            decimal.Decimal `json:"thetaBS"`
	VegaBS             decimal.Decimal `json:"vegaBS"`
	RealizedVolatility decimal.Decimal `json:"realVol"`
	BidVolatility      decimal.Decimal `json:"bidVol"`
	AskVolatility      decimal.Decimal `json:"askVol"`
	MarkVolatility     decimal.Decimal `json:"markVol"`
	Leverage           decimal.Decimal `json:"lever"`
	VolatilityLevel    decimal.Decimal `json:"volLv"`
	ForwardPrice       decimal.Decimal `json:"fwdPx"`
	Distance           decimal.Decimal `json:"distance"`
	BuyAPR             decimal.Decimal `json:"buyApr"`
	SellAPR            decimal.Decimal `json:"sellApr"`
	Timestamp          time.Time       `json:"ts"`
}

// GetDiscountRateInterestFreeQuotaService -- GET /api/v5/public/discount-rate-interest-free-quota (public)
//
// Returns the discount-rate tiers and interest-free quota per currency.
type GetDiscountRateInterestFreeQuotaService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetDiscountRateInterestFreeQuotaService() *GetDiscountRateInterestFreeQuotaService {
	return &GetDiscountRateInterestFreeQuotaService{c: c, params: map[string]string{}}
}

// SetCcy filters by currency.
func (s *GetDiscountRateInterestFreeQuotaService) SetCcy(ccy string) *GetDiscountRateInterestFreeQuotaService {
	s.params["ccy"] = ccy
	return s
}

// SetDiscountLv filters by discount level.
func (s *GetDiscountRateInterestFreeQuotaService) SetDiscountLv(discountLv string) *GetDiscountRateInterestFreeQuotaService {
	s.params["discountLv"] = discountLv
	return s
}

func (s *GetDiscountRateInterestFreeQuotaService) Do(ctx context.Context) ([]DiscountRateInterestFreeQuota, error) {
	req := request.Get(ctx, s.c, "/api/v5/public/discount-rate-interest-free-quota", s.params)
	return request.DoList[DiscountRateInterestFreeQuota](req)
}

// DiscountRateInterestFreeQuota is a currency's interest-free quota and its
// tiered discount-rate schedule.
type DiscountRateInterestFreeQuota struct {
	Currency           string             `json:"ccy"`
	Amount             decimal.Decimal    `json:"amt"`
	DiscountLevel      string             `json:"discountLv"`
	MinDiscountRate    decimal.Decimal    `json:"minDiscountRate"`
	CollateralRestrict bool               `json:"collateralRestrict"`
	ColRes             string             `json:"colRes"`
	Details            []DiscountRateTier `json:"details"`
}

// DiscountRateTier is one tier of a currency's discount-rate schedule.
type DiscountRateTier struct {
	DiscountRate           decimal.Decimal `json:"discountRate"`
	LiquidationPenaltyRate decimal.Decimal `json:"liqPenaltyRate"`
	DiscountCurrencyEquity decimal.Decimal `json:"disCcyEq"`
	MinAmount              decimal.Decimal `json:"minAmt"`
	MaxAmount              decimal.Decimal `json:"maxAmt"`
	Tier                   string          `json:"tier"`
}

// GetMarkPriceService -- GET /api/v5/public/mark-price (public)
//
// Returns the mark prices of instruments in a product line.
type GetMarkPriceService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetMarkPriceService(instType InstType) *GetMarkPriceService {
	return &GetMarkPriceService{c: c, params: map[string]string{"instType": string(instType)}}
}

// SetUly filters by underlying.
func (s *GetMarkPriceService) SetUly(uly string) *GetMarkPriceService {
	s.params["uly"] = uly
	return s
}

// SetInstFamily filters by instrument family.
func (s *GetMarkPriceService) SetInstFamily(instFamily string) *GetMarkPriceService {
	s.params["instFamily"] = instFamily
	return s
}

// SetInstId filters by a single instrument id.
func (s *GetMarkPriceService) SetInstId(instId string) *GetMarkPriceService {
	s.params["instId"] = instId
	return s
}

func (s *GetMarkPriceService) Do(ctx context.Context) ([]MarkPrice, error) {
	req := request.Get(ctx, s.c, "/api/v5/public/mark-price", s.params)
	return request.DoList[MarkPrice](req)
}

// MarkPrice is an instrument's mark price.
type MarkPrice struct {
	InstrumentType InstType        `json:"instType"`
	InstrumentID   string          `json:"instId"`
	MarkPrice      decimal.Decimal `json:"markPx"`
	Timestamp      time.Time       `json:"ts"`
}

// GetPositionTiersService -- GET /api/v5/public/position-tiers (public)
//
// Returns the position-tier (leverage / margin requirement) schedule for an
// instrument family or underlying.
type GetPositionTiersService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetPositionTiersService(instType InstType, tdMode TdMode) *GetPositionTiersService {
	return &GetPositionTiersService{c: c, params: map[string]string{
		"instType": string(instType),
		"tdMode":   string(tdMode),
	}}
}

// SetUly filters by underlying (comma-separated for multiple).
func (s *GetPositionTiersService) SetUly(uly string) *GetPositionTiersService {
	s.params["uly"] = uly
	return s
}

// SetInstFamily filters by instrument family (comma-separated for multiple).
func (s *GetPositionTiersService) SetInstFamily(instFamily string) *GetPositionTiersService {
	s.params["instFamily"] = instFamily
	return s
}

// SetInstId filters by a single instrument id (SPOT/MARGIN).
func (s *GetPositionTiersService) SetInstId(instId string) *GetPositionTiersService {
	s.params["instId"] = instId
	return s
}

// SetCcy filters by margin currency (MARGIN cross only).
func (s *GetPositionTiersService) SetCcy(ccy string) *GetPositionTiersService {
	s.params["ccy"] = ccy
	return s
}

// SetTier filters by a single tier.
func (s *GetPositionTiersService) SetTier(tier string) *GetPositionTiersService {
	s.params["tier"] = tier
	return s
}

func (s *GetPositionTiersService) Do(ctx context.Context) ([]PositionTier, error) {
	req := request.Get(ctx, s.c, "/api/v5/public/position-tiers", s.params)
	return request.DoList[PositionTier](req)
}

// PositionTier is one row of the position-tier schedule.
type PositionTier struct {
	Underlying         string          `json:"uly"`
	InstrumentFamily   string          `json:"instFamily"`
	InstrumentID       string          `json:"instId"`
	Tier               string          `json:"tier"`
	MinSize            decimal.Decimal `json:"minSz"`
	MaxSize            decimal.Decimal `json:"maxSz"`
	MMR                decimal.Decimal `json:"mmr"`
	IMR                decimal.Decimal `json:"imr"`
	MaxLeverage        decimal.Decimal `json:"maxLever"`
	OptionMarginFactor decimal.Decimal `json:"optMgnFactor"`
	QuoteMaxLoan       decimal.Decimal `json:"quoteMaxLoan"`
	BaseMaxLoan        decimal.Decimal `json:"baseMaxLoan"`
}

// GetInterestRateLoanQuotaService -- GET /api/v5/public/interest-rate-loan-quota (public)
//
// Returns the basic and VIP interest rates plus the margin-loan quota schedule.
type GetInterestRateLoanQuotaService struct {
	c *Client
}

func (c *Client) NewGetInterestRateLoanQuotaService() *GetInterestRateLoanQuotaService {
	return &GetInterestRateLoanQuotaService{c: c}
}

func (s *GetInterestRateLoanQuotaService) Do(ctx context.Context) (*InterestRateLoanQuota, error) {
	req := request.Get(ctx, s.c, "/api/v5/public/interest-rate-loan-quota")
	return request.DoOne[InterestRateLoanQuota](req)
}

// InterestRateLoanQuota is the platform's interest-rate and loan-quota schedule.
type InterestRateLoanQuota struct {
	Basic              []InterestRateBasic  `json:"basic"`
	VIP                []InterestRateLevel  `json:"vip"`
	Regular            []InterestRateLevel  `json:"regular"`
	Config             []InterestRateConfig `json:"config"`
	ConfigCurrencyList []InterestRateCcy    `json:"configCcyList"`
}

// InterestRateBasic is the basic per-currency interest rate and borrow quota.
type InterestRateBasic struct {
	Currency string          `json:"ccy"`
	Rate     decimal.Decimal `json:"rate"`
	Quota    decimal.Decimal `json:"quota"`
}

// InterestRateLevel is a VIP/regular tier's interest discount and quota
// coefficient.
type InterestRateLevel struct {
	Level                string          `json:"level"`
	LoanQuotaCoefficient decimal.Decimal `json:"loanQuotaCoef"`
	InterestRateDiscount decimal.Decimal `json:"irDiscount"`
}

// InterestRateConfig is one currency's per-level loan-quota config entry.
type InterestRateConfig struct {
	Currency     string          `json:"ccy"`
	Level        string          `json:"level"`
	Quota        decimal.Decimal `json:"quota"`
	StrategyType string          `json:"stgyType"`
}

// InterestRateCcy is a currency's configured base interest rate.
type InterestRateCcy struct {
	Currency string          `json:"ccy"`
	Rate     decimal.Decimal `json:"rate"`
}

// GetUnderlyingService -- GET /api/v5/public/underlying (public)
//
// Returns the underlyings available for a derivatives product line. The data is
// an array containing a single array of underlying strings.
type GetUnderlyingService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetUnderlyingService(instType InstType) *GetUnderlyingService {
	return &GetUnderlyingService{c: c, params: map[string]string{"instType": string(instType)}}
}

func (s *GetUnderlyingService) Do(ctx context.Context) ([][]string, error) {
	req := request.Get(ctx, s.c, "/api/v5/public/underlying", s.params)
	return request.DoList[[]string](req)
}

// GetInsuranceFundService -- GET /api/v5/public/insurance-fund (public)
//
// Returns the insurance-fund balance and its recent change records for a product
// line.
type GetInsuranceFundService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetInsuranceFundService(instType InstType) *GetInsuranceFundService {
	return &GetInsuranceFundService{c: c, params: map[string]string{"instType": string(instType)}}
}

// SetType filters by record type (e.g. liquidation_balance_deposit,
// bankruptcy_loss, platform_revenue).
func (s *GetInsuranceFundService) SetType(typ string) *GetInsuranceFundService {
	s.params["type"] = typ
	return s
}

// SetUly filters by underlying.
func (s *GetInsuranceFundService) SetUly(uly string) *GetInsuranceFundService {
	s.params["uly"] = uly
	return s
}

// SetInstFamily filters by instrument family.
func (s *GetInsuranceFundService) SetInstFamily(instFamily string) *GetInsuranceFundService {
	s.params["instFamily"] = instFamily
	return s
}

// SetCcy filters by currency (MARGIN/OPTION).
func (s *GetInsuranceFundService) SetCcy(ccy string) *GetInsuranceFundService {
	s.params["ccy"] = ccy
	return s
}

// SetAfter paginates to records earlier than the given time (older).
func (s *GetInsuranceFundService) SetAfter(t time.Time) *GetInsuranceFundService {
	s.params["after"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetBefore paginates to records later than the given time (newer).
func (s *GetInsuranceFundService) SetBefore(t time.Time) *GetInsuranceFundService {
	s.params["before"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetInsuranceFundService) SetLimit(limit int) *GetInsuranceFundService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetInsuranceFundService) Do(ctx context.Context) (*InsuranceFund, error) {
	req := request.Get(ctx, s.c, "/api/v5/public/insurance-fund", s.params)
	return request.DoOne[InsuranceFund](req)
}

// InsuranceFund is the insurance-fund summary and its change records.
type InsuranceFund struct {
	InstrumentType   InstType              `json:"instType"`
	InstrumentFamily string                `json:"instFamily"`
	Total            decimal.Decimal       `json:"total"`
	Details          []InsuranceFundDetail `json:"details"`
}

// InsuranceFundDetail is one insurance-fund change record.
type InsuranceFundDetail struct {
	Type                string          `json:"type"`
	Currency            string          `json:"ccy"`
	Amount              decimal.Decimal `json:"amt"`
	Balance             decimal.Decimal `json:"balance"`
	MaxBalance          decimal.Decimal `json:"maxBal"`
	MaxBalanceTimestamp time.Time       `json:"maxBalTs"`
	DecRate             decimal.Decimal `json:"decRate"`
	ADLType             string          `json:"adlType"`
	Timestamp           time.Time       `json:"ts"`
}

// GetConvertContractCoinService -- GET /api/v5/public/convert-contract-coin (public)
//
// Converts between a contract quantity and its coin amount for a derivatives
// instrument.
type GetConvertContractCoinService struct {
	c      *Client
	params map[string]string
}

// NewGetConvertContractCoinService builds the conversion. typ is "1" (coin -> to
// contract) or "2" (contract -> to coin); instId is the derivatives instrument;
// sz is the quantity to convert.
func (c *Client) NewGetConvertContractCoinService(typ, instId, sz string) *GetConvertContractCoinService {
	return &GetConvertContractCoinService{c: c, params: map[string]string{
		"type":   typ,
		"instId": instId,
		"sz":     sz,
	}}
}

// SetPx sets the order price (required for inverse contracts, ignored for
// linear).
func (s *GetConvertContractCoinService) SetPx(px string) *GetConvertContractCoinService {
	s.params["px"] = px
	return s
}

// SetUnit sets the coin unit ("coin" or "usds") when converting to coin.
func (s *GetConvertContractCoinService) SetUnit(unit string) *GetConvertContractCoinService {
	s.params["unit"] = unit
	return s
}

func (s *GetConvertContractCoinService) Do(ctx context.Context) (*ConvertContractCoin, error) {
	req := request.Get(ctx, s.c, "/api/v5/public/convert-contract-coin", s.params)
	return request.DoOne[ConvertContractCoin](req)
}

// ConvertContractCoin is the result of a contract/coin conversion.
type ConvertContractCoin struct {
	Type         string          `json:"type"`
	InstrumentID string          `json:"instId"`
	Price        decimal.Decimal `json:"px"`
	Size         decimal.Decimal `json:"sz"`
	Unit         string          `json:"unit"`
}

// GetInstrumentTickBandsService -- GET /api/v5/public/instrument-tick-bands (public)
//
// Returns the tick-size bands of OPTION instruments by instrument family.
type GetInstrumentTickBandsService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetInstrumentTickBandsService(instType InstType) *GetInstrumentTickBandsService {
	return &GetInstrumentTickBandsService{c: c, params: map[string]string{"instType": string(instType)}}
}

// SetInstFamily filters by instrument family.
func (s *GetInstrumentTickBandsService) SetInstFamily(instFamily string) *GetInstrumentTickBandsService {
	s.params["instFamily"] = instFamily
	return s
}

func (s *GetInstrumentTickBandsService) Do(ctx context.Context) ([]InstrumentTickBands, error) {
	req := request.Get(ctx, s.c, "/api/v5/public/instrument-tick-bands", s.params)
	return request.DoList[InstrumentTickBands](req)
}

// InstrumentTickBands is the tick-size band schedule of an instrument family.
type InstrumentTickBands struct {
	InstrumentType   InstType   `json:"instType"`
	InstrumentFamily string     `json:"instFamily"`
	SeriesID         string     `json:"seriesId"`
	TickBand         []TickBand `json:"tickBand"`
}

// TickBand is one price range and its applicable tick size.
type TickBand struct {
	MinPrice decimal.Decimal `json:"minPx"`
	MaxPrice decimal.Decimal `json:"maxPx"`
	TickSize decimal.Decimal `json:"tickSz"`
}

// GetPremiumHistoryService -- GET /api/v5/public/premium-history (public)
//
// Returns the premium-index history (spot vs perpetual mid) of an instrument.
type GetPremiumHistoryService struct {
	c      *Client
	params map[string]string
}

func (c *Client) NewGetPremiumHistoryService(instId string) *GetPremiumHistoryService {
	return &GetPremiumHistoryService{c: c, params: map[string]string{"instId": instId}}
}

// SetAfter paginates to records earlier than the given time (older).
func (s *GetPremiumHistoryService) SetAfter(t time.Time) *GetPremiumHistoryService {
	s.params["after"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetBefore paginates to records later than the given time (newer).
func (s *GetPremiumHistoryService) SetBefore(t time.Time) *GetPremiumHistoryService {
	s.params["before"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (max 100).
func (s *GetPremiumHistoryService) SetLimit(limit int) *GetPremiumHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetPremiumHistoryService) Do(ctx context.Context) ([]PremiumHistory, error) {
	req := request.Get(ctx, s.c, "/api/v5/public/premium-history", s.params)
	return request.DoList[PremiumHistory](req)
}

// PremiumHistory is one premium-index record.
type PremiumHistory struct {
	InstrumentID string          `json:"instId"`
	Premium      decimal.Decimal `json:"premium"`
	Timestamp    time.Time       `json:"ts"`
}
