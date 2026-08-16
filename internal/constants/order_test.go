package constants

import "testing"

func TestCanOrderTransition(t *testing.T) {
	tests := []struct {
		name string
		from OrderStatus
		to   OrderStatus
		want bool
	}{
		{"pending to shipped", OrderPending, OrderShipped, true},
		{"pending to cancelled", OrderPending, OrderCancelled, true},
		{"shipped to completed", OrderShipped, OrderCompleted, true},
		{"completed to cancelled", OrderCompleted, OrderCancelled, false},
		{"cancelled to shipped", OrderCancelled, OrderShipped, false},
	}
	for _, tt := range tests {
		if got := CanOrderTransition(tt.from, tt.to); got != tt.want {
			t.Errorf("CanOrderTransition(%s,%s) = %v, want %v", tt.from, tt.to, got, tt.want)
		}
	}
}
