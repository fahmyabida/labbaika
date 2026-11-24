import axios from "axios";
import { BACKEND_URL } from ".";

export interface Cashbook {
    id?: string;
    type: "IN" | "OUT";
    amount: number;
    category: string;
    description: string;
    created_at?: string;
    receipt_url?: string;
}

export interface CategorySummary {
    category: string;
    type: "IN" | "OUT";
    total: number;
}

const CASHBOOK_URL = `${BACKEND_URL}/api/v1/cashbook`;

export async function listCashbooks(startDate?: string, endDate?: string): Promise<Cashbook[]> {
    const params: Record<string, string> = {};
    if (startDate) params.start_date = startDate;
    if (endDate) params.end_date = endDate;

    const result = await axios.get(CASHBOOK_URL, { params });
    return result.data;
}

export async function createCashbook(data: Cashbook): Promise<Cashbook> {
    const result = await axios.post(CASHBOOK_URL, data);
    return result.data;
}

export async function updateCashbook(id: string, data: Cashbook): Promise<Cashbook> {
    const result = await axios.put(`${CASHBOOK_URL}/${id}`, data);
    return result.data;
}

export async function deleteCashbook(id: string): Promise<void> {
    await axios.delete(`${CASHBOOK_URL}/${id}`);
}

export async function uploadReceipt(id: string, file: File): Promise<Cashbook> {
    const formData = new FormData();
    formData.append("file", file);
    const result = await axios.post(`${CASHBOOK_URL}/${id}/receipt`, formData, {
        headers: { "Content-Type": "multipart/form-data" },
    });
    return result.data;
}

export async function getCashbookSummary(startDate?: string, endDate?: string): Promise<CategorySummary[]> {
    const params: Record<string, string> = {};
    if (startDate) params.start_date = startDate;
    if (endDate) params.end_date = endDate;

    const result = await axios.get(`${CASHBOOK_URL}/summary`, { params });
    return result.data;
}
