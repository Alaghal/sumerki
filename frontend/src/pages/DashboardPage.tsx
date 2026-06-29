import { useEffect, useState } from 'react';

import {
  Building,
  BuildingType,
  getMyBuildings,
  getMyResources,
  getMyRuler,
  Resources,
  Ruler,
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
  independent: 'Independent',
  empire_of_dusk: 'Empire of Dusk',
  old_pact: 'Old Pact',
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

  async function refreshCity() {
    await Promise.all([loadRuler(), loadResources(), loadBuildings()]);
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
                  {kingdom.patron ? patronLabels[kingdom.patron] : 'Без покровителя'}
                </dd>
              </div>
              <div className="flex justify-between gap-4">
                <dt className="text-stone-400">Player</dt>
                <dd className="max-w-[14rem] truncate text-right text-stone-100">{user.email}</dd>
              </div>
            </dl>
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
          <Card title="Army">Militia and scouts are not trained yet.</Card>
          <Card title="Reports">Mission and raid reports will be listed here.</Card>
        </div>
      </div>
    </AppShell>
  );
}
