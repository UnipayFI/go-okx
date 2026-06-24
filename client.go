// Package okx is a Go SDK for the OKX v5 API.
//
// OKX v5 is a single unified-account API: spot, margin, perpetual swaps,
// futures, options, funding, sub-accounts, earn and copy-trading are all served
// from one domain (https://www.okx.com) under the /api/v5/* path namespace and
// signed with one HMAC-SHA256 scheme. This package exposes every REST endpoint
// (and, in the ws_*.go files, the v5 WebSocket streams) through a fluent,
// per-endpoint Service API.
//
// Install: go get github.com/UnipayFI/go-okx
// Import:  import "github.com/UnipayFI/go-okx"
//
// Authentication uses the OK-ACCESS-KEY / OK-ACCESS-SIGN / OK-ACCESS-TIMESTAMP /
// OK-ACCESS-PASSPHRASE headers; the sign is base64(HMAC-SHA256(secret, prehash))
// over timestamp + method + requestPath + body, with an ISO-8601 millisecond
// UTC timestamp.
//
// Quick start:
//
//	c := okx.NewClient(client.WithAuth(apiKey, apiSecret, passphrase))
//	if err := c.SyncServerTime(ctx); err != nil { /* ... */ }
//
//	// Public market data (no auth).
//	tickers, err := c.NewGetTickersService(okx.InstTypeSpot).Do(ctx)
//
//	// Private account data.
//	bal, err := c.NewGetBalanceService().Do(ctx)
package okx

import (
	"context"

	"github.com/UnipayFI/go-okx/client"
	"github.com/UnipayFI/go-okx/request"
)

var _ request.Client = (*Client)(nil)

// Client is the REST client for OKX's v5 unified-account (/api/v5/*) endpoints.
// It embeds the shared core client, which holds the signing/transport machinery
// and the client/server clock offset.
type Client struct {
	*client.Client
}

// NewClient constructs an OKX v5 REST client.
func NewClient(options ...client.Options) *Client {
	return &Client{client.NewClient(options...)}
}

// SyncServerTime measures the client/server clock offset and stores it so that
// signed requests carry a timestamp the server accepts. OKX rejects requests
// whose OK-ACCESS-TIMESTAMP drifts more than 30s from its own clock, so call
// this once at startup (and periodically for long-lived processes).
func (c *Client) SyncServerTime(ctx context.Context) error {
	localBefore := c.Now().UnixMilli()
	resp, err := c.NewGetSystemTimeService().Do(ctx)
	if err != nil {
		return err
	}
	localAfter := c.Now().UnixMilli()
	local := (localBefore + localAfter) / 2
	c.SetTimeOffset(local - resp.Timestamp.UnixMilli())
	c.GetLogger().Infof("Time sync: local=%d, server=%d, offset=%dms",
		local, resp.Timestamp.UnixMilli(), c.GetTimeOffsetMs())
	return nil
}
