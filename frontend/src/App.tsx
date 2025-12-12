import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { Layout } from '@/components/layout';
import { ErrorBoundary } from '@/components/ErrorBoundary';
import { ToastProvider } from '@/components/ui';
import { Dashboard, MapPage, NodesPage, CablesPage, CustomersPage, SplicingPage } from '@/pages';

// Create query client
const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            staleTime: 30000, // 30 seconds
            retry: 2,
            refetchOnWindowFocus: false,
        },
    },
});

// Placeholder pages for routes not yet implemented

function SettingsPage() {
    return (
        <div className="p-6">
            <h1 className="text-3xl font-bold text-white mb-2">Settings</h1>
            <p className="text-slate-400">Application settings coming soon...</p>
        </div>
    );
}

function App() {
    return (
        <ErrorBoundary>
            <QueryClientProvider client={queryClient}>
                <ToastProvider>
                    <BrowserRouter>
                        <Routes>
                            <Route path="/" element={<Layout />}>
                                <Route index element={<Dashboard />} />
                                <Route path="map" element={<MapPage />} />
                                <Route path="nodes" element={<NodesPage />} />
                                <Route path="cables" element={<CablesPage />} />
                                <Route path="splicing" element={<SplicingPage />} />
                                <Route path="customers" element={<CustomersPage />} />
                                <Route path="settings" element={<SettingsPage />} />
                            </Route>
                        </Routes>
                    </BrowserRouter>
                </ToastProvider>
            </QueryClientProvider>
        </ErrorBoundary>
    );
}

export default App;

