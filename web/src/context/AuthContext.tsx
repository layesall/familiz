import React, { createContext, useState, useEffect, useCallback, type ReactNode } from 'react';
import type { User } from '../types';

interface AuthContextType {
  token: string | null;
  user: User | null;
  login: (token: string, user: User) => void;
  logout: () => void;
  isAuthenticated: boolean;
  isAdmin: boolean;
}

export const AuthContext = createContext<AuthContextType | null>(null);

export const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [token, setToken] = useState<string | null>(localStorage.getItem('familiz_token'));
  const [user, setUser] = useState<User | null>(() => {
    const saved = localStorage.getItem('familiz_user');
    return saved ? JSON.parse(saved) : null;
  });

  useEffect(() => {
    if (token) localStorage.setItem('familiz_token', token);
    else localStorage.removeItem('familiz_token');
  }, [token]);

  useEffect(() => {
    if (user) localStorage.setItem('familiz_user', JSON.stringify(user));
    else localStorage.removeItem('familiz_user');
  }, [user]);

  const login = useCallback((newToken: string, newUser: User) => {
    setToken(newToken);
    setUser(newUser);
  }, []);

  const logout = useCallback(() => {
    setToken(null);
    setUser(null);
    localStorage.removeItem('familiz_token');
    localStorage.removeItem('familiz_user');
  }, []);

  const value = {
    token,
    user,
    login,
    logout,
    isAuthenticated: !!token,
    isAdmin: user?.role === 'admin',
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};