package domain

import "testing"

func TestNewISBN_ValidISBN13(t *testing.T) {
	isbn, err := NewISBN("978-0-306-40615-7")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if isbn.String() != "9780306406157" {
		t.Errorf("got %s, want 9780306406157", isbn.String())
	}
}

func TestNewISBN_InvalidLength(t *testing.T) {
	_, err := NewISBN("123")
	if err == nil {
		t.Fatal("expected error for short ISBN")
	}
}

func TestNewISBN_InvalidChecksum(t *testing.T) {
	_, err := NewISBN("9780306406158")
	if err == nil {
		t.Fatal("expected error for invalid checksum")
	}
}
