package domain

import (
	"errors"
	"fmt"
)

// Money represents an amount in the smallest currency unit (e.g. cents).
type Money struct {
	amount   int
	currency string
}

func NewMoney(amount int, currency string) (Money, error) {
	if currency == "" {
		return Money{}, errors.New("currency must not be empty")
	}
	if len(currency) != 3 {
		return Money{}, fmt.Errorf("currency must be a 3-letter ISO code, got %q", currency)
	}
	return Money{amount: amount, currency: currency}, nil
}

func (m Money) Amount() int      { return m.amount }
func (m Money) Currency() string { return m.currency }

// Add returns a new Money with the sum of both amounts.
// Panics if currencies do not match.
func (m Money) Add(other Money) Money {
	if m.currency != other.currency {
		panic(fmt.Sprintf("cannot add %s to %s", m.currency, other.currency))
	}
	return Money{amount: m.amount + other.amount, currency: m.currency}
}

// Subtract returns a new Money with the difference.
// Panics if currencies do not match.
func (m Money) Subtract(other Money) Money {
	if m.currency != other.currency {
		panic(fmt.Sprintf("cannot subtract %s from %s", other.currency, m.currency))
	}
	return Money{amount: m.amount - other.amount, currency: m.currency}
}

// MultiplyScalar returns a new Money scaled by the given factor.
// Uses integer arithmetic — the factor is expressed as a percentage (e.g. 110 = 1.10x).
func (m Money) MultiplyPercent(percent int) Money {
	return Money{
		amount:   m.amount * percent / 100,
		currency: m.currency,
	}
}

// Display formats the money for human display, e.g. "12.99 EUR".
func (m Money) Display() string {
	whole := m.amount / 100
	cents := m.amount % 100
	if cents < 0 {
		cents = -cents
	}
	return fmt.Sprintf("%d.%02d %s", whole, cents, m.currency)
}

// IsZero returns true if the amount is zero.
func (m Money) IsZero() bool {
	return m.amount == 0
}
