package repository

import (
	"context"
	"fmt"
	"strings"

	"spectra-backend/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CableRepository handles database operations for cables
type CableRepository struct {
	pool *pgxpool.Pool
}

// NewCableRepository creates a new CableRepository
func NewCableRepository(pool *pgxpool.Pool) *CableRepository {
	return &CableRepository{pool: pool}
}

// Create inserts a new cable into the database
func (r *CableRepository) Create(ctx context.Context, req *models.CreateCableRequest) (*models.Cable, error) {
	// Build path geometry from coordinates if provided
	var pathGeometrySQL string
	if len(req.PathCoordinates) >= 2 {
		points := make([]string, len(req.PathCoordinates))
		for i, coord := range req.PathCoordinates {
			if len(coord) >= 2 {
				points[i] = fmt.Sprintf("%f %f", coord[0], coord[1])
			}
		}
		pathGeometrySQL = fmt.Sprintf("ST_SetSRID(ST_MakeLine(ARRAY[%s]::geometry[]), 4326)",
			"ST_MakePoint("+strings.Join(points, "), ST_MakePoint(")+")")
	}

	colorHex := "#000000"
	if req.ColorHex != nil {
		colorHex = *req.ColorHex
	}

	status := models.CableStatusActive
	if req.Status != "" {
		status = req.Status
	}

	var query string
	var args []interface{}

	if pathGeometrySQL != "" {
		query = fmt.Sprintf(`
			INSERT INTO cables (name, type, core_count, length_meter, origin_node_id, dest_node_id, path_geometry, color_hex, status)
			VALUES ($1, $2, $3, $4, $5, $6, %s, $7, $8)
			RETURNING id, name, type, core_count, length_meter, origin_node_id, dest_node_id, color_hex, status, created_at, updated_at
		`, pathGeometrySQL)
		args = []interface{}{
			req.Name,
			req.Type,
			req.CoreCount,
			req.LengthMeter,
			req.OriginNodeID,
			req.DestNodeID,
			colorHex,
			status,
		}
	} else {
		query = `
			INSERT INTO cables (name, type, core_count, length_meter, origin_node_id, dest_node_id, color_hex, status)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id, name, type, core_count, length_meter, origin_node_id, dest_node_id, color_hex, status, created_at, updated_at
		`
		args = []interface{}{
			req.Name,
			req.Type,
			req.CoreCount,
			req.LengthMeter,
			req.OriginNodeID,
			req.DestNodeID,
			colorHex,
			status,
		}
	}

	cable := &models.Cable{}
	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&cable.ID,
		&cable.Name,
		&cable.Type,
		&cable.CoreCount,
		&cable.LengthMeter,
		&cable.OriginNodeID,
		&cable.DestNodeID,
		&cable.ColorHex,
		&cable.Status,
		&cable.CreatedAt,
		&cable.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create cable: %w", err)
	}

	// Auto-generate cable cores
	if err := r.generateCores(ctx, cable.ID, req.CoreCount); err != nil {
		// Log warning but don't fail
		fmt.Printf("Warning: failed to generate cores for cable %d: %v\n", cable.ID, err)
	}

	return cable, nil
}

// generateCores creates cable core entries for a new cable
func (r *CableRepository) generateCores(ctx context.Context, cableID int64, coreCount int) error {
	cores := models.GenerateCoresForCable(cableID, coreCount)

	batch := &pgx.Batch{}
	for _, core := range cores {
		batch.Queue(`
			INSERT INTO cable_cores (cable_id, core_index, tube_color, core_color, status)
			VALUES ($1, $2, $3, $4, $5)
		`, core.CableID, core.CoreIndex, core.TubeColor, core.CoreColor, core.Status)
	}

	results := r.pool.SendBatch(ctx, batch)
	defer results.Close()

	for i := 0; i < batch.Len(); i++ {
		if _, err := results.Exec(); err != nil {
			return fmt.Errorf("failed to insert core %d: %w", i+1, err)
		}
	}

	return nil
}

