import React, { useState, useEffect } from 'react';
import type { Member } from '../../types';

interface MemberFormProps {
  initialData?: Member;
  onSubmit: (data: Omit<Member, 'id' | 'archived'>) => void;
  onCancel: () => void;
  isEditing: boolean;
}

// On définit le type du formulaire (sans id ni archived)
type MemberFormData = Omit<Member, 'id' | 'archived'>;

export const MemberForm: React.FC<MemberFormProps> = ({ initialData, onSubmit, onCancel, isEditing }) => {
  // On type explicitement le state avec MemberFormData
  const [formData, setFormData] = useState<MemberFormData>({
    first_name: '',
    last_name: '',
    birth_date: '',
    marital_status: 'single', // maintenant TypeScript comprend que c'est une valeur possible du union
  });

  useEffect(() => {
    if (initialData) {
      setFormData({
        first_name: initialData.first_name || '',
        last_name: initialData.last_name || '',
        birth_date: initialData.birth_date || '',
        marital_status: initialData.marital_status || 'single',
      });
    }
  }, [initialData]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(formData);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Prénom</label>
          <input
            type="text"
            name="first_name"
            value={formData.first_name}
            onChange={handleChange}
            required
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            placeholder="Prénom"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Nom</label>
          <input
            type="text"
            name="last_name"
            value={formData.last_name}
            onChange={handleChange}
            required
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            placeholder="Nom"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Date de naissance</label>
          <input
            type="date"
            name="birth_date"
            value={formData.birth_date}
            onChange={handleChange}
            required
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
          />
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Situation familiale</label>
          <select
            name="marital_status"
            value={formData.marital_status}
            onChange={handleChange}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
          >
            <option value="single">Célibataire</option>
            <option value="married">Marié(e)</option>
            <option value="minor">Mineur(e)</option>
          </select>
        </div>
      </div>
      <div className="flex justify-end gap-3 pt-2">
        <button type="button" onClick={onCancel} className="bg-gray-200 hover:bg-gray-300 text-gray-800 px-4 py-2 rounded-lg font-semibold">
          Annuler
        </button>
        <button type="submit" className="bg-gray-700 hover:bg-gray-600 text-white px-4 py-2 rounded-lg font-semibold">
          {isEditing ? 'Mettre à jour' : 'Créer'}
        </button>
      </div>
    </form>
  );
};