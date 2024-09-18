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
		GetUserOrganizationID(ctx context.Context, username string) (string, string, error)
		GetTender(ctx context.Context, tenderID string) (domain.TenderDTO, error)
		UpdateTender(ctx context.Context, tenderID string, tenderDTO domain.TenderDTO) (domain.TenderDTO, error)
	}

	// Handler struct for EditTender handler.
	Handler struct {
		repo repository
	}
)

// New returns new EditTender handler.
func New(repo repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

// EditTender edits tender.
func (h *Handler) EditTender(ctx context.Context, username string, tenderID string, tender domain.TenderEditRequest) (domain.TenderResponse, error) {
	uid, organizationID, err := h.repo.GetUserOrganizationID(ctx, username)
	if err != nil {
		return domain.TenderResponse{}, err
	}

	tenderDB, err := h.repo.GetTender(ctx, tenderID)
	if err != nil {
		return domain.TenderResponse{}, fmt.Errorf("repo.EditTender failed: %w", err)
	}

	if uid != tenderDB.UserID {
		return domain.TenderResponse{}, app_errors.ErrUserPermissions
	}
	if organizationID != tenderDB.OrganizationID {
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

	tenderUpd, err := h.repo.UpdateTender(ctx, tenderID, tenderDTO)
	if err != nil {
		return domain.TenderResponse{}, fmt.Errorf("repo.EditTender failed: %w", err)
	}

	tenderResp := handlers.ConvertTenderDTOToTenderResponse(tenderUpd)

	return tenderResp, nil
}
