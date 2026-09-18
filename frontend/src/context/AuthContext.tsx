import React, {
  createContext,
  useContext,
  useEffect,
  useState,
} from 'react';

import { User } from '../api/auth';

interface AuthContextType {
  user: User | null;
  login: (user: User, token: string) => void;
  logout: () => void;
  isAuthenticated: boolean;
  isLoading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{
  children: React.ReactNode;
}> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  /**
   * Restore authentication session
   */
  useEffect(() => {
    const restoreSession = () => {
      try {
        const token = localStorage.getItem('token');
        const savedUser = localStorage.getItem('user');

        if (token && savedUser) {
          const parsedUser: User = JSON.parse(savedUser);

          setUser(parsedUser);
          setIsAuthenticated(true);
        }
      } catch (error) {
        console.error(
          'Erreur lors de la restauration de la session:',
          error
        );

        // Nettoyage si les données sont invalides
        localStorage.removeItem('token');
        localStorage.removeItem('user');

        setUser(null);
        setIsAuthenticated(false);
      } finally {
        setIsLoading(false);
      }
    };

    restoreSession();
  }, []);

  /**
   * Login
   */
  const login = (user: User, token: string) => {
    localStorage.setItem('token', token);
    localStorage.setItem('user', JSON.stringify(user));

    setUser(user);
    setIsAuthenticated(true);
  };

  /**
   * Logout
   */
  const logout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');

    setUser(null);
    setIsAuthenticated(false);
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        login,
        logout,
        isAuthenticated,
        isLoading,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

/**
 * Authentication hook
 */
export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);

  if (!context) {
    throw new Error(
      'useAuth must be used within an AuthProvider'
    );
  }

  return context;
};
