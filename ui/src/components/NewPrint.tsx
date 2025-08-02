import { useState, type FormEvent, type ChangeEvent, type DragEvent } from "react";
import { useNavigate } from "react-router-dom";
import { createPrint, getPrintPreview } from "../util/prints";
import { Model3DPreview } from "./Model3DPreview";
import { Navbar } from "./Navbar";

export const NewPrint = () => {
    const [file, setFile] = useState<File | null>(null);
    const [error, setError] = useState("");
    const [loading, setLoading] = useState(false);
    const [uploadProgress, setUploadProgress] = useState(0);
    // const [thumbnailUrl, setThumbnailUrl] = useState<string | null>(null);
    const [stlData, setStlData] = useState<string | null>(null);
    const [uploadingPreview, setUploadingPreview] = useState(false);
    const [customColor, setCustomColor] = useState("#ffffff");
    const [isDragOver, setIsDragOver] = useState(false);
    const navigate = useNavigate();

    const validateFile = (file: File): boolean => {
        const validExtensions = ['.stl', '.3mf'];
        const fileName = file.name.toLowerCase();
        return validExtensions.some(ext => fileName.endsWith(ext));
    };

    const processFile = async (selectedFile: File) => {
        if (!validateFile(selectedFile)) {
            setError("Please select a valid .stl or .3mf file.");
            return;
        }

        setFile(selectedFile);
        // setThumbnailUrl(null);
        setStlData(null);
        setError("");

        await uploadPreview(selectedFile);
    };

    const handleFileChange = async (e: ChangeEvent<HTMLInputElement>) => {
        if (!e.target.files || e.target.files.length === 0) return;
        const selectedFile = e.target.files[0];
        await processFile(selectedFile);
    };

    const handleDragOver = (e: DragEvent<HTMLLabelElement>) => {
        e.preventDefault();
        setIsDragOver(true);
    };

    const handleDragLeave = (e: DragEvent<HTMLLabelElement>) => {
        e.preventDefault();
        setIsDragOver(false);
    };

    const handleDrop = async (e: DragEvent<HTMLLabelElement>) => {
        e.preventDefault();
        setIsDragOver(false);

        const droppedFiles = e.dataTransfer.files;
        if (droppedFiles.length > 0) {
            await processFile(droppedFiles[0]);
            const input = e.currentTarget.querySelector('input[type="file"]') as HTMLInputElement;
            if (input) {
                const dataTransfer = new DataTransfer();
                dataTransfer.items.add(droppedFiles[0]);
                input.files = dataTransfer.files;
            }
        }
    };

    const uploadPreview = async (file: File) => {
        setUploadingPreview(true);
        try {
            const data = await getPrintPreview(file);
            if (data.type === "stl" && data.model_data) {
                setStlData(data.model_data);
                // setThumbnailUrl(null);
            } else if (data.type === "3mf" && data.model_data) {
                setStlData(data.model_data);
                // setThumbnailUrl(null);
            } else {
                setError("Preview not available for this file.");
                setStlData(null);
                // setThumbnailUrl(null);
            }
        } catch (err: any) {
            setError(err.message || "Preview generation failed.");
            setStlData(null);
            // setThumbnailUrl(null);
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
        setUploadProgress(0);

        try {
            await createPrint(file, customColor, (progress) => {
                setUploadProgress(progress);
            });
            navigate("/dashboard");
        } catch (err: any) {
            setError(err.message || "Upload failed.");
        } finally {
            setLoading(false);
            setUploadProgress(0);
        }
    };

    return (
        <div className="flex flex-col h-screen w-full">
            <Navbar />
            <div className="flex-1 flex flex-col items-center justify-center overflow-auto p-4">
                <form
                    onSubmit={handleSubmit}
                    className="w-full max-w-2xl space-y-6 p-6"
                >
                    <label
                        className={`border-2 border-dashed rounded-xl p-12 text-center cursor-pointer transition-all duration-200 block min-h-[200px] flex flex-col items-center justify-center ${
                            isDragOver 
                                ? "border-spooler-orange bg-orange-50 border-solid" 
                                : file 
                                    ? "border-green-300 bg-green-50" 
                                    : "border-gray-300 bg-gray-50 hover:bg-gray-100 hover:border-gray-400"
                        }`}
                        onDragOver={handleDragOver}
                        onDragLeave={handleDragLeave}
                        onDrop={handleDrop}
                    >
                        <input
                            type="file"
                            accept=".stl, .3mf"
                            className="hidden"
                            onChange={handleFileChange}
                        />
                        {file ? (
                            <div className="space-y-2">
                                <div className="text-green-600 text-4xl">✓</div>
                                <p className="text-lg font-medium text-gray-700">
                                    {file.name}
                                </p>
                                <p className="text-sm text-gray-500">
                                    Click to select a different file or drag and drop to replace
                                </p>
                            </div>
                        ) : (
                            <div className="space-y-4">
                                <div className="text-gray-400 text-6xl">📁</div>
                                <div>
                                    <p className="text-lg font-medium text-gray-700 mb-2">
                                        Drop your 3D model here
                                    </p>
                                    <p className="text-sm text-gray-500">
                                        Drag and drop a .stl or .3mf file, or click to browse
                                    </p>
                                </div>
                            </div>
                        )}
                    </label>

                    {file?.name.toLowerCase().endsWith(".stl") && (
                        <div className="bg-white p-4 rounded-lg border">
                            <label className="block text-sm font-medium text-gray-700 mb-2">
                                Custom Color
                            </label>
                            <div className="flex items-center gap-3">
                                <input
                                    type="color"
                                    value={customColor}
                                    onChange={(e) => setCustomColor(e.target.value)}
                                    className="w-12 h-12 p-1 border rounded-lg cursor-pointer"
                                />
                                <span className="text-sm text-gray-600 font-mono">
                                    {customColor.toUpperCase()}
                                </span>
                            </div>
                        </div>
                    )}

                    {uploadingPreview && (
                        <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                            <p className="text-sm text-blue-700 text-center">Generating preview...</p>
                        </div>
                    )}

                    {stlData && (
                        <div className="rounded-lg overflow-hidden border bg-gray-50 flex items-center justify-center" style={{ height: "300px" }}>
                            <Model3DPreview 
                                modelData={stlData} 
                                color={customColor} 
                                fileType={file?.name.toLowerCase().endsWith(".3mf") ? "3mf" : "stl"}
                                width={400}
                                height={300}
                            />
                        </div>
                    )}

                    <button
                        type="submit"
                        className="w-full bg-spooler-orange hover:bg-spooler-orange-light hover:cursor-pointer text-white py-3 px-6 rounded-lg transition text-lg font-medium disabled:opacity-50 disabled:cursor-not-allowed"
                        disabled={loading || !file}
                    >
                        {loading ? "Submitting..." : "Submit Print"}
                    </button>

                    {loading && (
                        <div className="w-full bg-white p-4 rounded-lg border">
                            <div className="flex justify-between text-sm text-gray-600 mb-2">
                                <span>Uploading...</span>
                                <span>{uploadProgress.toFixed(0)}%</span>
                            </div>
                            <div className="w-full bg-gray-200 rounded-full h-3">
                                <div 
                                    className="bg-spooler-orange h-3 rounded-full transition-all duration-300"
                                    style={{ width: `${uploadProgress}%` }}
                                />
                            </div>
                        </div>
                    )}

                    {error && (
                        <div className="bg-red-50 border border-red-200 rounded-lg p-4">
                            <div className="text-center text-red-700 font-medium">{error}</div>
                        </div>
                    )}
                </form>
            </div>
        </div>
    );
};