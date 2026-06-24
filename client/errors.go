package client

import "fmt"

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
func IsAPIError(err error) bool {
	_, ok := err.(*APIError)
	return ok
}
