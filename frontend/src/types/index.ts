// Types matching backend models

export type NodeType = 'OLT' | 'ODC' | 'ODP' | 'CLOSURE' | 'POLE' | 'CUSTOMER';
export type NodeStatus = 'ACTIVE' | 'MAINTENANCE' | 'PLAN' | 'INACTIVE';
export type CableType = 'ADSS' | 'DUCT' | 'DROP';
export type CableStatus = 'ACTIVE' | 'MAINTENANCE' | 'PLAN' | 'INACTIVE';
export type CoreStatus = 'VACANT' | 'USED' | 'RESERVED' | 'DAMAGED';
export type ConnectionType = 'CORE' | 'PORT';
export type CustomerStatus = 'ONLINE' | 'OFFLINE' | 'LOS' | 'POWER_OFF';

export interface Node {
    id: number;
    name: string;
    type: NodeType;
    latitude: number;
    longitude: number;
    address?: string;
    capacity_ports: number;
    used_ports: number;
    model?: string;
    status: NodeStatus;
    created_at: string;
    updated_at: string;
}

export interface CreateNodeRequest {
    name: string;
    type: NodeType;
    latitude: number;
    longitude: number;
    address?: string;
    capacity_ports?: number;
    model?: string;
    status?: NodeStatus;
}

export interface UpdateNodeRequest {
    name?: string;
    type?: NodeType;
    latitude?: number;
    longitude?: number;
    address?: string;
    capacity_ports?: number;
    used_ports?: number;
    model?: string;
    status?: NodeStatus;
}

export interface Cable {
    id: number;
    name?: string;
    type: CableType;
    core_count: number;
    length_meter?: number;
    origin_node_id?: number;
    dest_node_id?: number;
    color_hex: string;
    status: CableStatus;
    created_at: string;
    updated_at: string;
    path_coordinates?: number[][];
}

export interface CreateCableRequest {
    name?: string;
    type: CableType;
    core_count: number;
    length_meter?: number;
    origin_node_id?: number;
    dest_node_id?: number;
    path_coordinates?: number[][];
    color_hex?: string;
    status?: CableStatus;
}

export interface CableCore {
    id: number;
    cable_id: number;
    core_index: number;
    tube_color?: string;
    core_color?: string;
    status: CoreStatus;
    created_at: string;
    updated_at: string;
}

export interface Connection {
    id: number;
    location_node_id?: number;
    input_cable_id: number;
    input_core_id: number;
    output_cable_id: number;
    output_core_id: number;
    splice_loss?: number;
    notes?: string;
    created_at: string;
    updated_at: string;
}

export interface CreateConnectionRequest {
    location_node_id?: number;
    input_cable_id: number;
    input_core_id: number;
    output_cable_id: number;
    output_core_id: number;
    splice_loss?: number;
    notes?: string;
}

export interface Customer {
    id: number;
    node_id?: number;
    name: string;
    ont_sn?: string;
    phone?: string;
    email?: string;
    current_status: CustomerStatus;
    last_rx_power?: number;
    subscription_type?: string;
    created_at: string;
    updated_at: string;
}

export interface CreateCustomerRequest {
    node_id?: number;
    name: string;
    ont_sn?: string;
    phone?: string;
    email?: string;
    current_status?: CustomerStatus;
    subscription_type?: string;
}

// API Response types
export interface ApiResponse<T> {
    success: boolean;
    message?: string;
    data?: T;
    error?: string;
}

export interface PaginatedResponse<T> {
    success: boolean;
    data: T[];
    pagination: {
        total: number;
        limit: number;
        offset: number;
        has_more: boolean;
    };
}

// GeoJSON types
export interface GeoJSONPoint {
    type: 'Point';
    coordinates: [number, number]; // [lng, lat]
}

export interface GeoJSONLineString {
    type: 'LineString';
    coordinates: number[][];
}

export interface GeoJSONFeature<G = GeoJSONPoint, P = Record<string, unknown>> {
    type: 'Feature';
    geometry: G;
    properties: P;
}

export interface GeoJSONFeatureCollection<F = GeoJSONFeature> {
    type: 'FeatureCollection';
    features: F[];
}

// Node symbology for map
export const NODE_COLORS: Record<NodeType, string> = {
    OLT: '#1e3a8a', // Dark Blue
    ODC: '#166534', // Dark Green
    ODP: '#3b82f6', // Light Blue
    CLOSURE: '#f97316', // Orange
    POLE: '#78716c', // Gray
    CUSTOMER: '#000000', // Black
};

export const NODE_ICONS: Record<NodeType, string> = {
    OLT: '🏢',
    ODC: '📦',
    ODP: '⚪',
    CLOSURE: '💊',
    POLE: '📍',
    CUSTOMER: '🏠',
};

export const CABLE_COLORS: Record<CableType, string> = {
    ADSS: '#ef4444', // Red - Backbone
    DUCT: '#3b82f6', // Blue - Distribution
    DROP: '#000000', // Black - Drop cable
};

export const STATUS_COLORS: Record<CustomerStatus, string> = {
    ONLINE: '#22c55e', // Green
    OFFLINE: '#6b7280', // Gray
    LOS: '#ef4444', // Red
    POWER_OFF: '#f97316', // Orange
};
