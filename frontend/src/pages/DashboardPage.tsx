import { useEffect, useState } from 'react';

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

const cultureLabels = {
  northern_principality: 'Северные Княжества',
  lizard_grad: 'Ящерские Грады',
  free_posad: 'Вольные Посады',
};

const patronLabels = {
  independent: 'Независимость',
  empire_of_dusk: 'Империя Заката',
  old_pact: 'Старый Договор',
};

const standingLabels = {
  hostile: 'Враждебно',
  cold: 'Холодно',
  neutral: 'Нейтрально',
  warm: 'Тепло',
  loyal: 'Верно',
};

const pressureStatusLabels = {
  none: 'Спокойно',
  warning: 'Предупреждение',
  active: 'Кризис',
  delayed: 'Отсрочка',
};

const healthLabels = {
  healthy: 'Здоров',
  wounded: 'Ранен',
  sick: 'Болен',
};

const rulerStats = [
  ['Власть', 'authority'],
  ['Храбрость', 'courage'],
  ['Хитрость', 'cunning'],
  ['Честь', 'honor'],
  ['Жестокость', 'cruelty'],
  ['Амбиция', 'ambition'],
  ['Паранойя', 'paranoia'],
] as const;

const resourceRows = [
  ['Золото', 'gold'],
  ['Еда', 'food'],
  ['Дерево', 'wood'],
  ['Камень', 'stone'],
  ['Население', 'population'],
] as const;

const costRows = [
  ['Золото', 'gold'],
  ['Дерево', 'wood'],
  ['Камень', 'stone'],
] as const;

const armyCostRows = [
  ['Золото', 'gold'],
  ['Еда', 'food'],
  ['Дерево', 'wood'],
  ['Камень', 'stone'],
  ['Население', 'population'],
] as const;

const unitStatRows = [
  ['Атака', 'attack'],
  ['Защита', 'defense'],
  ['Скорость', 'speed'],
  ['Снабжение', 'supply'],
] as const;

const unitTypes: UnitType[] = ['militia', 'spearmen', 'archers', 'cavalry', 'scouts'];

const missionTypeLabels = {
  expedition: 'Экспедиция',
  scouting: 'Разведка',
};

const missionStatusLabels = {
  active: 'В пути',
  completed: 'Завершено',
};

const missionResultLabels = {
  success: 'Успех',
  partial_success: 'Частичный успех',
  failure: 'Провал',
  attacker_success: 'Набег успешен',
  defender_success: 'Защитник отбился',
  bloody_stalemate: 'Кровавая ничья',
  repelled_by_protection: 'Набег сорван',
};

const reportTypeLabels = {
  pve_mission: 'Поход',
  pvp_raid_attacker: 'Набег',
  pvp_raid_defender: 'Защита',
  event: 'Событие',
};

const eventCategoryLabels = {
  economy: 'Хозяйство',
  ruler: 'Правитель',
  military: 'Войско',
  patron: 'Покровитель',
  dark_omen: 'Тёмное знамение',
};

const eventStatusLabels = {
  active: 'Доступно',
  resolved: 'Решено',
  expired: 'Истекло',
};

const powerEstimateLabels = {
  much_weaker: 'Намного слабее',
  weaker: 'Слабее',
  similar: 'Сравним',
  stronger: 'Сильнее',
  much_stronger: 'Намного сильнее',
};

const raidBlockedReasonLabels = {
  target_newbie_protected: 'Владение под защитой новичка',
  target_too_weak: 'Цель слишком слаба',
  raid_cooldown_active: 'Недавно уже был набег',
  target_under_protection: 'Цель под временной защитой',
};

