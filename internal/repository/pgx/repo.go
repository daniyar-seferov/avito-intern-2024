package db_pgx_repo

import (
	"avito/tender/internal/domain"
	app_errors "avito/tender/internal/errors"
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

type Repo struct {
	conn *pgx.Conn
}

func NewRepo(conn *pgx.Conn) *Repo {
	return &Repo{
		conn: conn,
	}
}

func (r *Repo) AddTender(ctx context.Context, tender domain.TenderAddDTO) (string, error) {
	const query = `
	INSERT INTO tender (organization_id, user_id, name, description, type)
	VALUES ($1, $2, $3, $4, $5) RETURNING id;`

	var tenderId string

	err := r.conn.QueryRow(ctx, query, tender.OrganizationId, tender.UserId, tender.Name, tender.Description, tender.ServiceType).Scan(&tenderId)
	if err != nil {
		return "", err
	}

	return tenderId, nil
}

func (r *Repo) GetUserOrganizationId(ctx context.Context, username string) (string, string, error) {
	const query = `
	SELECT employee.id, organization_responsible.organization_id FROM employee 
	LEFT JOIN organization_responsible ON employee.id=organization_responsible.user_id
	WHERE employee.username=$1;`

	var uid, organizationId string

	err := r.conn.QueryRow(ctx, query, username).Scan(&uid, &organizationId)
	if err == pgx.ErrNoRows {
		return "", "", app_errors.ErrInvalidUser
	}
	if err != nil {
		return "", "", err
	}

	return uid, organizationId, nil
}

func (r *Repo) GetTender(ctx context.Context, tenderId string) (domain.TenderAddDTO, error) {
	const query = `
	SELECT name, description, status, type, version, created_at 
	FROM tender WHERE id=$1;`

	tender := domain.TenderAddDTO{}

	err := r.conn.QueryRow(ctx, query, tenderId).Scan(&tender.Name, &tender.Description, &tender.Status, &tender.ServiceType, &tender.Version, &tender.CreatedAt)
	if err != nil {
		return tender, err
	}
	tender.ID = tenderId

	return tender, nil
}

func (r *Repo) GetTenderList(ctx context.Context, serviceTypes []string, limit int, offset int) ([]domain.TenderAddDTO, error) {
	var query = `
	SELECT id, name, description, status, type, version, created_at 
	FROM tender`

	var (
		args    []interface{}
		tenders []domain.TenderAddDTO
	)

	if len(serviceTypes) > 0 {
		placeholders := make([]string, len(serviceTypes))
		for pos, serviceType := range serviceTypes {
			placeholders[pos] = fmt.Sprintf("$%d", pos+1)
			args = append(args, serviceType)
		}
		query += " WHERE type IN (" + strings.Join(placeholders, ", ") + ")"
	}

	if limit > 0 {
		query += " LIMIT $" + fmt.Sprint(len(args)+1)
		args = append(args, limit)
	}
	if offset > 0 {
		query += " OFFSET $" + fmt.Sprint(len(args)+1)
		args = append(args, offset)
	}

	rows, err := r.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		tender := domain.TenderAddDTO{}
		err := rows.Scan(&tender.ID, &tender.Name, &tender.Description, &tender.Status, &tender.ServiceType, &tender.Version, &tender.CreatedAt)
		if err != nil {
			fmt.Printf("GetTenderList rows.Scan error: %v", err)
			continue
		}
		tenders = append(tenders, tender)
	}
	rows.Close()

	return tenders, nil
}

func (r *Repo) GetUsersTenders(ctx context.Context, uid string, limit int, offset int) ([]domain.TenderAddDTO, error) {
	var query = `
	SELECT id, name, description, status, type, version, created_at 
	FROM tender 
	WHERE user_id=$1`

	var (
		args    []interface{}
		tenders []domain.TenderAddDTO
	)
	args = append(args, uid)

	if limit > 0 {
		query += " LIMIT $" + fmt.Sprint(len(args)+1)
		args = append(args, limit)
	}
	if offset > 0 {
		query += " OFFSET $" + fmt.Sprint(len(args)+1)
		args = append(args, offset)
	}

	rows, err := r.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		tender := domain.TenderAddDTO{}
		err := rows.Scan(&tender.ID, &tender.Name, &tender.Description, &tender.Status, &tender.ServiceType, &tender.Version, &tender.CreatedAt)
		if err != nil {
			fmt.Printf("GetUsersTenders rows.Scan error: %v", err)
			continue
		}
		tenders = append(tenders, tender)
	}
	rows.Close()

	return tenders, nil
}
