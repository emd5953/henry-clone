import { useState } from 'react';
import { FileText, Search, PlusCircle, ClipboardCheck } from 'lucide-react';
import { DealList } from './components/DealList';
import { DealCreator } from './components/DealCreator';
import { DeckEditor } from './components/DeckEditor';
import { ReviewQueue } from './components/ReviewQueue';
import { cn } from './lib/utils';
import type { Deal } from './types';

type View = 'list' | 'create' | 'editor' | 'review' | 'review-editor';

const navItems = [
  { name: 'Deals', view: 'list' as View, icon: Search },
  { name: 'New Deal', view: 'create' as View, icon: PlusCircle },
  { name: 'QC Review', view: 'review' as View, icon: ClipboardCheck },
];

export default function App() {
  const [view, setView] = useState<View>('list');
  const [activeDeal, setActiveDeal] = useState<Deal | null>(null);

  return (
    <div className="min-h-screen bg-[#F9F8F4]">
      <link href="https://fonts.googleapis.com/css2?family=Playfair+Display:wght@400;600;700&family=Source+Sans+3:wght@300;400;500;600;700&display=swap" rel="stylesheet" />

      {/* Top Nav */}
      <nav className="sticky top-0 z-50 border-b border-[#E6E2DA] bg-white/80 backdrop-blur-xl">
        <div className="max-w-6xl mx-auto px-6 h-16 flex items-center justify-between">
          <button onClick={() => setView('list')} className="flex items-center gap-2.5">
            <div className="w-8 h-8 rounded-lg bg-[#8C9A84] flex items-center justify-center">
              <FileText className="w-4 h-4 text-white" />
            </div>
            <span className="text-lg font-bold text-[#2D3A31] tracking-tight" style={{ fontFamily: "'Playfair Display', serif" }}>
              Henry<span className="text-[#C27B66]">Clone</span>
            </span>
          </button>

          <div className="flex items-center gap-1">
            {navItems.map((item) => (
              <button key={item.name} onClick={() => setView(item.view)}
                className={cn(
                  'flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all',
                  (view === item.view || (item.view === 'list' && view === 'editor'))
                    ? 'bg-[#8C9A84] text-white'
                    : 'text-[#7A8578] hover:text-[#2D3A31] hover:bg-[#F2F0EB]'
                )}>
                <item.icon className="w-4 h-4" />
                {item.name}
              </button>
            ))}
          </div>

          <div className="w-8 h-8 rounded-full bg-[#DCCFC2] flex items-center justify-center text-[#2D3A31] text-xs font-bold">E</div>
        </div>
      </nav>

      <main className="max-w-6xl mx-auto px-6 py-8">
        {view === 'list' && <DealList onSelect={(d) => { setActiveDeal(d); setView('editor'); }} />}
        {view === 'create' && <DealCreator onCreated={(d) => { setActiveDeal(d); setView('editor'); }} />}
        {view === 'editor' && activeDeal && <DeckEditor deal={activeDeal} onBack={() => setView('list')} />}
        {view === 'review' && <ReviewQueue onReview={(d) => { setActiveDeal(d); setView('review-editor'); }} />}
        {view === 'review-editor' && activeDeal && <DeckEditor deal={activeDeal} onBack={() => setView('review')} reviewMode />}
      </main>
    </div>
  );
}
