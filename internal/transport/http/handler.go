package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/labstack/gommon/log"

	"github.com/diianpro/stellar-scope/internal/service"
	"github.com/diianpro/stellar-scope/internal/transport/http/response"
)

const (
	limitQueryParam  = "limit"
	offsetQueryParam = "offset"
	dateParamKey     = "date"
)

type Handler struct {
	srv service.Picture
}

func NewHandler(srv service.Picture) *Handler {
	return &Handler{
		srv: srv,
	}
}

func (h Handler) GetByDate(w http.ResponseWriter, r *http.Request) {
	d := chi.URLParam(r, dateParamKey)
	date, err := time.Parse(time.DateOnly, d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	picture, err := h.srv.GetByDate(r.Context(), date)
	if err != nil {
		log.Errorf("GetByDate: picture error: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err = encodeJSONResponse(w, http.StatusOK, picture); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	var err error
	var resp response.GetAll

	limit, err := getIntQueryParam(r, limitQueryParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	offset, err := getIntQueryParam(r, offsetQueryParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp.Pictures, err = h.srv.GetAll(r.Context(), limit, offset)
	if err != nil {
		log.Errorf("GetAll: pictures error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = encodeJSONResponse(w, http.StatusOK, resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
