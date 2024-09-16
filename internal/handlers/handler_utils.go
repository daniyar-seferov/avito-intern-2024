package handlers

import (
	"avito/tender/internal/domain"
	"time"
)

func ConvertTenderDTOToTenderResponse(tender domain.TenderDTO) domain.TenderResponse {
	return domain.TenderResponse{
		ID:          tender.ID,
		Name:        tender.Name,
		Description: tender.Description,
		Status:      domain.TenderStatusMap[tender.Status],
		ServiceType: domain.ServiceTypeMap[tender.ServiceType],
		Version:     tender.Version,
		CreatedAt:   tender.CreatedAt.Format(time.RFC3339),
	}
}
