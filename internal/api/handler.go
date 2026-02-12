package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sergekukharev/agent-test-writer-validator/internal/domain"
	"github.com/sergekukharev/agent-test-writer-validator/internal/storage"
)

type Handler struct {
	repo *storage.BookRepository
}

func NewHandler(repo *storage.BookRepository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /books", h.ListBooks)
	mux.HandleFunc("GET /books/{isbn}", h.GetBook)
	mux.HandleFunc("POST /books", h.CreateBook)
	mux.HandleFunc("DELETE /books/{isbn}", h.DeleteBook)
	return mux
}

func (h *Handler) ListBooks(w http.ResponseWriter, r *http.Request) {
	books := h.repo.FindAll()

	resp := ListResponse{
		Books: make([]BookResponse, 0, len(books)),
		Count: len(books),
	}
	for _, b := range books {
		resp.Books = append(resp.Books, toBookResponse(b))
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetBook(w http.ResponseWriter, r *http.Request) {
	isbn := r.PathValue("isbn")
	book, err := h.repo.FindByISBN(isbn)
	if err != nil {
		writeError(w, http.StatusNotFound, "book not found")
		return
	}
	writeJSON(w, http.StatusOK, toBookResponse(book))
}

type CreateBookRequest struct {
	ISBN      string `json:"isbn"`
	Title     string `json:"title"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	PriceCents int   `json:"price_cents"`
	Currency  string `json:"currency"`
	Genre     string `json:"genre"`
}

func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var req CreateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	isbn, err := domain.NewISBN(req.ISBN)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	author, err := domain.NewAuthor(req.FirstName, req.LastName)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	price, err := domain.NewMoney(req.PriceCents, req.Currency)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	book, err := domain.NewBook(isbn, req.Title, author, price, time.Now(), domain.Genre(req.Genre))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.repo.Save(book)
	writeJSON(w, http.StatusCreated, toBookResponse(book))
}

func (h *Handler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	isbn := r.PathValue("isbn")
	if err := h.repo.Delete(isbn); err != nil {
		writeError(w, http.StatusNotFound, "book not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func toBookResponse(b domain.Book) BookResponse {
	return BookResponse{
		ISBN:      b.ISBN().String(),
		Title:     b.Title(),
		Author:    b.Author().FullName(),
		Price:     b.Price().Display(),
		Genre:     string(b.Genre()),
		IsClassic: b.IsClassic(),
	}
}
