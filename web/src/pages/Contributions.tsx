import React, { useState, useEffect, useCallback } from 'react';
import { useSearchParams } from 'react-router-dom';
import { contributionsAPI } from '../api/contributions';
import { membersAPI } from '../api/members';
import { ContributionsForm } from '../components/contributions/ContributionsForm';
import { ContributionsTable } from '../components/contributions/ContributionsTable';
import { ArchivedToggle } from '../components/common/ArchivedToggle';
import { LoadingSpinner } from '../components/common/LoadingSpinner';
import { ErrorMessage } from '../components/common/ErrorMessage';
import { ConfirmDialog } from '../components/common/ConfirmDialog';
import { formatCurrency } from '../utils/helpers';
import type { Contribution, Member } from '../types';

export const Contributions: React.FC = () => {
  const [searchParams] = useSearchParams();
  const [contributions, setContributions] = useState<Contribution[]>([]);
  const [members, setMembers] = useState<Member[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string>('');
  const [showArchived, setShowArchived] = useState(false);
  const [showModal, setShowModal] = useState(false);
  const [editingTx, setEditingTx] = useState<Contribution | null>(null);
  const [deleteTarget, setDeleteTarget] = useState<number | null>(null);
  const [filterMember, setFilterMember] = useState<string>('');

  const memberIdFromUrl = searchParams.get('member_id') || '';

  useEffect(() => {
    if (memberIdFromUrl) {
      setFilterMember(memberIdFromUrl);
    }
  }, [memberIdFromUrl]);

  const fetchData = useCallback(async () => {
    setLoading(true);
    setError('');
    try {
      const [txRes, membersRes] = await Promise.all([
        contributionsAPI.getAll(filterMember ? parseInt(filterMember) : undefined, showArchived),
        membersAPI.getAll(),
      ]);
      setContributions(txRes.data || []);
      setMembers(membersRes.data || []);
    } catch (err) {
      setError('Erreur lors du chargement des contributions');
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, [showArchived, filterMember]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const handleCreate = async (data: Omit<Contribution, 'id' | 'archived' | 'created_at'>) => {
    try {
      await contributionsAPI.create(data);
      setShowModal(false);
      fetchData();
    } catch (err) {
      setError("Erreur lors de la création");
      console.error(err);
    }
  };

  const handleUpdate = async (data: Partial<Omit<Contribution, 'id' | 'archived' | 'created_at'>>) => {
    if (!editingTx) return;
    try {
      await contributionsAPI.update(editingTx.id, data);
      setShowModal(false);
      setEditingTx(null);
      fetchData();
    } catch (err) {
      setError("Erreur lors de la mise à jour");
      console.error(err);
    }
  };

  const handleDelete = async () => {
    if (!deleteTarget) return;
    try {
      await contributionsAPI.delete(deleteTarget);
      setDeleteTarget(null);
      fetchData();
    } catch (err) {
      setError("Erreur lors de la suppression");
      console.error(err);
    }
  };

  const openModal = (tx?: Contribution) => {
    setEditingTx(tx || null);
    setShowModal(true);
  };

  const closeModal = () => {
    setShowModal(false);
    setEditingTx(null);
  };

  const totalAmount = contributions.reduce((sum, t) => sum + (t.amount || 0), 0);

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <h1 className="text-2xl md:text-3xl font-bold text-gray-900">💰 Contributions</h1>
        <div className="flex items-center gap-3 flex-wrap">
          <select
            value={filterMember}
            onChange={(e) => setFilterMember(e.target.value)}
            className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 bg-white text-sm"
          >
            <option value="">Tous les membres</option>
            {members.map(m => (
              <option key={m.id} value={m.id}>
                {m.first_name} {m.last_name}
              </option>
            ))}
          </select>
          <ArchivedToggle showArchived={showArchived} onChange={setShowArchived} />
          <button
            onClick={() => openModal()}
            className="bg-gray-700 hover:bg-gray-600 text-white px-4 py-2 rounded-lg font-semibold text-sm whitespace-nowrap"
          >
            + Ajouter
          </button>
        </div>
      </div>

      <div className="bg-white rounded-xl shadow-sm p-4 border border-gray-100">
        <p className="text-4xl font-bold text-green-600">{formatCurrency(totalAmount)}</p>
        <p className="text-sm text-gray-500">Total des contributions</p>
      </div>

      {error && <ErrorMessage message={error} onRetry={fetchData} />}

      {loading ? (
        <LoadingSpinner />
      ) : contributions.length === 0 ? (
        <div className="bg-white rounded-xl shadow-sm p-6 border border-gray-100 text-center py-12">
          <p className="text-gray-400 text-lg">
            {showArchived ? 'Aucune contribution archivée' : 'Aucune contribution'}
          </p>
        </div>
      ) : (
        <ContributionsTable
          Contributions={contributions}
          members={members}
          onEdit={openModal}
          onDelete={(id) => setDeleteTarget(id)}
          onView={(memberId) => window.location.href = `/members/${memberId}`}
        />
      )}

      {showModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm">
          <div className="bg-white rounded-xl shadow-xl max-w-2xl w-full mx-4 p-6 max-h-[90vh] overflow-y-auto">
            <h2 className="text-lg font-bold text-gray-900 mb-4">
              {editingTx ? 'Modifier la contribution' : 'Nouvelle contribution'}
            </h2>
            <ContributionsForm
              initialData={editingTx || undefined}
              onSubmit={editingTx ? handleUpdate : handleCreate}
              onCancel={closeModal}
              isEditing={!!editingTx}
              members={members}
            />
          </div>
        </div>
      )}

      <ConfirmDialog
        isOpen={!!deleteTarget}
        title="Supprimer la contribution"
        message="Êtes-vous sûr de vouloir supprimer cette contribution ?"
        onConfirm={handleDelete}
        onCancel={() => setDeleteTarget(null)}
        confirmLabel="Supprimer"
      />
    </div>
  );
};