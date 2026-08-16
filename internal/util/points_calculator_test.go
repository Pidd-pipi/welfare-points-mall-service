package util

import "testing"

func TestExchangeCost(t *testing.T) {
	tests := []struct {
		name       string
		pointsCost int
		quantity   int
		want       int
	}{
		{"single", 100, 1, 100},
		{"multiple", 100, 3, 300},
		{"zero qty", 100, 0, 0},
		{"negative qty", 100, -1, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExchangeCost(tt.pointsCost, tt.quantity); got != tt.want {
				t.Errorf("ExchangeCost(%d,%d) = %d, want %d", tt.pointsCost, tt.quantity, got, tt.want)
			}
		})
	}
}

func TestCanAfford(t *testing.T) {
	tests := []struct {
		name    string
		balance int
		cost    int
		want    bool
	}{
		{"enough", 500, 400, true},
		{"equal", 400, 400, true},
		{"not enough", 300, 400, false},
	}
	for _, tt := range tests {
		if got := CanAfford(tt.balance, tt.cost); got != tt.want {
			t.Errorf("CanAfford(%d,%d) = %v, want %v", tt.balance, tt.cost, got, tt.want)
		}
	}
}

func TestDeductPoints(t *testing.T) {
	tests := []struct {
		name    string
		balance int
		amount  int
		want    int
	}{
		{"normal", 500, 200, 300},
		{"insufficient", 100, 200, 100},
	}
	for _, tt := range tests {
		if got := DeductPoints(tt.balance, tt.amount); got != tt.want {
			t.Errorf("DeductPoints(%d,%d) = %d, want %d", tt.balance, tt.amount, got, tt.want)
		}
	}
}
