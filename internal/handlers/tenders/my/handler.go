package tenders_my

import (
	"avito/tender/internal/domain"
	"avito/tender/internal/handlers"
	"context"
)

type (
	repository interface {
		GetUserOrganizationID(ctx context.Context, username string) (string, string, error)
		GetUsersTenders(ctx context.Context, username string, limit int, offset int) ([]domain.TenderDTO, error)
	}

	// Handler struct for MyTenders.
	Handler struct {
		repo repository
	}
)

// New returns new MyTenders handler.
func New(repo repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

// MyTenders returns user's tenders.
func (h *Handler) MyTenders(ctx context.Context, username string, limit int, offset int) ([]domain.TenderResponse, error) {
	uid, _, err := h.repo.GetUserOrganizationID(ctx, username)
	if err != nil {
		return nil, err
	}

	tenders, err := h.repo.GetUsersTenders(ctx, uid, limit, offset)
	if err != nil {
		return nil, err
	}

	tendersResp := make([]domain.TenderResponse, 0, len(tenders))
	for _, tenderDB := range tenders {
		tenderResp := handlers.ConvertTenderDTOToTenderResponse(tenderDB)
		tendersResp = append(tendersResp, tenderResp)
	}

	return tendersResp, nil
}
