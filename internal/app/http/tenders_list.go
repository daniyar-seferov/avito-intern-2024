package http

import (
	"avito/tender/internal/domain"
	app_errors "avito/tender/internal/errors"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"gopkg.in/validator.v2"
)

type (
	listTenderCommand interface {
		ListTender(ctx context.Context, serviceType []string, limit int, offset int) ([]domain.TenderResponse, error)
	}

	// ListHandler tenders' list struct.
	ListHandler struct {
		name              string
		listTenderCommand listTenderCommand
	}
)

// NewTendersListHandler returns new ListHandler.
func NewTendersListHandler(command listTenderCommand, name string) *ListHandler {
	return &ListHandler{
		name:              name,
		listTenderCommand: command,
	}
}

func (h *ListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		ctx           = r.Context()
		err           error
		limit, offset int64
	)

	queryFromURL := r.URL.Query()
	serviceTypeArr := queryFromURL["service_type"]
	limitQuery := queryFromURL.Get("limit")
	offsetQuery := queryFromURL.Get("offset")

	for _, serviceType := range serviceTypeArr {
		if err = validator.Valid(serviceType, "servicetype"); err != nil {
			GetErrorResponse(w, h.name, err, http.StatusBadRequest)
			return
		}
	}

	if limitQuery != "" {
		limit, err = strconv.ParseInt(limitQuery, 10, 32)
		if err != nil || limit <= 0 {
			GetErrorResponse(w, h.name, fmt.Errorf("invalid limit"), http.StatusBadRequest)
			return
		}
	}

	if offsetQuery != "" {
		offset, err = strconv.ParseInt(offsetQuery, 10, 32)
		if err != nil || offset < 0 {
			GetErrorResponse(w, h.name, fmt.Errorf("invalid offset"), http.StatusBadRequest)
			return
		}
	}

	resp, err := h.listTenderCommand.ListTender(
		ctx,
		serviceTypeArr,
		int(limit),
		int(offset),
	)
	if err != nil {
		fmt.Println(err)
		GetErrorResponse(w, h.name, app_errors.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	respData, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("response marshal error: %v", err)
		GetErrorResponse(w, h.name, app_errors.ErrInternalServer, http.StatusInternalServerError)
	}

	GetSuccessResponseWithBody(w, respData)
}
