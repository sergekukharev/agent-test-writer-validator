package storage

import (
	"strings"

	"github.com/sergekukharev/agent-test-writer-validator/internal/domain"
)

// FilterFunc returns true for books that match the filter criteria.
type FilterFunc func(domain.Book) bool

// ByGenre returns a filter that matches books of the given genre.
func ByGenre(genre domain.Genre) FilterFunc {
	return func(b domain.Book) bool {
		return b.Genre() == genre
	}
}

// ByAuthorLastName returns a filter that matches books by the author's last name (case-insensitive).
func ByAuthorLastName(name string) FilterFunc {
	lower := strings.ToLower(name)
	return func(b domain.Book) bool {
		return strings.ToLower(b.Author().LastName()) == lower
	}
}

// ByPriceRange returns a filter matching books within the given price range (inclusive, in cents).
func ByPriceRange(minCents, maxCents int) FilterFunc {
	return func(b domain.Book) bool {
		p := b.Price().Amount()
		return p >= minCents && p <= maxCents
	}
}

// ByTitleContains returns a filter that matches books whose title contains the substring (case-insensitive).
func ByTitleContains(substr string) FilterFunc {
	lower := strings.ToLower(substr)
	return func(b domain.Book) bool {
		return strings.Contains(strings.ToLower(b.Title()), lower)
	}
}

// Apply runs the given filters on a book slice, returning only books that match all filters.
func Apply(books []domain.Book, filters ...FilterFunc) []domain.Book {
	var result []domain.Book
	for _, b := range books {
		if matchesAll(b, filters) {
			result = append(result, b)
		}
	}
	return result
}

func matchesAll(b domain.Book, filters []FilterFunc) bool {
	for _, f := range filters {
		if !f(b) {
			return false
		}
	}
	return true
}
