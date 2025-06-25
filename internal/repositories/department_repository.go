package repositories

import (
	"database/sql"
	"errors"
	"time"

	"github.com/aliwert/go-hospital-management/internal/models"
)

type DepartmentRepository struct {
	db *sql.DB
}

func NewDepartmentRepository(db *sql.DB) *DepartmentRepository {
	return &DepartmentRepository{db: db}
}

func (r *DepartmentRepository) Create(department *models.Department) error {
	// Check if department with same name exists
	var count int
	checkQuery := "SELECT COUNT(*) FROM departments WHERE name = $1 AND deleted_at IS NULL"
	err := r.db.QueryRow(checkQuery, department.Name).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("department with this name already exists")
	}

	// Insert new department
	query := `
		INSERT INTO departments (name, description, head_doctor_id, status, location,
		                        floor_number, phone_number, email, capacity, open_time,
		                        close_time, staff_count, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id, created_at, updated_at`

	err = r.db.QueryRow(
		query,
		department.Name,
		department.Description,
		department.HeadDoctorID,
		department.Status,
		department.Location,
		department.FloorNumber,
		department.PhoneNumber,
		department.Email,
		department.Capacity,
		department.OpenTime,
		department.CloseTime,
		department.StaffCount,
		department.IsActive,
		time.Now(),
		time.Now(),
	).Scan(&department.ID, &department.CreatedAt, &department.UpdatedAt)

	return err
}

func (r *DepartmentRepository) FindById(id uint) (*models.Department, error) {
	query := `
		SELECT d.id, d.name, d.description, d.head_doctor_id, d.status, d.location,
		       d.floor_number, d.phone_number, d.email, d.capacity, d.open_time,
		       d.close_time, d.staff_count, d.is_active, d.created_at, d.updated_at,
		       doc.name as head_doctor_name
		FROM departments d
		LEFT JOIN doctors doc ON d.head_doctor_id = doc.id
		WHERE d.id = $1 AND d.deleted_at IS NULL`

	var department models.Department
	var headDoctorName sql.NullString
	var headDoctorID sql.NullInt64

	err := r.db.QueryRow(query, id).Scan(
		&department.ID,
		&department.Name,
		&department.Description,
		&headDoctorID,
		&department.Status,
		&department.Location,
		&department.FloorNumber,
		&department.PhoneNumber,
		&department.Email,
		&department.Capacity,
		&department.OpenTime,
		&department.CloseTime,
		&department.StaffCount,
		&department.IsActive,
		&department.CreatedAt,
		&department.UpdatedAt,
		&headDoctorName,
	)

	if err != nil {
		return nil, err
	}

	// Set head doctor info
	if headDoctorID.Valid {
		department.HeadDoctorID = uint(headDoctorID.Int64)
		if headDoctorName.Valid {
			department.HeadDoctor.Name = headDoctorName.String
			department.HeadDoctor.ID = uint(headDoctorID.Int64)
		}
	}

	return &department, nil
}

func (r *DepartmentRepository) FindAll() ([]models.Department, error) {
	query := `
		SELECT d.id, d.name, d.description, d.head_doctor_id, d.status, d.location,
		       d.floor_number, d.phone_number, d.email, d.capacity, d.open_time,
		       d.close_time, d.staff_count, d.is_active, d.created_at, d.updated_at,
		       doc.name as head_doctor_name
		FROM departments d
		LEFT JOIN doctors doc ON d.head_doctor_id = doc.id
		WHERE d.deleted_at IS NULL
		ORDER BY d.created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var departments []models.Department
	for rows.Next() {
		var department models.Department
		var headDoctorName sql.NullString
		var headDoctorID sql.NullInt64

		err := rows.Scan(
			&department.ID,
			&department.Name,
			&department.Description,
			&headDoctorID,
			&department.Status,
			&department.Location,
			&department.FloorNumber,
			&department.PhoneNumber,
			&department.Email,
			&department.Capacity,
			&department.OpenTime,
			&department.CloseTime,
			&department.StaffCount,
			&department.IsActive,
			&department.CreatedAt,
			&department.UpdatedAt,
			&headDoctorName,
		)
		if err != nil {
			return nil, err
		}

		// Set head doctor info
		if headDoctorID.Valid {
			department.HeadDoctorID = uint(headDoctorID.Int64)
			if headDoctorName.Valid {
				department.HeadDoctor.Name = headDoctorName.String
				department.HeadDoctor.ID = uint(headDoctorID.Int64)
			}
		}

		departments = append(departments, department)
	}

	return departments, rows.Err()
}

func (r *DepartmentRepository) Update(department *models.Department) error {
	query := `
		UPDATE departments
		SET name = $2, description = $3, head_doctor_id = $4, status = $5, location = $6,
		    floor_number = $7, phone_number = $8, email = $9, capacity = $10, open_time = $11,
		    close_time = $12, staff_count = $13, is_active = $14, updated_at = $15
		WHERE id = $1 AND deleted_at IS NULL`

	_, err := r.db.Exec(
		query,
		department.ID,
		department.Name,
		department.Description,
		department.HeadDoctorID,
		department.Status,
		department.Location,
		department.FloorNumber,
		department.PhoneNumber,
		department.Email,
		department.Capacity,
		department.OpenTime,
		department.CloseTime,
		department.StaffCount,
		department.IsActive,
		time.Now(),
	)

	return err
}

func (r *DepartmentRepository) Delete(id uint) error {
	query := `UPDATE departments SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	_, err := r.db.Exec(query, time.Now(), id)
	return err
}
