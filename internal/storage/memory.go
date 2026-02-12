package storage

import (
	"fmt"
	"sync"

	"github.com/sergekukharev/agent-test-writer-validator/internal/domain"
)

// BookRepository stores books in memory.
type BookRepository struct {
	mu    sync.RWMutex
	books map[string]domain.Book // keyed by ISBN string
}

func NewBookRepository() *BookRepository {
	return &BookRepository{books: make(map[string]domain.Book)}
}

func (r *BookRepository) Save(book domain.Book) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.books[book.ISBN().String()] = book
}

func (r *BookRepository) FindByISBN(isbn string) (domain.Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, ok := r.books[isbn]
	if !ok {
		return domain.Book{}, fmt.Errorf("book %s not found", isbn)
	}
	return b, nil
}

func (r *BookRepository) FindAll() []domain.Book {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]domain.Book, 0, len(r.books))
	for _, b := range r.books {
		result = append(result, b)
	}
	return result
}

func (r *BookRepository) Delete(isbn string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.books[isbn]; !ok {
		return fmt.Errorf("book %s not found", isbn)
	}
	delete(r.books, isbn)
	return nil
}

func (r *BookRepository) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.books)
}
