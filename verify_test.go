package okx

import (
	"context"
	"errors"
	"maps"
	"sort"
	"testing"
	"time"

	"github.com/UnipayFI/go-okx/client"
	"github.com/UnipayFI/go-okx/common"
	"github.com/UnipayFI/go-okx/request"
)

// doRawWithRetry runs a raw fetch, retrying on OKX's transient rate-limit code
// 50011 ("Too Many Requests") with a short backoff. Some OKX endpoints
// (exchange-rate, announcement-types, ...) have very tight per-endpoint limits
// that the call-then-fetchRaw test pattern can trip; this keeps the live suite
// deterministic without weakening the real client.
func doRawWithRetry(t *testing.T, label string, build func() *request.Request) []byte {
	t.Helper()
	var lastErr error
	for attempt := 0; attempt < 5; attempt++ {
		raw, err := request.DoRawData(build())
		if err == nil {
			return raw
		}
		lastErr = err
		var apiErr *client.APIError
		if errors.As(err, &apiErr) && apiErr.Code == "50011" {
			time.Sleep(time.Duration(attempt+1) * 700 * time.Millisecond)
			continue
		}
		break
	}
	t.Fatalf("raw %s: %v", label, lastErr)
	return nil
}

// tolerable reports whether err is an expected "this account lacks the
// capability or simply has no data for these params" OKX response rather than a
// code bug, so capability-gated read tests (copy-trading / broker / earn /
// sub-account) can treat it as a pass: the request path and signing were
// correct, the account just isn't enrolled in that product or has no data.
func tolerable(t *testing.T, label string, err error, codes ...string) bool {
	t.Helper()
	var apiErr *client.APIError
	if errors.As(err, &apiErr) {
		for _, code := range codes {
			if apiErr.Code == code {
				t.Logf("%s: account lacks this capability/data (code=%s) — endpoint+signing OK", label, apiErr.Code)
				return true
			}
		}
	}
	return false
}

// fetchRawGet returns the raw JSON of the "data" array for a GET endpoint, used
// to diff the real response shape against the typed structs.
func fetchRawGet(t *testing.T, c *Client, ctx context.Context, path string, params map[string]string, sign bool) []byte {
	t.Helper()
	return doRawWithRetry(t, "GET "+path, func() *request.Request {
		req := request.Get(ctx, c, path, params)
		if sign {
			req = req.WithSign()
		}
		return req
	})
}

// fetchRawPost mirrors fetchRawGet for POST endpoints.
func fetchRawPost(t *testing.T, c *Client, ctx context.Context, path string, body map[string]any, sign bool) []byte {
	t.Helper()
	return doRawWithRetry(t, "POST "+path, func() *request.Request {
		req := request.Post(ctx, c, path, body)
		if sign {
			req = req.WithSign()
		}
		return req
	})
}

// assertCovers checks that every JSON key present in the real response (raw) is
// also produced by marshaling the typed value. It compares key *sets* (not
// values), recursing into nested objects and merging array elements, so a
// missing struct field surfaces as an uncovered key path. This is the backbone
// of the "fields match the real API" guarantee.
//
// OKX always returns the envelope's "data" as a JSON array (a single-object
// result is wrapped in a one-element array). The typed value may be a slice (for
// list endpoints) or a single struct/pointer (for single-object endpoints); in
// the latter case it is wrapped into a one-element array so both sides line up.
func assertCovers(t *testing.T, label string, raw []byte, typed any) {
	t.Helper()
	var rawAny any
	if err := common.JSONUnmarshal(raw, &rawAny); err != nil {
		t.Fatalf("%s: unmarshal raw: %v", label, err)
	}
	typedBytes, err := common.JSONMarshal(typed)
	if err != nil {
		t.Fatalf("%s: marshal typed: %v", label, err)
	}
	var haveAny any
	if err := common.JSONUnmarshal(typedBytes, &haveAny); err != nil {
		t.Fatalf("%s: unmarshal typed: %v", label, err)
	}
	// OKX usually returns "data" as an array; a single-object endpoint marshals
	// to an object, so wrap it to line up with the raw array. The handful of
	// object-data endpoints (e.g. rubik support-coin) return a raw object too —
	// then both sides are objects and no wrapping is needed.
	if _, rawIsArr := rawAny.([]any); rawIsArr {
		if _, ok := haveAny.(map[string]any); ok {
			haveAny = []any{haveAny}
		}
	}

	var missing []string
	diffKeys(rawAny, haveAny, "", &missing)
	if len(missing) > 0 {
		sort.Strings(missing)
		t.Errorf("%s: %d field(s) in real response NOT captured by struct:\n  %v", label, len(missing), missing)
		return
	}
	t.Logf("%s: OK, all response keys covered by struct", label)
}

// diffKeys walks raw and records the paths of keys absent from have.
func diffKeys(raw, have any, path string, out *[]string) {
	switch r := raw.(type) {
	case map[string]any:
		h, ok := have.(map[string]any)
		if !ok {
			*out = append(*out, path+" (expected object)")
			return
		}
		for k, rv := range r {
			child := path + "/" + k
			hv, present := h[k]
			if !present {
				*out = append(*out, child)
				continue
			}
			diffKeys(rv, hv, child, out)
		}
	case []any:
		h, ok := have.([]any)
		if !ok || len(r) == 0 || len(h) == 0 {
			return
		}
		// Merge keys across all raw elements so optional fields present only on
		// some rows are still checked against the (single-shape) struct.
		merged := map[string]any{}
		for _, e := range r {
			if em, ok := e.(map[string]any); ok {
				maps.Copy(merged, em)
			}
		}
		if len(merged) > 0 {
			diffKeys(merged, h[0], path+"[]", out)
		} else {
			diffKeys(r[0], h[0], path+"[]", out)
		}
	}
}
