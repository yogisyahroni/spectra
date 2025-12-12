import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
    Plus,
    Trash2,
    Users,
    Wifi,
    WifiOff,
    AlertTriangle,
    Signal,
    Search,
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
import { api } from '@/services/api';
import type { Customer, CreateCustomerRequest, CustomerStatus } from '@/types';
import { STATUS_COLORS } from '@/types';
import { cn, formatRxPower, getRxPowerStatus } from '@/lib/utils';

const statusOptions: SelectOption[] = [
    { value: 'ONLINE', label: 'Online' },
    { value: 'OFFLINE', label: 'Offline' },
    { value: 'LOS', label: 'LOS (Loss of Signal)' },
    { value: 'POWER_OFF', label: 'Power Off' },
];

const filterOptions: SelectOption[] = [
    { value: 'ALL', label: 'All Status' },
    { value: 'ONLINE', label: 'Online' },
    { value: 'OFFLINE', label: 'Offline' },
    { value: 'LOS', label: 'LOS' },
    { value: 'POWER_OFF', label: 'Power Off' },
];

export function CustomersPage() {
    const queryClient = useQueryClient();
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [deleteCustomer, setDeleteCustomer] = useState<Customer | null>(null);
    const [statusFilter, setStatusFilter] = useState('ALL');

    // Form state
    const [formData, setFormData] = useState({
        name: '',
        ont_sn: '',
        phone: '',
        email: '',
        current_status: 'OFFLINE' as CustomerStatus,
        subscription_type: '',
    });

    // Fetch customers
    const { data: customersData, isLoading } = useQuery({
        queryKey: ['customers'],
        queryFn: () => api.getCustomers({ limit: 1000 }),
    });

    // Fetch LOS customers
    const { data: losCustomers } = useQuery({
        queryKey: ['customers-los'],
        queryFn: () => api.getLOSCustomers(),
    });

    const allCustomers = customersData?.data || [];
    const customers = statusFilter === 'ALL'
        ? allCustomers
        : allCustomers.filter(c => c.current_status === statusFilter);

    // Create mutation
    const createMutation = useMutation({
        mutationFn: (data: CreateCustomerRequest) => api.createCustomer(data),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['customers'] });
            setIsModalOpen(false);
            resetForm();
        },
    });

    // Delete mutation
    const deleteMutation = useMutation({
        mutationFn: (id: number) => api.deleteCustomer(id),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['customers'] });
            setDeleteCustomer(null);
        },
    });

    const resetForm = () => {
        setFormData({
            name: '',
            ont_sn: '',
            phone: '',
            email: '',
            current_status: 'OFFLINE',
            subscription_type: '',
        });
    };

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        createMutation.mutate(formData);
    };

    const getStatusBadge = (status: CustomerStatus) => {
        const variants: Record<CustomerStatus, 'success' | 'outline' | 'error' | 'warning'> = {
            ONLINE: 'success',
            OFFLINE: 'outline',
            LOS: 'error',
            POWER_OFF: 'warning',
        };
        const icons: Record<CustomerStatus, React.ReactNode> = {
            ONLINE: <Wifi className="h-3 w-3" />,
            OFFLINE: <WifiOff className="h-3 w-3" />,
            LOS: <AlertTriangle className="h-3 w-3" />,
            POWER_OFF: <WifiOff className="h-3 w-3" />,
        };
        return (
            <Badge variant={variants[status]} className="flex items-center gap-1.5">
                {icons[status]}
                {status}
            </Badge>
        );
    };

    const getRxBadge = (rxPower: number | undefined | null) => {
        const status = getRxPowerStatus(rxPower);
        const colors: Record<string, string> = {
            good: 'text-green-400',
            warning: 'text-yellow-400',
            critical: 'text-red-400',
            unknown: 'text-slate-400',
        };
        return (
            <span className={cn('font-mono text-sm', colors[status])}>
                {formatRxPower(rxPower)}
            </span>
        );
    };

    const columns: Column<Customer>[] = [
        {
            key: 'name',
            header: 'Name',
            sortable: true,
            render: (customer) => (
                <span className="font-medium text-white">{customer.name}</span>
            ),
        },
        {
            key: 'ont_sn',
            header: 'ONT S/N',
            render: (customer) => (
                <span className="font-mono text-sm text-slate-400">
                    {customer.ont_sn || '-'}
                </span>
            ),
        },
        {
            key: 'current_status',
            header: 'Status',
            sortable: true,
            render: (customer) => getStatusBadge(customer.current_status),
        },
        {
            key: 'last_rx_power',
            header: 'Rx Power',
            sortable: true,
            render: (customer) => getRxBadge(customer.last_rx_power),
        },
        {
            key: 'phone',
            header: 'Phone',
            render: (customer) => (
                <span className="text-slate-400">{customer.phone || '-'}</span>
            ),
        },
        {
            key: 'subscription_type',
            header: 'Plan',
            render: (customer) => (
                <Badge variant="outline">{customer.subscription_type || 'N/A'}</Badge>
            ),
        },
        {
            key: 'actions',
            header: 'Actions',
            width: '80px',
            render: (customer) => (
                <Button
                    size="icon"
                    variant="ghost"
                    onClick={(e) => {
                        e.stopPropagation();
                        setDeleteCustomer(customer);
                    }}
                >
                    <Trash2 className="h-4 w-4 text-red-400" />
                </Button>
            ),
        },
    ];

    // Stats
    const onlineCount = allCustomers.filter(c => c.current_status === 'ONLINE').length;
    const offlineCount = allCustomers.filter(c => c.current_status === 'OFFLINE').length;
    const losCount = losCustomers?.length || 0;
    const powerOffCount = allCustomers.filter(c => c.current_status === 'POWER_OFF').length;

    return (
        <div className="p-6">
            {/* Header */}
            <div className="flex items-center justify-between mb-6">
                <div>
                    <h1 className="text-3xl font-bold text-white flex items-center gap-3">
                        <Users className="h-8 w-8 text-emerald-400" />
                        Customers Management
                    </h1>
                    <p className="text-slate-400 mt-1">
                        Monitor customer status and manage subscribers
                    </p>
                </div>
                <Button onClick={() => setIsModalOpen(true)}>
                    <Plus className="h-4 w-4 mr-2" />
                    Add Customer
                </Button>
            </div>

            {/* Stats Cards */}
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
                <Card variant="glass" className="p-4">
                    <div className="flex items-center gap-3">
                        <div className="p-2 rounded-lg bg-green-500/20">
                            <Wifi className="h-5 w-5 text-green-400" />
                        </div>
                        <div>
                            <p className="text-2xl font-bold text-white">{onlineCount}</p>
                            <p className="text-xs text-slate-400">Online</p>
                        </div>
                    </div>
                </Card>
                <Card variant="glass" className="p-4">
                    <div className="flex items-center gap-3">
                        <div className="p-2 rounded-lg bg-slate-500/20">
                            <WifiOff className="h-5 w-5 text-slate-400" />
                        </div>
                        <div>
                            <p className="text-2xl font-bold text-white">{offlineCount}</p>
                            <p className="text-xs text-slate-400">Offline</p>
                        </div>
                    </div>
                </Card>
                <Card variant="glass" className="p-4 border-red-500/30">
                    <div className="flex items-center gap-3">
                        <div className="p-2 rounded-lg bg-red-500/20">
                            <AlertTriangle className="h-5 w-5 text-red-400" />
                        </div>
                        <div>
                            <p className="text-2xl font-bold text-red-400">{losCount}</p>
                            <p className="text-xs text-slate-400">LOS Alerts</p>
                        </div>
                    </div>
                </Card>
                <Card variant="glass" className="p-4">
                    <div className="flex items-center gap-3">
                        <div className="p-2 rounded-lg bg-amber-500/20">
                            <Signal className="h-5 w-5 text-amber-400" />
                        </div>
                        <div>
                            <p className="text-2xl font-bold text-white">{powerOffCount}</p>
                            <p className="text-xs text-slate-400">Power Off</p>
                        </div>
                    </div>
                </Card>
            </div>

            {/* LOS Alert Panel */}
            {losCount > 0 && (
                <Card variant="glass" className="mb-6 border-red-500/30 bg-red-500/5">
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2 text-red-400">
                            <AlertTriangle className="h-5 w-5" />
                            Active LOS Alerts ({losCount})
                        </CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="flex flex-wrap gap-2">
                            {losCustomers?.slice(0, 10).map((customer) => (
                                <Badge key={customer.id} variant="error" className="text-sm">
                                    {customer.name} ({customer.ont_sn || 'No ONT'})
                                </Badge>
                            ))}
                            {losCount > 10 && (
                                <Badge variant="outline">+{losCount - 10} more</Badge>
                            )}
                        </div>
                    </CardContent>
                </Card>
            )}

            {/* Filter and Table */}
            <Card variant="glass">
                <CardHeader>
                    <div className="flex items-center justify-between">
                        <CardTitle className="flex items-center gap-2">
                            <Users className="h-5 w-5 text-emerald-400" />
                            All Customers
                        </CardTitle>
                        <div className="w-48">
                            <Select
                                options={filterOptions}
                                value={statusFilter}
                                onChange={(val) => setStatusFilter(val)}
                            />
                        </div>
                    </div>
                </CardHeader>
                <CardContent>
                    <DataTable
                        data={customers}
                        columns={columns}
                        keyField="id"
                        isLoading={isLoading}
                        searchable
                        searchPlaceholder="Search by name or ONT..."
                        emptyMessage="No customers found."
                    />
                </CardContent>
            </Card>

            {/* Create Modal */}
            <Modal
                isOpen={isModalOpen}
                onClose={() => setIsModalOpen(false)}
                title="Add Customer"
                description="Add a new customer subscriber."
                size="lg"
            >
                <form onSubmit={handleSubmit} className="space-y-4">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <Input
                            label="Customer Name"
                            placeholder="John Doe"
                            value={formData.name}
                            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                            required
                        />

                        <Input
                            label="ONT Serial Number"
                            placeholder="HWTC12345678"
                            value={formData.ont_sn}
                            onChange={(e) => setFormData({ ...formData, ont_sn: e.target.value })}
                        />

                        <Input
                            label="Phone"
                            placeholder="+62 812 3456 7890"
                            value={formData.phone}
                            onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
                        />

                        <Input
                            label="Email"
                            type="email"
                            placeholder="customer@email.com"
                            value={formData.email}
                            onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                        />

                        <Input
                            label="Subscription Plan"
                            placeholder="e.g., 20 Mbps"
                            value={formData.subscription_type}
                            onChange={(e) => setFormData({ ...formData, subscription_type: e.target.value })}
                        />

                        <Select
                            label="Initial Status"
                            options={statusOptions}
                            value={formData.current_status}
                            onChange={(val) => setFormData({ ...formData, current_status: val as CustomerStatus })}
                        />
                    </div>

                    <div className="flex justify-end gap-3 pt-4 border-t border-slate-700/50">
                        <Button type="button" variant="outline" onClick={() => setIsModalOpen(false)}>
                            Cancel
                        </Button>
                        <Button type="submit" isLoading={createMutation.isPending}>
                            Add Customer
                        </Button>
                    </div>
                </form>
            </Modal>

            {/* Delete Confirmation */}
            <ConfirmDialog
                isOpen={!!deleteCustomer}
                onClose={() => setDeleteCustomer(null)}
                onConfirm={() => deleteCustomer && deleteMutation.mutate(deleteCustomer.id)}
                title="Delete Customer"
                message={`Are you sure you want to delete "${deleteCustomer?.name}"?`}
                confirmText="Delete"
                variant="danger"
                isLoading={deleteMutation.isPending}
            />
        </div>
    );
}
