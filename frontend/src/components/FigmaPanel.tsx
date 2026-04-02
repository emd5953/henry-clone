import { useState } from 'react';
import { linkFigmaFile, postFigmaComment } from '../api';
import type { Deal } from '../types';

interface Props {
  deal: Deal;
  onLinked: (fileKey: string, fileUrl: string) => void;
}

export function FigmaPanel({ deal, onLinked }: Props) {
  const [figmaUrl, setFigmaUrl] = useState('');
  const [comment, setComment] = useState('');
  const [linking, setLinking] = useState(false);
  const [posting, setPosting] = useState(false);

  const isLinked = !!deal.figma_file_key;

  // Extract file key from Figma URL
  // https://www.figma.com/design/ABC123/File-Name -> ABC123
  const extractFileKey = (url: string): string => {
    const match = url.match(/figma\.com\/(?:design|file)\/([a-zA-Z0-9]+)/);
    return match ? match[1] : url; // fallback to raw input if not a URL
  };

  const handleLink = async () => {
    setLinking(true);
    try {
      const fileKey = extractFileKey(figmaUrl);
      const result = await linkFigmaFile(deal.id, fileKey);
      onLinked(result.file_key, result.file_url);
    } finally {
      setLinking(false);
    }
  };

  const handleComment = async () => {
    if (!comment.trim()) return;
    setPosting(true);
    try {
      await postFigmaComment(deal.id, comment);
      setComment('');
    } finally {
      setPosting(false);
    }
  };

  return (
    <div className="figma-panel">
      <h3>Figma</h3>

      {isLinked ? (
        <>
          <a
            href={deal.figma_file_url}
            target="_blank"
            rel="noopener noreferrer"
            className="figma-link"
          >
            Open in Figma ↗
          </a>
          <p className="figma-hint">
            Edit the deck visually in Figma. Changes are made directly in the
            design file.
          </p>

          <div className="figma-comment">
            <textarea
              value={comment}
              onChange={(e) => setComment(e.target.value)}
              placeholder="Add a comment to the Figma file..."
              rows={2}
            />
            <button onClick={handleComment} disabled={posting || !comment.trim()}>
              {posting ? 'Posting...' : 'Comment'}
            </button>
          </div>
        </>
      ) : (
        <div className="figma-connect">
          <p className="figma-hint">
            Link a Figma file to edit this deck visually. Paste a Figma URL or
            file key.
          </p>
          <input
            value={figmaUrl}
            onChange={(e) => setFigmaUrl(e.target.value)}
            placeholder="https://www.figma.com/design/..."
          />
          <button
            onClick={handleLink}
            disabled={linking || !figmaUrl.trim()}
            className="btn-primary"
          >
            {linking ? 'Linking...' : 'Link Figma File'}
          </button>
        </div>
      )}
    </div>
  );
}
