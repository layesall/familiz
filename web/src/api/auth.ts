import api from './axios';
import type { AuthResponse } from '../types'

export const authAPI = {
 register: (data: {
    email: string;
    password: string;
    first_name: string;
    last_name: string;
    birth_date: string;
    marital_status: string;
  }) => api.post<{ message: string; member_id: number; email: string }>('/register', data),
  
  login: (data: { email: string; password: string }) =>
    api.post<AuthResponse>('/login', data),
};