import * as React from 'react';
import { ChevronUp, ChevronDown, Search } from 'lucide-react';
import { cn } from '@/lib/utils';
import { Input } from './Input';

export interface Column<T> {
    key: keyof T | string;
    header: string;
    sortable?: boolean;
    render?: (item: T) => React.ReactNode;
    width?: string;
}

export interface DataTableProps<T> {
    data: T[];
    columns: Column<T>[];
    keyField: keyof T;
    isLoading?: boolean;
    searchable?: boolean;
    searchPlaceholder?: string;
    onRowClick?: (item: T) => void;
    emptyMessage?: string;
}

export function DataTable<T extends object>({
    data,
    columns,
    keyField,
    isLoading = false,
    searchable = false,
    searchPlaceholder = 'Search...',
    onRowClick,
    emptyMessage = 'No data found',
}: DataTableProps<T>) {
    const [sortKey, setSortKey] = React.useState<string | null>(null);
    const [sortOrder, setSortOrder] = React.useState<'asc' | 'desc'>('asc');
    const [searchTerm, setSearchTerm] = React.useState('');

    // Handle sorting
    const handleSort = (key: string) => {
        if (sortKey === key) {
            setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
        } else {
            setSortKey(key);
            setSortOrder('asc');
        }
    };

    // Filter and sort data
    const processedData = React.useMemo(() => {
        let result = [...data];

        // Search filter
        if (searchTerm) {
            const lower = searchTerm.toLowerCase();
            result = result.filter((item) =>
                Object.values(item).some((value) =>
                    String(value).toLowerCase().includes(lower)
                )
            );
        }

        // Sort
        if (sortKey) {
            result.sort((a, b) => {
                const aVal = a[sortKey as keyof T];
                const bVal = b[sortKey as keyof T];

                if (aVal === bVal) return 0;
                if (aVal === null || aVal === undefined) return 1;
                if (bVal === null || bVal === undefined) return -1;

                const comparison = aVal < bVal ? -1 : 1;
                return sortOrder === 'asc' ? comparison : -comparison;
            });
        }

        return result;
    }, [data, searchTerm, sortKey, sortOrder]);

    return (
        <div className="space-y-4">
            {/* Search */}
            {searchable && (
                <div className="relative max-w-sm">
                    <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-400" />
                    <Input
                        placeholder={searchPlaceholder}
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                        className="pl-10"
                    />
                </div>
            )}

            {/* Table */}
            <div className="overflow-hidden rounded-xl border border-slate-700/50">
                <div className="overflow-x-auto">
                    <table className="w-full">
                        <thead>
                            <tr className="bg-slate-800/50 border-b border-slate-700/50">
                                {columns.map((column) => (
                                    <th
                                        key={String(column.key)}
                                        className={cn(
                                            'px-4 py-3 text-left text-sm font-medium text-slate-300',
                                            column.sortable && 'cursor-pointer hover:text-white select-none',
                                            column.width
                                        )}
                                        style={column.width ? { width: column.width } : undefined}
                                        onClick={() => column.sortable && handleSort(String(column.key))}
                                    >
                                        <div className="flex items-center gap-2">
                                            {column.header}
                                            {column.sortable && sortKey === column.key && (
                                                sortOrder === 'asc' ? (
                                                    <ChevronUp className="h-4 w-4" />
                                                ) : (
                                                    <ChevronDown className="h-4 w-4" />
                                                )
                                            )}
                                        </div>
                                    </th>
                                ))}
                            </tr>
                        </thead>
                        <tbody>
                            {isLoading ? (
                                // Loading skeleton
                                Array.from({ length: 5 }).map((_, i) => (
                                    <tr key={i} className="border-b border-slate-700/30">
                                        {columns.map((column, j) => (
                                            <td key={j} className="px-4 py-3">
                                                <div className="h-5 skeleton rounded w-3/4" />
                                            </td>
                                        ))}
                                    </tr>
                                ))
                            ) : processedData.length === 0 ? (
                                <tr>
                                    <td
                                        colSpan={columns.length}
                                        className="px-4 py-12 text-center text-slate-400"
                                    >
                                        {emptyMessage}
                                    </td>
                                </tr>
                            ) : (
                                processedData.map((item) => (
                                    <tr
                                        key={String(item[keyField])}
                                        className={cn(
                                            'border-b border-slate-700/30 transition-colors',
                                            onRowClick && 'cursor-pointer hover:bg-slate-800/50'
                                        )}
                                        onClick={() => onRowClick?.(item)}
                                    >
                                        {columns.map((column) => (
                                            <td
                                                key={String(column.key)}
                                                className="px-4 py-3 text-sm text-white"
                                            >
                                                {column.render
                                                    ? column.render(item)
                                                    : String(item[column.key as keyof T] ?? '-')}
                                            </td>
                                        ))}
                                    </tr>
                                ))
                            )}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    );
}
