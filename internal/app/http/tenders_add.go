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
	addTenderCommand interface {
		AddTender(ctx context.Context, tender domain.TenderAddRequest) (domain.TenderResponse, error)
	}

	// AddHandler add handler struct.
	AddHandler struct {
		name             string
		addTenderCommand addTenderCommand
	}
)

// NewTendersAddHandler returns new AddHandler.
func NewTendersAddHandler(command addTenderCommand, name string) *AddHandler {
	return &AddHandler{
		name:             name,
		addTenderCommand: command,
	}
}

func (h *AddHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		ctx     = r.Context()
		request domain.TenderAddRequest
		err     error
	)

	request = domain.TenderAddRequest{}
	err = json.NewDecoder(r.Body).Decode(&request)
	if err == io.EOF {
		GetErrorResponse(w, h.name, fmt.Errorf("invalid json"), http.StatusBadRequest)
		return
	}
	if err != nil {
		GetErrorResponse(w, h.name, err, http.StatusBadRequest)
		return
	}
	_ = r.Body.Close()

	if err = validator.Validate(request); err != nil {
		GetErrorResponse(w, h.name, err, http.StatusBadRequest)
		return
	}

	resp, err := h.addTenderCommand.AddTender(
		ctx,
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
