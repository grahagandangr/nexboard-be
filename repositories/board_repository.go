package repositories

import (
	"database/sql"

	"github.com/grahagandangr/nexboard-be/models"
)

type BoardRepository struct {
	DB *sql.DB
}

func NewBoardRepository(db *sql.DB) *BoardRepository {
	return &BoardRepository{DB: db}
}

// CreateBoard inserts a new board into the database
func (r *BoardRepository) CreateBoard(board *models.Board) error {
	query := `
		INSERT INTO boards (external_id, workspace_id, created_by_id, name, description)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	return r.DB.QueryRow(query, board.ExternalID, board.WorkspaceID, board.CreatedByID, board.Name, board.Description).
		Scan(&board.ID, &board.CreatedAt)
}

// GetBoardsByWorkspaceID retrieves all active boards for a given workspace ID
func (r *BoardRepository) GetBoardsByWorkspaceID(workspaceID int) ([]*models.Board, error) {
	query := `
		SELECT b.id, b.external_id, b.workspace_id, b.created_by_id, b.name, b.description, b.active_status, b.created_at, b.modified_at, w.external_id
		FROM boards b
		JOIN workspaces w ON b.workspace_id = w.id
		WHERE b.workspace_id = $1 AND b.active_status = 1
	`
	rows, err := r.DB.Query(query, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var boards []*models.Board
	for rows.Next() {
		b := &models.Board{}
		if err := rows.Scan(
			&b.ID,
			&b.ExternalID,
			&b.WorkspaceID,
			&b.CreatedByID,
			&b.Name,
			&b.Description,
			&b.ActiveStatus,
			&b.CreatedAt,
			&b.ModifiedAt,
			&b.WorkspaceExternalID,
		); err != nil {
			return nil, err
		}
		boards = append(boards, b)
	}

	return boards, nil
}

// GetBoardByExternalID retrieves a single board by its external ID
func (r *BoardRepository) GetBoardByExternalID(externalID string) (*models.Board, error) {
	b := &models.Board{}
	query := `
		SELECT b.id, b.external_id, b.workspace_id, b.created_by_id, b.name, b.description, b.active_status, b.created_at, b.modified_at, w.external_id
		FROM boards b
		JOIN workspaces w ON b.workspace_id = w.id
		WHERE b.external_id = $1 AND b.active_status = 1
	`
	err := r.DB.QueryRow(query, externalID).Scan(
		&b.ID,
		&b.ExternalID,
		&b.WorkspaceID,
		&b.CreatedByID,
		&b.Name,
		&b.Description,
		&b.ActiveStatus,
		&b.CreatedAt,
		&b.ModifiedAt,
		&b.WorkspaceExternalID,
	)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// UpdateBoard updates the name and description of a board
func (r *BoardRepository) UpdateBoard(b *models.Board) error {
	query := `
		UPDATE boards
		SET name = $1, description = $2, modified_at = NOW()
		WHERE id = $3
		RETURNING modified_at
	`
	return r.DB.QueryRow(query, b.Name, b.Description, b.ID).Scan(&b.ModifiedAt)
}

// DeleteBoard hard deletes a board
func (r *BoardRepository) DeleteBoard(id int) error {
	query := `DELETE FROM boards WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	return err
}
