import axios from "axios";
import { BACKEND_URL } from ".";

export interface Category {
    id?: string;
    name: string;
}

const CATEGORY_URL = `${BACKEND_URL}/api/v1/category`;

export async function listCategories(): Promise<Category[]> {
    const result = await axios.get(CATEGORY_URL);
    return result.data;
}

export async function createCategory(data: Category): Promise<Category> {
    const result = await axios.post(CATEGORY_URL, data);
    return result.data;
}

export async function updateCategory(id: string, data: Category): Promise<Category> {
    const result = await axios.put(`${CATEGORY_URL}/${id}`, data);
    return result.data;
}

export async function deleteCategory(id: string): Promise<void> {
    await axios.delete(`${CATEGORY_URL}/${id}`);
}
