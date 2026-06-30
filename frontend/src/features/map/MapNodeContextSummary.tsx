import { useTranslation } from 'react-i18next';
import type { ReactNode } from 'react';

import type {
  AvailableMission,
  Kingdom,
  KingdomEvent,
  Mission,
  Neighbor,
  PatronPressure,
  PatronStatus,
  Resources,
} from '../../api/client';
import type { DashboardFormatters } from '../dashboard/shared';
import type { SelectedContext } from './types';

type MapNodeContextSummaryProps = {
  activeEvents: KingdomEvent[];
  availableMissions: AvailableMission[];
  context: SelectedContext;
  kingdom: Kingdom;
  missions: Mission[];
  neighbors: Neighbor[];
  patronPressure: PatronPressure | null;
  patronStatus: PatronStatus | null;
  resources: Resources | null;
  unreadReportsCount: number;
  formatters: Pick<DashboardFormatters, 'resourceList'>;
};

export function MapNodeContextSummary({
  activeEvents,
  availableMissions,
  context,
  formatters,
  kingdom,
  missions,
  neighbors,
  patronPressure,
  patronStatus,
  resources,
  unreadReportsCount,
}: MapNodeContextSummaryProps) {
  const { t } = useTranslation(['common', 'game', 'kingdom', 'map', 'missions', 'patrons', 'raids']);

  if (context.kind === 'mission') {
    const activeMission = missions.find((mission) => mission.missionKey === context.missionKey && mission.status === 'active');
    const availableMission = availableMissions.find((mission) => mission.key === context.missionKey);
    return (
      <SummaryFrame title={t(`missions:${context.missionKey}.name`, { defaultValue: availableMission?.label ?? t('game:context.mission.title') })}>
        <SummaryLine label={t('game:context.mission.title')} value={availableMission ? t(`missions:types.${availableMission.type}`) : t('game:context.mission.available')} />
        <SummaryLine
          label={t('game:context.selected')}
          value={activeMission ? t('game:context.mission.active') : availableMission ? t('game:context.mission.available') : t('common:states.unknown')}
        />
        <p className="text-sm text-dusk-gold">{t('game:context.mission.suggestedAction')}</p>
      </SummaryFrame>
    );
  }

  if (context.kind === 'missions') {
    return (
      <SummaryFrame title={t('game:navigation.missions')}>
        <p className="text-sm text-stone-400">{t('game:modes.missions.summary')}</p>
      </SummaryFrame>
    );
  }

  if (context.kind === 'neighbor') {
    const neighbor = context.kingdomId ? neighbors.find((candidate) => candidate.kingdomId === context.kingdomId) : context.slot === 'neighbor_1' ? neighbors[0] : neighbors[1];
    return (
      <SummaryFrame title={neighbor?.name ?? t('game:context.neighbor.unknown')}>
        {neighbor ? (
          <>
            <SummaryLine label={t('game:dashboard.culture')} value={t(`kingdom:cultures.${neighbor.culture}.name`)} />
            <SummaryLine label={t('raids:powerLabel')} value={t(`raids:power.${neighbor.powerEstimate}`)} />
            <SummaryLine
              label={t('game:context.selected')}
              value={neighbor.canRaid ? t('game:context.neighbor.canRaid') : neighbor.blockedReason ? t(`raids:blockedReasons.${neighbor.blockedReason}`) : t('game:context.neighbor.cannotRaid')}
            />
          </>
        ) : (
          <p className="text-sm text-stone-400">{t('game:context.neighbor.unknown')}</p>
        )}
      </SummaryFrame>
    );
  }

  if (context.kind === 'raids') {
    return (
      <SummaryFrame title={t('game:navigation.raids')}>
        <p className="text-sm text-stone-400">{t('game:modes.raids.summary')}</p>
      </SummaryFrame>
    );
  }

  if (context.kind === 'patron') {
    const patron = patronStatus?.patron;
    return (
      <SummaryFrame title={t('game:context.patron.title')}>
        <SummaryLine label={t('game:dashboard.patron')} value={patron ? t(`patrons:${patron.key}.name`) : t('game:context.patron.noPatron')} />
        {patron ? <SummaryLine label={t('patrons:favor')} value={`${patron.favor}`} /> : null}
        {patronPressure ? (
          <SummaryLine label={t('game:context.patron.pressure')} value={t(`patrons:pressure.status.${patronPressure.crisisStatus}`)} />
        ) : null}
      </SummaryFrame>
    );
  }

  if (context.kind === 'events') {
    return (
      <SummaryFrame title={t('game:context.events.title')}>
        <SummaryLine label={t('game:context.events.activeCount')} value={`${activeEvents.length}`} />
        <SummaryLine label={t('game:context.reports.unreadCount')} value={`${unreadReportsCount}`} />
        {activeEvents[0] ? <p className="text-sm text-stone-400">{activeEvents[0].title}</p> : null}
      </SummaryFrame>
    );
  }

  if (context.kind === 'army') {
    return (
      <SummaryFrame title={t('game:navigation.army')}>
        <p className="text-sm text-stone-400">{t('game:modes.army.summary')}</p>
      </SummaryFrame>
    );
  }

  if (context.kind === 'reports') {
    return (
      <SummaryFrame title={t('game:navigation.reports')}>
        <SummaryLine label={t('game:context.reports.unreadCount')} value={`${unreadReportsCount}`} />
      </SummaryFrame>
    );
  }

  return (
    <SummaryFrame title={kingdom.name}>
      <SummaryLine label={t('game:dashboard.culture')} value={t(`kingdom:cultures.${kingdom.culture}.name`)} />
      <SummaryLine label={t('game:dashboard.patron')} value={patronStatus?.patron ? t(`patrons:${patronStatus.patron.key}.name`) : t('game:dashboard.noPatron')} />
      <p className="text-sm text-stone-400">{resources ? formatters.resourceList(resources) : t('common:states.loading')}</p>
    </SummaryFrame>
  );
}

type SummaryFrameProps = {
  children: ReactNode;
  title: string;
};

function SummaryFrame({ children, title }: SummaryFrameProps) {
  const { t } = useTranslation('game');

  return (
    <section className="grid min-w-0 max-w-full gap-3 overflow-hidden rounded border border-dusk-gold/40 bg-dusk-950 p-3">
      <div className="min-w-0">
        <p className="break-words text-xs uppercase tracking-normal text-dusk-gold">{t('context.title')}</p>
        <h3 className="mt-1 break-words font-semibold text-stone-100">{title}</h3>
      </div>
      <div className="grid min-w-0 gap-2 break-words">{children}</div>
    </section>
  );
}

type SummaryLineProps = {
  label: string;
  value: string;
};

function SummaryLine({ label, value }: SummaryLineProps) {
  return (
    <div className="flex min-w-0 flex-wrap justify-between gap-x-4 gap-y-1 text-sm">
      <span className="min-w-0 break-words text-stone-500">{label}</span>
      <span className="min-w-0 break-words text-left text-stone-100 sm:text-right">{value}</span>
    </div>
  );
}
