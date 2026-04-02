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
    <div className="mt-4 pt-4 border-t border-white/[0.06]">
      <p className="text-[10px] font-semibold uppercase tracking-wider text-[#555] mb-2 px-1">Figma</p>
      {isLinked ? (
        <div className="space-y-2">
          <a href={deal.figma_file_url} target="_blank" rel="noopener noreferrer"
            className="flex items-center justify-center gap-1.5 w-full py-2 text-xs font-medium bg-[#1e1e1e] text-white rounded-lg hover:bg-[#2a2a2a] transition-all">
            <ExternalLink className="w-3 h-3" />Open in Figma
          </a>
          <div className="flex gap-1">
            <input value={comment} onChange={(e) => setComment(e.target.value)} placeholder="Comment..."
              className="flex-1 bg-white/[0.05] border border-white/[0.08] rounded-lg px-2 py-1.5 text-[11px] text-white placeholder-[#555] focus:outline-none" />
            <button onClick={async () => { setPosting(true); await postFigmaComment(deal.id, comment); setComment(''); setPosting(false); }}
              disabled={posting || !comment.trim()} className="p-1.5 text-[#666] hover:text-white transition-colors disabled:opacity-30">
              <Send className="w-3 h-3" />
            </button>
          </div>
        </div>
      ) : (
        <div className="space-y-2">
          <input value={figmaUrl} onChange={(e) => setFigmaUrl(e.target.value)} placeholder="Paste Figma URL..."
            className="w-full bg-white/[0.05] border border-white/[0.08] rounded-lg px-3 py-1.5 text-[11px] text-white placeholder-[#555] focus:outline-none focus:border-blue-500/50" />
          <button onClick={async () => { setLinking(true); const r = await linkFigmaFile(deal.id, extractFileKey(figmaUrl)); onLinked(r.file_key, r.file_url); setLinking(false); }}
            disabled={linking || !figmaUrl.trim()} className="w-full py-1.5 text-xs font-medium bg-white/[0.08] text-[#888] rounded-lg hover:bg-white/[0.12] hover:text-white transition-all disabled:opacity-30">
            {linking ? 'Linking...' : 'Link Figma File'}
          </button>
        </div>
      )}
    </div>
  );
}
