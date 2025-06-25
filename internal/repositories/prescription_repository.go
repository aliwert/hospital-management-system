package repositories

import (
	"database/sql"
	"time"

	"github.com/aliwert/go-hospital-management/internal/models"
	_ "github.com/lib/pq"
)

type PrescriptionRepository struct {
	db *sql.DB
}

func NewPrescriptionRepository(db *sql.DB) *PrescriptionRepository {
	return &PrescriptionRepository{db: db}
}

func (r *PrescriptionRepository) Create(prescription *models.Prescription) error {
	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert prescription
	query := `
		INSERT INTO prescriptions (
			patient_id, doctor_id, diagnosis, notes, issue_date, valid_until,
			status, pharmacy_id, filled_date, refill_count, max_refills,
			is_controlled, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, created_at, updated_at`

	now := time.Now()
	err = tx.QueryRow(
		query,
		prescription.PatientID, prescription.DoctorID, prescription.Diagnosis, prescription.Notes,
		prescription.IssueDate, prescription.ValidUntil, prescription.Status, prescription.PharmacyID,
		prescription.FilledDate, prescription.RefillCount, prescription.MaxRefills,
		prescription.IsControlled, now, now,
	).Scan(&prescription.ID, &prescription.CreatedAt, &prescription.UpdatedAt)

	if err != nil {
		return err
	}

	// Insert medications
	if len(prescription.Medications) > 0 {
		medicationQuery := `
			INSERT INTO prescription_medications (
				prescription_id, medicine_name, dosage, frequency, duration,
				instructions, quantity, substitution, created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

		for i := range prescription.Medications {
			medication := &prescription.Medications[i]
			medication.PrescriptionID = prescription.ID
			medication.CreatedAt = now
			medication.UpdatedAt = now

			_, err = tx.Exec(
				medicationQuery,
				medication.PrescriptionID, medication.MedicineName, medication.Dosage,
				medication.Frequency, medication.Duration, medication.Instructions,
				medication.Quantity, medication.Substitution, medication.CreatedAt, medication.UpdatedAt,
			)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *PrescriptionRepository) FindById(id uint) (*models.Prescription, error) {
	prescription := &models.Prescription{}

	query := `
		SELECT p.id, p.patient_id, p.doctor_id, p.diagnosis, p.notes, p.issue_date,
		       p.valid_until, p.status, p.pharmacy_id, p.filled_date, p.refill_count,
		       p.max_refills, p.is_controlled, p.created_at, p.updated_at,
		       -- Patient info
		       pat.id, pat.user_id,
		       -- Patient User info
		       pu.first_name, pu.last_name, pu.email, pu.phone,
		       -- Doctor info
		       d.id, d.user_id, d.specialization,
		       -- Doctor User info
		       du.first_name, du.last_name, du.email, du.phone
		FROM prescriptions p
		LEFT JOIN patients pat ON p.patient_id = pat.id
		LEFT JOIN users pu ON pat.user_id = pu.id
		LEFT JOIN doctors d ON p.doctor_id = d.id
		LEFT JOIN users du ON d.user_id = du.id
		WHERE p.id = $1 AND p.deleted_at IS NULL`

	var pharmacyID sql.NullInt64
	var filledDate sql.NullTime

	var patientID, patientUserID sql.NullInt64
	var patientFirstName, patientLastName, patientEmail, patientPhone sql.NullString

	var doctorID, doctorUserID sql.NullInt64
	var doctorSpecialization, doctorFirstName, doctorLastName sql.NullString
	var doctorEmail, doctorPhone sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&prescription.ID, &prescription.PatientID, &prescription.DoctorID, &prescription.Diagnosis,
		&prescription.Notes, &prescription.IssueDate, &prescription.ValidUntil, &prescription.Status,
		&pharmacyID, &filledDate, &prescription.RefillCount, &prescription.MaxRefills,
		&prescription.IsControlled, &prescription.CreatedAt, &prescription.UpdatedAt,
		// Patient fields
		&patientID, &patientUserID, &patientFirstName, &patientLastName, &patientEmail, &patientPhone,
		// Doctor fields
		&doctorID, &doctorUserID, &doctorSpecialization, &doctorFirstName, &doctorLastName,
		&doctorEmail, &doctorPhone,
	)

	if err != nil {
		return nil, err
	}

	// Set nullable fields
	if pharmacyID.Valid {
		pharmacyIDUint := uint(pharmacyID.Int64)
		prescription.PharmacyID = &pharmacyIDUint
	}
	if filledDate.Valid {
		prescription.FilledDate = &filledDate.Time
	}

	// Set patient info
	if patientID.Valid {
		prescription.Patient = models.Patient{
			ID:     uint(patientID.Int64),
			UserID: uint(patientUserID.Int64),
			User: models.User{
				Name:  patientFirstName.String,
				Email: patientEmail.String,
			},
		}
	}

	// Set doctor info
	if doctorID.Valid {
		prescription.Doctor = models.Doctor{
			ID:             uint(doctorID.Int64),
			UserID:         uint(doctorUserID.Int64),
			Specialization: doctorSpecialization.String,
			User: models.User{
				Name:  doctorFirstName.String,
				Email: doctorEmail.String,
			},
		}
	}

	// Load medications
	medications, err := r.findMedicationsByPrescriptionId(prescription.ID)
	if err != nil {
		return nil, err
	}
	prescription.Medications = medications

	return prescription, nil
}

func (r *PrescriptionRepository) FindByPatientId(patientId uint) ([]models.Prescription, error) {
	query := `
		SELECT p.id, p.patient_id, p.doctor_id, p.diagnosis, p.issue_date,
		       p.valid_until, p.status, p.refill_count, p.max_refills,
		       p.created_at, p.updated_at,
		       -- Doctor info
		       du.first_name, du.last_name, d.specialization
		FROM prescriptions p
		LEFT JOIN doctors d ON p.doctor_id = d.id
		LEFT JOIN users du ON d.user_id = du.id
		WHERE p.patient_id = $1 AND p.deleted_at IS NULL
		ORDER BY p.issue_date DESC`

	rows, err := r.db.Query(query, patientId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prescriptions []models.Prescription
	for rows.Next() {
		prescription := models.Prescription{}
		var doctorFirstName, doctorLastName, doctorSpecialization sql.NullString

		err := rows.Scan(
			&prescription.ID, &prescription.PatientID, &prescription.DoctorID, &prescription.Diagnosis,
			&prescription.IssueDate, &prescription.ValidUntil, &prescription.Status,
			&prescription.RefillCount, &prescription.MaxRefills, &prescription.CreatedAt,
			&prescription.UpdatedAt, &doctorFirstName, &doctorLastName, &doctorSpecialization,
		)
		if err != nil {
			return nil, err
		}

		// Set basic doctor info
		prescription.Doctor = models.Doctor{
			Specialization: doctorSpecialization.String,
			User: models.User{
				Name: doctorFirstName.String,
			},
		}

		// Load medications for this prescription
		medications, err := r.findMedicationsByPrescriptionId(prescription.ID)
		if err != nil {
			return nil, err
		}
		prescription.Medications = medications

		prescriptions = append(prescriptions, prescription)
	}

	return prescriptions, rows.Err()
}

func (r *PrescriptionRepository) FindAll() ([]models.Prescription, error) {
	query := `
		SELECT p.id, p.patient_id, p.doctor_id, p.diagnosis, p.issue_date,
		       p.valid_until, p.status, p.refill_count, p.max_refills,
		       p.created_at, p.updated_at,
		       -- Patient info
		       pu.first_name, pu.last_name,
		       -- Doctor info
		       du.first_name, du.last_name
		FROM prescriptions p
		LEFT JOIN patients pat ON p.patient_id = pat.id
		LEFT JOIN users pu ON pat.user_id = pu.id
		LEFT JOIN doctors d ON p.doctor_id = d.id
		LEFT JOIN users du ON d.user_id = du.id
		WHERE p.deleted_at IS NULL
		ORDER BY p.issue_date DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prescriptions []models.Prescription
	for rows.Next() {
		prescription := models.Prescription{}
		var patientFirstName, patientLastName sql.NullString
		var doctorFirstName, doctorLastName sql.NullString

		err := rows.Scan(
			&prescription.ID, &prescription.PatientID, &prescription.DoctorID, &prescription.Diagnosis,
			&prescription.IssueDate, &prescription.ValidUntil, &prescription.Status,
			&prescription.RefillCount, &prescription.MaxRefills, &prescription.CreatedAt,
			&prescription.UpdatedAt, &patientFirstName, &patientLastName,
			&doctorFirstName, &doctorLastName,
		)
		if err != nil {
			return nil, err
		}

		// Set basic patient and doctor names
		prescription.Patient = models.Patient{
			User: models.User{
				Name: patientFirstName.String,
			},
		}
		prescription.Doctor = models.Doctor{
			User: models.User{
				Name: doctorFirstName.String,
			},
		}

		prescriptions = append(prescriptions, prescription)
	}

	return prescriptions, rows.Err()
}

func (r *PrescriptionRepository) Update(prescription *models.Prescription) error {
	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update prescription
	query := `
		UPDATE prescriptions SET
			patient_id = $2, doctor_id = $3, diagnosis = $4, notes = $5,
			valid_until = $6, status = $7, pharmacy_id = $8, filled_date = $9,
			refill_count = $10, max_refills = $11, is_controlled = $12, updated_at = $13
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING updated_at`

	err = tx.QueryRow(
		query,
		prescription.ID, prescription.PatientID, prescription.DoctorID, prescription.Diagnosis,
		prescription.Notes, prescription.ValidUntil, prescription.Status, prescription.PharmacyID,
		prescription.FilledDate, prescription.RefillCount, prescription.MaxRefills,
		prescription.IsControlled, time.Now(),
	).Scan(&prescription.UpdatedAt)

	if err != nil {
		return err
	}

	// Update medications if provided
	if len(prescription.Medications) > 0 {
		// Delete existing medications
		_, err = tx.Exec("DELETE FROM prescription_medications WHERE prescription_id = $1", prescription.ID)
		if err != nil {
			return err
		}

		// Insert new medications
		medicationQuery := `
			INSERT INTO prescription_medications (
				prescription_id, medicine_name, dosage, frequency, duration,
				instructions, quantity, substitution, created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

		now := time.Now()
		for i := range prescription.Medications {
			medication := &prescription.Medications[i]
			medication.PrescriptionID = prescription.ID
			medication.CreatedAt = now
			medication.UpdatedAt = now

			_, err = tx.Exec(
				medicationQuery,
				medication.PrescriptionID, medication.MedicineName, medication.Dosage,
				medication.Frequency, medication.Duration, medication.Instructions,
				medication.Quantity, medication.Substitution, medication.CreatedAt, medication.UpdatedAt,
			)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *PrescriptionRepository) Delete(id uint) error {
	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now()

	// Soft delete prescription
	_, err = tx.Exec("UPDATE prescriptions SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL", now, id)
	if err != nil {
		return err
	}

	// Note: We don't delete medications as they are linked to prescription history

	return tx.Commit()
}

// Helper method to find medications by prescription ID
func (r *PrescriptionRepository) findMedicationsByPrescriptionId(prescriptionId uint) ([]models.PrescriptionMedication, error) {
	query := `
		SELECT id, prescription_id, medicine_name, dosage, frequency, duration,
		       instructions, quantity, substitution, created_at, updated_at
		FROM prescription_medications
		WHERE prescription_id = $1
		ORDER BY medicine_name`

	rows, err := r.db.Query(query, prescriptionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var medications []models.PrescriptionMedication
	for rows.Next() {
		medication := models.PrescriptionMedication{}

		err := rows.Scan(
			&medication.ID, &medication.PrescriptionID, &medication.MedicineName,
			&medication.Dosage, &medication.Frequency, &medication.Duration,
			&medication.Instructions, &medication.Quantity, &medication.Substitution,
			&medication.CreatedAt, &medication.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		medications = append(medications, medication)
	}

	return medications, rows.Err()
}

// Additional helper methods

func (r *PrescriptionRepository) FindByDoctorId(doctorId uint) ([]models.Prescription, error) {
	query := `
		SELECT p.id, p.patient_id, p.doctor_id, p.diagnosis, p.issue_date,
		       p.valid_until, p.status, p.created_at, p.updated_at,
		       pu.first_name, pu.last_name
		FROM prescriptions p
		LEFT JOIN patients pat ON p.patient_id = pat.id
		LEFT JOIN users pu ON pat.user_id = pu.id
		WHERE p.doctor_id = $1 AND p.deleted_at IS NULL
		ORDER BY p.issue_date DESC`

	rows, err := r.db.Query(query, doctorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prescriptions []models.Prescription
	for rows.Next() {
		prescription := models.Prescription{}
		var patientFirstName, patientLastName sql.NullString

		err := rows.Scan(
			&prescription.ID, &prescription.PatientID, &prescription.DoctorID, &prescription.Diagnosis,
			&prescription.IssueDate, &prescription.ValidUntil, &prescription.Status,
			&prescription.CreatedAt, &prescription.UpdatedAt, &patientFirstName, &patientLastName,
		)
		if err != nil {
			return nil, err
		}

		prescription.Patient = models.Patient{
			User: models.User{
				Name: patientFirstName.String,
			},
		}

		prescriptions = append(prescriptions, prescription)
	}

	return prescriptions, rows.Err()
}

func (r *PrescriptionRepository) FindExpiringSoon(days int) ([]models.Prescription, error) {
	query := `
		SELECT p.id, p.patient_id, p.doctor_id, p.diagnosis, p.valid_until,
		       p.status, pu.first_name, pu.last_name, du.first_name, du.last_name
		FROM prescriptions p
		LEFT JOIN patients pat ON p.patient_id = pat.id
		LEFT JOIN users pu ON pat.user_id = pu.id
		LEFT JOIN doctors d ON p.doctor_id = d.id
		LEFT JOIN users du ON d.user_id = du.id
		WHERE p.valid_until <= CURRENT_DATE + INTERVAL '%d days'
		  AND p.valid_until >= CURRENT_DATE
		  AND p.status = 'active'
		  AND p.deleted_at IS NULL
		ORDER BY p.valid_until ASC`

	rows, err := r.db.Query(query, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prescriptions []models.Prescription
	for rows.Next() {
		prescription := models.Prescription{}
		var patientFirstName, patientLastName sql.NullString
		var doctorFirstName, doctorLastName sql.NullString

		err := rows.Scan(
			&prescription.ID, &prescription.PatientID, &prescription.DoctorID, &prescription.Diagnosis,
			&prescription.ValidUntil, &prescription.Status, &patientFirstName, &patientLastName,
			&doctorFirstName, &doctorLastName,
		)
		if err != nil {
			return nil, err
		}

		prescription.Patient = models.Patient{
			User: models.User{
				Name: patientFirstName.String,
			},
		}
		prescription.Doctor = models.Doctor{
			User: models.User{
				Name: doctorFirstName.String,
			},
		}

		prescriptions = append(prescriptions, prescription)
	}

	return prescriptions, rows.Err()
}

func (r *PrescriptionRepository) MarkAsFilled(id uint, pharmacyId uint) error {
	query := `
		UPDATE prescriptions SET
			status = 'filled', pharmacy_id = $1, filled_date = $2,
			refill_count = refill_count + 1, updated_at = $3
		WHERE id = $4 AND deleted_at IS NULL`

	now := time.Now()
	_, err := r.db.Exec(query, pharmacyId, now, now, id)
	return err
}
