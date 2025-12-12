import React, { useState } from 'react';
import { Plus, Trash2, ArrowRight, Zap } from 'lucide-react';
import { cn } from '@/lib/utils';
import { Button, Badge, Input, Modal, ConfirmDialog } from '@/components/ui';
import { COLOR_HEX, FIBER_COLORS } from '@/components/cables/CoreVisualization';

interface Core {
    id: number;
    core_index: number;
    tube_color?: string;
    core_color?: string;
    status: string;
}

interface Connection {
    id: number;
    input_cable_id: number;
    input_core_id: number;
    input_core_index: number;
    output_cable_id: number;
    output_core_id: number;
    output_core_index: number;
    splice_loss?: number;
}

interface SpliceMatrixProps {
    nodeId: number;
    nodeName: string;
    inputCores: Core[];
    outputCores: Core[];
    connections: Connection[];
    inputCableId: number;
    outputCableId: number;
    onCreateConnection: (inputCoreId: number, outputCoreId: number, loss: number) => void;
    onDeleteConnection: (connectionId: number) => void;
    isLoading?: boolean;
}

export function SpliceMatrix({
    nodeId,
    nodeName,
    inputCores,
    outputCores,
    connections,
    inputCableId,
    outputCableId,
    onCreateConnection,
    onDeleteConnection,
    isLoading,
}: SpliceMatrixProps) {
    const [selectedInput, setSelectedInput] = useState<Core | null>(null);
    const [selectedOutput, setSelectedOutput] = useState<Core | null>(null);
    const [spliceLoss, setSpliceLoss] = useState('0.1');
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [deleteConnection, setDeleteConnection] = useState<Connection | null>(null);

    const getCoreColor = (core: Core) => {
        return core.core_color || FIBER_COLORS[(core.core_index - 1) % 12];
    };

    const isInputConnected = (coreId: number) => {
        return connections.some((c) => c.input_core_id === coreId);
    };

    const isOutputConnected = (coreId: number) => {
        return connections.some((c) => c.output_core_id === coreId);
    };

    const getConnectionForInput = (coreId: number) => {
        return connections.find((c) => c.input_core_id === coreId);
    };

    const handleInputClick = (core: Core) => {
        if (isInputConnected(core.id)) return;
        setSelectedInput(core);
        if (selectedOutput) {
            setIsModalOpen(true);
        }
    };

    const handleOutputClick = (core: Core) => {
        if (isOutputConnected(core.id)) return;
        setSelectedOutput(core);
        if (selectedInput) {
            setIsModalOpen(true);
        }
    };

    const handleCreateConnection = () => {
        if (selectedInput && selectedOutput) {
            onCreateConnection(selectedInput.id, selectedOutput.id, parseFloat(spliceLoss) || 0.1);
            setSelectedInput(null);
            setSelectedOutput(null);
            setSpliceLoss('0.1');
            setIsModalOpen(false);
        }
    };

    const cancelSelection = () => {
        setSelectedInput(null);
        setSelectedOutput(null);
        setIsModalOpen(false);
    };

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="flex items-center justify-between">
                <div>
                    <h3 className="text-lg font-semibold text-white flex items-center gap-2">
                        <Zap className="h-5 w-5 text-amber-400" />
                        Splice Matrix at {nodeName}
                    </h3>
                    <p className="text-sm text-slate-400">
                        Connect fiber cores between cables. Click an input core, then an output core.
                    </p>
                </div>
                <Badge variant="info">
                    {connections.length} connection{connections.length !== 1 ? 's' : ''}
                </Badge>
            </div>

            {/* Legend */}
            <div className="flex flex-wrap gap-4 text-xs p-3 rounded-lg bg-slate-800/50">
                <div className="flex items-center gap-1.5">
                    <div className="w-3 h-3 rounded-full bg-slate-500 opacity-50" />
                    <span className="text-slate-400">Vacant</span>
                </div>
                <div className="flex items-center gap-1.5">
                    <div className="w-3 h-3 rounded-full bg-green-500 ring-2 ring-green-500" />
                    <span className="text-slate-400">Connected</span>
                </div>
                <div className="flex items-center gap-1.5">
                    <div className="w-3 h-3 rounded-full bg-sky-500 ring-2 ring-sky-500 animate-pulse" />
                    <span className="text-slate-400">Selected</span>
                </div>
            </div>

            {/* Matrix View */}
            <div className="flex gap-8">
                {/* Input Cores (Left) */}
                <div className="flex-1">
                    <div className="text-sm font-medium text-slate-300 mb-3">
                        Input Cable (Source)
                    </div>
                    <div className="space-y-1.5">
                        {inputCores.map((core) => {
                            const connected = isInputConnected(core.id);
                            const connection = getConnectionForInput(core.id);
                            const isSelected = selectedInput?.id === core.id;
                            const coreColor = getCoreColor(core);

                            return (
                                <div
                                    key={core.id}
                                    className={cn(
                                        'flex items-center gap-3 p-2 rounded-lg transition-all',
                                        connected
                                            ? 'bg-green-500/10 border border-green-500/30'
                                            : isSelected
                                                ? 'bg-sky-500/20 border border-sky-500'
                                                : 'bg-slate-800/50 border border-slate-700/50 cursor-pointer hover:border-slate-600'
                                    )}
                                    onClick={() => !connected && handleInputClick(core)}
                                >
                                    <div
                                        className={cn(
                                            'w-6 h-6 rounded-full flex items-center justify-center text-[10px] font-bold',
                                            connected && 'ring-2 ring-green-500',
                                            isSelected && 'ring-2 ring-sky-500 animate-pulse'
                                        )}
                                        style={{ backgroundColor: COLOR_HEX[coreColor] }}
                                    >
                                        {core.core_index}
                                    </div>
                                    <span className="text-sm text-white flex-1">
                                        Core {core.core_index}
                                    </span>
                                    {connected && connection && (
                                        <div className="flex items-center gap-2">
                                            <ArrowRight className="h-4 w-4 text-green-400" />
                                            <span className="text-xs text-slate-400">
                                                → Core {connection.output_core_index}
                                            </span>
                                            <Badge variant="outline" className="text-xs">
                                                {connection.splice_loss?.toFixed(2) || '0.00'} dB
                                            </Badge>
                                            <Button
                                                size="icon"
                                                variant="ghost"
                                                className="h-6 w-6"
                                                onClick={(e) => {
                                                    e.stopPropagation();
                                                    setDeleteConnection(connection);
                                                }}
                                            >
                                                <Trash2 className="h-3 w-3 text-red-400" />
                                            </Button>
                                        </div>
                                    )}
                                </div>
                            );
                        })}
                    </div>
                </div>

                {/* Connection Lines (Center) */}
                <div className="w-16 flex items-center justify-center">
                    <div className="h-full w-px bg-gradient-to-b from-slate-700 via-amber-500 to-slate-700" />
                </div>

                {/* Output Cores (Right) */}
                <div className="flex-1">
                    <div className="text-sm font-medium text-slate-300 mb-3">
                        Output Cable (Destination)
                    </div>
                    <div className="space-y-1.5">
                        {outputCores.map((core) => {
                            const connected = isOutputConnected(core.id);
                            const isSelected = selectedOutput?.id === core.id;
                            const coreColor = getCoreColor(core);

                            return (
                                <div
                                    key={core.id}
                                    className={cn(
                                        'flex items-center gap-3 p-2 rounded-lg transition-all',
                                        connected
                                            ? 'bg-green-500/10 border border-green-500/30'
                                            : isSelected
                                                ? 'bg-sky-500/20 border border-sky-500'
                                                : 'bg-slate-800/50 border border-slate-700/50 cursor-pointer hover:border-slate-600'
                                    )}
                                    onClick={() => !connected && handleOutputClick(core)}
                                >
                                    <div
                                        className={cn(
                                            'w-6 h-6 rounded-full flex items-center justify-center text-[10px] font-bold',
                                            connected && 'ring-2 ring-green-500',
                                            isSelected && 'ring-2 ring-sky-500 animate-pulse'
                                        )}
                                        style={{ backgroundColor: COLOR_HEX[coreColor] }}
                                    >
                                        {core.core_index}
                                    </div>
                                    <span className="text-sm text-white">Core {core.core_index}</span>
                                    {connected && (
                                        <Badge variant="success" className="ml-auto">Connected</Badge>
                                    )}
                                </div>
                            );
                        })}
                    </div>
                </div>
            </div>

            {/* Selection indicator */}
            {(selectedInput || selectedOutput) && (
                <div className="p-3 rounded-lg bg-sky-500/10 border border-sky-500/30">
                    <div className="flex items-center justify-between">
                        <div className="flex items-center gap-4 text-sm">
                            {selectedInput ? (
                                <span className="text-white">
                                    Input: <strong>Core {selectedInput.core_index}</strong>
                                </span>
                            ) : (
                                <span className="text-slate-400">Select input core...</span>
                            )}
                            <ArrowRight className="h-4 w-4 text-slate-400" />
                            {selectedOutput ? (
                                <span className="text-white">
                                    Output: <strong>Core {selectedOutput.core_index}</strong>
                                </span>
                            ) : (
                                <span className="text-slate-400">Select output core...</span>
                            )}
                        </div>
                        <Button variant="outline" size="sm" onClick={cancelSelection}>
                            Cancel
                        </Button>
                    </div>
                </div>
            )}

            {/* Create Connection Modal */}
            <Modal
                isOpen={isModalOpen}
                onClose={cancelSelection}
                title="Create Splice Connection"
                size="sm"
            >
                <div className="space-y-4">
                    <div className="flex items-center justify-between p-3 rounded-lg bg-slate-800/50">
                        <div className="text-center">
                            <div
                                className="w-10 h-10 rounded-full mx-auto flex items-center justify-center text-sm font-bold"
                                style={{ backgroundColor: selectedInput ? COLOR_HEX[getCoreColor(selectedInput)] : '#64748b' }}
                            >
                                {selectedInput?.core_index || '?'}
                            </div>
                            <p className="text-xs text-slate-400 mt-1">Input</p>
                        </div>
                        <ArrowRight className="h-6 w-6 text-amber-400" />
                        <div className="text-center">
                            <div
                                className="w-10 h-10 rounded-full mx-auto flex items-center justify-center text-sm font-bold"
                                style={{ backgroundColor: selectedOutput ? COLOR_HEX[getCoreColor(selectedOutput)] : '#64748b' }}
                            >
                                {selectedOutput?.core_index || '?'}
                            </div>
                            <p className="text-xs text-slate-400 mt-1">Output</p>
                        </div>
                    </div>

                    <Input
                        label="Splice Loss (dB)"
                        type="number"
                        step="0.01"
                        value={spliceLoss}
                        onChange={(e) => setSpliceLoss(e.target.value)}
                        helperText="Typical range: 0.05 - 0.3 dB"
                    />

                    <div className="flex justify-end gap-3 pt-4 border-t border-slate-700/50">
                        <Button variant="outline" onClick={cancelSelection}>
                            Cancel
                        </Button>
                        <Button onClick={handleCreateConnection} isLoading={isLoading}>
                            <Plus className="h-4 w-4 mr-2" />
                            Create Connection
                        </Button>
                    </div>
                </div>
            </Modal>

            {/* Delete Confirmation */}
            <ConfirmDialog
                isOpen={!!deleteConnection}
                onClose={() => setDeleteConnection(null)}
                onConfirm={() => {
                    if (deleteConnection) {
                        onDeleteConnection(deleteConnection.id);
                        setDeleteConnection(null);
                    }
                }}
                title="Delete Connection"
                message={`Delete connection from Core ${deleteConnection?.input_core_index} to Core ${deleteConnection?.output_core_index}?`}
                confirmText="Delete"
                variant="danger"
                isLoading={isLoading}
            />
        </div>
    );
}
