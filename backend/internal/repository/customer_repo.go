package repository

import (
	"context"
	"fmt"

	"spectra-backend/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CustomerRepository handles database operations for customers
type CustomerRepository struct {
	pool *pgxpool.Pool
}

// NewCustomerRepository creates a new CustomerRepository
func NewCustomerRepository(pool *pgxpool.Pool) *CustomerRepository {
	return &CustomerRepository{pool: pool}
}

// Create inserts a new customer into the database
func (r *CustomerRepository) Create(ctx context.Context, req *models.CreateCustomerRequest) (*models.Customer, error) {
	status := models.CustomerStatusOffline
	if req.CurrentStatus != nil {
		status = *req.CurrentStatus
	}

	query := `
		INSERT INTO customers (node_id, name, ont_sn, phone, email, current_status, subscription_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, node_id, name, ont_sn, phone, email, current_status, last_rx_power, subscription_type, created_at, updated_at
	`

	customer := &models.Customer{}
	err := r.pool.QueryRow(ctx, query,
		req.NodeID,
		req.Name,
		req.ONTSN,
		req.Phone,
		req.Email,
		status,
		req.SubscriptionType,
	).Scan(
		&customer.ID,
		&customer.NodeID,
		&customer.Name,
		&customer.ONTSN,
		&customer.Phone,
		&customer.Email,
		&customer.CurrentStatus,
		&customer.LastRxPower,
		&customer.SubscriptionType,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	return customer, nil
}

// GetByID retrieves a customer by its ID
func (r *CustomerRepository) GetByID(ctx context.Context, id int64) (*models.Customer, error) {
	query := `
		SELECT id, node_id, name, ont_sn, phone, email, current_status, last_rx_power, subscription_type, created_at, updated_at
		FROM customers
		WHERE id = $1
	`

	customer := &models.Customer{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&customer.ID,
		&customer.NodeID,
		&customer.Name,
		&customer.ONTSN,
		&customer.Phone,
		&customer.Email,
		&customer.CurrentStatus,
		&customer.LastRxPower,
		&customer.SubscriptionType,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	return customer, nil
}

// GetByONTSN retrieves a customer by ONT serial number
func (r *CustomerRepository) GetByONTSN(ctx context.Context, ontSN string) (*models.Customer, error) {
	query := `
		SELECT id, node_id, name, ont_sn, phone, email, current_status, last_rx_power, subscription_type, created_at, updated_at
		FROM customers
		WHERE ont_sn = $1
	`

	customer := &models.Customer{}
	err := r.pool.QueryRow(ctx, query, ontSN).Scan(
		&customer.ID,
		&customer.NodeID,
		&customer.Name,
		&customer.ONTSN,
		&customer.Phone,
		&customer.Email,
		&customer.CurrentStatus,
		&customer.LastRxPower,
		&customer.SubscriptionType,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get customer by ONT SN: %w", err)
	}

	return customer, nil
}

// List retrieves customers with optional filters
func (r *CustomerRepository) List(ctx context.Context, filter *models.CustomerFilter) ([]models.Customer, int64, error) {
	baseQuery := "FROM customers WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if filter.NodeID != nil {
		baseQuery += fmt.Sprintf(" AND node_id = $%d", argIndex)
		args = append(args, *filter.NodeID)
		argIndex++
	}

	if filter.CurrentStatus != nil {
		baseQuery += fmt.Sprintf(" AND current_status = $%d", argIndex)
		args = append(args, *filter.CurrentStatus)
		argIndex++
	}

	if filter.Search != nil && *filter.Search != "" {
		baseQuery += fmt.Sprintf(" AND (name ILIKE $%d OR ont_sn ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+*filter.Search+"%")
		argIndex++
	}

	// Count total
	var total int64
	countQuery := "SELECT COUNT(*) " + baseQuery
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count customers: %w", err)
	}

	// Get data with pagination
	limit := 100
	offset := 0
	if filter.Limit > 0 {
		limit = filter.Limit
	}
	if filter.Offset > 0 {
		offset = filter.Offset
	}

	dataQuery := fmt.Sprintf(`
		SELECT id, node_id, name, ont_sn, phone, email, current_status, last_rx_power, subscription_type, created_at, updated_at
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, baseQuery, argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list customers: %w", err)
	}
	defer rows.Close()

	var customers []models.Customer
	for rows.Next() {
		var customer models.Customer
		err := rows.Scan(
			&customer.ID,
			&customer.NodeID,
			&customer.Name,
			&customer.ONTSN,
			&customer.Phone,
			&customer.Email,
			&customer.CurrentStatus,
			&customer.LastRxPower,
			&customer.SubscriptionType,
			&customer.CreatedAt,
			&customer.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan customer: %w", err)
		}
		customers = append(customers, customer)
	}

	return customers, total, nil
}

// Update updates an existing customer
func (r *CustomerRepository) Update(ctx context.Context, id int64, req *models.UpdateCustomerRequest) (*models.Customer, error) {
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.NodeID != nil {
		setParts = append(setParts, fmt.Sprintf("node_id = $%d", argIndex))
		args = append(args, *req.NodeID)
		argIndex++
	}
	if req.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *req.Name)
		argIndex++
	}
	if req.ONTSN != nil {
		setParts = append(setParts, fmt.Sprintf("ont_sn = $%d", argIndex))
		args = append(args, *req.ONTSN)
		argIndex++
	}
	if req.Phone != nil {
		setParts = append(setParts, fmt.Sprintf("phone = $%d", argIndex))
		args = append(args, *req.Phone)
		argIndex++
	}
	if req.Email != nil {
		setParts = append(setParts, fmt.Sprintf("email = $%d", argIndex))
		args = append(args, *req.Email)
		argIndex++
	}
	if req.CurrentStatus != nil {
		setParts = append(setParts, fmt.Sprintf("current_status = $%d", argIndex))
		args = append(args, *req.CurrentStatus)
		argIndex++
	}
	if req.LastRxPower != nil {
		setParts = append(setParts, fmt.Sprintf("last_rx_power = $%d", argIndex))
		args = append(args, *req.LastRxPower)
		argIndex++
	}
	if req.SubscriptionType != nil {
		setParts = append(setParts, fmt.Sprintf("subscription_type = $%d", argIndex))
		args = append(args, *req.SubscriptionType)
		argIndex++
	}

	if len(setParts) == 0 {
		return r.GetByID(ctx, id)
	}

	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE customers
		SET %s
		WHERE id = $%d
		RETURNING id, node_id, name, ont_sn, phone, email, current_status, last_rx_power, subscription_type, created_at, updated_at
	`, joinStrings(setParts, ", "), argIndex)

	customer := &models.Customer{}
	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&customer.ID,
		&customer.NodeID,
		&customer.Name,
		&customer.ONTSN,
		&customer.Phone,
		&customer.Email,
		&customer.CurrentStatus,
		&customer.LastRxPower,
		&customer.SubscriptionType,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update customer: %w", err)
	}

	return customer, nil
}

// Delete removes a customer by its ID
func (r *CustomerRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM customers WHERE id = $1"
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("customer not found")
	}

	return nil
}

