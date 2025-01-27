import {
  createContext,
  useContext,
  useState,
  useEffect,
  ReactNode,
} from "react";
import Cookies from "js-cookie";

interface AuthContextType {
  userEmail: string | null;
  setUserEmail: (email: string | null) => void;
  clearAuth: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [userEmail, setUserEmail] = useState<string | null>(null);

  // On mount, try to restore email from localStorage
  useEffect(() => {
    const storedEmail = localStorage.getItem("user_email");
    if (storedEmail) {
      setUserEmail(storedEmail);
    }
  }, []);

  const handleSetUserEmail = (email: string | null) => {
    setUserEmail(email);
    if (email) {
      localStorage.setItem("user_email", email);
    } else {
      localStorage.removeItem("user_email");
    }
  };

  const clearAuth = () => {
    setUserEmail(null);
    localStorage.removeItem("user_email");
    Cookies.remove("session_token");
  };

  return (
    <AuthContext.Provider
      value={{ userEmail, setUserEmail: handleSetUserEmail, clearAuth }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
