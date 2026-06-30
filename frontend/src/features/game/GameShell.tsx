import type { ReactNode } from 'react';

import { GameContextPanel } from './GameContextPanel';
import { GameModeNavigation } from './GameModeNavigation';
import type { GameMode } from './types';

type GameShellProps = {
  activity: ReactNode;
  context: ReactNode;
  currentMode: GameMode;
  hud: ReactNode;
  scene: ReactNode;
  onModeChange: (mode: GameMode) => void;
};

export function GameShell({ activity, context, currentMode, hud, onModeChange, scene }: GameShellProps) {
  return (
    <div className="grid min-w-0 gap-3">
      {hud}
      <div className="grid min-w-0 gap-3 xl:grid-cols-[220px_minmax(0,1fr)_360px]">
        <GameModeNavigation currentMode={currentMode} onModeChange={onModeChange} />
        <main className="min-w-0 overflow-hidden">{scene}</main>
        <GameContextPanel currentMode={currentMode}>{context}</GameContextPanel>
      </div>
      {activity}
    </div>
  );
}
