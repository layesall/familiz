import api from './axios';
import type { Event } from '../types';

export const eventsAPI = {
  getAll: (memberId?: number, archived = false) => {
    const params = new URLSearchParams();
    if (memberId) params.append('member_id', String(memberId));
    if (archived) params.append('archived', 'true');
    return api.get<Event[]>(`/events?${params.toString()}`);
  },

  create: (data: Omit<Event, 'id' | 'archived' | 'created_at'>) =>
    api.post<Event>('/events', data),

  update: (id: number, data: Partial<Omit<Event, 'id' | 'archived' | 'created_at'>>) =>
    api.put<Event>(`/events/${id}`, data),

  delete: (id: number) => api.delete<void>(`/events/${id}`),
};