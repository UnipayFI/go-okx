package okx

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/UnipayFI/go-okx/client"
)

// testClient builds an authenticated OKX client from environment variables.
// Tests that need private endpoints skip themselves when credentials are absent
// so the suite stays runnable without secrets.
func testClient(t *testing.T) *Client {
	t.Helper()
	apiKey := os.Getenv("OKX_API_KEY")
	apiSecret := os.Getenv("OKX_API_SECRET")
	passphrase := os.Getenv("OKX_PASSPHRASE")
	if apiKey == "" || apiSecret == "" || passphrase == "" {
		t.Skip("OKX_API_KEY/SECRET/PASSPHRASE not set; skipping private test")
	}
	opts := []client.Options{
		client.WithAuth(apiKey, apiSecret, passphrase),
	}
	if proxy := os.Getenv("OKX_PROXY"); proxy != "" {
		opts = append(opts, client.WithProxy(proxy))
	}
	if os.Getenv("OKX_DEMO") != "" {
		opts = append(opts, client.WithDemoTrading(true))
	}
	return NewClient(opts...)
}

func testPublicClient() *Client {
	opts := []client.Options{}
	if proxy := os.Getenv("OKX_PROXY"); proxy != "" {
		opts = append(opts, client.WithProxy(proxy))
	}
	return NewClient(opts...)
}

func ctx(t *testing.T) context.Context {
	t.Helper()
	c, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	t.Cleanup(cancel)
	return c
}

// TestSystemTime is the foundation smoke test: it exercises the full public
// request → envelope-decode → typed-struct path end to end. Per-module endpoint
// coverage lives in each module's own _test.go.
func TestSystemTime(t *testing.T) {
	c := testPublicClient()
	resp, err := c.NewGetSystemTimeService().Do(ctx(t))
	if err != nil {
		t.Fatalf("system time: %v", err)
	}
	t.Logf("serverTime=%s (%d)", resp.Timestamp, resp.Timestamp.UnixMilli())
	if resp.Timestamp.IsZero() {
		t.Fatal("server time is zero")
	}
}
