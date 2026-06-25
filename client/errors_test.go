package client

import (
	"errors"
	"fmt"
	"testing"
)

func TestAsAPIError(t *testing.T) {
	ptr := &APIError{Code: "51400", Message: "order does not exist"}
	tests := []struct {
		name     string
		err      error
		wantOK   bool
		wantCode string
	}{
		{"nil", nil, false, ""},
		{"direct pointer", ptr, true, "51400"},
		{"wrapped pointer", fmt.Errorf("place: %w", ptr), true, "51400"},
		// A caller may %w-wrap the value form (as the OKX adapter does); AsAPIError
		// must still find it via its value target.
		{"wrapped value", fmt.Errorf("place: %w", APIError{Code: "51008", Message: "insufficient"}), true, "51008"},
		{"double wrapped value", fmt.Errorf("outer: %w", fmt.Errorf("inner: %w", APIError{Code: "59000"})), true, "59000"},
		{"non api error", errors.New("boom"), false, ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := AsAPIError(tc.err)
			if ok != tc.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tc.wantOK)
			}
			if ok && got.Code != tc.wantCode {
				t.Fatalf("code = %q, want %q", got.Code, tc.wantCode)
			}
			if !ok && got != nil {
				t.Fatalf("got = %v, want nil when !ok", got)
			}
		})
	}
}

func TestIsCode(t *testing.T) {
	wrapped := fmt.Errorf("cancel: %w", &APIError{Code: "51401", Message: "already canceled"})
	if !IsCode(wrapped, "51401") {
		t.Errorf("IsCode(51401) = false, want true")
	}
	if IsCode(wrapped, "51400") {
		t.Errorf("IsCode(51400) = true, want false")
	}
	if IsCode(errors.New("boom"), "51400") {
		t.Errorf("IsCode(non-api) = true, want false")
	}
	if IsCode(nil, "0") {
		t.Errorf("IsCode(nil) = true, want false")
	}
}
