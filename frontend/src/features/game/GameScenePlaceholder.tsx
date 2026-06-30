import { useTranslation } from 'react-i18next';

import type { Building, Kingdom, KingdomEvent, Mission, Neighbor, Raid } from '../../api/client';
import type { GameMode } from './types';

type GameScenePlaceholderProps = {
  activeEvents: KingdomEvent[];
  activeMissions: Mission[];
  activeRaids: Raid[];
  currentMode: GameMode;
  kingdom: Kingdom;
  neighbors: Neighbor[];
  upgradingBuildings: Building[];
};

const nodeClasses = [
  'left-[46%] top-[46%] border-dusk-gold text-dusk-gold',
  'left-[18%] top-[28%] border-stone-700 text-stone-300',
  'left-[68%] top-[24%] border-stone-700 text-stone-300',
  'left-[24%] top-[68%] border-stone-700 text-stone-300',
  'left-[72%] top-[66%] border-stone-700 text-stone-300',
];

export function GameScenePlaceholder({
  activeEvents,
  activeMissions,
  activeRaids,
  currentMode,
  kingdom,
  neighbors,
  upgradingBuildings,
}: GameScenePlaceholderProps) {
  const { t } = useTranslation('game');
  const nodeLabels = [
    kingdom.name,
    t('shell.nodes.blackForest'),
    t('shell.nodes.oldKurgan'),
    t('shell.nodes.dryFord'),
    t('shell.nodes.neighbors', { count: neighbors.length }),
  ];

  return (
    <section className="relative min-h-[28rem] overflow-hidden rounded border border-stone-800 bg-dusk-900/70 p-4">
      <div className="relative z-10 max-w-2xl">
        <p className="text-xs uppercase tracking-normal text-dusk-gold">{t('shell.scene')}</p>
        <h2 className="mt-1 text-2xl font-semibold text-stone-100">
          {currentMode === 'map' ? t('shell.mapPlaceholderTitle') : t(`navigation.${currentMode}`)}
        </h2>
        <p className="mt-2 text-sm text-stone-400">
          {currentMode === 'map' ? t('shell.mapPlaceholderDescription') : t(`modes.${currentMode}.summary`)}
        </p>
      </div>

      <div className="absolute inset-0 opacity-70">
        <div className="absolute left-[10%] right-[10%] top-1/2 border-t border-dashed border-stone-700" />
        <div className="absolute bottom-[18%] left-1/2 top-[18%] border-l border-dashed border-stone-700" />
        {nodeLabels.map((label, index) => (
          <div
            className={`absolute max-w-[9rem] -translate-x-1/2 -translate-y-1/2 rounded-full border bg-dusk-950/95 px-3 py-2 text-center text-xs shadow ${nodeClasses[index]}`}
            key={label}
          >
            {label}
          </div>
        ))}
      </div>

      <div className="absolute bottom-4 left-4 right-4 z-10 grid gap-2 sm:grid-cols-4">
        <div className="rounded border border-stone-800 bg-dusk-950/90 p-3">
          <div className="text-xs text-stone-500">{t('activity.activeMissions')}</div>
          <div className="text-lg font-semibold text-stone-100">{activeMissions.length}</div>
        </div>
        <div className="rounded border border-stone-800 bg-dusk-950/90 p-3">
          <div className="text-xs text-stone-500">{t('activity.activeRaids')}</div>
          <div className="text-lg font-semibold text-stone-100">{activeRaids.length}</div>
        </div>
        <div className="rounded border border-stone-800 bg-dusk-950/90 p-3">
          <div className="text-xs text-stone-500">{t('activity.activeEvents')}</div>
          <div className="text-lg font-semibold text-stone-100">{activeEvents.length}</div>
        </div>
        <div className="rounded border border-stone-800 bg-dusk-950/90 p-3">
          <div className="text-xs text-stone-500">{t('activity.activeUpgrades')}</div>
          <div className="text-lg font-semibold text-stone-100">{upgradingBuildings.length}</div>
        </div>
      </div>
    </section>
  );
}
