import React, { Component, ErrorInfo, ReactNode } from 'react';
import { AlertCircle, RefreshCw } from 'lucide-react';
import { Button } from '@/components/ui';

interface Props {
    children: ReactNode;
}

interface State {
    hasError: boolean;
    error: Error | null;
}

export class ErrorBoundary extends Component<Props, State> {
    public state: State = {
        hasError: false,
        error: null,
    };

    public static getDerivedStateFromError(error: Error): State {
        return { hasError: true, error };
    }

    public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
        console.error('ErrorBoundary caught an error:', error, errorInfo);
    }

    private handleRetry = () => {
        this.setState({ hasError: false, error: null });
        window.location.reload();
    };

    public render() {
        if (this.state.hasError) {
            return (
                <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-950 via-slate-900 to-slate-950">
                    <div className="text-center max-w-md p-8 glass-card rounded-2xl">
                        <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-red-500/20 mb-6">
                            <AlertCircle className="w-8 h-8 text-red-400" />
                        </div>
                        <h1 className="text-2xl font-bold text-white mb-2">Something went wrong</h1>
                        <p className="text-slate-400 mb-6">
                            An unexpected error occurred. Please try refreshing the page.
                        </p>
                        {this.state.error && (
                            <pre className="text-left text-xs text-red-400 bg-slate-900 p-4 rounded-lg mb-6 overflow-auto max-h-40">
                                {this.state.error.message}
                            </pre>
                        )}
                        <Button onClick={this.handleRetry}>
                            <RefreshCw className="w-4 h-4 mr-2" />
                            Refresh Page
                        </Button>
                    </div>
                </div>
            );
        }

        return this.props.children;
    }
}
