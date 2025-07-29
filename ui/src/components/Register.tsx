import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { register, requestOTP, verifyOTP } from "../util/auth";
import { useAuth } from "../context/authContext";

function Register() {
    const [step, setStep] = useState<"register" | "otp" | "done">("register");
    const [firstName, setFirstName] = useState("");
    const [lastName, setLastName] = useState("");
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
        const savedStep = localStorage.getItem("spooler_step");
        const savedEmail = localStorage.getItem("spooler_email");
        if (savedStep === "otp" && savedEmail) {
            setStep("otp");
            setEmail(savedEmail);
        }
    }, []);

    useEffect(() => {
        localStorage.setItem("spooler_step", step);
        if (email) localStorage.setItem("spooler_email", email);
    }, [step, email]);

    const handleRegister = async (e: React.FormEvent) => {
        e.preventDefault();
        setError("");
        try {
            await register(email, firstName, lastName);
            await requestOTP(email);
            setStep("otp");
        } catch (err: any) {
            setError(err.response?.data?.error || "Registration failed");
        }
    };

    const handleVerify = async (e: React.FormEvent) => {
        e.preventDefault();
        setError("");
        try {
            await verifyOTP(email, otp);
            setStep("done");
            auth.setIsAuthenticated(true);
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
                onSubmit={step === "register" ? handleRegister : handleVerify}
            >
                <h2 className="text-xl font-semibold text-gray-800 text-center">
                    {step === "register" && "Register"}
                    {step === "otp" && "Enter OTP"}
                </h2>

                {step === "register" && (
                    <>
                        <div className="flex gap-x-2">
                            <div>
                                <label className="block text-sm text-gray-700">First Name</label>
                                <input
                                    type="text"
                                    placeholder="First Name"
                                    className="mt-1 w-full border px-3 py-2 rounded border-gray-300 focus:outline-none focus:ring focus:border-spooler-orange"
                                    value={firstName}
                                    onChange={e => setFirstName(e.target.value)}
                                    required
                                />
                            </div>
                            <div>
                                <label className="block text-sm text-gray-700">Last Name</label>
                                <input
                                    type="text"
                                    placeholder="Last Name"
                                    className="mt-1 w-full border px-3 py-2 rounded border-gray-300 focus:outline-none focus:ring focus:border-spooler-orange"
                                    value={lastName}
                                    onChange={e => setLastName(e.target.value)}
                                    required
                                />
                            </div>
                        </div>
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
                            Register
                        </button>
                    </>
                )}

                {step === "otp" && (
                    <>
                        <div>
                            <label className="block text-sm text-gray-700">OTP Code</label>
                            <input
                                type="text"
                                placeholder="Enter OTP"
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
                            Verify OTP
                        </button>
                    </>
                )}

                {error && (
                    <div className="text-center text-red-600">{error}</div>
                )}
            </form>

            <p className="mt-2 text-sm">
                Already have an account?{" "}
                <a href="/login" className="text-spooler-orange hover:underline hover:cursor-pointer">login</a>
            </p>
        </div>
    );
}

export default Register;