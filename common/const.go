package common

import "time"

const (
	GO_OKX_USER_AGENT = "go-okx/1.0"

	// REST endpoint. OKX serves every business line (trading account, funding,
	// market data, sub-account, earn, ...) from a single unified-account domain;
	// the product is encoded in the request path, not the host. Demo (paper)
	// trading runs on the same domain and is toggled with the
	// "x-simulated-trading: 1" header (see client.WithDemoTrading).
	DEFAULT_REST_BASE_URL = "https://www.okx.com"

	// WebSocket endpoints. OKX splits its v5 streams across three gateways:
	// public (market data, no login), private (account/orders/positions, login
	// required) and business (candles, algo orders, copy-trading, earn, ...).
	DEFAULT_WS_PUBLIC_URL   = "wss://ws.okx.com:8443/ws/v5/public"
	DEFAULT_WS_PRIVATE_URL  = "wss://ws.okx.com:8443/ws/v5/private"
	DEFAULT_WS_BUSINESS_URL = "wss://ws.okx.com:8443/ws/v5/business"

	// Demo (paper) trading WebSocket endpoints. Demo uses the dedicated
	// wspap.okx.com host plus the "x-simulated-trading: 1" connect header.
	DEFAULT_WS_DEMO_PUBLIC_URL   = "wss://wspap.okx.com:8443/ws/v5/public"
	DEFAULT_WS_DEMO_PRIVATE_URL  = "wss://wspap.okx.com:8443/ws/v5/private"
	DEFAULT_WS_DEMO_BUSINESS_URL = "wss://wspap.okx.com:8443/ws/v5/business"

	// TimestampLayout is the ISO-8601 UTC millisecond format OKX requires for the
	// OK-ACCESS-TIMESTAMP header, e.g. "2020-12-08T09:08:57.715Z". For UTC times
	// the trailing "Z07:00" renders as a literal "Z".
	TimestampLayout = "2006-01-02T15:04:05.000Z07:00"

	DEFAULT_KEEP_ALIVE_INTERVAL = 25 * time.Second
	DEFAULT_KEEP_ALIVE_TIMEOUT  = 60 * time.Second

	// DEFAULT_WS_READ_TIMEOUT bounds how long a steady-state stream read may
	// block with no inbound frame before it is treated as a dead (half-open)
	// connection. OKX disconnects a silent socket after ~30s and the client
	// pings every DEFAULT_KEEP_ALIVE_INTERVAL, so a healthy stream always sees
	// traffic well within this window; exceeding it means the peer is gone.
	DEFAULT_WS_READ_TIMEOUT = DEFAULT_KEEP_ALIVE_TIMEOUT

	// DEFAULT_WS_WRITE_TIMEOUT bounds a single WS write (subscribe op, keepalive
	// ping). A write that blocks past this on a half-open socket is treated as a
	// failed connection.
	DEFAULT_WS_WRITE_TIMEOUT = 10 * time.Second
)

// Network identifies which OKX environment a client targets. OKX has no separate
// testnet host; demo (paper) trading runs on the same REST domain and is toggled
// with the "x-simulated-trading: 1" header (see client.WithDemoTrading). The
// type is kept for forward symmetry with sibling SDKs and to leave room for a
// future dedicated environment.
type Network int

const (
	Mainnet Network = iota
)

// RestBaseURL returns the REST base URL for this network.
func (n Network) RestBaseURL() string {
	return DEFAULT_REST_BASE_URL
}

// WsPublicURL returns the public WebSocket URL for this network.
func (n Network) WsPublicURL() string {
	return DEFAULT_WS_PUBLIC_URL
}

// WsPrivateURL returns the private WebSocket URL for this network.
func (n Network) WsPrivateURL() string {
	return DEFAULT_WS_PRIVATE_URL
}

// WsBusinessURL returns the business WebSocket URL for this network.
func (n Network) WsBusinessURL() string {
	return DEFAULT_WS_BUSINESS_URL
}
