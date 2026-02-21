package storage

import (
	"sync"
	"testing"
	"time"

	"github.com/sergekukharev/agent-test-writer-validator/internal/domain"
)

func createTestBook(t *testing.T, isbnStr, title string) domain.Book {
	t.Helper()
	isbn, err := domain.NewISBN(isbnStr)
	if err != nil {
		t.Fatalf("failed to create ISBN: %v", err)
	}
	author, err := domain.NewAuthor("Test", "Author")
	if err != nil {
		t.Fatalf("failed to create author: %v", err)
	}
	price, err := domain.NewMoney(1999, "USD")
	if err != nil {
		t.Fatalf("failed to create money: %v", err)
	}
	book, err := domain.NewBook(isbn, title, author, price, time.Now(), domain.GenreFiction)
	if err != nil {
		t.Fatalf("failed to create book: %v", err)
	}
	return book
}

func TestNewBookRepository(t *testing.T) {
	repo := NewBookRepository()
	if repo == nil {
		t.Fatal("NewBookRepository returned nil")
	}
	if repo.Count() != 0 {
		t.Errorf("new repository should be empty, got count: %d", repo.Count())
	}
}

func TestBookRepository_Save(t *testing.T) {
	repo := NewBookRepository()
	book := createTestBook(t, "9780306406157", "Test Book")

	repo.Save(book)

	if repo.Count() != 1 {
		t.Errorf("expected count 1, got %d", repo.Count())
	}

	retrieved, err := repo.FindByISBN(book.ISBN().String())
	if err != nil {
		t.Fatalf("failed to find saved book: %v", err)
	}
	if retrieved.Title() != "Test Book" {
		t.Errorf("expected title 'Test Book', got %s", retrieved.Title())
	}
}

func TestBookRepository_Save_Overwrite(t *testing.T) {
	repo := NewBookRepository()
	isbn := "9780306406157"
	
	book1 := createTestBook(t, isbn, "First Title")
	book2 := createTestBook(t, isbn, "Second Title")

	repo.Save(book1)
	repo.Save(book2)

	if repo.Count() != 1 {
		t.Errorf("expected count 1 after overwrite, got %d", repo.Count())
	}

	retrieved, err := repo.FindByISBN(isbn)
	if err != nil {
		t.Fatalf("failed to find book: %v", err)
	}
	if retrieved.Title() != "Second Title" {
		t.Errorf("expected title 'Second Title', got %s", retrieved.Title())
	}
}

func TestBookRepository_FindByISBN_NotFound(t *testing.T) {
	repo := NewBookRepository()

	_, err := repo.FindByISBN("9780306406157")
	if err == nil {
		t.Fatal("expected error when finding non-existent book")
	}
	
	expectedMsg := "book 9780306406157 not found"
	if err.Error() != expectedMsg {
		t.Errorf("expected error message %q, got %q", expectedMsg, err.Error())
	}
}

func TestBookRepository_FindByISBN_Found(t *testing.T) {
	repo := NewBookRepository()
	book := createTestBook(t, "9780306406157", "Found Book")
	repo.Save(book)

	retrieved, err := repo.FindByISBN("9780306406157")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if retrieved.Title() != "Found Book" {
		t.Errorf("expected title 'Found Book', got %s", retrieved.Title())
	}
	if retrieved.ISBN().String() != "9780306406157" {
		t.Errorf("expected ISBN '9780306406157', got %s", retrieved.ISBN().String())
	}
}

func TestBookRepository_FindAll_Empty(t *testing.T) {
	repo := NewBookRepository()

	books := repo.FindAll()
	if len(books) != 0 {
		t.Errorf("expected empty slice, got %d books", len(books))
	}
	if books == nil {
		t.Error("FindAll should return non-nil slice")
	}
}

