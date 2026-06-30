import type { ReactNode } from 'react';
import { useTranslation } from 'react-i18next';

import type { GameMode } from './types';

type GameContextPanelProps = {
  children: ReactNode;
  currentMode: GameMode;
};

export function GameContextPanel({ children, currentMode }: GameContextPanelProps) {
  const { t } = useTranslation('game');

  return (
    <aside className="min-w-0 rounded border border-stone-800 bg-dusk-900/80 p-3 xl:max-h-[calc(100vh-11rem)] xl:overflow-y-auto">
      <div className="mb-3 flex flex-wrap items-center justify-between gap-2">
        <div>
          <p className="text-xs uppercase tracking-normal text-stone-500">{t('shell.context')}</p>
          <h2 className="text-lg font-semibold text-stone-100">{t(`navigation.${currentMode}`)}</h2>
        </div>
      </div>
      <div className="grid min-w-0 gap-3 break-words">{children}</div>
    </aside>
  );
}
