// Command okxraw signs and executes a single OKX v5 REST call and pretty prints
// the raw response. It is a development aid for capturing the exact shape of
// private endpoints (which cannot be curled without HMAC signing) so the typed
// response structs can be reconciled against reality.
//
// Usage:
//
//	OKX_API_KEY=... OKX_API_SECRET=... OKX_PASSPHRASE=... \
//	  go run ./cmd/okxraw GET  /api/v5/account/config
//	  go run ./cmd/okxraw GET  /api/v5/account/bills "instType=SPOT&limit=5"
//	  go run ./cmd/okxraw POST /api/v5/trade/order '{"instId":"BTC-USDT", ...}'
//
// The third argument is the query string (GET) or JSON body (POST). Set
// OKX_PROXY to route through a proxy.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/UnipayFI/go-okx/client"
	okxCommon "github.com/UnipayFI/go-okx/common"
	"github.com/UnipayFI/go-okx/request"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: okxraw <GET|POST> <path> [query-or-jsonbody]")
		os.Exit(2)
	}
	method := strings.ToUpper(os.Args[1])
	path := os.Args[2]
	arg := ""
	if len(os.Args) > 3 {
		arg = os.Args[3]
	}

	opts := []client.Options{
		client.WithAuth(
			os.Getenv("OKX_API_KEY"),
			os.Getenv("OKX_API_SECRET"),
			os.Getenv("OKX_PASSPHRASE"),
		),
	}
	if proxy := os.Getenv("OKX_PROXY"); proxy != "" {
		opts = append(opts, client.WithProxy(proxy))
	}
	if os.Getenv("OKX_DEMO") != "" {
		opts = append(opts, client.WithDemoTrading(true))
	}
	c := client.NewClient(opts...)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var req *request.Request
	switch method {
	case "GET":
		req = request.Get(ctx, c, path, parseQuery(arg)).WithSign()
	case "POST":
		body := map[string]any{}
		if arg != "" {
			if err := okxCommon.JSONUnmarshal([]byte(arg), &body); err != nil {
				fail("invalid json body: %v", err)
			}
		}
		req = request.Post(ctx, c, path, body).WithSign()
	default:
		fail("unsupported method %q", method)
	}

	body, err := request.DoRaw(req)
	if err != nil {
		fail("request error: %v", err)
	}
	fmt.Println(pretty(body))
}

func parseQuery(q string) map[string]string {
	out := map[string]string{}
	q = strings.TrimPrefix(q, "?")
	for pair := range strings.SplitSeq(q, "&") {
		if pair == "" {
			continue
		}
		k, v, _ := strings.Cut(pair, "=")
		out[k] = v
	}
	return out
}

func pretty(b []byte) string {
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return string(b)
	}
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return string(b)
	}
	return string(out)
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
