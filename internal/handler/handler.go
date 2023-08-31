package handler

import (
	"github.com/diianpro/stellar-scope/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"time"
)

type Handler struct {
	srv service.Picture
}

func New(srv *service.Service) *Handler {
	return &Handler{
		srv: srv,
	}
}

func (h Handler) GetByDate(c echo.Context, data time.Time) error {
	picture, err := h.srv.GetByDate(c.Request().Context(), data)
	if err != nil {
		log.Errorf("GetByDate: picture error: %v", err)
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, picture)
}

func (h Handler) GetAll(c echo.Context) error {
	pictures, err := h.srv.GetAll(c.Request().Context())
	if err != nil {
		log.Errorf("GetAll: pictures error: %v", err)
		return echo.NewHTTPError(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusOK, pictures)
}
