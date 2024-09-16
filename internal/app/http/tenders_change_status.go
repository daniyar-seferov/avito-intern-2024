package http

import (
	"avito/tender/internal/domain"
	app_errors "avito/tender/internal/errors"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type (
	changeStatusTenderCommand interface {
		ChangeStatusTender(ctx context.Context, username, tenderId, status string) (domain.TenderResponse, error)
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
		GetErrorResponse(w, h.name, app_errors.ErrInvalidUser, http.StatusBadRequest)
		return
	}
	if _, inMap := domain.TenderStatusMap[strings.ToUpper(status)]; !inMap {
		GetErrorResponse(w, h.name, app_errors.ErrInvalidStatus, http.StatusBadRequest)
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
		case app_errors.ErrInvalidTenderId:
			GetErrorResponse(w, h.name, err, http.StatusNotFound)
			return
		default:
			fmt.Println(err)
			GetErrorResponse(w, h.name, app_errors.ErrInternalServer, http.StatusInternalServerError)
			return
		}
	}

	respData, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("response marshal error: %v", err)
		GetErrorResponse(w, h.name, app_errors.ErrInternalServer, http.StatusInternalServerError)
	}

	GetSuccessResponseWithBody(w, respData)
}
