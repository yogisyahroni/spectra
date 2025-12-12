import * as React from 'react';
import { cn } from '@/lib/utils';

export interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
    variant?: 'default' | 'secondary' | 'outline' | 'ghost' | 'destructive';
    size?: 'sm' | 'md' | 'lg' | 'icon';
    isLoading?: boolean;
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
    ({ className, variant = 'default', size = 'md', isLoading, children, disabled, ...props }, ref) => {
        const baseStyles =
            'inline-flex items-center justify-center rounded-xl font-medium transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-slate-900 disabled:opacity-50 disabled:cursor-not-allowed btn-glow';

        const variants = {
            default:
                'bg-gradient-to-r from-sky-500 to-blue-600 text-white hover:from-sky-400 hover:to-blue-500 focus:ring-sky-500 shadow-lg shadow-sky-500/25',
            secondary:
                'bg-gradient-to-r from-violet-500 to-purple-600 text-white hover:from-violet-400 hover:to-purple-500 focus:ring-violet-500 shadow-lg shadow-violet-500/25',
            outline:
                'border-2 border-slate-600 bg-transparent text-slate-200 hover:bg-slate-800 hover:border-slate-500 focus:ring-slate-500',
            ghost: 'bg-transparent text-slate-300 hover:bg-slate-800 hover:text-white focus:ring-slate-500',
            destructive:
                'bg-gradient-to-r from-red-500 to-rose-600 text-white hover:from-red-400 hover:to-rose-500 focus:ring-red-500 shadow-lg shadow-red-500/25',
        };

        const sizes = {
            sm: 'h-8 px-3 text-sm gap-1.5',
            md: 'h-10 px-4 text-sm gap-2',
            lg: 'h-12 px-6 text-base gap-2.5',
            icon: 'h-10 w-10',
        };

        return (
            <button
                className={cn(baseStyles, variants[variant], sizes[size], className)}
                ref={ref}
                disabled={disabled || isLoading}
                {...props}
            >
                {isLoading && (
                    <svg
                        className="animate-spin -ml-1 mr-2 h-4 w-4"
                        xmlns="http://www.w3.org/2000/svg"
                        fill="none"
                        viewBox="0 0 24 24"
                    >
                        <circle
                            className="opacity-25"
                            cx="12"
                            cy="12"
                            r="10"
                            stroke="currentColor"
                            strokeWidth="4"
                        ></circle>
                        <path
                            className="opacity-75"
                            fill="currentColor"
                            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                        ></path>
                    </svg>
                )}
                {children}
            </button>
        );
    }
);

Button.displayName = 'Button';

export { Button };
