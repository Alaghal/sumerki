import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

import {
  Army,
  AvailableMission,
  breakPatron,
  Building,
  BuildingType,
  chooseEvent,
  choosePatronCrisis,
  getAvailableMissions,
  getMyArmy,
  getMyBuildings,
  getMyEvents,
  getMyMissions,
  getMyPatron,
  getMyRaids,
  getMyReports,
  getMyResources,
  getMyRuler,
  getNeighbors,
  getPatronOptions,
  getPatronPressure,
  getReport,
  joinPatron,
  KingdomEvent,
  markReportRead,
  Mission,
  MissionReport,
  Neighbor,
  PatronKey,
  PatronOption,
  PatronPressure,
  PatronStatus,
  Raid,
  Resources,
  ResourceValues,
  Ruler,
  payPatronTribute,
  startMission,
  startRaid,
  trainUnits,
  UnitType,
  upgradeBuilding,
} from '../api/client';
import { toUserMessage } from '../api/errors';
import { AppShell } from '../components/layout/AppShell';
import { Button } from '../components/ui/Button';
import { Card } from '../components/ui/Card';
import { useSession } from '../context/SessionContext';

const rulerStats = [
  'authority',
  'courage',
  'cunning',
  'honor',
  'cruelty',
  'ambition',
  'paranoia',
] as const;

const resourceRows = ['gold', 'food', 'wood', 'stone', 'population'] as const;

const costRows = ['gold', 'wood', 'stone'] as const;

const armyCostRows = ['gold', 'food', 'wood', 'stone', 'population'] as const;

const unitStatRows = ['attack', 'defense', 'speed', 'supply'] as const;

const unitTypes: UnitType[] = ['militia', 'spearmen', 'archers', 'cavalry', 'scouts'];

type ResourceKey = keyof ResourceValues;

