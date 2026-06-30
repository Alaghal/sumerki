import { useTranslation } from 'react-i18next';

import { gameModes } from './gameModes';
import type { GameMode } from './types';

type GameModeNavigationProps = {
  currentMode: GameMode;
  onModeChange: (mode: GameMode) => void;
};

export function GameModeNavigation({ currentMode, onModeChange }: GameModeNavigationProps) {
  const { t } = useTranslation('game');

  return (
    <nav className="grid min-w-0 max-w-full gap-2 overflow-hidden rounded border border-stone-800 bg-dusk-900/90 p-3 xl:sticky xl:top-20 xl:self-start">
      <div className="break-words px-2 text-xs uppercase tracking-normal text-stone-500">{t('shell.modes')}</div>
      <div className="grid min-w-0 grid-cols-2 gap-2 sm:grid-cols-4 xl:grid-cols-1">
        {gameModes.map((mode) => {
          const isActive = currentMode === mode;
          return (
            <button
              className={
                isActive
                  ? 'min-w-0 rounded border border-dusk-gold bg-dusk-800 px-3 py-2 text-left text-sm font-semibold text-dusk-gold'
                  : 'min-w-0 rounded border border-stone-800 bg-dusk-950 px-3 py-2 text-left text-sm text-stone-300 hover:border-stone-700 hover:bg-dusk-800 hover:text-stone-100'
              }
              key={mode}
              onClick={() => onModeChange(mode)}
              type="button"
            >
              <span className="block break-words">{t(`navigation.${mode}`)}</span>
            </button>
          );
        })}
      </div>
    </nav>
  );
}