func TestBookRepository_FindAll_Single(t *testing.T) {
	repo := NewBookRepository()
	book := createTestBook(t, "9780306406157", "Single Book")
	repo.Save(book)

	books := repo.FindAll()
	if len(books) != 1 {
		t.Fatalf("expected 1 book, got %d", len(books))
	}
	if books[0].Title() != "Single Book" {
		t.Errorf("expected title 'Single Book', got %s", books[0].Title())
	}
}

func TestBookRepository_FindAll_Multiple(t *testing.T) {
	repo := NewBookRepository()
	
	book1 := createTestBook(t, "9780306406157", "Book One")
	book2 := createTestBook(t, "9781234567897", "Book Two")
	book3 := createTestBook(t, "9783161484100", "Book Three")
	
	repo.Save(book1)
	repo.Save(book2)
	repo.Save(book3)

	books := repo.FindAll()
	if len(books) != 3 {
		t.Fatalf("expected 3 books, got %d", len(books))
	}

	titles := make(map[string]bool)
	for _, b := range books {
		titles[b.Title()] = true
	}

	expectedTitles := []string{"Book One", "Book Two", "Book Three"}
	for _, title := range expectedTitles {
		if !titles[title] {
			t.Errorf("expected to find book with title %q", title)
		}
	}
}

func TestBookRepository_Delete_Success(t *testing.T) {
	repo := NewBookRepository()
	book := createTestBook(t, "9780306406157", "To Delete")
	repo.Save(book)

	err := repo.Delete("9780306406157")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.Count() != 0 {
		t.Errorf("expected count 0 after delete, got %d", repo.Count())
	}

	_, err = repo.FindByISBN("9780306406157")
	if err == nil {
		t.Error("expected error when finding deleted book")
	}
}

func TestBookRepository_Delete_NotFound(t *testing.T) {
	repo := NewBookRepository()

	err := repo.Delete("9780306406157")
	if err == nil {
		t.Fatal("expected error when deleting non-existent book")
	}

	expectedMsg := "book 9780306406157 not found"
	if err.Error() != expectedMsg {
		t.Errorf("expected error message %q, got %q", expectedMsg, err.Error())
	}
}

func TestBookRepository_Delete_Multiple(t *testing.T) {
	repo := NewBookRepository()
	
	book1 := createTestBook(t, "9780306406157", "Book One")
	book2 := createTestBook(t, "9781234567897", "Book Two")
	
	repo.Save(book1)
	repo.Save(book2)

	err := repo.Delete("9780306406157")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.Count() != 1 {
		t.Errorf("expected count 1 after delete, got %d", repo.Count())
	}

	_, err = repo.FindByISBN("9780306406157")
	if err == nil {
		t.Error("expected error when finding deleted book")
	}

	retrieved, err := repo.FindByISBN("9781234567897")
	if err != nil {
		t.Fatalf("unexpected error finding remaining book: %v", err)
	}
	if retrieved.Title() != "Book Two" {
		t.Errorf("expected title 'Book Two', got %s", retrieved.Title())
	}
}

func TestBookRepository_Count_Empty(t *testing.T) {
	repo := NewBookRepository()

	if repo.Count() != 0 {
		t.Errorf("expected count 0, got %d", repo.Count())
	}
}

func TestBookRepository_Count_AfterOperations(t *testing.T) {
	repo := NewBookRepository()

	if repo.Count() != 0 {
		t.Errorf("expected count 0, got %d", repo.Count())
	}

	book1 := createTestBook(t, "9780306406157", "Book One")
	repo.Save(book1)
	if repo.Count() != 1 {
		t.Errorf("expected count 1 after save, got %d", repo.Count())
	}

	book2 := createTestBook(t, "9781234567897", "Book Two")
	repo.Save(book2)
	if repo.Count() != 2 {
		t.Errorf("expected count 2 after second save, got %d", repo.Count())
	}

	repo.Delete("9780306406157")
	if repo.Count() != 1 {
		t.Errorf("expected count 1 after delete, got %d", repo.Count())
	}

	repo.Delete("9781234567897")
	if repo.Count() != 0 {
		t.Errorf("expected count 0 after second delete, got %d", repo.Count())
	}
}

