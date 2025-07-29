import { useEffect, useState } from "react";
import { useAuth } from "../context/authContext";
import { Navbar } from "./Navbar";
import { useNavigate } from "react-router-dom";
import type { Print, PrintStatus } from "../types/print";
import { getAllPrints, updatePrintStatus, deletePrint } from "../util/prints";

const PRINT_STATUS_OPTIONS: { label: string; value: PrintStatus }[] = [
    { label: "Approval Pending", value: "approval_pending" },
    { label: "Pending Print", value: "pending_print" },
    { label: "Printing", value: "printing" },
    { label: "Denied", value: "denied" },
    { label: "Completed", value: "completed" },
    { label: "Failed", value: "failed" },
    { label: "Canceled", value: "canceled" },
    { label: "Paused", value: "paused" },
];

function Administrator() {
    const { isAuthenticated, role, loading } = useAuth();
    const navigate = useNavigate();
    const [prints, setPrints] = useState<Print[]>([]);
    const [loadingPrints, setLoadingPrints] = useState(true);
    const [error, setError] = useState("");
    const [actionLoading, setActionLoading] = useState(false);
    const [selected, setSelected] = useState<number[]>([]);
    const [sortKey, setSortKey] = useState<keyof Print>("CreatedAt");
    const [sortDir, setSortDir] = useState<"asc" | "desc">("desc");
    const [denyModal, setDenyModal] = useState<{ open: boolean; id: number | null }>({ open: false, id: null });
    const [denyReason, setDenyReason] = useState("");

    useEffect(() => {
        if (loading) return;
        if (!isAuthenticated) navigate("/login", { replace: true });
        else if (role && role !== "admin") navigate("/dashboard", { replace: true });
    }, [isAuthenticated, role, loading, navigate]);

    useEffect(() => {
        fetchPrints();
    }, []);

    const fetchPrints = () => {
        setLoadingPrints(true);
        getAllPrints()
            .then(setPrints)
            .catch(() => setError("Failed to load prints"))
            .finally(() => setLoadingPrints(false));
    };

    const sortedPrints = [...prints].sort((a, b) => {
        let vA = a[sortKey], vB = b[sortKey];
        if (sortKey === "CreatedAt") {
            vA = new Date(a.CreatedAt).getTime();
            vB = new Date(b.CreatedAt).getTime();
        }
        if (vA! < vB!) return sortDir === "asc" ? -1 : 1;
        if (vA! > vB!) return sortDir === "asc" ? 1 : -1;
        return 0;
    });

    const toggleSelect = (id: number) => {
        setSelected((prev) =>
            prev.includes(id) ? prev.filter((x) => x !== id) : [...prev, id]
        );
    };

    const selectAll = () => {
        if (selected.length === prints.length) setSelected([]);
        else setSelected(prints.map((p) => p.ID));
    };

    const batchUpdateStatus = async (status: PrintStatus) => {
        setActionLoading(true);
        setError("");
        try {
            await Promise.all(
                selected.map((id) =>
                    updatePrintStatus(id, status)
                )
            );
            setPrints((prev) =>
                prev.map((p) =>
                    selected.includes(p.ID)
                        ? { ...p, Status: status, DenialReason: status === "denied" ? "Manually denied" : undefined }
                        : p
                )
            );
            setSelected([]);
        } catch {
            setError("Batch update failed");
        } finally {
            setActionLoading(false);
        }
    };

    const batchDelete = async () => {
        setActionLoading(true);
        setError("");
        try {
            await Promise.all(
                selected.map((id) =>
                    deletePrint(id)
                )
            );
            setPrints((prev) => prev.filter((p) => !selected.includes(p.ID)));
            setSelected([]);
        } catch {
            setError("Batch delete failed");
        } finally {
            setActionLoading(false);
        }
    };

    const updateStatus = async (id: number, status: PrintStatus, denialReason?: string) => {
        setActionLoading(true);
        setError("");
        try {
            await updatePrintStatus(id, status, denialReason);
            setPrints((prev) =>
                prev.map((p) =>
                    p.ID === id
                        ? { ...p, Status: status, DenialReason: status === "denied" ? (denialReason || "Manually denied") : undefined }
                        : p
                )
            );
        } catch {
            setError("Failed to update status");
        } finally {
            setActionLoading(false);
        }
    };

    // Modal for deny reason
    const openDenyModal = (id: number) => {
        setDenyModal({ open: true, id });
        setDenyReason("");
    };
    const closeDenyModal = () => setDenyModal({ open: false, id: null });

    const handleDeny = async () => {
        if (denyModal.id !== null) {
            await updateStatus(denyModal.id, "denied", denyReason || "Manually denied");
            closeDenyModal();
        }
    };

    if (loading) return null;

    return (
        <div className="h-screen w-full flex flex-col">
            <Navbar />
            <main className="max-w-4xl mx-auto mt-10 px-4">
                <div className="flex items-center mb-4 gap-2">
                    <button
                        className={`bg-spooler-orange text-white px-3 py-1 rounded text-xs transition-colors ${
                            loadingPrints ? "cursor-not-allowed opacity-60" : "cursor-pointer hover:bg-orange-500"
                        }`}
                        onClick={selectAll}
                        disabled={loadingPrints}
                    >
                        {selected.length === prints.length ? "Unselect All" : "Select All"}
                    </button>
                    <select
                        className={`border px-2 py-1 rounded text-xs transition-colors ${
                            selected.length === 0 ||
                            actionLoading ||
                            prints
                                .filter((p) => selected.includes(p.ID))
                                .some(
                                    (p) =>
                                        p.Status === "denied" ||
                                        p.Status === "completed"
                                )
                                ? "cursor-not-allowed bg-gray-100 text-gray-400"
                                : "cursor-pointer hover:bg-gray-50"
                        }`}
                        onChange={(e) => batchUpdateStatus(e.target.value as PrintStatus)}
                        disabled={
                            selected.length === 0 ||
                            actionLoading ||
                            prints
                                .filter((p) => selected.includes(p.ID))
                                .some(
                                    (p) =>
                                        p.Status === "denied" ||
                                        p.Status === "completed"
                                )
                        }
                        value=""
                    >
                        <option value="" disabled>
                            Batch Update Status
                        </option>
                        {PRINT_STATUS_OPTIONS.map((opt) => (
                            <option
                                key={opt.value}
                                value={opt.value}
                                disabled={
                                    prints
                                        .filter((p) => selected.includes(p.ID))
                                        .some(
                                            (p) =>
                                                p.Status === "denied" ||
                                                (opt.value === "denied" && p.Status !== "approval_pending") ||
                                                (p.Status === "completed" && opt.value === "denied") ||
                                                (p.Status === "denied" && opt.value === "completed")
                                        )
                                }
                            >
                                {opt.label}
                            </option>
                        ))}
                    </select>
                    <button
                        className={`bg-red-500 text-white px-3 py-1 rounded text-xs transition-colors ${
                            selected.length === 0 || actionLoading
                                ? "cursor-not-allowed opacity-60"
                                : "cursor-pointer hover:bg-red-600"
                        }`}
                        onClick={batchDelete}
                        disabled={
                            selected.length === 0 ||
                            actionLoading
                        }
                    >
                        Delete Selected
                    </button>
                    <button
                        className={`ml-auto text-xs underline transition-colors ${
                            loadingPrints ? "cursor-not-allowed opacity-60" : "cursor-pointer hover:text-blue-600"
                        }`}
                        onClick={fetchPrints}
                        disabled={loadingPrints}
                    >
                        Refresh
                    </button>
                </div>
                {loadingPrints ? (
                    <p className="text-gray-500">Loading...</p>
                ) : error ? (
                    <p className="text-red-600">{error}</p>
                ) : prints.length === 0 ? (
                    <p className="text-gray-500">No prints found.</p>
                ) : (
                    <div className="overflow-x-auto border-l border-r border-t">
                        <table className="min-w-full text-sm">
                            <thead>
                                <tr className="bg-gray-100 text-left">
                                    <th className="px-2 py-2 border-b">
                                        <input
                                            type="checkbox"
                                            checked={selected.length === prints.length}
                                            onChange={selectAll}
                                        />
                                    </th>
                                    <th
                                        className="px-4 py-2 border-b cursor-pointer"
                                        onClick={() => {
                                            setSortKey("UploadedFileName");
                                            setSortDir(sortDir === "asc" ? "desc" : "asc");
                                        }}
                                    >
                                        File
                                    </th>
                                    <th
                                        className="px-4 py-2 border-b cursor-pointer"
                                        onClick={() => {
                                            setSortKey("Status");
                                            setSortDir(sortDir === "asc" ? "desc" : "asc");
                                        }}
                                    >
                                        Status
                                    </th>
                                    <th className="px-4 py-2 border-b">Color</th>
                                    <th
                                        className="px-4 py-2 border-b cursor-pointer"
                                        onClick={() => {
                                            setSortKey("CreatedAt");
                                            setSortDir(sortDir === "asc" ? "desc" : "asc");
                                        }}
                                    >
                                        Created
                                    </th>
                                    <th className="px-4 py-2 border-b">Denial Reason</th>
                                    <th className="px-4 py-2 border-b">Update Status</th>
                                    <th className="px-4 py-2 border-b">Actions</th>
                                </tr>
                            </thead>
                            <tbody>
                                {sortedPrints.map((print) => {
                                    return (
                                        <tr key={print.ID} className="even:bg-gray-50">
                                            <td className="px-2 py-2 border-b">
                                                <input
                                                    type="checkbox"
                                                    checked={selected.includes(print.ID)}
                                                    onChange={() => toggleSelect(print.ID)}
                                                />
                                            </td>
                                            <td className="px-4 py-2 border-b truncate">{print.StoredFileName.slice(0, 11 - 3) + '...'}</td>
                                            <td className="px-4 py-2 border-b capitalize">{print.Status}</td>
                                            <td className="px-4 py-2 border-b">
                                                <div className="flex items-center gap-2">
                                                    <span
                                                        className="inline-block w-4 h-4 rounded-full border"
                                                        style={{ background: print.RequestedFilamentColor }}
                                                    />
                                                    {print.RequestedFilamentColor}
                                                </div>
                                            </td>
                                            <td className="px-4 py-2 border-b">
                                                {new Date(print.CreatedAt).toLocaleString()}
                                            </td>
                                            <td className="px-4 py-2 border-b text-xs text-gray-700 whitespace-pre-wrap">
                                                {print.DenialReason || "-"}
                                            </td>
                                            <td className="px-4 py-2 border-b">
                                                <select
                                                    className={`border px-2 py-1 rounded text-xs transition-colors ${
                                                        actionLoading || print.Status === "completed" || print.Status === "denied"
                                                            ? "cursor-not-allowed bg-gray-100 text-gray-400"
                                                            : "cursor-pointer hover:bg-gray-50"
                                                    }`}
                                                    value={print.Status}
                                                    onChange={(e) =>
                                                        e.target.value === "denied"
                                                            ? openDenyModal(print.ID)
                                                            : updateStatus(print.ID, e.target.value as PrintStatus)
                                                    }
                                                    disabled={
                                                        actionLoading ||
                                                        print.Status === "completed" ||
                                                        print.Status === "denied"
                                                    }
                                                >
                                                    {PRINT_STATUS_OPTIONS.map((opt) => (
                                                        <option
                                                            key={opt.value}
                                                            value={opt.value}
                                                            disabled={
                                                                (opt.value === "denied" && print.Status !== "approval_pending") ||
                                                                (print.Status === "completed" && opt.value === "denied") ||
                                                                (print.Status === "denied" && opt.value === "completed")
                                                            }
                                                        >
                                                            {opt.label}
                                                        </option>
                                                    ))}
                                                </select>
                                            </td>
                                            <td className="px-4 py-2 border-b">
                                                {print.Status === "approval_pending" && (
                                                    <>
                                                    <button
                                                        className="mr-2 px-2 py-1 rounded text-xs transition-colors bg-green-500 text-white cursor-pointer hover:bg-green-600"
                                                        disabled={print.Status !== "approval_pending" || actionLoading}
                                                        onClick={() => updateStatus(print.ID, "pending_print")}
                                                    >
                                                        Approve
                                                    </button>

                                                    <button
                                                        className="px-2 py-1 rounded text-xs transition-colors bg-red-500 text-white cursor-pointer hover:bg-red-600"
                                                        disabled={print.Status !== "approval_pending" || actionLoading}
                                                        onClick={() => openDenyModal(print.ID)}
                                                    >
                                                        Deny
                                                    </button>
                                                    </>

                                                )}
                                                
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
                                    );
                                })}
                            </tbody>
                        </table>
                    </div>
                )}
                {denyModal.open && (
                    <div className="fixed inset-0 flex items-center justify-center z-50">
                        <div className="bg-white rounded shadow-lg p-6 w-full max-w-xs">
                            <h2 className="text-lg mb-2">Deny Print</h2>
                            <textarea
                                className="w-full border rounded p-2 mb-3 text-sm"
                                rows={3}
                                placeholder="Reason for denial"
                                value={denyReason}
                                onChange={e => setDenyReason(e.target.value)}
                                autoFocus
                            />
                            <div className="flex justify-end gap-2">
                                <button
                                    className="px-3 py-1 rounded text-xs bg-gray-200 hover:bg-gray-300"
                                    onClick={closeDenyModal}
                                >
                                    Cancel
                                </button>
                                <button
                                    className="px-3 py-1 rounded text-xs bg-red-500 text-white hover:bg-red-600"
                                    onClick={handleDeny}
                                    disabled={actionLoading || !denyReason.trim()}
                                >
                                    Deny
                                </button>
                            </div>
                        </div>
                    </div>
                )}
            </main>
        </div>
    );
}

export default Administrator;