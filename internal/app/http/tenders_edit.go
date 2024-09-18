package http

import (
	"avito/tender/internal/domain"
	app_errors "avito/tender/internal/errors"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gopkg.in/validator.v2"
)

type (
	editTenderCommand interface {
		EditTender(ctx context.Context, username string, tenderId string, tender domain.TenderEditRequest) (domain.TenderResponse, error)
	}

	EditHandler struct {
		name              string
		editTenderCommand editTenderCommand
	}
)

func NewTendersEditHandler(command editTenderCommand, name string) *EditHandler {
	return &EditHandler{
		name:              name,
		editTenderCommand: command,
	}
}

func (h *EditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		ctx     = r.Context()
		request domain.TenderEditRequest
		err     error
	)

	queryFromURL := r.URL.Query()
	username := queryFromURL.Get("username")
	tenderId := r.PathValue("tenderId")

	if username == "" {
		GetErrorResponse(w, h.name, app_errors.ErrInvalidUser, http.StatusBadRequest)
		return
	}

	request = domain.TenderEditRequest{}
	err = json.NewDecoder(r.Body).Decode(&request)
	if err == io.EOF {
		GetErrorResponse(w, h.name, fmt.Errorf("invalid json"), http.StatusBadRequest)
		return
	}
	if err != nil {
		GetErrorResponse(w, h.name, err, http.StatusBadRequest)
		return
	}
	r.Body.Close()

	if err = validator.Validate(request); err != nil {
		GetErrorResponse(w, h.name, err, http.StatusBadRequest)
		return
	}

	resp, err := h.editTenderCommand.EditTender(
		ctx,
		username,
		tenderId,
		request,
	)
	if err != nil {
		switch err {
		case app_errors.ErrInvalidUser:
			GetErrorResponse(w, h.name, err, http.StatusUnauthorized)
			return
		case app_errors.ErrInvalidOrganization:
			GetErrorResponse(w, h.name, err, http.StatusForbidden)
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
