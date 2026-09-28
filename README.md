# go-okx

[![Go Reference](https://pkg.go.dev/badge/github.com/UnipayFI/go-okx.svg)](https://pkg.go.dev/github.com/UnipayFI/go-okx)
[![Go 1.27+](https://img.shields.io/badge/Go-1.27%2B-00ADD8?logo=go)](go.mod)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

A Go SDK for the [OKX](https://www.okx.com/docs-v5/en/) v5 API.

OKX v5 is a single **unified account** API: spot, margin, perpetual swaps,
futures, options, funding, sub-accounts, earn and copy-trading are all served
from one domain under the `/api/v5/`* path namespace and signed with one
HMAC-SHA256 scheme. This SDK wraps the complete REST surface plus the v5
WebSocket streams.


| API                                                         | Aligned to                                        |
| ----------------------------------------------------------- | ------------------------------------------------- |
| `/api/v5` REST + v5 WebSocket (public / private / business) | [2026-09-15](https://www.okx.com/docs-v5/log_en/#upcoming-changes) |


Response structs are reconciled against the **live API** (not just the docs), so
endpoints stay in sync with the date above: every public endpoint and every
private read endpoint is tested against production and the real JSON keys are
diffed against the typed structs.

## Install

```bash
go get github.com/UnipayFI/go-okx@latest
```

## Highlights

- One signing/transport core for the whole v5 API; a single `okx.Client`.
- Fluent per-endpoint API: `NewXxxService(...).SetFoo(...).Do(ctx)`.
- Amounts as `decimal.Decimal`, timestamps as `time.Time` whose wire format is
declared by a `format` tag option (`json:"cTime,format:unixmilli"`) — OKX's
string-encoded numbers and `""`/`"0"`/`"-1"` "not set" sentinels are decoded for
you.
- OKX's always-an-array `data` envelope handled by typed helpers
(`DoList`/`DoOne`/`DoObject`); batch order results expose per-item `sCode`.
- WebSocket: typed subscribe services over the public/private/business gateways,
automatic login + ping/pong keepalive, and order entry over a persistent
connection.

## Quick start

```go
package main

import (
	"context"
	"fmt"

	"github.com/UnipayFI/go-okx"
	"github.com/UnipayFI/go-okx/client"
	"github.com/shopspring/decimal"
)

func main() {
	ctx := context.Background()

	c := okx.NewClient(
		client.WithAuth("apiKey", "apiSecret", "passphrase"),
		// client.WithProxy("socks5://127.0.0.1:7890"),
		// client.WithDemoTrading(true),
	)
	_ = c.SyncServerTime(ctx) // align clock to avoid signature drift

	// Public market data (no auth).
	tickers, _ := c.NewGetTickersService(okx.InstTypeSpot).Do(ctx)
	fmt.Println(len(tickers), "spot tickers")

	// Private account data.
	bal, _ := c.NewGetBalanceService().Do(ctx)
	fmt.Println("total equity:", bal.TotalEq)

	// Place a limit order.
	ref, err := c.NewPlaceOrderService("BTC-USDT", okx.TdModeCash, okx.SideBuy,
		okx.OrdTypeLimit, decimal.RequireFromString("0.0001")).
		SetPx(decimal.RequireFromString("30000")).
		Do(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("ordId:", ref.OrdId, "sCode:", ref.SCode)
}
```

## Authentication

Pass credentials from the OKX API-management page (the passphrase is the one set
when the key was created):

```go
c := okx.NewClient(client.WithAuth(apiKey, apiSecret, passphrase))
```

Requests are signed with HMAC-SHA256 over
`timestamp + method + requestPath(+ "?" + query) + body`, base64-encoded into the
`OK-ACCESS-SIGN` header, with an ISO-8601 millisecond UTC `OK-ACCESS-TIMESTAMP`.
For an RSA key or external signer, pass `client.WithSignFn(fn)`.

Other options: `WithProxy` (http/https/socks5), `WithBaseURL`, `WithDemoTrading`
(routes to OKX's paper-trading via the `x-simulated-trading` header),
`WithTimeOffset`, `WithLogger`, `WithHTTPClient`.

## Response envelope

OKX always returns `{"code":"0","msg":"","data":[ ... ]}` with `data` as an
array. Service methods return the natural Go shape:

- list endpoints → `([]T, error)`
- single-object endpoints (balance, config, place-order ack, …) → `(*T, error)`
- batch order place/cancel/amend → `([]T, error)` whose items carry `sCode`/`sMsg`

A request-level failure is returned as a `*client.APIError`: any non-`"0"`
code, except `"1"`/`"2"` from order and ack endpoints that still return data —
their items carry `sCode`/`sMsg`, which each item's `Err()` returns as a
`*client.APIError` (nil on success). Use `client.AsAPIError` / `client.IsCode`
to inspect wrapped errors. `APIError` has a value receiver, so under Go 1.27
`go vet` (which `go test` runs) rejects `fmt.Errorf("%w", apiErr)` for an
`apiErr` of type `*client.APIError` — wrap the original `error` instead.

## WebSocket

```go
ws := okx.NewWebSocketClient(
	client.WithWebSocketAuth(apiKey, apiSecret, passphrase), // private/business channels
)

// Public ticker (public gateway, no login).
done, _, _ := ws.NewSubscribeTickersService("BTC-USDT").
	Do(ctx, func(p *request.WsPush[[]okx.WsTicker], err error) {
		if err != nil {
			return
		}
		fmt.Println(p.Arg.InstID, p.Data[0].Last)
	})
close(done) // unsubscribe + close

// Private account (auto login).
ws.NewSubscribeAccountService().Do(ctx, func(p *request.WsPush[[]okx.WsAccount], err error) {
	// p.Data[0].TotalEq, ...
})
```

Each `Do` returns `(done chan<- struct{}, stop <-chan struct{}, err error)`:
close `done` to unsubscribe; `stop` closes when the reader exits. Ping/pong
keepalive is automatic.

Orders can also be placed over a persistent, logged-in connection — a
low-latency alternative to the REST trade endpoints:

```go
tc, _ := ws.DialTrade(ctx) // connect + login
defer tc.Close()
ack, _ := tc.PlaceOrder(ctx, okx.OrderArg{
	InstId: "BTC-USDT", TdMode: okx.TdModeCash, Side: okx.SideBuy,
	OrdType: okx.OrdTypeLimit, Sz: decimal.RequireFromString("0.0001"),
	Px: decimal.RequireFromString("30000"),
})
// tc.AmendOrder / tc.CancelOrder / tc.BatchPlaceOrders / ...
```

## JSON and timestamps

The SDK uses Go 1.27's `encoding/json/v2` (keep the default `jsonv2` GOEXPERIMENT
enabled; without it the SDK does not compile). Every `time.Time` field declares
the format OKX actually sends with the `format` tag option, e.g.
`json:"cTime,format:unixmilli"` — experimental in Go 1.27, and enabled by the
SDK's codec. `common.JSONMarshal` / `common.JSONUnmarshal` apply that option's
semantics plus OKX's quirks: numbers are read quoted or bare and written quoted
in whole units; `""`/`"0"`/`"-1"`/`null` read as the zero time, which is written
as `""`; and `estCompleteTime`'s UTC+8 wall-clock text uses the SDK-specific
`okxUTC8WallClock` format. Decoded times are in UTC — use `.Equal` to compare and
`.In(loc)` / `.Local()` to display.

A `time.Time` without a `format` option is RFC 3339, as in the standard library.
That includes your own types passed through `common.JSONMarshal` /
`common.JSONUnmarshal` or the `request` helpers, which earlier versions encoded
and decoded as milliseconds: tag such fields `format:unixmilli`, and put a
`time.Time` into a `request.Post` body map as
`strconv.FormatInt(t.UnixMilli(), 10)`.

Go 1.27 only honours `format` tags when the experimental
`ExperimentalSupportFormatTag` option is passed, so serialize SDK types with
`common.JSONMarshal` / `common.JSONUnmarshal`. Where you only need
`encoding/json/v2` to accept the tags (say, to store or log SDK values), you can
pass `github.com/go-json-experiment/json.ExperimentalSupportFormatTag(true)`
yourself; that applies the option's plain semantics (bare numbers, no
`""`/`"0"`/`"-1"` sentinels, no `okxUTC8WallClock`), so it cannot read OKX's
wire data. Plain `encoding/json` returns an error for structs with `format` tags,
and `log/slog`'s JSON handler logs `!ERROR:...` in place of such a value.

## Packages


| Area                                    | Files                                                                   |
| --------------------------------------- | ----------------------------------------------------------------------- |
| Public / market data                    | `public_data.go` `market.go` `market_index.go` `rubik.go` `status.go`   |
| Trading account                         | `account.go` `account_bills.go` `account_borrow.go` `account_config.go` |
| Trade                                   | `trade_order.go` `trade_fills.go` `trade_convert.go`                    |
| Algo / Grid / Recurring                 | `algo.go` `grid.go` `recurring.go`                                      |
| Funding / Convert / Sub-account         | `asset.go` `convert.go` `subaccount.go`                                 |
| Financial (earn)                        | `finance_savings.go` `finance_staking.go` `finance_loan.go`             |
| Copy / Block (RFQ) / Spread / Affiliate | `copytrading.go` `rfq.go` `sprd.go` `affiliate.go`                      |
| WebSocket                               | `ws_public.go` `ws_business.go` `ws_private.go` `ws_trade.go`           |



| Package              | Scope                                                                                |
| -------------------- | ------------------------------------------------------------------------------------ |
| `okx`                | the unified-account REST + WebSocket client (root package)                           |
| `client/` `request/` | REST client, options, HMAC signer, envelope decode, WS subscribe/login               |
| `common/`            | constants, `encoding/json/v2` codec: `format`-tagged `time.Time` + `decimal.Decimal` |
| `cmd/okxraw/`        | dev tool: sign + dump any endpoint's raw response                                    |


## Testing

Tests hit the live API and read credentials from the environment, skipping when
unset:

```bash
export OKX_API_KEY=...  OKX_API_SECRET=...  OKX_PASSPHRASE=...
export OKX_PROXY=socks5://127.0.0.1:7890   # optional
export OKX_DEMO=1                            # optional: paper trading

go test ./... -run TestAccount -v                 # one module at a time
OKX_TEST_WRITE=1 go test . -run TestOrderLifecycle # live order test (tiny, reversible)
```

- Run **per module** (`-run TestXxx`) — OKX rate-limits per endpoint, so the full
suite can trip `50011 Too Many Requests` (the test helpers auto-retry it).
- Capability-gated reads (copy-trading, spread, RFQ, fixed-loan, sub-account)
are tolerated when the account lacks the capability — signing is still
exercised.
- State-changing endpoints are implemented but never executed by the suite,
except the gated `TestOrderLifecycle` (a tiny far-below-market post_only order
on a large-cap pair that is immediately cancelled). **Withdrawal is
implemented but never tested.**

The `cmd/okxraw` helper dumps any endpoint's raw signed response:

```bash
go run ./cmd/okxraw GET /api/v5/account/config
go run ./cmd/okxraw GET /api/v5/account/bills "instType=SPOT&limit=5"
```

## Changelog

- **2026-06-24** — Initial release. Full OKX v5 REST coverage (trading account,
order-book/algo/grid/recurring trading, funding, convert, sub-account,
financial products, copy-trading, block (RFQ) & spread trading, public/market
data, trading statistics, status) plus the v5 WebSocket streams (public,
private and business gateways + WebSocket order entry). Aligned to the OKX v5
docs as of 2026-06-24.

## License

[MIT](LICENSE)