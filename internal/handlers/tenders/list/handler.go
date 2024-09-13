package tenders_list

import (
	"avito/tender/internal/domain"
	"context"
	"strings"
	"time"
)

type (
	repository interface {
		AddTender(ctx context.Context, tender domain.TenderAddDTO) (string, error)
		GetUserOrganizationId(ctx context.Context, username string) (string, string, error)
		GetTender(ctx context.Context, tenderId string) (domain.TenderAddDTO, error)
		GetTenderList(ctx context.Context, serviceTypes []string, limit int, offset int) ([]domain.TenderAddDTO, error)
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

func (h *Handler) ListTender(ctx context.Context, serviceType []string, limit int, offset int) ([]domain.TenderAddResponse, error) {
	servDBTypes := make([]string, 0, len(serviceType))
	for _, servType := range serviceType {
		servDBTypes = append(servDBTypes, strings.ToUpper(servType))
	}

	tenders, err := h.repo.GetTenderList(ctx, servDBTypes, limit, offset)
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
