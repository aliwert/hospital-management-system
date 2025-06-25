package repositories

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/aliwert/go-hospital-management/internal/models"
)

type AppointmentRepository struct {
	db *sql.DB
}

func NewAppointmentRepository(db *sql.DB) *AppointmentRepository {
	return &AppointmentRepository{db: db}
}

func (r *AppointmentRepository) Create(appointment *models.Appointment) error {
	query := `
		INSERT INTO appointments (patient_id, doctor_id, appointment_date, status, description, fee, payment_status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(
		query,
		appointment.PatientID,
		appointment.DoctorID,
		appointment.AppointmentDate,
		appointment.Status,
		appointment.Description,
		appointment.Fee,
		appointment.PaymentStatus,
		time.Now(),
		time.Now(),
	).Scan(&appointment.ID, &appointment.CreatedAt, &appointment.UpdatedAt)

	return err
}

func (r *AppointmentRepository) FindById(id uint) (*models.Appointment, error) {
	query := r.getBaseQuery() + " WHERE a.id = $1 AND a.deleted_at IS NULL"

	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	return r.scanAppointment(rows)
}

func (r *AppointmentRepository) FindAll() ([]models.Appointment, error) {
	query := r.getBaseQuery() + " WHERE a.deleted_at IS NULL ORDER BY a.appointment_date DESC"

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []models.Appointment
	for rows.Next() {
		appointment, err := r.scanAppointment(rows)
		if err != nil {
			return nil, err
		}
		appointments = append(appointments, *appointment)
	}

	return appointments, rows.Err()
}

func (r *AppointmentRepository) FindByPatientId(patientId uint) ([]models.Appointment, error) {
	query := r.getBaseQuery() + " WHERE a.patient_id = $1 AND a.deleted_at IS NULL ORDER BY a.appointment_date DESC"

	rows, err := r.db.Query(query, patientId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []models.Appointment
	for rows.Next() {
		appointment, err := r.scanAppointment(rows)
		if err != nil {
			return nil, err
		}
		appointments = append(appointments, *appointment)
	}

	return appointments, rows.Err()
}

func (r *AppointmentRepository) FindByDoctorId(doctorId uint) ([]models.Appointment, error) {
	query := r.getBaseQuery() + " WHERE a.doctor_id = $1 AND a.deleted_at IS NULL ORDER BY a.appointment_date DESC"

	rows, err := r.db.Query(query, doctorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []models.Appointment
	for rows.Next() {
		appointment, err := r.scanAppointment(rows)
		if err != nil {
			return nil, err
		}
		appointments = append(appointments, *appointment)
	}

	return appointments, rows.Err()
}

func (r *AppointmentRepository) Update(appointment *models.Appointment) error {
	query := `
		UPDATE appointments
		SET status = $2, description = $3, payment_status = $4, payment_date = $5,
		    cancelled_at = $6, cancel_reason = $7, notes = $8, updated_at = $9
		WHERE id = $1 AND deleted_at IS NULL`

	_, err := r.db.Exec(
		query,
		appointment.ID,
		appointment.Status,
		appointment.Description,
		appointment.PaymentStatus,
		appointment.PaymentDate,
		appointment.CancelledAt,
		appointment.CancelReason,
		appointment.Notes,
		time.Now(),
	)

	return err
}

func (r *AppointmentRepository) Delete(id uint) error {
	query := `UPDATE appointments SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	_, err := r.db.Exec(query, time.Now(), id)
	return err
}

