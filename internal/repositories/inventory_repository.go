package repositories

import (
	"database/sql"
	"errors"
	"time"

	"github.com/aliwert/go-hospital-management/internal/models"
	_ "github.com/lib/pq"
)

type InventoryRepository struct {
	db *sql.DB
}

func NewInventoryRepository(db *sql.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) Create(inventory *models.Inventory) error {
	query := `
		INSERT INTO inventories (item_name, item_code, description, category, quantity, unit_price,
		                        supplier_id, reorder_level, expiry_date, batch_number,
		                        location, is_active, status, minimum_quantity, maximum_quantity,
		                        last_order_date, last_count_date, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
		RETURNING id, created_at, updated_at`

	now := time.Now()
	err := r.db.QueryRow(
		query,
		inventory.ItemName,
		inventory.ItemCode,
		inventory.Description,
		inventory.Category,
		inventory.Quantity,
		inventory.UnitPrice,
		inventory.SupplierID,
		inventory.ReorderLevel,
		inventory.ExpiryDate,
		inventory.BatchNumber,
		inventory.Location,
		inventory.IsActive,
		inventory.Status,
		inventory.MinimumQuantity,
		inventory.MaximumQuantity,
		inventory.LastOrderDate,
		inventory.LastCountDate,
		inventory.Notes,
		now,
		now,
	).Scan(&inventory.ID, &inventory.CreatedAt, &inventory.UpdatedAt)

	return err
}

func (r *InventoryRepository) FindById(id uint) (*models.Inventory, error) {
	query := `
		SELECT i.id, i.item_name, i.item_code, i.description, i.category, i.quantity, i.unit_price,
		       i.supplier_id, i.reorder_level, i.expiry_date, i.batch_number,
		       i.location, i.is_active, i.status, i.minimum_quantity, i.maximum_quantity,
		       i.last_order_date, i.last_count_date, i.notes, i.created_at, i.updated_at,
		       s.name as supplier_name, s.code as supplier_code, s.email as supplier_email
		FROM inventories i
		LEFT JOIN suppliers s ON i.supplier_id = s.id
		WHERE i.id = $1 AND i.deleted_at IS NULL`

	var inventory models.Inventory
	var supplierName, supplierCode, supplierEmail sql.NullString
	var supplierID sql.NullInt64

	err := r.db.QueryRow(query, id).Scan(
		&inventory.ID,
		&inventory.ItemName,
		&inventory.ItemCode,
		&inventory.Description,
		&inventory.Category,
		&inventory.Quantity,
		&inventory.UnitPrice,
		&supplierID,
		&inventory.ReorderLevel,
		&inventory.ExpiryDate,
		&inventory.BatchNumber,
		&inventory.Location,
		&inventory.IsActive,
		&inventory.Status,
		&inventory.MinimumQuantity,
		&inventory.MaximumQuantity,
		&inventory.LastOrderDate,
		&inventory.LastCountDate,
		&inventory.Notes,
		&inventory.CreatedAt,
		&inventory.UpdatedAt,
		&supplierName,
		&supplierCode,
		&supplierEmail,
	)

	if err != nil {
		return nil, err
	}

	// Set nullable fields
	if supplierID.Valid {
		inventory.SupplierID = uint(supplierID.Int64)
		if supplierName.Valid {
			inventory.Supplier.Name = supplierName.String
		}
		if supplierCode.Valid {
			inventory.Supplier.Code = supplierCode.String
		}
		if supplierEmail.Valid {
			inventory.Supplier.Email = supplierEmail.String
		}
	}

	return &inventory, nil
}

