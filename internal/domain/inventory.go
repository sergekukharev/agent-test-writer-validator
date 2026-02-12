package domain

import (
	"errors"
	"fmt"
)

// StockEntry tracks the available and reserved copies for a single book.
type StockEntry struct {
	book     Book
	total    int
	reserved int
}

func NewStockEntry(book Book, total int) (StockEntry, error) {
	if total < 0 {
		return StockEntry{}, errors.New("total stock must not be negative")
	}
	return StockEntry{book: book, total: total}, nil
}

func (s StockEntry) Book() Book    { return s.book }
func (s StockEntry) Total() int    { return s.total }
func (s StockEntry) Reserved() int { return s.reserved }

func (s StockEntry) Available() int {
	return s.total - s.reserved
}

// Reserve attempts to reserve n copies. Returns an error if insufficient stock.
func (s *StockEntry) Reserve(n int) error {
	if n <= 0 {
		return errors.New("reservation quantity must be positive")
	}
	if s.Available() < n {
		return fmt.Errorf("insufficient stock: %d available, %d requested", s.Available(), n)
	}
	s.reserved += n
	return nil
}

// Release releases n previously reserved copies back to available stock.
func (s *StockEntry) Release(n int) error {
	if n <= 0 {
		return errors.New("release quantity must be positive")
	}
	if n > s.reserved {
		return fmt.Errorf("cannot release %d: only %d reserved", n, s.reserved)
	}
	s.reserved -= n
	return nil
}

// Restock adds n copies to total stock.
func (s *StockEntry) Restock(n int) error {
	if n <= 0 {
		return errors.New("restock quantity must be positive")
	}
	s.total += n
	return nil
}

// IsLowStock returns true if available copies are below the threshold.
func (s StockEntry) IsLowStock(threshold int) bool {
	return s.Available() < threshold
}

// Inventory manages stock for multiple books.
type Inventory struct {
	entries map[string]*StockEntry // keyed by ISBN string
}

func NewInventory() *Inventory {
	return &Inventory{entries: make(map[string]*StockEntry)}
}

func (inv *Inventory) Add(entry StockEntry) {
	inv.entries[entry.book.ISBN().String()] = &entry
}

func (inv *Inventory) Find(isbn ISBN) (*StockEntry, error) {
	e, ok := inv.entries[isbn.String()]
	if !ok {
		return nil, fmt.Errorf("book %s not found in inventory", isbn)
	}
	return e, nil
}

// LowStockBooks returns all books with available stock below the threshold.
func (inv *Inventory) LowStockBooks(threshold int) []Book {
	var result []Book
	for _, e := range inv.entries {
		if e.IsLowStock(threshold) {
			result = append(result, e.book)
		}
	}
	return result
}

// TotalValue calculates the total value of all stock (total copies * price).
func (inv *Inventory) TotalValue() int {
	var total int
	for _, e := range inv.entries {
		total += e.book.Price().Amount() * e.total
	}
	return total
}
