package domain

import (
	"testing"
)

func TestNewMoney(t *testing.T) {
	tests := []struct {
		name        string
		amount      int
		currency    string
		wantAmount  int
		wantCurr    string
		wantErr     bool
		errContains string
	}{
		{
			name:       "valid USD",
			amount:     1299,
			currency:   "USD",
			wantAmount: 1299,
			wantCurr:   "USD",
		},
		{
			name:       "valid EUR",
			amount:     0,
			currency:   "EUR",
			wantAmount: 0,
			wantCurr:   "EUR",
		},
		{
			name:       "negative amount is allowed",
			amount:     -500,
			currency:   "GBP",
			wantAmount: -500,
			wantCurr:   "GBP",
		},
		{
			name:        "empty currency",
			amount:      100,
			currency:    "",
			wantErr:     true,
			errContains: "currency must not be empty",
		},
		{
			name:        "currency too short",
			amount:      100,
			currency:    "US",
			wantErr:     true,
			errContains: "3-letter ISO code",
		},
		{
			name:        "currency too long",
			amount:      100,
			currency:    "USDX",
			wantErr:     true,
			errContains: "3-letter ISO code",
		},
		{
			name:        "single char currency",
			amount:      100,
			currency:    "U",
			wantErr:     true,
			errContains: "3-letter ISO code",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := NewMoney(tt.amount, tt.currency)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errContains != "" && !containsStr(err.Error(), tt.errContains) {
					t.Errorf("error %q should contain %q", err.Error(), tt.errContains)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if m.Amount() != tt.wantAmount {
				t.Errorf("Amount() = %d, want %d", m.Amount(), tt.wantAmount)
			}
			if m.Currency() != tt.wantCurr {
				t.Errorf("Currency() = %q, want %q", m.Currency(), tt.wantCurr)
			}
		})
	}
}

func TestMoney_Add(t *testing.T) {
	tests := []struct {
		name       string
		a          Money
		b          Money
		wantAmount int
	}{
		{
			name:       "positive plus positive",
			a:          mustMoney(1000, "USD"),
			b:          mustMoney(299, "USD"),
			wantAmount: 1299,
		},
		{
			name:       "positive plus zero",
			a:          mustMoney(500, "EUR"),
			b:          mustMoney(0, "EUR"),
			wantAmount: 500,
		},
		{
			name:       "zero plus zero",
			a:          mustMoney(0, "GBP"),
			b:          mustMoney(0, "GBP"),
			wantAmount: 0,
		},
		{
			name:       "positive plus negative",
			a:          mustMoney(1000, "USD"),
			b:          mustMoney(-300, "USD"),
			wantAmount: 700,
		},
		{
			name:       "negative plus negative",
			a:          mustMoney(-100, "USD"),
			b:          mustMoney(-200, "USD"),
			wantAmount: -300,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.a.Add(tt.b)
			if result.Amount() != tt.wantAmount {
				t.Errorf("Amount() = %d, want %d", result.Amount(), tt.wantAmount)
			}
			if result.Currency() != tt.a.Currency() {
				t.Errorf("Currency() = %q, want %q", result.Currency(), tt.a.Currency())
			}
		})
	}
}

func TestMoney_Add_PanicsOnCurrencyMismatch(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic on currency mismatch, got none")
		}
		msg, ok := r.(string)
		if !ok {
			t.Fatalf("expected string panic, got %T", r)
		}
		if !containsStr(msg, "cannot add") {
			t.Errorf("panic message %q should contain 'cannot add'", msg)
		}
	}()

	a := mustMoney(100, "USD")
	b := mustMoney(200, "EUR")
	a.Add(b)
}

func TestMoney_Subtract(t *testing.T) {
	tests := []struct {
		name       string
		a          Money
		b          Money
		wantAmount int
	}{
		{
			name:       "positive minus smaller positive",
			a:          mustMoney(1000, "USD"),
			b:          mustMoney(300, "USD"),
			wantAmount: 700,
		},
		{
			name:       "positive minus equal",
			a:          mustMoney(500, "EUR"),
			b:          mustMoney(500, "EUR"),
			wantAmount: 0,
		},
		{
			name:       "positive minus larger positive results in negative",
			a:          mustMoney(100, "GBP"),
			b:          mustMoney(300, "GBP"),
			wantAmount: -200,
		},
		{
			name:       "subtract zero",
			a:          mustMoney(100, "USD"),
			b:          mustMoney(0, "USD"),
			wantAmount: 100,
		},
		{
			name:       "subtract negative (effectively adds)",
			a:          mustMoney(100, "USD"),
			b:          mustMoney(-50, "USD"),
			wantAmount: 150,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.a.Subtract(tt.b)
			if result.Amount() != tt.wantAmount {
				t.Errorf("Amount() = %d, want %d", result.Amount(), tt.wantAmount)
			}
			if result.Currency() != tt.a.Currency() {
				t.Errorf("Currency() = %q, want %q", result.Currency(), tt.a.Currency())
			}
		})
	}
}