// scanAppointment is a helper function to scan appointment rows with patient and doctor info
func (r *AppointmentRepository) scanAppointment(rows *sql.Rows) (*models.Appointment, error) {
	var appointment models.Appointment
	var patientName, patientEmail sql.NullString
	var doctorName, doctorSpecialization sql.NullString
	var paymentDate, cancelledAt sql.NullTime
	var cancelReason, notes sql.NullString

	err := rows.Scan(
		&appointment.ID,
		&appointment.PatientID,
		&appointment.DoctorID,
		&appointment.AppointmentDate,
		&appointment.Status,
		&appointment.Description,
		&appointment.Fee,
		&appointment.PaymentStatus,
		&paymentDate,
		&cancelledAt,
		&cancelReason,
		&notes,
		&appointment.CreatedAt,
		&appointment.UpdatedAt,
		&patientName,
		&patientEmail,
		&doctorName,
		&doctorSpecialization,
	)
	if err != nil {
		return nil, err
	}

	// Set nullable fields
	if paymentDate.Valid {
		appointment.PaymentDate = &paymentDate.Time
	}
	if cancelledAt.Valid {
		appointment.CancelledAt = &cancelledAt.Time
	}
	if cancelReason.Valid {
		appointment.CancelReason = cancelReason.String
	}
	if notes.Valid {
		appointment.Notes = notes.String
	}

	// Set patient info
	if patientName.Valid {
		appointment.Patient.Name = patientName.String
	}
	if patientEmail.Valid {
		appointment.Patient.Email = patientEmail.String
	}
	appointment.Patient.ID = appointment.PatientID

	// Set doctor info
	if doctorName.Valid {
		appointment.Doctor.Name = doctorName.String
	}
	if doctorSpecialization.Valid {
		appointment.Doctor.Specialization = doctorSpecialization.String
	}
	appointment.Doctor.ID = appointment.DoctorID

	return &appointment, nil
}

// getBaseQuery returns the base SELECT query with JOINs
func (r *AppointmentRepository) getBaseQuery() string {
	return `
		SELECT
			a.id, a.patient_id, a.doctor_id, a.appointment_date, a.status, a.description,
			a.fee, a.payment_status, a.payment_date, a.cancelled_at, a.cancel_reason,
			a.notes, a.created_at, a.updated_at,
			p.name as patient_name, p.email as patient_email,
			d.name as doctor_name, d.specialization as doctor_specialization
		FROM appointments a
		LEFT JOIN users p ON a.patient_id = p.id
		LEFT JOIN doctors d ON a.doctor_id = d.id`
}

// FindWithFilters allows for more complex filtering with pagination
func (r *AppointmentRepository) FindWithFilters(filters map[string]interface{}, limit, offset int) ([]models.Appointment, int, error) {
	baseQuery := r.getBaseQuery()
	whereClause := "WHERE a.deleted_at IS NULL"
	args := []interface{}{}
	argIndex := 1

	// Build dynamic WHERE clause
	for key, value := range filters {
		switch key {
		case "patient_id":
			whereClause += fmt.Sprintf(" AND a.patient_id = $%d", argIndex)
			args = append(args, value)
			argIndex++
		case "doctor_id":
			whereClause += fmt.Sprintf(" AND a.doctor_id = $%d", argIndex)
			args = append(args, value)
			argIndex++
		case "status":
			whereClause += fmt.Sprintf(" AND a.status = $%d", argIndex)
			args = append(args, value)
			argIndex++
		case "payment_status":
			whereClause += fmt.Sprintf(" AND a.payment_status = $%d", argIndex)
			args = append(args, value)
			argIndex++
		case "date_from":
			whereClause += fmt.Sprintf(" AND a.appointment_date >= $%d", argIndex)
			args = append(args, value)
			argIndex++
		case "date_to":
			whereClause += fmt.Sprintf(" AND a.appointment_date <= $%d", argIndex)
			args = append(args, value)
			argIndex++
		}
	}

	// Count query for pagination
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM appointments a %s", whereClause)
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Main query with pagination
	orderBy := " ORDER BY a.appointment_date DESC"
	limitOffset := fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	query := baseQuery + " " + whereClause + orderBy + limitOffset

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var appointments []models.Appointment
	for rows.Next() {
		appointment, err := r.scanAppointment(rows)
		if err != nil {
			return nil, 0, err
		}
		appointments = append(appointments, *appointment)
	}

	return appointments, total, rows.Err()
}

