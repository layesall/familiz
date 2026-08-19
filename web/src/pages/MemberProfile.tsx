import React, { useState, useEffect, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { membersAPI } from '../api/members';
import { settingsAPI } from '../api/settings';
import { contributionsAPI } from '../api/contributions';
import { eventsAPI } from '../api/events';
import { MemberForm } from '../components/members/MemberForm';
import { ContributionsForm } from '../components/contributions/ContributionsForm';
import { EventForm } from '../components/events/EventForm';
import { ArchivedToggle } from '../components/common/ArchivedToggle';
import { LoadingSpinner } from '../components/common/LoadingSpinner';
import { ErrorMessage } from '../components/common/ErrorMessage';
import { ConfirmDialog } from '../components/common/ConfirmDialog';
import { formatDate, formatCurrency, getMaritalStatusLabel, getEventTypeLabel, getMonthLabel, downloadBlob } from '../utils/helpers';
import type { Member, Contribution, Event } from '../types';

// Composant indicateur de paiement
const PaymentStatus: React.FC<{ member: Member; contributions: Contribution[]; year: number }> = ({
  member,
  contributions,
  year,
}) => {
  const [settings, setSettings] = useState({ amount_single: 0, amount_married: 0, amount_minor: 0 });

  useEffect(() => {
    settingsAPI.getContributions().then(res => {
      setSettings(res.data || { amount_single: 0, amount_married: 0, amount_minor: 0 });
    }).catch(console.error);
  }, []);

  const getMonthlyAmount = () => {
    switch (member.marital_status) {
      case 'single': return settings.amount_single;
      case 'married': return settings.amount_married;
      case 'minor': return settings.amount_minor;
      default: return 0;
    }
  };

  const currentMonth = new Date().getMonth() + 1;
  const monthlyAmount = getMonthlyAmount();
  const expectedTotal = currentMonth * monthlyAmount;
  const paidTotal = contributions
    .filter(t => t.year === year)
    .reduce((sum, t) => sum + t.amount, 0);

  const isUpToDate = paidTotal >= expectedTotal;
  const diff = paidTotal - expectedTotal;

  return (
    <div className={`p-4 rounded-lg ${isUpToDate ? 'bg-green-50 border border-green-200' : 'bg-amber-50 border border-amber-200'}`}>
      <div className="flex items-center justify-between">
        <div>
          <p className="text-sm font-medium text-gray-700">Situation de paiement</p>
          <p className="text-lg font-bold">
            {isUpToDate ? (
              <span className="text-green-700">✅ À jour</span>
            ) : (
              <span className="text-amber-700">⚠️ Pas à jour</span>
            )}
          </p>
        </div>
        <div className="text-right">
          <p className="text-sm text-gray-600">
            Cotisation mensuelle : {formatCurrency(monthlyAmount)}
          </p>
          <p className="text-sm text-gray-600">
            Mois en cours : {getMonthLabel(currentMonth)} {year}
          </p>
          <p className={`text-sm font-medium ${isUpToDate ? 'text-green-700' : 'text-amber-700'}`}>
            {isUpToDate ? `Avance de ${formatCurrency(diff)}` : `Manque ${formatCurrency(Math.abs(diff))}`}
          </p>
        </div>
      </div>
    </div>
  );
};

export const MemberProfile: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');
  const [showArchived, setShowArchived] = useState(false);
  const [member, setMember] = useState<Member | null>(null);
  const [contributions, setContributions] = useState<Contribution[]>([]);
  const [events, setEvents] = useState<Event[]>([]);
  const [pdfLoading, setPdfLoading] = useState(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
  const [currentYear, setCurrentYear] = useState(new Date().getFullYear());

  // États pour les modales
  const [showEditModal, setShowEditModal] = useState(false);
  const [editingMember, setEditingMember] = useState<Member | null>(null);
  const [showContributionModal, setShowContributionModal] = useState(false);
  const [showEventModal, setShowEventModal] = useState(false);
  const [membersList, setMembersList] = useState<Member[]>([]);

  const memberId = parseInt(id || '0');

  const fetchData = useCallback(async () => {
    if (!memberId) return;
    setLoading(true);
    setError('');
    try {
      const [profileRes, settingsRes, membersRes] = await Promise.all([
        membersAPI.getById(memberId, showArchived),
        settingsAPI.getContributions(),
        membersAPI.getAll(),
      ]);
      const profile = profileRes.data;
      if (profile) {
        setMember(profile.member || null);
        setContributions(profile.transactions || []);
        setEvents(profile.events || []);
        setCurrentYear(settingsRes.data?.current_year || new Date().getFullYear());
        setMembersList(membersRes.data || []);
      } else {
        setMember(null);
        setContributions([]);
        setEvents([]);
      }
    } catch (err) {
      setError('Erreur lors du chargement du profil');
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, [memberId, showArchived]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const handleDownloadPDF = async () => {
    if (!member) return;
    setPdfLoading(true);
    try {
      const res = await membersAPI.getPDF(memberId, showArchived);
      const filename = `profil_${member.first_name}_${member.last_name}.pdf`;
      downloadBlob(res.data, filename);
    } catch (err) {
      setError('Erreur lors du téléchargement du PDF');
      console.error(err);
    } finally {
      setPdfLoading(false);
    }
  };

  const handleDelete = async () => {
    try {
      await membersAPI.delete(memberId);
      navigate('/members');
    } catch (err) {
      setError('Erreur lors de la suppression');
      console.error(err);
    }
  };

  // Mise à jour du membre
  const handleUpdateMember = async (data: Partial<Omit<Member, 'id' | 'archived'>>) => {
    if (!member) return;
    try {
      const res = await membersAPI.update(member.id, data);
      setMember(res.data);
      setShowEditModal(false);
      setEditingMember(null);
      await fetchData();
    } catch (err) {
      setError('Erreur lors de la mise à jour');
      console.error(err);
    }
  };

  // Création d'une contribution
  const handleCreateContribution = async (data: Omit<Contribution, 'id' | 'archived' | 'created_at'>) => {
    try {
      await contributionsAPI.create(data);
      setShowContributionModal(false);
      await fetchData();
    } catch (err) {
      setError("Erreur lors de la création de la contribution");
      console.error(err);
    }
  };

  // Création d'un événement
  const handleCreateEvent = async (data: Omit<Event, 'id' | 'archived' | 'created_at'>) => {
    try {
      await eventsAPI.create(data);
      setShowEventModal(false);
      await fetchData();
    } catch (err) {
      setError("Erreur lors de la création de l'événement");
      console.error(err);
    }
  };

  const openEditModal = () => {
    setEditingMember(member);
    setShowEditModal(true);
  };

  const closeEditModal = () => {
    setShowEditModal(false);
    setEditingMember(null);
  };

  if (loading) return <LoadingSpinner size="lg" />;

  if (!member) {
    return (
      <div className="bg-white rounded-xl shadow-sm p-6 border border-gray-100 text-center py-12">
        <p className="text-gray-400 text-lg">Membre non trouvé</p>
        <button onClick={() => navigate('/members')} className="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-lg font-semibold mt-4">
          Retour à la liste
        </button>
      </div>
    );
  }

  const totalContributions = contributions.reduce((sum, t) => sum + (t.amount || 0), 0);
  const totalEvents = events.reduce((sum, e) => sum + (e.amount_received || 0), 0);

  // Trier du plus récent au plus ancien
  const sortedContributions = [...contributions].sort((a, b) => {
    const da = new Date(a.created_at || 0);
    const db = new Date(b.created_at || 0);
    return db.getTime() - da.getTime();
  });
  const sortedEvents = [...events].sort((a, b) => {
    const da = new Date(a.created_at || 0);
    const db = new Date(b.created_at || 0);
    return db.getTime() - da.getTime();
  });

  return (
    <div className="space-y-6 pb-24">
      {error && <ErrorMessage message={error} onRetry={fetchData} />}

      {/* Header */}
      <div className="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
        <div className="flex flex-wrap items-start justify-between gap-4">
          <div>
            <div className="flex items-center gap-3 flex-wrap">
              <h1 className="text-2xl font-bold text-gray-900">
                {member.first_name} {member.last_name}
              </h1>
              <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-indigo-100 text-indigo-800`}>
                {getMaritalStatusLabel(member.marital_status)}
              </span>
              {member.archived && (
                <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-200 text-gray-700">Archivé</span>
              )}
            </div>
            <p className="text-gray-500 mt-1">
              Né(e) le {formatDate(member.birth_date)} • ID: {member.id}
            </p>
          </div>
          <div className="flex items-center gap-3 flex-wrap">
            <ArchivedToggle showArchived={showArchived} onChange={setShowArchived} />
            <button
              onClick={handleDownloadPDF}
              disabled={pdfLoading}
              className="bg-gray-200 hover:bg-gray-300 text-gray-800 px-4 py-2 rounded-lg font-semibold flex items-center gap-2"
            >
              {pdfLoading ? '⏳' : '📄'} PDF
            </button>
            <button onClick={() => navigate('/members')} className="bg-gray-200 hover:bg-gray-300 text-gray-800 px-4 py-2 rounded-lg font-semibold">
              Retour
            </button>
          </div>
        </div>
      </div>

      {/* Indicateur de paiement */}
      <PaymentStatus member={member} contributions={contributions} year={currentYear} />

      {/* Stats */}
      <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
        <div className="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
          <p className="text-sm text-gray-500">💳 Contributions</p>
          <p className="text-2xl font-bold text-gray-900">{contributions.length}</p>
          <p className="text-sm text-green-600">{formatCurrency(totalContributions)}</p>
        </div>
        <div className="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
          <p className="text-sm text-gray-500">🎉 Événements</p>
          <p className="text-2xl font-bold text-gray-900">{events.length}</p>
          <p className="text-sm text-purple-600">{formatCurrency(totalEvents)}</p>
        </div>
        <div className="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
          <p className="text-sm text-gray-500">Statut</p>
          <p className="text-lg font-medium text-gray-900">
            {member.archived ? 'Archivé' : 'Actif'}
          </p>
        </div>
      </div>

      {/* Historique des contributions */}
      <div className="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-bold text-gray-900">💳 Historique des contributions</h2>
          <button
            onClick={() => setShowContributionModal(true)}
            className="bg-emerald-600 hover:bg-emerald-700 text-white px-3 py-1.5 rounded-lg text-sm font-semibold flex items-center gap-1"
          >
            + Ajouter
          </button>
        </div>
        {sortedContributions.length === 0 ? (
          <p className="text-gray-400 text-center py-4">Aucune contribution</p>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-gray-200">
                  <th className="text-left py-2 font-medium text-gray-500">Mois</th>
                  <th className="text-left py-2 font-medium text-gray-500">Année</th>
                  <th className="text-left py-2 font-medium text-gray-500">Montant</th>
                  <th className="text-left py-2 font-medium text-gray-500">Note</th>
                  <th className="text-left py-2 font-medium text-gray-500">Date</th>
                </tr>
              </thead>
              <tbody>
                {sortedContributions.map((t) => (
                  <tr key={t.id} className="border-b border-gray-50 hover:bg-gray-50">
                    <td className="py-2">{getMonthLabel(t.month)}</td>
                    <td className="py-2">{t.year}</td>
                    <td className="py-2 text-left font-medium">{formatCurrency(t.amount)}</td>
                    <td className="py-2 text-left text-gray-500">{t.note || '-'}</td>
                    <td className="py-2 text-gray-400 text-xs">
                      {t.created_at ? new Date(t.created_at).toLocaleDateString('fr-FR') : '-'}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Historique des événements */}
      <div className="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-bold text-gray-900">🎉 Historique des événements</h2>
          <button
            onClick={() => setShowEventModal(true)}
            className="bg-purple-600 hover:bg-purple-700 text-white px-3 py-1.5 rounded-lg text-sm font-semibold flex items-center gap-1"
          >
            + Ajouter
          </button>
        </div>
        {sortedEvents.length === 0 ? (
          <p className="text-gray-400 text-center py-4">Aucun événement</p>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-gray-200">
                  <th className="text-left py-2 font-medium text-gray-500">Type</th>
                  <th className="text-left py-2 font-medium text-gray-500">Date</th>
                  <th className="text-right py-2 font-medium text-gray-500">Montant reçu</th>
                </tr>
              </thead>
              <tbody>
                {sortedEvents.map((e) => (
                  <tr key={e.id} className="border-b border-gray-50 hover:bg-gray-50">
                    <td className="py-2">
                      <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${e.type === 'wedding' ? 'bg-pink-100 text-pink-800' : 'bg-cyan-100 text-cyan-800'}`}>
                        {getEventTypeLabel(e.type)}
                      </span>
                    </td>
                    <td className="py-2">{formatDate(e.event_date)}</td>
                    <td className="py-2 text-right font-medium">{formatCurrency(e.amount_received)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Footer actions */}
      <div className="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 p-4 shadow-lg z-40">
        <div className="container mx-auto max-w-7xl flex flex-wrap items-center justify-end gap-3">
          <button
            onClick={openEditModal}
            className="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-lg font-semibold text-sm flex items-center gap-2"
          >
            ✏️ Modifier
          </button>
          <button
            onClick={() => setShowDeleteConfirm(true)}
            className="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-lg font-semibold text-sm flex items-center gap-2"
          >
            🗑️ Supprimer
          </button>
        </div>
      </div>

      {/* Modale d'édition du membre */}
      {showEditModal && editingMember && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm">
          <div className="bg-white rounded-xl shadow-xl max-w-2xl w-full mx-4 p-6 max-h-[90vh] overflow-y-auto">
            <h2 className="text-lg font-bold text-gray-900 mb-4">Modifier le membre</h2>
            <MemberForm
              initialData={editingMember}
              onSubmit={handleUpdateMember}
              onCancel={closeEditModal}
              isEditing={true}
            />
          </div>
        </div>
      )}

      {/* Modale d'ajout de contribution */}
      {showContributionModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm">
          <div className="bg-white rounded-xl shadow-xl max-w-2xl w-full mx-4 p-6 max-h-[90vh] overflow-y-auto">
            <h2 className="text-lg font-bold text-gray-900 mb-4">Nouvelle contribution</h2>
            <ContributionsForm
              initialData={undefined}
              onSubmit={handleCreateContribution}
              onCancel={() => setShowContributionModal(false)}
              isEditing={false}
              members={membersList}
              defaultMemberId={memberId}
            />
          </div>
        </div>
      )}

      {/* Modale d'ajout d'événement */}
      {showEventModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm">
          <div className="bg-white rounded-xl shadow-xl max-w-2xl w-full mx-4 p-6 max-h-[90vh] overflow-y-auto">
            <h2 className="text-lg font-bold text-gray-900 mb-4">Nouvel événement</h2>
            <EventForm
              initialData={undefined}
              onSubmit={handleCreateEvent}
              onCancel={() => setShowEventModal(false)}
              isEditing={false}
              members={membersList}
              defaultMemberId={memberId}
            />
          </div>
        </div>
      )}

      <ConfirmDialog
        isOpen={showDeleteConfirm}
        title="⚠️ Supprimer définitivement"
        message={`Êtes-vous sûr de vouloir supprimer ${member.first_name} ${member.last_name} ? Cette action est irréversible et supprimera également toutes ses contributions et événements associés.`}
        onConfirm={handleDelete}
        onCancel={() => setShowDeleteConfirm(false)}
        confirmLabel="Supprimer définitivement"
      />
    </div>
  );
};