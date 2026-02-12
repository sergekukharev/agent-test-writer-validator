package calc

import "github.com/sergekukharev/agent-test-writer-validator/internal/domain"

// DiscountTier defines a discount based on quantity purchased.
type DiscountTier struct {
	MinQuantity int
	Percent     int // discount percentage, e.g. 10 = 10% off
}

// StandardTiers are the default bulk discount tiers.
var StandardTiers = []DiscountTier{
	{MinQuantity: 10, Percent: 5},
	{MinQuantity: 25, Percent: 10},
	{MinQuantity: 50, Percent: 15},
	{MinQuantity: 100, Percent: 20},
}

// BulkDiscount calculates the discount percentage for a given quantity.
// Uses the highest matching tier.
func BulkDiscount(quantity int, tiers []DiscountTier) int {
	best := 0
	for _, t := range tiers {
		if quantity >= t.MinQuantity && t.Percent > best {
			best = t.Percent
		}
	}
	return best
}

// OrderTotal calculates the total price for ordering n copies of a book,
// applying the best matching bulk discount.
func OrderTotal(book domain.Book, quantity int, tiers []DiscountTier) domain.Money {
	unitPrice := book.Price()
	discount := BulkDiscount(quantity, tiers)
	effectivePercent := 100 - discount

	lineTotal, _ := domain.NewMoney(
		unitPrice.Amount()*quantity*effectivePercent/100,
		unitPrice.Currency(),
	)
	return lineTotal
}

// ClassicSurcharge adds a 25% surcharge if the book is a classic (published > 50 years ago).
// Classics are considered collector items.
func ClassicSurcharge(book domain.Book) domain.Money {
	if !book.IsClassic() {
		return book.Price()
	}
	return book.Price().MultiplyPercent(125)
}

// NewReleasePremium adds a 10% premium for books published within the last year.
func NewReleasePremium(book domain.Book) domain.Money {
	if !book.IsRecent() {
		return book.Price()
	}
	return book.Price().MultiplyPercent(110)
}
