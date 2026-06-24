package client

import (
	"time"

	"github.com/UnipayFI/go-okx/pkg/log"
	"github.com/go-resty/resty/v2"
)

// Client is the shared, product-agnostic REST core. OKX's v5 API is a single
// unified account, so every business line (trading account, funding,
// sub-account, earn, ...) is just a set of request paths layered on top of this
// same signing + transport machinery; the core carries no product-specific
// state.
type Client struct {
	client *resty.Client

	apiKey       string
	apiSecret    string
	passphrase   string
	demoTrading  bool
	logger       log.Logger
	signFn       SignFn
	timeOffsetMs int64
}

func NewClient(options ...Options) *Client {
	opt := defaultOption()
	for _, option := range options {
		option(opt)
	}

	baseURL := opt.network.RestBaseURL()
	if opt.baseURL != "" {
		baseURL = opt.baseURL
	}
	opt.client.SetBaseURL(baseURL)

	return &Client{
		client:       opt.client,
		apiKey:       opt.apiKey,
		apiSecret:    opt.apiSecret,
		passphrase:   opt.passphrase,
		demoTrading:  opt.demoTrading,
		logger:       opt.logger,
		signFn:       opt.signFn,
		timeOffsetMs: opt.timeOffsetMs,
	}
}

func (c *Client) GetHttpClient() *resty.Client { return c.client }

func (c *Client) GetAPIKey() string { return c.apiKey }

func (c *Client) GetAPISecret() string { return c.apiSecret }

func (c *Client) GetPassphrase() string { return c.passphrase }

func (c *Client) IsDemoTrading() bool { return c.demoTrading }

func (c *Client) GetLogger() log.Logger { return c.logger }

func (c *Client) GetSignFn() SignFn { return c.signFn }

func (c *Client) GetTimeOffsetMs() int64 { return c.timeOffsetMs }

func (c *Client) SetTimeOffset(offsetMs int64) { c.timeOffsetMs = offsetMs }

// Now returns the current request time adjusted by the configured client/server
// clock offset. The REST layer formats it as an ISO-8601 millisecond string for
// the OK-ACCESS-TIMESTAMP header; the WebSocket login uses its Unix epoch
// seconds. OKX rejects requests whose timestamp drifts more than 30s from its
// own clock, so the offset (set via SyncServerTime) keeps signed requests valid.
func (c *Client) Now() time.Time {
	return time.Now().Add(-time.Duration(c.timeOffsetMs) * time.Millisecond)
}
