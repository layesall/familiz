import React, { useState, useEffect } from 'react';
import type { Contribution, Member } from '../../types';

interface ContributionsFormProps {
  initialData?: Contribution;
  onSubmit: (data: Omit<Contribution, 'id' | 'archived' | 'created_at'>) => void;
  onCancel: () => void;
  isEditing: boolean;
  members: Member[];
  defaultMemberId?: number;
}

type ContributionFormData = Omit<Contribution, 'id' | 'archived' | 'created_at'>;

export const ContributionsForm: React.FC<ContributionsFormProps> = ({
  initialData,
  onSubmit,
  onCancel,
  isEditing,
  members,
  defaultMemberId,
}) => {
  const [formData, setFormData] = useState<ContributionFormData>({
    member_id: 0,
    month: new Date().getMonth() + 1,
    year: new Date().getFullYear(),
    amount: 0,
    note: '',
  });

  useEffect(() => {
    if (initialData) {
      setFormData({
        member_id: initialData.member_id || 0,
        month: initialData.month || new Date().getMonth() + 1,
        year: initialData.year || new Date().getFullYear(),
        amount: initialData.amount || 0,
        note: initialData.note || '',
      });
    } else if (defaultMemberId) {
      setFormData(prev => ({
        ...prev,
        member_id: defaultMemberId,
      }));
    }
  }, [initialData, defaultMemberId]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: name === 'amount' || name === 'member_id' ? parseFloat(value) || 0 : value,
    }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(formData);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Membre</label>
          <select
            name="member_id"
            value={formData.member_id}
            onChange={handleChange}
            required
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
          >
            <option value="">Sélectionner un membre</option>
            {members.map(m => (
              <option key={m.id} value={m.id}>
                {m.first_name} {m.last_name}
              </option>
            ))}
          </select>
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Montant (0 = auto-calcul)</label>
          <input
            type="number"
            name="amount"
            value={formData.amount}
            onChange={handleChange}
            step="0.01"
            min="0"
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            placeholder="0.00"
          />
          <p className="text-xs text-gray-400 mt-1">Laisser 0 pour auto-calcul selon le statut</p>
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Mois</label>
          <select
            name="month"
            value={formData.month}
            onChange={handleChange}
            required
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
          >
            {Array.from({ length: 12 }, (_, i) => i + 1).map(m => (
              <option key={m} value={m}>
                {new Date(0, m - 1).toLocaleString('fr', { month: 'long' })}
              </option>
            ))}
          </select>
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Année</label>
          <input
            type="number"
            name="year"
            value={formData.year}
            onChange={handleChange}
            required
            min="2000"
            max="2100"
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
          />
        </div>
        <div className="md:col-span-2">
          <label className="block text-sm font-medium text-gray-700 mb-1">Note</label>
          <input
            type="text"
            name="note"
            value={formData.note}
            onChange={handleChange}
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            placeholder="Note optionnelle"
          />
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