func TestMoney_Subtract_PanicsOnCurrencyMismatch(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic on currency mismatch, got none")
		}
		msg, ok := r.(string)
		if !ok {
			t.Fatalf("expected string panic, got %T", r)
		}
		if !containsStr(msg, "cannot subtract") {
			t.Errorf("panic message %q should contain 'cannot subtract'", msg)
		}
	}()

	a := mustMoney(100, "USD")
	b := mustMoney(200, "EUR")
	a.Subtract(b)
}

func TestMoney_MultiplyPercent(t *testing.T) {
	tests := []struct {
		name       string
		amount     int
		currency   string
		percent    int
		wantAmount int
	}{
		{
			name:       "100 percent leaves unchanged",
			amount:     1000,
			currency:   "USD",
			percent:    100,
			wantAmount: 1000,
		},
		{
			name:       "200 percent doubles",
			amount:     1000,
			currency:   "USD",
			percent:    200,
			wantAmount: 2000,
		},
		{
			name:       "50 percent halves",
			amount:     1000,
			currency:   "EUR",
			percent:    50,
			wantAmount: 500,
		},
		{
			name:       "110 percent adds 10%",
			amount:     1000,
			currency:   "USD",
			percent:    110,
			wantAmount: 1100,
		},
		{
			name:       "0 percent yields zero",
			amount:     1000,
			currency:   "USD",
			percent:    0,
			wantAmount: 0,
		},
		{
			name:       "negative percent negates",
			amount:     1000,
			currency:   "USD",
			percent:    -100,
			wantAmount: -1000,
		},
		{
			name:       "integer truncation",
			amount:     1,
			currency:   "USD",
			percent:    50,
			wantAmount: 0,
		},
		{
			name:       "integer truncation on 33 percent",
			amount:     100,
			currency:   "USD",
			percent:    33,
			wantAmount: 33,
		},
		{
			name:       "preserves currency",
			amount:     500,
			currency:   "GBP",
			percent:    100,
			wantAmount: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mustMoney(tt.amount, tt.currency)
			result := m.MultiplyPercent(tt.percent)
			if result.Amount() != tt.wantAmount {
				t.Errorf("Amount() = %d, want %d", result.Amount(), tt.wantAmount)
			}
			if result.Currency() != tt.currency {
				t.Errorf("Currency() = %q, want %q", result.Currency(), tt.currency)
			}
		})
	}
}

func TestMoney_Display(t *testing.T) {
	tests := []struct {
		name     string
		amount   int
		currency string
		want     string
	}{
		{
			name:     "standard amount",
			amount:   1299,
			currency: "EUR",
			want:     "12.99 EUR",
		},
		{
			name:     "zero",
			amount:   0,
			currency: "USD",
			want:     "0.00 USD",
		},
		{
			name:     "whole dollar no cents",
			amount:   500,
			currency: "USD",
			want:     "5.00 USD",
		},
		{
			name:     "only cents",
			amount:   42,
			currency: "GBP",
			want:     "0.42 GBP",
		},
		{
			name:     "single cent",
			amount:   1,
			currency: "USD",
			want:     "0.01 USD",
		},
		{
			name:     "large amount",
			amount:   1234567,
			currency: "EUR",
			want:     "12345.67 EUR",
		},
		{
			name:     "negative amount",
			amount:   -1299,
			currency: "USD",
			want:     "-12.99 USD",
		},
		{
			name:     "negative cents only",
			amount:   -5,
			currency: "USD",
			want:     "0.05 USD",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mustMoney(tt.amount, tt.currency)
			got := m.Display()
			if got != tt.want {
				t.Errorf("Display() = %q, want %q", got, tt.want)
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
		{name: "zero is zero", amount: 0, want: true},
		{name: "positive is not zero", amount: 1, want: false},
		{name: "negative is not zero", amount: -1, want: false},
		{name: "large positive is not zero", amount: 100000, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mustMoney(tt.amount, "USD")
			if got := m.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoney_Immutability(t *testing.T) {
	a := mustMoney(1000, "USD")
	b := mustMoney(500, "USD")

	_ = a.Add(b)
	if a.Amount() != 1000 {
		t.Error("Add modified the receiver")
	}

	_ = a.Subtract(b)
	if a.Amount() != 1000 {
		t.Error("Subtract modified the receiver")
	}

	_ = a.MultiplyPercent(200)
	if a.Amount() != 1000 {
		t.Error("MultiplyPercent modified the receiver")
	}
}

// mustMoney is a test helper that creates Money or panics on error.
func mustMoney(amount int, currency string) Money {
	m, err := NewMoney(amount, currency)
	if err != nil {
		panic(err)
	}
	return m
}

// containsStr checks if s contains substr (avoids importing strings package).
func containsStr(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
