import { useEffect, useState } from 'react';
import { listDeals } from '../api';
import { Building2, TrendingUp, ArrowRight } from 'lucide-react';
import type { Deal } from '../types';

interface Props { onSelect: (deal: Deal) => void; }

export function DealList({ onSelect }: Props) {
  const [deals, setDeals] = useState<Deal[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => { listDeals().then(setDeals).finally(() => setLoading(false)); }, []);

  if (loading) return <div className="text-[#7A8578] text-center py-20">Loading deals...</div>;

  if (deals.length === 0) {
    return (
      <div className="text-center py-24">
        <div className="w-16 h-16 rounded-2xl bg-[#8C9A84]/10 flex items-center justify-center mx-auto mb-4">
          <Building2 className="w-7 h-7 text-[#8C9A84]" />
        </div>
        <h2 className="text-xl font-bold text-[#2D3A31] mb-2" style={{ fontFamily: "'Playfair Display', serif" }}>No deals yet</h2>
        <p className="text-[#7A8578] text-sm">Create your first deal to generate a professional deck.</p>
      </div>
    );
  }

  return (
    <div>
      <div className="grid grid-cols-4 gap-4 mb-8">
        {[
          { label: 'Total Deals', value: deals.length },
          { label: 'Ready', value: deals.filter(d => d.status === 'ready' || d.status === 'approved').length },
          { label: 'In Review', value: deals.filter(d => d.status === 'in_review').length },
          { label: 'Avg NOI', value: deals.filter(d => d.analysis).length > 0 ? `$${(deals.filter(d => d.analysis).reduce((s, d) => s + (d.analysis?.noi || 0), 0) / deals.filter(d => d.analysis).length / 1000).toFixed(0)}K` : '—' },
        ].map((stat, i) => (
          <div key={i} className="rounded-xl border border-[#E6E2DA] bg-white p-5 text-center shadow-[0_4px_6px_-1px_rgba(45,58,49,0.05)]">
            <p className="text-2xl font-bold text-[#2D3A31]">{stat.value}</p>
            <p className="text-xs text-[#7A8578] mt-1">{stat.label}</p>
          </div>
        ))}
      </div>

      <div className="space-y-3">
        {deals.map((deal) => (
          <button key={deal.id} onClick={() => onSelect(deal)}
            className="w-full text-left flex items-center justify-between p-5 rounded-xl border border-[#E6E2DA] bg-white hover:shadow-[0_10px_15px_-3px_rgba(45,58,49,0.05)] transition-all group">
            <div className="flex items-center gap-4">
              <div className="w-10 h-10 rounded-xl bg-[#8C9A84]/10 flex items-center justify-center shrink-0">
                <Building2 className="w-5 h-5 text-[#8C9A84]" />
              </div>
              <div>
                <div className="flex items-center gap-2">
                  <h3 className="text-sm font-semibold text-[#2D3A31] group-hover:text-[#C27B66] transition-colors">{deal.property.name}</h3>
                  <span className={`text-[10px] font-semibold uppercase tracking-wider px-2 py-0.5 rounded-full ${
                    deal.status === 'ready' || deal.status === 'approved' ? 'bg-[#8C9A84]/15 text-[#8C9A84]'
                    : deal.status === 'in_review' ? 'bg-[#C27B66]/15 text-[#C27B66]'
                    : deal.status === 'failed' ? 'bg-red-100 text-red-600'
                    : 'bg-[#F2F0EB] text-[#7A8578]'
                  }`}>{deal.status.replace(/_/g, ' ')}</span>
                </div>
                <p className="text-xs text-[#7A8578] mt-0.5">{deal.property.address.city}, {deal.property.address.state} · {deal.deck_type.replace(/_/g, ' ')}</p>
              </div>
            </div>
            <div className="flex items-center gap-6">
              {deal.analysis && (
                <div className="flex gap-5 text-xs text-[#7A8578]">
                  <span className="flex items-center gap-1"><TrendingUp className="w-3 h-3 text-[#8C9A84]" />NOI ${(deal.analysis.noi / 1000).toFixed(0)}K</span>
                  <span>{(deal.analysis.occupancy_rate * 100).toFixed(0)}% occ</span>
                </div>
              )}
              <ArrowRight className="w-4 h-4 text-[#DCCFC2] group-hover:text-[#C27B66] transition-colors" />
            </div>
          </button>
        ))}
      </div>
    </div>
  );
}
