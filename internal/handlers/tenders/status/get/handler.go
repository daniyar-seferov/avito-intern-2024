package tenders_new

import (
	"avito/tender/internal/domain"
	app_errors "avito/tender/internal/errors"
	"context"
)

type (
	repository interface {
		GetUserOrganizationID(ctx context.Context, username string) (string, string, error)
		GetTender(ctx context.Context, tenderID string) (domain.TenderDTO, error)
	}

	// Handler struct for GET StatusTender.
	Handler struct {
		repo repository
	}
)

// New returns GET StatusTender handler.
func New(repo repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

// StatusTender returns tender status.
func (h *Handler) StatusTender(ctx context.Context, username string, tenderID string) (string, error) {
	uid, organizationID, err := h.repo.GetUserOrganizationID(ctx, username)
	if err != nil {
		return "", err
	}

	tenderDB, err := h.repo.GetTender(ctx, tenderID)
	if err != nil {
		return "", err
	}

	if uid != tenderDB.UserID || organizationID != tenderDB.OrganizationID {
		return "", app_errors.ErrUserPermissions
	}

	return domain.TenderStatusMap[tenderDB.Status], nil
}
