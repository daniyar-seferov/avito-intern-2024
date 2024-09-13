package tenders_new

import (
	"avito/tender/internal/domain"
	app_errors "avito/tender/internal/errors"
	"context"
)

type (
	repository interface {
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

func (h *Handler) StatusTender(ctx context.Context, username string, tenderId string) (string, error) {
	uid, organization_id, err := h.repo.GetUserOrganizationId(ctx, username)
	if err != nil {
		return "", err
	}

	tenderDB, err := h.repo.GetTender(ctx, tenderId)
	if err != nil {
		return "", err
	}

	if uid != tenderDB.UserId || organization_id != tenderDB.OrganizationId {
		return "", app_errors.ErrUserPermissions
	}

	return domain.TenderStatusMap[tenderDB.Status], nil
}
