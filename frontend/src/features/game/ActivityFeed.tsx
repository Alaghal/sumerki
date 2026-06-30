import { useTranslation } from 'react-i18next';

import type { Army, Building, KingdomEvent, Mission, MissionReport, PatronPressure, Raid } from '../../api/client';
import { getLocalizedEventTitle, getLocalizedReportTitle } from '../../utils/localizedContent';
import type { GameMode } from './types';

type ActivityFeedProps = {
  activeEvents: KingdomEvent[];
  activeMissions: Mission[];
  activeRaids: Raid[];
  army: Army | null;
  buildings: Building[];
  patronPressure: PatronPressure | null;
  reports: MissionReport[];
  unreadReportsCount: number;
  formatDate: (value: string) => string;
  onModeChange: (mode: GameMode) => void;
};

export function ActivityFeed({
  activeEvents,
  activeMissions,
  activeRaids,
  army,
  buildings,
  formatDate,
  onModeChange,
  patronPressure,
  reports,
  unreadReportsCount,
}: ActivityFeedProps) {
  const { t } = useTranslation(['buildings', 'events', 'game', 'missions', 'patrons', 'reports', 'units']);
  const activeUpgrades = buildings.filter((building) => building.isUpgrading);
  const activeTraining = army?.trainingOrders ?? [];
  const latestReport = reports[0];
  const hasPatronWarning = patronPressure?.crisisStatus === 'warning' || patronPressure?.crisisStatus === 'active' || patronPressure?.crisisStatus === 'delayed';
  const hasActivity =
    activeUpgrades.length + activeTraining.length + activeMissions.length + activeRaids.length + activeEvents.length + unreadReportsCount + (hasPatronWarning ? 1 : 0) > 0;

  const items = [
    ...activeUpgrades.slice(0, 2).map((building) => ({
      detail: building.upgradeFinishesAt ? t('game:activity.endsAt', { date: formatDate(building.upgradeFinishesAt) }) : t('buildings:upgrading'),
      mode: 'city' as const,
      title: t('game:activity.construction'),
      value: t(`buildings:${building.type}.name`),
    })),
    ...activeTraining.slice(0, 2).map((order) => ({
      detail: t('game:activity.endsAt', { date: formatDate(order.finishesAt) }),
      mode: 'army' as const,
      title: t('game:activity.training'),
      value: t(`units:${order.unitType}.name`),
    })),
    ...activeMissions.slice(0, 2).map((mission) => ({
      detail: t('game:activity.resolvesAt', { date: formatDate(mission.finishesAt) }),
      mode: 'missions' as const,
      title: t('game:activity.missions'),
      value: t(`missions:${mission.missionKey}.name`, { defaultValue: mission.missionLabel }),
    })),
    ...activeRaids.slice(0, 2).map((raid) => ({
      detail: t('game:activity.resolvesAt', { date: formatDate(raid.arrivesAt) }),
      mode: 'raids' as const,
      title: t('game:activity.raids'),
      value: raid.defenderKingdomName,
    })),
    ...(activeEvents[0]
      ? [
          {
            detail: t('game:activity.activeEventsCount', { count: activeEvents.length }),
            mode: 'events' as const,
            title: t('game:activity.events'),
            value: getLocalizedEventTitle(t, activeEvents[0]),
          },
        ]
      : []),
    ...(unreadReportsCount > 0
      ? [
          {
            detail: latestReport ? getLocalizedReportTitle(t, latestReport) : t('game:activity.unreadReportsCount', { count: unreadReportsCount }),
            mode: 'reports' as const,
            title: t('game:activity.reports'),
            value: t('game:activity.unreadReports', { count: unreadReportsCount }),
          },
        ]
      : []),
    ...(hasPatronWarning && patronPressure
      ? [
          {
            detail: patronPressure.nextTributeAt ? t('game:activity.nextTributeAt', { date: formatDate(patronPressure.nextTributeAt) }) : t(`patrons:pressure.status.${patronPressure.crisisStatus}`),
            mode: 'patron' as const,
            title: t('game:activity.patron'),
            value: t('game:activity.patronWarning'),
          },
        ]
      : []),
  ];

  return (
    <section className="min-w-0 max-w-full overflow-hidden rounded border border-stone-800 bg-dusk-900/90 p-3">
      <div className="mb-2 flex min-w-0 flex-wrap items-center justify-between gap-2">
        <h2 className="break-words font-semibold text-stone-100">{t('game:activity.title')}</h2>
        {!hasActivity ? <span className="break-words text-sm text-stone-500">{t('game:activity.empty')}</span> : null}
      </div>
      <div className="grid min-w-0 gap-2 sm:grid-cols-2 xl:grid-cols-4">
        {items.map((item) => (
          <button
            className="min-w-0 max-w-full rounded border border-stone-800 bg-dusk-950 px-3 py-2 text-left text-sm text-stone-300 hover:border-dusk-gold/60 hover:bg-dusk-800"
            key={`${item.mode}-${item.title}-${item.value}`}
            onClick={() => onModeChange(item.mode)}
            type="button"
          >
            <span className="block break-words text-xs uppercase tracking-normal text-stone-500">{item.title}</span>
            <span className="mt-1 block break-words font-semibold text-stone-100">{item.value}</span>
            <span className="mt-1 block break-words text-xs text-stone-400">{item.detail}</span>
          </button>
        ))}
      </div>
    </section>
  );
}
