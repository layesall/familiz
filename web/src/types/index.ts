// Types correspondant à l'API
export interface Member {
  id: number;
  first_name: string;
  last_name: string;
  birth_date: string; // format YYYY-MM-DD
  marital_status: 'single' | 'married' | 'minor';
  archived?: boolean;
}

export interface Contribution {
  id: number;
  member_id: number;
  month: number; // 1-12
  year: number;
  amount: number;
  note?: string;
  archived?: boolean;
  created_at?: string;
}

export interface Event {
  id: number;
  member_id: number;
  type: 'wedding' | 'baptism';
  amount_received: number;
  event_date: string; // YYYY-MM-DD
  archived?: boolean;
  created_at?: string;
}

export interface ContributionSettings {
  amount_single: number;
  amount_married: number;
  amount_minor: number;
  current_year: number;
}

export interface EventDefault {
  event_type: 'wedding' | 'baptism';
  default_amount: number;
}

export interface User {
  email: string;
  role: 'admin';
}

export interface AuthResponse {
  token: string;
  email: string;
  role: 'admin';
}