package http

import (
	"avito/tender/internal/domain"
	app_errors "avito/tender/internal/errors"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type (
	myTenderCommand interface {
		MyTenders(ctx context.Context, username string, limit int, offset int) ([]domain.TenderResponse, error)
	}

	MyHandler struct {
		name            string
		myTenderCommand myTenderCommand
	}
)

func NewTendersMyHandler(command myTenderCommand, name string) *MyHandler {
	return &MyHandler{
		name:            name,
		myTenderCommand: command,
	}
}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		ctx           = r.Context()
		err           error
		limit, offset int64
	)

	queryFromURL := r.URL.Query()
	username := queryFromURL.Get("username")
	limitQuery := queryFromURL.Get("limit")
	offsetQuery := queryFromURL.Get("offset")

	if username == "" {
		GetErrorResponse(w, h.name, app_errors.ErrInvalidUser, http.StatusBadRequest)
		return
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

	resp, err := h.myTenderCommand.MyTenders(
		ctx,
		username,
		int(limit),
		int(offset),
	)
	if err != nil {
		switch err {
		case app_errors.ErrInvalidUser:
			GetErrorResponse(w, h.name, err, http.StatusUnauthorized)
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
