import { useState } from 'react';
import { useLocation, useNavigate } from "react-router";
import { useAuth } from '../context/authContext';
import logo from "../assets/logo-1.svg"

export const Navbar = () => {
    const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
    const navigate = useNavigate();
    const location = useLocation();
    const { isAuthenticated, role, loading } = useAuth();

    const isActiveRoute = (path: string) => location.pathname === path;

    return (
        <nav className="w-full h-20 bg-white">
            <div className="max-w-screen-xl flex flex-wrap items-center justify-between mx-auto py-6 px-4">
                <a href="/" className="flex items-center space-x-3">
                    <img src={logo} className="h-8 w-32" alt="Spooler Logo" />
                </a>

                <div className="hidden md:flex space-x-8 font-medium min-w-0">
                    {loading ? (
                        <div className="flex space-x-8">
                            <div className="bg-gray-200 animate-pulse rounded h-5 w-20"></div>
                            <div className="bg-gray-200 animate-pulse rounded h-5 w-24"></div>
                        </div>
                    ) : (
                        <>
                            {isAuthenticated && (
                                <>
                                    <a href="/dashboard" className={isActiveRoute('/dashboard') ? 'text-spooler-orange' : 'text-black'}>Dashboard</a>
                                    {role === 'admin' && (
                                        <a href="/admin" className={isActiveRoute('/admin') ? 'text-spooler-orange' : 'text-black'}>Administrator</a>
                                    )}
                                </>
                            )}
                            <a href="#" className="text-black">Contact</a>
                        </>
                    )}
                </div>

                <div className="flex items-center space-x-4">
                    {loading ? (
                        <div className="bg-gray-200 animate-pulse py-2 px-6 rounded-sm">
                            <span className="invisible">New Print</span>
                        </div>
                    ) : isAuthenticated ? (
                        <>
                            <button
                                onClick={() => navigate('/new')}
                                className="bg-spooler-orange hover:bg-spooler-orange-light hover:cursor-pointer text-white py-2 px-6 rounded-sm flex items-center space-x-2"
                            >
                                <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24">
                                    <path stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 5v14m7-7H5" />
                                </svg>
                                <span>New Print</span>
                            </button>
                        </>
                    ) : (
                        <button
                            onClick={() => navigate('/login')}
                            className="bg-bambu-green hover:bg-bambu-green-light text-white py-2 px-6 rounded-md"
                        >
                            Login
                        </button>
                    )}

                    <button
                        onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
                        className="md:hidden p-2 w-10 h-10 text-gray-500 rounded-lg hover:bg-gray-100"
                    >
                        <svg className="w-5 h-5" fill="none" viewBox="0 0 17 14">
                            <path stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M1 1h15M1 7h15M1 13h15" />
                        </svg>
                    </button>
                </div>

                <div className={`w-full md:hidden ${isMobileMenuOpen ? 'block' : 'hidden'} bg-gray-50 mt-4`}>
                    <ul className="flex flex-col space-y-4 p-4">
                        {!loading && isAuthenticated && (
                            <>
                                <li><a href="/dashboard" onClick={() => setIsMobileMenuOpen(false)} className={isActiveRoute('/dashboard') ? 'text-spooler-orange' : 'text-black'}>Dashboard</a></li>
                                {role === 'admin' && (
                                    <li><a href="/admin" onClick={() => setIsMobileMenuOpen(false)} className={isActiveRoute('/admin') ? 'text-spooler-orange' : 'text-black'}>Administrator</a></li>
                                )}
                            </>
                        )}
                        <li><a href="#" onClick={() => setIsMobileMenuOpen(false)} className="text-black">Contact</a></li>
                    </ul>
                </div>
            </div>
        </nav>
    );
};