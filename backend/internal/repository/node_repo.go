package repository

import (
	"context"
	"fmt"

	"spectra-backend/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NodeRepository handles database operations for nodes
type NodeRepository struct {
	pool *pgxpool.Pool
}

// NewNodeRepository creates a new NodeRepository
func NewNodeRepository(pool *pgxpool.Pool) *NodeRepository {
	return &NodeRepository{pool: pool}
}

// Create inserts a new node into the database
func (r *NodeRepository) Create(ctx context.Context, req *models.CreateNodeRequest) (*models.Node, error) {
	query := `
		INSERT INTO nodes (name, type, latitude, longitude, address, capacity_ports, model, status)
		VALUES ($1, $2, $3, $4, $5, COALESCE($6, 8), $7, COALESCE($8, 'ACTIVE'))
		RETURNING id, name, type, latitude, longitude, address, capacity_ports, used_ports, model, status, created_at, updated_at
	`

	capacityPorts := 8
	if req.CapacityPorts != nil {
		capacityPorts = *req.CapacityPorts
	}

	status := models.NodeStatusActive
	if req.Status != "" {
		status = req.Status
	}

	node := &models.Node{}
	err := r.pool.QueryRow(ctx, query,
		req.Name,
		req.Type,
		req.Latitude,
		req.Longitude,
		req.Address,
		capacityPorts,
		req.Model,
		status,
	).Scan(
		&node.ID,
		&node.Name,
		&node.Type,
		&node.Latitude,
		&node.Longitude,
		&node.Address,
		&node.CapacityPorts,
		&node.UsedPorts,
		&node.Model,
		&node.Status,
		&node.CreatedAt,
		&node.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create node: %w", err)
	}

	return node, nil
}

// GetByID retrieves a node by its ID
func (r *NodeRepository) GetByID(ctx context.Context, id int64) (*models.Node, error) {
	query := `
		SELECT id, name, type, latitude, longitude, address, capacity_ports, used_ports, model, status, created_at, updated_at
		FROM nodes
		WHERE id = $1
	`

	node := &models.Node{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&node.ID,
		&node.Name,
		&node.Type,
		&node.Latitude,
		&node.Longitude,
		&node.Address,
		&node.CapacityPorts,
		&node.UsedPorts,
		&node.Model,
		&node.Status,
		&node.CreatedAt,
		&node.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get node: %w", err)
	}

	return node, nil
}

// List retrieves nodes with optional filters
func (r *NodeRepository) List(ctx context.Context, filter *models.NodeFilter) ([]models.Node, int64, error) {
	baseQuery := "FROM nodes WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if filter.Type != nil {
		baseQuery += fmt.Sprintf(" AND type = $%d", argIndex)
		args = append(args, *filter.Type)
		argIndex++
	}

	if filter.Status != nil {
		baseQuery += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, *filter.Status)
		argIndex++
	}

	// Count total
	var total int64
	countQuery := "SELECT COUNT(*) " + baseQuery
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count nodes: %w", err)
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
		SELECT id, name, type, latitude, longitude, address, capacity_ports, used_ports, model, status, created_at, updated_at
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, baseQuery, argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list nodes: %w", err)
	}
	defer rows.Close()

	var nodes []models.Node
	for rows.Next() {
		var node models.Node
		err := rows.Scan(
			&node.ID,
			&node.Name,
			&node.Type,
			&node.Latitude,
			&node.Longitude,
			&node.Address,
			&node.CapacityPorts,
			&node.UsedPorts,
			&node.Model,
			&node.Status,
			&node.CreatedAt,
			&node.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan node: %w", err)
		}
		nodes = append(nodes, node)
	}

	return nodes, total, nil
}

