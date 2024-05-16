package service

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/aksbuzz/library-project/store"
	"github.com/go-chi/chi/v5"
)

type Service struct {
	logger *slog.Logger
	store  *store.Store
}

func New(store *store.Store, logger *slog.Logger) *Service {
	return &Service{store: store, logger: logger}
}

func (s *Service) Register(r *chi.Mux) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/books", s.GetBooks)
		r.Post("/books", s.CreateBook)
		r.Get("/books/{id}", s.GetBook)
		r.Delete("/books/{id}", s.DeleteBook)

		r.Get("/members", s.GetMembers)
		r.Post("/members", s.CreateMember)
		r.Get("/members/{id}", s.GetMember)
		r.Put("/members/{id}", s.UpdateMember)
		r.Patch("/members/{id}", s.UpdateMembership)

		r.Get("/loans", s.GetLoans)
		r.Post("/loans", s.CreateLoan)
		r.Get("/loans/{id}", s.GetLoan)
		r.Patch("/loans/{id}", s.UpdateLoanReturnDate)
	})
}

func parseInt32(urlParam string, defaultVal ...int32) (int32, error) {
	if len(defaultVal) > 0 && urlParam == "" {
		return defaultVal[0], nil
	}

	param, err := strconv.Atoi(urlParam)
	if err != nil {
		return 0, fmt.Errorf("failed to parse param: %w", err)
	}

	return int32(param), nil
}
