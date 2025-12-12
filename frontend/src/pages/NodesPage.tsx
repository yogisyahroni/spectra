import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
    Plus,
    Pencil,
    Trash2,
    MapPin,
    Network,
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
    type Column,
} from '@/components/ui';
import { NodeForm } from '@/components/nodes/NodeForm';
import { api } from '@/services/api';
import type { Node, CreateNodeRequest, NodeType, NodeStatus } from '@/types';
import { NODE_COLORS } from '@/types';

export function NodesPage() {
    const queryClient = useQueryClient();
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [selectedNode, setSelectedNode] = useState<Node | null>(null);
    const [deleteNode, setDeleteNode] = useState<Node | null>(null);

    // Fetch nodes
    const { data: nodesData, isLoading } = useQuery({
        queryKey: ['nodes'],
        queryFn: () => api.getNodes({ limit: 1000 }),
    });

    const nodes = nodesData?.data || [];

    // Create mutation
    const createMutation = useMutation({
        mutationFn: (data: CreateNodeRequest) => api.createNode(data),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['nodes'] });
            setIsModalOpen(false);
        },
    });

    // Update mutation
    const updateMutation = useMutation({
        mutationFn: ({ id, data }: { id: number; data: CreateNodeRequest }) =>
            api.updateNode(id, data),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['nodes'] });
            setIsModalOpen(false);
            setSelectedNode(null);
        },
    });

    // Delete mutation
    const deleteMutation = useMutation({
        mutationFn: (id: number) => api.deleteNode(id),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['nodes'] });
            setDeleteNode(null);
        },
    });

    const handleCreate = () => {
        setSelectedNode(null);
        setIsModalOpen(true);
    };

    const handleEdit = (node: Node) => {
        setSelectedNode(node);
        setIsModalOpen(true);
    };

    const handleSubmit = (data: CreateNodeRequest) => {
        if (selectedNode) {
            updateMutation.mutate({ id: selectedNode.id, data });
        } else {
            createMutation.mutate(data);
        }
    };

    const getStatusBadge = (status: NodeStatus) => {
        const variants: Record<NodeStatus, 'success' | 'warning' | 'info' | 'outline'> = {
            ACTIVE: 'success',
            MAINTENANCE: 'warning',
            PLAN: 'info',
            INACTIVE: 'outline',
        };
        return <Badge variant={variants[status]}>{status}</Badge>;
    };

    const getTypeBadge = (type: NodeType) => {
        return (
            <div className="flex items-center gap-2">
                <div
                    className="w-3 h-3 rounded-full"
                    style={{ backgroundColor: NODE_COLORS[type] }}
                />
                <span>{type}</span>
            </div>
        );
    };

    const columns: Column<Node>[] = [
        {
            key: 'name',
            header: 'Name',
            sortable: true,
            render: (node) => (
                <span className="font-medium text-white">{node.name}</span>
            ),
        },
        {
            key: 'type',
            header: 'Type',
            sortable: true,
            render: (node) => getTypeBadge(node.type),
        },
        {
            key: 'status',
            header: 'Status',
            sortable: true,
            render: (node) => getStatusBadge(node.status),
        },
        {
            key: 'address',
            header: 'Address',
            render: (node) => (
                <span className="text-slate-400">{node.address || '-'}</span>
            ),
        },
        {
            key: 'capacity',
            header: 'Ports',
            render: (node) => (
                <span>
                    {node.used_ports}/{node.capacity_ports}
                </span>
            ),
        },
        {
            key: 'coordinates',
            header: 'Coordinates',
            render: (node) => (
                <span className="text-xs text-slate-400 font-mono">
                    {node.latitude.toFixed(4)}, {node.longitude.toFixed(4)}
                </span>
            ),
        },
        {
            key: 'actions',
            header: 'Actions',
            width: '120px',
            render: (node) => (
                <div className="flex items-center gap-1">
                    <Button
                        size="icon"
                        variant="ghost"
                        onClick={(e) => {
                            e.stopPropagation();
                            handleEdit(node);
                        }}
                    >
                        <Pencil className="h-4 w-4" />
                    </Button>
                    <Button
                        size="icon"
                        variant="ghost"
                        onClick={(e) => {
                            e.stopPropagation();
                            setDeleteNode(node);
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
                        <Network className="h-8 w-8 text-sky-400" />
                        Nodes Management
                    </h1>
                    <p className="text-slate-400 mt-1">
                        Manage network infrastructure nodes
                    </p>
                </div>
                <Button onClick={handleCreate}>
                    <Plus className="h-4 w-4 mr-2" />
                    Add Node
                </Button>
            </div>

            {/* Stats */}
            <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4 mb-6">
                {(['OLT', 'ODC', 'ODP', 'CLOSURE', 'POLE', 'CUSTOMER'] as NodeType[]).map(
                    (type) => {
                        const count = nodes.filter((n) => n.type === type).length;
                        return (
                            <Card key={type} variant="glass" className="p-4">
                                <div className="flex items-center gap-3">
                                    <div
                                        className="w-3 h-3 rounded-full"
                                        style={{ backgroundColor: NODE_COLORS[type] }}
                                    />
                                    <div>
                                        <p className="text-2xl font-bold text-white">{count}</p>
                                        <p className="text-xs text-slate-400">{type}</p>
                                    </div>
                                </div>
                            </Card>
                        );
                    }
                )}
            </div>

            {/* Table */}
            <Card variant="glass">
                <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                        <MapPin className="h-5 w-5 text-sky-400" />
                        All Nodes
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    <DataTable
                        data={nodes}
                        columns={columns}
                        keyField="id"
                        isLoading={isLoading}
                        searchable
                        searchPlaceholder="Search nodes..."
                        emptyMessage="No nodes found. Click 'Add Node' to create one."
                    />
                </CardContent>
            </Card>

            {/* Create/Edit Modal */}
            <Modal
                isOpen={isModalOpen}
                onClose={() => {
                    setIsModalOpen(false);
                    setSelectedNode(null);
                }}
                title={selectedNode ? 'Edit Node' : 'Create Node'}
                description={
                    selectedNode
                        ? 'Update the node information below.'
                        : 'Fill in the details to create a new node.'
                }
                size="lg"
            >
                <NodeForm
                    node={selectedNode}
                    onSubmit={handleSubmit}
                    onCancel={() => {
                        setIsModalOpen(false);
                        setSelectedNode(null);
                    }}
                    isLoading={createMutation.isPending || updateMutation.isPending}
                />
            </Modal>

            {/* Delete Confirmation */}
            <ConfirmDialog
                isOpen={!!deleteNode}
                onClose={() => setDeleteNode(null)}
                onConfirm={() => deleteNode && deleteMutation.mutate(deleteNode.id)}
                title="Delete Node"
                message={`Are you sure you want to delete "${deleteNode?.name}"? This action cannot be undone.`}
                confirmText="Delete"
                variant="danger"
                isLoading={deleteMutation.isPending}
            />
        </div>
    );
}
