import React from 'react';

interface Props {
  message: string | null;
  onRetry?: () => void;
}

export const ErrorMessage: React.FC<Props> = ({ message, onRetry }) => {
  if (!message) return null;
  return (
    <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg flex items-center justify-between">
      <span>{message}</span>
      {onRetry && (
        <button onClick={onRetry} className="text-red-600 hover:text-red-800 font-medium text-sm">
          Réessayer
        </button>
      )}
    </div>
  );
};