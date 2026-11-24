import React, { useEffect, useState } from 'react';
import {
    Box, Button, TextField, Table, TableBody, TableCell,
    TableContainer, TableHead, TableRow, Paper, IconButton, Stack,
} from '@mui/material';
import DeleteIcon from '@mui/icons-material/Delete';
import EditIcon from '@mui/icons-material/Edit';
import { Category, createCategory, updateCategory, deleteCategory, listCategories } from '../store/Category';
import { ErrorSnackbar } from '../components/ErrorSnackbar';
import { ConfirmDialog } from '../components/ConfirmDialog';
import { SortableTableCell } from '../components/SortableTableCell';
import { extractErrorMessage } from '../utils/apiError';
import { sortRows, SortDirection } from '../utils/sort';

export const CategoryPage = () => {
    const [name, setName] = useState('');
    const [editingId, setEditingId] = useState<string | null>(null);
    const [categories, setCategories] = useState<Category[]>([]);
    const [error, setError] = useState<string | null>(null);
    const [pendingDeleteId, setPendingDeleteId] = useState<string | null>(null);
    const [sortKey, setSortKey] = useState<keyof Category>('name');
    const [sortDir, setSortDir] = useState<SortDirection>('asc');

    const handleSort = (key: keyof Category) => {
        if (key === sortKey) {
            setSortDir(sortDir === 'asc' ? 'desc' : 'asc');
        } else {
            setSortKey(key);
            setSortDir('asc');
        }
    };

    const loadCategories = () => {
        listCategories()
            .then(setCategories)
            .catch((err) => setError(extractErrorMessage(err)));
    };

    useEffect(() => {
        loadCategories();
    }, []);

    const resetForm = () => {
        setEditingId(null);
        setName('');
    };

    const handleSubmit = () => {
        if (!name.trim()) return;
        const save = editingId ? updateCategory(editingId, { name }) : createCategory({ name });
        save
            .then(() => {
                resetForm();
                loadCategories();
            })
            .catch((err) => setError(extractErrorMessage(err)));
    };

    const handleEdit = (category: Category) => {
        setEditingId(category.id ?? null);
        setName(category.name);
    };

    const handleDelete = (id?: string) => {
        if (!id) return;
        setPendingDeleteId(id);
    };

    const confirmDelete = () => {
        if (!pendingDeleteId) return;
        deleteCategory(pendingDeleteId)
            .then(loadCategories)
            .catch((err) => setError(extractErrorMessage(err)))
            .finally(() => setPendingDeleteId(null));
    };

    return (
        <Box>
            <Stack direction="row" spacing={2} sx={{ mb: 2, flexWrap: 'wrap' }}>
                <TextField size="small" label="Category name" value={name} onChange={(e) => setName(e.target.value)} />
                <Button variant="contained" onClick={handleSubmit}>{editingId ? 'Save' : 'Add'}</Button>
                {editingId && <Button onClick={resetForm}>Cancel</Button>}
            </Stack>

            <TableContainer component={Paper}>
                <Table size="small">
                    <TableHead>
                        <TableRow>
                            <SortableTableCell<Category> label="Name" sortKey="name" activeKey={sortKey} direction={sortDir} onSort={handleSort} />
                            <TableCell />
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {sortRows(categories, sortKey, sortDir).map((category) => (
                            <TableRow key={category.id}>
                                <TableCell>{category.name}</TableCell>
                                <TableCell align="right">
                                    <IconButton size="small" onClick={() => handleEdit(category)}>
                                        <EditIcon fontSize="small" />
                                    </IconButton>
                                    <IconButton size="small" onClick={() => handleDelete(category.id)}>
                                        <DeleteIcon fontSize="small" />
                                    </IconButton>
                                </TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </TableContainer>

            <ErrorSnackbar message={error} onClose={() => setError(null)} />

            <ConfirmDialog
                open={!!pendingDeleteId}
                title="Delete category?"
                message="This category will be removed from the master list. Cashbook entries already using it keep their text as-is."
                onConfirm={confirmDelete}
                onCancel={() => setPendingDeleteId(null)}
            />
        </Box>
    );
};
