package domain

import "time"

type (
	// TenderAddRequest struct.
	TenderAddRequest struct {
		Name            string `json:"name" validate:"nonzero"`
		Description     string `json:"description" validate:"nonzero"`
		ServiceType     string `json:"serviceType" validate:"nonzero, servicetype"`
		OrganizationID  string `json:"organizationId" validate:"nonzero"`
		CreatorUsername string `json:"creatorUsername" validate:"nonzero"`
	}

	// TenderEditRequest tenders edit request.
	TenderEditRequest struct {
		Name        string `json:"name" validate:"nonzero"`
		Description string `json:"description" validate:"nonzero"`
		ServiceType string `json:"serviceType" validate:"nonzero, servicetype"`
	}

	// TenderResponse struct.
	TenderResponse struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Status      string `json:"status"`
		ServiceType string `json:"serviceType"`
		Version     int    `json:"version"`
		CreatedAt   string `json:"created_at"`
	}

	// TenderDTO struct.
	TenderDTO struct {
		ID             string
		Name           string
		Description    string
		Status         string
		ServiceType    string
		OrganizationID string
		UserID         string
		Version        int
		CreatedAt      time.Time
		UpdatedAt      time.Time
	}
)

var (
	// ServiceTypeMap matches service type from DB to request service type.
	ServiceTypeMap = map[string]string{
		"CONSTRUCTION": "Construction",
		"DELIVERY":     "Delivery",
		"MANUFACTURE":  "Manufacture",
	}

	// TenderStatusMap matches tender status from DB to request tender status.
	TenderStatusMap = map[string]string{
		"CREATED":   "Created",
		"PUBLISHED": "Published",
		"CLOSED":    "Closed",
	}
)
