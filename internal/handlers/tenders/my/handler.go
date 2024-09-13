package tenders_my

import (
	"avito/tender/internal/domain"
	"context"
	"time"
)

type (
	repository interface {
		GetUserOrganizationId(ctx context.Context, username string) (string, string, error)
		GetUsersTenders(ctx context.Context, username string, limit int, offset int) ([]domain.TenderAddDTO, error)
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

func (h *Handler) MyTenders(ctx context.Context, username string, limit int, offset int) ([]domain.TenderAddResponse, error) {
	uid, _, err := h.repo.GetUserOrganizationId(ctx, username)
	if err != nil {
		return nil, err
	}

	tenders, err := h.repo.GetUsersTenders(ctx, uid, limit, offset)
	if err != nil {
		return nil, err
	}

	tendersResp := make([]domain.TenderAddResponse, 0, len(tenders))
	for _, tenderDB := range tenders {
		tenderResp := domain.TenderAddResponse{
			ID:          tenderDB.ID,
			Name:        tenderDB.Name,
			Description: tenderDB.Description,
			Status:      domain.TenderStatusMap[tenderDB.Status],
			ServiceType: domain.ServiceTypeMap[tenderDB.ServiceType],
			Version:     tenderDB.Version,
			CreatedAt:   tenderDB.CreatedAt.Format(time.RFC3339),
		}
		tendersResp = append(tendersResp, tenderResp)
	}

	return tendersResp, nil
}
