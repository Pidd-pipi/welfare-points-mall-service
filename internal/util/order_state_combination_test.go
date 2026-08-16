package util

import (
	"testing"

	"github.com/ld/welfaremall/internal/constants"
)

func TestOrderStateMachineAllowsCancel(t *testing.T) {
	if !constants.CanOrderTransition(constants.OrderPending, constants.OrderCancelled) {
		t.Fatal("pending -> cancelled should be allowed")
	}
	if !constants.CanOrderTransition(constants.OrderShipped, constants.OrderCancelled) {
		t.Fatal("shipped -> cancelled should be allowed")
	}
}

func TestOrderStatusTextBoundary(t *testing.T) {
	if got := OrderStatusText(constants.OrderPending); got != "待发货" {
		t.Fatalf("OrderStatusText(pending) = %s, want 待发货", got)
	}
	if got := OrderStatusTag(constants.OrderCancelled); got != "danger" {
		t.Fatalf("OrderStatusTag(cancelled) = %s, want danger", got)
	}
}
