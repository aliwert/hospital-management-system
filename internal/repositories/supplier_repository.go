package repositories

import (
	"database/sql"
	"errors"
	"time"

	"github.com/aliwert/go-hospital-management/internal/models"
	_ "github.com/lib/pq"
)

type SupplierRepository struct {
	db *sql.DB
}

func NewSupplierRepository(db *sql.DB) *SupplierRepository {
	return &SupplierRepository{db: db}
}

func (r *SupplierRepository) Create(supplier *models.Supplier) error {
	// Check if supplier code already exists
	var count int
	checkQuery := "SELECT COUNT(*) FROM suppliers WHERE code = $1 AND deleted_at IS NULL"
	err := r.db.QueryRow(checkQuery, supplier.Code).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("supplier code already exists")
	}

	query := `
		INSERT INTO suppliers (
			name, code, email, phone, address, contact_person, contact_phone,
			tax_number, bank_account, payment_terms, delivery_terms, website,
			rating, status, notes, total_orders, is_verified, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
		RETURNING id, created_at, updated_at`

	now := time.Now()
	err = r.db.QueryRow(
		query,
		supplier.Name, supplier.Code, supplier.Email, supplier.Phone, supplier.Address,
		supplier.ContactPerson, supplier.ContactPhone, supplier.TaxNumber, supplier.BankAccount,
		supplier.PaymentTerms, supplier.DeliveryTerms, supplier.Website, supplier.Rating,
		supplier.Status, supplier.Notes, supplier.TotalOrders, supplier.IsVerified, now, now,
	).Scan(&supplier.ID, &supplier.CreatedAt, &supplier.UpdatedAt)

	return err
}

func (r *SupplierRepository) FindById(id uint) (*models.Supplier, error) {
	supplier := &models.Supplier{}

	query := `
		SELECT id, name, code, email, phone, address, contact_person, contact_phone,
		       tax_number, bank_account, payment_terms, delivery_terms, website,
		       rating, status, notes, last_order_date, total_orders, is_verified,
		       created_at, updated_at
		FROM suppliers
		WHERE id = $1 AND deleted_at IS NULL`

	var lastOrderDate sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&supplier.ID, &supplier.Name, &supplier.Code, &supplier.Email, &supplier.Phone,
		&supplier.Address, &supplier.ContactPerson, &supplier.ContactPhone, &supplier.TaxNumber,
		&supplier.BankAccount, &supplier.PaymentTerms, &supplier.DeliveryTerms, &supplier.Website,
		&supplier.Rating, &supplier.Status, &supplier.Notes, &lastOrderDate, &supplier.TotalOrders,
		&supplier.IsVerified, &supplier.CreatedAt, &supplier.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if lastOrderDate.Valid {
		supplier.LastOrderDate = &lastOrderDate.Time
	}

	return supplier, nil
}

