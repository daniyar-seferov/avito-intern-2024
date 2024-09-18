package http

import (
	"avito/tender/internal/domain"
	"avito/tender/internal/handlers"
	"errors"
	"reflect"

	"gopkg.in/validator.v2"
)

func ValidateServiceType(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	if st.Kind() != reflect.String {
		return validator.ErrUnsupported
	}

	serviceType := handlers.ConvertServiceTypeReqToServiceTypeDB(st.String())
	if _, inMap := domain.ServiceTypeMap[serviceType]; !inMap {
		return errors.New("invalid service type")
	}

	return nil
}

func ValidateTenderStatus(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	if st.Kind() != reflect.String {
		return validator.ErrUnsupported
	}

	serviceType := handlers.ConvertTenderStatusReqToTenderStatusDB(st.String())
	if _, inMap := domain.TenderStatusMap[serviceType]; !inMap {
		return errors.New("invalid status")
	}

	return nil
}
