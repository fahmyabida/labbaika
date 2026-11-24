import { TableCell, TableSortLabel } from '@mui/material';
import { SortDirection } from '../utils/sort';

interface Props<T> {
    label: string;
    sortKey: keyof T;
    activeKey: keyof T;
    direction: SortDirection;
    onSort: (key: keyof T) => void;
    align?: 'left' | 'right' | 'center';
}

export function SortableTableCell<T>({ label, sortKey, activeKey, direction, onSort, align }: Props<T>) {
    return (
        <TableCell align={align}>
            {/* TableCell align="right" sets flexDirection: row-reverse on itself,
                which TableSortLabel inherits — putting the icon BEFORE the label
                for right-aligned columns (MUI's Material-Design convention).
                Forced back to 'row' here so every column's icon sits after the
                label, consistently, regardless of its text alignment. */}
            <TableSortLabel
                active={activeKey === sortKey}
                direction={activeKey === sortKey ? direction : 'asc'}
                onClick={() => onSort(sortKey)}
                sx={{ flexDirection: 'row' }}
            >
                {label}
            </TableSortLabel>
        </TableCell>
    );
}
