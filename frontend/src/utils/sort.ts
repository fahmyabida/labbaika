export type SortDirection = 'asc' | 'desc';

export function sortRows<T>(rows: T[], key: keyof T, direction: SortDirection): T[] {
    return [...rows].sort((a, b) => {
        const av = a[key];
        const bv = b[key];
        if (av == null && bv == null) return 0;
        if (av == null) return 1;
        if (bv == null) return -1;

        if (typeof av === 'number' && typeof bv === 'number') {
            return direction === 'asc' ? av - bv : bv - av;
        }

        const as = String(av).toLowerCase();
        const bs = String(bv).toLowerCase();
        if (as < bs) return direction === 'asc' ? -1 : 1;
        if (as > bs) return direction === 'asc' ? 1 : -1;
        return 0;
    });
}
