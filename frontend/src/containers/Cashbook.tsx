import React, { useEffect, useState } from 'react';
import {
    Box, Button, MenuItem, Select, TextField, Table, TableBody, TableCell,
    TableContainer, TableHead, TableRow, Paper, IconButton, Stack, Autocomplete,
    Dialog, DialogContent,
} from '@mui/material';
import VisibilityIcon from '@mui/icons-material/Visibility';
import DeleteIcon from '@mui/icons-material/Delete';
import EditIcon from '@mui/icons-material/Edit';
import AttachFileIcon from '@mui/icons-material/AttachFile';
import AddIcon from '@mui/icons-material/Add';
import CloseIcon from '@mui/icons-material/Close';
import ArrowCircleDownIcon from '@mui/icons-material/ArrowCircleDown';
import ArrowCircleUpIcon from '@mui/icons-material/ArrowCircleUp';
import { Cashbook, createCashbook, updateCashbook, deleteCashbook, listCashbooks, uploadReceipt } from '../store/Cashbook';
import { listCategories } from '../store/Category';
import { BACKEND_URL } from '../store';
import { ErrorSnackbar } from '../components/ErrorSnackbar';
import { ConfirmDialog } from '../components/ConfirmDialog';
import { SortableTableCell } from '../components/SortableTableCell';
import { extractErrorMessage } from '../utils/apiError';
import { currentMonthRange } from '../utils/dates';
import { sortRows, SortDirection } from '../utils/sort';
import { saveDraftsText, loadDraftsText, clearDraftsText, saveDraftFile, loadDraftFile, clearDraftFile } from '../utils/draftCache';

// Indonesian grouping ("50.000"), no decimals — accounting-style display only.
const formatAmount = (value: number) =>
    new Intl.NumberFormat('id-ID', { maximumFractionDigits: 0 }).format(value);

// Live-format the amount input as the user types, while the underlying state
// (`amount`) stays plain digits — that's what actually gets sent/stored.
const formatAmountInput = (digits: string) => (digits ? formatAmount(Number(digits)) : '');
const stripToDigits = (value: string) => value.replace(/\D/g, '');

// IN = money coming in (green, down arrow — "dropping into" the account).
// OUT = money going out (red, up arrow — "leaving" the account).
const TYPE_META = {
    IN: { color: 'success.main' as const, Icon: ArrowCircleDownIcon },
    OUT: { color: 'error.main' as const, Icon: ArrowCircleUpIcon },
};

const TypeLabel = ({ type }: { type: 'IN' | 'OUT' }) => {
    const { color, Icon } = TYPE_META[type];
    return (
        <Box sx={{ display: 'flex', alignItems: 'center', color, gap: 0.5 }}>
            <Icon fontSize="small" />
            {type}
        </Box>
    );
};

interface Draft {
    id: string;
    type: 'IN' | 'OUT';
    amount: string;
    category: string;
    description: string;
    file: File | null;
}

const makeId = () => (crypto.randomUUID ? crypto.randomUUID() : `${Date.now()}-${Math.random()}`);
const emptyDraft = (): Draft => ({ id: makeId(), type: 'IN', amount: '', category: '', description: '', file: null });

// A draft with nothing typed in at all (e.g. an extra row the user added then
// didn't use) is exempt from validation — it's silently skipped on submit, so
// it shouldn't show error styling either.
const isDraftBlank = (d: Draft) => !d.amount && !d.category && !d.description && !d.file;
const isDraftValid = (d: Draft) => Number(d.amount) > 0 && d.category.trim() !== '';