export function DashboardPage() {
  const { i18n, t } = useTranslation([
    'game',
    'common',
    'kingdom',
    'resources',
    'buildings',
    'units',
    'missions',
    'reports',
    'patrons',
    'events',
    'raids',
  ]);
  const { token, user, kingdom } = useSession();
  const [ruler, setRuler] = useState<Ruler | null>(null);
  const [rulerLoading, setRulerLoading] = useState(true);
  const [rulerError, setRulerError] = useState('');
  const [resources, setResources] = useState<Resources | null>(null);
  const [resourcesLoading, setResourcesLoading] = useState(true);
  const [resourcesError, setResourcesError] = useState('');
  const [buildings, setBuildings] = useState<Building[]>([]);
  const [buildingsLoading, setBuildingsLoading] = useState(true);
  const [buildingsError, setBuildingsError] = useState('');
  const [upgradingType, setUpgradingType] = useState<BuildingType | null>(null);
  const [army, setArmy] = useState<Army | null>(null);
  const [armyLoading, setArmyLoading] = useState(true);
  const [armyError, setArmyError] = useState('');
  const [trainingType, setTrainingType] = useState<UnitType>('militia');
  const [trainingAmount, setTrainingAmount] = useState(5);
  const [isTraining, setIsTraining] = useState(false);
  const [availableMissions, setAvailableMissions] = useState<AvailableMission[]>([]);
  const [missions, setMissions] = useState<Mission[]>([]);
  const [missionsLoading, setMissionsLoading] = useState(true);
  const [missionsError, setMissionsError] = useState('');
  const [missionInputs, setMissionInputs] = useState<Record<string, Partial<Record<UnitType, number>>>>({});
  const [startingMissionKey, setStartingMissionKey] = useState<string | null>(null);
  const [reports, setReports] = useState<MissionReport[]>([]);
  const [unreadReportsCount, setUnreadReportsCount] = useState(0);
  const [selectedReportID, setSelectedReportID] = useState<string | null>(null);
  const [loadingReportID, setLoadingReportID] = useState<string | null>(null);
  const [markingReportID, setMarkingReportID] = useState<string | null>(null);
  const [reportsLoading, setReportsLoading] = useState(true);
  const [reportsError, setReportsError] = useState('');
  const [patronOptions, setPatronOptions] = useState<PatronOption[]>([]);
  const [patronStatus, setPatronStatus] = useState<PatronStatus | null>(null);
  const [patronPressure, setPatronPressure] = useState<PatronPressure | null>(null);
  const [patronLoading, setPatronLoading] = useState(true);
  const [patronError, setPatronError] = useState('');
  const [joiningPatron, setJoiningPatron] = useState<PatronKey | null>(null);
  const [isBreakingPatron, setIsBreakingPatron] = useState(false);
  const [isPayingTribute, setIsPayingTribute] = useState(false);
  const [crisisChoice, setCrisisChoice] = useState<'ask_delay' | 'break_patron' | null>(null);
  const [neighbors, setNeighbors] = useState<Neighbor[]>([]);
  const [raids, setRaids] = useState<Raid[]>([]);
  const [raidsLoading, setRaidsLoading] = useState(true);
  const [raidsError, setRaidsError] = useState('');
  const [selectedRaidTargetID, setSelectedRaidTargetID] = useState<string | null>(null);
  const [raidInputs, setRaidInputs] = useState<Partial<Record<UnitType, number>>>({});
  const [isStartingRaid, setIsStartingRaid] = useState(false);
  const [events, setEvents] = useState<KingdomEvent[]>([]);
  const [eventsLoading, setEventsLoading] = useState(true);
  const [eventsError, setEventsError] = useState('');
  const [choosingEventID, setChoosingEventID] = useState<string | null>(null);

  function formatDate(value: string) {
    return new Date(value).toLocaleString(i18n.language === 'en' ? 'en-US' : 'ru-RU');
  }

  function resourceLabel(key: ResourceKey) {
    return t(`resources:${key}.name`);
  }

  function unitLabel(unitType: UnitType) {
    return t(`units:${unitType}.name`);
  }

  function resourceList(values: Partial<ResourceValues>) {
    return resourceRows.map((key) => `${resourceLabel(key)}: ${values[key] ?? 0}`).join(', ');
  }

  function costList(values: Partial<ResourceValues>, keys: readonly ResourceKey[] = costRows) {
    return keys.map((key) => `${resourceLabel(key)}: ${values[key] ?? 0}`).join(', ');
  }

  function unitList<TUnit extends { unitType: UnitType; amountSent: number; amountLost: number; amountReturned: number }>(units: TUnit[]) {
    return units
      .map((unit) =>
        t('missions:unitSummary', {
          lost: unit.amountLost,
          returned: unit.amountReturned,
          sent: unit.amountSent,
          unit: unitLabel(unit.unitType),
        }),
      )
      .join(', ');
  }

  function eventChoiceLabel(event: KingdomEvent) {
    return event.choices.find((choice) => choice.key === event.selectedChoiceKey)?.label ?? t('common:states.unknown');
  }

  async function loadRuler() {
    if (!token || !kingdom) {
      return;
    }

    setRulerLoading(true);
    setRulerError('');

    try {
      const response = await getMyRuler(token);
      setRuler(response.ruler);
    } catch {
      setRulerError(t('game:ruler.error'));
    } finally {
      setRulerLoading(false);
    }
  }

  async function loadResources() {
    if (!token || !kingdom) {
      return;
    }

    setResourcesLoading(true);
    setResourcesError('');

    try {
      const response = await getMyResources(token);
      setResources(response.resources);
    } catch {
      setResourcesError(t('resources:loadError'));
    } finally {
      setResourcesLoading(false);
    }
  }

  async function loadBuildings() {
    if (!token || !kingdom) {
      return;
    }

    setBuildingsLoading(true);
    setBuildingsError('');

    try {
      const response = await getMyBuildings(token);
      setBuildings(response.buildings);
    } catch {
      setBuildingsError(t('buildings:loadError'));
    } finally {
      setBuildingsLoading(false);
    }
  }

  async function loadArmy() {
    if (!token || !kingdom) {
      return;
    }

    setArmyLoading(true);
    setArmyError('');

    try {
      const response = await getMyArmy(token);
      setArmy(response.army);
    } catch {
      setArmyError(t('units:loadError'));
    } finally {
      setArmyLoading(false);
    }
  }

  async function loadMissions() {
    if (!token || !kingdom) {
      return;
    }

    setMissionsLoading(true);
    setMissionsError('');

    try {
      const [availableResponse, currentResponse] = await Promise.all([getAvailableMissions(token), getMyMissions(token)]);
      setAvailableMissions(availableResponse.missions);
      setMissions(currentResponse.missions);
    } catch {
      setMissionsError(t('missions:loadError'));
    } finally {
      setMissionsLoading(false);
    }
  }

  async function loadReports() {
    if (!token || !kingdom) {
      return;
    }

    setReportsLoading(true);
    setReportsError('');

    try {
      const response = await getMyReports(token);
      setReports(response.reports);
      setUnreadReportsCount(response.unreadCount);
    } catch {
      setReportsError(t('reports:loadError'));
    } finally {
      setReportsLoading(false);
    }
  }

  async function loadPatron() {
    if (!token || !kingdom) {
      return;
    }

    setPatronLoading(true);
    setPatronError('');

    try {
      const [optionsResponse, statusResponse, pressureResponse] = await Promise.all([
        getPatronOptions(token),
        getMyPatron(token),
        getPatronPressure(token),
      ]);
      setPatronOptions(optionsResponse.patrons);
      setPatronStatus(statusResponse);
      setPatronPressure(pressureResponse.pressure);
    } catch {
      setPatronError(t('patrons:loadError'));
    } finally {
      setPatronLoading(false);
    }
  }

  async function loadRaids() {
    if (!token || !kingdom) {
      return;
    }

    setRaidsLoading(true);
    setRaidsError('');

    try {
      const [neighborsResponse, raidsResponse] = await Promise.all([getNeighbors(token), getMyRaids(token)]);
      setNeighbors(neighborsResponse.neighbors);
      setRaids(raidsResponse.raids);
      setSelectedRaidTargetID((current) => current ?? neighborsResponse.neighbors.find((neighbor) => neighbor.canRaid)?.kingdomId ?? null);
    } catch {
      setRaidsError(t('raids:loadError'));
    } finally {
      setRaidsLoading(false);
    }
  }

  async function loadEvents() {
    if (!token || !kingdom) {
      return;
    }

    setEventsLoading(true);
    setEventsError('');

    try {
      const response = await getMyEvents(token);
      setEvents(response.events);
    } catch {
      setEventsError(t('events:loadError'));
    } finally {
      setEventsLoading(false);
    }
  }

  async function refreshCity() {
    await Promise.all([loadRuler(), loadResources(), loadBuildings(), loadArmy(), loadMissions(), loadReports(), loadPatron(), loadRaids(), loadEvents()]);
  }

  async function handleUpgrade(buildingType: BuildingType) {
    if (!token) {
      return;
    }

    setUpgradingType(buildingType);
    setBuildingsError('');

    try {
      const response = await upgradeBuilding(buildingType, token);
      setResources(response.resources);
      await loadBuildings();
      await loadResources();
    } catch (caughtError) {
      setBuildingsError(toUserMessage(caughtError));
    } finally {
      setUpgradingType(null);
    }
  }

  async function handleTrain() {
    if (!token) {
      return;
    }

    setIsTraining(true);
    setArmyError('');

    try {
      const response = await trainUnits(trainingType, trainingAmount, token);
      setResources(response.resources);
      await loadArmy();
      await loadResources();
    } catch (caughtError) {
      setArmyError(toUserMessage(caughtError));
    } finally {
      setIsTraining(false);
    }
  }

  async function refreshMissions() {
    await Promise.all([loadMissions(), loadArmy(), loadResources(), loadReports()]);
  }

  async function refreshRaids() {
    await Promise.all([loadRaids(), loadArmy(), loadResources(), loadReports()]);
  }

  async function handleJoinPatron(patron: PatronKey) {
    if (!token) {
      return;
    }

    setJoiningPatron(patron);
    setPatronError('');

    try {
      await joinPatron(patron, token);
      await Promise.all([loadPatron(), loadResources()]);
    } catch (caughtError) {
      setPatronError(toUserMessage(caughtError));
    } finally {
      setJoiningPatron(null);
    }
  }

  async function handleBreakPatron() {
    if (!token) {
      return;
    }

    setIsBreakingPatron(true);
    setPatronError('');

    try {
      await breakPatron(token);
      await loadPatron();
    } catch (caughtError) {
      setPatronError(toUserMessage(caughtError));
    } finally {
      setIsBreakingPatron(false);
    }
  }

  async function handlePayTribute() {
    if (!token) {
      return;
    }

    setIsPayingTribute(true);
    setPatronError('');

    try {
      const response = await payPatronTribute(token);
      setPatronPressure(response.pressure);
      setResources(response.resources);
    } catch (caughtError) {
      setPatronError(toUserMessage(caughtError));
    } finally {
      setIsPayingTribute(false);
    }
  }

  async function handleCrisisChoice(choice: 'ask_delay' | 'break_patron') {
    if (!token) {
      return;
    }

    setCrisisChoice(choice);
    setPatronError('');

    try {
      const response = await choosePatronCrisis(choice, token);
      setPatronPressure(response.pressure);
      if (response.kingdom) {
        setPatronStatus((current) => (current ? { ...current, patron: null } : current));
      }
      await Promise.all([loadPatron(), loadResources()]);
    } catch (caughtError) {
      setPatronError(toUserMessage(caughtError));
    } finally {
      setCrisisChoice(null);
    }
  }

  function setRaidUnitAmount(unitType: UnitType, amount: number) {
    setRaidInputs((current) => ({
      ...current,
      [unitType]: amount,
    }));
  }

  async function handleStartRaid() {
    if (!token || !selectedRaidTargetID) {
      return;
    }

    const units = unitTypes
      .map((unitType) => ({
        unitType,
        amount: raidInputs[unitType] ?? 0,
      }))
      .filter((unit) => unit.amount > 0);

    if (units.length === 0 || units.some((unit) => !Number.isInteger(unit.amount) || unit.amount < 0)) {
      setRaidsError(t('raids:validation.choosePositiveUnits'));
      return;
    }

    setIsStartingRaid(true);
    setRaidsError('');

    try {
      const response = await startRaid(selectedRaidTargetID, units, token);
      setArmy(response.army);
      await refreshRaids();
    } catch (caughtError) {
      setRaidsError(toUserMessage(caughtError));
    } finally {
      setIsStartingRaid(false);
    }
  }

  async function handleChooseEvent(eventID: string, choiceKey: string) {
    if (!token) {
      return;
    }

    setChoosingEventID(eventID);
    setEventsError('');

    try {
      const response = await chooseEvent(eventID, choiceKey, token);
      setEvents((current) => current.map((event) => (event.id === eventID ? response.event : event)));
      if (response.resources) {
        setResources(response.resources);
      }
      if (response.army) {
        setArmy(response.army);
      }
      await Promise.all([loadEvents(), loadResources(), loadArmy(), loadPatron(), loadReports()]);
    } catch (caughtError) {
      setEventsError(toUserMessage(caughtError));
    } finally {
      setChoosingEventID(null);
    }
  }

  async function toggleReportDetails(reportID: string) {
    if (!token) {
      return;
    }
    if (selectedReportID === reportID) {
      setSelectedReportID(null);
      return;
    }

    setSelectedReportID(reportID);
    setLoadingReportID(reportID);
    setReportsError('');

    try {
      const response = await getReport(reportID, token);
      setReports((current) => current.map((report) => (report.id === reportID ? response.report : report)));
    } catch (caughtError) {
      setReportsError(toUserMessage(caughtError));
    } finally {
      setLoadingReportID(null);
    }
  }

  async function handleMarkReportRead(reportID: string) {
    if (!token) {
      return;
    }

    setMarkingReportID(reportID);
    setReportsError('');

    try {
      const wasUnread = reports.some((report) => report.id === reportID && !report.isRead);
      const response = await markReportRead(reportID, token);
      setReports((current) => current.map((report) => (report.id === reportID ? response.report : report)));
      if (wasUnread) {
        setUnreadReportsCount((current) => Math.max(0, current - 1));
      }
    } catch (caughtError) {
      setReportsError(toUserMessage(caughtError));
    } finally {
      setMarkingReportID(null);
    }
  }

  function setMissionUnitAmount(missionKey: string, unitType: UnitType, amount: number) {
    setMissionInputs((current) => ({
      ...current,
      [missionKey]: {
        ...current[missionKey],
        [unitType]: amount,
      },
    }));
  }

  async function handleStartMission(missionKey: string) {
    if (!token) {
      return;
    }

    const units = unitTypes
      .map((unitType) => ({
        unitType,
        amount: missionInputs[missionKey]?.[unitType] ?? 0,
      }))
      .filter((unit) => unit.amount > 0);

    if (units.length === 0 || units.some((unit) => !Number.isInteger(unit.amount) || unit.amount < 0)) {
      setMissionsError(t('missions:validation.choosePositiveUnits'));
      return;
    }

    setStartingMissionKey(missionKey);
    setMissionsError('');

    try {
      const response = await startMission(missionKey, units, token);
      setArmy(response.army);
      await refreshMissions();
    } catch (caughtError) {
      setMissionsError(toUserMessage(caughtError));
    } finally {
      setStartingMissionKey(null);
    }
  }

  useEffect(() => {
    let isActive = true;

    async function loadRuler() {
      if (!token || !kingdom) {
        return;
      }

      setRulerLoading(true);
      setRulerError('');

      try {
        const response = await getMyRuler(token);
        if (isActive) {
          setRuler(response.ruler);
        }
      } catch {
        if (isActive) {
          setRulerError(t('game:ruler.error'));
        }
      } finally {
        if (isActive) {
          setRulerLoading(false);
        }
      }
    }

    loadRuler();

    return () => {
      isActive = false;
    };
  }, [kingdom, token]);

  useEffect(() => {
    let isActive = true;

    async function loadInitialRaids() {
      if (!token || !kingdom) {
        return;
      }

      setRaidsLoading(true);
      setRaidsError('');

      try {
        const [neighborsResponse, raidsResponse] = await Promise.all([getNeighbors(token), getMyRaids(token)]);
        if (isActive) {
          setNeighbors(neighborsResponse.neighbors);
          setRaids(raidsResponse.raids);
          setSelectedRaidTargetID(neighborsResponse.neighbors.find((neighbor) => neighbor.canRaid)?.kingdomId ?? null);
        }
      } catch {
        if (isActive) {
          setRaidsError(t('raids:loadError'));
        }
      } finally {
        if (isActive) {
          setRaidsLoading(false);
        }
      }
    }

    loadInitialRaids();

    return () => {
      isActive = false;
    };
  }, [kingdom, token]);

  useEffect(() => {
    let isActive = true;

    async function loadInitialPatron() {
      if (!token || !kingdom) {
        return;
      }

      setPatronLoading(true);
      setPatronError('');

      try {
        const [optionsResponse, statusResponse, pressureResponse] = await Promise.all([
          getPatronOptions(token),
          getMyPatron(token),
          getPatronPressure(token),
        ]);
        if (isActive) {
          setPatronOptions(optionsResponse.patrons);
          setPatronStatus(statusResponse);
          setPatronPressure(pressureResponse.pressure);
        }
      } catch {
        if (isActive) {
          setPatronError(t('patrons:loadError'));
        }
      } finally {
        if (isActive) {
          setPatronLoading(false);
        }
      }
    }

    loadInitialPatron();

    return () => {
      isActive = false;
    };
  }, [kingdom, token]);

  useEffect(() => {
    let isActive = true;

    async function loadInitialResources() {
      if (!token || !kingdom) {
        return;
      }

      setResourcesLoading(true);
      setResourcesError('');

      try {
        const response = await getMyResources(token);
        if (isActive) {
          setResources(response.resources);
        }
      } catch {
        if (isActive) {
          setResourcesError(t('resources:loadError'));
        }
      } finally {
        if (isActive) {
          setResourcesLoading(false);
        }
      }
    }

    loadInitialResources();

    return () => {
      isActive = false;
    };
  }, [kingdom, token]);

  useEffect(() => {
    let isActive = true;

    async function loadInitialBuildings() {
      if (!token || !kingdom) {
        return;
      }

      setBuildingsLoading(true);
      setBuildingsError('');

      try {
        const response = await getMyBuildings(token);
        if (isActive) {
          setBuildings(response.buildings);
        }
      } catch {
        if (isActive) {
          setBuildingsError(t('buildings:loadError'));
        }
      } finally {
        if (isActive) {
          setBuildingsLoading(false);
        }
      }
    }

    loadInitialBuildings();

    return () => {
      isActive = false;
    };
  }, [kingdom, token]);

  useEffect(() => {
    let isActive = true;

    async function loadInitialArmy() {
      if (!token || !kingdom) {
        return;
      }

      setArmyLoading(true);
      setArmyError('');

      try {
        const response = await getMyArmy(token);
        if (isActive) {
          setArmy(response.army);
        }
      } catch {
        if (isActive) {
          setArmyError(t('units:loadError'));
        }
      } finally {
        if (isActive) {
          setArmyLoading(false);
        }
      }
    }

    loadInitialArmy();

    return () => {
      isActive = false;
    };
  }, [kingdom, token]);

  useEffect(() => {
    let isActive = true;

    async function loadInitialMissions() {
      if (!token || !kingdom) {
        return;
      }

      setMissionsLoading(true);
      setMissionsError('');

      try {
        const [availableResponse, currentResponse] = await Promise.all([getAvailableMissions(token), getMyMissions(token)]);
        if (isActive) {
          setAvailableMissions(availableResponse.missions);
          setMissions(currentResponse.missions);
        }
      } catch {
        if (isActive) {
          setMissionsError(t('missions:loadError'));
        }
      } finally {
        if (isActive) {
          setMissionsLoading(false);
        }
      }
    }

    loadInitialMissions();

    return () => {
      isActive = false;
    };
  }, [kingdom, token]);

  useEffect(() => {
    let isActive = true;

    async function loadInitialReports() {
      if (!token || !kingdom) {
        return;
      }

      setReportsLoading(true);
      setReportsError('');

      try {
        const response = await getMyReports(token);
        if (isActive) {
          setReports(response.reports);
          setUnreadReportsCount(response.unreadCount);
        }
      } catch {
        if (isActive) {
          setReportsError(t('reports:loadError'));
        }
      } finally {
        if (isActive) {
          setReportsLoading(false);
        }
      }
    }

    loadInitialReports();

    return () => {
      isActive = false;
    };
  }, [kingdom, token]);

  useEffect(() => {
    let isActive = true;

    async function loadInitialEvents() {
      if (!token || !kingdom) {
        return;
      }

      setEventsLoading(true);
      setEventsError('');

      try {
        const response = await getMyEvents(token);
        if (isActive) {
          setEvents(response.events);
        }
      } catch {
        if (isActive) {
          setEventsError(t('events:loadError'));
        }
      } finally {
        if (isActive) {
          setEventsLoading(false);
        }
      }
    }

    loadInitialEvents();

    return () => {
      isActive = false;
    };
  }, [kingdom, token]);

  if (!user || !kingdom) {
    return null;
  }

  return (
    <AppShell>
      <div className="grid gap-4">
        <div>
          <h1 className="text-2xl font-semibold text-stone-100">{t('game:dashboard.title', { name: kingdom.name })}</h1>
          <p className="mt-1 text-sm text-stone-400">{t('game:dashboard.subtitle', { email: user.email })}</p>
        </div>
        <Button className="justify-self-start" onClick={refreshCity} type="button">
          {t('game:dashboard.refreshCity')}
        </Button>
        <div className="grid gap-4 lg:grid-cols-2">
          <Card title={t('game:dashboard.kingdom')}>
            <dl className="grid gap-2">
              <div className="flex justify-between gap-4">
                <dt className="text-stone-400">{t('game:dashboard.culture')}</dt>
                <dd className="text-right text-stone-100">{t(`kingdom:cultures.${kingdom.culture}.name`)}</dd>
              </div>
              <div className="flex justify-between gap-4">
                <dt className="text-stone-400">{t('game:dashboard.patron')}</dt>
                <dd className="text-right text-stone-100">
                  {patronStatus?.patron
                    ? t(`patrons:${patronStatus.patron.key}.name`)
                    : kingdom.patron
                      ? t(`patrons:${kingdom.patron}.name`)
                      : t('game:dashboard.noPatron')}
                </dd>
              </div>
              <div className="flex justify-between gap-4">
                <dt className="text-stone-400">{t('game:dashboard.player')}</dt>
                <dd className="max-w-[14rem] truncate text-right text-stone-100">{user.email}</dd>
              </div>
            </dl>
          </Card>
          <Card title={t('patrons:section.title')}>
            <div className="grid gap-4">
              {patronLoading ? <p>{t('patrons:loading')}</p> : null}
              {patronError ? <p className="text-red-300">{patronError}</p> : null}
              {!patronLoading ? (
                <>
                  <div className="rounded border border-stone-800 bg-dusk-950 p-3">
                    <div className="flex flex-wrap items-start justify-between gap-3">
                      <div>
                        <h3 className="font-semibold text-stone-100">
                          {patronStatus?.patron ? t(`patrons:${patronStatus.patron.key}.name`) : t('patrons:notChosen')}
                        </h3>
                        <p className="mt-1 text-sm text-stone-400">
                          {patronPressure?.summary ?? t('patrons:choosePath')}
                        </p>
                      </div>
                      {patronStatus?.patron ? (
                        <div className="text-right text-sm text-stone-300">
                          <div>{t('patrons:favor')}: {patronStatus.patron.favor}</div>
                          <div>{t(`patrons:standing.${patronStatus.patron.standing}`)}</div>
                          <div>{formatDate(patronStatus.patron.joinedAt)}</div>
                        </div>
                      ) : null}
                    </div>
                    {patronStatus?.patron ? (
                      <Button className="mt-3 justify-self-start" disabled={isBreakingPatron} onClick={handleBreakPatron} type="button">
                        {isBreakingPatron ? t('patrons:breakingTie') : t('patrons:breakTie')}
                      </Button>
                    ) : null}
                  </div>

                  {patronPressure ? (
                    <div className="rounded border border-stone-800 bg-dusk-950 p-3">
                      <div className="flex flex-wrap items-start justify-between gap-3">
                        <div>
                          <h3 className="font-semibold text-stone-100">{t('patrons:pressureTitle')}</h3>
                          <p className="mt-1 text-sm text-stone-400">{patronPressure.summary}</p>
                        </div>
                        <div className="text-right text-sm text-stone-300">
                          <div>{t(`patrons:pressure.status.${patronPressure.crisisStatus}`)}</div>
                          <div>{t('patrons:pressureLevel', { level: patronPressure.pressureLevel })}</div>
                        </div>
                      </div>
                      <dl className="mt-3 grid gap-2 text-sm">
                        <div className="flex justify-between gap-4">
                          <dt className="text-stone-400">{t('patrons:tributeDebt')}</dt>
                          <dd className="text-right text-stone-100">
                            {resourceList(patronPressure.tributeDebt)}
                          </dd>
                        </div>
                        <div className="flex justify-between gap-4">
                          <dt className="text-stone-400">{t('patrons:contributionDebt')}</dt>
                          <dd className="text-right text-stone-100">{resourceList(patronPressure.contributionDebt)}</dd>
                        </div>
                        <div className="flex justify-between gap-4">
                          <dt className="text-stone-400">{t('patrons:nextTribute')}</dt>
                          <dd className="text-right text-stone-100">
                            {patronPressure.nextTributeAt ? formatDate(patronPressure.nextTributeAt) : t('patrons:noNextTribute')}
                          </dd>
                        </div>
                        {patronPressure.delayUntil ? (
                          <div className="flex justify-between gap-4">
                            <dt className="text-stone-400">{t('patrons:delayUntil')}</dt>
                            <dd className="text-right text-stone-100">{formatDate(patronPressure.delayUntil)}</dd>
                          </div>
                        ) : null}
                        <div className="flex justify-between gap-4">
                          <dt className="text-stone-400">{t('patrons:protectedMinimums')}</dt>
                          <dd className="text-right text-stone-100">
                            {patronPressure.protectedMinimums.gold ?? 0} / {patronPressure.protectedMinimums.food ?? 0} /{' '}
                            {patronPressure.protectedMinimums.wood ?? 0} / {patronPressure.protectedMinimums.stone ?? 0}
                          </dd>
                        </div>
                      </dl>
                      <div className="mt-3 flex flex-wrap gap-2">
                        {patronPressure.availableActions.includes('pay_tribute') ? (
                          <Button disabled={isPayingTribute} onClick={handlePayTribute} type="button">
                            {isPayingTribute ? t('patrons:payingTribute') : t('patrons:payTribute')}
                          </Button>
                        ) : null}
                        {patronPressure.availableActions.includes('ask_delay') ? (
                          <Button disabled={crisisChoice === 'ask_delay'} onClick={() => handleCrisisChoice('ask_delay')} type="button">
                            {crisisChoice === 'ask_delay' ? t('patrons:askingDelay') : t('patrons:askDelay')}
                          </Button>
                        ) : null}
                        {patronPressure.availableActions.includes('break_patron') ? (
                          <Button disabled={crisisChoice === 'break_patron'} onClick={() => handleCrisisChoice('break_patron')} type="button">
                            {crisisChoice === 'break_patron' ? t('patrons:breakingTie') : t('patrons:breakDuringCrisis')}
                          </Button>
                        ) : null}
                      </div>
                    </div>
                  ) : null}

                  <div className="flex flex-wrap gap-2">
                    <Button disabled={patronLoading} onClick={loadPatron} type="button">
                      {t('patrons:refresh')}
                    </Button>
                  </div>

                  <div className="grid gap-3">
                    {patronOptions.map((option) => {
                      const isCurrent = patronStatus?.patron?.key === option.key;
                      return (
                        <div className="rounded border border-stone-800 bg-dusk-950 p-3" key={option.key}>
                          <div className="flex flex-wrap items-start justify-between gap-3">
                            <div>
                              <h3 className="font-semibold text-stone-100">{t(`patrons:${option.key}.name`)}</h3>
                              <p className="mt-1 text-sm text-stone-400">{option.shortDescription}</p>
                              <p className="mt-2 text-sm text-stone-500">{option.flavor}</p>
                            </div>
                            <Button
                              disabled={joiningPatron === option.key || isCurrent}
                              onClick={() => handleJoinPatron(option.key)}
                              type="button"
                            >
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
          <Card title={t('resources:section.title')}>
            <div className="grid gap-4">
              {resourcesLoading ? <p>{t('resources:loading')}</p> : null}
              {resourcesError ? <p className="text-red-300">{resourcesError}</p> : null}
              {resources && !resourcesLoading && !resourcesError ? (
                <dl className="grid gap-2">
                  {resourceRows.map((key) => (
                    <div className="flex items-center justify-between gap-4" key={key}>
                      <dt className="text-stone-400">{resourceLabel(key)}</dt>
                      <dd className="text-right">
                        <div className="font-semibold text-stone-100">{resources[key]}</div>
                        <div className="text-xs text-dusk-gold">{t('resources:productionPerHour', { amount: resources.productionPerHour[key] })}</div>
                      </dd>
                    </div>
                  ))}
                </dl>
              ) : null}
              <Button className="justify-self-start" disabled={resourcesLoading} onClick={loadResources} type="button">
                {t('resources:refresh')}
              </Button>
            </div>
          </Card>
          <Card title={t('game:ruler.section.title')}>
            {rulerLoading ? <p>{t('game:ruler.loading')}</p> : null}
            {rulerError ? <p className="text-red-300">{rulerError}</p> : null}
            {ruler && !rulerLoading && !rulerError ? (
              <div className="grid gap-4">
                <dl className="grid gap-2">
                  <div className="flex justify-between gap-4">
                    <dt className="text-stone-400">{t('game:ruler.name')}</dt>
                    <dd className="text-right text-stone-100">{ruler.name}</dd>
                  </div>
                  <div className="flex justify-between gap-4">
                    <dt className="text-stone-400">{t('game:ruler.age')}</dt>
                    <dd className="text-right text-stone-100">{ruler.age}</dd>
                  </div>
                  <div className="flex justify-between gap-4">
                    <dt className="text-stone-400">{t('game:ruler.health')}</dt>
                    <dd className="text-right text-stone-100">{t(`game:ruler.healthStatus.${ruler.healthStatus}`)}</dd>
                  </div>
                  <div className="flex justify-between gap-4">
                    <dt className="text-stone-400">{t('game:dashboard.culture')}</dt>
                    <dd className="text-right text-stone-100">{t(`kingdom:cultures.${ruler.culture}.name`)}</dd>
                  </div>
                </dl>
                <dl className="grid gap-2">
                  {rulerStats.map((key) => (
                    <div className="flex justify-between gap-4" key={key}>
                      <dt className="text-stone-400">{t(`game:ruler.stats.${key}`)}</dt>
                      <dd className="text-right text-stone-100">{ruler[key]}</dd>
                    </div>
                  ))}
                </dl>
              </div>
            ) : null}
          </Card>
          <Card title={t('buildings:section.title')}>
            <div className="grid gap-4">
              {buildingsLoading ? <p>{t('buildings:loading')}</p> : null}
              {buildingsError ? <p className="text-red-300">{buildingsError}</p> : null}
              {!buildingsLoading && !buildingsError ? (
                <div className="grid gap-3">
                  {buildings.map((building) => (
                    <div
                      className="rounded border border-stone-800 bg-dusk-950 p-3"
                      data-building-type={building.type}
                      key={building.id}
                    >
                      <div className="flex flex-wrap items-start justify-between gap-3">
                        <div>
                          <h3 className="font-semibold text-stone-100">{t(`buildings:${building.type}.name`)}</h3>
                        </div>
                        <div className="text-right text-sm text-stone-300">
                          {t('buildings:level', { level: building.level, maxLevel: building.maxLevel })}
                        </div>
                      </div>
                      <div className="mt-3 grid gap-2">
                        {building.effects.map((effect) => (
                          <p className="text-sm text-stone-400" key={effect}>
                            {effect}
                          </p>
                        ))}
                        {building.isUpgrading ? (
                          <p className="text-dusk-gold">
                            {t('buildings:completedAt', {
                              date: building.upgradeFinishesAt ? formatDate(building.upgradeFinishesAt) : t('buildings:untilComplete'),
                            })}
                          </p>
                        ) : null}
                        {!building.isUpgrading && building.nextUpgrade ? (
                          <div className="grid gap-2">
                            {building.nextUpgrade.blockedReason === 'max_level' ? (
                              <p className="text-dusk-gold">{t('buildings:maxLevel')}</p>
                            ) : (
                              <>
                                <div className="grid gap-1">
                                  <p className="text-stone-400">
                                    {t('buildings:upgradeDetails', {
                                      level: building.nextUpgrade.targetLevel,
                                      seconds: building.nextUpgrade.durationSeconds,
                                    })}
                                  </p>
                                  <p className="text-stone-400">
                                    {t('buildings:cost')}: {costList(building.nextUpgrade.cost)}
                                  </p>
                                </div>
                                <Button
                                  className="justify-self-start"
                                  data-building-upgrade={building.type}
                                  disabled={upgradingType === building.type}
                                  onClick={() => handleUpgrade(building.type)}
                                  type="button"
                                >
                                  {upgradingType === building.type ? t('buildings:upgrading') : t('buildings:upgrade')}
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
          <Card title={t('units:section.title')}>
            <div className="grid gap-4">
              {armyLoading ? <p>{t('units:loading')}</p> : null}
              {armyError ? <p className="text-red-300">{armyError}</p> : null}
              {army && !armyLoading ? (
                <>
                  <dl className="grid gap-2">
                    <div className="flex justify-between gap-4">
                      <dt className="text-stone-400">{t('units:summary.total')}</dt>
                      <dd className="text-right text-stone-100">{army.summary.totalUnits}</dd>
                    </div>
                    <div className="flex justify-between gap-4">
                      <dt className="text-stone-400">{t('units:stats.attack')}</dt>
                      <dd className="text-right text-stone-100">{army.summary.totalAttack}</dd>
                    </div>
                    <div className="flex justify-between gap-4">
                      <dt className="text-stone-400">{t('units:stats.defense')}</dt>
                      <dd className="text-right text-stone-100">{army.summary.totalDefense}</dd>
                    </div>
                    <div className="flex justify-between gap-4">
                      <dt className="text-stone-400">{t('units:stats.supply')}</dt>
                      <dd className="text-right text-stone-100">{army.summary.totalSupply}</dd>
                    </div>
                  </dl>

                  <div className="grid gap-3">
                    {army.units.map((unit) => (
                      <div className="rounded border border-stone-800 bg-dusk-950 p-3" data-unit-type={unit.type} key={unit.type}>
                        <div className="flex flex-wrap items-start justify-between gap-3">
                          <div>
                            <h3 className="font-semibold text-stone-100">{unitLabel(unit.type)}</h3>
                          </div>
                          <div className="text-right text-sm text-stone-300">{unit.amount}</div>
                        </div>
                        <dl className="mt-3 grid gap-1 text-sm">
                          {unitStatRows.map((key) => (
                            <div className="flex justify-between gap-4" key={key}>
                              <dt className="text-stone-400">{t(`units:stats.${key}`)}</dt>
                              <dd className="text-right text-stone-100">{unit.stats[key]}</dd>
                            </div>
                          ))}
                        </dl>
                        <p className="mt-2 text-sm text-stone-400">
                          {t('units:cost')}: {costList(unit.cost, armyCostRows)}
                        </p>
                        <p className="mt-1 text-sm text-stone-400">{t('units:secondsPerUnit', { seconds: unit.secondsPerUnit })}</p>
                        <p className={unit.requirements.isMet ? 'mt-1 text-sm text-dusk-gold' : 'mt-1 text-sm text-red-300'}>
                          {unit.requirements.barracksLevel > 0
                            ? t('units:barracksRequirement', {
                                level: unit.requirements.barracksLevel,
                                status: unit.requirements.isMet ? t('units:requirementMet') : t('units:requirementNotMet'),
                              })
                            : t('units:requirementMet')}
                        </p>
                      </div>
                    ))}
                  </div>

                  <div className="grid gap-3 rounded border border-stone-800 bg-dusk-950 p-3">
                    <div className="grid gap-2 sm:grid-cols-[1fr_8rem_auto]">
                      <label className="grid gap-1 text-sm text-stone-400">
                        {t('units:unitType')}
                        <select
                          className="rounded border border-stone-700 bg-dusk-900 px-3 py-2 text-stone-100"
                          disabled={isTraining}
                          onChange={(event) => setTrainingType(event.target.value as UnitType)}
                          value={trainingType}
                        >
                          {army.units.map((unit) => (
                            <option key={unit.type} value={unit.type}>
                              {unitLabel(unit.type)}
                            </option>
                          ))}
                        </select>
                      </label>
                      <label className="grid gap-1 text-sm text-stone-400">
                        {t('units:amount')}
                        <input
                          className="rounded border border-stone-700 bg-dusk-900 px-3 py-2 text-stone-100"
                          disabled={isTraining}
                          max={50}
                          min={1}
                          onChange={(event) => setTrainingAmount(Number(event.target.value))}
                          type="number"
                          value={trainingAmount}
                        />
                      </label>
                      <Button className="self-end" disabled={isTraining} onClick={handleTrain} type="button">
                        {isTraining ? t('units:trainingNow') : t('units:train')}
                      </Button>
                    </div>
                    <p className="text-sm text-stone-400">{t('units:futureSystems')}</p>
                  </div>

                  <div className="grid gap-2">
                    <h3 className="font-semibold text-stone-100">{t('units:training')}</h3>
                    {army.trainingOrders.length === 0 ? <p className="text-sm text-stone-400">{t('units:noTraining')}</p> : null}
                    {army.trainingOrders.map((order) => (
                      <div className="rounded border border-stone-800 bg-dusk-950 p-3" key={order.id}>
                        <div className="flex flex-wrap justify-between gap-3">
                          <div>
                            <p className="font-semibold text-stone-100">{unitLabel(order.unitType)}</p>
                            <p className="text-sm text-stone-400">{t('units:amount')}: {order.amount}</p>
                          </div>
                          <div className="text-right text-sm text-dusk-gold">
                            {t('units:completedAt', { date: formatDate(order.finishesAt) })}
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                </>
              ) : null}
            </div>
          </Card>
          <Card title={t('missions:section.title')}>
            <div className="grid gap-4">
              {missionsLoading ? <p>{t('missions:loading')}</p> : null}
              {missionsError ? <p className="text-red-300">{missionsError}</p> : null}
              <Button className="justify-self-start" disabled={missionsLoading} onClick={refreshMissions} type="button">
                {t('missions:refresh')}
              </Button>
              {!missionsLoading ? (
                <>
                  <div className="grid gap-3">
                    <h3 className="font-semibold text-stone-100">{t('missions:availableMissions')}</h3>
                    {availableMissions.length === 0 ? <p className="text-sm text-stone-400">{t('missions:noAvailable')}</p> : null}
                    {availableMissions.map((mission) => (
                      <div className="rounded border border-stone-800 bg-dusk-950 p-3" data-mission-key={mission.key} key={mission.key}>
                        <div className="flex flex-wrap items-start justify-between gap-3">
                          <div>
                            <h3 className="font-semibold text-stone-100">
                              {t(`missions:${mission.key}.name`, { defaultValue: mission.label })}
                            </h3>
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
                          {t('missions:rewards')}: {resourceList(mission.baseRewards)}
                        </p>
                        <div className="mt-3 grid gap-2 sm:grid-cols-5">
                          {unitTypes.map((unitType) => (
                            <label className="grid gap-1 text-xs text-stone-400" key={unitType}>
                              {unitLabel(unitType)}
                              <input
                                className="rounded border border-stone-700 bg-dusk-900 px-2 py-2 text-stone-100"
                                min={0}
                                onChange={(event) => setMissionUnitAmount(mission.key, unitType, Number(event.target.value))}
                                type="number"
                                value={missionInputs[mission.key]?.[unitType] ?? 0}
                              />
                            </label>
                          ))}
                        </div>
                        <Button
                          className="mt-3 justify-self-start"
                          disabled={startingMissionKey === mission.key}
                          onClick={() => handleStartMission(mission.key)}
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
                              ? t('missions:completedAt', { date: formatDate(mission.completedAt) })
                              : t('missions:resolvesAt', { date: formatDate(mission.finishesAt) })}
                          </div>
                        </div>
                        <p className="mt-2 text-sm text-stone-400">
                          {t('missions:sentUnits')}: {unitList(mission.units)}
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
          <Card title={t('raids:section.title')}>
            <div className="grid gap-4">
              {raidsLoading ? <p>{t('raids:loading')}</p> : null}
              {raidsError ? <p className="text-red-300">{raidsError}</p> : null}
              <Button className="justify-self-start" disabled={raidsLoading} onClick={refreshRaids} type="button">
                {t('raids:refresh')}
              </Button>
              {!raidsLoading ? (
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
                            <div>{t('raids:dread')}: {neighbor.dread}</div>
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
                            <input
                              checked={selectedRaidTargetID === neighbor.kingdomId}
                              onChange={() => setSelectedRaidTargetID(neighbor.kingdomId)}
                              type="radio"
                            />
                            {t('raids:selectTarget')}
                          </label>
                        ) : null}
                      </div>
                    ))}
                  </div>

                  {selectedRaidTargetID ? (
                    <div className="grid gap-3 rounded border border-stone-800 bg-dusk-950 p-3">
                      <h3 className="font-semibold text-stone-100">{t('raids:sendParty')}</h3>
                      <div className="grid gap-2 sm:grid-cols-5">
                        {unitTypes.map((unitType) => (
                          <label className="grid gap-1 text-xs text-stone-400" key={unitType}>
                            {unitLabel(unitType)}
                            <input
                              className="rounded border border-stone-700 bg-dusk-900 px-2 py-2 text-stone-100"
                              min={0}
                              onChange={(event) => setRaidUnitAmount(unitType, Number(event.target.value))}
                              type="number"
                              value={raidInputs[unitType] ?? 0}
                            />
                          </label>
                        ))}
                      </div>
                      <p className="text-sm text-stone-400">{t('raids:minimumHint')}</p>
                      <Button className="justify-self-start" disabled={isStartingRaid} onClick={handleStartRaid} type="button">
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
                              ? t('raids:completedAt', { date: formatDate(raid.completedAt) })
                              : t('raids:arrivesAt', { date: formatDate(raid.arrivesAt) })}
                          </div>
                        </div>
                        <p className="mt-2 text-sm text-stone-400">
                          {t('raids:sentUnits')}: {unitList(raid.units)}
                        </p>
                        <p className="mt-1 text-sm text-stone-400">
                          {t('reports:loot')}: {resourceList(raid.loot)}
                        </p>
                      </div>
                    ))}
                  </div>
                </>
              ) : null}
            </div>
          </Card>
          <Card title={t('events:section.title')}>
            <div className="grid gap-3">
              <div className="flex flex-wrap items-center justify-between gap-3">
                <p className="text-sm text-stone-400">{t('events:subtitle')}</p>
                <Button className="justify-self-start" disabled={eventsLoading} onClick={loadEvents} type="button">
                  {t('events:refresh')}
                </Button>
              </div>
              {eventsLoading ? <p>{t('events:loading')}</p> : null}
              {eventsError ? <p className="text-red-300">{eventsError}</p> : null}
              {!eventsLoading && events.length === 0 ? (
                <p className="text-sm text-stone-400">{t('events:noEvents')}</p>
              ) : null}
              {!eventsLoading && events.filter((event) => event.status === 'active').length > 0 ? (
                <div className="grid gap-3">
                  {events
                    .filter((event) => event.status === 'active')
                    .map((event) => (
                      <div className="rounded border border-dusk-gold bg-dusk-950 p-3" key={event.id}>
                        <div className="flex flex-wrap items-start justify-between gap-3">
                          <div>
                            <p className="text-xs uppercase tracking-wide text-dusk-gold">{t(`events:categories.${event.category}`)}</p>
                            <h3 className="mt-1 font-semibold text-stone-100">{event.title}</h3>
                          </div>
                          <div className="text-right text-sm text-stone-400">
                            <div>{t(`events:status.${event.status}`)}</div>
                            <div>{t('events:expiresAt', { date: formatDate(event.expiresAt) })}</div>
                          </div>
                        </div>
                        <p className="mt-2 text-sm text-stone-400">{event.body}</p>
                        <div className="mt-3 grid gap-2">
                          {event.choices.map((choice) => (
                            <div className="rounded border border-stone-800 bg-dusk-900 p-3" key={`${event.id}-${choice.key}`}>
                              <div className="flex flex-wrap items-start justify-between gap-3">
                                <div>
                                  <h4 className="font-semibold text-stone-100">{choice.label}</h4>
                                  <p className="mt-1 text-sm text-stone-400">{choice.description}</p>
                                </div>
                                <Button
                                  disabled={choosingEventID === event.id}
                                  onClick={() => handleChooseEvent(event.id, choice.key)}
                                  type="button"
                                >
                                  {choosingEventID === event.id ? t('events:choosing') : t('events:choose')}
                                </Button>
                              </div>
                            </div>
                          ))}
                        </div>
                      </div>
                    ))}
                </div>
              ) : null}
              {!eventsLoading && events.filter((event) => event.status !== 'active').length > 0 ? (
                <div className="grid gap-3">
                  <h3 className="font-semibold text-stone-100">{t('events:resolvedEvents')}</h3>
                  {events
                    .filter((event) => event.status !== 'active')
                    .map((event) => (
                      <div className="rounded border border-stone-800 bg-dusk-950 p-3" key={event.id}>
                        <div className="flex flex-wrap items-start justify-between gap-3">
                          <div>
                            <p className="text-xs uppercase tracking-wide text-stone-500">{t(`events:categories.${event.category}`)}</p>
                            <h3 className="mt-1 font-semibold text-stone-100">{event.title}</h3>
                          </div>
                          <div className="text-right text-sm text-stone-400">
                            <div>{t(`events:status.${event.status}`)}</div>
                            {event.resolvedAt ? <div>{formatDate(event.resolvedAt)}</div> : null}
                          </div>
                        </div>
                        {event.selectedChoiceKey ? (
                          <p className="mt-2 text-sm text-stone-500">{t('events:selectedChoice', { choice: eventChoiceLabel(event) })}</p>
                        ) : null}
                        {event.result ? (
                          <div className="mt-2">
                            <h4 className="font-semibold text-stone-100">{event.result.title}</h4>
                            <p className="mt-1 text-sm text-stone-400">{event.result.body}</p>
                          </div>
                        ) : null}
                      </div>
                    ))}
                </div>
              ) : null}
            </div>
          </Card>
          <Card title={t('reports:section.title')}>
            <div className="grid gap-3">
              <div className="flex flex-wrap items-center justify-between gap-3">
                <p className="text-sm text-stone-400">{t('reports:unread', { count: unreadReportsCount })}</p>
                <Button className="justify-self-start" disabled={reportsLoading} onClick={loadReports} type="button">
                  {t('reports:refresh')}
                </Button>
              </div>
              {reportsLoading ? <p>{t('reports:loading')}</p> : null}
              {reportsError ? <p className="text-red-300">{reportsError}</p> : null}
              {!reportsLoading && reports.length === 0 ? (
                <p className="text-sm text-stone-400">{t('reports:noReports')}</p>
              ) : null}
              {!reportsLoading
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
                              <h3 className="font-semibold text-stone-100">{report.title}</h3>
                              <span className="text-xs text-stone-500">{t(`reports:types.${report.type}`)}</span>
                              <span className={report.isRead ? 'text-xs text-stone-500' : 'text-xs text-dusk-gold'}>
                                {report.isRead ? t('reports:read') : t('reports:new')}
                              </span>
                            </div>
                            <p className="text-sm text-dusk-gold">{t(`reports:results.${report.result}`)}</p>
                          </div>
                          <div className="text-right text-sm text-stone-400">{formatDate(report.createdAt)}</div>
                        </div>
                        <p className="mt-2 text-sm text-stone-400">{report.body}</p>
                        <div className="mt-3 flex flex-wrap gap-2">
                          <Button onClick={() => toggleReportDetails(report.id)} type="button">
                            {isExpanded ? t('reports:close') : t('reports:open')}
                          </Button>
                          <Button
                            disabled={markingReportID === report.id || report.isRead}
                            onClick={() => handleMarkReportRead(report.id)}
                            type="button"
                          >
                            {markingReportID === report.id ? t('reports:markingRead') : t('reports:markRead')}
                          </Button>
                        </div>
                        {loadingReportID === report.id ? <p className="mt-3 text-sm text-stone-400">{t('reports:loadingDetails')}</p> : null}
                        {isExpanded && loadingReportID !== report.id ? (
                          <div className="mt-4 grid gap-3">
                            <div className="grid gap-2">
                              {report.phases.length === 0 ? <p className="text-sm text-stone-400">{t('reports:noPhases')}</p> : null}
                              {report.phases.map((phase) => (
                                <div className="rounded border border-stone-800 bg-dusk-900 p-3" key={`${report.id}-${phase.title}`}>
                                  <h4 className="font-semibold text-stone-100">{phase.title}</h4>
                                  <p className="mt-1 text-sm text-stone-400">{phase.body}</p>
                                </div>
                              ))}
                            </div>
                            <p className="text-sm text-stone-400">
                              {t('reports:rewards')}: {resourceList(report.rewards)}
                            </p>
                            <p className="text-sm text-stone-400">
                              {t('reports:unitsLost')}:{' '}
                              {unitTypes
                                .map((unitType) => `${unitLabel(unitType)}: ${report.losses[unitType] ?? 0}`)
                                .join(', ')}
                            </p>
                          </div>
                        ) : null}
                      </div>
                    );
                  })
                : null}
            </div>
          </Card>
        </div>
      </div>
    </AppShell>
  );
}
