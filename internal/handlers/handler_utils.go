package handlers

import (
	"avito/tender/internal/domain"
	"strings"
	"time"
)

// ConvertTenderDTOToTenderResponse converts TenderDTO to TenderResponse.
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

// ConvertServiceTypeReqToServiceTypeDB converts ServiceTypeReq to ServiceTypeDB.
func ConvertServiceTypeReqToServiceTypeDB(serviceType string) string {
	return strings.ToUpper(serviceType)
}

// ConvertTenderStatusReqToTenderStatusDB converts TenderStatusReq to TenderStatusDB.
func ConvertTenderStatusReqToTenderStatusDB(status string) string {
	return strings.ToUpper(status)
}
