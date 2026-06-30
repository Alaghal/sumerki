import { useTranslation } from 'react-i18next';

import type {
  Army,
  AvailableMission,
  Building,
  BuildingType,
  Kingdom,
  KingdomEvent,
  Mission,
  MissionReport,
  Neighbor,
  PatronKey,
  PatronOption,
  PatronPressure,
  PatronStatus,
  Raid,
  Resources,
  Ruler,
  UnitType,
  User,
} from '../../api/client';
import { Button } from '../../components/ui/Button';
import { Card } from '../../components/ui/Card';
import {
  getLocalizedEventBody,
  getLocalizedEventChoiceDescription,
  getLocalizedEventChoiceLabel,
  getLocalizedEventResultBody,
  getLocalizedEventResultTitle,
  getLocalizedEventSelectedChoiceLabel,
  getLocalizedEventTitle,
  getLocalizedReportPhaseTitle,
  getLocalizedReportTitle,
} from '../../utils/localizedContent';
import { armyCostRows, DashboardFormatters, resourceRows, rulerStats, unitStatRows, unitTypes } from './shared';

type DashboardHeaderProps = {
  kingdom: Kingdom;
  user: User;
};

export function DashboardHeader({ kingdom, user }: DashboardHeaderProps) {
  const { t } = useTranslation('game');

  return (
    <div>
      <h1 className="text-2xl font-semibold text-stone-100">{t('dashboard.title', { name: kingdom.name })}</h1>
      <p className="mt-1 text-sm text-stone-400">{t('dashboard.subtitle', { email: user.email })}</p>
    </div>
  );
}

type DashboardRefreshButtonProps = {
  onRefresh: () => void;
};

export function DashboardRefreshButton({ onRefresh }: DashboardRefreshButtonProps) {
  const { t } = useTranslation('game');

  return (
    <Button className="justify-self-start" onClick={onRefresh} type="button">
      {t('dashboard.refreshCity')}
    </Button>
  );
}

type KingdomPanelProps = {
  kingdom: Kingdom;
  patronStatus: PatronStatus | null;
  user: User;
};

export function KingdomPanel({ kingdom, patronStatus, user }: KingdomPanelProps) {
  const { t } = useTranslation(['game', 'kingdom', 'patrons']);

  return (
    <Card title={t('game:dashboard.kingdom')}>
      <dl className="grid gap-2">
        <div className="flex flex-wrap justify-between gap-x-4 gap-y-1">
          <dt className="text-stone-400">{t('game:dashboard.culture')}</dt>
          <dd className="text-right text-stone-100">{t(`kingdom:cultures.${kingdom.culture}.name`)}</dd>
        </div>
        <div className="flex flex-wrap justify-between gap-x-4 gap-y-1">
          <dt className="text-stone-400">{t('game:dashboard.patron')}</dt>
          <dd className="text-right text-stone-100">
            {patronStatus?.patron
              ? t(`patrons:${patronStatus.patron.key}.name`)
              : kingdom.patron
                ? t(`patrons:${kingdom.patron}.name`)
                : t('game:dashboard.noPatron')}
          </dd>
        </div>
        <div className="flex flex-wrap justify-between gap-x-4 gap-y-1">
          <dt className="text-stone-400">{t('game:dashboard.player')}</dt>
          <dd className="min-w-0 max-w-full break-words text-left text-stone-100 sm:max-w-[14rem] sm:text-right">{user.email}</dd>
        </div>
      </dl>
    </Card>
  );
}

type PatronPanelProps = {
  crisisChoice: 'ask_delay' | 'break_patron' | null;
  formatters: Pick<DashboardFormatters, 'formatDate' | 'resourceList'>;
  isBreakingPatron: boolean;
  isPayingTribute: boolean;
  joiningPatron: PatronKey | null;
  loading: boolean;
  error: string;
  options: PatronOption[];
  pressure: PatronPressure | null;
  status: PatronStatus | null;
  onBreakPatron: () => void;
  onCrisisChoice: (choice: 'ask_delay' | 'break_patron') => void;
  onJoinPatron: (patron: PatronKey) => void;
  onPayTribute: () => void;
  onRefresh: () => void;
};

