package service

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/aksbuzz/library-project/store/db"
	"github.com/go-chi/chi/v5"
)

func (s *Service) GetMembers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cursor, err := parseInt32(r.URL.Query().Get("cursor"), 0)
	if err != nil {
		s.logger.ErrorContext(ctx, "error parsing cursor", slog.String("service.GetMembers", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	limit, err := parseInt32(r.URL.Query().Get("limit"), 10)
	if err != nil {
		s.logger.ErrorContext(ctx, "error parsing limit", slog.String("service.GetMembers", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	listmembersparams := db.ListMembersParams{
		Cursor: cursor,
		Limit:  limit,
	}

	list, err := s.store.DB.ListMembers(ctx, listmembersparams)
	if err != nil {
		s.logger.ErrorContext(ctx, "error listing members", slog.String("service.GetMembers", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	members, err := json.Marshal(list)
	if err != nil {
		s.logger.ErrorContext(ctx, "error marshalling members", slog.String("service.GetMembers", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(members)
	w.WriteHeader(http.StatusOK)
}

func (s *Service) GetMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseInt32(chi.URLParam(r, "id"))
	if err != nil {
		s.logger.ErrorContext(ctx, "error parsing id", slog.String("service.GetMember", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	member, err := s.store.DB.GetMember(ctx, id)
	if err != nil {
		s.logger.ErrorContext(ctx, "error getting member", slog.String("service.GetMember", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(&member); err != nil {
		s.logger.ErrorContext(ctx, "error encoding member", slog.String("service.GetMember", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Service) CreateMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var member db.Member
	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		s.logger.ErrorContext(ctx, "error decoding member", slog.String("service.CreateMember", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	create := db.CreateMemberParams{
		Name:  member.Name,
		Email: member.Email,
		Phone: member.Phone,
	}

	newMember, err := s.store.DB.CreateMember(ctx, create)
	if err != nil {
		s.logger.ErrorContext(ctx, "error creating member", slog.String("service.CreateMember", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(&newMember); err != nil {
		s.logger.ErrorContext(ctx, "error encoding member", slog.String("service.CreateMember", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Service) UpdateMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseInt32(chi.URLParam(r, "id"))
	if err != nil {
		s.logger.ErrorContext(ctx, "error parsing id", slog.String("service.UpdateMember", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var member db.Member
	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		s.logger.ErrorContext(ctx, "error decoding member", slog.String("service.UpdateMember", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	update := db.UpdateMemberParams{
		ID:    id,
		Name:  member.Name,
		Email: member.Email,
		Phone: member.Phone,
	}

	err = s.store.DB.UpdateMember(ctx, update)
	if err != nil {
		s.logger.ErrorContext(ctx, "error updating member", slog.String("service.UpdateMember", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) UpdateMembership(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseInt32(chi.URLParam(r, "id"))
	if err != nil {
		s.logger.ErrorContext(ctx, "error parsing id", slog.String("service.UpdateMembership", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var member db.Member
	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		s.logger.ErrorContext(ctx, "error decoding member", slog.String("service.UpdateMembership", err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	update := db.UpdateMembershipTypeParams{
		ID:             id,
		MembershipType: member.MembershipType,
	}

	err = s.store.DB.UpdateMembershipType(ctx, update)
	if err != nil {
		s.logger.ErrorContext(ctx, "error updating membership type", slog.String("service.UpdateMembership", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
