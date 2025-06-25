package repositories

import (
	"database/sql"
	"time"

	"github.com/aliwert/go-hospital-management/internal/models"
	_ "github.com/lib/pq"
)

type MedicalRecordRepository struct {
	db *sql.DB
}

func NewMedicalRecordRepository(db *sql.DB) *MedicalRecordRepository {
	return &MedicalRecordRepository{db: db}
}

func (r *MedicalRecordRepository) Create(record *models.MedicalRecord) error {
	query := `
		INSERT INTO medical_records (
			patient_id, doctor_id, visit_date, diagnosis, symptoms, treatment, notes,
			prescription_id, blood_pressure, temperature, weight, height, bmi,
			pulse_rate, respiratory_rate, allergies, medications, follow_up_date,
			is_follow_up, attachments, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
		RETURNING id`

	var prescriptionID sql.NullInt64
	if record.PrescriptionID > 0 {
		prescriptionID.Int64 = int64(record.PrescriptionID)
		prescriptionID.Valid = true
	}

	err := r.db.QueryRow(
		query,
		record.PatientID, record.DoctorID, record.VisitDate, record.Diagnosis, record.Symptoms,
		record.Treatment, record.Notes, prescriptionID, record.BloodPressure, record.Temperature,
		record.Weight, record.Height, record.BMI, record.PulseRate, record.RespiratoryRate,
		record.Allergies, record.Medications, record.FollowUpDate, record.IsFollowUp,
		record.Attachments, record.Status,
	).Scan(&record.ID)

	return err
}

