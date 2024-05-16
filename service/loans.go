package service

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/aksbuzz/library-project/store/db"
	"github.com/go-chi/chi/v5"
)

func (s *Service) GetLoans(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cursor, err := parseInt32(r.URL.Query().Get("cursor"), 0)
	if err != nil {
		s.logger.ErrorContext(ctx, "error parsing cursor", slog.String("service.GetLoans", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	limit, err := parseInt32(r.URL.Query().Get("limit"), 10)
	if err != nil {
		s.logger.ErrorContext(ctx, "error parsing limit", slog.String("service.GetLoans", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	listloansparams := db.ListLoansParams{
		Cursor: cursor,
		Limit:  limit,
	}

	list, err := s.store.DB.ListLoans(ctx, listloansparams)
	if err != nil {
		s.logger.ErrorContext(ctx, "error listing loans", slog.String("service.GetLoans", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	loans, err := json.Marshal(list)
	if err != nil {
		s.logger.ErrorContext(ctx, "error marshalling loans", slog.String("service.GetLoans", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(loans)
	w.WriteHeader(http.StatusOK)
}

func (s *Service) GetLoan(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseInt32(chi.URLParam(r, "id"))
	if err != nil {
		s.logger.ErrorContext(ctx, "error parsing id", slog.String("service.GetLoan", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	loan, err := s.store.DB.GetLoan(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "error getting loan", slog.String("service.GetLoan", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(&loan); err != nil {
		s.logger.ErrorContext(ctx, "error encoding loan", slog.String("service.GetLoan", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Service) CreateLoan(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var loan db.Loan
	if err := json.NewDecoder(r.Body).Decode(&loan); err != nil {
		s.logger.ErrorContext(ctx, "error decoding loan", slog.String("service.CreateLoan", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var newLoan db.Loan
	err := s.store.DB.ExecTransaction(ctx, func(qtx *db.Queries) error {
		create := db.CreateLoanParams{
			BookID:   loan.BookID,
			MemberID: loan.MemberID,
		}

		book, err := qtx.GetBookForUpdate(ctx, loan.BookID)
		if err != nil {
			return err
		}

		loan, err := qtx.CreateLoan(ctx, create)
		if err != nil {
			return err
		}

		update := db.UpdateBookQuantityParams{
			ID:       book.ID,
			Quantity: book.Quantity - 1,
		}
		err = qtx.UpdateBookQuantity(ctx, update)
		if err != nil {
			return err
		}

		newLoan = loan
		return nil
	})
	if err != nil {
		s.logger.ErrorContext(ctx, "error creating loan", slog.String("service.CreateLoan", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(&newLoan); err != nil {
		s.logger.ErrorContext(ctx, "error encoding loan", slog.String("service.CreateLoan", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Service) UpdateLoanReturnDate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseInt32(chi.URLParam(r, "id"))
	if err != nil {
		s.logger.ErrorContext(ctx, "error parsing id", slog.String("service.UpdateLoanReturnDate", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var body db.Loan
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		s.logger.ErrorContext(ctx, "error decoding loan", slog.String("service.UpdateLoanReturnDate", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	loan, err := s.store.DB.GetLoan(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "error getting loan", slog.String("service.UpdateLoanReturnDate", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if loan.ReturnDate.Valid {
		s.logger.ErrorContext(ctx, "error updating loan", slog.String("service.UpdateLoanReturnDate", "loan already returned"))
		http.Error(w, "loan already returned", http.StatusBadRequest)
		return
	}

	err = s.store.DB.ExecTransaction(ctx, func(qtx *db.Queries) error {
		book, err := qtx.GetBookForUpdate(ctx, loan.BookID)
		if err != nil {
			return err
		}

		update := db.UpdateLoanReturnDateParams{
			ID:         id,
			ReturnDate: body.ReturnDate,
		}
		err = qtx.UpdateLoanReturnDate(ctx, update)
		if err != nil {
			return err
		}

		updateBook := db.UpdateBookQuantityParams{
			ID:       loan.BookID,
			Quantity: book.Quantity + 1,
		}
		err = qtx.UpdateBookQuantity(ctx, updateBook)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		s.logger.ErrorContext(ctx, "error updating loan", slog.String("service.UpdateLoanReturnDate", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
