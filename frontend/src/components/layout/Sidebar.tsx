import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import {
    Map,
    Network,
    Users,
    Cable,
    LayoutDashboard,
    Settings,
    ChevronLeft,
    ChevronRight,
} from 'lucide-react';
import { cn } from '@/lib/utils';

interface SidebarProps {
    collapsed: boolean;
    onToggle: () => void;
}

const navItems = [
    { icon: LayoutDashboard, label: 'Dashboard', path: '/' },
    { icon: Map, label: 'Map', path: '/map' },
    { icon: Network, label: 'Nodes', path: '/nodes' },
    { icon: Cable, label: 'Cables', path: '/cables' },
    { icon: Users, label: 'Customers', path: '/customers' },
    { icon: Settings, label: 'Settings', path: '/settings' },
];

export function Sidebar({ collapsed, onToggle }: SidebarProps) {
    const location = useLocation();

    return (
        <aside
            className={cn(
                'fixed left-0 top-0 z-40 h-screen transition-all duration-300 ease-in-out',
                'bg-slate-900/95 backdrop-blur-xl border-r border-slate-700/50',
                collapsed ? 'w-16' : 'w-64'
            )}
        >
            {/* Logo */}
            <div className="flex h-16 items-center justify-between px-4 border-b border-slate-700/50">
                {!collapsed && (
                    <div className="flex items-center gap-2">
                        <div className="h-8 w-8 rounded-lg bg-gradient-to-br from-sky-400 to-blue-600 flex items-center justify-center">
                            <Network className="h-5 w-5 text-white" />
                        </div>
                        <span className="text-lg font-bold gradient-text">SPECTRA</span>
                    </div>
                )}
                {collapsed && (
                    <div className="h-8 w-8 mx-auto rounded-lg bg-gradient-to-br from-sky-400 to-blue-600 flex items-center justify-center">
                        <Network className="h-5 w-5 text-white" />
                    </div>
                )}
            </div>

            {/* Navigation */}
            <nav className="flex flex-col gap-1 p-3">
                {navItems.map((item) => {
                    const isActive = location.pathname === item.path;
                    return (
                        <Link
                            key={item.path}
                            to={item.path}
                            className={cn(
                                'flex items-center gap-3 rounded-xl px-3 py-2.5 transition-all duration-200',
                                'hover:bg-slate-800/80',
                                isActive
                                    ? 'bg-gradient-to-r from-sky-500/20 to-blue-500/20 text-sky-400 border border-sky-500/30'
                                    : 'text-slate-400 hover:text-white'
                            )}
                        >
                            <item.icon className={cn('h-5 w-5 flex-shrink-0', isActive && 'text-sky-400')} />
                            {!collapsed && <span className="font-medium">{item.label}</span>}
                        </Link>
                    );
                })}
            </nav>

            {/* Toggle button */}
            <button
                onClick={onToggle}
                className={cn(
                    'absolute -right-3 top-20 z-50',
                    'flex h-6 w-6 items-center justify-center rounded-full',
                    'bg-slate-800 border border-slate-600 text-slate-400',
                    'hover:bg-slate-700 hover:text-white transition-colors'
                )}
            >
                {collapsed ? <ChevronRight className="h-4 w-4" /> : <ChevronLeft className="h-4 w-4" />}
            </button>
        </aside>
    );
}
