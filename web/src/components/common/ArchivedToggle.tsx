import React from 'react';

interface Props {
  showArchived: boolean;
  onChange: (value: boolean) => void;
  label?: string;
}

export const ArchivedToggle: React.FC<Props> = ({ showArchived, onChange, label = 'Afficher les archives' }) => (
  <label className="flex items-center gap-3 cursor-pointer select-none">
    <div className="relative">
      <input
        type="checkbox"
        checked={showArchived}
        onChange={(e) => onChange(e.target.checked)}
        className="sr-only peer"
      />
      <div className="w-11 h-6 bg-gray-200 peer-focus:ring-2 peer-focus:ring-indigo-500 rounded-full peer peer-checked:after:translate-x-full after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-indigo-600" />
    </div>
    <span className="text-sm font-medium text-gray-700">{label}</span>
  </label>
);