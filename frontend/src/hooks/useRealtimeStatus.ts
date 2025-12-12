import { useEffect, useCallback, useRef } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { api } from '@/services/api';
import type { Customer, CustomerStatus } from '@/types';

interface StatusChange {
    customerId: number;
    customerName: string;
    oldStatus: CustomerStatus;
    newStatus: CustomerStatus;
    timestamp: Date;
}

interface UseRealtimeStatusOptions {
    enabled?: boolean;
    pollingInterval?: number; // in milliseconds
    onStatusChange?: (change: StatusChange) => void;
}

/**
 * Hook for real-time customer status monitoring using polling.
 * In production, this would use WebSockets for true real-time updates.
 */
export function useRealtimeStatus({
    enabled = true,
    pollingInterval = 10000, // 10 seconds default
    onStatusChange,
}: UseRealtimeStatusOptions = {}) {
    const queryClient = useQueryClient();
    const previousStatusRef = useRef<Map<number, CustomerStatus>>(new Map());
    const intervalRef = useRef<ReturnType<typeof setInterval> | null>(null);

    const checkForChanges = useCallback(async () => {
        try {
            const response = await api.getCustomers({ limit: 1000 });
            const customers = response.data;

            customers.forEach((customer: Customer) => {
                const previousStatus = previousStatusRef.current.get(customer.id);

                if (previousStatus && previousStatus !== customer.current_status) {
                    // Status changed!
                    const change: StatusChange = {
                        customerId: customer.id,
                        customerName: customer.name,
                        oldStatus: previousStatus,
                        newStatus: customer.current_status,
                        timestamp: new Date(),
                    };

                    onStatusChange?.(change);

                    // Invalidate relevant queries
                    queryClient.invalidateQueries({ queryKey: ['customers'] });
                    queryClient.invalidateQueries({ queryKey: ['customers-los'] });
                }

                previousStatusRef.current.set(customer.id, customer.current_status);
            });
        } catch (error) {
            console.error('Failed to check customer status:', error);
        }
    }, [queryClient, onStatusChange]);

    useEffect(() => {
        if (!enabled) {
            if (intervalRef.current) {
                clearInterval(intervalRef.current);
                intervalRef.current = null;
            }
            return;
        }

        // Initial fetch to populate the status map
        checkForChanges();

        // Set up polling
        intervalRef.current = setInterval(checkForChanges, pollingInterval);

        return () => {
            if (intervalRef.current) {
                clearInterval(intervalRef.current);
                intervalRef.current = null;
            }
        };
    }, [enabled, pollingInterval, checkForChanges]);

    return {
        checkNow: checkForChanges,
    };
}

// Simple toast notification for status changes
export function createStatusChangeNotification(change: StatusChange): string {
    const { customerName, oldStatus, newStatus } = change;

    if (newStatus === 'LOS') {
        return `⚠️ LOS Alert: ${customerName} lost signal!`;
    }

    if (oldStatus === 'LOS' && newStatus === 'ONLINE') {
        return `✅ Recovered: ${customerName} is back online`;
    }

    if (newStatus === 'ONLINE') {
        return `🟢 ${customerName} is now online`;
    }

    if (newStatus === 'OFFLINE') {
        return `⚫ ${customerName} went offline`;
    }

    return `${customerName}: ${oldStatus} → ${newStatus}`;
}