export function PatronPanel({
  crisisChoice,
  error,
  formatters,
  isBreakingPatron,
  isPayingTribute,
  joiningPatron,
  loading,
  onBreakPatron,
  onCrisisChoice,
  onJoinPatron,
  onPayTribute,
  onRefresh,
  options,
  pressure,
  status,
}: PatronPanelProps) {
  const { t } = useTranslation(['game', 'patrons']);

  return (
    <Card title={t('patrons:section.title')}>
      <div className="grid gap-4">
        {loading ? <p>{t('patrons:loading')}</p> : null}
        {error ? <p className="text-red-300">{error}</p> : null}
        {!loading ? (
          <>
            <div className="rounded border border-stone-800 bg-dusk-950 p-3">
              <div className="flex flex-wrap items-start justify-between gap-3">
                <div>
                  <h3 className="font-semibold text-stone-100">
                    {status?.patron ? t(`patrons:${status.patron.key}.name`) : t('patrons:notChosen')}
                  </h3>
                  <p className="mt-1 text-sm text-stone-400">{pressure?.summary ?? t('patrons:choosePath')}</p>
                </div>
                {status?.patron ? (
                  <div className="text-right text-sm text-stone-300">
                    <div>
                      {t('patrons:favor')}: {status.patron.favor}
                    </div>
                    <div>{t(`patrons:standing.${status.patron.standing}`)}</div>
                    <div>{formatters.formatDate(status.patron.joinedAt)}</div>
                  </div>
                ) : null}
              </div>
              {status?.patron ? (
                <Button className="mt-3 justify-self-start" disabled={isBreakingPatron} onClick={onBreakPatron} type="button">
                  {isBreakingPatron ? t('patrons:breakingTie') : t('patrons:breakTie')}
                </Button>
              ) : null}
            </div>

            {pressure ? (
              <div className="rounded border border-stone-800 bg-dusk-950 p-3">
                <div className="flex flex-wrap items-start justify-between gap-3">
                  <div>
                    <h3 className="font-semibold text-stone-100">{t('patrons:pressureTitle')}</h3>
                    <p className="mt-1 text-sm text-stone-400">{pressure.summary}</p>
                  </div>
                  <div className="text-right text-sm text-stone-300">
                    <div>{t(`patrons:pressure.status.${pressure.crisisStatus}`)}</div>
                    <div>{t('patrons:pressureLevel', { level: pressure.pressureLevel })}</div>
                  </div>
                </div>
                <dl className="mt-3 grid gap-2 text-sm">
                  <div className="flex flex-wrap justify-between gap-x-4 gap-y-1">
                    <dt className="text-stone-400">{t('patrons:tributeDebt')}</dt>
                    <dd className="text-right text-stone-100">{formatters.resourceList(pressure.tributeDebt)}</dd>
                  </div>
                  <div className="flex flex-wrap justify-between gap-x-4 gap-y-1">
                    <dt className="text-stone-400">{t('patrons:contributionDebt')}</dt>
                    <dd className="text-right text-stone-100">{formatters.resourceList(pressure.contributionDebt)}</dd>
                  </div>
                  <div className="flex flex-wrap justify-between gap-x-4 gap-y-1">
                    <dt className="text-stone-400">{t('patrons:nextTribute')}</dt>
                    <dd className="text-right text-stone-100">
                      {pressure.nextTributeAt ? formatters.formatDate(pressure.nextTributeAt) : t('patrons:noNextTribute')}
                    </dd>
                  </div>
                  {pressure.delayUntil ? (
                    <div className="flex flex-wrap justify-between gap-x-4 gap-y-1">
                      <dt className="text-stone-400">{t('patrons:delayUntil')}</dt>
                      <dd className="text-right text-stone-100">{formatters.formatDate(pressure.delayUntil)}</dd>
                    </div>
                  ) : null}
                  <div className="flex flex-wrap justify-between gap-x-4 gap-y-1">
                    <dt className="text-stone-400">{t('patrons:protectedMinimums')}</dt>
                    <dd className="text-right text-stone-100">
                      {pressure.protectedMinimums.gold ?? 0} / {pressure.protectedMinimums.food ?? 0} / {pressure.protectedMinimums.wood ?? 0} /{' '}
                      {pressure.protectedMinimums.stone ?? 0}
                    </dd>
                  </div>
                </dl>
                <div className="mt-3 flex flex-wrap gap-2">
                  {pressure.availableActions.includes('pay_tribute') ? (
                    <Button disabled={isPayingTribute} onClick={onPayTribute} type="button">
                      {isPayingTribute ? t('patrons:payingTribute') : t('patrons:payTribute')}
                    </Button>
                  ) : null}
                  {pressure.availableActions.includes('ask_delay') ? (
                    <Button disabled={crisisChoice === 'ask_delay'} onClick={() => onCrisisChoice('ask_delay')} type="button">
                      {crisisChoice === 'ask_delay' ? t('patrons:askingDelay') : t('patrons:askDelay')}
                    </Button>
                  ) : null}
                  {pressure.availableActions.includes('break_patron') ? (
                    <Button disabled={crisisChoice === 'break_patron'} onClick={() => onCrisisChoice('break_patron')} type="button">
                      {crisisChoice === 'break_patron' ? t('patrons:breakingTie') : t('patrons:breakDuringCrisis')}
                    </Button>
                  ) : null}
                </div>
              </div>
            ) : null}

            <div className="flex flex-wrap gap-2">
              <Button disabled={loading} onClick={onRefresh} type="button">
                {t('patrons:refresh')}
              </Button>
            </div>

            <div className="grid gap-3">
              {options.map((option) => {
                const isCurrent = status?.patron?.key === option.key;
                return (
                  <div className="rounded border border-stone-800 bg-dusk-950 p-3" key={option.key}>
                    <div className="flex flex-wrap items-start justify-between gap-3">
                      <div>
                        <h3 className="font-semibold text-stone-100">{t(`patrons:${option.key}.name`)}</h3>
                        <p className="mt-1 text-sm text-stone-400">{option.shortDescription}</p>
                        <p className="mt-2 text-sm text-stone-500">{option.flavor}</p>
                      </div>
                      <Button disabled={joiningPatron === option.key || isCurrent} onClick={() => onJoinPatron(option.key)} type="button">
                        {isCurrent ? t('patrons:chosen') : option.key === 'independent' ? t('patrons:choose') : t('patrons:join')}
                      </Button>
                    </div>
                    <div className="mt-3 grid gap-2 text-sm text-stone-400">
                      <div>
                        <p className="font-semibold text-stone-300">{t('patrons:currentEffects')}</p>
                        <ul className="mt-1 list-disc pl-5">
                          {option.currentEffects.map((effect) => (
                            <li key={effect}>{effect}</li>
                          ))}
                        </ul>
                      </div>
                      <div>
                        <p className="font-semibold text-stone-300">{t('patrons:futureEffects')}</p>
                        <ul className="mt-1 list-disc pl-5">
                          {option.futureEffects.map((effect) => (
                            <li key={effect}>{effect}</li>
                          ))}
                        </ul>
                      </div>
                    </div>
                  </div>
                );
              })}
            </div>
          </>
        ) : null}
      </div>
    </Card>
  );
}

