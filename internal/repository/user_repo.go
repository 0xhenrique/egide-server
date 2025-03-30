package repository

import (
	"database/sql"
	"errors"
	"time"

	"egide-server/internal/models"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create adds a new user to the database
func (r *UserRepository) Create(user *models.User) (int64, error) {
	query := `
		INSERT INTO users (github_id, username, email, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := r.db.Exec(
		query,
		user.GitHubID,
		user.Username,
		user.Email,
		now,
		now,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(id int64) (*models.User, error) {
	query := `
		SELECT id, github_id, username, email, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	var user models.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.GitHubID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// FindByGitHubID finds a user by GitHub ID
func (r *UserRepository) FindByGitHubID(githubID string) (*models.User, error) {
	query := `
		SELECT id, github_id, username, email, created_at, updated_at
		FROM users
		WHERE github_id = ?
	`

	var user models.User
	err := r.db.QueryRow(query, githubID).Scan(
		&user.ID,
		&user.GitHubID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// Update updates an existing user
func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE users
		SET username = ?, email = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(
		query,
		user.Username,
		user.Email,
		time.Now(),
		user.ID,
	)

	return err
}
