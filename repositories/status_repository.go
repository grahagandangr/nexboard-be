package repositories

import (
	"database/sql"

	"github.com/grahagandangr/nexboard-be/models"
)

type StatusRepository struct {
	DB *sql.DB
}

func NewStatusRepository(db *sql.DB) *StatusRepository {
	return &StatusRepository{DB: db}
}

// CreateStatus inserts a new status into the database
func (r *StatusRepository) CreateStatus(status *models.Status) error {
	query := `
		INSERT INTO statuses (external_id, name, color, position)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	return r.DB.QueryRow(query, status.ExternalID, status.Name, status.Color, status.Position).
		Scan(&status.ID, &status.CreatedAt)
}

// GetAllStatuses retrieves all active statuses safely
func (r *StatusRepository) GetAllStatuses() ([]*models.Status, error) {
	query := `
		SELECT id, external_id, name, color, position, active_status, created_at, modified_at
		FROM statuses
		WHERE active_status = 1
		ORDER BY position ASC
	`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statuses []*models.Status
	for rows.Next() {
		s := &models.Status{}
		if err := rows.Scan(
			&s.ID,
			&s.ExternalID,
			&s.Name,
			&s.Color,
			&s.Position,
			&s.ActiveStatus,
			&s.CreatedAt,
			&s.ModifiedAt,
		); err != nil {
			return nil, err
		}
		statuses = append(statuses, s)
	}
	return statuses, nil
}

// GetStatusByExternalID retrieves a single status
func (r *StatusRepository) GetStatusByExternalID(externalID string) (*models.Status, error) {
	s := &models.Status{}
	query := `
		SELECT id, external_id, name, color, position, active_status, created_at, modified_at
		FROM statuses
		WHERE external_id = $1 AND active_status = 1
	`
	err := r.DB.QueryRow(query, externalID).Scan(
		&s.ID,
		&s.ExternalID,
		&s.Name,
		&s.Color,
		&s.Position,
		&s.ActiveStatus,
		&s.CreatedAt,
		&s.ModifiedAt,
	)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// GetStatusByName checks if a status with the same name already exists
func (r *StatusRepository) GetStatusByName(name string) (*models.Status, error) {
	s := &models.Status{}
	query := `
		SELECT id, external_id, name
		FROM statuses
		WHERE name = $1 AND active_status = 1
	`
	err := r.DB.QueryRow(query, name).Scan(&s.ID, &s.ExternalID, &s.Name)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// UpdateStatus modifies an existing status
func (r *StatusRepository) UpdateStatus(s *models.Status) error {
	query := `
		UPDATE statuses
		SET name = $1, color = $2, position = $3, modified_at = NOW()
		WHERE id = $4
		RETURNING modified_at
	`
	return r.DB.QueryRow(query, s.Name, s.Color, s.Position, s.ID).Scan(&s.ModifiedAt)
}

// CheckIfReferenced checks if the status is being used by any active tasks
func (r *StatusRepository) CheckIfReferenced(statusID int) (bool, error) {
	var count int
	query := `
		SELECT COUNT(*)
		FROM tasks
		WHERE status_id = $1 AND active_status = 1
	`
	err := r.DB.QueryRow(query, statusID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// DeleteStatus performs a hard delete or soft delete
func (r *StatusRepository) DeleteStatus(id int) error {
	query := `DELETE FROM statuses WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	return err
}
