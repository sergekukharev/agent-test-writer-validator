package calc

import (
	"sort"

	"github.com/sergekukharev/agent-test-writer-validator/internal/domain"
)

// AveragePrice calculates the average price of a slice of books (in cents).
// Returns 0 if the slice is empty.
func AveragePrice(books []domain.Book) int {
	if len(books) == 0 {
		return 0
	}
	var total int
	for _, b := range books {
		total += b.Price().Amount()
	}
	return total / len(books)
}

// MedianPrice returns the median price of a slice of books (in cents).
// Returns 0 if the slice is empty.
func MedianPrice(books []domain.Book) int {
	if len(books) == 0 {
		return 0
	}

	prices := make([]int, len(books))
	for i, b := range books {
		prices[i] = b.Price().Amount()
	}
	sort.Ints(prices)

	mid := len(prices) / 2
	if len(prices)%2 == 0 {
		return (prices[mid-1] + prices[mid]) / 2
	}
	return prices[mid]
}

// PriceRange returns the minimum and maximum prices in a slice of books.
// Returns (0, 0) if the slice is empty.
func PriceRange(books []domain.Book) (min, max int) {
	if len(books) == 0 {
		return 0, 0
	}

	min = books[0].Price().Amount()
	max = min
	for _, b := range books[1:] {
		p := b.Price().Amount()
		if p < min {
			min = p
		}
		if p > max {
			max = p
		}
	}
	return min, max
}

// GenreBreakdown returns a map of genre to number of books.
func GenreBreakdown(books []domain.Book) map[domain.Genre]int {
	result := make(map[domain.Genre]int)
	for _, b := range books {
		result[b.Genre()]++
	}
	return result
}

// MostExpensive returns the n most expensive books, sorted by price descending.
func MostExpensive(books []domain.Book, n int) []domain.Book {
	if n <= 0 || len(books) == 0 {
		return nil
	}

	sorted := make([]domain.Book, len(books))
	copy(sorted, books)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Price().Amount() > sorted[j].Price().Amount()
	})

	if n > len(sorted) {
		n = len(sorted)
	}
	return sorted[:n]
}
