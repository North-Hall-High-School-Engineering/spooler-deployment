import axios from "axios";

const API_BASE_URL = import.meta.env.VITE_SERVER_URL || "http://localhost:8080";

export async function getPrints() {
    const res = await axios.get(`${API_BASE_URL}/me/prints`, { withCredentials: true });
    return res.data;
}

export async function createPrint(file: File, requestedFilamentColor: string) {
    const formData = new FormData();
    formData.append("file", file);
    formData.append("file_name", file.name);
    formData.append("requested_filament_color", requestedFilamentColor);

    const res = await axios.post(`${API_BASE_URL}/prints/new`, formData, {
        withCredentials: true,
        headers: { "Content-Type": "multipart/form-data" },
    });
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

export async function updatePrintStatus(id: number, status: string, denialReason?: string) {
    const res = await axios.put(
        `${API_BASE_URL}/prints/${id}/status`,
        { status, denialReason },
        { withCredentials: true }
    );
    return res.data;
}

export async function deletePrint(id: number) {
    const res = await axios.delete(`${API_BASE_URL}/prints/${id}`, { withCredentials: true });
    return res.data;
}

