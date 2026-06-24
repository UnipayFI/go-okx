package request

import (
	"errors"
	"fmt"

	"github.com/UnipayFI/go-okx/client"
	"github.com/UnipayFI/go-okx/common"
	"github.com/go-json-experiment/json/jsontext"
)

// apiResponse is OKX's uniform REST envelope. "code" is "0" on success; "data"
// is ALWAYS a JSON array (a single-object result is wrapped in a one-element
// array). "msg" carries the error message when code != "0".
type apiResponse[T any] struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
	Data    []T    `json:"data"`
}

// DoList executes the request and decodes the envelope's data array into []T. A
// non-"0" code is returned as a *client.APIError. Most OKX endpoints (both list
// and single-object) funnel through here; single-object services wrap this with
// DoOne.
func DoList[T any](r *Request) ([]T, error) {
	out, err := do[T](r)
	if err != nil {
		return nil, err
	}
	if out.Code != "0" {
		return nil, &client.APIError{Code: out.Code, Message: out.Message}
	}
	return out.Data, nil
}

// DoOne executes the request and returns the first element of the data array. It
// is used by endpoints whose data is a single object wrapped in a one-element
// array (account config, place-order ack, set-leverage, ...). Returns (nil, nil)
// when the success response carries an empty data array.
func DoOne[T any](r *Request) (*T, error) {
	list, err := DoList[T](r)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return &list[0], nil
}

// objResponse is the envelope shape for the handful of OKX endpoints whose
// "data" is a JSON OBJECT instead of the usual array (e.g.
// /rubik/stat/trading-data/support-coin).
type objResponse[T any] struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
	Data    T      `json:"data"`
}

// DoObject executes the request and decodes the envelope's data OBJECT into *T.
// It is used by the rare OKX endpoints that return "data" as a single JSON
// object rather than an array. A non-"0" code is returned as a *client.APIError.
func DoObject[T any](r *Request) (*T, error) {
	if r.err != nil {
		return nil, r.err
	}
	if err := r.prepare(); err != nil {
		return nil, err
	}

	r.client.GetLogger().Debugf("request: %s %s body=%s", r.method, r.r.URL, r.bodyJSON)
	response, err := r.r.Send()
	if err != nil {
		r.client.GetLogger().Errorf("request %s %s failed: %s", r.method, r.r.URL, err)
		return nil, err
	}
	body := response.Body()
	r.client.GetLogger().Debugf("response: %s", common.BytesToString(body))

	var out objResponse[T]
	if uerr := r.client.GetHttpClient().JSONUnmarshal(body, &out); uerr != nil {
		if apiErr := parseAPIError(r, body); apiErr != nil {
			return nil, apiErr
		}
		return nil, fmt.Errorf("request failed (status %d): %s", response.StatusCode(), common.BytesToString(body))
	}
	if out.Code != "0" {
		return nil, &client.APIError{Code: out.Code, Message: out.Message}
	}
	return &out.Data, nil
}

// DoListPartial is like DoList but tolerates OKX's batch-level result codes: "1"
// (all items failed) and "2" (partial success). For batch order endpoints the
// real per-item status lives in each data element's sCode/sMsg, so the data
// array is returned even when the top-level code is "1"/"2". Only a genuine
// request-level failure (auth, rate-limit, validation — codes outside 0/1/2, or
// any non-zero code with an empty data array) is returned as a *client.APIError.
func DoListPartial[T any](r *Request) ([]T, error) {
	out, err := do[T](r)
	if err != nil {
		return nil, err
	}
	switch out.Code {
	case "0":
		return out.Data, nil
	case "1", "2":
		if len(out.Data) == 0 {
			return nil, &client.APIError{Code: out.Code, Message: out.Message}
		}
		return out.Data, nil
	default:
		return nil, &client.APIError{Code: out.Code, Message: out.Message}
	}
}

