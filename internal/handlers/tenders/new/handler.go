package tenders_new

import (
	"avito/tender/internal/domain"
	app_errors "avito/tender/internal/errors"
	"context"
	"fmt"
	"strings"
	"time"
)

type (
	repository interface {
		AddTender(ctx context.Context, tender domain.TenderAddDTO) (string, error)
		GetUserOrganizationId(ctx context.Context, username string) (string, string, error)
		GetTender(ctx context.Context, tenderId string) (domain.TenderAddDTO, error)
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

func (h *Handler) AddTender(ctx context.Context, tender domain.TenderAddRequest) (domain.TenderAddResponse, error) {
	uid, organizationId, err := h.repo.GetUserOrganizationId(ctx, tender.CreatorUsername)
	if err != nil {
		return domain.TenderAddResponse{}, err
	}

	if organizationId != tender.OrganizationId {
		return domain.TenderAddResponse{}, app_errors.ErrInvalidOrganization
	}

	tenderDTO := domain.TenderAddDTO{
		Name:           tender.Name,
		Description:    tender.Description,
		ServiceType:    strings.ToUpper(tender.ServiceType),
		OrganizationId: tender.OrganizationId,
		UserId:         uid,
	}

	tenderId, err := h.repo.AddTender(ctx, tenderDTO)
	if err != nil {
		return domain.TenderAddResponse{}, fmt.Errorf("repo.AddTender failed: %w", err)
	}

	tenderDB, err := h.repo.GetTender(ctx, tenderId)
	if err != nil {
		return domain.TenderAddResponse{}, fmt.Errorf("repo.GetTender failed: %w", err)
	}

	tenderResp := domain.TenderAddResponse{
		ID:          tenderDB.ID,
		Name:        tenderDB.Name,
		Description: tenderDB.Description,
		Status:      domain.TenderStatusMap[tenderDB.Status],
		ServiceType: domain.ServiceTypeMap[tenderDB.ServiceType],
		Version:     tenderDB.Version,
		CreatedAt:   tenderDB.CreatedAt.Format(time.RFC3339),
	}

	return tenderResp, nil
}
