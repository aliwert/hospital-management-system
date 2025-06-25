package repositories

import (
	"database/sql"
	"time"

	"github.com/aliwert/go-hospital-management/internal/models"
)

type DoctorScheduleRepository struct {
	db *sql.DB
}

func NewDoctorScheduleRepository(db *sql.DB) *DoctorScheduleRepository {
	return &DoctorScheduleRepository{db: db}
}

func (r *DoctorScheduleRepository) Create(schedule *models.DoctorSchedule) error {
	query := `
		INSERT INTO doctor_schedules (doctor_id, week_day, start_time, end_time, is_available,
		                             max_appointments, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(
		query,
		schedule.DoctorID,
		schedule.WeekDay,
		schedule.StartTime,
		schedule.EndTime,
		schedule.IsActive,
		schedule.MaxAppointments,
		time.Now(),
		time.Now(),
	).Scan(&schedule.ID, &schedule.CreatedAt, &schedule.UpdatedAt)

	return err
}

func (r *DoctorScheduleRepository) FindById(id uint) (*models.DoctorSchedule, error) {
	query := `
		SELECT ds.id, ds.doctor_id, ds.week_day, ds.start_time, ds.end_time,
		       ds.is_available, ds.max_appointments, ds.created_at, ds.updated_at,
		       d.name as doctor_name, d.specialization as doctor_specialization
		FROM doctor_schedules ds
		LEFT JOIN doctors d ON ds.doctor_id = d.id
		WHERE ds.id = $1 AND ds.deleted_at IS NULL`

	var schedule models.DoctorSchedule
	var doctorName, doctorSpecialization sql.NullString
	var maxAppointments sql.NullInt64

	err := r.db.QueryRow(query, id).Scan(
		&schedule.ID,
		&schedule.DoctorID,
		&schedule.WeekDay,
		&schedule.StartTime,
		&schedule.EndTime,
		&schedule.IsActive,
		&maxAppointments,
		&schedule.CreatedAt,
		&schedule.UpdatedAt,
		&doctorName,
		&doctorSpecialization,
	)

	if err != nil {
		return nil, err
	}

	// Set nullable fields
	if maxAppointments.Valid {
		schedule.MaxAppointments = int(maxAppointments.Int64)
	}

	// Set doctor info
	if doctorName.Valid {
		schedule.Doctor.Name = doctorName.String
	}
	if doctorSpecialization.Valid {
		schedule.Doctor.Specialization = doctorSpecialization.String
	}
	schedule.Doctor.ID = schedule.DoctorID

	return &schedule, nil
}

func (r *DoctorScheduleRepository) FindByDoctorId(doctorId uint) ([]models.DoctorSchedule, error) {
	query := `
		SELECT ds.id, ds.doctor_id, ds.week_day, ds.start_time, ds.end_time,
		       ds.is_available, ds.max_appointments, ds.created_at, ds.updated_at,
		       d.name as doctor_name, d.specialization as doctor_specialization
		FROM doctor_schedules ds
		LEFT JOIN doctors d ON ds.doctor_id = d.id
		WHERE ds.doctor_id = $1 AND ds.deleted_at IS NULL
		ORDER BY ds.week_day, ds.start_time`

	rows, err := r.db.Query(query, doctorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.DoctorSchedule
	for rows.Next() {
		var schedule models.DoctorSchedule
		var doctorName, doctorSpecialization sql.NullString
		var maxAppointments sql.NullInt64

		err := rows.Scan(
			&schedule.ID,
			&schedule.DoctorID,
			&schedule.WeekDay,
			&schedule.StartTime,
			&schedule.EndTime,
			&schedule.IsActive,
			&maxAppointments,
			&schedule.CreatedAt,
			&schedule.UpdatedAt,
			&doctorName,
			&doctorSpecialization,
		)
		if err != nil {
			return nil, err
		}

		// Set nullable fields
		if maxAppointments.Valid {
			schedule.MaxAppointments = int(maxAppointments.Int64)
		}

		// Set doctor info
		if doctorName.Valid {
			schedule.Doctor.Name = doctorName.String
		}
		if doctorSpecialization.Valid {
			schedule.Doctor.Specialization = doctorSpecialization.String
		}
		schedule.Doctor.ID = schedule.DoctorID

		schedules = append(schedules, schedule)
	}

	return schedules, rows.Err()
}

func (r *DoctorScheduleRepository) Update(schedule *models.DoctorSchedule) error {
	query := `
		UPDATE doctor_schedules
		SET doctor_id = $2, week_day = $3, start_time = $4, end_time = $5,
		    is_available = $6, max_appointments = $7, updated_at = $8
		WHERE id = $1 AND deleted_at IS NULL`

	_, err := r.db.Exec(
		query,
		schedule.ID,
		schedule.DoctorID,
		schedule.WeekDay,
		schedule.StartTime,
		schedule.EndTime,
		schedule.IsActive,
		schedule.MaxAppointments,
		time.Now(),
	)

	return err
}

func (r *DoctorScheduleRepository) Delete(id uint) error {
	query := `UPDATE doctor_schedules SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	_, err := r.db.Exec(query, time.Now(), id)
	return err
}

// FindByWeekDay finds schedules for a specific day of week
func (r *DoctorScheduleRepository) FindByWeekDay(weekDay int) ([]models.DoctorSchedule, error) {
	query := `
		SELECT ds.id, ds.doctor_id, ds.week_day, ds.start_time, ds.end_time,
		       ds.is_available, ds.max_appointments, ds.created_at, ds.updated_at,
		       d.name as doctor_name, d.specialization as doctor_specialization
		FROM doctor_schedules ds
		LEFT JOIN doctors d ON ds.doctor_id = d.id
		WHERE ds.week_day = $1 AND ds.is_available = true AND ds.deleted_at IS NULL
		ORDER BY ds.start_time`

	rows, err := r.db.Query(query, weekDay)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.DoctorSchedule
	for rows.Next() {
		var schedule models.DoctorSchedule
		var doctorName, doctorSpecialization sql.NullString
		var maxAppointments sql.NullInt64

		err := rows.Scan(
			&schedule.ID,
			&schedule.DoctorID,
			&schedule.WeekDay,
			&schedule.StartTime,
			&schedule.EndTime,
			&schedule.IsActive,
			&maxAppointments,
			&schedule.CreatedAt,
			&schedule.UpdatedAt,
			&doctorName,
			&doctorSpecialization,
		)
		if err != nil {
			return nil, err
		}

		// Set nullable fields
		if maxAppointments.Valid {
			schedule.MaxAppointments = int(maxAppointments.Int64)
		}

		// Set doctor info
		if doctorName.Valid {
			schedule.Doctor.Name = doctorName.String
		}
		if doctorSpecialization.Valid {
			schedule.Doctor.Specialization = doctorSpecialization.String
		}
		schedule.Doctor.ID = schedule.DoctorID

		schedules = append(schedules, schedule)
	}

	return schedules, rows.Err()
}