type ResourcesPanelProps = {
  error: string;
  formatters: Pick<DashboardFormatters, 'resourceLabel'>;
  loading: boolean;
  resources: Resources | null;
  onRefresh: () => void;
};

export function ResourcesPanel({ error, formatters, loading, onRefresh, resources }: ResourcesPanelProps) {
  const { t } = useTranslation('resources');

  return (
    <Card title={t('section.title')}>
      <div className="grid gap-4">
        {loading ? <p>{t('loading')}</p> : null}
        {error ? <p className="text-red-300">{error}</p> : null}
        {resources && !loading && !error ? (
          <dl className="grid gap-2">
            {resourceRows.map((key) => (
              <div className="flex flex-wrap items-center justify-between gap-x-4 gap-y-1" key={key}>
                <dt className="text-stone-400">{formatters.resourceLabel(key)}</dt>
                <dd className="text-right">
                  <div className="font-semibold text-stone-100">{resources[key]}</div>
                  <div className="text-xs text-dusk-gold">{t('productionPerHour', { amount: resources.productionPerHour[key] })}</div>
                </dd>
              </div>
            ))}
          </dl>
        ) : null}
        <Button className="justify-self-start" disabled={loading} onClick={onRefresh} type="button">
          {t('refresh')}
        </Button>
      </div>
    </Card>
  );
}

type RulerPanelProps = {
  error: string;
  loading: boolean;
  ruler: Ruler | null;
};

export function RulerPanel({ error, loading, ruler }: RulerPanelProps) {
  const { t } = useTranslation(['game', 'kingdom']);

  return (
    <Card title={t('game:ruler.section.title')}>
      {loading ? <p>{t('game:ruler.loading')}</p> : null}
      {error ? <p className="text-red-300">{error}</p> : null}
      {ruler && !loading && !error ? (
        <div className="grid gap-4">
          <dl className="grid gap-2">
            <div className="flex flex-wrap justify-between gap-x-4 gap-y-1">
              <dt className="text-stone-400">{t('game:ruler.name')}</dt>
              <dd className="text-right text-stone-100">{ruler.name}</dd>
            </div>
            <div className="flex flex-wrap justify-between gap-x-4 gap-y-1">
              <dt className="text-stone-400">{t('game:ruler.age')}</dt>
              <dd className="text-right text-stone-100">{ruler.age}</dd>
            </div>
            <div className="flex flex-wrap justify-between gap-x-4 gap-y-1">
              <dt className="text-stone-400">{t('game:ruler.health')}</dt>
              <dd className="text-right text-stone-100">{t(`game:ruler.healthStatus.${ruler.healthStatus}`)}</dd>
            </div>
            <div className="flex flex-wrap justify-between gap-x-4 gap-y-1">
              <dt className="text-stone-400">{t('game:dashboard.culture')}</dt>
              <dd className="text-right text-stone-100">{t(`kingdom:cultures.${ruler.culture}.name`)}</dd>
            </div>
          </dl>
          <dl className="grid gap-2">
            {rulerStats.map((key) => (
              <div className="flex flex-wrap justify-between gap-x-4 gap-y-1" key={key}>
                <dt className="text-stone-400">{t(`game:ruler.stats.${key}`)}</dt>
                <dd className="text-right text-stone-100">{ruler[key]}</dd>
              </div>
            ))}
          </dl>
        </div>
      ) : null}
    </Card>
  );
}

type BuildingsPanelProps = {
  buildings: Building[];
  error: string;
  formatters: Pick<DashboardFormatters, 'costList' | 'formatDate'>;
  loading: boolean;
  upgradingType: BuildingType | null;
  onUpgrade: (buildingType: BuildingType) => void;
};

export function BuildingsPanel({ buildings, error, formatters, loading, onUpgrade, upgradingType }: BuildingsPanelProps) {
  const { t } = useTranslation('buildings');

  return (
    <Card title={t('section.title')}>
      <div className="grid gap-4">
        {loading ? <p>{t('loading')}</p> : null}
        {error ? <p className="text-red-300">{error}</p> : null}
        {!loading && !error ? (
          <div className="grid gap-3">
            {buildings.map((building) => (
              <div className="rounded border border-stone-800 bg-dusk-950 p-3" data-building-type={building.type} key={building.id}>
                <div className="flex flex-wrap items-start justify-between gap-3">
                  <div>
                    <h3 className="font-semibold text-stone-100">{t(`${building.type}.name`)}</h3>
                  </div>
                  <div className="text-right text-sm text-stone-300">{t('level', { level: building.level, maxLevel: building.maxLevel })}</div>
                </div>
                <div className="mt-3 grid gap-2">
                  {building.effects.map((effect) => (
                    <p className="text-sm text-stone-400" key={effect}>
                      {effect}
                    </p>
                  ))}
                  {building.isUpgrading ? (
                    <p className="text-dusk-gold">
                      {t('completedAt', {
                        date: building.upgradeFinishesAt ? formatters.formatDate(building.upgradeFinishesAt) : t('untilComplete'),
                      })}
                    </p>
                  ) : null}
                  {!building.isUpgrading && building.nextUpgrade ? (
                    <div className="grid gap-2">
                      {building.nextUpgrade.blockedReason === 'max_level' ? (
                        <p className="text-dusk-gold">{t('maxLevel')}</p>
                      ) : (
                        <>
                          <div className="grid gap-1">
                            <p className="text-stone-400">
                              {t('upgradeDetails', {
                                level: building.nextUpgrade.targetLevel,
                                seconds: building.nextUpgrade.durationSeconds,
                              })}
                            </p>
                            <p className="text-stone-400">
                              {t('cost')}: {formatters.costList(building.nextUpgrade.cost)}
                            </p>
                          </div>
                          <Button
                            className="justify-self-start"
                            data-building-upgrade={building.type}
                            disabled={upgradingType === building.type}
                            onClick={() => onUpgrade(building.type)}
                            type="button"
                          >
                            {upgradingType === building.type ? t('upgrading') : t('upgrade')}
                          </Button>
                        </>
                      )}
                    </div>
                  ) : null}
                </div>
              </div>
            ))}
          </div>
        ) : null}
      </div>
    </Card>
  );
}

