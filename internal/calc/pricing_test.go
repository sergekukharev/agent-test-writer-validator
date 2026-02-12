package calc

import "testing"

func TestBulkDiscount_NoDiscount(t *testing.T) {
	discount := BulkDiscount(5, StandardTiers)
	if discount != 0 {
		t.Errorf("expected 0%% discount for 5 items, got %d%%", discount)
	}
}

func TestBulkDiscount_FirstTier(t *testing.T) {
	discount := BulkDiscount(10, StandardTiers)
	if discount != 5 {
		t.Errorf("expected 5%% discount for 10 items, got %d%%", discount)
	}
}

func TestBulkDiscount_HighestTier(t *testing.T) {
	discount := BulkDiscount(100, StandardTiers)
	if discount != 20 {
		t.Errorf("expected 20%% discount for 100 items, got %d%%", discount)
	}
}
