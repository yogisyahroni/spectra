import React from 'react';
import { useQuery } from '@tanstack/react-query';
import {
    Network,
    Cable,
    Users,
    AlertTriangle,
    CheckCircle,
    MapPin,
    Activity,
} from 'lucide-react';
import { Card, CardHeader, CardTitle, CardContent, Badge } from '@/components/ui';
import { api } from '@/services/api';
import { cn } from '@/lib/utils';

export function Dashboard() {
    // Fetch data for dashboard
    const { data: nodesData, isLoading: nodesLoading } = useQuery({
        queryKey: ['nodes'],
        queryFn: () => api.getNodes({ limit: 1000 }),
    });

    const { data: cablesData, isLoading: cablesLoading } = useQuery({
        queryKey: ['cables'],
        queryFn: () => api.getCables({ limit: 1000 }),
    });

    const { data: customersData, isLoading: customersLoading } = useQuery({
        queryKey: ['customers'],
        queryFn: () => api.getCustomers({ limit: 1000 }),
    });

    const { data: losCustomers } = useQuery({
        queryKey: ['customers-los'],
        queryFn: () => api.getLOSCustomers(),
    });

    const nodes = nodesData?.data || [];
    const cables = cablesData?.data || [];
    const customers = customersData?.data || [];
    const losCount = losCustomers?.length || 0;

    // Calculate statistics
    const nodesByType = nodes.reduce((acc, node) => {
        acc[node.type] = (acc[node.type] || 0) + 1;
        return acc;
    }, {} as Record<string, number>);

    const onlineCustomers = customers.filter((c) => c.current_status === 'ONLINE').length;
    const offlineCustomers = customers.filter((c) => c.current_status === 'OFFLINE').length;

    const statCards = [
        {
            title: 'Total Nodes',
            value: nodes.length,
            icon: Network,
            color: 'from-sky-500 to-blue-600',
            subtext: `${nodesByType['OLT'] || 0} OLT, ${nodesByType['ODP'] || 0} ODP`,
        },
        {
            title: 'Total Cables',
            value: cables.length,
            icon: Cable,
            color: 'from-violet-500 to-purple-600',
            subtext: `${cables.reduce((sum, c) => sum + c.core_count, 0)} cores total`,
        },
        {
            title: 'Customers',
            value: customers.length,
            icon: Users,
            color: 'from-emerald-500 to-green-600',
            subtext: `${onlineCustomers} online, ${offlineCustomers} offline`,
        },
        {
            title: 'LOS Alerts',
            value: losCount,
            icon: AlertTriangle,
            color: losCount > 0 ? 'from-red-500 to-rose-600' : 'from-slate-500 to-slate-600',
            subtext: losCount > 0 ? 'Requires attention' : 'All systems normal',
        },
    ];

    return (
        <div className="p-6">
            {/* Header */}
            <div className="mb-8">
                <h1 className="text-3xl font-bold text-white mb-2">Dashboard</h1>
                <p className="text-slate-400">SPECTRA Network Management Overview</p>
            </div>

            {/* Stats Grid */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
                {statCards.map((stat, index) => (
                    <Card key={index} variant="glass" className="relative overflow-hidden">
                        <div className={cn('absolute inset-0 opacity-10 bg-gradient-to-br', stat.color)} />
                        <CardContent className="relative p-6">
                            <div className="flex items-center justify-between mb-4">
                                <div className={cn('p-3 rounded-xl bg-gradient-to-br', stat.color)}>
                                    <stat.icon className="h-6 w-6 text-white" />
                                </div>
                                {nodesLoading || cablesLoading || customersLoading ? (
                                    <div className="h-8 w-16 skeleton rounded" />
                                ) : (
                                    <span className="text-3xl font-bold text-white">{stat.value}</span>
                                )}
                            </div>
                            <h3 className="text-slate-400 text-sm font-medium">{stat.title}</h3>
                            <p className="text-slate-500 text-xs mt-1">{stat.subtext}</p>
                        </CardContent>
                    </Card>
                ))}
            </div>

            {/* Node Types Breakdown */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
                <Card variant="glass">
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2">
                            <MapPin className="h-5 w-5 text-sky-400" />
                            Node Distribution
                        </CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="space-y-4">
                            {['OLT', 'ODC', 'ODP', 'CLOSURE', 'POLE', 'CUSTOMER'].map((type) => {
                                const count = nodesByType[type] || 0;
                                const percentage = nodes.length > 0 ? (count / nodes.length) * 100 : 0;
                                return (
                                    <div key={type} className="flex items-center gap-4">
                                        <span className="w-20 text-sm text-slate-400">{type}</span>
                                        <div className="flex-1 h-2 bg-slate-700 rounded-full overflow-hidden">
                                            <div
                                                className="h-full bg-gradient-to-r from-sky-500 to-blue-500 rounded-full transition-all duration-500"
                                                style={{ width: `${percentage}%` }}
                                            />
                                        </div>
                                        <span className="w-12 text-sm text-white text-right">{count}</span>
                                    </div>
                                );
                            })}
                        </div>
                    </CardContent>
                </Card>

                <Card variant="glass">
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2">
                            <Activity className="h-5 w-5 text-emerald-400" />
                            System Status
                        </CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="space-y-4">
                            <div className="flex items-center justify-between p-4 bg-slate-800/50 rounded-xl">
                                <div className="flex items-center gap-3">
                                    <CheckCircle className="h-5 w-5 text-emerald-400" />
                                    <span className="text-white">Backend API</span>
                                </div>
                                <Badge variant="success">Online</Badge>
                            </div>
                            <div className="flex items-center justify-between p-4 bg-slate-800/50 rounded-xl">
                                <div className="flex items-center gap-3">
                                    <CheckCircle className="h-5 w-5 text-emerald-400" />
                                    <span className="text-white">Database</span>
                                </div>
                                <Badge variant="success">Connected</Badge>
                            </div>
                            <div className="flex items-center justify-between p-4 bg-slate-800/50 rounded-xl">
                                <div className="flex items-center gap-3">
                                    {losCount > 0 ? (
                                        <AlertTriangle className="h-5 w-5 text-amber-400" />
                                    ) : (
                                        <CheckCircle className="h-5 w-5 text-emerald-400" />
                                    )}
                                    <span className="text-white">Network Status</span>
                                </div>
                                <Badge variant={losCount > 0 ? 'warning' : 'success'}>
                                    {losCount > 0 ? `${losCount} Alerts` : 'Normal'}
                                </Badge>
                            </div>
                        </div>
                    </CardContent>
                </Card>
            </div>
        </div>
    );
}
