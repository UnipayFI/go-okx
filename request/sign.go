package request

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"

	"github.com/UnipayFI/go-okx/common"
)

// SignFn mirrors client.SignFn: it turns the prehash string into the
// OK-ACCESS-SIGN header value, given the configured secret.
type SignFn = func(secret, prehash string) (signature string, err error)

// HMACSign is OKX's default request signer:
//
//	OK-ACCESS-SIGN = base64( HMAC-SHA256( secretKey, prehash ) )
//
// where prehash = timestamp + method + requestPath + body (see Request.prepare).
// The timestamp is the ISO-8601 millisecond UTC string also sent in the
// OK-ACCESS-TIMESTAMP header.
func HMACSign(secret, prehash string) (string, error) {
	mac := hmac.New(sha256.New, common.StringToBytes(secret))
	if _, err := mac.Write(common.StringToBytes(prehash)); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}
