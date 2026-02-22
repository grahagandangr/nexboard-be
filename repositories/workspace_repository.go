package repositories

import (
	"database/sql"

	"github.com/grahagandangr/nexboard-be/models"
)

type WorkspaceRepository struct {
	DB *sql.DB
}

func NewWorkspaceRepository(db *sql.DB) *WorkspaceRepository {
	return &WorkspaceRepository{DB: db}
}

// CreateWorkspace creates a new workspace and adds the owner to the workspace_members table
func (r *WorkspaceRepository) CreateWorkspace(workspace *models.Workspace) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert workspace
	workspaceQuery := `
		INSERT INTO workspaces (external_id, name, description, owner_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	err = tx.QueryRow(workspaceQuery, workspace.ExternalID, workspace.Name, workspace.Description, workspace.OwnerID).
		Scan(&workspace.ID, &workspace.CreatedAt)
	if err != nil {
		return err
	}

	// Insert owner into workspace_members
	memberQuery := `
		INSERT INTO workspace_members (workspace_id, user_id, role)
		VALUES ($1, $2, 'owner')
	`
	_, err = tx.Exec(memberQuery, workspace.ID, workspace.OwnerID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetWorkspacesByUserID retrieves all workspaces a user is a member of
func (r *WorkspaceRepository) GetWorkspacesByUserID(userID int) ([]*models.Workspace, error) {
	query := `
		SELECT w.id, w.external_id, w.name, w.description, w.owner_id, w.active_status, w.created_at, w.modified_at, o.external_id
		FROM workspaces w
		JOIN workspace_members wm ON w.id = wm.workspace_id
		JOIN users o ON w.owner_id = o.id
		WHERE wm.user_id = $1 AND w.active_status = 1
	`
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workspaces []*models.Workspace
	for rows.Next() {
		w := &models.Workspace{}
		if err := rows.Scan(
			&w.ID,
			&w.ExternalID,
			&w.Name,
			&w.Description,
			&w.OwnerID,
			&w.ActiveStatus,
			&w.CreatedAt,
			&w.ModifiedAt,
			&w.OwnerExternalID,
		); err != nil {
			return nil, err
		}
		workspaces = append(workspaces, w)
	}

	return workspaces, nil
}

// GetWorkspaceByExternalID retrieves a workspace by its external ID
func (r *WorkspaceRepository) GetWorkspaceByExternalID(externalID string) (*models.Workspace, error) {
	w := &models.Workspace{}
	query := `
		SELECT w.id, w.external_id, w.name, w.description, w.owner_id, w.active_status, w.created_at, w.modified_at, o.external_id
		FROM workspaces w
		JOIN users o ON w.owner_id = o.id
		WHERE w.external_id = $1 AND w.active_status = 1
	`
	err := r.DB.QueryRow(query, externalID).Scan(
		&w.ID,
		&w.ExternalID,
		&w.Name,
		&w.Description,
		&w.OwnerID,
		&w.ActiveStatus,
		&w.CreatedAt,
		&w.ModifiedAt,
		&w.OwnerExternalID,
	)
	if err != nil {
		return nil, err
	}
	return w, nil
}

// UpdateWorkspace updates a workspace
func (r *WorkspaceRepository) UpdateWorkspace(w *models.Workspace) error {
	query := `
		UPDATE workspaces
		SET name = $1, description = $2, modified_at = NOW()
		WHERE id = $3
		RETURNING modified_at
	`
	return r.DB.QueryRow(query, w.Name, w.Description, w.ID).Scan(&w.ModifiedAt)
}

// DeleteWorkspace deletes a workspace
func (r *WorkspaceRepository) DeleteWorkspace(id int) error {
	query := `DELETE FROM workspaces WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	return err
}

// AddMember adds a user to a workspace
func (r *WorkspaceRepository) AddMember(wm *models.WorkspaceMember) error {
	query := `
		INSERT INTO workspace_members (workspace_id, user_id, role)
		VALUES ($1, $2, $3)
		RETURNING id, joined_at
	`
	return r.DB.QueryRow(query, wm.WorkspaceID, wm.UserID, wm.Role).Scan(&wm.ID, &wm.JoinedAt)
}

// GetMembers retrieves all members of a workspace
func (r *WorkspaceRepository) GetMembers(workspaceID int) ([]*models.WorkspaceMemberResponse, error) {
	query := `
		SELECT u.external_id, u.name, u.email, wm.role
		FROM workspace_members wm
		JOIN users u ON wm.user_id = u.id
		WHERE wm.workspace_id = $1
	`
	rows, err := r.DB.Query(query, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*models.WorkspaceMemberResponse
	for rows.Next() {
		m := &models.WorkspaceMemberResponse{}
		if err := rows.Scan(&m.UserExternalID, &m.Name, &m.Email, &m.Role); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}

// GetMemberRole gets a user's role in a workspace (used for permission checking)
func (r *WorkspaceRepository) GetMemberRole(workspaceID, userID int) (string, error) {
	var role string
	query := `
		SELECT role FROM workspace_members
		WHERE workspace_id = $1 AND user_id = $2
	`
	err := r.DB.QueryRow(query, workspaceID, userID).Scan(&role)
	return role, err
}

// UpdateMemberRole updates a user's role in a workspace
func (r *WorkspaceRepository) UpdateMemberRole(workspaceID, userID int, role string) error {
	query := `
		UPDATE workspace_members
		SET role = $1, modified_at = NOW()
		WHERE workspace_id = $2 AND user_id = $3
	`
	_, err := r.DB.Exec(query, role, workspaceID, userID)
	return err
}

// RemoveMember removes a user from a workspace
func (r *WorkspaceRepository) RemoveMember(workspaceID, userID int) error {
	query := `DELETE FROM workspace_members WHERE workspace_id = $1 AND user_id = $2`
	_, err := r.DB.Exec(query, workspaceID, userID)
	return err
}