func (r *SupplierRepository) FindAll() ([]models.Supplier, error) {
	query := `
		SELECT id, name, code, email, phone, address, contact_person, contact_phone,
		       tax_number, bank_account, payment_terms, delivery_terms, website,
		       rating, status, notes, last_order_date, total_orders, is_verified,
		       created_at, updated_at
		FROM suppliers
		WHERE deleted_at IS NULL
		ORDER BY name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suppliers []models.Supplier
	for rows.Next() {
		supplier := models.Supplier{}
		var lastOrderDate sql.NullTime

		err := rows.Scan(
			&supplier.ID, &supplier.Name, &supplier.Code, &supplier.Email, &supplier.Phone,
			&supplier.Address, &supplier.ContactPerson, &supplier.ContactPhone, &supplier.TaxNumber,
			&supplier.BankAccount, &supplier.PaymentTerms, &supplier.DeliveryTerms, &supplier.Website,
			&supplier.Rating, &supplier.Status, &supplier.Notes, &lastOrderDate, &supplier.TotalOrders,
			&supplier.IsVerified, &supplier.CreatedAt, &supplier.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if lastOrderDate.Valid {
			supplier.LastOrderDate = &lastOrderDate.Time
		}

		suppliers = append(suppliers, supplier)
	}

	return suppliers, rows.Err()
}

func (r *SupplierRepository) Update(supplier *models.Supplier) error {
	query := `
		UPDATE suppliers SET
			name = $2, email = $3, phone = $4, address = $5, contact_person = $6,
			contact_phone = $7, tax_number = $8, bank_account = $9, payment_terms = $10,
			delivery_terms = $11, website = $12, rating = $13, status = $14, notes = $15,
			total_orders = $16, is_verified = $17, updated_at = $18
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING updated_at`

	err := r.db.QueryRow(
		query,
		supplier.ID, supplier.Name, supplier.Email, supplier.Phone, supplier.Address,
		supplier.ContactPerson, supplier.ContactPhone, supplier.TaxNumber, supplier.BankAccount,
		supplier.PaymentTerms, supplier.DeliveryTerms, supplier.Website, supplier.Rating,
		supplier.Status, supplier.Notes, supplier.TotalOrders, supplier.IsVerified, time.Now(),
	).Scan(&supplier.UpdatedAt)

	return err
}

func (r *SupplierRepository) Delete(id uint) error {
	query := "UPDATE suppliers SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL"
	result, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("supplier not found or already deleted")
	}

	return nil
}

// Additional helper methods

func (r *SupplierRepository) FindByCode(code string) (*models.Supplier, error) {
	supplier := &models.Supplier{}

	query := `
		SELECT id, name, code, email, phone, address, contact_person, contact_phone,
		       tax_number, bank_account, payment_terms, delivery_terms, website,
		       rating, status, notes, last_order_date, total_orders, is_verified,
		       created_at, updated_at
		FROM suppliers
		WHERE code = $1 AND deleted_at IS NULL`

	var lastOrderDate sql.NullTime

	err := r.db.QueryRow(query, code).Scan(
		&supplier.ID, &supplier.Name, &supplier.Code, &supplier.Email, &supplier.Phone,
		&supplier.Address, &supplier.ContactPerson, &supplier.ContactPhone, &supplier.TaxNumber,
		&supplier.BankAccount, &supplier.PaymentTerms, &supplier.DeliveryTerms, &supplier.Website,
		&supplier.Rating, &supplier.Status, &supplier.Notes, &lastOrderDate, &supplier.TotalOrders,
		&supplier.IsVerified, &supplier.CreatedAt, &supplier.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if lastOrderDate.Valid {
		supplier.LastOrderDate = &lastOrderDate.Time
	}

	return supplier, nil
}

func (r *SupplierRepository) FindByStatus(status string) ([]models.Supplier, error) {
	query := `
		SELECT id, name, code, email, phone, address, contact_person, contact_phone,
		       tax_number, bank_account, payment_terms, delivery_terms, website,
		       rating, status, notes, last_order_date, total_orders, is_verified,
		       created_at, updated_at
		FROM suppliers
		WHERE status = $1 AND deleted_at IS NULL
		ORDER BY name`

	rows, err := r.db.Query(query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suppliers []models.Supplier
	for rows.Next() {
		supplier := models.Supplier{}
		var lastOrderDate sql.NullTime

		err := rows.Scan(
			&supplier.ID, &supplier.Name, &supplier.Code, &supplier.Email, &supplier.Phone,
			&supplier.Address, &supplier.ContactPerson, &supplier.ContactPhone, &supplier.TaxNumber,
			&supplier.BankAccount, &supplier.PaymentTerms, &supplier.DeliveryTerms, &supplier.Website,
			&supplier.Rating, &supplier.Status, &supplier.Notes, &lastOrderDate, &supplier.TotalOrders,
			&supplier.IsVerified, &supplier.CreatedAt, &supplier.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if lastOrderDate.Valid {
			supplier.LastOrderDate = &lastOrderDate.Time
		}

		suppliers = append(suppliers, supplier)
	}

	return suppliers, rows.Err()
}

func (r *SupplierRepository) UpdateLastOrderDate(id uint) error {
	query := `
		UPDATE suppliers SET
			last_order_date = $1, total_orders = total_orders + 1, updated_at = $2
		WHERE id = $3 AND deleted_at IS NULL`

	now := time.Now()
	_, err := r.db.Exec(query, now, now, id)
	return err
}

func (r *SupplierRepository) UpdateRating(id uint, rating float32) error {
	query := `
		UPDATE suppliers SET
			rating = $1, updated_at = $2
		WHERE id = $3 AND deleted_at IS NULL`

	_, err := r.db.Exec(query, rating, time.Now(), id)
	return err
}
