import React from 'react';
import type { Event, Member } from '../../types';
import { formatCurrency, formatDate, getEventTypeLabel } from '../../utils/helpers';
import { ActionIcons } from '../common/ActionIcons';

interface EventTableProps {
  events: Event[];
  members: Member[];
  onEdit: (evt: Event) => void;
  onDelete: (id: number) => void;
  onView: (memberId: number) => void;
}

export const EventTable: React.FC<EventTableProps> = ({
  events,
  members,
  onEdit,
  onDelete,
  onView,
}) => {
  const getMemberName = (memberId: number) => {
    const m = members.find(m => m.id === memberId);
    return m ? `${m.first_name} ${m.last_name}` : `ID: ${memberId}`;
  };

  return (
    <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
      <div className="hidden md:block overflow-x-auto">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              <th className="text-left py-3 px-4 font-medium text-gray-500">Membre</th>
              <th className="text-left py-3 px-4 font-medium text-gray-500">Type</th>
              <th className="text-left py-3 px-4 font-medium text-gray-500">Date</th>
              <th className="text-right py-3 px-4 font-medium text-gray-500">Montant reçu</th>
              <th className="text-left py-3 px-4 font-medium text-gray-500">Statut</th>
              <th className="text-right py-3 px-4 font-medium text-gray-500">Actions</th>
            </tr>
          </thead>
          <tbody>
            {events.map((evt) => (
              <tr key={evt.id} className="border-b border-gray-100 hover:bg-gray-50 transition-colors">
                <td className="py-3 px-4 font-medium text-gray-900">{getMemberName(evt.member_id)}</td>
                <td className="py-3 px-4">
                  <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${evt.type === 'wedding' ? 'bg-pink-100 text-pink-800' : 'bg-cyan-100 text-cyan-800'}`}>
                    {getEventTypeLabel(evt.type)}
                  </span>
                </td>
                <td className="py-3 px-4">{formatDate(evt.event_date)}</td>
                <td className="py-3 px-4 text-right font-medium">{formatCurrency(evt.amount_received)}</td>
                <td className="py-3 px-4">
                  {evt.archived ? (
                    <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-200 text-gray-700">Archivé</span>
                  ) : (
                    <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">Actif</span>
                  )}
                </td>
                <td className="py-3 px-4 text-right">
                  <ActionIcons
                    actions={[
                      { icon: 'view', onClick: () => onView(evt.member_id), label: 'Voir le membre', color: 'indigo' },
                      { icon: 'edit', onClick: () => onEdit(evt), label: 'Modifier', color: 'blue' },
                      { icon: 'delete', onClick: () => onDelete(evt.id), label: 'Supprimer', color: 'red' },
                    ]}
                  />
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Version mobile */}
      <div className="md:hidden divide-y divide-gray-100">
        {events.map((evt) => (
          <div key={evt.id} className="p-4 hover:bg-gray-50 transition-colors">
            <div className="flex items-start justify-between">
              <div>
                <p className="font-medium text-gray-900">{getMemberName(evt.member_id)}</p>
                <p className="text-sm text-gray-500">{getEventTypeLabel(evt.type)}</p>
                <p className="text-sm text-gray-500">{formatDate(evt.event_date)}</p>
                <p className="text-sm font-medium text-purple-600">{formatCurrency(evt.amount_received)}</p>
                <span className="text-xs text-gray-500">{evt.archived ? 'Archivé' : 'Actif'}</span>
              </div>
              <ActionIcons
                actions={[
                  { icon: 'view', onClick: () => onView(evt.member_id), label: 'Voir le membre', color: 'indigo' },
                  { icon: 'edit', onClick: () => onEdit(evt), label: 'Modifier', color: 'blue' },
                  { icon: 'delete', onClick: () => onDelete(evt.id), label: 'Supprimer', color: 'red' },
                ]}
                size="md"
              />
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};