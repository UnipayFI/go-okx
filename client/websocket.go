package client

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	okxCommon "github.com/UnipayFI/go-okx/common"
	"github.com/UnipayFI/go-okx/pkg/log"
	"github.com/gorilla/websocket"
)

// WebSocketClient holds the configuration shared by every OKX v5 stream. OKX
// splits its channels across three gateways: public (market data, no login),
// private (account/orders/positions, login required) and business (candles,
// algo/grid/copy-trading/earn, login required for the private business
// channels). Credentials are only needed for the private and business logins.
type WebSocketClient struct {
	publicURL    string
	privateURL   string
	businessURL  string
	apiKey       string
	apiSecret    string
	passphrase   string
	demoTrading  bool
	signFn       SignFn
	logger       log.Logger
	dialer       *websocket.Dialer
	readTimeout  time.Duration
	writeTimeout time.Duration
}

type WebSocketOption struct {
	network      okxCommon.Network
	publicURL    string
	privateURL   string
	businessURL  string
	apiKey       string
	apiSecret    string
	passphrase   string
	demoTrading  bool
	signFn       SignFn
	logger       log.Logger
	dialer       *websocket.Dialer
	readTimeout  time.Duration
	writeTimeout time.Duration
}

type WebSocketOptions func(*WebSocketOption)

func defaultWebSocketOption() *WebSocketOption {
	return &WebSocketOption{
		network:      okxCommon.Mainnet,
		logger:       log.GetDefaultLogger(),
		dialer:       defaultDialer(),
		readTimeout:  okxCommon.DEFAULT_WS_READ_TIMEOUT,
		writeTimeout: okxCommon.DEFAULT_WS_WRITE_TIMEOUT,
	}
}

func defaultDialer() *websocket.Dialer {
	return &websocket.Dialer{
		Proxy:             http.ProxyFromEnvironment,
		HandshakeTimeout:  45 * time.Second,
		EnableCompression: true,
	}
}

func NewWebSocketClient(options ...WebSocketOptions) *WebSocketClient {
	opt := defaultWebSocketOption()
	for _, option := range options {
		option(opt)
	}
	publicURL := opt.network.WsPublicURL()
	privateURL := opt.network.WsPrivateURL()
	businessURL := opt.network.WsBusinessURL()
	if opt.demoTrading {
		publicURL = okxCommon.DEFAULT_WS_DEMO_PUBLIC_URL
		privateURL = okxCommon.DEFAULT_WS_DEMO_PRIVATE_URL
		businessURL = okxCommon.DEFAULT_WS_DEMO_BUSINESS_URL
	}
	if opt.publicURL != "" {
		publicURL = opt.publicURL
	}
	if opt.privateURL != "" {
		privateURL = opt.privateURL
	}
	if opt.businessURL != "" {
		businessURL = opt.businessURL
	}
	return &WebSocketClient{
		publicURL:    publicURL,
		privateURL:   privateURL,
		businessURL:  businessURL,
		apiKey:       opt.apiKey,
		apiSecret:    opt.apiSecret,
		passphrase:   opt.passphrase,
		demoTrading:  opt.demoTrading,
		signFn:       opt.signFn,
		logger:       opt.logger,
		dialer:       opt.dialer,
		readTimeout:  opt.readTimeout,
		writeTimeout: opt.writeTimeout,
	}
}

func (c *WebSocketClient) GetPublicURL() string         { return c.publicURL }
func (c *WebSocketClient) GetPrivateURL() string        { return c.privateURL }
func (c *WebSocketClient) GetBusinessURL() string       { return c.businessURL }
func (c *WebSocketClient) GetAPIKey() string            { return c.apiKey }
func (c *WebSocketClient) GetAPISecret() string         { return c.apiSecret }
func (c *WebSocketClient) GetPassphrase() string        { return c.passphrase }
func (c *WebSocketClient) IsDemoTrading() bool          { return c.demoTrading }
func (c *WebSocketClient) GetSignFn() SignFn            { return c.signFn }
func (c *WebSocketClient) GetLogger() log.Logger        { return c.logger }
func (c *WebSocketClient) GetDialer() *websocket.Dialer { return c.dialer }