// GetByID retrieves a cable by its ID
func (r *CableRepository) GetByID(ctx context.Context, id int64) (*models.Cable, error) {
	query := `
		SELECT 
			id, name, type, core_count, length_meter, origin_node_id, dest_node_id, 
			color_hex, status, created_at, updated_at,
			ST_AsGeoJSON(path_geometry)::jsonb->'coordinates' as path_coords
		FROM cables
		WHERE id = $1
	`

	cable := &models.Cable{}
	var pathCoordsJSON []byte

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&cable.ID,
		&cable.Name,
		&cable.Type,
		&cable.CoreCount,
		&cable.LengthMeter,
		&cable.OriginNodeID,
		&cable.DestNodeID,
		&cable.ColorHex,
		&cable.Status,
		&cable.CreatedAt,
		&cable.UpdatedAt,
		&pathCoordsJSON,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get cable: %w", err)
	}

	return cable, nil
}

// List retrieves cables with optional filters
func (r *CableRepository) List(ctx context.Context, filter *models.CableFilter) ([]models.Cable, int64, error) {
	baseQuery := "FROM cables WHERE 1=1"
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

	if filter.OriginNodeID != nil {
		baseQuery += fmt.Sprintf(" AND origin_node_id = $%d", argIndex)
		args = append(args, *filter.OriginNodeID)
		argIndex++
	}

	if filter.DestNodeID != nil {
		baseQuery += fmt.Sprintf(" AND dest_node_id = $%d", argIndex)
		args = append(args, *filter.DestNodeID)
		argIndex++
	}

	// Count total
	var total int64
	countQuery := "SELECT COUNT(*) " + baseQuery
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count cables: %w", err)
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
		SELECT id, name, type, core_count, length_meter, origin_node_id, dest_node_id, color_hex, status, created_at, updated_at
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, baseQuery, argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list cables: %w", err)
	}
	defer rows.Close()

	var cables []models.Cable
	for rows.Next() {
		var cable models.Cable
		err := rows.Scan(
			&cable.ID,
			&cable.Name,
			&cable.Type,
			&cable.CoreCount,
			&cable.LengthMeter,
			&cable.OriginNodeID,
			&cable.DestNodeID,
			&cable.ColorHex,
			&cable.Status,
			&cable.CreatedAt,
			&cable.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan cable: %w", err)
		}
		cables = append(cables, cable)
	}

	return cables, total, nil
}

