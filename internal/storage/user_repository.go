package storage

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
    "github.com/mavuno/mavuno-backend/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func(r *UserRepository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, role, full_name, phone_number, created_at, updated_at, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.FullName,
		user.PhoneNumber,
		user.CreatedAt,
		user.UpdatedAt,
		user.IsActive,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func(r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, role, full_name, phone_number, created_at, updated_at, is_active
		FROM users
		WHERE email = $1 AND is_active = true
	`
	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.Role,
			&user.FullName,
			&user.PhoneNumber,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.IsActive,
		)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("User Not Found")
		}
		if err != nil {
			return nil, fmt.Errorf("Failed to get user by email: %w", err)
		}
		return user, nil
}

func (r *UserRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, role, full_name, phone_number, created_at, updated_at, is_active
		FROM users
		WHERE id = $1 AND is_active = true
	`
	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
        &user.Email,
        &user.PasswordHash,
        &user.Role,
        &user.FullName,
        &user.PhoneNumber,
        &user.CreatedAt,
        &user.UpdatedAt,
        &user.IsActive,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("User Not Found")
	}
	if err != nil {
		return nil, fmt.Errorf("Failed to get user by ID: %w", err)
	}
	return user, nil
}