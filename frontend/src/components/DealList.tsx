import { useEffect, useState } from 'react';
import { listDeals } from '../api';
import type { Deal } from '../types';

interface Props {
  onSelect: (deal: Deal) => void;
}

export function DealList({ onSelect }: Props) {
  const [deals, setDeals] = useState<Deal[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    listDeals()
      .then(setDeals)
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <div className="loading">Loading deals...</div>;

  if (deals.length === 0) {
    return (
      <div className="empty-state">
        <p>No deals yet. Create your first deal to get started.</p>
      </div>
    );
  }

  return (
    <div className="deal-list">
      <h2>Your Deals</h2>
      <div className="deal-grid">
        {deals.map((deal) => (
          <div
            key={deal.id}
            className="deal-card"
            onClick={() => onSelect(deal)}
          >
            <div className="deal-card-header">
              <span className={`status status-${deal.status}`}>
                {deal.status}
              </span>
              <span className="deck-type">{deal.deck_type.replace(/_/g, ' ')}</span>
            </div>
            <h3>{deal.property.name}</h3>
            <p className="address">
              {deal.property.address.street}, {deal.property.address.city},{' '}
              {deal.property.address.state}
            </p>
            {deal.analysis && (
              <div className="deal-metrics">
                <span>NOI: ${deal.analysis.noi.toLocaleString()}</span>
                <span>
                  Occupancy: {(deal.analysis.occupancy_rate * 100).toFixed(0)}%
                </span>
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
