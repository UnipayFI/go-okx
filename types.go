package okx

// This file holds the cross-cutting enum vocabulary used by many OKX v5
// endpoints. Enums that appear in only a single endpoint live in that
// endpoint's file. All OKX enums are plain JSON strings.

// InstType is an OKX instrument type. It selects the product line for market and
// account queries.
type InstType string

const (
	InstTypeSpot    InstType = "SPOT"
	InstTypeMargin  InstType = "MARGIN"
	InstTypeSwap    InstType = "SWAP"
	InstTypeFutures InstType = "FUTURES"
	InstTypeOption  InstType = "OPTION"
	InstTypeEvents  InstType = "EVENTS"
	InstTypeAny     InstType = "ANY"
)

// Side is an order side.
type Side string

const (
	SideBuy  Side = "buy"
	SideSell Side = "sell"
)

// PosSide is a position side. "net" is used in net mode; "long"/"short" in
// long/short mode.
type PosSide string

const (
	PosSideLong  PosSide = "long"
	PosSideShort PosSide = "short"
	PosSideNet   PosSide = "net"
)

// TdMode is the trade (margin) mode of an order.
type TdMode string

const (
	TdModeCash         TdMode = "cash"
	TdModeIsolated     TdMode = "isolated"
	TdModeCross        TdMode = "cross"
	TdModeSpotIsolated TdMode = "spot_isolated"
)

// OrdType is an order type. It doubles as the time-in-force selector OKX folds
// into the order type (fok / ioc / post_only).
type OrdType string

const (
	OrdTypeMarket          OrdType = "market"
	OrdTypeLimit           OrdType = "limit"
	OrdTypePostOnly        OrdType = "post_only"
	OrdTypeFOK             OrdType = "fok"
	OrdTypeIOC             OrdType = "ioc"
	OrdTypeOptimalLimitIOC OrdType = "optimal_limit_ioc"
	OrdTypeMMP             OrdType = "mmp"
	OrdTypeMMPAndPostOnly  OrdType = "mmp_and_post_only"
)

// OrdState is the lifecycle state of an order.
type OrdState string

const (
	OrdStateLive            OrdState = "live"
	OrdStatePartiallyFilled OrdState = "partially_filled"
	OrdStateFilled          OrdState = "filled"
	OrdStateCanceled        OrdState = "canceled"
	OrdStateMMPCanceled     OrdState = "mmp_canceled"
)

// MgnMode is a margin mode (used by positions, leverage and borrowing).
type MgnMode string

const (
	MgnModeIsolated MgnMode = "isolated"
	MgnModeCross    MgnMode = "cross"
)

// TgtCcy selects whether a spot market order's size is denominated in the base
// or the quote currency.
type TgtCcy string

const (
	TgtCcyBase  TgtCcy = "base_ccy"
	TgtCcyQuote TgtCcy = "quote_ccy"
)

// ExecType is the liquidity role of a fill.
type ExecType string

const (
	ExecTypeTaker ExecType = "T"
	ExecTypeMaker ExecType = "M"
)

// CtType is a futures/swap contract type.
type CtType string

const (
	CtTypeLinear  CtType = "linear"
	CtTypeInverse CtType = "inverse"
)

// OptType is an option type (call or put).
type OptType string

const (
	OptTypeCall OptType = "C"
	OptTypePut  OptType = "P"
)

// InstState is the listing state of an instrument.
type InstState string

const (
	InstStateLive    InstState = "live"
	InstStateSuspend InstState = "suspend"
	InstStatePreOpen InstState = "preopen"
	InstStateTest    InstState = "test"
)
