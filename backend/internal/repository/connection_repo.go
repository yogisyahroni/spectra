package repository

import (
	"context"
	"fmt"

	"spectra-backend/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ConnectionRepository handles database operations for connections (splicing)
type ConnectionRepository struct {
	pool *pgxpool.Pool
}

// NewConnectionRepository creates a new ConnectionRepository
func NewConnectionRepository(pool *pgxpool.Pool) *ConnectionRepository {
	return &ConnectionRepository{pool: pool}
}

// Create inserts a new connection into the database
func (r *ConnectionRepository) Create(ctx context.Context, req *models.CreateConnectionRequest) (*models.Connection, error) {
	query := `
		INSERT INTO connections (location_node_id, input_type, input_id, output_type, output_id, loss_db, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, location_node_id, input_type, input_id, output_type, output_id, loss_db, notes, created_at, updated_at
	`

	conn := &models.Connection{}
	err := r.pool.QueryRow(ctx, query,
		req.LocationNodeID,
		req.InputType,
		req.InputID,
		req.OutputType,
		req.OutputID,
		req.LossDB,
		req.Notes,
	).Scan(
		&conn.ID,
		&conn.LocationNodeID,
		&conn.InputType,
		&conn.InputID,
		&conn.OutputType,
		&conn.OutputID,
		&conn.LossDB,
		&conn.Notes,
		&conn.CreatedAt,
		&conn.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}

	// Update cable core status if input/output are cores
	if req.InputType == models.ConnectionTypeCore {
		r.updateCoreStatus(ctx, req.InputID, models.CoreStatusUsed)
	}
	if req.OutputType == models.ConnectionTypeCore {
		r.updateCoreStatus(ctx, req.OutputID, models.CoreStatusUsed)
	}

	return conn, nil
}

// updateCoreStatus updates the status of a cable core
func (r *ConnectionRepository) updateCoreStatus(ctx context.Context, coreID int64, status models.CoreStatus) error {
	query := "UPDATE cable_cores SET status = $1 WHERE id = $2"
	_, err := r.pool.Exec(ctx, query, status, coreID)
	return err
}

// GetByID retrieves a connection by its ID
func (r *ConnectionRepository) GetByID(ctx context.Context, id int64) (*models.Connection, error) {
	query := `
		SELECT id, location_node_id, input_type, input_id, output_type, output_id, loss_db, notes, created_at, updated_at
		FROM connections
		WHERE id = $1
	`

	conn := &models.Connection{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&conn.ID,
		&conn.LocationNodeID,
		&conn.InputType,
		&conn.InputID,
		&conn.OutputType,
		&conn.OutputID,
		&conn.LossDB,
		&conn.Notes,
		&conn.CreatedAt,
		&conn.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}

	return conn, nil
}

// List retrieves connections with optional filters
func (r *ConnectionRepository) List(ctx context.Context, filter *models.ConnectionFilter) ([]models.Connection, int64, error) {
	baseQuery := "FROM connections WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if filter.LocationNodeID != nil {
		baseQuery += fmt.Sprintf(" AND location_node_id = $%d", argIndex)
		args = append(args, *filter.LocationNodeID)
		argIndex++
	}

	if filter.InputType != nil {
		baseQuery += fmt.Sprintf(" AND input_type = $%d", argIndex)
		args = append(args, *filter.InputType)
		argIndex++
	}

	if filter.InputID != nil {
		baseQuery += fmt.Sprintf(" AND input_id = $%d", argIndex)
		args = append(args, *filter.InputID)
		argIndex++
	}

	if filter.OutputType != nil {
		baseQuery += fmt.Sprintf(" AND output_type = $%d", argIndex)
		args = append(args, *filter.OutputType)
		argIndex++
	}

	if filter.OutputID != nil {
		baseQuery += fmt.Sprintf(" AND output_id = $%d", argIndex)
		args = append(args, *filter.OutputID)
		argIndex++
	}

	// Count total
	var total int64
	countQuery := "SELECT COUNT(*) " + baseQuery
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count connections: %w", err)
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
		SELECT id, location_node_id, input_type, input_id, output_type, output_id, loss_db, notes, created_at, updated_at
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, baseQuery, argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list connections: %w", err)
	}
	defer rows.Close()

	var connections []models.Connection
	for rows.Next() {
		var conn models.Connection
		err := rows.Scan(
			&conn.ID,
			&conn.LocationNodeID,
			&conn.InputType,
			&conn.InputID,
			&conn.OutputType,
			&conn.OutputID,
			&conn.LossDB,
			&conn.Notes,
			&conn.CreatedAt,
			&conn.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan connection: %w", err)
		}
		connections = append(connections, conn)
	}

	return connections, total, nil
}