type ArmyPanelProps = {
  army: Army | null;
  error: string;
  formatters: Pick<DashboardFormatters, 'costList' | 'formatDate' | 'unitLabel'>;
  isTraining: boolean;
  loading: boolean;
  trainingAmount: number;
  trainingType: UnitType;
  onTrain: () => void;
  onTrainingAmountChange: (amount: number) => void;
  onTrainingTypeChange: (unitType: UnitType) => void;
};

export function ArmyPanel({
  army,
  error,
  formatters,
  isTraining,
  loading,
  onTrain,
  onTrainingAmountChange,
  onTrainingTypeChange,
  trainingAmount,
  trainingType,
}: ArmyPanelProps) {
  const { t } = useTranslation('units');

  return (
    <Card title={t('section.title')}>
      <div className="grid gap-4">
        {loading ? <p>{t('loading')}</p> : null}
        {error ? <p className="text-red-300">{error}</p> : null}
        {army && !loading ? (
          <>
            <dl className="grid gap-2">
              <div className="flex flex-wrap justify-between gap-x-4 gap-y-1">
                <dt className="text-stone-400">{t('summary.total')}</dt>
                <dd className="text-right text-stone-100">{army.summary.totalUnits}</dd>
              </div>
              <div className="flex flex-wrap justify-between gap-x-4 gap-y-1">
                <dt className="text-stone-400">{t('stats.attack')}</dt>
                <dd className="text-right text-stone-100">{army.summary.totalAttack}</dd>
              </div>
              <div className="flex flex-wrap justify-between gap-x-4 gap-y-1">
                <dt className="text-stone-400">{t('stats.defense')}</dt>
                <dd className="text-right text-stone-100">{army.summary.totalDefense}</dd>
              </div>
              <div className="flex flex-wrap justify-between gap-x-4 gap-y-1">
                <dt className="text-stone-400">{t('stats.supply')}</dt>
                <dd className="text-right text-stone-100">{army.summary.totalSupply}</dd>
              </div>
            </dl>

            <div className="grid gap-3">
              {army.units.map((unit) => (
                <div className="rounded border border-stone-800 bg-dusk-950 p-3" data-unit-type={unit.type} key={unit.type}>
                  <div className="flex flex-wrap items-start justify-between gap-3">
                    <div>
                      <h3 className="font-semibold text-stone-100">{formatters.unitLabel(unit.type)}</h3>
                    </div>
                    <div className="text-right text-sm text-stone-300">{unit.amount}</div>
                  </div>
                  <dl className="mt-3 grid gap-1 text-sm">
                    {unitStatRows.map((key) => (
                      <div className="flex flex-wrap justify-between gap-x-4 gap-y-1" key={key}>
                        <dt className="text-stone-400">{t(`stats.${key}`)}</dt>
                        <dd className="text-right text-stone-100">{unit.stats[key]}</dd>
                      </div>
                    ))}
                  </dl>
                  <p className="mt-2 text-sm text-stone-400">
                    {t('cost')}: {formatters.costList(unit.cost, armyCostRows)}
                  </p>
                  <p className="mt-1 text-sm text-stone-400">{t('secondsPerUnit', { seconds: unit.secondsPerUnit })}</p>
                  <p className={unit.requirements.isMet ? 'mt-1 text-sm text-dusk-gold' : 'mt-1 text-sm text-red-300'}>
                    {unit.requirements.barracksLevel > 0
                      ? t('barracksRequirement', {
                          level: unit.requirements.barracksLevel,
                          status: unit.requirements.isMet ? t('requirementMet') : t('requirementNotMet'),
                        })
                      : t('requirementMet')}
                  </p>
                </div>
              ))}
            </div>

            <div className="grid gap-3 rounded border border-stone-800 bg-dusk-950 p-3">
              <div className="grid min-w-0 gap-2 sm:grid-cols-2 lg:grid-cols-[minmax(0,1fr)_8rem_auto]">
                <label className="grid min-w-0 gap-1 text-sm text-stone-400">
                  {t('unitType')}
                  <select
                    className="min-w-0 rounded border border-stone-700 bg-dusk-900 px-3 py-2 text-stone-100"
                    disabled={isTraining}
                    onChange={(event) => onTrainingTypeChange(event.target.value as UnitType)}
                    value={trainingType}
                  >
                    {army.units.map((unit) => (
                      <option key={unit.type} value={unit.type}>
                        {formatters.unitLabel(unit.type)}
                      </option>
                    ))}
                  </select>
                </label>
                <label className="grid min-w-0 gap-1 text-sm text-stone-400">
                  {t('amount')}
                  <input
                    className="min-w-0 rounded border border-stone-700 bg-dusk-900 px-3 py-2 text-stone-100"
                    disabled={isTraining}
                    max={50}
                    min={1}
                    onChange={(event) => onTrainingAmountChange(Number(event.target.value))}
                    type="number"
                    value={trainingAmount}
                  />
                </label>
                <Button className="self-end" disabled={isTraining} onClick={onTrain} type="button">
                  {isTraining ? t('trainingNow') : t('train')}
                </Button>
              </div>
              <p className="text-sm text-stone-400">{t('futureSystems')}</p>
            </div>

            <div className="grid gap-2">
              <h3 className="font-semibold text-stone-100">{t('training')}</h3>
              {army.trainingOrders.length === 0 ? <p className="text-sm text-stone-400">{t('noTraining')}</p> : null}
              {army.trainingOrders.map((order) => (
                <div className="rounded border border-stone-800 bg-dusk-950 p-3" key={order.id}>
                  <div className="flex flex-wrap justify-between gap-3">
                    <div>
                      <p className="font-semibold text-stone-100">{formatters.unitLabel(order.unitType)}</p>
                      <p className="text-sm text-stone-400">
                        {t('amount')}: {order.amount}
                      </p>
                    </div>
                    <div className="text-right text-sm text-dusk-gold">{t('completedAt', { date: formatters.formatDate(order.finishesAt) })}</div>
                  </div>
                </div>
              ))}
            </div>
          </>
        ) : null}
      </div>
    </Card>
  );
}

