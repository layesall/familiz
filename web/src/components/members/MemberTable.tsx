import React from 'react';
import { useNavigate } from 'react-router-dom';
import type { Member } from '../../types';
import { formatDate, getMaritalStatusLabel, getMaritalStatusBadge } from '../../utils/helpers';
import { ActionIcons } from '../common/ActionIcons';

interface MemberTableProps {
  members: Member[];
  onAddTransaction: (memberId: number) => void;
}

export const MemberTable: React.FC<MemberTableProps> = ({ members, onAddTransaction }) => {
  const navigate = useNavigate();

  return (
    <div className="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
      <div className="hidden md:block overflow-x-auto">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              <th className="text-left py-3 px-4 font-medium text-gray-500">ID</th>
              <th className="text-left py-3 px-4 font-medium text-gray-500">Prénom</th>
              <th className="text-left py-3 px-4 font-medium text-gray-500">Nom</th>
              <th className="text-left py-3 px-4 font-medium text-gray-500">Statut</th>
              <th className="text-left py-3 px-4 font-medium text-gray-500">Naissance</th>
              <th className="text-left py-3 px-4 font-medium text-gray-500">État</th>
              <th className="text-right py-3 px-4 font-medium text-gray-500">Actions</th>
            </tr>
          </thead>
          <tbody>
            {members.map((m) => (
              <tr key={m.id} className="border-b border-gray-100 hover:bg-gray-50 transition-colors">
                <td className="py-3 px-4 font-mono text-xs text-gray-500">{m.id}</td>
                <td className="py-3 px-4 font-medium text-gray-900">{m.first_name}</td>
                <td className="py-3 px-4 font-medium text-gray-900">{m.last_name}</td>
                <td className="py-3 px-4">
                  <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getMaritalStatusBadge(m.marital_status)}`}>
                    {getMaritalStatusLabel(m.marital_status)}
                  </span>
                </td>
                <td className="py-3 px-4 text-gray-600">{formatDate(m.birth_date)}</td>
                <td className="py-3 px-4">
                  {m.archived ? (
                    <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-200 text-gray-700">Archivé</span>
                  ) : (
                    <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">Actif</span>
                  )}
                </td>
                <td className="py-3 px-4 text-right">
                  <div className="flex justify-end items-center gap-1">
                    <ActionIcons
                      actions={[
                        { icon: 'view', onClick: () => navigate(`/members/${m.id}`), label: 'Voir le profil', color: 'indigo' },
                        { icon: 'add', onClick: () => onAddTransaction(m.id), label: 'Ajouter une transaction', color: 'emerald' },
                      ]}
                    />
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Version mobile */}
      <div className="md:hidden divide-y divide-gray-100">
        {members.map((m) => (
          <div key={m.id} className="p-4 hover:bg-gray-50 transition-colors">
            <div className="flex items-start justify-between">
              <div>
                <p className="font-medium text-gray-900">{m.first_name} {m.last_name}</p>
                <p className="text-sm text-gray-500">ID: {m.id}</p>
                <div className="flex items-center gap-2 mt-1">
                  <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${getMaritalStatusBadge(m.marital_status)}`}>
                    {getMaritalStatusLabel(m.marital_status)}
                  </span>
                  {m.archived && (
                    <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-gray-200 text-gray-700">Archivé</span>
                  )}
                </div>
              </div>
              <div className="flex items-center gap-1">
                <ActionIcons
                  actions={[
                    { icon: 'view', onClick: () => navigate(`/members/${m.id}`), label: 'Voir le profil', color: 'indigo' },
                    { icon: 'add', onClick: () => onAddTransaction(m.id), label: 'Ajouter une transaction', color: 'emerald' },
                  ]}
                  size="md"
                />
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};