func (r *InventoryRepository) FindAll() ([]models.Inventory, error) {
	query := `
		SELECT i.id, i.item_name, i.item_code, i.description, i.category, i.quantity, i.unit_price,
		       i.supplier_id, i.reorder_level, i.expiry_date, i.batch_number,
		       i.location, i.is_active, i.status, i.minimum_quantity, i.maximum_quantity,
		       i.last_order_date, i.last_count_date, i.notes, i.created_at, i.updated_at,
		       s.name as supplier_name
		FROM inventories i
		LEFT JOIN suppliers s ON i.supplier_id = s.id
		WHERE i.deleted_at IS NULL
		ORDER BY i.item_name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inventories []models.Inventory
	for rows.Next() {
		var inventory models.Inventory
		var supplierName sql.NullString
		var supplierID sql.NullInt64

		err := rows.Scan(
			&inventory.ID,
			&inventory.ItemName,
			&inventory.ItemCode,
			&inventory.Description,
			&inventory.Category,
			&inventory.Quantity,
			&inventory.UnitPrice,
			&supplierID,
			&inventory.ReorderLevel,
			&inventory.ExpiryDate,
			&inventory.BatchNumber,
			&inventory.Location,
			&inventory.IsActive,
			&inventory.Status,
			&inventory.MinimumQuantity,
			&inventory.MaximumQuantity,
			&inventory.LastOrderDate,
			&inventory.LastCountDate,
			&inventory.Notes,
			&inventory.CreatedAt,
			&inventory.UpdatedAt,
			&supplierName,
		)
		if err != nil {
			return nil, err
		}

		// Set supplier info
		if supplierID.Valid {
			inventory.SupplierID = uint(supplierID.Int64)
			if supplierName.Valid {
				inventory.Supplier.Name = supplierName.String
			}
		}

		inventories = append(inventories, inventory)
	}

	return inventories, rows.Err()
}

func (r *InventoryRepository) Update(inventory *models.Inventory) error {
	query := `
		UPDATE inventories SET
			item_name = $2, item_code = $3, description = $4, category = $5, quantity = $6,
			unit_price = $7, supplier_id = $8, reorder_level = $9, expiry_date = $10,
			batch_number = $11, location = $12, is_active = $13, status = $14,
			minimum_quantity = $15, maximum_quantity = $16, last_order_date = $17,
			last_count_date = $18, notes = $19, updated_at = $20
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING updated_at`

	err := r.db.QueryRow(
		query,
		inventory.ID,
		inventory.ItemName,
		inventory.ItemCode,
		inventory.Description,
		inventory.Category,
		inventory.Quantity,
		inventory.UnitPrice,
		inventory.SupplierID,
		inventory.ReorderLevel,
		inventory.ExpiryDate,
		inventory.BatchNumber,
		inventory.Location,
		inventory.IsActive,
		inventory.Status,
		inventory.MinimumQuantity,
		inventory.MaximumQuantity,
		inventory.LastOrderDate,
		inventory.LastCountDate,
		inventory.Notes,
		time.Now(),
	).Scan(&inventory.UpdatedAt)

	return err
}

func (r *InventoryRepository) Delete(id uint) error {
	query := "UPDATE inventories SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL"
	result, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("inventory item not found or already deleted")
	}

	return nil
}

// Additional helper methods

func (r *InventoryRepository) FindLowStock() ([]models.Inventory, error) {
	query := `
		SELECT i.id, i.item_name, i.item_code, i.category, i.quantity, i.reorder_level,
		       i.minimum_quantity, i.status, s.name as supplier_name
		FROM inventories i
		LEFT JOIN suppliers s ON i.supplier_id = s.id
		WHERE i.quantity <= i.reorder_level
		  AND i.is_active = true
		  AND i.deleted_at IS NULL
		ORDER BY (i.quantity::float / NULLIF(i.reorder_level, 0)) ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inventories []models.Inventory
	for rows.Next() {
		var inventory models.Inventory
		var supplierName sql.NullString

		err := rows.Scan(
			&inventory.ID,
			&inventory.ItemName,
			&inventory.ItemCode,
			&inventory.Category,
			&inventory.Quantity,
			&inventory.ReorderLevel,
			&inventory.MinimumQuantity,
			&inventory.Status,
			&supplierName,
		)
		if err != nil {
			return nil, err
		}

		if supplierName.Valid {
			inventory.Supplier.Name = supplierName.String
		}

		inventories = append(inventories, inventory)
	}

	return inventories, rows.Err()
}

func (r *InventoryRepository) UpdateQuantity(id uint, quantity int) error {
	query := `
		UPDATE inventories SET
			quantity = $2, last_count_date = $3, updated_at = $4,
			status = CASE
				WHEN $2 = 0 THEN 'out_of_stock'
				WHEN $2 <= reorder_level THEN 'low_stock'
				ELSE 'in_stock'
			END
		WHERE id = $1 AND deleted_at IS NULL`

	now := time.Now()
	_, err := r.db.Exec(query, id, quantity, now, now)
	return err
}

func (r *InventoryRepository) FindBySupplier(supplierID uint) ([]models.Inventory, error) {
	query := `
		SELECT i.id, i.item_name, i.item_code, i.category, i.quantity, i.unit_price,
		       i.reorder_level, i.status, i.expiry_date
		FROM inventories i
		WHERE i.supplier_id = $1 AND i.deleted_at IS NULL
		ORDER BY i.item_name`

	rows, err := r.db.Query(query, supplierID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inventories []models.Inventory
	for rows.Next() {
		var inventory models.Inventory

		err := rows.Scan(
			&inventory.ID,
			&inventory.ItemName,
			&inventory.ItemCode,
			&inventory.Category,
			&inventory.Quantity,
			&inventory.UnitPrice,
			&inventory.ReorderLevel,
			&inventory.Status,
			&inventory.ExpiryDate,
		)
		if err != nil {
			return nil, err
		}

		inventories = append(inventories, inventory)
	}

	return inventories, rows.Err()
}

func (r *InventoryRepository) FindExpiringSoon(days int) ([]models.Inventory, error) {
	query := `
		SELECT i.id, i.item_name, i.item_code, i.batch_number, i.expiry_date,
		       i.quantity, i.location, s.name as supplier_name
		FROM inventories i
		LEFT JOIN suppliers s ON i.supplier_id = s.id
		WHERE i.expiry_date <= CURRENT_DATE + INTERVAL '%d days'
		  AND i.expiry_date >= CURRENT_DATE
		  AND i.quantity > 0
		  AND i.deleted_at IS NULL
		ORDER BY i.expiry_date ASC`

	rows, err := r.db.Query(query, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inventories []models.Inventory
	for rows.Next() {
		var inventory models.Inventory
		var supplierName sql.NullString

		err := rows.Scan(
			&inventory.ID,
			&inventory.ItemName,
			&inventory.ItemCode,
			&inventory.BatchNumber,
			&inventory.ExpiryDate,
			&inventory.Quantity,
			&inventory.Location,
			&supplierName,
		)
		if err != nil {
			return nil, err
		}

		if supplierName.Valid {
			inventory.Supplier.Name = supplierName.String
		}

		inventories = append(inventories, inventory)
	}

	return inventories, rows.Err()
}
