package handlers

import (
	"avito/tender/internal/domain"
	"strings"
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

func ConvertServiceTypeReqToServiceTypeDB(serviceType string) string {
	return strings.ToUpper(serviceType)
}

func ConvertTenderStatusReqToTenderStatusDB(status string) string {
	return strings.ToUpper(status)
}
