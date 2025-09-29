package database

import (
	"database/sql"
	"fmt"
	"time"

	"go.uber.org/fx"
)

// UserRepository handles user database operations
type UserRepository struct {
	db *DB
}

type UserRepositoryDeps struct {
	fx.In
	DB *DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(deps UserRepositoryDeps) *UserRepository {
	return &UserRepository{db: deps.DB}
}

// Create creates a new user
func (r *UserRepository) Create(user *User) error {
	query := `
		INSERT INTO users (email, phone, name, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	err := r.db.QueryRow(query, user.Email, user.Phone, user.Name, time.Now()).
		Scan(&user.ID, &user.CreatedAt)

	return err
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id int) (*User, error) {
	user := &User{}
	query := `SELECT id, email, phone, name, created_at, updated_at FROM users WHERE id = $1`

	err := r.db.Get(user, query, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*User, error) {
	user := &User{}
	query := `SELECT id, email, phone, name, created_at, updated_at FROM users WHERE email = $1`

	err := r.db.Get(user, query, email)
	if err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}

// Update updates a user
func (r *UserRepository) Update(user *User) error {
	query := `
		UPDATE users 
		SET phone = $2, name = $3, updated_at = $4
		WHERE id = $1
		RETURNING updated_at`

	err := r.db.QueryRow(query, user.ID, user.Phone, user.Name, time.Now()).
		Scan(&user.UpdatedAt)

	return err
}

// Delete deletes a user
func (r *UserRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d not found", id)
	}

	return nil
}

// List retrieves all users with pagination
func (r *UserRepository) List(limit, offset int) ([]*User, error) {
	users := []*User{}
	query := `
		SELECT id, email, phone, name, created_at, updated_at 
		FROM users 
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2`

	err := r.db.Select(&users, query, limit, offset)
	return users, err
}

// Count returns the total number of users
func (r *UserRepository) Count() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users`

	err := r.db.Get(&count, query)
	return count, err
}
