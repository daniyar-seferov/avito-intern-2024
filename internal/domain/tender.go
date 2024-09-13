package domain

import "time"

type (
	TenderAddRequest struct {
		Name            string `json:"name" validate:"nonzero"`
		Description     string `json:"description" validate:"nonzero"`
		ServiceType     string `json:"serviceType" validate:"nonzero"`
		OrganizationId  string `json:"organizationId" validate:"nonzero"`
		CreatorUsername string `json:"creatorUsername" validate:"nonzero"`
	}

	TenderAddResponse struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Status      string `json:"status"`
		ServiceType string `json:"serviceType"`
		Version     int    `json:"version"`
		CreatedAt   string `json:"created_at"`
	}

	TenderAddDTO struct {
		ID             string
		Name           string
		Description    string
		Status         string
		ServiceType    string
		OrganizationId string
		UserId         string
		Version        int
		CreatedAt      time.Time
	}
)

var (
	ServiceTypeMap = map[string]string{
		"CONSTRUCTION": "Construction",
		"DELIVERY":     "Delivery",
		"MANUFACTURE":  "Manufacture",
	}

	TenderStatusMap = map[string]string{
		"CREATED":   "Created",
		"PUBLISHED": "Published",
		"CLOSED":    "Closed",
	}
)
