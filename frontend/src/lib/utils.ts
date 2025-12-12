import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
    return twMerge(clsx(inputs));
}

export function formatDate(date: string | Date): string {
    return new Date(date).toLocaleDateString('id-ID', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
    });
}

export function formatDistance(meters: number): string {
    if (meters < 1000) {
        return `${Math.round(meters)} m`;
    }
    return `${(meters / 1000).toFixed(2)} km`;
}

export function formatRxPower(dbm: number | undefined | null): string {
    if (dbm === undefined || dbm === null) return 'N/A';
    return `${dbm.toFixed(2)} dBm`;
}

export function getRxPowerStatus(rxPower: number | undefined | null): 'good' | 'warning' | 'critical' | 'unknown' {
    if (rxPower === undefined || rxPower === null) return 'unknown';
    if (rxPower >= -25) return 'good';
    if (rxPower >= -27) return 'warning';
    return 'critical';
}
