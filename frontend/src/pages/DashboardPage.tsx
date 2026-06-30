import { type ReactNode, useEffect, useState } from 'react';
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
import { useSession } from '../context/SessionContext';
import {
  ArmyPanel,
  BuildingsPanel,
  DashboardRefreshButton,
  EventsPanel,
  KingdomPanel,
  MissionsPanel,
  PatronPanel,
  RaidsPanel,
  ReportsPanel,
  ResourcesPanel,
  RulerPanel,
} from '../features/dashboard/DashboardPanels';
import { costRows, ResourceKey, resourceRows, unitTypes } from '../features/dashboard/shared';
import { ActivityFeed } from '../features/game/ActivityFeed';
import { GameHud } from '../features/game/GameHud';
import { GameScenePlaceholder } from '../features/game/GameScenePlaceholder';
import { GameShell } from '../features/game/GameShell';
import type { GameMode } from '../features/game/types';

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
  const [currentMode, setCurrentMode] = useState<GameMode>('map');

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

  const formatters = {
    formatDate,
    resourceLabel,
    unitLabel,
    resourceList,
    costList,
    unitList,
  };

  const activeEvents = events.filter((event) => event.status === 'active');
  const activeMissions = missions.filter((mission) => mission.status === 'active');
  const activeRaids = raids.filter((raid) => raid.status === 'active');
  const upgradingBuildings = buildings.filter((building) => building.isUpgrading);

  const patronPanel = (
    <PatronPanel
      crisisChoice={crisisChoice}
      error={patronError}
      formatters={formatters}
      isBreakingPatron={isBreakingPatron}
      isPayingTribute={isPayingTribute}
      joiningPatron={joiningPatron}
      loading={patronLoading}
      onBreakPatron={handleBreakPatron}
      onCrisisChoice={handleCrisisChoice}
      onJoinPatron={handleJoinPatron}
      onPayTribute={handlePayTribute}
      onRefresh={loadPatron}
      options={patronOptions}
      pressure={patronPressure}
      status={patronStatus}
    />
  );

  const resourcesPanel = (
    <ResourcesPanel error={resourcesError} formatters={formatters} loading={resourcesLoading} onRefresh={loadResources} resources={resources} />
  );

  const buildingsPanel = (
    <BuildingsPanel
      buildings={buildings}
      error={buildingsError}
      formatters={formatters}
      loading={buildingsLoading}
      onUpgrade={handleUpgrade}
      upgradingType={upgradingType}
    />
  );

  const armyPanel = (
    <ArmyPanel
      army={army}
      error={armyError}
      formatters={formatters}
      isTraining={isTraining}
      loading={armyLoading}
      onTrain={handleTrain}
      onTrainingAmountChange={setTrainingAmount}
      onTrainingTypeChange={setTrainingType}
      trainingAmount={trainingAmount}
      trainingType={trainingType}
    />
  );

  const missionsPanel = (
    <MissionsPanel
      availableMissions={availableMissions}
      error={missionsError}
      formatters={formatters}
      loading={missionsLoading}
      missionInputs={missionInputs}
      missions={missions}
      onMissionUnitAmountChange={setMissionUnitAmount}
      onRefresh={refreshMissions}
      onStartMission={handleStartMission}
      startingMissionKey={startingMissionKey}
    />
  );

  const raidsPanel = (
    <RaidsPanel
      error={raidsError}
      formatters={formatters}
      isStartingRaid={isStartingRaid}
      loading={raidsLoading}
      neighbors={neighbors}
      onRaidUnitAmountChange={setRaidUnitAmount}
      onRefresh={refreshRaids}
      onSelectRaidTargetID={setSelectedRaidTargetID}
      onStartRaid={handleStartRaid}
      raidInputs={raidInputs}
      raids={raids}
      selectedRaidTargetID={selectedRaidTargetID}
    />
  );

  const eventsPanel = (
    <EventsPanel
      choosingEventID={choosingEventID}
      error={eventsError}
      eventChoiceLabel={eventChoiceLabel}
      events={events}
      formatDate={formatDate}
      loading={eventsLoading}
      onChooseEvent={handleChooseEvent}
      onRefresh={loadEvents}
    />
  );

  const reportsPanel = (
    <ReportsPanel
      error={reportsError}
      formatters={formatters}
      loading={reportsLoading}
      loadingReportID={loadingReportID}
      markingReportID={markingReportID}
      onMarkReportRead={handleMarkReportRead}
      onRefresh={loadReports}
      onToggleReportDetails={toggleReportDetails}
      reports={reports}
      selectedReportID={selectedReportID}
      unreadReportsCount={unreadReportsCount}
    />
  );

  const contextPanelByMode: Record<GameMode, ReactNode> = {
    army: armyPanel,
    city: (
      <>
        <KingdomPanel kingdom={kingdom} patronStatus={patronStatus} user={user} />
        <RulerPanel error={rulerError} loading={rulerLoading} ruler={ruler} />
        {buildingsPanel}
      </>
    ),
    events: eventsPanel,
    map: (
      <>
        <DashboardRefreshButton onRefresh={refreshCity} />
        <KingdomPanel kingdom={kingdom} patronStatus={patronStatus} user={user} />
        {resourcesPanel}
      </>
    ),
    missions: missionsPanel,
    patron: patronPanel,
    raids: raidsPanel,
    reports: reportsPanel,
  };

  return (
    <AppShell showSidebar={false}>
      <GameShell
        activity={
          <ActivityFeed
            activeEvents={activeEvents}
            activeMissions={activeMissions}
            activeRaids={activeRaids}
            army={army}
            buildings={buildings}
            unreadReportsCount={unreadReportsCount}
          />
        }
        context={contextPanelByMode[currentMode]}
        currentMode={currentMode}
        hud={
          <GameHud
            activeEventsCount={activeEvents.length}
            kingdom={kingdom}
            patronStatus={patronStatus}
            resourceLabel={resourceLabel}
            resources={resources}
            resourcesLoading={resourcesLoading}
            unreadReportsCount={unreadReportsCount}
          />
        }
        onModeChange={setCurrentMode}
        scene={
          <GameScenePlaceholder
            activeEvents={activeEvents}
            activeMissions={activeMissions}
            activeRaids={activeRaids}
            currentMode={currentMode}
            kingdom={kingdom}
            neighbors={neighbors}
            upgradingBuildings={upgradingBuildings}
          />
        }
      />
    </AppShell>
  );
}
