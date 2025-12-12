import * as React from 'react';
import { cn } from '@/lib/utils';

export interface BadgeProps extends React.HTMLAttributes<HTMLSpanElement> {
    variant?: 'default' | 'success' | 'warning' | 'error' | 'info' | 'outline';
    size?: 'sm' | 'md';
}

const Badge = React.forwardRef<HTMLSpanElement, BadgeProps>(
    ({ className, variant = 'default', size = 'md', ...props }, ref) => {
        const variants = {
            default: 'bg-slate-700 text-slate-200',
            success: 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/30',
            warning: 'bg-amber-500/20 text-amber-400 border border-amber-500/30',
            error: 'bg-red-500/20 text-red-400 border border-red-500/30',
            info: 'bg-sky-500/20 text-sky-400 border border-sky-500/30',
            outline: 'border border-slate-600 text-slate-300 bg-transparent',
        };

        const sizes = {
            sm: 'px-2 py-0.5 text-xs',
            md: 'px-2.5 py-1 text-xs',
        };

        return (
            <span
                ref={ref}
                className={cn(
                    'inline-flex items-center rounded-full font-medium',
                    variants[variant],
                    sizes[size],
                    className
                )}
                {...props}
            />
        );
    }
);

Badge.displayName = 'Badge';

export { Badge };
