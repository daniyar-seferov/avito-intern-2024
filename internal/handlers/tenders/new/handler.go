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
		GetUserOrganizationId(ctx context.Context, username string) (string, string, error)
		GetTender(ctx context.Context, tenderId string) (domain.TenderDTO, error)
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

func (h *Handler) AddTender(ctx context.Context, tender domain.TenderAddRequest) (domain.TenderResponse, error) {
	uid, organizationId, err := h.repo.GetUserOrganizationId(ctx, tender.CreatorUsername)
	if err != nil {
		return domain.TenderResponse{}, err
	}

	if organizationId != tender.OrganizationId {
		return domain.TenderResponse{}, app_errors.ErrInvalidOrganization
	}

	tenderDTO := domain.TenderDTO{
		Name:           tender.Name,
		Description:    tender.Description,
		ServiceType:    handlers.ConvertServiceTypeReqToServiceTypeDB(tender.ServiceType),
		OrganizationId: tender.OrganizationId,
		UserId:         uid,
	}

	tenderId, err := h.repo.AddTender(ctx, tenderDTO)
	if err != nil {
		return domain.TenderResponse{}, fmt.Errorf("repo.AddTender failed: %w", err)
	}

	tenderDB, err := h.repo.GetTender(ctx, tenderId)
	if err != nil {
		return domain.TenderResponse{}, fmt.Errorf("repo.GetTender failed: %w", err)
	}

	tenderResp := handlers.ConvertTenderDTOToTenderResponse(tenderDB)

	return tenderResp, nil
}