// Update updates an existing cable
func (r *CableRepository) Update(ctx context.Context, id int64, req *models.UpdateCableRequest) (*models.Cable, error) {
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
	if req.CoreCount != nil {
		setParts = append(setParts, fmt.Sprintf("core_count = $%d", argIndex))
		args = append(args, *req.CoreCount)
		argIndex++
	}
	if req.LengthMeter != nil {
		setParts = append(setParts, fmt.Sprintf("length_meter = $%d", argIndex))
		args = append(args, *req.LengthMeter)
		argIndex++
	}
	if req.OriginNodeID != nil {
		setParts = append(setParts, fmt.Sprintf("origin_node_id = $%d", argIndex))
		args = append(args, *req.OriginNodeID)
		argIndex++
	}
	if req.DestNodeID != nil {
		setParts = append(setParts, fmt.Sprintf("dest_node_id = $%d", argIndex))
		args = append(args, *req.DestNodeID)
		argIndex++
	}
	if req.ColorHex != nil {
		setParts = append(setParts, fmt.Sprintf("color_hex = $%d", argIndex))
		args = append(args, *req.ColorHex)
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
		UPDATE cables
		SET %s
		WHERE id = $%d
		RETURNING id, name, type, core_count, length_meter, origin_node_id, dest_node_id, color_hex, status, created_at, updated_at
	`, joinStrings(setParts, ", "), argIndex)

	cable := &models.Cable{}
	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&cable.ID,
		&cable.Name,
		&cable.Type,
		&cable.CoreCount,
		&cable.LengthMeter,
		&cable.OriginNodeID,
		&cable.DestNodeID,
		&cable.ColorHex,
		&cable.Status,
		&cable.CreatedAt,
		&cable.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update cable: %w", err)
	}

	return cable, nil
}

// Delete removes a cable by its ID
func (r *CableRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM cables WHERE id = $1"
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete cable: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("cable not found")
	}

	return nil
}

// GetCores retrieves all cores for a cable
func (r *CableRepository) GetCores(ctx context.Context, cableID int64) ([]models.CableCore, error) {
	query := `
		SELECT id, cable_id, core_index, tube_color, core_color, status, created_at, updated_at
		FROM cable_cores
		WHERE cable_id = $1
		ORDER BY core_index ASC
	`

	rows, err := r.pool.Query(ctx, query, cableID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cable cores: %w", err)
	}
	defer rows.Close()

	var cores []models.CableCore
	for rows.Next() {
		var core models.CableCore
		err := rows.Scan(
			&core.ID,
			&core.CableID,
			&core.CoreIndex,
			&core.TubeColor,
			&core.CoreColor,
			&core.Status,
			&core.CreatedAt,
			&core.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan core: %w", err)
		}
		cores = append(cores, core)
	}

	return cores, nil
}

// UpdateCore updates a specific cable core
func (r *CableRepository) UpdateCore(ctx context.Context, cableID, coreID int64, req *models.UpdateCableCoreRequest) (*models.CableCore, error) {
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.TubeColor != nil {
		setParts = append(setParts, fmt.Sprintf("tube_color = $%d", argIndex))
		args = append(args, *req.TubeColor)
		argIndex++
	}
	if req.CoreColor != nil {
		setParts = append(setParts, fmt.Sprintf("core_color = $%d", argIndex))
		args = append(args, *req.CoreColor)
		argIndex++
	}
	if req.Status != nil {
		setParts = append(setParts, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *req.Status)
		argIndex++
	}

	if len(setParts) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	args = append(args, cableID, coreID)

	query := fmt.Sprintf(`
		UPDATE cable_cores
		SET %s
		WHERE cable_id = $%d AND id = $%d
		RETURNING id, cable_id, core_index, tube_color, core_color, status, created_at, updated_at
	`, joinStrings(setParts, ", "), argIndex, argIndex+1)

	core := &models.CableCore{}
	err := r.pool.QueryRow(ctx, query, args...).Scan(
		&core.ID,
		&core.CableID,
		&core.CoreIndex,
		&core.TubeColor,
		&core.CoreColor,
		&core.Status,
		&core.CreatedAt,
		&core.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update core: %w", err)
	}

	return core, nil
}

// GetAllAsGeoJSON retrieves all cables as GeoJSON features
func (r *CableRepository) GetAllAsGeoJSON(ctx context.Context) ([]models.CableGeoJSON, error) {
	query := `
		SELECT 
			id, name, type, core_count, length_meter, origin_node_id, dest_node_id, 
			color_hex, status, created_at, updated_at,
			COALESCE(
				(SELECT json_agg(json_build_array(ST_X(geom), ST_Y(geom)))
				 FROM ST_DumpPoints(path_geometry) AS dp(path, geom)),
				'[]'::json
			) as path_coords
		FROM cables
		WHERE path_geometry IS NOT NULL
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get cables as geojson: %w", err)
	}
	defer rows.Close()

	var features []models.CableGeoJSON
	for rows.Next() {
		var cable models.Cable
		var pathCoordsJSON []byte

		err := rows.Scan(
			&cable.ID,
			&cable.Name,
			&cable.Type,
			&cable.CoreCount,
			&cable.LengthMeter,
			&cable.OriginNodeID,
			&cable.DestNodeID,
			&cable.ColorHex,
			&cable.Status,
			&cable.CreatedAt,
			&cable.UpdatedAt,
			&pathCoordsJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan cable: %w", err)
		}

		features = append(features, cable.ToGeoJSON())
	}

	return features, nil
}