type MissionsPanelProps = {
  availableMissions: AvailableMission[];
  error: string;
  formatters: Pick<DashboardFormatters, 'formatDate' | 'resourceList' | 'unitLabel' | 'unitList'>;
  loading: boolean;
  missionInputs: Record<string, Partial<Record<UnitType, number>>>;
  missions: Mission[];
  startingMissionKey: string | null;
  onMissionUnitAmountChange: (missionKey: string, unitType: UnitType, amount: number) => void;
  onRefresh: () => void;
  onStartMission: (missionKey: string) => void;
};

export function MissionsPanel({
  availableMissions,
  error,
  formatters,
  loading,
  missionInputs,
  missions,
  onMissionUnitAmountChange,
  onRefresh,
  onStartMission,
  startingMissionKey,
}: MissionsPanelProps) {
  const { t } = useTranslation(['missions', 'reports']);

  return (
    <Card title={t('missions:section.title')}>
      <div className="grid gap-4">
        {loading ? <p>{t('missions:loading')}</p> : null}
        {error ? <p className="text-red-300">{error}</p> : null}
        <Button className="justify-self-start" disabled={loading} onClick={onRefresh} type="button">
          {t('missions:refresh')}
        </Button>
        {!loading ? (
          <>
            <div className="grid gap-3">
              <h3 className="font-semibold text-stone-100">{t('missions:availableMissions')}</h3>
              {availableMissions.length === 0 ? <p className="text-sm text-stone-400">{t('missions:noAvailable')}</p> : null}
              {availableMissions.map((mission) => (
                <div className="rounded border border-stone-800 bg-dusk-950 p-3" data-mission-key={mission.key} key={mission.key}>
                  <div className="flex flex-wrap items-start justify-between gap-3">
                    <div>
                      <h3 className="font-semibold text-stone-100">{t(`missions:${mission.key}.name`, { defaultValue: mission.label })}</h3>
                      <p className="text-xs text-stone-500">{t(`missions:types.${mission.type}`)}</p>
                    </div>
                    <div className="text-right text-sm text-stone-300">{t('missions:duration', { seconds: mission.durationSeconds })}</div>
                  </div>
                  <p className="mt-2 text-sm text-stone-400">{mission.description}</p>
                  <p className="mt-2 text-sm text-stone-400">
                    {t('missions:risk')}: {mission.risk}. {t('missions:minimum')}:{' '}
                    {t('missions:minimumUnits', {
                      scouts: mission.minimumRequirements.scouts,
                      total: mission.minimumRequirements.totalUnits,
                    })}
                  </p>
                  <p className="mt-1 text-sm text-stone-400">
                    {t('missions:rewards')}: {formatters.resourceList(mission.baseRewards)}
                  </p>
                  <div className="mt-3 grid min-w-0 grid-cols-2 gap-2 sm:grid-cols-3 lg:grid-cols-5">
                    {unitTypes.map((unitType) => (
                      <label className="grid min-w-0 gap-1 text-xs text-stone-400" key={unitType}>
                        {formatters.unitLabel(unitType)}
                        <input
                          className="min-w-0 rounded border border-stone-700 bg-dusk-900 px-2 py-2 text-stone-100"
                          min={0}
                          onChange={(event) => onMissionUnitAmountChange(mission.key, unitType, Number(event.target.value))}
                          type="number"
                          value={missionInputs[mission.key]?.[unitType] ?? 0}
                        />
                      </label>
                    ))}
                  </div>
                  <Button
                    className="mt-3 justify-self-start"
                    disabled={startingMissionKey === mission.key}
                    onClick={() => onStartMission(mission.key)}
                    type="button"
                  >
                    {startingMissionKey === mission.key ? t('missions:sending') : t('missions:send')}
                  </Button>
                </div>
              ))}
            </div>

            <div className="grid gap-2">
              <h3 className="font-semibold text-stone-100">{t('missions:activeMissions')}</h3>
              {missions.length === 0 ? <p className="text-sm text-stone-400">{t('missions:noActive')}</p> : null}
              {missions.map((mission) => (
                <div className="rounded border border-stone-800 bg-dusk-950 p-3" key={mission.id}>
                  <div className="flex flex-wrap justify-between gap-3">
                    <div>
                      <p className="font-semibold text-stone-100">
                        {t(`missions:${mission.missionKey}.name`, { defaultValue: mission.missionLabel })}
                      </p>
                      <p className="text-sm text-stone-400">{t(`missions:status.${mission.status}`)}</p>
                    </div>
                    <div className="text-right text-sm text-dusk-gold">
                      {mission.completedAt
                        ? t('missions:completedAt', { date: formatters.formatDate(mission.completedAt) })
                        : t('missions:resolvesAt', { date: formatters.formatDate(mission.finishesAt) })}
                    </div>
                  </div>
                  <p className="mt-2 text-sm text-stone-400">
                    {t('missions:sentUnits')}: {formatters.unitList(mission.units)}
                  </p>
                  {mission.result ? (
                    <p className="mt-1 text-sm text-stone-400">
                      {t('missions:result')}: {t(`reports:results.${mission.result.result}`)}
                    </p>
                  ) : null}
                </div>
              ))}
            </div>
          </>
        ) : null}
      </div>
    </Card>
  );
}

