package domain

import (
	"fmt"
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
		{"too short - 1 char", "E"},
		{"too short - 2 chars", "EU"},
		{"too long - 4 chars", "EURO"},
		{"too long - 5 chars", "EUROS"},
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

func TestMoney_Add_SameCurrency(t *testing.T) {
	tests := []struct {
		name   string
		a, b   int
		expect int
	}{
		{"positive + positive", 1000, 299, 1299},
		{"positive + zero", 500, 0, 500},
		{"zero + zero", 0, 0, 0},
		{"positive + negative", 1000, -300, 700},
		{"negative + negative", -100, -200, -300},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, _ := NewMoney(tt.a, "EUR")
			b, _ := NewMoney(tt.b, "EUR")
			result := a.Add(b)
			if result.Amount() != tt.expect {
				t.Errorf("Add() amount = %d, want %d", result.Amount(), tt.expect)
			}
			if result.Currency() != "EUR" {
				t.Errorf("Add() currency = %q, want %q", result.Currency(), "EUR")
			}
		})
	}
}

func TestMoney_Add_DifferentCurrency_Panics(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic when adding different currencies")
		}
	}()
	a, _ := NewMoney(100, "EUR")
	b, _ := NewMoney(200, "USD")
	a.Add(b)
}

func TestMoney_Subtract_SameCurrency(t *testing.T) {
	tests := []struct {
		name   string
		a, b   int
		expect int
	}{
		{"larger minus smaller", 1000, 300, 700},
		{"equal amounts", 500, 500, 0},
		{"smaller minus larger", 100, 500, -400},
		{"subtract zero", 999, 0, 999},
		{"subtract from zero", 0, 100, -100},
		{"negative minus negative", -100, -300, 200},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, _ := NewMoney(tt.a, "USD")
			b, _ := NewMoney(tt.b, "USD")
			result := a.Subtract(b)
			if result.Amount() != tt.expect {
				t.Errorf("Subtract() amount = %d, want %d", result.Amount(), tt.expect)
			}
			if result.Currency() != "USD" {
				t.Errorf("Subtract() currency = %q, want %q", result.Currency(), "USD")
			}
		})
	}
}

func TestMoney_Subtract_DifferentCurrency_Panics(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic when subtracting different currencies")
		}
	}()
	a, _ := NewMoney(100, "EUR")
	b, _ := NewMoney(50, "GBP")
	a.Subtract(b)
}

func TestMoney_MultiplyPercent(t *testing.T) {
	tests := []struct {
		name    string
		amount  int
		percent int
		expect  int
	}{
		{"100% keeps same", 1000, 100, 1000},
		{"200% doubles", 1000, 200, 2000},
		{"50% halves", 1000, 50, 500},
		{"110% adds 10%", 1000, 110, 1100},
		{"0% yields zero", 1000, 0, 0},
		{"rounding truncates", 1099, 50, 549},
		{"negative amount", -1000, 150, -1500},
		{"small amount truncation", 1, 50, 0},
		{"33 percent of 100", 100, 33, 33},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, _ := NewMoney(tt.amount, "EUR")
			result := m.MultiplyPercent(tt.percent)
			if result.Amount() != tt.expect {
				t.Errorf("MultiplyPercent(%d) = %d, want %d", tt.percent, result.Amount(), tt.expect)
			}
			if result.Currency() != "EUR" {
				t.Errorf("MultiplyPercent() currency = %q, want %q", result.Currency(), "EUR")
			}
		})
	}
}

func TestMoney_Display(t *testing.T) {
	tests := []struct {
		name   string
		amount int
		curr   string
		expect string
	}{
		{"standard amount", 1299, "EUR", "12.99 EUR"},
		{"whole amount", 1000, "USD", "10.00 USD"},
		{"zero", 0, "GBP", "0.00 GBP"},
		{"single cent", 1, "EUR", "0.01 EUR"},
		{"only cents", 99, "USD", "0.99 USD"},
		{"large amount", 999999, "EUR", "9999.99 EUR"},
		{"negative amount", -1299, "EUR", "-12.99 EUR"},
		{"negative cents only", -5, "USD", "0.05 USD"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, _ := NewMoney(tt.amount, tt.curr)
			got := m.Display()
			if got != tt.expect {
				t.Errorf("Display() = %q, want %q", got, tt.expect)
			}
		})
	}
}

func TestMoney_IsZero(t *testing.T) {
	tests := []struct {
		name   string
		amount int
		want   bool
	}{
		{"zero is zero", 0, true},
		{"positive is not zero", 100, false},
		{"negative is not zero", -100, false},
		{"one cent is not zero", 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, _ := NewMoney(tt.amount, "EUR")
			if got := m.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMoney_ErrorMessages(t *testing.T) {
	t.Run("empty currency error message", func(t *testing.T) {
		_, err := NewMoney(100, "")
		if err == nil {
			t.Fatal("expected error")
		}
		expected := "currency must not be empty"
		if err.Error() != expected {
			t.Errorf("error = %q, want %q", err.Error(), expected)
		}
	})

	t.Run("invalid length currency error message", func(t *testing.T) {
		_, err := NewMoney(100, "AB")
		if err == nil {
			t.Fatal("expected error")
		}
		expected := fmt.Sprintf("currency must be a 3-letter ISO code, got %q", "AB")
		if err.Error() != expected {
			t.Errorf("error = %q, want %q", err.Error(), expected)
		}
	})
}

func TestMoney_Add_PreservesCurrency(t *testing.T) {
	a, _ := NewMoney(100, "JPY")
	b, _ := NewMoney(200, "JPY")
	result := a.Add(b)
	if result.Currency() != "JPY" {
		t.Errorf("expected currency JPY, got %q", result.Currency())
	}
}

func TestMoney_Subtract_PreservesCurrency(t *testing.T) {
	a, _ := NewMoney(500, "CHF")
	b, _ := NewMoney(200, "CHF")
	result := a.Subtract(b)
	if result.Currency() != "CHF" {
		t.Errorf("expected currency CHF, got %q", result.Currency())
	}
}
