package repositories

import (
	"database/sql"
	"time"

	"github.com/aliwert/go-hospital-management/internal/models"
	_ "github.com/lib/pq"
)

type TestResultRepository struct {
	db *sql.DB
}

func NewTestResultRepository(db *sql.DB) *TestResultRepository {
	return &TestResultRepository{db: db}
}

func (r *TestResultRepository) Create(result *models.TestResult) error {
	query := `
		INSERT INTO test_results (
			medical_record_id, test_name, result, unit, reference_range, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`

	now := time.Now()
	err := r.db.QueryRow(
		query,
		result.MedicalRecordID, result.TestName, result.Result, result.Unit,
		result.ReferenceRange, now, now,
	).Scan(&result.ID, &result.CreatedAt, &result.UpdatedAt)

	return err
}

func (r *TestResultRepository) FindById(id uint) (*models.TestResult, error) {
	result := &models.TestResult{}

	query := `
		SELECT id, medical_record_id, test_name, result, unit, reference_range,
		       created_at, updated_at
		FROM test_results
		WHERE id = $1 AND deleted_at IS NULL`

	err := r.db.QueryRow(query, id).Scan(
		&result.ID, &result.MedicalRecordID, &result.TestName, &result.Result,
		&result.Unit, &result.ReferenceRange, &result.CreatedAt, &result.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *TestResultRepository) FindByMedicalRecordId(medicalRecordId uint) ([]models.TestResult, error) {
	query := `
		SELECT id, medical_record_id, test_name, result, unit, reference_range,
		       created_at, updated_at
		FROM test_results
		WHERE medical_record_id = $1 AND deleted_at IS NULL
		ORDER BY test_name, created_at DESC`

	rows, err := r.db.Query(query, medicalRecordId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.TestResult
	for rows.Next() {
		result := models.TestResult{}

		err := rows.Scan(
			&result.ID, &result.MedicalRecordID, &result.TestName, &result.Result,
			&result.Unit, &result.ReferenceRange, &result.CreatedAt, &result.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	return results, rows.Err()
}

func (r *TestResultRepository) Update(result *models.TestResult) error {
	query := `
		UPDATE test_results SET
			medical_record_id = $2, test_name = $3, result = $4, unit = $5,
			reference_range = $6, updated_at = $7
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING updated_at`

	err := r.db.QueryRow(
		query,
		result.ID, result.MedicalRecordID, result.TestName, result.Result,
		result.Unit, result.ReferenceRange, time.Now(),
	).Scan(&result.UpdatedAt)

	return err
}

func (r *TestResultRepository) Delete(id uint) error {
	query := "UPDATE test_results SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL"
	_, err := r.db.Exec(query, time.Now(), id)
	return err
}

// Additional helper methods

func (r *TestResultRepository) FindAll() ([]models.TestResult, error) {
	query := `
		SELECT id, medical_record_id, test_name, result, unit, reference_range,
		       created_at, updated_at
		FROM test_results
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.TestResult
	for rows.Next() {
		result := models.TestResult{}

		err := rows.Scan(
			&result.ID, &result.MedicalRecordID, &result.TestName, &result.Result,
			&result.Unit, &result.ReferenceRange, &result.CreatedAt, &result.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	return results, rows.Err()
}

func (r *TestResultRepository) FindByTestName(testName string) ([]models.TestResult, error) {
	query := `
		SELECT id, medical_record_id, test_name, result, unit, reference_range,
		       created_at, updated_at
		FROM test_results
		WHERE test_name ILIKE $1 AND deleted_at IS NULL
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, "%"+testName+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.TestResult
	for rows.Next() {
		result := models.TestResult{}

		err := rows.Scan(
			&result.ID, &result.MedicalRecordID, &result.TestName, &result.Result,
			&result.Unit, &result.ReferenceRange, &result.CreatedAt, &result.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	return results, rows.Err()
}

func (r *TestResultRepository) FindByPatientId(patientId uint) ([]models.TestResult, error) {
	query := `
		SELECT tr.id, tr.medical_record_id, tr.test_name, tr.result, tr.unit,
		       tr.reference_range, tr.created_at, tr.updated_at
		FROM test_results tr
		INNER JOIN medical_records mr ON tr.medical_record_id = mr.id
		WHERE mr.patient_id = $1 AND tr.deleted_at IS NULL AND mr.deleted_at IS NULL
		ORDER BY tr.created_at DESC`

	rows, err := r.db.Query(query, patientId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.TestResult
	for rows.Next() {
		result := models.TestResult{}

		err := rows.Scan(
			&result.ID, &result.MedicalRecordID, &result.TestName, &result.Result,
			&result.Unit, &result.ReferenceRange, &result.CreatedAt, &result.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	return results, rows.Err()
}

func (r *TestResultRepository) FindAbnormalResults() ([]models.TestResult, error) {
	query := `
		SELECT id, medical_record_id, test_name, result, unit, reference_range,
		       created_at, updated_at
		FROM test_results
		WHERE (result ILIKE '%abnormal%' OR result ILIKE '%high%' OR result ILIKE '%low%'
		       OR result ILIKE '%critical%' OR result ILIKE '%elevated%')
		  AND deleted_at IS NULL
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.TestResult
	for rows.Next() {
		result := models.TestResult{}

		err := rows.Scan(
			&result.ID, &result.MedicalRecordID, &result.TestName, &result.Result,
			&result.Unit, &result.ReferenceRange, &result.CreatedAt, &result.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	return results, rows.Err()
}

func (r *TestResultRepository) BatchCreate(results []models.TestResult) error {
	if len(results) == 0 {
		return nil
	}

	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO test_results (
			medical_record_id, test_name, result, unit, reference_range, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	now := time.Now()
	for i := range results {
		result := &results[i]
		result.CreatedAt = now
		result.UpdatedAt = now

		_, err = tx.Exec(
			query,
			result.MedicalRecordID, result.TestName, result.Result, result.Unit,
			result.ReferenceRange, result.CreatedAt, result.UpdatedAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