func TestBookRepository_ConcurrentSave(t *testing.T) {
	repo := NewBookRepository()
	var wg sync.WaitGroup
	
	validISBNs := []string{
		"9780306406157",
		"9781234567897",
		"9783161484100",
		"9780262033848",
		"9780134685991",
		"9781449355739",
		"9780321573513",
		"9781593279509",
		"9780596517748",
		"9781491950357",
	}
	
	numGoroutines := len(validISBNs)
	wg.Add(numGoroutines)
	
	for i := 0; i < numGoroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			book := createTestBook(t, validISBNs[idx], "Concurrent Book")
			repo.Save(book)
		}(i)
	}
	
	wg.Wait()
	
	count := repo.Count()
	if count != numGoroutines {
		t.Errorf("expected count %d after concurrent saves, got %d", numGoroutines, count)
	}
}

func TestBookRepository_ConcurrentRead(t *testing.T) {
	repo := NewBookRepository()
	book := createTestBook(t, "9780306406157", "Concurrent Read")
	repo.Save(book)
	
	var wg sync.WaitGroup
	numGoroutines := 20
	wg.Add(numGoroutines)
	
	errors := make(chan error, numGoroutines)
	
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			_, err := repo.FindByISBN("9780306406157")
			if err != nil {
				errors <- err
			}
		}()
	}
	
	wg.Wait()
	close(errors)
	
	for err := range errors {
		t.Errorf("unexpected error during concurrent read: %v", err)
	}
}

func TestBookRepository_ConcurrentMixed(t *testing.T) {
	repo := NewBookRepository()
	
	book1 := createTestBook(t, "9780306406157", "Book One")
	book2 := createTestBook(t, "9781234567897", "Book Two")
	repo.Save(book1)
	repo.Save(book2)
	
	newISBNs := []string{
		"9783161484100",
		"9780262033848",
		"9780134685991",
	}
	
	var wg sync.WaitGroup
	
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			repo.FindAll()
		}()
	}
	
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			repo.FindByISBN("9780306406157")
		}()
	}
	
	for i := 0; i < len(newISBNs); i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			book := createTestBook(t, newISBNs[idx], "New Book")
			repo.Save(book)
		}(i)
	}
	
	wg.Wait()
	
	count := repo.Count()
	if count != 5 {
		t.Errorf("expected count 5 after concurrent operations, got %d", count)
	}
}

func TestBookRepository_ConcurrentDelete(t *testing.T) {
	repo := NewBookRepository()
	
	isbns := []string{
		"9780306406157",
		"9781234567897",
		"9783161484100",
		"9780262033848",
		"9780134685991",
	}
	
	for _, isbn := range isbns {
		book := createTestBook(t, isbn, "Book "+isbn)
		repo.Save(book)
	}
	
	var wg sync.WaitGroup
	wg.Add(len(isbns))
	
	for _, isbn := range isbns {
		go func(i string) {
			defer wg.Done()
			repo.Delete(i)
		}(isbn)
	}
	
	wg.Wait()
	
	if repo.Count() != 0 {
		t.Errorf("expected count 0 after concurrent deletes, got %d", repo.Count())
	}
}

func TestBookRepository_FindAll_DoesNotModifyInternal(t *testing.T) {
	repo := NewBookRepository()
	book := createTestBook(t, "9780306406157", "Original Book")
	repo.Save(book)

	books := repo.FindAll()
	if len(books) != 1 {
		t.Fatalf("expected 1 book, got %d", len(books))
	}

	books = append(books, createTestBook(t, "9781234567897", "New Book"))

	if repo.Count() != 1 {
		t.Errorf("modifying FindAll result should not affect repository, expected count 1, got %d", repo.Count())
	}
}
