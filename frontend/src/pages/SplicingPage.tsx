import React, { useState, useEffect } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Zap, Network, Cable, ChevronRight } from 'lucide-react';
import {
    Card,
    CardHeader,
    CardTitle,
    CardContent,
    Button,
    Badge,
    Select,
    type SelectOption,
} from '@/components/ui';
import { SpliceMatrix } from '@/components/splice/SpliceMatrix';
import { api } from '@/services/api';
import type { Node, Cable as CableType, CableCore, Connection } from '@/types';

export function SplicingPage() {
    const queryClient = useQueryClient();
    const [selectedNode, setSelectedNode] = useState<Node | null>(null);
    const [inputCable, setInputCable] = useState<CableType | null>(null);
    const [outputCable, setOutputCable] = useState<CableType | null>(null);
    const [inputCores, setInputCores] = useState<CableCore[]>([]);
    const [outputCores, setOutputCores] = useState<CableCore[]>([]);

    // Fetch nodes (only splice-capable: ODC, ODP, CLOSURE)
    const { data: nodesData } = useQuery({
        queryKey: ['nodes'],
        queryFn: () => api.getNodes({ limit: 1000 }),
    });

    const spliceNodes = (nodesData?.data || []).filter((n) =>
        ['ODC', 'ODP', 'CLOSURE'].includes(n.type)
    );

    // Fetch cables
    const { data: cablesData } = useQuery({
        queryKey: ['cables'],
        queryFn: () => api.getCables({ limit: 1000 }),
    });

    const cables = cablesData?.data || [];

    // Fetch connections for selected node
    const { data: connections = [], refetch: refetchConnections } = useQuery({
        queryKey: ['connections', selectedNode?.id],
        queryFn: () => (selectedNode ? api.getConnectionsByLocation(selectedNode.id) : []),
        enabled: !!selectedNode,
    });

    // Load cores when cables are selected
    useEffect(() => {
        if (inputCable) {
            api.getCableCores(inputCable.id).then(setInputCores).catch(() => setInputCores([]));
        } else {
            setInputCores([]);
        }
    }, [inputCable]);

    useEffect(() => {
        if (outputCable) {
            api.getCableCores(outputCable.id).then(setOutputCores).catch(() => setOutputCores([]));
        } else {
            setOutputCores([]);
        }
    }, [outputCable]);

    // Create connection mutation
    const createMutation = useMutation({
        mutationFn: (data: { input_core_id: number; output_core_id: number; splice_loss: number }) =>
            api.createConnection({
                location_node_id: selectedNode!.id,
                input_cable_id: inputCable!.id,
                input_core_id: data.input_core_id,
                output_cable_id: outputCable!.id,
                output_core_id: data.output_core_id,
                splice_loss: data.splice_loss,
            }),
        onSuccess: () => {
            refetchConnections();
            queryClient.invalidateQueries({ queryKey: ['cables'] });
        },
    });

    // Delete connection mutation
    const deleteMutation = useMutation({
        mutationFn: (id: number) => api.deleteConnection(id),
        onSuccess: () => {
            refetchConnections();
            queryClient.invalidateQueries({ queryKey: ['cables'] });
        },
    });

    const handleCreateConnection = (inputCoreId: number, outputCoreId: number, loss: number) => {
        createMutation.mutate({
            input_core_id: inputCoreId,
            output_core_id: outputCoreId,
            splice_loss: loss,
        });
    };

    const nodeOptions: SelectOption[] = spliceNodes.map((n) => ({
        value: String(n.id),
        label: `${n.name} (${n.type})`,
    }));

    const cableOptions: SelectOption[] = cables.map((c) => ({
        value: String(c.id),
        label: `${c.name || `Cable-${c.id}`} (${c.type}, ${c.core_count} cores)`,
    }));

    const handleNodeChange = (val: string) => {
        const node = spliceNodes.find((n) => n.id === parseInt(val));
        setSelectedNode(node || null);
        setInputCable(null);
        setOutputCable(null);
    };

    const handleInputCableChange = (val: string) => {
        const cable = cables.find((c) => c.id === parseInt(val));
        setInputCable(cable || null);
    };

    const handleOutputCableChange = (val: string) => {
        const cable = cables.find((c) => c.id === parseInt(val));
        setOutputCable(cable || null);
    };

    // Filter connections for current cable pair
    const relevantConnections = inputCable && outputCable
        ? connections.filter(
            (c) =>
                c.input_cable_id === inputCable.id && c.output_cable_id === outputCable.id
        )
        : [];

    // Map connections to include core indices
    const mappedConnections = relevantConnections.map((conn) => ({
        ...conn,
        input_core_index: inputCores.find((c) => c.id === conn.input_core_id)?.core_index || 0,
        output_core_index: outputCores.find((c) => c.id === conn.output_core_id)?.core_index || 0,
    }));

    return (
        <div className="p-6">
            {/* Header */}
            <div className="mb-6">
                <h1 className="text-3xl font-bold text-white flex items-center gap-3">
                    <Zap className="h-8 w-8 text-amber-400" />
                    Splicing Matrix
                </h1>
                <p className="text-slate-400 mt-1">
                    Manage fiber core connections at splice points
                </p>
            </div>

            {/* Selection Panel */}
            <Card variant="glass" className="mb-6">
                <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                        <Network className="h-5 w-5 text-sky-400" />
                        Select Splice Location
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4 items-end">
                        <div>
                            <Select
                                label="Splice Node"
                                placeholder="Select a node..."
                                options={nodeOptions}
                                value={selectedNode ? String(selectedNode.id) : ''}
                                onChange={handleNodeChange}
                            />
                        </div>

                        {selectedNode && (
                            <>
                                <div>
                                    <Select
                                        label="Input Cable (From)"
                                        placeholder="Select input cable..."
                                        options={cableOptions}
                                        value={inputCable ? String(inputCable.id) : ''}
                                        onChange={handleInputCableChange}
                                    />
                                </div>
                                <div>
                                    <Select
                                        label="Output Cable (To)"
                                        placeholder="Select output cable..."
                                        options={cableOptions.filter((c) => c.value !== String(inputCable?.id))}
                                        value={outputCable ? String(outputCable.id) : ''}
                                        onChange={handleOutputCableChange}
                                    />
                                </div>
                            </>
                        )}
                    </div>

                    {selectedNode && inputCable && outputCable && (
                        <div className="mt-4 p-3 rounded-lg bg-slate-800/50 flex items-center gap-3">
                            <Badge variant="info">{inputCable.name || `Cable-${inputCable.id}`}</Badge>
                            <ChevronRight className="h-4 w-4 text-slate-400" />
                            <Badge variant="warning">{selectedNode.name}</Badge>
                            <ChevronRight className="h-4 w-4 text-slate-400" />
                            <Badge variant="success">{outputCable.name || `Cable-${outputCable.id}`}</Badge>
                        </div>
                    )}
                </CardContent>
            </Card>

            {/* Splice Matrix */}
            {selectedNode && inputCable && outputCable && inputCores.length > 0 && outputCores.length > 0 ? (
                <Card variant="glass">
                    <CardContent className="p-6">
                        <SpliceMatrix
                            nodeId={selectedNode.id}
                            nodeName={selectedNode.name}
                            inputCores={inputCores}
                            outputCores={outputCores}
                            connections={mappedConnections}
                            inputCableId={inputCable.id}
                            outputCableId={outputCable.id}
                            onCreateConnection={handleCreateConnection}
                            onDeleteConnection={(id) => deleteMutation.mutate(id)}
                            isLoading={createMutation.isPending || deleteMutation.isPending}
                        />
                    </CardContent>
                </Card>
            ) : (
                <Card variant="glass">
                    <CardContent className="p-12 text-center">
                        <Cable className="h-12 w-12 text-slate-600 mx-auto mb-4" />
                        <h3 className="text-lg font-medium text-white mb-2">
                            Select Cables to Start
                        </h3>
                        <p className="text-slate-400 max-w-md mx-auto">
                            Choose a splice node (ODC, ODP, or Closure), then select the input and output cables to manage their fiber core connections.
                        </p>
                    </CardContent>
                </Card>
            )}
        </div>
    );
}