// Delete removes a connection by its ID
func (r *ConnectionRepository) Delete(ctx context.Context, id int64) error {
	// Get connection first to free up cores
	conn, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if conn == nil {
		return fmt.Errorf("connection not found")
	}

	// Delete connection
	query := "DELETE FROM connections WHERE id = $1"
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete connection: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("connection not found")
	}

	// Free up cores
	if conn.InputType == models.ConnectionTypeCore {
		r.updateCoreStatus(ctx, conn.InputID, models.CoreStatusVacant)
	}
	if conn.OutputType == models.ConnectionTypeCore {
		r.updateCoreStatus(ctx, conn.OutputID, models.CoreStatusVacant)
	}

	return nil
}

// GetByLocation retrieves all connections at a specific node location
func (r *ConnectionRepository) GetByLocation(ctx context.Context, nodeID int64) ([]models.Connection, error) {
	query := `
		SELECT id, location_node_id, input_type, input_id, output_type, output_id, loss_db, notes, created_at, updated_at
		FROM connections
		WHERE location_node_id = $1
		ORDER BY id ASC
	`

	rows, err := r.pool.Query(ctx, query, nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get connections by location: %w", err)
	}
	defer rows.Close()

	var connections []models.Connection
	for rows.Next() {
		var conn models.Connection
		err := rows.Scan(
			&conn.ID,
			&conn.LocationNodeID,
			&conn.InputType,
			&conn.InputID,
			&conn.OutputType,
			&conn.OutputID,
			&conn.LossDB,
			&conn.Notes,
			&conn.CreatedAt,
			&conn.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan connection: %w", err)
		}
		connections = append(connections, conn)
	}

	return connections, nil
}

// GetSpliceMatrix retrieves the splice matrix for a specific closure node
func (r *ConnectionRepository) GetSpliceMatrix(ctx context.Context, nodeID int64) (*models.SpliceMatrix, error) {
	// Get connections at this location
	connections, err := r.GetByLocation(ctx, nodeID)
	if err != nil {
		return nil, err
	}

	matrix := &models.SpliceMatrix{
		Connections: make([]models.MatrixConnection, 0),
	}

	for _, conn := range connections {
		if conn.InputType == models.ConnectionTypeCore && conn.OutputType == models.ConnectionTypeCore {
			matrix.Connections = append(matrix.Connections, models.MatrixConnection{
				InputCoreID:  conn.InputID,
				OutputCoreID: conn.OutputID,
				LossDB:       conn.LossDB,
			})
		}
	}

	return matrix, nil
}

// CalculateTotalLoss calculates the total loss from a customer to OLT
func (r *ConnectionRepository) CalculateTotalLoss(ctx context.Context, customerNodeID int64) (float64, error) {
	// This is a simplified implementation
	// A full trace would need to follow the connection chain
	query := `
		WITH RECURSIVE trace AS (
			-- Start from connections at customer node
			SELECT c.id, c.loss_db, c.output_id, c.output_type, 1 as depth
			FROM connections c
			JOIN customers cust ON cust.node_id = c.location_node_id
			WHERE cust.node_id = $1

			UNION ALL

			-- Follow the chain
			SELECT c2.id, c2.loss_db, c2.output_id, c2.output_type, t.depth + 1
			FROM connections c2
			JOIN trace t ON c2.input_id = t.output_id AND c2.input_type = t.output_type
			WHERE t.depth < 20  -- Prevent infinite loops
		)
		SELECT COALESCE(SUM(loss_db), 0) FROM trace
	`

	var totalLoss float64
	err := r.pool.QueryRow(ctx, query, customerNodeID).Scan(&totalLoss)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate total loss: %w", err)
	}

	return totalLoss, nil
}
