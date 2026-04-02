import type { Deal, Section } from './types';

const BASE = '/api';

export async function createDeal(formData: FormData): Promise<Deal> {
  const res = await fetch(`${BASE}/deals`, { method: 'POST', body: formData });
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

export async function listDeals(): Promise<Deal[]> {
  const res = await fetch(`${BASE}/deals`);
  return res.json();
}

export async function getDeal(id: string): Promise<Deal> {
  const res = await fetch(`${BASE}/deals/${id}`);
  if (!res.ok) throw new Error('Deal not found');
  return res.json();
}

export async function getDeckHTML(id: string): Promise<string> {
  const res = await fetch(`${BASE}/deals/${id}/deck`);
  return res.text();
}

export async function getSections(id: string): Promise<Section[]> {
  const res = await fetch(`${BASE}/deals/${id}/sections`);
  return res.json();
}

export async function updateSection(
  dealId: string,
  sectionIdx: number,
  data: { title?: string; content?: string }
): Promise<Section> {
  const res = await fetch(`${BASE}/deals/${dealId}/sections/${sectionIdx}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

export function getDeckPDFUrl(id: string): string {
  return `${BASE}/deals/${id}/deck.pdf`;
}

// Review queue APIs
export async function getReviewQueue(): Promise<Deal[]> {
  const res = await fetch(`${BASE}/reviews`);
  return res.json();
}

export async function startReview(dealId: string, reviewerId: string) {
  const res = await fetch(`${BASE}/deals/${dealId}/review/start`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ reviewer_id: reviewerId }),
  });
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

export async function completeReview(
  dealId: string,
  status: 'approved' | 'needs_revision',
  notes: string
) {
  const res = await fetch(`${BASE}/deals/${dealId}/review/complete`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ status, notes }),
  });
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

export async function reviewEdit(
  dealId: string,
  sectionIdx: number,
  data: { title?: string; content?: string }
): Promise<Section> {
  const res = await fetch(`${BASE}/deals/${dealId}/review/edit`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ section_idx: sectionIdx, ...data }),
  });
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

// Figma integration APIs
export async function linkFigmaFile(dealId: string, fileKey: string) {
  const res = await fetch(`${BASE}/deals/${dealId}/figma/link`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ file_key: fileKey }),
  });
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

export async function postFigmaComment(dealId: string, message: string) {
  const res = await fetch(`${BASE}/deals/${dealId}/figma/comment`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ message }),
  });
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

export function getFigmaExportUrl(dealId: string): string {
  return `${BASE}/deals/${dealId}/figma/export`;
}
