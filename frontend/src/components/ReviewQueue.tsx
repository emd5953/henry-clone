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

  if (loading) return <div className="text-[#7A8578] text-center py-20">Loading review queue...</div>;

  return (
    <div className="max-w-3xl space-y-4">
      <div className="flex items-center justify-between">
        <p className="text-xs text-[#7A8578]">{queue.length} deck{queue.length !== 1 ? 's' : ''} pending review</p>
        <div className="flex items-center gap-2">
          <span className="text-xs text-[#7A8578]">Reviewer:</span>
          <input value={reviewerName} onChange={(e) => setReviewerName(e.target.value)}
            className="bg-white border border-[#E6E2DA] rounded-lg px-3 py-1.5 text-xs text-[#2D3A31] w-36 focus:outline-none focus:border-[#8C9A84]" />
        </div>
      </div>

      {queue.length === 0 ? (
        <div className="text-center py-16">
          <Clock className="w-8 h-8 text-[#DCCFC2] mx-auto mb-3" />
          <p className="text-[#7A8578] text-sm">No decks pending review. All clear.</p>
        </div>
      ) : (
        <div className="space-y-2">
          {queue.map((deal) => (
            <div key={deal.id} className="flex items-center justify-between p-4 rounded-xl border border-[#E6E2DA] bg-white hover:shadow-[0_4px_6px_-1px_rgba(45,58,49,0.05)] transition-all">
              <div>
                <h3 className="text-sm font-semibold text-[#2D3A31]">{deal.property.name}</h3>
                <p className="text-xs text-[#7A8578]">{deal.property.address.city}, {deal.property.address.state} · {deal.deck_type.replace(/_/g, ' ')}</p>
                {deal.review?.reviewer_id && <p className="text-[10px] text-[#C27B66] mt-1">Claimed by {deal.review.reviewer_id}</p>}
              </div>
              <button onClick={() => deal.status === 'in_review' ? onReview(deal) : handleClaim(deal)}
                className="flex items-center gap-1.5 px-4 py-2 text-xs font-semibold bg-[#2D3A31] text-white rounded-lg hover:bg-[#3D4A41] transition-all">
                {deal.status === 'in_review' ? 'Continue' : 'Claim'} <ArrowRight className="w-3 h-3" />
              </button>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
