package okx

import (
	"testing"
	"time"

	"github.com/UnipayFI/go-okx/common"
)

// TestRowTimestampsUTC checks that the hand-written parsers for OKX's
// array-of-arrays payloads (candles, rubik rows) produce the same instant and
// the same UTC location as the JSON codec gives a format:unixmilli field.
func TestRowTimestampsUTC(t *testing.T) {
	const ms = "1790561220000"
	var want struct {
		T time.Time `json:"t,format:unixmilli"`
	}
	if err := common.JSONUnmarshal([]byte(`{"t":"`+ms+`"}`), &want); err != nil {
		t.Fatal(err)
	}
	row := []string{ms, "1", "2", "0.5", "1.5", "10", "20", "30", "1"}
	rows := [][]string{row}
	got := map[string]time.Time{
		"parseCandles":        parseCandles(rows)[0].Timestamp,
		"parseIndexCandles":   parseIndexCandles(rows)[0].Timestamp,
		"parseSprdCandles":    parseSprdCandles(rows)[0].Timestamp,
		"parseWsCandles":      parseWsCandles(rows)[0].Timestamp,
		"parseWsIndexCandles": parseWsIndexCandles(rows)[0].Timestamp,
		"parseWsSprdCandles":  parseWsSprdCandles(rows)[0].Timestamp,
		"parseRubikTs":        parseRubikTs(ms),
	}
	for name, ts := range got {
		if ts != want.T {
			t.Errorf("%s = %v (%v), want %v (%v)", name, ts, ts.Location(), want.T, want.T.Location())
		}
	}
}
