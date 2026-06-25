package client

import (
	"errors"
	"fmt"
)

// APIError is the error envelope OKX returns when a request fails. OKX codes are
// strings; "0" means success and anything else is an error. For batch order
// endpoints the top-level code can also be "1" (all failed) or "2" (partial),
// with the real per-order status carried in each data item's sCode/sMsg.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
}

// Error returns the error code and message.
func (e APIError) Error() string {
	return fmt.Sprintf("<APIError> code=%s, msg=%s", e.Code, e.Message)
}

// IsValid reports whether e represents an actual API-level error (a non-empty,
// non-success code).
func (e APIError) IsValid() bool {
	return e.Code != "" && e.Code != "0"
}

// IsAPIError reports whether err is an OKX *APIError.
//
// Deprecated: this is a bare type assertion and does not unwrap. Prefer
// AsAPIError, which walks the error chain.
func IsAPIError(err error) bool {
	_, ok := err.(*APIError)
	return ok
}

// AsAPIError extracts the OKX APIError carried by err, walking the error chain
// (errors.As). It matches whether the chain holds a *APIError (as the SDK's
// request layer returns) or a value APIError (as a caller's %w-wrapping may
// embed), so it is robust to either wrapping style. Returns (nil, false) when
// err carries no APIError.
func AsAPIError(err error) (*APIError, bool) {
	if err == nil {
		return nil, false
	}
	var ptr *APIError
	if errors.As(err, &ptr) {
		return ptr, true
	}
	var val APIError
	if errors.As(err, &val) {
		return &val, true
	}
	return nil, false
}

// IsCode reports whether err carries an OKX APIError with the given code (e.g.
// "51400" for "order does not exist"). It unwraps via AsAPIError, so it works on
// wrapped errors — the typical use is collapsing an idempotent failure (cancel
// of an already-gone order) to a no-op.
func IsCode(err error, code string) bool {
	apiErr, ok := AsAPIError(err)
	return ok && apiErr.Code == code
}
