import { useEffect, useState } from 'react';
import { getReviewQueue, startReview } from '../api';
import type { Deal } from '../types';

interface Props {
  onReview: (deal: Deal) => void;
}

export function ReviewQueue({ onReview }: Props) {
  const [queue, setQueue] = useState<Deal[]>([]);
  const [loading, setLoading] = useState(true);
  const [reviewerName, setReviewerName] = useState('QC Analyst');

  useEffect(() => {
    getReviewQueue()
      .then(setQueue)
      .finally(() => setLoading(false));
  }, []);

  const handleClaim = async (deal: Deal) => {
    await startReview(deal.id, reviewerName);
    // Refresh and open editor
    const updated = await getReviewQueue();
    setQueue(updated);
    const claimed = updated.find((d) => d.id === deal.id);
    if (claimed) onReview(claimed);
  };

  if (loading) return <div className="loading">Loading review queue...</div>;

  return (
    <div className="review-queue">
      <div className="review-header">
        <h2>QC Review Queue</h2>
        <label className="reviewer-input">
          Reviewer:
          <input
            value={reviewerName}
            onChange={(e) => setReviewerName(e.target.value)}
            placeholder="Your name"
          />
        </label>
      </div>

      {queue.length === 0 ? (
        <div className="empty-state">
          <p>No decks pending review. All clear.</p>
        </div>
      ) : (
        <div className="review-list">
          {queue.map((deal) => (
            <div key={deal.id} className="review-card">
              <div className="review-card-left">
                <h3>{deal.property.name}</h3>
                <p className="address">
                  {deal.property.address.city}, {deal.property.address.state}
                </p>
                <div className="review-meta">
                  <span className="deck-type">
                    {deal.deck_type.replace(/_/g, ' ')}
                  </span>
                  <span className={`status status-${deal.status}`}>
                    {deal.status.replace(/_/g, ' ')}
                  </span>
                  {deal.review?.reviewer_id && (
                    <span className="reviewer">
                      Claimed by: {deal.review.reviewer_id}
                    </span>
                  )}
                </div>
              </div>
              <div className="review-card-actions">
                {deal.status === 'in_review' && deal.review?.reviewer_id ? (
                  <button className="btn-primary" onClick={() => onReview(deal)}>
                    Continue Review
                  </button>
                ) : (
                  <button className="btn-primary" onClick={() => handleClaim(deal)}>
                    Claim & Review
                  </button>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