// Update updates an existing node
func (r *NodeRepository) Update(ctx context.Context, id int64, req *models.UpdateNodeRequest) (*models.Node, error) {
	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *req.Name)
		argIndex++
	}
	if req.Type != nil {
		setParts = append(setParts, fmt.Sprintf("type = $%d", argIndex))
		args = append(args, *req.Type)
		argIndex++
	}
	if req.Latitude != nil {
		setParts = append(setParts, fmt.Sprintf("latitude = $%d", argIndex))
		args = append(args, *req.Latitude)
		argIndex++
	}
	if req.Longitude != nil {
		setParts = append(setParts, fmt.Sprintf("longitude = $%d", argIndex))
		args = append(args, *req.Longitude)
		argIndex++
	}
	if req.Address != nil {
		setParts = append(setParts, fmt.Sprintf("address = $%d", argIndex))
		args = append(args, *req.Address)
		argIndex++
	}
	if req.CapacityPorts != nil {
		setParts = append(setParts, fmt.Sprintf("capacity_ports = $%d", argIndex))
		args = append(args, *req.CapacityPorts)
		argIndex++
	}
	if req.UsedPorts != nil {
		setParts = append(setParts, fmt.Sprintf("used_ports = $%d", argIndex))
		args = append(args, *req.UsedPorts)
		argIndex++
	}
	if req.Model != nil {
		setParts = append(setParts, fmt.Sprintf("model = $%d", argIndex))
		args = append(args, *req.Model)
		argIndex++
	}
	if req.Status != nil {
		setParts = append(setParts, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *req.Status)
		argIndex++
	}

	if len(setParts) == 0 {
		return r.GetByID(ctx, id)
	}

	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE nodes
		SET %s
		WHERE id = $%d
		RETURNING id, name, type, latitude, longitude, address, capacity_ports, used_ports, model, status, created_at, updated_at
	`, joinStrings(setParts, ", "), argIndex)

	node := &models.Node{}
	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&node.ID,
		&node.Name,
		&node.Type,
		&node.Latitude,
		&node.Longitude,
		&node.Address,
		&node.CapacityPorts,
		&node.UsedPorts,
		&node.Model,
		&node.Status,
		&node.CreatedAt,
		&node.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update node: %w", err)
	}

	return node, nil
}

// Delete removes a node by its ID
func (r *NodeRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM nodes WHERE id = $1"
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete node: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("node not found")
	}

	return nil
}

// GetNearby retrieves nodes within a specified radius using Haversine formula
// This works without PostGIS by calculating distance in the query
func (r *NodeRepository) GetNearby(ctx context.Context, query *models.NearbyQuery) ([]models.Node, error) {
	// Haversine formula in SQL for distance calculation (returns distance in km)
	// 6371 is Earth's radius in kilometers
	sqlQuery := `
		SELECT id, name, type, latitude, longitude, address, capacity_ports, used_ports, model, status, created_at, updated_at,
			   (6371 * acos(cos(radians($1)) * cos(radians(latitude)) * cos(radians(longitude) - radians($2)) + sin(radians($1)) * sin(radians(latitude)))) AS distance_km
		FROM nodes
		WHERE (6371 * acos(cos(radians($1)) * cos(radians(latitude)) * cos(radians(longitude) - radians($2)) + sin(radians($1)) * sin(radians(latitude)))) <= $3
	`

	args := []interface{}{query.Latitude, query.Longitude, query.RadiusKM}
	argIndex := 4

	if query.Type != nil {
		sqlQuery += fmt.Sprintf(" AND type = $%d", argIndex)
		args = append(args, *query.Type)
	}

	sqlQuery += " ORDER BY distance_km ASC"

	rows, err := r.pool.Query(ctx, sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get nearby nodes: %w", err)
	}
	defer rows.Close()

	var nodes []models.Node
	for rows.Next() {
		var node models.Node
		var distanceKM float64
		err := rows.Scan(
			&node.ID,
			&node.Name,
			&node.Type,
			&node.Latitude,
			&node.Longitude,
			&node.Address,
			&node.CapacityPorts,
			&node.UsedPorts,
			&node.Model,
			&node.Status,
			&node.CreatedAt,
			&node.UpdatedAt,
			&distanceKM,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan node: %w", err)
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

// GetAllAsGeoJSON retrieves all nodes as GeoJSON features
func (r *NodeRepository) GetAllAsGeoJSON(ctx context.Context, filter *models.NodeFilter) ([]models.NodeGeoJSON, error) {
	nodes, _, err := r.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	features := make([]models.NodeGeoJSON, len(nodes))
	for i, node := range nodes {
		features[i] = node.ToGeoJSON()
	}

	return features, nil
}

// Helper function to join strings
func joinStrings(parts []string, sep string) string {
	result := ""
	for i, part := range parts {
		if i > 0 {
			result += sep
		}
		result += part
	}
	return result
}
