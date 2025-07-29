import { useState, type FormEvent, type ChangeEvent } from "react";
import { useNavigate } from "react-router-dom";
import { createPrint, getPrintMetadata } from "../util/prints";
import { STLPreview } from "./StlPreview";
import logo from "../assets/logo-1.svg";
import { Navbar } from "./Navbar";

export const NewPrint = () => {
    const [file, setFile] = useState<File | null>(null);
    const [error, setError] = useState("");
    const [loading, setLoading] = useState(false);
    const [thumbnailUrl, setThumbnailUrl] = useState<string | null>(null);
    const [stlData, setStlData] = useState<string | null>(null);
    const [uploadingPreview, setUploadingPreview] = useState(false);
    const [customColor, setCustomColor] = useState("#ffffff");
    const navigate = useNavigate();

    const handleFileChange = async (e: ChangeEvent<HTMLInputElement>) => {
        if (!e.target.files || e.target.files.length === 0) return;
        const selectedFile = e.target.files[0];

        setFile(selectedFile);
        setThumbnailUrl(null);
        setStlData(null);
        setError("");

        await uploadThumbnail(selectedFile);
    };

    const uploadThumbnail = async (file: File) => {
        setUploadingPreview(true);
        try {
            const data = await getPrintMetadata(file);
            if (data.stl_file) {
                setStlData(data.stl_file);
            } else if (data.thumbnail) {
                setThumbnailUrl(data.thumbnail);
            } else {
                setError("Failed to generate preview.");
            }
        } catch (err: any) {
            setError(err.message || "Thumbnail upload failed.");
        } finally {
            setUploadingPreview(false);
        }
    };

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();
        if (!file) {
            setError("You must select a file.");
            return;
        }

        setLoading(true);
        setError("");

        try {
            await createPrint(file, customColor);
            navigate("/dashboard");
        } catch (err: any) {
            setError(err.message || "Upload failed.");
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="flex flex-col h-screen w-full">
            <Navbar />
            <div className="flex-1 flex flex-col items-center justify-center overflow-auto">
                <form
                    onSubmit={handleSubmit}
                    className="w-full max-w-sm space-y-4 p-6"
                >
                    <div className="flex justify-center">
                        <img className="h-12" src={logo} alt="spooler logo" />
                    </div>
                    <label
                        className="border-2 border-dashed border-gray-200 rounded-xl p-6 text-center cursor-pointer bg-gray-50 hover:bg-gray-100 transition block"
                    >
                        <input
                            type="file"
                            accept=".stl, .3mf"
                            className="hidden"
                            onChange={handleFileChange}
                            required
                        />
                        {file ? (
                            <p className="text-sm font-medium text-gray-700">
                                Selected: {file.name}
                            </p>
                        ) : (
                            <p className="text-sm text-gray-500">Drag or click to select a .stl or .3mf file</p>
                        )}
                    </label>

                    {file?.name.toLowerCase().endsWith(".stl") && (
                        <div>
                            <label className="block text-sm text-gray-700 mb-1">Custom Color</label>
                            <input
                                type="color"
                                value={customColor}
                                onChange={(e) => setCustomColor(e.target.value)}
                                className="w-16 h-10 p-1 border rounded"
                            />
                        </div>
                    )}

                    {uploadingPreview && (
                        <p className="text-sm text-gray-500 text-center">Generating preview...</p>
                    )}

                    {stlData && (
                        <div className="rounded overflow-hidden border h-64 bg-gray-50 flex items-center justify-center">
                            <STLPreview stlData={stlData} color={customColor} />
                        </div>
                    )}

                    {thumbnailUrl && !stlData && (
                        <div className="rounded overflow-hidden border">
                            <img src={thumbnailUrl} alt="Thumbnail" className="w-full" />
                        </div>
                    )}

                    <button
                        type="submit"
                        className="w-full bg-spooler-orange hover:bg-spooler-orange-light hover:cursor-pointer text-white py-2 rounded transition"
                        disabled={loading}
                    >
                        {loading ? "Submitting..." : "Submit Print"}
                    </button>

                    {error && (
                        <div className="text-center text-red-600">{error}</div>
                    )}
                </form>
            </div>
        </div>
    );

};