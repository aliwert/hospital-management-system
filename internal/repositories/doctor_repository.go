package repositories

import (
	"database/sql"
	"errors"
	"time"

	"github.com/aliwert/go-hospital-management/internal/models"
)

type DoctorRepository struct {
	db *sql.DB
}

func NewDoctorRepository(db *sql.DB) *DoctorRepository {
	return &DoctorRepository{db: db}
}

func (r *DoctorRepository) Create(doctor *models.Doctor) error {
	// Check if doctor with same name exists
	var count int
	checkNameQuery := "SELECT COUNT(*) FROM doctors WHERE name = $1 AND deleted_at IS NULL"
	err := r.db.QueryRow(checkNameQuery, doctor.Name).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("doctor with this name already exists")
	}

	// Check if license number is unique
	checkLicenseQuery := "SELECT COUNT(*) FROM doctors WHERE license_number = $1 AND deleted_at IS NULL"
	err = r.db.QueryRow(checkLicenseQuery, doctor.LicenseNumber).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("license number already exists")
	}

	// Insert new doctor
	query := `
		INSERT INTO doctors (user_id, name, specialization, license_number, experience, department,
		                    availability, consultation_fee, status, education, qualifications, languages,
		                    biography, rating, review_count, office_number, working_days, working_hours,
		                    max_patients, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
		RETURNING id, created_at, updated_at`

	err = r.db.QueryRow(
		query,
		doctor.UserID,
		doctor.Name,
		doctor.Specialization,
		doctor.LicenseNumber,
		doctor.Experience,
		doctor.Department,
		doctor.Availability,
		doctor.ConsultationFee,
		doctor.Status,
		doctor.Education,
		doctor.Qualifications,
		doctor.Languages,
		doctor.Biography,
		doctor.Rating,
		doctor.ReviewCount,
		doctor.OfficeNumber,
		doctor.WorkingDays,
		doctor.WorkingHours,
		doctor.MaxPatients,
		time.Now(),
		time.Now(),
	).Scan(&doctor.ID, &doctor.CreatedAt, &doctor.UpdatedAt)

	return err
}

func (r *DoctorRepository) FindById(id uint) (*models.Doctor, error) {
	query := `
		SELECT d.id, d.user_id, d.name, d.specialization, d.license_number, d.experience,
		       d.department, d.availability, d.consultation_fee, d.status, d.education,
		       d.qualifications, d.languages, d.biography, d.rating, d.review_count,
		       d.office_number, d.working_days, d.working_hours, d.max_patients,
		       d.created_at, d.updated_at,
		       u.name as user_name, u.email as user_email, u.role as user_role
		FROM doctors d
		LEFT JOIN users u ON d.user_id = u.id
		WHERE d.id = $1 AND d.deleted_at IS NULL`

	var doctor models.Doctor
	var userName, userEmail, userRole sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&doctor.ID,
		&doctor.UserID,
		&doctor.Name,
		&doctor.Specialization,
		&doctor.LicenseNumber,
		&doctor.Experience,
		&doctor.Department,
		&doctor.Availability,
		&doctor.ConsultationFee,
		&doctor.Status,
		&doctor.Education,
		&doctor.Qualifications,
		&doctor.Languages,
		&doctor.Biography,
		&doctor.Rating,
		&doctor.ReviewCount,
		&doctor.OfficeNumber,
		&doctor.WorkingDays,
		&doctor.WorkingHours,
		&doctor.MaxPatients,
		&doctor.CreatedAt,
		&doctor.UpdatedAt,
		&userName,
		&userEmail,
		&userRole,
	)

	if err != nil {
		return nil, err
	}

	// Set user info
	if userName.Valid {
		doctor.User.Name = userName.String
	}
	if userEmail.Valid {
		doctor.User.Email = userEmail.String
	}
	if userRole.Valid {
		doctor.User.Role = userRole.String
	}
	doctor.User.ID = doctor.UserID

	return &doctor, nil
}