// GetByStatus retrieves customers by their connection status
func (r *CustomerRepository) GetByStatus(ctx context.Context, status models.CustomerStatus) ([]models.Customer, error) {
	query := `
		SELECT id, node_id, name, ont_sn, phone, email, current_status, last_rx_power, subscription_type, created_at, updated_at
		FROM customers
		WHERE current_status = $1
		ORDER BY updated_at DESC
	`

	rows, err := r.pool.Query(ctx, query, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get customers by status: %w", err)
	}
	defer rows.Close()

	var customers []models.Customer
	for rows.Next() {
		var customer models.Customer
		err := rows.Scan(
			&customer.ID,
			&customer.NodeID,
			&customer.Name,
			&customer.ONTSN,
			&customer.Phone,
			&customer.Email,
			&customer.CurrentStatus,
			&customer.LastRxPower,
			&customer.SubscriptionType,
			&customer.CreatedAt,
			&customer.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan customer: %w", err)
		}
		customers = append(customers, customer)
	}

	return customers, nil
}

// GetLOSCustomers retrieves all customers with LOS (Loss of Signal) status
func (r *CustomerRepository) GetLOSCustomers(ctx context.Context) ([]models.Customer, error) {
	return r.GetByStatus(ctx, models.CustomerStatusLOS)
}

// UpdateStatus updates the status and Rx power of a customer
func (r *CustomerRepository) UpdateStatus(ctx context.Context, id int64, status models.CustomerStatus, rxPower *float64) error {
	query := `
		UPDATE customers
		SET current_status = $1, last_rx_power = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`

	_, err := r.pool.Exec(ctx, query, status, rxPower, id)
	if err != nil {
		return fmt.Errorf("failed to update customer status: %w", err)
	}

	return nil
}

// BulkUpdateStatus updates status for multiple customers by ONT SN
func (r *CustomerRepository) BulkUpdateStatus(ctx context.Context, updates map[string]models.CustomerStatus) error {
	if len(updates) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	for ontSN, status := range updates {
		batch.Queue(`
			UPDATE customers
			SET current_status = $1, updated_at = CURRENT_TIMESTAMP
			WHERE ont_sn = $2
		`, status, ontSN)
	}

	results := r.pool.SendBatch(ctx, batch)
	defer results.Close()

	for i := 0; i < batch.Len(); i++ {
		if _, err := results.Exec(); err != nil {
			return fmt.Errorf("failed to update customer status: %w", err)
		}
	}

	return nil
}
