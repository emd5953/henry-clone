import { useState } from 'react';
import { DealList } from './components/DealList';
import { DealCreator } from './components/DealCreator';
import { DeckEditor } from './components/DeckEditor';
import { ReviewQueue } from './components/ReviewQueue';
import type { Deal } from './types';

type View = 'list' | 'create' | 'editor' | 'review' | 'review-editor';

export default function App() {
  const [view, setView] = useState<View>('list');
  const [activeDeal, setActiveDeal] = useState<Deal | null>(null);

  const openEditor = (deal: Deal) => {
    setActiveDeal(deal);
    setView('editor');
  };

  const openReviewEditor = (deal: Deal) => {
    setActiveDeal(deal);
    setView('review-editor');
  };

  return (
    <div className="app">
      <header className="app-header">
        <h1 onClick={() => setView('list')} style={{ cursor: 'pointer' }}>
          Henry Clone
        </h1>
        <nav>
          <button
            className={view === 'list' ? 'active' : ''}
            onClick={() => setView('list')}
          >
            Deals
          </button>
          <button
            className={view === 'create' ? 'active' : ''}
            onClick={() => setView('create')}
          >
            New Deal
          </button>
          <button
            className={view === 'review' ? 'active' : ''}
            onClick={() => setView('review')}
          >
            QC Review
          </button>
        </nav>
      </header>

      <main>
        {view === 'list' && <DealList onSelect={openEditor} />}
        {view === 'create' && (
          <DealCreator
            onCreated={(deal) => {
              setActiveDeal(deal);
              setView('editor');
            }}
          />
        )}
        {view === 'editor' && activeDeal && (
          <DeckEditor deal={activeDeal} onBack={() => setView('list')} />
        )}
        {view === 'review' && <ReviewQueue onReview={openReviewEditor} />}
        {view === 'review-editor' && activeDeal && (
          <DeckEditor
            deal={activeDeal}
            onBack={() => setView('review')}
            reviewMode
          />
        )}
      </main>
    </div>
  );
}