export function DashboardPage() {
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
      setRulerError('Не удалось загрузить правителя.');
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
      setResourcesError('Не удалось загрузить ресурсы.');
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
      setBuildingsError('Не удалось загрузить здания.');
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
      setArmyError('Не удалось загрузить войско.');
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
      setMissionsError('Не удалось загрузить походы.');
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
      setReportsError('Не удалось загрузить отчёты.');
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
      setPatronError('Не удалось загрузить покровителя.');
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
      setRaidsError('Не удалось загрузить набеги.');
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
      setEventsError('Не удалось загрузить события.');
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
      setRaidsError('Выберите целое неотрицательное количество войск.');
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
      setMissionsError('Выберите целое неотрицательное количество войск.');
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
          setRulerError('Не удалось загрузить правителя.');
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
          setRaidsError('Не удалось загрузить набеги.');
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
          setPatronError('Не удалось загрузить покровителя.');
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
          setResourcesError('Не удалось загрузить ресурсы.');
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
          setBuildingsError('Не удалось загрузить здания.');
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
          setArmyError('Не удалось загрузить войско.');
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
          setMissionsError('Не удалось загрузить походы.');
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
          setReportsError('Не удалось загрузить отчёты.');
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
          setEventsError('Не удалось загрузить события.');
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
          <h1 className="text-2xl font-semibold text-stone-100">{kingdom.name}</h1>
          <p className="mt-1 text-sm text-stone-400">Settlement dashboard for {user.email}</p>
        </div>
        <Button className="justify-self-start" onClick={refreshCity} type="button">
          Обновить город
        </Button>
        <div className="grid gap-4 lg:grid-cols-2">
          <Card title="Kingdom">
            <dl className="grid gap-2">
              <div className="flex justify-between gap-4">
                <dt className="text-stone-400">Culture</dt>
                <dd className="text-right text-stone-100">{cultureLabels[kingdom.culture]}</dd>
              </div>
              <div className="flex justify-between gap-4">
                <dt className="text-stone-400">Patron</dt>
                <dd className="text-right text-stone-100">
                  {patronStatus?.patron?.label ?? (kingdom.patron ? patronLabels[kingdom.patron] : 'Без покровителя')}
                </dd>
              </div>
              <div className="flex justify-between gap-4">
                <dt className="text-stone-400">Player</dt>
                <dd className="max-w-[14rem] truncate text-right text-stone-100">{user.email}</dd>
              </div>
            </dl>
          </Card>
          <Card title="Patron">
            <div className="grid gap-4">
              {patronLoading ? <p>Загрузка покровителя...</p> : null}
              {patronError ? <p className="text-red-300">{patronError}</p> : null}
              {!patronLoading ? (
                <>
                  <div className="rounded border border-stone-800 bg-dusk-950 p-3">
                    <div className="flex flex-wrap items-start justify-between gap-3">
                      <div>
                        <h3 className="font-semibold text-stone-100">
                          {patronStatus?.patron ? patronStatus.patron.label : 'Покровитель не выбран'}
                        </h3>
                        <p className="mt-1 text-sm text-stone-400">
                          {patronPressure?.summary ?? 'Выберите политический путь для княжества.'}
                        </p>
                      </div>
                      {patronStatus?.patron ? (
                        <div className="text-right text-sm text-stone-300">
                          <div>Favor: {patronStatus.patron.favor}</div>
                          <div>{standingLabels[patronStatus.patron.standing]}</div>
                          <div>{new Date(patronStatus.patron.joinedAt).toLocaleString('ru-RU')}</div>
                        </div>
                      ) : null}
                    </div>
                    {patronStatus?.patron ? (
                      <Button className="mt-3 justify-self-start" disabled={isBreakingPatron} onClick={handleBreakPatron} type="button">
                        {isBreakingPatron ? 'Разрывается...' : 'Разорвать связь'}
                      </Button>
                    ) : null}
                  </div>

                  {patronPressure ? (
                    <div className="rounded border border-stone-800 bg-dusk-950 p-3">
                      <div className="flex flex-wrap items-start justify-between gap-3">
                        <div>
                          <h3 className="font-semibold text-stone-100">Давление</h3>
                          <p className="mt-1 text-sm text-stone-400">{patronPressure.summary}</p>
                        </div>
                        <div className="text-right text-sm text-stone-300">
                          <div>{pressureStatusLabels[patronPressure.crisisStatus]}</div>
                          <div>{patronPressure.pressureLevel} / 100</div>
                        </div>
                      </div>
                      <dl className="mt-3 grid gap-2 text-sm">
                        <div className="flex justify-between gap-4">
                          <dt className="text-stone-400">Долг дани</dt>
                          <dd className="text-right text-stone-100">
                            {patronPressure.tributeDebt.gold} золота, {patronPressure.tributeDebt.food} еды
                          </dd>
                        </div>
                        <div className="flex justify-between gap-4">
                          <dt className="text-stone-400">Обязательство</dt>
                          <dd className="text-right text-stone-100">{patronPressure.contributionDebt.food} еды</dd>
                        </div>
                        <div className="flex justify-between gap-4">
                          <dt className="text-stone-400">Следующий сбор</dt>
                          <dd className="text-right text-stone-100">
                            {patronPressure.nextTributeAt ? new Date(patronPressure.nextTributeAt).toLocaleString('ru-RU') : 'Нет'}
                          </dd>
                        </div>
                        {patronPressure.delayUntil ? (
                          <div className="flex justify-between gap-4">
                            <dt className="text-stone-400">Отсрочка до</dt>
                            <dd className="text-right text-stone-100">{new Date(patronPressure.delayUntil).toLocaleString('ru-RU')}</dd>
                          </div>
                        ) : null}
                        <div className="flex justify-between gap-4">
                          <dt className="text-stone-400">Неприкосновенный запас</dt>
                          <dd className="text-right text-stone-100">
                            {patronPressure.protectedMinimums.gold ?? 0} / {patronPressure.protectedMinimums.food ?? 0} /{' '}
                            {patronPressure.protectedMinimums.wood ?? 0} / {patronPressure.protectedMinimums.stone ?? 0}
                          </dd>
                        </div>
                      </dl>
                      <div className="mt-3 flex flex-wrap gap-2">
                        {patronPressure.availableActions.includes('pay_tribute') ? (
                          <Button disabled={isPayingTribute} onClick={handlePayTribute} type="button">
                            {isPayingTribute ? 'Платим...' : 'Заплатить'}
                          </Button>
                        ) : null}
                        {patronPressure.availableActions.includes('ask_delay') ? (
                          <Button disabled={crisisChoice === 'ask_delay'} onClick={() => handleCrisisChoice('ask_delay')} type="button">
                            {crisisChoice === 'ask_delay' ? 'Просим...' : 'Попросить отсрочку'}
                          </Button>
                        ) : null}
                        {patronPressure.availableActions.includes('break_patron') ? (
                          <Button disabled={crisisChoice === 'break_patron'} onClick={() => handleCrisisChoice('break_patron')} type="button">
                            {crisisChoice === 'break_patron' ? 'Разрываем...' : 'Разорвать через кризис'}
                          </Button>
                        ) : null}
                      </div>
                    </div>
                  ) : null}

                  <div className="flex flex-wrap gap-2">
                    <Button disabled={patronLoading} onClick={loadPatron} type="button">
                      Обновить покровителя
                    </Button>
                  </div>

                  <div className="grid gap-3">
                    {patronOptions.map((option) => {
                      const isCurrent = patronStatus?.patron?.key === option.key;
                      return (
                        <div className="rounded border border-stone-800 bg-dusk-950 p-3" key={option.key}>
                          <div className="flex flex-wrap items-start justify-between gap-3">
                            <div>
                              <h3 className="font-semibold text-stone-100">{option.label}</h3>
                              <p className="mt-1 text-sm text-stone-400">{option.shortDescription}</p>
                              <p className="mt-2 text-sm text-stone-500">{option.flavor}</p>
                            </div>
                            <Button
                              disabled={joiningPatron === option.key || isCurrent}
                              onClick={() => handleJoinPatron(option.key)}
                              type="button"
                            >
                              {isCurrent ? 'Выбрано' : option.key === 'independent' ? 'Выбрать' : 'Присоединиться'}
                            </Button>
                          </div>
                          <div className="mt-3 grid gap-2 text-sm text-stone-400">
                            <div>
                              <p className="font-semibold text-stone-300">Сейчас</p>
                              <ul className="mt-1 list-disc pl-5">
                                {option.currentEffects.map((effect) => (
                                  <li key={effect}>{effect}</li>
                                ))}
                              </ul>
                            </div>
                            <div>
                              <p className="font-semibold text-stone-300">Позже</p>
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
          <Card title="Resources">
            <div className="grid gap-4">
              {resourcesLoading ? <p>Загрузка ресурсов...</p> : null}
              {resourcesError ? <p className="text-red-300">{resourcesError}</p> : null}
              {resources && !resourcesLoading && !resourcesError ? (
                <dl className="grid gap-2">
                  {resourceRows.map(([label, key]) => (
                    <div className="flex items-center justify-between gap-4" key={key}>
                      <dt className="text-stone-400">{label}</dt>
                      <dd className="text-right">
                        <div className="font-semibold text-stone-100">{resources[key]}</div>
                        <div className="text-xs text-dusk-gold">+{resources.productionPerHour[key]} / час</div>
                      </dd>
                    </div>
                  ))}
                </dl>
              ) : null}
              <Button className="justify-self-start" disabled={resourcesLoading} onClick={loadResources} type="button">
                Обновить ресурсы
              </Button>
            </div>
          </Card>
          <Card title="Ruler">
            {rulerLoading ? <p>Загрузка правителя...</p> : null}
            {rulerError ? <p className="text-red-300">{rulerError}</p> : null}
            {ruler && !rulerLoading && !rulerError ? (
              <div className="grid gap-4">
                <dl className="grid gap-2">
                  <div className="flex justify-between gap-4">
                    <dt className="text-stone-400">Name</dt>
                    <dd className="text-right text-stone-100">{ruler.name}</dd>
                  </div>
                  <div className="flex justify-between gap-4">
                    <dt className="text-stone-400">Age</dt>
                    <dd className="text-right text-stone-100">{ruler.age}</dd>
                  </div>
                  <div className="flex justify-between gap-4">
                    <dt className="text-stone-400">Health</dt>
                    <dd className="text-right text-stone-100">{healthLabels[ruler.healthStatus]}</dd>
                  </div>
                  <div className="flex justify-between gap-4">
                    <dt className="text-stone-400">Culture</dt>
                    <dd className="text-right text-stone-100">{cultureLabels[ruler.culture]}</dd>
                  </div>
                </dl>
                <dl className="grid gap-2">
                  {rulerStats.map(([label, key]) => (
                    <div className="flex justify-between gap-4" key={key}>
                      <dt className="text-stone-400">{label}</dt>
                      <dd className="text-right text-stone-100">{ruler[key]}</dd>
                    </div>
                  ))}
                </dl>
              </div>
            ) : null}
          </Card>
          <Card title="Buildings">
            <div className="grid gap-4">
              {buildingsLoading ? <p>Загрузка зданий...</p> : null}
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
                          <h3 className="font-semibold text-stone-100">{building.label}</h3>
                          <p className="text-xs text-stone-500">{building.type}</p>
                        </div>
                        <div className="text-right text-sm text-stone-300">
                          Level {building.level}/{building.maxLevel}
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
                            Улучшается до{' '}
                            {building.upgradeFinishesAt
                              ? new Date(building.upgradeFinishesAt).toLocaleString('ru-RU')
                              : 'завершения'}
                          </p>
                        ) : null}
                        {!building.isUpgrading && building.nextUpgrade ? (
                          <div className="grid gap-2">
                            {building.nextUpgrade.blockedReason === 'max_level' ? (
                              <p className="text-dusk-gold">Максимальный уровень</p>
                            ) : (
                              <>
                                <div className="grid gap-1">
                                  <p className="text-stone-400">
                                    Upgrade to level {building.nextUpgrade.targetLevel},{' '}
                                    {building.nextUpgrade.durationSeconds} sec
                                  </p>
                                  <p className="text-stone-400">
                                    Cost:{' '}
                                    {costRows
                                      .map(([label, key]) => `${label}: ${building.nextUpgrade?.cost[key] ?? 0}`)
                                      .join(', ')}
                                  </p>
                                </div>
                                <Button
                                  className="justify-self-start"
                                  data-building-upgrade={building.type}
                                  disabled={upgradingType === building.type}
                                  onClick={() => handleUpgrade(building.type)}
                                  type="button"
                                >
                                  {upgradingType === building.type ? 'Запуск...' : 'Улучшить'}
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
          <Card title="Army">
            <div className="grid gap-4">
              {armyLoading ? <p>Загрузка войска...</p> : null}
              {armyError ? <p className="text-red-300">{armyError}</p> : null}
              {army && !armyLoading ? (
                <>
                  <dl className="grid gap-2">
                    <div className="flex justify-between gap-4">
                      <dt className="text-stone-400">Всего</dt>
                      <dd className="text-right text-stone-100">{army.summary.totalUnits}</dd>
                    </div>
                    <div className="flex justify-between gap-4">
                      <dt className="text-stone-400">Атака</dt>
                      <dd className="text-right text-stone-100">{army.summary.totalAttack}</dd>
                    </div>
                    <div className="flex justify-between gap-4">
                      <dt className="text-stone-400">Защита</dt>
                      <dd className="text-right text-stone-100">{army.summary.totalDefense}</dd>
                    </div>
                    <div className="flex justify-between gap-4">
                      <dt className="text-stone-400">Снабжение</dt>
                      <dd className="text-right text-stone-100">{army.summary.totalSupply}</dd>
                    </div>
                  </dl>

                  <div className="grid gap-3">
                    {army.units.map((unit) => (
                      <div className="rounded border border-stone-800 bg-dusk-950 p-3" data-unit-type={unit.type} key={unit.type}>
                        <div className="flex flex-wrap items-start justify-between gap-3">
                          <div>
                            <h3 className="font-semibold text-stone-100">{unit.label}</h3>
                            <p className="text-xs text-stone-500">{unit.type}</p>
                          </div>
                          <div className="text-right text-sm text-stone-300">{unit.amount}</div>
                        </div>
                        <dl className="mt-3 grid gap-1 text-sm">
                          {unitStatRows.map(([label, key]) => (
                            <div className="flex justify-between gap-4" key={key}>
                              <dt className="text-stone-400">{label}</dt>
                              <dd className="text-right text-stone-100">{unit.stats[key]}</dd>
                            </div>
                          ))}
                        </dl>
                        <p className="mt-2 text-sm text-stone-400">
                          Cost:{' '}
                          {armyCostRows.map(([label, key]) => `${label}: ${unit.cost[key]}`).join(', ')}
                        </p>
                        <p className="mt-1 text-sm text-stone-400">{unit.secondsPerUnit} sec / unit</p>
                        <p className={unit.requirements.isMet ? 'mt-1 text-sm text-dusk-gold' : 'mt-1 text-sm text-red-300'}>
                          {unit.requirements.barracksLevel > 0
                            ? `Требуется казарма уровня ${unit.requirements.barracksLevel}. ${
                                unit.requirements.isMet ? 'Требование выполнено' : 'Требование не выполнено'
                              }`
                            : 'Требование выполнено'}
                        </p>
                      </div>
                    ))}
                  </div>

                  <div className="grid gap-3 rounded border border-stone-800 bg-dusk-950 p-3">
                    <div className="grid gap-2 sm:grid-cols-[1fr_8rem_auto]">
                      <label className="grid gap-1 text-sm text-stone-400">
                        Unit
                        <select
                          className="rounded border border-stone-700 bg-dusk-900 px-3 py-2 text-stone-100"
                          disabled={isTraining}
                          onChange={(event) => setTrainingType(event.target.value as UnitType)}
                          value={trainingType}
                        >
                          {army.units.map((unit) => (
                            <option key={unit.type} value={unit.type}>
                              {unit.label}
                            </option>
                          ))}
                        </select>
                      </label>
                      <label className="grid gap-1 text-sm text-stone-400">
                        Amount
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
                        {isTraining ? 'Обучается...' : 'Обучить'}
                      </Button>
                    </div>
                    <p className="text-sm text-stone-400">Миссии и бои появятся позже.</p>
                  </div>

                  <div className="grid gap-2">
                    <h3 className="font-semibold text-stone-100">Обучается</h3>
                    {army.trainingOrders.length === 0 ? <p className="text-sm text-stone-400">Нет активного обучения.</p> : null}
                    {army.trainingOrders.map((order) => (
                      <div className="rounded border border-stone-800 bg-dusk-950 p-3" key={order.id}>
                        <div className="flex flex-wrap justify-between gap-3">
                          <div>
                            <p className="font-semibold text-stone-100">{order.unitLabel}</p>
                            <p className="text-sm text-stone-400">Amount: {order.amount}</p>
                          </div>
                          <div className="text-right text-sm text-dusk-gold">
                            Завершится {new Date(order.finishesAt).toLocaleString('ru-RU')}
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                </>
              ) : null}
            </div>
          </Card>
          <Card title="Missions">
            <div className="grid gap-4">
              {missionsLoading ? <p>Загрузка походов...</p> : null}
              {missionsError ? <p className="text-red-300">{missionsError}</p> : null}
              <Button className="justify-self-start" disabled={missionsLoading} onClick={refreshMissions} type="button">
                Обновить походы
              </Button>
              {!missionsLoading ? (
                <>
                  <div className="grid gap-3">
                    {availableMissions.map((mission) => (
                      <div className="rounded border border-stone-800 bg-dusk-950 p-3" data-mission-key={mission.key} key={mission.key}>
                        <div className="flex flex-wrap items-start justify-between gap-3">
                          <div>
                            <h3 className="font-semibold text-stone-100">{mission.label}</h3>
                            <p className="text-xs text-stone-500">{missionTypeLabels[mission.type]}</p>
                          </div>
                          <div className="text-right text-sm text-stone-300">{mission.durationSeconds} sec</div>
                        </div>
                        <p className="mt-2 text-sm text-stone-400">{mission.description}</p>
                        <p className="mt-2 text-sm text-stone-400">
                          Risk: {mission.risk}. Minimum: total {mission.minimumRequirements.totalUnits}, scouts{' '}
                          {mission.minimumRequirements.scouts}
                        </p>
                        <p className="mt-1 text-sm text-stone-400">
                          Rewards:{' '}
                          {resourceRows.map(([label, key]) => `${label}: ${mission.baseRewards[key]}`).join(', ')}
                        </p>
                        <div className="mt-3 grid gap-2 sm:grid-cols-5">
                          {unitTypes.map((unitType) => (
                            <label className="grid gap-1 text-xs text-stone-400" key={unitType}>
                              {army?.units.find((unit) => unit.type === unitType)?.label ?? unitType}
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
                          {startingMissionKey === mission.key ? 'Отправка...' : 'Отправить'}
                        </Button>
                      </div>
                    ))}
                  </div>

                  <div className="grid gap-2">
                    <h3 className="font-semibold text-stone-100">Текущие походы</h3>
                    {missions.length === 0 ? <p className="text-sm text-stone-400">Нет походов.</p> : null}
                    {missions.map((mission) => (
                      <div className="rounded border border-stone-800 bg-dusk-950 p-3" key={mission.id}>
                        <div className="flex flex-wrap justify-between gap-3">
                          <div>
                            <p className="font-semibold text-stone-100">{mission.missionLabel}</p>
                            <p className="text-sm text-stone-400">{missionStatusLabels[mission.status]}</p>
                          </div>
                          <div className="text-right text-sm text-dusk-gold">
                            {mission.completedAt
                              ? `Завершено ${new Date(mission.completedAt).toLocaleString('ru-RU')}`
                              : `Завершится ${new Date(mission.finishesAt).toLocaleString('ru-RU')}`}
                          </div>
                        </div>
                        <p className="mt-2 text-sm text-stone-400">
                          Units:{' '}
                          {mission.units
                            .map((unit) => `${unit.unitLabel}: ${unit.amountSent} sent, ${unit.amountLost} lost, ${unit.amountReturned} returned`)
                            .join(', ')}
                        </p>
                        {mission.result ? (
                          <p className="mt-1 text-sm text-stone-400">
                            Result: {missionResultLabels[mission.result.result]}
                          </p>
                        ) : null}
                      </div>
                    ))}
                  </div>
                </>
              ) : null}
            </div>
          </Card>
          <Card title="Набеги">
            <div className="grid gap-4">
              {raidsLoading ? <p>Загрузка набегов...</p> : null}
              {raidsError ? <p className="text-red-300">{raidsError}</p> : null}
              <Button className="justify-self-start" disabled={raidsLoading} onClick={refreshRaids} type="button">
                Обновить набеги
              </Button>
              {!raidsLoading ? (
                <>
                  <div className="grid gap-3">
                    <h3 className="font-semibold text-stone-100">Соседи</h3>
                    {neighbors.length === 0 ? <p className="text-sm text-stone-400">Нет доступных соседей.</p> : null}
                    {neighbors.map((neighbor) => (
                      <div className="rounded border border-stone-800 bg-dusk-950 p-3" key={neighbor.kingdomId}>
                        <div className="flex flex-wrap items-start justify-between gap-3">
                          <div>
                            <h3 className="font-semibold text-stone-100">{neighbor.name}</h3>
                            <p className="text-sm text-stone-400">{cultureLabels[neighbor.culture]}</p>
                            <p className="text-sm text-stone-400">
                              Покровитель: {neighbor.patron ? patronLabels[neighbor.patron] : 'Без покровителя'}
                            </p>
                          </div>
                          <div className="text-right text-sm text-stone-300">
                            <div>Страх: {neighbor.dread}</div>
                            <div>{powerEstimateLabels[neighbor.powerEstimate]}</div>
                          </div>
                        </div>
                        <p className={neighbor.canRaid ? 'mt-2 text-sm text-dusk-gold' : 'mt-2 text-sm text-stone-500'}>
                          {neighbor.canRaid
                            ? 'Можно начать набег'
                            : neighbor.blockedReason
                              ? raidBlockedReasonLabels[neighbor.blockedReason]
                              : 'Набег недоступен'}
                        </p>
                        {neighbor.canRaid ? (
                          <label className="mt-3 flex items-center gap-2 text-sm text-stone-400">
                            <input
                              checked={selectedRaidTargetID === neighbor.kingdomId}
                              onChange={() => setSelectedRaidTargetID(neighbor.kingdomId)}
                              type="radio"
                            />
                            Выбрать цель
                          </label>
                        ) : null}
                      </div>
                    ))}
                  </div>

                  {selectedRaidTargetID ? (
                    <div className="grid gap-3 rounded border border-stone-800 bg-dusk-950 p-3">
                      <h3 className="font-semibold text-stone-100">Отправить отряд</h3>
                      <div className="grid gap-2 sm:grid-cols-5">
                        {unitTypes.map((unitType) => (
                          <label className="grid gap-1 text-xs text-stone-400" key={unitType}>
                            {army?.units.find((unit) => unit.type === unitType)?.label ?? unitType}
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
                      <p className="text-sm text-stone-400">Минимум 3, максимум 100 воинов. Население не крадётся.</p>
                      <Button className="justify-self-start" disabled={isStartingRaid} onClick={handleStartRaid} type="button">
                        {isStartingRaid ? 'Отправка...' : 'Начать набег'}
                      </Button>
                    </div>
                  ) : null}

                  <div className="grid gap-2">
                    <h3 className="font-semibold text-stone-100">Текущие набеги</h3>
                    {raids.length === 0 ? <p className="text-sm text-stone-400">Набегов пока нет.</p> : null}
                    {raids.map((raid) => (
                      <div className="rounded border border-stone-800 bg-dusk-950 p-3" key={raid.id}>
                        <div className="flex flex-wrap justify-between gap-3">
                          <div>
                            <p className="font-semibold text-stone-100">
                              {raid.attackerKingdomName} → {raid.defenderKingdomName}
                            </p>
                            <p className="text-sm text-stone-400">{missionStatusLabels[raid.status]}</p>
                            {raid.result ? <p className="text-sm text-dusk-gold">{missionResultLabels[raid.result]}</p> : null}
                          </div>
                          <div className="text-right text-sm text-stone-400">
                            {raid.completedAt
                              ? `Завершено ${new Date(raid.completedAt).toLocaleString('ru-RU')}`
                              : `Дойдёт ${new Date(raid.arrivesAt).toLocaleString('ru-RU')}`}
                          </div>
                        </div>
                        <p className="mt-2 text-sm text-stone-400">
                          Units:{' '}
                          {raid.units
                            .map((unit) => `${unit.unitLabel}: ${unit.amountSent} sent, ${unit.amountLost} lost, ${unit.amountReturned} returned`)
                            .join(', ')}
                        </p>
                        <p className="mt-1 text-sm text-stone-400">
                          Loot: {resourceRows.map(([label, key]) => `${label}: ${raid.loot[key]}`).join(', ')}
                        </p>
                      </div>
                    ))}
                  </div>
                </>
              ) : null}
            </div>
          </Card>
          <Card title="События">
            <div className="grid gap-3">
              <div className="flex flex-wrap items-center justify-between gap-3">
                <p className="text-sm text-stone-400">Дворовые вести и тревожные выборы</p>
                <Button className="justify-self-start" disabled={eventsLoading} onClick={loadEvents} type="button">
                  Обновить события
                </Button>
              </div>
              {eventsLoading ? <p>Загрузка событий...</p> : null}
              {eventsError ? <p className="text-red-300">{eventsError}</p> : null}
              {!eventsLoading && events.length === 0 ? (
                <p className="text-sm text-stone-400">Событий пока нет. Обновите двор позже.</p>
              ) : null}
              {!eventsLoading && events.filter((event) => event.status === 'active').length > 0 ? (
                <div className="grid gap-3">
                  {events
                    .filter((event) => event.status === 'active')
                    .map((event) => (
                      <div className="rounded border border-dusk-gold bg-dusk-950 p-3" key={event.id}>
                        <div className="flex flex-wrap items-start justify-between gap-3">
                          <div>
                            <p className="text-xs uppercase tracking-wide text-dusk-gold">{eventCategoryLabels[event.category]}</p>
                            <h3 className="mt-1 font-semibold text-stone-100">{event.title}</h3>
                          </div>
                          <div className="text-right text-sm text-stone-400">
                            <div>{eventStatusLabels[event.status]}</div>
                            <div>До {new Date(event.expiresAt).toLocaleString('ru-RU')}</div>
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
                                  {choosingEventID === event.id ? 'Выбираем...' : 'Выбрать'}
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
                  <h3 className="font-semibold text-stone-100">Решённые события</h3>
                  {events
                    .filter((event) => event.status !== 'active')
                    .map((event) => (
                      <div className="rounded border border-stone-800 bg-dusk-950 p-3" key={event.id}>
                        <div className="flex flex-wrap items-start justify-between gap-3">
                          <div>
                            <p className="text-xs uppercase tracking-wide text-stone-500">{eventCategoryLabels[event.category]}</p>
                            <h3 className="mt-1 font-semibold text-stone-100">{event.title}</h3>
                          </div>
                          <div className="text-right text-sm text-stone-400">
                            <div>{eventStatusLabels[event.status]}</div>
                            {event.resolvedAt ? <div>{new Date(event.resolvedAt).toLocaleString('ru-RU')}</div> : null}
                          </div>
                        </div>
                        {event.selectedChoiceKey ? (
                          <p className="mt-2 text-sm text-stone-500">Выбор: {event.selectedChoiceKey}</p>
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
          <Card title="Reports">
            <div className="grid gap-3">
              <div className="flex flex-wrap items-center justify-between gap-3">
                <p className="text-sm text-stone-400">Непрочитано: {unreadReportsCount}</p>
                <Button className="justify-self-start" disabled={reportsLoading} onClick={loadReports} type="button">
                  Обновить отчёты
                </Button>
              </div>
              {reportsLoading ? <p>Загрузка отчётов...</p> : null}
              {reportsError ? <p className="text-red-300">{reportsError}</p> : null}
              {!reportsLoading && reports.length === 0 ? (
                <p className="text-sm text-stone-400">Отчётов пока нет. Отправьте отряд в первую экспедицию.</p>
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
                              <span className="text-xs text-stone-500">{reportTypeLabels[report.type]}</span>
                              <span className={report.isRead ? 'text-xs text-stone-500' : 'text-xs text-dusk-gold'}>
                                {report.isRead ? 'Прочитано' : 'Новое'}
                              </span>
                            </div>
                            <p className="text-sm text-dusk-gold">{missionResultLabels[report.result]}</p>
                          </div>
                          <div className="text-right text-sm text-stone-400">{new Date(report.createdAt).toLocaleString('ru-RU')}</div>
                        </div>
                        <p className="mt-2 text-sm text-stone-400">{report.body}</p>
                        <div className="mt-3 flex flex-wrap gap-2">
                          <Button onClick={() => toggleReportDetails(report.id)} type="button">
                            {isExpanded ? 'Скрыть детали' : 'Открыть отчёт'}
                          </Button>
                          <Button
                            disabled={markingReportID === report.id || report.isRead}
                            onClick={() => handleMarkReportRead(report.id)}
                            type="button"
                          >
                            {markingReportID === report.id ? 'Отмечается...' : 'Отметить прочитанным'}
                          </Button>
                        </div>
                        {loadingReportID === report.id ? <p className="mt-3 text-sm text-stone-400">Загрузка деталей...</p> : null}
                        {isExpanded && loadingReportID !== report.id ? (
                          <div className="mt-4 grid gap-3">
                            <div className="grid gap-2">
                              {report.phases.length === 0 ? <p className="text-sm text-stone-400">Подробные фазы не записаны.</p> : null}
                              {report.phases.map((phase) => (
                                <div className="rounded border border-stone-800 bg-dusk-900 p-3" key={`${report.id}-${phase.title}`}>
                                  <h4 className="font-semibold text-stone-100">{phase.title}</h4>
                                  <p className="mt-1 text-sm text-stone-400">{phase.body}</p>
                                </div>
                              ))}
                            </div>
                            <p className="text-sm text-stone-400">
                              Rewards: {resourceRows.map(([label, key]) => `${label}: ${report.rewards[key]}`).join(', ')}
                            </p>
                            <p className="text-sm text-stone-400">
                              Losses:{' '}
                              {unitTypes
                                .map(
                                  (unitType) =>
                                    `${army?.units.find((unit) => unit.type === unitType)?.label ?? unitType}: ${
                                      report.losses[unitType] ?? 0
                                    }`,
                                )
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
