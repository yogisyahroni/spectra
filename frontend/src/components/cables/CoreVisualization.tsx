import React from 'react';
import { cn } from '@/lib/utils';

// Standard fiber optic color coding (TIA-598)
const TUBE_COLORS = ['Blue', 'Orange', 'Green', 'Brown', 'Slate', 'White', 'Red', 'Black', 'Yellow', 'Violet', 'Rose', 'Aqua'];
const FIBER_COLORS = ['Blue', 'Orange', 'Green', 'Brown', 'Slate', 'White', 'Red', 'Black', 'Yellow', 'Violet', 'Rose', 'Aqua'];

// Color hex values for display
const COLOR_HEX: Record<string, string> = {
    Blue: '#3b82f6',
    Orange: '#f97316',
    Green: '#22c55e',
    Brown: '#92400e',
    Slate: '#64748b',
    White: '#f1f5f9',
    Red: '#ef4444',
    Black: '#1e293b',
    Yellow: '#eab308',
    Violet: '#8b5cf6',
    Rose: '#f43f5e',
    Aqua: '#06b6d4',
};

interface CoreVisualizationProps {
    cores: {
        id: number;
        core_index: number;
        tube_color?: string;
        core_color?: string;
        status: string;
    }[];
    coreCount: number;
    compact?: boolean;
    onCoreClick?: (coreIndex: number) => void;
}

export function CoreVisualization({
    cores,
    coreCount,
    compact = false,
    onCoreClick
}: CoreVisualizationProps) {
    // Group cores by tube (12 cores per tube)
    const coresPerTube = 12;
    const tubeCount = Math.ceil(coreCount / coresPerTube);

    const getCoreStatus = (index: number) => {
        const core = cores.find((c) => c.core_index === index);
        return core?.status || 'VACANT';
    };

    const getCoreColor = (index: number) => {
        const core = cores.find((c) => c.core_index === index);
        return core?.core_color || FIBER_COLORS[(index - 1) % 12];
    };

    const getTubeColor = (tubeIndex: number) => {
        return TUBE_COLORS[tubeIndex % 12];
    };

    const getStatusStyle = (status: string) => {
        switch (status) {
            case 'USED':
                return 'ring-2 ring-green-500';
            case 'RESERVED':
                return 'ring-2 ring-yellow-500';
            case 'DAMAGED':
                return 'ring-2 ring-red-500 opacity-50';
            default:
                return 'opacity-70';
        }
    };

    if (compact) {
        // Compact view - just show summary
        const usedCount = cores.filter((c) => c.status === 'USED').length;
        const vacantCount = cores.filter((c) => c.status === 'VACANT').length;
        const damagedCount = cores.filter((c) => c.status === 'DAMAGED').length;

        return (
            <div className="flex items-center gap-4 text-sm">
                <div className="flex items-center gap-1.5">
                    <div className="w-2 h-2 rounded-full bg-green-500" />
                    <span className="text-slate-300">{usedCount} used</span>
                </div>
                <div className="flex items-center gap-1.5">
                    <div className="w-2 h-2 rounded-full bg-slate-500" />
                    <span className="text-slate-300">{vacantCount} vacant</span>
                </div>
                {damagedCount > 0 && (
                    <div className="flex items-center gap-1.5">
                        <div className="w-2 h-2 rounded-full bg-red-500" />
                        <span className="text-slate-300">{damagedCount} damaged</span>
                    </div>
                )}
            </div>
        );
    }

    return (
        <div className="space-y-4">
            {/* Legend */}
            <div className="flex flex-wrap gap-4 text-xs">
                <div className="flex items-center gap-1.5">
                    <div className="w-3 h-3 rounded-full bg-slate-500 opacity-70" />
                    <span className="text-slate-400">Vacant</span>
                </div>
                <div className="flex items-center gap-1.5">
                    <div className="w-3 h-3 rounded-full bg-green-500 ring-2 ring-green-500" />
                    <span className="text-slate-400">Used</span>
                </div>
                <div className="flex items-center gap-1.5">
                    <div className="w-3 h-3 rounded-full bg-yellow-500 ring-2 ring-yellow-500" />
                    <span className="text-slate-400">Reserved</span>
                </div>
                <div className="flex items-center gap-1.5">
                    <div className="w-3 h-3 rounded-full bg-red-500 opacity-50" />
                    <span className="text-slate-400">Damaged</span>
                </div>
            </div>

            {/* Tubes */}
            <div className="space-y-3">
                {Array.from({ length: tubeCount }).map((_, tubeIndex) => {
                    const tubeColor = getTubeColor(tubeIndex);
                    const startCore = tubeIndex * coresPerTube + 1;
                    const endCore = Math.min((tubeIndex + 1) * coresPerTube, coreCount);

                    return (
                        <div
                            key={tubeIndex}
                            className="p-3 rounded-lg border"
                            style={{
                                borderColor: COLOR_HEX[tubeColor] + '50',
                                backgroundColor: COLOR_HEX[tubeColor] + '10'
                            }}
                        >
                            <div className="flex items-center gap-2 mb-2">
                                <div
                                    className="w-4 h-4 rounded"
                                    style={{ backgroundColor: COLOR_HEX[tubeColor] }}
                                />
                                <span className="text-sm font-medium text-white">
                                    Tube {tubeIndex + 1} ({tubeColor})
                                </span>
                                <span className="text-xs text-slate-400">
                                    Core {startCore}-{endCore}
                                </span>
                            </div>

                            {/* Cores in this tube */}
                            <div className="flex flex-wrap gap-1.5">
                                {Array.from({ length: endCore - startCore + 1 }).map((_, i) => {
                                    const coreIndex = startCore + i;
                                    const coreColor = getCoreColor(coreIndex);
                                    const status = getCoreStatus(coreIndex);

                                    return (
                                        <button
                                            key={coreIndex}
                                            className={cn(
                                                'w-6 h-6 rounded-full flex items-center justify-center text-[10px] font-bold transition-all',
                                                getStatusStyle(status),
                                                onCoreClick && 'cursor-pointer hover:scale-110'
                                            )}
                                            style={{ backgroundColor: COLOR_HEX[coreColor] }}
                                            onClick={() => onCoreClick?.(coreIndex)}
                                            title={`Core ${coreIndex} (${coreColor}) - ${status}`}
                                        >
                                            {coreIndex}
                                        </button>
                                    );
                                })}
                            </div>
                        </div>
                    );
                })}
            </div>
        </div>
    );
}

export { TUBE_COLORS, FIBER_COLORS, COLOR_HEX };
