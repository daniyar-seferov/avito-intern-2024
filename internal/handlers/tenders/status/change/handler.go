package tenders_new

import (
	"avito/tender/internal/domain"
	app_errors "avito/tender/internal/errors"
	"avito/tender/internal/handlers"
	"context"
)

type (
	repository interface {
		GetUserOrganizationID(ctx context.Context, username string) (string, string, error)
		GetTender(ctx context.Context, tenderID string) (domain.TenderDTO, error)
		UpdateTender(ctx context.Context, tenderID string, tenderDTO domain.TenderDTO) (domain.TenderDTO, error)
	}

	// Handler struct for ChangeStatusTender.
	Handler struct {
		repo repository
	}
)

// New returns new ChangeStatusTender handler.
func New(repo repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

// ChangeStatusTender changes tender status.
func (h *Handler) ChangeStatusTender(ctx context.Context, username string, tenderID string, status string) (domain.TenderResponse, error) {
	var tenderResp domain.TenderResponse
	uid, organizationID, err := h.repo.GetUserOrganizationID(ctx, username)
	if err != nil {
		return tenderResp, err
	}

	tenderDB, err := h.repo.GetTender(ctx, tenderID)
	if err != nil {
		return tenderResp, err
	}

	if uid != tenderDB.UserID || organizationID != tenderDB.OrganizationID {
		return tenderResp, app_errors.ErrUserPermissions
	}

	tenderDTO := domain.TenderDTO{
		Status: handlers.ConvertTenderStatusReqToTenderStatusDB(status),
	}
	tenderDB, err = h.repo.UpdateTender(ctx, tenderID, tenderDTO)
	if err != nil {
		return tenderResp, err
	}

	tenderResp = handlers.ConvertTenderDTOToTenderResponse(tenderDB)

	return tenderResp, nil
}
