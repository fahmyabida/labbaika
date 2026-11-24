import React, { useEffect, useState } from 'react';
import { Box, Stack, TextField, Button, Typography } from '@mui/material';
import { BarChart } from '@mui/x-charts/BarChart';
import { getCashbookSummary, CategorySummary } from '../store/Cashbook';
import { ErrorSnackbar } from '../components/ErrorSnackbar';
import { extractErrorMessage } from '../utils/apiError';
import { currentMonthRange } from '../utils/dates';

interface Row {
    label: string;
    IN: number;
    OUT: number;
    [key: string]: string | number;
}

// Collapse the flat (category, type, total) rows into one combined IN total
// and OUT total, ignoring category — single dataset row so the chart renders
// just two bars.
const aggregateTotals = (rows: CategorySummary[]): Row[] => {
    const totals: Row = { label: 'Total', IN: 0, OUT: 0 };
    for (const row of rows) {
        totals[row.type] += row.total;
    }
    return [totals];
};

// Pivot the flat (category, type, total) rows into one row per category with
// IN/OUT as columns — the shape @mui/x-charts' BarChart dataset wants.
const pivotByCategory = (rows: CategorySummary[]): Row[] => {
    const byCategory = new Map<string, Row>();
    for (const row of rows) {
        const entry = byCategory.get(row.category) ?? { label: row.category, IN: 0, OUT: 0 };
        entry[row.type] = row.total;
        byCategory.set(row.category, entry);
    }
    return Array.from(byCategory.values());
};

const isEmpty = (dataset: Row[]) => dataset.length === 0 || dataset.every((row) => row.IN === 0 && row.OUT === 0);

// Soft/pastel palette — deliberately muted, not the vivid theme success/error colors.
const PASTEL_GREEN = '#B5EAD7';
const PASTEL_CORAL = '#FFB7B2';

export const DashboardPage = () => {
    const [startDate, setStartDate] = useState(() => currentMonthRange().start);
    const [endDate, setEndDate] = useState(() => currentMonthRange().end);
    const [totalsDataset, setTotalsDataset] = useState<Row[]>([]);
    const [categoryDataset, setCategoryDataset] = useState<Row[]>([]);
    const [error, setError] = useState<string | null>(null);

    const series = [
        { dataKey: 'IN', label: 'IN', color: PASTEL_GREEN },
        { dataKey: 'OUT', label: 'OUT', color: PASTEL_CORAL },
    ];

    const loadSummary = (sd?: string, ed?: string) => {
        getCashbookSummary(sd, ed)
            .then((rows) => {
                setTotalsDataset(aggregateTotals(rows));
                setCategoryDataset(pivotByCategory(rows));
            })
            .catch((err) => setError(extractErrorMessage(err)));
    };

    useEffect(() => {
        loadSummary(startDate, endDate);
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    return (
        <Box>
            <Stack direction="row" spacing={2} sx={{ mb: 2, flexWrap: 'wrap' }}>
                <TextField
                    size="small"
                    label="Start date"
                    type="date"
                    InputLabelProps={{ shrink: true }}
                    value={startDate}
                    onChange={(e) => setStartDate(e.target.value)}
                />
                <TextField
                    size="small"
                    label="End date"
                    type="date"
                    InputLabelProps={{ shrink: true }}
                    value={endDate}
                    onChange={(e) => setEndDate(e.target.value)}
                />
                <Button variant="outlined" onClick={() => loadSummary(startDate, endDate)}>Filter</Button>
                <Button onClick={() => { setStartDate(''); setEndDate(''); loadSummary(); }}>Clear</Button>
            </Stack>

            <Typography variant="subtitle1" sx={{ mb: 1 }}>IN vs OUT</Typography>
            {isEmpty(totalsDataset) ? (
                <Typography color="text.secondary" sx={{ mb: 3 }}>No data in this range.</Typography>
            ) : (
                <BarChart
                    dataset={totalsDataset}
                    xAxis={[{ scaleType: 'band', dataKey: 'label' }]}
                    series={series}
                    height={300}
                />
            )}

            <Typography variant="subtitle1" sx={{ mt: 3, mb: 1 }}>IN vs OUT by category</Typography>
            {isEmpty(categoryDataset) ? (
                <Typography color="text.secondary">No data in this range.</Typography>
            ) : (
                <BarChart
                    dataset={categoryDataset}
                    xAxis={[{ scaleType: 'band', dataKey: 'label' }]}
                    series={series}
                    height={400}
                />
            )}

            <ErrorSnackbar message={error} onClose={() => setError(null)} />
        </Box>
    );
};
