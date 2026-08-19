// src/api/contributions.ts
import api from './axios';
import type { Contribution } from '../types';

export const contributionsAPI = {
  getAll: (memberId?: number, archived = false) => {
    const params = new URLSearchParams();
    if (memberId) params.append('member_id', String(memberId));
    if (archived) params.append('archived', 'true');
    return api.get<Contribution[]>(`/contributions?${params.toString()}`);
  },

  create: (data: Omit<Contribution, 'id' | 'archived' | 'created_at'>) => {
    return api.post<Contribution>('/contributions', data);
  },

  update: (id: number, data: Partial<Omit<Contribution, 'id' | 'archived' | 'created_at'>>) => {
    return api.put<Contribution>(`/contributions/${id}`, data);
  },

  delete: (id: number) => api.delete<void>(`/contributions/${id}`),
};