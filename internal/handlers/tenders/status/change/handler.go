package tenders_new

import (
	"avito/tender/internal/domain"
	app_errors "avito/tender/internal/errors"
	"avito/tender/internal/handlers"
	"context"
)

type (
	repository interface {
		GetUserOrganizationId(ctx context.Context, username string) (string, string, error)
		GetTender(ctx context.Context, tenderId string) (domain.TenderDTO, error)
		UpdateTender(ctx context.Context, tenderId string, tenderDTO domain.TenderDTO) (domain.TenderDTO, error)
	}

	Handler struct {
		repo repository
	}
)

func New(repo repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) ChangeStatusTender(ctx context.Context, username string, tenderId string, status string) (domain.TenderResponse, error) {
	var tenderResp domain.TenderResponse
	uid, organizationID, err := h.repo.GetUserOrganizationId(ctx, username)
	if err != nil {
		return tenderResp, err
	}

	tenderDB, err := h.repo.GetTender(ctx, tenderId)
	if err != nil {
		return tenderResp, err
	}

	if uid != tenderDB.UserId || organizationID != tenderDB.OrganizationId {
		return tenderResp, app_errors.ErrUserPermissions
	}

	tenderDTO := domain.TenderDTO{
		Status: handlers.ConvertTenderStatusReqToTenderStatusDB(status),
	}
	tenderDB, err = h.repo.UpdateTender(ctx, tenderId, tenderDTO)
	if err != nil {
		return tenderResp, err
	}

	tenderResp = handlers.ConvertTenderDTOToTenderResponse(tenderDB)

	return tenderResp, nil
}
