import api from './axios';

export const reportsAPI = {
  getAnnualReport: (year: number) =>
    api.get<Blob>(`/reports/pdf?year=${year}`, {
      responseType: 'blob',
    }),
};