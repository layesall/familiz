import React from 'react';
import type { Contribution, Member } from '../../types';
import { formatCurrency, getMonthLabel } from '../../utils/helpers';
import { ActionIcons } from '../common/ActionIcons';

interface ContributionsTableProps {
  Contributions: Contribution[];
  members: Member[];
  onEdit: (tx: Contribution) => void;
  onDelete: (id: number) => void;
  onView: (memberId: number) => void;
}

export const ContributionsTable: React.FC<ContributionsTableProps> = ({
  Contributions,
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
              <th className="text-left py-3 px-4 font-medium text-gray-500">Mois</th>
              <th className="text-left py-3 px-4 font-medium text-gray-500">Année</th>
              <th className="text-right py-3 px-4 font-medium text-gray-500">Montant</th>
              <th className="text-left py-3 px-4 font-medium text-gray-500">Note</th>
              <th className="text-left py-3 px-4 font-medium text-gray-500">Statut</th>
              <th className="text-right py-3 px-4 font-medium text-gray-500">Actions</th>
            </tr>
          </thead>
          <tbody>
            {Contributions.map((tx) => (
              <tr key={tx.id} className="border-b border-gray-100 hover:bg-gray-50 transition-colors">
                <td className="py-3 px-4 font-medium text-gray-900">{getMemberName(tx.member_id)}</td>
                <td className="py-3 px-4">{getMonthLabel(tx.month)}</td>
                <td className="py-3 px-4">{tx.year}</td>
                <td className="py-3 px-4 text-right font-medium">{formatCurrency(tx.amount)}</td>
                <td className="py-3 px-4 text-gray-500">{tx.note || '—'}</td>
                <td className="py-3 px-4">
                  {tx.archived ? (
                    <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-200 text-gray-700">Archivé</span>
                  ) : (
                    <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">Actif</span>
                  )}
                </td>
                <td className="py-3 px-4 text-right">
                  <ActionIcons
                    actions={[
                      { icon: 'view', onClick: () => onView(tx.member_id), label: 'Voir le membre', color: 'indigo' },
                      { icon: 'edit', onClick: () => onEdit(tx), label: 'Modifier', color: 'blue' },
                      { icon: 'delete', onClick: () => onDelete(tx.id), label: 'Supprimer', color: 'red' },
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
        {Contributions.map((tx) => (
          <div key={tx.id} className="p-4 hover:bg-gray-50 transition-colors">
            <div className="flex items-start justify-between">
              <div>
                <p className="font-medium text-gray-900">{getMemberName(tx.member_id)}</p>
                <p className="text-sm text-gray-500">
                  {getMonthLabel(tx.month)} {tx.year}
                </p>
                <p className="text-sm font-medium text-gray-900">{formatCurrency(tx.amount)}</p>
                <p className="text-xs text-gray-400">{tx.note || '—'}</p>
                <span className="text-xs text-gray-500">
                  {tx.archived ? 'Archivé' : 'Actif'}
                </span>
              </div>
              <ActionIcons
                actions={[
                  { icon: 'view', onClick: () => onView(tx.member_id), label: 'Voir le membre', color: 'indigo' },
                  { icon: 'edit', onClick: () => onEdit(tx), label: 'Modifier', color: 'blue' },
                  { icon: 'delete', onClick: () => onDelete(tx.id), label: 'Supprimer', color: 'red' },
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