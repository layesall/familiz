import React, { useState, useEffect, useCallback } from 'react';
import { settingsAPI } from '../api/settings';
import { reportsAPI } from '../api/reports';
import { LoadingSpinner } from '../components/common/LoadingSpinner';
import { ErrorMessage } from '../components/common/ErrorMessage';
import { downloadBlob, getEventTypeLabel } from '../utils/helpers';
import type { ContributionSettings, EventDefault } from '../types';

export const Settings: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');
  const [success, setSuccess] = useState<string>('');
  const [contributions, setContributions] = useState<ContributionSettings>({
    amount_single: 0,
    amount_married: 0,
    amount_minor: 0,
    current_year: new Date().getFullYear(),
  });
  const [eventsDefaults, setEventsDefaults] = useState<EventDefault[]>([]);
  const [archiveLoading, setArchiveLoading] = useState(false);

  const fetchSettings = useCallback(async () => {
    setLoading(true);
    setError('');
    try {
      const [contribRes, eventsRes] = await Promise.all([
        settingsAPI.getContributions(),
        settingsAPI.getEventsDefaults(),
      ]);
      setContributions(contribRes.data || { amount_single: 0, amount_married: 0, amount_minor: 0, current_year: new Date().getFullYear() });
      setEventsDefaults(eventsRes.data || []);
    } catch (err) {
      setError('Erreur lors du chargement des paramètres');
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchSettings();
  }, [fetchSettings]);

  const handleContributionsUpdate = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setSuccess('');
    try {
      await settingsAPI.updateContributions({
        amount_single: contributions.amount_single,
        amount_married: contributions.amount_married,
        amount_minor: contributions.amount_minor,
      });
      setSuccess('Montants de cotisation mis à jour avec succès');
      await fetchSettings();
    } catch (err) {
      setError('Erreur lors de la mise à jour des montants');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleEventDefaultUpdate = async (type: 'wedding' | 'baptism', value: string) => {
    setLoading(true);
    setError('');
    setSuccess('');
    try {
      await settingsAPI.updateEventDefault(type, { default_amount: parseFloat(value) || 0 });
      setSuccess(`Montant par défaut pour ${getEventTypeLabel(type)} mis à jour`);
      await fetchSettings();
    } catch (err) {
      setError('Erreur lors de la mise à jour');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleArchive = async () => {
    setArchiveLoading(true);
    setError('');
    setSuccess('');
    try {
      await settingsAPI.archive();
      setSuccess('✅ Année archivée avec succès');
      await fetchSettings();
    } catch (err) {
      setError('Erreur lors de l\'archivage');
      console.error(err);
    } finally {
      setArchiveLoading(false);
    }
  };

  const handleUnarchive = async () => {
    setArchiveLoading(true);
    setError('');
    setSuccess('');
    try {
      await settingsAPI.unarchive();
      setSuccess('✅ Désarchivage effectué avec succès');
      await fetchSettings();
    } catch (err) {
      setError('Erreur lors du désarchivage');
      console.error(err);
    } finally {
      setArchiveLoading(false);
    }
  };

  const handleDownloadReport = async () => {
    setLoading(true);
    try {
      const year = contributions.current_year || new Date().getFullYear();
      const res = await reportsAPI.getAnnualReport(year);
      downloadBlob(res.data, `rapport_annuel_${year}.pdf`);
      setSuccess(`📄 Rapport annuel ${year} téléchargé`);
    } catch (err) {
      setError('Erreur lors du téléchargement du rapport');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  if (loading && !success) return <LoadingSpinner size="lg" />;

  return (
    <div className="space-y-8">
      <h1 className="text-3xl font-bold text-gray-900">⚙️ Paramètres</h1>

      {error && <ErrorMessage message={error} />}
      {success && (
        <div className="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded-lg flex items-center justify-between">
          <span>{success}</span>
          <button onClick={() => setSuccess('')} className="text-green-600 hover:text-green-800">✕</button>
        </div>
      )}

      {/* Montants de cotisation */}
      <div className="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
        <h2 className="text-lg font-bold text-gray-900 mb-4">💰 Montants de cotisation</h2>
        <p className="text-sm text-gray-500 mb-4">Année en cours : <strong>{contributions.current_year}</strong></p>
        <form onSubmit={handleContributionsUpdate} className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Célibataire</label>
            <input
              type="number"
              step="0.01"
              min="0"
              value={contributions.amount_single}
              onChange={(e) => setContributions(prev => ({ ...prev, amount_single: parseFloat(e.target.value) || 0 }))}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Marié(e)</label>
            <input
              type="number"
              step="0.01"
              min="0"
              value={contributions.amount_married}
              onChange={(e) => setContributions(prev => ({ ...prev, amount_married: parseFloat(e.target.value) || 0 }))}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Mineur(e)</label>
            <input
              type="number"
              step="0.01"
              min="0"
              value={contributions.amount_minor}
              onChange={(e) => setContributions(prev => ({ ...prev, amount_minor: parseFloat(e.target.value) || 0 }))}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            />
          </div>
          <div className="md:col-span-3 flex justify-end">
            <button type="submit" className="bg-gray-700 hover:bg-gray-600 text-white px-6 py-2 rounded-lg font-semibold" disabled={loading}>
              Mettre à jour
            </button>
          </div>
        </form>
      </div>

      {/* Événements */}
      <div className="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
        <h2 className="text-lg font-bold text-gray-900 mb-4">🎉 Montants par défaut des événements</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {eventsDefaults.map((item) => (
            <div key={item.event_type} className="flex items-center gap-4 p-3 bg-gray-50 rounded-lg">
              <span className="font-medium text-gray-700 min-w-[100px]">
                {getEventTypeLabel(item.event_type)}
              </span>
              <input
                type="number"
                step="0.01"
                min="0"
                value={item.default_amount}
                onChange={(e) => handleEventDefaultUpdate(item.event_type, e.target.value)}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
          ))}
        </div>
      </div>

      {/* Archivage */}
      <div className="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
        <h2 className="text-lg font-bold text-gray-900 mb-4">📦 Archivage</h2>
        <div className="flex flex-wrap gap-4">
          <button
            onClick={handleArchive}
            disabled={archiveLoading}
            className="bg-gray-700 hover:bg-gray-600 text-white px-4 py-2 rounded-lg font-semibold flex items-center gap-2 disabled:opacity-70"
          >
            {archiveLoading ? '⏳' : '📦'} Archiver l'année en cours
          </button>
          <button
            onClick={handleUnarchive}
            disabled={archiveLoading}
            className="bg-gray-200 hover:bg-gray-300 text-gray-800 px-4 py-2 rounded-lg font-semibold flex items-center gap-2 disabled:opacity-70"
          >
            🔓 Désarchiver
          </button>
          <button
            onClick={handleDownloadReport}
            disabled={loading}
            className="bg-gray-200 hover:bg-gray-300 text-gray-800 px-4 py-2 rounded-lg font-semibold flex items-center gap-2 disabled:opacity-70"
          >
            📄 Rapport annuel
          </button>
        </div>
        <p className="text-sm text-gray-500 mt-3">
          L'archivage verrouille l'année en cours. Le désarchivage n'est possible que si l'année N est vide.
        </p>
      </div>
    </div>
  );
};