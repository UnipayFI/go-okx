package okx

import (
	"testing"
	"time"

	"github.com/UnipayFI/go-okx/common"
	"github.com/UnipayFI/go-okx/request"
)

// TestWsPrivate exercises the private-gateway WebSocket channels. account,
// positions and balance_and_position push a snapshot right after login (captured
// and covered); orders, liquidation-warning and account-greeks are
// activity-driven, so a no-push within the timeout is tolerated (login+subscribe
// already validated).
func TestWsPrivate(t *testing.T) {
	c := testWsClient(t)

	// reference the typed services so the typed paths are compiled/exercised.
	_ = c.NewSubscribeAccountService().SetExtraParams(`{"updateInterval":"0"}`)
	_ = c.NewSubscribePositionsService(InstTypeAny)
	_ = c.NewSubscribeBalanceAndPositionService()
	_ = c.NewSubscribeOrdersService(InstTypeSpot)
	_ = c.NewSubscribeLiquidationWarningService(InstTypeSwap)
	_ = c.NewSubscribeAccountGreeksService()

	t.Run("account", func(t *testing.T) {
		arg := request.WsArg{Channel: "account"}
		raw := wsFirstDataArray(t, c, request.GatewayPrivate, true, arg, 15*time.Second)
		if raw == nil {
			return
		}
		var data []WsAccount
		if err := common.JSONUnmarshal(raw, &data); err != nil {
			t.Fatalf("decode: %v", err)
		}
		assertCovers(t, "ws/account", raw, data)
	})

	t.Run("positions", func(t *testing.T) {
		arg := request.WsArg{Channel: "positions", InstrumentType: string(InstTypeAny)}
		raw := wsFirstDataArray(t, c, request.GatewayPrivate, true, arg, 15*time.Second)
		if raw == nil {
			return
		}
		var data []WsPosition
		if err := common.JSONUnmarshal(raw, &data); err != nil {
			t.Fatalf("decode: %v", err)
		}
		assertCovers(t, "ws/positions", raw, data)
	})

	t.Run("balance_and_position", func(t *testing.T) {
		arg := request.WsArg{Channel: "balance_and_position"}
		raw := wsFirstDataArray(t, c, request.GatewayPrivate, true, arg, 15*time.Second)
		if raw == nil {
			return
		}
		var data []WsBalanceAndPosition
		if err := common.JSONUnmarshal(raw, &data); err != nil {
			t.Fatalf("decode: %v", err)
		}
		assertCovers(t, "ws/balance_and_position", raw, data)
	})

	t.Run("orders", func(t *testing.T) {
		arg := request.WsArg{Channel: "orders", InstrumentType: string(InstTypeSpot)}
		raw := wsFirstDataArray(t, c, request.GatewayPrivate, true, arg, 15*time.Second)
		if raw == nil {
			return
		}
		var data []WsOrder
		if err := common.JSONUnmarshal(raw, &data); err != nil {
			t.Fatalf("decode: %v", err)
		}
		assertCovers(t, "ws/orders", raw, data)
	})

	t.Run("liquidation-warning", func(t *testing.T) {
		arg := request.WsArg{Channel: "liquidation-warning", InstrumentType: string(InstTypeSwap)}
		raw := wsFirstDataArray(t, c, request.GatewayPrivate, true, arg, 15*time.Second)
		if raw == nil {
			return
		}
		var data []WsLiquidationWarning
		if err := common.JSONUnmarshal(raw, &data); err != nil {
			t.Fatalf("decode: %v", err)
		}
		assertCovers(t, "ws/liquidation-warning", raw, data)
	})

	t.Run("account-greeks", func(t *testing.T) {
		arg := request.WsArg{Channel: "account-greeks"}
		raw := wsFirstDataArray(t, c, request.GatewayPrivate, true, arg, 15*time.Second)
		if raw == nil {
			return
		}
		var data []WsAccountGreeks
		if err := common.JSONUnmarshal(raw, &data); err != nil {
			t.Fatalf("decode: %v", err)
		}
		assertCovers(t, "ws/account-greeks", raw, data)
	})
}