export const CashbookPage = () => {
    const [type, setType] = useState<'IN' | 'OUT'>('IN');
    const [amount, setAmount] = useState('');
    const [category, setCategory] = useState('');
    const [description, setDescription] = useState('');
    const [file, setFile] = useState<File | null>(null);
    const [startDate, setStartDate] = useState(() => currentMonthRange().start);
    const [endDate, setEndDate] = useState(() => currentMonthRange().end);
    const [entries, setEntries] = useState<Cashbook[]>([]);
    const [categoryOptions, setCategoryOptions] = useState<string[]>([]);
    const [editingId, setEditingId] = useState<string | null>(null);
    const [viewingReceipt, setViewingReceipt] = useState<string | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [pendingDeleteId, setPendingDeleteId] = useState<string | null>(null);
    const [submitAttempted, setSubmitAttempted] = useState(false);
    // Text restores synchronously from localStorage so there's no window where
    // the persist-on-change effect could race and clobber it with blank
    // defaults. Files are patched in afterward (IndexedDB reads are async).
    const [drafts, setDrafts] = useState<Draft[]>(() => {
        const cached = loadDraftsText();
        return cached && cached.length > 0 ? cached.map((d) => ({ ...d, file: null })) : [emptyDraft()];
    });
    const [previewUrl, setPreviewUrl] = useState<string | null>(null);
    const [sortKey, setSortKey] = useState<keyof Cashbook>('created_at');
    const [sortDir, setSortDir] = useState<SortDirection>('desc');

    const handleSort = (key: keyof Cashbook) => {
        if (key === sortKey) {
            setSortDir(sortDir === 'asc' ? 'desc' : 'asc');
        } else {
            setSortKey(key);
            setSortDir('asc');
        }
    };

    const loadEntries = (sd?: string, ed?: string) => {
        listCashbooks(sd, ed)
            .then(setEntries)
            .catch((err) => setError(extractErrorMessage(err)));
    };

    const loadCategories = () => {
        listCategories()
            .then((categories) => setCategoryOptions(categories.map((c) => c.name)))
            .catch((err) => setError(extractErrorMessage(err)));
    };

    useEffect(() => {
        loadEntries(startDate, endDate);
        loadCategories();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    // Patch in each restored draft's cached file from IndexedDB (text itself
    // was already restored synchronously in useState above). Only re-renders
    // if a file was actually found, so this is a no-op for a fresh session.
    useEffect(() => {
        Promise.all(drafts.map(async (d) => ({ ...d, file: await loadDraftFile(d.id) })))
            .then((withFiles) => {
                if (withFiles.some((d) => d.file)) setDrafts(withFiles);
            });
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    // Persist text fields on every change. File blobs are saved immediately
    // when picked (see updateDraft) — not here, since this effect skips the
    // (non-serializable) file field anyway.
    useEffect(() => {
        saveDraftsText(drafts.map(({ id, type, amount, category, description }) => ({ id, type, amount, category, description })));
    }, [drafts]);

    const resetForm = () => {
        setEditingId(null);
        setType('IN');
        setAmount('');
        setCategory('');
        setDescription('');
        setFile(null);
        setSubmitAttempted(false);
    };

    const updateDraft = (index: number, patch: Partial<Draft>) => {
        setDrafts((prev) => prev.map((d, i) => (i === index ? { ...d, ...patch } : d)));
        if ('file' in patch) {
            saveDraftFile(drafts[index].id, patch.file ?? null);
        }
    };
    const addDraftRow = () => setDrafts((prev) => [...prev, emptyDraft()]);
    const removeDraftRow = (index: number) => {
        const removed = drafts[index];
        if (drafts.length <= 1) return;
        setDrafts((prev) => prev.filter((_, i) => i !== index));
        clearDraftFile(removed.id);
    };

    const openPreview = (file: File) => setPreviewUrl(URL.createObjectURL(file));
    const closePreview = () => {
        if (previewUrl) URL.revokeObjectURL(previewUrl);
        setPreviewUrl(null);
    };

    const handleSubmit = () => {
        if (editingId) {
            if (!(Number(amount) > 0 && category.trim() !== '')) {
                setSubmitAttempted(true);
                return;
            }
            const payload: Cashbook = { type, amount: Number(amount), category, description };
            const pendingFile = file;
            const id = editingId;

            updateCashbook(id, payload)
                .then(() => {
                    resetForm();
                    loadCategories();
                    return pendingFile ? uploadReceipt(id, pendingFile) : undefined;
                })
                .then(() => loadEntries(startDate, endDate))
                .catch((err) => setError(extractErrorMessage(err)));
            return;
        }

        // Any draft with SOME data must be fully valid — otherwise a
        // half-filled row would get silently dropped and the user would
        // think it saved.
        const touchedDrafts = drafts.filter((d) => !isDraftBlank(d));
        if (touchedDrafts.length === 0) return;
        if (touchedDrafts.some((d) => !isDraftValid(d))) {
            setSubmitAttempted(true);
            return;
        }
        const validDrafts = touchedDrafts;

        // Submitted one at a time (not Promise.all) so a failure partway
        // through still leaves earlier successful rows saved and visible —
        // Promise.all would fire all requests concurrently and make it
        // unclear which of several rows actually landed if one rejects.
        let chain: Promise<unknown> = Promise.resolve();
        validDrafts.forEach((d) => {
            chain = chain.then(() =>
                createCashbook({ type: d.type, amount: Number(d.amount), category: d.category, description: d.description })
                    .then((created) => (d.file && created.id ? uploadReceipt(created.id, d.file) : undefined))
            );
        });

        chain
            .then(() => {
                drafts.forEach((d) => clearDraftFile(d.id));
                clearDraftsText();
                setDrafts([emptyDraft()]);
                setSubmitAttempted(false);
                loadCategories();
            })
            .catch((err) => setError(extractErrorMessage(err)))
            .finally(() => loadEntries(startDate, endDate));
    };

    const handleEdit = (entry: Cashbook) => {
        setEditingId(entry.id ?? null);
        setType(entry.type);
        setAmount(String(entry.amount));
        setCategory(entry.category);
        setDescription(entry.description);
        setFile(null);
    };

    const handleDelete = (id?: string) => {
        if (!id) return;
        setPendingDeleteId(id);
    };

    const confirmDelete = () => {
        if (!pendingDeleteId) return;
        deleteCashbook(pendingDeleteId)
            .then(() => loadEntries(startDate, endDate))
            .catch((err) => setError(extractErrorMessage(err)))
            .finally(() => setPendingDeleteId(null));
    };

    const handleAttach = (id: string | undefined, selectedFile: File | undefined) => {
        if (!id || !selectedFile) return;
        uploadReceipt(id, selectedFile)
            .then(() => loadEntries(startDate, endDate))
            .catch((err) => setError(extractErrorMessage(err)));
    };

    return (
        <Box>
            {editingId ? (
                <Stack direction="row" spacing={2} sx={{ mb: 2, flexWrap: 'wrap' }}>
                    <Select value={type} onChange={(e) => setType(e.target.value as 'IN' | 'OUT')} size="small">
                        <MenuItem value="IN"><TypeLabel type="IN" /></MenuItem>
                        <MenuItem value="OUT"><TypeLabel type="OUT" /></MenuItem>
                    </Select>
                    <TextField
                        size="small"
                        label="Amount"
                        value={formatAmountInput(amount)}
                        onChange={(e) => setAmount(stripToDigits(e.target.value))}
                        error={submitAttempted && !(Number(amount) > 0)}
                        helperText={submitAttempted && !(Number(amount) > 0) ? 'Must be > 0' : ''}
                    />
                    <Autocomplete
                        size="small"
                        sx={{ width: 200 }}
                        options={categoryOptions}
                        value={category || null}
                        onChange={(_, value) => setCategory(value ?? '')}
                        renderInput={(params) => (
                            <TextField
                                {...params}
                                label="Category"
                                error={submitAttempted && category.trim() === ''}
                                helperText={submitAttempted && category.trim() === '' ? 'Required' : ''}
                            />
                        )}
                    />
                    <TextField size="small" label="Description" value={description} onChange={(e) => setDescription(e.target.value)} />
                    <Button variant="outlined" component="label" startIcon={<AttachFileIcon fontSize="small" />}>
                        {file ? file.name : 'Receipt'}
                        <input type="file" hidden accept="image/*" onChange={(e) => setFile(e.target.files?.[0] ?? null)} />
                    </Button>
                    <Button variant="contained" onClick={handleSubmit}>Save</Button>
                    <Button onClick={resetForm}>Cancel</Button>
                </Stack>
            ) : (
                <Box sx={{ mb: 2 }}>
                    <Stack spacing={1}>
                        {drafts.map((draft, index) => (
                            <Stack key={index} direction="row" spacing={2} sx={{ flexWrap: 'wrap', alignItems: 'center' }}>
                                <Select value={draft.type} onChange={(e) => updateDraft(index, { type: e.target.value as 'IN' | 'OUT' })} size="small">
                                    <MenuItem value="IN"><TypeLabel type="IN" /></MenuItem>
                                    <MenuItem value="OUT"><TypeLabel type="OUT" /></MenuItem>
                                </Select>
                                <TextField
                                    size="small"
                                    label="Amount"
                                    value={formatAmountInput(draft.amount)}
                                    onChange={(e) => updateDraft(index, { amount: stripToDigits(e.target.value) })}
                                    error={submitAttempted && !isDraftBlank(draft) && !(Number(draft.amount) > 0)}
                                    helperText={submitAttempted && !isDraftBlank(draft) && !(Number(draft.amount) > 0) ? 'Must be > 0' : ''}
                                />
                                <Autocomplete
                                    size="small"
                                    sx={{ width: 200 }}
                                    options={categoryOptions}
                                    value={draft.category || null}
                                    onChange={(_, value) => updateDraft(index, { category: value ?? '' })}
                                    renderInput={(params) => (
                                        <TextField
                                            {...params}
                                            label="Category"
                                            error={submitAttempted && !isDraftBlank(draft) && draft.category.trim() === ''}
                                            helperText={submitAttempted && !isDraftBlank(draft) && draft.category.trim() === '' ? 'Required' : ''}
                                        />
                                    )}
                                />
                                <TextField
                                    size="small"
                                    label="Description"
                                    value={draft.description}
                                    onChange={(e) => updateDraft(index, { description: e.target.value })}
                                />
                                <Button variant="outlined" component="label" startIcon={<AttachFileIcon fontSize="small" />}>
                                    {draft.file ? draft.file.name : 'Receipt'}
                                    <input
                                        type="file"
                                        hidden
                                        accept="image/*"
                                        onChange={(e) => updateDraft(index, { file: e.target.files?.[0] ?? null })}
                                    />
                                </Button>
                                {draft.file && (
                                    <IconButton size="small" onClick={() => openPreview(draft.file as File)}>
                                        <VisibilityIcon fontSize="small" />
                                    </IconButton>
                                )}
                                {drafts.length > 1 && (
                                    <IconButton size="small" onClick={() => removeDraftRow(index)}>
                                        <CloseIcon fontSize="small" />
                                    </IconButton>
                                )}
                            </Stack>
                        ))}
                    </Stack>
                    <Stack direction="row" spacing={2} sx={{ mt: 1 }}>
                        <Button startIcon={<AddIcon fontSize="small" />} onClick={addDraftRow}>Tambah entry lain</Button>
                        <Button variant="contained" onClick={handleSubmit}>{drafts.length > 1 ? `Add all (${drafts.length})` : 'Add'}</Button>
                    </Stack>
                </Box>
            )}

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
                <Button variant="outlined" onClick={() => loadEntries(startDate, endDate)}>Filter</Button>
                <Button onClick={() => { setStartDate(''); setEndDate(''); loadEntries(); }}>Clear</Button>
            </Stack>

            <TableContainer component={Paper}>
                <Table size="small">
                    <TableHead>
                        <TableRow>
                            <SortableTableCell<Cashbook> label="Date" sortKey="created_at" activeKey={sortKey} direction={sortDir} onSort={handleSort} />
                            <SortableTableCell<Cashbook> label="Type" sortKey="type" activeKey={sortKey} direction={sortDir} onSort={handleSort} />
                            <SortableTableCell<Cashbook> label="Amount" sortKey="amount" activeKey={sortKey} direction={sortDir} onSort={handleSort} align="right" />
                            <SortableTableCell<Cashbook> label="Category" sortKey="category" activeKey={sortKey} direction={sortDir} onSort={handleSort} />
                            <SortableTableCell<Cashbook> label="Description" sortKey="description" activeKey={sortKey} direction={sortDir} onSort={handleSort} />
                            <TableCell>Receipt</TableCell>
                            <TableCell />
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {sortRows(entries, sortKey, sortDir).map((entry) => (
                            <TableRow key={entry.id}>
                                <TableCell>{entry.created_at?.slice(0, 10)}</TableCell>
                                <TableCell><TypeLabel type={entry.type} /></TableCell>
                                <TableCell align="right">{formatAmount(entry.amount)}</TableCell>
                                <TableCell>{entry.category}</TableCell>
                                <TableCell>{entry.description}</TableCell>
                                <TableCell>
                                    {entry.receipt_url && (
                                        <IconButton size="small" onClick={() => setViewingReceipt(entry.receipt_url!)}>
                                            <VisibilityIcon fontSize="small" />
                                        </IconButton>
                                    )}
                                    <IconButton size="small" component="label">
                                        <AttachFileIcon fontSize="small" />
                                        <input
                                            type="file"
                                            hidden
                                            accept="image/*"
                                            onChange={(e) => handleAttach(entry.id, e.target.files?.[0])}
                                        />
                                    </IconButton>
                                </TableCell>
                                <TableCell>
                                    <IconButton size="small" onClick={() => handleEdit(entry)}>
                                        <EditIcon fontSize="small" />
                                    </IconButton>
                                    <IconButton size="small" onClick={() => handleDelete(entry.id)}>
                                        <DeleteIcon fontSize="small" />
                                    </IconButton>
                                </TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </TableContainer>

            <Dialog open={!!viewingReceipt} onClose={() => setViewingReceipt(null)} maxWidth="md">
                <DialogContent>
                    {viewingReceipt && (
                        <img src={`${BACKEND_URL}${viewingReceipt}`} alt="receipt" style={{ maxWidth: '100%', display: 'block' }} />
                    )}
                </DialogContent>
            </Dialog>

            <Dialog open={!!previewUrl} onClose={closePreview} maxWidth="md">
                <DialogContent>
                    {previewUrl && (
                        <img src={previewUrl} alt="receipt preview" style={{ maxWidth: '100%', display: 'block' }} />
                    )}
                </DialogContent>
            </Dialog>

            <ErrorSnackbar message={error} onClose={() => setError(null)} />

            <ConfirmDialog
                open={!!pendingDeleteId}
                title="Delete entry?"
                message="This cashbook entry will be deleted. This can't be undone."
                onConfirm={confirmDelete}
                onCancel={() => setPendingDeleteId(null)}
            />
        </Box>
    );
};
