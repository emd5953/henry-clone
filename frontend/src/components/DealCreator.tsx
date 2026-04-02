import { useState } from 'react';
import { createDeal } from '../api';
import type { Deal } from '../types';

interface Props {
  onCreated: (deal: Deal) => void;
}

export function DealCreator({ onCreated }: Props) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    const form = e.currentTarget;
    const formData = new FormData(form);

    try {
      const deal = await createDeal(formData);
      onCreated(deal);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create deal');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="deal-creator">
      <h2>Create New Deal</h2>
      <form onSubmit={handleSubmit}>
        <fieldset>
          <legend>Property Info</legend>
          <div className="form-row">
            <label>
              Property Name
              <input name="property_name" required placeholder="Bayshore Commerce Center" />
            </label>
            <label>
              Asset Class
              <select name="asset_class" required>
                <option value="office">Office</option>
                <option value="multifamily">Multifamily</option>
                <option value="retail">Retail</option>
                <option value="industrial">Industrial</option>
                <option value="mixed_use">Mixed Use</option>
              </select>
            </label>
          </div>
          <div className="form-row">
            <label>
              Street
              <input name="street" required placeholder="1250 Bayshore Blvd" />
            </label>
            <label>
              City
              <input name="city" required placeholder="San Francisco" />
            </label>
          </div>
          <div className="form-row">
            <label>
              State
              <input name="state" required placeholder="CA" maxLength={2} />
            </label>
            <label>
              Zip
              <input name="zip" required placeholder="94124" />
            </label>
          </div>
          <div className="form-row">
            <label>
              Units
              <input name="units" type="number" placeholder="10" />
            </label>
            <label>
              Sq Ft
              <input name="sq_ft" type="number" placeholder="11300" />
            </label>
            <label>
              Year Built
              <input name="year_built" type="number" placeholder="1985" />
            </label>
          </div>
        </fieldset>

        <fieldset>
          <legend>Deal Details</legend>
          <label>
            Deck Type
            <select name="deck_type">
              <option value="offering_memorandum">Offering Memorandum</option>
              <option value="broker_opinion_of_value">Broker Opinion of Value</option>
              <option value="investment_teaser">Investment Teaser</option>
              <option value="leasing_flyer">Leasing Flyer</option>
              <option value="syndication_deck">Syndication Deck</option>
            </select>
          </label>
          <label>
            Investment Thesis
            <textarea
              name="thesis"
              rows={3}
              placeholder="Strong value-add opportunity in emerging submarket..."
            />
          </label>
        </fieldset>

        <fieldset>
          <legend>Upload Files</legend>
          <div className="form-row">
            <label>
              Rent Roll (CSV or Excel)
              <input name="rent_roll" type="file" accept=".csv,.xlsx,.xls" />
            </label>
            <label>
              T12 Financials (CSV or Excel)
              <input name="t12" type="file" accept=".csv,.xlsx,.xls" />
            </label>
          </div>
        </fieldset>

        {error && <div className="error">{error}</div>}

        <button type="submit" disabled={loading} className="btn-primary">
          {loading ? 'Generating Deck...' : 'Generate Deck'}
        </button>
      </form>
    </div>
  );
}
