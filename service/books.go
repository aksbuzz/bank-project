package service

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/aksbuzz/library-project/store/db"
	"github.com/go-chi/chi/v5"
)

func (s *Service) GetBooks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cursor, err := parseInt32(r.URL.Query().Get("cursor"))
	if err != nil {
		s.logger.ErrorContext(ctx, "error parsing cursor", slog.String("service.Getbooks", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	limit, err := parseInt32(r.URL.Query().Get("limit"))
	if err != nil {
		s.logger.ErrorContext(ctx, "error parsing limit", slog.String("service.Getbooks", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	listbooksparams := db.ListBooksParams{
		Cursor: cursor,
		Limit:  limit,
	}

	list, err := s.store.DB.ListBooks(ctx, listbooksparams)
	if err != nil {
		s.logger.ErrorContext(ctx, "error listing books", slog.String("service.Getbooks", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	books, err := json.Marshal(list)
	if err != nil {
		s.logger.ErrorContext(ctx, "error marshalling books", slog.String("service.Getbooks", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(books)
	w.WriteHeader(http.StatusOK)
}

func (s *Service) GetBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseInt32(chi.URLParam(r, "id"))
	if err != nil {
		s.logger.ErrorContext(ctx, "error parsing id", slog.String("service.Getbook", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book, err := s.store.DB.GetBook(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "error getting book", slog.String("service.Getbook", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Service) CreateBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var book db.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		s.logger.ErrorContext(ctx, "error decoding book", slog.String("service.Createbook", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	create := db.CreateBookParams{
		Title:    book.Title,
		Author:   book.Author,
		Year:     book.Year,
		Quantity: book.Quantity,
	}

	newBook, err := s.store.DB.CreateBook(ctx, create)
	if err != nil {
		s.logger.ErrorContext(ctx, "error creating book", slog.String("service.Createbook", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(&newBook); err != nil {
		s.logger.ErrorContext(ctx, "error encoding book", slog.String("service.Createbook", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Service) DeleteBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseInt32(chi.URLParam(r, "id"))
	if err != nil {
		s.logger.ErrorContext(ctx, "error parsing id", slog.String("service.Deletebook", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.store.DB.DeleteBook(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "error deleting book", slog.String("service.Deletebook", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
