import api from './axios';
import type { Member, Contribution, Event } from '../types';

export const membersAPI = {
  getAll: (archived = false) =>
    api.get<Member[]>(`/members${archived ? '?archived=true' : ''}`),

  getById: (id: number, archived = false) =>
    api.get<{member:Member, transactions?: Contribution[]; events?: Event[] }>(
      `/profile/${id}${archived ? '?archived=true' : ''}`
    ),

  create: (data: Omit<Member, 'id' | 'archived'>) =>
    api.post<Member>('/members', data),

  update: (id: number, data: Partial<Omit<Member, 'id' | 'archived'>>) =>
    api.put<Member>(`/members/${id}`, data),

  delete: (id: number) => api.delete<void>(`/members/${id}`),

  getPDF: (id: number, archived = false) =>
    api.get<Blob>(`/members/${id}/pdf${archived ? '?archived=true' : ''}`, {
      responseType: 'blob',
    }),
};