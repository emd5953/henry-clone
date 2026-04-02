import { useEffect, useState, useRef, useCallback } from 'react';
import { getSections, updateSection, getDeckHTML, getDeckPDFUrl, reviewEdit, completeReview } from '../api';
import { ArrowLeft, Eye, Download, Bold, Italic, Underline, List, ListOrdered, Type, Eraser, Check, X } from 'lucide-react';
import { FigmaPanel } from './FigmaPanel';
import type { Deal, Section } from '../types';

interface Props { deal: Deal; onBack: () => void; reviewMode?: boolean; }

export function DeckEditor({ deal, onBack, reviewMode = false }: Props) {
  const [sections, setSections] = useState<Section[]>([]);
  const [activeIdx, setActiveIdx] = useState(0);
  const [previewHTML, setPreviewHTML] = useState('');
  const [showPreview, setShowPreview] = useState(false);
  const [saving, setSaving] = useState(false);
  const [reviewNotes, setReviewNotes] = useState('');
  const [completing, setCompleting] = useState(false);
  const [dealState, setDealState] = useState(deal);
  const iframeRef = useRef<HTMLIFrameElement>(null);
  const contentRef = useRef<HTMLDivElement>(null);
  const titleRef = useRef<HTMLHeadingElement>(null);

  useEffect(() => { getSections(deal.id).then(setSections); }, [deal.id]);

  const saveSection = useCallback(async () => {
    if (!contentRef.current || !titleRef.current) return;
    setSaving(true);
    try {
      const saveFn = reviewMode ? reviewEdit : updateSection;
      const updated = await saveFn(deal.id, activeIdx, {
        title: titleRef.current.innerText,
        content: contentRef.current.innerHTML,
      });
      const s = [...sections]; s[activeIdx] = updated; setSections(s);
    } finally { setSaving(false); }
  }, [deal.id, activeIdx, sections, reviewMode]);

  const handleBlur = useCallback(() => { saveSection(); }, [saveSection]);

  const loadPreview = async () => {
    const html = await getDeckHTML(deal.id);
    setPreviewHTML(html);
    setShowPreview(true);
  };

  useEffect(() => {
    if (showPreview && iframeRef.current && previewHTML) {
      const doc = iframeRef.current.contentDocument;
      if (doc) { doc.open(); doc.write(previewHTML); doc.close(); }
    }
  }, [showPreview, previewHTML]);

  const exec = (cmd: string, val?: string) => document.execCommand(cmd, false, val);
  const activeSection = sections[activeIdx];

  return (
    <div className="flex flex-col h-[calc(100vh-7rem)]">
      {/* Toolbar */}
      <div className="flex items-center gap-3 pb-4 border-b border-white/[0.06] mb-4">
        <button onClick={onBack} className="text-[#666] hover:text-white transition-colors"><ArrowLeft className="w-4 h-4" /></button>
        <h2 className="text-sm font-semibold text-white flex-1">
          {deal.property.name}
          {reviewMode && <span className="ml-2 text-[10px] font-semibold uppercase tracking-wider px-2 py-0.5 rounded-full bg-yellow-500/20 text-yellow-400">QC Review</span>}
        </h2>
        {saving && <span className="text-[11px] text-[#555] animate-pulse">Saving...</span>}
        <button onClick={() => { setShowPreview(!showPreview); if (!showPreview) loadPreview(); }}
          className="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-[#888] border border-white/[0.08] rounded-lg hover:bg-white/[0.05] transition-all">
          <Eye className="w-3.5 h-3.5" />{showPreview ? 'Edit' : 'Preview'}
        </button>
        <a href={getDeckPDFUrl(deal.id)} target="_blank" rel="noopener noreferrer"
          className="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium bg-white text-black rounded-lg hover:bg-white/90 transition-all">
          <Download className="w-3.5 h-3.5" />PDF
        </a>
      </div>

      {showPreview ? (
        <iframe ref={iframeRef} title="Deck Preview" className="flex-1 rounded-xl border border-white/[0.08] bg-white" />
      ) : (
        <div className="flex gap-4 flex-1 min-h-0">
          {/* Section Nav */}
          <div className="w-52 shrink-0 overflow-y-auto space-y-0.5">
            <p className="text-[10px] font-semibold uppercase tracking-wider text-[#555] mb-2 px-1">Sections</p>
            {sections.map((s, i) => (
              <button key={i} onClick={() => setActiveIdx(i)}
                className={`w-full text-left px-3 py-2 rounded-lg text-xs transition-all ${
                  i === activeIdx ? 'bg-white/10 text-white' : 'text-[#666] hover:bg-white/[0.04] hover:text-[#aaa]'
                }`}>
                <span className="block text-[9px] uppercase tracking-wider text-[#555]">{s.type.replace(/_/g, ' ')}</span>
                <span className="block font-medium truncate">{s.title}</span>
              </button>
            ))}

            {reviewMode && (
              <div className="mt-4 pt-4 border-t border-white/[0.06] space-y-2">
                <p className="text-[10px] font-semibold uppercase tracking-wider text-[#555] px-1">Review</p>
                <textarea value={reviewNotes} onChange={(e) => setReviewNotes(e.target.value)}
                  placeholder="QC notes..." rows={3}
                  className="w-full bg-white/[0.05] border border-white/[0.08] rounded-lg px-3 py-2 text-xs text-white placeholder-[#555] resize-none focus:outline-none focus:border-blue-500/50" />
                <div className="flex gap-1.5">
                  <button onClick={async () => { setCompleting(true); await completeReview(deal.id, 'approved', reviewNotes); setCompleting(false); onBack(); }}
                    disabled={completing} className="flex-1 flex items-center justify-center gap-1 py-1.5 text-xs font-semibold bg-emerald-500/20 text-emerald-400 rounded-lg hover:bg-emerald-500/30 transition-all disabled:opacity-50">
                    <Check className="w-3 h-3" />Approve
                  </button>
                  <button onClick={async () => { setCompleting(true); await completeReview(deal.id, 'needs_revision', reviewNotes); setCompleting(false); onBack(); }}
                    disabled={completing} className="flex-1 flex items-center justify-center gap-1 py-1.5 text-xs font-semibold bg-red-500/20 text-red-400 rounded-lg hover:bg-red-500/30 transition-all disabled:opacity-50">
                    <X className="w-3 h-3" />Reject
                  </button>
                </div>
              </div>
            )}

            <FigmaPanel deal={dealState} onLinked={(k, u) => setDealState({ ...dealState, figma_file_key: k, figma_file_url: u })} />
          </div>

          {/* Editor */}
          <div className="flex-1 rounded-xl border border-white/[0.08] bg-white/[0.02] overflow-y-auto">
            {activeSection && (
              <div className="p-6">
                <span className="text-[10px] font-semibold uppercase tracking-wider text-blue-400">{activeSection.type.replace(/_/g, ' ')}</span>
                <h3 ref={titleRef} contentEditable suppressContentEditableWarning onBlur={handleBlur}
                  className="text-lg font-semibold text-white mt-1 mb-4 outline-none border-b border-transparent focus:border-blue-500/30 pb-1 transition-colors">
                  {activeSection.title}
                </h3>

                {/* Formatting toolbar */}
                <div className="flex gap-0.5 mb-3 pb-3 border-b border-white/[0.06]">
                  {[
                    { icon: Bold, cmd: 'bold' }, { icon: Italic, cmd: 'italic' }, { icon: Underline, cmd: 'underline' },
                    null,
                    { icon: List, cmd: 'insertUnorderedList' }, { icon: ListOrdered, cmd: 'insertOrderedList' },
                    null,
                    { icon: Type, cmd: 'formatBlock', val: 'h3' }, { icon: Eraser, cmd: 'removeFormat' },
                  ].map((item, i) => item === null
                    ? <div key={i} className="w-px h-5 bg-white/[0.08] mx-1" />
                    : <button key={i} onClick={() => exec(item.cmd, item.val)}
                        className="p-1.5 rounded text-[#666] hover:text-white hover:bg-white/[0.08] transition-all">
                        <item.icon className="w-3.5 h-3.5" />
                      </button>
                  )}
                </div>

                <div ref={contentRef} contentEditable suppressContentEditableWarning onBlur={handleBlur}
                  dangerouslySetInnerHTML={{ __html: activeSection.content }}
                  className="prose prose-invert prose-sm max-w-none text-[#ccc] leading-relaxed outline-none min-h-[200px] [&_table]:w-full [&_table]:border-collapse [&_th]:text-left [&_th]:text-[10px] [&_th]:uppercase [&_th]:tracking-wider [&_th]:text-[#666] [&_th]:pb-2 [&_th]:border-b [&_th]:border-white/[0.08] [&_td]:py-2 [&_td]:border-b [&_td]:border-white/[0.04] [&_td]:text-sm" />
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