type RaidsPanelProps = {
  error: string;
  formatters: Pick<DashboardFormatters, 'formatDate' | 'resourceList' | 'unitLabel' | 'unitList'>;
  isStartingRaid: boolean;
  loading: boolean;
  neighbors: Neighbor[];
  raidInputs: Partial<Record<UnitType, number>>;
  raids: Raid[];
  selectedRaidTargetID: string | null;
  onRaidUnitAmountChange: (unitType: UnitType, amount: number) => void;
  onRefresh: () => void;
  onSelectRaidTargetID: (kingdomID: string) => void;
  onStartRaid: () => void;
};

export function RaidsPanel({
  error,
  formatters,
  isStartingRaid,
  loading,
  neighbors,
  onRaidUnitAmountChange,
  onRefresh,
  onSelectRaidTargetID,
  onStartRaid,
  raidInputs,
  raids,
  selectedRaidTargetID,
}: RaidsPanelProps) {
  const { t } = useTranslation(['game', 'kingdom', 'patrons', 'raids', 'reports', 'missions']);

  return (
    <Card title={t('raids:section.title')}>
      <div className="grid gap-4">
        {loading ? <p>{t('raids:loading')}</p> : null}
        {error ? <p className="text-red-300">{error}</p> : null}
        <Button className="justify-self-start" disabled={loading} onClick={onRefresh} type="button">
          {t('raids:refresh')}
        </Button>
        {!loading ? (
          <>
            <div className="grid gap-3">
              <h3 className="font-semibold text-stone-100">{t('raids:neighbors')}</h3>
              {neighbors.length === 0 ? <p className="text-sm text-stone-400">{t('raids:noNeighbors')}</p> : null}
              {neighbors.map((neighbor) => (
                <div className="rounded border border-stone-800 bg-dusk-950 p-3" key={neighbor.kingdomId}>
                  <div className="flex flex-wrap items-start justify-between gap-3">
                    <div>
                      <h3 className="font-semibold text-stone-100">{neighbor.name}</h3>
                      <p className="text-sm text-stone-400">{t(`kingdom:cultures.${neighbor.culture}.name`)}</p>
                      <p className="text-sm text-stone-400">
                        {t('game:dashboard.patron')}: {neighbor.patron ? t(`patrons:${neighbor.patron}.name`) : t('game:dashboard.noPatron')}
                      </p>
                    </div>
                    <div className="text-right text-sm text-stone-300">
                      <div>
                        {t('raids:dread')}: {neighbor.dread}
                      </div>
                      <div>{t(`raids:power.${neighbor.powerEstimate}`)}</div>
                    </div>
                  </div>
                  <p className={neighbor.canRaid ? 'mt-2 text-sm text-dusk-gold' : 'mt-2 text-sm text-stone-500'}>
                    {neighbor.canRaid
                      ? t('raids:canRaid')
                      : neighbor.blockedReason
                        ? t(`raids:blockedReasons.${neighbor.blockedReason}`)
                        : t('raids:cannotRaid')}
                  </p>
                  {neighbor.canRaid ? (
                    <label className="mt-3 flex items-center gap-2 text-sm text-stone-400">
                      <input checked={selectedRaidTargetID === neighbor.kingdomId} onChange={() => onSelectRaidTargetID(neighbor.kingdomId)} type="radio" />
                      {t('raids:selectTarget')}
                    </label>
                  ) : null}
                </div>
              ))}
            </div>

            {selectedRaidTargetID ? (
              <div className="grid gap-3 rounded border border-stone-800 bg-dusk-950 p-3">
                <h3 className="font-semibold text-stone-100">{t('raids:sendParty')}</h3>
                <div className="grid min-w-0 grid-cols-2 gap-2 sm:grid-cols-3 lg:grid-cols-5">
                  {unitTypes.map((unitType) => (
                    <label className="grid min-w-0 gap-1 text-xs text-stone-400" key={unitType}>
                      {formatters.unitLabel(unitType)}
                      <input
                        className="min-w-0 rounded border border-stone-700 bg-dusk-900 px-2 py-2 text-stone-100"
                        min={0}
                        onChange={(event) => onRaidUnitAmountChange(unitType, Number(event.target.value))}
                        type="number"
                        value={raidInputs[unitType] ?? 0}
                      />
                    </label>
                  ))}
                </div>
                <p className="text-sm text-stone-400">{t('raids:minimumHint')}</p>
                <Button className="justify-self-start" disabled={isStartingRaid} onClick={onStartRaid} type="button">
                  {isStartingRaid ? t('raids:startingRaid') : t('raids:startRaid')}
                </Button>
              </div>
            ) : null}

            <div className="grid gap-2">
              <h3 className="font-semibold text-stone-100">{t('raids:activeRaids')}</h3>
              {raids.length === 0 ? <p className="text-sm text-stone-400">{t('raids:noActiveRaids')}</p> : null}
              {raids.map((raid) => (
                <div className="rounded border border-stone-800 bg-dusk-950 p-3" key={raid.id}>
                  <div className="flex flex-wrap justify-between gap-3">
                    <div>
                      <p className="font-semibold text-stone-100">
                        {raid.attackerKingdomName} → {raid.defenderKingdomName}
                      </p>
                      <p className="text-sm text-stone-400">{t(`missions:status.${raid.status}`)}</p>
                      {raid.result ? <p className="text-sm text-dusk-gold">{t(`reports:results.${raid.result}`)}</p> : null}
                    </div>
                    <div className="text-right text-sm text-stone-400">
                      {raid.completedAt
                        ? t('raids:completedAt', { date: formatters.formatDate(raid.completedAt) })
                        : t('raids:arrivesAt', { date: formatters.formatDate(raid.arrivesAt) })}
                    </div>
                  </div>
                  <p className="mt-2 text-sm text-stone-400">
                    {t('raids:sentUnits')}: {formatters.unitList(raid.units)}
                  </p>
                  <p className="mt-1 text-sm text-stone-400">
                    {t('reports:loot')}: {formatters.resourceList(raid.loot)}
                  </p>
                </div>
              ))}
            </div>
          </>
        ) : null}
      </div>
    </Card>
  );
}

