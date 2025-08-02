import axios from "axios";

const API_BASE_URL = import.meta.env.VITE_SERVER_URL || "http://localhost:8080";

export async function getPrints() {
    const res = await axios.get(`${API_BASE_URL}/me/prints`, { withCredentials: true });
    return res.data;
}

export async function createPrint(file: File, requestedFilamentColor: string, onProgress?: (progress: number) => void) {
    const formData = new FormData();
    formData.append("file", file);
    formData.append("file_name", file.name);
    formData.append("requested_filament_color", requestedFilamentColor);

    const res = await axios.post(`${API_BASE_URL}/prints/new`, formData, {
        withCredentials: true,
        headers: { "Content-Type": "multipart/form-data" },
        timeout: 60_000, // 1 minute timeout
        onUploadProgress: (progressEvent) => {
            if (progressEvent.total && onProgress) {
                const progress = (progressEvent.loaded / progressEvent.total) * 100;
                onProgress(Math.min(progress, 99));
            }
        },
    });

    if (onProgress) onProgress(100);
    return res.data;
}

export async function getPrintMetadata(file: File) {
    const formData = new FormData();
    formData.append("file", file);
    formData.append("file_name", file.name);

    const res = await axios.post(`${API_BASE_URL}/metadata`, formData, {
        withCredentials: true,
        headers: { "Content-Type": "multipart/form-data" },
    });
    return res.data.metadata;
}

export async function getAllPrints() {
    const res = await axios.get(`${API_BASE_URL}/prints/all`, { withCredentials: true });
    return res.data;
}

export async function updatePrint(id: number, update: Partial<{
    status?: string;
    denial_reason?: string;
}>) {
    const res = await axios.put(
        `${API_BASE_URL}/prints/${id}`,
        update,
        { withCredentials: true }
    );
    return res.data;
}

export async function deletePrint(id: number) {
    const res = await axios.delete(`${API_BASE_URL}/prints/${id}`, { withCredentials: true });
    return res.data;
}

export const getPrintPreview = async (file: File) => {
    const formData = new FormData();
    formData.append("file", file);

    const response = await axios.post(`${API_BASE_URL}/preview`, formData, {
        headers: { "Content-Type": "multipart/form-data" },
        withCredentials: true,
    });

    return response.data;
};