import React, { useState, useEffect } from 'react';
import type { Event, Member } from '../../types';

interface EventFormProps {
  initialData?: Event;
  onSubmit: (data: Omit<Event, 'id' | 'archived' | 'created_at'>) => void;
  onCancel: () => void;
  isEditing: boolean;
  members: Member[];
  defaultMemberId?: number;
}

type EventFormData = Omit<Event, 'id' | 'archived' | 'created_at'>;

export const EventForm: React.FC<EventFormProps> = ({
  initialData,
  onSubmit,
  onCancel,
  isEditing,
  members,
  defaultMemberId,
}) => {
  const [formData, setFormData] = useState<EventFormData>({
    member_id: 0,
    type: 'wedding',
    amount_received: 0,
    event_date: new Date().toISOString().split('T')[0],
  });

  useEffect(() => {
    if (initialData) {
      setFormData({
        member_id: initialData.member_id || 0,
        type: initialData.type || 'wedding',
        amount_received: initialData.amount_received || 0,
        event_date: initialData.event_date || new Date().toISOString().split('T')[0],
      });
    } else if (defaultMemberId) {
      setFormData(prev => ({
        ...prev,
        member_id: defaultMemberId,
      }));
    }
  }, [initialData, defaultMemberId]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: name === 'amount_received' || name === 'member_id' ? parseFloat(value) || 0 : value,
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
          <label className="block text-sm font-medium text-gray-700 mb-1">Type d'événement</label>
          <select
            name="type"
            value={formData.type}
            onChange={handleChange}
            required
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
          >
            <option value="wedding">Mariage</option>
            <option value="baptism">Baptême</option>
          </select>
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Montant reçu (0 = auto-calcul)</label>
          <input
            type="number"
            name="amount_received"
            value={formData.amount_received}
            onChange={handleChange}
            step="0.01"
            min="0"
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            placeholder="0.00"
          />
          <p className="text-xs text-gray-400 mt-1">Laisser 0 pour auto-calcul</p>
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Date de l'événement</label>
          <input
            type="date"
            name="event_date"
            value={formData.event_date}
            onChange={handleChange}
            required
            className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
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