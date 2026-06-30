import { useTranslation } from 'react-i18next';

import type { Kingdom, PatronStatus, Resources } from '../../api/client';
import { resourceRows, type ResourceKey } from '../dashboard/shared';

type GameHudProps = {
  activeEventsCount: number;
  kingdom: Kingdom;
  patronStatus: PatronStatus | null;
  resources: Resources | null;
  resourcesLoading: boolean;
  unreadReportsCount: number;
  resourceLabel: (key: ResourceKey) => string;
};

export function GameHud({
  activeEventsCount,
  kingdom,
  patronStatus,
  resourceLabel,
  resources,
  resourcesLoading,
  unreadReportsCount,
}: GameHudProps) {
  const { t } = useTranslation(['common', 'game', 'kingdom', 'patrons']);
  const patron = patronStatus?.patron?.key ?? kingdom.patron;

  return (
    <section className="min-w-0 rounded border border-stone-800 bg-dusk-900/90 p-3">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div className="min-w-0">
          <h1 className="truncate text-xl font-semibold text-stone-100">{kingdom.name}</h1>
          <p className="text-sm text-stone-400">
            {t(`kingdom:cultures.${kingdom.culture}.name`)} · {patron ? t(`patrons:${patron}.name`) : t('game:dashboard.noPatron')}
          </p>
        </div>
        <div className="flex flex-wrap justify-start gap-2 sm:justify-end">
          <span className="rounded border border-stone-700 bg-dusk-950 px-2 py-1 text-xs text-stone-300">
            {t('game:hud.unreadReports', { count: unreadReportsCount })}
          </span>
          <span className="rounded border border-stone-700 bg-dusk-950 px-2 py-1 text-xs text-stone-300">
            {t('game:hud.activeEvents', { count: activeEventsCount })}
          </span>
        </div>
      </div>
      <div className="mt-3 flex flex-wrap gap-2">
        {resourceRows.map((key) => (
          <div className="min-w-[7rem] rounded border border-stone-800 bg-dusk-950 px-3 py-2" key={key}>
            <div className="text-xs text-stone-500">{resourceLabel(key)}</div>
            <div className="font-semibold text-stone-100">{resourcesLoading ? t('common:states.loading') : resources ? resources[key] : '0'}</div>
          </div>
        ))}
      </div>
    </section>
  );
}