func (r *DoctorRepository) FindAll() ([]models.Doctor, error) {
	query := `
		SELECT d.id, d.user_id, d.name, d.specialization, d.license_number, d.experience,
		       d.department, d.availability, d.consultation_fee, d.status, d.education,
		       d.qualifications, d.languages, d.biography, d.rating, d.review_count,
		       d.office_number, d.working_days, d.working_hours, d.max_patients,
		       d.created_at, d.updated_at,
		       u.name as user_name, u.email as user_email, u.role as user_role
		FROM doctors d
		LEFT JOIN users u ON d.user_id = u.id
		WHERE d.deleted_at IS NULL
		ORDER BY d.created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var doctors []models.Doctor
	for rows.Next() {
		var doctor models.Doctor
		var userName, userEmail, userRole sql.NullString

		err := rows.Scan(
			&doctor.ID,
			&doctor.UserID,
			&doctor.Name,
			&doctor.Specialization,
			&doctor.LicenseNumber,
			&doctor.Experience,
			&doctor.Department,
			&doctor.Availability,
			&doctor.ConsultationFee,
			&doctor.Status,
			&doctor.Education,
			&doctor.Qualifications,
			&doctor.Languages,
			&doctor.Biography,
			&doctor.Rating,
			&doctor.ReviewCount,
			&doctor.OfficeNumber,
			&doctor.WorkingDays,
			&doctor.WorkingHours,
			&doctor.MaxPatients,
			&doctor.CreatedAt,
			&doctor.UpdatedAt,
			&userName,
			&userEmail,
			&userRole,
		)
		if err != nil {
			return nil, err
		}

		// Set user info
		if userName.Valid {
			doctor.User.Name = userName.String
		}
		if userEmail.Valid {
			doctor.User.Email = userEmail.String
		}
		if userRole.Valid {
			doctor.User.Role = userRole.String
		}
		doctor.User.ID = doctor.UserID

		doctors = append(doctors, doctor)
	}

	return doctors, rows.Err()
}

func (r *DoctorRepository) Update(doctor *models.Doctor) error {
	query := `
		UPDATE doctors
		SET user_id = $2, name = $3, specialization = $4, license_number = $5, experience = $6,
		    department = $7, availability = $8, consultation_fee = $9, status = $10, education = $11,
		    qualifications = $12, languages = $13, biography = $14, rating = $15, review_count = $16,
		    office_number = $17, working_days = $18, working_hours = $19, max_patients = $20, updated_at = $21
		WHERE id = $1 AND deleted_at IS NULL`

	_, err := r.db.Exec(
		query,
		doctor.ID,
		doctor.UserID,
		doctor.Name,
		doctor.Specialization,
		doctor.LicenseNumber,
		doctor.Experience,
		doctor.Department,
		doctor.Availability,
		doctor.ConsultationFee,
		doctor.Status,
		doctor.Education,
		doctor.Qualifications,
		doctor.Languages,
		doctor.Biography,
		doctor.Rating,
		doctor.ReviewCount,
		doctor.OfficeNumber,
		doctor.WorkingDays,
		doctor.WorkingHours,
		doctor.MaxPatients,
		time.Now(),
	)

	return err
}

func (r *DoctorRepository) Delete(id uint) error {
	query := `UPDATE doctors SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	_, err := r.db.Exec(query, time.Now(), id)
	return err
}

// FindBySpecialization finds doctors by specialization
func (r *DoctorRepository) FindBySpecialization(specialization string) ([]models.Doctor, error) {
	query := `
		SELECT d.id, d.user_id, d.name, d.specialization, d.license_number, d.experience,
		       d.department, d.availability, d.consultation_fee, d.status, d.education,
		       d.qualifications, d.languages, d.biography, d.rating, d.review_count,
		       d.office_number, d.working_days, d.working_hours, d.max_patients,
		       d.created_at, d.updated_at,
		       u.name as user_name, u.email as user_email, u.role as user_role
		FROM doctors d
		LEFT JOIN users u ON d.user_id = u.id
		WHERE d.specialization = $1 AND d.deleted_at IS NULL AND d.availability = true
		ORDER BY d.rating DESC`

	rows, err := r.db.Query(query, specialization)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var doctors []models.Doctor
	for rows.Next() {
		var doctor models.Doctor
		var userName, userEmail, userRole sql.NullString

		err := rows.Scan(
			&doctor.ID,
			&doctor.UserID,
			&doctor.Name,
			&doctor.Specialization,
			&doctor.LicenseNumber,
			&doctor.Experience,
			&doctor.Department,
			&doctor.Availability,
			&doctor.ConsultationFee,
			&doctor.Status,
			&doctor.Education,
			&doctor.Qualifications,
			&doctor.Languages,
			&doctor.Biography,
			&doctor.Rating,
			&doctor.ReviewCount,
			&doctor.OfficeNumber,
			&doctor.WorkingDays,
			&doctor.WorkingHours,
			&doctor.MaxPatients,
			&doctor.CreatedAt,
			&doctor.UpdatedAt,
			&userName,
			&userEmail,
			&userRole,
		)
		if err != nil {
			return nil, err
		}

		// Set user info
		if userName.Valid {
			doctor.User.Name = userName.String
		}
		if userEmail.Valid {
			doctor.User.Email = userEmail.String
		}
		if userRole.Valid {
			doctor.User.Role = userRole.String
		}
		doctor.User.ID = doctor.UserID

		doctors = append(doctors, doctor)
	}

	return doctors, rows.Err()
}
