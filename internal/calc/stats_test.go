package calc

import (
	"testing"
	"time"

	"github.com/sergekukharev/agent-test-writer-validator/internal/domain"
)

// makeBook is a test helper that creates a Book with the given price (in cents) and genre.
// It uses fixed values for all other fields since they are irrelevant to stats calculations.
func makeBook(t *testing.T, priceCents int, genre domain.Genre) domain.Book {
	t.Helper()

	isbn, err := domain.NewISBN("978-0-306-40615-7")
	if err != nil {
		t.Fatalf("failed to create ISBN: %v", err)
	}
	author, err := domain.NewAuthor("Test", "Author")
	if err != nil {
		t.Fatalf("failed to create Author: %v", err)
	}
	money, err := domain.NewMoney(priceCents, "EUR")
	if err != nil {
		t.Fatalf("failed to create Money: %v", err)
	}
	book, err := domain.NewBook(isbn, "Test Book", author, money, time.Now(), genre)
	if err != nil {
		t.Fatalf("failed to create Book: %v", err)
	}
	return book
}

// makeBooks is a test helper that creates a slice of fiction books with the given prices.
func makeBooks(t *testing.T, prices ...int) []domain.Book {
	t.Helper()
	books := make([]domain.Book, len(prices))
	for i, p := range prices {
		books[i] = makeBook(t, p, domain.GenreFiction)
	}
	return books
}

