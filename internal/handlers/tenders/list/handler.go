package tenders_list

import (
	"avito/tender/internal/domain"
	"avito/tender/internal/handlers"
	"context"
	"strings"
)

type (
	repository interface {
		GetTenderList(ctx context.Context, serviceTypes []string, limit int, offset int) ([]domain.TenderDTO, error)
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

func (h *Handler) ListTender(ctx context.Context, serviceType []string, limit int, offset int) ([]domain.TenderResponse, error) {
	servDBTypes := make([]string, 0, len(serviceType))
	for _, servType := range serviceType {
		servDBTypes = append(servDBTypes, strings.ToUpper(servType))
	}

	tenders, err := h.repo.GetTenderList(ctx, servDBTypes, limit, offset)
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
