import { useEffect, useState, useRef, useCallback } from 'react';
import {
  getSections,
  updateSection,
  getDeckHTML,
  getDeckPDFUrl,
  reviewEdit,
  completeReview,
} from '../api';
import type { Deal, Section } from '../types';
import { FigmaPanel } from './FigmaPanel';

interface Props {
  deal: Deal;
  onBack: () => void;
  reviewMode?: boolean;
}

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

  useEffect(() => {
    getSections(deal.id).then(setSections);
  }, [deal.id]);

  const saveSection = useCallback(async () => {
    if (!contentRef.current || !titleRef.current) return;
    setSaving(true);

    const newTitle = titleRef.current.innerText;
    const newContent = contentRef.current.innerHTML;
    const data = { title: newTitle, content: newContent };

    try {
      const saveFn = reviewMode ? reviewEdit : updateSection;
      const updated = await saveFn(deal.id, activeIdx, data);
      const newSections = [...sections];
      newSections[activeIdx] = updated;
      setSections(newSections);
    } finally {
      setSaving(false);
    }
  }, [deal.id, activeIdx, sections, reviewMode]);

  // Auto-save on blur
  const handleBlur = useCallback(() => {
    saveSection();
  }, [saveSection]);

  const loadPreview = async () => {
    const html = await getDeckHTML(deal.id);
    setPreviewHTML(html);
    setShowPreview(true);
  };

  useEffect(() => {
    if (showPreview && iframeRef.current && previewHTML) {
      const doc = iframeRef.current.contentDocument;
      if (doc) {
        doc.open();
        doc.write(previewHTML);
        doc.close();
      }
    }
  }, [showPreview, previewHTML]);

  const handleApprove = async () => {
    setCompleting(true);
    try {
      await completeReview(deal.id, 'approved', reviewNotes);
      onBack();
    } finally {
      setCompleting(false);
    }
  };

  const handleReject = async () => {
    setCompleting(true);
    try {
      await completeReview(deal.id, 'needs_revision', reviewNotes);
      onBack();
    } finally {
      setCompleting(false);
    }
  };

  const activeSection = sections[activeIdx];

  return (
    <div className="deck-editor">
      <div className="editor-toolbar">
        <button onClick={onBack}>← Back</button>
        <h2>
          {deal.property.name}
          {reviewMode && <span className="review-badge">QC Review</span>}
        </h2>
        <div className="toolbar-actions">
          {saving && <span className="save-indicator">Saving...</span>}
          <button onClick={() => { setShowPreview(!showPreview); if (!showPreview) loadPreview(); }}>
            {showPreview ? 'Edit' : 'Preview'}
          </button>
          <a
            href={getDeckPDFUrl(deal.id)}
            className="btn-primary"
            target="_blank"
            rel="noopener noreferrer"
          >
            PDF
          </a>
        </div>
      </div>

      {showPreview ? (
        <div className="preview-container">
          <iframe ref={iframeRef} title="Deck Preview" className="deck-preview-iframe" />
        </div>
      ) : (
        <div className="editor-layout">
          <aside className="section-nav">
            <h3>Sections</h3>
            {sections.map((s, i) => (
              <button
                key={i}
                className={`section-nav-item ${i === activeIdx ? 'active' : ''}`}
                onClick={() => setActiveIdx(i)}
              >
                <span className="section-type">{s.type.replace(/_/g, ' ')}</span>
                <span className="section-title">{s.title}</span>
              </button>
            ))}

            {reviewMode && (
              <div className="review-panel">
                <h3>Review Actions</h3>
                <textarea
                  className="review-notes"
                  placeholder="QC notes..."
                  value={reviewNotes}
                  onChange={(e) => setReviewNotes(e.target.value)}
                  rows={4}
                />
                <div className="review-actions">
                  <button
                    className="btn-approve"
                    onClick={handleApprove}
                    disabled={completing}
                  >
                    ✓ Approve
                  </button>
                  <button
                    className="btn-reject"
                    onClick={handleReject}
                    disabled={completing}
                  >
                    ✗ Needs Revision
                  </button>
                </div>
              </div>
            )}

            <FigmaPanel
              deal={dealState}
              onLinked={(fileKey, fileUrl) => {
                setDealState({
                  ...dealState,
                  figma_file_key: fileKey,
                  figma_file_url: fileUrl,
                });
              }}
            />
          </aside>

          <div className="section-editor">
            {activeSection && (
              <>
                <div className="section-header">
                  <span className="section-type-badge">
                    {activeSection.type.replace(/_/g, ' ')}
                  </span>
                  <h3
                    ref={titleRef}
                    contentEditable
                    suppressContentEditableWarning
                    onBlur={handleBlur}
                    className="editable-title"
                  >
                    {activeSection.title}
                  </h3>
                </div>

                <div className="visual-editor-toolbar">
                  <button onClick={() => document.execCommand('bold')}>B</button>
                  <button onClick={() => document.execCommand('italic')}>I</button>
                  <button onClick={() => document.execCommand('underline')}>U</button>
                  <span className="separator" />
                  <button onClick={() => document.execCommand('insertUnorderedList')}>• List</button>
                  <button onClick={() => document.execCommand('insertOrderedList')}>1. List</button>
                  <span className="separator" />
                  <button onClick={() => document.execCommand('formatBlock', false, 'h3')}>H3</button>
                  <button onClick={() => document.execCommand('formatBlock', false, 'p')}>¶</button>
                  <span className="separator" />
                  <button onClick={() => document.execCommand('removeFormat')}>Clear</button>
                </div>

                <div
                  ref={contentRef}
                  className="visual-editor-content"
                  contentEditable
                  suppressContentEditableWarning
                  onBlur={handleBlur}
                  dangerouslySetInnerHTML={{ __html: activeSection.content }}
                />
              </>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
