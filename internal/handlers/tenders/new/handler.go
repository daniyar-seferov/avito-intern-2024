package tenders_new

import (
	"avito/tender/internal/domain"
	app_errors "avito/tender/internal/errors"
	"avito/tender/internal/handlers"
	"context"
	"fmt"
)

type (
	repository interface {
		AddTender(ctx context.Context, tender domain.TenderDTO) (string, error)
		GetUserOrganizationID(ctx context.Context, username string) (string, string, error)
		GetTender(ctx context.Context, tenderID string) (domain.TenderDTO, error)
	}

	// Handler struct for AddTender.
	Handler struct {
		repo repository
	}
)

// New returns new AddTender handler.
func New(repo repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

// AddTender adds new tender.
func (h *Handler) AddTender(ctx context.Context, tender domain.TenderAddRequest) (domain.TenderResponse, error) {
	uid, organizationID, err := h.repo.GetUserOrganizationID(ctx, tender.CreatorUsername)
	if err != nil {
		return domain.TenderResponse{}, err
	}

	if organizationID != tender.OrganizationID {
		return domain.TenderResponse{}, app_errors.ErrInvalidOrganization
	}

	tenderDTO := domain.TenderDTO{
		Name:           tender.Name,
		Description:    tender.Description,
		ServiceType:    handlers.ConvertServiceTypeReqToServiceTypeDB(tender.ServiceType),
		OrganizationID: tender.OrganizationID,
		UserID:         uid,
	}

	tenderID, err := h.repo.AddTender(ctx, tenderDTO)
	if err != nil {
		return domain.TenderResponse{}, fmt.Errorf("repo.AddTender failed: %w", err)
	}

	tenderDB, err := h.repo.GetTender(ctx, tenderID)
	if err != nil {
		return domain.TenderResponse{}, fmt.Errorf("repo.GetTender failed: %w", err)
	}

	tenderResp := handlers.ConvertTenderDTOToTenderResponse(tenderDB)

	return tenderResp, nil
}
