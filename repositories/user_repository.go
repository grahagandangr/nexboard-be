package repositories

import (
	"database/sql"

	"github.com/grahagandangr/nexboard-be/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (external_id, name, email, password)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	return r.DB.QueryRow(query, user.ExternalID, user.Name, user.Email, user.Password).
		Scan(&user.ID, &user.CreatedAt)
}

// GetUserByEmail retrieves a user by email
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, external_id, name, email, password, avatar_url, active_status, created_at, created_by, modified_at, modified_by
		FROM users
		WHERE email = $1 AND active_status = 1
	`
	err := r.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.ExternalID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.AvatarURL,
		&user.ActiveStatus,
		&user.CreatedAt,
		&user.CreatedBy,
		&user.ModifiedAt,
		&user.ModifiedBy,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByExternalID retrieves a user by their external ID
func (r *UserRepository) GetUserByExternalID(externalID string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, external_id, name, email, password, avatar_url, active_status, created_at, created_by, modified_at, modified_by
		FROM users
		WHERE external_id = $1 AND active_status = 1
	`
	err := r.DB.QueryRow(query, externalID).Scan(
		&user.ID,
		&user.ExternalID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.AvatarURL,
		&user.ActiveStatus,
		&user.CreatedAt,
		&user.CreatedBy,
		&user.ModifiedAt,
		&user.ModifiedBy,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByID retrieves a user by ID
func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, external_id, name, email, password, avatar_url, active_status, created_at, created_by, modified_at, modified_by
		FROM users
		WHERE id = $1 AND active_status = 1
	`
	err := r.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.ExternalID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.AvatarURL,
		&user.ActiveStatus,
		&user.CreatedAt,
		&user.CreatedBy,
		&user.ModifiedAt,
		&user.ModifiedBy,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUserProfile updates user's basic info
func (r *UserRepository) UpdateUserProfile(user *models.User) error {
	query := `
		UPDATE users
		SET name = $1, avatar_url = $2, modified_at = NOW()
		WHERE id = $3
		RETURNING modified_at
	`
	return r.DB.QueryRow(query, user.Name, user.AvatarURL, user.ID).Scan(&user.ModifiedAt)
}
