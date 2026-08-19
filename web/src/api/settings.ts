import api from './axios';
import type { ContributionSettings, EventDefault } from '../types';

export const settingsAPI = {
  getContributions: () => api.get<ContributionSettings>('/settings/contributions'),
  updateContributions: (data: Omit<ContributionSettings, 'current_year'>) =>
    api.put<ContributionSettings>('/settings/contributions', data),

  getEventsDefaults: () => api.get<EventDefault[]>('/settings/events'),
  updateEventDefault: (type: 'wedding' | 'baptism', data: { default_amount: number }) =>
    api.put<EventDefault>(`/settings/events/${type}`, data),

  archive: () => api.post<{ message: string }>('/settings/archive'),
  unarchive: () => api.post<{ message: string }>('/settings/unarchive'),
};