type EventsPanelProps = {
  choosingEventID: string | null;
  error: string;
  events: KingdomEvent[];
  formatDate: (value: string) => string;
  loading: boolean;
  onChooseEvent: (eventID: string, choiceKey: string) => void;
  onRefresh: () => void;
};

export function EventsPanel({ choosingEventID, error, events, formatDate, loading, onChooseEvent, onRefresh }: EventsPanelProps) {
  const { t } = useTranslation(['events', 'common']);
  const activeEvents = events.filter((event) => event.status === 'active');
  const resolvedEvents = events.filter((event) => event.status !== 'active');

  return (
    <Card title={t('section.title')}>
      <div className="grid gap-3">
        <div className="flex flex-wrap items-center justify-between gap-3">
          <p className="text-sm text-stone-400">{t('subtitle')}</p>
          <Button className="justify-self-start" disabled={loading} onClick={onRefresh} type="button">
            {t('refresh')}
          </Button>
        </div>
        {loading ? <p>{t('loading')}</p> : null}
        {error ? <p className="text-red-300">{error}</p> : null}
        {!loading && events.length === 0 ? <p className="text-sm text-stone-400">{t('noEvents')}</p> : null}
        {!loading && activeEvents.length > 0 ? (
          <div className="grid gap-3">
            {activeEvents.map((event) => (
              <div className="rounded border border-dusk-gold bg-dusk-950 p-3" key={event.id}>
                <div className="flex flex-wrap items-start justify-between gap-3">
                  <div>
                    <p className="text-xs uppercase tracking-wide text-dusk-gold">{t(`categories.${event.category}`)}</p>
                    <h3 className="mt-1 font-semibold text-stone-100">{getLocalizedEventTitle(t, event)}</h3>
                  </div>
                  <div className="text-right text-sm text-stone-400">
                    <div>{t(`status.${event.status}`)}</div>
                    <div>{t('expiresAt', { date: formatDate(event.expiresAt) })}</div>
                  </div>
                </div>
                <p className="mt-2 text-sm text-stone-400">{getLocalizedEventBody(t, event)}</p>
                <div className="mt-3 grid gap-2">
                  {event.choices.map((choice) => (
                    <div className="rounded border border-stone-800 bg-dusk-900 p-3" key={`${event.id}-${choice.key}`}>
                      <div className="flex flex-wrap items-start justify-between gap-3">
                        <div>
                          <h4 className="font-semibold text-stone-100">{getLocalizedEventChoiceLabel(t, event, choice)}</h4>
                          <p className="mt-1 text-sm text-stone-400">{getLocalizedEventChoiceDescription(t, event, choice)}</p>
                        </div>
                        <Button disabled={choosingEventID === event.id} onClick={() => onChooseEvent(event.id, choice.key)} type="button">
                          {choosingEventID === event.id ? t('choosing') : t('choose')}
                        </Button>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            ))}
          </div>
        ) : null}
        {!loading && resolvedEvents.length > 0 ? (
          <div className="grid gap-3">
            <h3 className="font-semibold text-stone-100">{t('resolvedEvents')}</h3>
            {resolvedEvents.map((event) => (
              <div className="rounded border border-stone-800 bg-dusk-950 p-3" key={event.id}>
                <div className="flex flex-wrap items-start justify-between gap-3">
                  <div>
                    <p className="text-xs uppercase tracking-wide text-stone-500">{t(`categories.${event.category}`)}</p>
                    <h3 className="mt-1 font-semibold text-stone-100">{getLocalizedEventTitle(t, event)}</h3>
                  </div>
                  <div className="text-right text-sm text-stone-400">
                    <div>{t(`status.${event.status}`)}</div>
                    {event.resolvedAt ? <div>{formatDate(event.resolvedAt)}</div> : null}
                  </div>
                </div>
                {event.selectedChoiceKey ? <p className="mt-2 text-sm text-stone-500">{t('events:selectedChoice', { choice: getLocalizedEventSelectedChoiceLabel(t, event) })}</p> : null}
                {event.result ? (
                  <div className="mt-2">
                    <h4 className="font-semibold text-stone-100">{getLocalizedEventResultTitle(t, event)}</h4>
                    <p className="mt-1 text-sm text-stone-400">{getLocalizedEventResultBody(t, event)}</p>
                  </div>
                ) : null}
              </div>
            ))}
          </div>
        ) : null}
      </div>
    </Card>
  );
}

type ReportsPanelProps = {
  error: string;
  formatters: Pick<DashboardFormatters, 'formatDate' | 'resourceList' | 'unitLabel'>;
  loading: boolean;
  loadingReportID: string | null;
  markingReportID: string | null;
  reports: MissionReport[];
  selectedReportID: string | null;
  unreadReportsCount: number;
  onMarkReportRead: (reportID: string) => void;
  onRefresh: () => void;
  onToggleReportDetails: (reportID: string) => void;
};

export function ReportsPanel({
  error,
  formatters,
  loading,
  loadingReportID,
  markingReportID,
  onMarkReportRead,
  onRefresh,
  onToggleReportDetails,
  reports,
  selectedReportID,
  unreadReportsCount,
}: ReportsPanelProps) {
  const { t } = useTranslation('reports');

  return (
    <Card title={t('section.title')}>
      <div className="grid gap-3">
        <div className="flex flex-wrap items-center justify-between gap-3">
          <p className="text-sm text-stone-400">{t('unread', { count: unreadReportsCount })}</p>
          <Button className="justify-self-start" disabled={loading} onClick={onRefresh} type="button">
            {t('refresh')}
          </Button>
        </div>
        {loading ? <p>{t('loading')}</p> : null}
        {error ? <p className="text-red-300">{error}</p> : null}
        {!loading && reports.length === 0 ? <p className="text-sm text-stone-400">{t('noReports')}</p> : null}
        {!loading
          ? reports.map((report) => {
              const isExpanded = selectedReportID === report.id;
              return (
                <div
                  className={report.isRead ? 'rounded border border-stone-800 bg-dusk-950 p-3' : 'rounded border border-dusk-gold bg-dusk-950 p-3'}
                  key={report.id}
                >
                  <div className="flex flex-wrap justify-between gap-3">
                    <div>
                      <div className="flex flex-wrap items-center gap-2">
                        <h3 className="font-semibold text-stone-100">{getLocalizedReportTitle(t, report)}</h3>
                        <span className="text-xs text-stone-500">{t(`types.${report.type}`)}</span>
                        <span className={report.isRead ? 'text-xs text-stone-500' : 'text-xs text-dusk-gold'}>
                          {report.isRead ? t('read') : t('new')}
                        </span>
                      </div>
                      <p className="text-sm text-dusk-gold">{t(`results.${report.result}`)}</p>
                    </div>
                    <div className="text-right text-sm text-stone-400">{formatters.formatDate(report.createdAt)}</div>
                  </div>
                  <p className="mt-2 text-sm text-stone-400">{report.body}</p>
                  <div className="mt-3 flex flex-wrap gap-2">
                    <Button onClick={() => onToggleReportDetails(report.id)} type="button">
                      {isExpanded ? t('close') : t('open')}
                    </Button>
                    <Button disabled={markingReportID === report.id || report.isRead} onClick={() => onMarkReportRead(report.id)} type="button">
                      {markingReportID === report.id ? t('markingRead') : t('markRead')}
                    </Button>
                  </div>
                  {loadingReportID === report.id ? <p className="mt-3 text-sm text-stone-400">{t('loadingDetails')}</p> : null}
                  {isExpanded && loadingReportID !== report.id ? (
                    <div className="mt-4 grid gap-3">
                      <div className="grid gap-2">
                        {report.phases.length === 0 ? <p className="text-sm text-stone-400">{t('noPhases')}</p> : null}
                        {report.phases.map((phase) => (
                          <div className="rounded border border-stone-800 bg-dusk-900 p-3" key={`${report.id}-${phase.title}`}>
                            <h4 className="font-semibold text-stone-100">{getLocalizedReportPhaseTitle(t, phase.title)}</h4>
                            <p className="mt-1 text-sm text-stone-400">{phase.body}</p>
                          </div>
                        ))}
                      </div>
                      <p className="text-sm text-stone-400">
                        {t('rewards')}: {formatters.resourceList(report.rewards)}
                      </p>
                      <p className="text-sm text-stone-400">
                        {t('unitsLost')}: {unitTypes.map((unitType) => `${formatters.unitLabel(unitType)}: ${report.losses[unitType] ?? 0}`).join(', ')}
                      </p>
                    </div>
                  ) : null}
                </div>
              );
            })
          : null}
      </div>
    </Card>
  );
}
