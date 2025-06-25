package repositories

import (
	"database/sql"
	"errors"
	"time"

	"github.com/aliwert/go-hospital-management/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	// Check if email already exists
	var count int
	checkQuery := "SELECT COUNT(*) FROM users WHERE email = $1 AND deleted_at IS NULL"
	err := r.db.QueryRow(checkQuery, user.Email).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("email already exists")
	}

	// Insert new user
	query := `
		INSERT INTO users (name, email, password, role, status, created_by, updated_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at`

	err = r.db.QueryRow(
		query,
		user.Name,
		user.Email,
		user.Password,
		user.Role,
		user.Status,
		user.CreatedBy,
		user.UpdatedBy,
		time.Now(),
		time.Now(),
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	return err
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, name, email, password, role, status, version, created_by, updated_by, last_login, created_at, updated_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL`

	var user models.User
	var version sql.NullInt64
	var createdBy, updatedBy sql.NullInt64
	var lastLogin sql.NullTime

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.Status,
		&version,
		&createdBy,
		&updatedBy,
		&lastLogin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Set nullable fields
	if version.Valid {
		user.Version = uint(version.Int64)
	}
	if createdBy.Valid {
		user.CreatedBy = uint(createdBy.Int64)
	}
	if updatedBy.Valid {
		user.UpdatedBy = uint(updatedBy.Int64)
	}
	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}

	return &user, nil
}

func (r *UserRepository) FindById(id uint) (*models.User, error) {
	query := `
		SELECT id, name, email, password, role, status, version, created_by, updated_by, last_login, created_at, updated_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL`

	var user models.User
	var version sql.NullInt64
	var createdBy, updatedBy sql.NullInt64
	var lastLogin sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.Status,
		&version,
		&createdBy,
		&updatedBy,
		&lastLogin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Set nullable fields
	if version.Valid {
		user.Version = uint(version.Int64)
	}
	if createdBy.Valid {
		user.CreatedBy = uint(createdBy.Int64)
	}
	if updatedBy.Valid {
		user.UpdatedBy = uint(updatedBy.Int64)
	}
	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}

	return &user, nil
}

func (r *UserRepository) FindAll() ([]models.User, error) {
	query := `
		SELECT id, name, email, role, status, version, created_by, updated_by, last_login, created_at, updated_at
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		var version sql.NullInt64
		var createdBy, updatedBy sql.NullInt64
		var lastLogin sql.NullTime

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Role,
			&user.Status,
			&version,
			&createdBy,
			&updatedBy,
			&lastLogin,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Set nullable fields
		if version.Valid {
			user.Version = uint(version.Int64)
		}
		if createdBy.Valid {
			user.CreatedBy = uint(createdBy.Int64)
		}
		if updatedBy.Valid {
			user.UpdatedBy = uint(updatedBy.Int64)
		}
		if lastLogin.Valid {
			user.LastLogin = &lastLogin.Time
		}

		users = append(users, user)
	}

	return users, rows.Err()
}

func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE users
		SET name = $2, email = $3, role = $4, status = $5, version = $6, updated_by = $7, updated_at = $8
		WHERE id = $1 AND deleted_at IS NULL`

	_, err := r.db.Exec(
		query,
		user.ID,
		user.Name,
		user.Email,
		user.Role,
		user.Status,
		user.Version,
		user.UpdatedBy,
		time.Now(),
	)

	return err
}

func (r *UserRepository) Delete(id uint) error {
	query := `UPDATE users SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	_, err := r.db.Exec(query, time.Now(), id)
	return err
}

// UpdateLastLogin updates the last login timestamp for a user
func (r *UserRepository) UpdateLastLogin(id uint) error {
	query := `UPDATE users SET last_login = $1 WHERE id = $2 AND deleted_at IS NULL`
	_, err := r.db.Exec(query, time.Now(), id)
	return err
}
