import React, { useState, useEffect, useCallback } from 'react';
import { Link } from 'react-router-dom';
import { membersAPI } from '../api/members';
import { contributionsAPI } from '../api/contributions';
import { eventsAPI } from '../api/events';
import { settingsAPI } from '../api/settings';
import { LoadingSpinner } from '../components/common/LoadingSpinner';
import { ErrorMessage } from '../components/common/ErrorMessage';
import { formatCurrency } from '../utils/helpers';
import type { Contribution, Event } from '../types';

// StatsCard avec lien optionnel
interface StatsCardProps {
  label: string;
  value: string | number;
  icon: string;
  color: string;
  link?: string;
}

const StatsCard: React.FC<StatsCardProps> = ({ label, value, icon, color, link }) => {
  const content = (
    <div className={`bg-white rounded-xl shadow-sm p-6 border-l-4 ${color} hover:shadow-md transition-shadow`}>
      <div className="flex items-center justify-between">
        <div>
          <p className="text-sm font-medium text-gray-500">{label}</p>
          <p className="text-2xl font-bold text-gray-900 mt-1">{value}</p>
        </div>
        <span className="text-3xl">{icon}</span>
      </div>
    </div>
  );

  if (link) {
    return (
      <Link to={link} className="block cursor-pointer">
        {content}
      </Link>
    );
  }

  return content;
};

// Carte de total (neutre)
const TotalCard: React.FC<{
  title: string;
  amount: number;
  count: number;
  label: string;
  icon: string;
  link: string;
  linkText: string;
}> = ({ title, amount, count, label, icon, link, linkText }) => (
  <div className="bg-white rounded-xl shadow-sm p-6 border border-gray-100 hover:shadow-md transition-shadow">
    <div className="flex items-center justify-between">
      <div>
        <p className="text-sm font-medium text-gray-500">{title}</p>
        <p className="text-2xl font-bold text-gray-900 mt-1">{formatCurrency(amount)}</p>
        <p className="text-sm text-gray-400 mt-1">
          Sur {count} {label}{count > 1 ? 's' : ''}
        </p>
      </div>
      <div className="bg-gray-50 p-4 rounded-full border border-gray-100">
        <span className="text-2xl">{icon}</span>
      </div>
    </div>
    <div className="mt-4">
      <Link to={link} className="text-sm text-indigo-600 hover:text-indigo-800 transition-colors">
        {linkText} →
      </Link>
    </div>
  </div>
);

const RecentList: React.FC<{ title: string; items: any[]; renderItem: (item: any) => React.ReactNode }> = ({
  title, items, renderItem
}) => (
  <div className="bg-white rounded-xl shadow-sm p-6 border border-gray-100">
    <h3 className="text-lg font-semibold text-gray-900 mb-4">{title}</h3>
    {items.length === 0 ? (
      <p className="text-gray-400 text-center py-4">Aucune donnée récente</p>
    ) : (
      <div className="divide-y divide-gray-100">
        {items.slice(0, 5).map((item, idx) => (
          <div key={idx} className="py-3 first:pt-0 last:pb-0">
            {renderItem(item)}
          </div>
        ))}
      </div>
    )}
  </div>
);

export const Dashboard: React.FC = () => {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');
  const [stats, setStats] = useState({
    members: 0,
    contributions: 0,
    totalContributions: 0,
    events: 0,
    totalEventsReceived: 0,
    currentYear: new Date().getFullYear(),
  });
  const [recentContributions, setRecentContributions] = useState<Contribution[]>([]);
  const [recentEvents, setRecentEvents] = useState<Event[]>([]);

  const fetchData = useCallback(async () => {
    try {
      const [membersRes, txRes, evtRes, settingsRes] = await Promise.all([
        membersAPI.getAll(),
        contributionsAPI.getAll(),
        eventsAPI.getAll(),
        settingsAPI.getContributions(),
      ]);

      const members = membersRes.data || [];
      const contributions = txRes.data || [];
      const events = evtRes.data || [];

      const totalContributions = contributions.reduce((sum, t) => sum + (t.amount || 0), 0);
      const totalEventsReceived = events.reduce((sum, e) => sum + (e.amount_received || 0), 0);

      const sortedTx = [...contributions].sort((a, b) => {
        const da = new Date(a.created_at || 0);
        const db = new Date(b.created_at || 0);
        return db.getTime() - da.getTime();
      });
      const sortedEvents = [...events].sort((a, b) => {
        const da = new Date(a.created_at || 0);
        const db = new Date(b.created_at || 0);
        return db.getTime() - da.getTime();
      });

      setStats({
        members: members.length,
        contributions: contributions.length,
        totalContributions,
        events: events.length,
        totalEventsReceived,
        currentYear: settingsRes.data?.current_year || new Date().getFullYear(),
      });

      setRecentContributions(sortedTx.slice(0, 5));
      setRecentEvents(sortedEvents.slice(0, 5));
    } catch (err) {
      setError('Erreur lors du chargement des données');
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  if (loading) return <LoadingSpinner size="lg" />;

  return (
    <div className="space-y-8">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold text-gray-900">📊 Dashboard</h1>
        <span className="text-sm text-gray-500 bg-gray-100 px-3 py-1 rounded-full">
          Année {stats.currentYear}
        </span>
      </div>

      {error && <ErrorMessage message={error} onRetry={fetchData} />}

      {/* Stats avec liens */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <StatsCard
          label="Membres"
          value={stats.members}
          icon="👥"
          color="border-red-600"
          link="/members"
        />
        <StatsCard
          label="Contributions"
          value={stats.contributions}
          icon="💰"
          color="border-yellow-400"
          link="/contributions"
        />
        <StatsCard
          label="Événements"
          value={stats.events}
          icon="🎉"
          color="border-green-500"
          link="/events"
        />
      </div>

      {/* Cartes Total contributions et Total événements (version neutre) */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <TotalCard
          title="Total des contributions"
          amount={stats.totalContributions}
          count={stats.contributions}
          label="contribution"
          icon="💰"
          link="/contributions"
          linkText="Voir toutes les contributions"
        />
        <TotalCard
          title="Total des événements"
          amount={stats.totalEventsReceived}
          count={stats.events}
          label="événement"
          icon="🎉"
          link="/events"
          linkText="Voir tous les événements"
        />
      </div>

      {/* Récentes contributions et événements */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <RecentList
          title="📋 Dernières contributions"
          items={recentContributions}
          renderItem={(tx: Contribution) => (
            <div className="flex justify-between items-center">
              <div>
                <p className="font-medium text-gray-900">
                  {formatCurrency(tx.amount)}
                </p>
                <p className="text-sm text-gray-500">
                  {new Date(tx.created_at || '').toLocaleDateString('fr-FR')}
                </p>
              </div>
              <span className="text-sm text-gray-400">{tx.note || '—'}</span>
            </div>
          )}
        />

        <RecentList
          title="🎉 Derniers événements"
          items={recentEvents}
          renderItem={(evt: Event) => (
            <div className="flex justify-between items-center">
              <div>
                <p className="font-medium text-gray-900">
                  {evt.type === 'wedding' ? '💒 Mariage' : '🌙 Baptême'}
                </p>
                <p className="text-sm text-gray-500">
                  {new Date(evt.event_date).toLocaleDateString('fr-FR')}
                </p>
              </div>
              <span className="text-sm font-medium text-gray-600">
                {formatCurrency(evt.amount_received)}
              </span>
            </div>
          )}
        />
      </div>
    </div>
  );
};