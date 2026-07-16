# go-okx

[![Go Reference](https://pkg.go.dev/badge/github.com/UnipayFI/go-okx.svg)](https://pkg.go.dev/github.com/UnipayFI/go-okx)
[![Go 1.26+](https://img.shields.io/badge/Go-1.26%2B-00ADD8?logo=go)](go.mod)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

A Go SDK for the [OKX](https://www.okx.com/docs-v5/en/) v5 API.

OKX v5 is a single **unified account** API: spot, margin, perpetual swaps,
futures, options, funding, sub-accounts, earn and copy-trading are all served
from one domain under the `/api/v5/`* path namespace and signed with one
HMAC-SHA256 scheme. This SDK wraps the complete REST surface plus the v5
WebSocket streams.


| API                                                         | Aligned to                                        |
| ----------------------------------------------------------- | ------------------------------------------------- |
| `/api/v5` REST + v5 WebSocket (public / private / business) | [2026-07-16](https://www.okx.com/docs-v5/log_en/#upcoming-changes) |


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
- Amounts as `decimal.Decimal`, ms timestamps as `time.Time` — OKX's
string-encoded numbers and `""`/`"0"`/`"-1"` "not set" sentinels are decoded
for you (no per-field format tags).
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



| Package              | Scope                                                                  |
| -------------------- | ---------------------------------------------------------------------- |
| `okx`                | the unified-account REST + WebSocket client (root package)             |
| `client/` `request/` | REST client, options, HMAC signer, envelope decode, WS subscribe/login |
| `common/`            | constants, global `time.Time` + `decimal.Decimal` JSON codec           |
| `cmd/okxraw/`        | dev tool: sign + dump any endpoint's raw response                      |


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