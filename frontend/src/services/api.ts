import axios, { AxiosInstance, AxiosError } from 'axios';
import type {
    ApiResponse,
    PaginatedResponse,
    Node,
    CreateNodeRequest,
    UpdateNodeRequest,
    Cable,
    CreateCableRequest,
    CableCore,
    Connection,
    CreateConnectionRequest,
    Customer,
    CreateCustomerRequest,
    GeoJSONFeatureCollection,
} from '@/types';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

class ApiClient {
    private client: AxiosInstance;

    constructor() {
        this.client = axios.create({
            baseURL: API_BASE_URL,
            headers: {
                'Content-Type': 'application/json',
            },
            timeout: 10000,
        });

        // Response interceptor for error handling
        this.client.interceptors.response.use(
            (response) => response,
            (error: AxiosError<ApiResponse<unknown>>) => {
                console.error('API Error:', error.response?.data?.error || error.message);
                throw error;
            }
        );
    }

    // ==================== NODES ====================
    async getNodes(params?: {
        type?: string;
        status?: string;
        limit?: number;
        offset?: number;
    }): Promise<PaginatedResponse<Node>> {
        const response = await this.client.get<PaginatedResponse<Node>>('/api/nodes', { params });
        return response.data;
    }

    async getNodeById(id: number): Promise<Node | null> {
        const response = await this.client.get<ApiResponse<Node>>(`/api/nodes/${id}`);
        return response.data.data || null;
    }

    async createNode(data: CreateNodeRequest): Promise<Node> {
        const response = await this.client.post<ApiResponse<Node>>('/api/nodes', data);
        if (!response.data.data) throw new Error('Failed to create node');
        return response.data.data;
    }

    async updateNode(id: number, data: UpdateNodeRequest): Promise<Node> {
        const response = await this.client.put<ApiResponse<Node>>(`/api/nodes/${id}`, data);
        if (!response.data.data) throw new Error('Failed to update node');
        return response.data.data;
    }

    async deleteNode(id: number): Promise<void> {
        await this.client.delete(`/api/nodes/${id}`);
    }

    async getNearbyNodes(lat: number, lng: number, radius: number, type?: string): Promise<Node[]> {
        const response = await this.client.get<ApiResponse<Node[]>>('/api/nodes/nearby', {
            params: { lat, lng, radius, type },
        });
        return response.data.data || [];
    }

    async getNodesGeoJSON(): Promise<GeoJSONFeatureCollection> {
        const response = await this.client.get<GeoJSONFeatureCollection>('/api/geojson/nodes');
        return response.data;
    }

    // ==================== CABLES ====================
    async getCables(params?: {
        type?: string;
        status?: string;
        origin_node_id?: number;
        dest_node_id?: number;
        limit?: number;
        offset?: number;
    }): Promise<PaginatedResponse<Cable>> {
        const response = await this.client.get<PaginatedResponse<Cable>>('/api/cables', { params });
        return response.data;
    }

    async getCableById(id: number): Promise<Cable | null> {
        const response = await this.client.get<ApiResponse<Cable>>(`/api/cables/${id}`);
        return response.data.data || null;
    }

    async createCable(data: CreateCableRequest): Promise<Cable> {
        const response = await this.client.post<ApiResponse<Cable>>('/api/cables', data);
        if (!response.data.data) throw new Error('Failed to create cable');
        return response.data.data;
    }

    async deleteCable(id: number): Promise<void> {
        await this.client.delete(`/api/cables/${id}`);
    }

    async getCableCores(cableId: number): Promise<CableCore[]> {
        const response = await this.client.get<ApiResponse<CableCore[]>>(`/api/cables/${cableId}/cores`);
        return response.data.data || [];
    }

    async getCablesGeoJSON(): Promise<GeoJSONFeatureCollection> {
        const response = await this.client.get<GeoJSONFeatureCollection>('/api/geojson/cables');
        return response.data;
    }

    // ==================== CONNECTIONS ====================
    async getConnections(params?: {
        location_node_id?: number;
        limit?: number;
        offset?: number;
    }): Promise<PaginatedResponse<Connection>> {
        const response = await this.client.get<PaginatedResponse<Connection>>('/api/connections', { params });
        return response.data;
    }

    async createConnection(data: CreateConnectionRequest): Promise<Connection> {
        const response = await this.client.post<ApiResponse<Connection>>('/api/connections', data);
        if (!response.data.data) throw new Error('Failed to create connection');
        return response.data.data;
    }

    async deleteConnection(id: number): Promise<void> {
        await this.client.delete(`/api/connections/${id}`);
    }

    async getConnectionsByLocation(nodeId: number): Promise<Connection[]> {
        const response = await this.client.get<ApiResponse<Connection[]>>(`/api/connections/location/${nodeId}`);
        return response.data.data || [];
    }

    // ==================== CUSTOMERS ====================
    async getCustomers(params?: {
        node_id?: number;
        status?: string;
        search?: string;
        limit?: number;
        offset?: number;
    }): Promise<PaginatedResponse<Customer>> {
        const response = await this.client.get<PaginatedResponse<Customer>>('/api/customers', { params });
        return response.data;
    }

    async getCustomerById(id: number): Promise<Customer | null> {
        const response = await this.client.get<ApiResponse<Customer>>(`/api/customers/${id}`);
        return response.data.data || null;
    }

    async createCustomer(data: CreateCustomerRequest): Promise<Customer> {
        const response = await this.client.post<ApiResponse<Customer>>('/api/customers', data);
        if (!response.data.data) throw new Error('Failed to create customer');
        return response.data.data;
    }

    async deleteCustomer(id: number): Promise<void> {
        await this.client.delete(`/api/customers/${id}`);
    }

    async getLOSCustomers(): Promise<Customer[]> {
        const response = await this.client.get<ApiResponse<Customer[]>>('/api/customers/los');
        return response.data.data || [];
    }

    // ==================== HEALTH ====================
    async healthCheck(): Promise<{ status: string; message: string }> {
        const response = await this.client.get<{ status: string; message: string }>('/api/health');
        return response.data;
    }
}

export const api = new ApiClient();