// do is the shared transport+decode step for the typed Do* helpers.
func do[T any](r *Request) (*apiResponse[T], error) {
	if r.err != nil {
		return nil, r.err
	}
	if err := r.prepare(); err != nil {
		return nil, err
	}

	r.client.GetLogger().Debugf("request: %s %s body=%s", r.method, r.r.URL, r.bodyJSON)
	response, err := r.r.Send()
	if err != nil {
		r.client.GetLogger().Errorf("request %s %s failed: %s", r.method, r.r.URL, err)
		return nil, err
	}
	body := response.Body()
	r.client.GetLogger().Debugf("response: %s", common.BytesToString(body))

	var out apiResponse[T]
	if uerr := r.client.GetHttpClient().JSONUnmarshal(body, &out); uerr != nil {
		// The body was not a well-formed envelope (gateway error, HTML, ...).
		if apiErr := parseAPIError(r, body); apiErr != nil {
			return nil, apiErr
		}
		return nil, fmt.Errorf("request failed (status %d): %s", response.StatusCode(), common.BytesToString(body))
	}
	return &out, nil
}

// DoRawData executes the request and returns the raw JSON bytes of the
// envelope's "data" array (after verifying the code is 0/1/2). Tests use it to
// diff the real response shape against the typed structs.
func DoRawData(r *Request) ([]byte, error) {
	if r.err != nil {
		return nil, r.err
	}
	if err := r.prepare(); err != nil {
		return nil, err
	}
	response, err := r.r.Send()
	if err != nil {
		return nil, err
	}
	body := response.Body()
	var env struct {
		Code    string         `json:"code"`
		Message string         `json:"msg"`
		Data    jsontext.Value `json:"data"`
	}
	if uerr := r.client.GetHttpClient().JSONUnmarshal(body, &env); uerr != nil {
		return nil, fmt.Errorf("request failed (status %d): %s", response.StatusCode(), common.BytesToString(body))
	}
	switch env.Code {
	case "0", "1", "2":
		return env.Data, nil
	default:
		return nil, &client.APIError{Code: env.Code, Message: env.Message}
	}
}

// DoRaw executes the request and returns the raw, undecoded response body. Use
// it for the rare endpoints whose payload shape is non-uniform.
func DoRaw(r *Request) ([]byte, error) {
	if r.err != nil {
		return nil, r.err
	}
	if err := r.prepare(); err != nil {
		return nil, err
	}
	r.client.GetLogger().Debugf("request: %s %s body=%s", r.method, r.r.URL, r.bodyJSON)
	response, err := r.r.Send()
	if err != nil {
		return nil, err
	}
	body := response.Body()
	if apiErr := parseAPIError(r, body); apiErr != nil {
		return nil, apiErr
	}
	return body, nil
}

// prepare finalizes the URL, body and (when private) the OK-ACCESS-* signing
// headers. The signed prehash is timestamp + method + requestPath + body, using
// the exact bytes that go on the wire. The timestamp is OKX's ISO-8601
// millisecond UTC string.
func (r *Request) prepare() error {
	r.r.URL = r.fullURL()
	r.r.Method = r.method
	if r.bodyJSON != "" {
		r.r.SetHeader("Content-Type", "application/json")
		r.r.SetBody(r.bodyJSON)
	}
	if r.client.IsDemoTrading() {
		r.r.SetHeader("x-simulated-trading", "1")
	}
	if !r.needSign {
		return nil
	}

	apiKey := r.client.GetAPIKey()
	secret := r.client.GetAPISecret()
	passphrase := r.client.GetPassphrase()
	if apiKey == "" || secret == "" || passphrase == "" {
		return errors.New("missing credentials: configure client.WithAuth(apiKey, apiSecret, passphrase)")
	}

	ts := r.client.Now().UTC().Format(common.TimestampLayout)
	prehash := ts + r.method + r.requestPath() + r.bodyJSON

	var (
		sign string
		err  error
	)
	if fn := r.client.GetSignFn(); fn != nil {
		sign, err = fn(secret, prehash)
	} else {
		sign, err = HMACSign(secret, prehash)
	}
	if err != nil {
		return err
	}

	r.r.SetHeader("Content-Type", "application/json")
	r.r.SetHeader("OK-ACCESS-KEY", apiKey)
	r.r.SetHeader("OK-ACCESS-SIGN", sign)
	r.r.SetHeader("OK-ACCESS-TIMESTAMP", ts)
	r.r.SetHeader("OK-ACCESS-PASSPHRASE", passphrase)
	return nil
}

// parseAPIError tries to decode body as an OKX error envelope, returning nil
// when it is not an API-level error.
func parseAPIError(r *Request, body []byte) error {
	apiErr := &client.APIError{}
	if e := r.client.GetHttpClient().JSONUnmarshal(body, apiErr); e != nil {
		return nil
	}
	if !apiErr.IsValid() {
		return nil
	}
	return apiErr
}
