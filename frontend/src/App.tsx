import { useState } from 'react';
import { LayoutDashboard, Plus, CheckCircle2, FileText } from 'lucide-react';
import { DealList } from './components/DealList';
import { DealCreator } from './components/DealCreator';
import { DeckEditor } from './components/DeckEditor';
import { ReviewQueue } from './components/ReviewQueue';
import { cn } from './lib/utils';
import type { Deal } from './types';

type View = 'list' | 'create' | 'editor' | 'review' | 'review-editor';

const navItems = [
  { name: 'Deals', icon: LayoutDashboard, view: 'list' as View },
  { name: 'New Deal', icon: Plus, view: 'create' as View },
  { name: 'QC Review', icon: CheckCircle2, view: 'review' as View },
];

export default function App() {
  const [view, setView] = useState<View>('list');
  const [activeDeal, setActiveDeal] = useState<Deal | null>(null);

  const openEditor = (deal: Deal) => { setActiveDeal(deal); setView('editor'); };
  const openReviewEditor = (deal: Deal) => { setActiveDeal(deal); setView('review-editor'); };

  const activeNavName = view === 'list' ? 'Deals'
    : view === 'create' ? 'New Deal'
    : view === 'review' ? 'QC Review'
    : view === 'editor' ? 'Deck Editor'
    : 'QC Review';

  return (
    <div className="min-h-screen bg-[#0a0a0a] text-foreground font-sans antialiased">
      {/* Sidebar */}
      <aside className="hidden md:fixed md:inset-y-0 md:flex md:w-60 md:flex-col border-r border-border/40 bg-[#0a0a0a]">
        <div className="flex items-center h-14 px-5 border-b border-border/40">
          <FileText className="w-5 h-5 text-blue-400 mr-2" />
          <h1 className="text-base font-semibold tracking-tight">Henry Clone</h1>
        </div>
        <nav className="flex-1 px-3 py-3 space-y-0.5">
          {navItems.map((item) => (
            <button
              key={item.name}
              onClick={() => setView(item.view)}
              className={cn(
                'w-full flex items-center gap-2.5 px-3 py-2 rounded-lg text-sm font-medium transition-all',
                (view === item.view || (item.view === 'list' && view === 'editor'))
                  ? 'bg-white/10 text-white'
                  : 'text-[#888] hover:bg-white/5 hover:text-white'
              )}
            >
              <item.icon className="w-4 h-4" />
              {item.name}
            </button>
          ))}
        </nav>
        <div className="p-3 border-t border-border/40">
          <div className="flex items-center gap-2.5 px-3 py-2">
            <div className="w-7 h-7 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center text-white text-xs font-semibold">E</div>
            <span className="text-sm text-[#888]">Enrin</span>
          </div>
        </div>
      </aside>

      {/* Main */}
      <div className="md:pl-60 flex flex-col min-h-screen">
        <header className="sticky top-0 z-10 flex items-center h-14 px-6 border-b border-border/40 bg-[#0a0a0a]/80 backdrop-blur-xl">
          <h2 className="text-sm font-semibold text-white">{activeNavName}</h2>
        </header>

        <main className="flex-1 p-6">
          {view === 'list' && <DealList onSelect={openEditor} />}
          {view === 'create' && <DealCreator onCreated={(deal) => { setActiveDeal(deal); setView('editor'); }} />}
          {view === 'editor' && activeDeal && <DeckEditor deal={activeDeal} onBack={() => setView('list')} />}
          {view === 'review' && <ReviewQueue onReview={openReviewEditor} />}
          {view === 'review-editor' && activeDeal && <DeckEditor deal={activeDeal} onBack={() => setView('review')} reviewMode />}
        </main>
      </div>
    </div>
  );
}
