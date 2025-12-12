import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
    Plus,
    Pencil,
    Trash2,
    Cable,
    Eye,
} from 'lucide-react';
import {
    Card,
    CardHeader,
    CardTitle,
    CardContent,
    Button,
    Badge,
    DataTable,
    Modal,
    ConfirmDialog,
    Input,
    Select,
    type Column,
    type SelectOption,
} from '@/components/ui';
import { CoreVisualization } from '@/components/cables/CoreVisualization';
import { api } from '@/services/api';
import type { Cable as CableType, CableCore, CreateCableRequest, CableType as CType, CableStatus } from '@/types';
import { CABLE_COLORS } from '@/types';

const cableTypeOptions: SelectOption[] = [
    { value: 'ADSS', label: 'ADSS - All-Dielectric Self-Supporting' },
    { value: 'DUCT', label: 'DUCT - Underground Duct Cable' },
    { value: 'DROP', label: 'DROP - Drop Cable' },
];

const statusOptions: SelectOption[] = [
    { value: 'ACTIVE', label: 'Active' },
    { value: 'MAINTENANCE', label: 'Maintenance' },
    { value: 'PLAN', label: 'Planned' },
    { value: 'INACTIVE', label: 'Inactive' },
];

const coreCountOptions: SelectOption[] = [
    { value: '4', label: '4 Core' },
    { value: '6', label: '6 Core' },
    { value: '8', label: '8 Core' },
    { value: '12', label: '12 Core' },
    { value: '24', label: '24 Core' },
    { value: '48', label: '48 Core' },
    { value: '96', label: '96 Core' },
    { value: '144', label: '144 Core' },
];

