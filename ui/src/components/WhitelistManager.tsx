import { useEffect, useState } from "react";
import { getWhitelist, addWhitelist, removeWhitelist } from "../util/auth";

export default function WhitelistManager() {
    const [emails, setEmails] = useState<string[]>([]);
    const [input, setInput] = useState("");
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");
    const [enabled, setEnabled] = useState(true)

    function parseEmails(input: string): string[] {
    return input
        .split(/[\s,;]+/)
        .map(e => e.trim())
        .filter(e => e.length > 0);
    }

    const handleAdd = async () => {
        setLoading(true);
        setError("");
        try {
            const emails = parseEmails(input);
            if (emails.length === 0) throw new Error("No valid emails");
            await addWhitelist(emails);
            setInput("");
            fetchList();
        } catch (err: any) {
            setError(err.message || "Failed to add email(s)");
        } finally {
            setLoading(false);
        }
    };

    const fetchList = async () => {
        setLoading(true);
        setError("");
        try {
            const list = await getWhitelist();
            setEmails(list.map((e: any) => e.Email));
        } catch (error) {
            if (error == "whitelist not enabled") {
                setEnabled(false);
            }
            setError("Failed to fetch whitelist");
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => { fetchList(); }, []);

    const handleRemove = async (email: string) => {
        setLoading(true);
        setError("");
        try {
            await removeWhitelist(email);
            fetchList();
        } catch {
            setError("Failed to remove email");
        } finally {
            setLoading(false);
        }
    };

    return (
        <>
            {enabled && (
                <div className="border rounded p-4 my-6">
                    <h2 className="text-lg font-semibold mb-2">Whitelist</h2>
                    <div className="flex gap-2 mb-4">
                        <textarea
                            className="border px-2 py-1 rounded text-sm w-full"
                            placeholder="Add email(s), separated by comma, space, or newline"
                            value={input}
                            onChange={e => setInput(e.target.value)}
                            disabled={loading}
                        />
                        <button
                            className="bg-spooler-orange hover:bg-spooler-orange-light hover:cursor-pointer text-white px-3 py-1 rounded text-xs"
                            onClick={handleAdd}
                            disabled={loading || !input}
                        >
                            Add
                        </button>
                    </div>
                    {error && <div className="text-red-600 mb-2">{error}</div>}
                    <ul>
                        {emails.map(email => (
                            <li key={email} className="flex items-center justify-between border-b py-1">
                                <span>{email}</span>
                                <button
                                    className="text-xs text-red-600 hover:underline hover:cursor-pointer"
                                    onClick={() => handleRemove(email)}
                                    disabled={loading}
                                >
                                    Remove
                                </button>
                            </li>
                        ))}
                    </ul>
                </div>
            )}
        </>
    );
}