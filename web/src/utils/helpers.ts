export function formatDate(dateStr: string): string {
  if (!dateStr) return '-';
  const d = new Date(dateStr);
  return d.toLocaleDateString('fr-FR', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });
}

export function formatCurrency(amount: number): string {
  return new Intl.NumberFormat('fr-FR', {
    style: 'currency',
    currency: 'EUR',
  }).format(amount);
}

export function getMaritalStatusLabel(status: string): string {
  const map: Record<string, string> = {
    single: 'Célibataire',
    married: 'Marié(e)',
    minor: 'Mineur(e)',
  };
  return map[status] || status;
}

export function getMaritalStatusBadge(status: string): string {
  const map: Record<string, string> = {
    single: 'bg-blue-100 text-blue-800',
    married: 'bg-purple-100 text-purple-800',
    minor: 'bg-yellow-100 text-yellow-800',
  };
  return map[status] || '';
}

export function getEventTypeLabel(type: string): string {
  const map: Record<string, string> = {
    wedding: 'Mariage',
    baptism: 'Baptême',
  };
  return map[type] || type;
}

export function getEventTypeBadge(type: string): string {
  const map: Record<string, string> = {
    wedding: 'bg-pink-100 text-pink-800',
    baptism: 'bg-cyan-100 text-cyan-800',
  };
  return map[type] || '';
}

export function getMonthLabel(month: number): string {
  const months = [
    'Janvier', 'Février', 'Mars', 'Avril', 'Mai', 'Juin',
    'Juillet', 'Août', 'Septembre', 'Octobre', 'Novembre', 'Décembre'
  ];
  return months[month - 1] || String(month);
}

export function downloadBlob(blob: Blob, filename: string): void {
  const url = window.URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  window.URL.revokeObjectURL(url);
}