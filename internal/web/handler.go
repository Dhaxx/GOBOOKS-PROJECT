package web

import (
	"encoding/json"
	"gobooks/internal/service"
	"net/http"
	"strconv"
)

type BookHandlers struct {
	service *service.BookService
}

func NewBookHandlers(service *service.BookService) *BookHandlers {
	return &BookHandlers{service: service}
}

func (h *BookHandlers) GetBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.service.GetBooks()
	if err != nil {
		http.Error(w, "failed to get books", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func (h *BookHandlers) CreateBook(w http.ResponseWriter, r *http.Request) {
	var book service.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, "failed to decode book", http.StatusBadRequest)
		return
	}
	err = h.service.Create(&book)
	if err != nil {
		http.Error(w, "failed to create book", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandlers) GetBook(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "failed to parse id", http.StatusBadRequest)
		return
	}

	book, err := h.service.GetBookId(id)
	if err != nil {
		http.Error(w, "failed to get book", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandlers) UpdateBook(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid book ID", http.StatusBadRequest)
		return
	}

	var book service.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	book.ID = id

	if err := h.service.UpdateBook(&book); err != nil {
		http.Error(w, "failed to update book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandlers) DeleteBook(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "failed to parse id", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteBook(id)
	if err != nil {
		http.Error(w, "failed to delete book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}