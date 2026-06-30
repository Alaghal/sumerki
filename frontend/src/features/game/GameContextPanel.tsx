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
    <aside className="min-w-0 max-w-full overflow-hidden rounded border border-stone-800 bg-dusk-900/80 p-3 xl:max-h-[calc(100vh-11rem)] xl:overflow-y-auto">
      <div className="mb-3 flex flex-wrap items-center justify-between gap-2">
        <div className="min-w-0">
          <p className="break-words text-xs uppercase tracking-normal text-stone-500">{t('shell.context')}</p>
          <h2 className="break-words text-lg font-semibold text-stone-100">{t(`navigation.${currentMode}`)}</h2>
        </div>
      </div>
      <div className="grid min-w-0 gap-3 break-words [overflow-wrap:anywhere]">{children}</div>
    </aside>
  );
}
