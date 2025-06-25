package repositories

import (
	"database/sql"
	"time"

	"github.com/aliwert/go-hospital-management/internal/models"
)

type PatientRepository struct {
	db *sql.DB
}

func NewPatientRepository(db *sql.DB) *PatientRepository {
	return &PatientRepository{db: db}
}

func (r *PatientRepository) Create(patient *models.Patient) error {
	query := `
		INSERT INTO patients (user_id, date_of_birth, gender, blood_type, address, phone_number,
		                     emergency_contact, emergency_phone, insurance, insurance_no, allergies,
		                     medical_history, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(
		query,
		patient.UserID,
		patient.DateOfBirth,
		patient.Gender,
		patient.BloodType,
		patient.Address,
		patient.PhoneNumber,
		patient.EmergencyContact,
		patient.EmergencyPhone,
		patient.Insurance,
		patient.InsuranceNo,
		patient.Allergies,
		patient.MedicalHistory,
		patient.Status,
		time.Now(),
		time.Now(),
	).Scan(&patient.ID, &patient.CreatedAt, &patient.UpdatedAt)

	return err
}

func (r *PatientRepository) FindById(id uint) (*models.Patient, error) {
	query := `
		SELECT p.id, p.user_id, p.date_of_birth, p.gender, p.blood_type, p.address,
		       p.phone_number, p.emergency_contact, p.emergency_phone, p.insurance,
		       p.insurance_no, p.allergies, p.medical_history, p.status, p.created_at, p.updated_at,
		       u.name as user_name, u.email as user_email, u.role as user_role
		FROM patients p
		LEFT JOIN users u ON p.user_id = u.id
		WHERE p.id = $1 AND p.deleted_at IS NULL`

	var patient models.Patient
	var userName, userEmail, userRole sql.NullString
	var bloodType, insurance, insuranceNo, allergies, medicalHistory sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&patient.ID,
		&patient.UserID,
		&patient.DateOfBirth,
		&patient.Gender,
		&bloodType,
		&patient.Address,
		&patient.PhoneNumber,
		&patient.EmergencyContact,
		&patient.EmergencyPhone,
		&insurance,
		&insuranceNo,
		&allergies,
		&medicalHistory,
		&patient.Status,
		&patient.CreatedAt,
		&patient.UpdatedAt,
		&userName,
		&userEmail,
		&userRole,
	)

	if err != nil {
		return nil, err
	}

	// Set nullable fields
	if bloodType.Valid {
		patient.BloodType = bloodType.String
	}
	if insurance.Valid {
		patient.Insurance = insurance.String
	}
	if insuranceNo.Valid {
		patient.InsuranceNo = insuranceNo.String
	}
	if allergies.Valid {
		patient.Allergies = allergies.String
	}
	if medicalHistory.Valid {
		patient.MedicalHistory = medicalHistory.String
	}

	// Set user info
	if userName.Valid {
		patient.User.Name = userName.String
	}
	if userEmail.Valid {
		patient.User.Email = userEmail.String
	}
	if userRole.Valid {
		patient.User.Role = userRole.String
	}
	patient.User.ID = patient.UserID

	return &patient, nil
}

func (r *PatientRepository) FindAll() ([]models.Patient, error) {
	query := `
		SELECT p.id, p.user_id, p.date_of_birth, p.gender, p.blood_type, p.address,
		       p.phone_number, p.emergency_contact, p.emergency_phone, p.insurance,
		       p.insurance_no, p.allergies, p.medical_history, p.status, p.created_at, p.updated_at,
		       u.name as user_name, u.email as user_email, u.role as user_role
		FROM patients p
		LEFT JOIN users u ON p.user_id = u.id
		WHERE p.deleted_at IS NULL
		ORDER BY p.created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []models.Patient
	for rows.Next() {
		var patient models.Patient
		var userName, userEmail, userRole sql.NullString
		var bloodType, insurance, insuranceNo, allergies, medicalHistory sql.NullString

		err := rows.Scan(
			&patient.ID,
			&patient.UserID,
			&patient.DateOfBirth,
			&patient.Gender,
			&bloodType,
			&patient.Address,
			&patient.PhoneNumber,
			&patient.EmergencyContact,
			&patient.EmergencyPhone,
			&insurance,
			&insuranceNo,
			&allergies,
			&medicalHistory,
			&patient.Status,
			&patient.CreatedAt,
			&patient.UpdatedAt,
			&userName,
			&userEmail,
			&userRole,
		)
		if err != nil {
			return nil, err
		}

		// Set nullable fields
		if bloodType.Valid {
			patient.BloodType = bloodType.String
		}
		if insurance.Valid {
			patient.Insurance = insurance.String
		}
		if insuranceNo.Valid {
			patient.InsuranceNo = insuranceNo.String
		}
		if allergies.Valid {
			patient.Allergies = allergies.String
		}
		if medicalHistory.Valid {
			patient.MedicalHistory = medicalHistory.String
		}

		// Set user info
		if userName.Valid {
			patient.User.Name = userName.String
		}
		if userEmail.Valid {
			patient.User.Email = userEmail.String
		}
		if userRole.Valid {
			patient.User.Role = userRole.String
		}
		patient.User.ID = patient.UserID

		patients = append(patients, patient)
	}

	return patients, rows.Err()
}

func (r *PatientRepository) Update(patient *models.Patient) error {
	query := `
		UPDATE patients
		SET user_id = $2, date_of_birth = $3, gender = $4, blood_type = $5, address = $6,
		    phone_number = $7, emergency_contact = $8, emergency_phone = $9, insurance = $10,
		    insurance_no = $11, allergies = $12, medical_history = $13, status = $14, updated_at = $15
		WHERE id = $1 AND deleted_at IS NULL`

	_, err := r.db.Exec(
		query,
		patient.ID,
		patient.UserID,
		patient.DateOfBirth,
		patient.Gender,
		patient.BloodType,
		patient.Address,
		patient.PhoneNumber,
		patient.EmergencyContact,
		patient.EmergencyPhone,
		patient.Insurance,
		patient.InsuranceNo,
		patient.Allergies,
		patient.MedicalHistory,
		patient.Status,
		time.Now(),
	)

	return err
}

func (r *PatientRepository) Delete(id uint) error {
	query := `UPDATE patients SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	_, err := r.db.Exec(query, time.Now(), id)
	return err
}