// FindUpcoming finds upcoming appointments for a doctor or patient with performance optimization
func (r *AppointmentRepository) FindUpcoming(userType string, userID uint, limit int) ([]models.Appointment, error) {
	baseQuery := r.getBaseQuery()
	var whereClause string

	switch userType {
	case "doctor":
		whereClause = "WHERE a.doctor_id = $1 AND a.appointment_date > $2 AND a.status IN ('pending', 'confirmed') AND a.deleted_at IS NULL"
	case "patient":
		whereClause = "WHERE a.patient_id = $1 AND a.appointment_date > $2 AND a.status IN ('pending', 'confirmed') AND a.deleted_at IS NULL"
	default:
		return nil, fmt.Errorf("invalid user type: %s", userType)
	}

	orderBy := " ORDER BY a.appointment_date ASC"
	limitClause := fmt.Sprintf(" LIMIT $3")

	query := baseQuery + " " + whereClause + orderBy + limitClause

	rows, err := r.db.Query(query, userID, time.Now(), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []models.Appointment
	for rows.Next() {
		appointment, err := r.scanAppointment(rows)
		if err != nil {
			return nil, err
		}
		appointments = append(appointments, *appointment)
	}

	return appointments, rows.Err()
}

// BatchUpdateStatus updates multiple appointments status efficiently
func (r *AppointmentRepository) BatchUpdateStatus(ids []uint, status string) error {
	if len(ids) == 0 {
		return nil
	}

	// Convert IDs to string for IN clause
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids)+2)

	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	args[len(ids)] = status
	args[len(ids)+1] = time.Now()

	query := fmt.Sprintf(`
		UPDATE appointments
		SET status = $%d, updated_at = $%d
		WHERE id IN (%s) AND deleted_at IS NULL`,
		len(ids)+1, len(ids)+2, strings.Join(placeholders, ","))

	_, err := r.db.Exec(query, args...)
	return err
}

// GetAppointmentStats returns appointment statistics for dashboard
func (r *AppointmentRepository) GetAppointmentStats(doctorID *uint, dateFrom, dateTo time.Time) (map[string]interface{}, error) {
	baseWhere := "WHERE a.deleted_at IS NULL AND a.appointment_date BETWEEN $1 AND $2"
	args := []interface{}{dateFrom, dateTo}
	argIndex := 3

	if doctorID != nil {
		baseWhere += fmt.Sprintf(" AND a.doctor_id = $%d", argIndex)
		args = append(args, *doctorID)
	}

	query := fmt.Sprintf(`
		SELECT
			COUNT(*) as total,
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending,
			COUNT(CASE WHEN status = 'confirmed' THEN 1 END) as confirmed,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed,
			COUNT(CASE WHEN status = 'cancelled' THEN 1 END) as cancelled,
			COALESCE(SUM(CASE WHEN payment_status = 'paid' THEN fee ELSE 0 END), 0) as total_revenue,
			COUNT(CASE WHEN payment_status = 'paid' THEN 1 END) as paid_appointments
		FROM appointments a %s`, baseWhere)

	var stats struct {
		Total            int     `json:"total"`
		Pending          int     `json:"pending"`
		Confirmed        int     `json:"confirmed"`
		Completed        int     `json:"completed"`
		Cancelled        int     `json:"cancelled"`
		TotalRevenue     float64 `json:"total_revenue"`
		PaidAppointments int     `json:"paid_appointments"`
	}

	err := r.db.QueryRow(query, args...).Scan(
		&stats.Total,
		&stats.Pending,
		&stats.Confirmed,
		&stats.Completed,
		&stats.Cancelled,
		&stats.TotalRevenue,
		&stats.PaidAppointments,
	)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"total":             stats.Total,
		"pending":           stats.Pending,
		"confirmed":         stats.Confirmed,
		"completed":         stats.Completed,
		"cancelled":         stats.Cancelled,
		"total_revenue":     stats.TotalRevenue,
		"paid_appointments": stats.PaidAppointments,
	}

	return result, nil
}
