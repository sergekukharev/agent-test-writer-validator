package domain

import (
	"testing"
)

func TestNewMoney_Valid(t *testing.T) {
	tests := []struct {
		name     string
		amount   int
		currency string
	}{
		{"positive amount", 1299, "EUR"},
		{"zero amount", 0, "USD"},
		{"negative amount", -500, "GBP"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := NewMoney(tt.amount, tt.currency)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if m.Amount() != tt.amount {
				t.Errorf("Amount() = %d, want %d", m.Amount(), tt.amount)
			}
			if m.Currency() != tt.currency {
				t.Errorf("Currency() = %q, want %q", m.Currency(), tt.currency)
			}
		})
	}
}

func TestNewMoney_EmptyCurrency(t *testing.T) {
	_, err := NewMoney(100, "")
	if err == nil {
		t.Fatal("expected error for empty currency")
	}
}

func TestNewMoney_InvalidCurrencyLength(t *testing.T) {
	tests := []struct {
		name     string
		currency string
	}{
		{"too short — 1 char", "E"},
		{"too short — 2 chars", "EU"},
		{"too long — 4 chars", "EURO"},
		{"too long — 5 chars", "EUROS"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewMoney(100, tt.currency)
			if err == nil {
				t.Fatalf("expected error for currency %q", tt.currency)
			}
		})
	}
}

func TestAdd_SameCurrency(t *testing.T) {
	a, _ := NewMoney(1000, "EUR")
	b, _ := NewMoney(299, "EUR")

	result := a.Add(b)

	if result.Amount() != 1299 {
		t.Errorf("Amount() = %d, want 1299", result.Amount())
	}
	if result.Currency() != "EUR" {
		t.Errorf("Currency() = %q, want %q", result.Currency(), "EUR")
	}
}

func TestAdd_ZeroAmount(t *testing.T) {
	a, _ := NewMoney(500, "USD")
	b, _ := NewMoney(0, "USD")

	result := a.Add(b)

	if result.Amount() != 500 {
		t.Errorf("Amount() = %d, want 500", result.Amount())
	}
}

func TestAdd_NegativeAmounts(t *testing.T) {
	a, _ := NewMoney(-300, "USD")
	b, _ := NewMoney(-200, "USD")

	result := a.Add(b)

	if result.Amount() != -500 {
		t.Errorf("Amount() = %d, want -500", result.Amount())
	}
}

func TestAdd_DifferentCurrencyPanics(t *testing.T) {
	a, _ := NewMoney(100, "EUR")
	b, _ := NewMoney(200, "USD")

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for mismatched currencies")
		}
	}()

	a.Add(b)
}

func TestSubtract_SameCurrency(t *testing.T) {
	a, _ := NewMoney(1000, "EUR")
	b, _ := NewMoney(299, "EUR")

	result := a.Subtract(b)

	if result.Amount() != 701 {
		t.Errorf("Amount() = %d, want 701", result.Amount())
	}
	if result.Currency() != "EUR" {
		t.Errorf("Currency() = %q, want %q", result.Currency(), "EUR")
	}
}

func TestSubtract_ResultsInNegative(t *testing.T) {
	a, _ := NewMoney(100, "GBP")
	b, _ := NewMoney(500, "GBP")

	result := a.Subtract(b)

	if result.Amount() != -400 {
		t.Errorf("Amount() = %d, want -400", result.Amount())
	}
}

func TestSubtract_ZeroAmount(t *testing.T) {
	a, _ := NewMoney(500, "USD")
	b, _ := NewMoney(0, "USD")

	result := a.Subtract(b)

	if result.Amount() != 500 {
		t.Errorf("Amount() = %d, want 500", result.Amount())
	}
}

func TestSubtract_DifferentCurrencyPanics(t *testing.T) {
	a, _ := NewMoney(100, "EUR")
	b, _ := NewMoney(50, "USD")

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for mismatched currencies")
		}
	}()

	a.Subtract(b)
}

func TestMultiplyPercent(t *testing.T) {
	tests := []struct {
		name    string
		amount  int
		percent int
		want    int
	}{
		{"100% keeps amount", 1000, 100, 1000},
		{"200% doubles", 1000, 200, 2000},
		{"50% halves", 1000, 50, 500},
		{"110% adds 10%", 1000, 110, 1100},
		{"0% zeroes out", 1000, 0, 0},
		{"small amount truncates", 1, 150, 1},
		{"negative percent", 1000, -100, -1000},
		{"negative amount with percent", -1000, 110, -1100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, _ := NewMoney(tt.amount, "EUR")
			result := m.MultiplyPercent(tt.percent)
			if result.Amount() != tt.want {
				t.Errorf("Amount() = %d, want %d", result.Amount(), tt.want)
			}
			if result.Currency() != "EUR" {
				t.Errorf("Currency() = %q, want %q", result.Currency(), "EUR")
			}
		})
	}
}

func TestDisplay(t *testing.T) {
	tests := []struct {
		name   string
		amount int
		want   string
	}{
		{"positive amount", 1299, "12.99 EUR"},
		{"zero amount", 0, "0.00 EUR"},
		{"whole amount no cents", 500, "5.00 EUR"},
		{"single digit cents", 1205, "12.05 EUR"},
		{"large amount", 100000, "1000.00 EUR"},
		{"one cent", 1, "0.01 EUR"},
		{"negative amount", -1299, "-12.99 EUR"},
		{"negative small amount", -5, "0.05 EUR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, _ := NewMoney(tt.amount, "EUR")
			got := m.Display()
			if got != tt.want {
				t.Errorf("Display() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIsZero(t *testing.T) {
	tests := []struct {
		name   string
		amount int
		want   bool
	}{
		{"zero is zero", 0, true},
		{"positive is not zero", 100, false},
		{"negative is not zero", -100, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, _ := NewMoney(tt.amount, "EUR")
			if m.IsZero() != tt.want {
				t.Errorf("IsZero() = %v, want %v", m.IsZero(), tt.want)
			}
		})
	}
}

func TestAdd_IsImmutable(t *testing.T) {
	a, _ := NewMoney(100, "EUR")
	b, _ := NewMoney(200, "EUR")

	result := a.Add(b)

	if a.Amount() != 100 {
		t.Errorf("original should be unchanged, got %d", a.Amount())
	}
	if result.Amount() != 300 {
		t.Errorf("result Amount() = %d, want 300", result.Amount())
	}
}

func TestSubtract_IsImmutable(t *testing.T) {
	a, _ := NewMoney(500, "EUR")
	b, _ := NewMoney(200, "EUR")

	result := a.Subtract(b)

	if a.Amount() != 500 {
		t.Errorf("original should be unchanged, got %d", a.Amount())
	}
	if result.Amount() != 300 {
		t.Errorf("result Amount() = %d, want 300", result.Amount())
	}
}

func TestMultiplyPercent_IsImmutable(t *testing.T) {
	m, _ := NewMoney(1000, "EUR")

	result := m.MultiplyPercent(200)

	if m.Amount() != 1000 {
		t.Errorf("original should be unchanged, got %d", m.Amount())
	}
	if result.Amount() != 2000 {
		t.Errorf("result Amount() = %d, want 2000", result.Amount())
	}
}
