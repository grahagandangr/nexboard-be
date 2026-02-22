package repositories

import (
	"database/sql"

	"github.com/grahagandangr/nexboard-be/models"
)

type TaskRepository struct {
	DB *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{DB: db}
}

// CreateTask adds a new task to a board
func (r *TaskRepository) CreateTask(task *models.Task) error {
	query := `
		INSERT INTO tasks (external_id, board_id, status_id, assigned_to, created_by_id, title, description, priority, due_date, position)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at
	`
	return r.DB.QueryRow(
		query,
		task.ExternalID,
		task.BoardID,
		task.StatusID,
		task.AssignedTo,
		task.CreatedByID,
		task.Title,
		task.Description,
		task.Priority,
		task.DueDate,
		task.Position,
	).Scan(&task.ID, &task.CreatedAt)
}

// GetTasksByBoardID gets all active tasks for a specific board
func (r *TaskRepository) GetTasksByBoardID(boardID int) ([]*models.TaskResponse, error) {
	query := `
		SELECT 
			t.external_id,
			b.external_id AS board_external_id,
			s.external_id AS status_external_id,
			s.name AS status_name,
			s.color AS status_color,
			u.external_id AS assignee_external_id,
			u.name AS assignee_name,
			t.title,
			t.description,
			t.priority,
			t.due_date,
			t.position,
			t.created_at,
			t.modified_at
		FROM tasks t
		JOIN boards b ON t.board_id = b.id
		JOIN statuses s ON t.status_id = s.id
		LEFT JOIN users u ON t.assigned_to = u.id
		WHERE t.board_id = $1 AND t.active_status = 1
		ORDER BY s.position ASC, t.position ASC
	`
	rows, err := r.DB.Query(query, boardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.TaskResponse
	for rows.Next() {
		var (
			assigneeExtID *string
			assigneeName  *string
			statusColor   *string
		)
		tr := &models.TaskResponse{}

		if err := rows.Scan(
			&tr.ExternalID,
			&tr.BoardExternalID,
			&tr.Status.ExternalID,
			&tr.Status.Name,
			&statusColor,
			&assigneeExtID,
			&assigneeName,
			&tr.Title,
			&tr.Description,
			&tr.Priority,
			&tr.DueDate,
			&tr.Position,
			&tr.CreatedAt,
			&tr.ModifiedAt,
		); err != nil {
			return nil, err
		}

		tr.Status.Color = statusColor

		if assigneeExtID != nil {
			tr.AssignedTo = &models.TaskAssigneeInfo{
				ExternalID: *assigneeExtID,
				Name:       *assigneeName,
			}
		}

		tasks = append(tasks, tr)
	}

	return tasks, nil
}

// GetTaskByExternalID retrieves details of a specific task
func (r *TaskRepository) GetTaskByExternalID(externalID string) (*models.Task, error) {
	query := `
		SELECT id, external_id, board_id, status_id, assigned_to, created_by_id, title, description, priority, due_date, position, active_status, created_at, modified_at
		FROM tasks
		WHERE external_id = $1 AND active_status = 1
	`
	t := &models.Task{}
	err := r.DB.QueryRow(query, externalID).Scan(
		&t.ID,
		&t.ExternalID,
		&t.BoardID,
		&t.StatusID,
		&t.AssignedTo,
		&t.CreatedByID,
		&t.Title,
		&t.Description,
		&t.Priority,
		&t.DueDate,
		&t.Position,
		&t.ActiveStatus,
		&t.CreatedAt,
		&t.ModifiedAt,
	)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// UpdateTask modifies a task
func (r *TaskRepository) UpdateTask(t *models.Task) error {
	query := `
		UPDATE tasks
		SET title = $1, description = $2, priority = $3, due_date = $4, status_id = $5, assigned_to = $6, position = $7, modified_at = NOW()
		WHERE id = $8
		RETURNING modified_at
	`
	return r.DB.QueryRow(query, t.Title, t.Description, t.Priority, t.DueDate, t.StatusID, t.AssignedTo, t.Position, t.ID).
		Scan(&t.ModifiedAt)
}

// DeleteTask removes a task
func (r *TaskRepository) DeleteTask(id int) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	return err
}