export function CablesPage() {
    const queryClient = useQueryClient();
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [viewCable, setViewCable] = useState<CableType | null>(null);
    const [deleteCable, setDeleteCable] = useState<CableType | null>(null);
    const [cableCores, setCableCores] = useState<CableCore[]>([]);

    // Form state
    const [formData, setFormData] = useState({
        name: '',
        type: 'ADSS' as CType,
        core_count: 24,
        length_meter: 0,
        color_hex: '#ef4444',
        status: 'ACTIVE' as CableStatus,
    });

    // Fetch cables
    const { data: cablesData, isLoading } = useQuery({
        queryKey: ['cables'],
        queryFn: () => api.getCables({ limit: 1000 }),
    });

    const cables = cablesData?.data || [];

    // Create mutation
    const createMutation = useMutation({
        mutationFn: (data: CreateCableRequest) => api.createCable(data),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['cables'] });
            setIsModalOpen(false);
            resetForm();
        },
    });

    // Delete mutation
    const deleteMutation = useMutation({
        mutationFn: (id: number) => api.deleteCable(id),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['cables'] });
            setDeleteCable(null);
        },
    });

    const resetForm = () => {
        setFormData({
            name: '',
            type: 'ADSS',
            core_count: 24,
            length_meter: 0,
            color_hex: '#ef4444',
            status: 'ACTIVE',
        });
    };

    const handleCreate = () => {
        resetForm();
        setIsModalOpen(true);
    };

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        createMutation.mutate({
            name: formData.name,
            type: formData.type,
            core_count: formData.core_count,
            length_meter: formData.length_meter,
            color_hex: formData.color_hex,
            status: formData.status,
        });
    };

    const handleViewCores = async (cable: CableType) => {
        setViewCable(cable);
        try {
            const cores = await api.getCableCores(cable.id);
            setCableCores(cores);
        } catch (error) {
            console.error('Failed to fetch cores:', error);
            setCableCores([]);
        }
    };

    const getTypeBadge = (type: CType) => {
        return (
            <div className="flex items-center gap-2">
                <div
                    className="w-3 h-3 rounded"
                    style={{ backgroundColor: CABLE_COLORS[type] }}
                />
                <span>{type}</span>
            </div>
        );
    };

    const getStatusBadge = (status: CableStatus) => {
        const variants: Record<CableStatus, 'success' | 'warning' | 'info' | 'outline'> = {
            ACTIVE: 'success',
            MAINTENANCE: 'warning',
            PLAN: 'info',
            INACTIVE: 'outline',
        };
        return <Badge variant={variants[status]}>{status}</Badge>;
    };

    const columns: Column<CableType>[] = [
        {
            key: 'name',
            header: 'Name',
            sortable: true,
            render: (cable) => (
                <span className="font-medium text-white">{cable.name || `Cable-${cable.id}`}</span>
            ),
        },
        {
            key: 'type',
            header: 'Type',
            sortable: true,
            render: (cable) => getTypeBadge(cable.type),
        },
        {
            key: 'core_count',
            header: 'Cores',
            sortable: true,
            render: (cable) => (
                <Badge variant="outline">{cable.core_count} Core</Badge>
            ),
        },
        {
            key: 'status',
            header: 'Status',
            sortable: true,
            render: (cable) => getStatusBadge(cable.status),
        },
        {
            key: 'length_meter',
            header: 'Length',
            render: (cable) => (
                <span className="text-slate-400">
                    {cable.length_meter ? `${cable.length_meter.toFixed(0)} m` : '-'}
                </span>
            ),
        },
        {
            key: 'actions',
            header: 'Actions',
            width: '120px',
            render: (cable) => (
                <div className="flex items-center gap-1">
                    <Button
                        size="icon"
                        variant="ghost"
                        onClick={(e) => {
                            e.stopPropagation();
                            handleViewCores(cable);
                        }}
                        title="View Cores"
                    >
                        <Eye className="h-4 w-4 text-sky-400" />
                    </Button>
                    <Button
                        size="icon"
                        variant="ghost"
                        onClick={(e) => {
                            e.stopPropagation();
                            setDeleteCable(cable);
                        }}
                    >
                        <Trash2 className="h-4 w-4 text-red-400" />
                    </Button>
                </div>
            ),
        },
    ];

    return (
        <div className="p-6">
            {/* Header */}
            <div className="flex items-center justify-between mb-6">
                <div>
                    <h1 className="text-3xl font-bold text-white flex items-center gap-3">
                        <Cable className="h-8 w-8 text-violet-400" />
                        Cables Management
                    </h1>
                    <p className="text-slate-400 mt-1">
                        Manage fiber optic cables and cores
                    </p>
                </div>
                <Button onClick={handleCreate}>
                    <Plus className="h-4 w-4 mr-2" />
                    Add Cable
                </Button>
            </div>

            {/* Stats */}
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
                {(['ADSS', 'DUCT', 'DROP'] as CType[]).map((type) => {
                    const count = cables.filter((c) => c.type === type).length;
                    const totalCores = cables
                        .filter((c) => c.type === type)
                        .reduce((sum, c) => sum + c.core_count, 0);
                    return (
                        <Card key={type} variant="glass" className="p-4">
                            <div className="flex items-center gap-3">
                                <div
                                    className="w-3 h-3 rounded"
                                    style={{ backgroundColor: CABLE_COLORS[type] }}
                                />
                                <div>
                                    <p className="text-2xl font-bold text-white">{count}</p>
                                    <p className="text-xs text-slate-400">{type} ({totalCores} cores)</p>
                                </div>
                            </div>
                        </Card>
                    );
                })}
                <Card variant="glass" className="p-4">
                    <div>
                        <p className="text-2xl font-bold text-white">
                            {cables.reduce((sum, c) => sum + c.core_count, 0)}
                        </p>
                        <p className="text-xs text-slate-400">Total Cores</p>
                    </div>
                </Card>
            </div>

            {/* Table */}
            <Card variant="glass">
                <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                        <Cable className="h-5 w-5 text-violet-400" />
                        All Cables
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    <DataTable
                        data={cables}
                        columns={columns}
                        keyField="id"
                        isLoading={isLoading}
                        searchable
                        searchPlaceholder="Search cables..."
                        emptyMessage="No cables found. Click 'Add Cable' to create one."
                    />
                </CardContent>
            </Card>

            {/* Create Modal */}
            <Modal
                isOpen={isModalOpen}
                onClose={() => setIsModalOpen(false)}
                title="Create Cable"
                description="Fill in the details to create a new fiber cable."
                size="lg"
            >
                <form onSubmit={handleSubmit} className="space-y-4">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <Input
                            label="Cable Name"
                            placeholder="e.g., Feeder-JKT-01"
                            value={formData.name}
                            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                        />

                        <Select
                            label="Cable Type"
                            options={cableTypeOptions}
                            value={formData.type}
                            onChange={(val) => setFormData({ ...formData, type: val as CType })}
                        />

                        <Select
                            label="Core Count"
                            options={coreCountOptions}
                            value={String(formData.core_count)}
                            onChange={(val) => setFormData({ ...formData, core_count: parseInt(val) })}
                        />

                        <Input
                            label="Length (meters)"
                            type="number"
                            placeholder="500"
                            value={formData.length_meter || ''}
                            onChange={(e) => setFormData({ ...formData, length_meter: parseFloat(e.target.value) || 0 })}
                        />

                        <Select
                            label="Status"
                            options={statusOptions}
                            value={formData.status}
                            onChange={(val) => setFormData({ ...formData, status: val as CableStatus })}
                        />

                        <div className="space-y-1.5">
                            <label className="block text-sm font-medium text-slate-300">Cable Color</label>
                            <div className="flex gap-2">
                                <input
                                    type="color"
                                    value={formData.color_hex}
                                    onChange={(e) => setFormData({ ...formData, color_hex: e.target.value })}
                                    className="h-10 w-20 rounded-lg border border-slate-700/50 bg-slate-800/50 cursor-pointer"
                                />
                                <Input
                                    value={formData.color_hex}
                                    onChange={(e) => setFormData({ ...formData, color_hex: e.target.value })}
                                    className="flex-1"
                                />
                            </div>
                        </div>
                    </div>

                    <div className="flex justify-end gap-3 pt-4 border-t border-slate-700/50">
                        <Button type="button" variant="outline" onClick={() => setIsModalOpen(false)}>
                            Cancel
                        </Button>
                        <Button type="submit" isLoading={createMutation.isPending}>
                            Create Cable
                        </Button>
                    </div>
                </form>
            </Modal>

            {/* View Cores Modal */}
            <Modal
                isOpen={!!viewCable}
                onClose={() => {
                    setViewCable(null);
                    setCableCores([]);
                }}
                title={`Cable Cores: ${viewCable?.name || `Cable-${viewCable?.id}`}`}
                description={`${viewCable?.type} cable with ${viewCable?.core_count} cores`}
                size="xl"
            >
                <CoreVisualization
                    cores={cableCores}
                    coreCount={viewCable?.core_count || 0}
                />
            </Modal>

            {/* Delete Confirmation */}
            <ConfirmDialog
                isOpen={!!deleteCable}
                onClose={() => setDeleteCable(null)}
                onConfirm={() => deleteCable && deleteMutation.mutate(deleteCable.id)}
                title="Delete Cable"
                message={`Are you sure you want to delete "${deleteCable?.name || `Cable-${deleteCable?.id}`}"? This will also delete all associated cores.`}
                confirmText="Delete"
                variant="danger"
                isLoading={deleteMutation.isPending}
            />
        </div>
    );
}
