package tenders_edit

import (
	"avito/tender/internal/domain"
	app_errors "avito/tender/internal/errors"
	"avito/tender/internal/handlers"
	"context"
	"fmt"
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

func (h *Handler) EditTender(ctx context.Context, username string, tenderId string, tender domain.TenderEditRequest) (domain.TenderResponse, error) {
	uid, organizationId, err := h.repo.GetUserOrganizationId(ctx, username)
	if err != nil {
		return domain.TenderResponse{}, err
	}

	tenderDB, err := h.repo.GetTender(ctx, tenderId)
	if err != nil {
		return domain.TenderResponse{}, fmt.Errorf("repo.EditTender failed: %w", err)
	}

	if uid != tenderDB.UserId {
		return domain.TenderResponse{}, app_errors.ErrUserPermissions
	}
	if organizationId != tenderDB.OrganizationId {
		return domain.TenderResponse{}, app_errors.ErrInvalidOrganization
	}

	tenderDTO := domain.TenderDTO{}
	if tender.Name != "" {
		tenderDTO.Name = tender.Name
	}
	if tender.Description != "" {
		tenderDTO.Description = tender.Description
	}
	if tender.ServiceType != "" {
		tenderDTO.ServiceType = handlers.ConvertServiceTypeReqToServiceTypeDB(tender.ServiceType)
	}

	tenderUpd, err := h.repo.UpdateTender(ctx, tenderId, tenderDTO)
	if err != nil {
		return domain.TenderResponse{}, fmt.Errorf("repo.EditTender failed: %w", err)
	}

	tenderResp := handlers.ConvertTenderDTOToTenderResponse(tenderUpd)

	return tenderResp, nil
}
