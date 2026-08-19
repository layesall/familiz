import React from 'react';

interface Props {
  size?: 'sm' | 'md' | 'lg';
}

const sizes = {
  sm: 'h-6 w-6',
  md: 'h-10 w-10',
  lg: 'h-16 w-16',
};

export const LoadingSpinner: React.FC<Props> = ({ size = 'md' }) => (
  <div className="flex items-center justify-center p-4">
    <div className={`${sizes[size]} animate-spin rounded-full border-4 border-indigo-200 border-t-indigo-600`} />
  </div>
);