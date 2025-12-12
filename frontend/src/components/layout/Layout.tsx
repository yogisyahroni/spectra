import React, { useState } from 'react';
import { Outlet } from 'react-router-dom';
import { Sidebar } from './Sidebar';
import { cn } from '@/lib/utils';

export function Layout() {
    const [sidebarCollapsed, setSidebarCollapsed] = useState(false);

    return (
        <div className="min-h-screen bg-gradient-to-br from-slate-950 via-slate-900 to-slate-950">
            <Sidebar collapsed={sidebarCollapsed} onToggle={() => setSidebarCollapsed(!sidebarCollapsed)} />

            <main
                className={cn(
                    'min-h-screen transition-all duration-300 ease-in-out',
                    sidebarCollapsed ? 'ml-16' : 'ml-64'
                )}
            >
                <Outlet />
            </main>
        </div>
    );
}
