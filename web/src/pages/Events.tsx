import React, { useState, useEffect, useCallback } from 'react';
import { useSearchParams } from 'react-router-dom';
import { eventsAPI } from '../api/events';
import { membersAPI } from '../api/members';
import { EventForm } from '../components/events/EventForm';
import { EventTable } from '../components/events/EventTable';
import { ArchivedToggle } from '../components/common/ArchivedToggle';
import { LoadingSpinner } from '../components/common/LoadingSpinner';
import { ErrorMessage } from '../components/common/ErrorMessage';
import { ConfirmDialog } from '../components/common/ConfirmDialog';
import { formatCurrency } from '../utils/helpers';
import type { Event, Member } from '../types';

export const Events: React.FC = () => {
  const [searchParams] = useSearchParams();
  const [events, setEvents] = useState<Event[]>([]);
  const [members, setMembers] = useState<Member[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string>('');
  const [showArchived, setShowArchived] = useState(false);
  const [showModal, setShowModal] = useState(false);
  const [editingEvent, setEditingEvent] = useState<Event | null>(null);
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
      const [evtRes, membersRes] = await Promise.all([
        eventsAPI.getAll(filterMember ? parseInt(filterMember) : undefined, showArchived),
        membersAPI.getAll(),
      ]);
      setEvents(evtRes.data || []);
      setMembers(membersRes.data || []);
    } catch (err) {
      setError('Erreur lors du chargement des événements');
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, [showArchived, filterMember]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const handleCreate = async (data: Omit<Event, 'id' | 'archived' | 'created_at'>) => {
    try {
      await eventsAPI.create(data);
      setShowModal(false);
      fetchData();
    } catch (err) {
      setError("Erreur lors de la création");
      console.error(err);
    }
  };

  const handleUpdate = async (data: Partial<Omit<Event, 'id' | 'archived' | 'created_at'>>) => {
    if (!editingEvent) return;
    try {
      await eventsAPI.update(editingEvent.id, data);
      setShowModal(false);
      setEditingEvent(null);
      fetchData();
    } catch (err) {
      setError("Erreur lors de la mise à jour");
      console.error(err);
    }
  };

  const handleDelete = async () => {
    if (!deleteTarget) return;
    try {
      await eventsAPI.delete(deleteTarget);
      setDeleteTarget(null);
      fetchData();
    } catch (err) {
      setError("Erreur lors de la suppression");
      console.error(err);
    }
  };

  const openModal = (evt?: Event) => {
    setEditingEvent(evt || null);
    setShowModal(true);
  };

  const closeModal = () => {
    setShowModal(false);
    setEditingEvent(null);
  };

  const totalReceived = events.reduce((sum, e) => sum + (e.amount_received || 0), 0);

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <h1 className="text-2xl md:text-3xl font-bold text-gray-900">🎉 Événements</h1>
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
        <p className="text-4xl font-bold text-red-600">{formatCurrency(totalReceived)}</p>
        <p className="text-sm text-gray-500">Total sortie des événements</p>
      </div>

      {error && <ErrorMessage message={error} onRetry={fetchData} />}

      {loading ? (
        <LoadingSpinner />
      ) : events.length === 0 ? (
        <div className="bg-white rounded-xl shadow-sm p-6 border border-gray-100 text-center py-12">
          <p className="text-gray-400 text-lg">
            {showArchived ? 'Aucun événement archivé' : 'Aucun événement'}
          </p>
        </div>
      ) : (
        <EventTable
          events={events}
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
              {editingEvent ? "Modifier l'événement" : 'Nouvel événement'}
            </h2>
            <EventForm
              initialData={editingEvent || undefined}
              onSubmit={editingEvent ? handleUpdate : handleCreate}
              onCancel={closeModal}
              isEditing={!!editingEvent}
              members={members}
            />
          </div>
        </div>
      )}

      <ConfirmDialog
        isOpen={!!deleteTarget}
        title="Supprimer l'événement"
        message="Êtes-vous sûr de vouloir supprimer cet événement ?"
        onConfirm={handleDelete}
        onCancel={() => setDeleteTarget(null)}
        confirmLabel="Supprimer"
      />
    </div>
  );
};