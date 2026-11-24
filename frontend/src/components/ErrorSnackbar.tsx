import { Snackbar, Alert } from '@mui/material';

interface Props {
    message: string | null;
    onClose: () => void;
}

export const ErrorSnackbar = ({ message, onClose }: Props) => (
    <Snackbar open={!!message} autoHideDuration={5000} onClose={onClose}>
        <Alert severity="error" onClose={onClose} sx={{ width: '100%' }}>
            {message}
        </Alert>
    </Snackbar>
);
