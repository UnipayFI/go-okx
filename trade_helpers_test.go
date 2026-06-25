package okx

import (
	"strings"
	"testing"

	"github.com/UnipayFI/go-okx/client"
)

func TestOrderResultErr(t *testing.T) {
	if err := (OrderResult{SCode: "0"}).Err(); err != nil {
		t.Errorf("SCode 0: got %v, want nil", err)
	}
	if err := (OrderResult{}).Err(); err != nil {
		t.Errorf("empty SCode: got %v, want nil", err)
	}
	err := (OrderResult{SCode: "51400", SMsg: "gone"}).Err()
	if err == nil {
		t.Fatal("SCode 51400: got nil, want error")
	}
	if !client.IsCode(err, "51400") {
		t.Errorf("IsCode(51400) = false on %v", err)
	}
}

func TestAmendResultErr(t *testing.T) {
	if err := (AmendResult{SCode: "0"}).Err(); err != nil {
		t.Errorf("SCode 0: got %v, want nil", err)
	}
	err := (AmendResult{SCode: "51000", SMsg: "bad param"}).Err()
	if !client.IsCode(err, "51000") {
		t.Errorf("IsCode(51000) = false on %v", err)
	}
}

func TestValidateClOrdID(t *testing.T) {
	valid := []string{"", "abc123", "T10421ZL3ZBUYZ42", strings.Repeat("a", MaxClOrdIDLen)}
	for _, v := range valid {
		if err := ValidateClOrdID(v); err != nil {
			t.Errorf("ValidateClOrdID(%q) = %v, want nil", v, err)
		}
	}
	invalid := []string{strings.Repeat("a", MaxClOrdIDLen+1), "has-hyphen", "has_underscore", "has space", "tag😀"}
	for _, v := range invalid {
		if err := ValidateClOrdID(v); err == nil {
			t.Errorf("ValidateClOrdID(%q) = nil, want error", v)
		}
	}
}

func TestMaxClOrdIDLen(t *testing.T) {
	if MaxClOrdIDLen != 32 {
		t.Errorf("MaxClOrdIDLen = %d, want 32", MaxClOrdIDLen)
	}
}

func TestMarketBarValues(t *testing.T) {
	cases := map[MarketBar]string{
		MarketBar1m:    "1m",
		MarketBar15m:   "15m",
		MarketBar1H:    "1H",
		MarketBar4H:    "4H",
		MarketBar1D:    "1D",
		MarketBar1W:    "1W",
		MarketBar1M:    "1M",
		MarketBar1Dutc: "1Dutc",
	}
	for bar, want := range cases {
		if string(bar) != want {
			t.Errorf("MarketBar = %q, want %q", string(bar), want)
		}
	}
}
