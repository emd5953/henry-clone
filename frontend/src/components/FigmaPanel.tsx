import { useState } from 'react';
import { linkFigmaFile, postFigmaComment } from '../api';
import { ExternalLink, Send } from 'lucide-react';
import type { Deal } from '../types';

interface Props { deal: Deal; onLinked: (fileKey: string, fileUrl: string) => void; }

export function FigmaPanel({ deal, onLinked }: Props) {
  const [figmaUrl, setFigmaUrl] = useState('');
  const [comment, setComment] = useState('');
  const [linking, setLinking] = useState(false);
  const [posting, setPosting] = useState(false);
  const isLinked = !!deal.figma_file_key;

  const extractFileKey = (url: string) => {
    const match = url.match(/figma\.com\/(?:design|file)\/([a-zA-Z0-9]+)/);
    return match ? match[1] : url;
  };

  return (
    <div className="mt-4 pt-4 border-t border-[#E6E2DA]">
      <p className="text-[10px] font-bold uppercase tracking-widest text-[#C27B66] mb-2 px-1">Figma</p>
      {isLinked ? (
        <div className="space-y-2">
          <a href={deal.figma_file_url} target="_blank" rel="noopener noreferrer"
            className="flex items-center justify-center gap-1.5 w-full py-2 text-xs font-medium bg-[#2D3A31] text-white rounded-lg hover:bg-[#3D4A41] transition-all">
            <ExternalLink className="w-3 h-3" />Open in Figma
          </a>
          <div className="flex gap-1">
            <input value={comment} onChange={(e) => setComment(e.target.value)} placeholder="Comment..."
              className="flex-1 bg-white border border-[#E6E2DA] rounded-lg px-2 py-1.5 text-[11px] text-[#2D3A31] placeholder-[#B5B0A8] focus:outline-none" />
            <button onClick={async () => { setPosting(true); await postFigmaComment(deal.id, comment); setComment(''); setPosting(false); }}
              disabled={posting || !comment.trim()} className="p-1.5 text-[#B5B0A8] hover:text-[#2D3A31] transition-colors disabled:opacity-30">
              <Send className="w-3 h-3" />
            </button>
          </div>
        </div>
      ) : (
        <div className="space-y-2">
          <input value={figmaUrl} onChange={(e) => setFigmaUrl(e.target.value)} placeholder="Paste Figma URL..."
            className="w-full bg-white border border-[#E6E2DA] rounded-lg px-3 py-1.5 text-[11px] text-[#2D3A31] placeholder-[#B5B0A8] focus:outline-none focus:border-[#8C9A84]" />
          <button onClick={async () => { setLinking(true); const r = await linkFigmaFile(deal.id, extractFileKey(figmaUrl)); onLinked(r.file_key, r.file_url); setLinking(false); }}
            disabled={linking || !figmaUrl.trim()} className="w-full py-1.5 text-xs font-medium border border-[#E6E2DA] text-[#7A8578] rounded-lg hover:bg-[#F2F0EB] hover:text-[#2D3A31] transition-all disabled:opacity-30">
            {linking ? 'Linking...' : 'Link Figma File'}
          </button>
        </div>
      )}
    </div>
  );
}