func TestAveragePrice(t *testing.T) {
	tests := []struct {
		name   string
		prices []int
		want   int
	}{
		{
			name:   "empty slice returns 0",
			prices: nil,
			want:   0,
		},
		{
			name:   "single book returns its price",
			prices: []int{1000},
			want:   1000,
		},
		{
			name:   "two books with same price",
			prices: []int{500, 500},
			want:   500,
		},
		{
			name:   "multiple books",
			prices: []int{1000, 2000, 3000},
			want:   2000,
		},
		{
			name:   "integer division truncates",
			prices: []int{100, 200},
			want:   150,
		},
		{
			name:   "integer division truncates remainder",
			prices: []int{100, 200, 200},
			want:   166, // 500 / 3 = 166
		},
		{
			name:   "includes zero-priced books",
			prices: []int{0, 0, 300},
			want:   100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			books := makeBooks(t, tt.prices...)
			got := AveragePrice(books)
			if got != tt.want {
				t.Errorf("AveragePrice() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestMedianPrice(t *testing.T) {
	tests := []struct {
		name   string
		prices []int
		want   int
	}{
		{
			name:   "empty slice returns 0",
			prices: nil,
			want:   0,
		},
		{
			name:   "single book returns its price",
			prices: []int{1500},
			want:   1500,
		},
		{
			name:   "odd count returns middle element",
			prices: []int{100, 200, 300},
			want:   200,
		},
		{
			name:   "even count returns average of two middle elements",
			prices: []int{100, 200, 300, 400},
			want:   250, // (200 + 300) / 2
		},
		{
			name:   "unsorted input is sorted internally",
			prices: []int{300, 100, 200},
			want:   200,
		},
		{
			name:   "even count with integer division truncation",
			prices: []int{100, 200},
			want:   150, // (100 + 200) / 2
		},
		{
			name:   "even count truncates odd sum",
			prices: []int{100, 201},
			want:   150, // (100 + 201) / 2 = 150
		},
		{
			name:   "five elements returns third",
			prices: []int{500, 100, 300, 200, 400},
			want:   300, // sorted: 100, 200, 300, 400, 500
		},
		{
			name:   "duplicate values",
			prices: []int{200, 200, 200},
			want:   200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			books := makeBooks(t, tt.prices...)
			got := MedianPrice(books)
			if got != tt.want {
				t.Errorf("MedianPrice() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestMedianPrice_DoesNotMutateInput(t *testing.T) {
	books := makeBooks(t, 300, 100, 200)
	original := make([]domain.Book, len(books))
	copy(original, books)

	MedianPrice(books)

	for i, b := range books {
		if b.Price().Amount() != original[i].Price().Amount() {
			t.Errorf("MedianPrice mutated input: index %d changed from %d to %d",
				i, original[i].Price().Amount(), b.Price().Amount())
		}
	}
}

func TestPriceRange(t *testing.T) {
	tests := []struct {
		name    string
		prices  []int
		wantMin int
		wantMax int
	}{
		{
			name:    "empty slice returns 0 0",
			prices:  nil,
			wantMin: 0,
			wantMax: 0,
		},
		{
			name:    "single book returns same min and max",
			prices:  []int{1000},
			wantMin: 1000,
			wantMax: 1000,
		},
		{
			name:    "two different prices",
			prices:  []int{500, 1500},
			wantMin: 500,
			wantMax: 1500,
		},
		{
			name:    "min at start max at end",
			prices:  []int{100, 200, 300},
			wantMin: 100,
			wantMax: 300,
		},
		{
			name:    "min at end max at start",
			prices:  []int{300, 200, 100},
			wantMin: 100,
			wantMax: 300,
		},
		{
			name:    "min and max in middle",
			prices:  []int{200, 100, 300, 250},
			wantMin: 100,
			wantMax: 300,
		},
		{
			name:    "all same prices",
			prices:  []int{500, 500, 500},
			wantMin: 500,
			wantMax: 500,
		},
		{
			name:    "includes zero price",
			prices:  []int{0, 100, 200},
			wantMin: 0,
			wantMax: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			books := makeBooks(t, tt.prices...)
			gotMin, gotMax := PriceRange(books)
			if gotMin != tt.wantMin || gotMax != tt.wantMax {
				t.Errorf("PriceRange() = (%d, %d), want (%d, %d)",
					gotMin, gotMax, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestGenreBreakdown(t *testing.T) {
	tests := []struct {
		name   string
		genres []domain.Genre
		want   map[domain.Genre]int
	}{
		{
			name:   "empty slice returns empty map",
			genres: nil,
			want:   map[domain.Genre]int{},
		},
		{
			name:   "single genre",
			genres: []domain.Genre{domain.GenreFiction},
			want:   map[domain.Genre]int{domain.GenreFiction: 1},
		},
		{
			name:   "multiple books same genre",
			genres: []domain.Genre{domain.GenreScience, domain.GenreScience, domain.GenreScience},
			want:   map[domain.Genre]int{domain.GenreScience: 3},
		},
		{
			name: "multiple genres",
			genres: []domain.Genre{
				domain.GenreFiction,
				domain.GenreScience,
				domain.GenreFiction,
				domain.GenreBiography,
				domain.GenreScience,
				domain.GenreFiction,
			},
			want: map[domain.Genre]int{
				domain.GenreFiction:   3,
				domain.GenreScience:   2,
				domain.GenreBiography: 1,
			},
		},
		{
			name: "all genres represented",
			genres: []domain.Genre{
				domain.GenreFiction,
				domain.GenreNonFiction,
				domain.GenreScience,
				domain.GenreBiography,
				domain.GenreChildren,
			},
			want: map[domain.Genre]int{
				domain.GenreFiction:    1,
				domain.GenreNonFiction: 1,
				domain.GenreScience:    1,
				domain.GenreBiography:  1,
				domain.GenreChildren:   1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			books := make([]domain.Book, len(tt.genres))
			for i, g := range tt.genres {
				books[i] = makeBook(t, 100, g)
			}

			got := GenreBreakdown(books)

			if len(got) != len(tt.want) {
				t.Fatalf("GenreBreakdown() returned %d genres, want %d", len(got), len(tt.want))
			}
			for genre, wantCount := range tt.want {
				if gotCount, ok := got[genre]; !ok {
					t.Errorf("GenreBreakdown() missing genre %q", genre)
				} else if gotCount != wantCount {
					t.Errorf("GenreBreakdown()[%q] = %d, want %d", genre, gotCount, wantCount)
				}
			}
		})
	}
}

func TestMostExpensive(t *testing.T) {
	tests := []struct {
		name       string
		prices     []int
		n          int
		wantPrices []int
	}{
		{
			name:       "empty slice returns nil",
			prices:     nil,
			n:          3,
			wantPrices: nil,
		},
		{
			name:       "n is zero returns nil",
			prices:     []int{100, 200},
			n:          0,
			wantPrices: nil,
		},
		{
			name:       "n is negative returns nil",
			prices:     []int{100, 200},
			n:          -1,
			wantPrices: nil,
		},
		{
			name:       "n equals length returns all sorted descending",
			prices:     []int{100, 300, 200},
			n:          3,
			wantPrices: []int{300, 200, 100},
		},
		{
			name:       "n greater than length returns all sorted descending",
			prices:     []int{100, 200},
			n:          5,
			wantPrices: []int{200, 100},
		},
		{
			name:       "n less than length returns top n",
			prices:     []int{100, 500, 300, 200, 400},
			n:          3,
			wantPrices: []int{500, 400, 300},
		},
		{
			name:       "single book with n=1",
			prices:     []int{999},
			n:          1,
			wantPrices: []int{999},
		},
		{
			name:       "top 1 of many",
			prices:     []int{100, 300, 200},
			n:          1,
			wantPrices: []int{300},
		},
		{
			name:       "duplicate prices",
			prices:     []int{200, 200, 100},
			n:          2,
			wantPrices: []int{200, 200},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			books := makeBooks(t, tt.prices...)
			got := MostExpensive(books, tt.n)

			if tt.wantPrices == nil {
				if got != nil {
					t.Fatalf("MostExpensive() = %v, want nil", got)
				}
				return
			}

			if len(got) != len(tt.wantPrices) {
				t.Fatalf("MostExpensive() returned %d books, want %d", len(got), len(tt.wantPrices))
			}
			for i, want := range tt.wantPrices {
				if got[i].Price().Amount() != want {
					t.Errorf("MostExpensive()[%d].Price().Amount() = %d, want %d",
						i, got[i].Price().Amount(), want)
				}
			}
		})
	}
}

func TestMostExpensive_DoesNotMutateInput(t *testing.T) {
	books := makeBooks(t, 300, 100, 200)
	originalPrices := make([]int, len(books))
	for i, b := range books {
		originalPrices[i] = b.Price().Amount()
	}

	MostExpensive(books, 2)

	for i, b := range books {
		if b.Price().Amount() != originalPrices[i] {
			t.Errorf("MostExpensive mutated input: index %d changed from %d to %d",
				i, originalPrices[i], b.Price().Amount())
		}
	}
}
