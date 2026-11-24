import { get as idbGet, set as idbSet, del as idbDel } from 'idb-keyval';

// Text fields go in localStorage (small, synchronous). Files go in IndexedDB
// via idb-keyval — File/Blob objects aren't JSON-serializable and wouldn't
// fit localStorage's small per-origin quota anyway.
const DRAFTS_KEY = 'labbaika:cashbookDrafts';
const fileKey = (draftId: string) => `labbaika:cashbookDraftFile:${draftId}`;

export interface CacheableDraft {
    id: string;
    type: 'IN' | 'OUT';
    amount: string;
    category: string;
    description: string;
}

export function saveDraftsText(drafts: CacheableDraft[]): void {
    try {
        localStorage.setItem(DRAFTS_KEY, JSON.stringify(drafts));
    } catch {
        // storage full/blocked (private mode, etc) — caching is a convenience, not critical.
    }
}

export function loadDraftsText(): CacheableDraft[] | null {
    try {
        const raw = localStorage.getItem(DRAFTS_KEY);
        return raw ? JSON.parse(raw) : null;
    } catch {
        return null;
    }
}

export function clearDraftsText(): void {
    try {
        localStorage.removeItem(DRAFTS_KEY);
    } catch {
        // ignore
    }
}

export async function saveDraftFile(draftId: string, file: File | null): Promise<void> {
    try {
        if (file) {
            await idbSet(fileKey(draftId), file);
        } else {
            await idbDel(fileKey(draftId));
        }
    } catch {
        // IndexedDB unavailable — the file just won't survive a reload.
    }
}

export async function loadDraftFile(draftId: string): Promise<File | null> {
    try {
        const file = await idbGet(fileKey(draftId));
        return file ?? null;
    } catch {
        return null;
    }
}

export async function clearDraftFile(draftId: string): Promise<void> {
    try {
        await idbDel(fileKey(draftId));
    } catch {
        // ignore
    }
}
