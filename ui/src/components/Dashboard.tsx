import { useAuth } from "../context/authContext";
import { Navigate } from "react-router-dom";
import { Navbar } from "./Navbar";
import { useEffect, useState } from "react";
import axios from "axios";
import type { Print } from "../types/print";

function Dashboard() {
    const { isAuthenticated, loading } = useAuth();
    const [prints, setPrints] = useState<Print[]>([]);
    const [loadingPrints, setLoadingPrints] = useState(true);
    const [error, setError] = useState("");

    useEffect(() => {
        if (!loading && isAuthenticated) {
            setLoadingPrints(true);
            axios
                .get("/me/prints", { withCredentials: true, baseURL: import.meta.env.VITE_SERVER_URL || "http://localhost:8080" })
                .then(res => setPrints(res.data))
                .catch(() => setError("Failed to load prints"))
                .finally(() => setLoadingPrints(false));
        }
    }, [isAuthenticated, loading]);

    if (loading) return null;
    if (!isAuthenticated) return <Navigate to="/login" replace />;

    return (
        <>
            <Navbar />
            <div className="max-w-4xl mx-auto mt-10 px-4">
                {loadingPrints ? (
                    <div className="text-gray-500">Loading...</div>
                ) : error ? (
                    <div className="text-red-600">{error}</div>
                ) : prints.length === 0 ? (
                    <div className="text-gray-500">No prints found.</div>
                ) : (
                    <div className="overflow-x-auto border-l border-r border-t">
                        <table className="min-w-full text-sm">
                            <thead>
                                <tr className="bg-gray-100 text-left">
                                    <th className="px-4 py-2 border-b">File</th>
                                    <th className="px-4 py-2 border-b">Status</th>
                                    <th className="px-4 py-2 border-b">Color</th>
                                    <th className="px-4 py-2 border-b">Created</th>
                                    <th className="px-4 py-2 border-b">Denial Reason</th>
                                    <th className="px-4 py-2 border-b">Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                {prints.map(print => (
                                    <tr key={print.ID} className="even:bg-gray-50">
                                        <td className="px-4 py-2 border-b truncate">{print.UploadedFileName}</td>
                                        <td className="px-4 py-2 border-b capitalize">{print.Status}</td>
                                        <td className="px-4 py-2 border-b">
                                            <span className="inline-block w-5 h-5 rounded-full border mr-2 align-middle" style={{ background: print.RequestedFilamentColor }} />
                                            {print.RequestedFilamentColor}
                                        </td>
                                        <td className="px-4 py-2 border-b">{new Date(print.CreatedAt).toLocaleString()}</td>
                                        <td className="px-4 py-2 border-b text-xs text-gray-700 whitespace-pre-wrap">{print.DenialReason || "-"}</td>
                                        <td className="px-4 py-2 border-b">                
                                                <a
                                                    href={`${import.meta.env.VITE_SERVER_URL || "http://localhost:8080"}/bucket/${encodeURIComponent(print.StoredFileName)}`}
                                                    className="ml-2 px-2 py-1 rounded text-xs bg-blue-500 text-white hover:bg-blue-600 transition-colors cursor-pointer"
                                                    target="_blank"
                                                    rel="noopener noreferrer"
                                                    download
                                                >
                                                    Download
                                                </a>
                                            
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>
                )}
            </div>
        </>
    );
}

export default Dashboard;