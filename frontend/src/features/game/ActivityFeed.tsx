import { useTranslation } from 'react-i18next';

import type { Army, Building, KingdomEvent, Mission, Raid } from '../../api/client';

type ActivityFeedProps = {
  activeEvents: KingdomEvent[];
  activeMissions: Mission[];
  activeRaids: Raid[];
  army: Army | null;
  buildings: Building[];
  unreadReportsCount: number;
};

export function ActivityFeed({ activeEvents, activeMissions, activeRaids, army, buildings, unreadReportsCount }: ActivityFeedProps) {
  const { t } = useTranslation('game');
  const activeUpgrades = buildings.filter((building) => building.isUpgrading).length;
  const activeTraining = army?.trainingOrders.length ?? 0;
  const hasActivity = activeUpgrades + activeTraining + activeMissions.length + activeRaids.length + activeEvents.length + unreadReportsCount > 0;

  const items = [
    t('activity.activeUpgradesCount', { count: activeUpgrades }),
    t('activity.activeTrainingCount', { count: activeTraining }),
    t('activity.activeMissionsCount', { count: activeMissions.length }),
    t('activity.activeRaidsCount', { count: activeRaids.length }),
    t('activity.activeEventsCount', { count: activeEvents.length }),
    t('activity.unreadReportsCount', { count: unreadReportsCount }),
  ];

  return (
    <section className="rounded border border-stone-800 bg-dusk-900/90 p-3">
      <div className="mb-2 flex flex-wrap items-center justify-between gap-2">
        <h2 className="font-semibold text-stone-100">{t('activity.title')}</h2>
        {!hasActivity ? <span className="text-sm text-stone-500">{t('activity.noActivity')}</span> : null}
      </div>
      <div className="flex flex-wrap gap-2">
        {items.map((item) => (
          <span className="rounded border border-stone-800 bg-dusk-950 px-3 py-2 text-sm text-stone-300" key={item}>
            {item}
          </span>
        ))}
      </div>
    </section>
  );
}
