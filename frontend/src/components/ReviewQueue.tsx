import { useEffect, useState } from 'react';
import { getReviewQueue, startReview } from '../api';
import { Clock, ArrowRight } from 'lucide-react';
import type { Deal } from '../types';

interface Props { onReview: (deal: Deal) => void; }

export function ReviewQueue({ onReview }: Props) {
  const [queue, setQueue] = useState<Deal[]>([]);
  const [loading, setLoading] = useState(true);
  const [reviewerName, setReviewerName] = useState('QC Analyst');

  useEffect(() => { getReviewQueue().then(setQueue).finally(() => setLoading(false)); }, []);

  const handleClaim = async (deal: Deal) => {
    await startReview(deal.id, reviewerName);
    const updated = await getReviewQueue();
    setQueue(updated);
    const claimed = updated.find((d) => d.id === deal.id);
    if (claimed) onReview(claimed);
  };

  if (loading) return <div className="text-[#666] text-center py-20">Loading review queue...</div>;

  return (
    <div className="max-w-3xl space-y-4">
      <div className="flex items-center justify-between">
        <div>
          <p className="text-xs text-[#555]">{queue.length} deck{queue.length !== 1 ? 's' : ''} pending review</p>
        </div>
        <div className="flex items-center gap-2">
          <span className="text-xs text-[#555]">Reviewer:</span>
          <input value={reviewerName} onChange={(e) => setReviewerName(e.target.value)}
            className="bg-white/[0.05] border border-white/[0.08] rounded-lg px-3 py-1.5 text-xs text-white w-36 focus:outline-none focus:border-blue-500/50" />
        </div>
      </div>

      {queue.length === 0 ? (
        <div className="text-center py-16">
          <Clock className="w-8 h-8 text-[#333] mx-auto mb-3" />
          <p className="text-[#666] text-sm">No decks pending review. All clear.</p>
        </div>
      ) : (
        <div className="space-y-2">
          {queue.map((deal) => (
            <div key={deal.id} className="flex items-center justify-between p-4 rounded-xl border border-white/[0.08] bg-white/[0.02] hover:bg-white/[0.04] transition-all">
              <div>
                <h3 className="text-sm font-semibold text-white">{deal.property.name}</h3>
                <p className="text-xs text-[#666]">{deal.property.address.city}, {deal.property.address.state} · {deal.deck_type.replace(/_/g, ' ')}</p>
                {deal.review?.reviewer_id && <p className="text-[10px] text-blue-400 mt-1">Claimed by {deal.review.reviewer_id}</p>}
              </div>
              <button onClick={() => deal.status === 'in_review' ? onReview(deal) : handleClaim(deal)}
                className="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium bg-white text-black rounded-lg hover:bg-white/90 transition-all">
                {deal.status === 'in_review' ? 'Continue' : 'Claim'} <ArrowRight className="w-3 h-3" />
              </button>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
