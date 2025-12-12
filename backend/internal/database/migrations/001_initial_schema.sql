-- Migration: 001_initial_schema.sql
-- Description: Create initial database schema for SPECTRA
-- Note: This version works WITHOUT PostGIS. Spatial queries use lat/lng columns.
-- Try to enable PostGIS if available (ignore error if not)
DO $$ BEGIN CREATE EXTENSION IF NOT EXISTS postgis;
EXCEPTION
WHEN OTHERS THEN RAISE NOTICE 'PostGIS extension not available, using fallback spatial queries';
END $$;
-- =====================================================
-- NODES TABLE (Map Points: OLT, ODC, ODP, CLOSURE, POLE, CUSTOMER)
-- =====================================================
CREATE TABLE IF NOT EXISTS nodes (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (
        type IN (
            'OLT',
            'ODC',
            'ODP',
            'CLOSURE',
            'POLE',
            'CUSTOMER'
        )
    ),
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    address TEXT,
    capacity_ports INT DEFAULT 8,
    used_ports INT DEFAULT 0,
    model VARCHAR(50),
    status VARCHAR(20) DEFAULT 'ACTIVE' CHECK (
        status IN ('ACTIVE', 'MAINTENANCE', 'PLAN', 'INACTIVE')
    ),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- Trigger to auto-update timestamp
CREATE OR REPLACE FUNCTION update_timestamp() RETURNS TRIGGER AS $$ BEGIN NEW.updated_at := CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER trigger_update_node_timestamp BEFORE
UPDATE ON nodes FOR EACH ROW EXECUTE FUNCTION update_timestamp();
-- =====================================================
-- CABLES TABLE
-- =====================================================
CREATE TABLE IF NOT EXISTS cables (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100),
    type VARCHAR(20) NOT NULL CHECK (type IN ('ADSS', 'DUCT', 'DROP')),
    core_count INT NOT NULL CHECK (core_count > 0),
    length_meter FLOAT,
    origin_node_id BIGINT REFERENCES nodes(id) ON DELETE
    SET NULL,
        dest_node_id BIGINT REFERENCES nodes(id) ON DELETE
    SET NULL,
        path_coordinates JSONB,
        -- Store as [[lng, lat], [lng, lat], ...]
        color_hex VARCHAR(7) DEFAULT '#000000',
        status VARCHAR(20) DEFAULT 'ACTIVE' CHECK (
            status IN ('ACTIVE', 'MAINTENANCE', 'PLAN', 'INACTIVE')
        ),
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE TRIGGER trigger_update_cables_timestamp BEFORE
UPDATE ON cables FOR EACH ROW EXECUTE FUNCTION update_timestamp();
-- =====================================================
-- CABLE_CORES TABLE (Internal Cable Details)
-- =====================================================
CREATE TABLE IF NOT EXISTS cable_cores (
    id BIGSERIAL PRIMARY KEY,
    cable_id BIGINT NOT NULL REFERENCES cables(id) ON DELETE CASCADE,
    core_index INT NOT NULL CHECK (core_index > 0),
    tube_color VARCHAR(30),
    core_color VARCHAR(30),
    status VARCHAR(20) DEFAULT 'VACANT' CHECK (
        status IN ('VACANT', 'USED', 'RESERVED', 'DAMAGED')
    ),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(cable_id, core_index)
);
CREATE TRIGGER trigger_update_cable_cores_timestamp BEFORE
UPDATE ON cable_cores FOR EACH ROW EXECUTE FUNCTION update_timestamp();
-- =====================================================
-- CONNECTIONS TABLE (Splicing Links)
-- =====================================================
CREATE TABLE IF NOT EXISTS connections (
    id BIGSERIAL PRIMARY KEY,
    location_node_id BIGINT REFERENCES nodes(id) ON DELETE
    SET NULL,
        input_type VARCHAR(10) NOT NULL CHECK (input_type IN ('CORE', 'PORT')),
        input_id BIGINT NOT NULL,
        output_type VARCHAR(10) NOT NULL CHECK (output_type IN ('CORE', 'PORT')),
        output_id BIGINT NOT NULL,
        loss_db FLOAT,
        notes TEXT,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE TRIGGER trigger_update_connections_timestamp BEFORE
UPDATE ON connections FOR EACH ROW EXECUTE FUNCTION update_timestamp();
-- =====================================================
-- CUSTOMERS TABLE
-- =====================================================
CREATE TABLE IF NOT EXISTS customers (
    id BIGSERIAL PRIMARY KEY,
    node_id BIGINT REFERENCES nodes(id) ON DELETE
    SET NULL,
        name VARCHAR(100) NOT NULL,
        ont_sn VARCHAR(50) UNIQUE,
        phone VARCHAR(20),
        email VARCHAR(100),
        current_status VARCHAR(20) DEFAULT 'OFFLINE' CHECK (
            current_status IN ('ONLINE', 'OFFLINE', 'LOS', 'POWER_OFF')
        ),
        last_rx_power FLOAT,
        subscription_type VARCHAR(50),
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE TRIGGER trigger_update_customers_timestamp BEFORE
UPDATE ON customers FOR EACH ROW EXECUTE FUNCTION update_timestamp();
-- =====================================================
-- INDEXES
-- =====================================================
CREATE INDEX IF NOT EXISTS idx_nodes_type ON nodes(type);
CREATE INDEX IF NOT EXISTS idx_nodes_status ON nodes(status);
CREATE INDEX IF NOT EXISTS idx_nodes_latlon ON nodes(latitude, longitude);
CREATE INDEX IF NOT EXISTS idx_cables_origin ON cables(origin_node_id);
CREATE INDEX IF NOT EXISTS idx_cables_dest ON cables(dest_node_id);
CREATE INDEX IF NOT EXISTS idx_cable_cores_cable ON cable_cores(cable_id);
CREATE INDEX IF NOT EXISTS idx_connections_location ON connections(location_node_id);
CREATE INDEX IF NOT EXISTS idx_customers_node ON customers(node_id);
CREATE INDEX IF NOT EXISTS idx_customers_status ON customers(current_status);