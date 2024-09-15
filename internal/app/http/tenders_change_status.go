package http

import (
	app_errors "avito/tender/internal/errors"
	"context"
	"fmt"
	"net/http"
)

type (
	changeStatusTenderCommand interface {
		ChangeStatusTender(ctx context.Context, username, tenderId, status string) (string, error)
	}

	ChangeStatusHandler struct {
		name                      string
		changeStatusTenderCommand changeStatusTenderCommand
	}
)

func NewTendersChangeStatusHandler(command changeStatusTenderCommand, name string) *ChangeStatusHandler {
	return &ChangeStatusHandler{
		name:                      name,
		changeStatusTenderCommand: command,
	}
}

func (h *ChangeStatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		err error
	)

	queryFromURL := r.URL.Query()
	username := queryFromURL.Get("username")
	status := queryFromURL.Get("status")
	tenderId := r.PathValue("tenderId")

	if username == "" {
		GetErrorResponse(w, h.name, app_errors.ErrInvalidUser, http.StatusUnauthorized)
		return
	}

	resp, err := h.changeStatusTenderCommand.ChangeStatusTender(
		ctx,
		username,
		tenderId,
		status,
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
