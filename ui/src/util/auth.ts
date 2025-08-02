import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_SERVER_URL || 'http://localhost:8080';

export async function register(email: string, firstName: string, lastName: string) {
  const res = await axios.post(`${API_BASE_URL}/register`, {
    email,
    first_name: firstName,
    last_name: lastName,
  });
  return res.data;
}

export async function requestOTP(email: string) {
  const res = await axios.post(`${API_BASE_URL}/otp/request`, { email });
  return res.data;
}

export async function verifyOTP(email: string, code: string) {
  const res = await axios.post(`${API_BASE_URL}/otp/verify`, { email, code }, { withCredentials: true });
  return res.data; // contains { message, token }
}

export async function checkAuth() {
  try {
    const res = await axios.get(`${API_BASE_URL}/me`, { withCredentials: true });
    return res.data;
  } catch {
    return null
  }
}

export async function getWhitelist() {
    const res = await fetch(`${API_BASE_URL}/whitelist`, {
        credentials: "include",
    });
    if (!res.ok) throw new Error("Failed to fetch whitelist");
    return res.json();
}

export async function addWhitelist(emails: string[]) {
    const res = await fetch(`${API_BASE_URL}/whitelist`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ emails }),
    });
    if (!res.ok) throw new Error("Failed to add email(s)");
    return res.json();
}

export async function removeWhitelist(email: string) {
    const res = await fetch(`${API_BASE_URL}/whitelist`, {
        method: "DELETE",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
        body: JSON.stringify({ email }),
    });
    if (!res.ok) throw new Error("Failed to remove email");
    return res.json();
}

export async function getUserById(id: number) {
    const res = await axios.get(`${API_BASE_URL}/users/${id}`, {
        withCredentials: true,
    });
    return res.data;
}