func (r *MedicalRecordRepository) FindById(id uint) (*models.MedicalRecord, error) {
	record := &models.MedicalRecord{}

	query := `
		SELECT mr.id, mr.patient_id, mr.doctor_id, mr.visit_date, mr.diagnosis, mr.symptoms,
		       mr.treatment, mr.notes, mr.prescription_id, mr.blood_pressure, mr.temperature,
		       mr.weight, mr.height, mr.bmi, mr.pulse_rate, mr.respiratory_rate, mr.allergies,
		       mr.medications, mr.follow_up_date, mr.is_follow_up, mr.attachments, mr.status,
		       -- Patient info
		       p.id, p.user_id, p.date_of_birth, p.gender, p.blood_type,
		       p.emergency_contact, p.emergency_phone, p.insurance, p.insurance_no,
		       p.created_at, p.updated_at,
		       -- Patient User info
		       pu.id, pu.name, pu.email, pu.created_at, pu.updated_at,
		       -- Doctor info
		       d.id, d.user_id, d.department_id, d.specialization, d.license_number,
		       d.experience_years, d.education, d.created_at, d.updated_at,
		       -- Doctor User info
		       du.id, du.name, du.email, du.created_at, du.updated_at
		FROM medical_records mr
		LEFT JOIN patients p ON mr.patient_id = p.id
		LEFT JOIN users pu ON p.user_id = pu.id
		LEFT JOIN doctors d ON mr.doctor_id = d.id
		LEFT JOIN users du ON d.user_id = du.id
		WHERE mr.id = $1 AND mr.deleted_at IS NULL`

	var prescriptionID sql.NullInt64
	var followUpDate sql.NullTime

	var patientID, patientUserID sql.NullInt64
	var patientGender, patientBloodType, patientEmergencyContact sql.NullString
	var patientEmergencyPhone, patientInsurance, patientInsuranceNo sql.NullString
	var patientDateOfBirth sql.NullTime
	var patientCreatedAt, patientUpdatedAt sql.NullTime
	var patientUserName, patientUserEmail sql.NullString
	var patientUserCreatedAt, patientUserUpdatedAt sql.NullTime
	var patientUserIDScan sql.NullInt64

	var doctorID, doctorUserID, doctorDepartmentID sql.NullInt64
	var doctorSpecialization, doctorLicenseNumber, doctorEducation sql.NullString
	var doctorExperienceYears sql.NullInt64
	var doctorCreatedAt, doctorUpdatedAt sql.NullTime
	var doctorUserName, doctorUserEmail sql.NullString
	var doctorUserCreatedAt, doctorUserUpdatedAt sql.NullTime
	var doctorUserIDScan sql.NullInt64

	err := r.db.QueryRow(query, id).Scan(
		&record.ID, &record.PatientID, &record.DoctorID, &record.VisitDate, &record.Diagnosis,
		&record.Symptoms, &record.Treatment, &record.Notes, &prescriptionID, &record.BloodPressure,
		&record.Temperature, &record.Weight, &record.Height, &record.BMI, &record.PulseRate,
		&record.RespiratoryRate, &record.Allergies, &record.Medications, &followUpDate,
		&record.IsFollowUp, &record.Attachments, &record.Status,
		// Patient fields
		&patientID, &patientUserID, &patientDateOfBirth, &patientGender,
		&patientBloodType, &patientEmergencyContact, &patientEmergencyPhone, &patientInsurance,
		&patientInsuranceNo, &patientCreatedAt, &patientUpdatedAt,
		// Patient User fields
		&patientUserIDScan, &patientUserName, &patientUserEmail, &patientUserCreatedAt, &patientUserUpdatedAt,
		// Doctor fields
		&doctorID, &doctorUserID, &doctorDepartmentID, &doctorSpecialization, &doctorLicenseNumber,
		&doctorExperienceYears, &doctorEducation, &doctorCreatedAt, &doctorUpdatedAt,
		// Doctor User fields
		&doctorUserIDScan, &doctorUserName, &doctorUserEmail, &doctorUserCreatedAt, &doctorUserUpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Set nullable fields
	if prescriptionID.Valid {
		record.PrescriptionID = uint(prescriptionID.Int64)
	}
	if followUpDate.Valid {
		record.FollowUpDate = &followUpDate.Time
	}

	// Set patient info
	if patientID.Valid {
		record.Patient = models.Patient{
			ID:               uint(patientID.Int64),
			UserID:           uint(patientUserID.Int64),
			Gender:           patientGender.String,
			BloodType:        patientBloodType.String,
			EmergencyContact: patientEmergencyContact.String,
			EmergencyPhone:   patientEmergencyPhone.String,
			Insurance:        patientInsurance.String,
			InsuranceNo:      patientInsuranceNo.String,
			CreatedAt:        patientCreatedAt.Time,
			UpdatedAt:        patientUpdatedAt.Time,
		}
		if patientDateOfBirth.Valid {
			record.Patient.DateOfBirth = patientDateOfBirth.Time
		}

		// Set patient user info
		if patientUserIDScan.Valid {
			record.Patient.User = models.User{
				ID:        uint(patientUserIDScan.Int64),
				Name:      patientUserName.String,
				Email:     patientUserEmail.String,
				CreatedAt: patientUserCreatedAt.Time,
				UpdatedAt: patientUserUpdatedAt.Time,
			}
		}
	}

	// Set doctor info
	if doctorUserIDScan.Valid {
		record.Doctor.User = models.User{
			ID:        uint(doctorUserIDScan.Int64),
			Name:      doctorUserName.String,
			Email:     doctorUserEmail.String,
			CreatedAt: doctorUserCreatedAt.Time,
			UpdatedAt: doctorUserUpdatedAt.Time,
		}
	}

	return record, nil
}

func (r *MedicalRecordRepository) FindByPatientId(patientId uint) ([]models.MedicalRecord, error) {
	query := `
		SELECT mr.id, mr.patient_id, mr.doctor_id, mr.visit_date, mr.diagnosis, mr.symptoms,
		       mr.treatment, mr.notes, mr.prescription_id, mr.blood_pressure, mr.temperature,
		       mr.weight, mr.height, mr.bmi, mr.pulse_rate, mr.respiratory_rate, mr.allergies,
		       mr.medications, mr.follow_up_date, mr.is_follow_up, mr.attachments, mr.status,
		       -- Doctor info
		       d.id, d.specialization,
		       -- Doctor User info
		       du.name
		FROM medical_records mr
		LEFT JOIN doctors d ON mr.doctor_id = d.id
		LEFT JOIN users du ON d.user_id = du.id
		WHERE mr.patient_id = $1 AND mr.deleted_at IS NULL
		ORDER BY mr.visit_date DESC`

	rows, err := r.db.Query(query, patientId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.MedicalRecord
	for rows.Next() {
		record := models.MedicalRecord{}
		var prescriptionID sql.NullInt64
		var followUpDate sql.NullTime
		var doctorID sql.NullInt64
		var doctorSpecialization, doctorName sql.NullString

		err := rows.Scan(
			&record.ID, &record.PatientID, &record.DoctorID, &record.VisitDate, &record.Diagnosis,
			&record.Symptoms, &record.Treatment, &record.Notes, &prescriptionID, &record.BloodPressure,
			&record.Temperature, &record.Weight, &record.Height, &record.BMI, &record.PulseRate,
			&record.RespiratoryRate, &record.Allergies, &record.Medications, &followUpDate,
			&record.IsFollowUp, &record.Attachments, &record.Status,
			&doctorID, &doctorSpecialization, &doctorName,
		)
		if err != nil {
			return nil, err
		}

		if prescriptionID.Valid {
			record.PrescriptionID = uint(prescriptionID.Int64)
		}
		if followUpDate.Valid {
			record.FollowUpDate = &followUpDate.Time
		}

		// Set basic doctor info
		if doctorID.Valid {
			record.Doctor = models.Doctor{
				ID:             uint(doctorID.Int64),
				Specialization: doctorSpecialization.String,
				User: models.User{
					Name: doctorName.String,
				},
			}
		}

		records = append(records, record)
	}

	return records, rows.Err()
}

func (r *MedicalRecordRepository) FindAll() ([]models.MedicalRecord, error) {
	query := `
		SELECT mr.id, mr.patient_id, mr.doctor_id, mr.visit_date, mr.diagnosis, mr.symptoms,
		       mr.treatment, mr.notes, mr.prescription_id, mr.blood_pressure, mr.temperature,
		       mr.weight, mr.height, mr.bmi, mr.pulse_rate, mr.respiratory_rate, mr.allergies,
		       mr.medications, mr.follow_up_date, mr.is_follow_up, mr.attachments, mr.status,
		       -- Patient info
		       pu.name,
		       -- Doctor info
		       du.name
		FROM medical_records mr
		LEFT JOIN patients p ON mr.patient_id = p.id
		LEFT JOIN users pu ON p.user_id = pu.id
		LEFT JOIN doctors d ON mr.doctor_id = d.id
		LEFT JOIN users du ON d.user_id = du.id
		WHERE mr.deleted_at IS NULL
		ORDER BY mr.visit_date DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.MedicalRecord
	for rows.Next() {
		record := models.MedicalRecord{}
		var prescriptionID sql.NullInt64
		var followUpDate sql.NullTime
		var patientName, doctorName sql.NullString

		err := rows.Scan(
			&record.ID, &record.PatientID, &record.DoctorID, &record.VisitDate, &record.Diagnosis,
			&record.Symptoms, &record.Treatment, &record.Notes, &prescriptionID, &record.BloodPressure,
			&record.Temperature, &record.Weight, &record.Height, &record.BMI, &record.PulseRate,
			&record.RespiratoryRate, &record.Allergies, &record.Medications, &followUpDate,
			&record.IsFollowUp, &record.Attachments, &record.Status,
			&patientName, &doctorName,
		)
		if err != nil {
			return nil, err
		}

		if prescriptionID.Valid {
			record.PrescriptionID = uint(prescriptionID.Int64)
		}
		if followUpDate.Valid {
			record.FollowUpDate = &followUpDate.Time
		}

		// Set basic patient and doctor names
		record.Patient = models.Patient{
			User: models.User{
				Name: patientName.String,
			},
		}
		record.Doctor = models.Doctor{
			User: models.User{
				Name: doctorName.String,
			},
		}

		records = append(records, record)
	}

	return records, rows.Err()
}

func (r *MedicalRecordRepository) Update(record *models.MedicalRecord) error {
	query := `
		UPDATE medical_records SET
			patient_id = $2, doctor_id = $3, visit_date = $4, diagnosis = $5, symptoms = $6,
			treatment = $7, notes = $8, prescription_id = $9, blood_pressure = $10, temperature = $11,
			weight = $12, height = $13, bmi = $14, pulse_rate = $15, respiratory_rate = $16,
			allergies = $17, medications = $18, follow_up_date = $19, is_follow_up = $20,
			attachments = $21, status = $22
		WHERE id = $1 AND deleted_at IS NULL`

	var prescriptionID sql.NullInt64
	if record.PrescriptionID > 0 {
		prescriptionID.Int64 = int64(record.PrescriptionID)
		prescriptionID.Valid = true
	}

	_, err := r.db.Exec(
		query,
		record.ID, record.PatientID, record.DoctorID, record.VisitDate, record.Diagnosis,
		record.Symptoms, record.Treatment, record.Notes, prescriptionID, record.BloodPressure,
		record.Temperature, record.Weight, record.Height, record.BMI, record.PulseRate,
		record.RespiratoryRate, record.Allergies, record.Medications, record.FollowUpDate,
		record.IsFollowUp, record.Attachments, record.Status,
	)

	return err
}

func (r *MedicalRecordRepository) Delete(id uint) error {
	query := "UPDATE medical_records SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL"
	_, err := r.db.Exec(query, time.Now(), id)
	return err
}

// Additional helper methods

func (r *MedicalRecordRepository) FindByDoctorId(doctorId uint) ([]models.MedicalRecord, error) {
	query := `
		SELECT mr.id, mr.patient_id, mr.doctor_id, mr.visit_date, mr.diagnosis, mr.treatment,
		       mr.follow_up_date, mr.status,
		       pu.name
		FROM medical_records mr
		LEFT JOIN patients p ON mr.patient_id = p.id
		LEFT JOIN users pu ON p.user_id = pu.id
		WHERE mr.doctor_id = $1 AND mr.deleted_at IS NULL
		ORDER BY mr.visit_date DESC`

	rows, err := r.db.Query(query, doctorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.MedicalRecord
	for rows.Next() {
		record := models.MedicalRecord{}
		var followUpDate sql.NullTime
		var patientName sql.NullString

		err := rows.Scan(
			&record.ID, &record.PatientID, &record.DoctorID, &record.VisitDate, &record.Diagnosis,
			&record.Treatment, &followUpDate, &record.Status, &patientName,
		)
		if err != nil {
			return nil, err
		}

		if followUpDate.Valid {
			record.FollowUpDate = &followUpDate.Time
		}

		record.Patient = models.Patient{
			User: models.User{
				Name: patientName.String,
			},
		}

		records = append(records, record)
	}

	return records, rows.Err()
}

func (r *MedicalRecordRepository) FindUpcomingFollowUps() ([]models.MedicalRecord, error) {
	query := `
		SELECT mr.id, mr.patient_id, mr.doctor_id, mr.visit_date, mr.diagnosis,
		       mr.follow_up_date, mr.status,
		       pu.name, du.name
		FROM medical_records mr
		LEFT JOIN patients p ON mr.patient_id = p.id
		LEFT JOIN users pu ON p.user_id = pu.id
		LEFT JOIN doctors d ON mr.doctor_id = d.id
		LEFT JOIN users du ON d.user_id = du.id
		WHERE mr.follow_up_date IS NOT NULL
		  AND mr.follow_up_date >= CURRENT_DATE
		  AND mr.follow_up_date <= CURRENT_DATE + INTERVAL '30 days'
		  AND mr.deleted_at IS NULL
		ORDER BY mr.follow_up_date ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.MedicalRecord
	for rows.Next() {
		record := models.MedicalRecord{}
		var followUpDate sql.NullTime
		var patientName, doctorName sql.NullString

		err := rows.Scan(
			&record.ID, &record.PatientID, &record.DoctorID, &record.VisitDate, &record.Diagnosis,
			&followUpDate, &record.Status, &patientName, &doctorName,
		)
		if err != nil {
			return nil, err
		}

		if followUpDate.Valid {
			record.FollowUpDate = &followUpDate.Time
		}

		record.Patient = models.Patient{
			User: models.User{
				Name: patientName.String,
			},
		}
		record.Doctor = models.Doctor{
			User: models.User{
				Name: doctorName.String,
			},
		}

		records = append(records, record)
	}

	return records, rows.Err()
}
