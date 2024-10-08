package db_pgx_repo

import (
	"avito/tender/internal/domain"
	app_errors "avito/tender/internal/errors"
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

// Repo struct contains db connection.
type Repo struct {
	conn *pgx.Conn
}

// NewRepo creates new Repo.
func NewRepo(conn *pgx.Conn) *Repo {
	return &Repo{
		conn: conn,
	}
}

// AddTender adds new tender.
func (r *Repo) AddTender(ctx context.Context, tender domain.TenderDTO) (string, error) {
	const query = `
	INSERT INTO tender (organization_id, user_id, name, description, type)
	VALUES ($1, $2, $3, $4, $5) RETURNING id;`

	var tenderID string

	err := r.conn.QueryRow(ctx, query, tender.OrganizationID, tender.UserID, tender.Name, tender.Description, tender.ServiceType).Scan(&tenderID)
	if err != nil {
		return "", err
	}

	return tenderID, nil
}

// GetUserOrganizationID gets user id and user's organization id.
func (r *Repo) GetUserOrganizationID(ctx context.Context, username string) (string, string, error) {
	const query = `
	SELECT employee.id, organization_responsible.organization_id FROM employee 
	LEFT JOIN organization_responsible ON employee.id=organization_responsible.user_id
	WHERE employee.username=$1;`

	var uid, organizationID string

	err := r.conn.QueryRow(ctx, query, username).Scan(&uid, &organizationID)
	if err == pgx.ErrNoRows {
		return "", "", app_errors.ErrInvalidUser
	}
	if err != nil {
		return "", "", err
	}

	return uid, organizationID, nil
}

// GetTender gets tender by tender id.
func (r *Repo) GetTender(ctx context.Context, tenderID string) (domain.TenderDTO, error) {
	const query = `
	SELECT organization_id, user_id, name, description, status, type, version, created_at, updated_at 
	FROM tender WHERE id=$1;`

	tender := domain.TenderDTO{}

	err := r.conn.QueryRow(ctx, query, tenderID).Scan(&tender.OrganizationID, &tender.UserID, &tender.Name, &tender.Description, &tender.Status, &tender.ServiceType, &tender.Version, &tender.CreatedAt, &tender.UpdatedAt)
	if err == pgx.ErrNoRows {
		return tender, app_errors.ErrInvalidTenderID
	}
	if err != nil {
		return tender, err
	}
	tender.ID = tenderID

	return tender, nil
}

// GetTenderList gets tenders.
func (r *Repo) GetTenderList(ctx context.Context, serviceTypes []string, limit int, offset int) ([]domain.TenderDTO, error) {
	var query = `
	SELECT id, name, description, status, type, version, created_at 
	FROM tender`

	var (
		args    []interface{}
		tenders []domain.TenderDTO
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
		tender := domain.TenderDTO{}
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

// GetUsersTenders gets user tenders.
func (r *Repo) GetUsersTenders(ctx context.Context, uid string, limit int, offset int) ([]domain.TenderDTO, error) {
	var query = `
	SELECT id, name, description, status, type, version, created_at 
	FROM tender 
	WHERE user_id=$1`

	var (
		args    []interface{}
		tenders []domain.TenderDTO
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
		tender := domain.TenderDTO{}
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

// InTx transaction helper.
func (r *Repo) InTx(ctx context.Context, f func(tx pgx.Tx) error) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			fmt.Printf("transaction rollback error: %v \n", err)
		}
	}(tx, ctx)

	err = f(tx)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// SetTenderRevision saves tender previous version.
func (r *Repo) SetTenderRevision(ctx context.Context, tender domain.TenderDTO) error {
	const query = `
	INSERT INTO tender_revision (tender_id, organization_id, user_id, name, description, status, type, version, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`

	_, err := r.conn.Exec(ctx, query, tender.ID, tender.OrganizationID, tender.UserID, tender.Name, tender.Description, tender.Status, tender.ServiceType, tender.Version, tender.CreatedAt, tender.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// UpdateTender transactional updates tender and saves old version.
func (r *Repo) UpdateTender(ctx context.Context, tenderID string, tenderDTO domain.TenderDTO) (domain.TenderDTO, error) {
	var tenderDB domain.TenderDTO

	err := r.InTx(ctx, func(_ pgx.Tx) error {
		var err error
		tenderOld, err := r.GetTender(ctx, tenderID)
		if err != nil {
			return err
		}

		err = r.SetTenderRevision(ctx, tenderOld)
		if err != nil {
			return err
		}

		tenderDB, err = r.UpdateTenderByID(ctx, tenderID, tenderDTO)
		if err != nil {
			return err
		}

		return nil
	})

	return tenderDB, err
}

// UpdateTenderByID updates tender by tender id.
func (r *Repo) UpdateTenderByID(ctx context.Context, tenderID string, tenderDTO domain.TenderDTO) (domain.TenderDTO, error) {
	var (
		args     []interface{}
		tenderDB domain.TenderDTO
	)
	args = append(args, tenderID)

	var query = `
	UPDATE tender SET version=version + 1, updated_at=CURRENT_TIMESTAMP`

	if tenderDTO.Name != "" {
		query += ", name=$" + fmt.Sprint(len(args)+1)
		args = append(args, tenderDTO.Name)
	}
	if tenderDTO.Description != "" {
		query += ", description=$" + fmt.Sprint(len(args)+1)
		args = append(args, tenderDTO.Description)
	}
	if tenderDTO.Status != "" {
		query += ", status=$" + fmt.Sprint(len(args)+1)
		args = append(args, tenderDTO.Status)
	}
	if tenderDTO.ServiceType != "" {
		query += ", type=$" + fmt.Sprint(len(args)+1)
		args = append(args, tenderDTO.ServiceType)
	}

	query += `
	WHERE id=$1
	RETURNING id, organization_id, user_id, name, description, status, type, version, created_at, updated_at;`

	err := r.conn.QueryRow(ctx, query, args...).
		Scan(&tenderDB.ID, &tenderDB.OrganizationID, &tenderDB.UserID, &tenderDB.Name, &tenderDB.Description, &tenderDB.Status, &tenderDB.ServiceType, &tenderDB.Version, &tenderDB.CreatedAt, &tenderDB.UpdatedAt)
	if err != nil {
		return tenderDB, err
	}

	return tenderDB, nil
}
