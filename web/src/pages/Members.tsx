import React, { useState, useEffect, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { membersAPI } from '../api/members';
import { MemberForm } from '../components/members/MemberForm';
import { MemberTable } from '../components/members/MemberTable';
import { ArchivedToggle } from '../components/common/ArchivedToggle';
import { LoadingSpinner } from '../components/common/LoadingSpinner';
import { ErrorMessage } from '../components/common/ErrorMessage';
import { ConfirmDialog } from '../components/common/ConfirmDialog';
import type { Member } from '../types';

export const Members: React.FC = () => {
  const navigate = useNavigate();
  const [members, setMembers] = useState<Member[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string>('');
  const [showArchived, setShowArchived] = useState(false);
  const [showModal, setShowModal] = useState(false);
  const [editingMember, setEditingMember] = useState<Member | null>(null);
  const [deleteTarget, setDeleteTarget] = useState<number | null>(null);

  const fetchMembers = useCallback(async () => {
    setLoading(true);
    setError('');
    try {
      const res = await membersAPI.getAll(showArchived);
      setMembers(res.data || []);
    } catch (err) {
      setError('Erreur lors du chargement des membres');
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, [showArchived]);

  useEffect(() => {
    fetchMembers();
  }, [fetchMembers]);

  const handleCreate = async (data: Omit<Member, 'id' | 'archived'>) => {
    try {
      await membersAPI.create(data);
      setShowModal(false);
      fetchMembers();
    } catch (err) {
      setError("Erreur lors de la création");
      console.error(err);
    }
  };

  const handleUpdate = async (data: Partial<Omit<Member, 'id' | 'archived'>>) => {
    if (!editingMember) return;
    try {
      await membersAPI.update(editingMember.id, data);
      setShowModal(false);
      setEditingMember(null);
      fetchMembers();
    } catch (err) {
      setError("Erreur lors de la mise à jour");
      console.error(err);
    }
  };

  const handleDelete = async () => {
    if (!deleteTarget) return;
    try {
      await membersAPI.delete(deleteTarget);
      setDeleteTarget(null);
      fetchMembers();
    } catch (err) {
      setError("Erreur lors de la suppression");
      console.error(err);
    }
  };

  const openModal = (member?: Member) => {
    setEditingMember(member || null);
    setShowModal(true);
  };

  const closeModal = () => {
    setShowModal(false);
    setEditingMember(null);
  };

  const handleAddTransaction = (memberId: number) => {
    navigate(`/transactions?member_id=${memberId}`);
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <h1 className="text-2xl md:text-3xl font-bold text-gray-900">👥 Membres</h1>
        <div className="flex items-center gap-3 flex-wrap">
          <ArchivedToggle showArchived={showArchived} onChange={setShowArchived} />
          <button
            onClick={() => openModal()}
            className="bg-gray-700 hover:bg-gray-600 text-white px-4 py-2 rounded-lg font-semibold text-sm whitespace-nowrap"
          >
            + Ajouter
          </button>
        </div>
      </div>

      {error && <ErrorMessage message={error} onRetry={fetchMembers} />}

      {loading ? (
        <LoadingSpinner />
      ) : members.length === 0 ? (
        <div className="bg-white rounded-xl shadow-sm p-6 border border-gray-100 text-center py-12">
          <p className="text-gray-400 text-lg">
            {showArchived ? 'Aucun membre archivé' : 'Aucun membre enregistré'}
          </p>
          <button
            onClick={() => openModal()}
            className="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-lg font-semibold mt-4"
          >
            Ajouter le premier membre
          </button>
        </div>
      ) : (
        <MemberTable members={members} onAddTransaction={handleAddTransaction} />
      )}

      {/* Modal d'ajout / édition */}
      {showModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm">
          <div className="bg-white rounded-xl shadow-xl max-w-2xl w-full mx-4 p-6 max-h-[90vh] overflow-y-auto">
            <h2 className="text-lg font-bold text-gray-900 mb-4">
              {editingMember ? 'Modifier le membre' : 'Nouveau membre'}
            </h2>
            <MemberForm
              initialData={editingMember || undefined}
              onSubmit={editingMember ? handleUpdate : handleCreate}
              onCancel={closeModal}
              isEditing={!!editingMember}
            />
          </div>
        </div>
      )}

      <ConfirmDialog
        isOpen={!!deleteTarget}
        title="Supprimer le membre"
        message="Êtes-vous sûr de vouloir supprimer ce membre ? Cette action est irréversible."
        onConfirm={handleDelete}
        onCancel={() => setDeleteTarget(null)}
        confirmLabel="Supprimer"
      />
    </div>
  );
};