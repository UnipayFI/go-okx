package okx

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

// TestOrderLifecycle exercises the full single-order write path against the live
// account with a tiny, far-below-market SPOT post_only limit buy that rests
// unfilled and is then cancelled — fully reversible, costing only the briefly
// reserved ~3 USDT. It also verifies the Order struct against a real order
// (assertCovers), which the read-only TestTradeOrder cannot do on an account
// with no order history. Gated behind OKX_TEST_WRITE=1.
//
// Flow: ticker (compute a safe far-below px) -> place (post_only) -> get-order
// (assertCovers) -> orders-pending (find ours) -> amend px -> get-order (verify)
// -> cancel -> orders-history (find ours).
func TestOrderLifecycle(t *testing.T) {
	if os.Getenv("OKX_TEST_WRITE") == "" {
		t.Skip("set OKX_TEST_WRITE=1 to run live order tests")
	}
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := ctx(t)

	const instId = "BTC-USDT"

	// A far-below-market limit price so the post_only buy rests as maker and
	// never fills. Derive it from the live last price, truncated to an integer
	// (BTC-USDT tickSz is 0.1, so an integer price is always valid).
	tk, err := c.NewGetTickerService(instId).Do(cx)
	if err != nil {
		t.Fatalf("ticker: %v", err)
	}
	px := tk.Last.Mul(decimal.RequireFromString("0.5")).Truncate(0)
	qty := decimal.RequireFromString("0.0001") // ~min size; ~px*qty USDT reserved
	clOrdId := "gookx" + strconv.FormatInt(time.Now().UnixNano(), 10)
	if len(clOrdId) > 32 {
		clOrdId = clOrdId[:32]
	}

	// 1) Place a resting post_only limit buy.
	ref, err := c.NewPlaceOrderService(instId, TdModeCash, SideBuy, OrdTypePostOnly, qty).
		SetPx(px).
		SetClOrdId(clOrdId).
		Do(cx)
	if err != nil {
		t.Fatalf("place-order: %v", err)
	}
	if ref.SCode != "0" {
		t.Fatalf("place-order rejected: sCode=%s sMsg=%s", ref.SCode, ref.SMsg)
	}
	t.Logf("placed ordId=%s clOrdId=%s px=%s", ref.OrderID, ref.ClientOrderID, px)
	if ref.OrderID == "" {
		t.Fatal("place-order returned empty ordId")
	}

	// Always clean up, even if a later step fails.
	defer func() {
		_, _ = c.NewCancelOrderService(instId).SetOrdId(ref.OrderID).Do(ctx(t))
	}()

	// 2) get-order — validate the Order struct against the real order.
	order, err := c.NewGetOrderService(instId).SetOrdId(ref.OrderID).Do(cx)
	if err != nil {
		t.Fatalf("get-order: %v", err)
	}
	t.Logf("order: state=%s px=%s sz=%s side=%s ordType=%s", order.State, order.Price, order.Size, order.Side, order.OrderType)
	raw := fetchRawGet(t, c, cx, "/api/v5/trade/order", map[string]string{"instId": instId, "ordId": ref.OrderID}, true)
	assertCovers(t, "trade/order", raw, order)

	// 3) orders-pending should include it.
	open, err := c.NewGetOrdersPendingService().SetInstId(instId).Do(cx)
	if err != nil {
		t.Fatalf("orders-pending: %v", err)
	}
	found := false
	for _, o := range open {
		if o.OrderID == ref.OrderID {
			found = true
		}
	}
	t.Logf("orders-pending: %d (found ours=%v)", len(open), found)
	if !found {
		t.Errorf("placed order %s not present in orders-pending", ref.OrderID)
	}

	// 4) amend the price (still far below market).
	newPx := px.Sub(decimal.RequireFromString("100"))
	if _, err := c.NewAmendOrderService(instId).SetOrdId(ref.OrderID).SetNewPx(newPx).Do(cx); err != nil {
		t.Fatalf("amend-order: %v", err)
	}
	t.Logf("amended px -> %s", newPx)

	// 5) get-order reflects the new price.
	time.Sleep(500 * time.Millisecond)
	order2, err := c.NewGetOrderService(instId).SetOrdId(ref.OrderID).Do(cx)
	if err != nil {
		t.Fatalf("get-order(2): %v", err)
	}
	t.Logf("after amend: px=%s state=%s", order2.Price, order2.State)
	if !order2.Price.Equal(newPx) {
		t.Errorf("amend did not take effect: px=%s want=%s", order2.Price, newPx)
	}

	// 6) cancel.
	if _, err := c.NewCancelOrderService(instId).SetOrdId(ref.OrderID).Do(cx); err != nil {
		t.Fatalf("cancel-order: %v", err)
	}
	t.Logf("cancelled %s", ref.OrderID)

	// 7) orders-history should now include the cancelled order.
	time.Sleep(500 * time.Millisecond)
	hist, err := c.NewGetOrdersHistoryService(InstTypeSpot).SetInstId(instId).Do(cx)
	if err != nil {
		t.Fatalf("orders-history: %v", err)
	}
	inHist := false
	for _, o := range hist {
		if o.OrderID == ref.OrderID {
			inHist = true
		}
	}
	t.Logf("orders-history: %d (found ours=%v)", len(hist), inHist)
}
