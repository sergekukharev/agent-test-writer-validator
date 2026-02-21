package calc

import (
	"testing"
	"time"

	"github.com/sergekukharev/agent-test-writer-validator/internal/domain"
)

func mustBook(price int) domain.Book {
	isbn, _ := domain.NewISBN("978-0-306-40615-7")
	author, _ := domain.NewAuthor("Jane", "Doe")
	money, _ := domain.NewMoney(price, "USD")
	book, _ := domain.NewBook(isbn, "Test Book", author, money, time.Now(), domain.GenreFiction)
	return book
}

func mustBookWithGenre(price int, genre domain.Genre) domain.Book {
	isbn, _ := domain.NewISBN("978-0-306-40615-7")
	author, _ := domain.NewAuthor("Jane", "Doe")
	money, _ := domain.NewMoney(price, "USD")
	book, _ := domain.NewBook(isbn, "Test Book", author, money, time.Now(), genre)
	return book
}

// AveragePrice

func TestAveragePrice_EmptySlice(t *testing.T) {
	if got := AveragePrice(nil); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}

func TestAveragePrice(t *testing.T) {
	tests := []struct {
		name   string
		prices []int
		want   int
	}{
		{"single book", []int{1000}, 1000},
		{"two equal", []int{500, 500}, 500},
		{"two different", []int{400, 600}, 500},
		{"three books", []int{100, 200, 300}, 200},
		{"integer truncation", []int{100, 200}, 150},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			books := make([]domain.Book, len(tc.prices))
			for i, p := range tc.prices {
				books[i] = mustBook(p)
			}
			if got := AveragePrice(books); got != tc.want {
				t.Errorf("AveragePrice() = %d, want %d", got, tc.want)
			}
		})
	}
}

// MedianPrice

func TestMedianPrice_EmptySlice(t *testing.T) {
	if got := MedianPrice(nil); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}

func TestMedianPrice(t *testing.T) {
	tests := []struct {
		name   string
		prices []int
		want   int
	}{
		{"single book", []int{500}, 500},
		{"odd count sorted", []int{100, 300, 500}, 300},
		{"odd count unsorted", []int{500, 100, 300}, 300},
		{"even count", []int{100, 200, 300, 400}, 250},
		{"even count unsorted", []int{400, 100, 300, 200}, 250},
		{"even count integer truncation", []int{100, 301}, 200},
		{"two equal", []int{500, 500}, 500},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			books := make([]domain.Book, len(tc.prices))
			for i, p := range tc.prices {
				books[i] = mustBook(p)
			}
			if got := MedianPrice(books); got != tc.want {
				t.Errorf("MedianPrice() = %d, want %d", got, tc.want)
			}
		})
	}
}

// PriceRange

func TestPriceRange_EmptySlice(t *testing.T) {
	min, max := PriceRange(nil)
	if min != 0 || max != 0 {
		t.Errorf("expected (0, 0), got (%d, %d)", min, max)
	}
}

func TestPriceRange(t *testing.T) {
	tests := []struct {
		name    string
		prices  []int
		wantMin int
		wantMax int
	}{
		{"single book", []int{500}, 500, 500},
		{"two books", []int{200, 800}, 200, 800},
		{"already sorted", []int{100, 300, 500}, 100, 500},
		{"reverse sorted", []int{500, 300, 100}, 100, 500},
		{"all equal", []int{400, 400, 400}, 400, 400},
		{"many books", []int{300, 100, 500, 200, 400}, 100, 500},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			books := make([]domain.Book, len(tc.prices))
			for i, p := range tc.prices {
				books[i] = mustBook(p)
			}
			gotMin, gotMax := PriceRange(books)
			if gotMin != tc.wantMin || gotMax != tc.wantMax {
				t.Errorf("PriceRange() = (%d, %d), want (%d, %d)", gotMin, gotMax, tc.wantMin, tc.wantMax)
			}
		})
	}
}

// GenreBreakdown

func TestGenreBreakdown_EmptySlice(t *testing.T) {
	result := GenreBreakdown(nil)
	if len(result) != 0 {
		t.Errorf("expected empty map, got %v", result)
	}
}

