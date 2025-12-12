import React from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Button, Input, Select, type SelectOption } from '@/components/ui';
import type { Node, NodeType, NodeStatus, CreateNodeRequest } from '@/types';

// Validation schema
const nodeSchema = z.object({
    name: z.string().min(1, 'Name is required').max(100, 'Name too long'),
    type: z.enum(['OLT', 'ODC', 'ODP', 'CLOSURE', 'POLE', 'CUSTOMER'] as const),
    latitude: z.coerce.number().min(-90).max(90, 'Invalid latitude'),
    longitude: z.coerce.number().min(-180).max(180, 'Invalid longitude'),
    address: z.string().optional(),
    capacity_ports: z.coerce.number().int().min(1).max(1000).optional(),
    model: z.string().optional(),
    status: z.enum(['ACTIVE', 'MAINTENANCE', 'PLAN', 'INACTIVE'] as const).optional(),
});

type NodeFormData = z.infer<typeof nodeSchema>;

interface NodeFormProps {
    node?: Node | null;
    onSubmit: (data: CreateNodeRequest) => void;
    onCancel: () => void;
    isLoading?: boolean;
}

const typeOptions: SelectOption[] = [
    { value: 'OLT', label: 'OLT - Optical Line Terminal' },
    { value: 'ODC', label: 'ODC - Optical Distribution Cabinet' },
    { value: 'ODP', label: 'ODP - Optical Distribution Point' },
    { value: 'CLOSURE', label: 'CLOSURE - Splice Closure' },
    { value: 'POLE', label: 'POLE - Utility Pole' },
    { value: 'CUSTOMER', label: 'CUSTOMER - Customer Location' },
];

const statusOptions: SelectOption[] = [
    { value: 'ACTIVE', label: 'Active' },
    { value: 'MAINTENANCE', label: 'Maintenance' },
    { value: 'PLAN', label: 'Planned' },
    { value: 'INACTIVE', label: 'Inactive' },
];

export function NodeForm({ node, onSubmit, onCancel, isLoading }: NodeFormProps) {
    const {
        register,
        handleSubmit,
        setValue,
        watch,
        formState: { errors },
    } = useForm({
        resolver: zodResolver(nodeSchema) as any,
        defaultValues: {
            name: node?.name || '',
            type: (node?.type || 'ODP') as NodeType,
            latitude: node?.latitude || 0,
            longitude: node?.longitude || 0,
            address: node?.address || '',
            capacity_ports: node?.capacity_ports || 8,
            model: node?.model || '',
            status: (node?.status || 'ACTIVE') as NodeStatus,
        },
    });

    const selectedType = watch('type');
    const selectedStatus = watch('status');

    const onFormSubmit = (data: NodeFormData) => {
        onSubmit({
            name: data.name,
            type: data.type as NodeType,
            latitude: data.latitude,
            longitude: data.longitude,
            address: data.address,
            capacity_ports: data.capacity_ports,
            model: data.model,
            status: data.status as NodeStatus,
        });
    };

    return (
        <form onSubmit={handleSubmit(onFormSubmit)} className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <Input
                    label="Node Name"
                    placeholder="e.g., ODP-JATI-01"
                    {...register('name')}
                    error={errors.name?.message}
                />

                <Select
                    label="Node Type"
                    options={typeOptions}
                    value={selectedType}
                    onChange={(val) => setValue('type', val as NodeType)}
                    error={errors.type?.message}
                />

                <Input
                    label="Latitude"
                    type="number"
                    step="any"
                    placeholder="-6.2088"
                    {...register('latitude')}
                    error={errors.latitude?.message}
                />

                <Input
                    label="Longitude"
                    type="number"
                    step="any"
                    placeholder="106.8456"
                    {...register('longitude')}
                    error={errors.longitude?.message}
                />

                <Input
                    label="Address"
                    placeholder="Jl. Sudirman No. 1"
                    {...register('address')}
                    error={errors.address?.message}
                />

                <Input
                    label="Capacity Ports"
                    type="number"
                    placeholder="8"
                    {...register('capacity_ports')}
                    error={errors.capacity_ports?.message}
                />

                <Input
                    label="Model"
                    placeholder="e.g., Huawei MA5800-X7"
                    {...register('model')}
                    error={errors.model?.message}
                />

                <Select
                    label="Status"
                    options={statusOptions}
                    value={selectedStatus}
                    onChange={(val) => setValue('status', val as NodeStatus)}
                    error={errors.status?.message}
                />
            </div>

            <div className="flex justify-end gap-3 pt-4 border-t border-slate-700/50">
                <Button type="button" variant="outline" onClick={onCancel}>
                    Cancel
                </Button>
                <Button type="submit" isLoading={isLoading}>
                    {node ? 'Update Node' : 'Create Node'}
                </Button>
            </div>
        </form>
    );
}
