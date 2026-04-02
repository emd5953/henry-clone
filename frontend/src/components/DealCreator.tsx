import { useState } from 'react';
import { createDeal } from '../api';
import { Upload, Loader2 } from 'lucide-react';
import type { Deal } from '../types';

interface Props { onCreated: (deal: Deal) => void; }

export function DealCreator({ onCreated }: Props) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    try {
      const deal = await createDeal(new FormData(e.currentTarget));
      onCreated(deal);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create deal');
    } finally {
      setLoading(false);
    }
  };

  const inputClass = "w-full bg-white/[0.05] border border-white/[0.08] rounded-lg px-3 py-2 text-sm text-white placeholder-[#555] focus:outline-none focus:border-blue-500/50 focus:ring-1 focus:ring-blue-500/20 transition-all";
  const labelClass = "block text-xs font-medium text-[#888] mb-1.5";

  return (
    <div className="max-w-2xl">
      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Property Info */}
        <div className="rounded-xl border border-white/[0.08] bg-white/[0.02] p-5">
          <h3 className="text-xs font-semibold uppercase tracking-wider text-[#666] mb-4">Property Info</h3>
          <div className="grid grid-cols-2 gap-3">
            <div className="col-span-2 sm:col-span-1">
              <label className={labelClass}>Property Name</label>
              <input name="property_name" required placeholder="Bayshore Commerce Center" className={inputClass} />
            </div>
            <div>
              <label className={labelClass}>Asset Class</label>
              <select name="asset_class" required className={inputClass}>
                <option value="office">Office</option>
                <option value="multifamily">Multifamily</option>
                <option value="retail">Retail</option>
                <option value="industrial">Industrial</option>
                <option value="mixed_use">Mixed Use</option>
              </select>
            </div>
            <div className="col-span-2">
              <label className={labelClass}>Street</label>
              <input name="street" required placeholder="1250 Bayshore Blvd" className={inputClass} />
            </div>
            <div>
              <label className={labelClass}>City</label>
              <input name="city" required placeholder="San Francisco" className={inputClass} />
            </div>
            <div className="grid grid-cols-2 gap-3">
              <div>
                <label className={labelClass}>State</label>
                <input name="state" required placeholder="CA" maxLength={2} className={inputClass} />
              </div>
              <div>
                <label className={labelClass}>Zip</label>
                <input name="zip" required placeholder="94124" className={inputClass} />
              </div>
            </div>
            <div>
              <label className={labelClass}>Units</label>
              <input name="units" type="number" placeholder="10" className={inputClass} />
            </div>
            <div>
              <label className={labelClass}>Sq Ft</label>
              <input name="sq_ft" type="number" placeholder="11300" className={inputClass} />
            </div>
            <div>
              <label className={labelClass}>Year Built</label>
              <input name="year_built" type="number" placeholder="1985" className={inputClass} />
            </div>
            <div>
              <label className={labelClass}>Deck Type</label>
              <select name="deck_type" className={inputClass}>
                <option value="offering_memorandum">Offering Memorandum</option>
                <option value="broker_opinion_of_value">Broker Opinion of Value</option>
                <option value="investment_teaser">Investment Teaser</option>
                <option value="leasing_flyer">Leasing Flyer</option>
              </select>
            </div>
          </div>
        </div>

        {/* Thesis */}
        <div className="rounded-xl border border-white/[0.08] bg-white/[0.02] p-5">
          <h3 className="text-xs font-semibold uppercase tracking-wider text-[#666] mb-4">Investment Thesis</h3>
          <textarea name="thesis" rows={3} placeholder="Below-market rents with 20% vacancy create clear lease-up upside..." className={inputClass + " resize-none"} />
        </div>

        {/* File Uploads */}
        <div className="rounded-xl border border-white/[0.08] bg-white/[0.02] p-5">
          <h3 className="text-xs font-semibold uppercase tracking-wider text-[#666] mb-4">Upload Files</h3>
          <div className="grid grid-cols-2 gap-3">
            <div>
              <label className={labelClass}>Rent Roll</label>
              <div className="relative">
                <Upload className="absolute left-3 top-2.5 w-3.5 h-3.5 text-[#555]" />
                <input name="rent_roll" type="file" accept=".csv,.xlsx,.xls" className={inputClass + " pl-9 file:hidden text-xs"} />
              </div>
            </div>
            <div>
              <label className={labelClass}>T12 Financials</label>
              <div className="relative">
                <Upload className="absolute left-3 top-2.5 w-3.5 h-3.5 text-[#555]" />
                <input name="t12" type="file" accept=".csv,.xlsx,.xls" className={inputClass + " pl-9 file:hidden text-xs"} />
              </div>
            </div>
          </div>
          <div className="mt-3">
            <label className={labelClass}>Property Photos</label>
            <input name="photos" type="file" accept="image/*" multiple className={inputClass + " text-xs file:hidden"} />
            <p className="text-[10px] text-[#555] mt-1">AI analyzes photos to match the deck's color palette and style to the property.</p>
          </div>
        </div>

        {error && <div className="text-sm text-red-400 bg-red-500/10 border border-red-500/20 rounded-lg px-4 py-2">{error}</div>}

        <button
          type="submit"
          disabled={loading}
          className="w-full h-10 bg-white text-black font-semibold text-sm rounded-lg hover:bg-white/90 disabled:opacity-50 disabled:cursor-not-allowed transition-all flex items-center justify-center gap-2"
        >
          {loading ? <><Loader2 className="w-4 h-4 animate-spin" /> Generating Deck...</> : 'Generate Deck'}
        </button>
      </form>
    </div>
  );
}
