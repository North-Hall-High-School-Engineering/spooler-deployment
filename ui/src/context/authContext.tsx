import { createContext, useContext, useEffect, useState } from "react";
import { checkAuth } from "../util/auth";

interface AuthContextType {
  isAuthenticated: boolean;
  role: string;
  setIsAuthenticated: (isAuthenticated: boolean) => void;
  loading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const [role, setRole] = useState("user")
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        async function fetchAuthStatus() {
            const data = await checkAuth();
            if (data && typeof data === "object") {
                setIsAuthenticated(true);
                setRole
                setRole(data.role)
            } else {
                setIsAuthenticated(false);
            }
            setLoading(false);
        }
        fetchAuthStatus();
    }, []);

    return (
        <AuthContext.Provider value={{ isAuthenticated, setIsAuthenticated, role, loading }}>
            {children}
        </AuthContext.Provider>
    );
}

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (context === undefined) {
      throw new Error("useAuth must be used within an AuthProvider");
    }
    return context;
}