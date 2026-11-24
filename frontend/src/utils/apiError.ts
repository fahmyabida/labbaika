// Backend error responses are shaped { error_code, message } (see
// pkg/http/middleware/error_middleware.go). Falls back to the generic axios
// message, then a fixed string, so callers never show "undefined".
export function extractErrorMessage(err: unknown): string {
    const anyErr = err as any;
    return anyErr?.response?.data?.message || anyErr?.message || 'Something went wrong';
}
