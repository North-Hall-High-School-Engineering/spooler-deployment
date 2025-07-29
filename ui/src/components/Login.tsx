import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { requestOTP, verifyOTP } from "../util/auth";
import { useAuth } from "../context/authContext";

function Login() {
    const [step, setStep] = useState<"email" | "otp" | "done">("email");
    const [email, setEmail] = useState("");
    const [otp, setOtp] = useState("");
    const [error, setError] = useState("");
    const auth = useAuth();
    const navigate = useNavigate();

    useEffect(() => {
        if (!auth.loading && auth.isAuthenticated) {
            navigate("/dashboard", { replace: true });
        }
    }, [auth.isAuthenticated, auth.loading, navigate]);

    useEffect(() => {
        const savedStep = localStorage.getItem("spooler_login_step");
        const savedEmail = localStorage.getItem("spooler_login_email");
        if (savedStep === "otp" && savedEmail) {
            setStep("otp");
            setEmail(savedEmail);
        }
    }, []);

    useEffect(() => {
        localStorage.setItem("spooler_login_step", step);
        if (email) localStorage.setItem("spooler_login_email", email);
    }, [step, email]);

    const handleRequestOTP = async (e: React.FormEvent) => {
        e.preventDefault();
        setError("");
        try {
            await requestOTP(email);
            setStep("otp");
        } catch (err: any) {
            setError(err.response?.data?.error || "Failed to send OTP");
        }
    };

    const handleVerify = async (e: React.FormEvent) => {
        e.preventDefault();
        setError("");
        try {
            await verifyOTP(email, otp);
            setStep("done");
            auth.setIsAuthenticated(true); // update context
            navigate("/dashboard", { replace: true });
        } catch (err: any) {
            setError(err.response?.data?.error || "OTP verification failed");
        }
    };

    if (auth.loading) return null;

    return (
        <div className="min-h-screen flex flex-col items-center justify-center bg-gray-50">
            <form
                className="w-full max-w-sm space-y-4"
                onSubmit={step === "email" ? handleRequestOTP : handleVerify}
            >
                <h2 className="text-xl font-semibold text-gray-800 text-center">
                    {step === "email" && "Login"}
                    {step === "otp" && "Enter OTP"}
                </h2>

                {step === "email" && (
                    <>
                        <div>
                            <label className="block text-sm text-gray-700">Email</label>
                            <input
                                type="email"
                                placeholder="Email"
                                className="mt-1 w-full border px-3 py-2 rounded border-gray-300 focus:outline-none focus:ring focus:border-spooler-orange"
                                value={email}
                                onChange={e => setEmail(e.target.value)}
                                required
                            />
                        </div>
                        <button
                            type="submit"
                            className="w-full bg-spooler-orange hover:bg-spooler-orange-light hover:cursor-pointer text-white py-2 rounded transition"
                        >
                            Request Verification Code
                        </button>
                    </>
                )}

                {step === "otp" && (
                    <>
                        <div>
                            <label className="block text-sm text-gray-700">Verification Code</label>
                            <input
                                type="text"
                                placeholder="Enter Verification Code"
                                className="mt-1 w-full border px-3 py-2 rounded border-gray-300 focus:outline-none focus:ring focus:border-spooler-orange"
                                value={otp}
                                onChange={e => setOtp(e.target.value)}
                                required
                            />
                        </div>
                        <button
                            type="submit"
                            className="w-full bg-spooler-orange hover:bg-spooler-orange-light hover:cursor-pointer text-white py-2 rounded transition"
                        >
                            Verify Code
                        </button>
                    </>
                )}

                {error && (
                    <div className="text-center text-red-600">{error}</div>
                )}
            </form>

            <p className="mt-2 text-sm">
                Don't have an account?{" "}
                <a href="/register" className="text-spooler-orange hover:underline hover:cursor-pointer">register</a>
            </p>
        </div>
    );
}

export default Login;