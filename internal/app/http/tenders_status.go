package http

import (
	app_errors "avito/tender/internal/errors"
	"context"
	"fmt"
	"net/http"
)

type (
	statusTenderCommand interface {
		StatusTender(ctx context.Context, username string, tenderId string) (string, error)
	}

	StatusHandler struct {
		name                string
		statusTenderCommand statusTenderCommand
	}
)

func NewTendersStatusHandler(command statusTenderCommand, name string) *StatusHandler {
	return &StatusHandler{
		name:                name,
		statusTenderCommand: command,
	}
}

func (h *StatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		err error
	)

	queryFromURL := r.URL.Query()
	username := queryFromURL.Get("username")
	tenderId := r.PathValue("tenderId")

	if username == "" {
		GetErrorResponse(w, h.name, app_errors.ErrInvalidUser, http.StatusUnauthorized)
		return
	}

	resp, err := h.statusTenderCommand.StatusTender(
		ctx,
		username,
		tenderId,
	)
	if err != nil {
		switch err {
		case app_errors.ErrInvalidUser:
			GetErrorResponse(w, h.name, err, http.StatusUnauthorized)
			return
		case app_errors.ErrUserPermissions:
			GetErrorResponse(w, h.name, err, http.StatusForbidden)
			return
		case app_errors.ErrTenderId:
			GetErrorResponse(w, h.name, err, http.StatusNotFound)
			return
		default:
			fmt.Println(err)
			GetErrorResponse(w, h.name, app_errors.ErrInternalServer, http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(resp))
}
