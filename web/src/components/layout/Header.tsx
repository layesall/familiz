import React from 'react';
import { useLocation } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';
import { useNavigate } from 'react-router-dom';

interface HeaderProps {
  onMenuClick: () => void;
}

// Mapping des routes vers des titres lisibles
const pageTitles: Record<string, string> = {
  '/': 'Dashboard',
  '/members': 'Membres',
  '/contributions': 'Contributions',
  '/events': 'Événements',
  '/settings': 'Paramètres',
};

export const Header: React.FC<HeaderProps> = ({ onMenuClick }) => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  // Récupérer le nom de la page active
  const currentPage = pageTitles[location.pathname] || 'Dashboard';
  
  // Récupérer le nom de l'admin à partir de l'email
  const adminName = user?.email ? user.email.split('@')[0] : 'Admin';
  const displayName = adminName.charAt(0).toUpperCase() + adminName.slice(1);

  return (
    <header className="bg-white border-b border-gray-200 px-4 py-3 md:px-6 md:py-4 flex items-center justify-between shadow-sm">
      <div className="flex items-center gap-3">
        <button
          onClick={onMenuClick}
          className="lg:hidden text-gray-600 hover:text-gray-900"
          aria-label="Menu"
        >
          <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
          </svg>
        </button>
        <h2 className="text-lg font-semibold text-gray-700 hidden sm:block">
          {currentPage}
        </h2>
      </div>

      <div className="flex items-center gap-3 md:gap-4">
        <div className="flex items-center gap-2">
          <div className="w-8 h-8 rounded-full bg-indigo-100 flex items-center justify-center text-indigo-600 font-semibold text-sm">
            {displayName.charAt(0)}
          </div>
          <div className="hidden sm:block">
            <p className="text-sm font-medium text-gray-700">{displayName}</p>
            <p className="text-xs text-gray-400 capitalize">{user?.role || 'Admin'}</p>
          </div>
        </div>
        <button
          onClick={handleLogout}
          className="text-sm text-gray-600 hover:text-red-600 transition-colors"
        >
          Déconnexion
        </button>
      </div>
    </header>
  );
};