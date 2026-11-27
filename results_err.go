package okx

import "github.com/UnipayFI/go-okx/client"

// scodeErr folds a batch endpoint's per-item sCode/sMsg into a *client.APIError,
// or nil when the item succeeded (sCode "" or "0"). OKX batch endpoints report
// the real per-item outcome in sCode/sMsg even when the top-level request code is
// "0"/"1"/"2"; the Err methods below surface that as a standard error so it
// classifies like any other request error (errors.As, client.IsCode).
func scodeErr(sCode, sMsg string) error {
	if sCode == "" || sCode == "0" {
		return nil
	}
	return &client.APIError{Code: sCode, Message: sMsg}
}

// Err returns the per-item API error each batch result carries, or nil when that
// item succeeded.
func (r OrderResult) Err() error       { return scodeErr(r.SCode, r.SMsg) }
func (r AmendResult) Err() error       { return scodeErr(r.SCode, r.SMsg) }
func (r AlgoResult) Err() error        { return scodeErr(r.SCode, r.SMsg) }
func (r GridResult) Err() error        { return scodeErr(r.SCode, r.SMsg) }
func (r RfqCancelAck) Err() error      { return scodeErr(r.SCode, r.SMsg) }
func (r RfqQuoteCancelAck) Err() error { return scodeErr(r.SCode, r.SMsg) }
func (r SprdOrderAck) Err() error      { return scodeErr(r.SCode, r.SMsg) }
func (r SprdAmendAck) Err() error      { return scodeErr(r.SCode, r.SMsg) }
func (r SprdAlgoAck) Err() error       { return scodeErr(r.SCode, r.SMsg) }
func (r RecurringResult) Err() error   { return scodeErr(r.SCode, r.SMsg) }
