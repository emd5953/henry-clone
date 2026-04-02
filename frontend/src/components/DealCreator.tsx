import { useState } from 'react';
import { createDeal } from '../api';
import { Loader2, FileUp, Image } from 'lucide-react';
import type { Deal } from '../types';

interface Props { onCreated: (deal: Deal) => void; }

export function DealCreator({ onCreated }: Props) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault(); setLoading(true); setError('');
    try { onCreated(await createDeal(new FormData(e.currentTarget))); }
    catch (err) { setError(err instanceof Error ? err.message : 'Failed'); }
    finally { setLoading(false); }
  };

  const inputCls = "w-full bg-white border border-[#E6E2DA] rounded-xl px-4 py-3 text-sm text-[#2D3A31] placeholder-[#B5B0A8] focus:outline-none focus:border-[#8C9A84] focus:ring-2 focus:ring-[#8C9A84]/10 transition-all";

  return (
    <div className="max-w-2xl mx-auto">
      <div className="text-center mb-8">
        <h1 className="text-2xl font-bold text-[#2D3A31] mb-2" style={{ fontFamily: "'Playfair Display', serif" }}>Create a new deal</h1>
        <p className="text-[#7A8578] text-sm">Upload your data and we'll generate a professional deck in minutes.</p>
      </div>

      <form onSubmit={handleSubmit} className="space-y-6">
        <div className="rounded-2xl border border-[#E6E2DA] bg-white p-6 space-y-4 shadow-[0_4px_6px_-1px_rgba(45,58,49,0.05)]">
          <h3 className="text-xs font-bold uppercase tracking-widest text-[#C27B66]">Property Details</h3>
          <div className="grid grid-cols-2 gap-4">
            <div className="col-span-2"><input name="property_name" required placeholder="Property Name" className={inputCls} /></div>
            <select name="asset_class" required className={inputCls}>
              <option value="office">Office</option><option value="multifamily">Multifamily</option>
              <option value="retail">Retail</option><option value="industrial">Industrial</option>
              <option value="mixed_use">Mixed Use</option>
            </select>
            <select name="deck_type" className={inputCls}>
              <option value="offering_memorandum">Offering Memorandum</option>
              <option value="broker_opinion_of_value">Broker Opinion of Value</option>
              <option value="investment_teaser">Investment Teaser</option>
              <option value="leasing_flyer">Leasing Flyer</option>
            </select>
            <div className="col-span-2"><input name="street" required placeholder="Street Address" className={inputCls} /></div>
            <input name="city" required placeholder="City" className={inputCls} />
            <div className="grid grid-cols-2 gap-4">
              <input name="state" required placeholder="State" maxLength={2} className={inputCls} />
              <input name="zip" required placeholder="Zip" className={inputCls} />
            </div>
            <input name="units" type="number" placeholder="Units" className={inputCls} />
            <input name="sq_ft" type="number" placeholder="Sq Ft" className={inputCls} />
            <input name="year_built" type="number" placeholder="Year Built" className={inputCls} />
          </div>
        </div>

        <div className="rounded-2xl border border-[#E6E2DA] bg-white p-6 shadow-[0_4px_6px_-1px_rgba(45,58,49,0.05)]">
          <h3 className="text-xs font-bold uppercase tracking-widest text-[#C27B66] mb-3">Investment Thesis</h3>
          <textarea name="thesis" rows={3} placeholder="Below-market rents with 20% vacancy create clear lease-up upside..." className={inputCls + " resize-none"} />
          <p className="text-[11px] text-[#B5B0A8] mt-2">2-3 sentences on why this deal is compelling. AI will expand this.</p>
        </div>

        <div className="rounded-2xl border border-[#E6E2DA] bg-white p-6 space-y-4 shadow-[0_4px_6px_-1px_rgba(45,58,49,0.05)]">
          <h3 className="text-xs font-bold uppercase tracking-widest text-[#C27B66]">Upload Files</h3>
          <div className="grid grid-cols-2 gap-4">
            <label className="flex flex-col items-center justify-center p-6 rounded-xl border-2 border-dashed border-[#E6E2DA] hover:border-[#8C9A84] bg-[#F9F8F4] cursor-pointer transition-all group">
              <FileUp className="w-6 h-6 text-[#DCCFC2] group-hover:text-[#8C9A84] mb-2 transition-colors" />
              <span className="text-xs font-medium text-[#7A8578] group-hover:text-[#2D3A31]">Rent Roll</span>
              <span className="text-[10px] text-[#B5B0A8] mt-1">CSV or Excel</span>
              <input name="rent_roll" type="file" accept=".csv,.xlsx,.xls" className="hidden" />
            </label>
            <label className="flex flex-col items-center justify-center p-6 rounded-xl border-2 border-dashed border-[#E6E2DA] hover:border-[#8C9A84] bg-[#F9F8F4] cursor-pointer transition-all group">
              <FileUp className="w-6 h-6 text-[#DCCFC2] group-hover:text-[#8C9A84] mb-2 transition-colors" />
              <span className="text-xs font-medium text-[#7A8578] group-hover:text-[#2D3A31]">T12 Financials</span>
              <span className="text-[10px] text-[#B5B0A8] mt-1">CSV or Excel</span>
              <input name="t12" type="file" accept=".csv,.xlsx,.xls" className="hidden" />
            </label>
          </div>
          <label className="flex flex-col items-center justify-center p-6 rounded-xl border-2 border-dashed border-[#E6E2DA] hover:border-[#8C9A84] bg-[#F9F8F4] cursor-pointer transition-all group">
            <Image className="w-6 h-6 text-[#DCCFC2] group-hover:text-[#8C9A84] mb-2 transition-colors" />
            <span className="text-xs font-medium text-[#7A8578] group-hover:text-[#2D3A31]">Property Photos</span>
            <span className="text-[10px] text-[#B5B0A8] mt-1">AI matches deck colors to your property</span>
            <input name="photos" type="file" accept="image/*" multiple className="hidden" />
          </label>
        </div>

        {error && <div className="text-sm text-red-600 bg-red-50 border border-red-200 rounded-xl px-4 py-3">{error}</div>}

        <button type="submit" disabled={loading}
          className="w-full h-12 bg-[#2D3A31] text-white font-semibold text-sm rounded-xl hover:bg-[#3D4A41] disabled:opacity-50 transition-all flex items-center justify-center gap-2 shadow-[0_4px_6px_-1px_rgba(45,58,49,0.15)]">
          {loading ? <><Loader2 className="w-4 h-4 animate-spin" />Generating Deck...</> : 'Generate Deck'}
        </button>
      </form>
    </div>
  );
}