// GetReadTimeout is the max idle time a steady-state stream read may block with
// no inbound frame before the connection is considered dead. Zero means no
// read deadline is applied.
func (c *WebSocketClient) GetReadTimeout() time.Duration { return c.readTimeout }

// GetWriteTimeout bounds a single WS write (subscribe op, keepalive ping). Zero
// means no write deadline is applied.
func (c *WebSocketClient) GetWriteTimeout() time.Duration { return c.writeTimeout }

// WithWebSocketAuth sets the credentials used to log in to the private and
// business streams.
func WithWebSocketAuth(apiKey, apiSecret, passphrase string) WebSocketOptions {
	return func(opt *WebSocketOption) {
		opt.apiKey = apiKey
		opt.apiSecret = apiSecret
		opt.passphrase = passphrase
	}
}

func WithWebSocketNetwork(network okxCommon.Network) WebSocketOptions {
	return func(opt *WebSocketOption) { opt.network = network }
}

// WithWebSocketDemoTrading points the streams at OKX's demo (paper) trading
// gateways (wspap.okx.com).
func WithWebSocketDemoTrading(enabled bool) WebSocketOptions {
	return func(opt *WebSocketOption) { opt.demoTrading = enabled }
}

// WithWebSocketPublicURL overrides the public stream URL. Empty is ignored.
func WithWebSocketPublicURL(u string) WebSocketOptions {
	return func(opt *WebSocketOption) { opt.publicURL = u }
}

// WithWebSocketPrivateURL overrides the private stream URL. Empty is ignored.
func WithWebSocketPrivateURL(u string) WebSocketOptions {
	return func(opt *WebSocketOption) { opt.privateURL = u }
}

// WithWebSocketBusinessURL overrides the business stream URL. Empty is ignored.
func WithWebSocketBusinessURL(u string) WebSocketOptions {
	return func(opt *WebSocketOption) { opt.businessURL = u }
}

func WithWebSocketLogger(logger log.Logger) WebSocketOptions {
	return func(opt *WebSocketOption) { opt.logger = logger }
}

// WithWebSocketSignFn overrides the default HMAC login signer.
func WithWebSocketSignFn(signFn SignFn) WebSocketOptions {
	return func(opt *WebSocketOption) { opt.signFn = signFn }
}

// WithWebSocketReadTimeout overrides the steady-state read idle deadline (the
// max time a read may block with no inbound frame before the connection is
// treated as dead). A non-positive value disables the read deadline.
func WithWebSocketReadTimeout(d time.Duration) WebSocketOptions {
	return func(opt *WebSocketOption) { opt.readTimeout = d }
}

// WithWebSocketWriteTimeout overrides the per-write deadline applied to the
// subscribe op and keepalive pings. A non-positive value disables it.
func WithWebSocketWriteTimeout(d time.Duration) WebSocketOptions {
	return func(opt *WebSocketOption) { opt.writeTimeout = d }
}

// WithWebSocketProxy routes the stream dial through the given proxy (http,
// https, socks5, socks5h). Invalid URLs are logged and skipped.
func WithWebSocketProxy(proxyURL string) WebSocketOptions {
	return func(opt *WebSocketOption) {
		if proxyURL == "" {
			return
		}
		u, err := url.Parse(proxyURL)
		if err != nil {
			opt.logger.Errorf("WithWebSocketProxy: invalid proxy URL %q: %v", proxyURL, err)
			return
		}
		switch strings.ToLower(u.Scheme) {
		case "http", "https":
			opt.dialer.Proxy = http.ProxyURL(u)
			opt.dialer.NetDialContext = nil
		case "socks5", "socks5h":
			dialCtx, err := socks5DialContext(u)
			if err != nil {
				opt.logger.Errorf("WithWebSocketProxy: socks5 setup failed: %v", err)
				return
			}
			opt.dialer.Proxy = nil
			opt.dialer.NetDialContext = dialCtx
		default:
			opt.logger.Errorf("WithWebSocketProxy: unsupported scheme %q", u.Scheme)
		}
	}
}
