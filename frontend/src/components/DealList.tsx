import { useEffect, useState } from 'react';
import { listDeals } from '../api';
import { Building2, TrendingUp } from 'lucide-react';
import type { Deal } from '../types';

interface Props { onSelect: (deal: Deal) => void; }

export function DealList({ onSelect }: Props) {
  const [deals, setDeals] = useState<Deal[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => { listDeals().then(setDeals).finally(() => setLoading(false)); }, []);

  if (loading) return <div className="text-[#666] text-center py-20">Loading deals...</div>;

  if (deals.length === 0) {
    return (
      <div className="text-center py-20">
        <Building2 className="w-10 h-10 text-[#333] mx-auto mb-3" />
        <p className="text-[#666]">No deals yet. Create your first deal to get started.</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="grid gap-3 md:grid-cols-2 lg:grid-cols-3">
        {deals.map((deal) => (
          <button
            key={deal.id}
            onClick={() => onSelect(deal)}
            className="text-left p-5 rounded-xl border border-white/[0.08] bg-white/[0.02] hover:bg-white/[0.05] transition-all group"
          >
            <div className="flex items-center justify-between mb-3">
              <span className={`text-[10px] font-semibold uppercase tracking-wider px-2 py-0.5 rounded-full ${
                deal.status === 'ready' || deal.status === 'approved' ? 'bg-emerald-500/20 text-emerald-400'
                : deal.status === 'in_review' ? 'bg-yellow-500/20 text-yellow-400'
                : deal.status === 'failed' ? 'bg-red-500/20 text-red-400'
                : 'bg-white/10 text-[#888]'
              }`}>
                {deal.status.replace(/_/g, ' ')}
              </span>
              <span className="text-[11px] text-[#555] capitalize">{deal.deck_type.replace(/_/g, ' ')}</span>
            </div>
            <h3 className="text-sm font-semibold text-white mb-1 group-hover:text-blue-400 transition-colors">
              {deal.property.name}
            </h3>
            <p className="text-xs text-[#666] mb-3">
              {deal.property.address.city}, {deal.property.address.state}
            </p>
            {deal.analysis && (
              <div className="flex gap-4 pt-3 border-t border-white/[0.06]">
                <div className="flex items-center gap-1.5">
                  <TrendingUp className="w-3 h-3 text-emerald-400" />
                  <span className="text-xs text-[#888]">NOI ${(deal.analysis.noi / 1000).toFixed(0)}K</span>
                </div>
                <span className="text-xs text-[#888]">{(deal.analysis.occupancy_rate * 100).toFixed(0)}% occ</span>
              </div>
            )}
          </button>
        ))}
      </div>
    </div>
  );
}