func TestGenreBreakdown_SingleGenre(t *testing.T) {
	books := []domain.Book{
		mustBookWithGenre(100, domain.GenreFiction),
		mustBookWithGenre(200, domain.GenreFiction),
	}
	result := GenreBreakdown(books)
	if result[domain.GenreFiction] != 2 {
		t.Errorf("expected 2 fiction books, got %d", result[domain.GenreFiction])
	}
	if len(result) != 1 {
		t.Errorf("expected 1 genre, got %d", len(result))
	}
}

func TestGenreBreakdown_MultipleGenres(t *testing.T) {
	books := []domain.Book{
		mustBookWithGenre(100, domain.GenreFiction),
		mustBookWithGenre(200, domain.GenreFiction),
		mustBookWithGenre(300, domain.GenreScience),
		mustBookWithGenre(400, domain.GenreBiography),
	}
	result := GenreBreakdown(books)

	if result[domain.GenreFiction] != 2 {
		t.Errorf("expected 2 fiction, got %d", result[domain.GenreFiction])
	}
	if result[domain.GenreScience] != 1 {
		t.Errorf("expected 1 science, got %d", result[domain.GenreScience])
	}
	if result[domain.GenreBiography] != 1 {
		t.Errorf("expected 1 biography, got %d", result[domain.GenreBiography])
	}
	if len(result) != 3 {
		t.Errorf("expected 3 genres, got %d", len(result))
	}
}

// MostExpensive

func TestMostExpensive_EmptySlice(t *testing.T) {
	if got := MostExpensive(nil, 3); got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestMostExpensive_ZeroN(t *testing.T) {
	books := []domain.Book{mustBook(500)}
	if got := MostExpensive(books, 0); got != nil {
		t.Errorf("expected nil for n=0, got %v", got)
	}
}

func TestMostExpensive_NegativeN(t *testing.T) {
	books := []domain.Book{mustBook(500)}
	if got := MostExpensive(books, -1); got != nil {
		t.Errorf("expected nil for negative n, got %v", got)
	}
}

func TestMostExpensive(t *testing.T) {
	tests := []struct {
		name        string
		prices      []int
		n           int
		wantPrices  []int
	}{
		{
			name:       "top 1",
			prices:     []int{300, 100, 500, 200},
			n:          1,
			wantPrices: []int{500},
		},
		{
			name:       "top 2",
			prices:     []int{300, 100, 500, 200},
			n:          2,
			wantPrices: []int{500, 300},
		},
		{
			name:       "n exceeds length returns all sorted",
			prices:     []int{200, 100, 300},
			n:          10,
			wantPrices: []int{300, 200, 100},
		},
		{
			name:       "n equals length",
			prices:     []int{200, 100, 300},
			n:          3,
			wantPrices: []int{300, 200, 100},
		},
		{
			name:       "single book",
			prices:     []int{999},
			n:          1,
			wantPrices: []int{999},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			books := make([]domain.Book, len(tc.prices))
			for i, p := range tc.prices {
				books[i] = mustBook(p)
			}
			got := MostExpensive(books, tc.n)
			if len(got) != len(tc.wantPrices) {
				t.Fatalf("MostExpensive() returned %d books, want %d", len(got), len(tc.wantPrices))
			}
			for i, b := range got {
				if b.Price().Amount() != tc.wantPrices[i] {
					t.Errorf("got[%d].Price() = %d, want %d", i, b.Price().Amount(), tc.wantPrices[i])
				}
			}
		})
	}
}

func TestMostExpensive_DoesNotMutateInput(t *testing.T) {
	books := []domain.Book{
		mustBook(300),
		mustBook(100),
		mustBook(500),
	}
	original := make([]domain.Book, len(books))
	copy(original, books)

	MostExpensive(books, 2)

	for i, b := range books {
		if b.Price().Amount() != original[i].Price().Amount() {
			t.Errorf("input slice was mutated at index %d", i)
		}
	}
}
